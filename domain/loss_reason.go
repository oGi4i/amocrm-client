package domain

type LossReason struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
