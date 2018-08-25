package amocrm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

func (c *clientInfo) AddContact(contact Contact) (int, error) {
	if contact.Name == "" {
		return 0, errors.New("Name is empty")
	}

	url := c.Url + apiUrls["contacts"]
	resp, err := c.DoPost(url, ContactSetRequest{Add: []Contact{contact}})
	if err != nil {
		return 0, err
	}
	return c.GetResponseID(resp)
}

func (c *clientInfo) GetContact(reqParams RequestParams) ([]ContactResponse, error) {
	addValues := map[string]string{}
	contacts := ContactGetResponse{}
	var err error

	if reqParams.ID != 0 {
		addValues["id"] = strconv.Itoa(reqParams.ID)

	} else {
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
	}

	url := c.Url + apiUrls["contacts"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}
	rbody, _ := ioutil.ReadAll(body)
	err = json.Unmarshal(rbody, &contacts)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(contacts)
	return contacts._embedded.Items, err
}
