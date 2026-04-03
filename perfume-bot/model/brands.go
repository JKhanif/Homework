package model

import "time"

type Brand struct {
	ID          int
	Title       string
	Description string
	CreatedAt   time.Time
}
