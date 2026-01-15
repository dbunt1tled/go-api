package app

import (
	"net/http"

	"github.com/dbunt1tled/go-api/internal/lib/view"
	"github.com/labstack/echo/v5"
)

func Router(application *App) {
	WebRoutes(application)
	ApiRoutes(application)
}
func ApiRoutes(application *App) {
	app := application.Engine()
	api := app.Group("api")
	apiRoutes(api, application)
}
func WebRoutes(application *App) {
	app := application.Engine()
	app.GET("", func(c *echo.Context) error {
		return c.Render(
			http.StatusOK,
			"general/home.gohtml",
			view.MakeTemplateData(map[string]any{}),
		)
	})
}

func apiRoutes(api *echo.Group, app *App) {
	user := api.Group("users")
	usersRoutes(user, app)
}

func usersRoutes(user *echo.Group, app *App) {
	user.GET("", app.UserController.List)
}
