# 哨兵 Gatus（防自盲）示例

> 运维部署检查清单（异主机、构建本 fork、验 `/health`、告警演练）：[`docs/OPS-RUNBOOK.md`](../../OPS-RUNBOOK.md) §2。

本目录提供一套**轻量第二节点**配置：在**另一台主机 / 另一网络**上部署 Gatus，专门监视主实例与少量核心业务探活。它**不能替代**主 Gatus（无完整监控清单、无生产 Admin 工作流），只解决「主站挂了没人知道」的自盲问题。

## 监视什么

| 端点 | 环境变量 | 说明 |
|------|----------|------|
| 主 Gatus 健康 | `GATUS_PRIMARY_URL`（如 `https://status.example.com`） | 请求 `${GATUS_PRIMARY_URL}/health` |
| 核心业务 1 | `GATUS_SENTINEL_CRITICAL_URL_1` | 你最关心的业务探活 URL |
| 核心业务 2 | `GATUS_SENTINEL_CRITICAL_URL_2` | 可选第二业务 URL |

告警走企业微信：`GATUS_WECOM_WEBHOOK_URL`。探测间隔默认 **1 分钟**。存储默认 **memory**（可改 sqlite，见 `config.yaml` 注释）。

## 为何要另一台机器

若哨兵与主实例同机同网，机房断电 / 宿主机宕机会一起挂掉，仍会自盲。请放到：

- 不同宿主机，或
- 不同可用区 / 机房，或至少
- 不同 VPC / 出口网络

## 快速启动

在本 fork 仓库中构建镜像（**禁止**用官方 `twinproduction/gatus` / `ghcr.io/twin/gatus` 作为本 fork 生产镜像）：

```bash
# 仓库根目录
docker build -t gatus:local .

cd docs/examples/sentinel
export GATUS_PRIMARY_URL='https://status.example.com'
export GATUS_WECOM_WEBHOOK_URL='https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=REPLACE'
export GATUS_SENTINEL_CRITICAL_URL_1='https://app.example.com/health'
export GATUS_SENTINEL_CRITICAL_URL_2='https://api.example.com/health'
docker compose up -d --build
```

`compose.yaml` 的 `build.context` 指向仓库根（`../../..`），也可先构建 `gatus:local` 再只用 `image`。

默认映射宿主机 `8081 → 8080`，避免与主实例端口冲突。

## 环境变量一览

| 变量 | 必填 | 说明 |
|------|------|------|
| `GATUS_PRIMARY_URL` | 是 | 主 Gatus 对外基址（无尾斜杠），哨兵会请求 `/health` |
| `GATUS_WECOM_WEBHOOK_URL` | 是 | 企业微信机器人 webhook（勿提交真实值） |
| `GATUS_SENTINEL_CRITICAL_URL_1` | 是 | 核心业务探活 URL |
| `GATUS_SENTINEL_CRITICAL_URL_2` | 建议 | 第二核心业务探活 URL |

仓库内文件**仅含占位符**，不含真实密钥或客户域名。

## 与主实例的关系

- 哨兵：少端点、短间隔、外网可达视角、独立告警通道（可与主站共用 WeCom，但建议单独机器人便于区分）。
- 主实例：完整端点清单、Postgres、Admin、证书页等。
- **运维仍须真正部署哨兵**；本目录只提供 in-repo 交付物与文档。

更多背景见 [`docs/DEPLOYMENT.md`](../../DEPLOYMENT.md)「防自盲 / 哨兵」、[`docs/MILESTONES.md`](../../MILESTONES.md) 阶段 1，以及可勾选清单 [`docs/OPS-RUNBOOK.md`](../../OPS-RUNBOOK.md) §2。
