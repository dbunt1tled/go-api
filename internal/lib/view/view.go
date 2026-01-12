package view

import (
	"fmt"
	"html/template"
	"time"

	"github.com/dbunt1tled/go-api/internal/config"
	"golang.org/x/text/language"
)

func MakeTemplateData(data map[string]any) map[string]any {
	_, ok := data["Locale"]
	if !ok {
		data["Locale"] = language.English.String()
	}
	data["AppStaticLink"] = fmt.Sprintf("%s/%s", config.Get().URL, config.Get().Static.URL)
	data["AppStaticImageLink"] = fmt.Sprintf("%s/%s/%s", config.Get().URL, config.Get().Static.URL, "images")
	data["AppLink"] = config.Get().URL
	data["AppName"] = config.Get().Name
	data["Year"] = time.Now().UTC().Year()
	return data
}

func GetTemplate(templ string) (*template.Template, error) {
	basePath := "./resources/templates/"
	return template.New(templ).ParseFiles([]string{
		basePath + "base/header.gohtml",
		basePath + "base/footer.gohtml",
		basePath + "base/layout/l_header.gohtml",
		basePath + "base/layout/l_footer.gohtml",
		basePath + templ,
	}...)
}
