package amocrm

import (
	"encoding/json"
	"errors"
)

func (c *ClientInfo) AddEvent(event *EventPost) (*Event, error) {
	emptyResponse := new(Event)
	if event.PhoneNumber == "" {
		return emptyResponse, errors.New("phoneNumber is empty")
	}
	if event.Type != "phone_call" {
		if event.PhoneNumber == "" {
			return emptyResponse, errors.New("type not valid")
		}
	}
	url := c.Url + apiUrls["events"]
	resp, err := c.DoPost(url, &AddEventRequest{Add: []*EventPost{event}})
	if err != nil {
		return emptyResponse, err
	}
	response := new(AddEventResponse)
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
