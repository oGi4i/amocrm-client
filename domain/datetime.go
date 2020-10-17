package domain

type (
	DatetimeSettings struct {
		DatePattern      string `json:"date_pattern" validate:"required"`
		ShortDatePattern string `json:"short_date_pattern" validate:"required"`
		ShortTimePattern string `json:"short_time_pattern" validate:"required"`
		DateFormat       string `json:"date_format" validate:"required"`
		TimeFormat       string `json:"time_format" validate:"required"`
		Timezone         string `json:"timezone" validate:"required"`
		TimezoneOffset   string `json:"timezone_offset" validate:"required"`
	}
)
