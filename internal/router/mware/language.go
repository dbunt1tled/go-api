package mware

import (
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func Language(locale *i18n.Localizer, bundle *i18n.Bundle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var l language.Tag
			var err error
			l = language.English
			lang := c.FormValue("lang")
			if lang != "" {
				l, err = language.Parse(lang)
				if err != nil {
					lang = c.Request().Header.Get("Accept-Language")
					l, err = language.Parse(lang)
					if err != nil {
						l = language.English
					}
				}
			}
			locale = i18n.NewLocalizer(bundle, l.String())

			return next(c)
		}
	}
}
