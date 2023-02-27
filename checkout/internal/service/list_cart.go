package service

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

func (m *Service) ListCart(ctx context.Context, user int64) (model.Cart, error) {
	skus := []uint32{1076963, 4288068, 24808287}

	items := make([]*model.Product, 0, len(skus))
	totalPrice := uint32(0)
	for _, sku := range skus {
		product, err := m.productServiceClient.GetProduct(ctx, sku)
		if err != nil {
			return model.Cart{Items: []*model.Product{}, TotalPrice: 0}, errors.WithMessage(err, "getting product")
		}
		items = append(items, &product)
		totalPrice += product.Price * product.Count
	}

	return model.Cart{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
