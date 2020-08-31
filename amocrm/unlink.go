package amocrm

type Unlink struct {
	LeadsID    []int `json:"leads_id,omitempty" validate:"omitempty,gt=0,dive,required"`
	ContactsID []int `json:"contacts_id,omitempty" validate:"omitempty,gt=0,dive,required"`
	CompanyID  int   `json:"company_id,omitempty" validate:"omitempty"`
}
