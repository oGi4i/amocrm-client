package client

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ogi4i/amocrm-client/domain"
)

type (
	AddPipelineStatusesRequestData struct {
		Name      string                     `json:"name,omitempty" validate:"omitempty"`
		Sort      uint64                     `json:"sort,omitempty" validate:"omitempty"`
		Color     domain.PipelineStatusColor `json:"color,omitempty" validate:"omitempty,oneof=#fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 #d6eaff #c1e0ff #98cbff #ebffb1 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea"`
		RequestID string                     `json:"request_id,omitempty" validate:"omitempty"`
	}

	AddPipelineStatusesRequest struct {
		Add []*AddPipelineStatusesRequestData `validate:"required,gt=0,dive,required"`
	}

	AddPipelineStatusesResponse struct {
		TotalItems    uint64                   `json:"_total_items" validate:"required"`
		Embedded      *domain.PipelineEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError         `json:"response" validate:"omitempty"`
	}

	UpdatePipelineStatusRequest struct {
		Name  string                     `json:"name,omitempty" validate:"omitempty"`
		Sort  uint64                     `json:"sort,omitempty" validate:"omitempty"`
		Color domain.PipelineStatusColor `json:"color,omitempty" validate:"omitempty,oneof=#fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 #d6eaff #c1e0ff #98cbff #ebffb1 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea"`
	}
)

func (c *Client) GetPipelineStatuses(ctx context.Context, pipelineID uint64) ([]*domain.PipelineStatus, error) {
	if pipelineID == 0 {
		return nil, ErrInvalidPipelineID
	}

	body, err := c.doGet(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10)+"/statuses", nil)
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

	return response.Embedded.Pipelines[0].Embedded.Statuses, nil
}

func (c *Client) GetPipelineStatusByID(ctx context.Context, pipelineID, statusID uint64) (*domain.PipelineStatus, error) {
	if pipelineID == 0 {
		return nil, ErrInvalidPipelineID
	}

	if statusID == 0 {
		return nil, ErrInvalidPipelineStatusID
	}

	body, err := c.doGet(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10)+"/statuses/"+strconv.FormatUint(statusID, 10), nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(domain.PipelineStatus)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) AddPipelineStatuses(ctx context.Context, pipelineID uint64, req *AddPipelineStatusesRequest) ([]*domain.PipelineStatus, error) {
	if pipelineID == 0 {
		return nil, ErrInvalidPipelineID
	}

	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10)+"/statuses", http.MethodPost, req.Add)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(AddPipelineStatusesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response.Embedded.Statuses, nil
}

func (c *Client) UpdatePipelineStatus(ctx context.Context, pipelineID, statusID uint64, req *UpdatePipelineStatusRequest) (*domain.PipelineStatus, error) {
	if pipelineID == 0 {
		return nil, ErrInvalidPipelineID
	}

	if statusID == 0 {
		return nil, ErrInvalidPipelineStatusID
	}

	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10)+"/statuses/"+strconv.FormatUint(statusID, 10), http.MethodPatch, req)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(domain.PipelineStatus)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DeletePipelineStatus(ctx context.Context, pipelineID, statusID uint64) error {
	if pipelineID == 0 {
		return ErrInvalidPipelineID
	}

	if statusID == 0 {
		return ErrInvalidPipelineStatusID
	}

	_, err := c.do(ctx, c.baseURL+pipelinesURI+"/"+strconv.FormatUint(pipelineID, 10)+"/statuses/"+strconv.FormatUint(statusID, 10), http.MethodDelete, nil)
	if err != nil {
		return err
	}
	return nil
}
