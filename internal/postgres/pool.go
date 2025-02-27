package postgres

import (
	"context"
	"fmt"
	"github.com/IndexStorm/hit-my-bet-back/internal/config"
	"github.com/IndexStorm/hit-my-bet-back/internal/telemetry"
	"github.com/IndexStorm/hit-my-bet-back/pkg/db"
	"github.com/jackc/pgx/v5/pgxpool"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"time"
)

func NewPgxPoolWithOtel(
	ctx context.Context,
	dbConfig config.Database,
	env config.Environment,
) (*pgxpool.Pool, error) {
	sslMode := "require"
	if env == config.EnvironmentLocal {
		sslMode = "disable"
	}
	if dbConfig.SSLMode != "" {
		sslMode = dbConfig.SSLMode
	}
	database, err := db.NewPgxConnection(ctx,
		fmt.Sprintf(
			"postgresql://%s:%s@%s/%s?sslmode=%s",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Database,
			sslMode,
		),
		telemetry.NewPgxTracer(
			semconv.DBSystemPostgreSQL,
			semconv.DBNamespace(dbConfig.Username+"@"+dbConfig.Host+"/"+dbConfig.Database),
			// semconv.ServerAddress(config.Host),
			// semconv.ServerPort(int(config.Port)),
			// semconv.UserName(config.User),
			// semconv.DBNamespace(config.Database),
		),
		nil,
		time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err = database.Ping(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return database, nil
}

func NewPgxPool(
	ctx context.Context,
	dbConfig config.Database,
	env config.Environment,
) (*pgxpool.Pool, error) {
	sslMode := "require"
	if env == config.EnvironmentLocal {
		sslMode = "disable"
	}
	if dbConfig.SSLMode != "" {
		sslMode = dbConfig.SSLMode
	}
	database, err := db.NewPgxConnection(ctx,
		fmt.Sprintf(
			"postgresql://%s:%s@%s/%s?sslmode=%s",
			dbConfig.Username,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Database,
			sslMode,
		),
		nil,
		nil,
		time.Second*10)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err = database.Ping(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return database, nil
}
