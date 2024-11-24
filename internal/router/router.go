package router

import (
	"fmt"
	"go_echo/app/auth/authhandle"
	"go_echo/internal/config/env"
	apiServer "go_echo/internal/router/handler/server"
	"go_echo/internal/router/mware"
	"go_echo/internal/util/hash"
	"go_echo/internal/util/rand"
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
	systemRouter.Use(mware.SystemAuth)
	systemRouter.GET("/helm", apiServer.Helm)
	authRoutes(server)
	// group.GET("", h.HandlerShowUsers)
	// group.GET("/details/:id", h.HandlerShowUserById)

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}

func authRoutes(server *echo.Echo) {
	authRouter := server.Group("/auth")
	authRouter.POST("/login", authhandle.Login)
}

func setGeneralMiddlewares(server *echo.Echo, cfg *env.Config) {
	server.Use(middleware.Recover())
	server.Use(middleware.Logger())
	server.Use(middleware.Gzip())
	server.Use(mware.Base)
	server.Use(mware.Language())
	server.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			var u string
			ub, e := hash.UUIDVv7()
			if e != nil {
				u = strconv.FormatInt(time.Now().UnixMicro(), 10)
			} else {
				u = ub.String()
			}
			r, e := rand.String(4) //nolint:mnd //small random part
			if e != nil {
				r = strconv.FormatInt(time.Now().Unix(), 10)
			}
			return fmt.Sprintf("%s:%s", u, r)
		},
	}))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(cfg.CORS.AccessControlAllowOrigin, ","),
		AllowMethods:     strings.Split(cfg.CORS.AccessControlAllowMethods, ","),
		AllowHeaders:     strings.Split(cfg.CORS.AccessControlAllowHeaders, ","),
		AllowCredentials: false,
		MaxAge:           300, //nolint:mnd // Maximum value not ignored by any of major browsers
	}))
}
