package amocrm

import (
	"context"
	"encoding/json"
	"strconv"
)

type (
	NoteRequestType string

	NoteRequestParams struct {
		Type        NoteRequestType `validate:"required,oneof=lead contact company task"`
		ID          []int           `validate:"omitempty,gt=0,dive,required"`
		LimitRows   int             `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset int             `validate:"omitempty"`
		ElementID   []int           `validate:"omitempty,gt=0,dive,required"`
		NoteType    []int           `validate:"omitempty,gt=0,dive,required"`
	}

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
		Attachment        string          `json:"attachment" validate:"omitempty"`
		Parameters        *NoteParameters `json:"params,omitempty" validate:"omitempty"`
		Links             *Links          `json:"_links" validate:"required"`
	}

	NoteParameters struct {
		Text    string `json:"TEXT" validate:"required"`
		Service string `json:"service" validate:"omitempty"`
		HTML    string `json:"HTML" validate:"omitempty"`
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

const (
	ContactNoteType NoteRequestType = "contact"
	LeadNoteType    NoteRequestType = "lead"
	CompanyNoteType NoteRequestType = "company"
	TaskNoteType    NoteRequestType = "task"
)

func (c *Client) AddNote(ctx context.Context, note *NoteAdd) (int, error) {
	if err := c.validator.Struct(note); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+notesURI, &AddNoteRequest{Add: []*NoteAdd{note}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetNotes(ctx context.Context, reqParams *NoteRequestParams) ([]*Note, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(map[string]string)
	addValues["type"] = string(reqParams.Type)
	if reqParams.ID != nil {
		addValues["id"] = joinIntSlice(reqParams.ID)
	}
	if reqParams.LimitRows != 0 {
		addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
		if reqParams.LimitOffset != 0 {
			addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
		}
	}
	if reqParams.ElementID != nil {
		addValues["element_id"] = joinIntSlice(reqParams.ElementID)
	}
	if reqParams.NoteType != nil {
		addValues["note_type"] = joinIntSlice(reqParams.NoteType)
	}

	body, err := c.doGet(ctx, c.baseURL+notesURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	noteResponse := new(GetNoteResponse)
	err = json.Unmarshal(body, noteResponse)
	if err != nil {
		return nil, err
	}

	if noteResponse.Response != nil {
		return nil, noteResponse.Response
	}

	err = c.validator.Struct(noteResponse)
	if err != nil {
		return nil, err
	}

	if len(noteResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponseItems
	}

	return noteResponse.Embedded.Items, nil
}
