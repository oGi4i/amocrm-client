package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
)

type (
	GetContactsRequestWith string

	GetContactsRequestParams struct {
		With   []GetContactsRequestWith  `validate:"omitempty,dive,oneof=catalog_elements leads customers"` // Дополнительные параметры запроса, позволяющие получить больше данных в ответе
		Page   uint64                    `validate:"omitempty"`                                             // Страница выборки
		Limit  uint64                    `validate:"omitempty,lte=250"`                                     // Количество возвращаемых сущностей за один запрос (максимум - 250)
		Query  string                    `validate:"omitempty"`                                             // Поисковый запрос (осуществляет поиск по заполненным полям сущности)
		Filter *GetContactsRequestFilter `validate:"omitempty"`                                             // Фильтр
		Order  *GetContactsOrder         `validate:"omitempty"`                                             // Сортировка результатов
	}

	GetContactsRequestFilter struct {
		ID                *request.Filter   `validate:"omitempty"`                    // Фильтр по ID контактов
		Name              *request.Filter   `validate:"omitempty"`                    // Фильтр по названию контакта
		CreatedBy         *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, создавшего контакта
		UpdatedBy         *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, изменившего контакта
		ResponsibleUserID *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, ответственного за контакт
		CreatedAt         *request.Filter   `validate:"omitempty"`                    // Фильтр по дате создания контакта
		UpdatedAt         *request.Filter   `validate:"omitempty"`                    // Фильтр по дате изменения контакта
		ClosestTaskAt     *request.Filter   `validate:"omitempty"`                    // Фильтр по дате ближайшей задачи к выполнению по контакту
		CustomFieldValues []*request.Filter `validate:"omitempty,gt=0,dive,required"` // Фильтр по дополнительным полям, привязанным к контакту
	}

	GetContactsOrderBy string

	GetContactsOrder struct {
		By     GetContactsOrderBy  `validate:"required,oneof=id updated_at"`
		Method request.OrderMethod `validate:"required,oneof=asc desc"`
	}

	ContactRequestDataEmbedded struct {
		Tags []*domain.Tag `json:"tags" validate:"omitempty,dive,required"`
	}

	AddContactRequestData struct {
		Name               string                      `json:"name,omitempty" validate:"omitempty"`
		FirstName          string                      `json:"first_name,omitempty" validate:"omitempty"`
		LastName           string                      `json:"last_name,omitempty" validate:"omitempty"`
		ResponsibleUserID  uint64                      `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CreatedBy          uint64                      `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy          uint64                      `json:"updated_by,omitempty" validate:"omitempty"`
		CreatedAt          uint64                      `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt          uint64                      `json:"updated_at,omitempty" validate:"omitempty"`
		CustomFieldsValues []*domain.UpdateCustomField `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded           *ContactRequestDataEmbedded `json:"_embedded,omitempty" validate:"omitempty"`
		RequestID          string                      `json:"request_id,omitempty" validate:"omitempty"`
	}

	AddContactsRequest struct {
		Add []*AddContactRequestData `validate:"required,gt=0,dive,required"`
	}

	AddContactsResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		RequestID string        `json:"request_id" validate:"required"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	AddContactsResponseEmbedded struct {
		Contacts []*AddContactsResponseItem `json:"contacts" validate:"omitempty,gt=0,dive,required"`
	}

	AddContactsResponse struct {
		Links         *domain.Links                `json:"_links" validate:"required"`
		Embedded      *AddContactsResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError             `json:"response,omitempty" validate:"omitempty"`
	}

	UpdateContactsRequestData struct {
		ID                 uint64                      `json:"id" validate:"required"`
		Name               string                      `json:"name,omitempty" validate:"omitempty"`
		FirstName          string                      `json:"first_name,omitempty" validate:"omitempty"`
		LastName           string                      `json:"last_name,omitempty" validate:"omitempty"`
		ResponsibleUserID  uint64                      `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CreatedBy          uint64                      `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy          uint64                      `json:"updated_by,omitempty" validate:"omitempty"`
		CreatedAt          uint64                      `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt          uint64                      `json:"updated_at,omitempty" validate:"omitempty"`
		CustomFieldsValues []*domain.UpdateCustomField `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded           *ContactRequestDataEmbedded `json:"_embedded,omitempty" validate:"omitempty"`
		RequestID          string                      `json:"request_id,omitempty" validate:"omitempty"`
	}

	UpdateContactsRequest struct {
		Update []*UpdateContactsRequestData `validate:"required,gt=0,dive,required"`
	}

	UpdateContactsResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		Name      string        `json:"name" validate:"required"`
		UpdatedAt uint64        `json:"updated_at" validate:"required"`
		RequestID string        `json:"request_id,omitempty" validate:"omitempty"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	UpdateContactsResponseEmbedded struct {
		Contacts []*UpdateContactsResponseItem `json:"contacts" validate:"omitempty,gt=0,dive,required"`
	}

	UpdateContactsResponse struct {
		Links         *domain.Links                   `json:"_links" validate:"required"`
		Embedded      *UpdateContactsResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError                `json:"response,omitempty" validate:"omitempty"`
	}

	GetContactsResponseEmbedded struct {
		Contacts []*domain.Contact `json:"contacts" validate:"omitempty,dive,required"`
	}

	GetContactsResponse struct {
		Page          uint64                       `json:"_page" validate:"required"`
		Links         *domain.Links                `json:"_links" validate:"required"`
		Embedded      *GetContactsResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError             `json:"response,omitempty" validate:"omitempty"`
	}
)

