package amocrm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	ContactNoteType = "contact"
	LeadNoteType    = "lead"
	CompanyNoteType = "company"
	TaskNoteType    = "task"
)

func (c *ClientInfo) AddNote(note *NoteAdd) (int, error) {
	if err := Validate.Struct(note); err != nil {
		return 0, err
	}

	url := c.Url + apiUrls["notes"]
	fmt.Println(note)
	resp, err := c.DoPost(url, &AddNoteRequest{Add: []*NoteAdd{note}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}

func (c *ClientInfo) GetNote(reqParams *NoteRequestParams) ([]*Note, error) {
	addValues := make(map[string]string)
	noteResponse := new(GetNoteResponse)

	if err := Validate.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues["type"] = reqParams.Type
	if reqParams.ID != nil {
		addValues["id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ID)), ","), "[]")
	}
	if reqParams.LimitRows != 0 {
		addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
		if reqParams.LimitOffset != 0 {
			addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
		}
	}
	if reqParams.ElementID != nil {
		addValues["element_id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ElementID)), ","), "[]")
	}
	if reqParams.NoteType != nil {
		addValues["note_type"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.NoteType)), ","), "[]")
	}

	url := c.Url + apiUrls["notes"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, noteResponse)
	if err != nil {
		return nil, err
	}

	if noteResponse.Response != nil {
		return nil, noteResponse.Response
	}

	if err := Validate.Struct(noteResponse); err != nil {
		return nil, err
	}

	return noteResponse.Embedded.Items, err
}

func (c *ClientInfo) DownloadAttachment(attachment string) ([]byte, error) {
	url := c.Url + apiUrls["download"] + attachment
	response, err := c.DoGet(url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
