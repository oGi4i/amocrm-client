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
	GetLeadsRequestWith string

	GetLeadsRequestParams struct {
		With   []GetLeadsRequestWith  `validate:"omitempty,dive,oneof=catalog_elements is_price_modified_by_robot loss_reason contacts only_deleted source_id"` // Дополнительные параметры запроса, позволяющие получить больше данных в ответе
		Page   uint64                 `validate:"omitempty"`                                                                                                    // Страница выборки
		Limit  uint64                 `validate:"omitempty,lte=250"`                                                                                            // Количество возвращаемых сущностей за один запрос (максимум - 250)
		Query  string                 `validate:"omitempty"`                                                                                                    // Поисковый запрос (осуществляет поиск по заполненным полям сущности)
		Filter *GetLeadsRequestFilter `validate:"omitempty"`                                                                                                    // Фильтр
		Order  *GetLeadsOrder         `validate:"omitempty"`                                                                                                    // Сортировка результатов
	}

	GetLeadsRequestFilter struct {
		ID                *request.Filter   `validate:"omitempty"`                    // Фильтр по ID сделки
		Name              *request.Filter   `validate:"omitempty"`                    // Фильтр по названию сделки
		Price             *request.Filter   `validate:"omitempty"`                    // Фильтр по бюджету сделки
		Statuses          *request.Filter   `validate:"omitempty"`                    // Фильтр по статусам сделки
		PipelineID        *request.Filter   `validate:"omitempty"`                    // Фильтр по ID воронки
		CreatedBy         *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, создавшего сделку
		UpdatedBy         *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, изменившего сделку
		ResponsibleUserID *request.Filter   `validate:"omitempty"`                    // Фильтр по ID пользователя, ответственного за сделку
		CreatedAt         *request.Filter   `validate:"omitempty"`                    // Фильтр по дате создания сделки
		UpdatedAt         *request.Filter   `validate:"omitempty"`                    // Фильтр по дате изменения сделки
		ClosedAt          *request.Filter   `validate:"omitempty"`                    // Фильтр по дате закрытия сделки
		ClosestTaskAt     *request.Filter   `validate:"omitempty"`                    // Фильтр по дате ближайшей задачи к выполнению по сделке
		CustomFieldValues []*request.Filter `validate:"omitempty,gt=0,dive,required"` // Фильтр по дополнительным полям, привязанным к сделке
	}

	GetLeadsOrderBy string

	GetLeadsOrder struct {
		By     GetLeadsOrderBy     `validate:"required,oneof=id created_at updated_at"`
		Method request.OrderMethod `validate:"required,oneof=asc desc"`
	}

	AddLeadRequestDataEmbedded struct {
		Tags []*domain.Tag `json:"tags" validate:"omitempty"`
	}

	AddLeadRequestData struct {
		Name               string                      `json:"name,omitempty" validate:"omitempty"`
		Price              uint64                      `json:"price,omitempty" validate:"omitempty"`
		StatusID           uint64                      `json:"status_id,omitempty" validate:"omitempty"`
		PipelineID         uint64                      `json:"pipeline_id,omitempty" validate:"omitempty"`
		CreatedBy          uint64                      `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy          uint64                      `json:"updated_by,omitempty" validate:"omitempty"`
		ClosedAt           uint64                      `json:"closed_at,omitempty" validate:"omitempty"`
		CreatedAt          uint64                      `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt          uint64                      `json:"updated_at,omitempty" validate:"omitempty"`
		LossReasonID       uint64                      `json:"loss_reason_id,omitempty" validate:"omitempty"`
		ResponsibleUserID  uint64                      `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CustomFieldsValues []*domain.UpdateCustomField `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded           *AddLeadRequestDataEmbedded `json:"_embedded,omitempty" validate:"omitempty"`
		RequestID          string                      `json:"request_id,omitempty" validate:"omitempty"`
	}

	UpdateLeadRequestDataEmbedded struct {
		Tags []*domain.Tag `json:"tags" validate:"omitempty,dive,required"`
	}

	UpdateLeadRequestData struct {
		ID                 uint64                         `json:"id" validate:"required"`
		Name               string                         `json:"name,omitempty" validate:"omitempty"`
		Price              uint64                         `json:"price,omitempty" validate:"omitempty"`
		StatusID           uint64                         `json:"status_id,omitempty" validate:"omitempty"`
		PipelineID         uint64                         `json:"pipeline_id,omitempty" validate:"omitempty"`
		CreatedBy          uint64                         `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy          uint64                         `json:"updated_by,omitempty" validate:"omitempty"`
		ClosedAt           uint64                         `json:"closed_at,omitempty" validate:"omitempty"`
		CreatedAt          uint64                         `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt          uint64                         `json:"updated_at,omitempty" validate:"omitempty"`
		LossReasonID       uint64                         `json:"loss_reason_id,omitempty" validate:"omitempty"`
		ResponsibleUserID  uint64                         `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CustomFieldsValues []*domain.UpdateCustomField    `json:"custom_fields_values,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded           *UpdateLeadRequestDataEmbedded `json:"_embedded,omitempty" validate:"omitempty"`
		RequestID          string                         `json:"request_id,omitempty" validate:"omitempty"`
	}

	AddLeadsResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		RequestID string        `json:"request_id,omitempty" validate:"omitempty"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	AddLeadsResponseEmbedded struct {
		Leads []*AddLeadsResponseItem `json:"leads" validate:"omitempty,gt=0,dive,required"`
	}

	AddLeadsResponse struct {
		Links         *domain.Links             `json:"_links" validate:"required"`
		Embedded      *AddLeadsResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError          `json:"response" validate:"omitempty"`
	}

	UpdateLeadsResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		UpdatedAt uint64        `json:"updated_at" validate:"required"`
		RequestID string        `json:"request_id,omitempty" validate:"required"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	UpdateLeadsResponseEmbedded struct {
		Leads []*UpdateLeadsResponseItem `json:"leads" validate:"required,dive,required"`
	}

	UpdateLeadsResponse struct {
		Links         *domain.Links                `json:"_links" validate:"required"`
		Embedded      *UpdateLeadsResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError             `json:"response" validate:"omitempty"`
	}

	GetLeadsResponseEmbedded struct {
		Leads []*domain.Lead `json:"leads" validate:"omitempty,dive,required"`
	}

	GetLeadsResponse struct {
		Page          uint64                    `json:"_page" validate:"required"`
		Links         *domain.Links             `json:"_links" validate:"required"`
		Embedded      *GetLeadsResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError          `json:"response" validate:"omitempty"`
	}
)

