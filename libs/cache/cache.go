package cache

import (
	"context"
	"time"
)

type Cache[T any] struct {
	bucketsCount uint32
	buckets      []bucket[T]
}

func New[T any](ctx context.Context, bucketsCount uint32, ttl time.Duration) *Cache[T] {
	var c Cache[T]

	c.buckets = make([]bucket[T], bucketsCount)
	for i := range c.buckets {
		c.buckets[i].init(ttl)
	}

	go c.refreshCron(ctx)

	return &c
}

func (c *Cache[T]) Set(key uint32, value T) {
	idx := key % c.bucketsCount
	c.buckets[idx].set(key, value)
}

func (c *Cache[T]) Get(key uint32) (value *T, exists bool) {
	idx := key % c.bucketsCount
	return c.buckets[idx].get(key)
}

func (c *Cache[T]) refreshCron(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			for _, b := range c.buckets {
				b.refresh()
			}
		case <-ctx.Done():
			return
		}
	}
}