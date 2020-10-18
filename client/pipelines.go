package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ogi4i/amocrm-client/domain"
)

type (
	PipelinesResponseEmbedded struct {
		Pipelines []*domain.Pipeline `json:"pipelines" validate:"omitempty,dive,required"`
	}

	PipelinesResponse struct {
		TotalItems    uint64                     `json:"_total_items" validate:"required"`
		Links         *domain.Links              `json:"_links" validate:"required"`
		Embedded      *PipelinesResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError           `json:"response" validate:"omitempty"`
	}

	AddPipelinesRequestDataEmbedded struct {
		Statuses []*domain.EmbeddedPipelineStatus `json:"statuses" validate:"required,gt=0,dive,required"`
	}

	AddPipelinesRequestData struct {
		Name         string                           `json:"name" validate:"omitempty"`
		Sort         uint64                           `json:"sort" validate:"omitempty"`
		IsMain       bool                             `json:"is_main" validate:"omitempty"`
		IsUnsortedOn bool                             `json:"is_unsorted_on" validate:"omitempty"`
		RequestID    string                           `json:"request_id,omitempty" validate:"omitempty"`
		Embedded     *AddPipelinesRequestDataEmbedded `json:"_embedded" validate:"omitempty,gt=0,dive,required"`
	}

	AddPipelinesRequest struct {
		Add []*AddPipelinesRequestData `validate:"required,gt=0,dive,required"`
	}

	AddPipelinesResponseEmbedded struct {
		Pipelines []*domain.Pipeline `json:"pipelines" validate:"required,gt=0,dive,required"`
	}

	AddPipelinesResponse struct {
		TotalItems    uint64                        `json:"_total_items" validate:"required"`
		Links         *domain.Links                 `json:"_links" validate:"required"`
		Embedded      *AddPipelinesResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError              `json:"response" validate:"omitempty"`
	}

	UpdatePipelineRequest struct {
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
		return nil, ErrEmptyResponse
	}

	response := new(PipelinesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	if len(response.Embedded.Pipelines) == 0 {
		return nil, ErrEmptyResponse
	}

	return response.Embedded.Pipelines, nil
}

func (c *Client) GetPipelineByID(ctx context.Context, pipelineID uint64) (*domain.Pipeline, error) {
	if pipelineID == 0 {
		return nil, ErrInvalidPipelineID
	}

	body, err := c.doGet(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10), nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
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

func (c *Client) AddPipelines(ctx context.Context, req *AddPipelinesRequest) ([]*domain.Pipeline, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+pipelinesURI, http.MethodPost, req.Add)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(AddPipelinesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response.Embedded.Pipelines, nil
}

func (c *Client) UpdatePipeline(ctx context.Context, pipelineID uint64, req *UpdatePipelineRequest) (*domain.Pipeline, error) {
	if pipelineID == 0 {
		return nil, ErrInvalidPipelineID
	}

	body, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10), http.MethodPatch, req)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
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
	if pipelineID == 0 {
		return ErrInvalidPipelineID
	}

	_, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10), http.MethodDelete, nil)
	if err != nil {
		return err
	}
	return nil
}
