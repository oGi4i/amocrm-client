package client

import (
	"context"
	"encoding/json"
	"github.com/ogi4i/amocrm-client/domain"
	"net/http"
	"net/url"
)

func (c *Client) AddNote(ctx context.Context, note *domain.NoteAdd) (int, error) {
	if err := c.validator.Struct(note); err != nil {
		return 0, err
	}

	resp, err := c.do(ctx, c.baseURL+notesURI, http.MethodPost, &domain.AddNoteRequest{Add: []*domain.NoteAdd{note}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetNotes(ctx context.Context, reqParams *domain.NoteRequestParams) ([]*domain.Note, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(url.Values)
	//addValues["type"] = string(reqParams.Type)
	//if reqParams.ID != nil {
	//	addValues["id"] = joinIntSlice(reqParams.ID)
	//}
	//if reqParams.LimitRows != 0 {
	//	addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
	//	if reqParams.LimitOffset != 0 {
	//		addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
	//	}
	//}
	//if reqParams.ElementID != nil {
	//	addValues["element_id"] = joinIntSlice(reqParams.ElementID)
	//}
	//if reqParams.NoteType != nil {
	//	addValues["note_type"] = joinIntSlice(reqParams.NoteType)
	//}

	body, err := c.doGet(ctx, c.baseURL+notesURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	noteResponse := new(domain.GetNoteResponse)
	err = json.Unmarshal(body, noteResponse)
	if err != nil {
		return nil, err
	}

	if noteResponse.ErrorResponse != nil {
		return nil, noteResponse.ErrorResponse
	}

	err = c.validator.Struct(noteResponse)
	if err != nil {
		return nil, err
	}

	if len(noteResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponse
	}

	return noteResponse.Embedded.Items, nil
}
