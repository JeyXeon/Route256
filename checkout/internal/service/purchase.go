package service

import (
	"context"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientOrder = errors.New("insufficient order")
)

func (m *Service) Purchase(ctx context.Context, user int64) (int64, error) {
	cartItems, err := m.itemsRepository.GetItems(ctx, user)
	if err != nil {
		return 0, err
	}

	orderId, err := m.lomsClient.CreateOrder(ctx, user, cartItems)
	if err != nil {
		return 0, errors.WithMessage(err, "creating order")
	}
	if orderId == 0 {
		return 0, ErrInsufficientOrder
	}

	return orderId, nil
}
