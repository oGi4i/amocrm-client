package domain

type (
	ContactEmbedded struct {
		Tags            []*Tag              `json:"tags,omitempty" validate:"omitempty,dive,required"`             // Теги, привязанные к контакту
		Companies       []*EmbeddedCompany  `json:"companies,omitempty" validate:"omitempty,dive,required"`        // Данные компании, привязанной к контакту
		Leads           []*EmbeddedLead     `json:"leads,omitempty" validate:"omitempty,dive,required"`            // Данные сделок, привязанных к контакту (требует параметра `with=leads` в запросе)
		Customers       []*EmbeddedCustomer `json:"customers,omitempty" validate:"omitempty,dive,required"`        // Данные покупателей, привязанных к контакту (требует параметра `with=customers` в запросе)
		CatalogElements []*CatalogElement   `json:"catalog_elements,omitempty" validate:"omitempty,dive,required"` // Данные элементов списков, привязанных к контакту (требует параметра `with=catalog_elements` в запросе)
	}

	Contact struct {
		ID                 uint64           `json:"id" validate:"required"`                                            // ID контакта
		Name               string           `json:"name" validate:"required"`                                          // Название контакта
		FirstName          string           `json:"first_name" validate:"omitempty"`                                   // Имя контакта
		LastName           string           `json:"last_name" validate:"omitempty"`                                    // Фамилия контакта
		ResponsibleUserID  uint64           `json:"responsible_user_id" validate:"required"`                           // ID пользователя, ответственного за контакт
		GroupID            uint64           `json:"group_id" validate:"required"`                                      // ID группы, в которой состоит ответственный за контакт
		CreatedBy          uint64           `json:"created_by" validate:"required"`                                    // ID пользователя, создавшего контакт
		UpdatedBy          uint64           `json:"updated_by" validate:"required"`                                    // ID пользователя, изменившего контакт
		CreatedAt          uint64           `json:"created_at" validate:"required"`                                    // Дата создания контакта, передаётся в Unix Timestamp
		UpdatedAt          uint64           `json:"updated_at" validate:"required"`                                    // Дата изменения контакта, передаётся в Unix Timestamp
		ClosestTaskAt      uint64           `json:"closest_task_at,omitempty" validate:"omitempty"`                    // Дата ближайшей задачи к выполнению, передаётся в Unix Timestamp
		CustomFieldsValues []*CustomField   `json:"custom_fields_values,omitempty" validate:"omitempty,dive,required"` // Массив дополнительных полей, заданных для контакта
		AccountID          uint64           `json:"account_id" validate:"required"`                                    // ID аккаунта, в котором находится контакт
		Embedded           *ContactEmbedded `json:"_embedded" validate:"omitempty,gt=0,dive,required"`
		Links              *Links           `json:"_links" validate:"required"`
	}
)
