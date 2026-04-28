package repository

// CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY CATEGORY

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
