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
	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
	"github.com/stretchr/testify/assert"
)

func TestJoinGetCompaniesRequestWith(t *testing.T) {
	testCases := []struct {
		name   string
		params []GetCompaniesRequestWith
		want   string
	}{
		{
			name:   "Ни одного параметра",
			params: []GetCompaniesRequestWith{},
			want:   "",
		},
		{
			name:   "Один параметр",
			params: []GetCompaniesRequestWith{LeadsGetCompaniesRequestWith},
			want:   "leads",
		},
		{
			name:   "Два параметра",
			params: []GetCompaniesRequestWith{LeadsGetCompaniesRequestWith, CatalogElemetsGetCompaniesRequestWith},
			want:   "leads,catalog_elements",
		},
		{
			name: "Все параметры",
			params: []GetCompaniesRequestWith{
				CatalogElemetsGetCompaniesRequestWith,
				LeadsGetCompaniesRequestWith,
				CustomersGetCompaniesRequestWith,
				ContactsGetCompaniesRequestWith,
			},
			want: "catalog_elements,leads,customers,contacts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, joinGetCompaniesRequestWith(tc.params))
		})
	}
}

func TestGetCompaniesRequestFilter(t *testing.T) {
	testCases := []struct {
		name   string
		filter *GetCompaniesRequestFilter
		want   url.Values
	}{
		{
			name:   "Один ID",
			filter: &GetCompaniesRequestFilter{ID: request.CreateSimpleFilter("id", "1")},
			want:   url.Values{"filter[id]": []string{"1"}},
		},
		{
			name:   "Несколько ID",
			filter: &GetCompaniesRequestFilter{ID: request.CreateMultipleFilter("id", []string{"1", "2", "3"})},
			want:   url.Values{"filter[id][0]": []string{"1", "2", "3"}},
		},
		{
			name:   "Name",
			filter: &GetCompaniesRequestFilter{Name: request.CreateSimpleFilter("name", "name_value")},
			want:   url.Values{"filter[name]": []string{"name_value"}},
		},
		{
			name:   "Несколько Name",
			filter: &GetCompaniesRequestFilter{Name: request.CreateMultipleFilter("name", []string{"name_value_1", "name_value_2"})},
			want:   url.Values{"filter[name][0]": []string{"name_value_1", "name_value_2"}},
		},
		{
			name:   "CreatedBy",
			filter: &GetCompaniesRequestFilter{CreatedBy: request.CreateSimpleFilter("created_by", "123")},
			want:   url.Values{"filter[created_by]": []string{"123"}},
		},
		{
			name:   "Несколько CreatedBy",
			filter: &GetCompaniesRequestFilter{CreatedBy: request.CreateMultipleFilter("created_by", []string{"234", "345"})},
			want:   url.Values{"filter[created_by][0]": []string{"234", "345"}},
		},
		{
			name:   "UpdatedBy",
			filter: &GetCompaniesRequestFilter{UpdatedBy: request.CreateSimpleFilter("updated_by", "123")},
			want:   url.Values{"filter[updated_by]": []string{"123"}},
		},
		{
			name:   "Несколько UpdatedBy",
			filter: &GetCompaniesRequestFilter{UpdatedBy: request.CreateMultipleFilter("updated_by", []string{"234", "345"})},
			want:   url.Values{"filter[updated_by][0]": []string{"234", "345"}},
		},
		{
			name:   "ResponsibleUserID",
			filter: &GetCompaniesRequestFilter{ResponsibleUserID: request.CreateSimpleFilter("responsible_user_id", "123")},
			want:   url.Values{"filter[responsible_user_id]": []string{"123"}},
		},
		{
			name:   "Несколько ResponsibleUserID",
			filter: &GetCompaniesRequestFilter{ResponsibleUserID: request.CreateMultipleFilter("responsible_user_id", []string{"234", "345"})},
			want:   url.Values{"filter[responsible_user_id][0]": []string{"234", "345"}},
		},
		{
			name:   "Интервал CreatedAt",
			filter: &GetCompaniesRequestFilter{CreatedAt: request.CreateIntervalFilter("created_at", "12345678", "23456789")},
			want:   url.Values{"filter[created_at][from]": []string{"12345678"}, "filter[created_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал UpdatedAt",
			filter: &GetCompaniesRequestFilter{UpdatedAt: request.CreateIntervalFilter("updated_at", "12345678", "23456789")},
			want:   url.Values{"filter[updated_at][from]": []string{"12345678"}, "filter[updated_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал ClosestTaskAt",
			filter: &GetCompaniesRequestFilter{ClosestTaskAt: request.CreateIntervalFilter("closest_task_at", "12345678", "23456789")},
			want:   url.Values{"filter[closest_task_at][from]": []string{"12345678"}, "filter[closest_task_at][to]": []string{"23456789"}},
		},
		{
			name:   "Простой CustomField",
			filter: &GetCompaniesRequestFilter{CustomFieldValues: []*request.Filter{request.CreateSimpleCustomFieldFilter("123", "custom_field_value")}},
			want:   url.Values{"filter[custom_fields_values][123][]": []string{"custom_field_value"}},
		},
		{
			name:   "Диапазон CustomField",
			filter: &GetCompaniesRequestFilter{CustomFieldValues: []*request.Filter{request.CreateIntervalCustomFieldFilter("123", "12345678", "23456789")}},
			want:   url.Values{"filter[custom_fields_values][123][from]": []string{"12345678"}, "filter[custom_fields_values][123][to]": []string{"23456789"}},
		},
		{
			name: "Несколько CustomField",
			filter: &GetCompaniesRequestFilter{CustomFieldValues: []*request.Filter{
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

func TestGetCompaniesRequestParamsValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров with в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{With: []GetCompaniesRequestWith{}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Один параметр with в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{With: []GetCompaniesRequestWith{LeadsGetCompaniesRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Несколько параметров with в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{With: []GetCompaniesRequestWith{LeadsGetCompaniesRequestWith, ContactsGetCompaniesRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Невалидный параметр with в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{With: []GetCompaniesRequestWith{"with"}}
		assert.EqualError(t, v.Struct(req), "Key: 'GetCompaniesRequestParams.With[0]' Error:Field validation for 'With[0]' failed on the 'oneof' tag")
	})

	t.Run("Превышен лимит элементов в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Limit: 260}
		assert.EqualError(t, v.Struct(req), "Key: 'GetCompaniesRequestParams.Limit' Error:Field validation for 'Limit' failed on the 'lte' tag")
	})

	t.Run("Диапазонный фильтр по ID в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{ID: request.CreateIntervalFilter("id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ID filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по Name в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{Name: request.CreateIntervalFilter("name", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "Name filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по CreatedBy в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{CreatedBy: request.CreateIntervalFilter("created_by", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "CreatedBy filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по UpdatedBy в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{UpdatedBy: request.CreateIntervalFilter("updated_by", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedBy filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по ResponsibleUserID в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{ResponsibleUserID: request.CreateIntervalFilter("responsible_user_id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ResponsibleUserID filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по CreatedAt в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{CreatedAt: request.CreateMultipleFilter("created_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "CreatedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по UpdatedAt в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{UpdatedAt: request.CreateMultipleFilter("updated_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по ClosestTaskAt в запросе", func(t *testing.T) {
		req := &GetCompaniesRequestParams{Filter: &GetCompaniesRequestFilter{ClosestTaskAt: request.CreateMultipleFilter("closest_task_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "ClosestTaskAt filter must be simple or interval type")
	})
}

func TestAddCompaniesRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров CustomFieldsValues в запросе", func(t *testing.T) {
		data := &AddCompaniesRequestData{CustomFieldsValues: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(data), "Key: 'AddCompaniesRequestData.CustomFieldsValues' Error:Field validation for 'CustomFieldsValues' failed on the 'gt' tag")
	})
}

func TestUpdateCompaniesRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Нет ID сделки в запросе", func(t *testing.T) {
		data := &UpdateCompaniesRequestData{}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateCompaniesRequestData.ID' Error:Field validation for 'ID' failed on the 'required' tag")
	})

	t.Run("Пустой массив параметров CustomFieldsValues в запросе", func(t *testing.T) {
		data := &UpdateCompaniesRequestData{ID: 1, CustomFieldsValues: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateCompaniesRequestData.CustomFieldsValues' Error:Field validation for 'CustomFieldsValues' failed on the 'gt' tag")
	})
}

func TestGetCompaniesResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &GetCompaniesResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'GetCompaniesResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag
Key: 'GetCompaniesResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Leads в ответе", func(t *testing.T) {
		data := &GetCompaniesResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetCompaniesResponseEmbedded{Companies: []*domain.Company{}}}
		assert.NoError(t, v.Struct(data))
	})

	t.Run("Ни одного обязательного параметра Lead в ответе", func(t *testing.T) {
		data := &GetCompaniesResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetCompaniesResponseEmbedded{Companies: []*domain.Company{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'GetCompaniesResponse.Embedded.Companies[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].ResponsibleUserID' Error:Field validation for 'ResponsibleUserID' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].GroupID' Error:Field validation for 'GroupID' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].UpdatedBy' Error:Field validation for 'UpdatedBy' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'GetCompaniesResponse.Embedded.Companies[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetCompanies(t *testing.T) {
	const sampleGetCompaniesResponseBody = `{"_page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies?limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/companies?limit=2&page=2"}},"_embedded":{"companies":[{"id":7767077,"name":"Компания Васи","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1586359618,"updated_at":1586359618,"closest_task_at":null,"custom_fields_values":null,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/7767077"}},"_embedded":{"tags":[]}},{"id":7767457,"name":"Example","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1586360394,"updated_at":1586360394,"closest_task_at":null,"custom_fields_values":null,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/7767457"}},"_embedded":{"tags":[]}}]}}`

	requestParamsWant := url.Values{
		"with":              []string{"leads,contacts"},
		"page":              []string{"1"},
		"limit":             []string{"2"},
		"filter[id][0]":     []string{"7767077", "7767457"},
		"order[updated_at]": []string{"asc"},
		"query":             []string{"query_value"},
	}

	sampleGetCompaniesRequestParams := &GetCompaniesRequestParams{
		With:  []GetCompaniesRequestWith{LeadsGetCompaniesRequestWith, ContactsGetCompaniesRequestWith},
		Page:  1,
		Limit: 2,
		Query: "query_value",
		Filter: &GetCompaniesRequestFilter{
			ID: request.CreateMultipleFilter("id", []string{"7767077", "7767457"}),
		},
		Order: &GetCompaniesOrder{
			By:     UpdatedAtGetCompaniesOrderBy,
			Method: request.AscendingOrderMethod,
		},
	}

	responseWant := []*domain.Company{
		{
			ID:                7767077,
			Name:              "Компания Васи",
			ResponsibleUserID: 504141,
			GroupID:           1,
			CreatedBy:         504141,
			UpdatedBy:         504141,
			CreatedAt:         1586359618,
			UpdatedAt:         1586359618,
			AccountID:         28805383,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/7767077"}},
			Embedded: &domain.CompanyEmbedded{
				Tags: []*domain.Tag{},
			},
		},
		{
			ID:                7767457,
			Name:              "Example",
			ResponsibleUserID: 504141,
			GroupID:           1,
			CreatedBy:         504141,
			UpdatedBy:         504141,
			CreatedAt:         1586360394,
			UpdatedAt:         1586360394,
			AccountID:         28805383,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/7767457"}},
			Embedded: &domain.CompanyEmbedded{
				Tags: []*domain.Tag{},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetCompaniesResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetCompanies(ctx, sampleGetCompaniesRequestParams)
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

		responseGot, err := client.GetCompanies(ctx, sampleGetCompaniesRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив в ответе", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies?limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/companies?limit=2&page=2"}},"_embedded":{"companies":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetCompanies(ctx, sampleGetCompaniesRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies?limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/companies?limit=2&page=2"}},"_embedded":{"companies":[{"id":7767077,"name":"Компания Васи","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1586359618,"updated_at":1586359618,"closest_task_at":null,"custom_fields_values":null,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/7767077"}},"_embedded":{"tags":[]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetCompanies(ctx, sampleGetCompaniesRequestParams)
		assert.EqualError(t, err, "Key: 'GetCompaniesResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestGetCompanyByID(t *testing.T) {
	const sampleGetCompanyByIDResponseBody = `{"id":1,"name":"АО Рога и копыта","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1582117331,"updated_at":1586361223,"closest_task_at":null,"custom_fields_values":[{"field_id":3,"field_name":"Телефон","field_code":"PHONE","field_type":"textarea","values":[{"value":"123213","enum_id":1,"enum":"WORK"}]}],"account_id":28805383,"_links":{"self":{"href":"https://exmaple.amocrm.ru/api/v4/companies/1"}},"_embedded":{"tags":[]}}`

	requestParamsWant := url.Values{
		"with": []string{"leads,catalog_elements"},
	}

	sampleGetCompaniesRequestParams := &GetCompaniesRequestParams{
		With: []GetCompaniesRequestWith{LeadsGetCompaniesRequestWith, CatalogElemetsGetCompaniesRequestWith},
	}

	responseWant := &domain.Company{
		ID:                1,
		Name:              "АО Рога и копыта",
		ResponsibleUserID: 504141,
		GroupID:           1,
		CreatedBy:         504141,
		UpdatedBy:         504141,
		CreatedAt:         1582117331,
		UpdatedAt:         1586361223,
		CustomFieldsValues: []*domain.CustomField{
			{ID: 3, Name: "Телефон", Code: "PHONE", Type: domain.TextAreaCustomFieldType, Values: []*domain.CustomFieldValue{{Value: "123213", EnumID: 1, Enum: "WORK"}}},
		},
		AccountID: 28805383,
		Links:     &domain.Links{Self: &domain.Link{Href: "https://exmaple.amocrm.ru/api/v4/companies/1"}},
		Embedded: &domain.CompanyEmbedded{
			Tags: []*domain.Tag{},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetCompanyByIDResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetCompanyByID(ctx, 1, sampleGetCompaniesRequestParams.With)
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

		responseGot, err := client.GetCompanyByID(ctx, 1, sampleGetCompaniesRequestParams.With)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		responseGot, err := client.GetCompanyByID(ctx, 0, sampleGetCompaniesRequestParams.With)
		assert.EqualError(t, err, ErrInvalidCompanyID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":1,"name":"АО Рога и копыта","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1582117331,"updated_at":1586361223,"closest_task_at":null,"custom_fields_values":[{"field_id":3,"field_name":"Телефон","field_code":"PHONE","field_type":"textarea","values":[{"value":"123213","enum_id":1,"enum":"WORK"}]}],"account_id":28805383,"_links":{"self":{"href":"https://exmaple.amocrm.ru/api/v4/companies/1"}},"_embedded":{"tags":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetCompanyByID(ctx, 1, sampleGetCompaniesRequestParams.With)
		assert.EqualError(t, err, "Key: 'Company.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestAddCompanies(t *testing.T) {
	const (
		requestBodyWant                = `[{"name":"АО Рога и Копыта","custom_fields_values":[{"field_code":"PHONE","values":[{"enum_code":"WORK","value":"+7912322222"}]}]}]`
		sampleAddCompaniesResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies"}},"_embedded":{"companies":[{"id":11090825,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/11090825"}}}]}}`
	)

	sampleAddCompaniesRequest := &AddCompaniesRequest{
		Add: []*AddCompaniesRequestData{
			{
				Name: "АО Рога и Копыта",
				CustomFieldsValues: []*domain.UpdateCustomField{
					{Code: "PHONE", Values: []interface{}{map[string]string{"value": "+7912322222", "enum_code": "WORK"}}},
				},
			},
		},
	}

	responseWant := []*AddCompaniesResponseItem{
		{
			ID:        11090825,
			RequestID: "0",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/11090825"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAddCompaniesResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddCompanies(ctx, sampleAddCompaniesRequest)
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

		responseGot, err := client.AddCompanies(ctx, sampleAddCompaniesRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies"}},"_embedded":{}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddCompanies(ctx, sampleAddCompaniesRequest)
		assert.EqualError(t, err, "Key: 'AddCompaniesResponse.Embedded.Companies' Error:Field validation for 'Companies' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateCompanies(t *testing.T) {
	const (
		requestBodyWant                   = `[{"id":11090825,"name":"Новое название компании","custom_fields_values":[{"field_code":"EMAIL","values":[{"enum_code":"WORK","value":"test@example.com"}]}]}]`
		sampleUpdateCompaniesResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies"}},"_embedded":{"companies":[{"id":11090825,"name":"Новое название компании","updated_at":1590998669,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/11090825"}}}]}}`
	)

	sampleUpdateCompaniesRequest := &UpdateCompaniesRequest{
		Update: []*UpdateCompaniesRequestData{
			{
				ID:   11090825,
				Name: "Новое название компании",
				CustomFieldsValues: []*domain.UpdateCustomField{
					{Code: "EMAIL", Values: []interface{}{map[string]string{"value": "test@example.com", "enum_code": "WORK"}}},
				},
			},
		},
	}

	responseWant := []*UpdateCompaniesResponseItem{
		{
			ID:        11090825,
			Name:      "Новое название компании",
			UpdatedAt: 1590998669,
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/11090825"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateCompaniesResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateCompanies(ctx, sampleUpdateCompaniesRequest)
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

		responseGot, err := client.UpdateCompanies(ctx, sampleUpdateCompaniesRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies"}},"_embedded":{}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateCompanies(ctx, sampleUpdateCompaniesRequest)
		assert.EqualError(t, err, "Key: 'UpdateCompaniesResponse.Embedded.Companies' Error:Field validation for 'Companies' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateCompany(t *testing.T) {
	const (
		requestBodyWant                 = `{"id":11090825,"name":"Новое название компании","custom_fields_values":[{"field_code":"EMAIL","values":[{"enum_code":"WORK","value":"test@example.com"}]}]}`
		sampleUpdateCompanyResponseBody = `{"id":11090825,"name":"Новое название компании","updated_at":1590998669,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/11090825"}}}`
	)

	sampleUpdateCompanyRequest := &UpdateCompaniesRequestData{
		ID:   11090825,
		Name: "Новое название компании",
		CustomFieldsValues: []*domain.UpdateCustomField{
			{Code: "EMAIL", Values: []interface{}{map[string]string{"value": "test@example.com", "enum_code": "WORK"}}},
		},
	}

	responseWant := &UpdateCompaniesResponseItem{
		ID:        11090825,
		Name:      "Новое название компании",
		UpdatedAt: 1590998669,
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/11090825"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateCompanyResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateCompany(ctx, 11090825, sampleUpdateCompanyRequest)
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

		responseGot, err := client.UpdateCompany(ctx, 11090825, sampleUpdateCompanyRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		responseGot, err := client.UpdateCompany(ctx, 0, sampleUpdateCompanyRequest)
		assert.EqualError(t, err, ErrInvalidCompanyID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":11090825,"name":"Новое название компании","updated_at":1590998669,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/11090825"}}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateCompany(ctx, 11090825, sampleUpdateCompanyRequest)
		assert.EqualError(t, err, "Key: 'UpdateCompaniesResponseItem.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}
