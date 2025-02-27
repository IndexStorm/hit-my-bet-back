package prediction

import (
	"context"
	"github.com/IndexStorm/hit-my-bet-back/pkg/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	db.BaseRepository
}

func NewPostgres(pool *pgxpool.Pool) Repository {
	return &postgres{
		BaseRepository: db.NewPostgresBaseRepository(pool),
	}
}

func (p *postgres) CreateMarket(ctx context.Context, market Market) error {
	const CreateMarketQuery = `INSERT INTO prediction.markets
(id,
 chain_status,
 title,
 creator_pubkey,
 resolver_pubkey,
 resolution,
 description,
 created_at,
 open_through)
VALUES (@id,
        @chain_status,
        @title,
        @creator_pubkey,
        @resolver_pubkey,
        @resolution,
        @description,
        @created_at,
        @open_through);`
	conn := p.GetConnectionFromCtx(ctx)
	_, err := conn.Exec(ctx, CreateMarketQuery, pgx.NamedArgs{
		"id":              market.ID,
		"chain_status":    market.ChainStatus,
		"title":           market.Title,
		"description":     market.Description,
		"creator_pubkey":  market.CreatorPubkey,
		"resolver_pubkey": market.ResolverPubkey,
		"resolution":      market.Resolution,
		"created_at":      market.CreatedAt,
		"open_through":    market.OpenThrough,
	})
	return err
}

func (p *postgres) SetMarketInitialized(ctx context.Context, market string) error {
	const SetMarketInitializedQuery = `UPDATE prediction.markets
SET
  chain_status = $1
WHERE
  id = $2;`
	conn := p.GetConnectionFromCtx(ctx)
	_, err := conn.Exec(ctx, SetMarketInitializedQuery, MarketChainStatusConfirmed, market)
	return err
}
