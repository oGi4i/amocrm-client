package domain

type (
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
