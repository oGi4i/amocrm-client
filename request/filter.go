package request

import (
	"fmt"
	"net/url"
)

type (
	filterType uint8

	Filter struct {
		name       string     `validate:"omitempty"`
		filterType filterType `validate:"required"`
		values     []string   `validate:"required,gt=0,dive,required"`
	}
)

const (
	simpleFilterType filterType = iota + 1
	multipleFilterType
	intervalFilterType
	statusFilterType
	simpleCustomFieldFilterType
	multipleCustomFieldFilterType
	intervalCustomFieldFilterType
)

func CreateSimpleFilter(name, value string) *Filter {
	return &Filter{
		name:       name,
		filterType: simpleFilterType,
		values:     []string{value},
	}
}

func CreateMultipleFilter(name string, values []string) *Filter {
	return &Filter{
		name:       name,
		filterType: multipleFilterType,
		values:     values,
	}
}

func CreateIntervalFilter(name, from, to string) *Filter {
	return &Filter{
		name:       name,
		filterType: intervalFilterType,
		values:     []string{from, to},
	}
}

func CreateStatusFilter(pipelineID, statusID string) *Filter {
	return &Filter{
		filterType: statusFilterType,
		values:     []string{pipelineID, statusID},
	}
}

func CreateSimpleCustomFieldFilter(fieldID, value string) *Filter {
	return &Filter{
		name:       fieldID,
		filterType: simpleCustomFieldFilterType,
		values:     []string{value},
	}
}

func CreateMultipleCustomFieldFilter(fieldID string, values []string) *Filter {
	return &Filter{
		name:       fieldID,
		filterType: multipleCustomFieldFilterType,
		values:     values,
	}
}

func CreateIntervalCustomFieldFilter(fieldID, from, to string) *Filter {
	return &Filter{
		name:       fieldID,
		filterType: intervalCustomFieldFilterType,
		values:     []string{from, to},
	}
}

func (f *Filter) AppendToQuery(params url.Values) {
	switch f.filterType {
	case simpleFilterType:
		params.Add(fmt.Sprintf("filter[%s]", f.name), f.values[0])
	case multipleFilterType:
		for _, value := range f.values {
			params.Add(fmt.Sprintf("filter[%s][0]", f.name), value)
		}
	case intervalFilterType:
		params.Add(fmt.Sprintf("filter[%s][from]", f.name), f.values[0])
		params.Add(fmt.Sprintf("filter[%s][to]", f.name), f.values[1])
	case statusFilterType:
		params.Add("filter[statuses][0][pipeline_id]", f.values[0])
		params.Add("filter[statuses][0][status_id]", f.values[1])
	case simpleCustomFieldFilterType:
		params.Add(fmt.Sprintf("filter[custom_fields_values][%s][]", f.name), f.values[0])
	case multipleCustomFieldFilterType:
		for _, value := range f.values {
			params.Add(fmt.Sprintf("filter[custom_fields_values][%s][]", f.name), value)
		}
	case intervalCustomFieldFilterType:
		params.Add(fmt.Sprintf("filter[custom_fields_values][%s][from]", f.name), f.values[0])
		params.Add(fmt.Sprintf("filter[custom_fields_values][%s][to]", f.name), f.values[1])
	}
}

func (f *Filter) IsSimpleFilter() bool {
	return f.filterType == simpleFilterType
}

func (f *Filter) IsMultipleFilter() bool {
	return f.filterType == multipleFilterType
}

func (f *Filter) IsIntervalFilter() bool {
	return f.filterType == intervalFilterType
}

func (f *Filter) IsStatusFilter() bool {
	return f.filterType == statusFilterType
}

func (f *Filter) IsSimpleCustomFieldFilter() bool {
	return f.filterType == simpleCustomFieldFilterType
}

func (f *Filter) IsMultipleCustomFieldFilter() bool {
	return f.filterType == multipleCustomFieldFilterType
}

func (f *Filter) IsIntervalCustomFieldFilter() bool {
	return f.filterType == intervalCustomFieldFilterType
}
