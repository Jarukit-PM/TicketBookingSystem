# Cinema Ticket Booking System

## Overview

A full-stack cinema ticket booking platform where customers browse movies, pick showtimes, select seats on a live seat map, and complete bookings — with email confirmations, digital tickets, and an admin console for operations. The system centers on **real-time seat availability**: multiple users can browse the same showtime without double-booking, using short-lived seat holds during checkout and instant updates when seats are taken or released.

The customer app is a **Vue 3 + Vite** SPA (`app/`). The **Go (Gin)** API, **MongoDB**, **Redis**, and **WebSocket** real-time layer are defined in `architecture-context.md`.

## Goals

1. Let authenticated users discover movies and showtimes, then book seats through a clear, low-friction flow.
2. Keep seat inventory accurate under concurrency with real-time updates and temporary holds during checkout.
3. Deliver reliable booking confirmations via email, with scannable digital tickets.
4. Give admins a dashboard to manage movies, showtimes, halls, and bookings, with audit logs for support and compliance.
5. Let customers self-serve via booking history and on-demand access to their tickets.

## Core User Flow

### Customer (booking)

1. User signs up or signs in.
2. User picks a **cinema**, then browses **Now Showing** (movies with future showtimes) or **Coming Soon** (`COMING_SOON` teasers even without showtimes — see `CONTEXT.md`).
3. User selects a movie → picks a **showtime** (date, time, screen, format).
4. User opens the **seat map** for that showtime (no sign-in required to view); taken seats are disabled; available seats can be selected. Other users’ selections and bookings appear in near real time.
5. User **signs in** if needed, then selects one or more seats → system places a **temporary hold** with a visible **5-minute countdown** (timer resets when they add another seat) while they review the order.
6. User confirms booking → booking is **confirmed**, holds convert to sold seats, confirmation email is sent with ticket details and **QR code**.
7. User opens **My Bookings** to view upcoming and past tickets, re-open digital tickets, and download or print as needed.

### Admin

1. Admin signs in with elevated role.
2. Admin manages **movies**, **screens/halls**, **seat layouts**, and **showtimes**.
3. Admin monitors **today’s schedule**, occupancy, and recent bookings on the dashboard.
4. Admin searches bookings and reviews **audit logs** (who changed what, booking events, email delivery status). No cancel or refund actions in MVP.

## Features

### Authentication

- Email/password and Google OAuth sign-in.
- Roles: **Customer** and **Admin**.
- Protected routes: booking checkout, My Bookings, and all admin paths.
- Session management and secure API access for both SPA and admin.

### Seat Map (Real-time)

- Visual layout per screen: rows, seat numbers, aisles, wheelchair/companion seats, VIP zones.
- Live status per seat: `AVAILABLE`, `HELD` (by another user or self), `SOLD`, `BLOCKED` (admin maintenance).
- WebSocket or SSE (or equivalent) so all clients viewing the same showtime stay in sync.
- Optimistic UI with server authority to prevent double booking (transaction or row-level lock on confirm).

### Seat Hold TTL + Countdown

- Selected seats reserved for the current user with a **5-minute TTL**.
- **TTL refreshes** each time the user adds another seat on the same showtime (all held seats get a new 5-minute window). Deselecting a seat releases it immediately; remaining holds keep their current expiry.
- Visible **countdown timer** in the booking UI, driven by server `expiresAt`.
- Auto-release seats on hold timeout or explicit abandon. Navigating away without abandon leaves holds in place until TTL expires (no disconnect-based release in MVP).
- Held seats shown as unavailable to other users on the real-time seat map.

### Booking Flow

- Showtime selection → seat selection → order summary (seats, price breakdown) → confirm.
- Idempotent **confirm booking** so retries do not create duplicate orders.
- **Seat holds** live in Redis only (not bookings). Persisted **booking** status (MVP): `CONFIRMED` only. Hold timeout is a **hold expiry** — no row in `bookings`. No cancellation.

### My Bookings + Booking History

- **My Bookings** page: upcoming and past bookings for the signed-in user.
- Booking detail: movie, cinema, screen, showtime, seats, booking reference, status.
- Filter or tab by upcoming vs history; sort by showtime or booking date.
- Re-open confirmed tickets from history without contacting support.

### QR / Digital Ticket

- Each confirmed booking gets a unique **digital ticket** with a scannable **QR code** (booking reference + validation payload).
- Ticket view in-app: QR, show details, seat list; suitable for door check-in on mobile.
- QR included in confirmation email; optional PDF or print-friendly layout.
- Admin can **scan a ticket QR** to jump to that customer's **booking history** (support lookup — not door check-in).

### Admin Side

- **Dashboard**: booking counts, today’s showtimes, occupancy %, recent activity.
- **Catalog**: movies (title, poster, duration, rating, synopsis), showtimes, pricing tiers.
- **Venue**: screens, seat maps (layout template per screen); blocked seats set in layout (`type: blocked`) — applies to all showtimes on that screen.
- **Bookings**: search by user, email, booking ID, showtime (read-only in MVP).
- **Logs**: structured audit log (admin actions, booking lifecycle, integration errors) and email delivery log.

### Notification via Email

- **Booking confirmation** only: movie, cinema, screen, showtime, seats, total, booking reference, QR or link to digital ticket.
- Templates with plain-text + HTML; track send status for admin visibility.

## Scope

### In Scope (MVP)

- User authentication with customer vs admin roles.
- **Multi-cinema** catalog: multiple venues; browse and filter by cinema.
- Global movie catalog (admin-managed); showtimes scheduled per cinema/screen.
- Per-showtime seat map with real-time availability.
- 5-minute seat hold TTL (refreshed on each new seat) with visible countdown and auto-release.
- End-to-end booking: select seats → hold → confirm → persist booking.
- My Bookings and booking history for customers.
- QR / digital ticket generation and in-app display.
- Email confirmation on successful booking (including ticket / QR).
- Admin dashboard (summary metrics), booking management, and audit logs.
- Vue 3 frontend scaffold in `app/`; API and data model per `architecture-context.md` and feature specs.

### Out of Scope (MVP)

- **Payment integration** (e.g. Stripe): confirm-only bookings without online payment for MVP.
- Native mobile apps (responsive web only).
- Concessions / food ordering.
- Full loyalty program or membership tiers.
- Complex dynamic pricing (surge, demand-based) beyond fixed tiers per showtime.
- Public API for third-party aggregators.
- Offline box-office POS hardware integration.
- Promo codes, refunds, **booking cancellation**, SMS notifications, guest checkout.
- Door check-in validation (showtime window, pass/fail scan UI) — admin scan only opens booking history in MVP.

## Success Criteria

- Two users cannot book the same seat for the same showtime under concurrent load.
- Seat map updates within an acceptable latency (target: sub-second for connected clients).
- 5-minute hold countdown is accurate, refreshes when a new seat is added, and seats release when TTL expires.
- Confirmed bookings trigger a deliverable confirmation email with ticket access.
- Customers can view upcoming and past bookings and open a valid digital ticket with QR.
- Admins can create a showtime, see bookings, and trace actions in logs without database access.
- Customer can complete browse → seat select → confirm without leaving the SPA for the happy path.

## Related Context

- `architecture-context.md` — stack, domain model, real-time and booking invariants (to be aligned with this project).
- `feature-specs/` — per-feature specs (auth, seat map, booking, admin, email); legacy NovelCraft specs should be replaced or removed as this project progresses.
- `progress-tracker.md` — implementation status.
