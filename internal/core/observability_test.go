package core

import (
	"errors"
	testx "github.com/lcylpzls/testx"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lcylpzls/logx"
)

// testLogger 构造写入丢弃目标的日志器。
func testLogger() logx.Logger {
	logger, err := logx.NewBuilder().EnableWriter(io.Discard, logx.InfoLevel).Build()
	if err != nil {
		panic(err)
	}
	return logger
}

// TestObservability 覆盖日志与指标回调路径。
func TestObservability(t *testing.T) {
	var generated, rejected, waitMS atomic.Int64
	m := Metrics{
		Generated: func(_ int64, delta int) { generated.Add(int64(delta)) },
		Rejected:  func(_ int64, _ error) { rejected.Add(1) },
		WaitMS:    func(_ int64, ms int64) { waitMS.Add(ms) },
	}
	base := time.Date(2026, 8, 9, 10, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: base}
	g, err := New(DefaultConfig(), WithClock(clock.get), WithLogger(testLogger()), WithMetrics(m))
	testx.RequireNoError(t, err)

	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(time.Second)
	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
	if generated.Load() != 2 {
		t.Fatalf("生成计数应为 2：%d", generated.Load())
	}
	// 回拨等待成功。
	done := make(chan error, 1)
	clock.advance(-time.Millisecond)
	go func() {
		_, err := g.Next()
		done <- err
	}()
	time.Sleep(time.Millisecond)
	clock.advance(time.Millisecond)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if waitMS.Load() < 0 {
		t.Fatal("等待耗时指标应非负")
	}
	// 回拨拒绝（Reject 策略）计数。
	cfg := DefaultConfig()
	cfg.Backward = StrategyReject
	g2, err := New(cfg, WithClock(clock.get), WithLogger(testLogger()), WithMetrics(m))
	testx.RequireNoError(t, err)

	if _, err := g2.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(time.Second)
	if _, err := g2.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(-time.Millisecond)
	if _, err := g2.Next(); !errors.Is(err, ErrClockBackward) {
		t.Fatalf("应拒绝：%v", err)
	}
	if rejected.Load() != 1 {
		t.Fatalf("拒绝计数应为 1：%d", rejected.Load())
	}
	// 等待超时计数。
	cfg2 := DefaultConfig()
	cfg2.MaxWait = time.Millisecond
	g3, err := New(cfg2, WithClock(clock.get), WithLogger(testLogger()), WithMetrics(m))
	testx.RequireNoError(t, err)

	if _, err := g3.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(time.Second)
	if _, err := g3.Next(); err != nil {
		t.Fatal(err)
	}
	clock.advance(-20 * time.Millisecond)
	if _, err := g3.Next(); !errors.Is(err, ErrWaitTimeout) {
		t.Fatalf("应超时：%v", err)
	}
	if rejected.Load() != 2 {
		t.Fatalf("拒绝计数应为 2：%d", rejected.Load())
	}
}

// TestLoggerNil 覆盖无日志器路径。
func TestLoggerNil(t *testing.T) {
	g, err := New(DefaultConfig())
	testx.RequireNoError(t, err)

	if _, err := g.Next(); err != nil {
		t.Fatal(err)
	}
}

// TestLayoutBoundaries 覆盖极端位布局。
func TestLayoutBoundaries(t *testing.T) {
	epoch := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fixedClock{now: epoch.Add(time.Millisecond)} // 时间戳 = 1ms。
	cases := []Config{
		{TimestampBits: 61, NodeBits: 1, SequenceBits: 1, NodeID: 1},
		{TimestampBits: 1, NodeBits: 31, SequenceBits: 31, NodeID: 1<<31 - 1},
	}
	for _, cfg := range cases {
		cfg.Epoch = epoch
		g, err := New(cfg, WithClock(clock.get))
		testx.RequireNoError(t, err)

		id, err := g.Next()
		testx.RequireNoError(t, err)

		parts, err := g.Parse(id)
		testx.RequireNoError(t, err)

		testx.RequireEqual(t, parts.NodeID, cfg.NodeID)

	}
}

// TestConcurrentMetrics 覆盖并发指标回调无竞态（配合 race）。
func TestConcurrentMetrics(t *testing.T) {
	var wg sync.WaitGroup
	g, err := New(DefaultConfig(), WithMetrics(Metrics{
		Generated: func(int64, int) {},
	}))
	testx.RequireNoError(t, err)

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				_, _ = g.Next()
			}
		}()
	}
	wg.Wait()
}

// TestEpochNearNow 覆盖纪元贴近当前时间的边界。
func TestEpochNearNow(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Epoch = time.Now().Add(-time.Millisecond)
	g, err := New(cfg)
	testx.RequireNoError(t, err)

	id, err := g.Next()
	testx.RequireNoError(t, err)

	parts, err := g.Parse(id)
	testx.RequireNoError(t, err)

	if parts.Timestamp.IsZero() {
		t.Fatal("时间戳解析不应为零")
	}
}
