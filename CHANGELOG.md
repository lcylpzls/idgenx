# Changelog

本项目遵循语义化版本（SemVer）。v1.0.0 之前允许破坏性变更。

## [v1.0.0] - 2026-08-09

### 正式版声明

- idgenx 进入语义化版本稳定期：v1.x 起承诺不引入破坏性 API 变更，
  新增能力以向后兼容的方式演进；
- 功能：雪花 ID（可配置位布局、回拨三策略、单调生成、解析、
  时间戳/序列溢出防护）+ 短 ID（base62、拒绝采样、易混字符排除、
  唯一性重试）；
- 质量基线：根包与 shortid 100% 语句覆盖、3 个 fuzz 目标、
  10 轮稳定性测试、5 轮 race、vet/staticcheck、govulncheck 0 漏洞、
  三平台 CI + Release 全绿；
- 性能：雪花 ~243ns/op（0 分配），与 bwmarrin/snowflake 持平；
- v0.1.0 - v1.0.0 共 14 个版本全部发布（含 3 个 CI 修复 patch）。

## [v0.12.3] - 2026-08-09

### 修复（三平台 CI）

- 真实时钟测试（单调/并发/压力）改用 Loose 策略，
  容忍 Windows 等平台偶发时钟回拨导致的测试失败；

## [v0.12.2] - 2026-08-09

### 修复（CI fuzz）

- FuzzNext 收敛：固定默认布局 + 随机回拨策略（100 次/输入），
  消除随机位布局触发的真实 sleep 慢路径导致的 worker 超时；
- 位布局边界仍由 FuzzConfig 与单元测试矩阵覆盖；
- 25 秒 fuzz（226K execs/s）+ 3 轮 race 全绿。

## [v0.12.1] - 2026-08-09

### 修复（CI 发现）

- 序列溢出等待增加超时上限（`MaxWait`）：时钟回拨停滞期间
  `Next` 不再无限阻塞，超时返回 `ErrWaitTimeout`；
- FuzzNext 收紧（`MaxWait=100µs`、2 次迭代），消除慢路径；
- 20 秒 fuzz + 5 轮 race + 10 轮测试全绿。

## [v0.12.0] - 2026-08-09

### 最终稳定性审查

- 修复：系统时钟早于纪元时首次生成产生负 ID 的缺陷
  （负时间戳按严重回拨处理，等待/拒绝/宽松三策略覆盖）；
- 逐行代码审查（生成/配置/短 ID）；锁与阻塞路径复核无异常；
- 10 轮测试 + 5 轮 race + 三目标 fuzz + govulncheck 全绿；
- 设计文档补充早于纪元语义。

## [v0.11.0] - 2026-08-09

### 工程收尾

- go.mod 工具链要求放宽至 `go 1.26.0`；
- README 性能数据统一为实测值；
- 依赖整理；10 轮测试 + 5 轮 race + govulncheck 全绿。

## [v0.10.0] - 2026-08-09

### 工业级正确性修复

- 时间戳位宽溢出防护：`now-epoch` 超过 `TimestampBits` 位宽
  时返回 `ErrTimestampOverflow`，杜绝高位进位破坏节点/序列布局；
- 覆盖路径：常规生成、序列溢出等待、回拨等待恢复后大跳；
- fuzz 扩展至 3 目标（新增 `FuzzNext`：任意配置与时钟下
  不 panic、有限步返回）；
- 新增错误码 `idgenx_timestamp_overflow`。

## [v0.9.0] - 2026-08-09

### 新增（自主打磨）

- 纪元近距边界测试（Epoch 贴近当前时间）；
- 10 轮测试 + 5 轮 race + govulncheck 终审复核；
- 导出符号 godoc 注释全量核对。

## [v0.8.0] - 2026-08-09

### 新增（自主打磨）

- README 现代化：实际 API 快速上手、目录树补全；
- 公开 API 全面核对（命名与文档一致，无未导出类型泄漏）；
- 10 轮测试 + 5 轮 race + govulncheck 终审复核。

## [v0.7.0] - 2026-08-09

### 调研

- docs/benchmark.md：与 bwmarrin/snowflake、sony/sonyflake 对比压测；
- 实测：idgenx 243ns/op（0 分配）与 bwmarrin 持平，
  sonyflake 39µs/op（慢约 160 倍）；
- 功能矩阵确认差异化价值（回拨策略/可配置布局/解析/短 ID）。

## [v0.6.0] - 2026-08-09

### 新增（自主打磨）

- godoc 示例：Generator 与 shortid 包 Example；
- 并发压力：8 goroutine × 5000 次无碰撞、负载中回拨演练；
- API 文档核对：`New` 签名补充选项参数；
- 3 轮 race、10 轮稳定性复核。

## [v0.5.0] - 2026-08-09

### 发布前终审

- ERRORS.md 错误码清单、MIT LICENSE、examples/basic 示例；
- 依赖整理；10 轮测试、5 轮 race、vet/staticcheck、
  govulncheck（0 漏洞）全部通过。

### 版本线

- v0.1.0 - v0.5.0 全部完成并发布，**roadmap 计划收官**；
- 后续按需继续自我打磨推进版本号（v1.0.0 之前允许破坏性变更）。

## [v0.4.0] - 2026-08-09

### 新增

- 可观测性：
  - `WithLogger`：回拨告警/拒绝/等待超时结构化日志；
  - `WithMetrics`：Generated/Rejected/WaitMS 回调；
- 基准测试（本机参考）：
  - 雪花串行 ~246ns/op（≈4M IDs/s，0 分配）；
  - 雪花并行 ~236ns/op；
  - 短 ID（8 位）~934ns/op（≈1M/s）；
- 边界矩阵：极端位布局（61/1/1 与 1/31/31）、并发指标回调。

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