const (
	CatalogElementsGetLeadsRequestWith        GetLeadsRequestWith = "catalog_elements"           // Добавляет в ответ данные элементов списков, привязанных к сделке
	IsPriceModifiedByRobotGetLeadsRequestWith GetLeadsRequestWith = "is_price_modified_by_robot" // Добавляет в ответ признак изменён ли в последний раз бюджет сделки роботом
	LossReasonGetLeadsRequestWith             GetLeadsRequestWith = "loss_reason"                // Добавляет в ответ причину отказа сделки
	ContactsGetLeadsRequestWith               GetLeadsRequestWith = "contacts"                   // Добавляет в ответ данные контактов, привязанных к сделке
	OnlyDeletedGetLeadsRequestWith            GetLeadsRequestWith = "only_deleted"               // Добавляет в ответ удалённые сделки, которые ещё находятся в корзине
	SourceIDGetLeadsRequestWith               GetLeadsRequestWith = "source_id"                  // Добавляет в ответ ID источника сделки
)

const (
	IDGetLeadsOrderBy        GetLeadsOrderBy = "id"         // Сортировка по ID сделки
	CreatedAtGetLeadsOrderBy GetLeadsOrderBy = "created_at" // Сортировка по дате создания сделки
	UpdatedAtGetLeadsOrderBy GetLeadsOrderBy = "updated_at" // Сортировка по дате изменения сделки
)

func (l GetLeadsRequestWith) String() string {
	return string(l)
}

func (o *GetLeadsOrder) appendToQuery(params url.Values) {
	params.Add(fmt.Sprintf("order[%s]", string(o.By)), string(o.Method))
}

