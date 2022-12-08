package gp

import (
	"time"
)

type Pool struct {
	workers     chan worker
	idleRecycle time.Duration
}

// New create a new goroutine pool.
// The pool size is n, means that it will keep at most n goroutines in the pool.
// The dur parameter controls the idle recycle behaviour. If the goroutine in the pool is idle for a while, it will be recycled.
func New(n int, dur time.Duration) *Pool {
	return &Pool{
		workers:     make(chan worker, n),
		idleRecycle: dur,
	}
}

// Run execute the function in a seperate goroutine,
func (p *Pool) Go(f func()) {
	var w worker
	// Take a worker out from the pool
	select {
	case w = <-p.workers:
	default:
		w = worker{
			ch:   make(chan func()),
			Pool: p,
		}
		go workerGoroutine(w, p.idleRecycle)
	}
	// Let the worker run the task
	w.run(f)
}

type worker struct {
	ch chan func()
	*Pool
}

func (w worker) run(f func()) {
	w.ch <- f
}

// worker is bind with a goroutine,
func workerGoroutine(w worker, dur time.Duration) {
	t := time.NewTimer(dur)
	for {
		select {
		case f := <-w.ch:
			f()
			select {
			case w.Pool.workers <- w:
				if !t.Stop() {
					<-t.C
				}
				t.Reset(dur)
			default:
				return
			}
		case <-t.C:
			return
		}
	}
}
