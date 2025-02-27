package prediction

import (
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"time"
)

type MarketChainStatus string
type MarketResolution string

const (
	MarketChainStatusPending   MarketChainStatus = "PENDING"
	MarketChainStatusNeedRetry MarketChainStatus = "NEED_RETRY"
	MarketChainStatusConfirmed MarketChainStatus = "CONFIRMED"

	MarketResolutionUnresolved MarketResolution = "UNRESOLVED"
	MarketResolutionTie        MarketResolution = "TIE"
	MarketResolutionYes        MarketResolution = "YES"
	MarketResolutionNo         MarketResolution = "NO"
)

type Market struct {
	ID             string            `db:"id" json:"id,omitempty"`
	ChainStatus    MarketChainStatus `db:"chain_status" json:"chain_status,omitempty"`
	Title          string            `db:"title" json:"title,omitempty"`
	Description    zeronull.Text     `db:"description" json:"description,omitempty"`
	CreatorPubkey  string            `db:"creator_pubkey" json:"creator_pubkey,omitempty"`
	ResolverPubkey string            `db:"resolver_pubkey" json:"resolver_pubkey,omitempty"`
	MarketPubkey   string            `db:"market_pubkey" json:"market_pubkey,omitempty"`
	Resolution     MarketResolution  `db:"resolution" json:"resolution,omitempty"`
	CreatedAt      time.Time         `db:"created_at" json:"created_at,omitempty"`
	OpenThrough    time.Time         `db:"open_through" json:"open_through,omitempty"`
}
