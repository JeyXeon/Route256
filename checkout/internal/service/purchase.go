package service

import (
	"context"
	"route256/checkout/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrInsufficientOrder = errors.New("insufficient order")
)

func (m *Service) Purchase(ctx context.Context, user int64) (int64, error) {
	orderId, err := m.lomsClient.CreateOrder(ctx, user, []*model.OrderItem{{
		SKU:   5,
		Count: 3,
	}})
	if err != nil {
		return 0, errors.WithMessage(err, "creating order")
	}
	if orderId == 0 {
		return 0, ErrInsufficientOrder
	}
	return orderId, nil
}
