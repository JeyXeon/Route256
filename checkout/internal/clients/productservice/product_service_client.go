package productservice

import (
	"context"
	"google.golang.org/grpc"
	productserviceapi "route256/checkout/pkg/productservice"
)

type ProductServiceClient interface {
	GetProduct(ctx context.Context, in *productserviceapi.GetProductRequest, opts ...grpc.CallOption) (*productserviceapi.GetProductResponse, error)
}

type client struct {
	productServiceClient ProductServiceClient
	token                string
}

func New(productServiceClient ProductServiceClient, token string) *client {
	return &client{
		productServiceClient: productServiceClient,
		token:                token,
	}
}
