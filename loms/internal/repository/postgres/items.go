package postgres

import (
	"context"
	"route256/loms/internal/model"
)

type itemsRepository struct {
	queryEngine QueryEngine
}

func NewItemsRepository(queryEngine QueryEngine) *itemsRepository {
	return &itemsRepository{
		queryEngine: queryEngine,
	}
}

func (itemsRepository *itemsRepository) GetOrderItems(ctx context.Context, orderId int64) ([]*model.OrderItem, error) {
	return nil, nil
}

func (itemsRepository *itemsRepository) GetReservations(ctx context.Context, orderId int64) ([]*model.Reservation, error) {
	return nil, nil
}

func (itemsRepository *itemsRepository) AddReservations(ctx context.Context, orderItems []*model.Reservation) error {
	return nil
}

func (itemsRepository *itemsRepository) RemoveReservations(ctx context.Context, orderId int64) error {
	return nil
}

func (itemsRepository *itemsRepository) GetStocks(ctx context.Context, sku uint32) ([]*model.Stock, error) {
	return nil, nil
}

func (itemsRepository *itemsRepository) AddStocks(ctx context.Context, stocks []*model.Stock) error {
	return nil
}

func (itemsRepository *itemsRepository) ReserveStocks(ctx context.Context, stocks []*model.OrderItem) ([]*model.Stock, error) {
	return nil, nil
}
