package amocrm

type (
	Parameters struct {
		UNIQ       string `json:"uniq"`
		LINK       string `json:"link"`
		PHONE      string `json:"phone"`
		DURATION   int    `json:"duation"`
		SRC        string `json:"src"`
		FROM       string `json:"from"`
		CallStatus int    `json:"call_status"`
		CallResult string `json:"call_result,omitempty"`
		Text       string `json:"text,omitempty"`
	}
	Note struct {
		ElementID         int        `json:"element_id"`
		ElementType       int        `json:"element_type"`
		Text              string     `json:"text"`
		NoteType          int        `json:"note_type"`
		CreatedAt         string     `json:"created_at,omitempty"`
		UpdatedAt         int        `json:"updated_at,omitempty"`
		ResponsibleUserID int        `json:"responsible_user_id,omitempty"`
		Params            Parameters `json:"params,omitempty"`
	}
	NoteSetRequest struct {
		Add []Note `json:"add"`
	}
)
