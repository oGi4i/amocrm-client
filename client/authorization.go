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
		ClientID     string        `json:"client_id"`
		ClientSecret string        `json:"client_secret"`
		GrantType    AuthGrantType `json:"grant_type"`
		AuthCode     string        `json:"code,omitempty"`
		RefreshToken string        `json:"refresh_token,omitempty"`
		RedirectURI  string        `json:"redirect_uri,omitempty"`
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
	c.token.mu.RLock()
	authRequest := &AuthRequest{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		GrantType:    refreshTokenAuthGrantType,
		RefreshToken: c.token.RefreshToken,
	}
	c.token.mu.RUnlock()

	body, err := json.Marshal(authRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+authURI, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	addApplicationJsonContentType(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("http status not ok: " + strconv.Itoa(resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
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
	defer c.token.mu.Unlock()

	c.token.AccessToken = authResponse.AccessToken
	c.token.RefreshToken = authResponse.RefreshToken
	c.token.ExpiresAt = time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)

	go c.refreshAuthTokens(ctx)

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
			fmt.Printf("failed to refresh tokens with error: %v", err)
			return
		}
	case <-ctx.Done():
		return
	}
}

func (c *Client) authorizeWithRetry(ctx context.Context, backOff time.Duration) error {
	var err error
	for i := 0; i < authRetryCount; i++ {
		err = c.Authorize(ctx)
		if err != nil {
			ticker := time.NewTicker(backOff)
			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return err
			}
		}
	}

	return err
}

func (c *Client) withAuthToken(req *http.Request) {
	c.token.mu.RLock()
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", bearerAuthTokenType, c.token.AccessToken))
	c.token.mu.RUnlock()
}
