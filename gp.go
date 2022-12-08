package gp

import (
	"time"
)

type Pool struct {
	workers     chan worker
	idleRecycle time.Duration
}

func New(n int, dur time.Duration) *Pool {
	return &Pool{
		workers:     make(chan worker, n),
		idleRecycle: dur,
	}
}

func (p *Pool) Go(f func()) {
	var w worker
	// Get a worker from the worker pool
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
