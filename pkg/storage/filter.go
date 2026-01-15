package storage

type Oper string

const (
	OpEqual              Oper = "="
	OpNotEqual           Oper = "!="
	OpGreaterThan        Oper = ">"
	OpGreaterThanOrEqual Oper = ">="
	OpLessThan           Oper = "<"
	OpLessThanOrEqual    Oper = "<="
	OpLike               Oper = "LIKE"
	OpILike              Oper = "ILIKE"
	OpIn                 Oper = "IN"
	OpNotIn              Oper = "NOT IN"
	OpIsNull             Oper = "IS NULL"
	OpIsNotNull          Oper = "IS NOT NULL"
	OpContains           Oper = "CONTAINS"      // array (@>)
	OpContainedBy        Oper = "CONTAINED_BY"  // array (<@)
	OpOverlaps           Oper = "OVERLAPS"      // array (&&)
	OpJsonContains       Oper = "JSON_CONTAINS" // JSONB (@>)
	OpJsonExists         Oper = "JSON_EXISTS"   // JSONB (?)
)

type Rule struct {
	Field     string
	Operation Oper
	Value     interface{}
}

type Filter struct {
	Rules   []Rule
	OrderBy []Sort
	Limit   int
	Offset  int
}

type PaginationInfo interface {
	GetTotal() int
	GetPage() int
	GetPerPage() int
	GetTotalPages() int
	GetHasNext() bool
	GetHasPrev() bool
}

type Paginator[T any] struct {
	Items      []T  `json:"items"`
	Total      int  `json:"total"`
	Page       int  `json:"page"`
	PerPage    int  `json:"perPage"`
	TotalPages int  `json:"totalPages"`
	HasNext    bool `json:"next"`
	HasPrev    bool `json:"prev"`
}

func (p *Paginator[T]) GetTotal() int {
	return p.Total
}

func (p *Paginator[T]) GetPage() int {
	return p.Page
}

func (p *Paginator[T]) GetPerPage() int {
	return p.PerPage
}

func (p *Paginator[T]) GetTotalPages() int {
	return p.TotalPages
}

func (p *Paginator[T]) GetHasNext() bool {
	return p.HasNext
}

func (p *Paginator[T]) GetHasPrev() bool {
	return p.HasPrev
}

type Sort struct {
	Field      string `json:"field"`
	Descending bool   `json:"descending"`
}

type QueryOption func(*queryConfig)

type queryConfig struct {
	rules   []Rule
	orderBy []Sort
	limit   int
	offset  int
}

func WithFilter(rules ...Rule) QueryOption {
	return func(cfg *queryConfig) {
		cfg.rules = append(cfg.rules, rules...)
	}
}

func WithSort(field string, sort string) QueryOption {
	return func(cfg *queryConfig) {
		cfg.orderBy = append(cfg.orderBy, Sort{
			Field:      field,
			Descending: sort == "desc",
		})
	}
}

func WithSortAsc(field string) QueryOption {
	return func(cfg *queryConfig) {
		cfg.orderBy = append(cfg.orderBy, Sort{
			Field:      field,
			Descending: false,
		})
	}
}

func WithSortDesc(field string) QueryOption {
	return func(cfg *queryConfig) {
		cfg.orderBy = append(cfg.orderBy, Sort{
			Field:      field,
			Descending: true,
		})
	}
}

func WithLimit(limit int) QueryOption {
	limit = NormalizePerPage(limit)
	return func(cfg *queryConfig) {
		cfg.limit = limit
	}
}

func WithOffset(offset int) QueryOption {
	offset = NormalizeOffset(offset)
	return func(cfg *queryConfig) {
		cfg.offset = offset
	}
}

func WithLimitOffset(limit int, offset int) QueryOption {
	return func(cfg *queryConfig) {
		cfg.limit = limit
		cfg.offset = offset
	}
}

func WithPagination(page int, perPage int) QueryOption {
	return func(cfg *queryConfig) {
		page, perPage = NormalizePagination(page, perPage)
		cfg.limit = perPage
		cfg.offset = (page - 1) * perPage
	}
}

func NewRule(field string, oper Oper, value interface{}) Rule {
	return Rule{
		Field:     field,
		Operation: oper,
		Value:     value,
	}
}

func NormalizePagination(page int, perPage int) (int, int) {
	page = NormalizePage(page)
	perPage = NormalizePerPage(perPage)
	return page, perPage
}

func NormalizePerPage(perPage int) int {
	if perPage < 1 {
		perPage = 10
	}

	if perPage > 5000 { //nolint:mnd // max limit
		perPage = 5000
	}
	return perPage
}

func NormalizePage(page int) int {
	if page < 1 {
		page = 1
	}

	if page > 1000 { //nolint:mnd // max limit
		page = 1000
	}
	return page
}

func NormalizeOffset(offset int) int {
	if offset < 0 {
		offset = 0
	}
	return offset
}
