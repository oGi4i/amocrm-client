package amocrm

import "encoding/json"

func (c *ClientInfo) GetAccount(reqParams *AccountRequestParams) (*AccountResponse, error) {
	addValues := map[string]string{}
	account := new(AccountResponse)
	var err error
	if reqParams.With != "" {
		addValues["with"] = reqParams.With
	}
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
	return account, nil
}
