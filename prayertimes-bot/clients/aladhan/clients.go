package aladhan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"prayertimes/models"
	"time"
)

type Client struct {
	httpClient http.Client
}

func New() *Client {
	return &Client{
		httpClient: http.Client{},
	}
}

func (c *Client) GetTodayPrayerTimesByCity(ctx context.Context, city string) (models.AladhanResponse, error) {
	host := "https://api.aladhan.com/v1/timingsByCity"

	u, _ := url.Parse(host)
	u.JoinPath(time.Now().Format("02-01-2006"))
	q := u.Query()
	q.Set("city", city)
	q.Set("country", "RU")
	q.Set("method", "4")
	q.Set("timezonestring", "Europe/Moscow") // TODO: определять по стране
	q.Set("calendarMethod", "UAQ")

	u.RawQuery = q.Encode()
	fmt.Println("Final URL: ", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return models.AladhanResponse{}, fmt.Errorf("error create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.AladhanResponse{}, fmt.Errorf("error do request: %w", err)
	}

	var res models.AladhanResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return models.AladhanResponse{}, fmt.Errorf("decode error: %w", err)
	}

	return res, nil
}
