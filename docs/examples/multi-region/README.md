# 多地域探测（external-endpoints）示例

本目录演示如何用 **主实例 + 卫星推送** 做多地域汇总，不要求卫星跑完整 Gatus Admin。

## 架构

```
[卫星 EU] --POST Bearer--> [主 Gatus] 看板汇总
[卫星 US] --POST Bearer-->
```

- **主实例**：配置 `external-endpoints`，对外暴露 `POST /api/v1/endpoints/{key}/external`
- **卫星**：在本地域探测目标，用 curl（或 gatus-cli）把结果推到主实例

`key` 规则与站内一致：`group_name`（空格等会按 Gatus key 规则转换，通常为 `region-eu_website`）。

## 主实例

见 [`primary-config.yaml`](./primary-config.yaml)：

1. 设置 `GATUS_EXTERNAL_TOKEN_EU` / `GATUS_EXTERNAL_TOKEN_US`（随机长 token）
2. `ui.default-sort-by: group` 便于按地域分组查看
3. 用本 fork 镜像构建部署（勿用官方 twin 镜像）

## 卫星推送

[`satellite-push.sh`](./satellite-push.sh) 最小示例：

```bash
export GATUS_PRIMARY_URL='https://status.example.com'
export GATUS_EXTERNAL_KEY='region-eu_website'
export GATUS_EXTERNAL_TOKEN='replace-with-eu-token'
export GATUS_TARGET_URL='https://app.example.com/health'
chmod +x satellite-push.sh
# cron 每分钟
* * * * * /opt/gatus-sat/satellite-push.sh
```

等价 curl：

```bash
curl -X POST \
  -H "Authorization: Bearer $GATUS_EXTERNAL_TOKEN" \
  "$GATUS_PRIMARY_URL/api/v1/endpoints/region-eu_website/external?success=true&duration=120ms"
```

失败时加 `&error=timeout`（需 URL 编码）。

## 与哨兵的区别

| | 哨兵 (`docs/examples/sentinel`) | 多地域 |
|--|--|--|
| 目的 | 防主站自盲 | 从多出口探测同一业务 |
| 告警 | 卫星本地 WeCom 即可 | 通常汇总到主站告警 |
| 数据 | 独立面板 | 推送到主站 `external-endpoints` |

仓库内仅占位符，无真实密钥。真正跨区域部署属运维事项。
