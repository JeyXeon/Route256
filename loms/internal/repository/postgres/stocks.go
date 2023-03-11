package postgres

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type stocksRepository struct {
	queryEngineProvider QueryEngineProvider
}

const (
	stockTable = "stock"

	stockSkuColumn         = "sku"
	stockWarehouseIdColumn = "warehouse_id"
	stockCountColumn       = "count"
)

func NewStocksRepository(queryEngineProvider QueryEngineProvider) *stocksRepository {
	return &stocksRepository{
		queryEngineProvider: queryEngineProvider,
	}
}

func (r *stocksRepository) GetStocks(ctx context.Context, skus []uint32) ([]*model.Stock, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := sq.
		Select(stockSkuColumn, stockWarehouseIdColumn, stockCountColumn).
		From(stockTable).
		Where(sq.Eq{stockSkuColumn: skus}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var stocks []*schema.Stock
	if err := pgxscan.Select(ctx, db, &stocks, query, args...); err != nil {
		return nil, err
	}

	return converter.ToStockListModel(stocks), nil
}

func (r *stocksRepository) WriteOffStocks(ctx context.Context, stocks []*model.Stock) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	for _, stock := range stocks {
		query, args, err := sq.
			Update(stockTable).
			Set(stockCountColumn, sq.ConcatExpr(stockCountColumn, sq.Expr(" - ?", stock.Count))).
			Where(sq.Eq{stockSkuColumn: stock.Sku, stockWarehouseIdColumn: stock.WareHouseId}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return err
		}

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		rows.Close()
	}

	return nil
}

func (r *stocksRepository) RevertReservations(ctx context.Context, reservations []*model.Reservation) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	for _, reservation := range reservations {
		query, args, err := sq.
			Update(stockTable).
			Set(stockCountColumn, sq.ConcatExpr(stockCountColumn, sq.Expr(" + ?", reservation.Count))).
			Where(sq.Eq{stockSkuColumn: reservation.Sku, stockWarehouseIdColumn: reservation.WareHouseId}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return err
		}

		rows, err := db.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		rows.Close()
	}

	return nil
}
