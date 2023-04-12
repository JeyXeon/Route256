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

	c.bucketsCount = bucketsCount
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

func (c *Cache[T]) Get(key uint32) (*T, bool) {
	CacheRequestsTotal.Inc()

	timeStart := time.Now()

	idx := key % c.bucketsCount
	value, exists := c.buckets[idx].get(key)

	elapsed := time.Since(timeStart)
	CacheHistogramResponseTime.WithLabelValues("cached").Observe(elapsed.Seconds())

	if exists {
		CacheHitCount.Inc()
	}

	return value, exists
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
