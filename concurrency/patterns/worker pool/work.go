package main

import "sync"

type Worker interface {
	Task()
}

type Pool struct {
	tasks chan Worker
	wg    sync.WaitGroup
}

// NewPool creates a new work pool
func NewPool(maxWorkers int) *Pool {
	p := Pool{
		tasks: make(chan Worker),
	}

	for i := 0; i < maxWorkers; i++ {
		p.wg.Add(1)
		go func() {
			for w := range p.tasks {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Add submits work to the pool
func (p *Pool) Add(w Worker) {
	p.tasks <- w
}

// Shutdown waits for all the gorutines to shutdown
func (p *Pool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}
