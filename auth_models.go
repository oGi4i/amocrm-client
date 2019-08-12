package amocrm

type (
	AuthResponse struct {
		Response struct {
			Auth       bool           `json:"auth" validate:"required"`
			Accounts   []*AuthAccount `json:"accounts" validate:"omitempty,dive,required"`
			User       *AuthUser      `json:"user" validate:"required"`
			ServerTime int            `json:"server_time" validate:"required"`
			Error      string         `json:"error" validate:"omitempty"`
		} `json:"response" validate:"required"`
	}
)
