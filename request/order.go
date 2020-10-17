package request

import (
	"fmt"
	"net/url"
)

type (
	OrderBy string

	OrderMethod string

	Order struct {
		By     OrderBy     `validate:"required,oneof=id created_at updated_at"`
		Method OrderMethod `validate:"required,oneof=asc desc"`
	}
)

const (
	AscendingOrderMethod  OrderMethod = "asc"
	DescendingOrderMethod OrderMethod = "desc"

	IDRequestOrderBy        OrderBy = "id"
	CreatedAtRequestOrderBy OrderBy = "created_at"
	UpdatedAtRequestOrderBy OrderBy = "updated_at"
)

func (o *Order) AppendToQuery(params url.Values) {
	params.Add(fmt.Sprintf("order[%s]", string(o.By)), string(o.Method))
}
