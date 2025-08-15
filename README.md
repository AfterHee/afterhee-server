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

## Test

```sh
go test ./...
```

## Deploy

### Build Image

```sh
docker build -t afterhee:latest -t afterhee:v0.0.1 .

# podman
podman build -t afterhee:latest -t afterhee:v0.0.1 .
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
