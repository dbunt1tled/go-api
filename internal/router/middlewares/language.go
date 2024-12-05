package middlewares

import (
	"go_echo/internal/config/locale"

	"github.com/labstack/echo/v4"
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
			locale.InitLocalizerInstance(l, l)

			return next(c)
		}
	}
}
