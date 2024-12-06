package helper

import (
	"encoding/json"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"html/template"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/jsonapi"
	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
	"golang.org/x/text/language"
)

func isVarType(value interface{}, targetType reflect.Type) bool {
	return reflect.TypeOf(value) == targetType
}

func IsSliceVarOfType(slice interface{}, elemType reflect.Type) bool {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return false
	}
	return t.Elem() == elemType
}

func ValidationErrorString(validationErrors validator.ValidationErrors) map[string]interface{} {
	var fName string
	localizer := locale.GetLocalizerInstance()
	errors := make(map[string]interface{})
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

func JSONAPIModel(r *echo.Response, models interface{}, status int) error {
	r.Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	r.WriteHeader(status)
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

func GetTemplate(templ string) *template.Template {
	basePath := "./resources/templates/"
	return Must(template.New(templ).ParseFiles([]string{
		basePath + "base/header.gohtml",
		basePath + "base/footer.gohtml",
		basePath + templ,
	}...))
}
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
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
	jsonString, _ := json.Marshal(data)
	err := json.Unmarshal(jsonString, &structType)
	if err != nil {
		panic(err)
	}
	return structType
}

func MakeMailTemplateData(data map[string]any) interface{} {
	_, ok := data["Locale"]
	if !ok {
		data["Locale"] = language.English.String()
	}
	cfg := env.GetConfigInstance()
	data["AppStaticLink"] = cfg.AppURL + "/static/images/"
	data["AppLink"] = cfg.AppURL
	data["AppName"] = cfg.AppName
	return MakeStruct(data)
}
