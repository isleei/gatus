# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Gatus 是一个面向开发者的健康状态监控面板，支持使用 HTTP、ICMP、TCP、DNS、gRPC、WebSocket、SSH 等多种协议监控服务健康状态。通过条件表达式（如 `[STATUS] == 200`、`[RESPONSE_TIME] < 150`）评估健康状况，并支持 40+ 种告警渠道发送通知。

- 模块路径：`github.com/TwiN/gatus/v5`
- Go 版本：1.25.5
- HTTP 框架：Fiber v2
- 前端：Vue 3（通过 `//go:embed` 嵌入）

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
```

## 架构

### 启动流程（main.go）

`main()` → 加载配置 → 初始化存储 → 启动 HTTP 服务器（goroutine）→ 启动 watchdog 监控（每个端点一个 goroutine）→ 启动配置文件监听 → 等待 SIGINT/SIGTERM。

### 核心包

| 包 | 职责 |
|---|------|
| `config/` | YAML 配置加载、验证与热重载（30 秒轮询）。支持目录下多 YAML 文件通过 `deepmerge.YAML` 深度合并。通过 `os.ExpandEnv` 展开环境变量。 |
| `config/endpoint/` | 核心领域对象：`Endpoint` 结构体、`EvaluateHealth()`、条件解析与评估、占位符解析（`[STATUS]`、`[BODY]`、`[RESPONSE_TIME]` 等） |
| `config/suite/` | Suite（端点套件）定义与执行 |
| `watchdog/` | 监控引擎。为每个端点/Suite 启动 goroutine，通过 `semaphore` 管理并发，处理告警触发与恢复 |
| `alerting/provider/` | `AlertProvider` 接口 + 40 余种实现（slack、discord、telegram、pagerduty 等） |
| `storage/store/` | `Store` 接口，含 `memory/` 和 `sql/`（SQLite/PostgreSQL）两种实现。通过 `store.Get()`/`store.Initialize()` 单例访问 |
| `api/` | REST API 处理器，注册在 Fiber 上。路由：`/api/v1/endpoints/statuses`、徽章、图表、外部端点推送、SPA 服务 |
| `controller/` | HTTP 服务器生命周期管理（Fiber 应用创建、中间件、监听） |
| `security/` | 认证：Basic Auth 和 OIDC |
| `client/` | HTTP 客户端配置：代理、TLS、OAuth2、IAP、DNS 解析器、SSH 隧道 |
| `metrics/` | Prometheus 指标注册与发布 |
| `web/` | 嵌入式前端静态文件（`embed.FS`），Vue 3 源码位于 `web/app/` |

### 设计模式

- **`ValidateAndSetDefaults()` 约定**：几乎所有配置结构体都实现此方法，用于解析后验证并填充默认值。新增配置类型时应遵循此模式。
- **编译时接口检查**：使用 `var _ Interface = (*Impl)(nil)` 确保 `AlertProvider` 和 `Store` 实现的正确性。
- **单例存储**：通过 `store.Get()` 访问，由 `store.Initialize()` 初始化一次。
- **并发控制**：`golang.org/x/sync/semaphore` 限制并发健康检查数量（默认 3，通过 `concurrency` 配置项调整，0 = 无限制）。
- **配置热重载**：`listenToConfigurationFileChanges()` 每 30 秒轮询文件修改时间；检测到变更时：停止 → 保存 → 重新加载 → 启动。

### 配置系统

配置路径解析顺序：`GATUS_CONFIG_PATH` 环境变量 → `config/config.yaml` → `config/config.yml`。路径可以是目录（所有 `.yml`/`.yaml` 文件会被深度合并）。配置值支持 `$ENV_VAR` 环境变量展开（`$$` 表示字面量 `$`）。

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
3. 添加编译时检查：`var _ provider.AlertProvider = (*AlertProvider)(nil)`。
4. 在 `alerting.Config` 结构体中添加 provider 字段，并在 provider 查找中注册（使用反射）。
