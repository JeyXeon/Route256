package postgres

import (
	"context"
	"route256/loms/internal/converter"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type reservationsRepository struct {
	queryEngineProvider QueryEngineProvider
}

const (
	reservationTable = "reservation"

	reservationOrderIdColumn     = "order_id"
	reservationSkuColumn         = "sku"
	reservationWarehouseIdColumn = "warehouse_id"
	reservationCountColumn       = "count"
)

func NewReservationsRepository(queryEngineProvider QueryEngineProvider) *reservationsRepository {
	return &reservationsRepository{
		queryEngineProvider: queryEngineProvider,
	}
}

func (r *reservationsRepository) GetReservations(ctx context.Context, orderId int64) ([]*model.Reservation, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Select(reservationOrderIdColumn, reservationSkuColumn, reservationWarehouseIdColumn, reservationCountColumn).
		From(reservationTable).
		Where(sq.Eq{reservationOrderIdColumn: orderId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var reservations []*schema.Reservation
	if err := pgxscan.Select(ctx, db, &reservations, query, args...); err != nil {
		return nil, err
	}

	result := converter.SchemaToReservationListModel(reservations)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *reservationsRepository) AddReservations(ctx context.Context, orderItems []*model.Reservation) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	queryBuilder := queryBuilder().
		Insert(reservationTable).
		Columns(reservationSkuColumn, reservationOrderIdColumn, reservationWarehouseIdColumn, reservationCountColumn)

	for _, orderItem := range orderItems {
		queryBuilder = queryBuilder.Values(
			orderItem.Sku,
			orderItem.OrderId,
			orderItem.WareHouseId,
			orderItem.Count,
		)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *reservationsRepository) RemoveReservations(ctx context.Context, orderId int64) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Delete(reservationTable).
		Where(sq.Eq{reservationOrderIdColumn: orderId}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *reservationsRepository) RemoveReservationsByOrderIds(ctx context.Context, orderIds []int64) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Delete(reservationTable).
		Where(sq.Eq{reservationOrderIdColumn: orderIds}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
