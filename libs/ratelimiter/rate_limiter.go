package ratelimiter

import (
	"context"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	ErrRequestsAmountLimitReached = errors.New("requests amount limit reached")
)

type Limiter interface {
	Wait(ctx context.Context) error
}

type limiter struct {
	secondsLimit int
	tokensLimit  int

	tokens int
	mu     *sync.Mutex
	ticker *time.Ticker
}

func NewLimiter(ctx context.Context, secondsLimit int, tokensLimit int) Limiter {
	l := &limiter{
		secondsLimit: secondsLimit,
		tokensLimit:  tokensLimit,
		mu:           &sync.Mutex{},
	}

	go l.processing(ctx)

	return l
}

func (l *limiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Лочим мбютекс для чтения кол-ва доступных токенов
	l.mu.Lock()

	// Если токенов не осталось, заворачиваем с ошибкой
	if l.tokens <= 0 {
		l.mu.Unlock()
		return ErrRequestsAmountLimitReached
	}
	// Если остались, уменьшаем кол-во доступных токенов
	l.tokens -= 1

	l.mu.Unlock()

	return nil
}

func (l *limiter) processing(ctx context.Context) {
	//Создаем тикер, который будет раз в secondsLimit секунд рефрешить количество доступных токенов
	l.ticker = time.NewTicker(time.Second * time.Duration(l.secondsLimit))

	for {
		select {
		// По каждому тику лочим мьютекс, рефрешим кол-во токенов и разлочиваем
		case <-l.ticker.C:
			l.mu.Lock()
			l.tokens = l.tokensLimit
			l.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}

}
