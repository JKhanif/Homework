package repository

import (
	"context"
	"errors"
	"fmt"
	"perfume-bot/model/api_model"
	db_model "perfume-bot/model/db_model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateBrand(ctx context.Context, b api_model.CreateBrandRequest) (int, error) {
	var id int

	err := r.db.QueryRow(ctx, queryCreateBrand, b.Title, b.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error db.QueryRow: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateBrand(ctx context.Context, id int, req api_model.UpdateBrandRequest) error {
	query := "UPDATE brands SET "
	args := make([]interface{}, 0)
	argPos := 1

	if req.Title != nil {
		query += fmt.Sprintf("title=$%d,", argPos)
		args = append(args, *req.Title)
		argPos++
	}

	if req.Description != nil {
		query += fmt.Sprintf("description=$%d,", argPos)
		args = append(args, *req.Description)
		argPos++
	}

	// если ничего не передали
	if len(args) == 0 {
		return nil
	}

	// убрать последнюю запятую
	query = query[:len(query)-1]

	query += fmt.Sprintf(" WHERE id=$%d", argPos)
	args = append(args, id)

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update brand: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("Бренд не найден")
	}

	return nil
}

func (r *Repository) DeleteBrand(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, queryDeleteBrand, id)
	if err != nil {
		return fmt.Errorf("Error db.Exec: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Нет в базе данных")
	}

	return nil
}

func (r *Repository) GetBrandByID(ctx context.Context, id int) (api_model.BrandResponse, error) {
	var b api_model.BrandResponse
	err := r.db.QueryRow(ctx, queryGetBrandByID, id).Scan(&b.ID, &b.Title, &b.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return b, fmt.Errorf("Не найдено")
		}
		return b, fmt.Errorf("db error: %w", err)
	}

	return b, nil
}

func (r *Repository) GetAllBrands(ctx context.Context) ([]db_model.Brand, error) {
	rows, err := r.db.Query(ctx, queryGetAllBrands)
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
	rows, err := r.db.Query(ctx, queryGetProductsByBrandID, brandID)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]db_model.Product, 0)

	for rows.Next() {
		var p db_model.Product
		var b db_model.Brand
		var brandIDPTR *int
		var brandTitle, brandDesc *string
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Price,
			&brandIDPTR,
			&brandTitle,
			&brandDesc,
			&p.MainPhotoFailID,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		if brandIDPTR != nil {
			b.ID = *brandIDPTR
			b.Title = *brandTitle
			b.Description = brandDesc
		}
		p.Brand = b
		products = append(products, p)
	}

	return products, nil
}
