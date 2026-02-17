package main

import (
	"crypto/tls"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dbunt1tled/go-api/internal/config"
	centrifugo "github.com/dbunt1tled/go-api/internal/grpc"
	"github.com/dbunt1tled/go-api/internal/grpc/proxyproto"
	"github.com/dbunt1tled/go-api/internal/modules/auth"
	"github.com/dbunt1tled/go-api/internal/modules/user"
	"github.com/dbunt1tled/go-api/internal/modules/user_notification"
	"github.com/dbunt1tled/go-api/pkg/f"
	"github.com/dbunt1tled/go-api/pkg/grpc/interceptor"
	"github.com/dbunt1tled/go-api/pkg/hasher"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/dbunt1tled/go-api/pkg/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		lis         net.Listener
		srv         *grpc.Server
		err         error
		cert        tls.Certificate
		cred        credentials.TransportCredentials
		opts        []grpc.ServerOption
		middlewares []grpc.UnaryServerInterceptor
	)
	config.Load()
	log.Load(config.Get().Name, config.Get().Env, config.Get().Log.Level, config.Get().Log.File)
	db := postgres.New(
		config.Get().DB.Main.DSN,
		config.Get().Debug,
	)

	middlewares = append(middlewares, interceptor.RecoverInterceptor())
	if config.Get().Debug {
		middlewares = append(middlewares, interceptor.LoggingInterceptor())
	}
	opts = append(opts, grpc.ChainUnaryInterceptor(middlewares...))
	if config.Get().Server.HTTP.TLS.IsSet() {
		certData := config.Get().Server.HTTP.TLS.GetCertData()
		keyData := config.Get().Server.HTTP.TLS.GetKeyData()
		cert, err = tls.X509KeyPair(certData, keyData)
		if err != nil {
			log.Logger().Error("failed to create TLS certificate", err)
			return
		}
		cred = credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})
		opts = append(opts, grpc.Creds(cred))
	}
	srv = grpc.NewServer(opts...)
	lis, err = net.Listen("tcp", config.Get().Centrifugo.ServerUrl)
	if err != nil {
		log.Logger().Error("failed to listen", err)
	}

	hashService := f.Must(hasher.NewHasher(
		config.Get().Server.JWT.Algorithm,
		config.Get().Server.JWT.PublicKey,
		config.Get().Server.JWT.PrivateKey,
	))
	authService := auth.NewAuthService(hashService)
	userService := user.NewUserService(db.DB())
	userNotificationService := user_notification.NewUserNotificationService(db.DB())

	proxyproto.RegisterCentrifugoProxyServer(
		srv,
		centrifugo.NewCentrifugoServer(
			authService,
			userService,
			userNotificationService,
		),
	)
	reflection.Register(srv)
	log.Logger().Info("Start GRPC listening...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		log.Logger().Info("gRPC server is running on :" + config.Get().Centrifugo.ServerUrl)
		if err = srv.Serve(lis); err != nil {
			log.Logger().Error("failed to serve", err)
		}
	}()

	code := <-stop
	log.Logger().Info("System received signal: " + code.String())
	log.Logger().Info("Shutting down gRPC server...")

	f.MultiRunFunc("Shutdown", log.Logger(),
		func() error {
			log.Logger().Info("㋡ Quit: closing database connection")
			return db.DB().Close()
		},
		func() error {
			log.Logger().Info("㋡ Quit: closing GRPC server")
			srv.GracefulStop()
			return nil
		},
		func() error {
			log.Logger().Warn("｡◕‿‿◕｡ Quit: shutdown completed")
			return nil
		})
}
