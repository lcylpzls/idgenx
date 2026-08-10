package core

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/lcylpzls/errx"
	"github.com/lcylpzls/testx"
)

func TestRandomHex(t *testing.T) {
	id, err := RandomHex(16)
	testx.RequireNoError(t, err)
	testx.RequireLen(t, id, 32)
	b, err := hex.DecodeString(id)
	testx.RequireNoError(t, err)
	testx.RequireLen(t, b, 16)
	if id != strings.ToLower(id) {
		t.Fatalf("RandomHex 应返回小写：%s", id)
	}
	other, err := RandomHex(16)
	testx.RequireNoError(t, err)
	if id == other {
		t.Fatalf("两次生成不应相同：%s", id)
	}
}

func TestRandomBase64URL(t *testing.T) {
	id, err := RandomBase64URL(32)
	testx.RequireNoError(t, err)
	if strings.ContainsAny(id, "=+/") {
		t.Fatalf("base64url 无填充应不含 =+/：%s", id)
	}
	b, err := base64.RawURLEncoding.DecodeString(id)
	testx.RequireNoError(t, err)
	testx.RequireLen(t, b, 32)
}

func TestRandomIDsInvalidLength(t *testing.T) {
	for _, n := range []int{0, -1, maxRandomBytes + 1} {
		_, err := RandomHex(n)
		testx.RequireErrCode(t, err, CodeInvalidConfig)
		_, err = RandomBase64URL(n)
		testx.RequireErrCode(t, err, CodeInvalidConfig)
	}
}

func TestRandomIDsRandFailure(t *testing.T) {
	orig := randRead
	randRead = func(n int) ([]byte, error) { return nil, errors.New("随机源故障") }
	defer func() { randRead = orig }()
	_, err := RandomHex(16)
	testx.RequireErrCode(t, err, CodeRandomFailure)
	_, err = RandomBase64URL(16)
	testx.RequireErrCode(t, err, CodeRandomFailure)
	if errx.KindOf(err) != errx.KindUnavailable {
		t.Fatalf("随机源失败分类应为 unavailable，当前 %s", errx.KindOf(err))
	}
}
