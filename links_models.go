package amocrm

type (
	Links struct {
		Self struct {
			Href   string `json:"href" validate:"required"`
			Method string `json:"method" validate:"required"`
		} `json:"self" validate:"required"`
	}
)
