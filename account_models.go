package amocrm

type (
	//RequestParams параметры GET запроса
	AccountRequestParams struct {
		With []string
	}

	AccountResponse struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Subdomain      string `json:"subdomain"`
		Currency       string `json:"currency"`
		Timezone       string `json:"timezone"`
		TimezoneOffset string `json:"timezone_offset"`
		Language       string `json:"language"`
		DatePattern    struct {
			Date     string `json:"date"`
			Time     string `json:"time"`
			DateTime string `json:"date_time"`
			TimeFull string `json:"time_full"`
		} `json:"date_pattern"`
		CurrentUser int `json:"current_user"`
		Embedded    struct {
			Users        map[string]*User `json:"users"`
			CustomFields struct {
				Contacts  map[string]*ContactCustomField `json:"contacts"`
				Leads     map[string]*LeadCustomField    `json:"leads,omitempty"`
				Companies map[string]*CompanyCustomField `json:"companies,omitempty"`
				Customers []interface{}                  `json:"customers,omitempty"`
			} `json:"custom_fields"`
			NoteTypes map[string]*NoteType `json:"note_types"`
			Groups    map[string]*Group    `json:"groups"`
			TaskTypes map[string]*TaskType `json:"task_types"`
			Pipelines map[string]*Pipeline `json:"pipelines"`
		} `json:"_embedded"`
	}

	AuthAccount struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Subdomain string `json:"subdomain"`
		Language  string `json:"language"`
		Timezone  string `json:"timezone"`
	}

	AccountStatus struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Color      string `json:"color"`
		Sort       int    `json:"sort"`
		IsEditable bool   `json:"is_editable"`
	}
)
