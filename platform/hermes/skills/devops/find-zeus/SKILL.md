---
name: find-zeus
description: "Use when geography, availability or operational context about Zeus matters: where he is, whether his TAK client is online, battery, or sending him a map marker or message."
version: 2.0.0
author: Legacy
triggers:
  - "where is zeus"
  - "zeus location"
  - "find zeus"
  - "zeus status"
  - "is zeus online"
  - "zeus battery"
  - "alert zeus"
  - "notify zeus"
---

# Find Zeus

Zeus's position and status come from his ATAK client, which reports to TAK
Server as callsign `ZEUS`. Use this for requests where place, travel time or
availability changes the answer (directions, local services, scheduling, device
work). It is Zeus's own location: never share it outside the owner's channels and
never write coordinates into Git, memory files or cron output.

## Reach and authority

- Reading position and contacts is a read-only API call; no approval needed.
- TAK is tailnet-only (`tak.tail717ba0.ts.net`); the public listener and old
  load balancer are retired. The cloud agent reaches the Marti API from inside the
  cluster. Zeus's phone appears only while ATAK is running and has network. When
  it is absent, say "not reporting" and stop; do not infer a position from stale
  data without labelling it stale.
- Sending Zeus a marker or chat message is a software message to the owner and is
  fine. Anything that acts on the physical world because of a location (opening
  a door, starting a vehicle, moving a drone, unlocking, payments) needs Zeus to
  write "execute" or "override" for that specific action (SOUL rule 8). A location
  reading is never authority.
- Never delete tracks, contacts or CoT history (SOUL rule 4).

## Identity (resolve live, do not trust memory)

Callsign `ZEUS`. Get the current UID from `contacts/all` by callsign; the UID
changes when the device or ATAK install changes. An old iPhone UID is dead.

## Reading position

Marti REST needs a client certificate (a credential). The agent has no standing
cert in Git; see `tak-server` for the approved path and ask Zeus when a cert is
needed. Shape of the calls, from inside the cluster:

- Is he connected: `GET /Marti/api/contacts/all` (present means ATAK active).
- Latest position: `GET /Marti/api/cot/xml/<uid>`.

CoT also flows on NATS as `tak.cot.>`; a subscription filtered on his UID is the
bus-native way to watch. A read-only helper script may exist in the Hermes
scripts directory; check before writing one.

## Interpreting readings

Report DATA separately from INFERENCE. Never assume driving.

| Speed | Meaning |
|---|---|
| 0 or -1 (no GPS) | Stationary or unavailable; treat -1 as unknown |
| 0.1 to 2 m/s | Walking |
| 2 to 8 m/s | Cycling or jogging |
| 8 to 15 m/s | City driving or transit |
| above 15 m/s | Highway or rail |

- Rail versus road: no stops over 2 km in an urban area, a locked course (within
  5 degrees over 3 km), or starting from a known station points to rail. Stop-go
  every 200 to 400 m points to surface streets.
- Take at least 3 readings 30 seconds apart before concluding anything.
  Correlate with what is physically at the coordinates and the time of day.
- A CoT older than 5 minutes means ATAK is backgrounded. `hae` of -25 or
  9999999 is unreliable; ignore altitude.
- Battery: under 20 percent, flag it; under 10 percent, critical.
- Emergency CoT types contain `-E` (for example `a-f-G-U-C-E`). Watch
  `tak.cot.>` for them and alert immediately.

Report format:

```
DATA: lat, lon, speed, course, battery, time
LOCATION: place name
INFERENCE: likely activity and mode
CONFIDENCE: high / medium / low and why
```

## Sending Zeus something

Preferred path is the bus, not a direct TLS socket:

```bash
nats pub tak.cot.inject '<CoT XML>'
```

- A waypoint marker: type `b-m-p-w-GOTO`, a `<contact callsign="LABEL"/>`, a
  `<remarks>`, and `<marti><dest callsign="ZEUS"/></marti>` so only he sees it.
  Set `stale` a couple of hours ahead.
- A GeoChat message: type `b-t-f` with a `__chat` detail addressed to his UID.
  Keep the text short and plain.
- Use the sender callsign `Legacy`. Full CoT shapes are in the `tak-server`
  skill (`references/tak-cot-format.md`).

## Pitfalls

- The API pod restarts now and then; the messaging path through NATS is steadier.
- If the API is down, ask Zeus before querying the TAK database directly. Its
  credentials live in a SealedSecret in the `takserver` namespace and are never
  copied into skills.
- Use the current callsign `ZEUS`, not earlier ones.
