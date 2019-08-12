package amocrm

import (
	"fmt"
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
