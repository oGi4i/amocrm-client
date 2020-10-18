package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func defaultTestClient() (*Client, error) {
	return NewClient("localhost:1234", "clientID", "clientSecret", "refreshToken")
}

func defaultTestClientWithURL(baseURL string) (*Client, error) {
	return NewClient(baseURL, "clientID", "clientSecret", "refreshToken")
}

func TestNewClient(t *testing.T) {
	t.Run("Невалидный URL", func(t *testing.T) {
		client, err := defaultTestClientWithURL("url")
		assert.EqualError(t, err, `parse "url": invalid URI for request`)
		assert.Empty(t, client)
	})

	t.Run("Пустой clientID", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "", "clientSecret", "refreshToken")
		assert.EqualError(t, err, "empty_client_id")
		assert.Empty(t, client)
	})

	t.Run("Пустой clientSecret", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "clientID", "", "refreshToken")
		assert.EqualError(t, err, "empty_client_secret")
		assert.Empty(t, client)
	})

	t.Run("Пустой refreshToken", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "clientID", "clientSecret", "")
		assert.EqualError(t, err, "empty_refresh_token")
		assert.Empty(t, client)
	})

	t.Run("Клиент по-умолчанию", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)
		assert.Equal(t, "clientID", client.clientID)
		assert.Equal(t, "localhost:1234", client.baseURL)
		assert.Equal(t, "clientSecret", client.clientSecret)
		assert.Equal(t, "refreshToken", client.token.RefreshToken)
		assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
	})

	t.Run("Клиент с HTTP timeout", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "clientID", "clientSecret", "refreshToken", WithHTTPTimeout(time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "clientID", client.clientID)
		assert.Equal(t, "localhost:1234", client.baseURL)
		assert.Equal(t, "clientSecret", client.clientSecret)
		assert.Equal(t, "refreshToken", client.token.RefreshToken)
		assert.Equal(t, time.Second, client.httpClient.Timeout)
	})
}
