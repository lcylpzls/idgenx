package shortid

import (
	"errors"
	"strings"
	"testing"

	"github.com/lcylpzls/idgenx"
)

// TestGenerate 覆盖默认字母表生成。
func TestGenerate(t *testing.T) {
	code, err := Generate(8)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 8 {
		t.Fatalf("长度不符：%q", code)
	}
	for _, r := range code {
		if !strings.ContainsRune(AlphabetBase62, r) {
			t.Fatalf("非法字符：%q", r)
		}
	}
	if !IsValid(code, AlphabetBase62) {
		t.Fatal("生成结果应合法")
	}
}

// TestGenerateWithAlphabet 覆盖自定义字母表与参数校验。
func TestGenerateWithAlphabet(t *testing.T) {
	code, err := GenerateWithAlphabet(AlphabetNoConfusable, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(code) != 10 {
		t.Fatalf("长度不符：%q", code)
	}
	if strings.ContainsAny(code, "0O1lI") {
		t.Fatalf("不应包含易混字符：%q", code)
	}
	for _, tc := range []struct {
		name     string
		alphabet string
		length   int
	}{
		{"过短长度", AlphabetBase62, MinLength - 1},
		{"过长长度", AlphabetBase62, MaxLength + 1},
		{"字母表过短", "abc", 8},
		{"字母表重复", "aabbccddeeffgghh", 8},
	} {
		if _, err := GenerateWithAlphabet(tc.alphabet, tc.length); !errors.Is(err, idgenx.ErrInvalidConfig) {
			t.Fatalf("%s 应报配置错误，实际：%v", tc.name, err)
		}
	}
}

// TestGenerateRandFailure 覆盖随机源失败。
func TestGenerateRandFailure(t *testing.T) {
	orig := randRead
	randRead = func(b []byte) (int, error) { return 0, errors.New("随机源故障") }
	defer func() { randRead = orig }()
	if _, err := Generate(8); !errors.Is(err, idgenx.ErrRandomFailure) {
		t.Fatalf("随机源失败应报错，实际：%v", err)
	}
}

// TestGenerateDeterministic 覆盖拒绝采样确定性（注入固定随机值）。
func TestGenerateDeterministic(t *testing.T) {
	orig := randRead
	randRead = func(b []byte) (int, error) {
		for i := range b {
			b[i] = byte(i % 10)
		}
		return len(b), nil
	}
	defer func() { randRead = orig }()
	code, err := Generate(6)
	if err != nil {
		t.Fatal(err)
	}
	// 注入值 0-9 均小于 max（248），首 6 个字符为字母表[0..5]。
	if code != AlphabetBase62[:6] {
		t.Fatalf("确定性输出不符：%q", code)
	}
}

// TestGenerateUnique 覆盖碰撞重试与透传。
func TestGenerateUnique(t *testing.T) {
	if _, err := GenerateUnique(8, nil); !errors.Is(err, idgenx.ErrInvalidConfig) {
		t.Fatalf("空回调应报错，实际：%v", err)
	}
	calls := 0
	code, err := GenerateUnique(8, func(s string) (bool, error) {
		calls++
		if calls == 1 {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || code == "" {
		t.Fatalf("应在第二次重试成功：calls=%d code=%q", calls, code)
	}
	if _, err := GenerateUnique(8, func(string) (bool, error) {
		return false, nil
	}); !errors.Is(err, idgenx.ErrCollision) {
		t.Fatalf("8 次冲突应报碰撞，实际：%v", err)
	}
	if _, err := GenerateUnique(8, func(string) (bool, error) {
		return false, errors.New("查询故障")
	}); err == nil || !strings.Contains(err.Error(), "查询故障") {
		t.Fatalf("回调错误应透传，实际：%v", err)
	}
	// 生成失败透传。
	orig := randRead
	randRead = func(b []byte) (int, error) { return 0, errors.New("随机源故障") }
	defer func() { randRead = orig }()
	if _, err := GenerateUnique(8, func(string) (bool, error) { return true, nil }); !errors.Is(err, idgenx.ErrRandomFailure) {
		t.Fatalf("生成失败应透传，实际：%v", err)
	}
}

// TestIsValid 覆盖校验分支。
func TestIsValid(t *testing.T) {
	if !IsValid("AbC123", AlphabetBase62) {
		t.Fatal("合法短 ID 应通过")
	}
	if IsValid("", AlphabetBase62) || IsValid("abc", AlphabetBase62) {
		t.Fatal("空/过短应失败")
	}
	if IsValid("abcXYZ", "") {
		t.Fatal("空字母表应失败")
	}
	if IsValid("abc$%^", AlphabetBase62) {
		t.Fatal("非法字符应失败")
	}
	if IsValid(strings.Repeat("a", MaxLength+1), AlphabetBase62) {
		t.Fatal("过长应失败")
	}
}
