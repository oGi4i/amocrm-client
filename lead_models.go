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

	LeadAdd struct {
		Name              string             `json:"name"`
		CreatedAt         int                `json:"created_at,string,omitempty"`
		UpdatedAt         int                `json:"updated_at,string,omitempty"`
		StatusID          int                `json:"status_id,string"`
		PipelineID        int                `json:"pipeline_id,string,omitempty"`
		ResponsibleUserID int                `json:"responsible_user_id,string,omitempty"`
		Sale              int                `json:"sale,string,omitempty"`
		Tags              string             `json:"tags,omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty"`
		ContactsID        []string           `json:"contacts_id,omitempty"`
		CompanyID         int                `json:"company_id,string,omitempty"`
		RequestID         int                `json:"request_id,string,omitempty"`
	}

	LeadUpdate struct {
		ID                int                `json:"id,string"`
		Name              string             `json:"name,omitempty"`
		CreatedAt         int                `json:"created_at,string,omitempty"`
		UpdatedAt         int                `json:"updated_at,string"`
		StatusID          int                `json:"status_id,string,omitempty"`
		PipelineID        int                `json:"pipeline_id,string,omitempty"`
		ResponsibleUserID int                `json:"responsible_user_id,string,omitempty"`
		Sale              int                `json:"sale,string,omitempty"`
		Tags              string             `json:"tags,omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty"`
		ContactsID        []string           `json:"contacts_id,omitempty"`
		CompanyID         int                `json:"company_id,string,omitempty"`
		RequestID         int                `json:"request_id,string,omitempty"`
		Unlink            *Unlink            `json:"unlink,omitempty"`
	}

	AddLeadRequest struct {
		Add []*LeadAdd `json:"add"`
	}

	UpdateLeadRequest struct {
		Update []*LeadUpdate `json:"update"`
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
