# Self-review Stages 3-4

Branch: feat/stage3-4-product-and-optional
Repo: isleei/gatus
Date: 2026-09-10 (Asia/Shanghai)

## Fully implemented

- 3.2 Suite alerts: suites[].alerts; watchdog HandleSuiteAlerting; Admin suite alerts; tests in watchdog/suite_alerting_test.go
- 4.2 Audit retention: storage.admin-audit-max-age; daily cleanup; DELETE /api/v1/admin/audit-logs
- 3.1 Import YAML + WeCom/custom validation/help; CertificateMonitor 72h threshold

## Documented / examples

- 3.3 public vs admin + ui.default-sort-by: group
- 3.4 docs/examples/prometheus/
- 4.1 docs/examples/multi-region/
- 4.3 OIDC Chinese section in DEPLOYMENT.md

## Deferred

- Multi-step notification wizard
- Suite alert persistence reload on restart
- Prod secret rotation / deploy / IdP setup

## Remaining human ops

1. Deploy multi-region satellites
2. Set admin-audit-max-age in prod
3. Configure OIDC IdP
4. Restrict /metrics and Admin at proxy
5. Stage 0.2 credential rotation

## Fork features preserved

Admin v2, WeCom, certificates, groups, tamper, i18n.
