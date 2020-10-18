package domain

type (
	TaskType uint8

	TaskResult struct {
		Text string `json:"text" validate:"required"` // Тексту результата выполнения задачи
	}

	Task struct {
		ID                uint64      `json:"id" validate:"required"`                                          // ID задачи
		CreatedBy         uint64      `json:"created_by" validate:"required"`                                  // ID пользователя, создавшего задачу
		UpdatedBy         uint64      `json:"updated_by" validate:"required"`                                  // ID пользователя, изменившего задачу
		CreatedAt         uint64      `json:"created_at" validate:"required"`                                  // Дата создания задачи, перадаётся в Unix Timestamp
		UpdatedAt         uint64      `json:"updated_at" validate:"required"`                                  // Дата изменения задачи, передаётся в Unix Timestamp
		ResponsibleUserID uint64      `json:"responsible_user_id" validate:"required"`                         // ID пользователя, ответственного за задачу
		GroupID           uint64      `json:"group_id" validate:"required"`                                    // ID группы, в которой состоит ответственный
		EntityID          uint64      `json:"entity_id" validate:"required"`                                   // ID сущности, к которой привязана задача
		EntityType        EntityType  `json:"entity_type" validate:"oneof=leads contacts companies customers"` // Тип сущности, к которой привязана задача
		IsCompleted       bool        `json:"is_completed" validate:"omitempty"`                               // Признак выполнена ли задача
		TaskTypeID        TaskType    `json:"task_type_id" validate:"oneof=1 2"`                               // Тип задачи
		Text              string      `json:"text" validate:"omitempty"`                                       // Описание задачи
		Duration          uint64      `json:"duration" validate:"omitempty"`                                   // Длительность задачи в секундах
		CompleteTill      uint64      `json:"complete_till" validate:"required"`                               // Дата, когда задача должна быть завершена, передаётся в Unix Timestamp
		Result            *TaskResult `json:"result" validate:"omitempty"`                                     // Результат выполнения задачи
		AccountID         uint64      `json:"account_id" validate:"required"`                                  // ID аккаунта, в котором находится задача
		Links             *Links      `json:"_links" validate:"required"`
	}

	TaskTypeInfo struct {
		ID   uint64 `json:"id" validate:"required"`
		Name string `json:"name" validate:"required"`
		Code string `json:"code" validate:"required"`
	}
)

const (
	CallTaskType TaskType = iota + 1
	MeetingTaskType
)
