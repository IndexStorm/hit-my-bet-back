package main

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	"math/rand/v2"
	"time"
)

var solanaClient *req.Client

func init() {
	solanaClient = req.C().
		SetTimeout(time.Second * 15).
		SetJsonMarshal(json.Marshal).
		SetJsonUnmarshal(json.Unmarshal)
}

func (s *server) relayTx(c *fiber.Ctx) error {
	type Request struct {
		TxData string `json:"txData"`
	}
	var request Request
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return fmt.Errorf("unmarshal request: %w", err)
	}
	txHash, err := s.relayTxData(request.TxData)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"tx_hash": txHash})
}

func (s *server) relayTxData(data string) (string, error) {
	type RpcParams struct {
		Encoding   string `json:"encoding"`
		MaxRetries uint8  `json:"maxRetries"`
	}
	type RpcRequest struct {
		Jsonrpc string        `json:"jsonrpc"`
		Id      int           `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}
	rpcRequest := RpcRequest{
		Jsonrpc: "2.0",
		Id:      rand.Int(),
		Method:  "sendTransaction",
		Params:  []interface{}{data, RpcParams{Encoding: "base64", MaxRetries: 1}},
	}
	resp, err := solanaClient.R().SetBodyJsonMarshal(&rpcRequest).Post("https://api.devnet.solana.com")
	if err != nil {
		return "", fmt.Errorf("relay to solana: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, resp.String())
	}
	type Response struct {
		Jsonrpc string          `json:"jsonrpc"`
		Result  string          `json:"result"`
		Error   json.RawMessage `json:"error"`
	}
	var response Response
	if err = json.Unmarshal(resp.Bytes(), &response); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}
	if response.Error != nil {
		errData, err := json.Marshal(response.Error)
		if err != nil {
			return "", fmt.Errorf("failed: %+v", response.Error)
		}
		return "", errors.New(string(errData))
	}
	return response.Result, nil
}
