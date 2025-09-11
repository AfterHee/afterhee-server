# afterhee-server

## Requirements

- Golang 1.24.0
- DuckDB cli

## Preparing database

```sh
pushd database
duckdb db.duckdb < build.sql
popd
```

## Instruction

```sh
export AFTERHEE_NEIS_API_KEY="YOUR_NEIS_API_KEY"
go run .
```

## Update Swagger Document

```sh
swag init # https://github.com/gofiber/swagger?tab=readme-ov-file#usage
```

## Envirionments

- `AFTERHEE_PORT`: 서버 포트. 지정하지 않으면 `8080`
- `AFTERHEE_DUCKDB_FILENAME`: DuckDB 파일 경로. 지정하지 않으면 `database/db.duckdb`
- `AFTERHEE_NEIS_API_KEY`: NEIS API 키. [나이스 교육정보 개방 포털](https://open.neis.go.kr/)에서 발급
- `AFTERHEE_GEMINI_API_KEY`: Gemini API 키.
- `AFTERHEE_PROFILE`: 서버 모드. `PRODUCTION`이면 운영모드. `DEV`이면 개발모드.

## Test

```sh
go test ./...
```

## Deploy

### Build Image

```sh
docker build -t afterhee:latest -t afterhee:v0.0.1 .

# 서버용 빌드
docker build --platform linux/amd64 -t afterhee:latest -t afterhee:v0.0.1 .

# podman
podman build -t afterhee:latest -t afterhee:v0.0.1 .
```

### Export Image

```sh
podman save --format docker-archive -o afterhee-v0.0.1.tar afterhee:v0.0.1 afterhee:latest
```

### Replace a certain server instance

```sh
source ./export_env.sh
docker compose up -d --no-deps --force-recreate "afterhee1" # -> Health Check
docker compose up -d --no-deps --force-recreate "afterhee2" # -> Health Check
docker compose up -d --no-deps --force-recreate "afterhee3"

# podman
podman compose up -d --no-deps --force-recreate "afterhee1" # -> Health Check
podman compose up -d --no-deps --force-recreate "afterhee2" # -> Health Check
podman compose up -d --no-deps --force-recreate "afterhee3"
```

### 서버에서 헬스체크 스크립트 1분마다 실행

```
crontab -e
* * * * * /path/to/health_check.sh
```

### 로그 웹훅 서비스 설치

```sh
sudo cp ./nginx-loghook.sh /usr/local/bin/nginx-loghook.sh
sudo install -m 0755 /usr/local/bin/nginx-loghook.sh /usr/local/bin/nginx-loghook.sh

# 환경변수 세팅
sudo cat >/etc/nginx-loghook.env <<'EOF'
WEBHOOK_URL="https://discord.com/api/webhooks/xxx/yyy"
CONTAINER="afterhee-nginx"
PROVIDER="docker"
RATE_LIMIT_PER_SEC="10"           # 필요 시 조정
EOF
sudo chmod 640 /etc/nginx-loghook.env

sudo cp ./nginx-loghook.service /etc/systemd/system/nginx-loghook.service

# 서비스 시작
sudo systemctl daemon-reload
sudo systemctl enable --now nginx-loghook
```
