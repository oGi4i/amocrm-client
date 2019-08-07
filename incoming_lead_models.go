package amocrm

type (
	IncomingLeadRequestParams struct {
		SourceName string
		SourceUID  string
	}

	IncomingLead struct {
		SourceName       string            `json:"source_name"`
		SourceUID        string            `json:"source_uid"`
		CreatedAt        int64             `json:"created_at,omitempty"`
		PipelineID       string            `json:"pipeline_id,omitempty"`
		IncomingEntities *IncomingEntities `json:"incoming_entities"`
		IncomingLeadInfo *IncomingLeadInfo `json:"incoming_lead_info,omitempty"`
	}

	IncomingLeadRequest struct {
		Add []*IncomingLead `json:"add"`
	}

	IncomingEntities struct {
		Leads     []*IncomingLeadParams  `json:"leads"`
		Contacts  []*IncomingLeadContact `json:"contacts,omitempty"`
		Companies []*IncomingLeadCompany `json:"companies,omitempty"`
	}

	IncomingLeadParams struct {
		Name              string `json:"name"`
		CreatedAt         string `json:"created_at,omitempty"`
		StatusID          string `json:"status_id,omitempty"`
		ResponsibleUserID string `json:"responsible_user_id,omitempty"`
		Price             string `json:"price,omitempty"`
		Tags              string `json:"tags,omitempty"`
		Notes             []struct {
			NoteType    string `json:"note_type,omitempty"`
			ElementType string `json:"element_type,omitempty"`
			Text        string `json:"text,omitempty"`
		} `json:"notes,omitempty"`
		CustomFields []*CustomField `json:"custom_fields,omitempty"`
	}

	IncomingLeadContact struct {
		Name              string         `json:"name"`
		CustomFields      []*CustomField `json:"custom_fields"`
		ResponsibleUserID string         `json:"responsible_user_id"`
		DateCreate        string         `json:"date_create"`
	}

	IncomingLeadCompany struct {
		Name string `json:"name"`
	}

	IncomingLeadInfo struct {
		To          string `json:"to"`
		From        string `json:"from"`
		DateCall    int64  `json:"date_call"`
		Duration    string `json:"duration"`
		Link        string `json:"link"`
		ServiceCode string `json:"service_code"`
		Uniq        string `json:"uniq"`
		AddNote     bool   `json:"add_note,omitempty"`
	}

	IncomingLeadResponse struct {
		Status    string   `json:"status"`
		Data      []string `json:"data"`
		Links     *Links   `json:"_links"`
		Error     string   `json:"error,omitempty"`
		ErrorCode int      `json:"error_code,omitempty"`
	}
)
