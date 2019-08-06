package amocrm

type (
	//RequestParams параметры GET запроса
	AccountRequestParams struct {
		With string
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
			Users        map[string]AccountUser `json:"users"`
			CustomFields struct {
				Contacts  map[string]AccountContact   `json:"contacts"`
				Leads     map[string]AccountLead      `json:"leads,omitempty"`
				Companies map[string]AccountCompanies `json:"companies,omitempty"`
				Customers []interface{}               `json:"customers,omitempty"`
			} `json:"custom_fields"`
			NoteTypes map[string]AccountNoteType `json:"note_types"`
			Groups    map[string]AccountGroup    `json:"groups"`
			TaskTypes map[string]AccountTaskType `json:"task_types"`
			Pipelines map[string]AccountPipeline `json:"pipelines"`
		} `json:"_embedded"`
	}

	AccountUser struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Login    string `json:"login"`
		Language string `json:"language"`
		GroupID  int    `json:"group_id"`
		IsActive bool   `json:"is_active"`
		IsFree   bool   `json:"is_free"`
		IsAdmin  bool   `json:"is_admin"`
		Rights   struct {
			Mail          string `json:"mail"`
			IncomingLeads string `json:"incoming_leads"`
			Catalogs      string `json:"catalogs"`
			LeadAdd       string `json:"lead_add"`
			LeadView      string `json:"lead_view"`
			LeadEdit      string `json:"lead_edit"`
			LeadDelete    string `json:"lead_delete"`
			LeadExport    string `json:"lead_export"`
			ContactAdd    string `json:"contact_add"`
			ContactView   string `json:"contact_view"`
			ContactEdit   string `json:"contact_edit"`
			ContactDelete string `json:"contact_delete"`
			ContactExport string `json:"contact_export"`
			CompanyAdd    string `json:"company_add"`
			CompanyView   string `json:"company_view"`
			CompanyEdit   string `json:"company_edit"`
			CompanyDelete string `json:"company_delete"`
			CompanyExport string `json:"company_export"`
			TaskEdit      string `json:"task_edit"`
			TaskDelete    string `json:"task_delete"`
		} `json:"rights"`
	}

	AccountContact struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FieldType   int    `json:"field_type"`
		Sort        int    `json:"sort"`
		Code        string `json:"code"`
		IsMultiple  bool   `json:"is_multiple"`
		IsSystem    bool   `json:"is_system"`
		IsEditable  bool   `json:"is_editable"`
		IsRequired  bool   `json:"is_required"`
		IsDeletable bool   `json:"is_deletable"`
		IsVisible   bool   `json:"is_visible"`
		Params      struct {
		} `json:"params"`
		Enums map[string]string `json:"enums"`
	}

	AccountLead struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FieldType   int    `json:"field_type"`
		Sort        int    `json:"sort"`
		IsMultiple  bool   `json:"is_multiple"`
		IsSystem    bool   `json:"is_system"`
		IsEditable  bool   `json:"is_editable"`
		IsRequired  bool   `json:"is_required"`
		IsDeletable bool   `json:"is_deletable"`
		IsVisible   bool   `json:"is_visible"`
		Params      struct {
		} `json:"params"`
		Enums map[string]string `json:"enums"`
	}

	AccountCompanies struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FieldType   int    `json:"field_type"`
		Sort        int    `json:"sort"`
		IsMultiple  bool   `json:"is_multiple"`
		IsSystem    bool   `json:"is_system"`
		IsEditable  bool   `json:"is_editable"`
		IsRequired  bool   `json:"is_required"`
		IsDeletable bool   `json:"is_deletable"`
		IsVisible   bool   `json:"is_visible"`
		Params      struct {
		} `json:"params"`
		Enums map[string]string `json:"enums"`
	}

	AccountNoteType struct {
		ID         int    `json:"id"`
		Code       string `json:"code"`
		IsEditable bool   `json:"is_editable"`
	}

	AccountGroup struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	AccountTaskType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	AccountStatus struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Color      string `json:"color"`
		Sort       int    `json:"sort"`
		IsEditable bool   `json:"is_editable"`
	}

	AccountPipeline struct {
		ID       int                      `json:"id"`
		Name     string                   `json:"name"`
		Sort     int                      `json:"sort"`
		IsMain   bool                     `json:"is_main"`
		Statuses map[string]AccountStatus `json:"statuses"`
		Links    struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
	}
)
