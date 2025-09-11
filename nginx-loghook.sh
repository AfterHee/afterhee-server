#!/usr/bin/env bash
set -euo pipefail

: "${WEBHOOK_URL:?WEBHOOK_URL not set}"   # 환경변수로 주입
CONTAINER="${CONTAINER:-afterhee-nginx}"           # docker 컨테이너 이름/ID
PROVIDER="${PROVIDER:-docker}"            # docker|stdin
RATE_LIMIT_PER_SEC="${RATE_LIMIT_PER_SEC:-5}"  # 초당 전송 상한 (단순 레이트 리밋)

# 라인 버퍼링
export STDOUT_LINE_BUFFERED=1

send_webhook() {
  local text="$1"
  # 필요한 경우 이곳에서 escape
  text="${text//$'\r'/ }"
  text="${text//\"/\\\"}"

  # Discord 예시 (Slack이면 {"text": "..."} 로 바꾸면 됨)
  curl -sS -m 5 -H 'Content-Type: application/json' \
    -d "{\"content\":\"${text}\"}" \
    "$WEBHOOK_URL" >/dev/null || true
}

# 초당 단순 레이트리밋
can_send() {
  local now sec countfile
  now="$(date +%s)"
  countfile="/tmp/nginx-loghook.$now.count"
  if [ ! -f "$countfile" ]; then
    printf "0" >"$countfile"
  fi
  local cnt
  cnt=$(cat "$countfile" 2>/dev/null || echo 0)
  if [ "$cnt" -lt "$RATE_LIMIT_PER_SEC" ]; then
    cnt=$((cnt+1))
    printf "%s" "$cnt" >"$countfile"
    find /tmp -maxdepth 1 -type f -name 'nginx-loghook.*.count' -mmin +2 -delete 2>/dev/null || true
    return 0
  fi
  return 1
}

# 입력 소스 열기
open_stream() {
  if [ "$PROVIDER" = "docker" ]; then
    # 과거 로그 무시, 지금부터 팔로우
    docker logs --tail 0 -f "$CONTAINER"
    # 또는 시간기준:
    # docker logs --since=0s -f "$CONTAINER"
  else
    tail -n0 -F "$LOG_FILE"
  fi
}

open_stream | stdbuf -oL grep -E --line-buffered '" (4|5)[0-9]{2} ' | \
while IFS= read -r line; do
  # 예시 라인:
  # 140.248.29.2 - - [11/Sep/2025:02:00:35 +0000] "GET /api/v1/... HTTP/1.1" 500 69 "-" "UA..."
  # 또는 docker 로그 프리픽스가 있는 경우도 그대로 매칭됨

  # 상태코드/메서드/경로 파싱 (공백 포함 요청라인 안전추출)
  status=""
  method=""
  uri=""

  # 상태코드: 닫는 따옴표 뒤의 세 자리
  if [[ "$line" =~ \"\ ([0-9]{3})\  ]]; then
    status="${BASH_REMATCH[1]}"
  fi

  # === 여기서 404는 무시 ===
  if [[ "$status" == "404" ]]; then
    continue
  fi

  # 요청라인: "METHOD URI HTTP/x"
  if [[ "$line" =~ \"([A-Z]+)\ ([^[:space:]]+)\ HTTP/[0-9.]+\" ]]; then
    method="${BASH_REMATCH[1]}"
    uri="${BASH_REMATCH[2]}"
  fi

  # UA (선택)
  ua=""
  if [[ "$line" =~ \"([^\"]+)\"$ ]]; then
    ua="${BASH_REMATCH[1]}"
  fi

  # 호스트/IP (선택)
  ip=""
  if [[ "$line" =~ ^([0-9a-fA-F\:\.]+)\  ]]; then
    ip="${BASH_REMATCH[1]}"
  fi

  # 최종 메시지
  msg="$(printf '%s' "$line")"
  # 간결 버전 원하면 아래 사용
  # msg="$(printf '%s [%s] %s %s ip=%s ua=%s' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$status" "$method" "$uri" "$ip" "$ua")"

  # 레이트리밋 체크 후 전송
  if can_send; then
    send_webhook "$msg"
  fi
done
