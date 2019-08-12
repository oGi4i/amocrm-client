package amocrm

type (
	AuthUser struct {
		ID       int    `json:"id" validate:"required"`
		Language string `json:"language" validate:"required"`
	}

	User struct {
		ID       int    `json:"id" validate:"required"`
		Name     string `json:"name" validate:"required"`
		LastName string `json:"last_name,omitempty" validate:"omitempty"`
		Login    string `json:"login" validate:"required"`
		Language string `json:"language" validate:"required"`
		GroupID  int    `json:"group_id" validate:"omitempty"`
		IsActive bool   `json:"is_active" validate:"omitempty"`
		IsFree   bool   `json:"is_free" validate:"omitempty"`
		IsAdmin  bool   `json:"is_admin" validate:"omitempty"`
		Rights   struct {
			Mail          string `json:"mail" validate:"required"`
			IncomingLeads string `json:"incoming_leads" validate:"required"`
			Catalogs      string `json:"catalogs" validate:"required"`
			LeadAdd       string `json:"lead_add" validate:"required"`
			LeadView      string `json:"lead_view" validate:"required"`
			LeadEdit      string `json:"lead_edit" validate:"required"`
			LeadDelete    string `json:"lead_delete" validate:"required"`
			LeadExport    string `json:"lead_export" validate:"required"`
			ContactAdd    string `json:"contact_add" validate:"required"`
			ContactView   string `json:"contact_view" validate:"required"`
			ContactEdit   string `json:"contact_edit" validate:"required"`
			ContactDelete string `json:"contact_delete" validate:"required"`
			ContactExport string `json:"contact_export" validate:"required"`
			CompanyAdd    string `json:"company_add" validate:"required"`
			CompanyView   string `json:"company_view" validate:"required"`
			CompanyEdit   string `json:"company_edit" validate:"required"`
			CompanyDelete string `json:"company_delete" validate:"required"`
			CompanyExport string `json:"company_export" validate:"required"`
			TaskEdit      string `json:"task_edit" validate:"required"`
			TaskDelete    string `json:"task_delete" validate:"required"`
		} `json:"rights" validate:"required"`
	}
)
