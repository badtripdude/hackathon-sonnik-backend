package config

import (
	"log"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerPort string `env:"SERVER_HTTP_PORT" envDefault:"8082"`

	ServerReadTimeout       time.Duration `env:"SERVER_HTTP_READ_TIMEOUT" envDefault:"5s"`
	ServerWriteTimeout      time.Duration `env:"SERVER_HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	ServerIdleTimeout       time.Duration `env:"SERVER_HTTP_IDLE_TIMEOUT" envDefault:"120s"`
	ServerReadHeaderTimeout time.Duration `env:"SERVER_HTTP_READ_HEADER_TIMEOUT" envDefault:"5s"`

	PgDSN             string        `env:"PG_DSN" envDefault:"postgresql://app:app@postgres:5432/app?sslmode=disable"`
	PgMaxOpenConns    int           `env:"PG_MAX_OPEN_CONNS" envDefault:"10"`
	PgMaxIdleConns    int           `env:"PG_MAX_IDLE_CONNS" envDefault:"5"`
	PgConnMaxLifetime time.Duration `env:"PG_CONN_MAX_LIFETIME" envDefault:"1h"`
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		c := &Config{}
		if err := env.Parse(c); err != nil {
			log.Fatalf("failed to parse env: %v", err)
		}

		cfg = c
		log.Printf("config loaded: server_port=%s", cfg.ServerPort)
	})

	return cfg
}
