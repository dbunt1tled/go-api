package helper

import (
	"bytes"
	"fmt"
	"go_echo/internal/config/env"
	"go_echo/internal/config/logger"
	"go_echo/internal/util/builder/page"
	"go_echo/internal/util/type/checker"
	jf "go_echo/internal/util/type/json"
	"html/template"
	"log/slog"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"golang.org/x/text/language"
)

func ValidationErrorString(
	ctx echo.Context,
	validationErrors validator.ValidationErrors,
) map[string]interface{} {
	var fName string
	errors := make(map[string]interface{})
	localizer := GetLocalizer(ctx)
	if localizer == nil {
		errors["localizer"] = "Localizer not setup"
		return errors
	}
	defaultMessage := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "default_validation_message",
		TemplateData: map[string]interface{}{
			"Field": "{{.Field}}",
			"Tag":   "{{.Tag}}",
		},
	})
	for _, err := range validationErrors {
		fName = localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: err.Field(),
		})
		errors[err.Field()] = localizer.MustLocalize(&i18n.LocalizeConfig{
			TemplateData: map[string]interface{}{
				"Field": fName,
				"Tag":   err.ActualTag(),
			},
			DefaultMessage: &i18n.Message{
				ID:    err.ActualTag(),
				One:   defaultMessage,
				Other: defaultMessage,
			},
		})
	}
	return errors
}

func JSONAPIModel[
	T interface{} |
		*map[string]interface{} |
		[]*map[string]interface{} |
		page.Paginate[map[string]interface{}],
](r *echo.Response, models T, status int) error {
	r.Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	r.WriteHeader(status)
	pg := checker.VarToPaginate(models)
	if pg != nil {
		m := (*pg).GetModels()
		p, e := jsonapi.Marshal(m)
		if e != nil {
			return e
		}
		payload, _ := p.(*jsonapi.ManyPayload)
		payload.Meta = &jsonapi.Meta{
			"total":       (*pg).GetTotal(),
			"perPage":     (*pg).GetPerPage(),
			"currentPage": (*pg).GetCurrentPage(),
			"totalPages":  (*pg).GetTotalPages(),
		}
		e = sonic.ConfigDefault.NewEncoder(r).Encode(payload)
		if e != nil {
			return e
		}
		return nil
	}
	e := jsonapi.MarshalPayload(r, models)
	if e != nil {
		panic(e.Error()) // TODO logging
	}
	return nil
}

func AssignFromSlice[T any](slice []T, vars ...*T) {
	for i := 0; i < len(slice) && i < len(vars); i++ {
		v := reflect.ValueOf(vars[i])
		if v.Kind() == reflect.Ptr && v.Elem().CanSet() {
			v.Elem().Set(reflect.ValueOf(slice[i]))
		}
	}
}

func GetLocalizer(c echo.Context) *i18n.Localizer {
	if localizer, ok := c.Get("localizer").(*i18n.Localizer); ok {
		return localizer
	}
	return nil
}

func GetTemplate(templ string) *template.Template {
	basePath := "./resources/templates/"
	return Must(template.New(templ).ParseFiles([]string{
		basePath + "base/header.gohtml",
		basePath + "base/footer.gohtml",
		basePath + "base/layout/l_header.gohtml",
		basePath + "base/layout/l_footer.gohtml",
		basePath + templ,
	}...))
}
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func MustErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MakeStruct(data map[string]any) interface{} {
	sc := dynamicstruct.NewStruct()
	for k, v := range data {
		sc.AddField(
			strcase.ToCamel(k),
			v,
			"",
		)
	}
	structType := sc.Build().New()
	jsonString, _ := sonic.Marshal(data)
	err := sonic.Unmarshal(jsonString, &structType)
	if err != nil {
		panic(err)
	}
	return structType
}

func RuntimeStatistics(startTime time.Time, showName bool) string {
	pc, _, _, _ := runtime.Caller(1)
	name := ""
	if showName {
		funcObj := runtime.FuncForPC(pc)
		runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
		name = runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	}
	return strings.TrimSpace(fmt.Sprintf(
		"%s processed %s (%s)",
		name,
		time.Since(startTime).Round(time.Second).String(),
		MemoryUsage(),
	))
}
func MemoryUsage() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return fmt.Sprintf(
		"TotalAlloc: %v MB, Sys: %v MB",
		memStats.TotalAlloc/1024/1024, //nolint:mnd // Convert to MB
		memStats.Sys/1024/1024,        //nolint:mnd // Convert to MB
	)
}

func SubStr(stack string, needle string) string {
	index := strings.Index(stack, needle)
	if index != -1 {
		return stack[:index]
	}
	return stack
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	newMap := make(map[string]interface{})
	data, err := sonic.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = sonic.Unmarshal(data, &newMap)
	return newMap, err
}
func StructToJSONField(obj interface{}) (jf.JsonField, error) {
	m, err := StructToMap(obj)
	if err != nil {
		return nil, err
	}
	return m, err
}

func MapToByte(obj map[string]interface{}) ([]byte, error) {
	data, err := sonic.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetVarValue(v any) any {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		return val.Elem().Interface()
	}
	return v
}

func RequestID(c echo.Context) string {
	id := c.Request().Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = c.Response().Header().Get(echo.HeaderXRequestID)
	}
	return id
}

func AnyToBytesBuffer(i interface{}) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := sonic.ConfigDefault.NewEncoder(buf).Encode(i)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func IfThenElse[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func AnyToString(a any) string {
	switch value := a.(type) {
	case nil:
		return ""
	case string:
		return value
	case int:
		return strconv.Itoa(value)
	case int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(value).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(value).Uint(), 10)
	case []byte:
		return string(value)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case time.Time:
		return value.Format(time.RFC3339)
	default:
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			elem := reflect.ValueOf(value).Elem()
			if !elem.IsValid() {
				return ""
			}
			return AnyToString(elem.Interface())
		}
		return fmt.Sprintf("%v", value)
	}
}

func IsA(a any, typeToAssert reflect.Kind) bool {
	return typeToAssert == reflect.ValueOf(a).Kind()
}

func MakeTemplateData(data map[string]any) interface{} {
	_, ok := data["Locale"]
	if !ok {
		data["Locale"] = language.English.String()
	}
	cfg := env.GetConfigInstance()
	data["AppStaticImageLink"] = cfg.Static.URL
	data["AppLink"] = cfg.AppURL
	data["AppName"] = cfg.AppName
	data["Year"] = time.Now().UTC().Year()
	return MakeStruct(data)
}

func ToPointer[T any](t T) *T {
	return &t
}

func ToPointerOrNil[T comparable](t T) *T {
	if z, ok := any(t).(interface{ IsZero() bool }); ok {
		if z.IsZero() {
			return nil
		}
		return &t
	}

	var v T
	if t == v {
		return nil
	}
	return &t
}

func ValueFromPointer[T any](t *T) T {
	if t == nil {
		var v T
		return v
	}
	return *t
}

func GracefulShutdown(log *logger.AppLogger, ops ...func() error) {
	for _, op := range ops {
		if err := op(); err != nil {
			log.Error(
				"(ツ)_/¯ Graceful Shutdown op failed",
				slog.Any("error", err),
			)
			panic(err)
		}
	}
}

func Ter[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}

func IsNil(x interface{}) bool {
	if x == nil {
		return true
	}

	return reflect.ValueOf(x).IsNil()
}
