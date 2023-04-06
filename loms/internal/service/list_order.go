package service

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrListingOrderFailed = errors.New("listing order failed")
)

func (s *Service) ListOrder(ctx context.Context, orderId int64) (*model.Order, error) {
	var order *model.Order

	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		result, err := s.orderRepository.GetOrder(ctxTX, orderId)
		if err != nil {
			return err
		}

		reservations, err := s.reservationsRepository.GetReservations(ctx, orderId)
		if err != nil {
			return err
		}

		result.Items = reservationsToOrderItems(reservations)
		order = result
		return nil
	})
	if err != nil {
		return nil, ErrListingOrderFailed
	}

	return order, nil
}

func reservationsToOrderItems(reservations []*model.Reservation) []*model.OrderItem {
	if reservations == nil {
		return nil
	}

	countBySku := make(map[uint32]uint32)
	for _, reservation := range reservations {
		count, exists := countBySku[reservation.Sku]
		if !exists {
			countBySku[reservation.Sku] = reservation.Count
		} else {
			countBySku[reservation.Sku] = reservation.Count + count
		}
	}

	result := make([]*model.OrderItem, 0, len(countBySku))
	for sku, count := range countBySku {
		result = append(result, &model.OrderItem{
			Sku:   sku,
			Count: count,
		})
	}

	return result
}
