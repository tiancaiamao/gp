package gp

type Pool chan worker

func New(n int) Pool {
	wp := make(chan worker, n)
	for i := 0; i < n; i++ {
		w := worker{
			ch:         make(chan Runnable),
			Pool: wp,
		}
		go workerGoroutine(w)
		wp <- w
	}
	return wp
}

type Runnable interface {
	Run()
}

type worker struct {
	ch chan Runnable
	Pool
}

func (w worker) Run(f Runnable) {
	w.ch <- f
}

// worker bind with a goroutine,
func workerGoroutine(w worker) {
	for f := range w.ch {
		f.Run()
		// What's special about it is that the worker itself is put back to the pool,
		// after handling one task.
		w.Pool <- w
	}
}
