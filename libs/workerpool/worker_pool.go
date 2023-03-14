package workerpool

import (
	"context"
	"sync"
)

type TaskResult[Out any] struct {
	Result Out
	Error  error
}

// Task содержит колбэк, возвращающий результат и ошибку, и аргумент для этого колбека
// Не супер удобная реализация, тк работает только с аргументами одного типа, зато дженерно
type Task[In, Out any] struct {
	Callback func(context.Context, In) (Out, error)
	InArgs   In
}

type Pool[In, Out any] interface {
	SubmitTasks(context.Context, []Task[In, Out])
	SubmitTask(ctx context.Context, task Task[In, Out])
	Close()
}

type p[In, Out any] struct {
	amountWorkers int

	tasksWg   sync.WaitGroup
	workersWg sync.WaitGroup

	taskSource chan Task[In, Out]
	outSink    chan TaskResult[Out]
}

func NewPool[In, Out any](ctx context.Context, amountWorkers int) (Pool[In, Out], <-chan TaskResult[Out]) {
	pool := &p[In, Out]{
		amountWorkers: amountWorkers,
	}

	pool.bootstrap(ctx)

	return pool, pool.outSink
}

func (pool *p[In, Out]) Close() {

	// Дожидаемся, пока горутина в методе submit доразгребет все переданные таски
	// (Чтобы не закрыть канал pool.taskSource раньше времени)
	pool.tasksWg.Wait()

	// Больше задач не будет
	close(pool.taskSource)

	// Дожидаемся, пока все воркеры закончат работы
	pool.workersWg.Wait()

	// Закрываем канал на выход, чтобы потребители могли выйти из := range
	close(pool.outSink)
}

func (pool *p[In, Out]) SubmitTask(ctx context.Context, task Task[In, Out]) {
	pool.SubmitTasks(ctx, []Task[In, Out]{task})
}

func (pool *p[In, Out]) SubmitTasks(ctx context.Context, tasks []Task[In, Out]) {
	pool.tasksWg.Add(1)

	// Льём все переданные таски в канал
	go func() {
		defer pool.tasksWg.Done()
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return

			case pool.taskSource <- task:
			}
		}
	}()
}

func (pool *p[In, Out]) bootstrap(ctx context.Context) {
	// Инициализируем каналы для поступающих тасок и полученных результатов их выполнения
	pool.taskSource = make(chan Task[In, Out], pool.amountWorkers)
	pool.outSink = make(chan TaskResult[Out], pool.amountWorkers)

	// Запускаем заданное при создании количество воркеров
	for i := 0; i < pool.amountWorkers; i++ {
		pool.workersWg.Add(1)
		go func() {
			defer pool.workersWg.Done()
			worker(ctx, pool.taskSource, pool.outSink)
		}()
	}
}

func worker[In, Out any](
	ctx context.Context,
	taskSource <-chan Task[In, Out],
	resultSink chan<- TaskResult[Out],
) {
	// Читаем таски из канала taskSource, выполняем их (processTask) и заливаем результат в pool.outSink
	for task := range taskSource {
		select {
		case <-ctx.Done():
			return
		case resultSink <- processTask[In, Out](ctx, task):
		}
	}

	return
}

func processTask[In, Out any](ctx context.Context, task Task[In, Out]) TaskResult[Out] {
	result, err := task.Callback(ctx, task.InArgs)
	return TaskResult[Out]{
		Result: result,
		Error:  err,
	}
}
