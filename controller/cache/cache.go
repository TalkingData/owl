package cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = -1

	DefaultExpiration time.Duration = 0
)

type Item struct {
	Object     interface{}
	Expiration int64
}

func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}

	return time.Now().UnixNano() > item.Expiration
}

type Cache struct {
	*cache
}

type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{}) //callback
	janitor           *janitor
}

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
	c.mu.Unlock()

}

func (c *cache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.items[k] = Item{
		Object:     x,
		Expiration: e,
	}
}

func (c *cache) get(k string) (interface{}, bool) {
	item, found := c.items[k]
	if !found {
		return nil, false
	}
	//check item expired
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return nil, false
	}
	return item.Object, true
}

func (c *cache) Add(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	//check exit
	_, found := c.get(k)
	if found {
		err := fmt.Errorf("Item:%s has already exit", k)
		return err
	}
	//set
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *cache) Replace(k string, x interface{}, d time.Duration) error {
	c.mu.Lock()
	//check exit
	_, found := c.get(k)
	if !found {
		c.mu.Unlock()
		err := fmt.Errorf("Item:%s dosen't exit", k)
		return err
	}
	c.set(k, x, d)
	c.mu.Unlock()
	return nil
}

func (c *cache) Increment(k string, n int64) error {
	c.mu.Lock()
	v, found := c.items[k]
	if !found || v.Expired() {
		c.mu.Unlock()
		return fmt.Errorf("Item not found or expired")
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) + int(n)
	default:
		c.mu.Unlock()
		return fmt.Errorf("not support value tyepe")
	}
	c.items[k] = v
	c.mu.Unlock()
	return nil
}

func (c *cache) Get(k string) (interface{}, bool) {
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		c.mu.Unlock()
		return nil, false
	}
	//check item expired
	if item.Expiration < time.Now().UnixNano() {
		c.mu.Unlock()
		return nil, false
	}
	c.mu.Unlock()
	return item.Object, true

}

func (c *cache) GetItems() map[string]Item {
	return c.items
}

func (c *cache) delete(k string) (interface{}, bool) {
	//
	if c.onEvicted != nil {
		if v, found := c.items[k]; found {
			delete(c.items, k)
			return v.Object, true
		}
	}
	delete(c.items, k)
	return nil, false
}

func (c *cache) Delete(k string) {
	c.mu.Lock()
	v, evicted := c.delete(k)
	c.mu.Unlock()
	if evicted {
		c.onEvicted(k, v)
	}

}

type kv struct {
	key   string
	value interface{}
}

func (c *cache) DeleteExpired() {
	var evictedItems []kv
	timeNow := time.Now().UnixNano()
	c.mu.Lock()
	for k, v := range c.items {
		if v.Expiration > 0 && v.Expiration < timeNow {
			v, evicted := c.delete(k)
			if evicted {
				evictedItems = append(evictedItems, kv{k, v})
			}
		}
	}
	c.mu.Unlock()
	for _, v := range evictedItems {
		c.onEvicted(v.key, v.value)
	}
}

func (c *cache) Item() map[string]Item {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.items
}

func (c *cache) ItemCount() int {
	c.mu.Lock()
	n := len(c.items)
	c.mu.Unlock()
	return n

}

func (c *cache) OnEvicted(f func(string, interface{})) {
	c.mu.Lock()
	c.onEvicted = f
	c.mu.Unlock()
}

func (c *cache) Flush() {
	c.mu.Lock()
	c.items = map[string]Item{}
	c.mu.Unlock()

}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *cache) {
	j.stop = make(chan bool)
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			//delete expired
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
	}

	c.janitor = j

	go j.Run(c)

}

func newCache(de time.Duration, m map[string]Item) *cache {
	if de == 0 {
		de = -1
	}
	c := &cache{
		defaultExpiration: de,
		items:             m,
	}
	return c
}

func newCacheWithJanitor(defaultExpiration time.Duration, cleanupInterval time.Duration, items map[string]Item) *Cache {
	c := newCache(defaultExpiration, items)
	//init Cache
	C := &Cache{c}
	if cleanupInterval > 0 {
		runJanitor(c, cleanupInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, items)
}
