package amocrm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Информация о подключении к аккаунту
type clientInfo struct {
	userLogin         string
	apiHash           string
	Timezone          string
	accountWebAddress *url.URL
	Timeout           time.Duration
}

//AuthResponse Структура ответа при авторизации
type AuthResponse struct {
	Response struct {
		Auth     bool `json:"auth"`
		Accounts []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Subdomain string `json:"subdomain"`
			Language  string `json:"language"`
			Timezone  string `json:"timezone"`
		} `json:"accounts"`
		ServerTime int    `json:"server_time"`
		Error      string `json:"error"`
	} `json:"response"`
}

//New Открытия соединения и авторизация
func New(accountURL string, login string, hash string, timeout time.Duration) (*clientInfo, error) {
	var err error

	c := &clientInfo{
		userLogin: login,
		apiHash:   hash,
		Timeout:   timeout,
	}
	c.accountWebAddress, err = url.Parse(accountURL)
	if err != nil {
		return nil, err
	}
	requestURL := c.accountWebAddress
	requestURL.Path = apiUrls["auth"]
	params := requestURL.Query()
	params.Set("type", "json")
	requestURL.RawQuery = params.Encode()
	client := &http.Client{
		Timeout: c.Timeout,
	}

	resp, err := client.PostForm(
		requestURL.String(),
		url.Values{"USER_LOGIN": {c.userLogin}, "USER_HASH": {c.apiHash}},
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return nil, err
	}
	if len(authResponse.Response.Accounts) > 0 {
		c.Timezone = authResponse.Response.Accounts[0].Timezone
	}
	if !authResponse.Response.Auth {
		return nil, errors.New(authResponse.Response.Error)
	}
	return c, nil
}
