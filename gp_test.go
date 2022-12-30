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

func TestPanic(t *testing.T) {
	pool := New(2, 0)
	// Fill the pool
	for i:=0; i<10; i++ {
		pool.Go(func() {})
	}
	time.Sleep(100*time.Millisecond)
	for i := 0; i < 100; i++ {
		pool.Go(func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			panic(1)
		})
	}
	time.Sleep(100*time.Millisecond)
	if len(pool.count) != 2 {
		t.Fail()
	}
}

func TestIdleRecycle(t *testing.T) {
	pool := New(10, 100*time.Millisecond)
	for i:=0; i<1000; i++ {
		pool.Go(func(){})
	}
	time.Sleep(300*time.Millisecond)
	if len(pool.count) != 0 {
		t.Fail()
	}
	pool.Close()

	pool = New(10, 0)
	for i:=0; i<1000; i++ {
		pool.Go(func(){})
	}
	time.Sleep(300*time.Millisecond)
	if len(pool.count) == 0 {
		t.Fail()
	}
	pool.Close()
}
