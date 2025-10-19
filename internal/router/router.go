package router

import (
	"strconv"
	"strings"
	"time"

	"github.com/dbunt1tled/go-api/app/auth/authhandler"
	"github.com/dbunt1tled/go-api/app/general/generalhandler"
	"github.com/dbunt1tled/go-api/app/user/userhandler"
	"github.com/dbunt1tled/go-api/app/usernotification/usernotificationhandler"
	"github.com/dbunt1tled/go-api/internal/config/env"
	apiServer "github.com/dbunt1tled/go-api/internal/router/handler/server"
	"github.com/dbunt1tled/go-api/internal/router/middlewares"
	"github.com/dbunt1tled/go-api/internal/util/hasher"

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
	userRoutes(server)
	UserNotificationRoutes(server)
}
func generalRoutes(server *echo.Echo) {
	generalRouter := server.Group("/")
	generalRouter.GET("", generalhandler.Home)
}

func authRoutes(server *echo.Echo) {
	authRouter := server.Group("/auth")
	authRouter.POST("/login", authhandler.Login)
	authRouter.POST("/register", authhandler.Register)
	authRouter.GET("/confirm", authhandler.Confirm)
}

func userRoutes(server *echo.Echo) {
	userNotificationsRouter := server.Group("/users")
	userNotificationsRouter.Use(middlewares.AuthBearer)
	userNotificationsRouter.GET("/profile", userhandler.Profile)
}

func UserNotificationRoutes(server *echo.Echo) {
	userNotificationsRouter := server.Group("/notifications")
	userNotificationsRouter.Use(middlewares.AuthBearer)
	userNotificationsRouter.GET("", usernotificationhandler.UserNotificationList)
}

func setGeneralMiddlewares(server *echo.Echo, cfg *env.Config) {
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
	server.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 6}))
	server.Use(middleware.Recover())
	server.Use(middlewares.Language())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(cfg.CORS.AccessControlAllowOrigin, ","),
		AllowMethods:     strings.Split(cfg.CORS.AccessControlAllowMethods, ","),
		AllowHeaders:     strings.Split(cfg.CORS.AccessControlAllowHeaders, ","),
		AllowCredentials: false,
		MaxAge:           300, //nolint:mnd // Maximum value not ignored by any of major browsers
	}))
	server.Use(middlewares.LogRequest)
	server.Use(middlewares.Base)
}
