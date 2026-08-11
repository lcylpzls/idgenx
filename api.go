package idgenx

import (
	"time"

	"github.com/lcylpzls/idgenx/internal/core"
	"github.com/lcylpzls/logx"
)

// Version 是当前库版本，与 git tag 保持一致。
const Version = core.Version

type (
	Metrics          = core.Metrics
	Parts            = core.Parts
	Option           = core.Option
	Generator        = core.Generator
	BackwardStrategy = core.BackwardStrategy
	Config           = core.Config
)

const (
	CodeInvalidConfig     = core.CodeInvalidConfig
	CodeClockBackward     = core.CodeClockBackward
	CodeWaitTimeout       = core.CodeWaitTimeout
	CodeNodeInvalid       = core.CodeNodeInvalid
	CodeInvalidID         = core.CodeInvalidID
	CodeRandomFailure     = core.CodeRandomFailure
	CodeCollision         = core.CodeCollision
	CodeTimestampOverflow = core.CodeTimestampOverflow
)

var (
	ErrInvalidConfig     = core.ErrInvalidConfig
	ErrClockBackward     = core.ErrClockBackward
	ErrWaitTimeout       = core.ErrWaitTimeout
	ErrNodeInvalid       = core.ErrNodeInvalid
	ErrInvalidID         = core.ErrInvalidID
	ErrRandomFailure     = core.ErrRandomFailure
	ErrCollision         = core.ErrCollision
	ErrTimestampOverflow = core.ErrTimestampOverflow
)

func New(cfg Config, opts ...Option) (*Generator, error) {
	return core.New(cfg, opts...)
}
func DefaultConfig() Config                 { return core.DefaultConfig() }
func WithClock(now func() time.Time) Option { return core.WithClock(now) }
func WithLogger(logger logx.Logger) Option  { return core.WithLogger(logger) }
func WithMetrics(m Metrics) Option          { return core.WithMetrics(m) }
func RandomHex(n int) (string, error)       { return core.RandomHex(n) }
func RandomBase64URL(n int) (string, error) { return core.RandomBase64URL(n) }
