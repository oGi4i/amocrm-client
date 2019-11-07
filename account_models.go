package amocrm

type (
	AccountRequestParams struct {
		With []string `validate:"omitempty,dive,oneof=custom_fields users messenger notifications pipelines groups note_types task_types"`
	}

	AccountResponse struct {
		ID             int    `json:"id,omitempty" validate:"required,omitempty"`
		Name           string `json:"name" validate:"required"`
		Subdomain      string `json:"subdomain" validate:"required"`
		Currency       string `json:"currency" validate:"required"`
		Timezone       string `json:"timezone" validate:"required"`
		TimezoneOffset string `json:"timezone_offset" validate:"required"`
		Language       string `json:"language" validate:"required"`
		DatePattern    struct {
			Date     string `json:"date" validate:"required"`
			Time     string `json:"time" validate:"required"`
			DateTime string `json:"date_time" validate:"required"`
			TimeFull string `json:"time_full" validate:"required"`
		} `json:"date_pattern" validate:"required"`
		CurrentUser int `json:"current_user" validate:"required"`
		Embedded    struct {
			Users        map[string]*User `json:"users" validate:"omitempty,dive,required"`
			CustomFields struct {
				Contacts  map[string]*GetCustomField `json:"contacts" validate:"omitempty,dive,required"`
				Leads     map[string]*GetCustomField `json:"leads,omitempty" validate:"omitempty,dive,required"`
				Companies map[string]*GetCustomField `json:"companies,omitempty" validate:"omitempty,dive,required"`
				Customers []interface{}              `json:"customers,omitempty" validate:"omitempty,dive,required"`
			} `json:"custom_fields" validate:"omitempty"`
			NoteTypes map[string]*NoteType `json:"note_types" validate:"omitempty,dive,required"`
			Groups    map[string]*Group    `json:"groups" validate:"omitempty,dive,required"`
			TaskTypes map[string]*TaskType `json:"task_types" validate:"omitempty,dive,required"`
			Pipelines map[string]*Pipeline `json:"pipelines" validate:"omitempty,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
	}

	AuthAccount struct {
		ID        int    `json:"id" validate:"required"`
		Name      string `json:"name" validate:"required"`
		Subdomain string `json:"subdomain" validate:"required"`
		Language  string `json:"language" validate:"required"`
		Timezone  string `json:"timezone" validate:"required"`
	}
)
