package builder

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dbunt1tled/go-api/internal/dto"
	e "github.com/dbunt1tled/go-api/internal/err"
	"github.com/dbunt1tled/go-api/internal/storage"
	"github.com/dbunt1tled/go-api/internal/util/builder/page"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func GetDB() *sql.DB {
	return storage.GetInstance().GetDB()
}

func BuildSQLQuery(
	tableName string,
	filters *[]page.FilterCondition,
	sorts *[]page.SortOrder,
	limit int,
	offset int,
	asCount bool,
) (string, []any) {
	var whereClauses []string
	var orderClauses []string
	var args []any

	if filters != nil && len(*filters) > 0 {
		for _, filter := range *filters {
			if filter.Value == nil {
				continue
			}
			switch filter.Type {
			case page.Equal, page.NotEqual, page.Like:
				whereClauses = append(whereClauses, fmt.Sprintf("%s %s ?", filter.Field, filter.Type))
				args = append(args, filter.Value)
			case page.NotNull, page.IsNull:
				whereClauses = append(whereClauses, fmt.Sprintf("%s %s", filter.Field, filter.Type))
			case page.In, page.NotIn:
				rt := reflect.TypeOf(filter.Value)
				switch rt.Kind() {
				case reflect.Slice, reflect.Array:
					v := reflect.ValueOf(filter.Value)
					l := v.Len()
					placeholders := strings.Repeat("?, ", l)
					placeholders = strings.TrimRight(placeholders, ", ")
					whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", filter.Field, placeholders))
					for i := range l {
						args = append(args, v.Index(i).Interface())
					}
				default:
					panic("Builder.FilterData invalid type IN, NotIn: " + rt.String())
				}
			}
		}
	}
	if sorts != nil && len(*sorts) > 0 {
		for _, sort := range *sorts {
			orderClauses = append(orderClauses, fmt.Sprintf("%s %s", sort.Field, sort.Order))
		}
	}
	var query strings.Builder
	query.WriteString("SELECT ")
	if asCount {
		query.WriteString("COUNT(*)")
	} else {
		query.WriteString("*")
	}
	query.WriteString(" FROM ")
	query.WriteString(tableName)
	if len(whereClauses) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(strings.Join(whereClauses, " AND "))
	}

	if len(orderClauses) > 0 {
		query.WriteString(" ORDER BY ")
		query.WriteString(strings.Join(orderClauses, ", "))
	}

	if limit > 0 {
		query.WriteString(" LIMIT ")
		query.WriteString(strconv.Itoa(limit))

		if offset > 0 {
			query.WriteString(" OFFSET ")
			query.WriteString(strconv.Itoa(offset))
		}
	}

	return query.String(), args
}

func ValidateFilter(filters []page.FilterCondition, validFields map[string]bool) error {
	for _, filter := range filters {
		if !validFields[filter.Field] {
			return fmt.Errorf("invalid field: %s", filter.Field)
		}
	}
	return nil
}

func Eq[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.Equal,
		Value: value,
	}
}
func Inc[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.In,
		Value: value,
	}
}
func NInc[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.In,
		Value: value,
	}
}
func Lk[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.Like,
		Value: value,
	}
}
func NotEq[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.NotEqual,
		Value: value,
	}
}
func IsNl[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.IsNull,
		Value: value,
	}
}
func NotNl[T any](field string, value T) page.FilterCondition {
	return page.FilterCondition{
		Field: field,
		Type:  page.NotNull,
		Value: value,
	}
}

func GetSortOrder(sort dto.Sorting) *[]page.SortOrder {
	res := make([]page.SortOrder, 1)
	if sort.Field == nil {
		return nil
	}
	if sort.Order == nil {
		sort.Order = new(string)
		*sort.Order = string(page.Asc)
	}
	order := page.SortOrder{
		Field: *sort.Field,
		Order: page.OrderType(strings.ToLower(*sort.Order)),
	}
	res[0] = order
	return &res
}

func GetPagination(p dto.PaginationQuery) *page.Pagination {
	if p.Page == nil || p.Limit == nil {
		return nil
	}
	res := page.Pagination{
		Page:    *p.Page,
		PerPage: *p.Limit,
	}

	return &res
}

func Count(ctx context.Context, table string, filter *[]page.FilterCondition) (int, error) {
	var (
		cnt int
		err error
		res *sql.Row
	)
	query, args := BuildSQLQuery(table, filter, nil, 1, 0, true)
	res = GetDB().QueryRowContext(ctx, query, args...)
	err = res.Scan(&cnt)
	if err != nil {
		return 0, errors.Wrap(err, table+" count cast error")
	}
	return cnt, nil
}

func ByID[T any](
	ctx context.Context,
	table string,
	id int64,
	mapper func(res *sql.Row) (*T, error),
) (*T, error) {
	qb := strings.Builder{}
	qb.WriteString("SELECT * FROM ")
	qb.WriteString(table)
	qb.WriteString(" WHERE id = ? LIMIT 1")
	return ExecuteSQLRow(ctx, qb.String(), mapper, id)
}

