package domain

import (
	"context"

	"github.com/pkg/errors"
)

type Product struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Cart struct {
	Items      []Product
	TotalPrice uint32
}

func (m *Model) ListCart(ctx context.Context, user int64) (Cart, error) {
	skus := []uint32{1076963, 4288068, 24808287}

	items := make([]Product, 0, len(skus))
	totalPrice := uint32(0)
	for _, sku := range skus {
		product, err := m.productServiceClient.GetProduct(ctx, sku)
		if err != nil {
			return Cart{Items: []Product{}, TotalPrice: 0}, errors.WithMessage(err, "getting product")
		}
		items = append(items, product)
		totalPrice += product.Price * uint32(product.Count)
	}

	return Cart{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
