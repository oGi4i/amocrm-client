package amocrm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type (
	// Информация о подключении к аккаунту
	clientInfo struct {
		userLogin         string
		apiHash           string
		Timezone          string
		accountWebAddress *url.URL
		Timeout           time.Duration
	}
	//AuthResponse Структура ответа при авторизации
	AuthResponse struct {
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
	//respID стандартный ответ
	respID struct {
		Embedded struct {
			Items []struct {
				ID int `json:"id"`
			} `json:"items"`
		} `json:"_embedded"`
		Response struct {
			Error string `json:"error"`
		} `json:"response"`
	}
	//RequestParams параметры GET запроса
	RequestParams struct {
		ID                int
		LimitRows         int
		LimitOffset       int
		ResponsibleUserID int
		Query             string
	}
)

//New Открытия соединения и авторизация
func New(accountURL string, login string, hash string, timeout time.Duration) (*clientInfo, error) {
	var err error

	if login == "" {
		return nil, errors.New("login is empty")
	}
	if hash == "" {
		return nil, errors.New("hash is empty")
	}
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

// Установка URL и параметров по умолчанию
func (c *clientInfo) SetURL(path string, addValues map[string]string) string {
	requestURL := *c.accountWebAddress
	requestURL.Path = apiUrls[path]
	values := requestURL.Query()
	values.Set("USER_LOGIN", c.userLogin)
	values.Set("USER_HASH", c.apiHash)

	if addValues != nil {
		for key, value := range addValues {
			values.Set(key, value)
		}
	}
	requestURL.RawQuery = values.Encode()
	return requestURL.String()
}

func (c *clientInfo) DoGet(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(&result)
}

func (c *clientInfo) DoPost(url string, data interface{}) (*http.Response, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(buf)
	err := enc.Encode(data)
	fmt.Println(buf)
	if err != nil {
		return nil, err
	}
	return http.Post(url, "application/json", buf)
}

func (c *clientInfo) GetResponseID(resp *http.Response) (int, error) {
	result := respID{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&result)
	fmt.Println(result)
	if err != nil {
		return 0, err
	}
	if len(result.Embedded.Items) == 0 {
		if result.Response.Error != "" {
			return 0, errors.New(result.Response.Error)
		}
		return 0, errors.New("No Items")
	}
	return result.Embedded.Items[0].ID, nil
}
