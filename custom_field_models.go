package amocrm

type (
	CustomField struct {
		ID     int           `json:"id"`
		Values []CustomValue `json:"values"`
	}

	CustomValue struct {
		Value   string `json:"value"`
		Enum    int    `json:"enum"`
		Subtype string `json:"subtype,omitempty"`
	}
)
