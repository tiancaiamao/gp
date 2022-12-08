package gp

import (
	"sync"
	"testing"
	"time"
)

func TestT(t *testing.T) {
	var wg sync.WaitGroup
	const N = 30
	wg.Add(30)

	pool := New(3, 10*time.Second)
	a := 3
	for i := 0; i < N; i++ {
		pool.Go(func() {
			a++
			wg.Done()
		})
	}

	wg.Wait()
	if a != N+3 {
		t.Fail()
	}
}
