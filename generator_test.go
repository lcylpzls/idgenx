package idgenx

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// fixedClock 可推进的测试时钟。
type fixedClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fixedClock) get() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fixedClock) advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

// TestNextMonotonic 覆盖严格递增与解析一致性。
func TestNextMonotonic(t *testing.T) {
	g, err := New(DefaultConfig(), WithClock(time.Now))
	if err != nil {
		t.Fatal(err)
	}
	var prev int64
	for i := 0; i < 10000; i++ {
		id, err := g.Next()
		if err != nil {
			t.Fatal(err)
		}
		if id <= prev {
			t.Fatalf("ID 应严格递增：%d <= %d", id, prev)
		}
		prev = id
		parts, err := g.Parse(id)
		if err != nil {
			t.Fatal(err)
		}
		if parts.NodeID != 0 || parts.Timestamp.IsZero() {
			t.Fatalf("解析不符：%+v", parts)
		}
	}
}

// TestNextSequenceOverflow 覆盖同毫秒序列溢出等待。
func TestNextSequenceOverflow(t *testing.T) {
	orig := randRead
	randRead = func(b []byte) (int, error) {
		for i := range b {
			b[i] = 0
		}
		return len(b), nil
	}
	defer func() { randRead = orig }()
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	g, err := New(DefaultConfig(), WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	// 同毫秒生成 maxSequence+1 个（溢出点）。
	for i := 0; i <= int(g.maxSequence); i++ {
		if _, err := g.Next(); err != nil {
			t.Fatal(err)
		}
	}
	clock.advance(time.Millisecond)
	id, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	parts, err := g.Parse(id)
	if err != nil {
		t.Fatal(err)
	}
	if parts.Sequence != 0 {
		t.Fatalf("溢出后序列应归零：%d", parts.Sequence)
	}
	if !parts.Timestamp.After(base) {
		t.Fatalf("溢出后应进入下一毫秒：%v", parts.Timestamp)
	}
}

// TestBackwardWait 覆盖回拨等待成功。
func TestBackwardWait(t *testing.T) {
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	g, err := New(DefaultConfig(), WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(3 * time.Millisecond)
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	// 回拨 1ms（仍在 wait 上限内）。
	clock.advance(-time.Millisecond)
	done := make(chan error, 1)
	go func() {
		_, err := g.Next()
		done <- err
	}()
	time.Sleep(time.Millisecond) // 让等待循环开始轮询。
	clock.advance(time.Millisecond)
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("小回拨应等待成功：%v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("回拨等待未完成")
	}
}

// TestBackwardTimeout 覆盖回拨等待超时。
func TestBackwardTimeout(t *testing.T) {
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	g, err := New(DefaultConfig(), WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(100 * time.Millisecond)
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	// 回拨 20ms，超过 5ms 上限。
	clock.advance(-20 * time.Millisecond)
	if _, err := g.Next(); !errors.Is(err, ErrWaitTimeout) {
		t.Fatalf("大回拨应超时，实际：%v", err)
	}
}

// TestParseErrors 覆盖解析错误分支。
func TestParseErrors(t *testing.T) {
	g, err := New(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Parse(-1); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("负 ID 应报错，实际：%v", err)
	}
}

// TestConcurrent 覆盖并发生成无重复。
func TestConcurrent(t *testing.T) {
	g, err := New(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	const goroutines = 8
	const perGoroutine = 2000
	ids := make(chan int64, goroutines*perGoroutine)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				id, err := g.Next()
				if err != nil {
					t.Error(err)
					return
				}
				ids <- id
			}
		}()
	}
	wg.Wait()
	close(ids)
	seen := make(map[int64]struct{}, goroutines*perGoroutine)
	for id := range ids {
		if _, ok := seen[id]; ok {
			t.Fatalf("重复 ID：%d", id)
		}
		seen[id] = struct{}{}
	}
	if len(seen) != goroutines*perGoroutine {
		t.Fatalf("ID 数量不符：%d", len(seen))
	}
}

// TestWithClockNil 覆盖空时钟注入保持默认。
func TestWithClockNil(t *testing.T) {
	g, err := New(DefaultConfig(), WithClock(nil))
	if err != nil {
		t.Fatal(err)
	}
	if g.now == nil {
		t.Fatal("默认时间源不应为空")
	}
}

// TestBackwardReject 覆盖 Reject 策略。
func TestBackwardReject(t *testing.T) {
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	cfg := DefaultConfig()
	cfg.Backward = StrategyReject
	g, err := New(cfg, WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(time.Second)
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(-time.Millisecond)
	if _, err := g.Next(); !errors.Is(err, ErrClockBackward) {
		t.Fatalf("Reject 策略应拒绝，实际：%v", err)
	}
}

// TestBackwardLoose 覆盖 Loose 策略（沿用时间戳继续递增）。
func TestBackwardLoose(t *testing.T) {
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	cfg := DefaultConfig()
	cfg.Backward = StrategyLoose
	g, err := New(cfg, WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	id1, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	clock.advance(time.Second)
	id2, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	// 回拨 1ms。
	clock.advance(-time.Millisecond)
	id3, err := g.Next()
	if err != nil {
		t.Fatalf("Loose 策略不应报错：%v", err)
	}
	if id3 <= id2 {
		t.Fatalf("Loose 回拨后 ID 应继续递增：%d <= %d", id3, id2)
	}
	p2, _ := g.Parse(id2)
	p3, _ := g.Parse(id3)
	if !p3.Timestamp.Equal(p2.Timestamp) || p3.Sequence != p2.Sequence+1 {
		t.Fatalf("Loose 应沿用时间戳并递增序列：%+v %+v", p2, p3)
	}
	_ = id1
}

// TestNewRandFailure 覆盖随机源失败。
func TestNewRandFailure(t *testing.T) {
	orig := randRead
	randRead = func(b []byte) (int, error) { return 0, errors.New("随机源故障") }
	defer func() { randRead = orig }()
	if _, err := New(DefaultConfig()); !errors.Is(err, ErrRandomFailure) {
		t.Fatalf("随机源失败应报错，实际：%v", err)
	}
}

// TestSequenceRandomStart 覆盖序列随机起点。
func TestSequenceRandomStart(t *testing.T) {
	orig := randRead
	randRead = func(b []byte) (int, error) {
		for i := range b {
			b[i] = 0
		}
		b[len(b)-1] = 0xFF // 序列低 8 位 = 255。
		return len(b), nil
	}
	defer func() { randRead = orig }()
	g, err := New(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	id, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	parts, err := g.Parse(id)
	if err != nil {
		t.Fatal(err)
	}
	if parts.Sequence != 255 {
		t.Fatalf("随机起点应保留注入值 255：%d", parts.Sequence)
	}
}

// TestMaxWaitConfig 覆盖自定义等待上限。
func TestMaxWaitConfig(t *testing.T) {
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	cfg := DefaultConfig()
	cfg.MaxWait = time.Millisecond
	g, err := New(cfg, WithClock(clock.get))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(100 * time.Millisecond)
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(-20 * time.Millisecond)
	if _, err := g.Next(); !errors.Is(err, ErrWaitTimeout) {
		t.Fatalf("自定义 1ms 上限应超时，实际：%v", err)
	}
}
