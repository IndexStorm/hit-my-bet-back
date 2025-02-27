package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepository interface {
	RunInTx(ctx context.Context, fn func(context.Context) error) error
	GetConnectionFromCtx(ctx context.Context) PgxConnection
}

func NewPostgresBaseRepository(db *pgxpool.Pool) BaseRepository {
	return &postgresBaseRepository{db: db}
}

type postgresBaseRepository struct {
	db *pgxpool.Pool
}

func (r *postgresBaseRepository) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	ctx = context.WithValue(ctx, PgxConnectionCtxKey{}, tx)
	err = fn(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *postgresBaseRepository) GetConnectionFromCtx(ctx context.Context) PgxConnection {
	conn, ok := ctx.Value(PgxConnectionCtxKey{}).(PgxConnection)
	if !ok {
		return r.db
	}
	return conn
}
