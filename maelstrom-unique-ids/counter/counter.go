package counter

import "sync"

type Counter struct {
	value uint64
	mu    sync.Mutex
}

func (c *Counter) Incr() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

func NewCounter() *Counter {
	return &Counter{value: 0}
}
