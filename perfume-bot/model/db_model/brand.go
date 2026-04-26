package db_model

type Brand struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}
