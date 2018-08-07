package amocrm

type (
	Parameters struct {
		UNIQ       string `json:"UNIQ"`
		LINK       string `json:"LINK"`
		PHONE      string `json:"PHONE"`
		DURATION   int    `json:"DURATION"`
		SRC        string `json:"SRC"`
		FROM       string `json:"FROM"`
		CallStatus int    `json:"call_status"`
	}
	Note struct {
		ElementID         int        `json:"element_id"`
		ElementType       int        `json:"element_type"`
		Text              string     `json:"text"`
		NoteType          int        `json:"note_type"`
		CreatedAt         string     `json:"created_at"`
		ResponsibleUserID int        `json:"responsible_user_id"`
		Params            Parameters `json:"params"`
	}
	NoteSetRequest struct {
		Add []Note `json:"add"`
	}
	NoteSetRequestRoot struct {
		Add []NoteSetRequest `json:"request"`
	}
	CustomFiledValue struct {
		value   string
		enum    string
		subtype string
	}
	CustomFiled struct {
		id     int
		values []CustomFiledValue
	}
	Tag struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		ElementType int    `json:"element_type"`
	}
	Contact struct {
		Id                int           `json:"id"`
		Name              string        `json:"name"`
		RequestId         string        `json:"request_id"`
		LastModified      int           `json:"last_modified"`
		AccountId         int           `json:"account_id"`
		DateCreate        int           `json:"date_create"`
		CreatedUserId     int           `json:"created_user_id"`
		ModifiedUserId    int           `json:"modified_user_id"`
		ResponsibleUserId int           `json:"responsible_user_id"`
		GroupId           int           `json:"group_id"`
		ClosestTask       int           `json:"closest_task"`
		LinkedCompanyId   string        `json:"linked_company_id"`
		CompanyName       string        `json:"company_name"`
		Tags              []Tag         `json:"tags"`
		LinkedLeadsId     interface{}   `json:"linked_leads_id"`
		CustomFields      []CustomFiled `json:"custom_fields"`
	}
	ContactListResponse struct {
		Contacts   []Contact `json:"contacts"`
		ServerTime int       `json:"server_time"`
	}
	ContactListResponseRoot struct {
		Response ContactListResponse `json:"response"`
	}
	ContactsSetResponseRoot struct {
		Response ContactsSetResponseContacts `json:"response"`
	}

	ContactsSetResponseContacts struct {
		Contacts   ContactsSetResponse `json:"contacts"`
		ServerTime int                 `json:"server_time"`
	}

	ContactsSetResponse struct {
		Add    []ContactsSetResponseAdd    `json:"add"`
		Update []ContactsSetResponseUpdate `json:"update"`
	}
	ContactsSetResponseAdd struct {
		Id        int `json:"id"`
		RequestId int `json:"request_id"`
	}

	ContactsSetResponseUpdate struct {
		Id        int `json:"id"`
		RequestId int `json:"request_id"`
	}
)
