package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
)

type (
	LeadGetRequestWith string

	LeadGetRequestParams struct {
		With   []LeadGetRequestWith  `validate:"omitempty,dive,oneof=catalog_elements is_price_modified_by_robot loss_reason contacts only_deleted source_id"` // Дополнительные параметры запроса, позволяющие получить больше данных в ответе
		Page   uint64                `validate:"omitempty"`                                                                                                    // Страница выборки
		Limit  uint64                `validate:"omitempty,lte=250"`                                                                                            // Количество возвращаемых сущностей за один запрос (максимум - 250)
		Query  string                `validate:"omitempty"`                                                                                                    // Поисковый запрос (осуществляет поиск по заполненным полям сущности)
		Filter *LeadGetRequestFilter `validate:"omitempty"`                                                                                                    // Фильтр
		Order  *request.Order        `validate:"omitempty"`                                                                                                    // Сортировка результатов
	}

	LeadGetRequestFilter struct {
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

	LeadUpdateEmbedded struct {
		Tags []*domain.Tag `json:"tags,omitempty" validate:"omitempty"`
	}

	LeadUpdateData struct {
		Name              string                      `json:"name,omitempty" validate:"omitempty"`
		Price             uint64                      `json:"price,omitempty" validate:"omitempty"`
		StatusID          uint64                      `json:"status_id,omitempty" validate:"omitempty"`
		PipelineID        uint64                      `json:"pipeline_id,omitempty" validate:"omitempty"`
		CreatedBy         uint64                      `json:"created_by,omitempty" validate:"omitempty"`
		UpdatedBy         uint64                      `json:"updated_by,omitempty" validate:"omitempty"`
		ClosedAt          uint64                      `json:"closed_at,omitempty" validate:"omitempty"`
		CreatedAt         uint64                      `json:"created_at,omitempty" validate:"omitempty"`
		UpdatedAt         uint64                      `json:"updated_at,omitempty" validate:"omitempty"`
		LossReasonID      uint64                      `json:"loss_reason_id,omitempty" validate:"omitempty"`
		ResponsibleUserID uint64                      `json:"responsible_user_id,omitempty" validate:"omitempty"`
		CustomFields      []*domain.UpdateCustomField `json:"custom_fields,omitempty" validate:"omitempty,gt=0,dive,required"`
		Embedded          *LeadUpdateEmbedded         `json:"_embedded,omitempty" validate:"omitempty"`
		RequestID         string                      `json:"request_id,omitempty" validate:"omitempty"`
	}

	LeadUpdateResponseItem struct {
		ID        uint64        `json:"id" validate:"required"`
		RequestID string        `json:"request_id,omitempty" validate:"omitempty"`
		Links     *domain.Links `json:"_links" validate:"required"`
	}

	LeadUpdateResponseEmbedded struct {
		Leads []*LeadUpdateResponseItem `json:"leads" validate:"required,dive,required"`
	}

	LeadUpdateResponse struct {
		Links    *domain.Links               `json:"_links" validate:"required"`
		Embedded *LeadUpdateResponseEmbedded `json:"_embedded" validate:"required"`
	}

	LeadGetResponseEmbedded struct {
		Leads []*domain.Lead `json:"leads" validate:"omitempty,gt=0,dive,required"`
	}

	LeadGetResponse struct {
		Page          uint64                   `json:"_page" validate:"required"`
		Links         *domain.Links            `json:"_links" validate:"required"`
		Embedded      *LeadGetResponseEmbedded `json:"_embedded" validate:"required"`
		ErrorResponse *domain.AmoError         `json:"response" validate:"omitempty"`
	}
)

const (
	CatalogElementsLeadRequestWith        LeadGetRequestWith = "catalog_elements"           // Добавляет в ответ данные элементов списков, привязанных к сделке
	IsPriceModifiedByRobotLeadRequestWith LeadGetRequestWith = "is_price_modified_by_robot" // Добавляет в ответ признак изменён ли в последний раз бюджет сделки роботом
	LossReasonLeadRequestWith             LeadGetRequestWith = "loss_reason"                // Добавляет в ответ причину отказа сделки
	ContactsLeadRequestWith               LeadGetRequestWith = "contacts"                   // Добавляет в ответ данные контактов, привязанных к сделке
	OnlyDeletedLeadRequestWith            LeadGetRequestWith = "only_deleted"               // Добавляет в ответ удалённые сделки, которые ещё находятся в корзине
	SourceIDLeadRequestWith               LeadGetRequestWith = "source_id"                  // Добавляет в ответ ID источника сделки
)

func (l LeadGetRequestWith) String() string {
	return string(l)
}

func (f *LeadGetRequestFilter) AppendGetRequestFilter(params url.Values) {
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

func (c *Client) AddLead(ctx context.Context, leads []*LeadUpdateData) ([]*LeadUpdateResponseItem, error) {
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
		return nil, domain.ErrEmptyResponse
	}

	resp := new(LeadUpdateResponse)
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

func (c *Client) UpdateLead(ctx context.Context, leads []*LeadUpdateData) ([]*LeadUpdateResponseItem, error) {
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
		return nil, domain.ErrEmptyResponse
	}

	resp := new(LeadUpdateResponse)
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

func (c *Client) GetLeadByID(ctx context.Context, leadID uint64, with []LeadGetRequestWith) (*domain.Lead, error) {
	params := make(url.Values)
	if with != nil {
		params.Add("with", joinLeadRequestWithSlice(with))
	}

	body, err := c.doGet(ctx, c.baseURL+leadsURI+"/"+strconv.FormatUint(leadID, 10), params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, domain.ErrEmptyResponse
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

func (c *Client) GetLeads(ctx context.Context, reqParams *LeadGetRequestParams) ([]*domain.Lead, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	params := make(url.Values)
	if reqParams.With != nil {
		params.Add("with", joinLeadRequestWithSlice(reqParams.With))
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
		reqParams.Filter.AppendGetRequestFilter(params)
	}
	if reqParams.Order != nil {
		reqParams.Order.AppendToQuery(params)
	}

	body, err := c.doGet(ctx, c.baseURL+leadsURI, params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	response := new(LeadGetResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	err = c.validator.Struct(response)
	if err != nil {
		return nil, err
	}

	if len(response.Embedded.Leads) == 0 {
		return nil, domain.ErrEmptyResponse
	}

	return response.Embedded.Leads, nil
}

func joinLeadRequestWithSlice(with []LeadGetRequestWith) string {
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
