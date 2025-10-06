package config

import (
	"Effective_Mobile/internal/logger"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DsnPG      DsnPG      `envPrefix:"DSN_"`
	HTTPServer HTTPServer `envPrefix:"HTTP_"`
	Debug      bool       `env:"DEBUG"`
}

type DsnPG struct {
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"admin"`
	Password string `env:"PASSWORD" env-default:"adm_123"`
	Name     string `env:"NAME" env-default:"myapp"`
	Host     string `env:"HOST" env-default:"localhost"`
}

type HTTPServer struct {
	Address     string        `env:"ADDR" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	User        string        `env:"USER" env-required:"true"`
	Password    string        `env:"HTTP_SERVER_PASSWORD" env-required:"true"`
}

func MustLoad() *Config {
	const op = "config.MustLoad"

	if err := godotenv.Load(); err != nil {
		logger.Info("%s: .env file not found or couldn't load, using environment variables only", op)
	} else {
		logger.Info("%s: .env file successfully loaded", op)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.Error("%s: failed to parse environment variables: %v", op, err)
		os.Exit(1)
	}

	logger.DebugEnabled = cfg.Debug
	if cfg.Debug {
		logger.Info("%s: debug mode enabled", op)
	}

	logger.Info("%s: configuration loaded successfully", op)
	return &cfg
}
