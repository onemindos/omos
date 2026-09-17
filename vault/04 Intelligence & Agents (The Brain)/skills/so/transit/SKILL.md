---
name: transit
description: "Plans NYC transit routes using MTA GTFS-RT real-time data and Google Maps, auto-locating Zeus via Bouncie/HA before routing. Use when Zeus says 'how do I get home', 'subway home', 'next train', 'get me to office', 'how long to [place]', 'transit', or 'MTA'."
allowed-tools:
  - Bash(curl:*)
  - Bash(jq:*)
  - Bash(python3:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Transit Skill

NYC subway, bus, and walking directions powered by MTA GTFS-RT + Google Maps.

## Environment Variables

```bash
MTA_API_KEY            # api.mta.info key (free, register at api.mta.info)
GOOGLE_MAPS_API_KEY    # Google Maps Platform — Directions + Places APIs
BOUNCIE_ACCESS_TOKEN   # Vehicle GPS (optional — for car location)
HOME_ASSISTANT_URL     # HA instance URL
HOME_ASSISTANT_TOKEN   # HA long-lived token
```

> Credentials in 1Password → `00-24 SO` → items `MTA API` and `Google Maps API`

---

## The 4-Step Transit Flow

When Zeus asks for transit help, run these steps **in order**. Ask for missing info rather than assuming.

---

### STEP 1 — Intent & Destination

Ask: **"Where to?"**

Resolve destination from shorthand:

| Shorthand         | Resolves to                                                |
| ----------------- | ---------------------------------------------------------- |
| `home`            | Greenwich, CT — 40 Field Point Circle, Greenwich, CT 06830 |
| `nyc` / `office`  | NYC Office — ask for specific address if not stored        |
| `upstate`         | Upstate Homestead — ask for address                        |
| `mexico` / `farm` | Mexico Farm — international, no transit                    |
| Custom address    | Use as-is                                                  |

Known destinations (update as Zeus adds more):

```
HOME:    40 Field Point Circle, Greenwich, CT 06830
OFFICE:  [ask Zeus — not yet stored]
```

---

### STEP 2 — Get Current Location

Try sources **in this priority order** until one returns a valid lat/lon:

#### Source A: Bouncie (vehicle GPS — most precise if driving/parked)

```bash
curl -s "https://api.bouncie.dev/v1/vehicles" \
  -H "Authorization: $BOUNCIE_ACCESS_TOKEN" | \
  jq '[.[] | {vehicle: "\(.year) \(.make) \(.model)", lat: .stats.location.lat, lon: .stats.location.lon, speed: .stats.location.speed, ts: .stats.location.timestamp}]'
```

If speed > 0 the vehicle is moving. Use this location.

#### Source B: Home Assistant presence sensor

```bash
# Get Zeus's person/device_tracker entity
curl -s "$HOME_ASSISTANT_URL/api/states/person.zeus" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" | \
  jq '{state, lat: .attributes.latitude, lon: .attributes.longitude, source: .attributes.source}'

# Or phone device tracker directly
curl -s "$HOME_ASSISTANT_URL/api/states/device_tracker.zeus_phone" \
  -H "Authorization: Bearer $HOME_ASSISTANT_TOKEN" | \
  jq '{state, lat: .attributes.latitude, lon: .attributes.longitude}'
```

HA presence zones resolve to:

- `home` → Greenwich home coordinates
- `not_home` → use GPS coordinates from attributes
- Named zone → use zone center coordinates

#### Source C: Ask Zeus directly

If both APIs are unavailable or return stale data, ask:

> "What's your current location? (address, intersection, or neighborhood)"

---

### STEP 3 — Find Nearby Stations

#### Google Places — Subway stations near current location

```bash
LAT="40.7128"
LON="-74.0060"
RADIUS=800   # meters (~0.5 mile walk)

curl -s "https://maps.googleapis.com/maps/api/place/nearbysearch/json?\
location=$LAT,$LON&radius=$RADIUS&type=subway_station&key=$GOOGLE_MAPS_API_KEY" | \
  jq '[.results[] | {
    name,
    vicinity,
    distance: "TODO",
    place_id
  }]'
```

#### Walking time to each station (Google Directions)

```bash
ORIGIN_LAT="40.7128"
ORIGIN_LON="-74.0060"
STATION_ADDRESS="Grand Central Terminal, New York, NY"

curl -s "https://maps.googleapis.com/maps/api/directions/json?\
origin=$ORIGIN_LAT,$ORIGIN_LON\
&destination=$(python3 -c "import urllib.parse; print(urllib.parse.quote('$STATION_ADDRESS'))")\
&mode=walking\
&key=$GOOGLE_MAPS_API_KEY" | \
  jq '.routes[0].legs[0] | {
    duration: .duration.text,
    distance: .distance.text,
    steps: [.steps[] | .html_instructions | gsub("<[^>]+>"; "")]
  }'
```

---

### STEP 4 — MTA Real-Time Departures

#### MTA GTFS-RT Feed URLs (New York)

| Lines         | Feed URL                                                                             |
| ------------- | ------------------------------------------------------------------------------------ |
| 1 2 3 4 5 6 S | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/123456S` |
| A C E         | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/ACE`     |
| B D F M       | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/BDFM`    |
| G             | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/G`       |
| J Z           | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/JZ`      |
| L             | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/L`       |
| N Q R W       | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/NQRW`    |
| 7             | `https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/7`       |

#### Real-time arrivals via GTFS-RT protobuf

```python
#!/usr/bin/env python3
"""
Get real-time subway arrivals for a station.
Requires: pip install gtfs-realtime-bindings requests
"""
import os, time, requests
from google.transit import gtfs_realtime_pb2

MTA_KEY = os.environ["MTA_API_KEY"]
FEED_URLS = {
    "123456S": "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/123456S",
    "ACE":     "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/ACE",
    "BDFM":    "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/BDFM",
    "NQRW":    "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/NQRW",
    "L":       "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/L",
    "7":       "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/7",
    "JZ":      "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/JZ",
    "G":       "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/line/G",
}

# GTFS-RT protobuf feed (for arrival times)
GTFS_RT_URLS = {
    "123456S": "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs",
    "ACE":     "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-ace",
    "BDFM":    "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-bdfm",
    "NQRW":    "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-nqrw",
    "L":       "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-l",
    "7":       "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-7",
    "JZ":      "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-jz",
    "G":       "https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-g",
}

def get_arrivals(stop_id: str, feed_key: str = "123456S", limit: int = 5):
    """
    stop_id: MTA stop ID (e.g., "142N" = 72nd St uptown, "R16N" = 34th-Herald Sq uptown)
    Get station stop IDs from: https://gtfs.mta.info/static/stops.txt
    """
    url = f"https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs"
    if feed_key != "123456S":
        url = f"https://api-endpoint.mta.info/Datastore/datamine/api/xml/SubwayStatus/gtfs-{feed_key.lower()}"

    resp = requests.get(url, headers={"x-api-key": MTA_KEY})
    feed = gtfs_realtime_pb2.FeedMessage()
    feed.ParseFromString(resp.content)

    now = time.time()
    arrivals = []

    for entity in feed.entity:
        if not entity.HasField("trip_update"):
            continue
        for stop_time in entity.trip_update.stop_time_update:
            if stop_time.stop_id == stop_id:
                arr_time = stop_time.arrival.time if stop_time.arrival.time else stop_time.departure.time
                mins = round((arr_time - now) / 60)
                if 0 <= mins <= 60:
                    arrivals.append({
                        "route": entity.trip_update.trip.route_id,
                        "minutes": mins,
                        "trip_id": entity.trip_update.trip.trip_id
                    })

    return sorted(arrivals, key=lambda x: x["minutes"])[:limit]

if __name__ == "__main__":
    # Example: Grand Central (42 St-GCT, 4/5/6 uptown)
    arrivals = get_arrivals("640N", "123456S")
    for a in arrivals:
        print(f"  {a['route']} train — {a['minutes']} min")
```

---

## Common MTA Stop IDs (NYC)

| Station              | Lines             | Stop ID (N=uptown, S=downtown)    |
| -------------------- | ----------------- | --------------------------------- |
| Grand Central–42 St  | 4 5 6             | `640N` / `640S`                   |
| Times Sq–42 St       | N Q R W           | `R16N` / `R16S`                   |
| Penn Station–34 St   | 1 2 3             | `137N` / `137S`                   |
| Fulton St            | 4 5 A C           | `A36N` / `A36S`                   |
| 72 St                | 1 2 3             | `142N` / `142S`                   |
| 86 St                | 1                 | `140N` / `140S`                   |
| 96 St                | 1 2 3             | `A27N` / `A27S`                   |
| 14 St–Union Sq       | 4 5 6 L N Q R W   | `635N` / `635S`                   |
| Atlantic Av–Barclays | 2 3 4 5 B D N Q R | `235N` / `235S`                   |
| Jamaica (LIRR)       | —                 | for Metro-North use separate feed |

Full stop list: `https://gtfs.mta.info/static/stops.txt`

---

## Full Trip Plan — Google Directions (Transit Mode)

```bash
ORIGIN="40.7580,-73.9855"          # current location
DESTINATION="Greenwich+CT+06830"   # destination
DEPART_TIME=$(date +%s)            # now (Unix timestamp)

curl -s "https://maps.googleapis.com/maps/api/directions/json?\
origin=$ORIGIN\
&destination=$DESTINATION\
&mode=transit\
&transit_mode=subway|train\
&departure_time=$DEPART_TIME\
&key=$GOOGLE_MAPS_API_KEY" | \
  jq '.routes[0] | {
    summary: .summary,
    duration: .legs[0].duration.text,
    arrival: .legs[0].arrival_time.text,
    departure: .legs[0].departure_time.text,
    fare: .fare.text,
    steps: [.legs[0].steps[] | {
      mode: .travel_mode,
      instruction: (.html_instructions | gsub("<[^>]+>"; "")),
      duration: .duration.text,
      distance: .distance.text,
      transit: (if .transit_details then {
        line: .transit_details.line.short_name,
        from: .transit_details.departure_stop.name,
        to: .transit_details.arrival_stop.name,
        departs: .transit_details.departure_time.text,
        arrives: .transit_details.arrival_time.text,
        stops: .transit_details.num_stops
      } else null end)
    }]
  }'
```

---

## Home Route: NYC → Greenwich CT

Zeus's most common route: NYC → Home (Greenwich, CT)

**Options ranked by transit:**

1. **Metro-North (New Haven Line)** — Grand Central → Greenwich Station (~47 min express)
   - Grand Central Terminal, 42nd St & Park Ave
   - Trains run every 30-60 min peak, hourly off-peak
   - Check: `https://as0.mta.info/mnr/mobiledev/DepartureCurrentHTML.cfm?stationID=1&destinationID=67`

2. **Subway + Metro-North** — Local subway to Grand Central → Metro-North
   - From Midtown: Any 4/5/6 to Grand Central (42 St)
   - From Downtown: A/C/E to 42 St, walk to GCT

3. **Express Bus (BxM series)** — Less common for Greenwich route

### Metro-North Departures (Grand Central → Greenwich)

```bash
# MTA Metro-North real-time departures
curl -s "https://api-endpoint.mta.info/Datastore/datamine/api/xml/mnr/mnrDepart.xml?\
station=1&destination=67&key=$MTA_API_KEY" | \
  python3 -c "
import sys, xml.etree.ElementTree as ET
tree = ET.parse(sys.stdin)
for trip in tree.findall('.//TripList/Trip')[:5]:
    time = trip.findtext('ScheduledDepartureTime','?')
    track = trip.findtext('Track','?')
    status = trip.findtext('Status','On Time')
    est = trip.findtext('EstimatedDepartureTime', time)
    print(f'  {time} → {est} | Track {track} | {status}')
"
```

---

## Agent Response Format

When reporting transit options to Zeus, use this format:

```
📍 Current: [location name] → 🏁 [destination]

🚇 Option 1: [subway/train line]
   Depart: [time] from [station] (walk [X] min)
   Arrive: [time] at [destination]
   Total:  [X] min

🚇 Option 2: [alternate route]
   ...

⚡ Fastest: Option [N] — [X] min total
```

---

## GTFS Static Data

Download and parse full NYC subway stops:

```bash
# Download GTFS static (stops, routes, etc.)
curl -s "https://gtfs.mta.info/static/stops.txt" | \
  python3 -c "
import sys, csv
reader = csv.DictReader(sys.stdin)
for row in reader:
    if 'Grand Central' in row.get('stop_name',''):
        print(row['stop_id'], row['stop_name'], row['stop_lat'], row['stop_lon'])
"
```

---

## Install Dependencies

```bash
pip3 install gtfs-realtime-bindings requests
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
