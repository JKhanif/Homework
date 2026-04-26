package db_model

type Product_category struct {
	Id          int `db:"id"`
	Product_id  int `db:"product_id"`
	Category_id int
}
