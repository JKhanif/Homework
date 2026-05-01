package repository

// CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY

const queryCreateCategory = `INSERT INTO categories (title)
							VALUES ($1)
							RETURNING id`

const queryDeleteCategory = `DELETE FROM categories WHERE id = $1`

const queryGetCategoryByID = `SELECT id, title FROM categories WHERE id = $1`

const queryGetAllCategories = `SELECT id, title FROM categories`

const queryGetProductsByCategoryID = `SELECT 
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
										AND pp.is_main = true;`

const queryDeleteProductCategories = `DELETE FROM product_categories WHERE product_id = $1`

const queryInsertProductCategories = `INSERT INTO product_categories (product_id, category_id) SELECT $1, unnest($2::int[])`

// BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND BRAND

const queryCreateBrand = `INSERT INTO brands (title, description)
							VALUES ($1, $2)
							RETURNING id`

const queryDeleteBrand = `DELETE FROM brands WHERE id = $1`

const queryGetBrandByID = `SELECT FROM brands WHERE id = $1`

const queryGetAllBrands = `SELECT id, title, description FROM brands`

const queryGetProductsByBrandID = `SELECT
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
										`

// PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT PRODUCT

const queryGetAllProduct = `SELECT 
								p.id, 
								p.title, 
								p.price, 
								p.brand_id,
								p.description,
								p.created_at, 
								b.id, 
								b.title, 
								b.description,
								pp.tg_file_id
							FROM products p 
							JOIN brands b ON b.id = p.brand_id
							JOIN product_photos pp 
							  ON pp.product_id = p.id 
							  AND pp.is_main = true;`

const queryGetProductByID = `SELECT
								p.id, p.title, p.description, p.price, p.brand_id, p.created_at,
								b.id, b.title, b.description
							FROM products p
							JOIN brands b ON b.id = p.brand_id
							WHERE p.id = $1`

const queryCreateProduct = `INSERT INTO products (title, description, price, brand_id)
							VALUES ($1, $2, $3, $4)
							RETURNING id`

const queryDeleteProduct = `DELETE FROM products WHERE id=$1`
