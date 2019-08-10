package amocrm

type (
	AuthResponse struct {
		Response struct {
			Auth       bool           `json:"auth"`
			Accounts   []*AuthAccount `json:"accounts"`
			User       *AuthUser      `json:"user"`
			ServerTime int            `json:"server_time"`
			Error      string         `json:"error"`
		} `json:"response"`
	}
)
