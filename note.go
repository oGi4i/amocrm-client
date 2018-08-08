package amocrm

import "errors"

func (c *clientInfo) AddNote(note Note) (int, error) {
	if note.Text == "" {
		return 0, errors.New("Text is empty")
	}
	if note.ElementID == 0 {
		return 0, errors.New("ElementID is empty")
	}
	if note.ElementType == 0 {
		return 0, errors.New("ElementType is empty")
	}
	if note.NoteType == 0 {
		return 0, errors.New("NoteType is empty")
	}
	url := c.SetURL("note", nil)
	resp, err := c.DoPost(url, NoteSetRequest{Add: []Note{note}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}
