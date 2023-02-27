package domain

import "context"

type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (int64, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, sku uint32) (Product, error)
}

type Model struct {
	lomsClient           LomsClient
	productServiceClient ProductServiceClient
}

func New(stocksChecker LomsClient, productServiceClient ProductServiceClient) *Model {
	return &Model{
		lomsClient:           stocksChecker,
		productServiceClient: productServiceClient,
	}
}
