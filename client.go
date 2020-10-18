package amocrm

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type (
	ClientOption func(c *Client)

	AuthToken struct {
		mu           sync.RWMutex
		AccessToken  string
		RefreshToken string
		ExpiresAt    time.Time
	}

	Client struct {
		clientID     string
		clientSecret string
		baseURL      string
		client       *http.Client
		validator    *validator.Validate
		token        *AuthToken
	}

	PostResponse struct {
		ID        int `json:"id" validate:"omitempty"`
		RequestID int `json:"request_id" validate:"omitempty"`
		Embedded  struct {
			Items []struct {
				ID int `json:"id" validate:"omitempty"`
			} `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response" validate:"omitempty"`
	}
)

const (
	authURI      = "/oauth2/access_token"
	notesURI     = "/api/v2/notes"
	contactsURI  = "/api/v2/contacts"
	accountURI   = "/api/v2/account"
	leadsURI     = "/api/v2/leads"
	tasksURI     = "/api/v2/tasks"
	pipelinesURI = "/api/v2/pipelines"
	downloadURI  = "/download/"

	defaultHTTPTimeout = 5 * time.Second
)

func NewClient(clientID, clientSecret, refreshToken, accountURL string, opts ...ClientOption) (*Client, error) {
	if clientID == "" {
		return nil, ErrEmptyClientID
	}
	if clientSecret == "" {
		return nil, ErrEmptyClientSecret
	}
	if refreshToken == "" {
		return nil, ErrEmptyRefreshToken
	}

	_, err := url.Parse(accountURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		clientID:     clientID,
		clientSecret: clientID,
		baseURL:      accountURL,
		token: &AuthToken{
			RefreshToken: refreshToken,
		},
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   defaultHTTPTimeout,
		},
		validator: validator.New(),
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

func WithHTTPTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = d
	}
}

func (c *Client) DownloadAttachment(ctx context.Context, attachment string) ([]byte, error) {
	return c.doGet(ctx, c.baseURL+downloadURI+attachment, nil)
}

func (c *Client) doGet(ctx context.Context, url string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c.withAuthToken(req)

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) doPost(ctx context.Context, url string, data interface{}) (*http.Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	c.withAuthToken(req)

	return c.client.Do(req.WithContext(ctx))
}

func (c *Client) getResponseID(resp *http.Response) (int, error) {
	result := new(PostResponse)
	dec := json.NewDecoder(resp.Body)

	err := dec.Decode(result)
	if err != nil {
		amoError := new(AmoError)
		err = dec.Decode(amoError)
		if err != nil {
			return 0, err
		}

		return 0, amoError
	}

	if len(result.Embedded.Items) == 0 {
		if result.Response != nil {
			return 0, result.Response
		}
		return 0, ErrEmptyResponseItems
	}

	return result.Embedded.Items[0].ID, nil
}
