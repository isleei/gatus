# Gatus 部署文档

> 本项目为 **Gatus** — 面向开发者的服务健康状态监控面板。
> 支持 HTTP、ICMP、TCP、DNS、gRPC、WebSocket、SSH 等多种协议，提供 38 种告警渠道。

> 演进规划见 [`docs/MILESTONES.md`](./MILESTONES.md)。  
> Stages 0–4 仓内完成后的**人工运维检查清单**见 [`docs/OPS-RUNBOOK.md`](./OPS-RUNBOOK.md)（凭据轮换、哨兵部署、备份演练、生产加固）。

---

## 目录

- [环境要求](#环境要求)
- [配置文件](#配置文件)
- [部署方式一：本地二进制构建](#部署方式一本地二进制构建)
- [部署方式二：Docker](#部署方式二docker)
- [部署方式三：Docker Compose（推荐）](#部署方式三docker-compose推荐)
  - [基础版（内存存储）](#基础版内存存储)
  - [生产版（PostgreSQL 存储）](#生产版postgresql-存储)
- [部署方式四：Kubernetes](#部署方式四kubernetes)
- [存储配置](#存储配置)
- [安全配置](#安全配置)
- [关键环境变量](#关键环境变量)
- [健康检查与监控](#健康检查与监控)
- [防自盲 / 哨兵](#防自盲--哨兵)
- [Postgres 备份与恢复演练](#postgres-备份与恢复演练)
- [运维手册（人工清单）](#运维手册人工清单)
- [常见问题](#常见问题)

---

## 环境要求

| 组件 | 版本要求 |
|------|---------|
| Go | ≥ 1.22（如从源码构建） |
| Docker | ≥ 20.10 |
| Docker Compose | ≥ 2.0（使用 `compose.yaml`） |
| PostgreSQL | ≥ 13（可选，使用数据库存储时） |
| SQLite | 内置支持（CGO_ENABLED=0，使用 modernc.org/sqlite） |

---

## 配置文件

Gatus 使用 YAML 配置文件。**配置路径解析优先级**：

1. 环境变量 `GATUS_CONFIG_PATH`（文件路径或目录）
2. `config/config.yaml`
3. `config/config.yml`

> **目录模式**：当 `GATUS_CONFIG_PATH` 指向目录时，目录下所有 `.yaml` / `.yml` 文件会被**深度合并**，便于按模块拆分配置。

> **密钥与 overlay**：仓库内 `config.yaml` 仅为占位示例（`$GATUS_*`）。生产 managed overlay（`.gatus-managed-overlay.json`）含真实端点与 webhook，**不得**提交公开仓；可参考 [`docs/examples/gatus-managed-overlay.example.json`](./examples/gatus-managed-overlay.example.json)，并通过 `GATUS_MANAGED_OVERLAY_PATH` 指向私有路径。


### 最小配置示例

```yaml
# config/config.yaml
endpoints:
  - name: 示例服务
    url: "https://example.com/health"
    interval: 1m
    conditions:
      - "[STATUS] == 200"
      - "[RESPONSE_TIME] < 500"
```

### 完整配置示例（含存储与安全）

```yaml
# config/config.yaml

# 存储（使用 PostgreSQL）
storage:
  type: postgres
  path: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable"
  caching: true

# 安全认证（Basic Auth）— 用户名与 bcrypt 摘要均来自环境变量，勿提交真实值
security:
  basic:
    username: "$GATUS_ADMIN_USER"
    # 用本仓库源码生成：go run ./cmd/passwd
    # 将输出写入环境变量 GATUS_ADMIN_PASSWORD_BCRYPT_BASE64（勿把真实 hash 写进 git）
    password-bcrypt-base64: "$GATUS_ADMIN_PASSWORD_BCRYPT_BASE64"

# 告警配置（以企业微信为例）
alerting:
  wecom:
    webhook-url: "$GATUS_WECOM_WEBHOOK_URL"
    title: "Gatus 监控"
    # 可选：中文告警正文模板（留空则使用英文默认，行为与历史一致）
    # 占位符: [ENDPOINT] [ENDPOINT_NAME] [ENDPOINT_GROUP] [ALERT_DESCRIPTION]
    #         [FAILURE_COUNT] [SUCCESS_COUNT] [RESULT_CONDITIONS] [RESULT_ERRORS]
    text-triggered: |
      **触发** [ENDPOINT]
      失败阈值: [FAILURE_COUNT]
      描述: [ALERT_DESCRIPTION]
      [RESULT_CONDITIONS]
    text-resolved: |
      **恢复** [ENDPOINT]
      成功阈值: [SUCCESS_COUNT]
      描述: [ALERT_DESCRIPTION]
    default-alert:
      failure-threshold: 3
      success-threshold: 2

# 监控端点
endpoints:
  - name: 前端服务
    group: 核心服务
    url: "https://your-domain.com/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 500"
    alerts:
      - type: wecom
        failure-threshold: 3
        success-threshold: 2

  - name: 证书检查
    url: "https://your-domain.com/"
    interval: 1h
    conditions:
      - "[CERTIFICATE_EXPIRATION] > 168h"  # 7天预警
    alerts:
      - type: wecom

  - name: 域名到期检查
    url: "https://your-domain.com/"
    interval: 24h
    conditions:
      - "[DOMAIN_EXPIRATION] > 720h"  # 30天预警
    alerts:
      - type: wecom
```

### 生成 bcrypt 密码

**不要**把 `gatus-passwd` 二进制或真实 bcrypt hash 提交进仓库。用源码本地生成：

```bash
go run ./cmd/passwd
# 按提示输入密码，输出 bcrypt+base64 编码结果
# 导出为环境变量（推荐），或写入未入库的本地 config：
#   export GATUS_ADMIN_PASSWORD_BCRYPT_BASE64='<output>'
```

---

## 部署方式一：本地二进制构建

### 1. 构建

```bash
# 克隆本 fork（含 Admin / WeCom 等定制；生产请用本仓库）
git clone https://github.com/isleei/gatus.git
cd gatus

# 安装前端依赖并构建（前端已通过 go:embed 嵌入）
npm --prefix web/app install
npm --prefix web/app run build

# 编译 Go 二进制
go build -v -o gatus .
```

### 2. 运行

```bash
# 生产模式
GATUS_CONFIG_PATH=./config/config.yaml ./gatus

# 开发模式（启用 CORS，便于前端热调试）
ENVIRONMENT=dev GATUS_CONFIG_PATH=./config.yaml go run main.go
```

### 3. 访问

浏览器打开 `http://localhost:8080`

---

## 部署方式二：Docker

### 自建镜像（本 fork 生产必用）

> **重要**：本 fork 的 Admin v2、企业微信、证书页、分组鉴权等**不在**官方镜像中。  
> 生产环境请**始终**用本仓库 `Dockerfile` 构建，**不要**使用 `twinproduction/gatus` / `ghcr.io/twin/gatus` 作为本 fork 的生产镜像。

```bash
# 在本仓库根目录构建（或 make docker-build，同样打 gatus:local 标签）
docker build -t gatus:local .

# 运行（通过环境变量注入密钥；config 可用占位符 YAML）
docker run -d \
  --name gatus \
  --restart unless-stopped \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  -e GATUS_LOG_LEVEL=INFO \
  -e GATUS_DB_URL \
  -e GATUS_ADMIN_USER \
  -e GATUS_ADMIN_PASSWORD_BCRYPT_BASE64 \
  -e GATUS_WECOM_WEBHOOK_URL \
  gatus:local
```

> **注意**：最终镜像基于 `scratch`，体积极小，仅包含二进制文件和 CA 证书。

### 上游官方镜像（仅作对比 / 无 fork 定制时）

若你只需要上游功能、不需要本 fork 定制，才可考虑官方镜像；**部署本 fork 时请跳过本节**：

```bash
docker run -d \
  --name gatus \
  --restart unless-stopped \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  twinproduction/gatus:latest
```

---

## 部署方式三：Docker Compose（推荐）

### 基础版（内存存储）

适用于轻量部署，重启后历史数据**不保留**。

```yaml
# compose.yaml
services:
  gatus:
    # 本 fork：先 docker build -t gatus:local . 再 compose up
    image: gatus:local
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./config:/config
    environment:
      - GATUS_LOG_LEVEL=INFO
      - GATUS_ADMIN_USER=${GATUS_ADMIN_USER}
      - GATUS_ADMIN_PASSWORD_BCRYPT_BASE64=${GATUS_ADMIN_PASSWORD_BCRYPT_BASE64}
      - GATUS_WECOM_WEBHOOK_URL=${GATUS_WECOM_WEBHOOK_URL}
```

```bash
docker build -t gatus:local .
docker compose up -d
```

---

### 生产版（PostgreSQL 存储）

历史数据持久化，支持重启恢复。

**目录结构**：

```
deploy/
├── compose.yaml
├── config/
│   └── config.yaml
└── data/
    └── db/          # PostgreSQL 数据目录（自动创建）
```

**`compose.yaml`**：

```yaml
services:
  postgres:
    image: postgres:15-alpine
    restart: unless-stopped
    volumes:
      - ./data/db:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=${POSTGRES_DB:-gatus}
      - POSTGRES_USER=${POSTGRES_USER:-gatus}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    networks:
      - gatus-net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER:-gatus}"]
      interval: 10s
      timeout: 5s
      retries: 5

  gatus:
    # 本 fork：使用自建镜像，勿用 twinproduction/gatus
    image: gatus:local
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_USER=${POSTGRES_USER:-gatus}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB:-gatus}
      - GATUS_DB_URL=postgres://${POSTGRES_USER:-gatus}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB:-gatus}?sslmode=disable
      - GATUS_ADMIN_USER=${GATUS_ADMIN_USER:-admin}
      - GATUS_ADMIN_PASSWORD_BCRYPT_BASE64=${GATUS_ADMIN_PASSWORD_BCRYPT_BASE64}
      - GATUS_WECOM_WEBHOOK_URL=${GATUS_WECOM_WEBHOOK_URL}
      - GATUS_LOG_LEVEL=INFO
    volumes:
      - ./config:/config
    networks:
      - gatus-net
    depends_on:
      postgres:
        condition: service_healthy

networks:
  gatus-net:
```

**`.env` 文件**（与 `compose.yaml` 同目录，不要提交到 Git）：

```dotenv
POSTGRES_USER=gatus
POSTGRES_PASSWORD=your_strong_password_here
POSTGRES_DB=gatus
GATUS_ADMIN_USER=admin
GATUS_ADMIN_PASSWORD_BCRYPT_BASE64=  # go run ./cmd/passwd 生成后填入
GATUS_WECOM_WEBHOOK_URL=             # 企业微信机器人 webhook，勿提交真实 URL
```

**启动**：

```bash
docker compose up -d

# 查看日志
docker compose logs -f gatus

# 停止
docker compose down

# 停止并删除数据卷（谨慎！）
docker compose down -v
```

**`config/config.yaml`**：

```yaml
storage:
  type: postgres
  path: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable"
  caching: true

security:
  basic:
    username: "$GATUS_ADMIN_USER"
    password-bcrypt-base64: "$GATUS_ADMIN_PASSWORD_BCRYPT_BASE64"

endpoints:
  - name: 示例服务
    url: "https://example.com/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
```

---

## 部署方式四：Kubernetes

### 快速部署（单文件）

```bash
kubectl apply -f https://raw.githubusercontent.com/TwiN/gatus/master/.examples/kubernetes/gatus.yaml
```

### 完整生产级部署

**`gatus-namespace.yaml`** — 命名空间：

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: monitoring
```

**`gatus-configmap.yaml`** — 配置：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gatus-config
  namespace: monitoring
data:
  config.yaml: |
    storage:
      type: postgres
      path: "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres-svc:5432/gatus?sslmode=disable"
      caching: true
    endpoints:
      - name: 示例服务
        url: "https://example.com/health"
        interval: 5m
        conditions:
          - "[STATUS] == 200"
```

**`gatus-secret.yaml`** — 敏感信息：

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: gatus-secret
  namespace: monitoring
type: Opaque
stringData:
  POSTGRES_USER: "gatus"
  POSTGRES_PASSWORD: "your_strong_password"
```

**`gatus-deployment.yaml`** — 部署：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gatus
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gatus
  template:
    metadata:
      labels:
        app: gatus
    spec:
      terminationGracePeriodSeconds: 10
      containers:
        - name: gatus
          image: gatus:local  # 先构建本 fork 镜像并推到你的 registry
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
              name: http
          envFrom:
            - secretRef:
                name: gatus-secret
          resources:
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 250m
              memory: 256Mi
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
            failureThreshold: 3
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
            failureThreshold: 5
          volumeMounts:
            - mountPath: /config
              name: gatus-config
      volumes:
        - name: gatus-config
          configMap:
            name: gatus-config
---
apiVersion: v1
kind: Service
metadata:
  name: gatus-svc
  namespace: monitoring
spec:
  selector:
    app: gatus
  ports:
    - name: http
      port: 8080
      targetPort: 8080
```

**部署命令**：

```bash
kubectl apply -f gatus-namespace.yaml
kubectl apply -f gatus-secret.yaml
kubectl apply -f gatus-configmap.yaml
kubectl apply -f gatus-deployment.yaml

# 查看状态
kubectl -n monitoring get pods
kubectl -n monitoring logs -f deployment/gatus
```

---

## 存储配置

### 内存存储（默认）

```yaml
# 无需配置，或显式指定
storage:
  type: memory
```

- ✅ 零配置，开箱即用
- ❌ 重启后历史数据丢失

### SQLite 存储

```yaml
storage:
  type: sqlite
  path: "/data/gatus.db"
  caching: true
```

- ✅ 单文件，轻量持久化
- ⚠️ 需挂载持久化卷

### PostgreSQL 存储（推荐生产环境）

```yaml
storage:
  type: postgres
  path: "postgres://user:password@host:5432/dbname?sslmode=disable"
  caching: true
```

- ✅ 生产级持久化，支持高可用
- ✅ 支持审计日志、历史趋势查询

---

## 安全配置

### Basic Auth

```yaml
security:
  basic:
    username: "admin"
    password-bcrypt-base64: "<bcrypt+base64 编码的密码>"
```

生成密码：

```bash
go run ./cmd/passwd
export GATUS_ADMIN_PASSWORD_BCRYPT_BASE64='<output>'
```

### OIDC（单点登录）

上游已内置 OIDC。可与 Basic Auth **二选一或并存**（按 `security` 配置生效）。推荐生产用 OIDC，本地调试仍可用 Basic Auth。

```yaml
security:
  # 可选：保留 Basic 作为应急账号
  basic:
    username: "$GATUS_ADMIN_USER"
    password-bcrypt-base64: "$GATUS_ADMIN_PASSWORD_BCRYPT_BASE64"
  oidc:
    issuer-url: "https://your-oidc-provider.com"          # IdP issuer，勿提交真实租户
    redirect-url: "https://your-gatus-domain.com/authorization-code/callback"
    client-id: "your-client-id"
    client-secret: "${OIDC_CLIENT_SECRET}"               # 仅环境变量
    scopes: ["openid", "profile", "email"]
    # 允许登录的 subject（通常为 email / sub）；留空策略取决于上游实现，生产务必收紧
    allowed-subjects:
      - "ops@example.com"
```

运维注意：

1. 在 IdP 注册回调 URL，与 `redirect-url` 完全一致。
2. `client-secret`、issuer、client-id **不要**写入公开仓库。
3. OIDC / Basic 保护的是 **Admin 与需鉴权的 API**；公开状态页路由仍可读（见下节）。
4. 轮换 IdP 客户端密钥属运维操作，本仓库无法代劳。

---

## 公开状态页 vs Admin

| 区域 | 路径 | 鉴权 |
|------|------|------|
| 公开状态页 | `/`、`/endpoints/...`、`/suites/...`、`/certificates` | 默认公开（可按网络层限制） |
| Admin UI | `/admin` | `security.basic` / `security.oidc` |
| Admin API | `/api/v1/admin/*` | 同上（`protectedAPIRouter`） |
| 外部推送 | `/api/v1/endpoints/:key/external` | Bearer token（external-endpoints） |
| Metrics | `/metrics` | 默认公开；务必用反代/网络策略限制 |

`security.basic` / OIDC **不会**给整站所有路由加锁；公开看板与 Admin 已分离。部署时请：

- 勿将 Admin 密码或 OIDC 密钥提交到公开仓
- 对 `/metrics`、Admin 入口做来源 IP / 内网限制
- 需要「整站登录墙」时在反代层另行配置（非 Gatus 默认行为）

### 状态页默认按分组排序

```yaml
ui:
  default-sort-by: group   # name | group | health
```

首页搜索栏会读取该配置（用户本地 `localStorage` 可覆盖）。

### Suite 级告警

可在 suite 上配置与端点相同的 `alerts`（wecom / slack / custom 等）。套件失败达到 `failure-threshold` 后触发，恢复达到 `success-threshold` 后解除：

```yaml
suites:
  - name: checkout
    group: critical
    interval: 5m
    alerts:
      - type: wecom
    endpoints:
      - name: login
        url: "https://example.com/login"
        conditions:
          - "[STATUS] == 200"
```

Suite 合成消息可能较长。若使用企业微信，建议在 `alerting.wecom` 配置简短的 `text-triggered` / `text-resolved` 中文模板（见上文告警配置示例），用 `[ENDPOINT]`、`[ALERT_DESCRIPTION]`、`[FAILURE_COUNT]` 等占位符控制正文，避免默认英文长文刷屏。

本 fork 会将 WeCom markdown `content` **硬限制在 ≤4096 UTF-8 字节**（企业微信机器人文档上限）；超长时优先保留标题/导语与 `failed steps:` 摘要，并追加截断说明，避免因超限被机器人静默拒收。

### Admin 审计保留

```yaml
storage:
  type: postgres
  path: "$GATUS_DB_URL"
  admin-audit-max-age: 720h   # 30 天；0 / 省略 = 不自动清理
```

启用后进程内每日清理一次。也可手动：

```bash
curl -u admin:password -X DELETE \
  'http://localhost:8080/api/v1/admin/audit-logs?days=30'
# 或 ?maxAge=720h
```

---

## 关键环境变量

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `GATUS_CONFIG_PATH` | `config/config.yaml` | 配置文件或目录路径 |
| `GATUS_LOG_LEVEL` | `INFO` | 日志级别（`DEBUG`/`INFO`/`WARN`/`ERROR`） |
| `GATUS_CONFIG_WATCH_INTERVAL` | `5s` | 配置文件热重载轮询间隔 |
| `GATUS_DELAY_START_SECONDS` | `0` | 延迟启动秒数（等待依赖就绪） |
| `GATUS_MANAGED_OVERLAY_PATH` | 配置同目录 | Managed Overlay 文件路径 |
| `PORT` | `8080` | HTTP 监听端口 |
| `ENVIRONMENT` | — | 设为 `dev` 时启用 CORS（开发模式） |

---

## 健康检查与监控

### 内置健康端点

```
GET /health
```

响应示例：

```json
{"status": "UP"}
```

### Prometheus 指标

在配置中启用：

```yaml
metrics: true
```

指标暴露地址：`http://localhost:8080/metrics`

Prometheus 抓取配置（完整片段见 [`docs/examples/prometheus/`](./examples/prometheus/)）：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: gatus
    metrics_path: /metrics
    static_configs:
      - targets: ['gatus:8080']
```

### 配置热重载

配置文件变更后 Gatus 会自动重载（默认每 5 秒检测）。也可通过 Admin API 立即触发：

```bash
curl -X POST http://admin:password@localhost:8080/api/v1/admin/reload
```

---

## 防自盲 / 哨兵

主 Gatus 与业务同机房时，一旦本机或出口网络故障，监控自身也无法告警（自盲）。阶段 1 要求在**另一台主机 / 另一网络**部署轻量「哨兵」实例：

- 监视主实例 `GET ${GATUS_PRIMARY_URL}/health`
- 另加 1–2 个核心业务探活 URL
- 告警走企业微信（`$GATUS_WECOM_WEBHOOK_URL`）
- 短间隔（示例 1m）、memory 或 sqlite 即可

**In-repo 示例**（仅占位符，无真实密钥）：[`docs/examples/sentinel/`](./examples/sentinel/)

```bash
# 仓库根目录构建本 fork 镜像后
cd docs/examples/sentinel
# 设置 GATUS_PRIMARY_URL / GATUS_WECOM_WEBHOOK_URL / GATUS_SENTINEL_CRITICAL_URL_*
docker compose up -d --build
```

哨兵**不替代**主实例；完整端点清单、Admin、Postgres 仍由主站负责。真正上线部署仍属运维事项，见 [`MILESTONES.md`](./MILESTONES.md) 阶段 1。

多地域汇总（主实例 `external-endpoints` + 卫星推送）见 [`docs/examples/multi-region/`](./examples/multi-region/)。

---

## Postgres 备份与恢复演练

生产使用 PostgreSQL 时，请定期备份并至少做过一次恢复演练。以下为**提纲**（无真实凭据；密码与主机用环境变量 / `.env`）。

### 备份（pg_dump）

```bash
# 从运行 Postgres 的 compose 项目目录，或任意能连库的机器
export PGHOST="${POSTGRES_HOST:-localhost}"
export PGPORT="${POSTGRES_PORT:-5432}"
export PGUSER="${POSTGRES_USER:-gatus}"
export PGDATABASE="${POSTGRES_DB:-gatus}"
# PGPASSWORD 从密钥管理 / .env 注入，勿写入 git

mkdir -p ./backups
pg_dump --format=custom --file="./backups/gatus-$(date -u +%Y%m%dT%H%M%SZ).dump"
```

建议：保留最近 N 份；异地再存一份；备份任务失败要有告警。

### 恢复演练（pg_restore）

在**非生产**库或临时实例上验证，勿在未确认的生产库上直接覆盖：

```bash
# 1) 准备空库（示例）
# createdb -h "$PGHOST" -U "$PGUSER" gatus_restore_drill

# 2) 恢复
pg_restore --clean --if-exists --no-owner   -h "$PGHOST" -U "$PGUSER" -d gatus_restore_drill   ./backups/gatus-YYYYMMDDTHHMMSSZ.dump

# 3) 验收：表存在、端点历史条数合理、Gatus 指向演练库可启动
```

### 演练检查清单

- [ ] 备份任务定时跑通，产物可下载
- [ ] 用最近一份 dump 在演练库 `pg_restore` 成功
- [ ] 抽查 `endpoints` / 历史结果表行数与主库量级一致
- [ ] 临时把 Gatus `storage.path` 指到演练库能启动并看到数据
- [ ] 记录 RTO/RPO 预期，并更新本团队 runbook

---

## 运维手册（人工清单）

仓内 Stages 0–4（PR #4–#9；#10/#11 为套件告警持久化与 Admin 通知计数跟进）已交付文档与示例后，仍须人工完成的操作（凭据轮换、哨兵真实部署、Postgres 恢复演练、反代加固等）统一收敛在：

**[`docs/OPS-RUNBOOK.md`](./OPS-RUNBOOK.md)**

请按该手册编号清单勾选；本页保留详细部署与备份命令，runbook 只给可执行检查项与交叉链接。

---

## 常见问题

### Q: 容器启动后无法连接 PostgreSQL

检查以下项目：
1. PostgreSQL 容器是否健康（`docker compose ps`）
2. `config.yaml` 中连接字符串的 host 是否为服务名（`postgres`，非 `localhost`）
3. 环境变量是否正确注入（`docker compose config` 查看展开后的配置）

### Q: 密码如何更新

```bash
# 重新生成密码（写入环境变量 / 本地未入库配置，勿提交真实 hash）
go run ./cmd/passwd
export GATUS_ADMIN_PASSWORD_BCRYPT_BASE64='<new-output>'

# 热重载或手动触发
curl -X POST http://"$GATUS_ADMIN_USER":"$OLD_PASSWORD"@localhost:8080/api/v1/admin/reload
```

### Q: 如何查看运行日志

```bash
# Docker
docker logs -f gatus

# Docker Compose
docker compose logs -f gatus

# Kubernetes
kubectl -n monitoring logs -f deployment/gatus
```

### Q: 历史数据丢失

必须配置持久化存储（SQLite 或 PostgreSQL），且需挂载持久化卷。内存存储不支持数据持久化。

### Q: ICMP/Ping 监控需要特权

ICMP 探测需要 root 权限或 `NET_RAW` capability：

```yaml
# Docker Compose
services:
  gatus:
    cap_add:
      - NET_RAW
```

```yaml
# Kubernetes SecurityContext
securityContext:
  capabilities:
    add:
      - NET_RAW
```

---

*文档更新：2026-09-10（+ OPS-RUNBOOK 人工运维清单链接）*
