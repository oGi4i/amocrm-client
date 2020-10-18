package client

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/ogi4i/amocrm-client/domain"
)

type (
	AccountGetRequestWith string

	AccountGetRequestParams struct {
		With []AccountGetRequestWith `validate:"omitempty,dive,oneof=amojo_id amojo_rights user_groups task_types version entity_names datetime_settings"`
	}
)

const (
	AmojoIDAccountGetRequestWith          AccountGetRequestWith = "amojo_id"          // Добавляет в ответ ID аккаунта в сервисе чатов
	AmojoRightsAccountGetRequestWith      AccountGetRequestWith = "amojo_rights"      // Добавляет в ответ информацию о доступности функционала создания групповых и использования директ чатов пользователями
	UserGroupsAccountGetRequestWith       AccountGetRequestWith = "user_groups"       // Добавляет в ответ информацию о доступных группах пользователей аккаунта
	TaskTypesAccountGetRequestWith        AccountGetRequestWith = "task_types"        // Добавляет в ответ информацию о доступных типах задач в аккаунта
	VersionAccountGetRequestWith          AccountGetRequestWith = "version"           // Добавляет в ответ информацию о текущей версии amoCRM
	EntityNamesAccountGetRequestWith      AccountGetRequestWith = "entity_names"      // Добавляет в ответ названия сущностей с их переводами и форматами чисел
	DatetimeSettingsAccountGetRequestWith AccountGetRequestWith = "datetime_settings" // Добавляет в ответ информацию о текущих настройках форматов даты и времени аккаунта
)

func (a AccountGetRequestWith) String() string {
	return string(a)
}

func (c *Client) GetAccount(ctx context.Context, reqParams *AccountGetRequestParams) (*domain.Account, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	params := make(url.Values)
	if reqParams.With != nil {
		params.Add("with", joinAccountRequestWith(reqParams.With))
	}

	body, err := c.doGet(ctx, c.baseURL+accountURI, params)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, ErrEmptyResponse
	}

	response := new(domain.Account)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}

	if err := c.validator.Struct(response); err != nil {
		return nil, err
	}

	return response, nil
}

func joinAccountRequestWith(with []AccountGetRequestWith) string {
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