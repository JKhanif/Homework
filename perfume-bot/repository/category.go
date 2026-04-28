package repository

import (
	"context"
	"fmt"
	db_model "perfume-bot/model/db_model"
)

func (r *Repository) GetAllCategories(ctx context.Context) ([]db_model.Category, error) {
	rows, err := r.db.Query(ctx, queryGetAllCategories)
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
	rows, err := r.db.Query(ctx, queryGetProductsByCategoryID, categoryID)
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

func (r *Repository) SetProductCategories(ctx context.Context, productID int, categoryIDs []int) error {
	// Удаление всех категорий
	_, err := r.db.Exec(ctx, queryDeleteProductCategories)
	if err != nil {
		return fmt.Errorf("failed to delete product categories before update: %w", err)
	}

	//Добавление новых категорий
	_, err = r.db.Exec(ctx, queryInsertProductCategories, productID, categoryIDs)
	if err != nil {
		return fmt.Errorf("failed to insert product category: %w", err)
	}

	return nil
}
