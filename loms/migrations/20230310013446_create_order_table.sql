-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status as enum (
    'NEW',
    'AWAITING_PAYMENT',
    'FAILED',
    'PAYED',
    'CANCELLED'
    );

CREATE TABLE IF NOT EXISTS user_order(
    order_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    status order_status NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_order;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
