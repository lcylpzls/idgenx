package idgenx

import (
	"testing"

	"github.com/lcylpzls/errx"
)

// TestErrVarKinds 保证预定义错误值通过 NewCode 构造后分类正确
// （注册必须先于包级变量初始化）。
func TestErrVarKinds(t *testing.T) {
	cases := map[string]struct {
		err  error
		kind errx.Kind
	}{
		"配置非法":  {ErrInvalidConfig, errx.KindInvalid},
		"时钟回拨":  {ErrClockBackward, errx.KindUnavailable},
		"等待超时":  {ErrWaitTimeout, errx.KindTimeout},
		"节点越界":  {ErrNodeInvalid, errx.KindInvalid},
		"ID 非法": {ErrInvalidID, errx.KindInvalid},
		"随机源失败": {ErrRandomFailure, errx.KindUnavailable},
		"碰撞耗尽":  {ErrCollision, errx.KindConflict},
		"时间戳溢出": {ErrTimestampOverflow, errx.KindUnavailable},
	}
	for name, tc := range cases {
		if got := errx.KindOf(tc.err); got != tc.kind {
			t.Errorf("%s: Kind = %v,want %v", name, got, tc.kind)
		}
	}
}
