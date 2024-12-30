package env

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppName     string `env:"APP_NAME" env-required:"true"`
	AppURL      string `env:"APP_URL" env-default:"http://localhost"`
	Env         string `env:"ENV" env-default:"dev"`
	Debug       Debug
	Profiling   bool   `env:"PROFILING" env-default:"false"`
	DatabaseDSN string `env:"DATABASE_DSN" env-required:"true"`
	HTTPServer  HTTPServer
	CORS        CORS
	JWT         JWT
	Mail        Mail
	Static      Static
	RabbitMQ    RabbitMQ
	Redis       Redis
	Centrifugo  Centrifugo
}

type Redis struct {
	Host     string `env:"REDIS_HOST" env-default:"127.0.0.1"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB" env-default:"0"`
}

type Debug struct {
	Debug        bool `env:"DEBUG" env-default:"false"`
	DebugRequest bool `env:"DEBUG_REQUEST" env-default:"false"`
	DebugBody    bool `env:"DEBUG_BODY" env-default:"false"`
	Profiling    bool `env:"PROFILING" env-default:"false"`
}

type Centrifugo struct {
	ServerURL string `env:"SERVER_CENTRIFUGO_URL" env-default:"localhost:5000"`
	APIURL    string `env:"CENTRIFUGO_API_URL" env-default:"http://127.0.0.1:8000"`
	APIKey    string `env:"CENTRIFUGO_API_KEY" env-required:"true"`
}

type RabbitMQ struct {
	Host     string `env:"RABBITMQ_HOST" env-required:"true"`
	Port     int    `env:"RABBITMQ_PORT" env-required:"true"`
	User     string `env:"RABBITMQ_USER" env-required:"true"`
	Password string `env:"RABBITMQ_PASSWORD" env-required:"true"`
	Vhost    string `env:"RABBITMQ_VHOST" env-required:"true"`
}

type Static struct {
	Enable    bool   `env:"HTTP_STATIC" env-default:"false"`
	Directory string `env:"HTTP_STATIC_DIR" env-default:"HTTP_STATIC_DIR"`
	URL       string `env:"HTTP_STATIC_URL" env-default:"HTTP_STATIC_URL"`
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
	ConfirmLifeTime time.Duration `env:"TOKEN_CONFIRM_LIFE_TIME_SECONDS" env-default:"7200s"`
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

var (
	instance *Config   //nolint:gochecknoglobals // singleton
	m        sync.Once //nolint:gochecknoglobals // singleton
)

func GetConfigInstance() *Config {
	m.Do(func() {
		instance = MustLoadConfig()
	})
	return instance
}
