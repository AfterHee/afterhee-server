#!/bin/sh

# Guard
set -euo pipefail

URL="${URL:-http://localhost/api/v1/healthcheck}"   # 필요시 환경변수 URL로 오버라이드
WEBHOOK_URL="${WEBHOOK_URL:-}"               # 디스코드 웹훅 URL (필수)
TIMEOUT="${TIMEOUT:-3}"                      # 타임아웃 초
HOSTNAME="$(hostname)"
TIMESTAMP="$(date +'%Y-%m-%d %H:%M:%S %z')"

if [[ -z "${WEBHOOK_URL}" ]]; then
  echo "[ERROR] WEBHOOK_URL not set" >&2
  exit 2
fi

# --- 헬스체크 ---
echo "connection failed" > "./hc_body.$$"

HTTP_CODE="$(curl -sS -o ./hc_body.$$ -w '%{http_code}' --connect-timeout "${TIMEOUT}" -m "${TIMEOUT}" "${URL}" || echo 000)"

if [[ "${HTTP_CODE}" != "200" ]]; then
  BODY="$(head -c 500 ./hc_body.$$ | sed 's/"/\\"/g')"  # 메시지에 넣기 위해 따옴표 이스케이프
  CONTENT=":rotating_light: Healthcheck FAILED
- host: ${HOSTNAME}
- url: ${URL}
- http_code: ${HTTP_CODE}
- time: ${TIMESTAMP}
- body: ${BODY}"

  # 디스코드 전송 (간단히 content만)
  curl -sS -H 'Content-Type: application/json' \
       -d "$(jq -n --arg content "$CONTENT" '{content: $content}')" \
       "${WEBHOOK_URL}" >/dev/null || true
fi

rm -f ./hc_body.$$
