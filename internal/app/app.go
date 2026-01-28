package app

import (
	"context"
	"errors"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dbunt1tled/go-api/internal/config"
	"github.com/dbunt1tled/go-api/internal/modules/auth"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/pkg/f"
	"github.com/dbunt1tled/go-api/pkg/hasher"
	h "github.com/dbunt1tled/go-api/pkg/http"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/validation"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type App struct {
	cfg            *config.ServiceConfig
	engine         *echo.Echo
	AuthController *auth.Controller
	UserController *user.Controller
}

func NewApp(cfg *config.ServiceConfig) *App {

	hashService, err := hasher.NewHasher(
		config.Get().Server.JWT.Algorithm,
		config.Get().Server.JWT.PublicKey,
		config.Get().Server.JWT.PrivateKey,
	)
	if err != nil {
		panic(err)
	}

	engine := engineSetup(cfg)

	engine.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			var u string
			ub, e := hashService.UUIDVv7()
			if e != nil {
				u = strconv.FormatInt(time.Now().UnixMicro(), 10)
			} else {
				u = ub.String()
			}
			return u
		},
	}))
	engine.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 6})) //nolint:nolintlint,mnd
	engine.Use(middleware.Recover())
	engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     strings.Split(config.Get().Server.HTTP.CORS.AllowOrigins, ","),
		AllowMethods:     strings.Split(config.Get().Server.HTTP.CORS.AllowMethods, ","),
		AllowHeaders:     strings.Split(config.Get().Server.HTTP.CORS.AllowHeaders, ","),
		ExposeHeaders:    strings.Split(config.Get().Server.HTTP.CORS.ExposeHeaders, ","),
		AllowCredentials: false,
		MaxAge:           300, //nolint:mnd // Maximum value not ignored by any of major browsers
	}))

	if config.Get().Static.URL != "" && config.Get().Static.Directory != "" {
		engine.Static(config.Get().Static.URL, config.Get().Static.Directory)
	}

	userService := user.NewUserService(cfg.DB.DB())

	return &App{
		cfg:            cfg,
		engine:         engine,
		AuthController: auth.NewAuthController(auth.NewAuthService(hashService), cfg.RMProducer, userService),
		UserController: user.NewUserController(userService),
	}
}

func (a *App) Engine() *echo.Echo {
	return a.engine
}

func (a *App) Run(ctx context.Context) {
	ct, stop := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt,
	)
	defer stop()

	go func() {
		sc := echo.StartConfig{
			Address:    config.Get().Server.HTTP.Host + ":" + strconv.Itoa(config.Get().Server.HTTP.Port),
			HideBanner: true,

			GracefulTimeout: 10 * time.Second,
		}
		if config.Get().Server.HTTP.TLS.IsSet() {
			log.Logger().Debug("(っ◕‿◕)っ Start Server TLS listening on address: " + config.Get().URL)
			err := sc.StartTLS(
				ct,
				a.engine,
				config.Get().Server.HTTP.TLS.GetCertData(),
				config.Get().Server.HTTP.TLS.GetKeyData(),
			)
			if err != nil {
				log.Logger().Errorf("¯\\_(͡° ͜ʖ ͡°)_/¯Shutting down the server: %v", err)
			}
		} else {
			log.Logger().Debug("(/◔◡◔)/ Start Server listening on address: " + config.Get().URL)
			if err := sc.Start(ct, a.engine); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Logger().Fatalf("¯\\_(͡° ͜ʖ ͡°)_/¯Shutting down the server: %v", err)
			}
		}
	}()
	<-ct.Done()
	var cancel context.CancelFunc
	ct, cancel = context.WithTimeout(context.Background(), 10*time.Second) //nolint:mnd // 10 seconds timeout
	defer cancel()
	log.Logger().Warn("Quit: shutting down ...")
	defer log.Logger().Warn("｡◕‿‿◕｡ Quit: shutdown completed")
	f.MultiRunFunc("Shutdown", log.Logger(),
		func() error {
			log.Logger().Info("㋡ Quit: closing database connection")
			return a.cfg.DB.Close()
		},
		func() error {
			log.Logger().Info("㋡ Quit: closing RabbitMQ connection")
			return a.cfg.RMProducer.Close()
		},
		func() error {
			log.Logger().Info("㋡ Quit: closing mailer")
			return a.cfg.Mailer.Close()
		})
}

func engineSetup(cfg *config.ServiceConfig) *echo.Echo {
	v, err := validation.Validator(cfg.DB.DB())
	if err != nil {
		panic(err)
	}
	engine := echo.NewWithConfig(echo.Config{
		Logger:           log.Logger().Base(),
		Validator:        v,
		HTTPErrorHandler: h.APIErrorHandler,
		JSONSerializer:   &h.JSONSerializer{},
		Renderer: &echo.TemplateRenderer{
			Template: template.Must(template.ParseGlob("./resources/templates/**/*.gohtml")),
		},
	})

	return engine
}