func List[T any](
	ctx context.Context,
	table string,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
	mapper func(res *sql.Rows) (*T, error),
	paginator *page.Pagination,
) ([]*T, error) {
	limit := 0
	offset := 0
	if paginator != nil {
		limit = paginator.PerPage
		offset = (paginator.Page - 1) * paginator.PerPage
	}
	query, args := BuildSQLQuery(table, filter, sorts, limit, offset, false)
	return ExecuteSQLQuery(ctx, query, mapper, args...)
}

func One[T any](
	ctx context.Context,
	table string,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
	mapper func(res *sql.Row) (*T, error),
) (*T, error) {
	query, args := BuildSQLQuery(table, filter, sorts, 1, 0, false)
	return ExecuteSQLRow(ctx, query, mapper, args...)
}

func ExecuteSQLQuery[T any](
	ctx context.Context,
	query string,
	mapper func(res *sql.Rows) (*T, error),
	args ...interface{},
) ([]*T, error) {
	var (
		u   *T
		res *sql.Rows
		smt *sql.Stmt
		err error
	)
	models := make([]*T, 0)
	smt, err = GetDB().PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "query prepare error")
	}
	defer smt.Close()
	res, err = smt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, " query error")
	}
	defer res.Close()
	for res.Next() {
		u, err = mapper(res)
		if err != nil {
			return nil, errors.Wrap(err, "query cast error")
		}
		models = append(models, u)
	}
	return models, nil
}

func ExecuteSQLRow[T any](
	ctx context.Context,
	query string,
	mapper func(res *sql.Row) (*T, error),
	args ...interface{},
) (*T, error) {
	var (
		res *sql.Row
	)
	res = GetDB().QueryRowContext(ctx, query, args...)
	return mapper(res)
}

func Paginator[T any](
	ctx context.Context,
	table string,
	filter *[]page.FilterCondition,
	sorts *[]page.SortOrder,
	paginator *page.Pagination,
	mapper func(res *sql.Rows) (*T, error),
) (page.Paginate[T], error) {
	var (
		wg             sync.WaitGroup
		rows           []*T
		count          int
		receivedErrors []error
	)
	resErr := errors.New("Paginate error")
	rowsChan := make(chan []*T, 1)
	countChan := make(chan int)
	errChan := make(chan error, 2)
	rStop := false
	cStop := false
	eStop := false

	wg.Add(2) //nolint:nolintlint,mnd 2 requests list and count
	fCount := func(ctx context.Context, filter *[]page.FilterCondition, wg *sync.WaitGroup, countChan chan int) {
		defer wg.Done()
		c, err := Count(ctx, table, filter)
		if err != nil {
			errChan <- err
			close(countChan)
			return
		}
		countChan <- c
		close(countChan)
	}
	fList := func(
		ctx context.Context,
		filter *[]page.FilterCondition,
		sorts *[]page.SortOrder,
		paginator *page.Pagination,
		wg *sync.WaitGroup,
		rowsChan chan []*T,
	) {
		defer wg.Done()
		res, err := List(ctx, table, filter, sorts, mapper, paginator)
		if err != nil {
			errChan <- err
			close(rowsChan)
			return
		}
		rowsChan <- res
		close(rowsChan)
	}

	go fCount(ctx, filter, &wg, countChan)
	go fList(ctx, filter, sorts, paginator, &wg, rowsChan)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for {
		if (rStop && cStop) || eStop {
			break
		}
		select {
		case r := <-rowsChan:
			if rStop {
				break
			}
			rows = r
			rStop = true
		case c := <-countChan:
			if cStop {
				break
			}
			count = c
			cStop = true
		case e := <-errChan:
			if eStop {
				break
			}
			if e != nil {
				receivedErrors = append(receivedErrors, e)
				eStop = true
			}
		}
	}
	if len(receivedErrors) > 0 {
		for _, err := range receivedErrors {
			resErr = errors.Wrap(resErr, err.Error())
		}
		return page.Paginate[T]{}, resErr
	}

	result := page.Paginate[T]{
		Total:       count,
		CurrentPage: paginator.Page,
		PerPage:     paginator.PerPage,
		TotalPages:  int(math.Ceil(float64(count) / float64(paginator.PerPage))),
		Models:      rows,
	}

	return result, nil
}

func ScanStructRows[T any](st T, rows *sql.Rows) (*T, error) {
	s := reflect.ValueOf(&st).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := range numCols {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}

	err := rows.Scan(columns...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.NotFoundErr
		}
		return nil, err
	}
	return &st, nil
}

func ScanStructRow[T any](st T, rows *sql.Row) (*T, error) {
	s := reflect.ValueOf(&st).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)
	for i := range numCols {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}

	err := rows.Scan(columns...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.NotFoundErr
		}
		return nil, err
	}
	return &st, nil
}
