package client

import (
	"context"
	"encoding/json"
	"github.com/ogi4i/amocrm-client/domain"
	"net/http"
	"net/url"
)

func (c *Client) AddTask(ctx context.Context, task *domain.TaskAdd) (int, error) {
	if err := c.validator.Struct(task); err != nil {
		return 0, err
	}

	resp, err := c.do(ctx, c.baseURL+tasksURI, http.MethodPost, &domain.AddTaskRequest{Add: []*domain.TaskAdd{task}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) UpdateTask(ctx context.Context, task *domain.TaskUpdate) (int, error) {
	if err := c.validator.Struct(task); err != nil {
		return 0, err
	}

	resp, err := c.do(ctx, c.baseURL+tasksURI, http.MethodPost, &domain.UpdateTaskRequest{Update: []*domain.TaskUpdate{task}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetTasks(ctx context.Context, reqParams *domain.TaskRequestParams) ([]*domain.Task, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(url.Values)
	//if reqParams.ID != nil {
	//	addValues["id"] = joinIntSlice(reqParams.ID)
	//}
	//if reqParams.LimitRows != 0 {
	//	addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
	//	if reqParams.LimitOffset != 0 {
	//		addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
	//	}
	//}
	//if reqParams.ElementID != nil {
	//	addValues["element_id"] = joinIntSlice(reqParams.ElementID)
	//}
	//if reqParams.ResponsibleUserID != 0 {
	//	addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
	//}
	//if reqParams.Type != "" {
	//	addValues["type"] = string(reqParams.Type)
	//}

	body, err := c.doGet(ctx, c.baseURL+tasksURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	taskResponse := new(domain.GetTaskResponse)
	err = json.Unmarshal(body, taskResponse)
	if err != nil {
		return nil, err
	}

	if taskResponse.ErrorResponse != nil {
		return nil, taskResponse.ErrorResponse
	}

	err = c.validator.Struct(taskResponse)
	if err != nil {
		return nil, err
	}

	if len(taskResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponse
	}

	return taskResponse.Embedded.Items, nil
}
