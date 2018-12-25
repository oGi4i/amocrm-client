package amocrm

type (
	PipelineRequestParams struct {
		Id string
	}
	PipelineResponse struct {
		Links struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
		Embedded struct {
			Items map[string]Pipeline `json:"items"`
		} `json:"_embedded"`
	}

	Pipeline struct {
		ID       int                       `json:"id"`
		Name     string                    `json:"name"`
		Sort     int                       `json:"sort"`
		IsMain   bool                      `json:"is_main"`
		Statuses map[string]PipelineStatus `json:"statuses"`
		Links    struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
	}

	PipelineStatus struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Color      string `json:"color"`
		Sort       int    `json:"sort"`
		IsEditable bool   `json:"is_editable"`
	}
)
