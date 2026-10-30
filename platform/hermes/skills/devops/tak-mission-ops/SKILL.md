---
name: tak-mission-ops
description: "Use to run TAK missions on OneMind: mission and ExCheck APIs, the signal protocol with the owner, CoT message templates (waypoint, GeoChat, emergency) injected over NATS, tracking, client choice and degraded comms."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [tak, missions, c2, tracking, geochat, excheck, meshtastic, comms, nats]
    related_skills: [tak-server, tak-mission-runner, mission-replay, nats-operations, android-adb-management]
triggers:
  - "TAK mission"
  - "drop a waypoint"
  - "GeoChat"
  - "ExCheck"
  - "mission log"
  - "MEDEVAC"
---

# TAK mission operations

Rewritten 2026-10-29 from a July 2026 guide for TAK 5.7 behind a public load
balancer. That guide held the server address, the admin login and certificate
passphrases and used a raw TLS socket for injection. This version follows the
current topology: TAK Server 5.8, **tailnet only** (`tak.tail717ba0.ts.net`),
and every inject and read goes through the **NATS bus**. Topology, REST and
operational pitfalls are in `tak-server`; read it first.

## Rules of the road

- **Inject** = publish CoT XML to `tak.cot.inject`. **Read** = subscribe to
  `tak.cot.out` (stream `COT` also keeps `tak.cot.>`). Never publish
  `evt.tak.cot.*`; do not invent subjects. Mission lifecycle events for the
  rest of the fabric go on `evt.mission.{mission}.{event}` (grammar:
  `evt.{domain}.{entity}.{event}`).
- A publish is not proof of display. Prove it in stages: published, appears on
  `tak.cot.out`/the `COT` stream, then confirmed on a client (ask the owner or
  check the contact list). Say which stage you verified.
- Certificates, enrollment and user creation are owner-led (TAK Portal at
  `enroll.tail717ba0.ts.net`, runbook `platform/takserver/DEPLOY-RUNBOOK.md`).
  Credentials live in SealedSecrets (`tak-server-stores`, `tak-ca-credentials`);
  never ask for, print or guess a passphrase.
- Sending CoT to someone's map is communication on the owner's behalf: do it for
  the missions and people the owner named. Emergency alerts are never a test.
- Anything physical (drones, doors, vehicles) triggered from a map marker needs
  "execute" or "override" from the owner for that action.

## Server facts to rely on (verify with kubectl)

Namespace `takserver`; API service `takserver-api-service` (8443 mutual TLS,
8446 WebTAK), messaging `takserver-messaging-service` (8089), PostgreSQL
`takserver-postgresql` (database `cot`). Clients reach it over the tailnet on
8089, 8443, 8446 and 64738 (voice). The `cot-retention` CronJob keeps 24 hours of
`cot_router`.

## Missions through the REST API

Done from a host that holds an owner-provisioned client certificate (this pod
does not by default; ask). Base `https://takserver-api-service.takserver.svc.cluster.local:8443`.

```text
GET  /Marti/api/contacts/all | subscriptions/all | groups/all | version/config
GET  /Marti/api/missions?full=true
POST /Marti/api/missions?name=X&tool=public&group=__ANON__         # create
GET  /Marti/api/missions/{name}/{log|cot|contents|changes|subscriptions}
POST /Marti/api/missions/{name}/log        {content, dtg, entryUid, creatorUid}
PUT  /Marti/api/missions/{name}/subscription?uid=<uid>
PUT  /Marti/api/missions/{name}/uid?uid=<cot uid>                  # attach a CoT UID
POST /Marti/api/excheck/template           # XML body, POST-only
GET  /Marti/api/sync/search?keywords=X ; POST /Marti/sync/upload ; GET /Marti/sync/content?hash=<sha256>
GET|POST /Marti/api/video                  # JSON feed registration (not /Marti/vcm)
```

Mission names with spaces are URL-encoded (`Legacy%20Ops`). Deleting a mission
is a deletion: it needs "execute" or "override". If the API is down but
messaging is up, CoT inject and GeoChat still work over the bus; ExCheck, mission
CRUD and file sync do not (use the signal checklist below instead).

An ExCheck template looks like
`<MissionChecklist><checklistDetails name= description= startTime= endTime=><checklistItem uid="1" checked="false"><checklistDetails name="Step 1"/></checklistItem></checklistDetails></MissionChecklist>`.

## CoT message templates (publish to `tak.cot.inject`)

All times are UTC ISO 8601 (`2026-10-29T14:00:00.000Z`). Keep `stale` short for
live items (60 s for a track, hours for a waypoint) and use the `uid` of the
thing reported. Remove a marker by re-sending the same UID with `stale` equal
to `time`.

```bash
# Waypoint to a named contact (callsign from contacts/all, never a hard-coded personal UID)
nats pub tak.cot.inject '<event version="2.0" uid="LEGACY-WP-001" type="b-m-p-w-GOTO" time="NOW" start="NOW" stale="LATER" how="h-g-i-g-o"><point lat="LAT" lon="LON" hae="0" ce="10" le="10"/><detail><contact callsign="LABEL"/><remarks>TEXT</remarks><color argb="-16711936"/><marti><dest callsign="TARGET-CALLSIGN"/></marti></detail></event>'
```

GeoChat to everyone: type `b-t-f`, uid `GeoChat.<sender-uid>.All Chat Rooms.<msgid>`,
`<__chat chatroom="All Chat Rooms" id="All Chat Rooms" senderCallsign="Legacy" ...>`
with `<chatgrp uid0=<sender-uid> uid1="All Chat Rooms" id="All Chat Rooms"/>`,
`<remarks to="All Chat Rooms">text</remarks>` and
`<marti><dest callsign="All Chat Rooms"/></marti>`. For a direct message, set the
chatroom, `id` and `dest` to the recipient's callsign and UID. Use
`nats pub subject "$(cat file)"`, never `< file` (publishes zero bytes).

