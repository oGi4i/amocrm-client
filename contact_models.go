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

	ContactPost struct {
		Name              string         `json:"name"`
		CreatedAt         int            `json:"created_at,omitempty"`
		UpdatedAt         int            `json:"updated_at,omitempty"`
		ResponsibleUserID int            `json:"responsible_user_id,omitempty"`
		CreatedBy         int            `json:"created_by,omitempty"`
		CompanyName       string         `json:"company_name,omitempty"`
		Tags              string         `json:"tags,omitempty"`
		LeadsID           string         `json:"leads_id,omitempty"`
		CustomersID       string         `json:"customers_id,omitempty"`
		CompanyID         string         `json:"company_id,omitempty"`
		CustomFields      []*CustomField `json:"custom_fields,omitempty"`
	}

	AddContactRequest struct {
		Add []*ContactPost `json:"add"`
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
)
