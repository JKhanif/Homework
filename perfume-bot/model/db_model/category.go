package db_model

type Category struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
}
