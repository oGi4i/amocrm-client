package client

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/ogi4i/amocrm-client/domain"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPipelinesResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &PipelinesResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'PipelinesResponse.TotalItems' Error:Field validation for 'TotalItems' failed on the 'required' tag
Key: 'PipelinesResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag
Key: 'PipelinesResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Pipelines в ответе", func(t *testing.T) {
		data := &PipelinesResponse{TotalItems: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &PipelinesResponseEmbedded{Pipelines: []*domain.Pipeline{}}}
		assert.NoError(t, v.Struct(data))
	})

	t.Run("Ни одного обязательного параметра в Pipeline в ответе", func(t *testing.T) {
		data := &PipelinesResponse{TotalItems: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &PipelinesResponseEmbedded{Pipelines: []*domain.Pipeline{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'PipelinesResponse.Embedded.Pipelines[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'PipelinesResponse.Embedded.Pipelines[0].Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'PipelinesResponse.Embedded.Pipelines[0].Sort' Error:Field validation for 'Sort' failed on the 'required' tag
Key: 'PipelinesResponse.Embedded.Pipelines[0].AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'PipelinesResponse.Embedded.Pipelines[0].Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag
Key: 'PipelinesResponse.Embedded.Pipelines[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetPipelines(t *testing.T) {
	const sampleGetPipelinesResponseBody = `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ffc8c8","type":1,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}},{"id":32392159,"name":"Первичный контакт","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#ffdc7f","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}}},{"id":32392165,"name":"Принимают решение","sort":30,"is_editable":true,"pipeline_id":3177727,"color":"#ccc8f9","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}}},{"id":142,"name":"Успешно реализовано","sort":10000,"is_editable":false,"pipeline_id":3177727,"color":"#c1e0ff","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}}},{"id":143,"name":"Закрыто и не реализовано","sort":11000,"is_editable":false,"pipeline_id":3177727,"color":"#e6e8ea","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}}}]}}]}}`

	responseWant := []*domain.Pipeline{
		{
			ID:           3177727,
			Name:         "Воронка",
			Sort:         1,
			IsMain:       true,
			IsUnsortedOn: true,
			AccountID:    12345678,
			Links:        &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},
			Embedded: &domain.PipelineEmbedded{
				Statuses: []*domain.PipelineStatus{
					{
						ID:         32392156,
						Name:       "Неразобранное",
						Sort:       10,
						PipelineID: 3177727,
						Color:      domain.YourPinkPipelineStatusColor,
						Type:       domain.RegularPipelineStatusType,
						AccountID:  12345678,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}},
					},
					{
						ID:         32392159,
						Name:       "Первичный контакт",
						Sort:       20,
						IsEditable: true,
						PipelineID: 3177727,
						Color:      domain.SalomiePipelineStatusColor,
						AccountID:  12345678,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}},
					},
					{
						ID:         32392165,
						Name:       "Принимают решение",
						Sort:       30,
						IsEditable: true,
						PipelineID: 3177727,
						Color:      domain.LavanderBluePipelineStatusColor,
						AccountID:  12345678,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}},
					},
					{
						ID:         142,
						Name:       "Успешно реализовано",
						Sort:       10000,
						PipelineID: 3177727,
						Color:      domain.MediumPattensBluePipelineStatusColor,
						AccountID:  12345678,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}},
					},
					{
						ID:         143,
						Name:       "Закрыто и не реализовано",
						Sort:       11000,
						PipelineID: 3177727,
						Color:      domain.SolitudePipelineStatusColor,
						AccountID:  12345678,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}},
					},
				},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetPipelinesResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelines(ctx)
		assert.NoError(t, err)
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelines(ctx)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Pipelines в ответе", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelines(ctx)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Statuses в ответе", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelines(ctx)
		assert.EqualError(t, err, "Key: 'PipelinesResponse.Embedded.Pipelines[0].Embedded.Statuses' Error:Field validation for 'Statuses' failed on the 'gt' tag")
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ffc8c8","type":1,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}}]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatuses(ctx, 3177727)
		assert.EqualError(t, err, "Key: 'PipelinesResponse.TotalItems' Error:Field validation for 'TotalItems' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestGetPipelineByID(t *testing.T) {
	const sampleGetPipelineByIDResponseBody = `{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ffc8c8","type":1,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}},{"id":32392159,"name":"Первичный контакт","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#ffdc7f","type":0,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}}},{"id":32392165,"name":"Принимают решение","sort":30,"is_editable":true,"pipeline_id":3177727,"color":"#ccc8f9","type":0,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}}},{"id":142,"name":"Успешно реализовано","sort":10000,"is_editable":false,"pipeline_id":3177727,"color":"#c1e0ff","type":0,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}}},{"id":143,"name":"Закрыто и не реализовано","sort":11000,"is_editable":false,"pipeline_id":3177727,"color":"#e6e8ea","type":0,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}}}]}}`

	responseWant := &domain.Pipeline{
		ID:           3177727,
		Name:         "Воронка",
		Sort:         1,
		IsMain:       true,
		IsUnsortedOn: true,
		AccountID:    28847170,
		Links:        &domain.Links{Self: &domain.Link{Href: "https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727"}},
		Embedded: &domain.PipelineEmbedded{
			Statuses: []*domain.PipelineStatus{
				{
					ID:         32392156,
					Name:       "Неразобранное",
					Sort:       10,
					PipelineID: 3177727,
					Color:      domain.YourPinkPipelineStatusColor,
					Type:       domain.RegularPipelineStatusType,
					AccountID:  28847170,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}},
				},
				{
					ID:         32392159,
					Name:       "Первичный контакт",
					Sort:       20,
					IsEditable: true,
					PipelineID: 3177727,
					Color:      domain.SalomiePipelineStatusColor,
					AccountID:  28847170,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}},
				},
				{
					ID:         32392165,
					Name:       "Принимают решение",
					Sort:       30,
					IsEditable: true,
					PipelineID: 3177727,
					Color:      domain.LavanderBluePipelineStatusColor,
					AccountID:  28847170,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}},
				},
				{
					ID:         142,
					Name:       "Успешно реализовано",
					Sort:       10000,
					PipelineID: 3177727,
					Color:      domain.MediumPattensBluePipelineStatusColor,
					AccountID:  28847170,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}},
				},
				{
					ID:         143,
					Name:       "Закрыто и не реализовано",
					Sort:       11000,
					PipelineID: 3177727,
					Color:      domain.SolitudePipelineStatusColor,
					AccountID:  28847170,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}},
				},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetPipelineByIDResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineByID(ctx, 3177727)
		assert.NoError(t, err)
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineByID(ctx, 3177727)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Statuses в ответе", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineByID(ctx, 3177727)
		assert.EqualError(t, err, "Key: 'Pipeline.Embedded.Statuses' Error:Field validation for 'Statuses' failed on the 'gt' tag")
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ffc8c8","type":1,"account_id":28847170,"_links":{"self":{"href":"https://shard152.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineByID(ctx, 3177727)
		assert.EqualError(t, err, "Key: 'Pipeline.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestAddPipelines(t *testing.T) {
	const (
		sampleAddPipelinesResponseBody = `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3270358,"name":"Воронка для примера","sort":1,"is_main":true,"is_unsorted_on":false,"account_id":1415131,"request_id":"123","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358"}},"_embedded":{"statuses":[{"id":3304,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3270358,"color":"#ffc8c8","type":1,"account_id":1415131,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/3304"}}},{"id":3303,"name":"Первичный контакт","sort":20,"is_editable":true,"pipeline_id":3270358,"color":"#ffdc7f","type":0,"account_id":1415131,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/3303"}}},{"id":142,"name":"Мое название для успешных сделок","sort":10000,"is_editable":false,"pipeline_id":3270358,"color":"#c1e0ff","type":0,"account_id":1415131,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/142"}}},{"id":143,"name":"Закрыто и не реализовано","sort":11000,"is_editable":false,"pipeline_id":3270358,"color":"#e6e8ea","type":0,"account_id":1415131,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/143"}}}]}}]}}`
		requestBodyWant                = `[{"name":"Воронка доп продаж","sort":20,"is_main":false,"is_unsorted_on":true,"request_id":"123","_embedded":{"statuses":[{"id":142,"name":"Мое название для успешных сделок"},{"name":"Первичный контакт","sort":10,"color":"#fffeb2"}]}}]`
	)

	sampleAddPipelinesRequest := &AddPipelinesRequest{
		Add: []*AddPipelinesRequestData{
			{
				Name:         "Воронка доп продаж",
				Sort:         20,
				IsUnsortedOn: true,
				RequestID:    "123",
				Embedded: &AddPipelinesRequestDataEmbedded{
					Statuses: []*domain.EmbeddedPipelineStatus{
						{
							ID:   142,
							Name: "Мое название для успешных сделок",
						},
						{
							Name:  "Первичный контакт",
							Sort:  10,
							Color: domain.ShalimarPipelineStatusColor,
						},
					},
				},
			},
		},
	}

	responseWant := []*domain.Pipeline{
		{
			ID:        3270358,
			Name:      "Воронка для примера",
			Sort:      1,
			IsMain:    true,
			AccountID: 1415131,
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270358"}},
			Embedded: &domain.PipelineEmbedded{
				Statuses: []*domain.PipelineStatus{
					{
						ID:         3304,
						Name:       "Неразобранное",
						Sort:       10,
						PipelineID: 3270358,
						Color:      domain.YourPinkPipelineStatusColor,
						Type:       domain.RegularPipelineStatusType,
						AccountID:  1415131,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/3304"}},
					},
					{
						ID:         3303,
						Name:       "Первичный контакт",
						Sort:       20,
						IsEditable: true,
						PipelineID: 3270358,
						Color:      domain.SalomiePipelineStatusColor,
						AccountID:  1415131,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/3303"}},
					},
					{
						ID:         142,
						Name:       "Мое название для успешных сделок",
						Sort:       10000,
						PipelineID: 3270358,
						Color:      domain.MediumPattensBluePipelineStatusColor,
						AccountID:  1415131,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/142"}},
					},
					{
						ID:         143,
						Name:       "Закрыто и не реализовано",
						Sort:       11000,
						PipelineID: 3270358,
						Color:      domain.SolitudePipelineStatusColor,
						AccountID:  1415131,
						Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/143"}},
					},
				},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAddPipelinesResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddPipelines(ctx, sampleAddPipelinesRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddPipelines(ctx, sampleAddPipelinesRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Pipelines в ответе", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddPipelines(ctx, sampleAddPipelinesRequest)
		assert.EqualError(t, err, "Key: 'AddPipelinesResponse.Embedded.Pipelines' Error:Field validation for 'Pipelines' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Statuses в ответе", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3270358,"name":"Воронка для примера","sort":1,"is_main":true,"is_unsorted_on":false,"account_id":1415131,"request_id":"123","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358"}},"_embedded":{"statuses":[]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddPipelines(ctx, sampleAddPipelinesRequest)
		assert.EqualError(t, err, "Key: 'AddPipelinesResponse.Embedded.Pipelines[0].Embedded.Statuses' Error:Field validation for 'Statuses' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3270358,"name":"Воронка для примера","sort":1,"is_main":true,"is_unsorted_on":false,"account_id":1415131,"request_id":"123","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358"}},"_embedded":{"statuses":[{"id":3304,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3270358,"color":"#ffc8c8","type":1,"account_id":1415131,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270358/statuses/3304"}}}]}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.AddPipelines(ctx, sampleAddPipelinesRequest)
		assert.EqualError(t, err, "Key: 'AddPipelinesResponse.TotalItems' Error:Field validation for 'TotalItems' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestUpdatePipeline(t *testing.T) {
	const (
		sampleUpdatePipelineResponseBody = `{"id":3177727,"name":"Новое название для воронки","sort":1000,"is_main":false,"is_unsorted_on":false,"is_archive":false,"account_id":12345678,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392159,"name":"Первичный контакт","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#ffc8c8","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}}},{"id":32392165,"name":"Принимают решение","sort":30,"is_editable":true,"pipeline_id":3177727,"color":"#ffdc7f","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}}},{"id":142,"name":"Успешно реализовано","sort":10000,"is_editable":false,"pipeline_id":3177727,"color":"#c1e0ff","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}}},{"id":143,"name":"Закрыто и не реализовано","sort":11000,"is_editable":false,"pipeline_id":3177727,"color":"#e6e8ea","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}}}]}}`
		requestBodyWant                  = `{"name":"Новое название для воронки","sort":100,"is_main":false,"is_unsorted_on":false}`
	)

	sampleUpdatePipelineRequest := &UpdatePipelineRequest{
		Name:         "Новое название для воронки",
		IsMain:       false,
		IsUnsortedOn: false,
		Sort:         100,
	}

	responseWant := &domain.Pipeline{
		ID:        3177727,
		Name:      "Новое название для воронки",
		Sort:      1000,
		AccountID: 12345678,
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},
		Embedded: &domain.PipelineEmbedded{
			Statuses: []*domain.PipelineStatus{
				{
					ID:         32392159,
					Name:       "Первичный контакт",
					Sort:       20,
					IsEditable: true,
					PipelineID: 3177727,
					Color:      domain.YourPinkPipelineStatusColor,
					AccountID:  12345678,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}},
				},
				{
					ID:         32392165,
					Name:       "Принимают решение",
					Sort:       30,
					IsEditable: true,
					PipelineID: 3177727,
					Color:      domain.SalomiePipelineStatusColor,
					AccountID:  12345678,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}},
				},
				{
					ID:         142,
					Name:       "Успешно реализовано",
					Sort:       10000,
					PipelineID: 3177727,
					Color:      domain.MediumPattensBluePipelineStatusColor,
					AccountID:  12345678,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}},
				},
				{
					ID:         143,
					Name:       "Закрыто и не реализовано",
					Sort:       11000,
					PipelineID: 3177727,
					Color:      domain.SolitudePipelineStatusColor,
					AccountID:  12345678,
					Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}},
				},
			},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdatePipelineResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdatePipeline(ctx, 3177727, sampleUpdatePipelineRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdatePipeline(ctx, 3177727, sampleUpdatePipelineRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Statuses в ответе", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"id":3177727,"name":"Новое название для воронки","sort":1000,"is_main":false,"is_unsorted_on":false,"is_archive":false,"account_id":12345678,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdatePipeline(ctx, 3177727, sampleUpdatePipelineRequest)
		assert.EqualError(t, err, "Key: 'Pipeline.Embedded.Statuses' Error:Field validation for 'Statuses' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":3177727,"name":"Новое название для воронки","sort":1000,"is_main":false,"is_unsorted_on":false,"is_archive":false,"account_id":12345678,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392159,"name":"Первичный контакт","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#e6e8ea","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}}}]}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.UpdatePipeline(ctx, 3177727, sampleUpdatePipelineRequest)
		assert.EqualError(t, err, "Key: 'Pipeline.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestDeletePipeline(t *testing.T) {
	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		err = client.DeletePipeline(ctx, 3177727)
		assert.NoError(t, err)
	})

	t.Run("Ошибка при обработке", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, errorContentType)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, "")
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		err = client.DeletePipeline(ctx, 3177727)
		assert.EqualError(t, err, "http status not ok: 400")
	})
}
