package service

import (
	"context"
	"route256/checkout/internal/model"
)

type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) ([]model.Stock, error)
	CreateOrder(ctx context.Context, user int64, items []model.OrderItem) (int64, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, sku uint32) (model.Product, error)
}

type Service struct {
	lomsClient           LomsClient
	productServiceClient ProductServiceClient
}

func New(stocksChecker LomsClient, productServiceClient ProductServiceClient) *Service {
	return &Service{
		lomsClient:           stocksChecker,
		productServiceClient: productServiceClient,
	}
}
