[![Gatus](.github/assets/logo-with-dark-text.png)](https://gatus.io)

![test](https://github.com/TwiN/gatus/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/TwiN/gatus?)](https://goreportcard.com/report/github.com/TwiN/gatus)
[![codecov](https://codecov.io/gh/TwiN/gatus/branch/master/graph/badge.svg)](https://codecov.io/gh/TwiN/gatus)
[![Go version](https://img.shields.io/github/go-mod/go-version/TwiN/gatus.svg)](https://github.com/TwiN/gatus)
[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/gatus.svg)](https://cloud.docker.com/repository/docker/twinproduction/gatus)
[![Follow TwiN](https://img.shields.io/github/followers/TwiN?label=Follow&style=social)](https://github.com/TwiN)

Gatus 是一个面向开发者的健康监控仪表盘，它能够让你使用 HTTP、ICMP、TCP 甚至 DNS
查询来监控你的服务，并通过一系列条件来评估查询结果，这些条件可以基于状态码、
响应时间、证书过期时间、响应体等多种值。锦上添花的是，每一项健康检查都可以
配合 Slack、Teams、PagerDuty、Discord、Twilio 等多种方式进行告警通知。

我个人将它部署在我的 Kubernetes 集群中，用于监控核心应用的状态：https://status.twin.sh/

_正在寻找托管解决方案？请查看 [Gatus.io](https://gatus.io)。_

<details>
  <summary><b>快速开始</b></summary>

```console
docker run -p 8080:8080 --name gatus ghcr.io/twin/gatus:stable
```

如果你更喜欢使用 Docker Hub：
```console
docker run -p 8080:8080 --name gatus twinproduction/gatus:stable
```
更多详情请参见[使用方法](#usage)
</details>

> ❤ 喜欢这个项目吗？请考虑[赞助我](https://github.com/sponsors/TwiN)。

![Gatus 仪表盘](.github/assets/dashboard-dark.jpg)

有任何反馈或问题？[创建一个讨论](https://github.com/TwiN/gatus/discussions/new)。


## 目录
- [目录](#table-of-contents)
- [为什么选择 Gatus？](#why-gatus)
- [功能特性](#features)
- [使用方法](#usage)
- [配置](#configuration)
  - [端点](#endpoints)
  - [外部端点](#external-endpoints)
  - [套件 (ALPHA)](#suites-alpha)
  - [条件](#conditions)
    - [占位符](#placeholders)
    - [函数](#functions)
  - [Web 配置](#web)
  - [UI 配置](#ui)
  - [公告](#announcements)
  - [存储](#storage)
  - [客户端配置](#client-configuration)
  - [隧道](#tunneling)
  - [告警](#alerting)
    - [配置 AWS SES 告警](#configuring-aws-ses-alerts)
    - [配置 ClickUp 告警](#configuring-clickup-alerts)
    - [配置 Datadog 告警](#configuring-datadog-alerts)
    - [配置 Discord 告警](#configuring-discord-alerts)
    - [配置邮件告警](#configuring-email-alerts)
    - [配置 Gitea 告警](#configuring-gitea-alerts)
    - [配置 GitHub 告警](#configuring-github-alerts)
    - [配置 GitLab 告警](#configuring-gitlab-alerts)
    - [配置 Google Chat 告警](#configuring-google-chat-alerts)
    - [配置 Gotify 告警](#configuring-gotify-alerts)
    - [配置 HomeAssistant 告警](#configuring-homeassistant-alerts)
    - [配置 IFTTT 告警](#configuring-ifttt-alerts)
    - [配置 Ilert 告警](#configuring-ilert-alerts)
    - [配置 Incident.io 告警](#configuring-incidentio-alerts)
    - [配置 Line 告警](#configuring-line-alerts)
    - [配置 Matrix 告警](#configuring-matrix-alerts)
    - [配置 Mattermost 告警](#configuring-mattermost-alerts)
    - [配置 Messagebird 告警](#configuring-messagebird-alerts)
    - [配置 n8n 告警](#configuring-n8n-alerts)
    - [配置 New Relic 告警](#configuring-new-relic-alerts)
    - [配置 Ntfy 告警](#configuring-ntfy-alerts)
    - [配置 Opsgenie 告警](#configuring-opsgenie-alerts)
    - [配置 PagerDuty 告警](#configuring-pagerduty-alerts)
    - [配置 Plivo 告警](#configuring-plivo-alerts)
    - [配置 Pushover 告警](#configuring-pushover-alerts)
    - [配置 Rocket.Chat 告警](#configuring-rocketchat-alerts)
    - [配置 SendGrid 告警](#configuring-sendgrid-alerts)
    - [配置 Signal 告警](#configuring-signal-alerts)
    - [配置 SIGNL4 告警](#configuring-signl4-alerts)
    - [配置 Slack 告警](#configuring-slack-alerts)
    - [配置 Splunk 告警](#configuring-splunk-alerts)
    - [配置 Squadcast 告警](#configuring-squadcast-alerts)
    - [配置 Teams 告警 *(已弃用)*](#configuring-teams-alerts-deprecated)
    - [配置 Teams Workflow 告警](#configuring-teams-workflow-alerts)
    - [配置 Telegram 告警](#configuring-telegram-alerts)
    - [配置 Twilio 告警](#configuring-twilio-alerts)
    - [配置 Vonage 告警](#configuring-vonage-alerts)
    - [配置 Webex 告警](#configuring-webex-alerts)
    - [配置 Zapier 告警](#configuring-zapier-alerts)
    - [配置 Zulip 告警](#configuring-zulip-alerts)
    - [配置自定义告警](#configuring-custom-alerts)
    - [设置默认告警](#setting-a-default-alert)
  - [维护](#maintenance)
  - [安全](#security)
    - [基本身份验证](#basic-authentication)
    - [OIDC](#oidc)
  - [TLS 加密](#tls-encryption)
  - [指标](#metrics)
    - [自定义标签](#custom-labels)
  - [连接](#connectivity)
  - [远程实例 (实验性)](#remote-instances-experimental)
- [部署](#deployment)
  - [Docker](#docker)
  - [Helm Chart](#helm-chart)
  - [Terraform](#terraform)
    - [Kubernetes](#kubernetes)
- [运行测试](#running-the-tests)
- [生产环境使用](#using-in-production)
- [常见问题](#faq)
  - [发送 GraphQL 请求](#sending-a-graphql-request)
  - [推荐间隔时间](#recommended-interval)
  - [默认超时时间](#default-timeouts)
  - [监控 TCP 端点](#monitoring-a-tcp-endpoint)
  - [监控 UDP 端点](#monitoring-a-udp-endpoint)
  - [监控 SCTP 端点](#monitoring-a-sctp-endpoint)
  - [监控 WebSocket 端点](#monitoring-a-websocket-endpoint)
  - [使用 gRPC 监控端点](#monitoring-an-endpoint-using-grpc)
  - [使用 ICMP 监控端点](#monitoring-an-endpoint-using-icmp)
  - [使用 DNS 查询监控端点](#monitoring-an-endpoint-using-dns-queries)
  - [使用 SSH 监控端点](#monitoring-an-endpoint-using-ssh)
  - [使用 STARTTLS 监控端点](#monitoring-an-endpoint-using-starttls)
  - [使用 TLS 监控端点](#monitoring-an-endpoint-using-tls)
  - [监控域名过期](#monitoring-domain-expiration)
  - [并发](#concurrency)
  - [动态重新加载配置](#reloading-configuration-on-the-fly)
  - [端点分组](#endpoint-groups)
  - [如何默认按分组排序？](#how-do-i-sort-by-group-by-default)
  - [在自定义路径暴露 Gatus](#exposing-gatus-on-a-custom-path)
  - [在自定义端口暴露 Gatus](#exposing-gatus-on-a-custom-port)
  - [在配置文件中使用环境变量](#use-environment-variables-in-config-files)
  - [配置启动延迟](#configuring-a-startup-delay)
  - [保持配置文件精简](#keeping-your-configuration-small)
  - [代理客户端配置](#proxy-client-configuration)
  - [如何修复 431 Request Header Fields Too Large 错误](#how-to-fix-431-request-header-fields-too-large-error)
  - [徽章](#badges)
    - [可用率](#uptime)
    - [健康状态](#health)
    - [健康状态 (Shields.io)](#health-shieldsio)
    - [响应时间](#response-time)
    - [响应时间 (图表)](#response-time-chart)
      - [如何更改响应时间徽章的颜色阈值](#how-to-change-the-color-thresholds-of-the-response-time-badge)
  - [API](#api)
    - [以编程方式与 API 交互](#interacting-with-the-api-programmatically)
    - [原始数据](#raw-data)
      - [可用率](#uptime-1)
      - [响应时间](#response-time-1)
  - [作为二进制文件安装](#installing-as-binary)
  - [高层设计概览](#high-level-design-overview)


## 为什么选择 Gatus？
在深入细节之前，我想先回答一个最常见的问题：
> 既然我可以使用 Prometheus 的 Alertmanager、Cloudwatch 甚至 Splunk，为什么还要使用 Gatus？

如果没有客户端主动调用端点，以上任何一种工具都无法告诉你系统存在问题。
换句话说，这是因为监控指标主要依赖于已有流量，这实际上意味着除非
你的客户端已经在遇到问题，否则你不会收到通知。

而 Gatus 允许你为每个功能配置健康检查，从而在任何客户端受到影响之前
监控这些功能并及时向你发出告警。

一个判断你是否需要 Gatus 的信号，就是简单地问自己：如果你的负载均衡器
现在宕机了，你是否会收到告警？你现有的任何告警会被触发吗？如果没有流量能到达你的应用，
你的指标不会报告错误增加。这会让你陷入一种境地：由客户端来通知你服务降级，
而不是你在他们发现问题之前就已经在着手修复。


## 功能特性
Gatus 的主要功能特性：

- **高度灵活的健康检查条件**：虽然检查响应状态对某些场景来说已经足够，但 Gatus 走得更远，允许你对响应时间、响应体甚至 IP 地址添加条件。
- **可用于用户验收测试**：得益于上述特性，你可以利用此应用创建自动化用户验收测试。
- **非常易于配置**：配置不仅设计为尽可能可读，添加新服务或新监控端点也极其简单。
- **告警**：虽然拥有美观的可视化仪表盘对于跟踪应用状态很有用，但你可能不想整天盯着它。因此，开箱即支持通过 Slack、Mattermost、Messagebird、PagerDuty、Twilio、Google Chat 和 Teams 发送通知，并且可以配置自定义告警提供者以满足任何需求，无论是不同的提供者还是管理自动回滚的自定义应用。
- **指标**
- **低资源消耗**：与大多数 Go 应用一样，此应用所需的资源占用微乎其微。
- **[徽章](#badges)**：![可用率 7d](https://status.twin.sh/api/v1/endpoints/core_blog-external/uptimes/7d/badge.svg) ![响应时间 24h](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/24h/badge.svg)
- **深色模式**

![Gatus 仪表盘条件](.github/assets/dashboard-conditions.jpg)


## 使用方法

```console
docker run -p 8080:8080 --name gatus ghcr.io/twin/gatus:stable
```

如果你更喜欢使用 Docker Hub：
```console
docker run -p 8080:8080 --name gatus twinproduction/gatus:stable
```
如果你想创建自己的配置，请参见 [Docker](#docker) 了解如何挂载配置文件。

以下是一个简单的示例：
```yaml
endpoints:
  - name: website                 # Name of your endpoint, can be anything
    url: "https://twin.sh/health"
    interval: 5m                  # Duration to wait between every status check (default: 60s)
    conditions:
      - "[STATUS] == 200"         # Status must be 200
      - "[BODY].status == UP"     # The json path "$.status" must be equal to UP
      - "[RESPONSE_TIME] < 300"   # Response time must be under 300ms

  - name: make-sure-header-is-rendered
    url: "https://example.org/"
    interval: 60s
    conditions:
      - "[STATUS] == 200"                          # Status must be 200
      - "[BODY] == pat(*<h1>Example Domain</h1>*)" # Body must contain the specified header
```

此示例看起来类似于：

![简单示例](.github/assets/example.jpg)

如果你想在本地测试，请参见 [Docker](#docker)。

## 配置
默认情况下，配置文件预期位于 `config/config.yaml`。

你可以通过设置 `GATUS_CONFIG_PATH` 环境变量来指定自定义路径。

如果 `GATUS_CONFIG_PATH` 指向一个目录，该目录及其子目录中的所有 `*.yaml` 和 `*.yml` 文件将按以下方式合并：
- 所有映射/对象会进行深度合并（即你可以在一个文件中定义 `alerting.slack`，在另一个文件中定义 `alerting.pagerduty`）
- 所有切片/数组会追加合并（即你可以在多个文件中定义 `endpoints`，每个端点都会被添加到最终的端点列表中）
- 具有基本类型值的参数（如 `metrics`、`alerting.slack.webhook-url` 等）只能定义一次，以强制避免任何歧义
    - 需要说明的是，这也意味着你不能在两个文件中为 `alerting.slack.webhook-url` 定义不同的值。所有文件在处理前会先合并为一个。这是设计如此。

> 💡 你也可以在配置文件中使用环境变量（例如 `$DOMAIN`、`${DOMAIN}`）
>
> ⚠️ 当你的配置参数包含 `$` 符号时，你需要使用 `$$` 来转义 `$`。
>
> 参见[在配置文件中使用环境变量](#use-environment-variables-in-config-files)或 [examples/docker-compose-postgres-storage/config/config.yaml](.examples/docker-compose-postgres-storage/config/config.yaml) 获取示例。

如果你想在本地测试，请参见 [Docker](#docker)。


## 配置
| 参数                           | 描述                                                                                                                                       | 默认值        |
|:-----------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------|:--------------|
| `metrics`                    | 是否在 `/metrics` 暴露指标。                                                                                                               | `false`       |
| `storage`                    | [存储配置](#storage)。                                                                                                                      | `{}`          |
| `alerting`                   | [告警配置](#alerting)。                                                                                                                     | `{}`          |
| `announcements`              | [公告配置](#announcements)。                                                                                                                | `[]`          |
| `endpoints`                  | [端点配置](#endpoints)。                                                                                                                    | Required `[]` |
| `external-endpoints`         | [外部端点配置](#external-endpoints)。                                                                                                       | `[]`          |
| `security`                   | [安全配置](#security)。                                                                                                                     | `{}`          |
| `concurrency`                | 最大并发监控端点/套件数量。设置为 `0` 表示无限制。参见[并发](#concurrency)。                                                                  | `3`           |
| `disable-monitoring-lock`    | 是否[禁用监控锁](#disable-monitoring-lock)。**已弃用**：请改用 `concurrency: 0`。                                                            | `false`       |
| `skip-invalid-config-update` | 是否忽略无效的配置更新。<br />参见[动态重新加载配置](#reloading-configuration-on-the-fly)。                                                    | `false`       |
| `web`                        | [Web 配置](#web)。                                                                                                                         | `{}`          |
| `ui`                         | [UI 配置](#ui)。                                                                                                                           | `{}`          |
| `maintenance`                | [维护配置](#maintenance)。                                                                                                                  | `{}`          |

如果你需要更详细的日志，可以将 `GATUS_LOG_LEVEL` 环境变量设置为 `DEBUG`。
相反，如果你需要更简洁的日志，可以将上述环境变量设置为 `WARN`、`ERROR` 或 `FATAL`。
`GATUS_LOG_LEVEL` 的默认值为 `INFO`。

### 端点
端点是你想要监控的 URL、应用或服务。每个端点都有一组条件，
这些条件按你定义的时间间隔进行评估。如果任何条件失败，该端点将被视为不健康。
然后你可以配置告警，在端点不健康达到一定阈值时触发。

| 参数                                              | 描述                                                                                                                                          | 默认值                     |
|:------------------------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------|:---------------------------|
| `endpoints`                                     | 要监控的端点列表。                                                                                                                              | Required `[]`              |
| `endpoints[].enabled`                           | 是否监控该端点。                                                                                                                                | `true`                     |
| `endpoints[].name`                              | 端点名称。可以是任意值。                                                                                                                         | Required `""`              |
| `endpoints[].group`                             | 分组名称。用于在仪表盘上将多个端点分组显示。<br />参见[端点分组](#endpoint-groups)。                                                                | `""`                       |
| `endpoints[].url`                               | 发送请求的 URL。                                                                                                                                | Required `""`              |
| `endpoints[].method`                            | 请求方法。                                                                                                                                      | `GET`                      |
| `endpoints[].conditions`                        | 用于判定端点健康状态的条件。<br />参见[条件](#conditions)。                                                                                       | `[]`                       |
| `endpoints[].interval`                          | 每次状态检查之间的等待时间。                                                                                                                      | `60s`                      |
| `endpoints[].graphql`                           | 是否将请求体包装在 query 参数中（`{"query":"$body"}`）。                                                                                          | `false`                    |
| `endpoints[].body`                              | 请求体。                                                                                                                                        | `""`                       |
| `endpoints[].headers`                           | 请求头。                                                                                                                                        | `{}`                       |
| `endpoints[].dns`                               | DNS 类型端点的配置。<br />参见[使用 DNS 查询监控端点](#monitoring-an-endpoint-using-dns-queries)。                                                  | `""`                       |
| `endpoints[].dns.query-type`                    | 查询类型（例如 MX）。                                                                                                                            | `""`                       |
| `endpoints[].dns.query-name`                    | 查询名称（例如 example.com）。                                                                                                                   | `""`                       |
| `endpoints[].ssh`                               | SSH 类型端点的配置。<br />参见[使用 SSH 监控端点](#monitoring-an-endpoint-using-ssh)。                                                              | `""`                       |
| `endpoints[].ssh.username`                      | SSH 用户名（例如 example）。                                                                                                                     | Required `""`              |
| `endpoints[].ssh.password`                      | SSH 密码（例如 password）。                                                                                                                      | Required `""`              |
| `endpoints[].alerts`                            | 给定端点的所有告警列表。<br />参见[告警](#alerting)。                                                                                              | `[]`                       |
| `endpoints[].maintenance-windows`               | 给定端点的所有维护窗口列表。<br />参见[维护](#maintenance)。                                                                                       | `[]`                       |
| `endpoints[].client`                            | [客户端配置](#client-configuration)。                                                                                                           | `{}`                       |
| `endpoints[].tamper`                            | 基于响应体大小漂移的防篡改配置。                                                                                                                    | `{}`                       |
| `endpoints[].tamper.enabled`                    | 是否启用响应体大小漂移防篡改检测。                                                                                                                     | `false`                    |
| `endpoints[].tamper.baseline-samples`           | 用于建立响应体大小基线的最近样本数。                                                                                                                      | `20`                       |
| `endpoints[].tamper.drift-threshold-percent`    | 判定为一次漂移触发前允许的漂移百分比。                                                                                                                    | `20`                       |
| `endpoints[].tamper.consecutive-breaches`       | 将检查判定为不健康前需要连续触发的次数。                                                                                                                      | `3`                        |
| `endpoints[].ui`                                | 端点级别的 UI 配置。                                                                                                                             | `{}`                       |
| `endpoints[].ui.hide-conditions`                | 是否在结果中隐藏条件。注意这只会隐藏启用此选项后评估的条件。                                                                                        | `false`                    |
| `endpoints[].ui.hide-hostname`                  | 是否在结果中隐藏主机名。                                                                                                                         | `false`                    |
| `endpoints[].ui.hide-port`                      | 是否在结果中隐藏端口。                                                                                                                           | `false`                    |
| `endpoints[].ui.hide-url`                       | 是否在结果中隐藏 URL。当 URL 包含令牌时很有用。                                                                                                    | `false`                    |
| `endpoints[].ui.hide-errors`                    | 是否在结果中隐藏错误。                                                                                                                           | `false`                    |
| `endpoints[].ui.dont-resolve-failed-conditions` | 是否在 UI 中解析失败的条件。                                                                                                                     | `false`                    |
| `endpoints[].ui.resolve-successful-conditions`  | 是否在 UI 中解析成功的条件（有助于在检查通过时也展示响应体断言）。                                                                                   | `false`                    |
| `endpoints[].ui.badge.response-time`            | 响应时间阈值列表。每当达到一个阈值时，徽章会显示不同的颜色。                                                                                        | `[50, 200, 300, 500, 750]` |
| `endpoints[].extra-labels`                      | 添加到指标的额外标签。用于将端点分组。                                                                                                             | `{}`                       |
| `endpoints[].always-run`                        | （仅限套件）即使套件中之前的端点失败，是否仍执行此端点。                                                                                            | `false`                    |
| `endpoints[].store`                             | （仅限套件）从响应中提取并存储到套件上下文中的值映射（即使失败也会存储）。                                                                            | `{}`                       |

你可以在请求体（`endpoints[].body`）中使用以下占位符：
- `[ENDPOINT_NAME]`（从 `endpoints[].name` 解析）
- `[ENDPOINT_GROUP]`（从 `endpoints[].group` 解析）
- `[ENDPOINT_URL]`（从 `endpoints[].url` 解析）
- `[LOCAL_ADDRESS]`（解析为本地 IP 和端口，如 `192.0.2.1:25` 或 `[2001:db8::1]:80`）
- `[RANDOM_STRING_N]`（解析为长度为 N 的随机字母数字字符串（最大：8192））

### 外部端点
与常规端点不同，外部端点不由 Gatus 监控，而是通过编程方式推送状态。
这允许你监控任何你想要的内容，即使你要检查的内容位于 Gatus 通常无法访问的环境中。

例如：
- 你可以创建自己的代理，驻留在私有网络中，将服务状态推送到公开暴露的 Gatus 实例
- 你可以监控 Gatus 不支持的服务
- 你可以实现自己的监控系统，同时使用 Gatus 作为仪表盘

| 参数                                        | 描述                                                                                                                        | 默认值         |
|:------------------------------------------|:----------------------------------------------------------------------------------------------------------------------------------|:---------------|
| `external-endpoints`                      | 要监控的端点列表。                                                                                                               | `[]`           |
| `external-endpoints[].enabled`            | 是否监控该端点。                                                                                                                 | `true`         |
| `external-endpoints[].name`               | 端点名称。可以是任意值。                                                                                                          | Required `""`  |
| `external-endpoints[].group`              | 分组名称。用于在仪表盘上将多个端点分组显示。<br />参见[端点分组](#endpoint-groups)。                                                  | `""`           |
| `external-endpoints[].token`              | 推送状态所需的 Bearer 令牌。                                                                                                      | Required `""`  |
| `external-endpoints[].alerts`             | 给定端点的所有告警列表。<br />参见[告警](#alerting)。                                                                               | `[]`           |
| `external-endpoints[].heartbeat`          | 心跳配置，用于监控外部端点何时停止发送更新。                                                                                        | `{}`           |
| `external-endpoints[].heartbeat.interval` | 预期的更新间隔。如果在此间隔内未收到更新，将触发告警。最小值为 10s。                                                                  | `0`（已禁用）  |

示例：
```yaml
external-endpoints:
  - name: ext-ep-test
    group: core
    token: "potato"
    heartbeat:
      interval: 30m  # Automatically create a failure if no update is received within 30 minutes
    alerts:
      - type: discord
        description: "healthcheck failed"
        send-on-resolved: true
```

要推送外部端点的状态，你可以使用 [gatus-cli](https://github.com/TwiN/gatus-cli)：
```
gatus-cli external-endpoint push --url https://status.example.org --key "core_ext-ep-test" --token "potato" --success
```

或发送 HTTP 请求：
```
POST /api/v1/endpoints/{key}/external?success={success}&error={error}&duration={duration}
```
其中：
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都会被替换为 `-`。
  - 使用上面的示例配置，key 将是 `core_ext-ep-test`。
- `{success}` 是一个布尔值（`true` 或 `false`），表示健康检查是否成功。
- `{error}`（可选）：描述健康检查失败原因的字符串。如果 {success} 为 false，这应该包含错误消息；如果检查成功，此值将被忽略。
- `{duration}`（可选）：请求耗时，格式为持续时间字符串（例如 10s）。

你还必须在 `Authorization` 请求头中以 `Bearer` 令牌的形式传递令牌。


### 套件 (ALPHA)
套件是按顺序执行的端点集合，共享一个上下文。
这允许你创建复杂的监控场景，其中一个端点的结果可以在后续端点中使用，从而实现工作流式监控。

以下是一些套件可能有用的场景：
- 测试多步骤身份验证流程（登录 -> 访问受保护资源 -> 登出）
- 需要链式请求的 API 工作流（创建资源 -> 更新 -> 验证 -> 删除）
- 监控跨多个服务的业务流程
- 验证多个端点之间的数据一致性

| 参数                                | 描述                                                                                              | 默认值        |
|:----------------------------------|:----------------------------------------------------------------------------------------------------|:--------------|
| `suites`                          | 要监控的套件列表。                                                                                    | `[]`          |
| `suites[].enabled`                | 是否监控该套件。                                                                                      | `true`        |
| `suites[].name`                   | 套件名称。必须唯一。                                                                                   | Required `""` |
| `suites[].group`                  | 分组名称。用于在仪表盘上将多个套件分组显示。                                                              | `""`          |
| `suites[].interval`               | 套件执行之间的等待时间。                                                                                | `10m`         |
| `suites[].timeout`                | 整个套件执行的最大持续时间。                                                                             | `5m`          |
| `suites[].context`                | 可被端点引用的初始上下文值。                                                                             | `{}`          |
| `suites[].ui`                     | 套件中所有端点的 UI 配置默认值（与 `endpoints[].ui` 相同的字段）。                                        | `{}`          |
| `suites[].endpoints`              | 要按顺序执行的端点列表。                                                                                | Required `[]` |
| `suites[].endpoints[].store`      | 从响应中提取并存储到套件上下文中的值映射（即使失败也会存储）。                                              | `{}`          |
| `suites[].endpoints[].always-run` | 即使套件中之前的端点失败，是否仍执行此端点。                                                              | `false`       |

**注意**：套件级别的告警尚不支持。请在套件内的各个端点上配置告警。

#### 在端点中使用上下文
一旦值存储在上下文中，就可以在后续端点中引用它们：
- 在 URL 中：`https://api.example.com/users/[CONTEXT].user_id`
- 在请求头中：`Authorization: Bearer [CONTEXT].auth_token`
- 在请求体中：`{"user_id": "[CONTEXT].user_id"}`
- 在条件中：`[BODY].server_ip == [CONTEXT].server_ip`

注意，上下文/存储的键仅限于 A-Z、a-z、0-9、下划线（`_`）和连字符（`-`）。

#### 套件配置示例
```yaml
suites:
  - name: item-crud-workflow
    group: api-tests
    interval: 5m
    context:
      price: "19.99"  # Initial static value in context
    endpoints:
      # Step 1: Create an item and store the item ID
      - name: create-item
        url: https://api.example.com/items
        method: POST
        body: '{"name": "Test Item", "price": "[CONTEXT].price"}'
        conditions:
          - "[STATUS] == 201"
          - "len([BODY].id) > 0"
          - "[BODY].price == [CONTEXT].price"
        store:
          itemId: "[BODY].id"
        alerts:
          - type: slack
            description: "Failed to create item"

      # Step 2: Update the item using the stored item ID
      - name: update-item
        url: https://api.example.com/items/[CONTEXT].itemId
        method: PUT
        body: '{"price": "24.99"}'
        conditions:
          - "[STATUS] == 200"
        alerts:
          - type: slack
            description: "Failed to update item"

      # Step 3: Fetch the item and validate the price
      - name: get-item
        url: https://api.example.com/items/[CONTEXT].itemId
        method: GET
        conditions:
          - "[STATUS] == 200"
          - "[BODY].price == 24.99"
        alerts:
          - type: slack
            description: "Item price did not update correctly"

      # Step 4: Delete the item (always-run: true to ensure cleanup even if step 2 or 3 fails)
      - name: delete-item
        url: https://api.example.com/items/[CONTEXT].itemId
        method: DELETE
        always-run: true
        conditions:
          - "[STATUS] == 204"
        alerts:
          - type: slack
            description: "Failed to delete item"
```

只有当所有必需端点都通过其条件时，套件才被视为成功。


### 条件
以下是你可以使用的一些条件示例：

| 条件                               | 描述                                                  | 通过的值                     | 失败的值         |
|:---------------------------------|:----------------------------------------------------|:---------------------------|------------------|
| `[STATUS] == 200`                | 状态码必须等于 200                                     | 200                        | 201, 404, ...    |
| `[STATUS] < 300`                 | 状态码必须小于 300                                     | 200, 201, 299              | 301, 302, ...    |
| `[STATUS] <= 299`                | 状态码必须小于或等于 299                                | 200, 201, 299              | 301, 302, ...    |
| `[STATUS] > 400`                 | 状态码必须大于 400                                     | 401, 402, 403, 404         | 400, 200, ...    |
| `[STATUS] == any(200, 429)`      | 状态码必须为 200 或 429                                | 200, 429                   | 201, 400, ...    |
| `[CONNECTED] == true`            | 必须成功连接到主机                                     | true                       | false            |
| `[RESPONSE_TIME] < 500`          | 响应时间必须低于 500ms                                 | 100ms, 200ms, 300ms        | 500ms, 501ms     |
| `[IP] == 127.0.0.1`              | 目标 IP 必须为 127.0.0.1                               | 127.0.0.1                  | 0.0.0.0          |
| `[BODY] == 1`                    | 响应体必须等于 1                                       | 1                          | `{}`, `2`, ...   |
| `[BODY].user.name == john`       | JSONPath `$.user.name` 的值等于 `john`                 | `{"user":{"name":"john"}}` |                  |
| `[BODY].data[0].id == 1`         | JSONPath `$.data[0].id` 的值等于 1                     | `{"data":[{"id":1}]}`      |                  |
| `[BODY].age == [BODY].id`        | JSONPath `$.age` 的值等于 JSONPath `$.id`              | `{"age":1,"id":1}`         |                  |
| `len([BODY].data) < 5`           | JSONPath `$.data` 的数组元素少于 5 个                   | `{"data":[{"id":1}]}`      |                  |
| `len([BODY].name) == 8`          | JSONPath `$.name` 的字符串长度为 8                      | `{"name":"john.doe"}`      | `{"name":"bob"}` |
| `has([BODY].errors) == false`    | JSONPath `$.errors` 不存在                             | `{"name":"john.doe"}`      | `{"errors":[]}`  |
| `has([BODY].users) == true`      | JSONPath `$.users` 存在                                | `{"users":[]}`             | `{}`             |
| `[BODY].name == pat(john*)`      | JSONPath `$.name` 的字符串匹配模式 `john*`              | `{"name":"john.doe"}`      | `{"name":"bob"}` |
| `[BODY].id == any(1, 2)`         | JSONPath `$.id` 的值等于 `1` 或 `2`                    | 1, 2                       | 3, 4, 5          |
| `[CERTIFICATE_EXPIRATION] > 48h` | 证书过期时间距现在超过 48 小时                           | 49h, 50h, 123h             | 1h, 24h, ...     |
| `[DOMAIN_EXPIRATION] > 720h`     | 域名过期时间必须超过 720 小时                            | 4000h                      | 1h, 24h, ...     |


#### 占位符
| 占位符                       | 描述                                                                                        | 解析值示例                                   |
|:---------------------------|:------------------------------------------------------------------------------------------|:---------------------------------------------|
| `[STATUS]`                 | 解析为请求的 HTTP 状态码                                                                     | `404`                                        |
| `[RESPONSE_TIME]`          | 解析为请求所花费的响应时间，单位为毫秒                                                        | `10`                                         |
| `[IP]`                     | 解析为目标主机的 IP 地址                                                                     | `192.168.0.232`                              |
| `[BODY]`                   | 解析为响应体。支持 JSONPath。                                                                 | `{"name":"john.doe"}`                        |
| `[CONNECTED]`              | 解析为是否能够建立连接                                                                        | `true`                                       |
| `[CERTIFICATE_EXPIRATION]` | 解析为证书过期前的持续时间（有效单位为 "s"、"m"、"h"）                                          | `24h`、`48h`、0（如果协议不支持证书）          |
| `[DOMAIN_EXPIRATION]`      | 解析为域名过期前的持续时间（有效单位为 "s"、"m"、"h"）                                          | `24h`、`48h`、`1234h56m78s`                  |
| `[DNS_RCODE]`              | 解析为响应的 DNS 状态码                                                                       | `NOERROR`                                    |


#### 函数
| 函数     | 描述                                                                                                                                                                                                                                  | 示例                               |
|:---------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:-----------------------------------|
| `len`    | 如果给定路径指向一个数组，返回其长度。否则，将给定路径的 JSON 压缩并转换为字符串，返回结果的字符数。仅适用于 `[BODY]` 占位符。                                                                                                               | `len([BODY].username) > 8`         |
| `has`    | 根据给定路径是否有效，返回 `true` 或 `false`。仅适用于 `[BODY]` 占位符。                                                                                                                                                                | `has([BODY].errors) == false`      |
| `pat`    | 指定作为参数传递的字符串应被评估为模式。仅适用于 `==` 和 `!=`。                                                                                                                                                                          | `[IP] == pat(192.168.*)`           |
| `any`    | 指定作为参数传递的任何一个值都是有效值。仅适用于 `==` 和 `!=`。                                                                                                                                                                          | `[BODY].ip == any(127.0.0.1, ::1)` |

> 💡 仅在需要时使用 `pat`。`[STATUS] == pat(2*)` 的开销比 `[STATUS] < 300` 大得多。

### Web 配置
允许你配置仪表盘的服务方式和位置。

| 参数                         | 描述                                                                                          | 默认值    |
|:---------------------------|:--------------------------------------------------------------------------------------------|:----------|
| `web`                      | Web 配置                                                                                      | `{}`      |
| `web.address`              | 监听地址。                                                                                     | `0.0.0.0` |
| `web.port`                 | 监听端口。                                                                                     | `8080`    |
| `web.read-buffer-size`     | 从连接读取请求的缓冲区大小。同时也是最大请求头大小的限制。                                         | `8192`    |
| `web.tls.certificate-file` | 可选的 PEM 格式 TLS 公钥证书文件。                                                              | `""`      |
| `web.tls.private-key-file` | 可选的 PEM 格式 TLS 私钥文件。                                                                  | `""`      |

### UI 配置
允许你配置仪表盘 UI 的应用级默认设置。其中一些参数可以被用户通过浏览器的本地存储在本地覆盖。

| 参数                          | 描述                                                                                                                                       | 默认值                                              |
|:--------------------------|:-----------------------------------------------------------------------------------------------------------------------------------------|:----------------------------------------------------|
| `ui`                      | UI 配置                                                                                                                                    | `{}`                                                |
| `ui.title`                | [文档标题](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title)。                                                               | `Health Dashboard ǀ Gatus`                          |
| `ui.description`          | 页面的 meta 描述。                                                                                                                          | `Gatus is an advanced...`                           |
| `ui.dashboard-heading`    | 仪表盘标题，位于页眉和端点之间                                                                                                               | `Health Dashboard`                                  |
| `ui.dashboard-subheading` | 仪表盘描述，位于页眉和端点之间                                                                                                               | `Monitor the health of your endpoints in real-time` |
| `ui.header`               | 仪表盘顶部的页眉。                                                                                                                          | `Gatus`                                             |
| `ui.logo`                 | 要显示的 logo 的 URL。                                                                                                                      | `""`                                                |
| `ui.link`                 | 点击 logo 时打开的链接。                                                                                                                     | `""`                                                |
| `ui.favicon.default`      | 在浏览器标签页或地址栏中显示的默认收藏图标。                                                                                                   | `/favicon.ico`                                      |
| `ui.favicon.size16x16`    | 在浏览器中显示的 16x16 尺寸收藏图标。                                                                                                         | `/favicon-16x16.png`                                |
| `ui.favicon.size32x32`    | 在浏览器中显示的 32x32 尺寸收藏图标。                                                                                                         | `/favicon-32x32.png`                                |
| `ui.buttons`              | 显示在页眉下方的按钮列表。                                                                                                                    | `[]`                                                |
| `ui.buttons[].name`       | 按钮上显示的文本。                                                                                                                           | Required `""`                                       |
| `ui.buttons[].link`       | 点击按钮时打开的链接。                                                                                                                        | Required `""`                                       |
| `ui.custom-css`           | 自定义 CSS                                                                                                                                  | `""`                                                |
| `ui.dark-mode`            | 是否默认启用深色模式。注意此设置会被用户操作系统的主题偏好所覆盖。                                                                                | `true`                                              |
| `ui.default-sort-by`      | 仪表盘中端点的默认排序方式。可选值为 `name`、`group` 或 `health`。注意用户偏好会覆盖此设置。                                                      | `name`                                              |
| `ui.default-filter-by`    | 仪表盘中端点的默认筛选方式。可选值为 `none`、`failing` 或 `unstable`。注意用户偏好会覆盖此设置。                                                   | `none`                                              |

### 公告
系统级公告允许你在状态页面顶部显示重要消息。这些公告可用于通知用户计划维护、正在进行的问题或一般信息。你可以使用 markdown 来格式化公告内容。

这本质上就是一些状态页面所称的"事件通信"。

| 参数                          | 描述                                                                                                                       | 默认值   |
|:----------------------------|:-------------------------------------------------------------------------------------------------------------------------|:---------|
| `announcements`             | 要显示的公告列表                                                                                                             | `[]`     |
| `announcements[].timestamp` | 公告发布时的 UTC 时间戳（RFC3339 格式）                                                                                       | Required |
| `announcements[].type`      | 公告类型。有效值：`outage`、`warning`、`information`、`operational`、`none`                                                    | `"none"` |
| `announcements[].message`   | 显示给用户的消息                                                                                                              | Required |
| `announcements[].archived`  | 是否归档该公告。已归档的公告会显示在状态页面底部而非顶部。                                                                        | `false`  |

类型说明：
- **outage**：表示服务中断或严重问题（红色主题）
- **warning**：表示潜在问题或重要通知（黄色主题）
- **information**：一般信息或更新（蓝色主题）
- **operational**：表示已解决的问题或正常运行（绿色主题）
- **none**：无特定严重级别的中性公告（灰色主题，未指定时的默认值）

配置示例：
```yaml
announcements:
  - timestamp: 2025-11-07T14:00:00Z
    type: outage
    message: "Scheduled maintenance on database servers from 14:00 to 16:00 UTC"
  - timestamp: 2025-11-07T16:15:00Z
    type: operational
    message: "Database maintenance completed successfully. All systems operational."
  - timestamp: 2025-11-07T12:00:00Z
    type: information
    message: "New monitoring dashboard features will be deployed next week"
  - timestamp: 2025-11-06T09:00:00Z
    type: warning
    message: "Elevated API response times observed for US customers"
    archived: true
```

如果至少有一个公告被归档，状态页面底部将渲染一个**历史公告**部分：
![Gatus 历史公告部分](.github/assets/past-announcements.jpg)


### 存储
| 参数                                  | 描述                                                                                                                                         | 默认值     |
|:------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------------|:-----------|
| `storage`                           | 存储配置                                                                                                                                        | `{}`       |
| `storage.path`                      | 持久化数据的路径。仅支持 `sqlite` 和 `postgres` 类型。                                                                     | `""`       |
| `storage.type`                      | 存储类型。有效类型：`memory`、`sqlite`、`postgres`。                                                                                      | `"memory"` |
| `storage.caching`                   | 是否使用写穿透缓存。可改善大型仪表盘的加载时间。<br />仅在 `storage.type` 为 `sqlite` 或 `postgres` 时支持 | `false`    |
| `storage.maximum-number-of-results` | 端点可保存的最大结果数量                                                                                            | `100`      |
| `storage.maximum-number-of-events`  | 端点可保存的最大事件数量                                                                                             | `50`       |

每个端点健康检查的结果以及正常运行时间和历史事件的数据必须被持久化，
以便在仪表盘上显示。这些参数允许你配置相关的存储。

- 如果 `storage.type` 为 `memory`（默认值）：
```yaml
# 请注意这是默认值，你可以完全省略存储配置以获得相同的效果。
# 因为数据存储在内存中，数据在重启后不会保留。
storage:
  type: memory
  maximum-number-of-results: 200
  maximum-number-of-events: 5
```
- 如果 `storage.type` 为 `sqlite`，`storage.path` 不能为空：
```yaml
storage:
  type: sqlite
  path: data.db
```
参见 [examples/docker-compose-sqlite-storage](.examples/docker-compose-sqlite-storage) 示例。

- 如果 `storage.type` 为 `postgres`，`storage.path` 必须为连接 URL：
```yaml
storage:
  type: postgres
  path: "postgres://user:password@127.0.0.1:5432/gatus?sslmode=disable"
```
参见 [examples/docker-compose-postgres-storage](.examples/docker-compose-postgres-storage) 示例。


### 客户端配置
为了支持各种不同的环境，每个受监控的端点都有一个独立的客户端配置，用于发送请求。

| 参数                              | 描述                                                                   | 默认值         |
|:---------------------------------------|:------------------------------------------------------------------------------|:----------------|
| `client.insecure`                      | 是否跳过验证服务器的证书链和主机名。       | `false`         |
| `client.ignore-redirect`               | 是否忽略重定向（true）或跟随重定向（false，默认）。           | `false`         |
| `client.timeout`                       | 超时时间。                                                   | `10s`           |
| `client.dns-resolver`                  | 使用 `{proto}://{host}:{port}` 格式覆盖 DNS 解析器。         | `""`            |
| `client.oauth2`                        | OAuth2 客户端配置。                                                  | `{}`            |
| `client.oauth2.token-url`              | Token 端点 URL                                                        | required `""`   |
| `client.oauth2.client-id`              | 用于 `Client credentials flow` 的客户端 ID          | required `""`   |
| `client.oauth2.client-secret`          | 用于 `Client credentials flow` 的客户端密钥      | required `""`   |
| `client.oauth2.scopes[]`               | 用于 `Client credentials flow` 的 `scopes` 列表。    | required `[""]` |
| `client.proxy-url`                     | 客户端使用的代理 URL                                    | `""`            |
| `client.identity-aware-proxy`          | Google Identity-Aware-Proxy 客户端配置。                             | `{}`            |
| `client.identity-aware-proxy.audience` | Identity-Aware-Proxy 的 audience。（IAP oauth2 凭据的 client-id）   | required `""`   |
| `client.tls.certificate-file`          | 用于 mTLS 配置的客户端证书路径（PEM 格式）。         | `""`            |
| `client.tls.private-key-file`          | 用于 mTLS 配置的客户端私钥路径（PEM 格式）。         | `""`            |
| `client.tls.renegotiation`             | 提供的重新协商支持类型。（`never`、`freely`、`once`）。        | `"never"`       |
| `client.network`                       | 用于 ICMP 端点客户端的网络类型（`ip`、`ip4` 或 `ip6`）。           | `"ip"`          |
| `client.tunnel`                        | 用于此端点的 SSH 隧道名称。参见[隧道](#tunneling)。 | `""`            |


> 📝 其中一些参数会根据端点类型被忽略。例如，ICMP 请求（ping）不涉及证书，
> 因此将该类型端点的 `client.insecure` 设置为 `true` 不会产生任何效果。

默认配置如下：

```yaml
client:
  insecure: false
  ignore-redirect: false
  timeout: 10s
```

请注意，此配置仅在 `endpoints[]`、`alerting.mattermost` 和 `alerting.custom` 下可用。

以下是 `endpoints[]` 下客户端配置的示例：

```yaml
endpoints:
  - name: website
    url: "https://twin.sh/health"
    client:
      insecure: false
      ignore-redirect: false
      timeout: 10s
    conditions:
      - "[STATUS] == 200"
```

此示例展示如何指定自定义 DNS 解析器：

```yaml
endpoints:
  - name: with-custom-dns-resolver
    url: "https://your.health.api/health"
    client:
      dns-resolver: "tcp://8.8.8.8:53"
    conditions:
      - "[STATUS] == 200"
```

此示例展示如何使用 `client.oauth2` 配置通过 `Bearer token` 查询后端 API：

```yaml
endpoints:
  - name: with-custom-oauth2
    url: "https://your.health.api/health"
    client:
      oauth2:
        token-url: https://your-token-server/token
        client-id: 00000000-0000-0000-0000-000000000000
        client-secret: your-client-secret
        scopes: ['https://your.health.api/.default']
    conditions:
      - "[STATUS] == 200"
```

此示例展示如何使用 `client.identity-aware-proxy` 配置通过 Google Identity-Aware-Proxy 以 `Bearer token` 查询后端 API：

```yaml
endpoints:
  - name: with-custom-iap
    url: "https://my.iap.protected.app/health"
    client:
      identity-aware-proxy:
        audience: "XXXXXXXX-XXXXXXXXXXXX.apps.googleusercontent.com"
    conditions:
      - "[STATUS] == 200"
```

> 📝 请注意，Gatus 将使用其运行环境中的 [gcloud 默认凭据](https://cloud.google.com/docs/authentication/application-default-credentials) 来生成令牌。

此示例展示如何使用 `client.tls` 配置对后端 API 执行 mTLS 查询：

```yaml
endpoints:
  - name: website
    url: "https://your.mtls.protected.app/health"
    client:
      tls:
        certificate-file: /path/to/user_cert.pem
        private-key-file: /path/to/user_key.pem
        renegotiation: once
    conditions:
      - "[STATUS] == 200"
```

> 📝 请注意，如果在容器中运行，你必须将证书和密钥通过卷挂载到容器中。

### 隧道
Gatus 支持 SSH 隧道，可通过跳板机或堡垒服务器监控内部服务。
这在监控从 Gatus 部署位置无法直接访问的服务时特别有用。

SSH 隧道在 `tunneling` 部分全局定义，然后在端点客户端配置中通过名称引用。

| 参数                             | 描述                                                 | 默认值       |
|:--------------------------------------|:------------------------------------------------------------|:--------------|
| `tunneling`                           | SSH 隧道配置                                   | `{}`          |
| `tunneling.<tunnel-name>`             | 命名 SSH 隧道的配置                        | `{}`          |
| `tunneling.<tunnel-name>.type`        | 隧道类型（目前仅支持 `SSH`）          | Required `""` |
| `tunneling.<tunnel-name>.host`        | SSH 服务器主机名或 IP 地址                           | Required `""` |
| `tunneling.<tunnel-name>.port`        | SSH 服务器端口                                             | `22`          |
| `tunneling.<tunnel-name>.username`    | SSH 用户名                                                | Required `""` |
| `tunneling.<tunnel-name>.password`    | SSH 密码（与 private-key 二选一使用）               | `""`          |
| `tunneling.<tunnel-name>.private-key` | PEM 格式的 SSH 私钥（与 password 二选一使用） | `""`          |
| `client.tunnel`                       | 用于此端点的隧道名称                 | `""`          |

```yaml
tunneling:
  production:
    type: SSH
    host: "jumphost.example.com"
    username: "monitoring"
    private-key: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEpAIBAAKCAQEA...
      -----END RSA PRIVATE KEY-----

endpoints:
  - name: "internal-api"
    url: "http://internal-api.example.com:8080/health"
    client:
      tunnel: "production"
    conditions:
      - "[STATUS] == 200"
```

> ⚠️ **警告**：隧道可能会引入额外的延迟，特别是在隧道连接频繁重试的情况下。
> 这可能导致响应时间测量不准确。


### 告警
Gatus 支持多种告警提供商，如 Slack 和 PagerDuty，并支持为每个端点配置不同的告警，
具有可配置的描述和阈值。

告警在端点级别进行配置，如下所示：

| 参数                            | 描述                                                                                                                                               | 默认值       |
|:-------------------------------------|:----------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------|
| `alerts`                             | 给定端点的所有告警列表。                                                                                                                  | `[]`          |
| `alerts[].type`                      | 告警类型。<br />有效类型请参见下表。                                                                                                 | Required `""` |
| `alerts[].enabled`                   | 是否启用该告警。                                                                                                                              | `true`        |
| `alerts[].failure-threshold`         | 触发告警所需的连续失败次数。                                                                                           | `3`           |
| `alerts[].success-threshold`         | 将正在进行的事件标记为已解决所需的连续成功次数。                                                                            | `2`           |
| `alerts[].minimum-reminder-interval` | 告警提醒之间的最小时间间隔。例如 `"30m"`、`"1h45m30s"` 或 `"24h"`。如果为空或 `0`，则禁用提醒。不能低于 `5m`。 | `0`           |
| `alerts[].send-on-resolved`          | 当触发的告警被标记为已解决时，是否发送通知。                                                                              | `false`       |
| `alerts[].description`               | 告警描述。将包含在发送的告警中。                                                                                             | `""`          |
| `alerts[].provider-override`         | 针对给定告警类型的告警提供商配置覆盖                                                                                         | `{}`          |

以下是端点级别告警配置的示例：
```yaml
endpoints:
  - name: example
    url: "https://example.org"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack
        description: "healthcheck failed"
        send-on-resolved: true
```

你还可以使用 `alerts[].provider-override` 覆盖全局提供商配置，如下所示：
```yaml
endpoints:
  - name: example
    url: "https://example.org"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack
        provider-override:
          webhook-url: "https://hooks.slack.com/services/**********/**********/**********"
```

> 📝 如果告警提供商未正确配置，所有使用该提供商类型配置的告警将被忽略。

| 参数                  | 描述                                                                                                                             | 默认值 |
|:---------------------------|:----------------------------------------------------------------------------------------------------------------------------------------|:--------|
| `alerting.awsses`          | `awsses` 类型告警的配置。<br />参见[配置 AWS SES 告警](#configuring-aws-ses-alerts)。                         | `{}`    |
| `alerting.clickup`         | `clickup` 类型告警的配置。<br />参见[配置 ClickUp 告警](#configuring-clickup-alerts)。                        | `{}`    |
| `alerting.custom`          | 失败或告警时自定义操作的配置。<br />参见[配置自定义告警](#configuring-custom-alerts)。               | `{}`    |
| `alerting.datadog`         | `datadog` 类型告警的配置。<br />参见[配置 Datadog 告警](#configuring-datadog-alerts)。                        | `{}`    |
| `alerting.discord`         | `discord` 类型告警的配置。<br />参见[配置 Discord 告警](#configuring-discord-alerts)。                        | `{}`    |
| `alerting.email`           | `email` 类型告警的配置。<br />参见[配置 Email 告警](#configuring-email-alerts)。                              | `{}`    |
| `alerting.gitea`           | `gitea` 类型告警的配置。<br />参见[配置 Gitea 告警](#configuring-gitea-alerts)。                              | `{}`    |
| `alerting.github`          | `github` 类型告警的配置。<br />参见[配置 GitHub 告警](#configuring-github-alerts)。                           | `{}`    |
| `alerting.gitlab`          | `gitlab` 类型告警的配置。<br />参见[配置 GitLab 告警](#configuring-gitlab-alerts)。                           | `{}`    |
| `alerting.googlechat`      | `googlechat` 类型告警的配置。<br />参见[配置 Google Chat 告警](#configuring-google-chat-alerts)。             | `{}`    |
| `alerting.gotify`          | `gotify` 类型告警的配置。<br />参见[配置 Gotify 告警](#configuring-gotify-alerts)。                           | `{}`    |
| `alerting.homeassistant`   | `homeassistant` 类型告警的配置。<br />参见[配置 HomeAssistant 告警](#configuring-homeassistant-alerts)。      | `{}`    |
| `alerting.ifttt`           | `ifttt` 类型告警的配置。<br />参见[配置 IFTTT 告警](#configuring-ifttt-alerts)。                              | `{}`    |
| `alerting.ilert`           | `ilert` 类型告警的配置。<br />参见[配置 ilert 告警](#configuring-ilert-alerts)。                              | `{}`    |
| `alerting.incident-io`     | `incident-io` 类型告警的配置。<br />参见[配置 Incident.io 告警](#configuring-incidentio-alerts)。             | `{}`    |
| `alerting.line`            | `line` 类型告警的配置。<br />参见[配置 Line 告警](#configuring-line-alerts)。                                 | `{}`    |
| `alerting.matrix`          | `matrix` 类型告警的配置。<br />参见[配置 Matrix 告警](#configuring-matrix-alerts)。                           | `{}`    |
| `alerting.mattermost`      | `mattermost` 类型告警的配置。<br />参见[配置 Mattermost 告警](#configuring-mattermost-alerts)。               | `{}`    |
| `alerting.messagebird`     | `messagebird` 类型告警的配置。<br />参见[配置 Messagebird 告警](#configuring-messagebird-alerts)。            | `{}`    |
| `alerting.n8n`             | `n8n` 类型告警的配置。<br />参见[配置 n8n 告警](#configuring-n8n-alerts)。                                    | `{}`    |
| `alerting.newrelic`        | `newrelic` 类型告警的配置。<br />参见[配置 New Relic 告警](#configuring-new-relic-alerts)。                   | `{}`    |
| `alerting.ntfy`            | `ntfy` 类型告警的配置。<br />参见[配置 Ntfy 告警](#configuring-ntfy-alerts)。                                 | `{}`    |
| `alerting.opsgenie`        | `opsgenie` 类型告警的配置。<br />参见[配置 Opsgenie 告警](#configuring-opsgenie-alerts)。                     | `{}`    |
| `alerting.pagerduty`       | `pagerduty` 类型告警的配置。<br />参见[配置 PagerDuty 告警](#configuring-pagerduty-alerts)。                  | `{}`    |
| `alerting.plivo`           | `plivo` 类型告警的配置。<br />参见[配置 Plivo 告警](#configuring-plivo-alerts)。                              | `{}`    |
| `alerting.pushover`        | `pushover` 类型告警的配置。<br />参见[配置 Pushover 告警](#configuring-pushover-alerts)。                     | `{}`    |
| `alerting.rocketchat`      | `rocketchat` 类型告警的配置。<br />参见[配置 Rocket.Chat 告警](#configuring-rocketchat-alerts)。              | `{}`    |
| `alerting.sendgrid`        | `sendgrid` 类型告警的配置。<br />参见[配置 SendGrid 告警](#configuring-sendgrid-alerts)。                     | `{}`    |
| `alerting.signal`          | `signal` 类型告警的配置。<br />参见[配置 Signal 告警](#configuring-signal-alerts)。                           | `{}`    |
| `alerting.signl4`          | `signl4` 类型告警的配置。<br />参见[配置 SIGNL4 告警](#configuring-signl4-alerts)。                           | `{}`    |
| `alerting.slack`           | `slack` 类型告警的配置。<br />参见[配置 Slack 告警](#configuring-slack-alerts)。                              | `{}`    |
| `alerting.splunk`          | `splunk` 类型告警的配置。<br />参见[配置 Splunk 告警](#configuring-splunk-alerts)。                           | `{}`    |
| `alerting.squadcast`       | `squadcast` 类型告警的配置。<br />参见[配置 Squadcast 告警](#configuring-squadcast-alerts)。                  | `{}`    |
| `alerting.teams`           | `teams` 类型告警的配置。*(已弃用)* <br />参见[配置 Teams 告警](#configuring-teams-alerts-deprecated)。    | `{}`    |
| `alerting.teams-workflows` | `teams-workflows` 类型告警的配置。<br />参见[配置 Teams Workflow 告警](#configuring-teams-workflow-alerts)。  | `{}`    |
| `alerting.telegram`        | `telegram` 类型告警的配置。<br />参见[配置 Telegram 告警](#configuring-telegram-alerts)。                     | `{}`    |
| `alerting.twilio`          | `twilio` 类型告警的设置。<br />参见[配置 Twilio 告警](#configuring-twilio-alerts)。                                | `{}`    |
| `alerting.vonage`          | `vonage` 类型告警的配置。<br />参见[配置 Vonage 告警](#configuring-vonage-alerts)。                           | `{}`    |
| `alerting.webex`           | `webex` 类型告警的配置。<br />参见[配置 Webex 告警](#configuring-webex-alerts)。                              | `{}`    |
| `alerting.zapier`          | `zapier` 类型告警的配置。<br />参见[配置 Zapier 告警](#configuring-zapier-alerts)。                           | `{}`    |
| `alerting.zulip`           | `zulip` 类型告警的配置。<br />参见[配置 Zulip 告警](#configuring-zulip-alerts)。                              | `{}`    |


#### 配置 AWS SES 告警
| 参数                            | 描述                                                                                | 默认值       |
|:-------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.aws-ses`                   | `aws-ses` 类型告警的设置                                                      | `{}`          |
| `alerting.aws-ses.access-key-id`     | AWS 访问密钥 ID                                                                          | Optional `""` |
| `alerting.aws-ses.secret-access-key` | AWS 秘密访问密钥                                                                      | Optional `""` |
| `alerting.aws-ses.region`            | AWS 区域                                                                                 | Required `""` |
| `alerting.aws-ses.from`              | 发送邮件的邮箱地址（应在 SES 中注册）                    | Required `""` |
| `alerting.aws-ses.to`                | 逗号分隔的通知邮箱地址列表                                            | Required `""` |
| `alerting.aws-ses.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A           |
| `alerting.aws-ses.overrides`         | 可能优先于默认配置的覆盖列表                   | `[]`          |
| `alerting.aws-ses.overrides[].group` | 此配置将覆盖其配置的端点组        | `""`          |
| `alerting.aws-ses.overrides[].*`     | 参见 `alerting.aws-ses.*` 参数                                                        | `{}`          |

```yaml
alerting:
  aws-ses:
    access-key-id: "..."
    secret-access-key: "..."
    region: "us-east-1"
    from: "status@example.com"
    to: "user@example.com"

endpoints:
  - name: website
    interval: 30s
    url: "https://twin.sh/health"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: aws-ses
        failure-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```

如果未定义 `access-key-id` 和 `secret-access-key`，Gatus 将回退使用 IAM 认证。

请确保你有使用 `ses:SendEmail` 的权限。


#### 配置 ClickUp 告警

| 参数                          | 描述                                                                                | 默认值       |
| :--------------------------------- | :----------------------------------------------------------------------------------------- | :------------ |
| `alerting.clickup`                 | `clickup` 类型告警的配置                                                 | `{}`          |
| `alerting.clickup.list-id`         | 将创建任务的 ClickUp 列表 ID                                                | Required `""` |
| `alerting.clickup.token`           | ClickUp API 令牌                                                                          | Required `""` |
| `alerting.clickup.api-url`         | 自定义 API URL                   | `https://api.clickup.com/api/v2`          |
| `alerting.clickup.assignees`       | 要分配任务的用户 ID 列表                                                        | `[]`          |
| `alerting.clickup.status`          | 创建任务的初始状态                                                           | `""`          |
| `alerting.clickup.priority`        | 优先级：`urgent`、`high`、`normal`、`low` 或 `none`                               | `normal`      |
| `alerting.clickup.notify-all`      | 创建任务时是否通知所有受理人                                       | `true`        |
| `alerting.clickup.name`            | 自定义任务名称模板（支持占位符）                                          | `Health Check: [ENDPOINT_GROUP]:[ENDPOINT_NAME]`          |
| `alerting.clickup.content`         | 自定义任务内容模板（支持占位符）                                       | `Triggered: [ENDPOINT_GROUP] - [ENDPOINT_NAME] - [ALERT_DESCRIPTION] - [RESULT_ERRORS]`          |
| `alerting.clickup.default-alert`   | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A           |
| `alerting.clickup.overrides`       | 可能优先于默认配置的覆盖列表                   | `[]`          |
| `alerting.clickup.overrides[].group` | 此配置将覆盖其配置的端点组      | `""`          |
| `alerting.clickup.overrides[].*`   | 参见 `alerting.clickup.*` 参数                                                        | `{}`          |

ClickUp 告警提供商在告警触发时会在 ClickUp 列表中创建任务。如果端点告警设置了 `send-on-resolved` 为 `true`，当告警解决时任务将自动关闭。

`name` 和 `content` 中支持以下占位符：

-   `[ENDPOINT_GROUP]` - 从 `endpoints[].group` 解析
-   `[ENDPOINT_NAME]` - 从 `endpoints[].name` 解析
-   `[ALERT_DESCRIPTION]` - 从 `endpoints[].alerts[].description` 解析
-   `[RESULT_ERRORS]` - 从健康检查评估错误中解析

```yaml
alerting:
  clickup:
    list-id: "123456789"
    token: "pk_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    assignees:
      - "12345"
      - "67890"
    status: "in progress"
    priority: high
    name: "Health Check Alert: [ENDPOINT_GROUP] - [ENDPOINT_NAME]"
    content: "Alert triggered for [ENDPOINT_GROUP] - [ENDPOINT_NAME] - [ALERT_DESCRIPTION] - [RESULT_ERRORS]"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: clickup
        send-on-resolved: true
```

要获取你的 ClickUp API 令牌，请参考：[生成或重新生成个人 API 令牌](https://developer.clickup.com/docs/authentication#:~:text=the%20API%20docs.-,Generate%20or%20regenerate%20a%20Personal%20API%20Token,-Log%20in%20to)

要查找你的列表 ID：

1. 打开你想要创建任务的 ClickUp 列表
2. 列表 ID 在 URL 中：`https://app.clickup.com/{workspace_id}/v/l/li/{list_id}`

要查找受理人 ID：

1. 前往 `https://app.clickup.com/{workspace_id}/teams-pulse/teams/people`
2. 将鼠标悬停在团队成员上
3. 点击三个点（溢出菜单）
3. 点击 `Copy member ID`

#### 配置 Datadog 告警

> ⚠️ **警告**：此告警提供商尚未经过测试。如果你已测试并确认其正常工作，请删除此警告并创建 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供商是否按预期工作。感谢你的配合。

| 参数                            | 描述                                                                                | 默认值           |
|:-------------------------------------|:-------------------------------------------------------------------------------------------|:------------------|
| `alerting.datadog`                   | `datadog` 类型告警的配置                                                 | `{}`              |
| `alerting.datadog.api-key`           | Datadog API 密钥                                                                            | Required `""`     |
| `alerting.datadog.site`              | Datadog 站点（例如 datadoghq.com、datadoghq.eu）                                           | `"datadoghq.com"` |
| `alerting.datadog.tags`              | 要包含的附加标签                                                                 | `[]`              |
| `alerting.datadog.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A               |
| `alerting.datadog.overrides`         | 可能优先于默认配置的覆盖列表                   | `[]`              |
| `alerting.datadog.overrides[].group` | 此配置将覆盖其配置的端点组        | `""`              |
| `alerting.datadog.overrides[].*`     | 参见 `alerting.datadog.*` 参数                                                        | `{}`              |

```yaml
alerting:
  datadog:
    api-key: "YOUR_API_KEY"
    site: "datadoghq.com"  # or datadoghq.eu for EU region
    tags:
      - "environment:production"
      - "team:platform"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: datadog
        send-on-resolved: true
```


#### 配置 Discord 告警
| 参数                            | 描述                                                                                | 默认值                             |
|:-------------------------------------|:-------------------------------------------------------------------------------------------|:------------------------------------|
| `alerting.discord`                   | `discord` 类型告警的配置                                                 | `{}`                                |
| `alerting.discord.webhook-url`       | Discord Webhook URL                                                                        | Required `""`                       |
| `alerting.discord.title`             | 通知标题                                                                  | `":helmet_with_white_cross: Gatus"` |
| `alerting.discord.message-content`   | 在嵌入内容之前发送的消息内容（可用于 @ 提及用户/角色，例如 `<@123>`）   | `""`                                |
| `alerting.discord.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A                                 |
| `alerting.discord.overrides`         | 可能优先于默认配置的覆盖列表                   | `[]`                                |
| `alerting.discord.overrides[].group` | 此配置将覆盖其配置的端点组        | `""`                                |
| `alerting.discord.overrides[].*`     | 参见 `alerting.discord.*` 参数                                                        | `{}`                                |

```yaml
alerting:
  discord:
    webhook-url: "https://discord.com/api/webhooks/**********/**********"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: discord
        description: "healthcheck failed"
        send-on-resolved: true
```


#### 配置 Email 告警
| 参数                          | 描述                                                                                   | 默认值       |
|:-----------------------------------|:----------------------------------------------------------------------------------------------|:--------------|
| `alerting.email`                   | `email` 类型告警的配置                                                      | `{}`          |
| `alerting.email.from`              | 用于发送告警的邮箱地址                                                                  | Required `""` |
| `alerting.email.username`          | 用于发送告警的 SMTP 服务器用户名。如果为空，则使用 `alerting.email.from`。     | `""`          |
| `alerting.email.password`          | 用于发送告警的 SMTP 服务器密码。如果为空，则不执行认证。 | `""`          |
| `alerting.email.host`              | 邮件服务器主机（例如 `smtp.gmail.com`）                                               | Required `""` |
| `alerting.email.port`              | 邮件服务器监听端口（例如 `587`）                                             | Required `0`  |
| `alerting.email.to`                | 发送告警的目标邮箱地址                                                                | Required `""` |
| `alerting.email.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert)    | N/A           |
| `alerting.email.client.insecure`   | 是否跳过 TLS 验证                                                              | `false`       |
| `alerting.email.overrides`         | 可能优先于默认配置的覆盖列表                      | `[]`          |
| `alerting.email.overrides[].group` | 此配置将覆盖其配置的端点组           | `""`          |
| `alerting.email.overrides[].*`     | 参见 `alerting.email.*` 参数                                                             | `{}`          |

```yaml
alerting:
  email:
    from: "from@example.com"
    username: "from@example.com"
    password: "hunter2"
    host: "mail.example.com"
    port: 587
    to: "recipient1@example.com,recipient2@example.com"
    client:
      insecure: false
    # 你还可以添加特定组的 to 键，
    # 这将覆盖上面指定组的 to 键
    overrides:
      - group: "core"
        to: "recipient3@example.com,recipient4@example.com"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: email
        description: "healthcheck failed"
        send-on-resolved: true

  - name: back-end
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[CERTIFICATE_EXPIRATION] > 48h"
    alerts:
      - type: email
        description: "healthcheck failed"
        send-on-resolved: true
```

> ⚠ 某些邮件服务器速度可能非常慢。


#### 配置 Gitea 告警

| 参数                       | 描述                                                                                                | 默认值       |
|:--------------------------------|:-----------------------------------------------------------------------------------------------------------|:--------------|
| `alerting.gitea`                | `gitea` 类型告警的配置                                                                   | `{}`          |
| `alerting.gitea.repository-url` | Gitea 仓库 URL（例如 `https://gitea.com/TwiN/example`）                                               | Required `""` |
| `alerting.gitea.token`          | 用于认证的个人访问令牌。<br />至少需要对 issues 有读写权限，对 metadata 有只读权限。 | Required `""` |
| `alerting.gitea.default-alert`  | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert)。                | N/A           |

Gitea 告警提供商会为每个告警创建一个以 `alert(gatus):` 为前缀、以端点显示名称为后缀的 issue。
如果端点告警设置了 `send-on-resolved` 为 `true`，当告警解决时该 issue 将自动关闭。

```yaml
alerting:
  gitea:
    repository-url: "https://gitea.com/TwiN/test"
    token: "349d63f16......"

endpoints:
  - name: example
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 75"
    alerts:
      - type: gitea
        failure-threshold: 2
        success-threshold: 3
        send-on-resolved: true
        description: "Everything's burning AAAAAHHHHHHHHHHHHHHH"
```

![Gitea 告警](.github/assets/gitea-alerts.png)


#### 配置 GitHub 告警

| 参数                        | 描述                                                                                                | 默认值       |
|:---------------------------------|:-----------------------------------------------------------------------------------------------------------|:--------------|
| `alerting.github`                | `github` 类型告警的配置                                                                  | `{}`          |
| `alerting.github.repository-url` | GitHub 仓库 URL（例如 `https://github.com/TwiN/example`）                                             | Required `""` |
| `alerting.github.token`          | 用于认证的个人访问令牌。<br />至少需要对 issues 有读写权限，对 metadata 有只读权限。 | Required `""` |
| `alerting.github.default-alert`  | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert)。                | N/A           |

GitHub 告警提供商会为每个告警创建一个以 `alert(gatus):` 为前缀、以端点显示名称为后缀的 issue。
如果端点告警设置了 `send-on-resolved` 为 `true`，当告警解决时该 issue 将自动关闭。

```yaml
alerting:
  github:
    repository-url: "https://github.com/TwiN/test"
    token: "github_pat_12345..."

endpoints:
  - name: example
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 75"
    alerts:
      - type: github
        failure-threshold: 2
        success-threshold: 3
        send-on-resolved: true
        description: "Everything's burning AAAAAHHHHHHHHHHHHHHH"
```

![GitHub 告警](.github/assets/github-alerts.png)


#### 配置 GitLab 告警
| 参数                           | 描述                                                                                                         | 默认值       |
|:------------------------------------|:--------------------------------------------------------------------------------------------------------------------|:--------------|
| `alerting.gitlab`                   | `gitlab` 类型告警的配置                                                                           | `{}`          |
| `alerting.gitlab.webhook-url`       | GitLab 告警 Webhook URL（例如 `https://gitlab.com/yourusername/example/alerts/notify/gatus/xxxxxxxxxxxxxxxx.json`） | Required `""` |
| `alerting.gitlab.authorization-key` | GitLab 告警授权密钥。                                                                                     | Required `""` |
| `alerting.gitlab.severity`          | 覆盖默认严重级别（critical），可选值为 `critical, high, medium, low, info, unknown`                    | `""`          |
| `alerting.gitlab.monitoring-tool`   | 覆盖监控工具名称（gatus）                                                                           | `"gatus"`     |
| `alerting.gitlab.environment-name`  | 设置 GitLab 环境名称。在仪表盘上显示告警时需要此项。                                           | `""`          |
| `alerting.gitlab.service`           | 覆盖端点显示名称                                                                                      | `""`          |
| `alerting.gitlab.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert)。                         | N/A           |

GitLab 告警提供商会为每个告警创建一个以 `alert(gatus):` 为前缀、以端点显示名称为后缀的告警。
如果端点告警设置了 `send-on-resolved` 为 `true`，当告警解决时该告警将自动关闭。参见
https://docs.gitlab.com/ee/operations/incident_management/integrations.html#configuration 以配置端点。

```yaml
alerting:
  gitlab:
    webhook-url: "https://gitlab.com/hlidotbe/example/alerts/notify/gatus/xxxxxxxxxxxxxxxx.json"
    authorization-key: "12345"

endpoints:
  - name: example
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 75"
    alerts:
      - type: gitlab
        failure-threshold: 2
        success-threshold: 3
        send-on-resolved: true
        description: "Everything's burning AAAAAHHHHHHHHHHHHHHH"
```

![GitLab 告警](.github/assets/gitlab-alerts.png)


#### 配置 Google Chat 告警
| 参数                               | 描述                                                                                 | 默认值       |
|:----------------------------------------|:--------------------------------------------------------------------------------------------|:--------------|
| `alerting.googlechat`                   | `googlechat` 类型告警的配置                                               | `{}`          |
| `alerting.googlechat.webhook-url`       | Google Chat Webhook URL                                                                     | Required `""` |
| `alerting.googlechat.client`            | 客户端配置。<br />参见[客户端配置](#client-configuration)。              | `{}`          |
| `alerting.googlechat.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert)。 | N/A           |
| `alerting.googlechat.overrides`         | 可能优先于默认配置的覆盖列表                    | `[]`          |
| `alerting.googlechat.overrides[].group` | 此配置将覆盖其配置的端点组         | `""`          |
| `alerting.googlechat.overrides[].*`     | 参见 `alerting.googlechat.*` 参数                                                      | `{}`          |

```yaml
alerting:
  googlechat:
    webhook-url: "https://chat.googleapis.com/v1/spaces/*******/messages?key=**********&token=********"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: googlechat
        description: "healthcheck failed"
        send-on-resolved: true
```


#### 配置 Gotify 告警
| 参数                                     | 描述                                                                                 | 默认值               |
|:----------------------------------------------|:--------------------------------------------------------------------------------------------|:----------------------|
| `alerting.gotify`                             | `gotify` 类型告警的配置                                                   | `{}`                  |
| `alerting.gotify.server-url`                  | Gotify 服务器 URL                                                                           | Required `""`         |
| `alerting.gotify.token`                       | 用于认证的令牌。                                                      | Required `""`         |
| `alerting.gotify.priority`                    | 根据 Gotify 标准设置的告警优先级。                                        | `5`                   |
| `alerting.gotify.title`                       | 通知标题                                                                   | `"Gatus: <endpoint>"` |
| `alerting.gotify.default-alert`               | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert)。 | N/A                   |

```yaml
alerting:
  gotify:
    server-url: "https://gotify.example"
    token: "**************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: gotify
        description: "healthcheck failed"
        send-on-resolved: true
```

以下是通知的示例截图：

![Gotify 通知](.github/assets/gotify-alerts.png)


#### 配置 HomeAssistant 告警
| 参数                                  | 描述                                                                            | 默认值 |
|:-------------------------------------------|:---------------------------------------------------------------------------------------|:--------------|
| `alerting.homeassistant.url`               | HomeAssistant 实例 URL                                                             | Required `""` |
| `alerting.homeassistant.token`             | HomeAssistant 的长期访问令牌                                             | Required `""` |
| `alerting.homeassistant.default-alert`     | 用于具有相应类型告警的端点的默认告警配置 | `{}`          |
| `alerting.homeassistant.overrides`         | 可能优先于默认配置的覆盖列表               | `[]`          |
| `alerting.homeassistant.overrides[].group` | 此配置将覆盖其配置的端点组    | `""`          |
| `alerting.homeassistant.overrides[].*`     | 参见 `alerting.homeassistant.*` 参数                                              | `{}`          |

```yaml
alerting:
  homeassistant:
    url: "http://homeassistant:8123"  # URL of your HomeAssistant instance
    token: "YOUR_LONG_LIVED_ACCESS_TOKEN"  # Long-lived access token from HomeAssistant

endpoints:
  - name: my-service
    url: "https://my-service.com"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: homeassistant
        enabled: true
        send-on-resolved: true
        description: "My service health check"
        failure-threshold: 3
        success-threshold: 2
```

告警将作为事件发送到 HomeAssistant，事件类型为 `gatus_alert`。事件数据包括：
- `status`：`"triggered"` 或 `"resolved"`
- `endpoint`：被监控端点的名称
- `description`：告警描述（如果提供）
- `conditions`：条件列表及其结果
- `failure_count`：连续失败次数（触发时）
- `success_count`：连续成功次数（解决时）

你可以在 HomeAssistant 自动化中使用这些事件来：
- 发送通知
- 控制设备
- 触发场景
- 记录到历史
- 以及更多

HomeAssistant 自动化示例：
```yaml
automation:
  - alias: "Gatus Alert Handler"
    trigger:
      platform: event
      event_type: gatus_alert
    action:
      - service: notify.notify
        data_template:
          title: "Gatus Alert: {{ trigger.event.data.event_data.endpoint }}"
          message: >
            Status: {{ trigger.event.data.event_data.status }}
            {% if trigger.event.data.event_data.description %}
            Description: {{ trigger.event.data.event_data.description }}
            {% endif %}
            {% for condition in trigger.event.data.event_data.conditions %}
            {{ '✅' if condition.success else '❌' }} {{ condition.condition }}
            {% endfor %}
```

要获取你的 HomeAssistant 长期访问令牌：
1. 打开 HomeAssistant
2. 点击你的个人资料名称（左下角）
3. 向下滚动到 "Long-Lived Access Tokens"
4. 点击 "Create Token"
5. 给它起一个名字（例如 "Gatus"）
6. 复制令牌 - 你只能看到一次！


#### 配置 IFTTT 告警

> ⚠️ **警告**：此告警提供商尚未经过测试。如果你已测试并确认其正常工作，请删除此警告并创建 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供商是否按预期工作。感谢你的配合。

| 参数                          | 描述                                                                                | 默认值       |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.ifttt`                   | `ifttt` 类型告警的配置                                                   | `{}`          |
| `alerting.ifttt.webhook-key`       | IFTTT Webhook 密钥                                                                          | Required `""` |
| `alerting.ifttt.event-name`        | IFTTT 事件名称                                                                           | Required `""` |
| `alerting.ifttt.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A           |
| `alerting.ifttt.overrides`         | 可能优先于默认配置的覆盖列表                   | `[]`          |
| `alerting.ifttt.overrides[].group` | 此配置将覆盖其配置的端点组        | `""`          |
| `alerting.ifttt.overrides[].*`     | 参见 `alerting.ifttt.*` 参数                                                          | `{}`          |

```yaml
alerting:
  ifttt:
    webhook-key: "YOUR_WEBHOOK_KEY"
    event-name: "gatus_alert"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: ifttt
        send-on-resolved: true
```


#### 配置 Ilert 告警
| 参数                          | 描述                                                                                | 默认值 |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:--------|
| `alerting.ilert`                   | `ilert` 类型告警的配置                                                   | `{}`    |
| `alerting.ilert.integration-key`   | ilert 告警源集成密钥                                                         | `""`    |
| `alerting.ilert.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A     |
| `alerting.ilert.overrides`         | 可能优先于默认配置的覆盖列表                   | `[]`    |
| `alerting.ilert.overrides[].group` | 此配置将覆盖其配置的端点组        | `""`    |
| `alerting.ilert.overrides[].*`     | 参见 `alerting.ilert.*` 参数                                                          | `{}`    |

强烈建议将 `ilert` 类型告警的 `endpoints[].alerts[].send-on-resolved` 设置为 `true`，
因为与其他告警不同，将该参数设置为 `true` 所产生的操作不会创建另一个告警，
而是在 ilert 上将告警标记为已解决。

行为：
- 默认情况下，使用 `alerting.ilert.integration-key` 作为集成密钥
- 如果被评估的端点属于某个组（`endpoints[].group`），且该组与 `alerting.ilert.overrides[].group` 的值匹配，则提供商将使用该覆盖的集成密钥，而非 `alerting.ilert.integration-key` 的值

```yaml
alerting:
  ilert:
    integration-key: "********************************"
    # 你还可以添加特定组的集成密钥，
    # 这将覆盖上面指定组的集成密钥
    overrides:
      - group: "core"
        integration-key: "********************************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: ilert
        failure-threshold: 3
        success-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Incident.io 告警
| 参数                                | 描述                                                                                | 默认值       |
|:-----------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.incident-io`                   | `incident-io` 类型告警的配置                                             | `{}`          |
| `alerting.incident-io.url`               | 触发告警事件的 URL。                                                             | Required `""` |
| `alerting.incident-io.auth-token`        | 用于认证的令牌。                                                     | Required `""` |
| `alerting.incident-io.source-url`        | 来源 URL                                                                                 | `""`          |
| `alerting.incident-io.default-alert`     | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A           |
| `alerting.incident-io.overrides`         | 可能优先于默认配置的覆盖列表                   | `[]`          |
| `alerting.incident-io.overrides[].group` | 此配置将覆盖其配置的端点组        | `""`          |
| `alerting.incident-io.overrides[].*`     | 参见 `alerting.incident-io.*` 参数                                                    | `{}`          |

```yaml
alerting:
  incident-io:
    url: "*****************"
    auth-token: "********************************************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: incident-io
        description: "healthcheck failed"
        send-on-resolved: true
```
要获取所需的告警源配置 ID 和认证令牌，你必须配置 HTTP 告警源。

> **_注意：_** 源配置 ID 的格式为 `https://api.incident.io/v2/alert_events/http/$ID`，令牌应作为 Bearer 令牌传递，格式如下：`Authorization: Bearer $TOKEN`


#### 配置 Line 告警

| 参数                            | 描述                                                                                | 默认值       |
|:-------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.line`                      | `line` 类型告警的配置                                                    | `{}`          |
| `alerting.line.channel-access-token` | Line Messaging API 频道访问令牌                                                    | Required `""` |
| `alerting.line.user-ids`             | 要发送消息的 Line 用户 ID 列表（可以是用户 ID、房间 ID 或群组 ID）    | Required `[]` |
| `alerting.line.default-alert`        | 默认告警配置。<br />参见[设置默认告警](#setting-a-default-alert) | N/A           |
| `alerting.line.overrides`            | 可能优先于默认配置的覆盖列表                   | `[]`          |
| `alerting.line.overrides[].group`    | 此配置将覆盖其配置的端点组        | `""`          |
| `alerting.line.overrides[].*`        | 参见 `alerting.line.*` 参数                                                           | `{}`          |

```yaml
alerting:
  line:
    channel-access-token: "YOUR_CHANNEL_ACCESS_TOKEN"
    user-ids:
      - "U1234567890abcdef" # This can be a group id, room id or user id
      - "U2345678901bcdefg"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: line
        send-on-resolved: true
```

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: ifttt
        send-on-resolved: true
```


#### 配置 Ilert 告警
| 参数                                | 描述                                                                                     | 默认值  |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:--------|
| `alerting.ilert`                   | `ilert` 类型告警的配置                                                                      | `{}`    |
| `alerting.ilert.integration-key`   | ilert 告警源集成密钥                                                                        | `""`    |
| `alerting.ilert.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A     |
| `alerting.ilert.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`    |
| `alerting.ilert.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`    |
| `alerting.ilert.overrides[].*`     | 参见 `alerting.ilert.*` 参数                                                                | `{}`    |

强烈建议将 `ilert` 类型告警的 `endpoints[].alerts[].send-on-resolved` 设置为 `true`，因为与其他告警不同，将该参数设置为 `true` 所产生的操作不会创建另一个告警，而是在 ilert 上将该告警标记为已解决。

行为：
- 默认情况下，使用 `alerting.ilert.integration-key` 作为集成密钥
- 如果被评估的端点属于某个组（`endpoints[].group`），且该组匹配 `alerting.ilert.overrides[].group` 的值，则提供者将使用该覆盖配置的集成密钥，而不是 `alerting.ilert.integration-key` 的值

```yaml
alerting:
  ilert:
    integration-key: "********************************"
    # 你也可以添加特定组的集成密钥，
    # 这将覆盖上面指定组的集成密钥
    overrides:
      - group: "core"
        integration-key: "********************************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: ilert
        failure-threshold: 3
        success-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Incident.io 告警
| 参数                                    | 描述                                                                                     | 默认值         |
|:-----------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.incident-io`                   | `incident-io` 类型告警的配置                                                                | `{}`          |
| `alerting.incident-io.url`               | 触发告警事件的 URL。                                                                        | 必填 `""`      |
| `alerting.incident-io.auth-token`        | 用于身份验证的令牌。                                                                         | 必填 `""`      |
| `alerting.incident-io.source-url`        | 来源 URL                                                                                   | `""`          |
| `alerting.incident-io.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.incident-io.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.incident-io.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.incident-io.overrides[].*`     | 参见 `alerting.incident-io.*` 参数                                                          | `{}`          |

```yaml
alerting:
  incident-io:
    url: "*****************"
    auth-token: "********************************************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: incident-io
        description: "healthcheck failed"
        send-on-resolved: true
```
要获取所需的告警源配置 ID 和身份验证令牌，你必须配置一个 HTTP 告警源。

> **_注意：_** 来源配置 ID 的格式为 `https://api.incident.io/v2/alert_events/http/$ID`，令牌需要作为 Bearer 令牌传递，格式如下：`Authorization: Bearer $TOKEN`


#### 配置 Line 告警

| 参数                                  | 描述                                                                                     | 默认值         |
|:-------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.line`                      | `line` 类型告警的配置                                                                       | `{}`          |
| `alerting.line.channel-access-token` | Line Messaging API 频道访问令牌                                                             | 必填 `""`      |
| `alerting.line.user-ids`             | 要发送消息的 Line 用户 ID 列表（可以是用户 ID、房间 ID 或群组 ID）                                | 必填 `[]`      |
| `alerting.line.default-alert`        | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.line.overrides`            | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.line.overrides[].group`    | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.line.overrides[].*`        | 参见 `alerting.line.*` 参数                                                                 | `{}`          |

```yaml
alerting:
  line:
    channel-access-token: "YOUR_CHANNEL_ACCESS_TOKEN"
    user-ids:
      - "U1234567890abcdef" # This can be a group id, room id or user id
      - "U2345678901bcdefg"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: line
        send-on-resolved: true
```


#### 配置 Matrix 告警
| 参数                                    | 描述                                                                                     | 默认值                              |
|:-----------------------------------------|:-------------------------------------------------------------------------------------------|:-----------------------------------|
| `alerting.matrix`                        | `matrix` 类型告警的配置                                                                     | `{}`                               |
| `alerting.matrix.server-url`             | Homeserver URL                                                                             | `https://matrix-client.matrix.org` |
| `alerting.matrix.access-token`           | 机器人用户访问令牌（参见 https://webapps.stackexchange.com/q/131056）                          | 必填 `""`                           |
| `alerting.matrix.internal-room-id`       | 发送告警的房间内部 ID（可在房间设置 > 高级中找到）                                               | 必填 `""`                           |
| `alerting.matrix.default-alert`          | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A                                |
| `alerting.matrix.overrides`              | 可以优先于默认配置的覆盖列表                                                                  | `[]`                               |
| `alerting.matrix.overrides[].group`      | 将使用此配置覆盖默认配置的端点组                                                               | `""`                               |
| `alerting.matrix.overrides[].*`          | 参见 `alerting.matrix.*` 参数                                                               | `{}`                               |

```yaml
alerting:
  matrix:
    server-url: "https://matrix-client.matrix.org"
    access-token: "123456"
    internal-room-id: "!example:matrix.org"

endpoints:
  - name: website
    interval: 5m
    url: "https://twin.sh/health"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: matrix
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Mattermost 告警
| 参数                                         | 描述                                                                                      | 默认值         |
|:----------------------------------------------|:--------------------------------------------------------------------------------------------|:--------------|
| `alerting.mattermost`                         | `mattermost` 类型告警的配置                                                                  | `{}`          |
| `alerting.mattermost.webhook-url`             | Mattermost Webhook URL                                                                      | 必填 `""`      |
| `alerting.mattermost.channel`                 | Mattermost 频道名称覆盖（可选）                                                               | `""`          |
| `alerting.mattermost.client`                  | 客户端配置。<br />参见 [客户端配置](#client-configuration)。                                    | `{}`          |
| `alerting.mattermost.default-alert`           | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)。                              | N/A           |
| `alerting.mattermost.overrides`               | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.mattermost.overrides[].group`       | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.mattermost.overrides[].*`           | 参见 `alerting.mattermost.*` 参数                                                            | `{}`          |

```yaml
alerting:
  mattermost:
    webhook-url: "http://**********/hooks/**********"
    client:
      insecure: true

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: mattermost
        description: "healthcheck failed"
        send-on-resolved: true
```

以下是通知的示例效果：

![Mattermost 通知](.github/assets/mattermost-alerts.png)


#### 配置 Messagebird 告警
| 参数                                  | 描述                                                                                     | 默认值         |
|:-------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.messagebird`               | `messagebird` 类型告警的配置                                                                 | `{}`          |
| `alerting.messagebird.access-key`    | Messagebird 访问密钥                                                                        | 必填 `""`      |
| `alerting.messagebird.originator`    | 消息发送者                                                                                  | 必填 `""`      |
| `alerting.messagebird.recipients`    | 消息接收者                                                                                  | 必填 `""`      |
| `alerting.messagebird.default-alert` | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |

使用 Messagebird 发送 **SMS** 短信告警的示例：
```yaml
alerting:
  messagebird:
    access-key: "..."
    originator: "31619191918"
    recipients: "31619191919,31619191920"

endpoints:
  - name: website
    interval: 5m
    url: "https://twin.sh/health"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: messagebird
        failure-threshold: 3
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 New Relic 告警

> **警告**：此告警提供者尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供者是否按预期工作。感谢你的配合。

| 参数                                   | 描述                                                                                     | 默认值         |
|:--------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.newrelic`                   | `newrelic` 类型告警的配置                                                                    | `{}`          |
| `alerting.newrelic.api-key`           | New Relic API 密钥                                                                          | 必填 `""`      |
| `alerting.newrelic.account-id`        | New Relic 账户 ID                                                                           | 必填 `""`      |
| `alerting.newrelic.region`            | 区域（US 或 EU）                                                                            | `"US"`        |
| `alerting.newrelic.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.newrelic.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.newrelic.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.newrelic.overrides[].*`     | 参见 `alerting.newrelic.*` 参数                                                              | `{}`          |

```yaml
alerting:
  newrelic:
    api-key: "YOUR_API_KEY"
    account-id: "1234567"
    region: "US"  # or "EU" for European region

endpoints:
  - name: example
    url: "https://example.org"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: newrelic
        send-on-resolved: true
```


#### 配置 n8n 告警
| 参数                              | 描述                                                                                     | 默认值         |
|:---------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.n8n`                   | `n8n` 类型告警的配置                                                                        | `{}`          |
| `alerting.n8n.webhook-url`       | n8n webhook URL                                                                            | 必填 `""`      |
| `alerting.n8n.title`             | 发送到 n8n 的告警标题                                                                        | `""`          |
| `alerting.n8n.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.n8n.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.n8n.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.n8n.overrides[].*`     | 参见 `alerting.n8n.*` 参数                                                                  | `{}`          |

[n8n](https://n8n.io/) 是一个工作流自动化平台，允许你使用 Webhook 在不同的应用程序和服务之间自动执行任务。

参见 [n8n-nodes-gatus-trigger](https://github.com/TwiN/n8n-nodes-gatus-trigger) 了解可用作触发器的 n8n 社区节点。

示例：
```yaml
alerting:
  n8n:
    webhook-url: "https://your-n8n-instance.com/webhook/your-webhook-id"
    title: "Gatus Monitoring"
    default-alert:
      send-on-resolved: true

endpoints:
  - name: example
    url: "https://example.org"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: n8n
        description: "Health check alert"
```

发送到 n8n webhook 的 JSON 负载将包含：
- `title`：配置的标题
- `endpoint_name`：端点名称
- `endpoint_group`：端点所属组（如有）
- `endpoint_url`：被监控的 URL
- `alert_description`：自定义告警描述
- `resolved`：布尔值，指示告警是否已解决
- `message`：人类可读的告警消息
- `condition_results`：条件结果数组，包含其成功状态


#### 配置 Ntfy 告警
| 参数                                  | 描述                                                                                                                                          | 默认值             |
|:-------------------------------------|:---------------------------------------------------------------------------------------------------------------------------------------------|:------------------|
| `alerting.ntfy`                      | `ntfy` 类型告警的配置                                                                                                                          | `{}`              |
| `alerting.ntfy.topic`                | 发送告警的主题                                                                                                                                  | 必填 `""`          |
| `alerting.ntfy.url`                  | 目标服务器的 URL                                                                                                                                | `https://ntfy.sh` |
| `alerting.ntfy.token`                | 受限主题的[访问令牌](https://docs.ntfy.sh/publish/#access-tokens)                                                                                | `""`              |
| `alerting.ntfy.email`                | 用于额外电子邮件通知的电子邮件地址                                                                                                                 | `""`              |
| `alerting.ntfy.click`                | 点击通知时打开的网站                                                                                                                             | `""`              |
| `alerting.ntfy.priority`             | 告警的优先级                                                                                                                                    | `3`               |
| `alerting.ntfy.disable-firebase`     | 是否禁用通过 Firebase 的消息推送。[ntfy.sh 默认启用](https://docs.ntfy.sh/publish/#disable-firebase)                                                | `false`           |
| `alerting.ntfy.disable-cache`        | 是否禁用服务器端消息缓存。[ntfy.sh 默认启用](https://docs.ntfy.sh/publish/#message-caching)                                                        | `false`           |
| `alerting.ntfy.default-alert`        | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                                                                                  | N/A               |
| `alerting.ntfy.overrides`            | 可以优先于默认配置的覆盖列表                                                                                                                      | `[]`              |
| `alerting.ntfy.overrides[].group`    | 将使用此配置覆盖默认配置的端点组                                                                                                                   | `""`              |
| `alerting.ntfy.overrides[].*`        | 参见 `alerting.ntfy.*` 参数                                                                                                                     | `{}`              |

[ntfy](https://github.com/binwiederhier/ntfy) 是一个出色的项目，允许你订阅桌面和移动通知，使其成为 Gatus 的绝佳补充。

示例：
```yaml
alerting:
  ntfy:
    topic: "gatus-test-topic"
    priority: 2
    token: faketoken
    default-alert:
      failure-threshold: 3
      send-on-resolved: true
    # 你也可以添加特定组的密钥，
    # 这将覆盖上面指定组的密钥
    overrides:
      - group: "other"
        topic: "gatus-other-test-topic"
        priority: 4
        click: "https://example.com"

endpoints:
  - name: website
    interval: 5m
    url: "https://twin.sh/health"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: ntfy
  - name: other example
    group: other
    interval: 30m
    url: "https://example.com"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
    alerts:
      - type: ntfy
        description: example
```


#### 配置 Opsgenie 告警
| 参数                               | 描述                                                                                     | 默认值                 |
|:----------------------------------|:-------------------------------------------------------------------------------------------|:---------------------|
| `alerting.opsgenie`               | `opsgenie` 类型告警的配置                                                                    | `{}`                 |
| `alerting.opsgenie.api-key`       | Opsgenie API 密钥                                                                          | 必填 `""`             |
| `alerting.opsgenie.priority`      | 告警的优先级级别。                                                                           | `P1`                 |
| `alerting.opsgenie.source`        | 告警的来源字段。                                                                             | `gatus`              |
| `alerting.opsgenie.entity-prefix` | 实体字段前缀。                                                                               | `gatus-`             |
| `alerting.opsgenie.alias-prefix`  | 别名字段前缀。                                                                               | `gatus-healthcheck-` |
| `alerting.opsgenie.tags`          | 告警标签。                                                                                   | `[]`                 |
| `alerting.opsgenie.default-alert` | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A                  |

Opsgenie 提供者将自动打开和关闭告警。

```yaml
alerting:
  opsgenie:
    api-key: "00000000-0000-0000-0000-000000000000"
```


#### 配置 PagerDuty 告警
| 参数                                    | 描述                                                                                     | 默认值  |
|:---------------------------------------|:-------------------------------------------------------------------------------------------|:--------|
| `alerting.pagerduty`                   | `pagerduty` 类型告警的配置                                                                   | `{}`    |
| `alerting.pagerduty.integration-key`   | PagerDuty Events API v2 集成密钥                                                            | `""`    |
| `alerting.pagerduty.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A     |
| `alerting.pagerduty.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`    |
| `alerting.pagerduty.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`    |
| `alerting.pagerduty.overrides[].*`     | 参见 `alerting.pagerduty.*` 参数                                                            | `{}`    |

强烈建议将 `pagerduty` 类型告警的 `endpoints[].alerts[].send-on-resolved` 设置为 `true`，因为与其他告警不同，将该参数设置为 `true` 所产生的操作不会创建另一个事件，而是在 PagerDuty 上将该事件标记为已解决。

行为：
- 默认情况下，使用 `alerting.pagerduty.integration-key` 作为集成密钥
- 如果被评估的端点属于某个组（`endpoints[].group`），且该组匹配 `alerting.pagerduty.overrides[].group` 的值，则提供者将使用该覆盖配置的集成密钥，而不是 `alerting.pagerduty.integration-key` 的值

```yaml
alerting:
  pagerduty:
    integration-key: "********************************"
    # 你也可以添加特定组的集成密钥，
    # 这将覆盖上面指定组的集成密钥
    overrides:
      - group: "core"
        integration-key: "********************************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: pagerduty
        failure-threshold: 3
        success-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"

  - name: back-end
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[CERTIFICATE_EXPIRATION] > 48h"
    alerts:
      - type: pagerduty
        failure-threshold: 3
        success-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Plivo 告警

> **警告**：此告警提供者尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供者是否按预期工作。感谢你的配合。

| 参数                                | 描述                                                                                     | 默认值         |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.plivo`                   | `plivo` 类型告警的配置                                                                      | `{}`          |
| `alerting.plivo.auth-id`           | Plivo 认证 ID                                                                              | 必填 `""`      |
| `alerting.plivo.auth-token`        | Plivo 认证令牌                                                                              | 必填 `""`      |
| `alerting.plivo.from`              | 发送短信的电话号码                                                                           | 必填 `""`      |
| `alerting.plivo.to`                | 接收短信的电话号码列表                                                                        | 必填 `[]`      |
| `alerting.plivo.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.plivo.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.plivo.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.plivo.overrides[].*`     | 参见 `alerting.plivo.*` 参数                                                                | `{}`          |

```yaml
alerting:
  plivo:
    auth-id: "MAXXXXXXXXXXXXXXXXXX"
    auth-token: "your-auth-token"
    from: "+1234567890"
    to:
      - "+0987654321"
      - "+1122334455"

endpoints:
  - name: website
    interval: 30s
    url: "https://twin.sh/health"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: plivo
        failure-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Pushover 告警
| 参数                                   | 描述                                                                                                       | 默认值                  |
|:--------------------------------------|:-------------------------------------------------------------------------------------------------------------|:----------------------|
| `alerting.pushover`                   | `pushover` 类型告警的配置                                                                                     | `{}`                  |
| `alerting.pushover.application-token` | Pushover 应用程序令牌                                                                                         | `""`                  |
| `alerting.pushover.user-key`          | 用户或群组密钥                                                                                                | `""`                  |
| `alerting.pushover.title`             | 通过 Pushover 发送的所有消息的固定标题                                                                          | `"Gatus: <endpoint>"` |
| `alerting.pushover.priority`          | 所有消息的优先级，范围从 -2（极低）到 2（紧急）                                                                   | `0`                   |
| `alerting.pushover.resolved-priority` | 已解决消息的优先级覆盖，范围从 -2（极低）到 2（紧急）                                                              | `0`                   |
| `alerting.pushover.sound`             | 所有消息的提示音<br />参见 [sounds](https://pushover.net/api#sounds) 了解所有有效选项。                            | `""`                  |
| `alerting.pushover.ttl`               | 设置消息的生存时间，超时后将自动从 Pushover 通知中删除                                                             | `0`                   |
| `alerting.pushover.device`            | 发送消息的目标设备（可选）<br/>参见 [devices](https://pushover.net/api#identifiers) 了解详情                        | `""` (所有设备)        |
| `alerting.pushover.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                                                | N/A                   |

```yaml
alerting:
  pushover:
    application-token: "******************************"
    user-key: "******************************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: pushover
        failure-threshold: 3
        success-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Rocket.Chat 告警

> **警告**：此告警提供者尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供者是否按预期工作。感谢你的配合。

| 参数                                     | 描述                                                                                     | 默认值         |
|:----------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.rocketchat`                   | `rocketchat` 类型告警的配置                                                                  | `{}`          |
| `alerting.rocketchat.webhook-url`       | Rocket.Chat 传入 Webhook URL                                                               | 必填 `""`      |
| `alerting.rocketchat.channel`           | 可选的频道覆盖                                                                               | `""`          |
| `alerting.rocketchat.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.rocketchat.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.rocketchat.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.rocketchat.overrides[].*`     | 参见 `alerting.rocketchat.*` 参数                                                           | `{}`          |

```yaml
alerting:
  rocketchat:
    webhook-url: "https://your-rocketchat.com/hooks/YOUR_WEBHOOK_ID/YOUR_TOKEN"
    channel: "#alerts"  # Optional

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: rocketchat
        send-on-resolved: true
```


#### 配置 SendGrid 告警

> **警告**：此告警提供者尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供者是否按预期工作。感谢你的配合。

| 参数                                   | 描述                                                                                     | 默认值         |
|:--------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.sendgrid`                   | `sendgrid` 类型告警的配置                                                                    | `{}`          |
| `alerting.sendgrid.api-key`           | SendGrid API 密钥                                                                          | 必填 `""`      |
| `alerting.sendgrid.from`              | 发件人电子邮件地址                                                                           | 必填 `""`      |
| `alerting.sendgrid.to`                | 接收告警的电子邮件地址（多个收件人用逗号分隔）                                                    | 必填 `""`      |
| `alerting.sendgrid.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.sendgrid.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.sendgrid.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.sendgrid.overrides[].*`     | 参见 `alerting.sendgrid.*` 参数                                                             | `{}`          |

```yaml
alerting:
  sendgrid:
    api-key: "SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    from: "alerts@example.com"
    to: "admin@example.com,ops@example.com"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: sendgrid
        send-on-resolved: true
```


#### 配置 Signal 告警

> **警告**：此告警提供者尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供者是否按预期工作。感谢你的配合。

| 参数                                 | 描述                                                                                     | 默认值         |
|:------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.signal`                   | `signal` 类型告警的配置                                                                     | `{}`          |
| `alerting.signal.api-url`           | Signal API URL（例如 signal-cli-rest-api 实例）                                              | 必填 `""`      |
| `alerting.signal.number`            | 发送者电话号码                                                                               | 必填 `""`      |
| `alerting.signal.recipients`        | 接收者电话号码列表                                                                           | 必填 `[]`      |
| `alerting.signal.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.signal.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.signal.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.signal.overrides[].*`     | 参见 `alerting.signal.*` 参数                                                               | `{}`          |

```yaml
alerting:
  signal:
    api-url: "http://localhost:8080"
    number: "+1234567890"
    recipients:
      - "+0987654321"
      - "+1122334455"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: signal
        send-on-resolved: true
```


#### 配置 SIGNL4 告警

SIGNL4 是一个移动告警和事件管理服务，通过移动推送、短信、语音呼叫和电子邮件向团队成员发送关键告警。

| 参数                                 | 描述                                                                                     | 默认值         |
|:------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.signl4`                   | `signl4` 类型告警的配置                                                                     | `{}`          |
| `alerting.signl4.team-secret`       | SIGNL4 团队密钥（Webhook URL 的一部分）                                                       | 必填 `""`      |
| `alerting.signl4.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.signl4.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.signl4.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.signl4.overrides[].*`     | 参见 `alerting.signl4.*` 参数                                                               | `{}`          |

```yaml
alerting:
  signl4:
    team-secret: "your-team-secret-here"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: signl4
        send-on-resolved: true
```


#### 配置 Slack 告警
| 参数                                | 描述                                                                                     | 默认值                                |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:------------------------------------|
| `alerting.slack`                   | `slack` 类型告警的配置                                                                      | `{}`                                |
| `alerting.slack.webhook-url`       | Slack Webhook URL                                                                          | 必填 `""`                            |
| `alerting.slack.title`             | 通知标题                                                                                    | `":helmet_with_white_cross: Gatus"` |
| `alerting.slack.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A                                 |
| `alerting.slack.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`                                |
| `alerting.slack.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`                                |
| `alerting.slack.overrides[].*`     | 参见 `alerting.slack.*` 参数                                                                | `{}`                                |

```yaml
alerting:
  slack:
    webhook-url: "https://hooks.slack.com/services/**********/**********/**********"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: slack
        description: "healthcheck failed 3 times in a row"
        send-on-resolved: true
      - type: slack
        failure-threshold: 5
        description: "healthcheck failed 5 times in a row"
        send-on-resolved: true
```

以下是通知的示例效果：

![Slack 通知](.github/assets/slack-alerts.png)


#### 配置 Splunk 告警

| 参数                                 | 描述                                                                                     | 默认值           |
|:------------------------------------|:-------------------------------------------------------------------------------------------|:----------------|
| `alerting.splunk`                   | `splunk` 类型告警的配置                                                                     | `{}`            |
| `alerting.splunk.hec-url`           | Splunk HEC（HTTP 事件收集器）URL                                                             | 必填 `""`        |
| `alerting.splunk.hec-token`         | Splunk HEC 令牌                                                                            | 必填 `""`        |
| `alerting.splunk.source`            | 事件来源                                                                                    | `"gatus"`       |
| `alerting.splunk.sourcetype`        | 事件来源类型                                                                                 | `"gatus:alert"` |
| `alerting.splunk.index`             | Splunk 索引                                                                                 | `""`            |
| `alerting.splunk.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A             |
| `alerting.splunk.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`            |
| `alerting.splunk.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`            |
| `alerting.splunk.overrides[].*`     | 参见 `alerting.splunk.*` 参数                                                               | `{}`            |

```yaml
alerting:
  splunk:
    hec-url: "https://splunk.example.com:8088"
    hec-token: "YOUR_HEC_TOKEN"
    index: "main"  # Optional

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: splunk
        send-on-resolved: true
```


#### 配置 Squadcast 告警

> **警告**：此告警提供者尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 Pull Request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 上评论该提供者是否按预期工作。感谢你的配合。

| 参数                                    | 描述                                                                                     | 默认值         |
|:---------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.squadcast`                   | `squadcast` 类型告警的配置                                                                   | `{}`          |
| `alerting.squadcast.webhook-url`       | Squadcast webhook URL                                                                      | 必填 `""`      |
| `alerting.squadcast.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.squadcast.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`          |
| `alerting.squadcast.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`          |
| `alerting.squadcast.overrides[].*`     | 参见 `alerting.squadcast.*` 参数                                                            | `{}`          |

```yaml
alerting:
  squadcast:
    webhook-url: "https://api.squadcast.com/v3/incidents/api/YOUR_API_KEY"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: squadcast
        send-on-resolved: true
```


#### 配置 Teams 告警 *（已弃用）*

> [!CAUTION]
> **已弃用：** Microsoft Teams 中的 Office 365 连接器正在停用（[来源：Microsoft DevBlog](https://devblogs.microsoft.com/microsoft365dev/retirement-of-office-365-connectors-within-microsoft-teams/)）。
> 现有连接器将继续工作到 2025 年 12 月。应使用新的 [Teams Workflow 告警](#configuring-teams-workflow-alerts) 配合 Microsoft Workflows，而不是此旧版配置。

| 参数                                | 描述                                                                                     | 默认值                |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:--------------------|
| `alerting.teams`                   | `teams` 类型告警的配置                                                                      | `{}`                |
| `alerting.teams.webhook-url`       | Teams Webhook URL                                                                          | 必填 `""`            |
| `alerting.teams.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A                 |
| `alerting.teams.title`             | 通知标题                                                                                    | `"&#x1F6A8; Gatus"` |
| `alerting.teams.client.insecure`   | 是否跳过 TLS 验证                                                                           | `false`             |
| `alerting.teams.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`                |
| `alerting.teams.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`                |
| `alerting.teams.overrides[].*`     | 参见 `alerting.teams.*` 参数                                                                | `{}`                |

```yaml
alerting:
  teams:
    webhook-url: "https://********.webhook.office.com/webhookb2/************"
    client:
      insecure: false
    # 你也可以添加特定组的密钥，
    # 这将覆盖上面指定组的密钥
    overrides:
      - group: "core"
        webhook-url: "https://********.webhook.office.com/webhookb3/************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: teams
        description: "healthcheck failed"
        send-on-resolved: true

  - name: back-end
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[CERTIFICATE_EXPIRATION] > 48h"
    alerts:
      - type: teams
        description: "healthcheck failed"
        send-on-resolved: true
```

以下是通知的示例效果：

![Teams 通知](.github/assets/teams-alerts.png)


#### 配置 Teams Workflow 告警

> [!NOTE]
> 此告警兼容 Microsoft Teams 的 Workflows。要设置工作流并获取 Webhook URL，请参阅 [Microsoft 文档](https://support.microsoft.com/en-us/office/create-incoming-webhooks-with-workflows-for-microsoft-teams-8ae491c7-0394-4861-ba59-055e33f75498)。

| 参数                                          | 描述                                                                                     | 默认值               |
|:---------------------------------------------|:-------------------------------------------------------------------------------------------|:-------------------|
| `alerting.teams-workflows`                   | `teams` 类型告警的配置                                                                      | `{}`               |
| `alerting.teams-workflows.webhook-url`       | Teams Webhook URL                                                                          | 必填 `""`           |
| `alerting.teams-workflows.title`             | 通知标题                                                                                    | `"&#x26D1; Gatus"` |
| `alerting.teams-workflows.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A                |
| `alerting.teams-workflows.overrides`         | 可以优先于默认配置的覆盖列表                                                                  | `[]`               |
| `alerting.teams-workflows.overrides[].group` | 将使用此配置覆盖默认配置的端点组                                                               | `""`               |
| `alerting.teams-workflows.overrides[].*`     | 参见 `alerting.teams-workflows.*` 参数                                                      | `{}`               |

```yaml
alerting:
  teams-workflows:
    webhook-url: "https://********.webhook.office.com/webhookb2/************"
    # 你也可以添加特定组的密钥，
    # 这将覆盖上面指定组的密钥
    overrides:
      - group: "core"
        webhook-url: "https://********.webhook.office.com/webhookb3/************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: teams-workflows
        description: "healthcheck failed"
        send-on-resolved: true

  - name: back-end
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[CERTIFICATE_EXPIRATION] > 48h"
    alerts:
      - type: teams-workflows
        description: "healthcheck failed"
        send-on-resolved: true
```

以下是通知的示例效果：

![Teams Workflow 通知](.github/assets/teams-workflows-alerts.png)
```yaml
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: teams
        description: "healthcheck failed"
        send-on-resolved: true

  - name: back-end
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[CERTIFICATE_EXPIRATION] > 48h"
    alerts:
      - type: teams
        description: "healthcheck failed"
        send-on-resolved: true
```

以下是通知效果的示例：

![Teams 通知](.github/assets/teams-alerts.png)


#### 配置 Teams Workflow 告警

> [!NOTE]
> 此告警兼容 Microsoft Teams 的 Workflows。要设置工作流并获取 webhook URL，请参阅 [Microsoft 文档](https://support.microsoft.com/en-us/office/create-incoming-webhooks-with-workflows-for-microsoft-teams-8ae491c7-0394-4861-ba59-055e33f75498)。

| 参数                                         | 描述                                                                                       | 默认值             |
|:---------------------------------------------|:-------------------------------------------------------------------------------------------|:-------------------|
| `alerting.teams-workflows`                   | `teams` 类型告警的配置                                                                      | `{}`               |
| `alerting.teams-workflows.webhook-url`       | Teams Webhook URL                                                                          | 必填 `""`          |
| `alerting.teams-workflows.title`             | 通知标题                                                                                    | `"&#x26D1; Gatus"` |
| `alerting.teams-workflows.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                             | N/A                |
| `alerting.teams-workflows.overrides`         | 可优先于默认配置的覆盖列表                                                                   | `[]`               |
| `alerting.teams-workflows.overrides[].group` | 将被此配置覆盖的端点组                                                                       | `""`               |
| `alerting.teams-workflows.overrides[].*`     | 参见 `alerting.teams-workflows.*` 参数                                                      | `{}`               |

```yaml
alerting:
  teams-workflows:
    webhook-url: "https://********.webhook.office.com/webhookb2/************"
    # 你也可以添加特定组的 key，这将
    # 覆盖上面指定组的 key
    overrides:
      - group: "core"
        webhook-url: "https://********.webhook.office.com/webhookb3/************"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: teams-workflows
        description: "healthcheck failed"
        send-on-resolved: true

  - name: back-end
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[CERTIFICATE_EXPIRATION] > 48h"
    alerts:
      - type: teams-workflows
        description: "healthcheck failed"
        send-on-resolved: true
```

以下是通知效果的示例：

![Teams Workflow 通知](.github/assets/teams-workflows-alerts.png)


#### 配置 Telegram 告警
| 参数                                  | 描述                                                                                       | 默认值                     |
|:--------------------------------------|:-------------------------------------------------------------------------------------------|:---------------------------|
| `alerting.telegram`                   | `telegram` 类型告警的配置                                                                    | `{}`                       |
| `alerting.telegram.token`             | Telegram 机器人 Token                                                                       | 必填 `""`                  |
| `alerting.telegram.id`                | Telegram 聊天 ID                                                                            | 必填 `""`                  |
| `alerting.telegram.topic-id`          | 群组中的 Telegram 话题 ID，对应 Telegram API 中的 `message_thread_id`                         | `""`                       |
| `alerting.telegram.api-url`           | Telegram API URL                                                                           | `https://api.telegram.org` |
| `alerting.telegram.client`            | 客户端配置。<br />参见 [客户端配置](#client-configuration)。                                   | `{}`                       |
| `alerting.telegram.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A                        |
| `alerting.telegram.overrides`         | 可优先于默认配置的覆盖列表                                                                    | `[]`                       |
| `alerting.telegram.overrides[].group` | 将被此配置覆盖的端点组                                                                        | `""`                       |
| `alerting.telegram.overrides[].*`     | 参见 `alerting.telegram.*` 参数                                                              | `{}`                       |

```yaml
alerting:
  telegram:
    token: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
    id: "0123456789"
    topic-id: "7"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
    alerts:
      - type: telegram
        send-on-resolved: true
```

以下是通知效果的示例：

![Telegram 通知](.github/assets/telegram-alerts.png)


#### 配置 Twilio 告警
| 参数                            | 描述                                                                                       | 默认值        |
|:--------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.twilio`               | `twilio` 类型告警的设置                                                                      | `{}`          |
| `alerting.twilio.sid`           | Twilio 账户 SID                                                                             | 必填 `""`     |
| `alerting.twilio.token`         | Twilio 认证令牌                                                                              | 必填 `""`     |
| `alerting.twilio.from`          | 发送 Twilio 告警的号码                                                                       | 必填 `""`     |
| `alerting.twilio.to`            | 接收 Twilio 告警的号码                                                                       | 必填 `""`     |
| `alerting.twilio.default-alert` | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |

通过以下附加参数支持自定义消息模板：

| 参数                                    | 描述                                                                                       | 默认值  |
|:----------------------------------------|:-------------------------------------------------------------------------------------------|:--------|
| `alerting.twilio.text-twilio-triggered` | 触发告警的自定义消息模板。支持 `[ENDPOINT]`、`[ALERT_DESCRIPTION]`                            | `""`    |
| `alerting.twilio.text-twilio-resolved`  | 恢复告警的自定义消息模板。支持 `[ENDPOINT]`、`[ALERT_DESCRIPTION]`                            | `""`    |

```yaml
alerting:
  twilio:
    sid: "..."
    token: "..."
    from: "+1-234-567-8901"
    to: "+1-234-567-8901"
    # 使用占位符的自定义消息模板（可选）
    # 同时支持旧格式 {endpoint}/{description} 和新格式 [ENDPOINT]/[ALERT_DESCRIPTION]
    text-twilio-triggered: "🚨 ALERT: [ENDPOINT] is down! [ALERT_DESCRIPTION]"
    text-twilio-resolved: "✅ RESOLVED: [ENDPOINT] is back up! [ALERT_DESCRIPTION]"

endpoints:
  - name: website
    interval: 30s
    url: "https://twin.sh/health"
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: twilio
        failure-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Vonage 告警

> ⚠️ **警告**：此告警提供方尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 pull request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 中评论该提供方是否按预期工作。感谢你的配合。

| 参数                                | 描述                                                                                       | 默认值        |
|:------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.vonage`                   | `vonage` 类型告警的配置                                                                      | `{}`          |
| `alerting.vonage.api-key`           | Vonage API 密钥                                                                             | 必填 `""`     |
| `alerting.vonage.api-secret`        | Vonage API 密钥                                                                             | 必填 `""`     |
| `alerting.vonage.from`              | 发送者名称或电话号码                                                                          | 必填 `""`     |
| `alerting.vonage.to`                | 接收者电话号码                                                                                | 必填 `""`     |
| `alerting.vonage.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.vonage.overrides`         | 可优先于默认配置的覆盖列表                                                                    | `[]`          |
| `alerting.vonage.overrides[].group` | 将被此配置覆盖的端点组                                                                        | `""`          |
| `alerting.vonage.overrides[].*`     | 参见 `alerting.vonage.*` 参数                                                                | `{}`          |

```yaml
alerting:
  vonage:
    api-key: "YOUR_API_KEY"
    api-secret: "YOUR_API_SECRET"
    from: "Gatus"
    to: "+1234567890"
```

发送告警到 Vonage 的示例：
```yaml
endpoints:
  - name: website
    url: "https://example.org"
    alerts:
      - type: vonage
        failure-threshold: 5
        send-on-resolved: true
        description: "healthcheck failed"
```


#### 配置 Webex 告警

> ⚠️ **警告**：此告警提供方尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 pull request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 中评论该提供方是否按预期工作。感谢你的配合。

| 参数                               | 描述                                                                                       | 默认值        |
|:-----------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.webex`                   | `webex` 类型告警的配置                                                                       | `{}`          |
| `alerting.webex.webhook-url`       | Webex Teams webhook URL                                                                    | 必填 `""`     |
| `alerting.webex.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.webex.overrides`         | 可优先于默认配置的覆盖列表                                                                    | `[]`          |
| `alerting.webex.overrides[].group` | 将被此配置覆盖的端点组                                                                        | `""`          |
| `alerting.webex.overrides[].*`     | 参见 `alerting.webex.*` 参数                                                                 | `{}`          |

```yaml
alerting:
  webex:
    webhook-url: "https://webexapis.com/v1/webhooks/incoming/YOUR_WEBHOOK_ID"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: webex
        send-on-resolved: true
```


#### 配置 Zapier 告警

> ⚠️ **警告**：此告警提供方尚未经过测试。如果你已测试并确认其正常工作，请移除此警告并创建一个 pull request，或在 [#1223](https://github.com/TwiN/gatus/discussions/1223) 中评论该提供方是否按预期工作。感谢你的配合。

| 参数                                | 描述                                                                                       | 默认值        |
|:--------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.zapier`               | `zapier` 类型告警的配置                                                                      | `{}`          |
| `alerting.zapier.webhook-url`   | Zapier webhook URL                                                                         | 必填 `""`     |
| `alerting.zapier.default-alert` | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.zapier.overrides`     | 可优先于默认配置的覆盖列表                                                                    | `[]`          |
| `alerting.zapier.overrides[].group` | 将被此配置覆盖的端点组                                                                    | `""`          |
| `alerting.zapier.overrides[].*` | 参见 `alerting.zapier.*` 参数                                                               | `{}`          |

```yaml
alerting:
  zapier:
    webhook-url: "https://hooks.zapier.com/hooks/catch/YOUR_WEBHOOK_ID/"

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: zapier
        send-on-resolved: true
```


#### 配置 Zulip 告警
| 参数                               | 描述                                                                        | 默认值        |
|:-----------------------------------|:----------------------------------------------------------------------------|:--------------|
| `alerting.zulip`                   | `zulip` 类型告警的配置                                                       | `{}`          |
| `alerting.zulip.bot-email`         | 机器人邮箱                                                                   | 必填 `""`     |
| `alerting.zulip.bot-api-key`       | 机器人 API 密钥                                                              | 必填 `""`     |
| `alerting.zulip.domain`            | 完整的组织域名（例如：yourZulipDomain.zulipchat.com）                         | 必填 `""`     |
| `alerting.zulip.channel-id`        | Gatus 发送告警的频道 ID                                                      | 必填 `""`     |
| `alerting.zulip.overrides`         | 可优先于默认配置的覆盖列表                                                    | `[]`          |
| `alerting.zulip.overrides[].group` | 将被此配置覆盖的端点组                                                        | `""`          |
| `alerting.zulip.overrides[].*`     | 参见 `alerting.zulip.*` 参数                                                 | `{}`          |

```yaml
alerting:
  zulip:
    bot-email: gatus-bot@some.zulip.org
    bot-api-key: "********************************"
    domain: some.zulip.org
    channel-id: 123456

endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: zulip
        description: "healthcheck failed"
        send-on-resolved: true
```


#### 配置自定义告警
| 参数                                | 描述                                                                                       | 默认值        |
|:------------------------------------|:-------------------------------------------------------------------------------------------|:--------------|
| `alerting.custom`                   | 失败时自定义动作或告警的配置                                                                  | `{}`          |
| `alerting.custom.url`               | 自定义告警请求 URL                                                                           | 必填 `""`     |
| `alerting.custom.method`            | 请求方法                                                                                     | `GET`         |
| `alerting.custom.body`              | 自定义告警请求正文                                                                            | `""`          |
| `alerting.custom.headers`           | 自定义告警请求头                                                                              | `{}`          |
| `alerting.custom.client`            | 客户端配置。<br />参见 [客户端配置](#client-configuration)。                                   | `{}`          |
| `alerting.custom.default-alert`     | 默认告警配置。<br />参见 [设置默认告警](#setting-a-default-alert)                              | N/A           |
| `alerting.custom.overrides`         | 可优先于默认配置的覆盖列表                                                                    | `[]`          |
| `alerting.custom.overrides[].group` | 将被此配置覆盖的端点组                                                                        | `""`          |
| `alerting.custom.overrides[].*`     | 参见 `alerting.custom.*` 参数                                                                | `{}`          |

虽然它们被称为告警，但你可以使用此功能调用任何内容。

例如，你可以通过让一个应用程序跟踪新部署来自动化回滚，借助 Gatus，你可以让 Gatus 在端点开始失败时调用该应用程序端点。你的应用程序随后会检查开始失败的端点是否属于最近部署的应用程序的一部分，如果是，则自动回滚。

此外，你可以在请求正文（`alerting.custom.body`）和 URL（`alerting.custom.url`）中使用以下占位符：
- `[ALERT_DESCRIPTION]`（从 `endpoints[].alerts[].description` 解析）
- `[ENDPOINT_NAME]`（从 `endpoints[].name` 解析）
- `[ENDPOINT_GROUP]`（从 `endpoints[].group` 解析）
- `[ENDPOINT_URL]`（从 `endpoints[].url` 解析）
- `[RESULT_ERRORS]`（从给定健康检查的健康评估中解析）
- `[RESULT_CONDITIONS]`（从给定健康检查的健康评估中的条件结果解析）
-
如果你使用 `custom` 提供方且 `send-on-resolved` 设置为 `true` 的告警，可以使用
`[ALERT_TRIGGERED_OR_RESOLVED]` 占位符来区分通知。
上述占位符将根据情况被替换为 `TRIGGERED` 或 `RESOLVED`，但可以修改
（详情见本节末尾）。

出于所有目的，我们将使用 Slack webhook 配置自定义告警，但你可以调用任何你想要的内容。
```yaml
alerting:
  custom:
    url: "https://hooks.slack.com/services/**********/**********/**********"
    method: "POST"
    body: |
      {
        "text": "[ALERT_TRIGGERED_OR_RESOLVED]: [ENDPOINT_GROUP] - [ENDPOINT_NAME] - [ALERT_DESCRIPTION] - [RESULT_ERRORS]"
      }
endpoints:
  - name: website
    url: "https://twin.sh/health"
    interval: 30s
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 300"
    alerts:
      - type: custom
        failure-threshold: 10
        success-threshold: 3
        send-on-resolved: true
        description: "health check failed"
```

请注意，你可以像这样自定义 `[ALERT_TRIGGERED_OR_RESOLVED]` 占位符的解析值：
```yaml
alerting:
  custom:
    placeholders:
      ALERT_TRIGGERED_OR_RESOLVED:
        TRIGGERED: "partial_outage"
        RESOLVED: "operational"
```
因此，本节第一个示例中请求正文里的 `[ALERT_TRIGGERED_OR_RESOLVED]` 在告警触发时将被替换为
`partial_outage`，在告警恢复时将被替换为 `operational`。


#### 设置默认告警
| 参数                                         | 描述                                                                  | 默认值  |
|:---------------------------------------------|:----------------------------------------------------------------------|:--------|
| `alerting.*.default-alert.enabled`           | 是否启用告警                                                          | N/A     |
| `alerting.*.default-alert.failure-threshold` | 触发告警前需要连续失败的次数                                           | N/A     |
| `alerting.*.default-alert.success-threshold` | 将进行中的事件标记为已恢复前需要连续成功的次数                          | N/A     |
| `alerting.*.default-alert.send-on-resolved`  | 触发的告警被标记为已恢复后是否发送通知                                  | N/A     |
| `alerting.*.default-alert.description`       | 告警描述。将包含在发送的告警中                                         | N/A     |

> ⚠ 即使你设置了提供方的默认告警，你仍然必须在端点配置中指定告警的 `type`。

虽然你可以直接在端点定义中指定告警配置，但这很繁琐，可能导致配置文件非常长。

为避免此问题，你可以使用每个提供方配置中的 `default-alert` 参数：
```yaml
alerting:
  slack:
    webhook-url: "https://hooks.slack.com/services/**********/**********/**********"
    default-alert:
      description: "health check failed"
      send-on-resolved: true
      failure-threshold: 5
      success-threshold: 5
```

这样，你的 Gatus 配置看起来更加整洁：
```yaml
endpoints:
  - name: example
    url: "https://example.org"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack

  - name: other-example
    url: "https://example.com"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack
```

它还允许你做这样的事情：
```yaml
endpoints:
  - name: example
    url: "https://example.org"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack
        failure-threshold: 5
      - type: slack
        failure-threshold: 10
      - type: slack
        failure-threshold: 15
```

当然，你也可以混合使用告警类型：
```yaml
alerting:
  slack:
    webhook-url: "https://hooks.slack.com/services/**********/**********/**********"
    default-alert:
      failure-threshold: 3
  pagerduty:
    integration-key: "********************************"
    default-alert:
      failure-threshold: 5

endpoints:
  - name: endpoint-1
    url: "https://example.org"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack
      - type: pagerduty

  - name: endpoint-2
    url: "https://example.org"
    conditions:
      - "[STATUS] == 200"
    alerts:
      - type: slack
      - type: pagerduty
```


### 维护
如果你有维护窗口，可能不希望被告警打扰。
为此，你需要使用维护配置：

| 参数                   | 描述                                                                                                                                                                                       | 默认值        |
|:-----------------------|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------|
| `maintenance.enabled`  | 是否启用维护期                                                                                                                                                                             | `true`        |
| `maintenance.start`    | 维护窗口开始时间，格式为 `hh:mm`（例如 `23:00`）                                                                                                                                            | 必填 `""`     |
| `maintenance.duration` | 维护窗口持续时间（例如 `1h`、`30m`）                                                                                                                                                        | 必填 `""`     |
| `maintenance.timezone` | 维护窗口的时区格式（例如 `Europe/Amsterdam`）。<br />更多信息请参见 [tz 数据库时区列表](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)                                          | `UTC`         |
| `maintenance.every`    | 维护期生效的日期（例如 `[Monday, Thursday]`）。<br />如果留空，维护窗口每天生效                                                                                                               | `[]`          |

以下是一个示例：
```yaml
maintenance:
  start: 23:00
  duration: 1h
  timezone: "Europe/Amsterdam"
  every: [Monday, Thursday]
```
请注意，你也可以在单独的行中指定每一天：
```yaml
maintenance:
  start: 23:00
  duration: 1h
  timezone: "Europe/Amsterdam"
  every:
    - Monday
    - Thursday
```
你也可以按端点指定维护窗口：
```yaml
endpoints:
  - name: endpoint-1
    url: "https://example.org"
    maintenance-windows:
      - start: "07:30"
        duration: 40m
        timezone: "Europe/Berlin"
      - start: "14:30"
        duration: 1h
        timezone: "Europe/Berlin"
```


### 安全
| 参数             | 描述                         | 默认值  |
|:-----------------|:-----------------------------|:--------|
| `security`       | 安全配置                     | `{}`    |
| `security.basic` | HTTP Basic 认证配置           | `{}`    |
| `security.oidc`  | OpenID Connect 配置           | `{}`    |
| `security.authorization` | 端点/套件状态 API 的组级授权配置 | `{}` |


#### 授权
| 参数                                     | 描述                                                                                        | 默认值  |
|:-----------------------------------------|:------------------------------------------------------------------------------------------------|:--------|
| `security.authorization.endpoint-groups` | 可在 `/api/v1/endpoints/*` 和 `/api/v1/groups` 中访问的端点组列表。为空表示允许所有组。 | `[]` |
| `security.authorization.suite-groups`    | 可在 `/api/v1/suites/*` 和 `/api/v1/groups` 中访问的套件组列表。为空表示允许所有组。      | `[]` |

```yaml
security:
  authorization:
    endpoint-groups: ["core", "partner"]
    suite-groups: ["smoke"]
```


#### Basic 认证
| 参数                                    | 描述                                                                       | 默认值        |
|:----------------------------------------|:---------------------------------------------------------------------------|:--------------|
| `security.basic`                        | HTTP Basic 认证配置                                                         | `{}`          |
| `security.basic.username`               | Basic 认证的用户名                                                          | 必填 `""`     |
| `security.basic.password-bcrypt-base64` | 使用 Bcrypt 哈希然后用 base64 编码的 Basic 认证密码                          | 必填 `""`     |

以下示例要求你使用用户名 `john.doe` 和密码 `hunter2` 进行认证：
```yaml
security:
  basic:
    username: "john.doe"
    password-bcrypt-base64: "JDJhJDEwJHRiMnRFakxWazZLdXBzRERQazB1TE8vckRLY05Yb1hSdnoxWU0yQ1FaYXZRSW1McmladDYu"
```

> ⚠ 请确保仔细选择 bcrypt 哈希的成本。成本越高，计算哈希所需的时间越长，
> 而 basic 认证会在每次请求时验证密码与哈希的匹配。截至 2023-01-06，我建议成本设为 9。


#### OIDC
| 参数                             | 描述                                                           | 默认值        |
|:---------------------------------|:---------------------------------------------------------------|:--------------|
| `security.oidc`                  | OpenID Connect 配置                                             | `{}`          |
| `security.oidc.issuer-url`       | 发行者 URL                                                      | 必填 `""`     |
| `security.oidc.redirect-url`     | 重定向 URL。必须以 `/authorization-code/callback` 结尾           | 必填 `""`     |
| `security.oidc.client-id`        | 客户端 ID                                                       | 必填 `""`     |
| `security.oidc.client-secret`    | 客户端密钥                                                       | 必填 `""`     |
| `security.oidc.scopes`           | 请求的范围。你唯一需要的范围是 `openid`。                         | 必填 `[]`     |
| `security.oidc.allowed-subjects` | 允许的主体列表。如果为空，则允许所有主体。                         | `[]`          |
| `security.oidc.session-ttl`      | 会话生存时间（例如 `8h`、`1h30m`、`2h`）。                        | `8h`          |

```yaml
security:
  oidc:
    issuer-url: "https://example.okta.com"
    redirect-url: "https://status.example.com/authorization-code/callback"
    client-id: "123456789"
    client-secret: "abcdefghijk"
    scopes: ["openid"]
    # 你可以选择性地指定允许的主体列表。如果未指定，则允许所有主体。
    #allowed-subjects: ["johndoe@example.com"]
    # 你可以选择性地指定会话生存时间。如果未指定，默认为 8 小时。
    #session-ttl: 8h
```

感到困惑？请阅读 [使用 Auth0 通过 OIDC 保护 Gatus](https://twin.sh/articles/56/securing-gatus-with-oidc-using-auth0)。


### TLS 加密
Gatus 支持使用 TLS 进行基本加密。要启用此功能，必须提供 PEM 格式的证书文件。

以下示例展示了一个使 Gatus 在 4443 端口响应 HTTPS 请求的配置示例：
```yaml
web:
  port: 4443
  tls:
    certificate-file: "certificate.crt"
    private-key-file: "private.key"
```


### 指标
要启用指标，你必须将 `metrics` 设置为 `true`。这样做将在你的应用程序配置运行的同一端口（`web.port`）上的 `/metrics` 端点暴露 Prometheus 友好的指标。

| 指标名称                                         | 类型    | 描述                                                                   | 标签                            | 相关端点类型        |
|:---------------------------------------------|:--------|:---------------------------------------------------------------------------|:--------------------------------|:------------------------|
| gatus_results_total                          | counter | 每个端点每种成功状态的结果数                                              | key, group, name, type, success | 全部                    |
| gatus_results_code_total                     | counter | 按状态码的结果总数                                                        | key, group, name, type, code    | DNS, HTTP               |
| gatus_results_connected_total                | counter | 成功建立连接的结果总数                                                    | key, group, name, type          | 全部                    |
| gatus_results_duration_seconds               | gauge   | 请求持续时间（秒）                                                        | key, group, name, type          | 全部                    |
| gatus_results_certificate_expiration_seconds | gauge   | 证书到期前的秒数                                                          | key, group, name, type          | HTTP, STARTTLS          |
| gatus_results_domain_expiration_seconds      | gauge   | 域名到期前的秒数                                                          | key, group, name, type          | HTTP, STARTTLS          |
| gatus_results_endpoint_success               | gauge   | 显示端点是否成功（0 失败，1 成功）                                        | key, group, name, type          | 全部                    |

更多文档和示例请参见 [examples/docker-compose-grafana-prometheus](.examples/docker-compose-grafana-prometheus)。

#### 自定义标签

你可以通过在 `extra-labels` 字段下定义键值对来为端点的 Prometheus 指标添加自定义标签。例如：

```yaml
endpoints:
  - name: front-end
    group: core
    url: "https://twin.sh/health"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
      - "[BODY].status == UP"
      - "[RESPONSE_TIME] < 150"
    extra-labels:
      environment: staging
```

### 连通性
| 参数                            | 描述                               | 默认值        |
|:--------------------------------|:-------------------------------------------|:--------------|
| `connectivity`                  | 连通性配置                                  | `{}`          |
| `connectivity.checker`          | 连通性检查器配置                             | 必填 `{}`     |
| `connectivity.checker.target`   | 用于验证连通性的目标主机                      | 必填 `""`     |
| `connectivity.checker.interval` | 验证连通性的间隔                              | `1m`          |

虽然 Gatus 用于监控其他服务，但 Gatus 本身也可能失去与互联网的连接。
为了防止 Gatus 在自身不健康时将端点报告为不健康，你可以配置
Gatus 定期检查互联网连通性。

当连通性检查器判断连接已断开时，所有端点执行将被跳过。

```yaml
connectivity:
  checker:
    target: 1.1.1.1:53
    interval: 60s
```


### 远程实例（实验性）
此功能允许你从远程 Gatus 实例获取端点状态。

此功能有两个主要用例：
- 你有多个运行在不同机器上的 Gatus 实例，希望通过单个仪表板以可视化方式展示状态
- 你有一个或多个不可公开访问的 Gatus 实例（例如在防火墙后面），希望获取状态

这是一个实验性功能。它可能随时被移除或以破坏性方式更新。此外，
此功能存在已知问题。如果你想提供反馈，请在 [#64](https://github.com/TwiN/gatus/issues/64) 中发表评论。
使用风险自负。

| 参数                               | 描述                                 | 默认值        |
|:-----------------------------------|:---------------------------------------------|:--------------|
| `remote`                           | 远程配置                                      | `{}`          |
| `remote.instances`                 | 远程实例列表                                   | 必填 `[]`     |
| `remote.instances.endpoint-prefix` | 所有端点名称的前缀字符串                        | `""`          |
| `remote.instances.url`             | 获取端点状态的 URL                              | 必填 `""`     |

```yaml
remote:
  instances:
    - endpoint-prefix: "status.example.org-"
      url: "https://status.example.org/api/v1/endpoints/statuses"
```


## 部署
许多示例可以在 [.examples](.examples) 文件夹中找到，但本节将重点介绍最流行的 Gatus 部署方式。


### Docker
使用 Docker 在本地运行 Gatus：
```console
docker run -p 8080:8080 --name gatus ghcr.io/twin/gatus:stable
```

除了使用 [.examples](.examples) 文件夹中提供的示例外，你还可以通过
创建一个配置文件（在此示例中我们将其命名为 `config.yaml`）并运行以下
命令来在本地试用：
```console
docker run -p 8080:8080 --mount type=bind,source="$(pwd)"/config.yaml,target=/config/config.yaml --name gatus ghcr.io/twin/gatus:stable
```

如果你使用的是 Windows，请将 `"$(pwd)"` 替换为当前目录的绝对路径，例如：
```console
docker run -p 8080:8080 --mount type=bind,source=C:/Users/Chris/Desktop/config.yaml,target=/config/config.yaml --name gatus ghcr.io/twin/gatus:stable
```

在本地构建镜像：
```console
docker build . -t ghcr.io/twin/gatus:stable
```


### Helm Chart
使用该 chart 必须安装 [Helm](https://helm.sh)。
请参阅 Helm 的 [文档](https://helm.sh/docs/) 开始使用。

Helm 正确设置后，按如下方式添加仓库：

```console
helm repo add twin https://twin.github.io/helm-charts
helm repo update
helm install gatus twin/gatus
```

要获取更多详细信息，请查看 [chart 的配置](https://github.com/TwiN/helm-charts/blob/master/charts/gatus/README.md)。


### Terraform

#### Kubernetes

Gatus 可以使用以下模块通过 Terraform 部署到 Kubernetes：[terraform-kubernetes-gatus](https://github.com/TwiN/terraform-kubernetes-gatus)。

## 运行测试
```console
go test -v ./...
```


## 在生产环境中使用
请参阅 [部署](#deployment) 章节。


## 常见问题
### 发送 GraphQL 请求
通过将 `endpoints[].graphql` 设置为 true，请求正文将自动被标准 GraphQL `query` 参数包装。

例如，以下配置：
```yaml
endpoints:
  - name: filter-users-by-gender
    url: http://localhost:8080/playground
    method: POST
    graphql: true
    body: |
      {
        users(gender: "female") {
          id
          name
          gender
          avatar
        }
      }
    conditions:
      - "[STATUS] == 200"
      - "[BODY].data.users[0].gender == female"
```

将发送一个 `POST` 请求到 `http://localhost:8080/playground`，请求正文如下：
```json
{"query":"      {\n        users(gender: \"female\") {\n          id\n          name\n          gender\n          avatar\n        }\n      }"}
```


### 推荐间隔
为确保 Gatus 提供可靠和准确的结果（即响应时间），Gatus 限制了可以同时评估的端点/套件数量。
换句话说，即使你有多个具有相同间隔的端点，它们也不保证会同时运行。

并发评估的数量由 `concurrency` 配置参数决定，默认为 `3`。

你可以自己测试，通过运行配置了多个具有非常短且不切实际间隔（如 1ms）的端点的 Gatus。你会注意到响应时间不会波动——这是因为虽然端点在不同的 goroutine 上评估，但有一个信号量控制着同时运行的端点/套件数量。

不幸的是，这有一个缺点。如果你有大量端点，其中一些非常慢或容易超时（默认超时为 10 秒），那些慢的评估可能会阻止其他端点/套件被评估。

间隔不包括请求本身的持续时间，这意味着如果一个端点的间隔为 30 秒，而请求需要 2 秒才能完成，那么两次评估之间的时间戳将是 32 秒，而不是 30 秒。

虽然这不会阻止 Gatus 对所有其他端点执行健康检查，但可能导致 Gatus 无法遵守配置的间隔，例如，假设 `concurrency` 设置为 `1`：
- 端点 A 的间隔为 5 秒，超时后需要 10 秒才能完成
- 端点 B 的间隔为 5 秒，只需要 1 毫秒即可完成
- 端点 B 将无法每 5 秒运行一次，因为端点 A 的健康评估时间超过了其间隔

总而言之，虽然 Gatus 可以处理你设置的任何间隔，但对于慢请求最好使用更长的间隔。

根据经验，我个人将更复杂的健康检查间隔设置为 `5m`（5 分钟），将用于告警（PagerDuty/Twilio）的简单健康检查设置为 `30s`。


### 默认超时
| 端点类型      | 超时    |
|:--------------|:--------|
| HTTP          | 10s     |
| TCP           | 10s     |
| ICMP          | 10s     |

要修改超时，请参见 [客户端配置](#client-configuration)。


### 监控 TCP 端点
通过为 `endpoints[].url` 添加 `tcp://` 前缀，你可以在非常基本的层面上监控 TCP 端点：
```yaml
endpoints:
  - name: redis
    url: "tcp://127.0.0.1:6379"
    interval: 30s
    conditions:
      - "[CONNECTED] == true"
```
如果设置了 `endpoints[].body`，它将被发送，响应的前 1024 字节将在 `[BODY]` 中。

占位符 `[STATUS]` 以及字段 `endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持 TCP 端点。

这适用于数据库（Postgres、MySQL 等）和缓存（Redis、Memcached 等）等应用程序。

> 📝 `[CONNECTED] == true` 不保证端点本身是健康的——它只保证在给定地址的给定端口上有
> 某些东西在监听，并且到该地址的连接已成功建立。


### 监控 UDP 端点
通过为 `endpoints[].url` 添加 `udp://` 前缀，你可以在非常基本的层面上监控 UDP 端点：
```yaml
endpoints:
  - name: example
    url: "udp://example.org:80"
    conditions:
      - "[CONNECTED] == true"
```

如果设置了 `endpoints[].body`，它将被发送，响应的前 1024 字节将在 `[BODY]` 中。

占位符 `[STATUS]` 以及字段 `endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持 UDP 端点。

这适用于基于 UDP 的应用程序。


### 监控 SCTP 端点
通过为 `endpoints[].url` 添加 `sctp://` 前缀，你可以在非常基本的层面上监控流控制传输协议（SCTP）端点：
```yaml
endpoints:
  - name: example
    url: "sctp://127.0.0.1:38412"
    conditions:
      - "[CONNECTED] == true"
```

占位符 `[STATUS]` 和 `[BODY]` 以及字段 `endpoints[].body`、`endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持 SCTP 端点。

这适用于基于 SCTP 的应用程序。
| ICMP          | 10s     |

要修改超时时间，请参阅[客户端配置](#client-configuration)。


### 监控 TCP 端点
通过在 `endpoints[].url` 前添加 `tcp://` 前缀，您可以在非常基础的层面上监控 TCP 端点：
```yaml
endpoints:
  - name: redis
    url: "tcp://127.0.0.1:6379"
    interval: 30s
    conditions:
      - "[CONNECTED] == true"
```
如果设置了 `endpoints[].body`，则会发送该内容，响应的前 1024 字节将存储在 `[BODY]` 中。

占位符 `[STATUS]` 以及字段 `endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持用于 TCP 端点。

这适用于数据库（Postgres、MySQL 等）和缓存（Redis、Memcached 等）等应用。

> 📝 `[CONNECTED] == true` 并不保证端点本身是健康的——它只保证在给定地址的给定端口上有某个服务在监听，
> 并且已成功建立到该地址的连接。


### 监控 UDP 端点
通过在 `endpoints[].url` 前添加 `udp://` 前缀，您可以在非常基础的层面上监控 UDP 端点：
```yaml
endpoints:
  - name: example
    url: "udp://example.org:80"
    conditions:
      - "[CONNECTED] == true"
```

如果设置了 `endpoints[].body`，则会发送该内容，响应的前 1024 字节将存储在 `[BODY]` 中。

占位符 `[STATUS]` 以及字段 `endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持用于 UDP 端点。

这适用于基于 UDP 的应用。


### 监控 SCTP 端点
通过在 `endpoints[].url` 前添加 `sctp://` 前缀，您可以在非常基础的层面上监控流控制传输协议（SCTP）端点：
```yaml
endpoints:
  - name: example
    url: "sctp://127.0.0.1:38412"
    conditions:
      - "[CONNECTED] == true"
```

占位符 `[STATUS]` 和 `[BODY]` 以及字段 `endpoints[].body`、`endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持用于 SCTP 端点。

这适用于基于 SCTP 的应用。


### 监控 WebSocket 端点
通过在 `endpoints[].url` 前添加 `ws://` 或 `wss://` 前缀，您可以监控 WebSocket 端点：
```yaml
endpoints:
  - name: example
    url: "wss://echo.websocket.org/"
    body: "status"
    conditions:
      - "[CONNECTED] == true"
      - "[BODY] == pat(*served by*)"
```

`[BODY]` 占位符包含查询的输出，`[CONNECTED]`
表示连接是否已成功建立。您可以使用 Go 模板语法。


### 使用 gRPC 监控端点
您可以通过在 `endpoints[].url` 前添加 `grpc://` 或 `grpcs://` 前缀来监控 gRPC 服务。
Gatus 会对目标执行标准的 `grpc.health.v1.Health/Check` RPC 调用。

```yaml
endpoints:
  - name: my-grpc
    url: grpc://localhost:50051
    interval: 30s
    conditions:
      - "[CONNECTED] == true"
      - "[BODY].status == SERVING"  # BODY 仅在被引用时才读取
    client:
      timeout: 5s
```

对于启用了 TLS 的服务器，使用 `grpcs://` 并根据需要配置客户端 TLS：

```yaml
endpoints:
  - name: my-grpcs
    url: grpcs://example.com:443
    conditions:
      - "[CONNECTED] == true"
      - "[BODY].status == SERVING"
    client:
      timeout: 5s
      insecure: false          # 设置为 true 可跳过证书验证（不推荐）
      tls:
        certificate-file: /path/to/cert.pem      # 可选的 mTLS 客户端证书
        private-key-file: /path/to/key.pem       # 可选的 mTLS 客户端密钥
```

注意事项：
- 健康检查针对默认服务（`service: ""`）。如有需要，后续可以添加对自定义服务名称的支持。
- 响应体仅在条件或测试套件存储映射需要时，才会以最小 JSON 对象的形式暴露，如 `{"status":"SERVING"}`。
- 超时、自定义 DNS 解析器和 SSH 隧道通过现有的[`客户端配置`](#client-configuration)来支持。


### 使用 ICMP 监控端点
通过在 `endpoints[].url` 前添加 `icmp://` 前缀，您可以使用 ICMP（通常被称为"ping"或"echo"）在非常基础的层面上监控端点：
```yaml
endpoints:
  - name: ping-example
    url: "icmp://example.com"
    conditions:
      - "[CONNECTED] == true"
```

ICMP 类型的端点仅支持 `[CONNECTED]`、`[IP]` 和 `[RESPONSE_TIME]` 占位符。
您可以指定以 `icmp://` 为前缀的域名，或以 `icmp://` 为前缀的 IP 地址。

如果您在 Linux 上运行 Gatus，遇到任何问题请阅读 [https://github.com/prometheus-community/pro-bing#linux] 上的 Linux 部分。

在 `v5.31.0` 之前，某些环境设置需要添加 `CAP_NET_RAW` 功能才能使 ping 正常工作。
从 `v5.31.0` 开始，这不再必要，ICMP 检查将使用非特权 ping 工作，除非以 root 身份运行。详见 #1346。


### 使用 DNS 查询监控端点
在端点中定义 `dns` 配置将自动将该端点标记为 DNS 类型的端点：
```yaml
endpoints:
  - name: example-dns-query
    url: "8.8.8.8" # 要使用的 DNS 服务器地址
    dns:
      query-name: "example.com"
      query-type: "A"
    conditions:
      - "[BODY] == 93.184.215.14"
      - "[DNS_RCODE] == NOERROR"
```

DNS 类型端点的条件中可以使用两个占位符：
- 占位符 `[BODY]` 解析为查询的输出。例如，类型为 `A` 的查询将返回一个 IPv4 地址。
- 占位符 `[DNS_RCODE]` 解析为与查询返回的响应代码关联的名称，例如
`NOERROR`、`FORMERR`、`SERVFAIL`、`NXDOMAIN` 等。


### 使用 SSH 监控端点
您可以通过在 `endpoints[].url` 前添加 `ssh://` 前缀来使用 SSH 监控端点：
```yaml
endpoints:
  # 基于密码的 SSH 示例
  - name: ssh-example-password
    url: "ssh://example.com:22" # 端口是可选的。默认为 22。
    ssh:
      username: "username"
      password: "password"
    body: |
      {
        "command": "echo '{\"memory\": {\"used\": 512}}'"
      }
    interval: 1m
    conditions:
      - "[CONNECTED] == true"
      - "[STATUS] == 0"
      - "[BODY].memory.used > 500"

  # 基于密钥的 SSH 示例
  - name: ssh-example-key
    url: "ssh://example.com:22" # 端口是可选的。默认为 22。
    ssh:
      username: "username"
      private-key: |
        -----BEGIN RSA PRIVATE KEY-----
        TESTRSAKEY...
        -----END RSA PRIVATE KEY-----
    interval: 1m
    conditions:
      - "[CONNECTED] == true"
      - "[STATUS] == 0"
```

您也可以通过不指定用户名、密码和私钥字段来使用无认证方式监控端点。

```yaml
endpoints:
  - name: ssh-example
    url: "ssh://example.com:22" # 端口是可选的。默认为 22。
    ssh:
      username: ""
      password: ""
      private-key: ""

    interval: 1m
    conditions:
      - "[CONNECTED] == true"
      - "[STATUS] == 0"
```

SSH 类型端点支持以下占位符：
- `[CONNECTED]` 如果 SSH 连接成功则解析为 `true`，否则为 `false`
- `[STATUS]` 解析为在远程服务器上执行的命令的退出代码（例如 `0` 表示成功）
- `[BODY]` 解析为在远程服务器上执行的命令的标准输出
- `[IP]` 解析为服务器的 IP 地址
- `[RESPONSE_TIME]` 解析为建立连接和执行命令所花费的时间


### 使用 STARTTLS 监控端点
如果您有一个邮件服务器，希望确保没有问题，通过 STARTTLS 进行监控将作为一个很好的初步指标：
```yaml
endpoints:
  - name: starttls-smtp-example
    url: "starttls://smtp.gmail.com:587"
    interval: 30m
    client:
      timeout: 5s
    conditions:
      - "[CONNECTED] == true"
      - "[CERTIFICATE_EXPIRATION] > 48h"
```


### 使用 TLS 监控端点
监控使用 SSL/TLS 加密的端点，例如基于 TLS 的 LDAP，可以帮助检测证书过期：
```yaml
endpoints:
  - name: tls-ldaps-example
    url: "tls://ldap.example.com:636"
    interval: 30m
    client:
      timeout: 5s
    conditions:
      - "[CONNECTED] == true"
      - "[CERTIFICATE_EXPIRATION] > 48h"
```

如果设置了 `endpoints[].body`，则会发送该内容，响应的前 1024 字节将存储在 `[BODY]` 中。

占位符 `[STATUS]` 以及字段 `endpoints[].headers`、
`endpoints[].method` 和 `endpoints[].graphql` 不支持用于 TLS 端点。


### 监控域名过期
您可以使用 `[DOMAIN_EXPIRATION]` 占位符监控域名的过期时间，适用于除 DNS 之外的所有端点类型：
```yaml
endpoints:
  - name: check-domain-and-certificate-expiration
    url: "https://example.org"
    interval: 1h
    conditions:
      - "[DOMAIN_EXPIRATION] > 720h"
      - "[CERTIFICATE_EXPIRATION] > 240h"
```

> ⚠ 使用 `[DOMAIN_EXPIRATION]` 占位符需要 Gatus 使用 RDAP，或作为备选方案，
> [通过一个库](https://github.com/TwiN/whois)向官方 IANA WHOIS 服务发送请求，
> 在某些情况下还需要向特定 TLD 的 WHOIS 服务器（例如 `whois.nic.sh`）发送二次请求。
> 为了防止 WHOIS 服务因您发送过多请求而限制您的 IP 地址，Gatus 会阻止您在间隔小于 `5m` 的端点上
> 使用 `[DOMAIN_EXPIRATION]` 占位符。


### 并发
默认情况下，Gatus 允许最多 3 个端点/套件同时进行监控。这在性能和资源使用之间提供了平衡，同时保持准确的响应时间测量。

您可以使用 `concurrency` 参数配置并发级别：

```yaml
# 允许 10 个端点/套件同时监控
concurrency: 10

# 允许无限制的并发监控
concurrency: 0

# 使用默认并发数（3）
# concurrency: 3
```

**重要注意事项：**
- 更高的并发数可以在您有许多端点时提高监控性能
- 使用 `[RESPONSE_TIME]` 占位符的条件在非常高的并发下可能不太准确，因为系统资源竞争
- 设置为 `0` 表示无限并发（等同于已弃用的 `disable-monitoring-lock: true`）

**适用于更高并发的场景：**
- 您有大量需要监控的端点
- 您希望以非常短的间隔（< 5s）监控端点
- 您正在使用 Gatus 进行负载测试场景

**旧版配置：**
`disable-monitoring-lock` 参数已弃用，但仍支持向后兼容。它等同于设置 `concurrency: 0`。


### 运行时重新加载配置
为了方便起见，如果在 Gatus 运行期间更新了加载的配置文件，Gatus 会自动实时重新加载配置。

默认情况下，如果更新的配置无效，应用程序将退出，但您可以通过将 `skip-invalid-config-update` 设置为 `true`
来配置 Gatus 在配置文件更新为无效配置时继续运行。

请记住，每次在 Gatus 运行期间对配置文件进行更新后，确保配置文件的有效性符合您的最佳利益。
请查看日志并确保您没有看到以下消息：
```
The configuration file was updated, but it is not valid. The old configuration will continue being used.
```
如果不这样做，可能会导致 Gatus 在因任何原因重启时无法启动。

我建议不要将 `skip-invalid-config-update` 设置为 `true` 以避免这种情况，但选择权在您。

**如果您没有使用文件存储**，在 Gatus 运行期间更新配置实际上等同于重启应用程序。

> 📝 如果绑定的是配置文件而不是配置文件夹，则可能无法检测到更新。参见 [#151](https://github.com/TwiN/gatus/issues/151)。


### 端点分组
端点分组用于在仪表板上将多个端点组合在一起。

```yaml
endpoints:
  - name: frontend
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"

  - name: backend
    group: core
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"

  - name: monitoring
    group: internal
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"

  - name: nas
    group: internal
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"

  - name: random endpoint that is not part of a group
    url: "https://example.org/"
    interval: 5m
    conditions:
      - "[STATUS] == 200"
```

上述配置在按分组排序时，仪表板将呈现如下效果：

![Gatus Endpoint Groups](.github/assets/endpoint-groups.jpg)


### 如何默认按分组排序？
在配置文件中将 `ui.default-sort-by` 设置为 `group`：
```yaml
ui:
  default-sort-by: group
```
请注意，如果用户已经按其他字段排序了仪表板，除非用户清除浏览器的 localstorage，否则不会应用默认排序。


### 在自定义路径上暴露 Gatus
目前，您可以使用完全限定域名（FQDN）来暴露 Gatus UI，例如 `status.example.org`。但是，它不支持基于路径的路由，这意味着您无法通过类似 `example.org/status/` 的 URL 来暴露它。

更多信息请参见 https://github.com/TwiN/gatus/issues/88。


### 在自定义端口上暴露 Gatus
默认情况下，Gatus 暴露在 `8080` 端口上，但您可以通过设置 `web.port` 参数来指定不同的端口：
```yaml
web:
  port: 8081
```

如果您使用的是像 Heroku 这样的 PaaS，不允许您设置自定义端口而是通过环境变量暴露端口，
请参阅[在配置文件中使用环境变量](#use-environment-variables-in-config-files)。

### 在配置文件中使用环境变量

您可以直接在配置文件中使用环境变量，它们将从环境中替换：
```yaml
web:
  port: ${PORT}

ui:
  title: $TITLE
```
⚠️ 当您的配置参数包含 `$` 符号时，您需要使用 `$$` 来转义 `$`。

### 配置启动延迟
如果出于任何原因，您需要 Gatus 在应用程序启动时等待一段时间再开始监控端点，您可以使用 `GATUS_DELAY_START_SECONDS` 环境变量使 Gatus 在启动时休眠。


### 保持配置文件精简
虽然这不是 Gatus 特有的，但您可以利用 YAML 锚点来创建默认配置。
如果您有一个大型配置文件，这应该能帮助您保持整洁。

<details>
  <summary>示例</summary>

```yaml
default-endpoint: &defaults
  group: core
  interval: 5m
  client:
    insecure: true
    timeout: 30s
  conditions:
    - "[STATUS] == 200"

endpoints:
  - name: anchor-example-1
    <<: *defaults               # 这将把 &defaults 下的配置合并到此端点
    url: "https://example.org"

  - name: anchor-example-2
    <<: *defaults
    group: example              # 这将覆盖 &defaults 中定义的 group
    url: "https://example.com"

  - name: anchor-example-3
    <<: *defaults
    url: "https://twin.sh/health"
    conditions:                # 这将覆盖 &defaults 中定义的 conditions
      - "[STATUS] == 200"
      - "[BODY].status == UP"
```
</details>


### 代理客户端配置
您可以通过在客户端配置中设置 `proxy-url` 参数来为客户端配置代理。

```yaml
endpoints:
  - name: website
    url: "https://twin.sh/health"
    client:
      proxy-url: http://proxy.example.com:8080
    conditions:
      - "[STATUS] == 200"
```


### 如何修复 431 Request Header Fields Too Large 错误
根据您的环境部署位置以及 Gatus 前面的中间件或反向代理类型，
您可能会遇到此问题。这可能是因为请求头过大，例如大型 cookie。

默认情况下，`web.read-buffer-size` 设置为 `8192`，但像这样增加此值将增大读取缓冲区大小：
```yaml
web:
  read-buffer-size: 32768
```

### 徽章
#### 正常运行时间
![Uptime 1h](https://status.twin.sh/api/v1/endpoints/core_blog-external/uptimes/1h/badge.svg)
![Uptime 24h](https://status.twin.sh/api/v1/endpoints/core_blog-external/uptimes/24h/badge.svg)
![Uptime 7d](https://status.twin.sh/api/v1/endpoints/core_blog-external/uptimes/7d/badge.svg)
![Uptime 30d](https://status.twin.sh/api/v1/endpoints/core_blog-external/uptimes/30d/badge.svg)

Gatus 可以自动为您监控的某个端点生成 SVG 徽章。
这允许您在各个应用程序的 README 中放置徽章，甚至可以根据需要创建自己的状态页面。

生成徽章的路径如下：
```
/api/v1/endpoints/{key}/uptimes/{duration}/badge.svg
```
其中：
- `{duration}` 为 `30d`、`7d`、`24h` 或 `1h`
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

例如，如果您想获取 `core` 分组中 `frontend` 端点最近 24 小时的正常运行时间，
URL 将如下所示：
```
https://example.com/api/v1/endpoints/core_frontend/uptimes/7d/badge.svg
```
如果您想显示不属于任何分组的端点，必须将分组值留空：
```
https://example.com/api/v1/endpoints/_frontend/uptimes/7d/badge.svg
```
示例：
```
![Uptime 24h](https://status.twin.sh/api/v1/endpoints/core_blog-external/uptimes/24h/badge.svg)
```
如果您想查看每个可用徽章的可视化示例，可以直接导航到端点的详情页面。


#### 健康状态
![Health](https://status.twin.sh/api/v1/endpoints/core_blog-external/health/badge.svg)

生成徽章的路径如下：
```
/api/v1/endpoints/{key}/health/badge.svg
```
其中：
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

例如，如果您想获取 `core` 分组中 `frontend` 端点的当前状态，
URL 将如下所示：
```
https://example.com/api/v1/endpoints/core_frontend/health/badge.svg
```


#### 健康状态（Shields.io）
![Health](https://img.shields.io/endpoint?url=https%3A%2F%2Fstatus.twin.sh%2Fapi%2Fv1%2Fendpoints%2Fcore_blog-external%2Fhealth%2Fbadge.shields)

生成徽章的路径如下：
```
/api/v1/endpoints/{key}/health/badge.shields
```
其中：
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

例如，如果您想获取 `core` 分组中 `frontend` 端点的当前状态，
URL 将如下所示：
```
https://example.com/api/v1/endpoints/core_frontend/health/badge.shields
```

有关 Shields.io 徽章端点的更多信息请参见[这里](https://shields.io/badges/endpoint-badge)。


#### 响应时间
![Response time 1h](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/1h/badge.svg)
![Response time 24h](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/24h/badge.svg)
![Response time 7d](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/7d/badge.svg)
![Response time 30d](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/30d/badge.svg)

生成徽章的端点如下：
```
/api/v1/endpoints/{key}/response-times/{duration}/badge.svg
```
其中：
- `{duration}` 为 `30d`、`7d`、`24h` 或 `1h`
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

#### 响应时间（图表）
![Response time 24h](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/24h/chart.svg)
![Response time 7d](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/7d/chart.svg)
![Response time 30d](https://status.twin.sh/api/v1/endpoints/core_blog-external/response-times/30d/chart.svg)

生成响应时间图表的端点如下：
```
/api/v1/endpoints/{key}/response-times/{duration}/chart.svg
```
其中：
- `{duration}` 为 `30d`、`7d` 或 `24h`
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

##### 如何更改响应时间徽章的颜色阈值
要更改响应时间徽章的阈值，可以在端点中添加相应的配置。
数组中的值对应级别 [极佳, 优秀, 良好, 及格, 差]，
所有五个值必须以毫秒（ms）为单位给出。

```yaml
endpoints:
- name: nas
  group: internal
  url: "https://example.org/"
  interval: 5m
  conditions:
    - "[STATUS] == 200"
  ui:
    badge:
      response-time:
        thresholds: [550, 850, 1350, 1650, 1750]
```


### API
Gatus 提供了一个简单的只读 API，可以通过查询来以编程方式确定端点状态和历史记录。

所有端点均可通过 GET 请求访问以下端点获取：
```
/api/v1/endpoints/statuses
````
示例：https://status.twin.sh/api/v1/endpoints/statuses

也可以使用以下模式查询特定端点：
```
/api/v1/endpoints/{group}_{endpoint}/statuses
```
示例：https://status.twin.sh/api/v1/endpoints/core_blog-home/statuses

如果 `Accept-Encoding` HTTP 请求头包含 `gzip`，将使用 Gzip 压缩。

API 将返回 JSON 负载，`Content-Type` 响应头设置为 `application/json`。
查询 API 不需要此请求头。


#### 以编程方式与 API 交互
参见 [TwiN/gatus-sdk](https://github.com/TwiN/gatus-sdk)


#### 原始数据
Gatus 暴露您监控端点之一的原始数据。
这允许您在自己的应用程序中跟踪和聚合受监控端点的数据。例如，如果您想跟踪超过 7 天的正常运行时间。

##### 正常运行时间
获取端点原始正常运行时间数据的路径为：
```
/api/v1/endpoints/{key}/uptimes/{duration}
```
其中：
- `{duration}` 为 `30d`、`7d`、`24h` 或 `1h`
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

例如，如果您想获取 `core` 分组中 `frontend` 端点最近 24 小时的原始正常运行时间数据，URL 将如下所示：
```
https://example.com/api/v1/endpoints/core_frontend/uptimes/24h
```

##### 响应时间
获取端点原始响应时间数据的路径为：
```
/api/v1/endpoints/{key}/response-times/{duration}
```
其中：
- `{duration}` 为 `30d`、`7d`、`24h` 或 `1h`
- `{key}` 的格式为 `<GROUP_NAME>_<ENDPOINT_NAME>`，其中两个变量中的 ` `、`/`、`_`、`,`、`.`、`#`、`+` 和 `&` 都被替换为 `-`。

例如，如果您想获取 `core` 分组中 `frontend` 端点最近 24 小时的原始响应时间数据，URL 将如下所示：
```
https://example.com/api/v1/endpoints/core_frontend/response-times/24h
```


### 以二进制方式安装
您可以使用以下命令将 Gatus 作为二进制文件下载：
```
go install github.com/TwiN/gatus/v5@latest
```


### 高层设计概览
![Gatus diagram](.github/assets/gatus-diagram.jpg)
