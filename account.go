package amocrm

import (
	"context"
	"encoding/json"
	"strings"
)

type (
	AccountWithType string

	AccountRequestParams struct {
		With []AccountWithType `validate:"omitempty,dive,oneof=custom_fields users messenger notifications pipelines groups note_types task_types"`
	}

	AccountResponse struct {
		ID             int    `json:"id,omitempty" validate:"required,omitempty"`
		Name           string `json:"name" validate:"required"`
		Subdomain      string `json:"subdomain" validate:"required"`
		Currency       string `json:"currency" validate:"required"`
		Timezone       string `json:"timezone" validate:"required"`
		TimezoneOffset string `json:"timezone_offset" validate:"required"`
		Language       string `json:"language" validate:"required"`
		DatePattern    struct {
			Date     string `json:"date" validate:"required"`
			Time     string `json:"time" validate:"required"`
			DateTime string `json:"date_time" validate:"required"`
			TimeFull string `json:"time_full" validate:"required"`
		} `json:"date_pattern" validate:"required"`
		CurrentUser int `json:"current_user" validate:"required"`
		Embedded    struct {
			Users        map[string]*User `json:"users" validate:"omitempty,dive,required"`
			CustomFields struct {
				Contacts  map[string]*CustomFieldInfo `json:"contacts" validate:"omitempty,dive,required"`
				Leads     map[string]*CustomFieldInfo `json:"leads,omitempty" validate:"omitempty,dive,required"`
				Companies map[string]*CustomFieldInfo `json:"companies,omitempty" validate:"omitempty,dive,required"`
				Customers []interface{}               `json:"customers,omitempty" validate:"omitempty,dive,required"`
			} `json:"custom_fields" validate:"omitempty"`
			NoteTypes map[string]*NoteType `json:"note_types" validate:"omitempty,dive,required"`
			Groups    map[string]*Group    `json:"groups" validate:"omitempty,dive,required"`
			TaskTypes map[string]*TaskType `json:"task_types" validate:"omitempty,dive,required"`
			Pipelines map[string]*Pipeline `json:"pipelines" validate:"omitempty,dive,required"`
		} `json:"_embedded" validate:"omitempty"`
	}

	AuthAccount struct {
		ID        int    `json:"id" validate:"required"`
		Name      string `json:"name" validate:"required"`
		Subdomain string `json:"subdomain" validate:"required"`
		Language  string `json:"language" validate:"required"`
		Timezone  string `json:"timezone" validate:"required"`
	}
)

const (
	AccountWithCustomFields  AccountWithType = "custom_fields"
	AccountWithUsers         AccountWithType = "users"
	AccountWithMessenger     AccountWithType = "messenger"
	AccountWithNotifications AccountWithType = "notifications"
	AccountWithPipelines     AccountWithType = "pipelines"
	AccountWithGroups        AccountWithType = "groups"
	AccountWithNoteTypes     AccountWithType = "note_types"
	AccountWithTaskTypes     AccountWithType = "task_types"
)

func (c *Client) GetAccount(ctx context.Context, reqParams *AccountRequestParams) (*AccountResponse, error) {
	if err := c.validator.Struct(reqParams); err != nil {
		return nil, err
	}

	body, err := c.doGet(ctx, c.baseURL+accountURI, map[string]string{"with": joinAccountWithTypes(reqParams.With, ",")})
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	accountResponse := new(AccountResponse)
	err = json.Unmarshal(body, accountResponse)
	if err != nil {
		amoError := new(AmoError)
		err = json.Unmarshal(body, amoError)
		if err != nil {
			return nil, err
		}

		return nil, amoError
	}

	if err := c.validator.Struct(accountResponse); err != nil {
		return nil, err
	}

	return accountResponse, nil
}

func joinAccountWithTypes(types []AccountWithType, sep string) string {
	if len(types) == 0 {
		return ""
	}

	out := new(strings.Builder)

	for i, t := range types {
		if i != len(types)-1 {
			out.WriteString(string(t) + sep)
		} else {
			out.WriteString(string(t))
		}
	}

	return out.String()
}
