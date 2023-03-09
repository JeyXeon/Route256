package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}
