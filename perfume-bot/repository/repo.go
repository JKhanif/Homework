package repository

import (
	"context"
	"fmt"
	"perfume-bot/model/api_model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UpdateProductPartial(ctx context.Context, id int, req api_model.UpdateProductRequest) error {
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
