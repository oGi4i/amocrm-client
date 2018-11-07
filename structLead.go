package amocrm

type (
	Lead struct {
		Name              string   `json:"name"`
		CreatedAt         string   `json:"created_at,omitempty"`
		UpdatedAt         string   `json:"updated_at,omitempty"`
		StatusID          string   `json:"status_id"`
		ResponsibleUserID string   `json:"responsible_user_id,omitempty"`
		Sale              string   `json:"sale,omitempty"`
		Tags              string   `json:"tags,omitempty"`
		ContactsID        []string `json:"contacts_id,omitempty"`
		CompanyID         string   `json:"company_id,omitempty"`
		RequestID         string   `json:"request_id,omitempty"`
	}
	LeadSetRequest struct {
		Add []Lead `json:"add"`
	}

	LeadGetResponse struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
		Embedded struct {
			Items []LeadResponse `json:"items"`
		} `json:"_embedded"`
	}

	LeadResponse struct {
		ID        int `json:"id"`
		RequestID int `json:"request_id"`
		Link      struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_link"`
	}
)
