# Feature 03 — Go API Scaffold

Read `AGENTS.md`, `context/architecture-context.md`, `context/code-standards.md`, and `golang-project-layout` skill before starting.

Scaffold the **Go + Gin** API and **asynq worker** entrypoints with **Viper** config, MongoDB/Redis connection stubs, structured logging, and `GET /api/health` (or `/healthz`). Handlers return JSON; no domain features yet.

**Depends on:** Feature 02 (compose network) recommended. Feature 04 can follow in parallel for models only.

## Objective

`api/` is a buildable Go module with standard layout, boots in Docker, and exposes a health route through nginx. Worker process starts and connects to Redis without crashing.

**Success looks like:** `go test ./...` passes; `docker compose up api worker` shows healthy logs; `curl http://localhost/api/health` returns `{"status":"ok"}`.

## Assumptions

1. Go **1.22+**; module path matches repo (e.g. `github.com/Jarukit-PM/TicketBookingSystem/api`).
2. Gin for HTTP; `gorilla/websocket` or `nhooyr.io/websocket` added later in spec 07.
3. Viper reads `config.yaml` + env overrides (`MONGO_URI`, `REDIS_URL`, `JWT_SECRET`, `PORT`).
4. API listens on `0.0.0.0:8080`.

## Commands

```bash
cd api
go mod init github.com/Jarukit-PM/TicketBookingSystem/api   # if new
go get github.com/gin-gonic/gin github.com/spf13/viper \
  go.mongodb.org/mongo-driver/mongo github.com/redis/go-redis/v9 \
  github.com/hibiken/asynq
go run ./cmd/server
go run ./cmd/worker
go test ./...
go vet ./...
```

## Project Structure

```
api/
├── cmd/
│   ├── server/main.go      # Gin router, graceful shutdown
│   └── worker/main.go      # asynq server stub (no tasks yet)
├── internal/
│   ├── config/config.go    # Viper load + typed Config struct
│   ├── db/mongo.go         # Connect, Ping, Disconnect
│   ├── db/redis.go         # Connect, Ping
│   └── middleware/         # request ID, recovery (auth in spec 05)
├── pkg/
│   └── httputil/response.go # JSON helpers, error shape
├── config.yaml.example
├── Dockerfile              # multi-stage; target server + worker
└── go.mod
```

## Code Style

Thin `main` — wire dependencies, pass to router:

```go
func main() {
    cfg := config.MustLoad()
    ctx := context.Background()
    mongoClient := db.MustConnectMongo(ctx, cfg.MongoURI)
    redisClient := db.MustConnectRedis(cfg.RedisURL)
    defer mongoClient.Disconnect(ctx)

    r := gin.New()
    r.Use(gin.Recovery(), middleware.RequestID())
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
    // graceful shutdown on SIGTERM
}
```

JSON errors (consistent for later specs):

```json
{ "error": { "code": "INTERNAL", "message": "..." } }
```

## Testing Strategy

- `internal/config`: test env override precedence with `t.Setenv`.
- `GET /api/health`: `httptest` returns 200.
- Mongo/Redis: integration tests optional behind `-tags=integration` (skip in CI until compose service).

## Boundaries

- **Always:** `context.Context` on DB calls; bind `0.0.0.0`; log port on start.
- **Ask first:** New top-level dependencies, changing Dockerfile base image.
- **Never:** Business logic in `cmd/`; hardcode secrets.

## Implementation notes

### Config (`internal/config`)

| Field | Env | Default |
| ----- | --- | ------- |
| `Port` | `PORT` | `8080` |
| `MongoURI` | `MONGO_URI` | `mongodb://mongo:27017/tbs` |
| `RedisURL` | `REDIS_URL` | `redis://redis:6379/0` |
| `JWTSecret` | `JWT_SECRET` | dev-only placeholder |
| `AppURL` | `APP_URL` | `http://localhost` |

### Worker stub

- Connect Redis; register empty `asynq.ServeMux`; log "worker ready".
- Task handlers added in spec 09.

### Dockerfile

- Stage 1: build `server` and `worker` binaries.
- Stage 2: distroless or alpine runtime; `CMD` overridden in compose for worker.

## Tasks

- [ ] Initialize `api/go.mod` and directory layout
- [ ] Implement config, mongo/redis connect with ping on health (optional deep check)
- [ ] `cmd/server` with Gin + `/api/health`
- [ ] `cmd/worker` asynq stub
- [ ] `api/Dockerfile`; update `docker-compose.yml` commands
- [ ] Extend health to verify mongo + redis when `?deep=1` (optional)

## Out of scope

- Auth middleware, domain routes, WebSocket
- Index migrations (spec 04)
- asynq email tasks (spec 09)

## Check when done

- [ ] `go test ./...` and `go vet ./...` pass in CI
- [ ] `curl http://localhost/api/health` returns 200 via nginx
- [ ] Worker container stays running
- [ ] `config.yaml.example` committed; secrets in `.env` only
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None.
