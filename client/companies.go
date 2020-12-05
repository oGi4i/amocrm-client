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
	GetCompaniesRequestWith string

	GetCompaniesOrderBy string

	GetCompaniesOrder struct {
		By     GetCompaniesOrderBy `validate:"required,oneof=id updated_at"`
		Method request.OrderMethod `validate:"required,oneof=asc desc"`
	}

	GetCompaniesRequestFilter struct {
		ID                *request.Filter   `validate:"omitempty"`                    // Фильтр по ID компании
		Name              *request.Filter   `validate:"omitempty"`                    // Фильтр по названию компании
		CreatedBy         *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, создавшего компанию
		UpdatedBy         *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, изменившего компанию
		ResponsibleUserID *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, ответственного за компанию
		CreatedAt         *request.Filter   `validate:"omitempty"`                    // Фильтр по дате создания компании
		UpdatedAt         *request.Filter   `validate:"omitempty"`                    // Фильтр по дате изменения компании
		ClosestTaskAt     *request.Filter   `validate:"omitempty"`                    // Фильтр по дате ближайшей задачи к выполнению по компании
		CustomFieldValues []*request.Filter `validate:"omitempty,gt=0,dive,required"` // Фильтр по дополнительным полям, привязанным к компании
	}

	GetCompaniesRequestParams struct {
		With   []GetCompaniesRequestWith  `validate:"omitempty,dive,oneof=catalog_elements leads customers contacts"` // Дополнительные параметры запроса, позволяющие получить больше данных в ответе
		Page   uint64                     `validate:"omitempty"`                                                      // Страница выборки
		Limit  uint64                     `validate:"omitempty,lte=250"`                                              // Количество возвращаемых сущностей за один запрос (максимум - 250)
		Query  string                     `validate:"omitempty"`                                                      // Поисковый запрос (осуществляет поиск по заполненным полям сущности)
		Filter *GetCompaniesRequestFilter `validate:"omitempty"`                                                      // Фильтр
		Order  *GetCompaniesOrder         `validate:"omitempty"`                                                      // Сортировка результатов
	}

	GetCompaniesResponseEmbedded struct {
		Companies []*domain.Company `json:"companies" validate:"omitempty,dive,required"`
	}

	GetCompaniesResponse struct {
		Page          uint64                        `json:"_page" validate:"required"`
		Links         *domain.Links                 `json:"_links" validate:"required"`
		Embedded      *GetCompaniesResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError              `json:"response" validate:"omitempty"`
	}

	ModifyCompaniesEmbedded struct {
		Tags []*domain.Tag `json:"tags" validate:"omitempty"`
	}

	AddCompaniesRequestData struct {
		Name               string                      `json:"name,omitempty" validate:"omitempty"`
		ResponsibleUserID  uint64                      `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CreatedBy          uint64                      `json:"create_by,omitempty" validate:"omitempty"`
		UpdatedBy          uint64                      `json:"updated_by,omitempty" validate:"omitempty"`
		CreatedAt          uint64                      `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt          uint64                      `json:"updated_at,omitempty" validate:"omitempty"`
		CustomFieldsValues []*domain.UpdateCustomField `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded           *ModifyCompaniesEmbedded    `json:"embedded,omitempty" validate:"omitempty,gt=0,dive,required"`
		RequestID          string                      `json:"request_id,omitempty" validate:"omitempty"`
	}

	AddCompaniesRequest struct {
		Add []*AddCompaniesRequestData `validate:"required,gt=0,dive,required"`
	}

	AddCompaniesResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		RequestID string        `json:"request_id,omitempty" validate:"omitempty"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	AddCompaniesResponseEmbedded struct {
		Companies []*AddCompaniesResponseItem `json:"companies" validate:"required,gt=0,dive,required"`
	}

	AddCompaniesResponse struct {
		Links         *domain.Links                 `json:"_links" validate:"required"`
		Embedded      *AddCompaniesResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError              `json:"response" validate:"omitempty"`
	}

	UpdateCompaniesRequestData struct {
		ID                 uint64                      `json:"id" validate:"required"`
		Name               string                      `json:"name,omitempty" validate:"omitempty"`
		ResponsibleUserID  uint64                      `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CreatedBy          uint64                      `json:"create_by,omitempty" validate:"omitempty"`
		UpdatedBy          uint64                      `json:"updated_by,omitempty" validate:"omitempty"`
		CreatedAt          uint64                      `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt          uint64                      `json:"updated_at,omitempty" validate:"omitempty"`
		CustomFieldsValues []*domain.UpdateCustomField `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded           *ModifyCompaniesEmbedded    `json:"embedded,omitempty" validate:"omitempty,gt=0,dive,required"`
		RequestID          string                      `json:"request_id,omitempty" validate:"omitempty"`
	}

	UpdateCompaniesRequest struct {
		Update []*UpdateCompaniesRequestData `validate:"required,gt=0,dive,required"`
	}

	UpdateCompaniesResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		Name      string        `json:"name" validate:"required"`
		UpdatedAt uint64        `json:"updated_at" validate:"required"`
		RequestID string        `json:"request_id,omitempty" validate:"omitempty"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	UpdateCompaniesResponseEmbedded struct {
		Companies []*UpdateCompaniesResponseItem `json:"companies" validate:"required,gt=0,dive,required"`
	}

	UpdateCompaniesResponse struct {
		Links         *domain.Links                    `json:"_links" validate:"required"`
		Embedded      *UpdateCompaniesResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError                 `json:"response" validate:"omitempty"`
	}
)

const (
	CatalogElemetsGetCompaniesRequestWith GetCompaniesRequestWith = "catalog_elements" // Добавляет в ответ связанные с компанией элементы списков
	LeadsGetCompaniesRequestWith          GetCompaniesRequestWith = "leads"            // Добавляет в ответ связанные с компанией сделки
	CustomersGetCompaniesRequestWith      GetCompaniesRequestWith = "customers"        // Добавляет в ответ связанных с компанией покупателей
	ContactsGetCompaniesRequestWith       GetCompaniesRequestWith = "contacts"         // Добавляет в ответ связанные с компанией контакты
)

const (
	IDGetCompaniesOrderBy        GetCompaniesOrderBy = "id"         // Сортировка по ID компании
	UpdatedAtGetCompaniesOrderBy GetCompaniesOrderBy = "updated_at" // Сортировка по дате изменения компании
)

func (c GetCompaniesRequestWith) String() string {
	return string(c)
}

func (o *GetCompaniesOrder) appendToQuery(params url.Values) {
	params.Add(fmt.Sprintf("order[%s]", string(o.By)), string(o.Method))
}

//nolint:dupl
func (f *GetCompaniesRequestFilter) validate() error {
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

	if f.CustomFieldValues != nil {
		for _, cf := range f.CustomFieldValues {
			if !cf.IsSimpleCustomFieldFilter() && !cf.IsMultipleCustomFieldFilter() && !cf.IsIntervalCustomFieldFilter() {
				return errors.New("CustomFieldValues filter must be custom field specific type")
			}
		}
	}

	return nil
}

func (f *GetCompaniesRequestFilter) appendGetRequestFilter(params url.Values) {
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

func (c *Client) AddCompanies(ctx context.Context, req *AddCompaniesRequest) ([]*AddCompaniesResponseItem, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+companiesURI, http.MethodPost, req.Add)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(AddCompaniesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response.Embedded.Companies, nil
}

func (c *Client) UpdateCompanies(ctx context.Context, req *UpdateCompaniesRequest) ([]*UpdateCompaniesResponseItem, error) {
	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+companiesURI, http.MethodPatch, req.Update)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(UpdateCompaniesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response.Embedded.Companies, nil
}

func (c *Client) UpdateCompany(ctx context.Context, companyID uint64, req *UpdateCompaniesRequestData) (*UpdateCompaniesResponseItem, error) {
	if companyID == 0 {
		return nil, ErrInvalidCompanyID
	}

	if err := c.validator.Struct(req); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+companiesURI+"/"+strconv.FormatUint(companyID, 10), http.MethodPatch, req)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(UpdateCompaniesResponseItem)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetCompanyByID(ctx context.Context, companyID uint64, with []GetCompaniesRequestWith) (*domain.Company, error) {
	if companyID == 0 {
		return nil, ErrInvalidCompanyID
	}

	params := make(url.Values)
	if with != nil {
		params.Add("with", joinGetCompaniesRequestWith(with))
	}

	body, err := c.doGet(ctx, c.baseURL+companiesURI+"/"+strconv.FormatUint(companyID, 10), params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(domain.Company)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//nolint:dupl
func (c *Client) GetCompanies(ctx context.Context, reqParams *GetCompaniesRequestParams) ([]*domain.Company, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	params := make(url.Values)
	if reqParams.With != nil {
		params.Add("with", joinGetCompaniesRequestWith(reqParams.With))
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
		err := reqParams.Filter.validate()
		if err != nil {
			return nil, err
		}

		reqParams.Filter.appendGetRequestFilter(params)
	}
	if reqParams.Order != nil {
		reqParams.Order.appendToQuery(params)
	}

	body, err := c.doGet(ctx, c.baseURL+companiesURI, params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(GetCompaniesResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	if len(response.Embedded.Companies) == 0 {
		return nil, ErrEmptyResponse
	}

	return response.Embedded.Companies, nil
}

func joinGetCompaniesRequestWith(with []GetCompaniesRequestWith) string {
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
