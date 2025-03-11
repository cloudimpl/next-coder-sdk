package polycode

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// entry represents a cache entry with a key, value, the last time it was accessed,
// and an optional fixed expiration time (if non-zero, the entry will not be removed
// due to sliding expiration until this time is reached).
type entry struct {
	key        string
	value      interface{}
	lastAccess time.Time
	expireAt   time.Time // zero value means no fixed expiration.
}

// LRUCache is a thread-safe cache that evicts least recently used items
// when a maximum capacity is exceeded and removes items that haven't been accessed
// within a specified TTL.
type LRUCache struct {
	capacity int
	ttl      time.Duration
	mu       sync.Mutex
	cache    map[string]*list.Element
	list     *list.List // Front = most-recently used, Back = least-recently used.
	quit     chan struct{}
}

// NewLRUCache creates a new LRUCache with the given maximum capacity,
// a time-to-live (TTL) for sliding expiration, and a cleanup interval for purging expired items.
func NewLRUCache(capacity int, ttl time.Duration, cleanupInterval time.Duration) *LRUCache {
	c := &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
		quit:     make(chan struct{}),
	}
	go c.startJanitor(cleanupInterval)
	return c
}

// startJanitor runs a background goroutine that periodically checks and removes expired items.
func (c *LRUCache) startJanitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.removeExpired()
		case <-c.quit:
			return
		}
	}
}

// removeExpired iterates over all entries and removes those that have expired.
// For entries with a fixed expiration (expireAt), if the current time is still before
// expireAt, they will not be removed regardless of sliding expiration.
func (c *LRUCache) removeExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	// Iterate over all items in the list.
	for elem := c.list.Back(); elem != nil; {
		prev := elem.Prev()
		ent := elem.Value.(*entry)
		// If a fixed expiration is set and it has not been reached, skip removal.
		if !ent.expireAt.IsZero() && now.Before(ent.expireAt) {
			// Do nothing; item is protected.
		} else if now.Sub(ent.lastAccess) > c.ttl {
			c.list.Remove(elem)
			delete(c.cache, ent.key)
		}
		elem = prev
	}
}

// Get retrieves a value from the cache by key.
// It returns false if the key is not found or has expired.
// If found, it updates the entry’s lastAccess time and moves it to the front.
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if elem, ok := c.cache[key]; ok {
		ent := elem.Value.(*entry)
		// If the entry has a fixed expiration and is still protected, return it.
		if !ent.expireAt.IsZero() && now.Before(ent.expireAt) {
			ent.lastAccess = now
			c.list.MoveToFront(elem)
			return ent.value, true
		}
		// Otherwise, use sliding expiration.
		if now.Sub(ent.lastAccess) > c.ttl {
			c.list.Remove(elem)
			delete(c.cache, key)
			return nil, false
		}
		ent.lastAccess = now
		c.list.MoveToFront(elem)
		return ent.value, true
	}
	return nil, false
}

// ComputeIfAbsent adds a key-value pair to the cache using sliding expiration only.
// If the key exists, it updates the value and resets its lastAccess time.
// If the cache is full, it evicts the least recently used item.
func (c *LRUCache) ComputeIfAbsent(key string, supplier func() (any, error)) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if elem, ok := c.cache[key]; ok {
		ent := elem.Value.(*entry)
		ent.lastAccess = now
		// For a normal Put, we clear any fixed expiration.
		ent.expireAt = time.Time{}
		c.list.MoveToFront(elem)
		return ent.value, nil
	}

	item, err := supplier()
	if err != nil {
		return nil, err
	}
	ent := &entry{
		key:        key,
		value:      item,
		lastAccess: now,
		// expireAt remains zero.
	}
	elem := c.list.PushFront(ent)
	c.cache[key] = elem

	if c.capacity != -1 && c.list.Len() > c.capacity {
		backElem := c.list.Back()
		if backElem != nil {
			backEnt := backElem.Value.(*entry)
			c.list.Remove(backElem)
			delete(c.cache, backEnt.key)
		}
	}

	return ent.value, nil
}