func (f *GetLeadsRequestFilter) validate() error {
	if f.ID != nil && !f.ID.IsSimpleFilter() && !f.ID.IsMultipleFilter() {
		return errors.New("ID filter must be simple or multiple type")
	}

	if f.Name != nil && !f.Name.IsSimpleFilter() && !f.Name.IsMultipleFilter() {
		return errors.New("Name filter must be simple or multiple type")
	}

	if f.Price != nil && !f.Price.IsSimpleFilter() && !f.Price.IsIntervalFilter() {
		return errors.New("Price filter must be simple or interval type")
	}

	if f.Statuses != nil && !f.Statuses.IsStatusFilter() {
		return errors.New("Statuses filter must be status type")
	}

	if f.PipelineID != nil && !f.PipelineID.IsSimpleFilter() && !f.PipelineID.IsMultipleFilter() {
		return errors.New("PipelineID filter must be simple or multiple type")
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

	if f.ClosedAt != nil && !f.ClosedAt.IsSimpleFilter() && !f.ClosedAt.IsIntervalFilter() {
		return errors.New("ClosedAt filter must be simple or interval type")
	}

	if f.ClosestTaskAt != nil && !f.ClosestTaskAt.IsSimpleFilter() && !f.ClosestTaskAt.IsIntervalFilter() {
		return errors.New("ClosestTaskAt filter must be simple or interval type")
	}

	return nil
}

func (f *GetLeadsRequestFilter) appendGetRequestFilter(params url.Values) {
	if f.ID != nil {
		f.ID.AppendToQuery(params)
	}

	if f.Name != nil {
		f.Name.AppendToQuery(params)
	}

	if f.Price != nil {
		f.Price.AppendToQuery(params)
	}

	if f.Statuses != nil {
		f.Statuses.AppendToQuery(params)
	}

	if f.PipelineID != nil {
		f.PipelineID.AppendToQuery(params)
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

	if f.ClosedAt != nil {
		f.ClosedAt.AppendToQuery(params)
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

func (c *Client) AddLeads(ctx context.Context, leads []*AddLeadRequestData) ([]*AddLeadsResponseItem, error) {
	for _, lead := range leads {
		if err := c.validator.Struct(lead); err != nil {
			return nil, err
		}
	}

	body, err := c.do(ctx, c.baseURL+leadsURI, http.MethodPost, leads)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(AddLeadsResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Leads, nil
}

func (c *Client) UpdateLeads(ctx context.Context, leads []*UpdateLeadRequestData) ([]*UpdateLeadsResponseItem, error) {
	for _, lead := range leads {
		if err := c.validator.Struct(lead); err != nil {
			return nil, err
		}
	}

	body, err := c.do(ctx, c.baseURL+leadsURI, http.MethodPatch, leads)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(UpdateLeadsResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Leads, nil
}

func (c *Client) UpdateLead(ctx context.Context, leadID uint64, reqData *UpdateLeadRequestData) (*UpdateLeadsResponseItem, error) {
	if err := c.validator.Struct(reqData); err != nil {
		return nil, err
	}

	body, err := c.do(ctx, c.baseURL+leadsURI+"/"+strconv.FormatUint(leadID, 10), http.MethodPatch, reqData)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(UpdateLeadsResponse)
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(resp)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Leads[0], nil
}

func (c *Client) GetLeadByID(ctx context.Context, leadID uint64, with []GetLeadsRequestWith) (*domain.Lead, error) {
	params := make(url.Values)
	if with != nil {
		params.Add("with", joinGetLeadsRequestWith(with))
	}

	body, err := c.doGet(ctx, c.baseURL+leadsURI+"/"+strconv.FormatUint(leadID, 10), params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	resp := new(domain.Lead)
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

func (c *Client) GetLeads(ctx context.Context, reqParams *GetLeadsRequestParams) ([]*domain.Lead, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	params := make(url.Values)
	if reqParams.With != nil {
		params.Add("with", joinGetLeadsRequestWith(reqParams.With))
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

	body, err := c.doGet(ctx, c.baseURL+leadsURI, params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(GetLeadsResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	if len(response.Embedded.Leads) == 0 {
		return nil, ErrEmptyResponse
	}

	return response.Embedded.Leads, nil
}

func joinGetLeadsRequestWith(with []GetLeadsRequestWith) string {
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
