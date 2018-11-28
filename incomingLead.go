package amocrm

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *clientInfo) AddIncomingLeadCall(incominglead IncomingLead) (string, error) {
	if incominglead.SourceName == "" {
		return "0", errors.New("SourceName is empty")
	}
	if incominglead.SourceUID == "" {
		return "0", errors.New("SourceUID is empty")
	}
	if incominglead.IncomingLeadInfo.To == "" {
		return "0", errors.New("IncomingLeadInfo.To is empty")
	}
	if incominglead.IncomingLeadInfo.From == "" {
		return "0", errors.New("IncomingLeadInfo.From is empty")
	}
	if incominglead.IncomingLeadInfo.DateCall == "" {
		return "0", errors.New("IncomingLeadInfo.DateCall is empty")
	}
	if incominglead.IncomingLeadInfo.Duration == "" {
		return "0", errors.New("IncomingLeadInfo.Duration is empty")
	}
	if incominglead.IncomingLeadInfo.Duration == "" {
		return "0", errors.New("IncomingLeadInfo.Duration is empty")
	}
	if incominglead.IncomingLeadInfo.Link == "" {
		return "0", errors.New("IncomingLeadInfo.Link is empty")
	}
	if incominglead.IncomingLeadInfo.ServiceCode == "" {
		return "0", errors.New("IncomingLeadInfo.ServiceCode is empty")
	}
	if incominglead.IncomingLeadInfo.Uniq == "" {
		return "0", errors.New("IncomingLeadInfo.Uniq is empty")
	}
	url := fmt.Sprint(
		c.Url,
		apiUrls["incomingleadsip"],
		"?login=",
		c.userLogin,
		"&api_key=",
		c.apiHash,
	)
	fmt.Println(incominglead)
	resp, err := c.DoPostWithoutCookie(url, IncomingLeadRequest{Add: []IncomingLead{incominglead}})
	if err != nil {
		return "0", err
	}
	response := IncomingLeadResponse{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		return "0", err
	}
	if len(response.Data) == 0 {
		return "0", errors.New("No Items")
	}
	return response.Data, nil
}
