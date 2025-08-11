# afterhee-server

## Requirements

- Golang 1.24.0
- DuckDB cli

## Preparing database

```sh
cd database
duckdb db.duckdb < build.sql
```

## Instruction

```sh
go run server.go
```
