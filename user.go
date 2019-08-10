package amocrm

type (
	AuthUser struct {
		ID       int    `json:"id"`
		Language string `json:"language"`
	}

	User struct {
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
)
