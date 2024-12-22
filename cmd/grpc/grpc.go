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
	"log/slog"
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
	logger.InitLogger(cfg.Env, cfg.Debug)
	log := logger.GetLoggerInstance()
	validate.GetValidateInstance()
	profiler.SetProfiler()
	storage.GetInstance()
	defer storage.Close()
	run(cfg, log)
}

func run(cfg *env.Config, log *slog.Logger) {
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
	srv.GracefulStop()
	log.Info("Server stopped gracefully")
}
