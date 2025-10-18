package middlewares

import (
	"github.com/dbunt1tled/go-api/internal/config/locale"

	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func Language() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var l language.Tag
			var err error
			l = language.English
			lang := c.FormValue("lang")
			if lang == "" {
				lang = c.Request().Header.Get("Accept-Language")
			}
			if lang != "" {
				l, err = language.Parse(lang)
				if err != nil {
					l = language.English
				}
			}
			localizerInstance := i18n.NewLocalizer(locale.GetLocaleBundleInstance(), l.String(), l.String())
			c.Set("localizer", localizerInstance)

			return next(c)
		}
	}
}