const (
	CatalogElementsGetContactsRequestWith GetContactsRequestWith = "catalog_elements" // Добавляет в ответ связанные с контактами элеметры списков
	LeadsGetContactsRequestWith           GetContactsRequestWith = "leads"            // Добавляет в ответ связанные с контактами сделки
	CustomersGetContactsRequestWith       GetContactsRequestWith = "customers"        // Добавляет в ответ связанных с контактами покупателей
)

const (
	IDGetContactsOrderBy        GetContactsOrderBy = "id"         // Сортировка по ID контакта
	UpdatedAtGetContactsOrderBy GetContactsOrderBy = "updated_at" // Сортировка по дате изменения контакта
)

func (c GetContactsRequestWith) String() string {
	return string(c)
}

func (o *GetContactsOrder) appendToQuery(params url.Values) {
	params.Add(fmt.Sprintf("order[%s]", string(o.By)), string(o.Method))
}

func (f *GetContactsRequestFilter) validate() error {
	if f.ID != nil && !f.ID.IsSimpleFilter() && !f.ID.IsMultipleFilter() {
		return errors.New("ID filter must be simple or multiple type")
	}

	if f.Name != nil && !f.Name.IsSimpleFilter() && !f.Name.IsMultipleFilter() {
		return errors.New("Name filter must be simple or multiple type")
	}

	if f.CreatedBy != nil && !f.CreatedBy.IsSimpleFilter() && !f.CreatedBy.IsMultipleFilter() {
		return errors.New("CreatedBy filter must be simple or multiple type")
	}

	if f.UpdatedBy != nil && !f.UpdatedBy.IsSimpleFilter() && !f.UpdatedBy.IsMultipleFilter() {
		return errors.New("UpdatedBy filter must be simple or multiple type")
	}

	if f.ResponsibleUserID != nil && !f.ResponsibleUserID.IsSimpleFilter() && !f.ResponsibleUserID.IsMultipleFilter() {
		return errors.New("ResponsibleUserID filter must be simple or multiple type")
	}

	if f.CreatedAt != nil && !f.CreatedAt.IsSimpleFilter() && !f.CreatedAt.IsIntervalFilter() {
		return errors.New("CreatedAt filter must be simple or interval type")
	}

	if f.UpdatedAt != nil && !f.UpdatedAt.IsSimpleFilter() && !f.UpdatedAt.IsIntervalFilter() {
		return errors.New("UpdatedAt filter must be simple or interval type")
	}

	if f.ClosestTaskAt != nil && !f.ClosestTaskAt.IsSimpleFilter() && !f.ClosestTaskAt.IsIntervalFilter() {
		return errors.New("ClosestTaskAt filter must be simple or interval type")
	}

	return nil
}

