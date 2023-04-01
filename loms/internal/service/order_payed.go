package service

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrPayingOrder = errors.New("paying order failed")
)

func (s *Service) PayOrder(ctx context.Context, orderId int64) error {
	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		if err := s.orderRepository.UpdateOrderStatus(ctxTX, orderId, model.Payed); err != nil {
			return err
		}

		reservations, err := s.reservationsRepository.GetReservations(ctxTX, orderId)
		if err != nil {
			return err
		}

		stocksToWriteOff := reservationsToStocks(reservations)
		if err = s.stocksRepository.WriteOffStocks(ctxTX, stocksToWriteOff); err != nil {
			return err
		}

		orderStateChangeRecord, err := model.NewOrderStatusChangeKafkaRecord(orderId, model.Payed)
		if err != nil {
			return err
		}

		err = s.outboxKafkaRepository.CreateKafkaRecord(ctxTX, orderStateChangeRecord)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return ErrPayingOrder
	}

	return nil
}

func reservationsToStocks(reservations []*model.Reservation) []*model.Stock {
	if reservations == nil {
		return nil
	}

	stocks := make([]*model.Stock, 0, len(reservations))
	for _, reservation := range reservations {
		stocks = append(stocks, &model.Stock{
			Sku:         reservation.Sku,
			WareHouseId: reservation.WareHouseId,
			Count:       uint64(reservation.Count),
		})
	}

	return stocks
}
