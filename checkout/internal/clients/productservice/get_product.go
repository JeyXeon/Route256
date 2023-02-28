package productservice

import (
	"context"
	"route256/checkout/internal/model"
	productserviceapi "route256/checkout/pkg/productservice"
)

func (c *client) GetProduct(ctx context.Context, SKU uint32) (*model.Product, error) {
	req := &productserviceapi.GetProductRequest{
		Token: c.token,
		Sku:   SKU,
	}

	res, err := c.productServiceClient.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}

	return &model.Product{
		SKU:   SKU,
		Count: 1,
		Name:  res.GetName(),
		Price: res.GetPrice(),
	}, nil
}
