# 错误码清单

所有错误统一使用 errx 语义（`errx.Is(err, Code)` 判断）。

| 错误码 | 含义 | Kind | 建议 HTTP |
| --- | --- | --- | --- |
| `idgenx_invalid_config` | 配置非法 | invalid_argument | 400 |
| `idgenx_clock_backward` | 检测到时钟回拨（Reject 策略） | unavailable | 503 |
| `idgenx_wait_timeout` | 等待时钟追平超时 | timeout | 504 |
| `idgenx_node_invalid` | 节点 ID 越界 | invalid_argument | 400 |
| `idgenx_invalid_id` | ID 解析失败 | invalid_argument | 400 |
| `idgenx_rand_failure` | 随机源失败 | unavailable | 503 |
| `idgenx_collision` | 短 ID 碰撞重试耗尽 | conflict | 409 |
| `idgenx_timestamp_overflow` | 时间戳超出位宽范围（生成不可用） | unavailable | 503 |
