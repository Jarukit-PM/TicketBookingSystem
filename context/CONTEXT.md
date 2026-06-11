# Cinema Ticket Booking

Domain language for the cinema ticket booking platform. Implementation details live in `architecture-context.md`.

## Catalog

**Movie**:
A global catalog entry (title, poster, synopsis, etc.) shared across all cinemas. Not scoped to a single venue.
_Avoid_: Per-cinema movie, venue-specific film record (MVP)

**Showtime**:
A scheduled screening of a movie on a specific screen at a specific time. Cinema-scoped via screen → cinema. Drives seat inventory and booking.
_Avoid_: Screening, performance (unless quoting external systems)

**Now showing (browse)**:
For a selected cinema, a movie appears in the **Now Showing** tab only if it has at least one future showtime at that cinema and is not `ARCHIVED`.
_Avoid_: Global status-only browse, show bookable film with zero local showtimes

**Coming soon (browse)**:
For a selected cinema, the **Coming Soon** tab lists movies with `status: COMING_SOON` even when there are **no showtimes** yet (teaser/detail only — no booking). Once showtimes exist, the film may also appear under Now Showing.
_Avoid_: Requiring showtimes for Coming Soon teasers

**Blocked seat**:
A seat permanently unavailable on a screen, defined in the screen layout (`type: blocked`). Applies to every showtime on that screen until an admin edits the layout. No per-showtime block overrides in MVP.
_Avoid_: Per-showtime block, maintenance hold (MVP)

**Price tier**:
A showtime-level price for a seat layout type (e.g. `standard`, `vip`, `wheelchair`). Order total is the sum of each selected seat's tier price from the layout `type` → `priceTiers` map.
_Avoid_: Flat showtime price (when types differ), per-seat price override (MVP)

## Booking lifecycle

**Seat hold**:
A short-lived reservation of one or more seats on a showtime, stored in Redis only. Not visible in My Bookings or admin booking search.
_Avoid_: Pending booking, temporary booking, cart

**Booking**:
A confirmed purchase persisted in MongoDB after successful checkout. Created only via confirm; always tied to a user, showtime, seats, and ticket.
_Avoid_: Hold, reservation (when meaning Redis hold)

**Hold expiry**:
When a seat hold's TTL ends without confirm. Seats return to available; nothing is written to the `bookings` collection.
_Avoid_: Expired booking

## Booking rules

**Multiple bookings per showtime**:
A user may confirm more than one booking for the same showtime. Each confirm creates a separate booking with its own `bookingRef`, ticket, and confirmation email.
_Avoid_: One booking per showtime, add-seats-to-booking (MVP)

**Seat deselect**:
Removing a seat from the user's checkout selection releases that seat immediately for others. Remaining held seats keep their current expiry — TTL does not refresh on remove.
_Avoid_: Refresh timer on deselect, block individual deselect (MVP)

**Navigate away**:
Leaving the seat map without confirming does not release holds early. Seats remain held until TTL expires, explicit abandon (`DELETE` holds), or successful confirm. WebSocket disconnect alone does not release holds.
_Avoid_: Release on disconnect, ghost holds until manual cleanup (MVP)

**Concurrent holds**:
A user may hold seats on multiple showtimes at the same time. Each showtime has its own hold set, TTL, and countdown. Confirming one showtime does not release holds on another.
_Avoid_: One checkout at a time, auto-release other showtimes on new hold (MVP)

**Confirm**:
Completes checkout for the user's entire hold set on one showtime. Books all currently held seats in one booking; to book fewer, deselect before confirm. No partial-confirm payload in MVP.
_Avoid_: Subset confirm, seatIds on confirm body (MVP)

**Idempotency key**:
Client-supplied key on confirm. If a prior request with the same key already succeeded, the server returns that booking again. If no successful booking exists and holds are gone, the server rejects with conflict — client must hold seats again with a new key.
_Avoid_: Retry confirm without holds, idempotent re-book after expiry

## Access

**Admin**:
A user with elevated role who can manage all cinemas, catalogs, showtimes, bookings, and audit logs. Global access in MVP — no per-cinema admin assignment.
_Avoid_: Venue manager, cinema-scoped admin (MVP)

**Session**:
Authenticated API access for the SPA via JWT in an `httpOnly` cookie (same origin behind nginx). No Bearer token in client JS for MVP.
_Avoid_: localStorage JWT, Authorization header in SPA (MVP)

**Account link**:
When Google OAuth returns a verified email that already belongs to a user, attach `googleId` to that user and sign in. One email = one account and one booking history.
_Avoid_: Duplicate accounts per email, block Google if password exists (MVP)

**Showtime cutoff**:
Customers cannot hold or confirm seats once the showtime has started (`startsAt <= now` in the cinema's timezone). Past screenings are not bookable in MVP.
_Avoid_: Late booking window, book after start (MVP)

**Seat limit**:
Maximum **10 seats** per hold and per booking for one user on one showtime. Server rejects hold/confirm above the cap.
_Avoid_: Unlimited seats per checkout (MVP)

**Booking reference**:
Human-readable ID for a confirmed booking (email, My Bookings, support). Format: `TBS-` + 6–8 alphanumeric characters (exclude ambiguous `0/O`, `1/I`). Distinct from the opaque `ticketToken` used in QR validation.
_Avoid_: UUID as customer-facing ref, reusing ticketToken as bookingRef

**Guest browse**:
Anyone may view movies, showtimes, and the live seat map without signing in. Authentication is required to hold seats or confirm a booking.
_Avoid_: Login before seat map, guest checkout (MVP)
