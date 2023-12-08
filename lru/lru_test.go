package lru

import (
	"slices"
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache, err := New(4, 2*time.Second)
	if err != nil {
		t.Errorf("Error creating cache: %s", err)
	}

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Put("key3", "value3")
	cache.Put("key4", "value4")
	cache.Put("key5", "value5")

	// check if key 1 is evicted due to capacity
	_, err = cache.Get("key1")
	if err == nil || err != ErrCacheNotFound {
		t.Errorf("cache should have been evicted due to capacity but got err %v", err)
	}

	if len(cache.items) != 4 || cache.queue.Len() != 4 {
		t.Errorf("expecting cache len 4 but got %v", cache.queue.Len())
	}

	time.Sleep(time.Second * 3)

	_, err = cache.Get("key3")
	if err == nil || err != ErrCacheExpired {
		t.Errorf("cache should have been evicted due to expiry but got err %v", err)
	}

	cache.Put("key6", "value6")
	cache.Put("key7", []int{1, 2, 3})

	val, err := cache.Get("key6")
	if err != nil {
		t.Error(err)
	}
	if val != "value6" {
		t.Errorf("expecting value6 but got %v", val)
	}

	val, err = cache.Get("key7")
	if err != nil {
		t.Error(err)
	}

	data, ok := val.([]int)
	if !ok {

	}

	if !slices.Equal[[]int](data, []int{1, 2, 3}) {
		t.Errorf("expecting value6 but got %v", val)
	}
}
