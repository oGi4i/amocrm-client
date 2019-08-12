package amocrm

type (
	Group struct {
		ID   int    `json:"id" validate:"required"`
		Name string `json:"name" validate:"required"`
	}
)
