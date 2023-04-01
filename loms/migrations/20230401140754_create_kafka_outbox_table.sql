-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox_record (
    outbox_record_id SERIAL PRIMARY KEY,
    topic TEXT NOT NULL,
    key TEXT NOT NULL,
    message TEXT NOT NULL,
    state INT NOT NULL,
    created_on TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE outbox_record;
-- +goose StatementEnd
