package domain

type (
	PipelineStatusType uint8

	PipelineStatusColor string

	EmbeddedPipelineStatus struct {
		ID         uint64              `json:"id,omitempty" validate:"omitempty"`
		Name       string              `json:"name,omitempty" validate:"omitempty"`
		Sort       uint64              `json:"sort,omitempty" validate:"omitempty"`
		IsEditable bool                `json:"is_editable,omitempty" validate:"omitempty"`
		PipelineID uint64              `json:"pipeline_id,omitempty" validate:"omitempty"`
		Color      PipelineStatusColor `json:"color,omitempty" validate:"omitempty,oneof=#fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 #d6eaff #c1e0ff #98cbff #ebffb1 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea"`
		Type       PipelineStatusType  `json:"type,omitempty" validate:"omitempty,oneof=0 1"`
		AccountID  uint64              `json:"account_id,omitempty" validate:"omitempty"`
	}

	PipelineStatus struct {
		ID         uint64              `json:"id,omitempty" validate:"required"`
		Name       string              `json:"name,omitempty" validate:"required"`
		Sort       uint64              `json:"sort,omitempty" validate:"required"`
		IsEditable bool                `json:"is_editable,omitempty" validate:"omitempty"`
		PipelineID uint64              `json:"pipeline_id,omitempty" validate:"required"`
		Color      PipelineStatusColor `json:"color,omitempty" validate:"omitempty,oneof=#fffeb2 #fffd7f #fff000 #ffeab2 #ffdc7f #ffce5a #ffdbdb #ffc8c8 #ff8f92 #d6eaff #c1e0ff #98cbff #ebffb1 #87f2c0 #f9deff #f3beff #ccc8f9 #eb93ff #f2f3f4 #e6e8ea"`
		Type       PipelineStatusType  `json:"type,omitempty" validate:"omitempty,oneof=0 1"`
		AccountID  uint64              `json:"account_id,omitempty" validate:"required"`
		Links      *Links              `json:"_links,omitempty" validate:"required"`
	}
)

const (
	UnsortedPipelineStatusType PipelineStatusType = iota
	RegularPipelineStatusType
)

const (
	ShalimarPipelineStatusColor          PipelineStatusColor = "#fffeb2"
	WitchHazePipelineStatusColor         PipelineStatusColor = "#fffd7f"
	YellowPipelineStatusColor            PipelineStatusColor = "#fff000"
	BananaManiaPipelineStatusColor       PipelineStatusColor = "#ffeab2"
	SalomiePipelineStatusColor           PipelineStatusColor = "#ffdc7f"
	KournikovaPipelineStatusColor        PipelineStatusColor = "#ffce5a"
	MistyRosePipelineStatusColor         PipelineStatusColor = "#ffdbdb"
	YourPinkPipelineStatusColor          PipelineStatusColor = "#ffc8c8"
	WewakPipelineStatusColor             PipelineStatusColor = "#ff8f92"
	LightPattensBluePipelineStatusColor  PipelineStatusColor = "#d6eaff"
	MediumPattensBluePipelineStatusColor PipelineStatusColor = "#c1e0ff"
	LightSkyBluePipelineStatusColor      PipelineStatusColor = "#98cbff"
	AustralianMintPipelineStatusColor    PipelineStatusColor = "#ebffb1"
	AquamarinePipelineStatusColor        PipelineStatusColor = "#87f2c0"
	SelagoPipelineStatusColor            PipelineStatusColor = "#f9deff"
	MauvePipelineStatusColor             PipelineStatusColor = "#f3beff"
	LavanderBluePipelineStatusColor      PipelineStatusColor = "#ccc8f9"
	VioletPipelineStatusColor            PipelineStatusColor = "#eb93ff"
	AliceBluePipelineStatusColor         PipelineStatusColor = "#f2f3f4"
	SolitudePipelineStatusColor          PipelineStatusColor = "#e6e8ea"
)
