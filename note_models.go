package amocrm

type (
	Parameters struct {
		UNIQ       string `json:"UNIQ"`
		LINK       string `json:"LINK"`
		PHONE      string `json:"PHONE"`
		DURATION   int    `json:"DURATION"`
		SRC        string `json:"SRC"`
		FROM       string `json:"FROM,omitempty"`
		CallStatus int    `json:"call_status"`
		CallResult string `json:"call_result,omitempty"`
		Text       string `json:"text,omitempty"`
	}

	Note struct {
		ElementID         int        `json:"element_id"`
		ElementType       int        `json:"element_type"`
		Text              string     `json:"text,omitempty"`
		NoteType          int        `json:"note_type"`
		CreatedAt         string     `json:"created_at,omitempty"`
		UpdatedAt         int        `json:"updated_at,omitempty"`
		ResponsibleUserID int        `json:"responsible_user_id,omitempty"`
		CreatedBy         int        `json:"created_by,omitempty"`
		Params            Parameters `json:"params,omitempty"`
	}

	NoteSetRequest struct {
		Add []Note `json:"add"`
	}

	NoteGetResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []NoteResponse `json:"items"`
		} `json:"_embedded"`
	}

	NoteResponse struct {
		ID                int `json:"id"`
		ResponsibleUserID int `json:"id"`
	}
)
