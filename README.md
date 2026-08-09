# idgenx

自研分布式 ID 生成库：雪花 ID（Snowflake）与短 ID，
与 errx / logx 生态打通。

> 当前状态：**v0.12.1 已发布**。CI 发现的溢出等待缺陷已修复，
> 达到工业级生产稳定候选（等待用户决定 v1）。

性能参考（本机）：雪花 ~243ns/op（0 分配）、短 ID ~934ns/op。

完整错误码清单见 [ERRORS.md](ERRORS.md)。

## 定位

idgenx **不依赖外部协调服务**，解决自用项目中每个业务都要重复的部分：

- 雪花 ID：64 位（时间戳 + 节点 + 序列），单调递增、并发安全；
- 位布局可配置：时间/节点/序列位数按业务规模调整；
- 时钟回拨防护：等待、拒绝、宽松三种策略；
- 短 ID：base62 随机码（邀请码、订单短号），易混字符可选排除；
- 可观测性：logx 结构化日志、外部注入 Metrics；
- 错误语义：统一 errx 错误码。

所有生成器并发安全，可在多个 goroutine 间共享。

## 目录

```
idgenx/
├── CHANGELOG.md          # 变更记录
├── ERRORS.md             # 错误码清单
├── LICENSE               # MIT 许可
├── docs/
│   ├── README.md          # 文档索引
│   ├── design.md          # 设计定版（位布局/回拨策略/并发模型/错误码）
│   ├── api.md             # API 定版（完整签名与语义）
│   ├── research.md        # 领域调研与设计取舍
│   ├── roadmap.md         # 版本路线
│   └── benchmark.md       # 竞品对比与压测报告
├── examples/basic/        # 基础示例
├── shortid/               # 短 ID 子包
└── README.md
```

## 快速上手

```go
import (
	"github.com/lcylpzls/idgenx"
	"github.com/lcylpzls/idgenx/shortid"
)

// 雪花 ID：默认布局 41/10/12，严格递增。
g, err := idgenx.New(idgenx.DefaultConfig())
id, err := g.Next()
parts, err := g.Parse(id) // 时间戳/节点/序列

// 短 ID：8 位 base62 随机码。
code, err := shortid.Generate(8)

// 唯一短 ID：业务回调校验（查库/查缓存）。
code, err = shortid.GenerateUnique(8, func(s string) (bool, error) {
	return true, nil
})
```

配置（位布局/纪元/节点/回拨策略）见 [docs/api.md](docs/api.md)。

## License

MIT © [lcylpzls](https://github.com/lcylpzls)
