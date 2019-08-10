package amocrm

type (
	EventPost struct {
		PhoneNumber string   `json:"phone_number"`
		Type        string   `json:"type"`
		Users       []string `json:"users,omitempty"`
	}

	AddEventRequest struct {
		Add []*EventPost `json:"add"`
	}

	AddEventResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*Event `json:"items"`
		} `json:"_embedded"`
	}

	Event struct {
		ElementType int    `json:"element_type"`
		ElementID   int    `json:"element_id"`
		UID         string `json:"uid"`
		PhoneNumber string `json:"phone_number"`
		QueryString string `json:"query_string"`
	}
)
