---
id: so-01-myzone-20260611
title: "MYZONE — Heart Rate & Effort Tracking Platform"
type: service
domain: SO
folder: "00-24 SO (Sovereign Operations)/01 The Fabric (Node Registry)"
status: active
created: 2026-06-11
modified: 2026-06-11
owner: zeus
tags: [domain/so, folder/01, service/myzone, health/heart-rate, health/fitness, hardware/wearable, api/partner]
---

# MYZONE — Heart Rate & Effort Tracking Platform

> Chest-strap heart rate monitoring system built around effort-based tracking (MEPs). Active ecosystem with partner API, webhooks, and Zapier integration.

---

## Hardware Owned

### MYZONE MZ-3 BLE (Model: MZ3-BLE)
- **Status: ✅ Active — Keep & Use**
- Chest strap HR monitor — BLE + ANT+ dual protocol
- 16-hour internal memory (stores workouts without phone nearby)
- Water resistant: 10m (IPX7)
- Battery: CR2032 (user-replaceable, ~6 months life)
- App confirmed compatible as of **MYZONE app v9.0.1 (November 2024)**
- Earns MEPs in the full MYZONE ecosystem
- Replacement straps still sold by MYZONE directly
- Superseded commercially by MZ-Switch (2022) but **not abandoned** — full app support continues

**Feature Icons (from box):**
- 16 HR Memory | MEPs | kcal | ANT+ | Water Resist 10m
- Bluetooth® | Bluetooth® 100% Live Display | Bluetooth® Upload | 5.3 kHz

---

## App

