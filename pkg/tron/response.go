package tron

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
