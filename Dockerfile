############################
# Build stage
############################
FROM golang:1.24.6-bookworm AS build

WORKDIR /build

# [소스 복사]
COPY . .

# [DB 초기화]
WORKDIR /build/database
RUN curl https://install.duckdb.org | sh
RUN ~/.duckdb/cli/latest/duckdb /build/database/db.duckdb < /build/database/build.sql

# [빌드]
WORKDIR /build

# 모듈 캐시 최적화
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 빌드 아규먼트(필요시 오버라이드)
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG BINARY_NAME=afterhee

# 빌드 캐시 활용
# duckdb 빌드를 위해 CGO 활성화 필요
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags "-s -w" -o /build/out/${BINARY_NAME} .

############################
# Runtime stage
############################
FROM debian:bookworm-slim

WORKDIR /app

RUN apt update \
 && apt install -y ca-certificates

COPY --from=build /build/out/afterhee /app/afterhee
COPY --from=build /build/database /app/database

EXPOSE 8080
ENTRYPOINT ["/app/afterhee"]
