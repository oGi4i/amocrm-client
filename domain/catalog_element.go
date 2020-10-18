package domain

type (
	CatalogElementMetadata struct {
		Quantity  uint64 `json:"quantity" validate:"required"`
		CatalogID uint64 `json:"catalog_id" validate:"required"`
	}

	CatalogElement struct {
		ID       uint64                  `json:"id" validate:"required"`
		Metadata *CatalogElementMetadata `json:"metadata" validate:"required"`
	}
)
