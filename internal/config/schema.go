package config

import (
	"encoding/base64"
	"log"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	Name       string           `koanf:"name"`
	URL        string           `koanf:"url"`
	Env        string           `koanf:"env"`
	Debug      bool             `koanf:"debug"`
	Profiling  bool             `koanf:"profiling"`
	Server     ServerConfig     `koanf:"server"`
	DB         DBConfig         `koanf:"db"`
	Redis      RedisConfig      `koanf:"redis"`
	Log        LogConfig        `koanf:"log"`
	Mailer     MailerConfig     `koanf:"mailer"`
	Static     StaticConfig     `koanf:"static"`
	AMQP       AMQPConfig       `koanf:"amqp"`
	Centrifugo CentrifugoConfig `koanf:"centrifugo"`
}
type ServerConfig struct {
	HTTP HTTPConfig `koanf:"http"`
	JWT  JWTConfig  `koanf:"jwt"`
}

type HTTPConfig struct {
	Host      string        `koanf:"host"`
	Port      int           `koanf:"port"`
	Timeout   time.Duration `koanf:"timeout"`
	Prefork   bool          `koanf:"prefork"`
	CORS      CORSConfig    `koanf:"cors"`
	BodyLimit int64         `koanf:"bodylimit"`
	TLS       TLSConfig     `koanf:"tls"`
}

type TLSConfig struct {
	Key  string `koanf:"key"`
	Cert string `koanf:"cert"`
}

func (t TLSConfig) IsSet() bool {
	return t.Key != "" && t.Cert != ""
}
func (t TLSConfig) GetCertData() []byte {
	var (
		err error
		res []byte
	)
	_, err = os.Stat(t.Cert)
	if err == nil {
		certBytes, err := ReadCert(t.Cert)
		if err != nil {
			log.Fatal(err)
		}
		return certBytes
	}
	res, err = base64.StdEncoding.DecodeString(t.Cert)
	if err == nil {
		return res
	}
	return []byte(t.Cert)
}

func (t TLSConfig) GetKeyData() []byte {
	var (
		err error
		res []byte
	)
	_, err = os.Stat(t.Key)
	if err == nil {
		certBytes, err := ReadCert(t.Key)
		if err != nil {
			log.Fatal(err)
		}
		return certBytes
	}
	res, err = base64.StdEncoding.DecodeString(t.Key)
	if err == nil {
		return res
	}
	return []byte(t.Key)
}

func ReadCert(path string) ([]byte, error) {
	return os.ReadFile(path)
}

type JWTConfig struct {
	PublicKey  string     `koanf:"public"`
	PrivateKey string     `koanf:"private"`
	Algorithm  string     `koanf:"algorithm"`
	Expire     ExpireConf `koanf:"expire"`
}

type ExpireConf struct {
	Access  time.Duration `koanf:"access"`
	Refresh time.Duration `koanf:"refresh"`
	Confirm time.Duration `koanf:"confirm"`
}

type LogConfig struct {
	Level slog.Level `koanf:"level"`
	File  string     `koanf:"file"`
}
type CORSConfig struct {
	AllowMethods  string `koanf:"allowmethods"`
	AllowHeaders  string `koanf:"allowheaders"`
	AllowOrigins  string `koanf:"alloworigins"`
	ExposeHeaders string `koanf:"exposeheaders"`
}

type DBConfig struct {
	Main MainDBConfig `koanf:"main"`
}

type RedisConfig struct {
	Addr string `koanf:"addr"`
}

type MainDBConfig struct {
	DSN string `koanf:"dsn"`
}

type MailerConfig struct {
	Address  string `koanf:"address"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type StaticConfig struct {
	URL       string `koanf:"url"`
	Directory string `koanf:"dir"`
}

type AMQPConfig struct {
	URL string `koanf:"url"`
}

type CentrifugoConfig struct {
	ServerUrl string `koanf:"server"`
	APIKey    string `koanf:"key"`
	APIUrl    string `koanf:"api"`
}
