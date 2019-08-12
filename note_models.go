package amocrm

type (
	NotePostParameters struct {
		UNIQ       string `json:"UNIQ" validate:"required"`
		LINK       string `json:"LINK" validate:"required"`
		PHONE      string `json:"PHONE" validate:"required"`
		DURATION   int    `json:"DURATION" validate:"required"`
		SRC        string `json:"SRC" validate:"required"`
		FROM       string `json:"FROM,omitempty" validate:"omitempty"`
		CallStatus int    `json:"call_status" validate:"required"`
		CallResult string `json:"call_result,omitempty" validate:"omitempty"`
		Text       string `json:"text,omitempty" validate:"omitempty"`
	}

	NoteAdd struct {
		ElementID         int                 `json:"element_id" validate:"required"`
		ElementType       int                 `json:"element_type" validate:"oneof=1 2 3 4 12"`
		Text              string              `json:"text,omitempty" validate:"omitempty"`
		NoteType          int                 `json:"note_type" validate:"omitempty"`
		CreatedAt         string              `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt         int                 `json:"updated_at,omitempty" validate:"omitempty"`
		ResponsibleUserID int                 `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CreatedBy         int                 `json:"created_by,omitempty" validate:"omitempty"`
		Params            *NotePostParameters `json:"params,omitempty" validate:"omitempty"`
	}

	AddNoteRequest struct {
		Add []*NoteAdd `json:"add" validate:"required"`
	}

	GetNoteResponse struct {
		Links    *Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items []*Note `json:"items" validate:"required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response" validate:"omitempty"`
	}

	Note struct {
		ID                int             `json:"id" validate:"required"`
		CreatedBy         int             `json:"created_by" validate:"required"`
		AccountID         int             `json:"account_id" validate:"required"`
		GroupID           int             `json:"group_id" validate:"omitempty"`
		IsEditable        bool            `json:"is_editable" validate:"omitempty"`
		ElementID         int             `json:"element_id" validate:"required"`
		ElementType       int             `json:"element_type" validate:"oneof=1 2 3 4 12"`
		Text              string          `json:"text" validate:"required"`
		NoteType          int             `json:"note_type" validate:"required"`
		CreatedAt         int             `json:"created_at" validate:"required"`
		UpdatedAt         int             `json:"updated_at" validate:"required"`
		ResponsibleUserID int             `json:"responsible_user_id" validate:"required"`
		Parameters        *NoteParameters `json:"params,omitempty" validate:"omitempty"`
		Links             *Links          `json:"_links" validate:"required"`
	}

	NoteParameters struct {
		Text    string `json:"text" validate:"required"`
		Service string `json:"service" validate:"required"`
	}

	NoteType struct {
		ID         int    `json:"id" validate:"required"`
		Code       string `json:"code" validate:"required"`
		IsEditable bool   `json:"is_editable" validate:"omitempty"`
	}

	NoteTask struct {
		ID    int    `json:"id" validate:"omitempty"`
		Text  string `json:"text" validate:"omitempty"`
		Links *Links `json:"_links" validate:"omitempty"`
	}
)
