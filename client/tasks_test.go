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

func TestGetTasksRequestFilter(t *testing.T) {
	testCases := []struct {
		name   string
		filter *GetTasksRequestFilter
		want   url.Values
	}{
		{
			name:   "Один ID",
			filter: &GetTasksRequestFilter{ID: request.CreateSimpleFilter("id", "1")},
			want:   url.Values{"filter[id]": []string{"1"}},
		},
		{
			name:   "Несколько ID",
			filter: &GetTasksRequestFilter{ID: request.CreateMultipleFilter("id", []string{"1", "2", "3"})},
			want:   url.Values{"filter[id][0]": []string{"1", "2", "3"}},
		},
		{
			name:   "ResponsibleUserID",
			filter: &GetTasksRequestFilter{ResponsibleUserID: request.CreateSimpleFilter("responsible_user_id", "name_value")},
			want:   url.Values{"filter[responsible_user_id]": []string{"name_value"}},
		},
		{
			name:   "Несколько ResponsibleUserID",
			filter: &GetTasksRequestFilter{ResponsibleUserID: request.CreateMultipleFilter("responsible_user_id", []string{"name_value_1", "name_value_2"})},
			want:   url.Values{"filter[responsible_user_id][0]": []string{"name_value_1", "name_value_2"}},
		},
		{
			name:   "IsCompleted",
			filter: &GetTasksRequestFilter{IsCompleted: request.CreateSimpleFilter("is_completed", "1")},
			want:   url.Values{"filter[is_completed]": []string{"1"}},
		},
		{
			name:   "TaskType",
			filter: &GetTasksRequestFilter{TaskType: request.CreateSimpleFilter("task_type", "2")},
			want:   url.Values{"filter[task_type]": []string{"2"}},
		},
		{
			name:   "Несколько TaskType",
			filter: &GetTasksRequestFilter{TaskType: request.CreateMultipleFilter("task_type", []string{"1", "2"})},
			want:   url.Values{"filter[task_type][0]": []string{"1", "2"}},
		},
		{
			name:   "EntityType",
			filter: &GetTasksRequestFilter{EntityType: request.CreateSimpleFilter("entity_type", "leads")},
			want:   url.Values{"filter[entity_type]": []string{"leads"}},
		},
		{
			name:   "Несколько EntityType",
			filter: &GetTasksRequestFilter{EntityType: request.CreateMultipleFilter("entity_type", []string{"leads", "contacts"})},
			want:   url.Values{"filter[entity_type][0]": []string{"leads", "contacts"}},
		},
		{
			name:   "EntityID",
			filter: &GetTasksRequestFilter{EntityID: request.CreateSimpleFilter("entity_id", "123")},
			want:   url.Values{"filter[entity_id]": []string{"123"}},
		},
		{
			name:   "Несколько EntityID",
			filter: &GetTasksRequestFilter{EntityID: request.CreateMultipleFilter("entity_id", []string{"123", "234", "345"})},
			want:   url.Values{"filter[entity_id][0]": []string{"123", "234", "345"}},
		},
		{
			name:   "Интервал UpdatedAt",
			filter: &GetTasksRequestFilter{UpdatedAt: request.CreateIntervalFilter("updated_at", "12345678", "23456789")},
			want:   url.Values{"filter[updated_at][from]": []string{"12345678"}, "filter[updated_at][to]": []string{"23456789"}},
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

func TestGetTasksRequestParamsValidation(t *testing.T) {
	v := validator.New()

	t.Run("Превышен лимит элементов в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Limit: 260}
		assert.EqualError(t, v.Struct(req), "Key: 'GetTasksRequestParams.Limit' Error:Field validation for 'Limit' failed on the 'lte' tag")
	})

	t.Run("Диапазонный фильтр по ID в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{ID: request.CreateIntervalFilter("id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ID filter must be simple or multiple type")
	})

	t.Run("Диапазонный фильтр по ResponsibleUserID в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{ResponsibleUserID: request.CreateIntervalFilter("responsible_user_id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "ResponsibleUserID filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по IsCompleted в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{IsCompleted: request.CreateMultipleFilter("is_completed", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "IsCompleted filter must be simple type")
	})

	t.Run("Диапазонный фильтр по IsCompleted в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{IsCompleted: request.CreateIntervalFilter("is_completed", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "IsCompleted filter must be simple type")
	})

	t.Run("Диапазонный фильтр по TaskType в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{TaskType: request.CreateIntervalFilter("task_type", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "TaskType filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по EntityType в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{EntityType: request.CreateMultipleFilter("entity_type", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "EntityType filter must be simple type")
	})

	t.Run("Диапазонный фильтр по EntityType в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{EntityType: request.CreateIntervalFilter("entity_type", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "EntityType filter must be simple type")
	})

	t.Run("Диапазонный фильтр по EntityID в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{EntityID: request.CreateIntervalFilter("entity_id", "123", "234")}}
		assert.EqualError(t, req.Filter.validate(), "EntityID filter must be simple or multiple type")
	})

	t.Run("Множественный фильтр по UpdatedAt в запросе", func(t *testing.T) {
		req := &GetTasksRequestParams{Filter: &GetTasksRequestFilter{UpdatedAt: request.CreateMultipleFilter("updated_at", []string{"123", "234"})}}
		assert.EqualError(t, req.Filter.validate(), "UpdatedAt filter must be simple or interval type")
	})
}

func TestAddTasksRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой запрос", func(t *testing.T) {
		data := &AddTasksRequest{}
		assert.EqualError(t, v.Struct(data), "Key: 'AddTasksRequest.Add' Error:Field validation for 'Add' failed on the 'required' tag")
	})

	t.Run("Пустой массив в запросе", func(t *testing.T) {
		data := &AddTasksRequest{Add: []*AddTasksRequestData{}}
		assert.EqualError(t, v.Struct(data), "Key: 'AddTasksRequest.Add' Error:Field validation for 'Add' failed on the 'gt' tag")
	})

	t.Run("Пустой объект в запросе", func(t *testing.T) {
		data := &AddTasksRequest{Add: []*AddTasksRequestData{{}}}
		assert.EqualError(t, v.Struct(data), `Key: 'AddTasksRequest.Add[0].Text' Error:Field validation for 'Text' failed on the 'required' tag
Key: 'AddTasksRequest.Add[0].CompleteTill' Error:Field validation for 'CompleteTill' failed on the 'required' tag`)
	})
}

func TestUpdateTasksRequestDataValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой запрос", func(t *testing.T) {
		data := &UpdateTasksRequest{}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateTasksRequest.Update' Error:Field validation for 'Update' failed on the 'required' tag")
	})

	t.Run("Пустой массив в запросе", func(t *testing.T) {
		data := &UpdateTasksRequest{Update: []*UpdateTasksRequestData{}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateTasksRequest.Update' Error:Field validation for 'Update' failed on the 'gt' tag")
	})

	t.Run("Пустой объект в запросе", func(t *testing.T) {
		data := &UpdateTasksRequest{Update: []*UpdateTasksRequestData{{}}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateTasksRequest.Update[0].ID' Error:Field validation for 'ID' failed on the 'required' tag")
	})

	t.Run("Невалидный TaskTypeID в запросе", func(t *testing.T) {
		data := &UpdateTasksRequest{Update: []*UpdateTasksRequestData{{ID: 1, TaskTypeID: 3}}}
		assert.EqualError(t, v.Struct(data), "Key: 'UpdateTasksRequest.Update[0].TaskTypeID' Error:Field validation for 'TaskTypeID' failed on the 'oneof' tag")
	})
}

func TestGetTasksResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		data := &GetTasksResponse{}
		assert.EqualError(t, v.Struct(data), `Key: 'GetTasksResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag
Key: 'GetTasksResponse.Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})

	t.Run("Пустой массив Leads в ответе", func(t *testing.T) {
		data := &GetTasksResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetTasksResponseEmbedded{Tasks: []*domain.Task{}}}
		assert.NoError(t, v.Struct(data))
	})

	t.Run("Ни одного обязательного параметра Lead в ответе", func(t *testing.T) {
		data := &GetTasksResponse{Page: 1, Links: &domain.Links{Self: &domain.Link{Href: "url"}}, Embedded: &GetTasksResponseEmbedded{Tasks: []*domain.Task{{}}}}
		assert.EqualError(t, v.Struct(data), `Key: 'GetTasksResponse.Embedded.Tasks[0].ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].UpdatedBy' Error:Field validation for 'UpdatedBy' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].ResponsibleUserID' Error:Field validation for 'ResponsibleUserID' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].GroupID' Error:Field validation for 'GroupID' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].EntityID' Error:Field validation for 'EntityID' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].EntityType' Error:Field validation for 'EntityType' failed on the 'oneof' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].TaskTypeID' Error:Field validation for 'TaskTypeID' failed on the 'oneof' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].CompleteTill' Error:Field validation for 'CompleteTill' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].AccountID' Error:Field validation for 'AccountID' failed on the 'required' tag
