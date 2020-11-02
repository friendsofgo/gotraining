package main

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrPoolSizeTooSmall = errors.New("pool size is too small")
	ErrPoolClosed       = errors.New("pool has been closed")
)

// FactoryFunc a func tha allow to creates new elements inside the pool with the user logic
type FactoryFunc func() (io.Closer, error)

type Pool struct {
	sync.Mutex
	resources chan io.Closer
	closed    bool
}

// NewPool creates a pool that manages resources
func NewPool(size int, fn FactoryFunc) (*Pool, error) {
	if size <= 0 {
		return nil, ErrPoolSizeTooSmall
	}

	resources := make(chan io.Closer, size)
	for i := 0; i < size; i++ {
		r, err := fn()
		if err != nil {
			return nil, err
		}
		resources <- r
	}

	return &Pool{
		resources: resources,
	}, nil
}

// Get retrieves a resource from the pool
func (p *Pool) Get() (io.Closer, error) {
	r, ok := <-p.resources
	if !ok {
		return nil, ErrPoolClosed
	}
	return r, nil
}

// Release places a new resource onto the pool
func (p *Pool) Release(r io.Closer) {
	p.Lock()
	defer p.Unlock()

	if p.closed {
		r.Close()
		return
	}

	p.resources <- r
}

func (p *Pool) Close() {
	p.Lock()
	defer p.Unlock()

	if p.closed {
		return
	}

	p.closed = true
	close(p.resources)
	for r := range p.resources {
		r.Close()
	}
}
