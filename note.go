package amocrm

import (
	"errors"
	"fmt"
)

func (c *ClientInfo) AddNote(note *Note) (int, error) {
	if note.ElementID == 0 {
		return 0, errors.New("elementID is empty")
	}
	if note.ElementType == 0 {
		return 0, errors.New("elementType is empty")
	}
	if note.NoteType == 0 {
		return 0, errors.New("noteType is empty")
	}
	url := c.Url + apiUrls["notes"]
	fmt.Println(note)
	resp, err := c.DoPost(url, &NoteSetRequest{Add: []*Note{note}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}
