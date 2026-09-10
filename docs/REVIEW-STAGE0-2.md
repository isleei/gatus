# Self-review · Stage 0–2（哨兵 + 可维护性）

> 对应分支 / PR：`chore/stage1-2-sentinel-and-hardening`  
> 日期：2026-09-10（Asia/Shanghai）

## Checklist

| 项 | 结果 |
|---|---|
| 新文件无真实密钥 / DSN / bcrypt / webhook | ✅ 哨兵与文档仅 `$GATUS_*` / `example.com` 占位 |
| 哨兵示例仅占位符 | ✅ `docs/examples/sentinel/{config,compose,README}` |
| `go build -o /tmp/gatus .` | ✅ 通过 |
| `go test ./api/ -run 'TestEndpointStatuses$|TestSuiteStatuses$'` | ✅ PASS |
| `go test ./api/`（整包） | ✅ PASS |
| `go test ./client/ -run 'TestPing|TestShouldRunPinger'` | ✅ `TestPing` 非 root → SKIP；特权探测 PASS |
| `go test ./storage/...` | ✅ PASS |
| 未 force-push / 未改写已推送历史 | ✅ |

## 本 PR 改了什么

1. **Stage 1 仓内**：哨兵示例目录；`DEPLOYMENT.md`「防自盲 / 哨兵」+ Postgres 备份演练提纲；`MILESTONES.md` 进度更新。
2. **Stage 2 仓内**：API 状态测试在 insert 前清零时间戳；SQL `getSuiteResults` 补回 `status`/`hostname` 且空结果返回 `[]`；ICMP `TestPing` 非 root `t.Skip`；里程碑注明 workflow scope 不阻塞。

## 仍需人类 / 运维完成

| 项 | 说明 |
|---|---|
| **0.2 凭据轮换** | 线上轮换 Postgres、Basic Auth、WeCom；作废旧 webhook（git 历史仍可能含旧值） |
| **1.x 真正部署哨兵** | 在另一主机/网络 `compose up`，填真实 env |
| **1.3 备份演练** | 按 `DEPLOYMENT.md` 跑一遍 pg_dump / pg_restore |
| **2.1 workflow scope** | `gh auth refresh -s workflow` 后补上游 workflows（本 PR 不阻塞） |

## 已知非阻塞差距

- 哨兵 README / compose 假设操作者会自行注入 env；未提供真实 `.env`。
- ICMP 在非特权环境跳过属预期；有 `CAP_NET_RAW`/root 时应再跑通。
