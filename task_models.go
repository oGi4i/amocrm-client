package amocrm

type (
	TaskRequestParams struct {
		ID                []int
		LimitRows         int
		LimitOffset       int
		ElementID         []int
		ResponsibleUserID int
		Type              string
		Filter            *TaskRequestFilter
	}

	TaskRequestFilter struct {
		Status   int
		TaskType []int
	}

	TaskPost struct {
		ElementID         string `json:"element_id"`
		ElementType       string `json:"element_type"`
		CompleteTill      string `json:"complete_till,omitempty"`
		TaskType          string `json:"task_type"`
		Text              string `json:"text"`
		CreatedAt         string `json:"created_at,omitempty"`
		UpdatedAt         string `json:"updated_at,omitempty"`
		ResponsibleUserID string `json:"responsible_user_id,omitempty"`
		IsCompleted       bool   `json:"is_completed,omitempty"`
		CreatedBy         string `json:"created_by,omitempty"`
		RequestID         string `json:"request_id,omitempty"`
	}

	AddTaskRequest struct {
		Add []*TaskPost `json:"add"`
	}

	GetTaskResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*Task `json:"items"`
		} `json:"_embedded"`
		Response *AmoError `json:"response"`
	}

	Task struct {
		ID                int           `json:"id"`
		ElementID         int           `json:"element_id"`
		ElementType       int           `json:"element_type"`
		CompleteTillAt    int           `json:"complete_till_at"`
		TaskType          int           `json:"task_type"`
		Text              string        `json:"text"`
		CreatedAt         int           `json:"created_at"`
		UpdatedAt         int           `json:"updated_at"`
		ResponsibleUserID int           `json:"responsible_user_id"`
		IsCompleted       bool          `json:"is_completed"`
		CreatedBy         int           `json:"created_by"`
		AccountID         int           `json:"account_id"`
		GroupID           int           `json:"group_id"`
		Result            []interface{} `json:"result"`
		Links             *Links        `json:"_links"`
	}

	TaskType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
