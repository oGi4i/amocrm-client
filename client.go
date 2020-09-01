package amocrm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type (
	ClientOption func(c *Client)

	Client struct {
		userLogin string
		apiHash   string
		timezone  string
		baseURL   string
		cookie    []*http.Cookie
		client    *http.Client
		validator *validator.Validate
		mu        sync.RWMutex
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

	AuthResponse struct {
		Response struct {
			Auth       bool           `json:"auth" validate:"required"`
			Accounts   []*AuthAccount `json:"accounts" validate:"omitempty,dive,required"`
			User       *AuthUser      `json:"user" validate:"required"`
			ServerTime int            `json:"server_time" validate:"required"`
			Error      string         `json:"error" validate:"omitempty"`
		} `json:"response" validate:"required"`
	}
)

const (
	authURI      = "/private/api/auth.php?type=json"
	notesURI     = "/api/v2/notes"
	contactsURI  = "/api/v2/contacts"
	accountURI   = "/api/v2/account"
	leadsURI     = "/api/v2/leads"
	tasksURI     = "/api/v2/tasks"
	pipelinesURI = "/api/v2/pipelines"
	downloadURI  = "/download/"

	defaultHTTPTimeout = 5 * time.Second
)

func NewClient(accountURL string, login string, hash string, opts ...ClientOption) (*Client, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	if hash == "" {
		return nil, ErrEmptyAPIHash
	}

	_, err := url.Parse(accountURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		userLogin: login,
		apiHash:   hash,
		baseURL:   accountURL,
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

func (c *Client) Authorize(ctx context.Context) error {
	values := url.Values{}
	values.Set("USER_LOGIN", c.userLogin)
	values.Set("USER_HASH", c.apiHash)

	req, err := http.NewRequest(http.MethodPost, c.baseURL+authURI, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		c.mu.Lock()
		defer c.mu.Unlock()
		c.cookie = resp.Cookies()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		authResponse := new(AuthResponse)
		err = json.Unmarshal(body, authResponse)
		if err != nil {
			return err
		}

		if len(authResponse.Response.Accounts) > 0 {
			c.timezone = authResponse.Response.Accounts[0].Timezone
		}

		if !authResponse.Response.Auth {
			return errors.New(authResponse.Response.Error)
		}

		if err := c.validator.Struct(authResponse); err != nil {
			return err
		}

		return nil
	}

	return errors.New("http status not ok: " + strconv.Itoa(resp.StatusCode))
}

func (c *Client) DownloadAttachment(ctx context.Context, attachment string) ([]byte, error) {
	return c.doGet(ctx, c.baseURL+downloadURI+attachment, nil)
}

func (c *Client) doGet(ctx context.Context, url string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c.mu.RLock()
	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}
	c.mu.RUnlock()

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

	c.mu.RLock()
	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}
	c.mu.RUnlock()

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
