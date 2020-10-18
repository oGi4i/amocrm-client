package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func defaultTestClient() (*Client, error) {
	return NewClient("localhost:1234", "clientID", "clientSecret", "authorizationCode")
}

func defaultTestClientWithURL(baseURL string) (*Client, error) {
	return NewClient(baseURL, "clientID", "clientSecret", "authorizationCode")
}

func TestNewClient(t *testing.T) {
	t.Run("Невалидный URL", func(t *testing.T) {
		client, err := defaultTestClientWithURL("url")
		assert.EqualError(t, err, `parse "url": invalid URI for request`)
		assert.Empty(t, client)
	})

	t.Run("Пустой clientID", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "", "clientSecret", "authorizationCode")
		assert.EqualError(t, err, "empty_client_id")
		assert.Empty(t, client)
	})

	t.Run("Пустой clientSecret", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "clientID", "", "authorizationCode")
		assert.EqualError(t, err, "empty_client_secret")
		assert.Empty(t, client)
	})

	t.Run("Пустой authorizationCode", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "clientID", "clientSecret", "")
		assert.EqualError(t, err, "empty_authorization_code")
		assert.Empty(t, client)
	})

	t.Run("Клиент по-умолчанию", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)
		assert.Equal(t, "clientID", client.clientID)
		assert.Equal(t, "localhost:1234", client.baseURL)
		assert.Equal(t, "clientSecret", client.clientSecret)
		assert.Equal(t, "authorizationCode", client.token.AuthorizationCode)
		assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
	})

	t.Run("Клиент с HTTP timeout", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "clientID", "clientSecret", "authorizationCode", WithHTTPTimeout(time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "clientID", client.clientID)
		assert.Equal(t, "localhost:1234", client.baseURL)
		assert.Equal(t, "clientSecret", client.clientSecret)
		assert.Equal(t, "authorizationCode", client.token.AuthorizationCode)
		assert.Equal(t, time.Second, client.httpClient.Timeout)
	})
}
