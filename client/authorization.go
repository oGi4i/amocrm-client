package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type (
	AuthGrantType string

	AuthRequest struct {
		ClientID     string        `json:"client_id" validate:"required"`
		ClientSecret string        `json:"client_secret" validate:"required"`
		GrantType    AuthGrantType `json:"grant_type" validate:"oneof=authorization_code refresh_token"`
		AuthCode     string        `json:"code,omitempty" validate:"required_if=GrantType authorization_code"`
		RefreshToken string        `json:"refresh_token,omitempty" validate:"required_if=GrantType refresh_token"`
		RedirectURI  string        `json:"redirect_uri,omitempty" validate:"omitempty"`
	}

	AuthTokenType string

	AuthResponse struct {
		TokenType    AuthTokenType `json:"token_type" validate:"required"`
		ExpiresIn    int64         `json:"expires_in" validate:"required"`
		AccessToken  string        `json:"access_token" validate:"required"`
		RefreshToken string        `json:"refresh_token" validate:"required"`
	}
)

const (
	authorizationCodeAuthGrantType AuthGrantType = "authorization_code"
	refreshTokenAuthGrantType      AuthGrantType = "refresh_token"

	bearerAuthTokenType AuthTokenType = "Bearer"

	authRetryCount = 5
)

func (c *Client) Authorize(ctx context.Context) error {
	err := c.authorizeWithCode(ctx)
	if err != nil {
		return err
	}

	go c.refreshAuthTokens(ctx)

	return nil
}

func (c *Client) authorizeWithCode(ctx context.Context) error {
	c.token.mu.RLock()
	authRequest := &AuthRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		GrantType:    authorizationCodeAuthGrantType,
		AuthCode:     c.token.AuthorizationCode,
	}
	c.token.mu.RUnlock()

	return c.authorize(ctx, authRequest)
}

func (c *Client) authorizeWithRefreshToken(ctx context.Context) error {
	c.token.mu.RLock()
	authRequest := &AuthRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		GrantType:    refreshTokenAuthGrantType,
		RefreshToken: c.token.RefreshToken,
	}
	c.token.mu.RUnlock()

	return c.authorize(ctx, authRequest)
}

func (c *Client) authorize(ctx context.Context, authRequest *AuthRequest) error {
	if err := c.validator.Struct(authRequest); err != nil {
		return err
	}

	reqBody, err := json.Marshal(authRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+authURI, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	addApplicationJSONContentType(req)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("http status not ok: " + strconv.Itoa(response.StatusCode))
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if len(respBody) == 0 {
		return ErrEmptyResponse
	}

	authResponse := new(AuthResponse)
	err = json.Unmarshal(respBody, authResponse)
	if err != nil {
		return err
	}

	if authResponse.TokenType != bearerAuthTokenType {
		return ErrInvalidAuthTokenType
	}

	if err := c.validator.Struct(authResponse); err != nil {
		return err
	}

	c.token.mu.Lock()
	c.token.AccessToken = authResponse.AccessToken
	c.token.RefreshToken = authResponse.RefreshToken
	c.token.ExpiresAt = time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)
	c.token.mu.Unlock()

	return nil
}

func (c *Client) refreshAuthTokens(ctx context.Context) {
	c.token.mu.RLock()
	expiresIn := time.Until(c.token.ExpiresAt)
	c.token.mu.RUnlock()

	// calculate when to start refreshing tokens as 10% till expiration time
	ticker := time.NewTicker(expiresIn - (expiresIn / 100 * 10))
	// calculate backOff on fail as 10% till expiration time divided by retry count
	backOff := (expiresIn / 100 * 10) / (authRetryCount + 1)

	var err error
	select {
	case <-ticker.C:
		err = c.authorizeWithRetry(ctx, backOff)
		if err != nil {
			fmt.Printf("failed to refresh tokens: %v\n", err)
			return
		}
	case <-ctx.Done():
		fmt.Println("context canceled while waiting to refresh tokens")
		return
	}

	go c.refreshAuthTokens(ctx)
}

func (c *Client) authorizeWithRetry(ctx context.Context, backOff time.Duration) error {
	var err error
	for i := 0; i < authRetryCount; i++ {
		err = c.authorizeWithRefreshToken(ctx)
		if err != nil {
			fmt.Printf("failed to authorize with refresh token: %v\n", err)
			ticker := time.NewTicker(backOff)
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return err
			}
		}
		break
	}

	return err
}

func (c *Client) withAuthToken(req *http.Request) {
	c.token.mu.RLock()
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", bearerAuthTokenType, c.token.AccessToken))
	c.token.mu.RUnlock()
}
