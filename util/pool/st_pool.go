// map_create_pool.go
package pool

import (
	"sync"
)

type st struct {
	a int
	b string
	c string
	d string
	e string
	f string
}

type stPool struct {
	iq      *intQueue
	mut     sync.RWMutex
	mapping map[int]int
	p       []st
}

type intQueue struct {
	mut     sync.Mutex
	popPos  int
	pushPos int
	bs      []int
}

func NewIntQueue(c int) *intQueue {
	return &intQueue{bs: make([]int, c, c)}
}

func (this *intQueue) Pop() (int, bool) {
	this.mut.Lock()
	if this.popPos == this.pushPos {
		this.mut.Unlock()
		return 0, false
	}
	v := this.bs[this.popPos]
	this.popPos++
	this.mut.Unlock()
	return v, true
}

func (this *intQueue) Push(v int) {
	this.mut.Lock()
	this.bs[this.pushPos] = v
	this.pushPos++
	if this.pushPos >= len(this.bs) {
		this.pushPos = 0
	}
	if this.pushPos == this.popPos {
		l := 2 * len(this.bs)
		bs := make([]int, l, l)
		copy(bs, this.bs)
		this.popPos = 0
		this.pushPos = len(this.bs)
		this.bs = bs
	}
	this.mut.Unlock()
}

func NewSTPool(c int) *stPool {
	mapping := make(map[int]int, c)
	p := make([]st, c, c)
	iq := NewIntQueue(c + 10)
	for i := 0; i < c; i++ {
		iq.Push(i)
	}
	return &stPool{mapping: mapping, p: p, iq: iq}
}

func (this *stPool) alloc(k int) (*st, bool) {
	this.mut.RLock()
	if v, ok := this.mapping[k]; ok {
		this.mut.RUnlock()
		return &this.p[v], true
	}
	this.mut.RUnlock()
	v, ok := this.iq.Pop()
	if !ok {
		return nil, false
	}
	this.mut.Lock()
	this.mapping[k] = v
	this.mut.Unlock()
	return &this.p[v], true
}

func (this *stPool) release(k int) {
	this.mut.Lock()
	if v, ok := this.mapping[k]; ok {
		delete(this.mapping, k)
		this.mut.Unlock()
		this.iq.Push(v)
	} else {
		this.mut.Unlock()
	}
}
