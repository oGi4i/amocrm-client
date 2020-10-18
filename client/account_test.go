package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/ogi4i/amocrm-client/domain"
)

func TestJoinAccountRequestWith(t *testing.T) {
	testCases := []struct {
		name   string
		params []AccountGetRequestWith
		want   string
	}{
		{
			name:   "Ни одного параметра",
			params: []AccountGetRequestWith{},
			want:   "",
		},
		{
			name:   "Один параметр",
			params: []AccountGetRequestWith{TaskTypesAccountGetRequestWith},
			want:   "task_types",
		},
		{
			name:   "Два параметра",
			params: []AccountGetRequestWith{TaskTypesAccountGetRequestWith, VersionAccountGetRequestWith},
			want:   "task_types,version",
		},
		{
			name: "Все параметры",
			params: []AccountGetRequestWith{
				AmojoIDAccountGetRequestWith,
				AmojoRightsAccountGetRequestWith,
				UserGroupsAccountGetRequestWith,
				TaskTypesAccountGetRequestWith,
				VersionAccountGetRequestWith,
				EntityNamesAccountGetRequestWith,
				DatetimeSettingsAccountGetRequestWith,
			},
			want: "amojo_id,amojo_rights,user_groups,task_types,version,entity_names,datetime_settings",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, joinAccountRequestWith(tc.params))
		})
	}
}

