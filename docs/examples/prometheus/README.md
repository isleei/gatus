# Prometheus 抓取示例

1. Gatus 配置：

```yaml
metrics: true
```

2. 将 [`scrape.yml`](./scrape.yml) 片段并入 Prometheus `scrape_configs`。

3. 验证：`curl -s http://gatus:8080/metrics | head`。

更多说明见 [`docs/DEPLOYMENT.md`](../../DEPLOYMENT.md)「Prometheus 指标」。

## Tamper body-size drift (fork)

When an endpoint has `tamper.enabled: true`, Gatus also exports:

| Metric | Type | Labels | Meaning |
|--------|------|--------|---------|
| `gatus_results_body_size_drift_percent` | gauge | `key`, `group`, `name`, `type` (+ extra) | Absolute drift % vs rolling baseline |
| `gatus_results_body_size_drift_breach_streak` | gauge | same | Consecutive threshold breaches in a row |
| `gatus_results_body_size_drift_breaches_total` | counter | same | Total threshold-breach evaluations |

Example PromQL — alert when any endpoint has a sustained drift streak:

```promql
gatus_results_body_size_drift_breach_streak > 0
```

Drift magnitude for a single endpoint:

```promql
gatus_results_body_size_drift_percent{key="security_home"}
```

