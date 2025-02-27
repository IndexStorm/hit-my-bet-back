package db

import (
	"context"
	"crypto/x509"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"runtime"
	"time"
)

type PgxConnection interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	CopyFrom(ctx context.Context, table pgx.Identifier, columns []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type PgxConnectionCtxKey struct{}

var ErrTlsConfigRequired = errors.New("pgx:TLSConfig is required with CertPool")

func NewPgxConnection(
	ctx context.Context,
	conn string,
	tracer pgx.QueryTracer,
	certPool *x509.CertPool,
	timeout time.Duration,
) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}
	config.ConnConfig.ConnectTimeout = timeout
	config.ConnConfig.Tracer = tracer
	if certPool != nil {
		if config.ConnConfig.TLSConfig == nil {
			return nil, ErrTlsConfigRequired
		}
		config.ConnConfig.TLSConfig.RootCAs = certPool
		config.ConnConfig.TLSConfig.InsecureSkipVerify = false
	}
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnLifetimeJitter = time.Second * 15
	config.MaxConnIdleTime = time.Minute
	config.MaxConns = int32(4 * runtime.NumCPU())
	config.MinConns = int32(runtime.NumCPU())
	config.HealthCheckPeriod = time.Minute
	return pgxpool.NewWithConfig(ctx, config)
}
