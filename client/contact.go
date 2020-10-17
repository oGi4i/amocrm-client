package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ogi4i/amocrm-client/domain"
)

type (
	ContactGetRequestParams struct {
		ID                []int  `validate:"omitempty,gt=0,dive,required"`
		LimitRows         int    `validate:"required_with=LimitOffset,lte=500"`
		LimitOffset       int    `validate:"omitempty"`
		ResponsibleUserID int    `validate:"omitempty"`
		Query             string `validate:"omitempty"`
	}

	ContactAddRequestData struct {
		Name              string                      `json:"name" validate:"required"`
		CreatedAt         int                         `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                         `json:"updated_at,string,omitempty" validate:"omitempty"`
		ResponsibleUserID int                         `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		CreatedBy         int                         `json:"created_by,string,omitempty" validate:"omitempty"`
		CompanyName       string                      `json:"company_name,omitempty" validate:"omitempty"`
		Tags              string                      `json:"tags,omitempty" validate:"omitempty"`
		LeadsID           []string                    `json:"leads_id,omitempty" validate:"omitempty,gt=0,dive,required"`
		CustomersID       int                         `json:"customers_id,string,omitempty" validate:"omitempty"`
		CompanyID         int                         `json:"company_id,string,omitempty" validate:"omitempty"`
		CustomFields      []*domain.UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,required"`
	}

	ContactUpdateRequestData struct {
		ID                int                         `json:"id,string" validate:"required"`
		Name              string                      `json:"name,omitempty" validate:"omitempty"`
		CreatedAt         int                         `json:"created_at,string,omitempty" validate:"omitempty"`
		UpdatedAt         int                         `json:"updated_at,string" validate:"required"`
		ResponsibleUserID int                         `json:"responsible_user_id,string,omitempty" validate:"omitempty"`
		CreatedBy         int                         `json:"created_by,string,omitempty" validate:"omitempty"`
		CompanyName       string                      `json:"company_name,omitempty" validate:"omitempty"`
		Tags              string                      `json:"tags,omitempty" validate:"omitempty"`
		LeadsID           []string                    `json:"leads_id,omitempty" validate:"omitempty,gt=0,required"`
		CustomersID       int                         `json:"customers_id,string,omitempty" validate:"omitempty"`
		CompanyID         int                         `json:"company_id,string,omitempty" validate:"omitempty"`
		CustomFields      []*domain.UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,required"`
		Unlink            *domain.Unlink              `json:"unlink,omitempty" validate:"omitempty"`
	}

	ContactAddRequest struct {
		Add []*ContactAddRequestData `json:"add" validate:"required,dive,required"`
	}

	ContactUpdateRequest struct {
		Update []*ContactUpdateRequestData `json:"update" validate:"required,dive,required"`
	}

	ContactGetResponse struct {
		Links    *domain.Links `json:"_links" validate:"omitempty"`
		Embedded struct {
			Items []*domain.Contact `json:"items" validate:"required,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
		ErrorResponse *domain.AmoError `json:"response,omitempty" validate:"omitempty"`
	}
)

var (
	contactArrayFields = [][]byte{
		[]byte("tags"),
		[]byte("custom_fields"),
	}
)

func (c *Client) AddContact(ctx context.Context, contact *ContactAddRequestData) (int, error) {
	if err := c.validator.Struct(contact); err != nil {
		return 0, err
	}

	resp, err := c.do(ctx, c.baseURL+contactsURI, http.MethodPost, &ContactAddRequest{Add: []*ContactAddRequestData{contact}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) UpdateContact(ctx context.Context, contact *ContactUpdateRequestData) (int, error) {
	if err := c.validator.Struct(contact); err != nil {
		return 0, err
	}

	resp, err := c.do(ctx, c.baseURL+contactsURI, http.MethodPost, &ContactUpdateRequest{Update: []*ContactUpdateRequestData{contact}})
	if err != nil {
		return 0, err
	}

	return c.getResponseID(resp)
}

func (c *Client) GetContacts(ctx context.Context, reqParams *ContactGetRequestParams) ([]*domain.Contact, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	addValues := make(url.Values)
	//if reqParams.ID != nil {
	//	addValues["id"] = joinIntSlice(reqParams.ID)
	//}
	//
	//if reqParams.LimitRows != 0 {
	//	addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
	//	if reqParams.LimitOffset != 0 {
	//		addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
	//	}
	//}
	//
	//if reqParams.ResponsibleUserID != 0 {
	//	addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
	//}
	//
	//if reqParams.Query != "" {
	//	addValues["query"] = reqParams.Query
	//}

	body, err := c.doGet(ctx, c.baseURL+contactsURI, addValues)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	contantsResponse := new(ContactGetResponse)
	err = json.Unmarshal(body, contantsResponse)
	if err != nil {
		// fix bad json serialization, where nil array is serialized as nil object
		body = fixBadArraySerialization(body, contactArrayFields)
		err = json.Unmarshal(body, contantsResponse)
		if err != nil {
			return nil, err
		}
	}

	if contantsResponse.ErrorResponse != nil {
		return nil, contantsResponse.ErrorResponse
	}

	err = c.validator.Struct(contantsResponse)
	if err != nil {
		return nil, err
	}

	if len(contantsResponse.Embedded.Items) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	return contantsResponse.Embedded.Items, nil
}
