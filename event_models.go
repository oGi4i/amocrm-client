package amocrm

type (
	Event struct {
		PhoneNumber string   `json:"phone_number"`
		Type        string   `json:"type"`
		Users       []string `json:"users,omitempty"`
	}

	EventSetRequest struct {
		Add []Event `json:"add"`
	}

	EventGetResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []EventResponse `json:"items"`
		} `json:"_embedded"`
	}

	EventResponse struct {
		ElementType int    `json:"element_type"`
		ElementID   int    `json:"element_id"`
		UID         string `json:"uid"`
		PhoneNumber string `json:"phone_number"`
		QueryString string `json:"query_string"`
	}
)
