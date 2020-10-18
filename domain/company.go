package domain

type (
	EmbeddedCompany struct {
		ID    uint64 `json:"id" validate:"required"`
		Links *Links `json:"_links" validate:"required"`
	}
)
