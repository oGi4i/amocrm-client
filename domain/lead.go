package domain

type (
	EmbeddedLead struct {
		ID    uint64 `json:"id" validate:"required"`
		Links *Links `json:"_links" validate:"required"`
	}

	LeadEmbedded struct {
		LossReasons     []*LossReason      `json:"loss_reason" validate:"omitempty,dive,required"`                // Причина отказа сделки (требует параметра `with=loss_reason` в запросе)
		Tags            []*Tag             `json:"tags,omitempty" validate:"omitempty,dive,required"`             // Теги, привязанные к сделке
		Contacts        []*Contact         `json:"contacts,omitempty" validate:"omitempty,dive,required"`         // Данные контактов, привязанных к сделке (требует параметра `with=contacts` в запросе)
		Companies       []*EmbeddedCompany `json:"companies,omitempty" validate:"omitempty,dive,required"`        // Данные компании, привязанной к сделке. В данном массиве всегда 1 элемент, так как у сделки может быть только 1 компания
		CatalogElements []*CatalogElement  `json:"catalog_elements,omitempty" validate:"omitempty,dive,required"` // Данные элементов списков, привязанных к сделке (требует параметра `with=catalog_elements` в запросе)
	}

	Lead struct {
		ID                     uint64         `json:"id" validate:"required"`                              // ID сделки
		Name                   string         `json:"name" validate:"required"`                            // Название сделки
		Price                  uint64         `json:"price" validate:"required"`                           // Бюджет сделки
		ResponsibleUserID      uint64         `json:"responsible_user_id" validate:"required"`             // ID пользователя, ответственного за сделку
		GroupID                uint64         `json:"group_id,omitempty" validate:"omitempty"`             // ID группы, в которой состоит пользователь, ответственный за сделку
		StatusID               uint64         `json:"status_id" validate:"required"`                       // ID статуса, в который добавляется сделка, по-умолчанию - первый этап главной воронки
		PipelineID             uint64         `json:"pipeline_id" validate:"required"`                     // ID воронки, в которую добавляется сделка
		LossReasonID           uint64         `json:"loss_reason_id,omitempty" validate:"omitempty"`       // ID причины отказа
		SourceID               uint64         `json:"source_id,omitempty" validate:"omitempty"`            // ID источника сделки (требует параметра `with=source_id` в запросе)
		CreatedBy              uint64         `json:"created_by" validate:"required"`                      // ID пользователя, создавшего сделку
		UpdatedBy              uint64         `json:"updated_by" validate:"required"`                      // ID пользователя, изменившего сделку
		CreatedAt              uint64         `json:"created_at" validate:"required"`                      // Дата закрытия сделки, передаётся в Unit Timestamp
		UpdatedAt              uint64         `json:"updated_at" validate:"required"`                      // Дата создания сделки, передаётся в Unit Timestamp
		ClosedAt               uint64         `json:"closed_at,omitempty" validate:"omitempty"`            // Дата изменения сделки, передаётся в Unit Timestamp
		ClosestTaskAt          uint64         `json:"closest_task_at,omitempty" validate:"omitempty"`      // Дата ближайшей задачи к выполнению, передаётся в Unit Timestamp
		IsDeleted              bool           `json:"is_deleted" validate:"omitempty"`                     // Признак удалена ли сделка
		CustomFieldsValues     []*CustomField `json:"custom_fields_values,omitempty" validate:"omitempty"` // // Массив дополнительных полей, заданных для сделки
		Score                  uint64         `json:"score" validate:"omitempty"`                          // Скоринг сделки
		AccountID              uint64         `json:"account_id" validate:"required"`                      // ID аккаунта, в котором находится сделка
		IsPriceModifiedByRobot bool           `json:"is_price_modified_by_robot" validate:"omitempty"`     // Признак изменён ли в последний раз бюджет сделки роботом (требует параметра `with=is_price_modified_by_robot` в запросе)
		Links                  *Links         `json:"_links" validate:"required"`
		Embedded               *LeadEmbedded  `json:"_embedded" validate:"omitempty"`
	}
)
