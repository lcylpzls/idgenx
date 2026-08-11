package idgenx_test

import (
	"testing"
	"time"

	"github.com/lcylpzls/idgenx"
)

// TestPublicAPI 黑盒冒烟测试：覆盖根包全部转发函数、类型别名与常量。
func TestPublicAPI(t *testing.T) {
	if idgenx.Version != "v1.5.2" {
		t.Fatalf("Version 不符：%s", idgenx.Version)
	}

	cfg := idgenx.DefaultConfig()
	if cfg.Epoch.IsZero() {
		t.Fatal("DefaultConfig 返回零值配置")
	}

	g, err := idgenx.New(cfg,
		idgenx.WithClock(time.Now),
		idgenx.WithLogger(nil),
		idgenx.WithMetrics(idgenx.Metrics{}),
	)
	if err != nil || g == nil {
		t.Fatalf("New 失败：%v", err)
	}

	_, err = idgenx.RandomHex(8)
	if err != nil {
		t.Fatalf("RandomHex 失败：%v", err)
	}
	_, err = idgenx.RandomBase64URL(8)
	if err != nil {
		t.Fatalf("RandomBase64URL 失败：%v", err)
	}

	var _ idgenx.Metrics
	var _ idgenx.Parts
	var _ idgenx.Option
	var _ idgenx.Generator
	var _ idgenx.BackwardStrategy
	var _ idgenx.Config
	_ = idgenx.CodeInvalidConfig
	_ = idgenx.CodeClockBackward
	_ = idgenx.CodeWaitTimeout
	_ = idgenx.CodeNodeInvalid
	_ = idgenx.CodeInvalidID
}
