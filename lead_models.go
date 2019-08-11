package amocrm

type (
	//RequestParams параметры GET запроса
	LeadRequestParams struct {
		ID                []int
		LimitRows         int
		LimitOffset       int
		ResponsibleUserID int
		Query             string
		Status            []int
		Filter            *LeadRequestFilter
	}

	LeadRequestFilter struct {
		Tasks  int
		Active int
	}

	LeadPost struct {
		Name              string             `json:"name"`
		CreatedAt         string             `json:"created_at,omitempty"`
		UpdatedAt         string             `json:"updated_at,omitempty"`
		StatusID          string             `json:"status_id"`
		PipelineID        string             `json:"pipeline_id,omitempty"`
		ResponsibleUserID string             `json:"responsible_user_id,omitempty"`
		Sale              string             `json:"sale,omitempty"`
		Tags              string             `json:"tags,omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty"`
		ContactsID        []string           `json:"contacts_id,omitempty"`
		CompanyID         string             `json:"company_id,omitempty"`
		RequestID         string             `json:"request_id,omitempty"`
	}

	AddLeadRequest struct {
		Add []*LeadPost `json:"add"`
	}

	GetLeadResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*Lead `json:"items"`
		} `json:"_embedded"`
		Response *AmoError `json:"response"`
	}

	Lead struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		ResponsibleUserID int    `json:"responsible_user_id"`
		CreatedBy         int    `json:"created_by"`
		CreatedAt         int    `json:"created_at"`
		UpdatedAt         int    `json:"updated_at"`
		AccountID         int    `json:"account_id"`
		IsDeleted         bool   `json:"is_deleted"`
		MainContact       struct {
			ID    int    `json:"id"`
			Links *Links `json:"_links"`
		} `json:"main_contact"`
		GroupID       int            `json:"group_id"`
		ClosedAt      int            `json:"closed_at"`
		ClosestTaskAt int            `json:"closest_task_at"`
		Tags          []*Tag         `json:"tags"`
		CustomFields  []*CustomField `json:"custom_fields,omitempty"`
		Contact       struct {
			ID    []int  `json:"id"`
			Links *Links `json:"_links"`
		} `json:"contacts"`
		StatusID int `json:"status_id"`
		Sale     int `json:"sale"`
		Pipeline struct {
			ID    int    `json:"id"`
			Links *Links `json:"_links"`
		} `json:"pipeline"`
		Links *Links `json:"_links"`
	}

	LeadCustomField struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FieldType   int    `json:"field_type"`
		Sort        int    `json:"sort"`
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
