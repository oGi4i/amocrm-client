package amocrm

type (
	//RequestParams параметры GET запроса
	ContactRequestParams struct {
		ID                int    `validate:"omitempty"`
		LimitRows         int    `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset       int    `validate:"omitempty"`
		ResponsibleUserID int    `validate:"omitempty"`
		Query             string `validate:"omitempty"`
	}

	ContactAdd struct {
		Name              string             `json:"name" validate:"required"`
		CreatedAt         int                `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                `json:"updated_at,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int                `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		CreatedBy         int                `json:"created_by,string,omitempty" validate:"omitempty"`
		CompanyName       string             `json:"company_name,omitempty" validate:"omitempty"`
		Tags              string             `json:"tags,omitempty" validate:"omitempty"`
		LeadsID           []string           `json:"leads_id,omitempty" validate:"omitempty,gt=0,dive,required"`
		CustomersID       int                `json:"customers_id,string,omitempty" validate:"omitempty"`
		CompanyID         int                `json:"company_id,string,omitempty" validate:"omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty" validate:"omitempty,gt=0,required"`
	}

	ContactUpdate struct {
		ID                int                `json:"id,string" validate:"required"`
		Name              string             `json:"name,omitempty" validate:"omitempty"`
		CreatedAt         int                `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                `json:"updated_at,string" validate:"required"`
		ResponsibleUserID int                `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		CreatedBy         int                `json:"created_by,string,omitempty" validate:"omitempty"`
		CompanyName       string             `json:"company_name,omitempty" validate:"omitempty"`
		Tags              string             `json:"tags,omitempty" validate:"omitempty"`
		LeadsID           []string           `json:"leads_id,omitempty" validate:"omitempty,gt=0,required"`
		CustomersID       int                `json:"customers_id,string,omitempty" validate:"omitempty"`
		CompanyID         int                `json:"company_id,string,omitempty" validate:"omitempty"`
		CustomFields      []*CustomFieldPost `json:"custom_fields,omitempty" validate:"omitempty,gt=0,required"`
		Unlink            *Unlink            `json:"unlink,omitempty" validate:"omitempty"`
	}

	AddContactRequest struct {
		Add []*ContactAdd `json:"add" validate:"required,dive,required"`
	}

	UpdateContactRequest struct {
		Update []*ContactUpdate `json:"update" validate:"required,dive,required"`
	}

	GetContactResponse struct {
		Links    *Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items []*Contact `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response,omitempty" validate:"omitempty"`
	}

	Contact struct {
		ID                int    `json:"id" validate:"required"`
		Name              string `json:"name" validate:"required"`
		ResponsibleUserID int    `json:"responsible_user_id" validate:"required"`
		CreatedBy         int    `json:"created_by" validate:"required"`
		CreatedAt         int    `json:"created_at" validate:"required"`
		UpdatedAt         int    `json:"updated_at" validate:"required"`
		AccountID         int    `json:"account_id" validate:"required"`
		UpdatedBy         int    `json:"updated_by" validate:"required"`
		GroupID           int    `json:"group_id,omitempty" validate:"omitempty"`
		Company           struct {
			ID    int    `json:"id" validate:"omitempty"`
			Name  string `json:"name" validate:"omitempty"`
			Links *Links `json:"_links" validate:"omitempty"`
		} `json:"company,omitempty" validate:"omitempty"`
		Leads struct {
			ID    []int  `json:"id" validate:"omitempty,dive,required"`
			Links *Links `json:"_links" validate:"omitempty"`
		} `json:"leads,omitempty" validate:"omitempty"`
		ClosestTaskAt int            `json:"closest_task_at,omitempty" validate:"omitempty"`
		Tags          []*Tag         `json:"tags,omitempty" validate:"omitempty,dive,required"`
		CustomFields  []*CustomField `json:"custom_fields,omitempty" validate:"omitempty,dive,required"`
		Customers     struct {
		} `json:"customers,omitempty" validate:"omitempty"`
		Links *Links `json:"_links" validate:"required"`
	}
)
