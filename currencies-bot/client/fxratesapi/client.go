package fxratesapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func New(token string) *FraClient {
	return &FraClient{
		token: token,
	}
}

type FraLatesResp struct {
	Rates FraRates `json:"rates"`
}

type FraRates struct {
	Rub float32 `json:"RUB"`
	Usd float32 `json:"USD"`
	Eur float32 `json:"EUR"`
	Sar float32 `json:"SAR"`
	Try float32 `json:"TRY"`
}

type FraClient struct {
	token string
}

func (c *FraClient) GetCurrencyRate(ctx context.Context, base string) (map[string]float32, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.fxratesapi.com/latest?api_key=%s&base=%s", c.token, base))
	if err != nil {
		return nil, fmt.Errorf("error get latest: %w", err)
	}

	var ratesResp FraLatesResp
	err = json.NewDecoder(resp.Body).Decode(&ratesResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal latest: %w", err)
	}

	return map[string]float32{
		"RUB": ratesResp.Rates.Rub,
		"USD": ratesResp.Rates.Usd,
		"EUR": ratesResp.Rates.Eur,
		"TRY": ratesResp.Rates.Try,
		"SAR": ratesResp.Rates.Sar,
	}, nil
}
