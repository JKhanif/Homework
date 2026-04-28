package repository

import (
	"context"
	"errors"
	"fmt"
	"perfume-bot/model/api_model"
	db_model "perfume-bot/model/db_model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetAllProducts(ctx context.Context) ([]db_model.Product, error) {
	rows, err := r.db.Query(ctx, queryGetAllProduct)
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
			&p.Price,
			&p.Brand.ID,
			&p.Description,
			&p.CreatedAt,
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

func (r *Repository) GetProductByID(ctx context.Context, id int) (db_model.Product, error) {
	var p db_model.Product
	var b db_model.Brand

	err := r.db.QueryRow(ctx, `
							SELECT
								p.id, p.title, p.description, p.price, p.brand_id, p.created_at,
								b.id, b.title, b.description
							FROM products p
							JOIN brands b ON b.id = p.brand_id
							WHERE p.id = $1
			 `, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.Price, &p.Brand.ID, &p.CreatedAt,
		&b.ID, &b.Title, &b.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, fmt.Errorf("Не найдено")
		}
		return p, fmt.Errorf("db error: %w", err)
	}

	p.Brand = b

	return p, nil
}

func (r *Repository) CreateProduct(ctx context.Context, p api_model.CreateProductRequest) (int, error) {
	var id int

	err := r.db.QueryRow(ctx, `
					INSERT INTO products (title, description, price, brand_id)
					VALUES ($1, $2, $3, $4)
					RETURNING id
					`, p.Title, p.Description, p.Price, p.BrandID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error db.QueryRow: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, p api_model.UpdateProductRequest) error {
	_, err := r.db.Exec(ctx, `
					UPDATE products
					SET title=$1, description=$2, price=$3, brand_id=$4
					WHERE id=$5
					`, p.Title, p.Description, p.Price, p.BrandID)
	if err != nil {
		return fmt.Errorf("Error db.Exec: %w", err)
	}

	return nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, `
		DELETE FROM products WHERE id=$1
	`, id)
	if err != nil {
		return fmt.Errorf("Error db.Exec: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Нет в базе данных")
	}

	return nil
}
