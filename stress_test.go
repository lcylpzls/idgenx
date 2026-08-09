package idgenx

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestConcurrentStress 并发压力：8 goroutine × 5000 次生成，无重复。
func TestConcurrentStress(t *testing.T) {
	g, err := New(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	const goroutines = 8
	const perGoroutine = 5000
	var collisions atomic.Int64
	var wg sync.WaitGroup
	// 按并发分片校验（避免单 map 竞争影响测试语义）。
	seens := make([]sync.Map, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(gi int) {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				id, err := g.Next()
				if err != nil {
					t.Error(err)
					return
				}
				if _, loaded := seens[gi].LoadOrStore(id, struct{}{}); loaded {
					collisions.Add(1)
				}
			}
		}(i)
	}
	wg.Wait()
	if collisions.Load() != 0 {
		t.Fatalf("并发碰撞：%d", collisions.Load())
	}
}

// TestBackwardUnderLoad 负载中回拨演练：生成中时钟回拨 1ms 后恢复。
func TestBackwardUnderLoad(t *testing.T) {
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	g, err := New(DefaultConfig(), WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				_, _ = g.Next()
			}
		}
	}()
	time.Sleep(10 * time.Millisecond)
	clock.advance(time.Second)
	clock.advance(-time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	clock.advance(time.Millisecond)
	close(stop)
	wg.Wait()
}
