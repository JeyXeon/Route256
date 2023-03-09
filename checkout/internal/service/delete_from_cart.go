package service

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrDeletingFromCart = errors.New("deleting from cart failed")
)

func (m *Service) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint32) error {
	err := m.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		cartItem, err := m.itemsRepository.GetItem(ctxTX, user, sku)
		if err != nil {
			return err
		}

		if cartItem.Count == count {
			if err := m.itemsRepository.DeleteItem(ctx, user, &model.CartItem{
				SKU:   sku,
				Count: count,
			}); err != nil {
				return err
			}
		} else {
			if err := m.itemsRepository.RemoveItems(ctx, user, &model.CartItem{
				SKU:   sku,
				Count: count,
			}); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return ErrDeletingFromCart
	}

	return nil
}
