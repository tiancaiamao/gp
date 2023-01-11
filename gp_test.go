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
	for i := 0; i < 10; i++ {
		pool.Go(func() {})
	}
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 100; i++ {
		pool.Go(func() {
			defer func() {
				if r := recover(); r != nil {
				}
			}()
			panic(1)
		})
	}
	time.Sleep(100 * time.Millisecond)
	if len(pool.count) != 2 {
		t.Fail()
	}
}

func TestIdleRecycle(t *testing.T) {
	pool := New(10, 100*time.Millisecond)
	for i := 0; i < 1000; i++ {
		pool.Go(func() {})
	}
	time.Sleep(300 * time.Millisecond)
	if len(pool.count) != 0 {
		t.Fail()
	}
	pool.Close()

	pool = New(10, 0)
	for i := 0; i < 1000; i++ {
		pool.Go(func() {})
	}
	time.Sleep(300 * time.Millisecond)
	if len(pool.count) == 0 {
		t.Fail()
	}
	pool.Close()
}

func TestManOrBoy(t *testing.T) {
	pool := New(5, 100*time.Millisecond)
	t.Run("Panic", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(5)
		for i := 0; i < 5; i++ {
			pool.Go(func() {
				defer func() {
					recover()
					wg.Done()
				}()
				panic("hey boy")
			})
		}
		wg.Wait()

		var x = 0
		wg.Add(1)
		pool.Go(func() {
			x = 42
			wg.Done()
		})
		wg.Wait()
		if x != 42 {
			t.Fail()
		}
	})

	t.Run("Sleep", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			pool.Go(func() {
				time.Sleep(time.Hour)
			})
		}
		time.Sleep(100 * time.Millisecond)
		var wg sync.WaitGroup
		var x = 0
		wg.Add(1)
		pool.Go(func() {
			x = 42
			wg.Done()
		})
		wg.Wait()
		if x != 42 {
			t.Fail()
		}
	})

	t.Run("Block", func(t *testing.T) {
		ch := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(5)
		for i := 0; i < 5; i++ {
			pool.Go(func() {
				<-ch
				wg.Done()
			})
		}
		pool.Go(func() {
			close(ch)
		})
		wg.Wait()
	})
}
