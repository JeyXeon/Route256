package service

import (
	"context"
	"log"
	"route256/libs/workerpool"
	"route256/loms/internal/model"
	"sync"
	"time"

	"github.com/pkg/errors"
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
				log.Println(result.Error.Error())
			} else if result.Result > 0 {
				log.Printf("cancelled %d timeouted orders", result.Result)
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
	var cancelledOrderIds []int64
	err := s.transactionManager.RunRepeatableRead(ctx, func(ctxTX context.Context) error {
		timeoutedOrderIds, err := s.orderRepository.GetTimeoutedPaymentOrderIds(ctxTX, t)
		if err != nil {
			return err
		}

		if timeoutedOrderIds == nil || len(timeoutedOrderIds) != 0 {
			if err := s.reservationsRepository.RemoveReservationsByOrderIds(ctxTX, timeoutedOrderIds); err != nil {
				return err
			}

			cancelledIds, err := s.orderRepository.UpdateOrdersStatuses(ctxTX, timeoutedOrderIds, model.Cancelled)
			if err != nil {
				return err
			}

			cancelledOrderIds = cancelledIds
		}

		return nil
	})
	if err != nil {
		log.Println(err.Error())
		return 0, ErrCancellingOrdersFailed
	}

	for _, cancelledOrderId := range cancelledOrderIds {
		err := s.orderStateChangeProducer.SendOrderStatusChange(cancelledOrderId, model.Cancelled)
		if err != nil {
			return 0, err
		}
	}

	return len(cancelledOrderIds), nil
}
