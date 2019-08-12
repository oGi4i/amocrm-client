package amocrm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	NoTasksLeadFilter         = 1
	InProgressTasksLeadFilter = 2

	ActiveLeadsLeadFilter = 1
)

var (
	leadArrayFields = []string{"tags", "custom_fields"}
)

func (c *ClientInfo) AddLead(lead *LeadAdd) (int, error) {
	if err := Validate.Struct(lead); err != nil {
		return 0, err
	}

	url := c.Url + apiUrls["leads"]
	resp, err := c.DoPost(url, &AddLeadRequest{Add: []*LeadAdd{lead}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) UpdateLead(lead *LeadUpdate) (int, error) {
	if err := Validate.Struct(lead); err != nil {
		return 0, err
	}

	url := c.Url + apiUrls["leads"]
	resp, err := c.DoPost(url, &UpdateLeadRequest{Update: []*LeadUpdate{lead}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) GetLead(reqParams *LeadRequestParams) ([]*Lead, error) {
	addValues := make(map[string]string)
	leadResponse := new(GetLeadResponse)

	if err := Validate.Struct(reqParams); err != nil {
		return nil, err
	}

	if reqParams.ID != nil {
		addValues["id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ID)), ","), "[]")
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
	if reqParams.Status != nil {
		addValues["status"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.Status)), ","), "[]")
	}

	url := c.Url + apiUrls["leads"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, leadResponse)
	if err != nil {
		// fix bad json serialization, where nil array is serialized as nil object
		stringBody := string(body)
		for _, s := range leadArrayFields {
			stringBody = strings.ReplaceAll(stringBody, "\""+s+"\":{}", "\""+s+"\":[]")
		}
		err = json.Unmarshal([]byte(stringBody), leadResponse)
		if err != nil {
			return nil, err
		}
	}

	if leadResponse.Response != nil {
		return nil, leadResponse.Response
	}

	if err := Validate.Struct(leadResponse); err != nil {
		return nil, err
	}

	return leadResponse.Embedded.Items, err
}
