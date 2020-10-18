package client

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequestValidation(t *testing.T) {
	v := validator.New()

	testCases := []struct {
		name        string
		authRequest *AuthRequest
		errString   string
	}{
		{
			name: "Отсутствует ClientID",
			authRequest: &AuthRequest{
				ClientID:     "",
				ClientSecret: "clientSecret",
				GrantType:    "authorization_code",
				AuthCode:     "code",
			},
			errString: "Key: 'AuthRequest.ClientID' Error:Field validation for 'ClientID' failed on the 'required' tag",
		},
		{
			name: "Отсутствует ClientSecret",
			authRequest: &AuthRequest{
				ClientID:     "clientID",
				ClientSecret: "",
				GrantType:    "authorization_code",
				AuthCode:     "code",
			},
			errString: "Key: 'AuthRequest.ClientSecret' Error:Field validation for 'ClientSecret' failed on the 'required' tag",
		},
		{
			name: "Невалидный GrantType",
			authRequest: &AuthRequest{
				ClientID:     "clientID",
				ClientSecret: "clientSecret",
				GrantType:    "unknown_grant_type",
			},
			errString: "Key: 'AuthRequest.GrantType' Error:Field validation for 'GrantType' failed on the 'oneof' tag",
		},
		{
			name: "Не передан AuthCode при GrantType = authorization_code",
			authRequest: &AuthRequest{
				ClientID:     "clientID",
				ClientSecret: "clientSecret",
				GrantType:    "authorization_code",
			},
			errString: "Key: 'AuthRequest.AuthCode' Error:Field validation for 'AuthCode' failed on the 'required_if' tag",
		},
		{
			name: "Не передан RefreshToken при GrantType = refresh_token",
			authRequest: &AuthRequest{
				ClientID:     "clientID",
				ClientSecret: "clientSecret",
				GrantType:    "refresh_token",
			},
			errString: "Key: 'AuthRequest.RefreshToken' Error:Field validation for 'RefreshToken' failed on the 'required_if' tag",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.EqualError(t, v.Struct(tc.authRequest), tc.errString)
		})
	}
}

func TestAuthorizeWithCode(t *testing.T) {
	const (
		sampleAuthorizeResponseBody = `{"token_type":"Bearer","expires_in":86400,"access_token":"ACCESS_TOKEN_VALUE","refresh_token":"REFRESH_TOKEN_VALUE"}`
		requestBodyWant             = `{"client_id":"clientID","client_secret":"clientSecret","grant_type":"authorization_code","code":"authorizationCode"}`
	)

	ctx := context.Background()

	t.Run("Успешная авторизация", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAuthorizeResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		err = client.authorizeWithCode(ctx)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Equal(t, "ACCESS_TOKEN_VALUE", client.token.AccessToken)
		assert.Equal(t, "REFRESH_TOKEN_VALUE", client.token.RefreshToken)
		assert.Equal(t, time.Now().Add(86400*time.Second).Round(time.Second), client.token.ExpiresAt.Round(time.Second))
	})

	//nolint:dupl
	t.Run("Пустое тело ответа", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, ``)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		err = client.authorizeWithCode(ctx)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Equal(t, "", client.token.AccessToken)
		assert.Equal(t, "", client.token.RefreshToken)
		assert.Equal(t, time.Time{}, client.token.ExpiresAt)
	})

	t.Run("Невалидный код ответа", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, ``)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		err = client.authorizeWithCode(ctx)
		assert.EqualError(t, err, "http status not ok: 400")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Equal(t, "", client.token.AccessToken)
		assert.Equal(t, "", client.token.RefreshToken)
		assert.Equal(t, time.Time{}, client.token.ExpiresAt)
	})

	//nolint:dupl
	t.Run("Невалидный тип токена в ответе", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"token_type":"InvalidTokenType","expires_in":86400,"access_token":"ACCESS_TOKEN_VALUE","refresh_token":"REFRESH_TOKEN_VALUE"}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		err = client.authorizeWithCode(ctx)
		assert.EqualError(t, err, ErrInvalidAuthTokenType.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Equal(t, "", client.token.AccessToken)
		assert.Equal(t, "", client.token.RefreshToken)
		assert.Equal(t, time.Time{}, client.token.ExpiresAt)
	})
}

