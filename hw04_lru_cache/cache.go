package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCacheEntry struct {
	key   Key
	value any
}

type lruCache struct {
	sync.Mutex

	capacity int
	queue    List

	items map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(k Key, v any) (ok bool) {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.items[k]; ok {
		c.queue.MoveToFront(item)
		item.Value = lruCacheEntry{k, v}
		return true
	}

	if c.queue.Len() == c.capacity && c.capacity > 0 {
		item := c.queue.Back()
		c.queue.Remove(item)
		delete(c.items, item.Value.(lruCacheEntry).key)
	}

	if c.capacity > 0 {
		c.items[k] = c.queue.PushFront(lruCacheEntry{k, v})
	}

	return false
}

func (c *lruCache) Get(k Key) (any, bool) {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.items[k]; ok {
		c.queue.MoveToFront(item)
		return item.Value.(lruCacheEntry).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
