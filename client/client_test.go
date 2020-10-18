package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("Невалидный URL", func(t *testing.T) {
		client, err := NewClient("url", "login", "hash")
		assert.EqualError(t, err, `parse "url": invalid URI for request`)
		assert.Empty(t, client)
	})

	t.Run("Пустой login", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "", "hash")
		assert.EqualError(t, err, "empty_login")
		assert.Empty(t, client)
	})

	t.Run("Пустой hash", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "")
		assert.EqualError(t, err, "empty_api_hash")
		assert.Empty(t, client)
	})

	t.Run("Клиент по-умолчанию", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "hash")
		assert.NoError(t, err)
		assert.Equal(t, "login", client.login)
		assert.Equal(t, "localhost:1234", client.baseURL)
		assert.Equal(t, "hash", client.apiHash)
		assert.Equal(t, 10*time.Second, client.httpClient.Timeout)
	})

	t.Run("Клиент с HTTP timeout", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "hash", WithHTTPTimeout(time.Second))
		assert.NoError(t, err)
		assert.Equal(t, "login", client.login)
		assert.Equal(t, "localhost:1234", client.baseURL)
		assert.Equal(t, "hash", client.apiHash)
		assert.Equal(t, time.Second, client.httpClient.Timeout)
	})
}
