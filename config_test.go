package idgenx

import (
	"testing"
	"time"

	"github.com/lcylpzls/errx"
)

// TestDefaultConfig 覆盖默认配置值。
func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TimestampBits != 41 || cfg.NodeBits != 10 || cfg.SequenceBits != 12 ||
		cfg.NodeID != 0 || cfg.Epoch.IsZero() {
		t.Fatalf("默认配置不符：%+v", cfg)
	}
}

// TestNewErrors 覆盖配置校验分支。
func TestNewErrors(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name string
		cfg  Config
		want errx.Code
	}{
		{"位总和不足", Config{TimestampBits: 40, NodeBits: 10, SequenceBits: 12}, CodeInvalidConfig},
		{"纪元晚于当前", Config{Epoch: now.Add(time.Hour), TimestampBits: 41, NodeBits: 10, SequenceBits: 12}, CodeInvalidConfig},
		{"节点越界", Config{TimestampBits: 41, NodeBits: 10, SequenceBits: 12, NodeID: 1 << 10}, CodeNodeInvalid},
		{"负节点", Config{TimestampBits: 41, NodeBits: 10, SequenceBits: 12, NodeID: -1}, CodeNodeInvalid},
		{"非法回拨策略", Config{TimestampBits: 41, NodeBits: 10, SequenceBits: 12, Backward: BackwardStrategy(99)}, CodeInvalidConfig},
		{"负等待上限", Config{TimestampBits: 41, NodeBits: 10, SequenceBits: 12, MaxWait: -time.Second}, CodeInvalidConfig},
	}
	for _, tc := range cases {
		if _, err := New(tc.cfg); err == nil || !errx.Is(err, tc.want) {
			t.Fatalf("%s 应报 %s，实际：%v", tc.name, tc.want, err)
		}
	}
}

// TestNewDefaults 覆盖零值字段填充默认。
func TestNewDefaults(t *testing.T) {
	g, err := New(Config{NodeID: 3})
	if err != nil {
		t.Fatal(err)
	}
	if g.cfg.TimestampBits != 41 || g.cfg.NodeBits != 10 || g.cfg.SequenceBits != 12 {
		t.Fatalf("默认填充失败：%+v", g.cfg)
	}
	id, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	parts, err := g.Parse(id)
	if err != nil {
		t.Fatal(err)
	}
	if parts.NodeID != 3 {
		t.Fatalf("节点不符：%d", parts.NodeID)
	}
}
