#!/usr/bin/env bash
# Satellite agent: probe a target locally, then push success/failure to primary.
# Requires: curl. Optional: gatus-cli if you prefer its helpers.
set -euo pipefail

PRIMARY_URL="${GATUS_PRIMARY_URL:?set GATUS_PRIMARY_URL e.g. https://status.example.com}"
EXTERNAL_KEY="${GATUS_EXTERNAL_KEY:?set key e.g. region-eu_website}"
TOKEN="${GATUS_EXTERNAL_TOKEN:?set Bearer token matching primary external-endpoints[].token}"
TARGET_URL="${GATUS_TARGET_URL:?URL to probe from this region}"

start_ns=$(date +%s%N)
http_code=0
if curl -fsS -o /dev/null -w "%{http_code}" --max-time 15 "$TARGET_URL" >/tmp/gatus-sat-code 2>/tmp/gatus-sat-err; then
  http_code=$(cat /tmp/gatus-sat-code)
else
  http_code=0
fi
end_ns=$(date +%s%N)
duration_ms=$(( (end_ns - start_ns) / 1000000 ))
duration="${duration_ms}ms"

success=false
error=""
if [[ "$http_code" == "200" ]]; then
  success=true
else
  error="region probe failed status=${http_code}"
fi

query="success=${success}&duration=${duration}"
if [[ -n "$error" ]]; then
  query="${query}&error=$(python3 -c 'import urllib.parse,sys; print(urllib.parse.quote(sys.argv[1]))' "$error")"
fi

curl -fsS -X POST \
  -H "Authorization: Bearer ${TOKEN}" \
  "${PRIMARY_URL}/api/v1/endpoints/${EXTERNAL_KEY}/external?${query}"

echo "pushed key=${EXTERNAL_KEY} success=${success} duration=${duration}"
