package amocrm

type (
	PipelineRequestParams struct {
		Id string
	}

	PipelineResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items map[string]Pipeline `json:"items"`
		} `json:"_embedded"`
		Response *AmoError `json:"response"`
	}

	Pipeline struct {
		ID       int                       `json:"id"`
		Name     string                    `json:"name"`
		Sort     int                       `json:"sort"`
		IsMain   bool                      `json:"is_main"`
		Statuses map[string]PipelineStatus `json:"statuses"`
		Links    *Links                    `json:"_links"`
	}

	PipelineStatus struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Color      string `json:"color"`
		Sort       int    `json:"sort"`
		IsEditable bool   `json:"is_editable"`
	}
)
