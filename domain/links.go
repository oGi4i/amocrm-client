package domain

type (
	Link struct {
		Href string `json:"href" validate:"required"`
	}

	Links struct {
		Self     *Link `json:"self" validate:"required"`
		Next     *Link `json:"next,omitempty" validate:"omitempty"`
		First    *Link `json:"first,omitempty" validate:"omitempty"`
		Previous *Link `json:"prev,omitempty" validate:"omitempty"`
	}
)
