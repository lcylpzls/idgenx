package idgenx

import "github.com/lcylpzls/errx"

// 错误码统一以 idgenx_ 为前缀。
const (
	// CodeInvalidConfig 配置非法。
	CodeInvalidConfig errx.Code = "idgenx_invalid_config"
	// CodeClockBackward 检测到时钟回拨（Reject 策略）。
	CodeClockBackward errx.Code = "idgenx_clock_backward"
	// CodeWaitTimeout 等待时钟追平超时。
	CodeWaitTimeout errx.Code = "idgenx_wait_timeout"
	// CodeNodeInvalid 节点 ID 越界。
	CodeNodeInvalid errx.Code = "idgenx_node_invalid"
	// CodeInvalidID ID 解析失败。
	CodeInvalidID errx.Code = "idgenx_invalid_id"
	// CodeRandomFailure 随机源失败。
	CodeRandomFailure errx.Code = "idgenx_rand_failure"
	// CodeCollision 短 ID 碰撞重试耗尽。
	CodeCollision errx.Code = "idgenx_collision"
	// CodeTimestampOverflow 时间戳超出位宽范围（生成不可用）。
	CodeTimestampOverflow errx.Code = "idgenx_timestamp_overflow"
)

// 预定义错误值，可用 errx.Is / errors.Is 判断。
var (
	// ErrInvalidConfig 配置非法。
	ErrInvalidConfig = errx.New(errx.KindInvalid, CodeInvalidConfig, "配置非法")
	// ErrClockBackward 检测到时钟回拨。
	ErrClockBackward = errx.New(errx.KindUnavailable, CodeClockBackward, "检测到时钟回拨")
	// ErrWaitTimeout 等待时钟追平超时。
	ErrWaitTimeout = errx.New(errx.KindTimeout, CodeWaitTimeout, "等待时钟追平超时")
	// ErrNodeInvalid 节点 ID 越界。
	ErrNodeInvalid = errx.New(errx.KindInvalid, CodeNodeInvalid, "节点 ID 越界")
	// ErrInvalidID ID 解析失败。
	ErrInvalidID = errx.New(errx.KindInvalid, CodeInvalidID, "ID 解析失败")
	// ErrRandomFailure 随机源失败。
	ErrRandomFailure = errx.New(errx.KindUnavailable, CodeRandomFailure, "随机源失败")
	// ErrCollision 短 ID 碰撞重试耗尽。
	ErrCollision = errx.New(errx.KindConflict, CodeCollision, "短 ID 碰撞重试耗尽")
	// ErrTimestampOverflow 时间戳超出位宽范围。
	ErrTimestampOverflow = errx.New(errx.KindUnavailable, CodeTimestampOverflow, "时间戳超出位宽范围")
)
