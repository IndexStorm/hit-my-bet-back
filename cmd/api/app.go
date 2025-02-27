package main

import (
	"context"
	"fmt"
	"github.com/IndexStorm/hit-my-bet-back/pkg/log"
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
	"io"
)

type application struct {
	config  appConfig
	logger  zerolog.Logger
	closers []io.Closer
}

func newApplication() (*application, error) {
	var cfg appConfig
	err := env.Parse(&cfg)
	return &application{
		config: cfg,
		logger: log.NewWithLevel(cfg.LogLevel),
	}, err
}

func (a *application) start(ctx context.Context) error {
	otel.SetTracerProvider(noop.NewTracerProvider())

	dependencies, err := a.buildDependencies(ctx)
	if err != nil {
		return fmt.Errorf("build dependencies: %w", err)
	}
	a.closers = append(a.closers, dependencies)
	if err = dependencies.server.start(":5050"); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	return nil
}

func (a *application) stop() {
	for _, closer := range a.closers {
		if err := closer.Close(); err != nil {
			a.logger.Err(err).Msg("stop:closer failed")
		}
	}
}

func (a *application) buildDependencies(ctx context.Context) (*applicationDependencies, error) {
	builder := newDependencyBuilder(a.config, a.logger)
	return builder.build(ctx)
}
