package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// TronGrid -> MVP level!

type Config struct {
	Network Network
	APIKey  string
}

type Client struct {
	http *http.Client
	cfg  Config
}

func NewClient(cfg Config) *Client {
	return &Client{http: &http.Client{
		Timeout: 15 * time.Second,
	}, cfg: cfg}
}

func (c *Client) raw(ctx context.Context, method, path string, by io.Reader) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, method, c.cfg.Network.HEndpoint+path, by)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Set("API-KEY", c.cfg.APIKey)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"tron-grid returned %d: %s",
			res.StatusCode,
			body,
		)
	}
	return body, nil
}

type TransResponse struct {
	Data []struct {
		TransactionId string `json:"transaction_id"`
		TokenInfo     struct {
			Symbol   string `json:"symbol"`
			Address  string `json:"address"`
			Decimals int    `json:"decimals"`
			Name     string `json:"name"`
		} `json:"token_info"`
		BlockTimestamp int64  `json:"block_timestamp"`
		From           string `json:"from"`
		To             string `json:"to"`
		Type           string `json:"type"`
		Value          string `json:"value"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At       int64 `json:"at"`
		PageSize int   `json:"page_size"`
	} `json:"meta"`
}

func (c *Client) Transactions(ctx context.Context, address string, minTimestamp int64) ([]Transaction, error) {

	path := fmt.Sprintf("/v1/accounts/%s/transactions/trc20?only_confirmed=true&only_to=true&contract_address=%s&min_timestamp=%d", address, c.cfg.Network.Contract, minTimestamp)

	body, rawErr := c.raw(ctx, http.MethodGet, path, nil)
	if rawErr != nil {
		return nil, rawErr
	}

	var trans TransResponse
	err := json.Unmarshal(body, &trans)
	if err != nil {
		return nil, err
	}

	var res []Transaction

	if trans.Success {
		for _, item := range trans.Data {
			value, err := decimal.NewFromString(item.Value)
			if err != nil {
				return nil, err
			}

			it := Transaction{
				Hash:      item.TransactionId,
				Sender:    item.From,
				Receiver:  item.To,
				Amount:    value.Shift(-6),
				Timestamp: time.UnixMilli(item.BlockTimestamp),
			}

			res = append(res, it)
		}
	}
	return res, nil
}
