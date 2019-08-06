package amocrm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
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
)

//New Открытия соединения и авторизация
func New(accountURL string, login string, hash string) (*ClientInfo, error) {
	var err error

	if login == "" {
		return nil, errors.New("login is empty")
	}
	if hash == "" {
		return nil, errors.New("hash is empty")
	}
	c := &ClientInfo{
		userLogin: login,
		apiHash:   hash,
	}
	_, err = url.Parse(accountURL)
	if err != nil {
		return nil, err
	}
	c.Url = accountURL
	values := url.Values{}
	values.Set("USER_LOGIN", c.userLogin)
	values.Set("USER_HASH", c.apiHash)
	reqbody := strings.NewReader(values.Encode())
	urlString := c.Url + apiUrls["auth"]
	resp, err := http.Post(urlString, "application/x-www-form-urlencoded", reqbody)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		c.Cookie = resp.Cookies()
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
	} else {
		err = errors.New("Wrong http status: " + string(resp.StatusCode))
		return nil, err
	}
}

func (c *ClientInfo) DoGet(url string, data map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for _, cookie := range c.Cookie {
		req.AddCookie(cookie)
	}
	q := req.URL.Query()
	for key, value := range data {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *ClientInfo) DoPost(url string, data interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	fmt.Println(url)
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range c.Cookie {
		req.AddCookie(cookie)
	}
	fmt.Println(req)
	client := &http.Client{}
	return client.Do(req)
}

func (c *ClientInfo) DoPostWithoutCookie(url string, data string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	return client.Do(req)
}

func (c *ClientInfo) GetResponseID(resp *http.Response) (int, error) {
	result := respID{}
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&result)
	if err != nil {
		return 0, err
	}
	if len(result.Embedded.Items) == 0 {
		if result.Response.Error != "" {
			return 0, errors.New(result.Response.Error)
		}
		return 0, errors.New("no Items")
	}
	return result.Embedded.Items[0].ID, nil
}
