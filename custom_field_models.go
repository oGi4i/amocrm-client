package amocrm

const (
	TextCustomFieldType = iota + 1
	NumericCustomFieldType
	CheckboxCustomFieldType
	SelectCustomFieldType
	MultiSelectCustomFieldType
	DateCustomFieldType
	UrlCustomFieldType
	MultiTextCustomFieldType
	TextAreaCustomFieldType
	RadioButtonCustomFieldType
	StreetAddressCustomFieldType
	SmartAddressCustomFieldType
	BirthDayCustomFieldType
	LegalEntityCustomFieldType
	ItemsCustomFieldType
	OrgLegalNameCustomFieldType
)

type (
	CustomField struct {
		ID       int            `json:"id" validate:"required"`
		Name     string         `json:"name" validate:"required"`
		Values   []*CustomValue `json:"values" validate:"required,dive,required"`
		IsSystem bool           `json:"is_system" validate:"omitempty"`
	}

	CustomValue struct {
		Value   string `json:"value" validate:"required"`
		Enum    int    `json:"enum,omitempty" validate:"omitempty"`
		Subtype string `json:"subtype,omitempty" validate:"omitempty"`
	}

	CustomFieldPost struct {
		ID     int           `json:"id" validate:"required"`
		Values []interface{} `json:"values" validate:"required,dive,required"`
	}

	CustomValuePost struct {
		Value   string `json:"value" validate:"required"`
		Subtype string `json:"subtype,omitempty" validate:"omitempty"`
	}

	GetCustomField struct {
		ID          int    `json:"id" validate:"required"`
		Name        string `json:"name" validate:"required"`
		FieldType   int    `json:"field_type" validate:"required"`
		Sort        int    `json:"sort" validate:"required"`
		Code        string `json:"code" validate:"omitempty"`
		IsMultiple  bool   `json:"is_multiple" validate:"omitempty"`
		IsSystem    bool   `json:"is_system" validate:"omitempty"`
		IsEditable  bool   `json:"is_editable" validate:"omitempty"`
		IsRequired  bool   `json:"is_required" validate:"omitempty"`
		IsDeletable bool   `json:"is_deletable" validate:"omitempty"`
		IsVisible   bool   `json:"is_visible" validate:"omitempty"`
		Params      struct {
		} `json:"params" validate:"omitempty"`
		Enums map[string]string `json:"enums" validate:"omitempty,dive,required"`
	}
)
