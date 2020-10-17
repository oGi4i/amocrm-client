package domain

type CatalogElement struct {
	ID       int64 `json:"id" validate:"required"`
	Metadata struct {
		Quantity  int64 `json:"quantity" validate:"required"`
		CatalogID int64 `json:"catalog_id" validate:"required"`
	} `json:"metadata" validate:"required"`
}
