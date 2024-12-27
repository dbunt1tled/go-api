package router

import (
	"go_echo/app/auth/authhandler"
	"go_echo/app/general/generalhandler"
	"go_echo/app/usernotification/usernotificationhandler"
	"go_echo/internal/config/env"
	apiServer "go_echo/internal/router/handler/server"
	"go_echo/internal/router/middlewares"
	"go_echo/internal/util/hasher"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(server *echo.Echo) {
	cfg := env.GetConfigInstance()
	setGeneralMiddlewares(server, cfg)
	systemRouter := server.Group("/system")
	systemRouter.Use(middlewares.SystemAuth)
	systemRouter.GET("/helm", apiServer.Helm)
	if cfg.Static.Enable {
		server.Static("/"+cfg.Static.URL, cfg.Static.Directory)
	}
	generalRoutes(server)
	authRoutes(server)
	UserNotificationRoutes(server)
}
func generalRoutes(server *echo.Echo) {
	generalRouter := server.Group("/")
	// generalRouter.Use(middlewares.AuthBearer)
	generalRouter.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	generalRouter.GET("home", generalhandler.Home)
}

func authRoutes(server *echo.Echo) {
	authRouter := server.Group("/auth")
	authRouter.POST("/login", authhandler.Login)
	authRouter.POST("/register", authhandler.Register)
	authRouter.GET("/confirm", authhandler.Confirm)
}

func UserNotificationRoutes(server *echo.Echo) {
	userNotificationsRouter := server.Group("/notifications")
	userNotificationsRouter.Use(middlewares.AuthBearer)
	userNotificationsRouter.GET("", usernotificationhandler.UserNotificationList)
}

func setGeneralMiddlewares(server *echo.Echo, cfg *env.Config) {
	server.Use(middleware.Recover())
	server.Use(middleware.Gzip())
	server.Use(middlewares.Base)
	server.Use(middlewares.Language())
	server.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			var u string
			ub, e := hasher.UUIDVv7()
			if e != nil {
				u = strconv.FormatInt(time.Now().UnixMicro(), 10)
			} else {
				u = ub.String()
			}

			return u
		},
	}))
	if cfg.Debug {
		server.Use(middleware.Logger())
	}
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(cfg.CORS.AccessControlAllowOrigin, ","),
		AllowMethods:     strings.Split(cfg.CORS.AccessControlAllowMethods, ","),
		AllowHeaders:     strings.Split(cfg.CORS.AccessControlAllowHeaders, ","),
		AllowCredentials: false,
		MaxAge:           300, //nolint:mnd // Maximum value not ignored by any of major browsers
	}))
}
