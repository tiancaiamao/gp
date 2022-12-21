package gp

import (
	"time"
)

type Pool struct {
	ch          chan func()
	count       chan struct{}
	idleRecycle time.Duration
	closed      chan struct{}
}

// New create a new goroutine pool.
// The pool size is n, means that it will keep at most n goroutines in the pool.
// The dur parameter controls the idle recycle behaviour. If the goroutine in the pool is idle for a while, it will be recycled.
func New(n int, dur time.Duration) *Pool {
	return &Pool{
		ch:          make(chan func()),
		count:       make(chan struct{}, n),
		idleRecycle: dur,
		closed:      make(chan struct{}),
	}
}

// Run execute the function in a seperate goroutine,
func (p *Pool) Go(f func()) {
	select {
	case p.ch <- f:
	case <-p.closed:
	default:
		go worker(p)
		p.ch <- f
	}
}

func worker(p *Pool) {
	t := time.NewTimer(p.idleRecycle)
	var inPool bool
	for {
		select {
		case f := <-p.ch:
			f()

			// When worker finish a task, it would decide whether to reuse.
			select {
			case p.count <- struct{}{}:
				if !t.Stop() {
					<-t.C
				}
				t.Reset(p.idleRecycle)
				inPool = true
			default:
				return
			}
		case <-t.C:
			if inPool {
				<-p.count
			}
			return
		case <-p.closed:
			return
		}
	}
}

// Close releases the goroutines in the Pool.
// After this operation, inflight tasks may still execute until finish.
// But all the new coming tasks will be simply ignored.
func (p *Pool) Close() {
	close(p.closed)
}
