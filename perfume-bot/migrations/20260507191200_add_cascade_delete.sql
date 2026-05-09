-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
  DROP CONSTRAINT products_brand_id_fkey,
  ADD CONSTRAINT products_brand_id_fkey
    FOREIGN KEY (brand_id) REFERENCES brands(id)
    ON DELETE SET NULL;

ALTER TABLE product_categories
  DROP CONSTRAINT product_categories_category_id_fkey,
  ADD CONSTRAINT product_categories_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES categories(id)
    ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
  DROP CONSTRAINT products_brand_id_fkey,
  ADD CONSTRAINT products_brand_id_fkey
    FOREIGN KEY (brand_id) REFERENCES brands(id);

ALTER TABLE product_categories
  DROP CONSTRAINT product_categories_category_id_fkey,
  ADD CONSTRAINT product_categories_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES categories(id);
-- +goose StatementEnd
