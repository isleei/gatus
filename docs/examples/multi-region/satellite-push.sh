#!/usr/bin/env bash
# Satellite agent: probe a target locally, then push success/failure to primary.
# Requires: curl. Optional: python3 for error URL-encoding (falls back to raw).
# Portable: no GNU date +%s%N (works on macOS / BusyBox).
set -euo pipefail

PRIMARY_URL="${GATUS_PRIMARY_URL:?set GATUS_PRIMARY_URL e.g. https://status.example.com}"
EXTERNAL_KEY="${GATUS_EXTERNAL_KEY:?set key e.g. region-eu_website}"
TOKEN="${GATUS_EXTERNAL_TOKEN:?set Bearer token matching primary external-endpoints[].token}"
TARGET_URL="${GATUS_TARGET_URL:?URL to probe from this region}"

# Portable millisecond duration (python3 preferred; awk fallback)
now_ms() {
  if command -v python3 >/dev/null 2>&1; then
    python3 -c 'import time; print(int(time.time() * 1000))'
  else
    # Seconds precision only (BusyBox / minimal shells without python3)
    echo $(( $(date +%s) * 1000 ))
  fi
}

tmp_dir="${TMPDIR:-/tmp}"
code_file=$(mktemp "${tmp_dir}/gatus-sat-code.XXXXXX")
err_file=$(mktemp "${tmp_dir}/gatus-sat-err.XXXXXX")
cleanup() { rm -f "$code_file" "$err_file"; }
trap cleanup EXIT

start_ms=$(now_ms)
http_code=0
if curl -fsS -o /dev/null -w "%{http_code}" --max-time 15 "$TARGET_URL" >"$code_file" 2>"$err_file"; then
  http_code=$(cat "$code_file")
else
  http_code=0
fi
end_ms=$(now_ms)
duration_ms=$(( end_ms - start_ms ))
if [[ "$duration_ms" -lt 0 ]]; then
  duration_ms=0
fi
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
  if command -v python3 >/dev/null 2>&1; then
    enc=$(python3 -c 'import urllib.parse,sys; print(urllib.parse.quote(sys.argv[1]))' "$error")
  else
    enc=$(printf '%s' "$error" | sed 's/ /%20/g')
  fi
  query="${query}&error=${enc}"
fi

curl -fsS -X POST \
  -H "Authorization: Bearer ${TOKEN}" \
  "${PRIMARY_URL}/api/v1/endpoints/${EXTERNAL_KEY}/external?${query}"

echo "pushed key=${EXTERNAL_KEY} success=${success} duration=${duration}"
