-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_photos (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id),
    is_main BOOLEAN DEFAULT false,
    url TEXT NOT NULL
);

CREATE UNIQUE INDEX one_primary_per_product ON product_photos (product_id) WHERE is_main = true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
