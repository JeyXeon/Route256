package cache

import (
	"sync"
	"time"
)

type chunk[T any] struct {
	key       uint32
	value     *T
	createdAt time.Time
}

type bucket[T any] struct {
	mu     sync.RWMutex
	chunks []*chunk[T]
	ttl    time.Duration
}

func (b *bucket[T]) init(ttl time.Duration) {
	b.chunks = make([]*chunk[T], 0)
	b.mu = sync.RWMutex{}
	b.ttl = ttl
}

func (b *bucket[T]) get(key uint32) (value *T, exists bool) {
	b.mu.RLock()

	for _, v := range b.chunks {
		if v.key == key && v.createdAt.Add(b.ttl).After(time.Now()) {
			b.mu.RUnlock()

			return v.value, true
		}
	}
	b.mu.RUnlock()

	return nil, false
}

func (b *bucket[T]) set(key uint32, value T) {
	b.mu.Lock()

	for _, v := range b.chunks {
		if key == v.key {
			v.value = &value
			v.createdAt = time.Now()

			b.mu.Unlock()
			return
		}
	}

	b.chunks = append(b.chunks, &chunk[T]{
		key:       key,
		value:     &value,
		createdAt: time.Now(),
	})

	b.mu.Unlock()
}

func (b *bucket[T]) refresh() {
	b.mu.Lock()
	refreshedChunks := make([]*chunk[T], 0, len(b.chunks))
	for _, v := range b.chunks {
		if !v.createdAt.Add(b.ttl).After(time.Now()) {
			refreshedChunks = append(refreshedChunks, v)
		}
	}

	b.chunks = refreshedChunks
	b.mu.Unlock()
}
