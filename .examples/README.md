# Upstream example compose layouts

These directories are largely inherited from upstream [TwiN/gatus](https://github.com/TwiN/gatus) sample deployments.

**Important for this fork (`isleei/gatus`):** many upstream samples pull `twinproduction/gatus` / `ghcr.io/twin/gatus`. That image does **not** include this fork’s Admin, WeCom, managed overlay, or other product changes. For local/production use of **this** fork, build your own image from the repo root:

```bash
docker build -t gatus:local .
```

Then point compose `image:` / `build:` at `gatus:local` (see `docs/examples/sentinel/compose.yaml` and `docs/DEPLOYMENT.md`). We intentionally do **not** rewrite every upstream `.examples/*/docker-compose*.yml` file here — treat them as layout references, not drop-in production stacks for the fork.
