package shortid

import (
	"testing"
)

// FuzzShortID 模糊测试短 ID 生成与校验，确保任意输入不 panic。
func FuzzShortID(f *testing.F) {
	f.Add(8, AlphabetBase62)
	f.Add(0, "")
	f.Add(100, AlphabetNoConfusable)
	f.Fuzz(func(t *testing.T, length int, alphabet string) {
		if len(alphabet) > 256 || length > 128 {
			t.Skip("输入过大")
		}
		code, err := GenerateWithAlphabet(alphabet, length)
		if err != nil {
			return
		}
		_ = IsValid(code, alphabet)
	})
}
