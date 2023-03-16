package service

import (
	"context"
	"route256/checkout/internal/model"
	"route256/libs/workerpool"
	"sync"

	"github.com/pkg/errors"
)

func (m *Service) ListCart(ctx context.Context, user int64) (*model.Cart, error) {
	cartItems, err := m.itemsRepository.GetItems(ctx, user)
	if err != nil {
		return nil, err
	}

	request := func(ctx context.Context, sku uint32) (*model.Product, error) {
		if err := m.productServiceLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		return m.productServiceClient.GetProduct(ctx, sku)
	}

	tasks := make([]workerpool.Task[uint32, *model.Product], 0, len(cartItems))
	for _, cartItem := range cartItems {
		sku := cartItem.SKU

		tasks = append(tasks, workerpool.Task[uint32, *model.Product]{
			Callback: request,
			InArgs:   sku,
		})
	}

	workerPool, responses := workerpool.NewPool[uint32, *model.Product](ctx, 5)
	workerPool.SubmitTasks(ctx, tasks)

	items := make([]*model.Product, 0, len(cartItems))
	totalPrice := uint32(0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for response := range responses {
			if response.Error != nil {
				err = response.Error
				break
			}

			product := response.Result
			items = append(items, product)
			totalPrice += product.Price * product.Count
		}
	}()
	if err != nil {
		return nil, errors.WithMessage(err, "getting product")
	}

	workerPool.Close()
	wg.Wait()

	return &model.Cart{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
