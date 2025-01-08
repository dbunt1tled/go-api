package main

import (
	"go_echo/app/centrifugo"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/validate"
	proxyproto "go_echo/internal/grpc"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/storage"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := env.GetConfigInstance()
	locale.GetLocaleBundleInstance()
	logger.InitLogger(cfg.Env, cfg.Debug.Debug)
	log := logger.GetLoggerInstance()
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	run(cfg, log)
}

func run(cfg *env.Config, log *logger.AppLogger) {
	lis, err := net.Listen("tcp", cfg.Centrifugo.ServerURL)
	if err != nil {
		log.Error("failed to listen: " + err.Error())
	}
	srv := grpc.NewServer()
	proxyproto.RegisterCentrifugoProxyServer(srv, &centrifugo.Server{})
	reflection.Register(srv)
	log.Info("Start GRPC listening...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		log.Info("gRPC server is running on :" + cfg.Centrifugo.ServerURL)
		if err = srv.Serve(lis); err != nil {
			log.Error("failed to serve: " + err.Error())
		}
	}()

	<-stop
	log.Info("Shutting down gRPC server...")
	gracefulShutdown(
		log,
		func() error {
			log.Info("㋡ Quit: closing database connection")
			return storage.Close()
		},
		func() error {
			log.Info("㋡ gRPC Server stopped")
			srv.GracefulStop()
			return nil
		},
		func() error {
			log.Warn("｡◕‿‿◕｡ Quit: shutdown completed")
			os.Exit(0)
			return nil
		},
	)
}
func gracefulShutdown(log *logger.AppLogger, ops ...func() error) {
	for _, op := range ops {
		if err := op(); err != nil {
			log.Error("(ツ)_/¯ gRPC Graceful Shutdown op failed", "error", err)
			panic(err)
		}
	}
}
