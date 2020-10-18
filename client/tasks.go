package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
)

type (
	GetTasksOrderBy string

	GetTasksOrder struct {
		By     GetTasksOrderBy     `validate:"required,oneof=id updated_at"`
		Method request.OrderMethod `validate:"required,oneof=asc desc"`
	}

	GetTasksRequestFilter struct {
		ID                *request.Filter `validate:"omitempty"`                        // Фильтр по ID контактов
		UpdatedAt         *request.Filter `validate:"omitempty"`                        // Фильтр по дате изменения контакта
		ResponsibleUserID *request.Filter `validate:"omitempty"`                        // Фильтр по ID пользователя, ответственного за контакт
		IsCompleted       *request.Filter `validate:"omitempty"`                        // Фильтр по статусу задачи
		TaskType          *request.Filter `validate:"omitempty"`                        // Фильтр по ID типа задачи
		EntityType        *request.Filter `validate:"omitempty,required_with=EntityID"` // Фильтр по типу привязанной к задаче сущности
		EntityID          *request.Filter `validate:"omitempty"`                        // Фильтр по ID привязанной к задаче сущности. Необходимо использовать вместе с фильтром по типу привязанной сущности
	}

	GetTasksRequestParams struct {
		Page   uint64                 `validate:"omitempty"`         // Страница выборки
		Limit  uint64                 `validate:"omitempty,lte=250"` // Количество возвращаемых сущностей за один запрос (максимум - 250)
		Filter *GetTasksRequestFilter `validate:"omitempty"`         // Фильтр
		Order  *GetTasksOrder         `validate:"omitempty"`         // Сортировка результатов
	}

	GetTasksResponseEmbedded struct {
		Tasks []*domain.Task `json:"tasks" validate:"omitempty,dive,required"`
	}

	GetTasksResponse struct {
		Page          uint64                    `json:"_page" validate:"required"`
		Links         *domain.Links             `json:"_links" validate:"required"`
		Embedded      *GetTasksResponseEmbedded `json:"_embedded" validate:"omitempty"`
		ErrorResponse *domain.AmoError          `json:"response" validate:"omitempty"`
	}

	AddTasksRequestData struct {
		ResponsibleUserID uint64             `json:"responsible_user_id,omitempty" validate:"omitempty"`
		EntityID          uint64             `json:"entity_id,omitempty" validate:"omitempty"`
		EntityType        domain.EntityType  `json:"entity_type,omitempty" validate:"omitempty,oneof=leads contacts companies customers"`
		IsCompleted       bool               `json:"is_completed,omitempty" validate:"omitempty"`
		TaskTypeID        domain.TaskType    `json:"task_type_id,omitempty" validate:"omitempty,oneof=1 2"`
		Text              string             `json:"text" validate:"required"`
		Duration          uint64             `json:"duration,omitempty" validate:"omitempty"`
		CompleteTill      uint64             `json:"complete_till" validate:"required"`
		Result            *domain.TaskResult `json:"result,omitempty" validate:"omitempty"`
		CratedBy          uint64             `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy         bool               `json:"updated_by,omitempty" validate:"omitempty"`
		CreatedAt         bool               `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt         bool               `json:"updated_at,omitempty" validate:"omitempty"`
		RequestID         string             `json:"request_id,omitempty" validate:"omitempty"`
	}

	AddTasksRequest struct {
		Add []*AddTasksRequestData `validate:"required,gt=0,dive,required"`
	}

	AddTasksResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		RequestID string        `json:"request_id" validate:"required"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	AddTasksResponseEmbedded struct {
		Tasks []*AddTasksResponseItem `json:"tasks" validate:"required,gt=0,dive,required"`
	}

	AddTasksResponse struct {
		Links    *domain.Links             `json:"_links" validate:"required"`
		Embedded *AddTasksResponseEmbedded `json:"_embedded" validate:"required"`
	}

	UpdateTasksRequestData struct {
		ID                uint64             `json:"id" validate:"required"`
		ResponsibleUserID uint64             `json:"responsible_user_id,omitempty" validate:"omitempty"`
		EntityID          uint64             `json:"entity_id,omitempty" validate:"omitempty"`
		EntityType        domain.EntityType  `json:"entity_type,omitempty" validate:"omitempty,oneof=leads contacts companies customers"`
		IsCompleted       bool               `json:"is_completed,omitempty" validate:"omitempty"`
		TaskTypeID        domain.TaskType    `json:"task_type_id,omitempty" validate:"omitempty,oneof=1 2"`
		Text              string             `json:"text,omitempty" validate:"omitempty"`
		Duration          uint64             `json:"duration,omitempty" validate:"omitempty"`
		CompleteTill      uint64             `json:"complete_till,omitempty" validate:"omitempty"`
		Result            *domain.TaskResult `json:"result,omitempty" validate:"omitempty"`
		CratedBy          uint64             `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy         bool               `json:"updated_by,omitempty" validate:"omitempty"`
		CreatedAt         bool               `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt         bool               `json:"updated_at,omitempty" validate:"omitempty"`
		RequestID         string             `json:"request_id,omitempty" validate:"omitempty"`
	}

	UpdateTasksRequest struct {
		Update []*UpdateTasksRequestData `validate:"required,gt=0,dive,required"`
	}

	UpdateTasksResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		UpdatedAt uint64        `json:"updated_at" validate:"required"`
		RequestID string        `json:"request_id" validate:"required"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	UpdateTasksResponseEmbedded struct {
		Tasks []*UpdateTasksResponseItem `json:"tasks" validate:"required,gt=0,dive,required"`
	}

	UpdateTasksResponse struct {
		Links    *domain.Links                `json:"_links" validate:"required"`
		Embedded *UpdateTasksResponseEmbedded `json:"_embedded" validate:"required"`
	}
)

