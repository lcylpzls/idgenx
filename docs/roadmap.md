# idgenx 版本路线

> 目标：v0.1.0 起每版完成即全自动 CI + Release，全部通过后进入下一版；
> 定版级质量贯穿全程（100% 覆盖、race、fuzz、三平台 CI、govulncheck）。

## v0.1.0 — 雪花核心（已发布）

- 默认位布局（41/10/12）+ `DefaultConfig`；
- `New`/`Next`/`Parse`；单调递增、序列溢出等待；
- `WithClock` 时钟注入；并发安全；
- errx 错误码全集（按能力分版启用）。

## v0.2.0 — 可配置与回拨策略（已发布）

- 位布局/纪元/节点配置校验；
- 时钟回拨三策略（Wait/Reject/Loose）与 `MaxWait`；
- 序列随机起点；节点 ID 校验。

## v0.3.0 — 短 ID

- `shortid` 子包：`Generate`/`GenerateWithAlphabet`/`IsValid`；
- 易混字符字母表、长度校验、碰撞重试回调；
- fuzz：配置边界、短 ID 参数。

## v0.4.0 — 可观测与性能

- `WithLogger`/`WithMetrics`；回拨告警、失败、碰撞日志；
- 基准测试：串行/并发生成、回拨等待、短 ID；
- 边界矩阵：位边界、序列溢出、回拨各策略、纪元边界。

## v0.5.0 — 发布前终审

- 依赖整理、govulncheck、静态检查全量；
- README / ERRORS.md / LICENSE / 示例定稿；
- 收口于 v0.5.0（roadmap 完成）。

## 质量门禁（每版）

```powershell
go test -count=1 ./...
go test -count=1 -coverprofile=coverage.out ./...   # 核心包 100%
go test -race -count=1 ./...
go vet ./... && staticcheck ./...
go test -run '^$' -fuzz '^FuzzConfig$' -fuzztime=10s .
govulncheck ./...
```

CI：ubuntu/windows/macos 三平台 + fuzz job + govulncheck job +
Release（tag 触发）。
