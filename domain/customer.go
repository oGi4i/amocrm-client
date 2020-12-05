package domain

type (
	EmbeddedCustomer struct {
		ID    uint64 `json:"id" validate:"required"`
		Links *Links `json:"_links" validate:"required"`
	}

	CustomerEmbedded struct {
		Tags            []*Tag             `json:"tags,omitempty" validate:"omitempty,dive,required"`             // Теги, привязанные к покупателю
		Contacts        []*Contact         `json:"contacts,omitempty" validate:"omitempty,dive,required"`         // Контакты, привязанные к покупателю (требует параметра `with=contacts` в запросе)
		Companies       []*EmbeddedCompany `json:"companies,omitempty" validate:"omitempty,dive,required"`        // Компании, привязанные к покупателю (требует параметра `with=companies` в запросе)
		CatalogElements []*CatalogElement  `json:"catalog_elements,omitempty" validate:"omitempty,dive,required"` // Данные элементов списков, привязанных к покупателю (требует параметра `with=catalog_elements` в запросе)
	}

	Customer struct {
		ID                 uint64            `json:"id" validate:"required"`                              // ID покупателя
		Name               string            `json:"name" validate:"required"`                            // Название покупателя
		NextPrice          uint64            `json:"next_price" validate:"omitempty"`                     // Ожидаемая сумма покупки
		NextData           uint64            `json:"next_date" validate:"omitempty"`                      // Ожидаемая дата следующей покупки, передаётся в Unix Timestamp
		ResponsibleUserID  uint64            `json:"responsible_user_id" validate:"omitempty"`            // ID пользователя, ответственного за покупателя
		Periodicity        uint64            `json:"periodicity" validate:"omitempty"`                    // Периодичность. Данные необходимы для покупателей, при включенном функционале периодичности
		CreateBy           uint64            `json:"created_by" validate:"omitempty"`                     // ID пользователя, создавшего покупателя
		UpdateBy           uint64            `json:"updated_by" validate:"omitempty"`                     // ID пользователя, изменивнего данные покупателя
		CreatedAt          uint64            `json:"created_at" validate:"omitempty"`                     // Дата создания покупателя, передаётся в Unix Timestamp
		UpdatedAt          uint64            `json:"updated_at" validate:"omitempty"`                     // Дата изменения покупателя, передаётся в Unix Timestamp
		ClosestTaskAt      uint64            `json:"closest_task_at" validate:"omitempty"`                // Дата ближайшей задачи к выполнению, передаётся в Unix Timestamp
		IsDeleted          bool              `json:"is_deleted" validate:"omitempty"`                     // Признак удалён ли покупателя
		CustomFieldsValues []*CustomField    `json:"custom_fields_values,omitempty" validate:"omitempty"` // Массив дополнительных полей, заданным для покупателя
		ITV                uint64            `json:"itv" validate:"omitempty"`                            // Сумма покупок
		PurchaseCount      uint64            `json:"purchase_count" validate:"omitempty"`                 // Количество покупок
		AverageCheck       uint64            `json:"average_check" validate:"omitempty"`                  // Средний чек
		AccountID          uint64            `json:"account_id" validate:"omitempty"`                     // ID аккаунта, в котором находится покупатель
		Embedded           *CustomerEmbedded `json:"_embedded" validate:"omitempty,gt=0,dive,required"`
		Links              *Links            `json:"_links" validate:"required"`
	}
)
