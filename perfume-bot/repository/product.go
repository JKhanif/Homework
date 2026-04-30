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

	err := r.db.QueryRow(ctx, queryGetProductByID, id).Scan(
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

	err := r.db.QueryRow(ctx, queryCreateProduct, p.Title, p.Description, p.Price, p.BrandID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error db.QueryRow: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, id int, req api_model.UpdateProductRequest) error {
	query := "UPDATE products SET "
	args := []interface{}{}
	i := 1

	if req.Title != nil {
		query += fmt.Sprintf("title=$%d,", i)
		args = append(args, *req.Title)
		i++
	}

	if req.Description != nil {
		query += fmt.Sprintf("description=$%d,", i)
		args = append(args, *req.Description)
		i++
	}

	if req.Price != nil {
		query += fmt.Sprintf("price=$%d,", i)
		args = append(args, *req.Price)
		i++
	}

	if req.BrandID != nil {
		query += fmt.Sprintf("brand_id=$%d,", i)
		args = append(args, *req.BrandID)
		i++
	}

	// убираем последнюю запятую
	query = query[:len(query)-1]

	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) DeleteProduct(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, queryDeleteProduct, id)
	if err != nil {
		return fmt.Errorf("Error db.Exec: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Нет в базе данных")
	}

	return nil
}
