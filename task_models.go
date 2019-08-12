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

	TaskAdd struct {
		ElementID         int    `json:"element_id,string"`
		ElementType       int    `json:"element_type,string"`
		CompleteTill      int    `json:"complete_till,omitempty"`
		TaskType          int    `json:"task_type,string"`
		Text              string `json:"text"`
		CreatedAt         int    `json:"created_at,string,omitempty"`
		UpdatedAt         int    `json:"updated_at,string,omitempty"`
		ResponsibleUserID int    `json:"responsible_user_id,string,omitempty"`
		IsCompleted       bool   `json:"is_completed,omitempty"`
		CreatedBy         int    `json:"created_by,string,omitempty"`
		RequestID         int    `json:"request_id,string,omitempty"`
	}

	AddTaskRequest struct {
		Add []*TaskAdd `json:"add"`
	}

	GetTaskResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*Task `json:"items"`
		} `json:"_embedded"`
		Response *AmoError `json:"response"`
	}

	Task struct {
		ID                int       `json:"id"`
		ElementID         int       `json:"element_id"`
		ElementType       int       `json:"element_type"`
		CompleteTillAt    int       `json:"complete_till_at"`
		TaskType          int       `json:"task_type"`
		Text              string    `json:"text"`
		CreatedAt         int       `json:"created_at"`
		UpdatedAt         int       `json:"updated_at"`
		ResponsibleUserID int       `json:"responsible_user_id"`
		IsCompleted       bool      `json:"is_completed"`
		CreatedBy         int       `json:"created_by"`
		AccountID         int       `json:"account_id"`
		GroupID           int       `json:"group_id"`
		Result            *NoteTask `json:"result"`
		Links             *Links    `json:"_links"`
	}

	TaskType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)
