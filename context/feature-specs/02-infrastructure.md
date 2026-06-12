# Feature 02 — Infrastructure (Docker, nginx, CI)

Read `AGENTS.md`, `context/architecture-context.md`, and `context/code-standards.md` before starting.

Local and CI foundation: **docker-compose** stack, **nginx** reverse proxy (single origin), and **GitHub Actions** CI. No business logic — services start and wire together.

**Depends on:** Feature 01 (design system) optional for CI `app` build. No Go API code required yet (api service may be a stub image or minimal health binary added in spec 03).

## Objective

Developers and CI can run the full stack with one command. nginx serves the Vue SPA, proxies `/api` and `/ws` to Gin, and all services share a compose network. CI validates Go and Vue on every PR.

**Users:** Engineers deploying and testing locally.

**Success looks like:** `docker compose up` brings up nginx + app build + api (stub OK) + worker (stub OK) + mongo + redis; `curl localhost/healthz` or `/api/health` returns OK; GitHub Actions passes on push.

## Assumptions

1. Single-machine local dev via Docker Compose (not Kubernetes).
2. Public port **80** on nginx; API internal on **8080**.
3. MongoDB and Redis use named volumes for persistence across restarts.
4. `api` and `worker` share one Docker image with different `command` (multi-stage build from `api/`).

## Commands

```bash
# Local stack
docker compose up --build

# Tear down (keep volumes)
docker compose down

# CI-equivalent (from repo root)
cd api && go vet ./... && go test ./...
cd app && npm ci && npm run lint && npm run type-check && npm run build
```

## Project Structure

```
TicketBookingSystem/
├── docker-compose.yml
├── nginx/
│   └── nginx.conf
├── .github/workflows/
│   └── ci.yml
├── app/                    # Vite build → dist (volume or copy into nginx)
└── api/                    # Go service (stub until spec 03)
```

## Code Style

nginx location blocks — API and WebSocket on same upstream:

```nginx
location /api/ {
    proxy_pass http://api:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
location /ws/ {
    proxy_pass http://api:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

## Testing Strategy

- **CI:** `go test ./...` (api may be empty module initially — `go test` still passes).
- **Manual:** `docker compose up` → browser `http://localhost` loads SPA; `curl -i http://localhost/api/health` returns 200 when api scaffold lands.
- No E2E in this spec.

## Boundaries

- **Always:** API binds `0.0.0.0:8080`; nginx is the only public entry on `:80`.
- **Ask first:** Adding new compose services, changing CI triggers, publishing images to a registry.
- **Never:** Commit `.env` with secrets; store real `JWT_SECRET` / `BREVO_API_KEY` in compose env files gitignored.

## docker-compose.yml

Services per `architecture-context.md`:

| Service | Image / build | Notes |
| ------- | ------------- | ----- |
| `nginx` | nginx:alpine | Mount `nginx.conf`; volume `app_dist` for SPA static |
| `app` | build `app/` | Build stage outputs `dist` to shared volume |
| `api` | build `api/` | `cmd/server`; depends on mongo, redis |
| `worker` | same as api | `cmd/worker`; depends on redis, mongo |
| `mongo` | mongo:7 | Volume `mongo_data` |
| `redis` | redis:7-alpine | Volume optional |

Env vars (compose `environment` or `.env.example`): `MONGO_URI`, `REDIS_URL`, `JWT_SECRET`, `APP_URL=http://localhost`.

## nginx.conf

- `/` → static files from `app` dist root; `try_files $uri $uri/ /index.html` for SPA routing.
- `/api/` → `api:8080`
- `/ws/` → WebSocket upgrade to `api:8080`

## GitHub Actions — `.github/workflows/ci.yml`

On `push` and `pull_request` to `main`:

1. **job: api** — setup-go, `go vet ./...`, `go test ./...` in `api/`
2. **job: app** — setup-node, `npm ci`, `npm run lint`, `npm run type-check`, `npm run build` in `app/`

Optional: cache npm and go modules.

## Tasks

- [ ] Add `docker-compose.yml` with all six services
  - Acceptance: `docker compose config` validates
  - Verify: `docker compose up --build` (api/worker may exit until spec 03 — document expected state)
  - Files: `docker-compose.yml`, `.env.example`

- [ ] Add `nginx/nginx.conf` and wire SPA volume
  - Acceptance: SPA loads at `http://localhost`
  - Verify: browser + `curl -I http://localhost`
  - Files: `nginx/nginx.conf`, compose volume mounts

- [ ] Add `.github/workflows/ci.yml`
  - Acceptance: workflow runs on PR
  - Verify: `act` locally or push branch
  - Files: `.github/workflows/ci.yml`

## Out of scope

- Production TLS, Render/Fly deploy manifests
- Playwright E2E in CI (phase 2)
- Brevo, Google OAuth secrets in compose (placeholders only)

## Check when done

- [ ] `docker compose up --build` starts mongo, redis, nginx; app static served
- [ ] nginx proxies `/api` and `/ws` paths (502 acceptable until api scaffold)
- [ ] `.github/workflows/ci.yml` runs Go + Vue jobs
- [ ] `.env.example` documents required env vars (no secrets)
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- None — resolved by `architecture-context.md` Docker Compose target.