Key: 'GetTasksResponse.Embedded.Tasks[0].Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})
}

func TestGetTasks(t *testing.T) {
	const sampleGetTasksResponseBody = `{"_page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks?filter[task_type][]=2&filter[is_completed][]=1&limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/tasks?filter[task_type][]=2&filter[is_completed][]=1&limit=2&page=2"}},"_embedded":{"tasks":[{"id":7087,"created_by":3987910,"updated_by":3987910,"created_at":1575364000,"updated_at":1575364851,"responsible_user_id":123123,"group_id":1,"entity_id":167353,"entity_type":"leads","duration":0,"is_completed":true,"task_type_id":2,"text":"Пригласить на бесплатную тренировку","complete_till":1575665940,"account_id":321321,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/7087"}}},{"id":215089,"created_by":3987910,"updated_by":3987910,"created_at":1576767879,"updated_at":1576767914,"responsible_user_id":123123,"group_id":1,"entity_id":1035487,"entity_type":"leads","duration":0,"is_completed":true,"task_type_id":2,"text":"Назначить встречу с клиентом","complete_till":1576768179,"account_id":321312,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/215089"}}}]}}`

	requestParamsWant := url.Values{
		"page":                 []string{"1"},
		"limit":                []string{"2"},
		"filter[task_type]":    []string{"2"},
		"filter[is_completed]": []string{"1"},
		"order[id]":            []string{"asc"},
	}

	sampleGetTasksRequestParams := &GetTasksRequestParams{
		Page:  1,
		Limit: 2,
		Filter: &GetTasksRequestFilter{
			TaskType:    request.CreateSimpleFilter("task_type", "2"),
			IsCompleted: request.CreateSimpleFilter("is_completed", "1"),
		},
		Order: &GetTasksOrder{
			By:     IDGetTasksOrderBy,
			Method: request.AscendingOrderMethod,
		},
	}

	responseWant := []*domain.Task{
		{
			ID:                7087,
			CreatedBy:         3987910,
			UpdatedBy:         3987910,
			CreatedAt:         1575364000,
			UpdatedAt:         1575364851,
			ResponsibleUserID: 123123,
			GroupID:           1,
			EntityID:          167353,
			EntityType:        "leads",
			Duration:          0,
			IsCompleted:       true,
			TaskTypeID:        2,
			Text:              "Пригласить на бесплатную тренировку",
			CompleteTill:      1575665940,
			AccountID:         321321,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/7087"}},
		},
		{
			ID:                215089,
			CreatedBy:         3987910,
			UpdatedBy:         3987910,
			CreatedAt:         1576767879,
			UpdatedAt:         1576767914,
			ResponsibleUserID: 123123,
			GroupID:           1,
			EntityID:          1035487,
			EntityType:        "leads",
			IsCompleted:       true,
			TaskTypeID:        2,
			Text:              "Назначить встречу с клиентом",
			CompleteTill:      1576768179,
			AccountID:         321312,
			Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/215089"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetTasksResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetTasks(ctx, sampleGetTasksRequestParams)
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

		responseGot, err := client.GetTasks(ctx, sampleGetTasksRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Пустой массив в ответе", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_page":2,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/leads?limit=2&page=2"}},"_embedded":{"tasks":[]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetTasks(ctx, sampleGetTasksRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"page":1,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks?filter[task_type][]=2&filter[is_completed][]=1&limit=2&page=1"},"next":{"href":"https://example.amocrm.ru/api/v4/tasks?filter[task_type][]=2&filter[is_completed][]=1&limit=2&page=2"}},"_embedded":{"tasks":[{"id":7087,"created_by":3987910,"updated_by":3987910,"created_at":1575364000,"updated_at":1575364851,"responsible_user_id":123123,"group_id":1,"entity_id":167353,"entity_type":"leads","duration":0,"is_completed":true,"task_type_id":2,"text":"Пригласить на бесплатную тренировку","complete_till":1575665940,"account_id":321321,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/7087"}}}]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetTasks(ctx, sampleGetTasksRequestParams)
		assert.EqualError(t, err, "Key: 'GetTasksResponse.Page' Error:Field validation for 'Page' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}

func TestGetTaskByID(t *testing.T) {
	const sampleGetTaskByIDResponseBody = `{"id":56981,"created_by":54224,"updated_by":3987910,"created_at":1575910123,"updated_at":1576767989,"responsible_user_id":123123,"group_id":1,"entity_id":180765,"entity_type":"leads","duration":0,"is_completed":true,"task_type_id":2,"text":"Назначить встречу с клиентом","result":{"text":"Результат есть"},"complete_till":1575910423,"account_id":321312,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/56981"}}}`

	responseWant := &domain.Task{
		ID:                56981,
		CreatedBy:         54224,
		UpdatedBy:         3987910,
		CreatedAt:         1575910123,
		UpdatedAt:         1576767989,
		ResponsibleUserID: 123123,
		GroupID:           1,
		EntityID:          180765,
		EntityType:        "leads",
		IsCompleted:       true,
		TaskTypeID:        2,
		Text:              "Назначить встречу с клиентом",
		Result:            &domain.TaskResult{Text: "Результат есть"},
		CompleteTill:      1575910423,
		AccountID:         321312,
		Links:             &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/56981"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetTaskByIDResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetTaskByID(ctx, 56981)
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

		responseGot, err := client.GetTaskByID(ctx, 56981)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetTaskByID(ctx, 0)
		assert.EqualError(t, err, ErrInvalidTaskID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":56981,"created_by":54224,"updated_by":3987910,"created_at":1575910123,"updated_at":1576767989,"responsible_user_id":123123,"group_id":1,"entity_id":180765,"entity_type":"leads","duration":0,"is_completed":true,"task_type_id":2,"text":"Назначить встречу с клиентом","result":{"text":"Результат есть"},"complete_till":1575910423,"account_id":321312,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/56981"}}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.GetTaskByID(ctx, 56981)
		assert.EqualError(t, err, "Key: 'Task.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Empty(t, responseGot)
	})
}

func TestAddTasks(t *testing.T) {
	const (
		requestBodyWant            = `[{"entity_id":9785993,"entity_type":"leads","task_type_id":1,"text":"Встретиться с клиентом Иван Иванов","complete_till":1588885140,"request_id":"example"}]`
		sampleAddTasksResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks"}},"_embedded":{"tasks":[{"id":4745251,"request_id":"example","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4745251"}}}]}}`
	)

	sampleAddTasksRequest := &AddTasksRequest{
		Add: []*AddTasksRequestData{
			{
				TaskTypeID:   domain.CallTaskType,
				Text:         "Встретиться с клиентом Иван Иванов",
				CompleteTill: 1588885140,
				EntityID:     9785993,
				EntityType:   domain.LeadsEntityType,
				RequestID:    "example",
			},
		},
	}

	responseWant := []*AddTasksResponseItem{
		{
			ID:        4745251,
			RequestID: "example",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/4745251"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleAddTasksResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddTasks(ctx, sampleAddTasksRequest)
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

		responseGot, err := client.AddTasks(ctx, sampleAddTasksRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks"}},"_embedded":{}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.AddTasks(ctx, sampleAddTasksRequest)
		assert.EqualError(t, err, "Key: 'AddTasksResponse.Embedded.Tasks' Error:Field validation for 'Tasks' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateTasks(t *testing.T) {
	const (
		requestBodyWant               = `[{"id":4745251,"task_type_id":2,"text":"Новое название для задачи","complete_till":1588885140},{"id":4747929,"task_type_id":1,"text":"Новое название для задачи 2","complete_till":1588885140}]`
		sampleUpdateTasksResponseBody = `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks"}},"_embedded":{"tasks":[{"id":4745251,"updated_at":1588760725,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4745251"}}},{"id":4747929,"updated_at":1588760725,"request_id":"1","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4747929"}}}]}}`
	)

	sampleUpdateTasksRequest := &UpdateTasksRequest{
		Update: []*UpdateTasksRequestData{
			{
				ID:           4745251,
				TaskTypeID:   domain.MeetingTaskType,
				Text:         "Новое название для задачи",
				CompleteTill: 1588885140,
			},
			{
				ID:           4747929,
				TaskTypeID:   domain.CallTaskType,
				Text:         "Новое название для задачи 2",
				CompleteTill: 1588885140,
			},
		},
	}

	responseWant := []*UpdateTasksResponseItem{
		{
			ID:        4745251,
			UpdatedAt: 1588760725,
			RequestID: "0",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/4745251"}},
		},
		{
			ID:        4747929,
			UpdatedAt: 1588760725,
			RequestID: "1",
			Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/4747929"}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateTasksResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateTasks(ctx, sampleUpdateTasksRequest)
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

		responseGot, err := client.UpdateTasks(ctx, sampleUpdateTasksRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks"}},"_embedded":{"tasks":[]}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateTasks(ctx, sampleUpdateTasksRequest)
		assert.EqualError(t, err, "Key: 'UpdateTasksResponse.Embedded.Tasks' Error:Field validation for 'Tasks' failed on the 'gt' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestUpdateTask(t *testing.T) {
	const (
		requestBodyWant              = `{"id":4745251,"task_type_id":2,"text":"Новое название для задачи","complete_till":1588885140}`
		sampleUpdateTaskResponseBody = `{"id":4745251,"updated_at":1588760725,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4745251"}}}`
	)

	sampleUpdateTaskRequest := &UpdateTasksRequestData{
		ID:           4745251,
		TaskTypeID:   domain.MeetingTaskType,
		Text:         "Новое название для задачи",
		CompleteTill: 1588885140,
	}

	responseWant := &UpdateTasksResponseItem{
		ID:        4745251,
		UpdatedAt: 1588760725,
		RequestID: "0",
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/4745251"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateTaskResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateTask(ctx, 4745251, sampleUpdateTaskRequest)
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

		responseGot, err := client.UpdateTask(ctx, 4745251, sampleUpdateTaskRequest)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный запрос", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateTask(ctx, 0, sampleUpdateTaskRequest)
		assert.EqualError(t, err, ErrInvalidTaskID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":4745251,"updated_at":1588760725,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4745251"}}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.UpdateTask(ctx, 4745251, sampleUpdateTaskRequest)
		assert.EqualError(t, err, "Key: 'UpdateTasksResponseItem.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}

func TestCompleteTask(t *testing.T) {
	const (
		requestBodyWant              = `{"id":4747929,"is_completed":true,"result":{"text":"Удалось связаться с клиентом"}}`
		sampleUpdateTaskResponseBody = `{"id":4747929,"updated_at":1588770600,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4747929"}}}`
	)

	//sampleUpdateTaskRequest := &UpdateTasksRequestData{
	//	ID:           4747929,
	//	IsCompleted: true,
	//	Result: &domain.TaskResult{Text: "Удалось связаться с клиентом"},
	//}

	responseWant := &UpdateTasksResponseItem{
		ID:        4747929,
		UpdatedAt: 1588770600,
		RequestID: "0",
		Links:     &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/tasks/4747929"}},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleUpdateTaskResponseBody)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.CompleteTask(ctx, 4747929, "Удалось связаться с клиентом")
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

		responseGot, err := client.CompleteTask(ctx, 4747929, "Удалось связаться с клиентом")
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный taskID в запросе", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.CompleteTask(ctx, 0, "Удалось связаться с клиентом")
		assert.EqualError(t, err, ErrInvalidTaskID.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный result в запросе", func(t *testing.T) {
		client, err := NewClient("localhost:1234", "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.CompleteTask(ctx, 4747929, "")
		assert.EqualError(t, err, ErrInvalidTaskResult.Error())
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		var requestBodyGot []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBodyGot, _ = ioutil.ReadAll(r.Body)
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"_id":4747929,"updated_at":1588770600,"request_id":"0","_links":{"self":{"href":"https://example.amocrm.ru/api/v4/tasks/4747929"}}}`)
		}))

		client, err := NewClient(server.URL, "login", "hash")
		assert.NoError(t, err)

		responseGot, err := client.CompleteTask(ctx, 4747929, "Удалось связаться с клиентом")
		assert.EqualError(t, err, "Key: 'UpdateTasksResponseItem.ID' Error:Field validation for 'ID' failed on the 'required' tag")
		assert.Equal(t, requestBodyWant, string(requestBodyGot))
		assert.Empty(t, responseGot)
	})
}
