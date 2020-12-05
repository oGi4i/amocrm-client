package domain

type LossReason struct {
	ID   uint64 `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}
