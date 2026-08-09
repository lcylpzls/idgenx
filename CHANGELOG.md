# Changelog

本项目遵循语义化版本（SemVer）。v1.0.0 之前允许破坏性变更。

## [v0.3.0] - 2026-08-09

### 新增

- shortid 子包：
  - `Generate`/`GenerateWithAlphabet`/`GenerateUnique`/`IsValid`；
  - `AlphabetBase62` 与 `AlphabetNoConfusable`（排除易混字符）；
  - 拒绝采样消除模偏差（均匀分布）；
  - `GenerateUnique` 碰撞重试（上限 8 次，回调可透传错误）；
  - fuzz 目标 `FuzzShortID`。

## [v0.2.0] - 2026-08-09

### 新增

- 可配置位布局与回拨策略：
  - `Config.Backward`：`StrategyWait`（默认）/`StrategyReject`/
    `StrategyLoose`；
  - `Config.MaxWait` 自定义回拨等待上限（默认 5ms）；
  - 序列随机起点（`crypto/rand`，首次生成保留随机序列，
    防重启后可预测性）；随机源失败返回 `ErrRandomFailure`；
- 修复：位偏移计算（nodeShift = 节点位 + 序列位）、
  Loose 策略沿用上一时间戳不再倒退。

## [v0.1.0] - 2026-08-09

### 新增

- 雪花 ID 核心：
  - 默认布局 41/10/12（纪元 2024-01-01T00:00:00Z）；
  - `New`/`Next`/`Parse`，严格递增、并发安全；
  - 同毫秒序列溢出自动等待下一毫秒；
  - 时钟回拨固定等待（5ms 上限，超时 `ErrWaitTimeout`）；
  - `WithClock` 时钟注入；零值字段自动填充默认；
- errx 错误码骨架（`idgenx_*` 全集，按能力分版启用）；
- fuzz 目标 `FuzzConfig`；CI：三平台 + race + fuzz + govulncheck + Release。
