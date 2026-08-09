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

// registerCodes 在错误值初始化前完成注册，保证 NewCode 自动分类生效
// （包级变量初始化先于 init 执行，故不用 init 注册）。
var _ = registerCodes()

func registerCodes() bool {
	errx.RegisterCode(CodeInvalidConfig, "配置非法")
	errx.RegisterCodeKind(CodeInvalidConfig, errx.KindInvalid)
	errx.RegisterCode(CodeClockBackward, "检测到时钟回拨")
	errx.RegisterCodeKind(CodeClockBackward, errx.KindUnavailable)
	errx.RegisterCode(CodeWaitTimeout, "等待时钟追平超时")
	errx.RegisterCodeKind(CodeWaitTimeout, errx.KindTimeout)
	errx.RegisterCode(CodeNodeInvalid, "节点 ID 越界")
	errx.RegisterCodeKind(CodeNodeInvalid, errx.KindInvalid)
	errx.RegisterCode(CodeInvalidID, "ID 解析失败")
	errx.RegisterCodeKind(CodeInvalidID, errx.KindInvalid)
	errx.RegisterCode(CodeRandomFailure, "随机源失败")
	errx.RegisterCodeKind(CodeRandomFailure, errx.KindUnavailable)
	errx.RegisterCode(CodeCollision, "短 ID 碰撞重试耗尽")
	errx.RegisterCodeKind(CodeCollision, errx.KindConflict)
	errx.RegisterCode(CodeTimestampOverflow, "时间戳超出位宽范围")
	errx.RegisterCodeKind(CodeTimestampOverflow, errx.KindUnavailable)
	return true
}

// 预定义错误值，可用 errx.Is / errors.Is 判断。
var (
	// ErrInvalidConfig 配置非法。
	ErrInvalidConfig = errx.NewCode(CodeInvalidConfig, "配置非法")
	// ErrClockBackward 检测到时钟回拨。
	ErrClockBackward = errx.NewCode(CodeClockBackward, "检测到时钟回拨")
	// ErrWaitTimeout 等待时钟追平超时。
	ErrWaitTimeout = errx.NewCode(CodeWaitTimeout, "等待时钟追平超时")
	// ErrNodeInvalid 节点 ID 越界。
	ErrNodeInvalid = errx.NewCode(CodeNodeInvalid, "节点 ID 越界")
	// ErrInvalidID ID 解析失败。
	ErrInvalidID = errx.NewCode(CodeInvalidID, "ID 解析失败")
	// ErrRandomFailure 随机源失败。
	ErrRandomFailure = errx.NewCode(CodeRandomFailure, "随机源失败")
	// ErrCollision 短 ID 碰撞重试耗尽。
	ErrCollision = errx.NewCode(CodeCollision, "短 ID 碰撞重试耗尽")
	// ErrTimestampOverflow 时间戳超出位宽范围。
	ErrTimestampOverflow = errx.NewCode(CodeTimestampOverflow, "时间戳超出位宽范围")
)
