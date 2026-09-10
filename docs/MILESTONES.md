# Gatus Fork 里程碑规划

> 仓库：`isleei/gatus`（fork from [TwiN/gatus](https://github.com/TwiN/gatus)）  
> 目标：自建服务探活，替代 Uptime（UptimeRobot 等 SaaS）；不选用 Uptime Kuma（性能不足）。  
> 原则：保留本 fork 定制能力（Admin v2、企业微信、证书监控、分组鉴权、防篡改、中文 i18n）；可持续合入上游。

---

## 背景与结论

| 项 | 说明 |
|---|---|
| 选型 | 继续使用本 fork，不迁移到 Uptime Kuma / SaaS |
| 优势 | Go 轻量、条件表达式强、企业微信、Admin UI、证书页、防篡改、已有 100+ 监控 + Postgres |
| 主要风险 | 自建单点自盲、公开仓密钥、合上游维护成本 |
| 上游同步 | 见 PR [#4](https://github.com/isleei/gatus/pull/4)（`merge/upstream-master-2026-09`） |

建议执行顺序：**阶段 0 → 1 → 2**，阶段 3 / 4 按痛点插队。不要先大改 Admin UI。

---

## 阶段 0 · 安全与基线（立刻）

**目标**：线上可安全运行、密钥不出仓、跑的是自己的镜像。

| # | 事项 | 验收标准 |
|---|---|---|
| 0.1 | 合并上游同步 PR [#4](https://github.com/isleei/gatus/pull/4) | `master` 包含上游修复（Postgres 索引、DNS TXT、cookie、JSONPath 等）与 fork 功能并存 |
| 0.2 | 轮换泄露凭据 | Postgres 密码、Basic Auth、企业微信 webhook 全部轮换；旧 webhook 作废 |
| 0.3 | 密钥移出公开仓 | `config.yaml`、`.gatus-managed-overlay.json`、文档示例仅用环境变量 / 占位符；生产 overlay 不进公开 git |
| 0.4 | 使用自建镜像 | 部署文档与流水线指向本仓库构建产物，**禁止**依赖 `ghcr.io/twin/gatus`（官方镜像无 Admin / WeCom 等定制） |

**相关文件（清理时重点）**：`config.yaml`、`.gatus-managed-overlay.json*`、`docs/DEPLOYMENT.md` 中的示例凭据。

---

## 阶段 1 · 防自盲（自建刚需）

**目标**：主实例或本机房挂掉时仍能收到告警。

| # | 事项 | 验收标准 |
|---|---|---|
| 1.1 | 第二节点（哨兵） | 另一台机器 / 另一机房部署轻量 Gatus（或等价探活） |
| 1.2 | 哨兵监控范围 | 至少监控：主 Gatus HTTP 健康、若干核心业务站点；告警走企业微信 |
| 1.3 | Postgres 备份 | 定期备份 + 恢复演练步骤写入 `docs/DEPLOYMENT.md` |

**说明**：哨兵不要求功能与主站对等；「能发现主站挂了」优先于功能完整。

---

## 阶段 2 · 可维护性

**目标**：合上游可预期，CI / 测试可信。

| # | 事项 | 验收标准 |
|---|---|---|
| 2.1 | GitHub `workflow` 权限 | `gh auth` 具备 `workflow` scope；补入此前因权限跳过的上游 `.github/workflows` 变更 |
| 2.2 | 合上游节奏 | 建议每季度一次；冲突高发区：Admin、`App.vue`、`go.mod`、`storage/store/sql/`、前端 `web/static`；操作手册 → [`UPSTREAM-MERGE.md`](./UPSTREAM-MERGE.md) |
| 2.3 | 测试修复 | 修复 API JSON 断言失败（时间戳 / `omitempty` / null vs `[]`）；ICMP 测试文档标明需 sudo |
| 2.4 | Overlay 策略 | 生产 overlay 私有化或本地挂载；公开仓不存真实客户域名与 webhook |

---

## 阶段 3 · 产品完善（相对 Uptime 体验）

**目标**：日常运营体验接近「能替 Uptime」，不牺牲 Gatus 的深度检查。

| # | 事项 | 说明 |
|---|---|---|
| 3.1 | Admin 增强 | 批量导入、告警渠道配置向导、证书页与告警阈值对齐 |
| 3.2 | Suite 级告警 | 上游仍弱，可按需自研 suite 失败统一通知 |
| 3.3 | 状态页 | 分组默认视图；对外只读状态页与 Admin 权限分离 |
| 3.4 | 可观测性 | 可选打开 Prometheus metrics，接入现有看板 |

---

## 阶段 4 · 可选增强

| # | 事项 | 说明 |
|---|---|---|
| 4.1 | 多地域探测 | 多台 Gatus + `external-endpoints` 汇总到主面板 |
| 4.2 | 审计保留 | Admin 审计日志保留 / 清理策略 |
| 4.3 | 身份认证 | OIDC 替代或补充纯 Basic Auth |

---

## 本 Fork 需长期保留的能力

合上游或重构时，以下能力不得丢失：

- Admin 控制台 v2（含审计日志、managed overlay、热重载）
- 企业微信（WeCom）告警
- 证书监控页 `/certificates`
- 分组鉴权 / groups API
- 端点防篡改（body-size drift / 关键词）
- 中文 i18n 与 `README_zh.md` / 部署文档
- Suite 级告警（`suites[].alerts`）

---

## 明确不做（当前）

- 迁移到 Uptime Kuma 或纯 SaaS UptimeRobot 作为主监控
- 用官方 `twin/gatus` 镜像替代本 fork 构建
- 在公开仓库继续存放生产 overlay / 明文密钥

---

## 进度跟踪

| 阶段 | 状态 | 备注 |
|---|---|---|
| 0 安全与基线 | 进行中 | 0.1/0.3/0.4 已完成；**0.2 仍待运维** → 见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §1 |
| 1 防自盲 | 进行中（仓内已完成） | 哨兵示例 + 备份演练文档已入库；**真实部署 / 演练仍属运维** → 见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §2–§3 |
| 2 可维护性 | 进行中 | API/ICMP 测试修复已入库；**2.1 workflow scope 仍需运维** → 见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §4.4 |
| 3 产品完善 | 仓内已完成 | 见细项；运维网络策略/真实 WeCom 仍属人工 |
| 4 可选增强 | 仓内已完成 | 多地域示例 + 审计保留 + OIDC 文档；部署/IdP 属运维 |

### 阶段 0 细项

| # | 状态 | 备注 |
|---|---|---|
| 0.1 合并上游 PR [#4](https://github.com/isleei/gatus/pull/4) | 已完成 | 已合入 `master`（含里程碑文档 PR [#5](https://github.com/isleei/gatus/pull/5)、密钥清理 PR [#6](https://github.com/isleei/gatus/pull/6)） |
| 0.2 轮换泄露凭据 | 待运维 | **无法在 git 完成**；检查清单见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §1（Postgres / Admin / WeCom；历史密钥仍须轮换） |
| 0.3 密钥移出公开仓 | 已完成 | `config.yaml` / overlay / 文档示例已占位；生产 overlay 已移出跟踪并加入 `.gitignore`；**git 历史仍含旧密钥，必须配合 0.2 轮换** |
| 0.4 使用自建镜像 | 已完成（文档） | `docs/DEPLOYMENT.md` 强调 `docker build -t gatus:local .`，禁止本 fork 生产依赖 `ghcr.io/twin/gatus` / `twinproduction/gatus` |

### 阶段 1 细项

| # | 状态 | 备注 |
|---|---|---|
| 1.1 第二节点（哨兵） | 仓内已完成 / 部署待运维 | 示例：[`docs/examples/sentinel/`](./examples/sentinel/)；**部署清单见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §2** |
| 1.2 哨兵监控范围 | 仓内已完成 | 主 `/health` + 2 个业务 URL 占位；WeCom；1m 间隔 |
| 1.3 Postgres 备份 | 仓内已完成 / 演练待运维 | `DEPLOYMENT.md` 含命令提纲；**演练清单见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §3** |

### 阶段 2 细项

| # | 状态 | 备注 |
|---|---|---|
| 2.1 GitHub `workflow` 权限 | 待运维 | 见 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §4.4：`gh auth refresh -s workflow` 后补入跳过的上游 workflows；**不阻塞其余仓内修复** |
| 2.2 合上游节奏 | 进行中 | 建议每季度；手册 [`UPSTREAM-MERGE.md`](./UPSTREAM-MERGE.md)；冲突高发区见上表 |
| 2.3 测试修复 | 已完成（仓内） | API `TestEndpointStatuses` / `TestSuiteStatuses`（时间戳先置零 + suite SQL 补 status/hostname + `results` 空数组）；ICMP `TestPing` 非 root 时 `t.Skip` |
| 2.4 Overlay 策略 | 已完成（与 0.3） | 公开仓仅 example；生产 overlay 本地/私有挂载 |

更新本表时请同步改「状态 / 备注」，并在相关 PR 描述中引用本文件对应章节。

---


### 阶段 3 细项

| # | 状态 | 备注 |
|---|---|---|
| 3.1 Admin 增强 | 已完成（仓内） | JSON/YAML 导入；WeCom/custom 表单校验与帮助文案；证书页可配置即将过期阈值（默认 72h，对齐 `[CERTIFICATE_EXPIRATION]`） |
| 3.2 Suite 级告警 | 已完成（仓内） | `suites[].alerts` + watchdog 触发/恢复；复用 wecom/slack/custom 等 provider；含单测；**重启/热重载后 suites[].alerts Triggered 状态恢复**（`initializeStorage` + `main_test.go`）；Admin 通知删除前计 suite/step 引用（`usedBySuites`，409） |
| 3.3 状态页 | 已完成（仓内） | 文档 `ui.default-sort-by: group`；公开看板 vs `/admin` / `/api/v1/admin/*` 分离说明 |
| 3.4 可观测性 | 已完成（仓内） | `metrics: true` + [`docs/examples/prometheus/`](./examples/prometheus/) |

### 阶段 4 细项

| # | 状态 | 备注 |
|---|---|---|
| 4.1 多地域探测 | 仓内已完成 / 部署待运维 | [`docs/examples/multi-region/`](./examples/multi-region/)（primary `external-endpoints` + curl 卫星推送 + 中文 README） |
| 4.2 审计保留 | 已完成（仓内） | `storage.admin-audit-max-age` + 日清理 ticker；`DELETE /api/v1/admin/audit-logs?days=` |
| 4.3 OIDC | 已完成（文档） | `docs/DEPLOYMENT.md` 中文 OIDC 专节（占位符）；IdP 注册与密钥轮换待运维 |

自审：[`docs/REVIEW-STAGE3-4.md`](./REVIEW-STAGE3-4.md)

## 相关链接

- 上游：https://github.com/TwiN/gatus
- 上游合并 PR：https://github.com/isleei/gatus/pull/4
- 部署文档：[`docs/DEPLOYMENT.md`](./DEPLOYMENT.md)
- 运维手册：[`docs/OPS-RUNBOOK.md`](./OPS-RUNBOOK.md)
