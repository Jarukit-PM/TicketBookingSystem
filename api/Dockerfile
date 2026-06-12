FROM golang:1.25-alpine AS base

WORKDIR /src

# Keep peak RAM low on small hosts (e.g. t3.small EC2). One compile at a time.
ENV GOMAXPROCS=1
ENV GOGC=50
ENV GOMEMLIMIT=1400MiB

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

FROM base AS build-server
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -p 1 -trimpath -ldflags="-s -w" -o /server ./cmd/server

FROM base AS build-worker
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -p 1 -trimpath -ldflags="-s -w" -o /worker ./cmd/worker

# Default image: both binaries (api + worker share this via compose).
FROM alpine:3.20 AS release

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build-server /server /server
COPY --from=build-worker /worker /worker

EXPOSE 8080

CMD ["/server"]

# Optional targets for building one binary at a time on very small instances:
#   docker build --target api -t tbs-api .
#   docker build --target worker -t tbs-worker .
FROM alpine:3.20 AS api

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build-server /server /server

EXPOSE 8080

CMD ["/server"]

FROM alpine:3.20 AS worker

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build-worker /worker /worker

CMD ["/worker"]
