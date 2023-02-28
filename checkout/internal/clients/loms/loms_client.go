package loms

import (
	"context"
	lomsapi "route256/checkout/pkg/loms"

	"google.golang.org/grpc"
)

type LomsClient interface {
	Stocks(ctx context.Context, in *lomsapi.StocksRequest, opts ...grpc.CallOption) (*lomsapi.StocksResponse, error)
	CreateOrder(ctx context.Context, in *lomsapi.CreateOrderRequest, opts ...grpc.CallOption) (*lomsapi.CreateOrderResponse, error)
}

type client struct {
	lomsClient LomsClient
}

func New(lomsClient LomsClient) *client {
	return &client{
		lomsClient: lomsClient,
	}
}
