---
name: uber
description: Request and manage Uber rides — fare estimates, ride requests, status tracking, cancellation.
metadata: { "clawdbot": { "emoji": "🚗" } }
---

# Uber

Book, track, and manage Uber rides.

## Tools

- `uber(action, ...)` — universal Uber tool

## Actions

| Action     | What it does                        | Needs                      |
| ---------- | ----------------------------------- | -------------------------- |
| `products` | List available products at location | start_lat, start_lng       |
| `estimate` | Fare + time estimate                | start_lat/lng, end_lat/lng |
| `request`  | Book a ride                         | start_lat/lng, end_lat/lng |
| `status`   | Track current ride                  | ride_id (optional)         |
| `cancel`   | Cancel ride                         | ride_id (optional)         |
| `history`  | Recent rides                        | —                          |
| `auth_url` | Get OAuth setup URL                 | —                          |

## Env

- `UBER_SERVER_TOKEN` — for estimates + products (required)
- `UBER_ACCESS_TOKEN` — for ride requests (after OAuth)
- `UBER_CLIENT_ID` / `UBER_CLIENT_SECRET` — for OAuth flow
- `UBER_USE_PRODUCTION` — set `true` for real rides (default: sandbox)

See `uber-api-setup.md` for full OAuth2 setup guide.

## Examples

- "Get me an Uber to the airport"
- "How much would an Uber to Manhattan cost?"
- "Where's my Uber?"
- "Cancel my Uber"
