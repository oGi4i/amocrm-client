package client

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
)

func TestJoinGetLeadsRequestWith(t *testing.T) {
	testCases := []struct {
		name   string
		params []GetLeadsRequestWith
		want   string
	}{
		{
			name:   "Ни одного параметра",
			params: []GetLeadsRequestWith{},
			want:   "",
		},
		{
			name:   "Один параметр",
			params: []GetLeadsRequestWith{SourceIDGetLeadsRequestWith},
			want:   "source_id",
		},
		{
			name:   "Два параметра",
			params: []GetLeadsRequestWith{SourceIDGetLeadsRequestWith, LossReasonGetLeadsRequestWith},
			want:   "source_id,loss_reason",
		},
		{
			name: "Все параметры",
			params: []GetLeadsRequestWith{
				CatalogElementsGetLeadsRequestWith,
				IsPriceModifiedByRobotGetLeadsRequestWith,
				LossReasonGetLeadsRequestWith,
				ContactsGetLeadsRequestWith,
				OnlyDeletedGetLeadsRequestWith,
				SourceIDGetLeadsRequestWith,
			},
			want: "catalog_elements,is_price_modified_by_robot,loss_reason,contacts,only_deleted,source_id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, joinGetLeadsRequestWith(tc.params))
		})
	}
}

