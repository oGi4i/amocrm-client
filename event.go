package amocrm

import (
	"encoding/json"
	"errors"
)

func (c *ClientInfo) AddEvent(event *Event) (*EventResponse, error) {
	emptyResponse := new(EventResponse)
	if event.PhoneNumber == "" {
		return emptyResponse, errors.New("phoneNumber is empty")
	}
	if event.Type != "phone_call" {
		if event.PhoneNumber == "" {
			return emptyResponse, errors.New("type not valid")
		}
	}
	url := c.Url + apiUrls["events"]
	resp, err := c.DoPost(url, &EventSetRequest{Add: []*Event{event}})
	if err != nil {
		return emptyResponse, err
	}
	response := new(EventGetResponse)
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&response)
	if err != nil {
		return emptyResponse, err
	}
	if len(response.Embedded.Items) == 0 {
		return emptyResponse, errors.New("no Items")
	}
	return response.Embedded.Items[0], nil
}
