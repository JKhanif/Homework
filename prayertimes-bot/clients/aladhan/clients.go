package aladhan

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"prayertimes/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	httpClient http.Client
	rdb        *redis.Client
}

func New(rdb *redis.Client) *Client {
	return &Client{
		httpClient: http.Client{},
		rdb:        rdb,
	}
}

func (c *Client) GetTodayPrayerTimesByCity(ctx context.Context, city string) (models.AladhanResponse, error) {
	resJSON, err := c.rdb.Get(ctx, fmt.Sprintf("prayers_%s", city)).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("error get cached prayer times:", err)
	}

	if resJSON != "" {
		var res models.AladhanResponse
		err := json.Unmarshal([]byte(resJSON), &res)
		if err != nil {
			fmt.Println("error marshal res for cache:", err)
		}
		return res, nil
	}

	host := "https://api.aladhan.com/v1/timingsByCity"

	u, _ := url.Parse(host)
	u = u.JoinPath(time.Now().Format("02-01-2006"))

	q := u.Query()
	q.Set("city", city)
	q.Set("country", "RU")
	q.Set("method", "4")
	q.Set("timezonestring", "Europe/Moscow") // TODO: определять по городу
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

	defer resp.Body.Close()

	var prayertimes models.AladhanResponse
	err = json.NewDecoder(resp.Body).Decode(&prayertimes)
	if err != nil {
		return models.AladhanResponse{}, fmt.Errorf("decode error: %w", err)
	}

	prayertimesJSON, err := json.Marshal(prayertimes)
	if err != nil {
		fmt.Println("error marshal prayertimesJSON for cache:", err)
		return models.AladhanResponse{}, err
	}

	date := time.Now().Format("02.01")

	cmd := c.rdb.Set(ctx, fmt.Sprintf("prayers_%s_%s", city, date), prayertimesJSON, 6*time.Hour)
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		return models.AladhanResponse{}, err
	}

	return prayertimes, nil
}
