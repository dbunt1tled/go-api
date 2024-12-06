package main

import (
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/mailer"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/graceful"
	"go_echo/internal/lib/handler"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/router"
	"go_echo/internal/storage"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	log := logger.GetLoggerInstance()
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	mailer.GetMailInstance()
	defer mailer.Close()
	httpServer := echo.New()
	httpServer.HideBanner = true
	httpServer.Debug = cfg.Debug
	httpServer.HTTPErrorHandler = handler.APIErrorHandler
	router.SetupRoutes(httpServer)
	done := graceful.ShutdownGraceful(log, httpServer)

	println(language.English.String())

	// m := mail.NewMsg()
	// if err := m.From(cfg.Mail.AddressFrom); err != nil {
	//
	// }
	// if err := m.To("d.balagov.bekey@gmail.com"); err != nil {
	//
	// }
	// m.Subject("Why are you not using go-mail yet?")
	// m.SetBodyString(mail.TypeTextPlain, "You won't need a sales pitch. It's FOSS.")
	//
	// // Your message-specific code here
	// if err := client.DialAndSend(m); err != nil {
	// 	panic("failed to send mail: " + err.Error())
	// }

	// u, err := service.UserRepository{}.List([]builder.FilterCondition{
	// 	{Field: "id", Type: builder.In, Value: []interface{}{1, 2}},
	// }, []builder.SortOrder{
	// 	{Field: "id", Order: builder.Desc},
	// })
	// if err != nil {
	// 	log.Error(err.Error())
	// } else {
	// 	log.Debug(fmt.Sprintf("%+v\n", u))
	// }

	go func() {
		log.Debug("Start listening on address: " + cfg.HTTPServer.Address)
		if err := httpServer.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) { //nolint:lll,govet
			log.Error("shutting down the server")
		}
	}()
	<-done
}
