package prediction

import (
	"context"
	"github.com/IndexStorm/hit-my-bet-back/pkg/db"
)

type Repository interface {
	db.BaseRepository

	CreateMarket(ctx context.Context, market Market) error
	SetMarketInitialized(ctx context.Context, market string) error
}
