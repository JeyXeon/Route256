package service

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

func (m *Service) ListCart(ctx context.Context, user int64) (*model.Cart, error) {
	cartItems, err := m.itemsRepository.GetItems(ctx, user)
	if err != nil {
		return nil, err
	}

	items := make([]*model.Product, 0, len(cartItems))
	totalPrice := uint32(0)
	for _, cartItem := range cartItems {
		sku := cartItem.SKU
		product, err := m.productServiceClient.GetProduct(ctx, sku)
		if err != nil {
			return nil, errors.WithMessage(err, "getting product")
		}

		items = append(items, product)
		totalPrice += product.Price * product.Count
	}

	return &model.Cart{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
