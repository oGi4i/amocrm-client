package domain

type (
	EntityNames struct {
		Leads map[string]*LanguageEntityNames `json:"leads" validate:"required,dive,keys,required,endkeys,required"`
	}

	EntityForm struct {
		Dative        string `json:"dative" validate:"omitempty"`
		Default       string `json:"default" validate:"required"`
		Genitive      string `json:"genitive" validate:"omitempty"`
		Accusative    string `json:"accusative" validate:"omitempty"`
		Instrumental  string `json:"instrumental" validate:"omitempty"`
		Prepositional string `json:"prepositional" validate:"omitempty"`
	}

	LanguageEntityNames struct {
		Gender       string      `json:"gender" validate:"required"`
		PluralForm   *EntityForm `json:"plural_form" validate:"required"`
		SingularForm *EntityForm `json:"singular_form" validate:"required"`
	}
)
