package idgenx_test

import (
	"sync"
	"testing"

	"github.com/lcylpzls/idgenx"
	"github.com/lcylpzls/idgenx/shortid"
)

// BenchmarkSnowflakeNext 测量雪花 ID 串行生成吞吐。
func BenchmarkSnowflakeNext(b *testing.B) {
	g, err := idgenx.New(idgenx.DefaultConfig())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = g.Next()
	}
}

// BenchmarkSnowflakeParallel 测量并发生成吞吐。
func BenchmarkSnowflakeParallel(b *testing.B) {
	g, err := idgenx.New(idgenx.DefaultConfig())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = g.Next()
		}
	})
}

// BenchmarkShortIDGenerate 测量短 ID 生成吞吐。
func BenchmarkShortIDGenerate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = shortid.Generate(8)
	}
}

// BenchmarkSnowflakeParallelVerify 并行生成并校验无重复（正确性基准）。
func BenchmarkSnowflakeParallelVerify(b *testing.B) {
	g, err := idgenx.New(idgenx.DefaultConfig())
	if err != nil {
		b.Fatal(err)
	}
	seen := make(map[int64]struct{}, 1<<16)
	var mu sync.Mutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id, err := g.Next()
			if err != nil {
				b.Fatal(err)
			}
			mu.Lock()
			if _, ok := seen[id]; ok {
				b.Fatal("重复 ID")
			}
			seen[id] = struct{}{}
			mu.Unlock()
		}
	})
}
