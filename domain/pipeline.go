package domain

type (
	PipelineEmbedded struct {
		Statuses []*PipelineStatus `json:"statuses" validate:"required,gt=0,dive,required"` // Данные статусов, имеющихся в воронке
	}

	Pipeline struct {
		ID           uint64            `json:"id,omitempty" validate:"required"`              // ID воронки
		Name         string            `json:"name,omitempty" validate:"required"`            // Название воронки
		Sort         uint64            `json:"sort,omitempty" validate:"required"`            // Сортировка воронки
		IsMain       bool              `json:"is_main,omitempty" validate:"omitempty"`        // Является ли воронка главной
		IsUnsortedOn bool              `json:"is_unsorted_on,omitempty" validate:"omitempty"` // Включено ли неразобранное в воронке
		IsArchive    bool              `json:"is_archive,omitempty" validate:"omitempty"`     // Является ли воронка архивной
		AccountID    uint64            `json:"account_id,omitempty" validate:"required"`      // ID аккаунта, в котором находится воронка
		Embedded     *PipelineEmbedded `json:"_embedded,omitempty" validate:"required"`
		Links        *Links            `json:"_links,omitempty" validate:"required"`
	}
)
