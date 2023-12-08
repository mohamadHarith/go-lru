package lru

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

const MAX_CAPACITY = 100

var (
	ErrAssertion            = errors.New("assertion error")
	ErrCacheExpired         = errors.New("cache expired")
	ErrCacheNotFound        = errors.New("cache not found")
	ErrInvalidCacheCapacity = errors.New("invalid cache capacity")
	ErrInvalidExpiry        = errors.New("invalid expiry")
)

type LRUCache struct {
	capacity uint32
	expiry   time.Duration
	queue    *list.List
	items    map[string]*list.Element
	locker   *sync.Mutex
}

type node struct {
	key      string
	value    any
	expireAt time.Time
}

func New(capacity uint32, expiry time.Duration) (*LRUCache, error) {
	if capacity > MAX_CAPACITY || capacity < 2 {
		return nil, ErrInvalidCacheCapacity
	}
	if expiry <= 0 {
		return nil, ErrInvalidExpiry
	}
	return &LRUCache{
		capacity: capacity,
		expiry:   expiry,
		queue:    list.New(),
		items:    make(map[string]*list.Element),
		locker:   &sync.Mutex{},
	}, nil
}

func (c *LRUCache) Put(key string, value any) {
	c.locker.Lock()
	defer c.locker.Unlock()

	nd := &node{key: key, value: value, expireAt: time.Now().Add(c.expiry)}

	e, ok := c.items[key]
	if ok {
		e.Value = nd
		c.queue.MoveToFront(e)
		c.items[key] = e
	} else {
		if len(c.items) >= int(c.capacity) {
			c.evict(c.queue.Back())
		}
		c.items[key] = c.queue.PushFront(nd)
	}
}

func (c *LRUCache) Get(key string) (any, error) {
	c.locker.Lock()
	defer c.locker.Unlock()

	v, ok := c.items[key]
	if !ok {
		return nil, ErrCacheNotFound
	}

	d, ok := v.Value.(*node)
	if !ok {
		return nil, ErrAssertion
	}

	if time.Now().After(d.expireAt) {
		c.evict(v)
		return nil, ErrCacheExpired
	}

	c.queue.MoveToFront(v)

	return d.value, nil
}

func (c *LRUCache) evict(e *list.Element) {
	n, ok := e.Value.(*node)
	if ok {
		delete(c.items, n.key)
	}
	c.queue.Remove(e)
}
