# Feature 05 — Authentication

Read `context/CONTEXT.md`, `context/architecture-context.md`, `security-and-hardening` skill, and `clerk-vue-patterns` only if useful for Vue guards — this project uses custom JWT.

Email/password registration and login, **Google OAuth 2.0**, **httpOnly JWT cookie** session, role-based access (`customer` | `admin`). Vue Router guards and Pinia auth store.

**Depends on:** Features 01 (UI primitives), 03 (API), 04 (`users` collection).

## Objective

Customers and admins sign in securely. SPA uses cookie-based sessions (`credentials: 'include'`). Protected API routes reject unauthenticated callers; admin routes require `admin` role. Google OAuth auto-links by verified email (`CONTEXT.md`).

**Success looks like:** Register → login → cookie set → `GET /api/bookings/mine` works; non-admin blocked from `/api/admin/*`; Vue redirects unauthenticated users from `/book/*`, `/my-bookings`, `/admin/*`.

## Assumptions

1. **httpOnly JWT cookie only** for MVP — no Bearer in client JS (`CONTEXT.md`).
2. Cookie: `Secure` in production, `SameSite=Lax`, `Path=/api` or `/` per implementation (must work for `/api` routes).
3. First admin via `ADMIN_EMAIL` env on bootstrap or seed.
4. bcrypt cost 12 for password hashes.
5. Google OAuth: authorization code flow; callback `GET /api/auth/google/callback`.

## Commands

```bash
cd api && go test ./internal/auth/...
cd app && npm run type-check && npm run build
# Manual: register, login, inspect Set-Cookie in DevTools
```

## API Surface

| Method | Path | Auth | Body / notes |
| ------ | ---- | ---- | ------------ |
| POST | `/api/auth/register` | Public | `{ email, password, name }` |
| POST | `/api/auth/login` | Public | `{ email, password }` → Set-Cookie |
| POST | `/api/auth/logout` | Customer | Clear cookie |
| GET | `/api/auth/me` | Customer | `{ id, email, name, role }` |
| GET | `/api/auth/google` | Public | Redirect to Google |
| GET | `/api/auth/google/callback` | Public | Set-Cookie, redirect to `APP_URL` |

JWT claims: `sub` (user id), `role`, `exp`.

## Vue Structure

```
app/src/
├── stores/auth.ts           # user, login, logout, fetchMe
├── api/client.ts            # fetch wrapper, credentials: 'include'
├── views/auth/
│   ├── LoginView.vue
│   └── RegisterView.vue
└── router/index.ts          # beforeEach guards
```

### Router guards

| Route pattern | Rule |
| ------------- | ---- |
| `/book/*`, `/my-bookings` | Require auth |
| `/admin/*` | Require auth + `role === 'admin'` |
| Catalog, seat map view | Public (hold triggers login — spec 07) |

### Login prompt for holds

When guest clicks a seat (spec 07), redirect to `/login?redirect=/book/...` — document contract here.

## Code Style

Pinia store — no token in localStorage:

```ts
async function login(email: string, password: string) {
  const res = await api.post('/auth/login', { email, password })
  user.value = res.user // if returned; else await fetchMe()
}
```

Go middleware:

```go
func RequireAuth() gin.HandlerFunc { /* parse cookie JWT, set user in context */ }
func RequireAdmin() gin.HandlerFunc { /* role == admin */ }
```

## Testing Strategy

- Go: table tests for password hash, JWT issue/parse, `RequireAdmin` rejects customer.
- Go: OAuth auto-link when email matches existing user.
- Vue: unit test auth store with mocked fetch.
- Manual: Google OAuth with dev credentials.

## Boundaries

- **Always:** Validate email format; rate-limit login (Redis) — basic counter OK.
- **Ask first:** Changing cookie name, JWT expiry duration, adding refresh tokens.
- **Never:** Store JWT in localStorage; log passwords; commit OAuth client secrets.

## Tasks

- [ ] `internal/auth`: register, login, bcrypt, JWT cookie issue/clear
- [ ] Google OAuth handlers + auto-link by email
- [ ] Gin middleware `RequireAuth`, `RequireAdmin`
- [ ] Bootstrap admin from `ADMIN_EMAIL`
- [ ] Pinia `auth` store + API client with credentials
- [ ] Login/Register views (design system components)
- [ ] Router guards for protected routes

## Out of scope

- Password reset email
- MFA
- Bearer token API access
- Per-cinema admin RBAC

## Check when done

- [ ] Register + login sets httpOnly cookie; `/api/auth/me` returns user
- [ ] Google OAuth creates or links user; one email = one account
- [ ] Customer gets 403 on `/api/admin/*`
- [ ] Vue guards block `/admin/*` for customers
- [ ] `go test ./internal/auth/...` passes
- [ ] `progress-tracker.md` updated when implementation lands

## Open Questions

- JWT expiry duration (recommend **7d** session) — confirm at implementation if not set in env `JWT_EXPIRY`.
