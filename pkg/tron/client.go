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

func (c *Client) Raw(ctx context.Context, method, path string) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, method, string(c.cfg.Network)+path, nil)
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

	path := fmt.Sprintf("/v1/accounts/%s/trc20/balance", address)

	body, _ := c.Raw(ctx, http.MethodGet, path)

	var ac AccountResponse
	_ = json.Unmarshal(body, &ac)

	des := decimal.Zero
	if ac.Success {
		for _, f := range ac.Data {
			if nileUsdt := f[NileUSDTContract]; nileUsdt != "" {
				usdt, err := decimal.NewFromString(nileUsdt)
				if err != nil {
					return decimal.Zero, err
				}
				des = usdt.Shift(-6)
			}
		}
	}
	return des, nil
}

func (c *Client) Test(ctx context.Context, address Address) (decimal.Decimal, error) {

	path := fmt.Sprintf("/v1/accounts/%s/trc20/transactions/trc20", address)
	body, _ := c.Raw(ctx, http.MethodGet, path)

	fmt.Println(string(body))

	return decimal.Zero, nil
}
