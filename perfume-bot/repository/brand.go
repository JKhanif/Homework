package repository

import (
	"context"
	"fmt"
	db_model "perfume-bot/model/db_model"
)

func (r *Repository) GetAllBrands(ctx context.Context) ([]db_model.Brand, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title, description FROM brands`)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	brands := make([]db_model.Brand, 0)

	for rows.Next() {
		var b db_model.Brand
		err := rows.Scan(&b.ID, &b.Title, &b.Description)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		brands = append(brands, b)
	}

	return brands, nil
}

func (r *Repository) GetProductsByBrandID(ctx context.Context, brandID string) ([]db_model.Product, error) {
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
								JOIN product_photos pp ON pp.product_id = p.id 
								WHERE p.brand_id = $1
								  AND pp.is_main = true;
									`, brandID)
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
