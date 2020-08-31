package amocrm

import (
	"context"
	"encoding/json"
	"strconv"
)

type (
	LeadRequestParams struct {
		ID                []int              `validate:"omitempty,gt=0,dive,required"`
		LimitRows         int                `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset       int                `validate:"omitempty"`
		ResponsibleUserID int                `validate:"omitempty"`
		Query             string             `validate:"omitempty"`
		Status            []int              `validate:"omitempty,gt=0,dive,required"`
		Filter            *LeadRequestFilter `validate:"omitempty"`
	}

	LeadRequestTasksFilter int

	LeadRequestActiveFilter int

	LeadRequestFilter struct {
		Tasks  LeadRequestTasksFilter  `validate:"omitempty,oneof=1 2"`
		Active LeadRequestActiveFilter `validate:"omitempty,eq=1"`
	}

	LeadAdd struct {
		Name              string               `json:"name" validate:"required"`
		CreatedAt         int                  `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                  `json:"updated_at,string,omitempty" validate:"omitempty"`
		StatusID          int                  `json:"status_id,string" validate:"required"`
		PipelineID        int                  `json:"pipeline_id,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int                  `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		Sale              int                  `json:"sale,string,omitempty" validate:"omitempty"`
		Tags              string               `json:"tags,omitempty" validate:"omitempty"`
		CustomFields      []*UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,dive,required"`
		ContactsID        []string             `json:"contacts_id,omitempty" validate:"omitempty,gt=0,dive,required"`
		CompanyID         int                  `json:"company_id,string,omitempty" validate:"omitempty"`
		RequestID         int                  `json:"request_id,string,omitempty" validate:"omitempty"`
	}

	LeadUpdate struct {
		ID                int                  `json:"id,string" validate:"required"`
		Name              string               `json:"name,omitempty" validate:"omitempty"`
		CreatedAt         int                  `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                  `json:"updated_at,string" validate:"required"`
		StatusID          int                  `json:"status_id,string,omitempty" validate:"omitempty"`
		PipelineID        int                  `json:"pipeline_id,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int                  `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		Sale              int                  `json:"sale,string,omitempty" validate:"omitempty"`
		Tags              string               `json:"tags,omitempty" validate:"omitempty"`
		CustomFields      []*UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,dive,required"`
		ContactsID        []string             `json:"contacts_id,omitempty" validate:"omitempty,gt=0,dive,required"`
		CompanyID         int                  `json:"company_id,string,omitempty" validate:"omitempty"`
		RequestID         int                  `json:"request_id,string,omitempty" validate:"omitempty"`
		Unlink            *Unlink              `json:"unlink,omitempty" validate:"omitempty"`
	}

	AddLeadRequest struct {
		Add []*LeadAdd `json:"add" validate:"required,dive,required"`
	}

	UpdateLeadRequest struct {
		Update []*LeadUpdate `json:"update" validate:"required,dive,required"`
	}

	GetLeadResponse struct {
		Links    *Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items []*Lead `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response" validate:"omitempty"`
	}

	Lead struct {
		ID                int    `json:"id" validate:"required"`
		Name              string `json:"name" validate:"required"`
		ResponsibleUserID int    `json:"responsible_user_id" validate:"required"`
		CreatedBy         int    `json:"created_by" validate:"required"`
		CreatedAt         int    `json:"created_at" validate:"required"`
		UpdatedAt         int    `json:"updated_at" validate:"required"`
		AccountID         int    `json:"account_id" validate:"required"`
		IsDeleted         bool   `json:"is_deleted" validate:"omitempty"`
		MainContact       struct {
			ID    int    `json:"id" validate:"omitempty"`
			Links *Links `json:"_links" validate:"omitempty"`
		} `json:"main_contact,omitempty" validate:"omitempty"`
		GroupID       int            `json:"group_id,omitempty" validate:"omitempty"`
		ClosedAt      int            `json:"closed_at,omitempty" validate:"omitempty"`
		ClosestTaskAt int            `json:"closest_task_at,omitempty" validate:"omitempty"`
		Tags          []*Tag         `json:"tags,omitempty" validate:"omitempty,dive,required"`
		CustomFields  []*CustomField `json:"custom_fields,omitempty" validate:"omitempty"`
		Contact       struct {
			ID    []int  `json:"id" validate:"omitempty,dive,required"`
			Links *Links `json:"_links" validate:"omitempty"`
		} `json:"contacts,omitempty" validate:"omitempty"`
		StatusID int `json:"status_id" validate:"required"`
		Sale     int `json:"sale,omitempty" validate:"omitempty"`
		Pipeline struct {
			ID    int    `json:"id" validate:"required"`
			Links *Links `json:"_links" validate:"required"`
		} `json:"pipeline" validate:"required"`
		Links *Links `json:"_links" validate:"required"`
	}
)

const (
	NoTasksLeadFilter         LeadRequestTasksFilter = 1
	InProgressTasksLeadFilter LeadRequestTasksFilter = 2

	ActiveLeadsLeadFilter LeadRequestActiveFilter = 1
)

var (
	leadArrayFields = [][]byte{
		[]byte("tags"),
		[]byte("custom_fields"),
	}
)

func (c *Client) AddLead(ctx context.Context, lead *LeadAdd) (int, error) {
	if err := c.validator.Struct(lead); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+leadsURI, &AddLeadRequest{Add: []*LeadAdd{lead}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) UpdateLead(ctx context.Context, lead *LeadUpdate) (int, error) {
	if err := c.validator.Struct(lead); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+leadsURI, &UpdateLeadRequest{Update: []*LeadUpdate{lead}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetLeads(ctx context.Context, reqParams *LeadRequestParams) ([]*Lead, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(map[string]string)
	if reqParams.ID != nil {
		addValues["id"] = joinIntSlice(reqParams.ID)
	}
	if reqParams.LimitRows != 0 {
		addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
		if reqParams.LimitOffset != 0 {
			addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
		}
	}
	if reqParams.ResponsibleUserID != 0 {
		addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
	}
	if reqParams.Query != "" {
		addValues["query"] = reqParams.Query
	}
	if reqParams.Status != nil {
		addValues["status"] = joinIntSlice(reqParams.Status)
	}

	body, err := c.doGet(ctx, c.baseURL+leadsURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	leadResponse := new(GetLeadResponse)
	err = json.Unmarshal(body, leadResponse)
	if err != nil {
		// fix bad json serialization, where nil array is serialized as nil object
		body = fixBadArraySerialization(body, leadArrayFields)
		err = json.Unmarshal(body, leadResponse)
		if err != nil {
			return nil, err
		}
	}

	if leadResponse.Response != nil {
		return nil, leadResponse.Response
	}

	err = c.validator.Struct(leadResponse)
	if err != nil {
		return nil, err
	}

	if len(leadResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponseItems
	}

	return leadResponse.Embedded.Items, nil
}
