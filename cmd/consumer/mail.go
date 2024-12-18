package main

import (
	"encoding/json"
	"fmt"
	"go_echo/app/user/model/user"
	"go_echo/app/user/service"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/mailservice"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/rmq"
	"go_echo/internal/storage"

	"github.com/rabbitmq/amqp091-go"
)

const (
	NumWorkers = 6
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()

	rc := rmq.GetRMQInstance(rmq.MailExchange)
	f := func(d amqp091.Delivery) error {
		var (
			job mailservice.MailJobMessage
			u   *user.User
			err error
		)
		log := logger.GetLoggerInstance()
		if err = json.Unmarshal(d.Body, &job); err != nil {
			return fmt.Errorf("failed to unmarshal message: %s", err.Error())
		}
		u, err = service.UserRepository{}.ByID(int64(job.UserID))
		if err != nil {
			log.Error(fmt.Sprintf("User: #%d not found. %s", job.UserID, err.Error()))
			return nil
		}
		switch job.Subject {
		case mailservice.ConfirmSubject:
			mailservice.SendUserConfirmEmail(u, job.Token)
		default:
			return fmt.Errorf("unknown subject: %s", job.Subject)
		}
		return nil
	}
	rc.Consume(rmq.MailExchange, rmq.MailQueue, f, NumWorkers)
	rc.Close()
}
