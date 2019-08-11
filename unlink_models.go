package amocrm

type (
	Unlink struct {
		LeadsID    []int `json:"leads_id,omitempty"`
		ContactsID []int `json:"contacts_id,omitempty"`
		CompanyID  int   `json:"company_id,omitempty"`
	}
)