const (
	IDGetTasksOrderBy           GetTasksOrderBy = "id"            // Сортировка по ID задачи
	CreatedAtGetTasksOrderBy    GetTasksOrderBy = "created_at"    // Сортировка по дате создания задачи
	CompleteTillGetTasksOrderBy GetTasksOrderBy = "complete_till" // Сортировка по сроку выполнения задачи
)

func (o *GetTasksOrder) appendToQuery(params url.Values) {
	params.Add(fmt.Sprintf("order[%s]", string(o.By)), string(o.Method))
}

func (f *GetTasksRequestFilter) validate() error {
	if f.ID != nil && !f.ID.IsSimpleFilter() && !f.ID.IsMultipleFilter() {
		return errors.New("ID filter must be simple or multiple type")
	}

	if f.UpdatedAt != nil && !f.UpdatedAt.IsSimpleFilter() && !f.UpdatedAt.IsIntervalFilter() {
		return errors.New("UpdatedAt filter must be simple or interval type")
	}

	if f.ResponsibleUserID != nil && !f.ResponsibleUserID.IsSimpleFilter() && !f.ResponsibleUserID.IsMultipleFilter() {
		return errors.New("ResponsibleUserID filter must be simple or multiple type")
	}

	if f.IsCompleted != nil && !f.IsCompleted.IsSimpleFilter() {
		return errors.New("IsCompleted filter must be simple type")
	}

	if f.TaskType != nil && !f.TaskType.IsSimpleFilter() && !f.TaskType.IsMultipleFilter() {
		return errors.New("TaskType filter must be simple or multiple type")
	}

	if f.EntityType != nil && !f.EntityType.IsSimpleFilter() {
		return errors.New("EntityType filter must be simple type")
	}

	if f.EntityID != nil && !f.EntityID.IsSimpleFilter() && !f.EntityID.IsMultipleFilter() {
		return errors.New("EntityID filter must be simple or multiple type")
	}

	return nil
}

func (f *GetTasksRequestFilter) appendGetRequestFilter(params url.Values) {
	if f.ID != nil {
		f.ID.AppendToQuery(params)
	}

	if f.UpdatedAt != nil {
		f.UpdatedAt.AppendToQuery(params)
	}

	if f.ResponsibleUserID != nil {
		f.ResponsibleUserID.AppendToQuery(params)
	}

	if f.IsCompleted != nil {
		f.IsCompleted.AppendToQuery(params)
	}

	if f.TaskType != nil {
		f.TaskType.AppendToQuery(params)
	}

	if f.EntityType != nil {
		f.EntityType.AppendToQuery(params)
	}

	if f.EntityID != nil {
		f.EntityID.AppendToQuery(params)
	}
}

func (c *Client) AddTasks(ctx context.Context, req *AddTasksRequest) ([]*AddTasksResponseItem, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+tasksURI, http.MethodPost, req.Add)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(AddTasksResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response.Embedded.Tasks, nil
}

func (c *Client) UpdateTasks(ctx context.Context, req *UpdateTasksRequest) ([]*UpdateTasksResponseItem, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+tasksURI, http.MethodPatch, req.Update)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(UpdateTasksResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response.Embedded.Tasks, nil
}

func (c *Client) UpdateTask(ctx context.Context, taskID uint64, req *UpdateTasksRequestData) (*UpdateTasksResponseItem, error) {
	if taskID == 0 {
		return nil, ErrInvalidTaskID
	}

	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+tasksURI+"/"+strconv.FormatUint(taskID, 10), http.MethodPatch, req)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(UpdateTasksResponseItem)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) CompleteTask(ctx context.Context, taskID uint64, result string) (*UpdateTasksResponseItem, error) {
	if taskID == 0 {
		return nil, ErrInvalidTaskID
	}

	if result == "" {
		return nil, ErrInvalidTaskResult
	}

	return c.UpdateTask(ctx, taskID, &UpdateTasksRequestData{ID: taskID, IsCompleted: true, Result: &domain.TaskResult{Text: result}})
}

func (c *Client) GetTasks(ctx context.Context, reqParams *GetTasksRequestParams) ([]*domain.Task, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	params := make(url.Values)
	if reqParams.Page != 0 {
		params.Add("page", strconv.FormatUint(reqParams.Page, 10))
	}
	if reqParams.Limit != 0 {
		params.Add("limit", strconv.FormatUint(reqParams.Limit, 10))
	}
	if reqParams.Filter != nil {
		err := reqParams.Filter.validate()
		if err != nil {
			return nil, err
		}

		reqParams.Filter.appendGetRequestFilter(params)
	}
	if reqParams.Order != nil {
		reqParams.Order.appendToQuery(params)
	}

	body, err := c.doGet(ctx, c.baseURL+tasksURI, params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(GetTasksResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	if len(response.Embedded.Tasks) == 0 {
		return nil, ErrEmptyResponse
	}

	return response.Embedded.Tasks, nil
}

func (c *Client) GetTaskByID(ctx context.Context, taskID uint64) (*domain.Task, error) {
	if taskID == 0 {
		return nil, ErrInvalidTaskID
	}

	body, err := c.doGet(ctx, c.baseURL+tasksURI+"/"+strconv.FormatUint(taskID, 10), nil)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(domain.Task)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
