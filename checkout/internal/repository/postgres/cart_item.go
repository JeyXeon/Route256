package postgres

import (
	"context"
	"route256/checkout/internal/model"
)

type CartItemRepository struct {
	queryEngine QueryEngine
}

func NewCartItemRepository(queryEngine QueryEngine) *CartItemRepository {
	return &CartItemRepository{
		queryEngine: queryEngine,
	}
}

func (cartItemRepository *CartItemRepository) AddItem(ctx context.Context, userId int64, item *model.OrderItem) error {
	return nil
}

func (cartItemRepository *CartItemRepository) DeleteItem(ctx context.Context, userId int64, item *model.OrderItem) error {
	return nil
}

func (cartItemRepository *CartItemRepository) GetItems(ctx context.Context, userId int64) ([]*model.OrderItem, error) {
	return nil, nil
}
