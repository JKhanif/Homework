-- +goose Up
-- +goose StatementBegin
ALTER TABLE product_photos
ADD COLUMN tg_file_id TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
