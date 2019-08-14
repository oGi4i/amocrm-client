package amocrm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	ContactTaskElementType  = 1
	LeadTaskElementType     = 2
	CompanyTaskElementType  = 3
	CustomerTaskElementType = 12

	CompletedTasksTaskFilterStatus  = 1
	InProgressTasksTaskFilterStatus = 0
)

func (c *ClientInfo) AddTask(task *TaskAdd) (int, error) {
	if err := Validate.Struct(task); err != nil {
		return 0, err
	}

	url := c.Url + apiUrls["tasks"]
	resp, err := c.DoPost(url, &AddTaskRequest{Add: []*TaskAdd{task}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) UpdateTask(task *TaskUpdate) (int, error) {
	if err := Validate.Struct(task); err != nil {
		return 0, err
	}

	url := c.Url + apiUrls["tasks"]
	resp, err := c.DoPost(url, &UpdateTaskRequest{Update: []*TaskUpdate{task}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) GetTask(reqParams *TaskRequestParams) ([]*Task, error) {
	addValues := make(map[string]string)
	tasks := new(GetTaskResponse)
	if err := Validate.Struct(reqParams); err != nil {
		return nil, err
	}

	if reqParams.ID != nil {
		addValues["id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ID)), ","), "[]")
	}
	if reqParams.LimitRows != 0 {
		addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
		if reqParams.LimitOffset != 0 {
			addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
		}
	}
	if reqParams.ElementID != nil {
		addValues["element_id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ElementID)), ","), "[]")
	}
	if reqParams.ResponsibleUserID != 0 {
		addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
	}
	if reqParams.Type != "" {
		addValues["type"] = reqParams.Type
	}

	url := c.Url + apiUrls["tasks"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, tasks)
	if err != nil {
		return nil, err
	}

	if tasks.Response != nil {
		return nil, tasks.Response
	}

	if err := Validate.Struct(tasks); err != nil {
		return nil, err
	}

	return tasks.Embedded.Items, err
}
