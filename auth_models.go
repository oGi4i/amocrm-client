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
			Auth     bool `json:"auth"`
			Accounts []struct {
				ID        int    `json:"id"`
				Name      string `json:"name"`
				Subdomain string `json:"subdomain"`
				Language  string `json:"language"`
				Timezone  string `json:"timezone"`
			} `json:"accounts"`
			User struct {
				ID       int    `json:"id"`
				Language string `json:"language"`
			} `json:"user"`
			ServerTime int    `json:"server_time"`
			Error      string `json:"error"`
		} `json:"response"`
	}
)
