package idgenx

import "time"

const (
	defaultEpoch         = "2024-01-01T00:00:00Z"
	defaultTimestampBits = 41
	defaultNodeBits      = 10
	defaultSequenceBits  = 12
)

// Config 雪花 ID 生成器配置。
type Config struct {
	// Epoch 时间戳纪元（默认 2024-01-01T00:00:00Z）。
	Epoch time.Time
	// TimestampBits 时间戳位数（默认 41）。
	TimestampBits uint8
	// NodeBits 节点位数（默认 10）。
	NodeBits uint8
	// SequenceBits 序列位数（默认 12）。
	SequenceBits uint8
	// NodeID 本实例节点 ID（默认 0，范围 [0, 2^NodeBits)）。
	NodeID int64
}

// DefaultConfig 返回默认布局（41/10/12，纪元 2024-01-01）。
func DefaultConfig() Config {
	epoch, _ := time.Parse(time.RFC3339, defaultEpoch)
	return Config{
		Epoch:         epoch,
		TimestampBits: defaultTimestampBits,
		NodeBits:      defaultNodeBits,
		SequenceBits:  defaultSequenceBits,
		NodeID:        0,
	}
}

// normalize 填充零值字段为默认值并计算派生量。
func (c Config) normalize() Config {
	if c.Epoch.IsZero() {
		c.Epoch, _ = time.Parse(time.RFC3339, defaultEpoch)
	}
	if c.TimestampBits == 0 {
		c.TimestampBits = defaultTimestampBits
	}
	if c.NodeBits == 0 {
		c.NodeBits = defaultNodeBits
	}
	if c.SequenceBits == 0 {
		c.SequenceBits = defaultSequenceBits
	}
	return c
}
