-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_item(
    user_id BIGSERIAL NOT NULL,
    sku INTEGER NOT NULL,
    count INTEGER NOT NULL,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_item;
-- +goose StatementEnd
