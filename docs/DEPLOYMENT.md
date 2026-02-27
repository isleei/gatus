# Gatus 部署文档

> 本项目为 **Gatus** — 面向开发者的服务健康状态监控面板。
> 支持 HTTP、ICMP、TCP、DNS、gRPC、WebSocket、SSH 等多种协议，提供 38 种告警渠道。

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

# 安全认证（Basic Auth）
security:
  basic:
    username: "admin"
    # 使用 bcrypt+base64 加密的密码，可通过 ./gatus-passwd 生成
    password-bcrypt-base64: "JDJhJDA4JC82VHZtbGtpdkJjT3pzSnN0V1o4U3VKaXMyNzlFWkNaZXI1MUNQYkJQZ2xVYUtYeC4zdFVp"

# 告警配置（以企业微信为例）
alerting:
  wecom:
    webhook-url: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=YOUR_KEY"
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

使用项目内置工具 `gatus-passwd` 生成 Basic Auth 密码：

```bash
./gatus-passwd
# 按提示输入密码，输出 bcrypt+base64 编码结果
# 将结果填入 security.basic.password-bcrypt-base64
```

---

## 部署方式一：本地二进制构建

### 1. 构建

```bash
# 克隆仓库
git clone https://github.com/TwiN/gatus.git
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

### 使用官方镜像

```bash
docker run -d \
  --name gatus \
  --restart unless-stopped \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  twinproduction/gatus:latest
```

### 自行构建镜像

```bash
# 构建镜像（多阶段构建：golang:alpine → scratch）
docker build -t my-gatus:latest .

# 运行
docker run -d \
  --name gatus \
  --restart unless-stopped \
  -p 8080:8080 \
  -v $(pwd)/config:/config \
  -e GATUS_LOG_LEVEL=INFO \
  my-gatus:latest
```

> **注意**：最终镜像基于 `scratch`，体积极小，仅包含二进制文件和 CA 证书。

---

## 部署方式三：Docker Compose（推荐）

### 基础版（内存存储）

适用于轻量部署，重启后历史数据**不保留**。

```yaml
# compose.yaml
services:
  gatus:
    image: twinproduction/gatus:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./config:/config
    environment:
      - GATUS_LOG_LEVEL=INFO
```

```bash
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
    image: twinproduction/gatus:latest
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - POSTGRES_USER=${POSTGRES_USER:-gatus}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB:-gatus}
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
    username: "admin"
    password-bcrypt-base64: "YOUR_BCRYPT_BASE64_PASSWORD"

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
          image: twinproduction/gatus:latest
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
./gatus-passwd
```

### OIDC（单点登录）

```yaml
security:
  oidc:
    issuer-url: "https://your-oidc-provider.com"
    redirect-url: "https://your-gatus-domain.com/authorization-code/callback"
    client-id: "your-client-id"
    client-secret: "${OIDC_CLIENT_SECRET}"
    scopes: ["openid"]
    allowed-subjects: ["user@example.com"]
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

Prometheus 抓取配置：

```yaml
# prometheus.yml
scrape_configs:
  - job_name: gatus
    static_configs:
      - targets: ['gatus:8080']
```

### 配置热重载

配置文件变更后 Gatus 会自动重载（默认每 5 秒检测）。也可通过 Admin API 立即触发：

```bash
curl -X POST http://admin:password@localhost:8080/api/v1/admin/reload
```

---

## 常见问题

### Q: 容器启动后无法连接 PostgreSQL

检查以下项目：
1. PostgreSQL 容器是否健康（`docker compose ps`）
2. `config.yaml` 中连接字符串的 host 是否为服务名（`postgres`，非 `localhost`）
3. 环境变量是否正确注入（`docker compose config` 查看展开后的配置）

### Q: 密码如何更新

```bash
# 重新生成密码
./gatus-passwd

# 更新 config.yaml 中的 password-bcrypt-base64 后，等待热重载或手动触发
curl -X POST http://admin:old_password@localhost:8080/api/v1/admin/reload
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

*文档生成时间：2026-02-27*
