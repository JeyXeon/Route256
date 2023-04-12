package service

import (
	"context"
	"route256/libs/logger"
	"route256/libs/workerpool"
	"route256/loms/internal/model"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrCancellingOrdersFailed = errors.New("cancelling orders failed")
)

func (s *Service) CheckPaymentTimeoutCron(ctx context.Context) {
	// Таска на очищение просроченных заказов гоняется раз в секунду
	ticker := time.NewTicker(time.Second * 1)

	workerPool, results := workerpool.NewPool[time.Time, int](ctx, 5)
	defer workerPool.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		for result := range results {
			if result.Error != nil {
				logger.Error("payment timeout check failed", zap.Error(result.Error))
			} else if result.Result > 0 {
				logger.Info("cancelled timeouted orders", zap.Int("cancelled", result.Result))
			}
		}
	}()

	// По каждому тику сабмитим таску с текущим временем в параметре
	for {
		select {
		case <-ticker.C:
			task := workerpool.Task[time.Time, int]{
				Callback: s.cancelTimeoutedOrders,
				InArgs:   time.Now(),
			}
			workerPool.SubmitTask(ctx, task)
		case <-ctx.Done():
			return
		}
	}
}

// Метод, используемый в качестве колбэка для таски
func (s *Service) cancelTimeoutedOrders(ctx context.Context, t time.Time) (int, error) {
	var cancelledOrdersAmount int
	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		timeoutedOrderIds, err := s.orderRepository.GetTimeoutedPaymentOrderIds(ctxTX, t)
		if err != nil {
			return err
		}

		if len(timeoutedOrderIds) != 0 {
			if err := s.reservationsRepository.RemoveReservationsByOrderIds(ctxTX, timeoutedOrderIds); err != nil {
				return err
			}

			cancelledIds, err := s.orderRepository.UpdateOrdersStatuses(ctxTX, timeoutedOrderIds, model.Cancelled)
			if err != nil {
				return err
			}

			for _, cancelledOrderId := range cancelledIds {
				orderStateChangeRecord, err := model.NewOrderStatusChangeKafkaRecord(cancelledOrderId, model.Cancelled)
				if err != nil {
					return err
				}

				err = s.outboxKafkaRepository.CreateKafkaRecord(ctxTX, orderStateChangeRecord)
				if err != nil {
					return err
				}
			}

			cancelledOrdersAmount = len(cancelledIds)
		}

		return nil
	})
	if err != nil {
		return 0, ErrCancellingOrdersFailed
	}

	return cancelledOrdersAmount, nil
}
