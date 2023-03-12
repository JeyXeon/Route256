package service

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
)

func (m *Service) AddToCart(ctx context.Context, user int64, sku uint32, count uint32) error {
	stocks, err := m.lomsClient.Stocks(ctx, sku)
	if err != nil {
		return errors.WithMessage(err, "checking stocks")
	}

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			err := m.itemsRepository.AddItem(ctx, user, &model.CartItem{
				SKU:   sku,
				Count: count,
			})
			if err != nil {
				return err
			}

			return nil
		}
	}

	return ErrInsufficientStocks
}
