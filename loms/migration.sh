goose -dir ./migrations postgres "postgres://postgres:password@localhost:5434/loms?sslmode=disable" status
goose -dir ./migrations postgres "postgres://postgres:password@localhost:5434/loms?sslmode=disable" up
