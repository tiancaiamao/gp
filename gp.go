package gp

type Pool chan worker

func New(n int) Pool {
	wp := make(chan worker, n)
	for i := 0; i < n; i++ {
		w := worker{
			ch:         make(chan func()),
			Pool: wp,
		}
		go workerGoroutine(w)
		wp <- w
	}
	return wp
}

type worker struct {
	ch chan func()
	Pool
}

func (w worker) Run(f func()) {
	w.ch <- f
}

// worker bind with a goroutine,
func workerGoroutine(w worker) {
	for f := range w.ch {
		f()
		// What's special about it is that the worker itself is put back to the pool,
		// after handling one task.
		w.Pool <- w
	}
}
