package amocrm

import "time"

type (
	//RequestParams параметры GET запроса
	LeadRequestParams struct {
		ID                int
		LimitRows         int
		LimitOffset       int
		ResponsibleUserID int
		Query             string
		Status            []int
		Filter            *LeadRequestFilter
	}

	LeadRequestFilter struct {
		Tasks  int
		Active int
	}

	LeadPost struct {
		Name              string   `json:"name"`
		CreatedAt         string   `json:"created_at,omitempty"`
		UpdatedAt         string   `json:"updated_at,omitempty"`
		StatusID          string   `json:"status_id"`
		ResponsibleUserID string   `json:"responsible_user_id,omitempty"`
		Sale              string   `json:"sale,omitempty"`
		Tags              string   `json:"tags,omitempty"`
		ContactsID        []string `json:"contacts_id,omitempty"`
		CompanyID         string   `json:"company_id,omitempty"`
		RequestID         string   `json:"request_id,omitempty"`
	}

	LeadGetRequest struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*LeadResponse `json:"items"`
		} `json:"_embedded"`
	}

	LeadSetRequest struct {
		Add []*LeadPost `json:"add"`
	}

	LeadGetResponse struct {
		Links    *Links `json:"_links"`
		Embedded struct {
			Items []*LeadResponse `json:"items"`
		} `json:"_embedded"`
	}

	LeadResponse struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		ResponsibleUserId int    `json:"responsible_user_id"`
		CreatedBy         int    `json:"created_by"`
		CreatedAt         int    `json:"created_at"`
		UpdatedAt         int    `json:"updated_at"`
		AccountId         int    `json:"account_id"`
		IsDeleted         bool   `json:"is_deleted"`
		MainContact       struct {
			ID    int    `json:"id"`
			Links *Links `json:"_links"`
		} `json:"main_contact"`
		GroupId       int            `json:"group_id"`
		ClosedAt      int            `json:"closed_at"`
		ClosestTaskAt int            `json:"closest_task_at"`
		Tags          *Tag           `json:"tags"`
		CustomFields  []*CustomField `json:"custom_fields"`
		Contact       struct {
			ID    []int  `json:"id"`
			Links *Links `json:"_links"`
		} `json:"contacts"`
		StatusId int `json:"status_id"`
		Sale     int `json:"sale"`
		Pipeline struct {
			ID    int    `json:"id"`
			Links *Links `json:"_links"`
		} `json:"pipeline"`
		Links *Links `json:"_links"`
	}

	LeadGet struct {
		ID                int       `json:"id"`
		Name              string    `json:"name"`
		ResponsibleUserId int       `json:"responsible_user_id"`
		CreatedBy         int       `json:"created_by"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		AccountId         int       `json:"account_id"`
		IsDeleted         bool      `json:"is_deleted"`
		MainContact       struct {
			ID    int    `json:"id"`
			Links *Links `json:"_links"`
		} `json:"main_contact"`
		GroupId       int       `json:"group_id"`
		ClosedAt      time.Time `json:"closed_at"`
		ClosestTaskAt time.Time `json:"closest_task_at"`
		Contact       struct {
			ID    []int  `json:"id"`
			Links *Links `json:"_links"`
		} `json:"contacts"`
		StatusId int `json:"status_id"`
		Sale     int `json:"sale"`
		Pipeline struct {
			ID    int    `json:"id"`
			Links *Links `json:"_links"`
		} `json:"pipeline"`
		Links *Links `json:"_links"`
	}
)
