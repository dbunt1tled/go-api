package reader

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Mapper struct {
	Fields        map[string]interface{}
	Values        []interface{}
	DynamicValues map[int]string
	DateFormat    string
	DynamicFields []string
}

func NewMapper(fields map[string]interface{}, dateFormat string, dynamicFields []string) *Mapper {
	return &Mapper{
		Fields:        transformFields(fields),
		Values:        []interface{}{},
		DynamicValues: map[int]string{},
		DateFormat:    dateFormat,
		DynamicFields: dynamicFields,
	}
}

func (m *Mapper) SetColumns(values []interface{}) bool {
	status := false
	for i, value := range values {
		strValue, ok := value.(string)
		if !ok {
			continue
		}
		strValue = strings.ToLower(strValue)

		if fieldName, found := m.Fields[strValue]; found {
			m.Fields[fieldName.(string)] = i
			status = true
		} else if len(m.DynamicFields) > 0 {
			dynamicMatch := regexp.MustCompile(strings.Join(m.DynamicFields, "|"))
			if dynamicMatch.MatchString(strValue) {
				m.DynamicValues[i] = strValue
			}
		}
	}
	return status
}

func (m *Mapper) GetValue(key string) (string, error) {
	index, ok := m.Fields[key].(int)
	if !ok || index >= len(m.Values) {
		return "", nil
	}

	switch value := m.Values[index].(type) {
	case string:
		return m.sanitizeString(value), nil
	case time.Time:
		return value.Format(time.RFC3339), nil
	default:
		return fmt.Sprintf("%v", value), nil
	}
}

func (m *Mapper) sanitizeString(input string) string {
	spaceReplacer := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(spaceReplacer.ReplaceAllString(input, " "))
}

func transformFields(fields map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range fields {
		switch v := value.(type) {
		case []string:
			for _, subValue := range v {
				result[strings.ToLower(subValue)] = key
			}
		case string:
			result[strings.ToLower(v)] = key
		}
	}
	return result
}

func (m *Mapper) ReturnInteger(value interface{}, field string) (int, error) {
	switch v := value.(type) {
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			return num, nil
		}
	case float64:
		return int(v), nil
	}
	return 0, fmt.Errorf("%s is not an integer", field)
}

// ReturnFloat возвращает число с плавающей точкой или ошибку
func (m *Mapper) ReturnFloat(value interface{}, field string) (float64, error) {
	switch v := value.(type) {
	case string:
		v = strings.ReplaceAll(v, ",", ".")
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, nil
		}
	case float64:
		return v, nil
	}
	return 0, fmt.Errorf("%s is not a float", field)
}

// ReturnDate парсит дату или возвращает ошибку
func (m *Mapper) ReturnDate(value interface{}, field string) (time.Time, error) {
	switch v := value.(type) {
	case string:
		if date, err := time.Parse(m.DateFormat, v); err == nil {
			return date, nil
		}
	case float64:
		// Пример для Excel timestamp
		startDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
		return startDate.Add(time.Duration(v) * 24 * time.Hour), nil
	case time.Time:
		return v, nil
	}
	return time.Time{}, fmt.Errorf("%s is not a date", field)
}
