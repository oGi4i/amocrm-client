package domain

type (
	EmbeddedCompany struct {
		ID    uint64 `json:"id" validate:"required"`
		Links *Links `json:"_links" validate:"required"`
	}

	CompanyEmbedded struct {
		Tags            []*Tag              `json:"tags,omitempty" validate:"omitempty,dive,required"`             // Теги, привязанные к компании
		Contacts        []*Contact          `json:"contacts,omitempty" validate:"omitempty,dive,required"`         // Данные контактов, привязанных к компании (требует параметра `with=contacts` в запросе)
		Customers       []*EmbeddedCustomer `json:"customers,omitempty" validate:"omitempty,dive,required"`        // Данные покупателей, привязанных к компании (требует параметра `with=customers` в запросе)
		Leads           []*EmbeddedLead     `json:"leads,omitempty" validate:"omitempty,dive,required"`            // Данные сделок, привязанных к компании (требует параметра `with=leads` в запросе)
		CatalogElements []*CatalogElement   `json:"catalog_elements,omitempty" validate:"omitempty,dive,required"` // Данные элементов списков, привязанных к компании (требует параметра `with=catalog_elements` в запросе)
	}

	Company struct {
		ID                 uint64           `json:"id" validate:"required"`                                                 // ID компании
		Name               string           `json:"name" validate:"required"`                                               // Название компании
		ResponsibleUserID  uint64           `json:"responsible_user_id" validate:"required"`                                // ID пользователя, ответственного за компанию
		GroupID            uint64           `json:"group_id" validate:"required"`                                           // ID группы, в которой состоит ответственный
		CreatedBy          uint64           `json:"created_by" validate:"required"`                                         // ID пользователя, создавшего компанию
		UpdatedBy          uint64           `json:"updated_by" validate:"required"`                                         // ID пользователя, изменившего компанию
		CreatedAt          uint64           `json:"created_at" validate:"required"`                                         // Дата создания компании, передаётся в Unix Timestamp
		UpdatedAt          uint64           `json:"updated_at" validate:"required"`                                         // Дата изменения компании, передаётся в Unix Timestamp
		ClosestTaskAt      uint64           `json:"closest_task_at" validate:"omitempty"`                                   // Дата ближайшей задачи к выполнению, передаётся в Unix Timestamp
		CustomFieldsValues []*CustomField   `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"` // // Массив дополнительных полей, заданных для компании
		AccountID          uint64           `json:"account_id" validate:"required"`                                         // ID аккаунта, в котором находится компания
		Links              *Links           `json:"_links" validate:"required"`
		Embedded           *CompanyEmbedded `json:"_embedded" validate:"omitempty"`
	}
)
