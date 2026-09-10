# Prometheus 抓取示例

1. Gatus 配置：

```yaml
metrics: true
```

2. 将 [`scrape.yml`](./scrape.yml) 片段并入 Prometheus `scrape_configs`。

3. 验证：`curl -s http://gatus:8080/metrics | head`。

更多说明见 [`docs/DEPLOYMENT.md`](../../DEPLOYMENT.md)「Prometheus 指标」。