func TestGetLeadsRequestFilter(t *testing.T) {
	testCases := []struct {
		name   string
		filter *GetLeadsRequestFilter
		want   url.Values
	}{
		{
			name:   "Один ID",
			filter: &GetLeadsRequestFilter{ID: request.CreateSimpleFilter("id", "1")},
			want:   url.Values{"filter[id]": []string{"1"}},
		},
		{
			name:   "Несколько ID",
			filter: &GetLeadsRequestFilter{ID: request.CreateMultipleFilter("id", []string{"1", "2", "3"})},
			want:   url.Values{"filter[id][0]": []string{"1", "2", "3"}},
		},
		{
			name:   "Name",
			filter: &GetLeadsRequestFilter{Name: request.CreateSimpleFilter("name", "name_value")},
			want:   url.Values{"filter[name]": []string{"name_value"}},
		},
		{
			name:   "Несколько Name",
			filter: &GetLeadsRequestFilter{Name: request.CreateMultipleFilter("name", []string{"name_value_1", "name_value_2"})},
			want:   url.Values{"filter[name][0]": []string{"name_value_1", "name_value_2"}},
		},
		{
			name:   "Price",
			filter: &GetLeadsRequestFilter{Price: request.CreateSimpleFilter("price", "100")},
			want:   url.Values{"filter[price]": []string{"100"}},
		},
		{
			name:   "Диапазон Price",
			filter: &GetLeadsRequestFilter{Price: request.CreateIntervalFilter("price", "100", "200")},
			want:   url.Values{"filter[price][from]": []string{"100"}, "filter[price][to]": []string{"200"}},
		},
		{
			name:   "Statuses",
			filter: &GetLeadsRequestFilter{Statuses: request.CreateStatusFilter("123", "234")},
			want:   url.Values{"filter[statuses][0][pipeline_id]": []string{"123"}, "filter[statuses][0][status_id]": []string{"234"}},
		},
		{
			name:   "PipelineID",
			filter: &GetLeadsRequestFilter{PipelineID: request.CreateSimpleFilter("pipeline_id", "123")},
			want:   url.Values{"filter[pipeline_id]": []string{"123"}},
		},
		{
			name:   "Несколько PipelineID",
			filter: &GetLeadsRequestFilter{PipelineID: request.CreateMultipleFilter("pipeline_id", []string{"123", "234", "345"})},
			want:   url.Values{"filter[pipeline_id][0]": []string{"123", "234", "345"}},
		},
		{
			name:   "CreatedBy",
			filter: &GetLeadsRequestFilter{CreatedBy: request.CreateSimpleFilter("created_by", "123")},
			want:   url.Values{"filter[created_by]": []string{"123"}},
		},
		{
			name:   "Несколько CreatedBy",
			filter: &GetLeadsRequestFilter{CreatedBy: request.CreateMultipleFilter("created_by", []string{"234", "345"})},
			want:   url.Values{"filter[created_by][0]": []string{"234", "345"}},
		},
		{
			name:   "UpdatedBy",
			filter: &GetLeadsRequestFilter{UpdatedBy: request.CreateSimpleFilter("updated_by", "123")},
			want:   url.Values{"filter[updated_by]": []string{"123"}},
		},
		{
			name:   "Несколько UpdatedBy",
			filter: &GetLeadsRequestFilter{UpdatedBy: request.CreateMultipleFilter("updated_by", []string{"234", "345"})},
			want:   url.Values{"filter[updated_by][0]": []string{"234", "345"}},
		},
		{
			name:   "ResponsibleUserID",
			filter: &GetLeadsRequestFilter{ResponsibleUserID: request.CreateSimpleFilter("responsible_user_id", "123")},
			want:   url.Values{"filter[responsible_user_id]": []string{"123"}},
		},
		{
			name:   "Несколько ResponsibleUserID",
			filter: &GetLeadsRequestFilter{ResponsibleUserID: request.CreateMultipleFilter("responsible_user_id", []string{"234", "345"})},
			want:   url.Values{"filter[responsible_user_id][0]": []string{"234", "345"}},
		},
		{
			name:   "Интервал CreatedAt",
			filter: &GetLeadsRequestFilter{CreatedAt: request.CreateIntervalFilter("created_at", "12345678", "23456789")},
			want:   url.Values{"filter[created_at][from]": []string{"12345678"}, "filter[created_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал UpdatedAt",
			filter: &GetLeadsRequestFilter{UpdatedAt: request.CreateIntervalFilter("updated_at", "12345678", "23456789")},
			want:   url.Values{"filter[updated_at][from]": []string{"12345678"}, "filter[updated_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал ClosedAt",
			filter: &GetLeadsRequestFilter{ClosedAt: request.CreateIntervalFilter("closed_at", "12345678", "23456789")},
			want:   url.Values{"filter[closed_at][from]": []string{"12345678"}, "filter[closed_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал ClosestTaskAt",
			filter: &GetLeadsRequestFilter{ClosestTaskAt: request.CreateIntervalFilter("closest_task_at", "12345678", "23456789")},
			want:   url.Values{"filter[closest_task_at][from]": []string{"12345678"}, "filter[closest_task_at][to]": []string{"23456789"}},
		},
		{
			name:   "Простой CustomField",
			filter: &GetLeadsRequestFilter{CustomFieldValues: []*request.Filter{request.CreateSimpleCustomFieldFilter("123", "custom_field_value")}},
			want:   url.Values{"filter[custom_fields_values][123][]": []string{"custom_field_value"}},
		},
		{
			name:   "Диапазон CustomField",
			filter: &GetLeadsRequestFilter{CustomFieldValues: []*request.Filter{request.CreateIntervalCustomFieldFilter("123", "12345678", "23456789")}},
			want:   url.Values{"filter[custom_fields_values][123][from]": []string{"12345678"}, "filter[custom_fields_values][123][to]": []string{"23456789"}},
		},
		{
			name: "Несколько CustomField",
			filter: &GetLeadsRequestFilter{CustomFieldValues: []*request.Filter{
				request.CreateSimpleCustomFieldFilter("123", "custom_field_value"),
				request.CreateMultipleCustomFieldFilter("234", []string{"custom_field_value_1", "custom_field_value_2"}),
				request.CreateIntervalCustomFieldFilter("345", "12345678", "23456789"),
			},
			},
			want: url.Values{
				"filter[custom_fields_values][123][]":     []string{"custom_field_value"},
				"filter[custom_fields_values][234][]":     []string{"custom_field_value_1", "custom_field_value_2"},
				"filter[custom_fields_values][345][from]": []string{"12345678"},
				"filter[custom_fields_values][345][to]":   []string{"23456789"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := make(url.Values)
			tc.filter.appendGetRequestFilter(params)
			assert.Equal(t, tc.want, params)
		})
	}
}

func TestGetLeadsRequestParamsValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров with в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{With: []GetLeadsRequestWith{}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Один параметр with в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{With: []GetLeadsRequestWith{LossReasonGetLeadsRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Несколько параметров with в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{With: []GetLeadsRequestWith{LossReasonGetLeadsRequestWith, ContactsGetLeadsRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Невалидный параметр with в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{With: []GetLeadsRequestWith{"with"}}
		assert.EqualError(t, v.Struct(req), "Key: 'GetLeadsRequestParams.With[0]' Error:Field validation for 'With[0]' failed on the 'oneof' tag")
	})

	t.Run("Превышен лимит элементов в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Limit: 260}
		assert.EqualError(t, v.Struct(req), "Key: 'GetLeadsRequestParams.Limit' Error:Field validation for 'Limit' failed on the 'lte' tag")
	})

	t.Run("Диапазонный фильтр по ID в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{ID: request.CreateIntervalFilter("id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ID filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по Name в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{Name: request.CreateIntervalFilter("name", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "Name filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по Price в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{Price: request.CreateMultipleFilter("price", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "Price filter must be simple or interval type")
	})

	t.Run("Простой фильтр по Status в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{Statuses: request.CreateSimpleFilter("status", "123")}}
		assert.EqualError(t, req.Filter.validate(), "Statuses filter must be status type")
	})

	t.Run("Множественный фильтр по Status в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{Statuses: request.CreateMultipleFilter("status", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "Statuses filter must be status type")
	})

	t.Run("Диапазонный фильтр по Status в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{Statuses: request.CreateIntervalFilter("status", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "Statuses filter must be status type")
	})

	t.Run("Диапазонный фильтр по PipelineID в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{PipelineID: request.CreateIntervalFilter("pipeline_id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "PipelineID filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по CreatedBy в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{CreatedBy: request.CreateIntervalFilter("created_by", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "CreatedBy filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по UpdatedBy в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{UpdatedBy: request.CreateIntervalFilter("updated_by", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedBy filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по ResponsibleUserID в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{ResponsibleUserID: request.CreateIntervalFilter("responsible_user_id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ResponsibleUserID filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по CreatedAt в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{CreatedAt: request.CreateMultipleFilter("created_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "CreatedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по UpdatedAt в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{UpdatedAt: request.CreateMultipleFilter("updated_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по ClosedAt в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{ClosedAt: request.CreateMultipleFilter("closed_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "ClosedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по ClosestTaskAt в запросе", func(t *testing.T) {
		req := &GetLeadsRequestParams{Filter: &GetLeadsRequestFilter{ClosestTaskAt: request.CreateMultipleFilter("closest_task_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "ClosestTaskAt filter must be simple or interval type")
	})
}

func TestAddLeadsRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров CustomFieldsValues в запросе", func(t *testing.T) {
		data := &AddLeadsRequestData{CustomFieldsValues: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(data), "Key: 'AddLeadsRequestData.CustomFieldsValues' Error:Field validation for 'CustomFieldsValues' failed on the 'gt' tag")
	})
}

func TestUpdateLeadsRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Нет ID сделки в запросе", func(t *testing.T) {
		data := &UpdateLeadsRequestData{}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateLeadsRequestData.ID' Error:Field validation for 'ID' failed on the 'required' tag")
	})

	t.Run("Пустой массив параметров CustomFieldsValues в запросе", func(t *testing.T) {
		data := &UpdateLeadsRequestData{ID: 1, CustomFieldsValues: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateLeadsRequestData.CustomFieldsValues' Error:Field validation for 'CustomFieldsValues' failed on the 'gt' tag")
	})
}

//nolint:dupl
func TestGetLeadsResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &GetLeadsResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'GetLeadsResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag
Key: 'GetLeadsResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Leads в ответе", func(t *testing.T) {
		data := &GetLeadsResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetLeadsResponseEmbedded{Leads: []*domain.Lead{}}}
		assert.NoError(t, v.Struct(data))
	})

	t.Run("Ни одного обязательного параметра Lead в ответе", func(t *testing.T) {
		data := &GetLeadsResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetLeadsResponseEmbedded{Leads: []*domain.Lead{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'GetLeadsResponse.Embedded.Leads[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].Price' Error:Field validation for 'Price' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].ResponsibleUserID' Error:Field validation for 'ResponsibleUserID' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].StatusID' Error:Field validation for 'StatusID' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].PipelineID' Error:Field validation for 'PipelineID' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].UpdatedBy' Error:Field validation for 'UpdatedBy' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'GetLeadsResponse.Embedded.Leads[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetLeads(t *testing.T) {
	const sampleGetLeadsResponseBody = `{"_page":2,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=2"},"next":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=3"},"first":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=1"},"prev":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=1"}},"_embedded":{"leads":[{"id":19619,"name":"Сделка для примера","price":46333,"responsible_user_id":123321,"group_id":625,"status_id":142,"pipeline_id":1300,"loss_reason_id":null,"source_id":null,"created_by":321123,"updated_by":321123,"created_at":1453279607,"updated_at":1502193501,"closed_at":1483005931,"closest_task_at":null,"is_deleted":false,"custom_fields_values":null,"score":null,"account_id":5135160,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/19619"}},"_embedded":{"tags":[],"companies":[]}},{"id":14460,"name":"Сделка для примера 2","price":655,"responsible_user_id":123321,"group_id":625,"status_id":142,"pipeline_id":1300,"loss_reason_id":null,"source_id":null,"created_by":321123,"updated_by":321123,"created_at":1453279607,"updated_at":1502193501,"closed_at":1483005931,"closest_task_at":null,"is_deleted":false,"custom_fields_values":null,"score":null,"account_id":1351360,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/14460"}},"_embedded":{"tags":[],"companies":[]}}]}}`

	requestParamsWant := url.Values{
		"with":          []string{"contacts,loss_reason"},
		"page":          []string{"2"},
		"limit":         []string{"2"},
		"filter[id][0]": []string{"14460", "19619"},
		"order[id]":     []string{"desc"},
		"query":         []string{"query_value"},
	}

	sampleGetLeadsRequestParams := &GetLeadsRequestParams{
		With:  []GetLeadsRequestWith{ContactsGetLeadsRequestWith, LossReasonGetLeadsRequestWith},
		Page:  2,
		Limit: 2,
		Query: "query_value",
		Filter: &GetLeadsRequestFilter{
			ID: request.CreateMultipleFilter("id", []string{"14460", "19619"}),
		},
		Order: &GetLeadsOrder{
			By:     IDGetLeadsOrderBy,
			Method: request.DescendingOrderMethod,
		},
	}

	responseWant := []*domain.Lead{
		{
			ID:                19619,
			Name:              "Сделка для примера",
			Price:             46333,
			ResponsibleUserID: 123321,
			GroupID:           625,
			StatusID:          142,
			PipelineID:        1300,
			CreatedBy:         321123,
			UpdatedBy:         321123,
			CreatedAt:         1453279607,
			UpdatedAt:         1502193501,
			ClosedAt:          1483005931,
			AccountID:         5135160,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/19619"}},
			Embedded: &domain.LeadEmbedded{
				Tags:      []*domain.Tag{},
				Companies: []*domain.EmbeddedCompany{},
			},
		},
		{
			ID:                14460,
			Name:              "Сделка для примера 2",
			Price:             655,
			ResponsibleUserID: 123321,
			GroupID:           625,
			StatusID:          142,
			PipelineID:        1300,
			CreatedBy:         321123,
			UpdatedBy:         321123,
			CreatedAt:         1453279607,
			UpdatedAt:         1502193501,
			ClosedAt:          1483005931,
			AccountID:         1351360,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/14460"}},
			Embedded: &domain.LeadEmbedded{
				Tags:      []*domain.Tag{},
				Companies: []*domain.EmbeddedCompany{},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetLeadsResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.NoError(t, err)
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив в ответе", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_page":2,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=2"}},"_embedded":{"leads":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"page":2,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=2"},"next":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=3"},"first":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=1"},"prev":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=1"}},"_embedded":{"leads":[{"id":19619,"name":"Сделка для примера","price":46333,"responsible_user_id":123321,"group_id":625,"status_id":142,"pipeline_id":1300,"loss_reason_id":null,"source_id":null,"created_by":321123,"updated_by":321123,"created_at":1453279607,"updated_at":1502193501,"closed_at":1483005931,"closest_task_at":null,"is_deleted":false,"custom_fields_values":null,"score":null,"account_id":5135160,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/19619"}},"_embedded":{"tags":[],"companies":[]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.EqualError(t, err, "Key: 'GetLeadsResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestGetLeadByID(t *testing.T) {
	const sampleGetLeadByIDResponseBody = `{"id":3912171,"name":"Example","price":12,"responsible_user_id":504141,"group_id":0,"status_id":143,"pipeline_id":3104455,"loss_reason_id":4203748,"source_id":null,"created_by":504141,"updated_by":504141,"created_at":1585299171,"updated_at":1590683337,"closed_at":1590683337,"closest_task_at":null,"is_deleted":false,"custom_fields_values":null,"score":null,"account_id":28805383,"is_price_modified_by_robot":false,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/3912171"}},"_embedded":{"tags":[{"id":100667,"name":"тест"}],"catalog_elements":[{"id":525439,"metadata":{"quantity":1,"catalog_id":4521}}],"loss_reason":[{"id":4203748,"name":"Пропала потребность","sort":100000,"created_at":1582117280,"updated_at":1582117280,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/loss_reasons/4203748"}}}],"companies":[{"id":10971463,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/10971463"}}}],"contacts":[]}}`

	requestParamsWant := url.Values{
		"with": []string{"loss_reason,catalog_elements"},
	}

	sampleGetLeadsRequestParams := &GetLeadsRequestParams{
		With: []GetLeadsRequestWith{LossReasonGetLeadsRequestWith, CatalogElementsGetLeadsRequestWith},
	}

	responseWant := &domain.Lead{
		ID:                3912171,
		Name:              "Example",
		Price:             12,
		ResponsibleUserID: 504141,
		GroupID:           0,
		StatusID:          143,
		PipelineID:        3104455,
		LossReasonID:      4203748,
		CreatedBy:         504141,
		UpdatedBy:         504141,
		CreatedAt:         1585299171,
		UpdatedAt:         1590683337,
		ClosedAt:          1590683337,
		AccountID:         28805383,
		Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/3912171"}},
		Embedded: &domain.LeadEmbedded{
			Tags: []*domain.Tag{{ID: 100667, Name: "тест"}},
			CatalogElements: []*domain.CatalogElement{
				{
					ID: 525439,
					Metadata: &domain.CatalogElementMetadata{
						Quantity:  1,
						CatalogID: 4521,
					},
				},
			},
			LossReasons: []*domain.LossReason{{ID: 4203748, Name: "Пропала потребность"}},
			Companies:   []*domain.EmbeddedCompany{{ID: 10971463, Links: &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/10971463"}}}},
			Contacts:    []*domain.Contact{},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetLeadByIDResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeadByID(ctx, 3912171, sampleGetLeadsRequestParams.With)
		assert.NoError(t, err)
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeadByID(ctx, 3912171, sampleGetLeadsRequestParams.With)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		responseGot, err := client.GetLeadByID(ctx, 0, sampleGetLeadsRequestParams.With)
		assert.EqualError(t, err, ErrInvalidLeadID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"id":3912171,"name":"Example","price":12,"responsible_user_id":504141,"group_id":0,"status_id":143,"pipeline_id":3104455,"loss_reason_id":4203748,"source_id":null,"created_by":504141,"updated_by":504141,"created_at":1585299171,"updated_at":1590683337,"closed_at":1590683337,"account_id":28805383}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetLeadByID(ctx, 3912171, sampleGetLeadsRequestParams.With)
		assert.EqualError(t, err, "Key: 'Lead.Links' Error:Field validation for 'Links' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestAddLeads(t *testing.T) {
	const (
		requestBodyWant            = `[{"name":"Сделка для примера 1","price":20000,"custom_fields_values":[{"field_id":294471,"values":[{"value":"Наш первый клиент"}]}]},{"name":"Сделка для примера 2","price":10000,"_embedded":{"tags":[{"id":2719}]}}]`
		sampleAddLeadsResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads"}},"_embedded":{"leads":[{"id":10185151,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/10185151"}}},{"id":10185153,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/10185153"}}}]}}`
	)

	sampleAddLeadsRequest := &AddLeadsRequest{
		Add: []*AddLeadsRequestData{
			{
				Name:      "Сделка для примера 1",
				Price:     20000,
				CreatedBy: 0,
				CustomFieldsValues: []*domain.UpdateCustomField{
					{ID: 294471, Values: []interface{}{map[string]string{"value": "Наш первый клиент"}}},
				},
			},
			{
				Name:     "Сделка для примера 2",
				Price:    10000,
				Embedded: &ModifyLeadsEmbedded{Tags: []*domain.Tag{{ID: 2719}}},
			},
		},
	}

	responseWant := []*AddLeadsResponseItem{
		{
			ID:        10185151,
			RequestID: "0",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/10185151"}},
		},
		{
			ID:        10185153,
			RequestID: "1",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/10185153"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAddLeadsResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddLeads(ctx, sampleAddLeadsRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустой ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddLeads(ctx, sampleAddLeadsRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads"}},"embedded":{}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddLeads(ctx, sampleAddLeadsRequest)
		assert.EqualError(t, err, "Key: 'AddLeadsResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateLeads(t *testing.T) {
	const (
		requestBodyWant               = `[{"id":54886,"status_id":143,"pipeline_id":47521,"closed_at":1589297221,"loss_reason_id":7323},{"id":54884,"price":50000,"status_id":525743,"pipeline_id":47521,"_embedded":{"tags":[]}}]`
		sampleUpdateLeadsResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads"}},"_embedded":{"leads":[{"id":54886,"updated_at":1589556420,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/54886"}}},{"id":54884,"updated_at":1589556420,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/54884"}}}]}}`
	)

	sampleUpdateLeadsRequest := &UpdateLeadsRequest{
		Update: []*UpdateLeadsRequestData{
			{
				ID:           54886,
				PipelineID:   47521,
				StatusID:     143,
				ClosedAt:     1589297221,
				LossReasonID: 7323,
			},
			{
				ID:         54884,
				Price:      50000,
				PipelineID: 47521,
				StatusID:   525743,
				Embedded: &ModifyLeadsEmbedded{
					Tags: []*domain.Tag{},
				},
			},
		},
	}

	responseWant := []*UpdateLeadsResponseItem{
		{
			ID:        54886,
			UpdatedAt: 1589556420,
			RequestID: "0",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/54886"}},
		},
		{
			ID:        54884,
			UpdatedAt: 1589556420,
			RequestID: "1",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/54884"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateLeadsResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateLeads(ctx, sampleUpdateLeadsRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустой ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateLeads(ctx, sampleUpdateLeadsRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads"}},"embedded":{}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateLeads(ctx, sampleUpdateLeadsRequest)
		assert.EqualError(t, err, "Key: 'UpdateLeadsResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateLead(t *testing.T) {
	const (
		requestBodyWant              = `{"id":54884,"price":50000,"status_id":525743,"pipeline_id":47521,"_embedded":{"tags":[]}}`
		sampleUpdateLeadResponseBody = `{"id":54884,"updated_at":1589556420,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/54884"}}}`
	)

	sampleUpdateLeadRequest := &UpdateLeadsRequestData{
		ID:         54884,
		Price:      50000,
		PipelineID: 47521,
		StatusID:   525743,
		Embedded: &ModifyLeadsEmbedded{
			Tags: []*domain.Tag{},
		},
	}

	responseWant := &UpdateLeadsResponseItem{
		ID:        54884,
		UpdatedAt: 1589556420,
		RequestID: "1",
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/54884"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateLeadResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, 54884, sampleUpdateLeadRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустой ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, 54884, sampleUpdateLeadRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, 0, sampleUpdateLeadRequest)
		assert.EqualError(t, err, ErrInvalidLeadID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":54884,"updated_at":1589556420,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/54884"}}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, 54884, sampleUpdateLeadRequest)
		assert.EqualError(t, err, "Key: 'UpdateLeadsResponseItem.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}
