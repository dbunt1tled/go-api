package builder

import (
	"database/sql"
	"fmt"
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
	NotEqual ConditionType = "!="
)

type FilterCondition struct {
	Field string
	Type  ConditionType
	Value interface{}
}

type OrderType string

const (
	Asc  OrderType = "ASC"
	Desc OrderType = "DESC"
)

type SortOrder struct {
	Field string
	Order OrderType
}

func GetDB() *sql.DB {
	return storage.GetInstance().GetDB()
}

func BuildSQLQuery(tableName string, filters []FilterCondition, sorts []SortOrder, addLimit bool) (string, []interface{}) {
	var whereClauses []string
	var orderClauses []string
	var args []interface{}

	for _, filter := range filters {
		switch filter.Type {
		case Equal, NotEqual, Like:
			whereClauses = append(whereClauses, fmt.Sprintf("%s %s ?", filter.Field, filter.Type))
			args = append(args, filter.Value)
		case NotNull, IsNull:
			whereClauses = append(whereClauses, fmt.Sprintf("%s %s", filter.Field, filter.Type))
		case In:
			placeholders := strings.Repeat("?, ", len(filter.Value.([]interface{})))
			placeholders = strings.TrimRight(placeholders, ", ")
			whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", filter.Field, placeholders))
			args = append(args, filter.Value.([]interface{})...)
		}
	}

	for _, sort := range sorts {
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", sort.Field, sort.Order))
	}

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if len(orderClauses) > 0 {
		query += " ORDER BY " + strings.Join(orderClauses, ", ")
	}

	if addLimit {
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
