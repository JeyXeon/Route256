-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reservation(
    order_id BIGINT NOT NULL REFERENCES user_order,
    sku INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    count INTEGER NOT NULL,
    PRIMARY KEY (order_id, sku, warehouse_id)
);

CREATE TABLE IF NOT EXISTS stock(
   sku INTEGER NOT NULL,
   warehouse_id INTEGER NOT NULL,
   count INTEGER NOT NULL,
   PRIMARY KEY (sku, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reservation;
DROP TABLE IF EXISTS stock;
-- +goose StatementEnd
