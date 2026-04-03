package repository

import (
	"context"
	"fmt"
	"perfume-bot/model"

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

func (r *Repository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.Query(ctx, `
							SELECT 
								p.id, 
								p.title, 
								p.price, 
								p.created_at, 
								b.id, 
								b.title, 
								b.description 
							FROM products p 
							JOIN brands b 
							ON p.brand_id = b.id;`)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}

	defer rows.Close()

	products := make([]model.Product, 0)

	for rows.Next() {
		var p model.Product
		var b model.Brand
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.PriceKopecks,
			&p.CreatedAt,
			&b.ID,
			&b.Title,
			&b.Description)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}
		p.Brand = b
		products = append(products, p)
	}

	return products, nil
}

// SELECT p.id, p.title, p.description, p.price, b.title, FROM products p JOIN producrs_categories pc ON p.id = pc.product_id JOIN brands ON b.id = p.brand_id WHERE pc.category_id = $1;

func (r *Repository) GetProductsByCategoryID(ctx context.Context, categoryID string) ([]model.Product, error) {
	rows, err := r.db.Query(ctx, `SELECT 
    								p.id,
    								p.title,
    								p.description,
    								p.price,
    								b.id,
									b.title,
									b.description
								FROM products p
								JOIN brands b ON b.id = p.brand_id
								JOIN product_categories pc ON pc.product_id = p.id
								WHERE pc.category_id = $1;`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}

	products := make([]model.Product, 0)

	for rows.Next() {
		var p model.Product
		var b model.Brand
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.PriceKopecks,
			&b.ID,
			&b.Title,
			&b.Description)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		p.Brand = b

		products = append(products, p)
	}

	return products, nil
}

func (r *Repository) GetAllCategories(ctx context.Context) ([]model.Categories, error) {
	rows, err := r.db.Query(ctx, `SELECT id, title FROM categories`)
	if err != nil {
		return nil, fmt.Errorf("Error db.Query: %w", err)
	}

	defer rows.Close()

	categories := make([]model.Categories, 0)

	for rows.Next() {
		var c model.Categories
		err := rows.Scan(&c.ID, &c.Title)
		if err != nil {
			return nil, fmt.Errorf("Error rows.Scan: %w", err)
		}

		categories = append(categories, c)
	}

	return categories, nil
}
