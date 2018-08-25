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
		Links struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
		Embedded struct {
			Items []ContactResponse `json:"items"`
		} `json:"_embedded"`
	}
	ContactResponse struct {
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
			Links struct {
				Self struct {
					Href   string `json:"href"`
					Method string `json:"method"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"company"`
		Leads struct {
			ID    []int `json:"id"`
			Links struct {
				Self struct {
					Href   string `json:"href"`
					Method string `json:"method"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"leads"`
		ClosestTaskAt int `json:"closest_task_at"`
		Tags          struct {
		} `json:"tags"`
		CustomFields []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Values []struct {
				Value string `json:"value"`
				Enum  string `json:"enum"`
			} `json:"values"`
			IsSystem bool `json:"is_system"`
		} `json:"custom_fields"`
		Customers struct {
		} `json:"customers"`
		Links struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
	}
)
