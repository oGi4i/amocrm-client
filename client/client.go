package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/ogi4i/amocrm-client/domain"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type (
	Option func(c *Client)

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
		ErrorResponse *domain.AmoError `json:"response" validate:"omitempty"`
	}

	AuthResponse struct {
		Response struct {
			Auth       bool             `json:"auth" validate:"required"`
			Accounts   []*AuthAccount   `json:"accounts" validate:"omitempty,dive,required"`
			User       *domain.AuthUser `json:"user" validate:"required"`
			ServerTime int              `json:"server_time" validate:"required"`
			Error      string           `json:"error" validate:"omitempty"`
		} `json:"response" validate:"required"`
	}
)

const (
	authURI      = "/private/api/auth.php?type=json"
	notesURI     = "/api/v2/note"
	contactsURI  = "/api/v2/contacts"
	accountURI   = "/api/v4/account"
	leadsURI     = "/api/v4/leads"
	tasksURI     = "/api/v2/tasks"
	pipelinesURI = "/api/v4/leads/pipelines"
	downloadURI  = "/download/"

	contentTypeHeader  = "ContentType"
	successContentType = "application/hal+json"
	errorContentType   = "application/problem+json"

	defaultHTTPTimeout = 10 * time.Second
)

func NewClient(accountURL string, login string, hash string, opts ...Option) (*Client, error) {
	if login == "" {
		return nil, domain.ErrEmptyLogin
	}
	if hash == "" {
		return nil, domain.ErrEmptyAPIHash
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

func WithHTTPTimeout(d time.Duration) Option {
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

	if resp.StatusCode != 200 {
		return errors.New("http status not ok: " + strconv.Itoa(resp.StatusCode))
	}

	c.mu.Lock()
	c.cookie = resp.Cookies()
	c.mu.Unlock()

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

func (c *Client) DownloadAttachment(ctx context.Context, attachment string) ([]byte, error) {
	return c.doGet(ctx, c.baseURL+downloadURI+attachment, nil)
}

func (c *Client) doGet(ctx context.Context, url string, params url.Values) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()

	c.mu.RLock()
	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}
	c.mu.RUnlock()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccessResponse(resp) {
		return nil, errors.New("http status not ok: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) do(ctx context.Context, url string, method string, data interface{}) ([]byte, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set(contentTypeHeader, "application/json")

	c.mu.RLock()
	for _, cookie := range c.cookie {
		req.AddCookie(cookie)
	}
	c.mu.RUnlock()

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccessResponse(resp) {
		return nil, errors.New("http status not ok: " + strconv.Itoa(resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) getResponseID(body []byte) (int, error) {
	result := new(PostResponse)
	err := json.Unmarshal(body, result)
	if err != nil {
		return 0, err
	}

	if len(result.Embedded.Items) == 0 {
		if result.ErrorResponse != nil {
			return 0, result.ErrorResponse
		}
		return 0, domain.ErrEmptyResponse
	}

	return result.Embedded.Items[0].ID, nil
}

func isSuccessResponse(resp *http.Response) bool {
	switch resp.Header.Get(contentTypeHeader) {
	case successContentType:
		return true
	case errorContentType:
		return false
	default:
		return false
	}
}
