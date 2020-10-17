package domain

type Company struct {
	ID int64 `json:"id" validate:"required"`
}
