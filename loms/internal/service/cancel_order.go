package service

import (
	"context"
	"log"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrCancellingOrderFailed = errors.New("cancelling order failed")
)

func (s *Service) CancelOrder(ctx context.Context, orderId int64) error {
	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		reservations, err := s.itemRepository.GetReservations(ctxTX, orderId)
		if err != nil {
			return err
		}

		if err := s.itemRepository.RemoveReservations(ctxTX, orderId); err != nil {
			return err
		}

		order, err := s.orderRepository.GetOrder(ctxTX, orderId)
		if err != nil {
			return err
		}

		if order.Status == model.Payed {
			err = s.stocksRepository.RevertReservations(ctxTX, reservations)
			if err != nil {
				return err
			}
		}

		err = s.orderRepository.UpdateOrderStatus(ctxTX, orderId, model.Cancelled)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err.Error())
		return ErrCancellingOrderFailed
	}

	return nil
}
