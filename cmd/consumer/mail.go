package main

import (
	"fmt"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/validate"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/rmq"
	"go_echo/internal/storage"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug)
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()

	var rc rmq.RabbitClient
	f := func(d amqp091.Delivery) error {
		time.Sleep(time.Second * 20)
		fmt.Println(d.Body)
		return nil
	}
	rc.Consume("bb", "aaaa", f, 6)
	rc.Close()
}
