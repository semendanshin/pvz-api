package inmemmory

import (
	"slices"
	"sync"
	"time"
)

type CacheItem[K comparable, V any] struct {
	Key   K
	Value V

	createTime time.Time

	lastAccessTime time.Time
	usages         int
}

type Cache[K comparable, V any] struct {
	items                map[K]*CacheItem[K, V]
	m                    sync.Mutex
	invalidationStrategy InvalidationStrategy[K, V]
	ttl                  time.Duration
	maxItems             int
}

func NewCache[K comparable, V any](ttl time.Duration, maxItems int, invalidationStrategy InvalidationStrategy[K, V]) *Cache[K, V] {
	return &Cache[K, V]{
		items:                make(map[K]*CacheItem[K, V]),
		invalidationStrategy: invalidationStrategy,
		ttl:                  ttl,
		maxItems:             maxItems,
	}
}

type InvalidationStrategy[K comparable, V any] interface {
	Invalidate(c *Cache[K, V], n int)
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m.Lock()
	defer c.m.Unlock()
	c.items[key] = NewCacheItem(key, value)
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	if len(c.items) > c.maxItems {
		c.invalidate()
	}

	var zeroValue V

	c.m.Lock()
	defer c.m.Unlock()

	item, ok := c.items[key]
	if !ok {
		return zeroValue, false
	}

	if time.Since(item.createTime) > c.ttl {
		delete(c.items, key)
		return zeroValue, false
	}

	item.usages++
	item.lastAccessTime = time.Now()

	return item.Value, true
}

func (c *Cache[K, V]) Delete(key K) {
	c.m.Lock()
	defer c.m.Unlock()
	delete(c.items, key)
}

func (c *Cache[K, V]) invalidate() {
	c.m.Lock()
	defer c.m.Unlock()
	c.invalidationStrategy.Invalidate(c, len(c.items)-c.maxItems)
}

func NewCacheItem[K comparable, V any](key K, value V) *CacheItem[K, V] {
	return &CacheItem[K, V]{
		Key:        key,
		Value:      value,
		createTime: time.Now(),
		usages:     0,
	}
}

type LRUInvalidationStrategy[K comparable, V any] struct {
}

func (s *LRUInvalidationStrategy[K, V]) Invalidate(c *Cache[K, V], n int) {
	items := make([]*CacheItem[K, V], 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	slices.SortFunc(items, func(a, b *CacheItem[K, V]) int {
		return int(a.lastAccessTime.Sub(b.lastAccessTime).Nanoseconds())
	})

	for i := 0; i < n; i++ {
		delete(c.items, items[i].Key)
	}
}

func NewLRUInvalidationStrategy[K comparable, V any]() *LRUInvalidationStrategy[K, V] {
	return &LRUInvalidationStrategy[K, V]{}
}

type LFUInvalidationStrategy[K comparable, V any] struct {
}

func (s *LFUInvalidationStrategy[K, V]) Invalidate(c *Cache[K, V], n int) {
	items := make([]*CacheItem[K, V], 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	slices.SortFunc(items, func(a, b *CacheItem[K, V]) int {
		return a.usages - b.usages
	})

	for i := 0; i < n; i++ {
		delete(c.items, items[i].Key)
	}
}

func NewLFUInvalidationStrategy[K comparable, V any]() *LFUInvalidationStrategy[K, V] {
	return &LFUInvalidationStrategy[K, V]{}
}
