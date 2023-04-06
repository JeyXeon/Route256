package service

import (
	"context"
	"route256/loms/internal/model"

	"github.com/pkg/errors"
)

var (
	ErrCreatingOrderFailed    = errors.New("creating order failed")
	ErrNotEnoughStocksForItem = errors.New("not enough stocks for item")
)

func (s *Service) CreateOrder(ctx context.Context, userId int64, orderItems []*model.OrderItem) (int64, error) {
	var orderId int64

	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		createdOrderId, err := s.orderRepository.CreateOrder(ctxTX, userId)
		if err != nil {
			return err
		}

		orderId = createdOrderId

		skus := orderItemsToSkus(orderItems)
		stocks, err := s.stocksRepository.GetStocks(ctxTX, skus)
		if err != nil {
			return err
		}

		reservedStocks, err := createReservedStocks(orderItems, stocks)
		if err != nil {
			if err := s.orderRepository.UpdateOrderStatus(ctxTX, orderId, model.Failed); err != nil {
				return err
			}

			orderStateChangeRecord, err := model.NewOrderStatusChangeKafkaRecord(orderId, model.Failed)
			if err != nil {
				return err
			}

			err = s.outboxKafkaRepository.CreateKafkaRecord(ctxTX, orderStateChangeRecord)
			if err != nil {
				return err
			}

			return nil
		}

		reservations := stocksToReservations(reservedStocks, createdOrderId)
		if err := s.reservationsRepository.AddReservations(ctxTX, reservations); err != nil {
			return err
		}

		if err := s.orderRepository.UpdateOrderStatus(ctxTX, orderId, model.AwaitingPayment); err != nil {
			return err
		}

		orderStateChangeRecord, err := model.NewOrderStatusChangeKafkaRecord(orderId, model.AwaitingPayment)
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
		return 0, ErrCreatingOrderFailed
	}

	return orderId, nil
}

func createReservedStocks(orderItems []*model.OrderItem, stocks []*model.Stock) ([]*model.Stock, error) {
	stocksBySku := make(map[uint32][]*model.Stock)
	for _, stock := range stocks {
		_, exists := stocksBySku[stock.Sku]
		if !exists {
			stocksBySku[stock.Sku] = make([]*model.Stock, 0)
		}
		stocksBySku[stock.Sku] = append(stocksBySku[stock.Sku], stock)
	}

	result := make([]*model.Stock, 0)
	for _, orderItem := range orderItems {
		count := orderItem.Count
		stocks := stocksBySku[orderItem.Sku]
		for _, stock := range stocks {
			reserveCount := min(stock.Count, uint64(count))
			result = append(result, &model.Stock{
				Sku:         stock.Sku,
				WareHouseId: stock.WareHouseId,
				Count:       reserveCount,
			})

			count -= uint32(reserveCount)
			if count == 0 {
				break
			}
		}

		if count > 0 {
			return nil, ErrNotEnoughStocksForItem
		}
	}

	return result, nil
}

func orderItemsToSkus(orderItems []*model.OrderItem) []uint32 {
	if orderItems == nil {
		return nil
	}

	skus := make(map[uint32]struct{})
	for _, orderItem := range orderItems {
		_, exists := skus[orderItem.Sku]
		if !exists {
			skus[orderItem.Sku] = struct{}{}
		}
	}

	result := make([]uint32, 0, len(skus))
	for sku := range skus {
		result = append(result, sku)
	}

	return result
}

func stocksToReservations(stocks []*model.Stock, orderId int64) []*model.Reservation {
	if stocks == nil {
		return nil
	}

	reservations := make([]*model.Reservation, 0, len(stocks))
	for _, stock := range stocks {
		reservations = append(reservations, &model.Reservation{
			OrderId:     orderId,
			Sku:         stock.Sku,
			WareHouseId: stock.WareHouseId,
			Count:       uint32(stock.Count),
		})
	}

	return reservations
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
