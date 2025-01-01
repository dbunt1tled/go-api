package data

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var specialRegex = regexp.MustCompile(`_#([a-zA-Z\s\d_]+)#_`)

type Mapper struct {
	Fields        map[string]interface{}
	MappedFields  map[string]int
	SpecialFields map[string]string
	Values        []interface{}
	DynamicValues map[int]string
	DynamicFields []string
}

func NewMapper(fields map[string]interface{}, dynamicFields []string) *Mapper {
	return &Mapper{
		Fields:        transformFields(fields),
		SpecialFields: getSpecialFields(fields), // TODO: fix, original code was hardcoded, currently it is reference to field
		MappedFields:  map[string]int{},
		Values:        []interface{}{},
		DynamicValues: map[int]string{},
		DynamicFields: dynamicFields,
	}
}

func (m *Mapper) SetColumns(values []string) bool {
	status := false
	for i, strValue := range values {
		strValue = strings.ToLower(strValue)

		if fieldName, found := m.Fields[strValue]; found {
			m.MappedFields[fieldName.(string)] = i //nolint:errcheck
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
	index, ok := m.MappedFields[key]
	if !ok || index >= len(m.Values) {
		k, o := m.SpecialFields[key]
		if !o {
			return "", nil
		}
		k = strings.TrimPrefix(k, "_#")
		k = strings.TrimSuffix(k, "#_")
		index, ok = m.MappedFields[k]
		if !ok || index >= len(m.Values) {
			return "", nil
		}
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

func getSpecialFields(fields map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range fields {
		switch v := value.(type) {
		case []string:
			for _, subValue := range v {
				if isSpecialField(subValue) {
					result[strings.ToLower(key)] = strings.ToLower(subValue)
				}
			}
		case string:
			if isSpecialField(v) {
				result[strings.ToLower(key)] = strings.ToLower(v)
			}
		}
	}
	return result
}

func isSpecialField(field string) bool {
	return specialRegex.MatchString(field)
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

func (m *Mapper) ReturnDate(value interface{}, field string, dateFormat string) (time.Time, error) {
	switch v := value.(type) {
	case string:
		if date, err := time.Parse(dateFormat, v); err == nil {
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

func SliceToSliceInterface[T any](value []T) []interface{} {
	interfaces := make([]interface{}, len(value))
	for i, v := range value {
		interfaces[i] = v
	}
	return interfaces
}
