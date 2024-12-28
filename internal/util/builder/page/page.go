package page

type ConditionType string

const (
	Equal    ConditionType = "="
	NotNull  ConditionType = "IS NOT NULL"
	IsNull   ConditionType = "IS NULL"
	Like     ConditionType = "LIKE"
	In       ConditionType = "IN"
	NotIn    ConditionType = "NOT IN"
	NotEqual ConditionType = "!="
)

type FilterCondition struct {
	Field string
	Type  ConditionType
	Value interface{}
}

type OrderType string

const (
	Asc  OrderType = "asc"
	Desc OrderType = "desc"
)

type SortOrder struct {
	Field string
	Order OrderType
}

type Pagination struct {
	Page    int
	PerPage int
}

type PaginateInterface interface {
	GetTotal() int
	GetCurrentPage() int
	GetPerPage() int
	GetTotalPages() int
	GetModels() any
}

type Paginate[T any] struct {
	Total       int
	CurrentPage int
	PerPage     int
	TotalPages  int
	Models      []*T
}

func (p Paginate[T]) GetTotal() int {
	return p.Total
}
func (p Paginate[T]) GetCurrentPage() int {
	return p.CurrentPage
}
func (p Paginate[T]) GetPerPage() int {
	return p.PerPage
}
func (p Paginate[T]) GetTotalPages() int {
	return p.TotalPages
}
func (p Paginate[T]) GetModels() any {
	return p.Models
}
