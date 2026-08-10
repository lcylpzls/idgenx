package idgenx

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/lcylpzls/errx"
)

// maxRandomBytes 是随机 ID 单次生成的最大字节数（防滥用/防 DoS）。
const maxRandomBytes = 4096

// RandomHex 生成 n 字节安全随机数的十六进制 ID（2n 位小写字符）。
// n 必须在 1..maxRandomBytes 范围内；随机源失败返回 ErrRandomFailure。
func RandomHex(n int) (string, error) {
	if n < 1 || n > maxRandomBytes {
		return "", errx.NewCodef(CodeInvalidConfig,
			"随机 ID 字节数必须在 1..%d 范围内，当前 %d", maxRandomBytes, n)
	}
	b, err := randRead(n)
	if err != nil {
		return "", errx.Wrap(err, errx.KindUnavailable, CodeRandomFailure, "随机 ID 生成失败")
	}
	return hex.EncodeToString(b), nil
}

// RandomBase64URL 生成 n 字节安全随机数的 base64url 无填充 ID。
// n 必须在 1..maxRandomBytes 范围内；随机源失败返回 ErrRandomFailure。
func RandomBase64URL(n int) (string, error) {
	if n < 1 || n > maxRandomBytes {
		return "", errx.NewCodef(CodeInvalidConfig,
			"随机 ID 字节数必须在 1..%d 范围内，当前 %d", maxRandomBytes, n)
	}
	b, err := randRead(n)
	if err != nil {
		return "", errx.Wrap(err, errx.KindUnavailable, CodeRandomFailure, "随机 ID 生成失败")
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
