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
go run .
```

## Update Swagger Document

```sh
swag init # https://github.com/gofiber/swagger?tab=readme-ov-file#usage
```

## Envirionments

- `AFTERHEE_PORT`: 서버 포트. 지정하지 않으면 `8080`
- `AFTERHEE_DUCKDB_FILENAME`: DuckDB 파일 경로. 지정하지 않으면 `database/db.duckdb`
