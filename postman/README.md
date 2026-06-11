# Postman

Import these files into [Postman](https://www.postman.com/) (Import → drag files or browse).

| File | Purpose |
|------|---------|
| `TicketBookingSystem.postman_collection.json` | All REST API endpoints |
| `TBS-Local.postman_environment.json` | Local `docker compose` variables |

## Quick start

1. `docker compose up --build`
2. Import collection + environment
3. Select **TBS Local (docker compose)** environment
4. Run **Auth → Login** (cookie session is stored automatically)
5. Copy IDs from responses into environment variables (`cinemaId`, `showtimeId`, etc.)

Admin routes require an admin user (`ADMIN_EMAIL` / `ADMIN_SEED_PASSWORD` from `.env.example`).
