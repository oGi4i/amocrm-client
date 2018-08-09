package amocrm

type (
	Contact struct {
		Name              string        `json:"name"`
		CreatedAt         int           `json:"created_at,omitempty"`
		UpdatedAt         int           `json:"updated_at,omitempty"`
		ResponsibleUserID int           `json:"responsible_user_id,omitempty"`
		CreatedBy         int           `json:"created_by,omitempty"`
		CompanyName       string        `json:"company_name,omitempty"`
		Tags              string        `json:"tags,omitempty"`
		LeadsID           string        `json:"leads_id,omitempty"`
		CustomersID       string        `json:"customers_id,omitempty"`
		CompanyID         string        `json:"company_id,omitempty"`
		CustomFields      []CustomField `json:"custom_fields,omitempty"`
	}
	ContactSetRequest struct {
		Add []Contact `json:"add"`
	}
	CustomField struct {
		ID     int           `json:"id"`
		Values []CustomValue `json:"values"`
	}
	CustomValue struct {
		Value   string `json:"value"`
		Enum    string `json:"enum"`
		Subtype string `json:"subtype"`
	}
	ContactGetResponse struct {
		Embedded struct {
			Items []ContactResponse `json:"items"`
		} `json:"_embedded"`
	}
	ContactResponse struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		CreatedAt         int    `json:"created_at,omitempty"`
		UpdatedAt         int    `json:"updated_at,omitempty"`
		ResponsibleUserID int    `json:"responsible_user_id,omitempty"`
		CreatedBy         int    `json:"created_by,omitempty"`
		AccountID         int    `json:"account_id,omitempty"`
		UpdatedBy         int    `json:"updated_by,omitempty"`
		GroupID           int    `json:"group_id,omitempty"`
		Company           struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"company,omitempty"`
		Leads struct {
			id map[string]int `json:"id"`
		} `json:"leads,omitempty"`
		ClosestTaskAt int `json:"closest_task_at"`
	}
)
