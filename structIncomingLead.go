package amocrm

type (
	IncomingLeadRequestParams struct {
		SourceName string
		SourceUID  string
	}
	IncomingLead struct {
		SourceName       string           `json:"source_name"`
		SourceUID        string           `json:"source_uid"`
		CreatedAt        string           `json:"created_at,omitempty"`
		PipelineID       string           `json:"pipeline_id,omitempty"`
		IncomingEntities IncomingEntities `json:"incoming_entities"`
		IncomingLeadInfo IncomingLeadInfo `json:"incoming_lead_info,omitempty"`
	}
	IncomingLeadRequest struct {
		Add []IncomingLead `json:"add"`
	}
	IncomingEntities struct {
		Leads     []IncomingLeadParams  `json:"leads"`
		Contacts  []IncomingLeadContact `json:"contacts"`
		Companies []IncomingLeadCompany `json:"companies"`
	}
	IncomingLeadParams struct {
		Name              string `json:"name"`
		CreatedAt         string `json:"created_at"`
		StatusID          string `json:"status_id"`
		ResponsibleUserID string `json:"responsible_user_id"`
		Price             string `json:"price"`
		Tags              string `json:"tags"`
		Notes             []struct {
			NoteType    string `json:"note_type"`
			ElementType string `json:"element_type"`
			Text        string `json:"text"`
		} `json:"notes"`
		CustomFields []struct {
			ID     string   `json:"id"`
			Values []string `json:"values"`
		} `json:"custom_fields"`
	}

	IncomingLeadContact struct {
		Name         string `json:"name"`
		CustomFields []struct {
			ID     string `json:"id"`
			Values []struct {
				Value string `json:"value"`
				Enum  string `json:"enum"`
			} `json:"values"`
		} `json:"custom_fields"`
		ResponsibleUserID string `json:"responsible_user_id"`
		DateCreate        string `json:"date_create"`
	}
	IncomingLeadCompany struct {
		Name string `json:"name"`
	}
	IncomingLeadInfo struct {
		To          string `json:"to"`
		From        string `json:"from"`
		DateCall    string `json:"date_call"`
		Duration    string `json:"duration"`
		Link        string `json:"link"`
		ServiceCode string `json:"service_code"`
		Uniq        string `json:"uniq"`
		AddNote     string `json:"add_note,omitempty"`
	}

	IncomingLeadResponse struct {
		Status string `json:"status"`
		Data   string `json:"data"`
		Links  struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
	}
)