// Put adds a key-value pair to the cache using sliding expiration only.
// If the key exists, it updates the value and resets its lastAccess time.
// If the cache is full, it evicts the least recently used item.
func (c *LRUCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if elem, ok := c.cache[key]; ok {
		ent := elem.Value.(*entry)
		ent.value = value
		ent.lastAccess = now
		// For a normal Put, we clear any fixed expiration.
		ent.expireAt = time.Time{}
		c.list.MoveToFront(elem)
		return
	}

	ent := &entry{
		key:        key,
		value:      value,
		lastAccess: now,
		// expireAt remains zero.
	}
	elem := c.list.PushFront(ent)
	c.cache[key] = elem

	if c.capacity != -1 && c.list.Len() > c.capacity {
		backElem := c.list.Back()
		if backElem != nil {
			backEnt := backElem.Value.(*entry)
			c.list.Remove(backElem)
			delete(c.cache, backEnt.key)
		}
	}
}

// PutWithTimestamp adds a key-value pair to the cache with a fixed expiration timestamp.
// The entry will not be removed by sliding expiration until the current time is after notExpireUntil.
// If the key already exists, its value, lastAccess time, and fixed expiration are updated.
func (c *LRUCache) PutWithTimestamp(key string, value interface{}, notExpireUntil time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if elem, ok := c.cache[key]; ok {
		ent := elem.Value.(*entry)
		ent.value = value
		ent.lastAccess = now
		ent.expireAt = notExpireUntil
		c.list.MoveToFront(elem)
		return
	}

	ent := &entry{
		key:        key,
		value:      value,
		lastAccess: now,
		expireAt:   notExpireUntil,
	}
	elem := c.list.PushFront(ent)
	c.cache[key] = elem

	if c.capacity != -1 && c.list.Len() > c.capacity {
		backElem := c.list.Back()
		if backElem != nil {
			backEnt := backElem.Value.(*entry)
			c.list.Remove(backElem)
			delete(c.cache, backEnt.key)
		}
	}
}

// Stop terminates the background janitor goroutine.
func (c *LRUCache) Stop() {
	close(c.quit)
}

func main() {
	// Create a cache with maximum 3 items, a TTL of 5 seconds,
	// and a cleanup interval of 2 seconds.
	cache := NewLRUCache(3, 5*time.Second, 2*time.Second)
	defer cache.Stop()

	// Normal Put: sliding expiration applies.
	cache.Put("a", 1)
	cache.Put("b", 2)

	// PutWithTimestamp: "c" will not expire until 10 seconds from now.
	cache.PutWithTimestamp("c", 3, time.Now().Add(10*time.Second))

	// Access "a" to reset its TTL.
	if v, ok := cache.Get("a"); ok {
		fmt.Println("Key a:", v)
	}

	// Sleep for 3 seconds and then access "b" to update its lastAccess time.
	time.Sleep(3 * time.Second)
	if v, ok := cache.Get("b"); ok {
		fmt.Println("Key b:", v)
	}

	// Sleep enough so that "a" (normal sliding) expires.
	time.Sleep(6 * time.Second)
	if _, ok := cache.Get("a"); !ok {
		fmt.Println("Key a has expired due to inactivity")
	}

	// "c" should still be available because its fixed expiration time is in the future.
	if v, ok := cache.Get("c"); ok {
		fmt.Println("Key c still exists:", v)
	}

	// Add a new key to exceed the maximum capacity.
	cache.Put("d", 4)
	// Check which keys remain in the cache.
	if _, ok := cache.Get("b"); ok {
		fmt.Println("Key b still exists")
	}
	if _, ok := cache.Get("c"); ok {
		fmt.Println("Key c still exists")
	}
	if _, ok := cache.Get("d"); ok {
		fmt.Println("Key d exists")
	}
}
