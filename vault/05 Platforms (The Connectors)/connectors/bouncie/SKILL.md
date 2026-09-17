---
name: bouncie
description: Track vehicles and trips via Bouncie GPS API. Use when you need current vehicle location, trip history, odometer, speed, battery status, or to set up geo-zones. Vehicle GPS data can feed FIND ZEUS MODE as a location source.
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Bash(python3:*)
  - Read
  - Write
---

# Bouncie Skill

Vehicle GPS tracking and trip data via [Bouncie API](https://docs.bouncie.dev).

## Environment Variables

```bash
BOUNCIE_ACCESS_TOKEN   # OAuth2 access token (from portal or refresh flow)
BOUNCIE_REFRESH_TOKEN  # OAuth2 refresh token
BOUNCIE_CLIENT_ID      # App client ID (Bouncie Developer Portal)
BOUNCIE_CLIENT_SECRET  # App client secret
```

> Credentials live in 1Password vault `00-24 SO (Sovereign Operations)` → item `Bouncie API`

## Base URL

```
https://api.bouncie.dev/v1
```

All calls require:

```bash
-H "Authorization: $BOUNCIE_ACCESS_TOKEN"
-H "Content-Type: application/json"
```

---

## Quick Reference

| Operation                | Endpoint                                                                          |
| ------------------------ | --------------------------------------------------------------------------------- |
| List all vehicles        | `GET /vehicles`                                                                   |
| Get vehicle by IMEI      | `GET /vehicles?imei={imei}`                                                       |
| Get vehicle by VIN       | `GET /vehicles?vin={vin}`                                                         |
| List trips               | `GET /trips?imei={imei}&gps-format=polyline&starts-after={iso}&ends-before={iso}` |
| List geo-zones           | `GET /geozones`                                                                   |
| Create geo-zone location | `POST /locations`                                                                 |
| Create geo-zone          | `POST /geozones`                                                                  |

---

## Operations

### List All Vehicles

```bash
curl -s "https://api.bouncie.dev/v1/vehicles" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | jq '.'
```

Response fields per vehicle:

- `imei` — device identifier (use this for trip queries)
- `vin` — vehicle identification number
- `make`, `model`, `year`, `color`
- `stats.lastTrip` — last trip info
- `stats.location` — `{lat, lon, heading, speed, timestamp}`
- `stats.odometer` — `{value, units}`
- `stats.battery` — `{status, value}` (good/low/critical)
- `stats.mil` — engine check light (`{status, description}`)

### Get Current Location

```bash
# Get all vehicles then extract location
curl -s "https://api.bouncie.dev/v1/vehicles" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | \
  jq '[.[] | {imei, vin, make, model, location: .stats.location}]'
```

### Get Latest Trip

```bash
IMEI="your-device-imei"
curl -s "https://api.bouncie.dev/v1/trips?imei=$IMEI&gps-format=polyline" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | jq '.[0]'
```

### Get Trips for Date Range

```bash
IMEI="your-device-imei"
START="2026-04-01T00:00:00Z"
END="2026-04-12T23:59:59Z"

curl -s "https://api.bouncie.dev/v1/trips?imei=$IMEI&gps-format=polyline&starts-after=$START&ends-before=$END" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | jq '[.[] | {
    tripId, distance, duration,
    startTime: .start.time, startAddress: .start.formattedAddress,
    endTime: .end.time, endAddress: .end.formattedAddress,
    avgSpeed: .avgSpeed, maxSpeed: .maxSpeed
  }]'
```

### Get Vehicle Health

```bash
curl -s "https://api.bouncie.dev/v1/vehicles" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | \
  jq '[.[] | {
    vehicle: "\(.year) \(.make) \(.model)",
    battery: .stats.battery.status,
    mil: .stats.mil.status,
    odometer: .stats.odometer
  }]'
```

### Refresh Access Token

```bash
curl -s -X POST "https://auth.bouncie.com/oauth/token" \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": "'"$BOUNCIE_CLIENT_ID"'",
    "client_secret": "'"$BOUNCIE_CLIENT_SECRET"'",
    "grant_type": "refresh_token",
    "refresh_token": "'"$BOUNCIE_REFRESH_TOKEN"'"
  }' | jq '{access_token, refresh_token, expires_in}'
```

---

## Geo-Zones (Application)

### Create a Location (for geo-zone)

```bash
# Circle geo-zone (lat/lon + radius in meters)
curl -s -X POST "https://api.bouncie.dev/v1/locations" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home - Greenwich",
    "type": "circle",
    "lat": 41.0534,
    "lon": -73.6287,
    "radius": 200
  }' | jq '{locationId, name}'
```

### Create a Geo-Zone

```bash
# After creating location, create the geo-zone
curl -s -X POST "https://api.bouncie.dev/v1/geozones" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Home Zone",
    "locationId": "LOCATION_ID_FROM_ABOVE",
    "direction": "both"
  }' | jq '.'
```

### List Geo-Zones

```bash
curl -s "https://api.bouncie.dev/v1/geozones" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | jq '.'
```

---

## FIND ZEUS MODE Integration

Bouncie is the **vehicle location source** for FIND ZEUS MODE.
Subject namespace: `zeus.presence.vehicle`

Publish vehicle location to NATS:

```python
import subprocess, json, os

def publish_vehicle_location():
    token = os.environ["BOUNCIE_ACCESS_TOKEN"]

    result = subprocess.run([
        "curl", "-s", "https://api.bouncie.dev/v1/vehicles",
        "-H", f"Authorization: {token}"
    ], capture_output=True)

    vehicles = json.loads(result.stdout)

    for v in vehicles:
        loc = v.get("stats", {}).get("location", {})
        if not loc:
            continue
        payload = {
            "source": "bouncie",
            "vehicle": f"{v.get('year')} {v.get('make')} {v.get('model')}",
            "lat": loc.get("lat"),
            "lon": loc.get("lon"),
            "speed": loc.get("speed"),
            "heading": loc.get("heading"),
            "timestamp": loc.get("timestamp")
        }
        # Publish to NATS
        subprocess.run([
            "nats", "pub", "zeus.presence.vehicle",
            json.dumps(payload)
        ])
```

---

## Webhook Events

Configure webhook URL on [Bouncie Developer Portal](https://www.bouncie.dev/).

| Event                | Trigger                            |
| -------------------- | ---------------------------------- |
| `tripStart`          | Vehicle starts moving              |
| `tripData`           | GPS pings during trip (continuous) |
| `tripMetrics`        | Summary at trip end                |
| `tripEnd`            | Vehicle stops                      |
| `batteryLow`         | Battery drops below threshold      |
| `deviceConnected`    | OBD device plugged in              |
| `deviceDisconnected` | OBD device unplugged               |
| `geoZoneEnter`       | Vehicle enters a geo-zone          |
| `geoZoneExit`        | Vehicle leaves a geo-zone          |

Webhook payload `Authorization` header = your security key.


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
