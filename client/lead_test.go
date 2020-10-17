package client

import (
	"context"
	"github.com/go-playground/validator/v10"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ogi4i/amocrm-client/domain"
	"github.com/ogi4i/amocrm-client/request"
)

func TestJoinLeadRequestWithSlice(t *testing.T) {
	testCases := []struct {
		name   string
		params []LeadGetRequestWith
		want   string
	}{
		{
			name:   "Ни одного параметра",
			params: []LeadGetRequestWith{},
			want:   "",
		},
		{
			name:   "Один параметр",
			params: []LeadGetRequestWith{SourceIDLeadRequestWith},
			want:   "source_id",
		},
		{
			name:   "Два параметра",
			params: []LeadGetRequestWith{SourceIDLeadRequestWith, LossReasonLeadRequestWith},
			want:   "source_id,loss_reason",
		},
		{
			name: "Все параметры",
			params: []LeadGetRequestWith{
				CatalogElementsLeadRequestWith,
				IsPriceModifiedByRobotLeadRequestWith,
				LossReasonLeadRequestWith,
				ContactsLeadRequestWith,
				OnlyDeletedLeadRequestWith,
				SourceIDLeadRequestWith,
			},
			want: "catalog_elements,is_price_modified_by_robot,loss_reason,contacts,only_deleted,source_id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, joinLeadRequestWithSlice(tc.params))
		})
	}
}

