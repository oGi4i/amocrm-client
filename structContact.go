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
		_embedded struct {
			Items []ContactResponse `json:"items"`
		} `json:"_embedded"`
		_links struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
	}
	ContactResponse struct {
		_links struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
		AccountID     int `json:"account_id"`
		ClosestTaskAt int `json:"closest_task_at"`
		Company       struct {
			_links struct {
				Self struct {
					Href   string `json:"href"`
					Method string `json:"method"`
				} `json:"self"`
			} `json:"_links"`
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"company"`
		CreatedAt    int `json:"created_at"`
		CreatedBy    int `json:"created_by"`
		CustomFields []struct {
			ID       int    `json:"id"`
			IsSystem bool   `json:"is_system"`
			Name     string `json:"name"`
			Values   []struct {
				Enum  string `json:"enum"`
				Value string `json:"value"`
			} `json:"values"`
		} `json:"custom_fields"`
		Customers struct{} `json:"customers"`
		GroupID   int      `json:"group_id"`
		ID        int      `json:"id"`
		Leads     struct {
			_links struct {
				Self struct {
					Href   string `json:"href"`
					Method string `json:"method"`
				} `json:"self"`
			} `json:"_links"`
			ID []int `json:"id"`
		} `json:"leads"`
		Name              string   `json:"name"`
		ResponsibleUserID int      `json:"responsible_user_id"`
		Tags              struct{} `json:"tags"`
		UpdatedAt         int      `json:"updated_at"`
		UpdatedBy         int      `json:"updated_by"`
	}
)
