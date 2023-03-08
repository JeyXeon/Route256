package postgres

import (
	"context"
	"route256/loms/internal/model"
)

type orderRepository struct {
	queryEngine QueryEngine
}

func NewOrderRepository(queryEngine QueryEngine) *orderRepository {
	return &orderRepository{
		queryEngine: queryEngine,
	}
}

func (orderRepository *orderRepository) CreateOrder(ctx context.Context, userId int64) (int64, error) {
	return 0, nil
}

func (orderRepository *orderRepository) GetOrder(ctx context.Context, orderId int64) (*model.Order, error) {
	return nil, nil
}

func (orderRepository *orderRepository) UpdateOrderStatus(ctx context.Context, orderId int64, newStatus model.OrderStatus) error {
	return nil
}
