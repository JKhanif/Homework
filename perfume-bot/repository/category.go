package repository

import (
	"context"
	"errors"
	"fmt"
	"perfume-bot/model/api_model"
	db_model "perfume-bot/model/db_model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateCategory(ctx context.Context, c api_model.CreateCategoryRequest) (int, error) {
	var id int

	err := r.db.QueryRow(ctx, queryCreateCategory, c.Title).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error db.QueryRow: %w", err)
	}

	return id, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, id int, req api_model.UpdateCategoryRequest) error {
	query := "UPDATE categories SET "
	args := make([]interface{}, 0)
	argPos := 1

	if req.Title != nil {
		query += fmt.Sprintf("title=$%d,", argPos)
		args = append(args, *req.Title)
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
		return fmt.Errorf("update category: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("Категория не найдена")
	}

	return nil
}

func (r *Repository) DeleteCategory(ctx context.Context, id int) error {
	result, err := r.db.Exec(ctx, queryDeleteCategory, id)
	if err != nil {
		return fmt.Errorf("Error db.Exec: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Нет в базе данных")
	}

	return nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, id int) (api_model.CategoryResponse, error) {
	var c api_model.CategoryResponse
	err := r.db.QueryRow(ctx, queryGetCategoryByID, id).Scan(&c.ID, &c.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c, fmt.Errorf("Не найдено")
		}
		return c, fmt.Errorf("db error: %w", err)
	}

	return c, nil
}

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
		var brandID *int
		var brandTitle, brandDesc *string
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Price,
			&brandID,
			&brandTitle,
			&brandDesc,
			&p.MainPhotoFailID,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		if brandID != nil {
			b.ID = *brandID
			b.Title = *brandTitle
			b.Description = brandDesc
		}
		p.Brand = b

		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) SetProductCategories(ctx context.Context, productID int, categoryIDs []int) error {
	// Удаление всех категорий продукта
	_, err := r.db.Exec(ctx, queryDeleteProductCategories, productID)
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
