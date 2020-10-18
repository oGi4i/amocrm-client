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

func TestAddPipelineStatusDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Невалидный цвет в запросе", func(t *testing.T) {
		data := &AddPipelineStatusData{Color: "#ffffff"}
		assert.EqualError(t, v.Struct(data), "Key: 'AddPipelineStatusData.Color' Error:Field validation for 'Color' failed on the 'oneof' tag")
	})
}

func TestUpdatePipelineStatusDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Невалидный цвет в запросе", func(t *testing.T) {
		data := &UpdatePipelineStatusData{Color: "#ffffff"}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdatePipelineStatusData.Color' Error:Field validation for 'Color' failed on the 'oneof' tag")
	})
}

func TestAddPipelineStatusResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &AddPipelineStatusResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'AddPipelineStatusResponse.TotalItems' Error:Field validation for 'TotalItems' failed on the 'required' tag
Key: 'AddPipelineStatusResponse.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Statuses в ответе", func(t *testing.T) {
		data := &AddPipelineStatusResponse{TotalItems: 1, Embedded: &domain.PipelineEmbedded{Statuses: []*domain.PipelineStatus{}}}
		assert.EqualError(t, v.Struct(data), "Key: 'AddPipelineStatusResponse.Embedded.Statuses' Error:Field validation for 'Statuses' failed on the 'gt' tag")
	})

	t.Run("Ни одного обязательного параметра в Status в ответе", func(t *testing.T) {
		data := &AddPipelineStatusResponse{TotalItems: 1, Embedded: &domain.PipelineEmbedded{Statuses: []*domain.PipelineStatus{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'AddPipelineStatusResponse.Embedded.Statuses[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'AddPipelineStatusResponse.Embedded.Statuses[0].Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'AddPipelineStatusResponse.Embedded.Statuses[0].Sort' Error:Field validation for 'Sort' failed on the 'required' tag
Key: 'AddPipelineStatusResponse.Embedded.Statuses[0].PipelineID' Error:Field validation for 'PipelineID' failed on the 'required' tag
Key: 'AddPipelineStatusResponse.Embedded.Statuses[0].AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'AddPipelineStatusResponse.Embedded.Statuses[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetPipelineStatuses(t *testing.T) {
	const sampleGetPipelineStatusesResponseBody = `{"_total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ffdc7f","type":1,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}},{"id":32392159,"name":"Первичный контакт","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#ccc8f9","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}}},{"id":32392165,"name":"Принимают решение","sort":30,"is_editable":true,"pipeline_id":3177727,"color":"#c1e0ff","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}}},{"id":142,"name":"Успешно реализовано","sort":10000,"is_editable":false,"pipeline_id":3177727,"color":"#f2f3f4","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/142"}}},{"id":143,"name":"Закрыто и не реализовано","sort":11000,"is_editable":false,"pipeline_id":3177727,"color":"#e6e8ea","type":0,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/143"}}}]}}]}}`

	responseWant := []*domain.PipelineStatus{
		{
			ID:         32392156,
			Name:       "Неразобранное",
			Sort:       10,
			PipelineID: 3177727,
			Color:      domain.SalomiePipelineStatusColor,
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
			Color:      domain.LavanderBluePipelineStatusColor,
			AccountID:  12345678,
			Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392159"}},
		},
		{
			ID:         32392165,
			Name:       "Принимают решение",
			Sort:       30,
			IsEditable: true,
			PipelineID: 3177727,
			Color:      domain.MediumPattensBluePipelineStatusColor,
			AccountID:  12345678,
			Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}},
		},
		{
			ID:         142,
			Name:       "Успешно реализовано",
			Sort:       10000,
			PipelineID: 3177727,
			Color:      domain.AliceBluePipelineStatusColor,
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
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetPipelineStatusesResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatuses(ctx, 3177727)
		assert.NoError(t, err)
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatuses(ctx, 3177727)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив Pipelines в ответе", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_total_items":1,"_links":{"self":{"href":"url"}},"_embedded":{"items":[]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatuses(ctx, 3177727)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"total_items":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines"}},"_embedded":{"pipelines":[{"id":3177727,"name":"Воронка","sort":1,"is_main":true,"is_unsorted_on":true,"is_archive":false,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727"}},"_embedded":{"statuses":[{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#f9deff","type":1,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}}]}}]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatuses(ctx, 3177727)
		assert.EqualError(t, err, "Key: 'PipelinesResponse.TotalItems' Error:Field validation for 'TotalItems' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestGetPipelineStatusByID(t *testing.T) {
	const sampleGetPipelineStatusByIDResponseBody = `{"id":32392156,"name":"Неразобранное","sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ccc8f9","type":1,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}}`

	responseWant := &domain.PipelineStatus{
		ID:         32392156,
		Name:       "Неразобранное",
		Sort:       10,
		PipelineID: 3177727,
		Color:      domain.LavanderBluePipelineStatusColor,
		Type:       domain.RegularPipelineStatusType,
		AccountID:  12345678,
		Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetPipelineStatusByIDResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatusByID(ctx, 3177727, 32392156)
		assert.NoError(t, err)
		assert.Exactly(t, responseWant, responseGot)
	})

	t.Run("Пустое тело ответа", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatusByID(ctx, 3177727, 32392156)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"id":32392156,"name":"Неразобранное","_sort":10,"is_editable":false,"pipeline_id":3177727,"color":"#ccc8f9","type":1,"account_id":12345678,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392156"}}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetPipelineStatusByID(ctx, 3177727, 32392156)
		assert.EqualError(t, err, "Key: 'PipelineStatus.Sort' Error:Field validation for 'Sort' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestAddPipelineStatuses(t *testing.T) {
	const (
		sampleAddPipelineStatusesResponseBody = `{"_total_items":2,"_embedded":{"statuses":[{"id":33035290,"name":"Новый этап","sort":60,"is_editable":true,"pipeline_id":3270355,"color":"#fffeb2","type":0,"account_id":1415131,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270355/statuses/33035290"}}},{"id":33035293,"name":"Новый этап 2","sort":70,"is_editable":true,"pipeline_id":3270355,"color":"#fffeb2","type":0,"account_id":1415131,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270355/statuses/33035293"}}}]}}`
		requestBodyWant                       = `[{"name":"Новый этап","sort":100,"color":"#fffeb2"},{"name":"Новый этап 2","sort":200,"color":"#fffeb2"}]`
	)

	sampleAddPipelineStatusesRequest := []*AddPipelineStatusData{
		{
			Name:  "Новый этап",
			Sort:  100,
			Color: domain.ShalimarPipelineStatusColor,
		},
		{
			Name:  "Новый этап 2",
			Sort:  200,
			Color: domain.ShalimarPipelineStatusColor,
		},
	}

	responseWant := []*domain.PipelineStatus{
		{
			ID:         33035290,
			Name:       "Новый этап",
			Sort:       60,
			IsEditable: true,
			PipelineID: 3270355,
			Color:      domain.ShalimarPipelineStatusColor,
			AccountID:  1415131,
			Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270355/statuses/33035290"}},
		},
		{
			ID:         33035293,
			Name:       "Новый этап 2",
			Sort:       70,
			IsEditable: true,
			PipelineID: 3270355,
			Color:      domain.ShalimarPipelineStatusColor,
			AccountID:  1415131,
			Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3270355/statuses/33035293"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var reqeustBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqeustBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAddPipelineStatusesResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddPipelineStatuses(ctx, 3270355, sampleAddPipelineStatusesRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(reqeustBodyGot))
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

		responseGot, err := client.AddPipelineStatuses(ctx, 3270355, sampleAddPipelineStatusesRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив статусов в ответе", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_total_items":2,"_embedded":{"statuses":[]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddPipelineStatuses(ctx, 3270355, sampleAddPipelineStatusesRequest)
		assert.EqualError(t, err, "Key: 'AddPipelineStatusResponse.Embedded.Statuses' Error:Field validation for 'Statuses' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"total_items":2,"_embedded":{"statuses":[{"id":33035290,"name":"Новый этап","sort":60,"is_editable":true,"pipeline_id":3270355,"color":"#fffeb2","type":0,"account_id":1415131,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270355/statuses/33035290"}}},{"id":33035293,"name":"Новый этап 2","sort":70,"is_editable":true,"pipeline_id":3270355,"color":"#fffeb2","type":0,"account_id":1415131,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3270355/statuses/33035293"}}}]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddPipelineStatuses(ctx, 3270355, sampleAddPipelineStatusesRequest)
		assert.EqualError(t, err, "Key: 'AddPipelineStatusResponse.TotalItems' Error:Field validation for 'TotalItems' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdatePipelineStatuse(t *testing.T) {
	const (
		sampleUpdatePipelineStatuseResponseBody = `{"id":32392165,"name":"Новое название для статуса","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#c1e0ff","type":0,"account_id":12345678,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}}}`
		requestBodyWant                         = `{"name":"Новое название для статуса","color":"#c1e0ff"}`
	)

	sampleUpdatePipelineStatusRequest := &UpdatePipelineStatusData{
		Name:  "Новое название для статуса",
		Color: domain.MediumPattensBluePipelineStatusColor,
	}

	responseWant := &domain.PipelineStatus{
		ID:         32392165,
		Name:       "Новое название для статуса",
		Sort:       20,
		IsEditable: true,
		PipelineID: 3177727,
		Color:      domain.MediumPattensBluePipelineStatusColor,
		AccountID:  12345678,
		Links:      &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var reqeustBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqeustBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdatePipelineStatuseResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdatePipelineStatus(ctx, 3177727, 32392165, sampleUpdatePipelineStatusRequest)
		assert.NoError(t, err)
		assert.Equal(t, requestBodyWant, string(reqeustBodyGot))
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

		responseGot, err := client.UpdatePipelineStatus(ctx, 3177727, 32392165, sampleUpdatePipelineStatusRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":32392165,"name":"Новое название для статуса","sort":20,"is_editable":true,"pipeline_id":3177727,"color":"#c1e0ff","type":0,"account_id":12345678,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads/pipelines/3177727/statuses/32392165"}}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdatePipelineStatus(ctx, 3177727, 32392165, sampleUpdatePipelineStatusRequest)
		assert.EqualError(t, err, "Key: 'PipelineStatus.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestDeletePipelineStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, "")
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		err = client.DeletePipelineStatus(ctx, 3177727, 32392165)
		assert.NoError(t, err)
	})

	t.Run("Ошибка при обработке", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, errorContentType)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, "")
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		err = client.DeletePipelineStatus(ctx, 3177727, 32392165)
		assert.EqualError(t, err, "http status not ok: 400")
	})
}
