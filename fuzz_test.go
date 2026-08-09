package idgenx

import (
	"testing"
	"time"
)

// FuzzConfig 模糊测试配置校验与生成路径，确保任意输入不 panic。
func FuzzConfig(f *testing.F) {
	f.Add(int64(0), uint8(41), uint8(10), uint8(12))
	f.Add(int64(-1), uint8(0), uint8(0), uint8(0))
	f.Add(int64(1<<20), uint8(63), uint8(63), uint8(63))
	f.Fuzz(func(t *testing.T, nodeID int64, tsBits, nodeBits, seqBits uint8) {
		cfg := Config{
			Epoch:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			TimestampBits: tsBits,
			NodeBits:      nodeBits,
			SequenceBits:  seqBits,
			NodeID:        nodeID,
		}
		g, err := New(cfg)
		if err != nil {
			return
		}
		_, _ = g.Next()
		_, _ = g.Parse(0)
	})
}

// FuzzNext 模糊测试生成路径：默认布局 + 随机回拨策略下
// 不 panic、有限步返回（位布局边界由 FuzzConfig 与单元测试覆盖）。
func FuzzNext(f *testing.F) {
	f.Add(uint8(0))
	f.Add(uint8(2))
	f.Fuzz(func(t *testing.T, strategy uint8) {
		cfg := Config{
			Epoch:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Backward: BackwardStrategy(strategy % 3),
			MaxWait:  time.Millisecond,
		}
		g, err := New(cfg)
		if err != nil {
			return
		}
		for i := 0; i < 100; i++ {
			_, _ = g.Next()
		}
	})
}
