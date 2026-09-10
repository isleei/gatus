# 合上游操作手册（本 Fork）

> 对应里程碑 **阶段 2.2**（见 [`MILESTONES.md`](./MILESTONES.md)）。  
> 目标：定期吸收 [TwiN/gatus](https://github.com/TwiN/gatus) 修复与能力，同时**不丢**「[本 Fork 需长期保留的能力](./MILESTONES.md#本-fork-需长期保留的能力)」。

建议节奏：**每季度一次**（或安全修复按需提前）。

---

## 步骤

```bash
git fetch upstream
git checkout -b merge/upstream-master-YYYY-MM master
git merge upstream/master
# 解决冲突 → 见下方清单与冲突高发区
go test ./...
make frontend-build          # 若改动了 web/app，务必重生成 web/static
git push -u origin HEAD
gh pr create --repo isleei/gatus --base master --title "merge: upstream master YYYY-MM"
# CI 绿后：gh pr merge <n> --repo isleei/gatus --merge
```

远程约定：`origin` = `isleei/gatus`，`upstream` = `TwiN/gatus`。

---

## 保留能力检查清单（合入前后勾选）

合上游或重构时，以下能力**不得丢失**（与里程碑「本 Fork 需长期保留的能力」对齐，并含 suite 告警）：

- [ ] **Admin 控制台 v2**（审计日志、managed overlay、热重载）
- [ ] **企业微信（WeCom）告警**（含 markdown 硬限 / 文本模板）
- [ ] **证书监控页** `/certificates`
- [ ] **分组鉴权 / groups API**
- [ ] **端点防篡改**（body-size drift / 关键词；相关 metrics）
- [ ] **中文 i18n** 与 `README_zh.md` / 部署文档
- [ ] **Suite 级告警**（`suites[].alerts` + 重启后 Triggered 恢复）

冲突高发区：`controller/admin*`、`web/app/src/App.vue`、`go.mod`、`storage/store/sql/`、已生成的 `web/static/*`。

---

## 前端静态资源

改动 `web/app` 后必须：

```bash
make frontend-build   # rebuild web/static from web/app
```

不要只提交 Vue 源码而漏掉 `web/static`（或反过来被上游静态文件覆盖后未重建）。

---

## Workflow 文件与 `workflow` scope

若 `gh` / push 因缺少 OAuth **`workflow` scope** 拒绝写入 `.github/workflows/*`：

1. 可先**跳过**这些文件，完成其余合入并开 PR（不阻塞主流程）。
2. 有权限后再：`gh auth refresh -s workflow`，按 [`OPS-RUNBOOK.md`](./OPS-RUNBOOK.md) §4.4 把跳过的上游 workflow 变更补入。

---

## 合入后快速自检

- [ ] `go test` 相关包通过（至少 alerting / api / watchdog / storage）
- [ ] Admin、WeCom、certificates、groups、tamper、zh 文案、suite alerts 冒烟或单测仍在
- [ ] 生产仍使用**本仓库镜像**，非 `ghcr.io/twin/gatus`
