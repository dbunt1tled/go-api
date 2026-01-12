package config

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var (
	cfg  atomic.Value
	once sync.Once
)

func Get() *Config {
	return cfg.Load().(*Config)
}

func Load() {
	once.Do(func() {
		err := load()
		if err != nil {
			panic(err)
		}
	})
}

func load() error {
	k := koanf.New(".")

	if err := k.Load(confmap.Provider(map[string]interface{}{
		"name":                  "fiber-api",
		"env":                   "develop",
		"debug":                 false,
		"server.http.host":      "localhost",
		"server.http.port":      8081, //nolint:mnd // default port
		"server.http.timeout":   "5s",
		"server.http.bodylimit": 4 * 1024 * 1024, //nolint:mnd // 4MB
	}, "."), nil); err != nil {
		return fmt.Errorf("error loading confmap: %w", err)
	}

	if err := k.Load(file.Provider(".env"), dotenv.ParserEnv("APP_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "APP_")), "_", ".", -1)
	})); err != nil {
		log.Printf("Warning: .env file error: %v", err)
	}

	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(k, "APP_")), "_", ".")
			if strings.Contains(v, " ") {
				return k, strings.Split(v, " ")
			}

			return k, v
		},
	}), nil); err != nil {
		return fmt.Errorf("error loading env variables: %w", err)
	}

	var c Config
	if err := k.Unmarshal("", &c); err != nil {
		return err
	}

	cfg.Store(&c)

	return nil
}
