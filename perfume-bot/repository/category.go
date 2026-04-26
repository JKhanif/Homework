package repository

import (
	"context"
	"fmt"
	db_model "perfume-bot/model/db_model"
)

func (r *Repository) GetAllCategories(ctx context.Context) ([]db_model.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title FROM categories`)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	categories := make([]db_model.Category, 0)

	for rows.Next() {
		var c db_model.Category
		err := rows.Scan(&c.ID, &c.Title)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func (r *Repository) GetProductsByCategoryID(ctx context.Context, categoryID string) ([]db_model.Product, error) {
	rows, err := r.db.Query(ctx, `SELECT 
    								p.id,
    								p.title,
    								p.description,
    								p.price,
    								b.id,
									b.title,
									b.description,
									pp.tg_file_id
								FROM products p
								JOIN brands b ON b.id = p.brand_id
								JOIN product_categories pc ON pc.product_id = p.id
								JOIN product_photos pp ON pp.product_id = p.id
								WHERE pc.category_id = $1
								  AND pp.is_main = true;`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]db_model.Product, 0)

	for rows.Next() {
		var p db_model.Product
		var b db_model.Brand
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Price,
			&b.ID,
			&b.Title,
			&b.Description,
			&p.MainPhotoFailID,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		p.Brand = b

		products = append(products, p)
	}

	return products, nil
}
