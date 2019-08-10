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

	NotePost struct {
		ElementID         int         `json:"element_id"`
		ElementType       int         `json:"element_type"`
		Text              string      `json:"text,omitempty"`
		NoteType          int         `json:"note_type"`
		CreatedAt         string      `json:"created_at,omitempty"`
		UpdatedAt         int         `json:"updated_at,omitempty"`
		ResponsibleUserID int         `json:"responsible_user_id,omitempty"`
		CreatedBy         int         `json:"created_by,omitempty"`
		Params            *Parameters `json:"params,omitempty"`
	}

	AddNoteRequest struct {
		Add []*NotePost `json:"add"`
	}

	GetNoteResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*Note `json:"items"`
		} `json:"_embedded"`
		Response *AmoError `json:"response"`
	}

	Note struct {
		ID                int `json:"id"`
		ResponsibleUserID int `json:"id"`
	}
)