func (f *GetContactsRequestFilter) appendGetRequestFilter(params url.Values) {
	if f.ID != nil {
		f.ID.AppendToQuery(params)
	}

	if f.Name != nil {
		f.Name.AppendToQuery(params)
	}

	if f.CreatedBy != nil {
		f.CreatedBy.AppendToQuery(params)
	}

	if f.UpdatedBy != nil {
		f.UpdatedBy.AppendToQuery(params)
	}

	if f.ResponsibleUserID != nil {
		f.ResponsibleUserID.AppendToQuery(params)
	}

	if f.CreatedAt != nil {
		f.CreatedAt.AppendToQuery(params)
	}

	if f.UpdatedAt != nil {
		f.UpdatedAt.AppendToQuery(params)
	}

	if f.ClosestTaskAt != nil {
		f.ClosestTaskAt.AppendToQuery(params)
	}

	if f.CustomFieldValues != nil {
		for _, f := range f.CustomFieldValues {
			f.AppendToQuery(params)
		}
	}
}

func (c *Client) AddContacts(ctx context.Context, req *AddContactsRequest) ([]*AddContactsResponseItem, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+contactsURI, http.MethodPost, req.Add)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(AddContactsResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Contacts, nil
}

func (c *Client) UpdateContacts(ctx context.Context, req *UpdateContactsRequest) ([]*UpdateContactsResponseItem, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+contactsURI, http.MethodPatch, req.Update)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(UpdateContactsResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Contacts, nil
}

func (c *Client) UpdateContact(ctx context.Context, contactID uint64, req *UpdateContactsRequestData) (*UpdateContactsResponseItem, error) {
	if contactID == 0 {
		return nil, ErrInvalidContactID
	}

	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+contactsURI+"/"+strconv.FormatUint(contactID, 10), http.MethodPatch, req)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(UpdateContactsResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Contacts[0], nil
}

func (c *Client) GetContactByID(ctx context.Context, contactID uint64, with []GetContactsRequestWith) (*domain.Contact, error) {
	if contactID == 0 {
		return nil, ErrInvalidContactID
	}

	params := make(url.Values)
	if with != nil {
		params.Add("with", joinGetContactsRequestWith(with))
	}

	body, err := c.doGet(ctx, c.baseURL+contactsURI+"/"+strconv.FormatUint(contactID, 10), params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(domain.Contact)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetContacts(ctx context.Context, reqParams *GetContactsRequestParams) ([]*domain.Contact, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	params := make(url.Values)
	if reqParams.With != nil {
		params.Add("with", joinGetContactsRequestWith(reqParams.With))
	}
	if reqParams.Page != 0 {
		params.Add("page", strconv.FormatUint(reqParams.Page, 10))
	}
	if reqParams.Limit != 0 {
		params.Add("limit", strconv.FormatUint(reqParams.Limit, 10))
	}
	if reqParams.Query != "" {
		params.Add("query", reqParams.Query)
	}
	if reqParams.Filter != nil {
		reqParams.Filter.appendGetRequestFilter(params)
	}
	if reqParams.Order != nil {
		reqParams.Order.appendToQuery(params)
	}

	body, err := c.doGet(ctx, c.baseURL+contactsURI, params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(GetContactsResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	if len(response.Embedded.Contacts) == 0 {
		return nil, ErrEmptyResponse
	}

	return response.Embedded.Contacts, nil
}

func joinGetContactsRequestWith(with []GetContactsRequestWith) string {
	if len(with) == 0 {
		return ""
	}

	out := new(strings.Builder)
	for i, s := range with {
		if i != len(with)-1 {
			out.WriteString(s.String() + ",")
		} else {
			out.WriteString(s.String())
		}
	}
	return out.String()
}
