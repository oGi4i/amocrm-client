package amocrm

type (
	PipelineRequestParams struct {
		ID int `validate:"omitempty"`
	}

	GetPipelineResponse struct {
		Links    *Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items map[string]*Pipeline `json:"items" validate:"required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response" validate:"omitempty"`
	}

	Pipeline struct {
		ID       int                        `json:"id" validate:"required"`
		Name     string                     `json:"name" validate:"required"`
		Sort     int                        `json:"sort" validate:"required"`
		IsMain   bool                       `json:"is_main" validate:"omitempty"`
		Statuses map[string]*PipelineStatus `json:"statuses" validate:"required,dive,required"`
		Links    *Links                     `json:"_links" validate:"required"`
	}

	PipelineStatus struct {
		ID         int    `json:"id" validate:"required"`
		Name       string `json:"name" validate:"required"`
		Color      string `json:"color" validate:"required"`
		Sort       int    `json:"sort" validate:"required"`
		IsEditable bool   `json:"is_editable" validate:"omitempty"`
	}
)