func TestGetAccountRequestValidation(t *testing.T) {
	v := validator.New()

	t.Run("Пустой массив параметров with в запросе", func(t *testing.T) {
		req := &AccountGetRequestParams{With: []AccountGetRequestWith{}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Один параметр with в запросе", func(t *testing.T) {
		req := &AccountGetRequestParams{With: []AccountGetRequestWith{VersionAccountGetRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Несколько параметров with в запросе", func(t *testing.T) {
		req := &AccountGetRequestParams{With: []AccountGetRequestWith{VersionAccountGetRequestWith, TaskTypesAccountGetRequestWith}}
		assert.NoError(t, v.Struct(req))
	})

	t.Run("Невалидный параметр with в запросе", func(t *testing.T) {
		req := &AccountGetRequestParams{With: []AccountGetRequestWith{"with"}}
		assert.EqualError(t, v.Struct(req), "Key: 'AccountGetRequestParams.With[0]' Error:Field validation for 'With[0]' failed on the 'oneof' tag")
	})
}

func TestGetAccountResponseValidation(t *testing.T) {
	v := validator.New()

	t.Run("Ни одного обязательного параметра в ответе", func(t *testing.T) {
		req := &domain.Account{}
		assert.EqualError(t, v.Struct(req), `Key: 'Account.ID' Error:Field validation for 'ID' failed on the 'required' tag
Key: 'Account.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'Account.Subdomain' Error:Field validation for 'Subdomain' failed on the 'required' tag
Key: 'Account.CreatedAt' Error:Field validation for 'CreatedAt' failed on the 'required' tag
Key: 'Account.CreatedBy' Error:Field validation for 'CreatedBy' failed on the 'required' tag
Key: 'Account.UpdatedAt' Error:Field validation for 'UpdatedAt' failed on the 'required' tag
Key: 'Account.UpdatedBy' Error:Field validation for 'UpdatedBy' failed on the 'required' tag
Key: 'Account.CurrentUserID' Error:Field validation for 'CurrentUserID' failed on the 'required' tag
Key: 'Account.Country' Error:Field validation for 'Country' failed on the 'required' tag
Key: 'Account.CustomersMode' Error:Field validation for 'CustomersMode' failed on the 'required' tag
Key: 'Account.ContactNameDisplayOrder' Error:Field validation for 'ContactNameDisplayOrder' failed on the 'required' tag
Key: 'Account.Embedded' Error:Field validation for 'Embedded' failed on the 'required' tag
Key: 'Account.Links' Error:Field validation for 'Links' failed on the 'required' tag`)
	})

	t.Run("Пустой массив UsersGroups в ответе", func(t *testing.T) {
		req := &domain.AccountEmbedded{UsersGroups: []*domain.UserGroup{}}
		assert.EqualError(t, v.Struct(req), "Key: 'AccountEmbedded.UsersGroups' Error:Field validation for 'UsersGroups' failed on the 'gt' tag")
	})

	t.Run("Пустой массив TaskTypes в ответе", func(t *testing.T) {
		req := &domain.AccountEmbedded{TaskTypes: []*domain.TaskTypeInfo{}}
		assert.EqualError(t, v.Struct(req), "Key: 'AccountEmbedded.TaskTypes' Error:Field validation for 'TaskTypes' failed on the 'gt' tag")
	})
}

func TestGetAccount(t *testing.T) {
	const sampleGetAccountResponseBody = `{"id":1231414,"name":"example","subdomain":"example","created_at":1585840134,"created_by":321321,"updated_at":1589472711,"updated_by":321321,"current_user_id":581651,"country":"RU","customers_mode":"segments","is_unsorted_on":true,"is_loss_reason_enabled":true,"is_helpbot_enabled":false,"is_technical_account":false,"contact_name_display_order":1,"amojo_id":"f3c6340d-410e-4ad1-9f7e-c5e663599909","uuid":"824f3a59-6154-4edf-ba90-0b5593715d07","version":11,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/account"}},"_embedded":{"amojo_rights":{"can_direct":true,"can_create_groups":true},"users_groups":[{"id":1,"name":"Отдел продаж","uuid":null}],"task_types":[{"id":1,"name":"Связаться","color":null,"icon_id":null,"code":"FOLLOW_UP"},{"id":2,"name":"Встреча","color":null,"icon_id":null,"code":"MEETING"}],"entity_names":{"leads":{"ru":{"gender":"m","plural_form":{"dative":"клиентам","default":"клиенты","genitive":"клиентов","accusative":"клиентов","instrumental":"клиентами","prepositional":"клиентах"},"singular_form":{"dative":"клиенту","default":"клиент","genitive":"клиента","accusative":"клиента","instrumental":"клиентом","prepositional":"клиенте"}},"en":{"singular_form":{"default":"lead"},"plural_form":{"default":"leads"},"gender":"f"},"es":{"singular_form":{"default":"acuerdo"},"plural_form":{"default":"acuerdos"},"gender":"m"}}},"datetime_settings":{"date_pattern":"d.m.Y H:i","short_date_pattern":"d.m.Y","short_time_pattern":"H:i","date_format":"d.m.Y","time_format":"H:i:s","timezone":"Europe/Moscow","timezone_offset":"+03:00"}}}`

	requestParamsWant := url.Values{
		"with": []string{"amojo_rights,user_groups,task_types,entity_names,datetime_settings"},
	}

	sampleGetAccountRequestParams := &AccountGetRequestParams{
		With: []AccountGetRequestWith{AmojoRightsAccountGetRequestWith, UserGroupsAccountGetRequestWith, TaskTypesAccountGetRequestWith, EntityNamesAccountGetRequestWith, DatetimeSettingsAccountGetRequestWith},
	}

	responseWant := &domain.Account{
		ID:                      1231414,
		Name:                    "example",
		Subdomain:               "example",
		CreatedAt:               1585840134,
		CreatedBy:               321321,
		UpdatedAt:               1589472711,
		UpdatedBy:               321321,
		CurrentUserID:           581651,
		Country:                 "RU",
		CustomersMode:           domain.SegmentsCustomersMode,
		IsUnsortedOn:            true,
		IsLossReasonEnabled:     true,
		ContactNameDisplayOrder: domain.NameSurnameContactNameDisplayOrder,
		AmojoID:                 "f3c6340d-410e-4ad1-9f7e-c5e663599909",
		Version:                 11,
		Links:                   &domain.Links{Self: &domain.Link{Href: "https://example.amocrm.ru/api/v4/account"}},
		Embedded: &domain.AccountEmbedded{
			AmojoRights: &domain.AmojoRights{CanDirect: true, CanCreateGroups: true},
			UsersGroups: []*domain.UserGroup{{ID: 1, Name: "Отдел продаж"}},
			TaskTypes: []*domain.TaskTypeInfo{
				{ID: 1, Name: "Связаться", Code: "FOLLOW_UP"},
				{ID: 2, Name: "Встреча", Code: "MEETING"},
			},
			DatetimeSettings: &domain.DatetimeSettings{
				DatePattern:      "d.m.Y H:i",
				ShortDatePattern: "d.m.Y",
				ShortTimePattern: "H:i",
				DateFormat:       "d.m.Y",
				TimeFormat:       "H:i:s",
				Timezone:         "Europe/Moscow",
				TimezoneOffset:   "+03:00",
			},
			EntityNames: &domain.EntityNames{Leads: map[string]*domain.LanguageEntityNames{
				"ru": {
					Gender:       "m",
					PluralForm:   &domain.EntityForm{Dative: "клиентам", Default: "клиенты", Genitive: "клиентов", Accusative: "клиентов", Instrumental: "клиентами", Prepositional: "клиентах"},
					SingularForm: &domain.EntityForm{Dative: "клиенту", Default: "клиент", Genitive: "клиента", Accusative: "клиента", Instrumental: "клиентом", Prepositional: "клиенте"},
				},
				"en": {Gender: "f", PluralForm: &domain.EntityForm{Default: "leads"}, SingularForm: &domain.EntityForm{Default: "lead"}},
				"es": {Gender: "m", PluralForm: &domain.EntityForm{Default: "acuerdos"}, SingularForm: &domain.EntityForm{Default: "acuerdo"}},
			}},
		},
	}

	ctx := context.Background()

	t.Run("Успешный обработка", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, sampleGetAccountResponseBody)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetAccount(ctx, sampleGetAccountRequestParams)
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

		responseGot, err := client.GetAccount(ctx, sampleGetAccountRequestParams)
		assert.EqualError(t, err, ErrEmptyResponse.Error())
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})

	t.Run("Невалидный ответ", func(t *testing.T) {
		requestParamsGot := make(url.Values)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestParamsGot = r.URL.Query()
			w.Header().Add(contentTypeHeader, successContentType)
			_, _ = io.WriteString(w, `{"id":1231414,"name":"example","subdomain":"example","created_at":1585840134,"created_by":321321,"updated_at":1589472711,"updated_by":321321,"current_user_id":581651,"country":"RU","customers_mode":"segments","is_unsorted_on":true,"is_loss_reason_enabled":true,"is_helpbot_enabled":false,"is_technical_account":false,"contact_name_display_order":1,"amojo_id":"f3c6340d-410e-4ad1-9f7e-c5e663599909","uuid":"824f3a59-6154-4edf-ba90-0b5593715d07","version":11,"_links":{"self":{"href":"https://example.amocrm.ru/api/v4/account"}},"_embedded":{"amojo_rights":{"can_direct":true,"can_create_groups":true},"users_groups":[{"id":1,"name":"Отдел продаж","uuid":null}],"task_types":[{"id":1,"name":"Связаться","color":null,"icon_id":null,"code":"FOLLOW_UP"},{"id":2,"name":"Встреча","color":null,"icon_id":null,"code":"MEETING"}],"entity_names":{"leads":{"ru":{"gender":"m","plural_form":{"dative":"клиентам","default":"клиенты","genitive":"клиентов","accusative":"клиентов","instrumental":"клиентами","prepositional":"клиентах"},"singular_form":{"dative":"клиенту","default":"клиент","genitive":"клиента","accusative":"клиента","instrumental":"клиентом","prepositional":"клиенте"}},"en":{"singular_form":{"default":"lead"},"plural_form":{"default":"leads"},"gender":"f"},"es":{"singular_form":{"default":"acuerdo"},"plural_form":{"default":"acuerdos"},"gender":"m"}}},"datetime_settings":{"date_pattern":"d.m.Y H:i","short_date_pattern":"d.m.Y","short_time_pattern":"H:i","date_format":"d.m.Y","time_format":"H:i:s","timezone":"Europe/Moscow"}}}`)
		}))

		client, err := defaultTestClientWithURL(server.URL)
		assert.NoError(t, err)

		responseGot, err := client.GetAccount(ctx, sampleGetAccountRequestParams)
		assert.EqualError(t, err, "Key: 'Account.Embedded.DatetimeSettings.TimezoneOffset' Error:Field validation for 'TimezoneOffset' failed on the 'required' tag")
		assert.Equal(t, requestParamsWant, requestParamsGot)
		assert.Empty(t, responseGot)
	})
}
