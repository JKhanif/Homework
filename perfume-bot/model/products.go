package model

import "time"

type Product struct {
	ID           int
	Title        string
	PriceKopecks int
	CreatedAt    time.Time
	Brand        Brand
	Description  string
}
