goose -dir ./migrations postgres "postgres://postgres:password@localhost:5433/checkout?sslmode=disable" status
goose -dir ./migrations postgres "postgres://postgres:password@localhost:5433/checkout?sslmode=disable" up
