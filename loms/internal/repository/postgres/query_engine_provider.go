package postgres

import (
	"context"
	"route256/libs/dbmanager"

	"github.com/jackc/pgx/v4"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) dbmanager.QueryEngine
}
