package types

type SORT_BY string

const SORT_BY_DATE SORT_BY = "created_at"
const SORT_BY_PRICE SORT_BY = "price"

type ORDER_BY string

const ORDER_BY_ASC ORDER_BY = "asc"
const ORDER_BY_DESC ORDER_BY = "desc"

type GetAdParams struct {
	Page     int
	MinPrice int
	MaxPrice int
	SortBy   SORT_BY
	OrderBy  ORDER_BY
}
