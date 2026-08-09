// Package idgenx 提供分布式 ID 生成：雪花 ID（Snowflake）与短 ID，
// 与 errx / logx 生态打通。
//
// 雪花 ID 为 64 位有符号整数：时间戳（毫秒，相对纪元）+ 节点 + 序列，
// 同节点单调递增、并发安全；时钟回拨可等待/拒绝/宽松处理。
// 短 ID 提供 base62 随机码（邀请码、短号），易混字符可选排除。
package idgenx
