package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/IndexStorm/hit-my-bet-back/internal/config"
	"github.com/IndexStorm/hit-my-bet-back/pkg/db"
	"github.com/IndexStorm/hit-my-bet-back/pkg/log"
	"github.com/caarlos0/env/v11"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"time"
)

type appConfig struct {
	Database     config.Database `envPrefix:"DB_" env:"notEmpty"`
	ForceVersion int             `env:"FORCE_VERSION"`
	SqlSchemaDir string          `env:"SQL_SCHEMA_DIR,notEmpty"`
}

type application struct {
	logger zerolog.Logger
	config appConfig
}

func newApplication() (*application, error) {
	var cfg appConfig
	err := env.Parse(&cfg)
	return &application{
		logger: log.NewWithLevel(zerolog.DebugLevel),
		config: cfg,
	}, err
}

func (a *application) migrate(ctx context.Context) error {
	a.logger.Info().Str("path", a.config.SqlSchemaDir).Msg("starting migration")
	db, err := a.connectDB(ctx)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()
	err = a.runMigrations(db)
	if errors.Is(err, migrate.ErrNoChange) {
		a.logger.Info().Msg("no changes detected")
		return nil
	} else if err != nil {
		return err
	}
	a.logger.Info().Msg("migrations applied")
	return nil
}

func (a *application) connectDB(ctx context.Context) (*pgxpool.Pool, error) {
	tracer := log.NewSQLTracer(a.logger.With().Str("sys", "sql").Logger())
	sslMode := "require"
	if a.config.Database.Plaintext {
		sslMode = "disable"
	}
	database, err := db.NewPgxConnection(context.Background(), fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s",
		a.config.Database.Username,
		a.config.Database.Password,
		a.config.Database.Host,
		a.config.Database.Database,
		sslMode,
	), tracer, nil, time.Second*10)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err = database.Ping(ctx); err != nil {
		database.Close()
		return nil, err
	}
	return database, nil
}

func (a *application) runMigrations(db *pgxpool.Pool) error {
	m, err := migrate.New(
		a.config.SqlSchemaDir,
		db.Config().ConnString(),
	)
	if err != nil {
		return err
	}
	if a.config.ForceVersion >= -1 {
		err = m.Force(a.config.ForceVersion)
		if err != nil {
			return err
		}
		err = m.Down()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}
	return m.Up()
}
