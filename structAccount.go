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
				Leads     map[string]AccountLead      `json:"leads"`
				Companies map[string]AccountCompanies `json:"companies"`
				Customers []interface{}               `json:"customers"`
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
		}
		Enums struct {
			Num361629 string `json:"361629"`
			Num361631 string `json:"361631"`
			Num361633 string `json:"361633"`
			Num361635 string `json:"361635"`
			Num361637 string `json:"361637"`
			Num361639 string `json:"361639"`
		} `json:"enums"`
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
		Enums struct {
			Num1054811 string `json:"1054811"`
			Num1054813 string `json:"1054813"`
			Num1054815 string `json:"1054815"`
			Num1054817 string `json:"1054817"`
			Num1054819 string `json:"1054819"`
			Num1054821 string `json:"1054821"`
			Num1054969 string `json:"1054969"`
			Num1054971 string `json:"1054971"`
			Num1060097 string `json:"1060097"`
		} `json:"enums"`
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
		Enums struct {
			Num361629 string `json:"361629"`
			Num361631 string `json:"361631"`
			Num361633 string `json:"361633"`
			Num361635 string `json:"361635"`
			Num361637 string `json:"361637"`
			Num361639 string `json:"361639"`
		} `json:"enums"`
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
	AccountPipeline struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Sort     int    `json:"sort"`
		IsMain   bool   `json:"is_main"`
		Statuses struct {
			Num142 struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Sort       int    `json:"sort"`
				IsEditable bool   `json:"is_editable"`
			} `json:"142"`
			Num143 struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Sort       int    `json:"sort"`
				IsEditable bool   `json:"is_editable"`
			} `json:"143"`
			Num19743178 struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Sort       int    `json:"sort"`
				IsEditable bool   `json:"is_editable"`
			} `json:"19743178"`
			Num19743181 struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Sort       int    `json:"sort"`
				IsEditable bool   `json:"is_editable"`
			} `json:"19743181"`
			Num19743184 struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Sort       int    `json:"sort"`
				IsEditable bool   `json:"is_editable"`
			} `json:"19743184"`
			Num19743187 struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Sort       int    `json:"sort"`
				IsEditable bool   `json:"is_editable"`
			} `json:"19743187"`
		} `json:"statuses"`
		Links struct {
			Self struct {
				Href   string `json:"href"`
				Method string `json:"method"`
			} `json:"self"`
		} `json:"_links"`
	}
)
