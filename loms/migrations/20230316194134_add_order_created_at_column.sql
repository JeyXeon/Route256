-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_order
ADD COLUMN created_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_order
DROP COLUMN created_at;
-- +goose StatementEnd
