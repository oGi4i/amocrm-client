package amocrm

type (
	IncomingLeadRequestParams struct {
		SourceName string
		SourceUID  string
	}
	IncomingLead struct {
		SourceName       string           `url:"source_name"`
		SourceUID        string           `url:"source_uid"`
		CreatedAt        int64            `url:"created_at,omitempty"`
		PipelineID       string           `url:"pipeline_id,omitempty"`
		IncomingEntities IncomingEntities `url:"incoming_entities"`
		IncomingLeadInfo IncomingLeadInfo `url:"incoming_lead_info,omitempty"`
	}
	IncomingLeadRequest struct {
		Add []IncomingLead `url:"add"`
	}
	IncomingEntities struct {
		Leads     []IncomingLeadParams  `url:"leads"`
		Contacts  []IncomingLeadContact `url:"contacts,omitempty"`
		Companies []IncomingLeadCompany `url:"companies,omitempty"`
	}
	IncomingLeadParams struct {
		Name              string `url:"name"`
		CreatedAt         string `url:"created_at,omitempty"`
		StatusID          string `url:"status_id,omitempty"`
		ResponsibleUserID string `url:"responsible_user_id,omitempty"`
		Price             string `url:"price,omitempty"`
		Tags              string `url:"tags,omitempty"`
		Notes             []struct {
			NoteType    string `url:"note_type,omitempty"`
			ElementType string `url:"element_type,omitempty"`
			Text        string `url:"text,omitempty"`
		} `url:"notes,omitempty"`
		CustomFields []struct {
			ID     string   `url:"id,omitempty"`
			Values []string `url:"values,omitempty"`
		} `url:"custom_fields,omitempty"`
	}

	IncomingLeadContact struct {
		Name         string `url:"name"`
		CustomFields []struct {
			ID     string `url:"id"`
			Values []struct {
				Value string `url:"value"`
				Enum  string `url:"enum"`
			} `url:"values"`
		} `url:"custom_fields"`
		ResponsibleUserID string `url:"responsible_user_id"`
		DateCreate        string `url:"date_create"`
	}
	IncomingLeadCompany struct {
		Name string `url:"name"`
	}
	IncomingLeadInfo struct {
		To          string `url:"to"`
		From        string `url:"from"`
		DateCall    int64  `url:"date_call"`
		Duration    string `url:"duration"`
		Link        string `url:"link"`
		ServiceCode string `url:"service_code"`
		Uniq        string `url:"uniq"`
		AddNote     string `url:"add_note,omitempty"`
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
		Error     string `json:"error"`
		ErrorCode int    `json:"error_code"`
	}
)
