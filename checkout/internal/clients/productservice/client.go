package productservice

import (
	"context"
	"route256/checkout/internal/model"
	productserviceapi "route256/checkout/pkg/productservice"

	"google.golang.org/grpc"
)

type Client interface {
	GetProduct(ctx context.Context, SKU uint32) (*model.Product, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, in *productserviceapi.GetProductRequest, opts ...grpc.CallOption) (*productserviceapi.GetProductResponse, error)
}

type client struct {
	productServiceClient ProductServiceClient
	token                string
}

func New(productServiceClient ProductServiceClient, token string) Client {
	return &client{
		productServiceClient: productServiceClient,
		token:                token,
	}
}
