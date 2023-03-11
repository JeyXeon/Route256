package postgres

import (
	"context"
	"route256/libs/dbmanager"

	"github.com/jackc/pgx/v4"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) dbmanager.QueryEngine
}
