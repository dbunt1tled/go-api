package main

import (
	"encoding/json"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/lib/centservice"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	log := logger.GetLoggerInstance()
	u, err := centservice.SendUserNotification(centservice.UserNotification{
		UserID:  3,
		Message: "Your account has been confirmed",
	})
	if err != nil {
		log.Error("Error send user notification: " + err.Error())
		return
	}
	d, err := json.Marshal(u)
	if err != nil {
		log.Error("Error marshal user notification: " + err.Error())
		return
	}
	log.Info("User notification sent: " + string(d))
}