func TestAppendGetRequestFilter(t *testing.T) {
	testCases := []struct {
		name   string
		filter *LeadGetRequestFilter
		want   url.Values
	}{
		{
			name:   "Один ID",
			filter: &LeadGetRequestFilter{ID: request.CreateSimpleFilter("id", "1")},
			want:   url.Values{"filter[id]": []string{"1"}},
		},
		{
			name:   "Несколько ID",
			filter: &LeadGetRequestFilter{ID: request.CreateMultipleFilter("id", []string{"1", "2", "3"})},
			want:   url.Values{"filter[id][0]": []string{"1", "2", "3"}},
		},
		{
			name:   "Name",
			filter: &LeadGetRequestFilter{Name: request.CreateSimpleFilter("name", "name_value")},
			want:   url.Values{"filter[name]": []string{"name_value"}},
		},
		{
			name:   "Несколько Name",
			filter: &LeadGetRequestFilter{ID: request.CreateMultipleFilter("name", []string{"name_value_1", "name_value_2"})},
			want:   url.Values{"filter[name][0]": []string{"name_value_1", "name_value_2"}},
		},
		{
			name:   "Price",
			filter: &LeadGetRequestFilter{ID: request.CreateSimpleFilter("price", "100")},
			want:   url.Values{"filter[price]": []string{"100"}},
		},
		{
			name:   "Диапазон Price",
			filter: &LeadGetRequestFilter{ID: request.CreateIntervalFilter("price", "100", "200")},
			want:   url.Values{"filter[price][from]": []string{"100"}, "filter[price][to]": []string{"200"}},
		},
		{
			name:   "Statuses",
			filter: &LeadGetRequestFilter{ID: request.CreateStatusFilter("123", "234")},
			want:   url.Values{"filter[statuses][0][pipeline_id]": []string{"123"}, "filter[statuses][0][status_id]": []string{"234"}},
		},
		{
			name:   "PipelineID",
			filter: &LeadGetRequestFilter{ID: request.CreateSimpleFilter("pipeline_id", "123")},
			want:   url.Values{"filter[pipeline_id]": []string{"123"}},
		},
		{
			name:   "Несколько PipelineID",
			filter: &LeadGetRequestFilter{ID: request.CreateMultipleFilter("pipeline_id", []string{"123", "234", "345"})},
			want:   url.Values{"filter[pipeline_id][0]": []string{"123", "234", "345"}},
		},
		{
			name:   "CreatedBy",
			filter: &LeadGetRequestFilter{CreatedBy: request.CreateSimpleFilter("created_by", "123")},
			want:   url.Values{"filter[created_by]": []string{"123"}},
		},
		{
			name:   "Несколько CreatedBy",
			filter: &LeadGetRequestFilter{CreatedBy: request.CreateMultipleFilter("created_by", []string{"234", "345"})},
			want:   url.Values{"filter[created_by][0]": []string{"234", "345"}},
		},
		{
			name:   "UpdatedBy",
			filter: &LeadGetRequestFilter{UpdatedBy: request.CreateSimpleFilter("updated_by", "123")},
			want:   url.Values{"filter[updated_by]": []string{"123"}},
		},
		{
			name:   "Несколько UpdatedBy",
			filter: &LeadGetRequestFilter{UpdatedBy: request.CreateMultipleFilter("updated_by", []string{"234", "345"})},
			want:   url.Values{"filter[updated_by][0]": []string{"234", "345"}},
		},
		{
			name:   "ResponsibleUserID",
			filter: &LeadGetRequestFilter{ResponsibleUserID: request.CreateSimpleFilter("responsible_user_id", "123")},
			want:   url.Values{"filter[responsible_user_id]": []string{"123"}},
		},
		{
			name:   "Несколько ResponsibleUserID",
			filter: &LeadGetRequestFilter{ResponsibleUserID: request.CreateMultipleFilter("responsible_user_id", []string{"234", "345"})},
			want:   url.Values{"filter[responsible_user_id][0]": []string{"234", "345"}},
		},
		{
			name:   "Интервал CreatedAt",
			filter: &LeadGetRequestFilter{CreatedAt: request.CreateIntervalFilter("created_at", "12345678", "23456789")},
			want:   url.Values{"filter[created_at][from]": []string{"12345678"}, "filter[created_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал UpdatedAt",
			filter: &LeadGetRequestFilter{UpdatedAt: request.CreateIntervalFilter("updated_at", "12345678", "23456789")},
			want:   url.Values{"filter[updated_at][from]": []string{"12345678"}, "filter[updated_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал ClosedAt",
			filter: &LeadGetRequestFilter{ClosedAt: request.CreateIntervalFilter("closed_at", "12345678", "23456789")},
			want:   url.Values{"filter[closed_at][from]": []string{"12345678"}, "filter[closed_at][to]": []string{"23456789"}},
		},
		{
			name:   "Интервал ClosestTaskAt",
			filter: &LeadGetRequestFilter{ClosestTaskAt: request.CreateIntervalFilter("closest_task_at", "12345678", "23456789")},
			want:   url.Values{"filter[closest_task_at][from]": []string{"12345678"}, "filter[closest_task_at][to]": []string{"23456789"}},
		},
		{
			name:   "Простой CustomField",
			filter: &LeadGetRequestFilter{CustomFieldValues: []*request.Filter{request.CreateSimpleCustomFieldFilter("123", "custom_field_value")}},
			want:   url.Values{"filter[custom_fields_values][123][]": []string{"custom_field_value"}},
		},
		{
			name:   "Диапазон CustomField",
			filter: &LeadGetRequestFilter{CustomFieldValues: []*request.Filter{request.CreateIntervalCustomFieldFilter("123", "12345678", "23456789")}},
			want:   url.Values{"filter[custom_fields_values][123][from]": []string{"12345678"}, "filter[custom_fields_values][123][to]": []string{"23456789"}},
		},
		{
			name: "Несколько CustomField",
			filter: &LeadGetRequestFilter{CustomFieldValues: []*request.Filter{
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
			tc.filter.AppendGetRequestFilter(params)
			assert.Equal(t, tc.want, params)
		})
	}
}

func TestLeadGetRequestParamsValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров with в запросе", func(t *testing.T) {
		req := &LeadGetRequestParams{With: []LeadGetRequestWith{}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Один параметр with в запросе", func(t *testing.T) {
		req := &LeadGetRequestParams{With: []LeadGetRequestWith{LossReasonLeadRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Несколько параметров with в запросе", func(t *testing.T) {
		req := &LeadGetRequestParams{With: []LeadGetRequestWith{LossReasonLeadRequestWith, ContactsLeadRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Невалидный параметр with в запросе", func(t *testing.T) {
		req := &LeadGetRequestParams{With: []LeadGetRequestWith{"with"}}
		assert.EqualError(t, v.Struct(req), "Key: 'LeadGetRequestParams.With[0]' Error:Field validation for 'With[0]' failed on the 'oneof' tag")
	})

	t.Run("Превышен лимит элементов в запросе", func(t *testing.T) {
		req := &LeadGetRequestParams{Limit: 260}
		assert.EqualError(t, v.Struct(req), "Key: 'LeadGetRequestParams.Limit' Error:Field validation for 'Limit' failed on the 'lte' tag")
	})
}

func TestLeadUpdateDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров CustomFields в запросе", func(t *testing.T) {
		req := &LeadUpdateData{CustomFields: []*domain.UpdateCustomField{}}
		assert.EqualError(t, v.Struct(req), "Key: 'LeadUpdateData.CustomFields' Error:Field validation for 'CustomFields' failed on the 'gt' tag")
	})
}

func TestLeadGetResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		req := &LeadGetResponse{}
		assert.EqualError(t, v.Struct(req), `Key: 'LeadGetResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag
Key: 'LeadGetResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'LeadGetResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Leads в ответе", func(t *testing.T) {
		req := &LeadGetResponseEmbedded{Leads: []*domain.Lead{}}
		assert.EqualError(t, v.Struct(req), "Key: 'LeadGetResponseEmbedded.Leads' Error:Field validation for 'Leads' failed on the 'gt' tag")
	})

	t.Run("Ни одного обязательного параметра Lead в ответе", func(t *testing.T) {
		req := &domain.Lead{}
		assert.EqualError(t, v.Struct(req), `Key: 'Lead.ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'Lead.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'Lead.Price' Error:Field validation for 'Price' failed on the 'required' tag
Key: 'Lead.ResponsibleUserID' Error:Field validation for 'ResponsibleUserID' failed on the 'required' tag
Key: 'Lead.StatusID' Error:Field validation for 'StatusID' failed on the 'required' tag
Key: 'Lead.PipelineID' Error:Field validation for 'PipelineID' failed on the 'required' tag
Key: 'Lead.CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag
Key: 'Lead.UpdatedBy' Error:Field validation for 'UpdatedBy' failed on the 'required' tag
Key: 'Lead.CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag
Key: 'Lead.UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'Lead.AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'Lead.Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetLeads(t *testing.T) {
	const (
		sampleGetLeadsResponseBody = `{
    "_page": 2,
    "_links": {
        "self": {
            "href": "https://example.amocrm.ru/api/v4/leads?limit=2&page=2"
        },
        "next": {
            "href": "https://example.amocrm.ru/api/v4/leads?limit=2&page=3"
        },
        "first": {
            "href": "https://example.amocrm.ru/api/v4/leads?limit=2&page=1"
        },
        "prev": {
            "href": "https://example.amocrm.ru/api/v4/leads?limit=2&page=1"
        }
    },
    "_embedded": {
        "leads": [
            {
                "id": 19619,
                "name": "Сделка для примера",
                "price": 46333,
                "responsible_user_id": 123321,
                "group_id": 625,
                "status_id": 142,
                "pipeline_id": 1300,
                "loss_reason_id": null,
                "source_id": null,
                "created_by": 321123,
                "updated_by": 321123,
                "created_at": 1453279607,
                "updated_at": 1502193501,
                "closed_at": 1483005931,
                "closest_task_at": null,
                "is_deleted": false,
                "custom_fields_values": null,
                "score": null,
                "account_id": 5135160,
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/19619"
                    }
                },
                "_embedded": {
                    "tags": [],
                    "companies": []
                }
            },
            {
                "id": 14460,
                "name": "Сделка для примера 2",
                "price": 655,
                "responsible_user_id": 123321,
                "group_id": 625,
                "status_id": 142,
                "pipeline_id": 1300,
                "loss_reason_id": null,
                "source_id": null,
                "created_by": 321123,
                "updated_by": 321123,
                "created_at": 1453279607,
                "updated_at": 1502193501,
                "closed_at": 1483005931,
                "closest_task_at": null,
                "is_deleted": false,
                "custom_fields_values": null,
                "score": null,
                "account_id": 1351360,
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/14460"
                    }
                },
                "_embedded": {
                    "tags": [],
                    "companies": []
                }
            }
        ]
    }
}`
	)

	requestParamsWant := url.Values{
		"with":          []string{"contacts,loss_reason"},
		"page":          []string{"2"},
		"limit":         []string{"2"},
		"filter[id][0]": []string{"14460", "19619"},
		"order[id]":     []string{"desc"},
		"query":         []string{"query_value"},
	}

	sampleGetLeadsRequestParams := &LeadGetRequestParams{
		With:  []LeadGetRequestWith{ContactsLeadRequestWith, LossReasonLeadRequestWith},
		Page:  2,
		Limit: 2,
		Query: "query_value",
		Filter: &LeadGetRequestFilter{
			ID: request.CreateMultipleFilter("id", []string{"14460", "19619"}),
		},
		Order: &request.Order{
			By:     request.IDRequestOrderBy,
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
				Companies: []*domain.Company{},
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
				Companies: []*domain.Company{},
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

		client, err := NewClient(server.URL, "login", "hash")
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.EqualError(t, err, domain.ErrEmptyResponse.Error())
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.EqualError(t, err, "Key: 'LeadGetResponse.Embedded.Leads' Error:Field validation for 'Leads' failed on the 'gt' tag")
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetLeads(ctx, sampleGetLeadsRequestParams)
		assert.EqualError(t, err, "Key: 'LeadGetResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestGetLeadByID(t *testing.T) {
	const (
		sampleGetLeadByIDResponseBody = `{
    "id": 3912171,
    "name": "Example",
    "price": 12,
    "responsible_user_id": 504141,
    "group_id": 0,
    "status_id": 143,
    "pipeline_id": 3104455,
    "loss_reason_id": 4203748,
    "source_id": null,
    "created_by": 504141,
    "updated_by": 504141,
    "created_at": 1585299171,
    "updated_at": 1590683337,
    "closed_at": 1590683337,
    "closest_task_at": null,
    "is_deleted": false,
    "custom_fields_values": null,
    "score": null,
    "account_id": 28805383,
    "is_price_modified_by_robot": false,
    "_links": {
        "self": {
            "href": "https://example.amocrm.ru/api/v4/leads/3912171"
        }
    },
    "_embedded": {
        "tags": [
            {
                "id": 100667,
                "name": "тест"
            }
        ],
        "catalog_elements": [
            {
                "id": 525439,
                "metadata": {
                    "quantity": 1,
                    "catalog_id": 4521
                }
            }
        ],
        "loss_reason": [
            {
                "id": 4203748,
                "name": "Пропала потребность",
                "sort": 100000,
                "created_at": 1582117280,
                "updated_at": 1582117280,
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/loss_reasons/4203748"
                    }
                }
            }
        ],
        "companies": [
            {
                "id": 10971463,
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/companies/10971463"
                    }
                }
            }
        ],
        "contacts": []
    }
}`
	)

	requestParamsWant := url.Values{
		"with": []string{"loss_reason,catalog_elements"},
	}

	sampleGetLeadsRequestParams := &LeadGetRequestParams{
		With: []LeadGetRequestWith{LossReasonLeadRequestWith, CatalogElementsLeadRequestWith},
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
					Metadata: struct {
						Quantity  int64 `json:"quantity" validate:"required"`
						CatalogID int64 `json:"catalog_id" validate:"required"`
					}{
						Quantity:  1,
						CatalogID: 4521,
					},
				},
			},
			LossReasons: []*domain.LossReason{{ID: 4203748, Name: "Пропала потребность"}},
			Companies:   []*domain.Company{{ID: 10971463}},
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

		client, err := NewClient(server.URL, "login", "hash")
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetLeadByID(ctx, 3912171, sampleGetLeadsRequestParams.With)
		assert.EqualError(t, err, domain.ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"id":3912171,"name":"Example","price":12,"responsible_user_id":504141,"group_id":0,"status_id":143,"pipeline_id":3104455,"loss_reason_id":4203748,"source_id":null,"created_by":504141,"updated_by":504141,"created_at":1585299171,"updated_at":1590683337,"closed_at":1590683337,"account_id":28805383}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetLeadByID(ctx, 3912171, sampleGetLeadsRequestParams.With)
		assert.EqualError(t, err, "Key: 'Lead.Links' Error:Field validation for 'Links' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestAddLead(t *testing.T) {
	const (
		requestBodyWant           = `[{"name":"Сделка для примера 1","price":20000,"custom_fields":[{"field_id":294471,"values":[{"value":"Наш первый клиент"}]}]},{"name":"Сделка для примера 2","price":10000,"_embedded":{"tags":[{"id":2719,"name":""}]}}]`
		sampleAddLeadResponseBody = `{
    "_links": {
        "self": {
            "href": "https://example.amocrm.ru/api/v4/leads"
        }
    },
    "_embedded": {
        "leads": [
            {
                "id": 10185151,
                "request_id": "0",
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/10185151"
                    }
                }
            },
            {
                "id": 10185153,
                "request_id": "1",
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/10185153"
                    }
                }
            }
        ]
    }
}`
	)

	sampleAddLeadRequest := []*LeadUpdateData{
		{
			Name:      "Сделка для примера 1",
			Price:     20000,
			CreatedBy: 0,
			CustomFields: []*domain.UpdateCustomField{
				{
					FieldID: 294471,
					Values: []interface{}{
						map[string]string{
							"value": "Наш первый клиент",
						},
					},
				},
			},
		},
		{
			Name:  "Сделка для примера 2",
			Price: 10000,
			Embedded: &LeadUpdateEmbedded{
				Tags: []*domain.Tag{
					{
						ID: 2719,
					},
				},
			},
		},
	}

	responseWant := []*LeadUpdateResponseItem{
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
			_, _ = io.WriteString(w, sampleAddLeadResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddLead(ctx, sampleAddLeadRequest)
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddLead(ctx, sampleAddLeadRequest)
		assert.EqualError(t, err, domain.ErrEmptyResponse.Error())
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddLead(ctx, sampleAddLeadRequest)
		assert.EqualError(t, err, "Key: 'LeadUpdateResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateLead(t *testing.T) {
	const (
		requestBodyWant           = `[{"name":"Сделка для примера 1","price":20000,"custom_fields":[{"field_id":294471,"values":[{"value":"Наш первый клиент"}]}]},{"name":"Сделка для примера 2","price":10000,"_embedded":{"tags":[{"id":2719,"name":""}]}}]`
		sampleAddLeadResponseBody = `{
    "_links": {
        "self": {
            "href": "https://example.amocrm.ru/api/v4/leads"
        }
    },
    "_embedded": {
        "leads": [
            {
                "id": 10185151,
                "request_id": "0",
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/10185151"
                    }
                }
            },
            {
                "id": 10185153,
                "request_id": "1",
                "_links": {
                    "self": {
                        "href": "https://example.amocrm.ru/api/v4/leads/10185153"
                    }
                }
            }
        ]
    }
}`
	)

	sampleAddLeadRequest := []*LeadUpdateData{
		{
			Name:      "Сделка для примера 1",
			Price:     20000,
			CreatedBy: 0,
			CustomFields: []*domain.UpdateCustomField{
				{
					FieldID: 294471,
					Values: []interface{}{
						map[string]string{
							"value": "Наш первый клиент",
						},
					},
				},
			},
		},
		{
			Name:  "Сделка для примера 2",
			Price: 10000,
			Embedded: &LeadUpdateEmbedded{
				Tags: []*domain.Tag{
					{
						ID: 2719,
					},
				},
			},
		},
	}

	responseWant := []*LeadUpdateResponseItem{
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
			_, _ = io.WriteString(w, sampleAddLeadResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, sampleAddLeadRequest)
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, sampleAddLeadRequest)
		assert.EqualError(t, err, domain.ErrEmptyResponse.Error())
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

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateLead(ctx, sampleAddLeadRequest)
		assert.EqualError(t, err, "Key: 'LeadUpdateResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}
