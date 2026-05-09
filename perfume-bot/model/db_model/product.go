package db_model

import "time"

type Product struct {
	ID              int       `db:"id"`
	Title           string    `db:"title"`
	Price           int       `db:"price"`
	CreatedAt       time.Time `db:"created_at"`
	Brand           Brand     `db:"brand_id"`
	Description     string
	MainPhotoFailID *string
	MainPhotoURL    *string   `db:"url"`
}

type ProductPhoto struct {
	ID        int64  `db:"id"`
	ProductID int64  `db:"product_id"`
	IsMain    bool   `db:"is_main"`
	URL       string `db:"url"`
	TgFileID  string `db:"tg_file_id"`
}
