package amocrm

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (c *ClientInfo) AddTask(task *TaskPost) (int, error) {
	if task.ElementID == "" {
		return 0, errors.New("elementID is empty")
	}
	if task.ElementType == "" {
		return 0, errors.New("elementType is empty")
	}
	if task.TaskType == "" {
		return 0, errors.New("taskType is empty")
	}
	if task.Text == "" {
		return 0, errors.New("text is empty")
	}

	url := c.Url + apiUrls["tasks"]
	resp, err := c.DoPost(url, &AddTaskRequest{Add: []*TaskPost{task}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) GetTask(reqParams *TaskRequestParams) ([]*Task, error) {
	addValues := make(map[string]string)
	tasks := new(GetTaskResponse)
	var err error

	if len(reqParams.ID) > 0 {
		addValues["id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ID)), ","), "[]")
	} else {
		if reqParams.LimitRows != 0 {
			addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
			if reqParams.LimitOffset != 0 {
				addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
			}
		}
		if len(reqParams.ElementID) > 0 {
			addValues["element_id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ElementID)), ","), "[]")
		}
		if reqParams.ResponsibleUserID != 0 {
			addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
		}
		if reqParams.Type != "" {
			addValues["type"] = reqParams.Type
		}
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

	return tasks.Embedded.Items, err
}
