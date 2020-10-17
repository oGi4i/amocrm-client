package domain

type (
	CustomersMode string

	ContactNameDisplayOrder uint8

	AccountEmbedded struct {
		AmojoRights      *AmojoRights      `json:"amojo_rights" validate:"omitempty"`                    // Объект настроек чатов (требует параметра `with=amojo_rights` в запросе)
		UsersGroups      []*UserGroup      `json:"users_groups" validate:"omitempty,gt=0,dive,required"` // Массив объектов группы пользователей (требует параметра `with=user_groups` в запросе)
		TaskTypes        []*TaskType       `json:"task_types" validate:"omitempty,gt=0,dive,required"`   // Типы задач, доступные в аккаунте (требует параметра `with=task_types` в запросе)
		DatetimeSettings *DatetimeSettings `json:"datetime_settings" validate:"omitempty"`               // Настройки и форматы даты и времени в аккаунте (требует параметра `with=datetime_settings` в запросе)
		EntityNames      *EntityNames      `json:"entity_names" validate:"omitempty"`                    // Настройки названия сущностей (требует параметра `with=entity_names` в запросе)
	}

	Account struct {
		ID                      uint64                  `json:"id,omitempty" validate:"required"`               // ID аккаунта
		Name                    string                  `json:"name" validate:"required"`                       // Название аккаунта
		Subdomain               string                  `json:"subdomain" validate:"required"`                  // Субдомен аккаунта
		CreatedAt               uint64                  `json:"created_at" validate:"required"`                 // Дата создания аккаунта в Unix Timestamp
		CreatedBy               uint64                  `json:"created_by" validate:"required"`                 // ID пользователя, который создал аккаунт
		UpdatedAt               uint64                  `json:"updated_at" validate:"required"`                 // Дата последнего изменения свойств аккаунта в Unix Timestamp
		UpdatedBy               uint64                  `json:"updated_by" validate:"required"`                 // ID пользователя, который последним изменял свойства аккаунта
		CurrentUserID           uint64                  `json:"current_user_id" validate:"required"`            // ID текущего пользователя
		Country                 string                  `json:"country" validate:"required"`                    // Страна, указанная в настройках аккаунта
		CustomersMode           CustomersMode           `json:"customers_mode" validate:"required"`             // Режим покупателей
		IsUnsortedOn            bool                    `json:"is_unsorted_on" validate:"omitempty"`            // Включен ли функционал "Неразобранного" в аккаунте
		IsLossReasonEnabled     bool                    `json:"is_loss_reason_enabled" validate:"omitempty"`    // Включен ли функционал причин отказа
		IsHelpBotEnabled        bool                    `json:"is_helpbot_enabled" validate:"omitempty"`        // Включен ли функционал Типовых вопросов
		IsTechnicalAccount      bool                    `json:"is_technical_account" validate:"omitempty"`      // Признак технического аккаунта
		ContactNameDisplayOrder ContactNameDisplayOrder `json:"contact_name_display_order" validate:"required"` // Порядок отображения имён контактов
		AmojoID                 string                  `json:"amojo_id" validate:"omitempty"`                  // Уникальный идентификато аккаунта для работы с сервисом чатов amoJo (требует параметра `with=amojo_id` в запросе)
		Version                 uint64                  `json:"version" validate:"omitempty"`                   // Текущая версия amoCRM (требует параметра `with=version` в запросе)
		Embedded                *AccountEmbedded        `json:"_embedded" validate:"required"`
		Links                   *Links                  `json:"_links" validate:"required"`
	}
)

const (
	UnavailableCustomersMode CustomersMode = "unavailable"
	DisabledCustomersMode    CustomersMode = "disabled"
	SegmentsCustomersMode    CustomersMode = "segments"
	DynamicCustomersMode     CustomersMode = "dynamic"
	PeriodicityCustomersMode CustomersMode = "periodicity"
)

const (
	NameSurnameContactNameDisplayOrder ContactNameDisplayOrder = iota + 1 // Имя, Фамилия
	SurnameNameContactNameDisplayOrder                                    // Фамилия, Имя
)