func TestRefreshAuthTokens(t *testing.T) {
	const (
		oldRefreshToken = "oldRefreshToken"
		oldAccessToken  = "oldAccessToken"
		newRefreshToken = "REFRESH_TOKEN_VALUE"
		newAccessToken  = "ACCESS_TOKEN_VALUE"

		sampleAuthorizeResponseBody = `{"token_type":"Bearer","expires_in":86400,"access_token":"ACCESS_TOKEN_VALUE","refresh_token":"REFRESH_TOKEN_VALUE"}`
		requestBodyWant             = `{"client_id":"clientID","client_secret":"clientSecret","grant_type":"refresh_token","refresh_token":"oldRefreshToken"}`
	)

	ctx := context.Background()

	t.Run("Успешное обновление токенов с первой попытки", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAuthorizeResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		expiresAt := time.Now().Add(2 * time.Second)
		client.token.AccessToken = oldAccessToken
		client.token.RefreshToken = oldRefreshToken
		client.token.ExpiresAt = expiresAt

		client.refreshAuthTokens(ctx)

		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Equal(t, newAccessToken, client.token.AccessToken)
		assert.Equal(t, newRefreshToken, client.token.RefreshToken)
		assert.Equal(t, time.Now().Add(86400*time.Second).Round(time.Second), client.token.ExpiresAt.Round(time.Second))
	})

	t.Run("Успешное обновление токенов со второй попытки", func(t *testing.T) {
		failedChan := make(chan struct{})
		failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, "")
			failedChan <- struct{}{}
		}))
		var successRequestBodyGot []byte
		successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			successRequestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAuthorizeResponseBody)
		}))

		client, err := defaultTestClientWithURL(failServer.URL)
		assert.NoError(t, err)

		expiresAt := time.Now().Add(time.Second)
		client.token.AccessToken = oldAccessToken
		client.token.RefreshToken = oldRefreshToken
		client.token.ExpiresAt = expiresAt

		failedAttempts := 0
		go func() {
			<-failedChan
			failedAttempts++
			client.token.mu.Lock()
			client.baseURL = successServer.URL
			client.token.mu.Unlock()
		}()

		client.refreshAuthTokens(ctx)

		assert.Equal(t, 1, failedAttempts)
		assert.Equal(t, requestBodyWant, string(successRequestBodyGot))
		assert.Equal(t, newAccessToken, client.token.AccessToken)
		assert.Equal(t, newRefreshToken, client.token.RefreshToken)
		assert.Equal(t, time.Now().Add(86400*time.Second).Round(time.Second), client.token.ExpiresAt.Round(time.Second))
	})

	t.Run("Успешное обновление токенов с последней попытки", func(t *testing.T) {
		failedChan := make(chan struct{})
		failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, "")
			failedChan <- struct{}{}
		}))
		var successRequestBodyGot []byte
		successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			successRequestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAuthorizeResponseBody)
		}))

		client, err := defaultTestClientWithURL(failServer.URL)
		assert.NoError(t, err)

		expiresAt := time.Now().Add(time.Second)
		client.token.AccessToken = oldAccessToken
		client.token.RefreshToken = oldRefreshToken
		client.token.ExpiresAt = expiresAt

		failedAttempts := 0
		go func() {
			for {
				if failedAttempts == authRetryCount-1 {
					client.token.mu.Lock()
					client.baseURL = successServer.URL
					client.token.mu.Unlock()
					return
				}

				<-failedChan
				failedAttempts++
			}
		}()

		client.refreshAuthTokens(ctx)

		assert.Equal(t, 4, failedAttempts)
		assert.Equal(t, requestBodyWant, string(successRequestBodyGot))
		assert.Equal(t, newAccessToken, client.token.AccessToken)
		assert.Equal(t, newRefreshToken, client.token.RefreshToken)
		assert.Equal(t, time.Now().Add(86400*time.Second).Round(time.Second), client.token.ExpiresAt.Round(time.Second))
	})

	t.Run("Неуспешное обновление токенов", func(t *testing.T) {
		var requestBodyGot []byte
		failedChan := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, "")
			failedChan <- struct{}{}
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		expiresAt := time.Now().Add(time.Second)
		client.token.AccessToken = oldAccessToken
		client.token.RefreshToken = oldRefreshToken
		client.token.ExpiresAt = expiresAt

		failedAttempts := 0
		go func() {
			for {
				if failedAttempts == authRetryCount {
					return
				}

				<-failedChan
				failedAttempts++
			}
		}()

		client.refreshAuthTokens(ctx)

		assert.Equal(t, authRetryCount, failedAttempts)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Equal(t, oldAccessToken, client.token.AccessToken)
		assert.Equal(t, oldRefreshToken, client.token.RefreshToken)
		assert.Equal(t, expiresAt, client.token.ExpiresAt)
	})

	t.Run("Отмена контекста во время обновления токенов", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		expiresAt := time.Now().Add(time.Second)
		client.token.AccessToken = oldAccessToken
		client.token.RefreshToken = oldRefreshToken
		client.token.ExpiresAt = expiresAt

		ctx, cancel := context.WithCancel(ctx)
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			<-ticker.C
			cancel()
		}()

		client.refreshAuthTokens(ctx)

		assert.Equal(t, oldAccessToken, client.token.AccessToken)
		assert.Equal(t, oldRefreshToken, client.token.RefreshToken)
		assert.Equal(t, expiresAt, client.token.ExpiresAt)
	})
}
