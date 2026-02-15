package app

import (
	"fmt"
	"net/http"

	"github.com/dbunt1tled/go-api/internal/lib/view"
	"github.com/labstack/echo/v5"
)

func Router(application *App) {
	WebRoutes(application)
	ApiRoutes(application)
	for _, r := range application.Engine().Router().Routes() {
		fmt.Println(fmt.Sprintf("name: %s, method: %s, path: %s", r.Name, r.Method, r.Path))
	}
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
	usersRoutes(api, app)
	authRoutes(api, app)
}

func usersRoutes(api *echo.Group, app *App) {
	user := api.Group("/users", app.AuthMiddleware.AuthBearer)
	user.GET("", app.UserController.List)
	user.GET("/profile", app.UserController.Profile)
	userNotificationGroup(user, app)
}

func userNotificationGroup(api *echo.Group, app *App) {
	userNotification := api.Group("/notifications", app.AuthMiddleware.AuthBearer)
	userNotification.GET("", app.UserNotificationController.List)
}

func authRoutes(api *echo.Group, app *App) {
	auth := api.Group("/auth")
	auth.POST("/register", app.AuthController.Register)
	auth.POST("/login", app.AuthController.Login)
	auth.POST("/refresh", app.AuthController.Refresh)
	auth.GET("/confirm/:token", app.AuthController.Confirm)
}
