package amocrm

import (
	"context"
	"encoding/json"
	"strconv"
)

type (
	TaskRequestType string

	TaskRequestParams struct {
		ID                []int              `validate:"omitempty,gt=0,dive,required"`
		LimitRows         int                `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset       int                `validate:"omitempty"`
		ElementID         []int              `validate:"omitempty,gt=0,dive,required"`
		ResponsibleUserID int                `validate:"omitempty"`
		Type              TaskRequestType    `validate:"omitempty,oneof=lead contact company customer"`
		Filter            *TaskRequestFilter `validate:"omitempty"`
	}

	TaskRequestStatusFilter int

	TaskRequestFilter struct {
		Status   TaskRequestStatusFilter `validate:"omitempty,oneof=1 0"`
		TaskType []int                   `validate:"omitempty,gt=0,dive,required"`
	}

	TaskElementType int

	TaskAdd struct {
		ElementID         int             `json:"element_id,string" validate:"required"`
		ElementType       TaskElementType `json:"element_type,string" validate:"oneof=1 2 3 12"`
		CompleteTill      int             `json:"complete_till,omitempty" validate:"omitempty"`
		TaskType          int             `json:"task_type,string" validate:"required"`
		Text              string          `json:"text,omitempty" validate:"omitempty"`
		CreatedAt         int             `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int             `json:"updated_at,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int             `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		IsCompleted       bool            `json:"is_completed,omitempty" validate:"omitempty"`
		CreatedBy         int             `json:"created_by,string,omitempty" validate:"omitempty"`
		RequestID         int             `json:"request_id,string,omitempty" validate:"omitempty"`
	}

	TaskUpdate struct {
		ID                int             `json:"id,string" validate:"required"`
		ElementID         int             `json:"element_id,string,omitempty" validate:"omitempty"`
		ElementType       TaskElementType `json:"element_type,string,omitempty" validate:"omitempty,oneof=1 2 3 12"`
		CompleteTill      int             `json:"complete_till,omitempty" validate:"omitempty"`
		TaskType          int             `json:"task_type,string,omitempty" validate:"omitempty"`
		Text              string          `json:"text" validate:"omitempty"`
		CreatedAt         int             `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int             `json:"updated_at,string" validate:"required"`
		ResponsibleUserID int             `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		IsCompleted       bool            `json:"is_completed,omitempty" validate:"omitempty"`
		CreatedBy         int             `json:"created_by,string,omitempty" validate:"omitempty"`
		RequestID         int             `json:"request_id,string,omitempty" validate:"omitempty"`
	}

	AddTaskRequest struct {
		Add []*TaskAdd `json:"add" validate:"required,dive,required"`
	}

	UpdateTaskRequest struct {
		Update []*TaskUpdate `json:"update" validate:"required,dive,required"`
	}

	GetTaskResponse struct {
		Links    *Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items []*Task `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response" validate:"omitempty"`
	}

	Task struct {
		ID                int             `json:"id" validate:"required"`
		ElementID         int             `json:"element_id" validate:"required"`
		ElementType       TaskElementType `json:"element_type" validate:"oneof=1 2 3 12"`
		CompleteTillAt    int             `json:"complete_till_at" validate:"required"`
		TaskType          int             `json:"task_type" validate:"required"`
		Text              string          `json:"text" validate:"omitempty"`
		CreatedAt         int             `json:"created_at" validate:"required"`
		UpdatedAt         int             `json:"updated_at" validate:"required"`
		ResponsibleUserID int             `json:"responsible_user_id" validate:"required"`
		IsCompleted       bool            `json:"is_completed" validate:"omitempty"`
		CreatedBy         int             `json:"created_by" validate:"required"`
		AccountID         int             `json:"account_id" validate:"required"`
		GroupID           int             `json:"group_id" validate:"omitempty"`
		Result            *NoteTask       `json:"result" validate:"omitempty"`
		Links             *Links          `json:"_links" validate:"required"`
	}

	TaskType struct {
		ID   int    `json:"id" validate:"required"`
		Name string `json:"name" validate:"required"`
	}
)

const (
	ContactTaskElementType  TaskElementType = 1
	LeadTaskElementType     TaskElementType = 2
	CompanyTaskElementType  TaskElementType = 3
	CustomerTaskElementType TaskElementType = 12

	ContactTaskType  TaskRequestType = "contact"
	LeadTaskType     TaskRequestType = "lead"
	CompanyTaskType  TaskRequestType = "company"
	CustomerTaskType TaskRequestType = "customer"

	CompletedStatusTaskFilter  TaskRequestStatusFilter = 1
	InProgressStatusTaskFilter TaskRequestStatusFilter = 0
)

func (c *Client) AddTask(ctx context.Context, task *TaskAdd) (int, error) {
	if err := c.validator.Struct(task); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+tasksURI, &AddTaskRequest{Add: []*TaskAdd{task}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) UpdateTask(ctx context.Context, task *TaskUpdate) (int, error) {
	if err := c.validator.Struct(task); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+tasksURI, &UpdateTaskRequest{Update: []*TaskUpdate{task}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetTasks(ctx context.Context, reqParams *TaskRequestParams) ([]*Task, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(map[string]string)
	if reqParams.ID != nil {
		addValues["id"] = joinIntSlice(reqParams.ID)
	}
	if reqParams.LimitRows != 0 {
		addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
		if reqParams.LimitOffset != 0 {
			addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
		}
	}
	if reqParams.ElementID != nil {
		addValues["element_id"] = joinIntSlice(reqParams.ElementID)
	}
	if reqParams.ResponsibleUserID != 0 {
		addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
	}
	if reqParams.Type != "" {
		addValues["type"] = string(reqParams.Type)
	}

	body, err := c.doGet(ctx, c.baseURL+tasksURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	taskResponse := new(GetTaskResponse)
	err = json.Unmarshal(body, taskResponse)
	if err != nil {
		return nil, err
	}

	if taskResponse.Response != nil {
		return nil, taskResponse.Response
	}

	err = c.validator.Struct(taskResponse)
	if err != nil {
		return nil, err
	}

	if len(taskResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponseItems
	}

	return taskResponse.Embedded.Items, nil
}
