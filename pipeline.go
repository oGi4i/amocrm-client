package amocrm

import (
	"context"
	"encoding/json"
	"strconv"
)

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

func (c *Client) GetPipelines(ctx context.Context, reqParams *PipelineRequestParams) (map[string]*Pipeline, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(map[string]string)
	if reqParams.ID != 0 {
		addValues["id"] = strconv.Itoa(reqParams.ID)
	}

	body, err := c.doGet(ctx, c.baseURL+pipelinesURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	pipelineResponse := new(GetPipelineResponse)
	err = json.Unmarshal(body, pipelineResponse)
	if err != nil {
		return nil, err
	}

	if pipelineResponse.Response != nil {
		return nil, pipelineResponse.Response
	}

	if err := c.validator.Struct(pipelineResponse); err != nil {
		return nil, err
	}

	if len(pipelineResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponseItems
	}

	return pipelineResponse.Embedded.Items, nil
}
