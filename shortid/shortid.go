// Package shortid 提供短 ID 生成：base62 随机码（邀请码、短号）。
package shortid

import (
	"strconv"

	"github.com/lcylpzls/cryptox"
	"github.com/lcylpzls/idgenx"
	"github.com/lcylpzls/validx"
)

const (
	// AlphabetBase62 默认字母表（0-9A-Za-z）。
	AlphabetBase62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// AlphabetNoConfusable 排除易混字符（0O1lI 等）的字母表。
	AlphabetNoConfusable = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"
	// MinLength 最短长度。
	MinLength = 6
	// MaxLength 最长长度。
	MaxLength = 64

	minAlphabetLength = 16
	maxCollisionRetry = 8
)

// randRead 可替换的随机源，便于测试注入失败场景。
var randRead = cryptox.RandomBytes

// Generate 生成 base62 短 ID。
func Generate(length int) (string, error) {
	return GenerateWithAlphabet(AlphabetBase62, length)
}

// GenerateWithAlphabet 按指定字母表生成短 ID（拒绝采样，均匀分布）。
func GenerateWithAlphabet(alphabet string, length int) (string, error) {
	if err := validateParams(alphabet, length); err != nil {
		return "", err
	}
	out := make([]byte, length)
	max := 256 - (256 % len(alphabet))
	for i := 0; i < length; {
		buf, err := randRead(length * 2)
		if err != nil {
			return "", idgenx.ErrRandomFailure
		}
		for _, v := range buf {
			if int(v) >= max {
				continue // 拒绝采样，消除模偏差。
			}
			out[i] = alphabet[int(v)%len(alphabet)]
			i++
			if i == length {
				break
			}
		}
	}
	return string(out), nil
}

// GenerateUnique 生成唯一短 ID：调用 isUnique 校验，冲突重试（上限 8 次）。
func GenerateUnique(length int, isUnique func(string) (bool, error)) (string, error) {
	if isUnique == nil {
		return "", idgenx.ErrInvalidConfig
	}
	for i := 0; i < maxCollisionRetry; i++ {
		code, err := Generate(length)
		if err != nil {
			return "", err
		}
		ok, err := isUnique(code)
		if err != nil {
			return "", err
		}
		if ok {
			return code, nil
		}
	}
	return "", idgenx.ErrCollision
}

// init 注册短 ID 规则到 validx 全局规则表，错误码保持 idgenx 语义。
func init() {
	_ = validx.RegisterRule("idgenx_shortid_params", func(value any, param, path string) error {
		// 内部调用保证 value 为字母表、param 为长度。
		alphabet := value.(string)
		length, _ := strconv.Atoi(param)
		if length < MinLength || length > MaxLength {
			return idgenx.ErrInvalidConfig
		}
		if len(alphabet) < minAlphabetLength {
			return idgenx.ErrInvalidConfig
		}
		seen := make(map[rune]struct{}, len(alphabet))
		for _, r := range alphabet {
			if _, ok := seen[r]; ok {
				return idgenx.ErrInvalidConfig
			}
			seen[r] = struct{}{}
		}
		return nil
	})
	_ = validx.RegisterRule("idgenx_shortid_valid", func(value any, param, path string) error {
		// 内部调用保证 value 为 ID、param 为字母表。
		id := value.(string)
		alphabet := param
		if id == "" || alphabet == "" || len(id) < MinLength || len(id) > MaxLength {
			return idgenx.ErrInvalidConfig
		}
		allowed := make(map[rune]struct{}, len(alphabet))
		for _, r := range alphabet {
			allowed[r] = struct{}{}
		}
		for _, r := range id {
			if _, ok := allowed[r]; !ok {
				return idgenx.ErrInvalidConfig
			}
		}
		return nil
	})
}

// IsValid 校验短 ID 长度与字符集（统一走 validx 规则）。
func IsValid(id, alphabet string) bool {
	return validx.ValidateField(id, "idgenx_shortid_valid="+alphabet) == nil
}

// validateParams 校验字母表与长度（统一走 validx 规则）。
func validateParams(alphabet string, length int) error {
	return validx.ValidateField(alphabet, "idgenx_shortid_params="+strconv.Itoa(length))
}
