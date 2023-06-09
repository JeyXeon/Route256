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

	query, args, err := queryBuilder().
		Select(stockSkuColumn, stockWarehouseIdColumn, stockCountColumn).
		From(stockTable).
		Where(sq.Eq{stockSkuColumn: skus}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var stocks []*schema.Stock
	if err := pgxscan.Select(ctx, db, &stocks, query, args...); err != nil {
		return nil, err
	}

	return converter.SchemaToStockListModel(stocks), nil
}

func (r *stocksRepository) WriteOffStocks(ctx context.Context, stocks []*model.Stock) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	for _, stock := range stocks {
		query, args, err := queryBuilder().
			Update(stockTable).
			Set(stockCountColumn, sq.ConcatExpr(stockCountColumn, sq.Expr(" - ?", stock.Count))).
			Where(sq.Eq{stockSkuColumn: stock.Sku, stockWarehouseIdColumn: stock.WareHouseId}).
			ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *stocksRepository) RevertReservations(ctx context.Context, reservations []*model.Reservation) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	for _, reservation := range reservations {
		query, args, err := queryBuilder().
			Update(stockTable).
			Set(stockCountColumn, sq.ConcatExpr(stockCountColumn, sq.Expr(" + ?", reservation.Count))).
			Where(sq.Eq{stockSkuColumn: reservation.Sku, stockWarehouseIdColumn: reservation.WareHouseId}).
			ToSql()
		if err != nil {
			return err
		}

		_, err = db.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
