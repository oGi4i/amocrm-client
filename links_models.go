package amocrm

type (
	Links struct {
		Self struct {
			Href   string `json:"href"`
			Method string `json:"method"`
		} `json:"self"`
	}
)
