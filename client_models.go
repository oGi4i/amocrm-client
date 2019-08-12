package amocrm

import "net/http"

type (
	ClientInfo struct {
		userLogin string         `validate:"required"`
		apiHash   string         `validate:"required"`
		Timezone  string         `validate:"required"`
		Url       string         `validate:"required"`
		Cookie    []*http.Cookie `validate:"required,dive,required"`
	}

	PostResponse struct {
		ID        int `json:"id" validate:"omitempty"`
		RequestID int `json:"request_id" validate:"omitempty"`
		Embedded  struct {
			Items []struct {
				ID int `json:"id" validate:"omitempty"`
			} `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response" validate:"omitempty"`
	}
)
