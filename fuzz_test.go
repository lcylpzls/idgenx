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
