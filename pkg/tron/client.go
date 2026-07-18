package tron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type Network string

const (
	Mainnet          Network = "https://api.trongrid.io"
	Nila             Network = "https://nile.trongrid.io"
	NileUSDTContract         = "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf"
)

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

	req, err := http.NewRequestWithContext(ctx, method, string(c.cfg.Network)+path, by)
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

func (c *Client) GetBalance(ctx context.Context, address Address) (decimal.Decimal, error) {

	path := fmt.Sprintf("/v1/accounts/%s/trc20/balance?contract_address=%s", address, NileUSDTContract)

	body, rawErr := c.raw(ctx, http.MethodGet, path, nil)
	if rawErr != nil {
		return decimal.Zero, rawErr
	}

	var ac AccountResponse
	err := json.Unmarshal(body, &ac)
	if err != nil {
		return decimal.Zero, err
	}

	des := decimal.Zero
	if ac.Success {
		for key, value := range ac.Data[0] {
			if key == NileUSDTContract {

				usdt, uErr := decimal.NewFromString(value)
				if uErr != nil {
					return decimal.Zero, uErr
				}
				des = usdt.Shift(-6)
			}

		}

	}
	return des, nil
}

func (c *Client) Transactions(ctx context.Context, address Address, minTimestamp int64) (TransResponse, error) {

	path := fmt.Sprintf("/v1/accounts/%s/transactions/trc20?only_confirmed=true&only_to=true&contract_address=%s&min_timestamp=%d", address, NileUSDTContract, minTimestamp)

	body, rawErr := c.raw(ctx, http.MethodGet, path, nil)
	if rawErr != nil {
		return TransResponse{}, rawErr
	}

	var tr TransResponse
	err := json.Unmarshal(body, &tr)
	if err != nil {
		return TransResponse{}, err
	}
	return tr, nil
}