Emergency (repeating types configured on the server: "911 Alert", "Ring The
Bell", "Geo-fence Breached", "Troops In Contact"): type `a-f-G-U-C-E` with
`<emergency type="Ring The Bell"/>`; cancel with `<emergency cancel="true"/>`.

| Type | Meaning |
|---|---|
| `a-f-G-U-C` | Friendly ground unit |
| `a-h-G` / `a-u-G` | Hostile / unknown ground |
| `b-m-p-w-GOTO` | Waypoint |
| `b-m-p-s-p-loc` | Search area |
| `b-t-f` | GeoChat |
| `b-i-x` | Sensor point of interest |
| `u-rb-a` | Range and bearing |

ARGB colours (signed ints): green `-16711936`, yellow `-256`, blue `-16776961`,
red `-65536`, cyan `-16711681`, white `-1`, orange `-23296`.

## Signal protocol with the owner

| Owner does on the map | You see | You do |
|---|---|---|
| Yellow marker with remarks | `tak.cot.out` / mission change | Research the location, send a waypoint back |
| Green marker "OBJ COMPLETE" | mission change | Log complete, acknowledge on the map |
| Red marker "HOLD" | `tak.cot.out` | Log hold, keep watching position |
| Chat to the Legacy room | `b-t-f` CoT | Process, reply in chat |
| Emergency button | `a-f-G-U-C-E` | Priority log, push support information, tell the owner on every channel |
| Team colour change | SA update | Treat as status |
| Blue marker "INTEL" | mission UID add | Log, cross-reference |

Live phase signals typed in chat (activating, moving, eyes on...) are in
`tak-mission-runner`. Review afterwards: `mission-replay`.

## Starting a mission and tracking

1. Subscribe yourself to the mission, log "MISSION START: <objective>", push a
   start waypoint to the owner (REST for the first two, bus for the third).
2. Watch `tak.cot.>`: `nats sub 'tak.cot.out' --count 20 --timeout 20s`.
3. Position of a person: the latest CoT for their device UID on `tak.cot.out`
   (filter on `uid`). Presence also appears in KV `zeus-presence` when something
   publishes it. A direct database query of `cot_router` (geometry `event_pt`,
   no callsign column, filter by `uid`) needs a read-only credential from a
   SealedSecret; ask for one rather than using an admin login.
4. Check the cadence: iTAK reports about every 30 s; do not poll faster.

## Signal-based checklist (when ExCheck is unavailable)

Define phase words (SP, MOVING, EYES ON, AT OBJ, PACKAGE, RTB, HOME). On each
word from the owner: look up the owner's position, drop a colour-coded waypoint
there, push a one-line GeoChat confirmation, acknowledge on the channel the word
arrived on, and publish `evt.mission.<mission>.phase` with the phase. The map
becomes the checklist.

## Clients (pick by role)

| Client | Role | Notes |
|---|---|---|
| ATAK / CivTAK (Android) | Field primary | Plugins, 9-line MEDEVAC, ExCheck, video, PTT; same binary |
| WinTAK / WinTAK-CIV | Fixed command post | Same capability as ATAK; multi-monitor, keyboard speed; wire it |
| WebTAK | Remote monitor | No plugins, no offline, needs connectivity |
| iTAK (iOS) | Limited | No plugins, MEDEVAC or video; text and markers only |
| TAK ICU | Admin panel | Not a map client |

See `references/atak-vs-itak-feature-gap.md` and
`references/atak-bidirectional-c2.md`. On ATAK the map is the interface: the
owner drops markers that you read as commands and you answer with markers.
Wearable builds (rugged Android tablet or SBC with ATAK, same server and COP;
a Garmin watch shows SA only) and the 24/7 WinTAK desk are owner-led hardware
decisions; offer the options, do not buy or flash anything.

## Degraded and off-grid comms

| Scenario | Path |
|---|---|
| LTE or Wi-Fi | EUD to server on 8089 over the tailnet, about 30 s SA |
| Satellite IP link | Same TLS path; the client reconnects by itself |
| Poor link | Client queues locally; server `smartRetry` can be enabled in CoreConfig (owner change) |
| Offline EUD | Keeps local SA, syncs on reconnect; pre-cache offline maps before a mission |
| Meshtastic LoRa | Needs at least two physical nodes and a bridge plugin; positions appear on the map |
| Federation | Mutual TLS between servers; FedHub is tailnet-only, ports not exposed |

## MEDEVAC 9-line (reference)

1 location, 2 radio frequency and callsign, 3 patients by precedence, 4 special
equipment, 5 patients by type (litter or ambulatory), 6 security at the pickup
site, 7 marking method, 8 patient nationality and status, 9 contamination.
ATAK has a 9-line form that fills line 1 from GPS. On iTAK type it into chat.
This is information only; real emergencies go to the emergency services.

## Pitfalls

- Marker not showing: stale already passed (clock skew), wrong group (users in
  different groups cannot see each other), or the client is not connected. Check
  `contacts/all` or ask the owner before re-sending in a loop.
- `callsign` is not a database column; it lives inside `detail`.
- Do not hold your own TLS connection to 8089; one socket per event is how the
  old workaround behaved, and a bus publish avoids the problem entirely.
- Do not run tap automation on ATAK; give the owner the manual steps.
