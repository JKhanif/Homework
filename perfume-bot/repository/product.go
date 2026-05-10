package repository

import (
	"context"
	"errors"
	"fmt"
	"perfume-bot/model/api_model"
	db_model "perfume-bot/model/db_model"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetProductsByBrandIDPage(ctx context.Context, brandID int, limit int, offset int) ([]db_model.Product, error) {
	rows, err := r.db.Query(ctx, queryGetProductsByBrandPage, brandID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]db_model.Product, 0)
	for rows.Next() {
		var p db_model.Product
		var brandID *int
		var brandTitle, brandDesc *string
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.Price,
			&brandID, &brandTitle, &brandDesc,
			&p.MainPhotoFailID,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}
		if brandID != nil {
			p.Brand.ID = *brandID
			p.Brand.Title = *brandTitle
			p.Brand.Description = brandDesc
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) GetProductsByCategoryIDPage(ctx context.Context, categoryID int, limit int, offset int) ([]db_model.Product, error) {
	rows, err := r.db.Query(ctx, queryGetProductsByCategoryPage, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]db_model.Product, 0)
	for rows.Next() {
		var p db_model.Product
		var brandID *int
		var brandTitle, brandDesc *string
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.Price,
			&brandID, &brandTitle, &brandDesc,
			&p.MainPhotoFailID,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}
		if brandID != nil {
			p.Brand.ID = *brandID
			p.Brand.Title = *brandTitle
			p.Brand.Description = brandDesc
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) GetAllProductsPage(ctx context.Context, limit int, offset int) ([]db_model.Product, error) {
	rows, err := r.db.Query(ctx, queryGetAllProductsPage, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]db_model.Product, 0)
	for rows.Next() {
		var p db_model.Product
		var pBrandID *int
		var brandID *int
		var brandTitle, brandDesc *string
		err := rows.Scan(
			&p.ID, &p.Title, &p.Price, &pBrandID, &p.Description,
			&brandID, &brandTitle, &brandDesc,
			&p.MainPhotoFailID, &p.MainPhotoURL,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}
		if brandID != nil {
			p.Brand.ID = *brandID
			p.Brand.Title = *brandTitle
			p.Brand.Description = brandDesc
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) GetAllProducts(ctx context.Context) ([]db_model.Product, error) {
	rows, err := r.db.Query(ctx, queryGetAllProduct)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	products := make([]db_model.Product, 0)

	for rows.Next() {
		var p db_model.Product
		var pBrandID *int
		var brandID *int
		var brandTitle, brandDesc *string
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Price,
			&pBrandID,
			&p.Description,
			&p.CreatedAt,
			&brandID,
			&brandTitle,
			&brandDesc,
			&p.MainPhotoFailID,
			&p.MainPhotoURL,
		)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}
		if brandID != nil {
			p.Brand.ID = *brandID
			p.Brand.Title = *brandTitle
			p.Brand.Description = brandDesc
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) GetProductByID(ctx context.Context, id int) (db_model.Product, error) {
	var p db_model.Product
	var pBrandID *int
	var brandID *int
	var brandTitle, brandDesc *string

	err := r.db.QueryRow(ctx, queryGetProductByID, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.Price, &pBrandID, &p.CreatedAt,
		&brandID, &brandTitle, &brandDesc)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, fmt.Errorf("Не найдено")
		}
		return p, fmt.Errorf("db error: %w", err)
	}

	if brandID != nil {
		p.Brand.ID = *brandID
		p.Brand.Title = *brandTitle
		p.Brand.Description = brandDesc
	}

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

func (r *Repository) CreateProductPhoto(ctx context.Context, productID int, url string, tgFileID string, isMain bool) (int, error) {
	var id int
	err := r.db.QueryRow(ctx, queryCreateProductPhoto, productID, url, tgFileID, isMain).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Error db.QueryRow create product_photo: %w", err)
	}
	return id, nil
}

func (r *Repository) UnsetMainPhoto(ctx context.Context, productID int) error {
	_, err := r.db.Exec(ctx, queryUnsetMainPhoto, productID)
	if err != nil {
		return fmt.Errorf("Error db.Exec unset main photo: %w", err)
	}
	return nil
}

func (r *Repository) GetPhotoByID(ctx context.Context, photoID int) (*db_model.ProductPhoto, error) {
	var ph db_model.ProductPhoto
	err := r.db.QueryRow(ctx, queryGetPhotoByID, photoID).Scan(&ph.ID, &ph.ProductID, &ph.URL, &ph.TgFileID, &ph.IsMain)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Не найдено")
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &ph, nil
}

func (r *Repository) GetProductPhotos(ctx context.Context, productID int) ([]db_model.ProductPhoto, error) {
	rows, err := r.db.Query(ctx, queryGetPhotosByProductID, productID)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}
	defer rows.Close()

	photos := make([]db_model.ProductPhoto, 0)
	for rows.Next() {
		var ph db_model.ProductPhoto
		err := rows.Scan(&ph.ID, &ph.ProductID, &ph.URL, &ph.TgFileID, &ph.IsMain)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}
		photos = append(photos, ph)
	}
	return photos, nil
}

func (r *Repository) DeletePhoto(ctx context.Context, photoID int) error {
	result, err := r.db.Exec(ctx, queryDeletePhoto, photoID)
	if err != nil {
		return fmt.Errorf("Error db.Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("Не найдено")
	}
	return nil
}

func (r *Repository) DeleteProductPhotos(ctx context.Context, productID int) error {
	_, err := r.db.Exec(ctx, queryDeleteProductPhotos, productID)
	if err != nil {
		return fmt.Errorf("Error db.Exec delete product photos: %w", err)
	}
	return nil
}

func (r *Repository) SetMainPhoto(ctx context.Context, photoID int, productID int) error {
	_, err := r.db.Exec(ctx, queryUnsetMainPhoto, productID)
	if err != nil {
		return fmt.Errorf("Error db.Exec unset main: %w", err)
	}
	_, err = r.db.Exec(ctx, querySetMainPhoto, photoID, productID)
	if err != nil {
		return fmt.Errorf("Error db.Exec set main: %w", err)
	}
	return nil
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
