package tron

// Balance

type AccountResponse struct {
	Success bool          `json:"success"`
	Data    []TokenAmount `json:"data"`
	Meta    Meta          `json:"meta"`
}

type TokenAmount map[string]string

type Meta struct {
	At       int64 `json:"at"`
	PageSize int   `json:"page_size"`
}

// Transactions

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
