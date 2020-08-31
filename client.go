package amocrm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

type (
	Client struct {
		userLogin string
		apiHash   string
		timezone  string
		baseURL   string
		cookie    []*http.Cookie
		client    *http.Client
		validator *validator.Validate
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
)

func NewClient(accountURL string, login string, hash string) (*Client, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	if hash == "" {
		return nil, ErrEmptyAPIHash
	}

	c := &Client{
		userLogin: login,
		apiHash:   hash,
	}

	_, err := url.Parse(accountURL)
	if err != nil {
		return nil, err
	}

	c.baseURL = accountURL

	values := url.Values{}
	values.Set("USER_LOGIN", c.userLogin)
	values.Set("USER_HASH", c.apiHash)

	body := strings.NewReader(values.Encode())
	resp, err := http.Post(c.baseURL+authURI, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		c.cookie = resp.Cookies()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		authResponse := new(AuthResponse)
		err = json.Unmarshal(body, authResponse)
		if err != nil {
			return nil, err
		}

		if len(authResponse.Response.Accounts) > 0 {
			c.timezone = authResponse.Response.Accounts[0].Timezone
		}

		if !authResponse.Response.Auth {
			return nil, errors.New(authResponse.Response.Error)
		}

		if err := c.validator.Struct(authResponse); err != nil {
			return nil, err
		}

		return c, nil
	}

	return nil, errors.New("http status not ok: " + strconv.Itoa(resp.StatusCode))
}

func (c *Client) DownloadAttachment(ctx context.Context, attachment string) ([]byte, error) {
	return c.doGet(ctx, c.baseURL+downloadURI+attachment, nil)
}

func (c *Client) doGet(ctx context.Context, url string, params map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}

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
	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}

	log.Printf("Request: Created: %s; URL: %s; Headers: %v; Body: %v", time.Now().Format(time.RFC3339), req.URL, req.Header, data)

	return c.client.Do(req.WithContext(ctx))
}

func (c *Client) DoPostWithoutCookie(ctx context.Context, url string, data string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
