# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Gatus 是一个面向开发者的健康状态监控面板，支持使用 HTTP、ICMP、TCP、DNS、gRPC、WebSocket、SSH、SCTP、UDP、STARTTLS、TLS 等多种协议监控服务健康状态。通过条件表达式（如 `[STATUS] == 200`、`[RESPONSE_TIME] < 150`）评估健康状况，并支持 38 种告警渠道发送通知。

- 模块路径：`github.com/TwiN/gatus/v5`
- Go 版本：1.25.5
- HTTP 框架：Fiber v2
- 前端：Vue 3 + Tailwind CSS（通过 `//go:embed` 嵌入）
- Docker 构建：多阶段（golang:alpine → scratch），CGO_ENABLED=0

## 常用命令

```bash
# 构建
go build -v -o gatus .

# 运行（开发模式，启用 CORS）
ENVIRONMENT=dev GATUS_CONFIG_PATH=./config.yaml go run main.go

# 运行全部测试
go test ./... -cover

# 运行单个测试
go test ./path/to/package -run TestFunctionName -cover

# 带竞态检测运行测试（ICMP 测试需要 sudo）
sudo go test ./... -race -cover

# 前端
npm --prefix web/app install        # 安装依赖
npm --prefix web/app run build       # 构建（输出到 web/static/）
npm --prefix web/app run serve       # 开发服务器
npm --prefix web/app run lint        # 前端 lint

# Docker
docker build -t twinproduction/gatus:latest .
docker run -p 8080:8080 --name gatus twinproduction/gatus:latest
```

## 架构

### 启动流程（main.go）

`main()` → 加载配置（含 Managed Overlay 合并）→ 初始化存储 → 启动 HTTP 服务器（goroutine）→ 启动 watchdog 监控（每个端点一个 goroutine）→ 启动配置文件监听 → 等待 SIGINT/SIGTERM。

### 核心包

| 包 | 职责 |
|---|------|
| `config/` | YAML 配置加载、验证与热重载。支持目录下多 YAML 文件通过 `deepmerge.YAML` 深度合并。通过 `os.ExpandEnv` 展开环境变量。 |
| `config/endpoint/` | 核心领域对象：`Endpoint` 结构体、`EvaluateHealth()`、条件解析与评估、占位符解析（`[STATUS]`、`[BODY]`、`[RESPONSE_TIME]` 等）、`TamperConfig` 防篡改检测 |
| `config/suite/` | Suite（端点套件）定义与执行。套件中端点按顺序执行，通过 `gontext.Gontext` 上下文传递数据（`store` 映射 + `[CONTEXT].path` 占位符） |
| `config/gontext/` | 线程安全的上下文存储，支持 Suite 中端点间数据传递 |
| `watchdog/` | 监控引擎。为每个端点/Suite 启动 goroutine，通过 `semaphore` 管理并发，处理告警触发与恢复，含防篡改检测逻辑 |
| `alerting/provider/` | `AlertProvider` 接口 + `Config[T]` 泛型接口 + 38 种实现（slack、discord、telegram、pagerduty、wecom 等） |
| `storage/store/` | `Store` 接口，含 `memory/` 和 `sql/`（SQLite/PostgreSQL）两种实现。通过 `store.Get()`/`store.Initialize()` 单例访问 |
| `api/` | REST API 处理器（Fiber）。含公开路由（徽章、图表、外部端点推送）和保护路由（状态查询、Admin API v2） |
| `controller/` | HTTP 服务器生命周期管理（Fiber 应用创建、中间件、监听） |
| `security/` | 认证：Basic Auth 和 OIDC |
| `client/` | HTTP 客户端配置：代理、TLS、OAuth2、IAP、DNS 解析器、SSH 隧道 |
| `metrics/` | Prometheus 指标注册与发布 |
| `web/` | 嵌入式前端静态文件（`embed.FS`），Vue 3 源码位于 `web/app/` |

### 设计模式

- **`ValidateAndSetDefaults()` 约定**：几乎所有配置结构体都实现此方法，用于解析后验证并填充默认值。新增配置类型时应遵循此模式。
- **编译时接口检查**：使用 `var _ Interface = (*Impl)(nil)` 确保 `AlertProvider`、`Config[T]` 和 `Store` 实现的正确性。
- **单例存储**：通过 `store.Get()` 访问，由 `store.Initialize()` 初始化一次。
- **并发控制**：`golang.org/x/sync/semaphore` 限制并发健康检查数量（默认 3，通过 `concurrency` 配置项调整，0 = 无限制）。
- **配置热重载**：`listenToConfigurationFileChanges()` 默认每 5 秒轮询文件修改时间（可通过 `GATUS_CONFIG_WATCH_INTERVAL` 环境变量调整）。同时监听 `config.ImmediateReloadRequests()` channel，Admin API 可通过 `POST /api/v1/admin/reload` 触发即时重载。检测到变更时：停止 → 保存 → 重新加载 → 启动。

### 配置系统

