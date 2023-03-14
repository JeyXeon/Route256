package postgres

import (
	"context"
	"fmt"
	"route256/loms/internal/converter"
	"route256/loms/internal/model"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type orderRepository struct {
	queryEngineProvider QueryEngineProvider
}

const (
	orderTable = "user_order"

	orderIdColumn = "order_id"
	userIdColumn  = "user_id"
	statusColumn  = "status"
)

func NewOrderRepository(queryEngineProvider QueryEngineProvider) *orderRepository {
	return &orderRepository{
		queryEngineProvider: queryEngineProvider,
	}
}

func (r *orderRepository) CreateOrder(ctx context.Context, userId int64) (int64, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Insert(orderTable).
		Columns(userIdColumn, statusColumn).
		Values(userId, schema.New).
		Suffix(fmt.Sprintf("RETURNING %s", orderIdColumn)).
		ToSql()
	if err != nil {
		return 0, err
	}

	row := db.QueryRow(ctx, query, args...)

	var orderId int64
	if err := row.Scan(&orderId); err != nil {
		return 0, err
	}

	return orderId, nil
}

func (r *orderRepository) GetOrder(ctx context.Context, orderId int64) (*model.Order, error) {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	query, args, err := queryBuilder().
		Select(orderIdColumn, userIdColumn, statusColumn).
		From(orderTable).
		Where(sq.Eq{orderIdColumn: orderId}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var order schema.Order
	if err := pgxscan.Get(ctx, db, &order, query, args...); err != nil {
		return nil, err
	}

	result, err := converter.SchemaToOrderModel(&order)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderId int64, newStatus model.OrderStatus) error {
	db := r.queryEngineProvider.GetQueryEngine(ctx)

	status, err := converter.ModelToOrderStatusSchema(newStatus)
	if err != nil {
		return err
	}

	query, args, err := queryBuilder().
		Update(orderTable).
		Set(statusColumn, status).
		Where(sq.Eq{orderIdColumn: orderId}).
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
