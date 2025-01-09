package main

import (
	"crypto/tls"
	"go_echo/app/centrifugo"
	"go_echo/internal/config/env"
	"go_echo/internal/config/locale"
	"go_echo/internal/config/logger"
	"go_echo/internal/config/validate"
	proxyproto "go_echo/internal/grpc"
	"go_echo/internal/lib/profiler"
	"go_echo/internal/storage"
	"go_echo/internal/util/helper"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	defer storage.Close()
	run(cfg, log)
}

func run(cfg *env.Config, log *logger.AppLogger) {
	var (
		lis  net.Listener
		srv  *grpc.Server
		err  error
		cert tls.Certificate
		cred credentials.TransportCredentials
	)
	if cfg.HTTPServer.TLS.IsSet() {
		certData := cfg.HTTPServer.TLS.GetCertData()
		keyData := cfg.HTTPServer.TLS.GetKeyData()

		if helper.IsA(certData, reflect.String) && helper.IsA(keyData, reflect.String) {
			cred, err = credentials.NewServerTLSFromFile(helper.AnyToString(certData), helper.AnyToString(keyData))
			if err != nil {
				log.Error("failed to create TLS credentials: " + err.Error())
				return
			}
		} else {
			cert, err = tls.X509KeyPair(certData.([]byte), keyData.([]byte))
			if err != nil {
				log.Error("failed to create TLS certificate: " + err.Error())
				return
			}
			cred = credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		}
		srv = grpc.NewServer(grpc.Creds(cred))
		lis, err = net.Listen("tcp", cfg.Centrifugo.ServerURL)
		if err != nil {
			log.Error("failed to listen: " + err.Error())
		}
	} else {
		lis, err = net.Listen("tcp", cfg.Centrifugo.ServerURL)
		if err != nil {
			log.Error("failed to listen: " + err.Error())
		}
	}

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
