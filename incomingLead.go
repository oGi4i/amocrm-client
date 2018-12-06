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
	if incominglead.IncomingLeadInfo.DateCall == 0 {
		return "0", errors.New("IncomingLeadInfo.DateCall is empty")
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
	if incominglead.IncomingEntities.Leads[0].Name == "" {
		return "0", errors.New("IncomingEntities.Leads[0].Name is empty")
	}
	url := fmt.Sprint(
		c.Url,
		apiUrls["incomingleadsip"],
		"?login=",
		c.userLogin,
		"&api_key=",
		c.apiHash,
	)
	payload := fmt.Sprint("add", "%5B0%5D%5B", "source_name", "%5D=", incominglead.SourceName)
	payload = fmt.Sprint(payload, "&add", "%5B0%5D%5B", "source_uid", "%5D=", incominglead.SourceUID)
	if incominglead.CreatedAt != 0 {
		payload = fmt.Sprint(payload, "&add", "%5B0%5D%5B", "created_at", "%5D=", incominglead.CreatedAt)
	}
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_entities",
		"%5D%5B",
		"leads",
		"%5D%5B0%5D%5B",
		"name",
		"%5D=",
		incominglead.IncomingEntities.Leads[0].Name,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"to",
		"%5D=",
		incominglead.IncomingLeadInfo.To,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"from",
		"%5D=",
		incominglead.IncomingLeadInfo.From,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"date_call",
		"%5D=",
		incominglead.IncomingLeadInfo.DateCall,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"duration",
		"%5D=",
		incominglead.IncomingLeadInfo.Duration,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"duration",
		"%5D=",
		incominglead.IncomingLeadInfo.Duration,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"link",
		"%5D=",
		incominglead.IncomingLeadInfo.Link,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"service_code",
		"%5D=",
		incominglead.IncomingLeadInfo.ServiceCode,
	)
	payload = fmt.Sprint(
		payload,
		"&add",
		"%5B0%5D%5B",
		"incoming_lead_info",
		"%5D%5B",
		"uid",
		"%5D=",
		incominglead.IncomingLeadInfo.Uniq,
	)
	if incominglead.IncomingLeadInfo.AddNote == true {
		payload = fmt.Sprint(
			payload,
			"&add",
			"%5B0%5D%5B",
			"incoming_lead_info",
			"%5D%5B",
			"add_note",
			"%5D=",
			incominglead.IncomingLeadInfo.AddNote,
		)
	}

	resp, err := c.DoPostWithoutCookie(url, payload)
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
		return response.Error, errors.New("No Items")
	}
	return response.Status, nil
}
