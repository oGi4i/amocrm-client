package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ogi4i/amocrm-client/domain"
)

type (
	GetPipelinesResponseEmbedded struct {
		Pipelines []*domain.Pipeline `json:"pipelines" validate:"required,gt=0,dive,required"`
	}

	GetPipelinesResponse struct {
		TotalItems    uint64                        `json:"_total_items" validate:"required"`
		Links         *domain.Links                 `json:"_links" validate:"required"`
		Embedded      *GetPipelinesResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError              `json:"response" validate:"omitempty"`
	}

	AddPipelineData struct {
		Name         string                   `json:"name" validate:"omitempty"`
		Sort         uint64                   `json:"sort" validate:"omitempty"`
		IsMain       bool                     `json:"is_main" validate:"omitempty"`
		IsUnsortedOn bool                     `json:"is_unsorted_on" validate:"omitempty"`
		RequestID    string                   `json:"request_id,omitempty" validate:"omitempty"`
		Embedded     *domain.PipelineEmbedded `json:"_embedded" validate:"omitempty,gt=0,dive,required"`
	}

	UpdatePipelineData struct {
		Name         string `json:"name" validate:"omitempty"`
		Sort         uint64 `json:"sort" validate:"omitempty"`
		IsMain       bool   `json:"is_main" validate:"omitempty"`
		IsUnsortedOn bool   `json:"is_unsorted_on" validate:"omitempty"`
	}
)

func (c *Client) GetPipelines(ctx context.Context) ([]*domain.Pipeline, error) {
	body, err := c.doGet(ctx, c.baseURL+pipelinesURI, nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	response := new(GetPipelinesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response.Embedded.Pipelines, nil
}

func (c *Client) GetPipelineByID(ctx context.Context, pipelineID uint64) (*domain.Pipeline, error) {
	body, err := c.doGet(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10), nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	response := new(domain.Pipeline)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) AddPipelines(ctx context.Context, pipelines []*AddPipelineData) ([]*domain.Pipeline, error) {
	body, err := c.do(ctx, c.baseURL+pipelinesURI, http.MethodPost, pipelines)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	response := new(GetPipelinesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response.Embedded.Pipelines, nil
}

func (c *Client) UpdatePipeline(ctx context.Context, pipelineID uint64, pipeline *UpdatePipelineData) (*domain.Pipeline, error) {
	body, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10), http.MethodPatch, pipeline)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	response := new(domain.Pipeline)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeletePipeline(ctx context.Context, pipelineID uint64) error {
	_, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10), http.MethodDelete, nil)
	if err != nil {
		return err
	}
	return nil
}
