package amocrm

type (
	TaskRequestParams struct {
		ID                []int              `validate:"omitempty,gt=0,dive,required"`
		LimitRows         int                `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset       int                `validate:"omitempty"`
		ElementID         []int              `validate:"omitempty,gt=0,dive,required"`
		ResponsibleUserID int                `validate:"omitempty"`
		Type              string             `validate:"omitempty,oneof=lead contact company customer"`
		Filter            *TaskRequestFilter `validate:"omitempty"`
	}

	TaskRequestFilter struct {
		Status   int   `validate:"omitempty,oneof=1 0"`
		TaskType []int `validate:"omitempty,gt=0,dive,required"`
	}

	TaskAdd struct {
		ElementID         int    `json:"element_id,string" validate:"required"`
		ElementType       int    `json:"element_type,string" validate:"oneof=1 2 3 12"`
		CompleteTill      int    `json:"complete_till,omitempty" validate:"omitempty"`
		TaskType          int    `json:"task_type,string" validate:"required"`
		Text              string `json:"text,omitempty" validate:"omitempty"`
		CreatedAt         int    `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int    `json:"updated_at,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int    `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		IsCompleted       bool   `json:"is_completed,omitempty" validate:"omitempty"`
		CreatedBy         int    `json:"created_by,string,omitempty" validate:"omitempty"`
		RequestID         int    `json:"request_id,string,omitempty" validate:"omitempty"`
	}

	TaskUpdate struct {
		ID                int    `json:"id,string" validate:"required"`
		ElementID         int    `json:"element_id,string,omitempty" validate:"omitempty"`
		ElementType       int    `json:"element_type,string,omitempty" validate:"omitempty,oneof=1 2 3 12"`
		CompleteTill      int    `json:"complete_till,omitempty" validate:"omitempty"`
		TaskType          int    `json:"task_type,string,omitempty" validate:"omitempty"`
		Text              string `json:"text" validate:"omitempty"`
		CreatedAt         int    `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int    `json:"updated_at,string" validate:"required"`
		ResponsibleUserID int    `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		IsCompleted       bool   `json:"is_completed,omitempty" validate:"omitempty"`
		CreatedBy         int    `json:"created_by,string,omitempty" validate:"omitempty"`
		RequestID         int    `json:"request_id,string,omitempty" validate:"omitempty"`
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
		ID                int       `json:"id" validate:"required"`
		ElementID         int       `json:"element_id" validate:"required"`
		ElementType       int       `json:"element_type" validate:"oneof=1 2 3 12"`
		CompleteTillAt    int       `json:"complete_till_at" validate:"required"`
		TaskType          int       `json:"task_type" validate:"required"`
		Text              string    `json:"text" validate:"omitempty"`
		CreatedAt         int       `json:"created_at" validate:"required"`
		UpdatedAt         int       `json:"updated_at" validate:"required"`
		ResponsibleUserID int       `json:"responsible_user_id" validate:"required"`
		IsCompleted       bool      `json:"is_completed" validate:"omitempty"`
		CreatedBy         int       `json:"created_by" validate:"required"`
		AccountID         int       `json:"account_id" validate:"required"`
		GroupID           int       `json:"group_id" validate:"omitempty"`
		Result            *NoteTask `json:"result" validate:"omitempty"`
		Links             *Links    `json:"_links" validate:"required"`
	}

	TaskType struct {
		ID   int    `json:"id" validate:"required"`
		Name string `json:"name" validate:"required"`
	}
)
