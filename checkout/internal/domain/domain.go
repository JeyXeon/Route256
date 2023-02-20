package domain

import "context"

type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (int64, error)
}

type Model struct {
	lomsClient LomsClient
}

func New(stocksChecker LomsClient) *Model {
	return &Model{
		lomsClient: stocksChecker,
	}
}
