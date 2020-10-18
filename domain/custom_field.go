package domain

type (
	CustomField struct {
		ID     uint64              `json:"field_id" validate:"required"`                                                                                                                                                       // ID поля
		Name   string              `json:"field_name" validate:"required"`                                                                                                                                                     // Название поля
		Code   string              `json:"field_code" validate:"required"`                                                                                                                                                     // Код поля, по которому можно обновлять значение в сущности, без передачи ID поля 		// Сортровка поля
		Type   CustomFieldType     `json:"field_type" validate:"oneof=text numeric checkbox select multiselect date url textarea radiobutton streetaddress smart_address birthday legal_entity datetime price category items"` // Тип поля
		Values []*CustomFieldValue `json:"values" validate:"required,dive,required"`
	}

	CustomFieldValue struct {
		Value  string `json:"value" validate:"required"`
		EnumID uint64 `json:"enum_id,omitempty" validate:"omitempty"`
		Enum   string `json:"enum,omitempty" validate:"omitempty"`
	}

	UpdateCustomField struct {
		ID     uint64        `json:"field_id,omitempty" validate:"required_without=Code"`
		Code   string        `json:"field_code,omitempty" validate:"required_without=ID"`
		Name   string        `json:"field_name,omitempty" validate:"omitempty"`
		Values []interface{} `json:"values" validate:"required,dive,required"`
	}

	CustomFieldType string

	CustomFieldRemind string

	CustomFieldEnum struct {
		ID    uint64 `json:"id" validate:"required"`    // ID значения
		Value string `json:"value" validate:"required"` // Значение
		Sort  uint64 `json:"sort" validate:"required"`  // Сортировка значения
	}

	CustomFieldNested struct {
		ID       uint64 `json:"id" validate:"required"`        // ID вложенного значения
		ParentID uint64 `json:"parent_id" validate:"required"` // ID родительского значения
		Value    string `json:"value" validate:"required"`     // Значение
		Sort     uint64 `json:"sort" validate:"required"`      // Сортировка значения
	}

	CustomFieldRequiredStatus struct {
		StatusID   uint64 `json:"status_id" validate:"required"`   // ID статуса, для перехода в который данное поле обязательно к заполнению
		PipelineID uint64 `json:"pipeline_id" validate:"required"` // ID воронки, для перехода в который данное поле обязательно к заполнению
	}

	CustomFieldInfo struct {
		ID               uint64                       `json:"id" validate:"required"`                                                                                                                                                       // ID поля
		Name             string                       `json:"name" validate:"required"`                                                                                                                                                     // Название поля
		Code             string                       `json:"code" validate:"required"`                                                                                                                                                     // Код поля, по которому можно обновлять значение в сущности, без передачи ID поля
		Sort             uint64                       `json:"sort" validate:"required"`                                                                                                                                                     // Сортровка поля
		Type             CustomFieldType              `json:"type" validate:"oneof=text numeric checkbox select multiselect date url textarea radiobutton streetaddress smart_address birthday legal_entity datetime price category items"` // Тип поля
		EntityType       EntityType                   `json:"entity_type" validate:"oneof=leads contacts companies segments customers catalogs"`                                                                                            // Тип сущности
		IsPredefined     bool                         `json:"is_predefined" validate:"omitempty"`                                                                                                                                           // Признак является ли поле предустановленным. Данный признак возвращается только при получении списка полей контактов и компаний
		IsDeletable      bool                         `json:"is_deletable" validate:"omitempty"`                                                                                                                                            // Признак доступно ли поле для удаления. Данный признак возвращается только при получении списка полей контактов, компаний и каталогов
		IsVisible        bool                         `json:"is_visible" validate:"omitempty"`                                                                                                                                              // Признак отображается ли поле в интерфейсе списка. Данный признак возвращается только при получении списка полей каталогов
		IsRequired       bool                         `json:"is_required" validate:"omitempty"`                                                                                                                                             // Признак обязательно ли поле для заполнения при создании элемента списка. Данный признак возвращается только при получении списка полей каталогов
		Settings         []string                     `json:"settings" validate:"omitempty"`                                                                                                                                                // Настройки поля. Данный признак возвращается только при получении списка полей каталогов
		Remind           CustomFieldRemind            `json:"remind" validate:"omitempty,oneof=never day week month"`                                                                                                                       // Когда напоминать о дне рождения. Значение доступно только для поля типа birthday. Данный признак возвращается только при получении списка полей контактов, сделок и компаний
		Enums            []*CustomFieldEnum           `json:"enums" validate:"omitempty,gt=0,dive,required"`                                                                                                                                // Доступные значения для поля. Значение доступно только для полей с поддержкой enum
		Nested           []*CustomFieldNested         `json:"nested" validate:"omitempty,gt=0,dive,required"`                                                                                                                               // Вложенные значения. Значение доступно только для поля category. Признак возвращается только при получении списка полей каталогов
		IsAPIOnly        bool                         `json:"is_api_only" validate:"omitempty"`                                                                                                                                             // Признак доступно ли поле для редактирования только через API. Признак возвращается только при получении списка полей контактов, сделок и компаний
		GroupID          uint64                       `json:"group_id" validate:"omitempty"`                                                                                                                                                // ID группы полей, в которой состоит поле. Признак возвращается только при получении списка полей контактов, сделок, покупателей и компаний
		RequiredStatuses []*CustomFieldRequiredStatus `json:"required_statuses" validate:"omitempty"`                                                                                                                                       // Обязательные поля для смены этапов. Признак возвращается только при получении списка полей контактов, сделок и компаний
	}
)

const (
	TextCustomFieldType          CustomFieldType = "text"          // Текст
	NumericCustomFieldType       CustomFieldType = "numeric"       // Число
	CheckboxCustomFieldType      CustomFieldType = "checkbox"      // Флаг
	SelectCustomFieldType        CustomFieldType = "select"        // Список
	MultiSelectCustomFieldType   CustomFieldType = "multiselect"   // Мультисписок
	DateCustomFieldType          CustomFieldType = "date"          // Дата
	URLCustomFieldType           CustomFieldType = "url"           // Ссылка
	TextAreaCustomFieldType      CustomFieldType = "textarea"      // Текстовая область
	RadioButtonCustomFieldType   CustomFieldType = "radiobutton"   // Переключатель
	StreetAddressCustomFieldType CustomFieldType = "streetaddress" // Короткий адрес
	SmartAddressCustomFieldType  CustomFieldType = "smart_address" // Адрес
	BirthDayCustomFieldType      CustomFieldType = "birthday"      // День рождения
	LegalEntityCustomFieldType   CustomFieldType = "legal_entity"  // Юр. лицо
	DatetimeCustomFieldType      CustomFieldType = "datetime"      // Дата и время
	PriceCustomFieldType         CustomFieldType = "price"         // Цена
	CategoryCustomFieldType      CustomFieldType = "category"      // Категория
	ItemsCustomFieldType         CustomFieldType = "items"         // Предметы
)

const (
	NeverCustomFieldRemind CustomFieldRemind = "never" // Никогда не напоминать о дне рождения
	DayCustomFieldRemind   CustomFieldRemind = "day"   // Напоминать за день до дня рождения
	WeekCustomFieldRemind  CustomFieldRemind = "week"  // Напоминать за неделю до дня рождения
	MonthCustomFieldRemind CustomFieldRemind = "month" // Напоминать за месяц до дня рождения
)