配置路径解析顺序：`GATUS_CONFIG_PATH` 环境变量 → `config/config.yaml` → `config/config.yml`。路径可以是目录（所有 `.yml`/`.yaml` 文件会被深度合并）。配置值支持 `$ENV_VAR` 环境变量展开（`$$` 表示字面量 `$`）。

**Managed Overlay**：`.gatus-managed-overlay.json` 文件（存储在配置文件同目录，路径可通过 `GATUS_MANAGED_OVERLAY_PATH` 覆盖）允许通过 Web UI/Admin API 动态管理端点、外部端点、套件和告警配置。Overlay 会在配置加载时合并到 YAML 基础配置之上，采用原子写入（先写临时文件再 rename）确保安全。

### 关键环境变量

| 环境变量 | 用途 |
|----------|------|
| `GATUS_CONFIG_PATH` | 配置文件/目录路径 |
| `GATUS_LOG_LEVEL` | 日志级别 |
| `GATUS_CONFIG_WATCH_INTERVAL` | 配置监听间隔（默认 5s） |
| `GATUS_DELAY_START_SECONDS` | 延迟启动秒数 |
| `GATUS_MANAGED_OVERLAY_PATH` | Managed Overlay 文件路径 |
| `ENVIRONMENT=dev` | 启用开发模式 CORS |
| `ROUTER_TEST=true` | 测试模式跳过实际网络监听 |
| `MOCK_ALERT_PROVIDER=true` | 模拟告警发送 |

### API 路由概览

**公开路由**：`GET /api/v1/config`（配置）、`GET /api/v1/endpoints/:key/health/badge.svg`（健康徽章）、`GET /api/v1/endpoints/:key/uptimes/:duration/badge.svg`（可用性徽章）、`GET /api/v1/endpoints/:key/response-times/:duration/chart.svg`（响应时间图表）、`POST /api/v1/endpoints/:key/external`（外部端点推送，需 Bearer token）

**保护路由**（需认证）：`GET /api/v1/endpoints/statuses`（端点状态）、`GET /api/v1/suites/statuses`（套件状态）、`GET /api/v1/groups`（分组列表）

**Admin API v2**（需认证）：`/api/v1/admin/` 下的端点/套件/外部端点 CRUD、`managed-config` CRUD、`reload`（即时重载）、`monitors`（统一监控列表）、`monitors/batch`（批量操作）、`export`/`import`（导入导出）、`notifications/:type`（通知渠道管理）、`audit-logs`（审计日志查询）

**SPA 路由**：`/`（首页）、`/endpoints/:key`、`/suites/:key`、`/admin/`（管理面板）、`/certificates`（证书监控）

### 防篡改检测（Tamper Detection）

`config/endpoint/tamper.go` 定义 `TamperConfig`，支持两种检测模式：
- **Body Size Drift**：基于历史响应体大小基线，检测当前响应偏移（默认：baseline-samples=20, drift-threshold-percent=20%, consecutive-breaches=3）
- **Content Checks**：`required-substrings`（必含关键词）和 `forbidden-substrings`（禁含关键词）

在 `watchdog/endpoint_tamper.go` 中，`applyBodySizeTamperDetection()` 在 `EvaluateHealth()` 之后调用。

### 审计日志系统

所有 Admin API 操作通过 `api/admin_common.go` 中的 `writeAdminAudit()` 记录审计日志，包含操作人、操作类型、实体、前后快照（敏感字段自动脱敏）等信息。`Store` 接口提供 `InsertAdminAuditLog()`、`GetAdminAuditLogs()`、`DeleteAdminAuditLogsOlderThan()` 方法。

### 前端架构

- **技术栈**：Vue 3 + Vue Router 4 + Tailwind CSS 3 + Chart.js 4 + lucide-vue-next
- **UI 组件**：shadcn/ui 风格，位于 `web/app/src/components/ui/`
- **i18n**：自定义实现（非 vue-i18n），支持 en-US 和 zh-CN，通过 `useI18n()` composable 提供 `t()` 翻译函数（`web/app/src/i18n/`）
- **构建工具**：@vue/cli-service，构建输出到 `web/static/`

### 测试说明

- 使用标准 `testing` 包，未使用外部测试框架。
- 测试文件与源文件同目录（`*_test.go`）。
- Controller 测试使用 `ROUTER_TEST=true` 环境变量跳过实际网络监听。
- 告警 Provider 测试使用 `MOCK_ALERT_PROVIDER=true` 环境变量模拟发送。
- ICMP/ping 测试需要 root 权限（`sudo`）。
- 测试证书位于 `testdata/` 目录。

### 新增 Alert Provider

1. 在 `alerting/provider/<name>/` 下创建新包。
2. 实现 `AlertProvider` 接口：`Validate()`、`Send()`、`GetDefaultAlert()`、`ValidateOverrides()`。
3. 实现 `Config[T]` 泛型接口：`Validate()`、`Merge(override *T)`。
4. 在 `alerting/provider/provider.go` 添加编译时检查：`var _ AlertProvider = (*AlertProvider)(nil)` 和 `var _ Config[<name>.Config] = (*<name>.Config)(nil)`。
5. 在 `alerting.Config` 结构体中添加 provider 字段，并在 provider 查找中注册。
