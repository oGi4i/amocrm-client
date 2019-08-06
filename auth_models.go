package amocrm

import "net/http"

type (
	// Информация о подключении к аккаунту
	ClientInfo struct {
		userLogin string
		apiHash   string
		Timezone  string
		Url       string
		Cookie    []*http.Cookie
	}

	//AuthResponse Структура ответа при авторизации
	AuthResponse struct {
		Response struct {
			Auth       bool       `json:"auth"`
			Accounts   []*Account `json:"accounts"`
			User       *User      `json:"user"`
			ServerTime int        `json:"server_time"`
			Error      string     `json:"error"`
		} `json:"response"`
	}

	Account struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Subdomain string `json:"subdomain"`
		Language  string `json:"language"`
		Timezone  string `json:"timezone"`
	}

	User struct {
		ID       int    `json:"id"`
		Language string `json:"language"`
	}
)
