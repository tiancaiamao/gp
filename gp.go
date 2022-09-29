package gp

type Pool chan worker

func New(n int) Pool {
	return make(chan worker, n)
}

func (p Pool) Go(f func()) {
	var w worker
	select {
	case w = <-p:
	default:
		w = worker{
			ch:   make(chan func()),
			Pool: p,
		}
		go workerGoroutine(w)
	}

	w.run(f)
}

type worker struct {
	ch chan func()
	Pool
}

func (w worker) run(f func()) {
	w.ch <- f
}

// worker bind with a goroutine,
func workerGoroutine(w worker) {
	for f := range w.ch {
		f()

		select {
		case w.Pool <- w:
		default:
			return
		}

	}
}
