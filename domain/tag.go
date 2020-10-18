package domain

type Tag struct {
	ID   uint64 `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
