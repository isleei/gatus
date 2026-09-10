# 运维手册（OPS Runbook）

> 仓库：`isleei/gatus`（本 fork）  
> 适用：Stages 0–4 仓内交付已完成后（PR [#4](https://github.com/isleei/gatus/pull/4)–[#8](https://github.com/isleei/gatus/pull/8)），**仍须人工完成**的线上操作。  
> 原则：清单可勾选；命令可复制；**密钥仅用占位符**，真实值只进环境变量 / 密钥库，永不进公开 git。

相关文档：[`DEPLOYMENT.md`](./DEPLOYMENT.md) · [`MILESTONES.md`](./MILESTONES.md) · [`examples/sentinel/`](./examples/sentinel/) · [`examples/multi-region/`](./examples/multi-region/)

---

## 1. 凭据轮换（里程碑 0.2）

**背景**：公开仓历史可能仍含旧密钥；即使当前 `config.yaml` 已是 `$GATUS_*` 占位符，**仍必须轮换**线上凭据并作废旧 webhook。

本 fork 实际使用的变量（见仓库根 `config.yaml` 与 `docs/DEPLOYMENT.md` 生产 compose）：

| 用途 | 环境变量 / 配置 |
|------|-----------------|
| Postgres（推荐单 DSN） | `GATUS_DB_URL`（如 `postgres://…`） |
| Postgres（compose 拆分） | `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` → 组装进 `GATUS_DB_URL` |
| Admin Basic Auth | `GATUS_ADMIN_USER`、`GATUS_ADMIN_PASSWORD_BCRYPT_BASE64`（`go run ./cmd/passwd`） |
| 企业微信 | `GATUS_WECOM_WEBHOOK_URL` |

### 1.1 生成新凭据（仅本地 / 密钥库，勿提交）

```bash
# --- Postgres 新密码（示例；按你们密钥策略生成）---
# 把结果写入密钥库，占位：REPLACE_NEW_POSTGRES_PASSWORD
openssl rand -base64 32

# --- Admin bcrypt（在本 fork 仓库根目录）---
go run ./cmd/passwd
# 按提示输入新明文密码；stdout 为 base64(bcrypt)
# 导出占位示例（勿把真实 hash 写进 git）：
#   export GATUS_ADMIN_USER='admin'
#   export GATUS_ADMIN_PASSWORD_BCRYPT_BASE64='REPLACE_NEW_BCRYPT_BASE64'

# --- 企业微信：在企微群机器人控制台新建 webhook ---
# 新 URL 形如：
#   https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=REPLACE_NEW_KEY
# 仅写入密钥库 / 运行环境，变量名：GATUS_WECOM_WEBHOOK_URL
```

组装 DSN 示例（compose 内常用写法）：

```bash
# 勿 echo 到日志或提交；仅示意结构
export GATUS_DB_URL="postgres://${POSTGRES_USER}:REPLACE_NEW_POSTGRES_PASSWORD@postgres:5432/${POSTGRES_DB:-gatus}?sslmode=disable"
```

### 1.2 只写入环境 / 密钥库并更新运行配置

- [ ] 新密码 / bcrypt / webhook **只**写入：`.env`（未入库）、systemd `EnvironmentFile`、K8s `Secret`、或云密钥管理
- [ ] 更新 compose / systemd / Deployment 中的对应 env（**不要**改公开仓里的占位符为真值）
- [ ] 若 Postgres 密码变更：先在库内改角色密码，再改 `GATUS_DB_URL` / `POSTGRES_PASSWORD`，避免启动窗口连不上

```bash
# Postgres 角色改密示例（在已认证的 psql 会话中；占位符）
# ALTER ROLE gatus WITH PASSWORD 'REPLACE_NEW_POSTGRES_PASSWORD';

# Docker Compose 示例：改完 .env 后
docker compose up -d

# systemd 示例
sudo systemctl daemon-reload && sudo systemctl restart gatus

# Kubernetes 示例
kubectl -n monitoring apply -f gatus-secret.yaml
kubectl -n monitoring rollout restart deployment/gatus
```

### 1.3 重启后验收

- [ ] `GET /health` 返回 UP
- [ ] Admin 用**新**用户名/密码可登录（旧密码失败）
- [ ] 触发一次 WeCom 测试告警（可临时把某测试端点弄失败，或走你们已有的测试通道），确认**新**机器人收到消息
- [ ] 在企业微信控制台 **删除/禁用旧 webhook**（作废旧 key）
- [ ] 确认公开 git 中无新密钥；历史泄露项已全部轮换

### 1.4 明确提醒

> 密钥可能仍存在于**公开 git 历史**。仓内 Stage 0 清理（占位符）不能替代线上轮换。完成 0.2 前，视为凭据仍可能被滥用。

---

## 2. 哨兵部署（阶段 1）

完整说明与文件：**[`docs/examples/sentinel/`](./examples/sentinel/)**  
（[`README.md`](./examples/sentinel/README.md) · [`compose.yaml`](./examples/sentinel/compose.yaml) · [`config.yaml`](./examples/sentinel/config.yaml)）

### 检查清单

1. [ ] 选定**另一台主机 / 另一可用区 / 至少不同出口网络**（与主实例同机同网仍会自盲）
2. [ ] 在该主机克隆**本 fork**，构建镜像（**禁止** `twinproduction/gatus` / `ghcr.io/twin/gatus`）：

```bash
git clone https://github.com/isleei/gatus.git
cd gatus
docker build -t gatus:local .
```

3. [ ] 按哨兵 README 导出环境变量（占位符）：

```bash
cd docs/examples/sentinel
export GATUS_PRIMARY_URL='https://status.example.com'
export GATUS_WECOM_WEBHOOK_URL='https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=REPLACE'
export GATUS_SENTINEL_CRITICAL_URL_1='https://app.example.com/health'
export GATUS_SENTINEL_CRITICAL_URL_2='https://api.example.com/health'
```

4. [ ] 启动：

```bash
docker compose up -d --build
```

5. [ ] 验证哨兵健康（默认宿主机 **8081 → 8080**）：

```bash
curl -sS http://127.0.0.1:8081/health
# 期望含 UP
```

6. [ ] 演练告警：短暂破坏主站可达性（停主实例或临时改 `GATUS_PRIMARY_URL` 指向不可达地址）→ 确认企业微信收到哨兵告警 → **立即恢复**主站/配置
7. [ ] 记录部署主机、镜像 tag、机器人名称，便于交接

---

## 3. Postgres dump / restore 演练

详细命令与说明见 [`DEPLOYMENT.md`「Postgres 备份与恢复演练」](./DEPLOYMENT.md#postgres-备份与恢复演练)。此处仅紧凑检查清单：

- [ ] **Dump**：用 `pg_dump --format=custom` 产出一份带时间戳的 `.dump`（`PGPASSWORD` 等从密钥库注入，勿进 git）
- [ ] **Restore to scratch**：在非生产空库（如 `gatus_restore_drill`）上 `pg_restore --clean --if-exists --no-owner`
- [ ] **Verify**：抽查端点 / 历史结果行数量级合理；可选将临时 Gatus `storage.path` 指到演练库能启动并看到数据
- [ ] **Note time**：记录 dump 耗时、restore 耗时、演练日期与操作人（团队 RTO/RPO）

```bash
# 时间戳记录示例（演练结束时填写）
# DUMP_STARTED=... DUMP_FINISHED=... RESTORE_FINISHED=... OPERATOR=...
```

---

## 4. 生产加固

### 4.1 Admin 审计保留

在生产 `storage` 中设置（或确认已设置）：

```yaml
storage:
  type: postgres
  path: "$GATUS_DB_URL"
  admin-audit-max-age: 720h   # 30 天；按合规可改为 168h / 2160h 等
```

- [ ] 配置已生效；可选手动清理：`DELETE /api/v1/admin/audit-logs?days=30`（需 Admin 鉴权）

### 4.2 反代限制敏感路径

将 `/metrics`、`/admin`、`/api/v1/admin/` 限制为内网 CIDR（示例 `10.0.0.0/8`；按实际改）。

**nginx：**

```nginx
# 内网 CIDR — 按实际替换
geo $gatus_admin_ok {
    default       0;
    10.0.0.0/8    1;
    192.168.0.0/16 1;
    127.0.0.1/32  1;
}

server {
    # ... 上游 proxy_pass 到 gatus:8080 ...

    location /metrics {
        if ($gatus_admin_ok = 0) { return 403; }
        proxy_pass http://gatus_upstream;
    }
    location /admin {
        if ($gatus_admin_ok = 0) { return 403; }
        proxy_pass http://gatus_upstream;
    }
    location /api/v1/admin/ {
        if ($gatus_admin_ok = 0) { return 403; }
        proxy_pass http://gatus_upstream;
    }
}
```

**Caddy：**

```caddy
handle /metrics* {
    @denied not remote_ip 10.0.0.0/8 192.168.0.0/16 127.0.0.1
    respond @denied 403
    reverse_proxy gatus:8080
}
handle /admin* {
    @denied not remote_ip 10.0.0.0/8 192.168.0.0/16 127.0.0.1
    respond @denied 403
    reverse_proxy gatus:8080
}
handle /api/v1/admin/* {
    @denied not remote_ip 10.0.0.0/8 192.168.0.0/16 127.0.0.1
    respond @denied 403
    reverse_proxy gatus:8080
}
```

公开状态页（`/` 等）可继续对业务需要的来源开放；整站登录墙需在反代另配（非 Gatus 默认）。

### 4.3 OIDC 与多地域

- OIDC 配置专节：[`DEPLOYMENT.md`「OIDC（单点登录）」](./DEPLOYMENT.md#oidc单点登录)
- 多地域 / 卫星推送：[`docs/examples/multi-region/`](./examples/multi-region/)

### 4.4（可选）补入跳过的上游 workflow 文件

此前合上游时若因缺少 `workflow` scope 跳过了 `.github/workflows` 变更：

```bash
# 可能需要交互式登录；在操作者本机执行
gh auth refresh -s workflow
# 再按你们的上游同步流程把跳过的 workflow 文件补入并开 PR
```

- [ ] 已评估是否需要；若需要，已由有权限的操作者完成（不阻塞其余运维项）

---

## 5. 完成定义（Done）

操作者勾选后可将对应里程碑标为「运维已完成」：

- [ ] **0.2** Postgres / Admin Basic Auth / WeCom 均已轮换；旧 WeCom webhook 已在控制台作废
- [ ] **阶段 1 哨兵** 已在异主机/异 AZ 部署；`/health` on 8081 正常；主站故障演练收到 WeCom 并已恢复
- [ ] **备份演练** dump → scratch restore → 行数/端点抽查通过；耗时已记录
- [ ] **加固** `admin-audit-max-age` 已设；反代已限制 `/metrics`、`/admin`、`/api/v1/admin/`
- [ ] （可选）OIDC / 多地域 / `workflow` scope 补文件已按需处理

完成后请同步更新 [`MILESTONES.md`](./MILESTONES.md) 对应细项状态。

---

*文档创建：2026-09-10（Stages 0–4 仓内完成后的人工运维清单）*
