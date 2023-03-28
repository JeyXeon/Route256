package loms

import (
	"context"
	"route256/checkout/internal/model"
	lomsapi "route256/checkout/pkg/loms"

	"google.golang.org/grpc"
)

type Client interface {
	CreateOrder(ctx context.Context, user int64, items []*model.CartItem) (int64, error)
	Stocks(ctx context.Context, sku uint32) ([]*model.Stock, error)
}

type LomsClient interface {
	Stocks(ctx context.Context, in *lomsapi.StocksRequest, opts ...grpc.CallOption) (*lomsapi.StocksResponse, error)
	CreateOrder(ctx context.Context, in *lomsapi.CreateOrderRequest, opts ...grpc.CallOption) (*lomsapi.CreateOrderResponse, error)
}

type client struct {
	lomsClient LomsClient
}

func New(lomsClient LomsClient) Client {
	return &client{
		lomsClient: lomsClient,
	}
}
