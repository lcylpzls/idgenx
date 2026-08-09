package idgenx

import (
	"sync"
	"time"

	"github.com/lcylpzls/errx"
)

const defaultMaxWait = 5 * time.Millisecond

// Parts 雪花 ID 解析结果。
type Parts struct {
	// Timestamp 生成时刻（纪元后毫秒换算）。
	Timestamp time.Time
	// NodeID 节点 ID。
	NodeID int64
	// Sequence 同毫秒序列号。
	Sequence int64
}

// Option 生成器配置项。
type Option func(*Generator)

// WithClock 注入时间源（测试用）。
func WithClock(now func() time.Time) Option {
	return func(g *Generator) {
		if now != nil {
			g.now = now
		}
	}
}

// Generator 雪花 ID 生成器（并发安全）。
type Generator struct {
	cfg          Config
	mu           sync.Mutex
	now          func() time.Time
	epochMS      int64
	sequenceBits uint8
	nodeShift    uint8
	nodeMask     int64
	seqMask      int64
	maxSequence  int64
	lastTs       int64
	sequence     int64
}

// New 构造雪花 ID 生成器。
func New(cfg Config, opts ...Option) (*Generator, error) {
	cfg = cfg.normalize()
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}
	g := &Generator{
		cfg:          cfg,
		now:          time.Now,
		epochMS:      cfg.Epoch.UnixMilli(),
		sequenceBits: cfg.SequenceBits,
		nodeShift:    cfg.NodeBits + cfg.SequenceBits,
		seqMask:      (int64(1) << cfg.SequenceBits) - 1,
		maxSequence:  (int64(1) << cfg.SequenceBits) - 1,
	}
	g.nodeMask = (int64(1) << cfg.NodeBits) - 1
	for _, opt := range opts {
		opt(g)
	}
	return g, nil
}

// Next 生成下一个严格递增的雪花 ID。
func (g *Generator) Next() (int64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.now().UnixMilli() - g.epochMS
	if now < g.lastTs {
		if err := g.waitBackward(now); err != nil {
			return 0, err
		}
		now = g.now().UnixMilli() - g.epochMS
	}
	if now == g.lastTs {
		g.sequence++
		if g.sequence > g.maxSequence {
			for now <= g.lastTs {
				now = g.now().UnixMilli() - g.epochMS
			}
			g.sequence = 0
		}
	} else {
		g.sequence = 0
	}
	g.lastTs = now
	return (now << g.nodeShift) | (g.cfg.NodeID << g.sequenceBits) | g.sequence, nil
}

// Parse 解析雪花 ID 的时间戳、节点与序列。
func (g *Generator) Parse(id int64) (Parts, error) {
	if id < 0 {
		return Parts{}, ErrInvalidID
	}
	ts := id >> g.nodeShift
	node := (id >> g.sequenceBits) & g.nodeMask
	seq := id & g.seqMask
	return Parts{
		Timestamp: time.UnixMilli(g.epochMS + ts),
		NodeID:    node,
		Sequence:  seq,
	}, nil
}

// waitBackward 等待时钟追平（固定上限 defaultMaxWait）。
func (g *Generator) waitBackward(now int64) error {
	start := time.Now()
	for now < g.lastTs {
		if time.Since(start) > defaultMaxWait {
			return ErrWaitTimeout
		}
		time.Sleep(100 * time.Microsecond)
		now = g.now().UnixMilli() - g.epochMS
	}
	return nil
}

// validateConfig 校验配置。
func validateConfig(cfg Config) error {
	if 1+int(cfg.TimestampBits)+int(cfg.NodeBits)+int(cfg.SequenceBits) != 64 {
		return errInvalid("位布局总和必须为 63（符号位 1 位）")
	}
	if cfg.Epoch.IsZero() || cfg.Epoch.After(time.Now()) {
		return errInvalid("纪元必须早于当前时间")
	}
	if cfg.NodeID < 0 || cfg.NodeID > (int64(1)<<cfg.NodeBits)-1 {
		return ErrNodeInvalid
	}
	return nil
}

// errInvalid 构造配置错误。
func errInvalid(msg string) error {
	return errx.New(errx.KindInvalid, CodeInvalidConfig, msg)
}
