package amocrm

var (
	apiUrls = map[string]string{
		"auth":            "/private/api/auth.php?type=json",
		"notes":           "/api/v2/notes",
		"contacts":        "/api/v2/contacts",
		"incomingleadsip": "/api/v2/incoming_leads/sip",
		"account":         "/api/v2/account",
	}
)
