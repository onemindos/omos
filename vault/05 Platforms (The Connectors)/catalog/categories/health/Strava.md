# Strava

> **Activity tracking platform for running, cycling, and endurance sports**

## Overview

Strava is the social network for athletes. Tracks GPS-based workouts (running, cycling, swimming, hiking), computes performance metrics, and stores a complete activity history. REST API enables agent integration for automatic HP domain logging.

## Status

| Attribute   | Value                        |
| ----------- | ---------------------------- |
| Status      | ✅ Active                     |
| Category    | Health / Fitness             |
| Type        | Activity Tracking            |
| Integration | REST API, iOS App, Webhooks  |

## Use Cases in OneMind

1. **Workout Logging** — Auto-sync GPS workouts to HP domain
2. **Performance Metrics** — Pace, power, HR zones, VO2 max trends
3. **Segment Racing** — Track PRs on recurring routes
4. **N8N Integration** — Webhook triggers new activity → N8N → HP metric update
5. **Composio Skill** — Access Strava data via Composio tool calls

## Key API Endpoints

```bash
# Get recent activities
GET https://www.strava.com/api/v3/athlete/activities

# Get activity details
GET https://www.strava.com/api/v3/activities/{id}

# Subscribe to webhooks
POST https://www.strava.com/api/v3/push_subscriptions
```

## Links

- Website: [strava.com](https://strava.com)
- API Docs: [developers.strava.com](https://developers.strava.com)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
