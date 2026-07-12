package tron

import (
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type Network string

const (
	Mainnet Network = "https://api.trongrid.io"
	Nila    Network = "https://nile.trongrid.io"
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

func (c *Client) Raw(method, url string) {
	req, _ := http.NewRequest(method, url, nil)
	res, _ := c.http.Do(req)
	defer res.Body.Close()
	_, _ = io.ReadAll(res.Body)

}
func (c *Client) GetBalance(address Address) (decimal.Decimal, error) {
	return decimal.Zero, nil
}
