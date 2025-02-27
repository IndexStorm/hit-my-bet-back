package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IndexStorm/hit-my-bet-back/internal/repository/prediction"
	"github.com/IndexStorm/hit-my-bet-back/pkg/nanoid"
	"github.com/gagliardetto/solana-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"time"
)

func (s *server) createMarket(c *fiber.Ctx) error {
	type MarketData struct {
		Title       string `json:"title"`
		Creator     string `json:"creator"`
		Description string `json:"description"`
		OpenThrough int64  `json:"openThrough"`
	}
	type Request struct {
		RawData   string `json:"rawData"`
		Signature []byte `json:"signature"`
	}
	var request Request
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return fmt.Errorf("unmarshal request: %w", err)
	}
	var marketData MarketData
	if err := json.Unmarshal([]byte(request.RawData), &marketData); err != nil {
		return fmt.Errorf("unmarshal market data: %w", err)
	}
	creatorPubkey, err := solana.PublicKeyFromBase58(marketData.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator pubkey: %w", err)
	}
	if !creatorPubkey.Verify([]byte(request.RawData), solana.SignatureFromBytes(request.Signature)) {
		return fiber.NewError(fiber.StatusUnauthorized, "signature is not valid")
	}
	openThrough := time.Unix(marketData.OpenThrough/1000, 0)
	if time.Now().After(openThrough) {
		return fiber.NewError(fiber.StatusBadRequest, "market is closed")
	}
	market := prediction.Market{
		ID:             nanoid.RandomID(),
		ChainStatus:    prediction.MarketChainStatusPending,
		Title:          marketData.Title,
		Description:    zeronull.Text(marketData.Description),
		CreatorPubkey:  marketData.Creator,
		ResolverPubkey: marketData.Creator,
		Resolution:     prediction.MarketResolutionUnresolved,
		CreatedAt:      time.Now(),
		OpenThrough:    openThrough,
	}
	ctx := c.UserContext()
	if err = s.predictionRepo.CreateMarket(ctx, market); err != nil {
		return fmt.Errorf("create market: %w", err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": market.ID})
}

func (s *server) initMarket(c *fiber.Ctx) error {
	type Request struct {
		MarketID string `json:"marketID"`
		TxData   string `json:"txData"`
	}
	var request Request
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return fmt.Errorf("unmarshal request: %w", err)
	}
	txHash, err := s.relayTxData(request.TxData)
	if err != nil {
		return fmt.Errorf("relay tx: %w", err)
	}
	if txHash == "" {
		return errors.New("tx hash is empty")
	}
	// TODO: Use chain monitoring to confirm markets
	// TODO: Validate market id from tx data
	ctx := c.UserContext()
	if err = s.predictionRepo.SetMarketInitialized(ctx, request.MarketID); err != nil {
		s.logger.Err(err).Msg("failed to set market initialized")
	}
	return c.SendStatus(fiber.StatusOK)
}
