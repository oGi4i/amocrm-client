package amocrm

type (
	CustomField struct {
		ID       int            `json:"id"`
		Name     string         `json:"name"`
		Values   []*CustomValue `json:"values"`
		IsSystem bool           `json:"is_system"`
	}

	CustomValue struct {
		Value   string `json:"value"`
		Enum    int    `json:"enum,omitempty"`
		Subtype string `json:"subtype,omitempty"`
	}

	CustomFieldPost struct {
		ID     int           `json:"id"`
		Values []interface{} `json:"values"`
	}

	CustomValuePost struct {
		Value   string `json:"value"`
		Subtype string `json:"subtype,omitempty"`
	}
)
