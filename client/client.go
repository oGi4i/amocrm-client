package client

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

	"github.com/go-playground/validator/v10"
)

type (
	Option func(c *Client)

	AuthToken struct {
		mu                sync.RWMutex
		AuthorizationCode string
		AccessToken       string
		RefreshToken      string
		ExpiresAt         time.Time
	}

	Client struct {
		clientID     string
		clientSecret string
		baseURL      string
		httpClient   *http.Client
		validator    *validator.Validate
		token        *AuthToken
	}
)

const (
	authURI      = "/oauth2/access_token"
	accountURI   = "/api/v4/account"
	companiesURI = "/api/v4/companies"
	contactsURI  = "/api/v4/contacts"
	downloadURI  = "/download/"
	leadsURI     = "/api/v4/leads"
	pipelinesURI = "/api/v4/leads/pipelines"
	tasksURI     = "/api/v4/tasks"
)

const (
	contentTypeHeader = "ContentType"

	applicationJSONContentType = "application/json"

	successContentType = "application/hal+json"
	errorContentType   = "application/problem+json"
)

const defaultHTTPTimeout = 10 * time.Second

func NewClient(baseURL, clientID, clientSecret, authorizationCode string, opts ...Option) (*Client, error) {
	if clientID == "" {
		return nil, ErrEmptyClientID
	}
	if clientSecret == "" {
		return nil, ErrEmptyClientSecret
	}
	if authorizationCode == "" {
		return nil, ErrEmptyAuthorizationCode
	}

	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
		token: &AuthToken{
			AuthorizationCode: authorizationCode,
		},
		httpClient: &http.Client{
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
		c.httpClient.Timeout = d
	}
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
	c.withAuthToken(req)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if !isSuccessResponse(response) {
		return nil, errors.New("http status not ok: " + strconv.Itoa(response.StatusCode))
	}

	body, err := ioutil.ReadAll(response.Body)
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

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	addApplicationJSONContentType(req)
	c.withAuthToken(req)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if !isSuccessResponse(response) {
		return nil, errors.New("http status not ok: " + strconv.Itoa(response.StatusCode))
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func isSuccessResponse(response *http.Response) bool {
	switch response.Header.Get(contentTypeHeader) {
	case successContentType:
		return true
	case errorContentType:
		return false
	default:
		return false
	}
}

func addApplicationJSONContentType(req *http.Request) {
	req.Header.Set(contentTypeHeader, applicationJSONContentType)
}
