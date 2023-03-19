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

	tasks := m.prepareGetProductTasks(cartItems)

	workerPool, responses := workerpool.NewPool[uint32, *model.Product](ctx, 5)
	workerPool.SubmitTasks(ctx, tasks)

	items := make([]*model.Product, 0, len(cartItems))
	totalPrice := uint32(0)

	// Вг чтобы дождаться, пока прочтутся все полученные результаты из канала, чтобы не вернуть ответ раньше времени
	wg := sync.WaitGroup{}
	wg.Add(1)

	var reqErr error
	// В отдельной горутине читаем из канала, чтобы не заблочить текущую
	go func() {
		defer wg.Done()
		for response := range responses {
			if response.Error != nil {
				reqErr = response.Error
				break
			}

			product := response.Result
			items = append(items, product)
			totalPrice += product.Price * product.Count
		}
	}()
	// Закрываем воркер пул, чтобы выйти из := range responses выше
	workerPool.Close()
	// Дожидаемся вг ответов
	wg.Wait()

	// Если была зафейлившаяся таска, вышли из := range в через break и возвращаемся с ошибкой
	if reqErr != nil {
		return nil, errors.WithMessage(reqErr, "getting product")
	}

	return &model.Cart{
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}

func (m *Service) prepareGetProductTasks(cartItems []*model.CartItem) []workerpool.Task[uint32, *model.Product] {
	// Подготавливаем колбэк для генерации тасок
	request := func(ctx context.Context, sku uint32) (*model.Product, error) {
		if err := m.productServiceLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		return m.productServiceClient.GetProduct(ctx, sku)
	}

	// Генерируем по таске для каждого sku в корзине
	tasks := make([]workerpool.Task[uint32, *model.Product], 0, len(cartItems))
	for _, cartItem := range cartItems {
		tasks = append(tasks, workerpool.Task[uint32, *model.Product]{
			Callback: request,
			InArgs:   cartItem.SKU,
		})
	}

	return tasks
}
