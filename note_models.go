package amocrm

type (
	NotePostParameters struct {
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
		ElementID         int                 `json:"element_id"`
		ElementType       int                 `json:"element_type"`
		Text              string              `json:"text,omitempty"`
		NoteType          int                 `json:"note_type"`
		CreatedAt         string              `json:"created_at,omitempty"`
		UpdatedAt         int                 `json:"updated_at,omitempty"`
		ResponsibleUserID int                 `json:"responsible_user_id,omitempty"`
		CreatedBy         int                 `json:"created_by,omitempty"`
		Params            *NotePostParameters `json:"params,omitempty"`
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
		ID                int             `json:"id"`
		CreatedBy         int             `json:"created_by"`
		AccountID         int             `json:"account_id"`
		GroupID           int             `json:"group_id"`
		IsEditable        bool            `json:"is_editable"`
		ElementID         int             `json:"element_id"`
		ElementType       int             `json:"element_type"`
		Text              string          `json:"text"`
		NoteType          int             `json:"note_type"`
		CreatedAt         int             `json:"created_at"`
		UpdatedAt         int             `json:"updated_at"`
		ResponsibleUserID int             `json:"responsible_user_id"`
		Parameters        *NoteParameters `json:"params,omitempty"`
		Links             *Links          `json:"_links"`
	}

	NoteParameters struct {
		Text    string `json:"text"`
		Service string `json:"service"`
	}

	NoteType struct {
		ID         int    `json:"id"`
		Code       string `json:"code"`
		IsEditable bool   `json:"is_editable"`
	}

	NoteTask struct {
		ID    int    `json:"id"`
		Text  string `json:"text"`
		Links *Links `json:"_links"`
	}
)
