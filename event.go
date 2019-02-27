package amocrm

import (
	"encoding/json"
	"errors"
)

func (c *ClientInfo) AddEvent(event Event) (EventResponse, error) {
	emptyResponse := EventResponse{}
	if event.PhoneNumber == "" {
		return emptyResponse, errors.New("PhoneNumber is empty")
	}
	if event.Type != "phone_call" {
		if event.PhoneNumber == "" {
			return emptyResponse, errors.New("Type not valid")
		}
	}
	url := c.Url + apiUrls["events"]
	resp, err := c.DoPost(url, EventSetRequest{Add: []Event{event}})
	if err != nil {
		return emptyResponse, err
	}
	response := EventGetResponse{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		return emptyResponse, err
	}
	if len(response.Embedded.Items) == 0 {
		return emptyResponse, errors.New("No Items")
	}
	return response.Embedded.Items[0], nil
}
