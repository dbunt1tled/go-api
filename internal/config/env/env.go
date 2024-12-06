package env

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppName     string `env:"APP_NAME" env-required:"true"`
	AppURL      string `env:"APP_URL" env-default:"http://localhost"`
	Env         string `env:"ENV" env-default:"dev"`
	Debug       bool   `env:"DEBUG" env-default:"false"`
	Profiling   bool   `env:"PROFILING" env-default:"false"`
	DatabaseDSN string `env:"DATABASE_DSN" env-required:"true"`
	HTTPServer  HTTPServer
	CORS        CORS
	JWT         JWT
	Mail        Mail
}

type Mail struct {
	Host        string `env:"MAIL_HOST" env-required:"true"`
	Port        int    `env:"MAIL_PORT" env-required:"true"`
	Username    string `env:"MAIL_USERNAME" env-required:"true"`
	Password    string `env:"MAIL_PASSWORD" env-required:"true"`
	AddressFrom string `env:"MAIL_FROM_ADDRESS" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080" env-required:"true"`
	Timeout     time.Duration `env:"HTTP_SERVER_TIMEOUT" env-required:"true"`
	IdleTimeout time.Duration `env:"HTTP_SERVER_IDLE_TIMEOUT" env-required:"true"`
}

type CORS struct {
	AccessControlAllowHeaders  string `env:"ACCESS_CONTROL_ALLOW_HEADERS" env-default:""`
	AccessControlExposeHeaders string `env:"ACCESS_CONTROL_EXPOSE_HEADERS" env-default:""`
	AccessControlAllowMethods  string `env:"ACCESS_CONTROL_ALLOW_METHODS" env-default:""`
	AccessControlAllowOrigin   string `env:"ACCESS_CONTROL_ALLOW_ORIGIN" env-default:""`
}

type JWT struct {
	PublicKey       string        `env:"JWT_PUBLIC_KEY" env-required:"true"`
	PrivateKey      string        `env:"JWT_PRIVATE_KEY" env-required:"true"`
	Algorithm       string        `env:"JWT_TOKEN_ALGORITHM" env-default:"HS256"`
	AccessLifeTime  time.Duration `env:"TOKEN_ACCESS_LIFE_TIME_SECONDS" env-default:"3600s"`
	RefreshLifeTime time.Duration `env:"TOKEN_REFRESH_LIFE_TIME_SECONDS" env-default:"7200s"`
	SystemAPIKey    string        `env:"SYSTEM_API_KEY" env-required:"true"`
}

func MustLoadConfig() *Config {
	var cfg Config
	configPath := ".env"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("Error load config enviroment: %s", err)
		}
	} else {
		err = cleanenv.ReadConfig(configPath, &cfg)
		if err != nil {
			log.Fatalf("Error load config file enviroment: %s", err)
		}
	}

	return &cfg
}

var instance *Config //nolint:gochecknoglobals // singleton

func GetConfigInstance() *Config {
	if instance == nil {
		instance = MustLoadConfig()
	}
	return instance
}