| Platform | Link | Last Updated |
|----------|------|-------------|
| iOS | [App Store](https://apps.apple.com/us/app/myzone/id496208945) | Nov 21, 2024 (v9.0.1) |
| Android | [Google Play](https://play.google.com/store/apps/details?id=com.myzone.myzone) | Active |

**Supported devices:** MZ-3 BLE, MZ-5, MZ-Switch, MZ-Open
**Requires:** MYZONE belt or Bluetooth HR monitor
**Free 30-day trial** — subscription required for full features

---

## MEPs — MYZONE Effort Points

MYZONE's proprietary effort metric. Calculated as heart rate as a % of max HR:
- Higher % of max HR = more MEPs
- Accumulate over time across all workouts
- Used for challenges, leaderboards, progress tracking
- 5 effort zones: **Gray → Blue → Green → Yellow → Red**

---

## Current Product Lineup (2026)

| Model | Status | Key Features | Price |
|-------|--------|-------------|-------|
| **MZ-3 BLE** | Discontinued (still app-supported) | Chest only, BLE + ANT+, 6-mo battery | ~$80-100 used |
| **MZ-Switch** | ✅ Current flagship | Chest OR wrist, optical HR, BLE + ANT+ + 5kHz, 4-mo battery | $149.99 |
| **MZ-Open** | ✅ Current | Open-ear headphone with HR monitoring | $119.99 |

**MZ-Switch upgrades over MZ-3:**
- Optical HR sensor (no gel required)
- Wrist-wearing option
- 5kHz analog signal for older gym cardio equipment
- Tradeoff: slightly shorter battery life

---

## API

**Type:** Partner API (REST, v2) — not public; requires partner approval

**Access:** Email **developer@myzone.org** with:
- Name + organization
- Intended use case
- OAuth 2.0 redirect URI(s)

Upon approval → receive `client_id` + `client_secret`

### Base URLs
| Environment | URL |
|------------|-----|
| Production | `https://api.myzone.org` |
| Sandbox | `https://sandbox.api.myzone.org` |

### Authentication
**OAuth 2.0** — two grant types:
- **Authorization Code** — user-level access (individual data with user consent)
  - Auth URL: `https://api.myzone.org/oauth/authorize`
  - Token URL: `https://api.myzone.org/oauth/token`
- **Client Credentials** — facility-level access (aggregated member data, no per-user auth)

**Scopes:** `read:user` | `read:workouts` | `read:meps` | `read:facility` | `read:activities`

### Endpoints (v2)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/v2/me` | Current authenticated user profile |
| GET | `/v2/user/{userId}/workout_summary` | Aggregated workout totals (calories, MEPs, effort zones) |
| GET | `/v2/user/{userId}/workout_sessions` | Detailed session data — **HR, calories, MEPs** |
| GET | `/v2/user/{userId}/workout_meps` | MEP data per workout |
| GET | `/v2/user/{userId}/workout_activity_summary` | Activity type breakdown |
| GET | `/v2/facility/{facilityId}/workout_summary` | All member workout data for a facility |
| GET | `/v2/facility/{facilityId}/workout_mep_summary` | MEP summary across facility members |
| GET | `/v2/facility/{facilityId}/activity_summary` | Activity breakdown for a facility |

### Data Available via API

| Data Type | Available |
|-----------|-----------|
| Heart Rate (bpm) | ✅ |
| Calories burned | ✅ |
| MEPs (Effort Points) | ✅ |
| Effort zone % (Gray/Blue/Green/Yellow/Red) | ✅ |
| Steps | ✅ |
| Distance | ✅ |
| User profile (name, DOB, sex) | ✅ |
| Facility/member aggregated data | ✅ |

---

## Webhooks (Native)

Real-time JSON POST to a partner-provided endpoint. Includes signature header for verification.

| Event | Trigger |
|-------|---------|
| `workout_completed` | User finishes a workout |
| `mep_milestone_reached` | User hits a MEP milestone |
| `user_created` | New user account created |
| `user_updated` | User profile updated |

---

## Zapier Integration

Available at: https://zapier.com/apps/myzone/integrations

**Triggers:**
- New Workout — fires on workout completion
- New MEP Milestone — fires on milestone
- New Member — fires when member joins a facility

**Actions:** None currently available

**Compatible with:** Google Sheets, Airtable, Slack, HubSpot, Mailchimp, Salesforce, Notion, Trello + 7,000 more

*Note: May require partner/business credentials*

---

## Other Integrations

- **Apple Health** — MYZONE app syncs workout data natively
- **Google Fit** — MYZONE app syncs natively
- **ANT+** — MZ-3 BLE broadcasts to ANT+ devices (Garmin watches, etc.)

---

## Documentation

| Resource | URL |
|----------|-----|
| Developer Docs | https://api.myzone.org/docs/ |
| API Reference (Apiary) | https://myzone.docs.apiary.io/ |
| Postman Collection | https://www.postman.com/myzone-api/myzone-api/overview |
| Getting Started | https://api.myzone.org/docs/getting-started.html |
| Concepts | https://api.myzone.org/docs/concepts.html |
| Authorization | https://api.myzone.org/docs/authorization.html |
| Developer Contact | developer@myzone.org |
| Support | https://support.myzone.org |
| Shop (straps/accessories) | https://shop.myzone.org |

---

## OneMind Use Cases

- **Personal fitness tracking** — MZ-3 BLE → MYZONE app → Apple Health pipeline
- **Biometric data ingestion** — Partner API to pull HR + MEPs into OneMind health layer (HP domain)
- **Automation** — Webhooks → NATS → OneMind event bus on `workout_completed`
- **Dashboards** — API data → WorldView or Homepage dashboard widgets

---

## Integration Status

- [x] Hardware confirmed working (MZ-3 BLE)
- [x] App support verified (v9.0.1, Nov 2024)
- [ ] Partner API credentials requested (developer@myzone.org)
- [ ] Webhook endpoint configured in OneMind
- [ ] NATS bridge for `workout_completed` events
- [ ] HP domain health data pipeline

---

## Manufacturer

**MYZONE Limited**
Falcon House, Ridgeway Street
Douglas, Isle of Man, British Isles IM1 1EL
Website: [myzone.org](https://myzone.org)
Trademarks: MYZONE® is a registered trademark of MYZONE Limited
