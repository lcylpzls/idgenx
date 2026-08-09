# idgenx API 定版

> 版本：v0.0.0（规划定稿） · 以下签名在实现阶段按本文执行；
> v0.1.0 前允许微调，v0.1.0 起冻结核心公开面。

## 1. 包结构

```go
idgenx           // 根包：雪花 ID 生成器、配置、错误
idgenx/shortid   // 短 ID 生成（base62）
```

## 2. 核心类型

### 2.1 BackwardStrategy

```go
type BackwardStrategy uint8

const (
	StrategyWait   BackwardStrategy = iota // 等待时钟追平（默认）
	StrategyReject                         // 拒绝并返回错误
	StrategyLoose                          // 沿用上一时间戳继续生成
)
```

### 2.2 Config

```go
type Config struct {
	Epoch         time.Time        // 默认 2024-01-01T00:00:00Z
	TimestampBits uint8            // 默认 41
	NodeBits      uint8            // 默认 10
	SequenceBits  uint8            // 默认 12
	NodeID        int64            // 默认 0
	Backward      BackwardStrategy // 默认 StrategyWait
	MaxWait       time.Duration    // 默认 5ms（仅 Wait 策略生效）
}

func DefaultConfig() Config
```

约束：`1+TimestampBits+NodeBits+SequenceBits == 64`；
`NodeID` 必须在 `[0, 2^NodeBits)`；`Epoch` 不得晚于当前时间；
`Backward` 必须合法；`MaxWait` 必须为正（零值填充默认）。

### 2.3 Parts

```go
type Parts struct {
	Timestamp time.Time
	NodeID    int64
	Sequence  int64
}
```

## 3. 雪花生成器

```go
func New(cfg Config, opts ...Option) (*Generator, error)
func (g *Generator) Next() (int64, error)
func (g *Generator) Parse(id int64) (Parts, error)
```

语义：

- `New`：配置校验失败返回 `ErrInvalidConfig`/`ErrNodeInvalid`；
- `Next`：返回严格递增的 64 位有符号 ID；时钟回拨按策略处理；
- `Parse`：解析时间戳/节点/序列；非法 ID（负值、位段越界）
  返回 `ErrInvalidID`；
- `Generator` 并发安全。

## 4. 选项（后续版本启用）

```go
func WithLogger(logger logx.Logger) Option    // v0.4.0
func WithMetrics(m Metrics) Option            // v0.4.0
func WithClock(now func() time.Time) Option   // v0.1.0
```

## 5. shortid 子包

```go
package shortid

const (
	AlphabetBase62        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	AlphabetNoConfusable  = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"
	MinLength = 6
	MaxLength = 64
)

func Generate(length int) (string, error)                       // base62
func GenerateWithAlphabet(alphabet string, length int) (string, error)
func GenerateUnique(length int, isUnique func(string) (bool, error)) (string, error) // v0.3.0
func IsValid(id string, alphabet string) bool                   // 长度与字符集校验
```

语义：

- 长度越界返回 `ErrInvalidConfig`；
- 字母表至少 16 字符且无重复，否则 `ErrInvalidConfig`；
- 随机源失败返回 `ErrRandomFailure`；
- `GenerateUnique`：调用 `isUnique` 校验，false 时重试（上限 8 次），
  仍冲突返回 `ErrCollision`；`isUnique` 返回错误时透传。

## 6. 错误值清单

```go
var (
	ErrInvalidConfig  = errx.NewCode(CodeInvalidConfig, "配置非法")
	ErrClockBackward  = errx.NewCode(CodeClockBackward, "检测到时钟回拨")
	ErrWaitTimeout    = errx.NewCode(CodeWaitTimeout, "等待时钟追平超时")
	ErrNodeInvalid    = errx.NewCode(CodeNodeInvalid, "节点 ID 越界")
	ErrInvalidID      = errx.NewCode(CodeInvalidID, "ID 解析失败")
	ErrRandomFailure  = errx.NewCode(CodeRandomFailure, "随机源失败")
	ErrCollision      = errx.NewCode(CodeCollision, "短 ID 碰撞重试耗尽")
	ErrTimestampOverflow = errx.NewCode(CodeTimestampOverflow, "时间戳超出位宽范围")
)
```

对应 `CodeXxx` 常量（`idgenx_*` 前缀），完整对照见
[design.md §7](design.md)。

## 7. 完整示例（规划）

```go
g, err := idgenx.New(idgenx.DefaultConfig())
if err != nil {
	panic(err)
}
id, err := g.Next()
parts, err := g.Parse(id)

code, err := shortid.Generate(8)
```
