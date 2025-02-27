package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/IndexStorm/hit-my-bet-back/internal/postgres"
	"github.com/IndexStorm/hit-my-bet-back/internal/repository/prediction"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type dependencyBuilder struct {
	config appConfig
	logger zerolog.Logger
	tracer trace.Tracer
}

func newDependencyBuilder(config appConfig, logger zerolog.Logger) *dependencyBuilder {
	return &dependencyBuilder{
		config: config,
		logger: logger,
		tracer: otel.Tracer("dependency-builder"),
	}
}

func (b *dependencyBuilder) build(ctx context.Context) (*applicationDependencies, error) {
	db, err := b.newDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("prepare database: %w", err)
	}
	dependencies := &applicationDependencies{database: db}
	predictionRepo := prediction.NewPostgres(db)
	dependencies.predictionRepo = predictionRepo

	appServer := newServer(b.logger, otel.Tracer("server"), predictionRepo)
	dependencies.server = appServer

	return &applicationDependencies{server: appServer}, nil
}

func (b *dependencyBuilder) newDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	ctx, span := b.tracer.Start(ctx, "db:connect")
	defer span.End()
	return postgres.NewPgxPoolWithOtel(ctx, b.config.Database, b.config.Environment.Value)
}

type applicationDependencies struct {
	database       *pgxpool.Pool
	predictionRepo prediction.Repository
	server         *server
}

func (d *applicationDependencies) Close() error {
	var errs []error
	if d.server != nil {
		if err := d.server.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close server: %w", err))
		}
	}
	if d.database != nil {
		d.database.Close()
	}
	return errors.Join(errs...)
}
