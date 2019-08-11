package amocrm

type (
	//RequestParams параметры GET запроса
	ContactRequestParams struct {
		ID                int
		LimitRows         int
		LimitOffset       int
		ResponsibleUserID int
		Query             string
	}

	ContactAdd struct {
		Name              string             `json:"name"`
		CreatedAt         string             `json:"created_at,omitempty"`
		UpdatedAt         string             `json:"updated_at,omitempty"`
		ResponsibleUserID string             `json:"responsible_user_id,omitempty"`
		CreatedBy         string             `json:"created_by,omitempty"`
		CompanyName       string             `json:"company_name,omitempty"`
		Tags              string             `json:"tags,omitempty"`
		LeadsID           []string           `json:"leads_id,omitempty"`
		CustomersID       string             `json:"customers_id,omitempty"`
		CompanyID         string             `json:"company_id,omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty"`
	}

	ContactUpdate struct {
		ID                string             `json:"id"`
		Name              string             `json:"name,omitempty"`
		CreatedAt         string             `json:"created_at,omitempty"`
		UpdatedAt         string             `json:"updated_at"`
		ResponsibleUserID string             `json:"responsible_user_id,omitempty"`
		CreatedBy         string             `json:"created_by,omitempty"`
		CompanyName       string             `json:"company_name,omitempty"`
		Tags              string             `json:"tags,omitempty"`
		LeadsID           string             `json:"leads_id,omitempty"`
		CustomersID       string             `json:"customers_id,omitempty"`
		CompanyID         string             `json:"company_id,omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty"`
		Unlink            *Unlink            `json:"unlink,omitempty"`
	}

	AddContactRequest struct {
		Add []*ContactAdd `json:"add"`
	}

	UpdateContactRequest struct {
		Update []*ContactUpdate `json:"update"`
	}

	GetContactResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*Contact `json:"items"`
		} `json:"_embedded"`
		Response *AmoError `json:"response,omitempty"`
	}

	Contact struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		ResponsibleUserID int    `json:"responsible_user_id"`
		CreatedBy         int    `json:"created_by"`
		CreatedAt         int    `json:"created_at"`
		UpdatedAt         int    `json:"updated_at"`
		AccountID         int    `json:"account_id"`
		UpdatedBy         int    `json:"updated_by"`
		GroupID           int    `json:"group_id"`
		Company           struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Links *Links `json:"_links"`
		} `json:"company"`
		Leads struct {
			ID    []int  `json:"id"`
			Links *Links `json:"_links"`
		} `json:"leads"`
		ClosestTaskAt int            `json:"closest_task_at"`
		Tags          []*Tag         `json:"tags"`
		CustomFields  []*CustomField `json:"custom_fields"`
		Customers     struct {
		} `json:"customers"`
		Links *Links `json:"_links"`
	}

	ContactCustomField struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FieldType   int    `json:"field_type"`
		Sort        int    `json:"sort"`
		Code        string `json:"code"`
		IsMultiple  bool   `json:"is_multiple"`
		IsSystem    bool   `json:"is_system"`
		IsEditable  bool   `json:"is_editable"`
		IsRequired  bool   `json:"is_required"`
		IsDeletable bool   `json:"is_deletable"`
		IsVisible   bool   `json:"is_visible"`
		Params      struct {
		} `json:"params"`
		Enums map[string]string `json:"enums"`
	}
)
