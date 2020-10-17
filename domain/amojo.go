package domain

type (
	AmojoRights struct {
		CanDirect       bool `json:"can_direct" validate:"omitempty"`
		CanCreateGroups bool `json:"can_create_groups" validate:"omitempty"`
	}
)
