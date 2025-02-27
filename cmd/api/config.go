package main

import (
	"github.com/IndexStorm/hit-my-bet-back/internal/config"
	"github.com/rs/zerolog"
)

type appConfig struct {
	Environment config.DefaultEnvironment
	LogLevel    zerolog.Level `env:"LOG_LEVEL,notEmpty"`
	Telemetry   config.Telemetry
	Database    config.Database `envPrefix:"DB_" env:"notEmpty"`
}
