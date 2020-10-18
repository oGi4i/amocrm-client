package client

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestJoinGetContactsRequestWithSlice(t *testing.T) {
	testCases := []struct {
		name   string
		params []GetContactsRequestWith
		want   string
	}{
		{
			name:   "Ни одного параметра",
			params: []GetContactsRequestWith{},
			want:   "",
		},
		{
			name:   "Один параметр",
			params: []GetContactsRequestWith{LeadsGetContactsRequestWith},
			want:   "leads",
		},
		{
			name:   "Два параметра",
			params: []GetContactsRequestWith{LeadsGetContactsRequestWith, CustomersGetContactsRequestWith},
			want:   "leads,customers",
		},
		{
			name: "Все параметры",
			params: []GetContactsRequestWith{
				CatalogElementsGetContactsRequestWith,
				LeadsGetContactsRequestWith,
				CustomersGetContactsRequestWith,
			},
			want: "catalog_elements,leads,customers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, joinGetContactsRequestWith(tc.params))
		})
	}
}

func TestGetContactsRequestFilter(t *testing.T) {
	testCases := []struct {
		name   string
		filter *GetContactsRequestFilter
		want   url.Values
	}{
		{
			name:   "Один ID",
			filter: &GetContactsRequestFilter{ID: request.CreateSimpleFilter("id", "1")},
			want:   url.Values{"filter[id]": []string{"1"}},
		},
		{
			name:   "Несколько ID",
			filter: &GetContactsRequestFilter{ID: request.CreateMultipleFilter("id", []string{"1", "2", "3"})},
			want:   url.Values{"filter[id][0]": []string{"1", "2", "3"}},
		},
		{
			name:   "Name",
			filter: &GetContactsRequestFilter{Name: request.CreateSimpleFilter("name", "name_value")},
			want:   url.Values{"filter[name]": []string{"name_value"}},
		},
		{
			name:   "Несколько Name",
			filter: &GetContactsRequestFilter{Name: request.CreateMultipleFilter("name", []string{"name_value_1", "name_value_2"})},
			want:   url.Values{"filter[name][0]": []string{"name_value_1", "name_value_2"}},
		},
		{
			name:   "CreatedBy",
			filter: &GetContactsRequestFilter{CreatedBy: request.CreateSimpleFilter("created_by", "123")},
			want:   url.Values{"filter[created_by]": []string{"123"}},
		},
		{
			name:   "Несколько CreatedBy",
			filter: &GetContactsRequestFilter{CreatedBy: request.CreateMultipleFilter("created_by", []string{"234", "345"})},
			want:   url.Values{"filter[created_by][0]": []string{"234", "345"}},
		},
		{
			name:   "UpdatedBy",
			filter: &GetContactsRequestFilter{UpdatedBy: request.CreateSimpleFilter("updated_by", "123")},
			want:   url.Values{"filter[updated_by]": []string{"123"}},
		},
		{
			name:   "Несколько UpdatedBy",
			filter: &GetContactsRequestFilter{UpdatedBy: request.CreateMultipleFilter("updated_by", []string{"234", "345"})},
			want:   url.Values{"filter[updated_by][0]": []string{"234", "345"}},
		},
		{
			name:   "ResponsibleUserID",
			filter: &GetContactsRequestFilter{ResponsibleUserID: request.CreateSimpleFilter("responsible_user_id", "123")},
			want:   url.Values{"filter[responsible_user_id]": []string{"123"}},
		},
		{
			name:   "Несколько ResponsibleUserID",
			filter: &GetContactsRequestFilter{ResponsibleUserID: request.CreateMultipleFilter("responsible_user_id", []string{"234", "345"})},
			want:   url.Values{"filter[responsible_user_id][0]": []string{"234", "345"}},
		},
		{
			name:   "Интервал CreatedAt",
			filter: &GetContactsRequestFilter{CreatedAt: request.CreateIntervalFilter("created_at", "12345678", "23456789")},
			want:   url.Values{"filter[created_at][from]": []string{"12345678"}, "filter[created_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал UpdatedAt",
			filter: &GetContactsRequestFilter{UpdatedAt: request.CreateIntervalFilter("updated_at", "12345678", "23456789")},
			want:   url.Values{"filter[updated_at][from]": []string{"12345678"}, "filter[updated_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал ClosestTaskAt",
			filter: &GetContactsRequestFilter{ClosestTaskAt: request.CreateIntervalFilter("closest_task_at", "12345678", "23456789")},
			want:   url.Values{"filter[closest_task_at][from]": []string{"12345678"}, "filter[closest_task_at][to]": []string{"23456789"}},
		},
		{
			name:   "Простой CustomField",
			filter: &GetContactsRequestFilter{CustomFieldValues: []*request.Filter{request.CreateSimpleCustomFieldFilter("123", "custom_field_value")}},
			want:   url.Values{"filter[custom_fields_values][123][]": []string{"custom_field_value"}},
		},
		{
			name:   "Диапазон CustomField",
			filter: &GetContactsRequestFilter{CustomFieldValues: []*request.Filter{request.CreateIntervalCustomFieldFilter("123", "12345678", "23456789")}},
			want:   url.Values{"filter[custom_fields_values][123][from]": []string{"12345678"}, "filter[custom_fields_values][123][to]": []string{"23456789"}},
		},
		{
			name: "Несколько CustomField",
			filter: &GetContactsRequestFilter{CustomFieldValues: []*request.Filter{
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

func TestGetContactsRequestParamsValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров with в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{With: []GetContactsRequestWith{}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Один параметр with в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{With: []GetContactsRequestWith{LeadsGetContactsRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Несколько параметров with в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{With: []GetContactsRequestWith{LeadsGetContactsRequestWith, CustomersGetContactsRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Невалидный параметр with в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{With: []GetContactsRequestWith{"with"}}
		assert.EqualError(t, v.Struct(req), "Key: 'GetContactsRequestParams.With[0]' Error:Field validation for 'With[0]' failed on the 'oneof' tag")
	})

	t.Run("Превышен лимит элементов в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Limit: 260}
		assert.EqualError(t, v.Struct(req), "Key: 'GetContactsRequestParams.Limit' Error:Field validation for 'Limit' failed on the 'lte' tag")
	})

	t.Run("Диапазонный фильтр по ID в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{ID: request.CreateIntervalFilter("id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ID filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по Name в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{Name: request.CreateIntervalFilter("name", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "Name filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по CreatedBy в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{CreatedBy: request.CreateIntervalFilter("created_by", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "CreatedBy filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по UpdatedBy в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{UpdatedBy: request.CreateIntervalFilter("updated_by", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedBy filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по ResponsibleUserID в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{ResponsibleUserID: request.CreateIntervalFilter("responsible_user_id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ResponsibleUserID filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по CreatedAt в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{CreatedAt: request.CreateMultipleFilter("created_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "CreatedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по UpdatedAt в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{UpdatedAt: request.CreateMultipleFilter("updated_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedAt filter must be simple or interval type")
	})

	t.Run("Множественный фильтр по ClosestTaskAt в запросе", func(t *testing.T) {
		req := &GetContactsRequestParams{Filter: &GetContactsRequestFilter{ClosestTaskAt: request.CreateMultipleFilter("closest_task_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "ClosestTaskAt filter must be simple or interval type")
	})
}

func TestAddContactRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров CustomFieldsValues в запросе", func(t *testing.T) {
		data := &AddContactRequestData{CustomFieldsValues: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(data), "Key: 'AddContactRequestData.CustomFieldsValues' Error:Field validation for 'CustomFieldsValues' failed on the 'gt' tag")
	})
}

func TestUpdateContactRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Нет ID контакта в запросе", func(t *testing.T) {
		data := &UpdateContactsRequestData{}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateContactsRequestData.ID' Error:Field validation for 'ID' failed on the 'required' tag")
	})

	t.Run("Пустой массив параметров CustomFieldsValues в запросе", func(t *testing.T) {
		data := &UpdateContactsRequestData{ID: 1, CustomFieldsValues: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateContactsRequestData.CustomFieldsValues' Error:Field validation for 'CustomFieldsValues' failed on the 'gt' tag")
	})
}

func TestGetContactsResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &GetContactsResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'GetContactsResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag
Key: 'GetContactsResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Contacts в ответе", func(t *testing.T) {
		data := &GetContactsResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetContactsResponseEmbedded{Contacts: []*domain.Contact{}}}
		assert.NoError(t, v.Struct(data))
	})

	t.Run("Ни одного обязательного параметра в Contact в ответе", func(t *testing.T) {
		data := &GetContactsResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetContactsResponseEmbedded{Contacts: []*domain.Contact{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'GetContactsResponse.Embedded.Contacts[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].ResponsibleUserID' Error:Field validation for 'ResponsibleUserID' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].GroupID' Error:Field validation for 'GroupID' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].UpdatedBy' Error:Field validation for 'UpdatedBy' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'GetContactsResponse.Embedded.Contacts[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestAddContactsResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &AddContactsResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'AddContactsResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'AddContactsResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Contacts в ответе", func(t *testing.T) {
		data := &AddContactsResponse{Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &AddContactsResponseEmbedded{Contacts: []*AddContactsResponseItem{}}}
		assert.EqualError(t, v.Struct(data), "Key: 'AddContactsResponse.Embedded.Contacts' Error:Field validation for 'Contacts' failed on the 'gt' tag")
	})

	t.Run("Ни одного обязательного параметра в Contact в ответе", func(t *testing.T) {
		data := &AddContactsResponse{Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &AddContactsResponseEmbedded{Contacts: []*AddContactsResponseItem{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'AddContactsResponse.Embedded.Contacts[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'AddContactsResponse.Embedded.Contacts[0].RequestID' Error:Field validation for 'RequestID' failed on the 'required' tag
Key: 'AddContactsResponse.Embedded.Contacts[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestUpdateContactsResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &UpdateContactsResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'UpdateContactsResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'UpdateContactsResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Contacts в ответе", func(t *testing.T) {
		data := &UpdateContactsResponse{Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &UpdateContactsResponseEmbedded{Contacts: []*UpdateContactsResponseItem{}}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateContactsResponse.Embedded.Contacts' Error:Field validation for 'Contacts' failed on the 'gt' tag")
	})

	t.Run("Ни одного обязательного параметра в Contact в ответе", func(t *testing.T) {
		data := &UpdateContactsResponse{Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &UpdateContactsResponseEmbedded{Contacts: []*UpdateContactsResponseItem{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'UpdateContactsResponse.Embedded.Contacts[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'UpdateContactsResponse.Embedded.Contacts[0].Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'UpdateContactsResponse.Embedded.Contacts[0].UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'UpdateContactsResponse.Embedded.Contacts[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetContacts(t *testing.T) {
	const sampleGetContactsResponseBody = `{"_page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts?limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/contacts?limit=2&page=2"}},"_embedded":{"contacts":[{"id":7143599,"name":"1","first_name":"","last_name":"","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1585758065,"updated_at":1585758065,"closest_task_at":null,"custom_fields_values":null,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/7143599"}},"_embedded":{"tags":[],"companies":[]}},{"id":7767065,"name":"dsgdsg","first_name":"","last_name":"","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1586359590,"updated_at":1586359590,"closest_task_at":null,"custom_fields_values":null,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/7767065"}},"_embedded":{"tags":[],"companies":[]}}]}}`

	requestParamsWant := url.Values{
		"with":          []string{"leads,customers"},
		"page":          []string{"1"},
		"limit":         []string{"2"},
		"filter[id][0]": []string{"7143599", "7767065"},
		"order[id]":     []string{"asc"},
		"query":         []string{"query_value"},
	}

	sampleGetContactsRequestParams := &GetContactsRequestParams{
		With:  []GetContactsRequestWith{LeadsGetContactsRequestWith, CustomersGetContactsRequestWith},
		Page:  1,
		Limit: 2,
		Query: "query_value",
		Filter: &GetContactsRequestFilter{
			ID: request.CreateMultipleFilter("id", []string{"7143599", "7767065"}),
		},
		Order: &GetContactsOrder{
			By:     IDGetContactsOrderBy,
			Method: request.AscendingOrderMethod,
		},
	}

	responseWant := []*domain.Contact{
		{
			ID:                7143599,
			Name:              "1",
			FirstName:         "",
			LastName:          "",
			ResponsibleUserID: 504141,
			GroupID:           1,
			CreatedBy:         504141,
			UpdatedBy:         504141,
			CreatedAt:         1585758065,
			UpdatedAt:         1585758065,
			AccountID:         28805383,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/7143599"}},
			Embedded: &domain.ContactEmbedded{
				Tags:      []*domain.Tag{},
				Companies: []*domain.EmbeddedCompany{},
			},
		},
		{
			ID:                7767065,
			Name:              "dsgdsg",
			FirstName:         "",
			LastName:          "",
			ResponsibleUserID: 504141,
			GroupID:           1,
			CreatedBy:         504141,
			UpdatedBy:         504141,
			CreatedAt:         1586359590,
			UpdatedAt:         1586359590,
			AccountID:         28805383,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/7767065"}},
			Embedded: &domain.ContactEmbedded{
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
			_, _ = io.WriteString(w, sampleGetContactsResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetContacts(ctx, sampleGetContactsRequestParams)
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

		responseGot, err := client.GetContacts(ctx, sampleGetContactsRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив в ответе", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts?limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/contacts?limit=2&page=2"}},"_embedded":{"contacts":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetContacts(ctx, sampleGetContactsRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts?limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/contacts?limit=2&page=2"}},"_embedded":{"contacts":[{"id":7143599,"name":"1","first_name":"","last_name":"","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1585758065,"updated_at":1585758065,"closest_task_at":null,"custom_fields_values":null,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/7143599"}},"_embedded":{"tags":[],"companies":[]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetContacts(ctx, sampleGetContactsRequestParams)
		assert.EqualError(t, err, "Key: 'GetContactsResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestGetContactByID(t *testing.T) {
	const sampleGetContactByIDResponseBody = `{"id":3,"name":"Иван Иванов","first_name":"Иван","last_name":"Иванов","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1582117331,"updated_at":1590943929,"closest_task_at":null,"custom_fields_values":[{"field_id":3,"field_name":"Телефон","field_code":"PHONE","field_type":"text","values":[{"value":"+79123","enum_id":1,"enum":"WORK"}]}],"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/3"}},"_embedded":{"tags":[],"leads":[{"id":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/1"}}},{"id":3916883,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/3916883"}}}],"customers":[{"id":134923,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/customers/134923"}}}],"catalog_elements":[],"companies":[{"id":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/companies/1"}}}]}}`

	requestParamsWant := url.Values{
		"with": []string{"leads,customers,catalog_elements"},
	}

	sampleGetContactByIDRequestParams := []GetContactsRequestWith{LeadsGetContactsRequestWith, CustomersGetContactsRequestWith, CatalogElementsGetContactsRequestWith}

	responseWant := &domain.Contact{
		ID:                3,
		Name:              "Иван Иванов",
		FirstName:         "Иван",
		LastName:          "Иванов",
		ResponsibleUserID: 504141,
		GroupID:           1,
		CreatedBy:         504141,
		UpdatedBy:         504141,
		CreatedAt:         1582117331,
		UpdatedAt:         1590943929,
		CustomFieldsValues: []*domain.CustomField{
			{ID: 3, Name: "Телефон", Code: "PHONE", Type: domain.TextCustomFieldType, Values: []*domain.CustomFieldValue{{Value: "+79123", EnumID: 1, Enum: "WORK"}}},
		},
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/3"}},
		AccountID: 28805383,
		Embedded: &domain.ContactEmbedded{
			Tags: []*domain.Tag{},
			Leads: []*domain.EmbeddedLead{
				{ID: 1, Links: &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/1"}}},
				{ID: 3916883, Links: &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/3916883"}}},
			},
			Customers: []*domain.EmbeddedCustomer{
				{ID: 134923, Links: &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/customers/134923"}}},
			},
			CatalogElements: []*domain.CatalogElement{},
			Companies: []*domain.EmbeddedCompany{
				{ID: 1, Links: &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/companies/1"}}},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetContactByIDResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetContactByID(ctx, 3, sampleGetContactByIDRequestParams)
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

		responseGot, err := client.GetContactByID(ctx, 3, sampleGetContactByIDRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		responseGot, err := client.GetContactByID(ctx, 0, sampleGetContactByIDRequestParams)
		assert.EqualError(t, err, ErrInvalidContactID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":3,"name":"Иван Иванов","first_name":"Иван","last_name":"Иванов","responsible_user_id":504141,"group_id":1,"created_by":504141,"updated_by":504141,"created_at":1582117331,"updated_at":1590943929,"account_id":28805383,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/3"}},"_embedded":{"tags":[],"leads":[],"customers":[],"catalog_elements":[],"companies":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetContactByID(ctx, 3, sampleGetContactByIDRequestParams)
		assert.EqualError(t, err, "Key: 'Contact.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestAddContacts(t *testing.T) {
	const (
		requestBodyWant               = `[{"first_name":"Петр","last_name":"Смирнов","custom_fields_values":[{"field_id":271316,"values":[{"value":"Директор"}]}]},{"name":"Владимир Смирнов","created_by":47272}]`
		sampleAddContactsResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts"}},"_embedded":{"contacts":[{"id":40401635,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/40401635"}}},{"id":40401636,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/40401636"}}}]}}`
	)

	sampleAddContactsRequest := &AddContactsRequest{
		Add: []*AddContactRequestData{
			{
				FirstName: "Петр",
				LastName:  "Смирнов",
				CustomFieldsValues: []*domain.UpdateCustomField{
					{ID: 271316, Values: []interface{}{map[string]string{"value": "Директор"}}},
				},
			},
			{
				Name:      "Владимир Смирнов",
				CreatedBy: 47272,
			},
		},
	}

	responseWant := []*AddContactsResponseItem{
		{
			ID:        40401635,
			RequestID: "0",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/40401635"}},
		},
		{
			ID:        40401636,
			RequestID: "1",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/40401636"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAddContactsResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddContacts(ctx, sampleAddContactsRequest)
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

		responseGot, err := client.AddContacts(ctx, sampleAddContactsRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts"}},"_embedded":{"contacts":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddContacts(ctx, sampleAddContactsRequest)
		assert.EqualError(t, err, "Key: 'AddContactsResponse.Embedded.Contacts' Error:Field validation for 'Contacts' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateContacts(t *testing.T) {
	const (
		requestBodyWant                  = `[{"id":3,"first_name":"Иван","last_name":"Иванов","custom_fields_values":[{"field_id":66192,"field_name":"Телефон","values":[{"enum_code":"WORK","value":"79999999999"}]}]}]`
		sampleUpdateContactsResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts"}},"_embedded":{"contacts":[{"id":3,"name":"Иван Иванов","updated_at":1590945248,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/3"}}}]}}`
	)

	sampleUpdateContactsRequest := &UpdateContactsRequest{
		Update: []*UpdateContactsRequestData{
			{
				ID:        3,
				FirstName: "Иван",
				LastName:  "Иванов",
				CustomFieldsValues: []*domain.UpdateCustomField{
					{ID: 66192, Name: "Телефон", Values: []interface{}{map[string]string{"value": "79999999999", "enum_code": "WORK"}}},
				},
			},
		},
	}

	responseWant := []*UpdateContactsResponseItem{
		{
			ID:        3,
			Name:      "Иван Иванов",
			UpdatedAt: 1590945248,
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/3"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateContactsResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateContacts(ctx, sampleUpdateContactsRequest)
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

		responseGot, err := client.UpdateContacts(ctx, sampleUpdateContactsRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts"}},"_embedded":{"contacts":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateContacts(ctx, sampleUpdateContactsRequest)
		assert.EqualError(t, err, "Key: 'UpdateContactsResponse.Embedded.Contacts' Error:Field validation for 'Contacts' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateContact(t *testing.T) {
	const (
		requestBodyWant                 = `{"id":3,"first_name":"Иван","last_name":"Иванов","custom_fields_values":[{"field_id":66192,"field_name":"Телефон","values":[{"enum_code":"WORK","value":"79999999999"}]}]}`
		sampleUpdateContactResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts"}},"_embedded":{"contacts":[{"id":3,"name":"Иван Иванов","updated_at":1590945248,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts/3"}}}]}}`
	)

	sampleUpdateContactRequest := &UpdateContactsRequestData{
		ID:        3,
		FirstName: "Иван",
		LastName:  "Иванов",
		CustomFieldsValues: []*domain.UpdateCustomField{
			{ID: 66192, Name: "Телефон", Values: []interface{}{map[string]string{"value": "79999999999", "enum_code": "WORK"}}},
		},
	}

	responseWant := &UpdateContactsResponseItem{
		ID:        3,
		Name:      "Иван Иванов",
		UpdatedAt: 1590945248,
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/contacts/3"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateContactResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateContact(ctx, 3, sampleUpdateContactRequest)
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

		responseGot, err := client.UpdateContact(ctx, 3, sampleUpdateContactRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := defaultTestClient()
		assert.NoError(t, err)

		responseGot, err := client.UpdateContact(ctx, 0, sampleUpdateContactRequest)
		assert.EqualError(t, err, ErrInvalidContactID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/contacts"}},"_embedded":{"contacts":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdateContact(ctx, 3, sampleUpdateContactRequest)
		assert.EqualError(t, err, "Key: 'UpdateContactsResponse.Embedded.Contacts' Error:Field validation for 'Contacts' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}
