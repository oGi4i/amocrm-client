package amocrm

import (
	"encoding/json"
	"strings"
)

const (
	AccountWithCustomFields  = "custom_fields"
	AccountWithUsers         = "users"
	AccountWithMessenger     = "messenger"
	AccountWithNotifications = "notifications"
	AccountWithPipelines     = "pipelines"
	AccountWithGroups        = "groups"
	AccountWithNoteTypes     = "note_types"
	AccountWithTaskTypes     = "task_types"
)

func (c *ClientInfo) GetAccount(reqParams *AccountRequestParams) (*AccountResponse, error) {
	addValues := map[string]string{}
	account := new(AccountResponse)

	if err := Validate.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues["with"] = strings.Join(reqParams.With, ",")

	url := c.Url + apiUrls["account"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, account)
	if err != nil {
		amoError := new(AmoError)
		err = json.Unmarshal(body, amoError)
		if err != nil {
			return nil, err
		}

		return nil, amoError
	}

	if err := Validate.Struct(account); err != nil {
		return nil, err
	}

	return account, nil
}
