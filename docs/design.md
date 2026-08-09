# idgenx 设计定版

> 版本：v0.0.0（规划定稿） · 状态：文档已定版，代码未开始

## 1. 定位

idgenx 是**分布式环境下的 ID 生成库**：雪花 ID 保证全局有序、单调、
高吞吐；短 ID 提供可读、紧凑的随机码。两者均不依赖外部协调服务，
节点身份由部署配置决定。

## 2. 范围边界（明确不做）

| 不做 | 原因与替代 |
| --- | --- |
| 外部协调分配节点 ID | 需要 zookeeper/etcd 时由业务或适配器实现；idgenx 提供节点配置与校验 |
| 跨进程唯一性自动保证 | 唯一性依赖节点 ID 唯一配置；文档明确部署要求 |
| 序列持久化 | 重启后序列随机起点，时间戳+节点保证无碰撞 |
| UUID 生成 | 标准库已有；idgenx 聚焦有序 ID 与短码 |
| 数据库自增主键替代 | 分库分表场景业务自行选择 |

## 3. 雪花 ID 位布局

64 位有符号整数（Go int64，最高位恒为 0）：

```
┌──────────────┬──────────┬──────────────┐
│ sign 1       │ timestamp│ node │ seq   │
│ (0)          │ N 位     │ M 位 │ K 位  │
└──────────────┴──────────┴──────────────┘
```

约束：`1 + N + M + K = 64`。

默认布局（`DefaultConfig`）：

- 时间戳 41 位：相对纪元 `2024-01-01T00:00:00Z` 的毫秒数，
  可用约 69 年（至 2093 年）；
- 节点 10 位：0-1023；
- 序列 12 位：同毫秒 0-4095。

`Config` 支持自定义 `Epoch`、`TimestampBits`、`NodeBits`、
`SequenceBits` 与 `NodeID`；位数越界、总和不为 63、节点越界、
纪元晚于当前时间均在构造时返回 `ErrInvalidConfig`。

## 4. 单调性与序列

- 同节点同毫秒：序列 +1，溢出后等待下一毫秒（忙等或短睡眠）；
- 跨毫秒：序列归零（从随机起点开始，避免重启后从 0 引发可预测性）；
- 同节点产出 ID 严格递增（时间或序列至少一个前进）；
- 并发模型：单互斥锁串行化生成（正确性优先，目标 ≥ 2M IDs/s）；
- 时钟源可注入（测试与回拨演练）。

## 5. 时钟回拨策略

检测：`now < lastTimestamp` 即视为回拨。三种策略：

| 策略 | 行为 | 适用 |
| --- | --- | --- |
| `StrategyWait`（默认） | 等待时钟追平，`MaxWait` 内成功，超时返回 `ErrWaitTimeout` | 生产默认，容忍毫秒级回拨 |
| `StrategyReject` | 立即返回 `ErrClockBackward` | 严格场景，回拨即不可用 |
| `StrategyLoose` | 沿用 `lastTimestamp` 并递增序列，不报错 | 测试/容忍短暂回拨，存在极小重复风险 |

`MaxWait` 默认 5ms，必须为正。

## 6. 短 ID

- 字母表：默认 base62（`0-9A-Za-z`）；
- 易混字符排除：可选 `AlphabetNoConfusable`（去掉 `0O1lI` 等）；
- 随机源：`crypto/rand`，失败返回 `ErrRandomFailure`；
- 长度校验：`ShortIDMinLength`（默认 6）至 `ShortIDMaxLength`（默认 64）；
- 碰撞处理：`GenerateUnique` 内置唯一性回调（业务查库/查缓存），
  返回 false 时重试，上限 8 次，仍冲突返回 `ErrCollision`。

## 7. 错误码（errx）

| 错误码 | 含义 | Kind | 建议 HTTP |
| --- | --- | --- | --- |
| `idgenx_invalid_config` | 配置非法 | invalid_argument | 400 |
| `idgenx_clock_backward` | 时钟回拨（Reject 策略） | unavailable | 503 |
| `idgenx_wait_timeout` | 回拨等待超时 | timeout | 504 |
| `idgenx_node_invalid` | 节点 ID 越界 | invalid_argument | 400 |
| `idgenx_invalid_id` | ID 解析失败 | invalid_argument | 400 |
| `idgenx_rand_failure` | 随机源失败 | unavailable | 503 |
| `idgenx_collision` | 短 ID 碰撞重试耗尽 | conflict | 409 |

预定义错误值（`ErrInvalidConfig` 等）支持 `errors.Is`；`errx.Is` 按码判断。

## 8. 可观测性

### 8.1 日志（logx）

结构化字段统一前缀 `idgenx_`：

| 事件 | 字段 |
| --- | --- |
| 生成失败 | `idgenx_node`、`error` |
| 回拨告警 | `idgenx_node`、`idgenx_backward_ms`、`idgenx_strategy` |
| 等待超时 | `idgenx_node`、`idgenx_max_wait` |

> 短 ID 碰撞重试不内建日志（薄库原则），结果由调用方感知
> （`ErrCollision` 或唯一性回调）。

### 8.2 Metrics（外部注入）

```go
type Metrics struct {
	Generated func(node int64, delta int)   // 生成计数
	Rejected  func(node int64, err error)   // 拒绝计数
	WaitMS    func(node int64, ms int64)    // 回拨等待耗时
}
```

全部回调可选（nil 跳过），Prometheus 适配由业务或 metricsx 完成。

## 9. 安全与健壮性

- `crypto/rand` 随机源，失败返回错误而非弱回退；
- 所有配置有上限校验，构造期失败（不 panic）；
- 时钟回拨不静默降级（除显式 Loose 策略）；
- 生成路径无阻塞外部调用；锁临界区仅整数运算；
- 时钟与随机源可注入（测试与故障演练）。

## 10. 性能目标

- 雪花生成：单机 ≥ 2M IDs/s（默认布局，8 核参考）；
- 短 ID 生成：≥ 500K/s（长度 8）；
- 基准测试：串行/并发生成、回拨等待、短 ID 碰撞重试。

> 目标基线以 v0.1.0 本机实测校准后写入 README。

## 11. 测试与质量门禁

- 核心包 100% 语句覆盖（含 shortid 子包）；
- `-race` 全绿、连续多轮无偶发竞态；
- fuzz：配置边界、短 ID 参数；
- 三平台 CI + govulncheck + Release 自动发布；
- 边界矩阵：位边界（0/63）、序列溢出、回拨各策略、纪元边界。

## 12. 依赖

```go
require (
	github.com/lcylpzls/errx v1.2.0
	github.com/lcylpzls/logx v1.0.0
)
```

除自身生态外 **0 第三方依赖**。

## 13. 开放问题（定稿结论）

| 问题 | 结论 |
| --- | --- |
| 默认纪元 | 2024-01-01T00:00:00Z |
| 回拨默认策略 | StrategyWait，MaxWait 5ms |
| 序列重启起点 | 随机起点（0-4095），防可预测性 |
| 短 ID 默认长度 | 8（base62，含易混字符） |
| 短 ID 碰撞处理 | 业务回调 + 重试上限，库不查库 |
