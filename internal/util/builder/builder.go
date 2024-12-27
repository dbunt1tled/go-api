package builder

import (
	"database/sql"
	"fmt"
	"go_echo/internal/dto"
	"go_echo/internal/storage"
	"strings"
)

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

type Pagination[T map[string]interface{}] struct {
	Total       int
	CurrentPage int
	PerPage     int
	TotalPages  int
	Models      []*T
}

func GetDB() *sql.DB {
	return storage.GetInstance().GetDB()
}

func BuildSQLQuery(
	tableName string,
	filters *[]FilterCondition,
	sorts *[]SortOrder,
	addLimit1 bool,
) (string, []interface{}) {
	var whereClauses []string
	var orderClauses []string
	var args []interface{}

	if filters != nil && len(*filters) > 0 {
		for _, filter := range *filters {
			if filter.Value == nil {
				continue
			}
			switch filter.Type {
			case Equal, NotEqual, Like:
				whereClauses = append(whereClauses, fmt.Sprintf("%s %s ?", filter.Field, filter.Type))
				args = append(args, filter.Value)
			case NotNull, IsNull:
				whereClauses = append(whereClauses, fmt.Sprintf("%s %s", filter.Field, filter.Type))
			case In, NotIn:
				placeholders := strings.Repeat("?, ", len(filter.Value.([]interface{})))
				placeholders = strings.TrimRight(placeholders, ", ")
				whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", filter.Field, placeholders))
				args = append(args, filter.Value.([]interface{})...)
			}
		}
	}
	if sorts != nil && len(*sorts) > 0 {
		for _, sort := range *sorts {
			orderClauses = append(orderClauses, fmt.Sprintf("%s %s", sort.Field, sort.Order))
		}
	}

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if len(orderClauses) > 0 {
		query += " ORDER BY " + strings.Join(orderClauses, ", ")
	}

	if addLimit1 {
		query += " LIMIT 1"
	}

	return query, args
}

func ValidateFilter(filters []FilterCondition, validFields map[string]bool) error {
	for _, filter := range filters {
		if !validFields[filter.Field] {
			return fmt.Errorf("invalid field: %s", filter.Field)
		}
	}
	return nil
}

func Eq[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  Equal,
		Value: value,
	}
}
func Inc[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  In,
		Value: value,
	}
}
func NInc[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  In,
		Value: value,
	}
}
func Lk[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  Like,
		Value: value,
	}
}
func NotEq[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  NotEqual,
		Value: value,
	}
}
func IsNl[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  IsNull,
		Value: value,
	}
}
func NotNl[T any](field string, value T) FilterCondition {
	return FilterCondition{
		Field: field,
		Type:  NotNull,
		Value: value,
	}
}

func GetSortOrder(sort dto.Sorting) *[]SortOrder {
	res := make([]SortOrder, 1)
	if sort.Field == nil {
		return nil
	}
	if sort.Order == nil {
		sort.Order = new(string)
		*sort.Order = string(Asc)
	}
	order := SortOrder{
		Field: *sort.Field,
		Order: OrderType(strings.ToLower(*sort.Order)),
	}
	res[0] = order
	return &res
}
