package amocrm

import (
	"errors"
	"fmt"
	"strconv"
)

func (c *clientInfo) AddContact(contact Contact) (int, error) {
	if contact.Name == "" {
		return 0, errors.New("Name is empty")
	}

	url := c.SetURL("contacts", nil)
	resp, err := c.DoPost(url, ContactSetRequest{Add: []Contact{contact}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}

func (c *clientInfo) GetContact(reqParams RequestParams) ([]ContactResponse, error) {
	fmt.Println("TEst1")
	addValues := map[string]string{}
	addValues["type"] = "json"
	contacts := ContactGetResponse{}
	var err error

	if reqParams.ID != 0 {
		addValues["id"] = strconv.Itoa(reqParams.ID)
		url := c.SetURL("contacts", addValues)
		err := c.DoGet(url, &contacts)
		return contacts.Embedded.Items, err
	}
	if reqParams.LimitRows != 0 {
		addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
		if reqParams.LimitOffset != 0 {
			addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
		}
	}
	if reqParams.ResponsibleUserID != 0 {
		addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
	}
	if reqParams.Query != "" {
		addValues["query"] = reqParams.Query
	}
	url := c.SetURL("contacts", addValues)
	err = c.DoGet(url, &contacts)
	if err != nil {
		return nil, err
	}
	return contacts.Embedded.Items, err
}
