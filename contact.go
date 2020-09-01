package amocrm

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

type (
	ContactRequestParams struct {
		ID                []int  `validate:"omitempty,gt=0,dive,required"`
		LimitRows         int    `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset       int    `validate:"omitempty"`
		ResponsibleUserID int    `validate:"omitempty"`
		Query             string `validate:"omitempty"`
	}

	ContactAdd struct {
		Name              string               `json:"name" validate:"required"`
		CreatedAt         int                  `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                  `json:"updated_at,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int                  `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		CreatedBy         int                  `json:"created_by,string,omitempty" validate:"omitempty"`
		CompanyName       string               `json:"company_name,omitempty" validate:"omitempty"`
		Tags              string               `json:"tags,omitempty" validate:"omitempty"`
		LeadsID           []string             `json:"leads_id,omitempty" validate:"omitempty,gt=0,dive,required"`
		CustomersID       int                  `json:"customers_id,string,omitempty" validate:"omitempty"`
		CompanyID         int                  `json:"company_id,string,omitempty" validate:"omitempty"`
		CustomFields      []*UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,required"`
	}

	ContactUpdate struct {
		ID                int                  `json:"id,string" validate:"required"`
		Name              string               `json:"name,omitempty" validate:"omitempty"`
		CreatedAt         int                  `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                  `json:"updated_at,string" validate:"required"`
		ResponsibleUserID int                  `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		CreatedBy         int                  `json:"created_by,string,omitempty" validate:"omitempty"`
		CompanyName       string               `json:"company_name,omitempty" validate:"omitempty"`
		Tags              string               `json:"tags,omitempty" validate:"omitempty"`
		LeadsID           []string             `json:"leads_id,omitempty" validate:"omitempty,gt=0,required"`
		CustomersID       int                  `json:"customers_id,string,omitempty" validate:"omitempty"`
		CompanyID         int                  `json:"company_id,string,omitempty" validate:"omitempty"`
		CustomFields      []*UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,required"`
		Unlink            *Unlink              `json:"unlink,omitempty" validate:"omitempty"`
	}

	AddContactRequest struct {
		Add []*ContactAdd `json:"add" validate:"required,dive,required"`
	}

	UpdateContactRequest struct {
		Update []*ContactUpdate `json:"update" validate:"required,dive,required"`
	}

	GetContactResponse struct {
		Links    *Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items []*Contact `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		Response *AmoError `json:"response,omitempty" validate:"omitempty"`
	}

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

var (
	contactArrayFields = [][]byte{
		[]byte("tags"),
		[]byte("custom_fields"),
	}
)

func (c *Client) AddContact(ctx context.Context, contact *ContactAdd) (int, error) {
	if err := c.validator.Struct(contact); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+contactsURI, &AddContactRequest{Add: []*ContactAdd{contact}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) UpdateContact(ctx context.Context, contact *ContactUpdate) (int, error) {
	if err := c.validator.Struct(contact); err != nil {
		return 0, err
	}

	resp, err := c.doPost(ctx, c.baseURL+contactsURI, &UpdateContactRequest{Update: []*ContactUpdate{contact}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetContacts(ctx context.Context, reqParams *ContactRequestParams) ([]*Contact, error) {
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

	body, err := c.doGet(ctx, c.baseURL+contactsURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	contantsResponse := new(GetContactResponse)
	err = json.Unmarshal(body, contantsResponse)
	if err != nil {
		// fix bad json serialization, where nil array is serialized as nil object
		body = fixBadArraySerialization(body, contactArrayFields)
		err = json.Unmarshal(body, contantsResponse)
		if err != nil {
			return nil, err
		}
	}

	if contantsResponse.Response != nil {
		return nil, contantsResponse.Response
	}

	err = c.validator.Struct(contantsResponse)
	if err != nil {
		return nil, err
	}

	if len(contantsResponse.Embedded.Items) == 0 {
		return nil, ErrEmptyResponseItems
	}

	return contantsResponse.Embedded.Items, nil
}

func joinIntSlice(s []int) string {
	out := new(strings.Builder)
	for i, n := range s {
		if i != len(s)-1 {
			out.WriteString(strconv.Itoa(n) + ",")
		} else {
			out.WriteString(strconv.Itoa(n))
		}
	}

	return out.String()
}

func fixBadArraySerialization(body []byte, fields [][]byte) []byte {
	var old, new []byte

	for _, field := range fields {
		old = append(append([]byte(`"`), field...), []byte(`":{}`)...)
		new = append(append([]byte(`"`), field...), []byte(`":[]`)...)
		body = bytes.ReplaceAll(body, old, new)
	}

	return body
}
