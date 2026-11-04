---
name: tak-mission-runner
description: "Use when the owner sends a live mission phase word (activating, moving, eyes on, going internal, jackpot, RTB, home): mark the phase on the TAK map over NATS and acknowledge in one line."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [tak, mission, signal, field-ops, gps, cot, nats]
    related_skills: [tak-mission-ops, tak-server, mission-replay, nats-operations]
triggers:
  - activating
  - moving
  - eyes on
  - going internal
  - jackpot
  - RTB
  - home
---

# TAK mission runner: live signal-response protocol

A live field protocol. When the owner is in a mission and sends a phase word,
act at once: mark it on the map, acknowledge in one line, no questions, no
filler. Rewritten 2026-10-29: the earlier version ran pod scripts with a
database login and a raw TLS socket, and was hard-wired to one past errand (a
named store address and a personal device UID). Those are gone. This version
publishes over NATS (`tak.cot.inject`), takes the objective from the owner for
each mission, and finds the owner's device UID from the live contact list.

Only start a mission protocol when the owner has said a mission is on, or sends
a phase word in a thread where one is active. A casual "I'm home" in
conversation is not a signal: ask once if unsure.

## Mission setup (once, from the owner)

Get these in the owner's first message or ask in one line:

- mission name (short, no spaces; used in `evt.mission.<name>.phase`),
- objective: label and lat/lon (the owner may drop a marker instead; read it),
- the owner's TAK callsign (find the UID in `contacts/all` or the latest
  `tak.cot.out` event for that callsign; do not store a UID in a skill).

Keep the active mission in the memory tool or `/opt/data/mission_active.json`
(name, objective, start time, current phase, active flag). Set `active:false` on HOME.

## Phase words

Match case-insensitively anywhere in the message. Natural speech ("Legacy, I'm
activating", "jackpot") counts.

| Word or phrase | Phase | Colour (ARGB int) |
|---|---|---|
| activating | SP, leaving base | green `-16711936` |
| moving, en route | EN ROUTE | yellow `-256` |
| eyes on | EYES ON target | orange `-23296` |
| going internal | AT OBJ, on station | blue `-16776961` |
| jackpot | PACKAGE secured | purple `-65281` |
| RTB, heading home, heading back | RTB | cyan `-16711681` |
| home, back at base | MISSION COMPLETE | green `-16711936` |

## Per signal, in order

1. **Position.** Latest CoT for the owner's UID from the bus (for example
   `nats sub tak.cot.out --count 20 --timeout 15s`, keep the newest event whose
   `uid` matches). If older than 120 s say "GPS stale Ns" in the acknowledgement.
2. **Distance and bearing** to the objective (haversine; cardinal from bearing):

   ```bash
   python3 - "$LAT" "$LON" "$OBJ_LAT" "$OBJ_LON" <<'PY'
   import sys, math
   a,b,c,d = map(float, sys.argv[1:5])
   p1,p2 = math.radians(a), math.radians(c); dl = math.radians(d-b)
   h = math.sin((p2-p1)/2)**2 + math.cos(p1)*math.cos(p2)*math.sin(dl/2)**2
   m = 2*6371000*math.asin(math.sqrt(h))
   y = math.sin(dl)*math.cos(p2); x = math.cos(p1)*math.sin(p2)-math.sin(p1)*math.cos(p2)*math.cos(dl)
   brg = (math.degrees(math.atan2(y,x))+360)%360
   print(f"{m:.0f}m {['N','NE','E','SE','S','SW','W','NW'][int((brg+22.5)//45)%8]} {brg:.0f}")
   PY
   ```

3. **Mark the map.** Publish a `b-m-p-w-GOTO` waypoint at the owner's position
   with the phase colour, uid `MISSION-<name>-<PHASE>`, callsign = phase name,
   stale 4 h, to `tak.cot.inject` (template in `tak-mission-ops`). On
   ACTIVATING also publish the objective marker (`a-f-G-U-C`, red, uid
   `MISSION-<name>-OBJ`).
4. **Tell the fabric.** `nats pub evt.mission.<name>.phase '{"phase":"<PHASE>","lat":..,"lon":..,"dist_m":..}'`.
5. **Acknowledge on the channel the word came from.** One line: phase, position,
   distance and bearing. Example: `Activating. 41.019, -73.641. OBJ 374m N. Tracking.`
   If a step failed, say which and why in one line.

On HOME, add the summary: total time from activating, markers placed, any
failures. Then set the mission inactive.

## Heartbeat (optional, owner creates it)

A per-minute position line during an active mission is a Hermes cron job with
a no-agent script that stays silent when no mission is active. Proposing it is
fine; creating a recurring job and its delivery target is a standing change, so
ask first. Cron scripts must sit in the directory the cron runner reads
(`hermes cron --help` shows it; the predecessor used `/opt/data/scripts/`).

## Pre-mission readiness (read-only, all must pass)

1. Bus reachable: `nats sub 'tak.cot.out' --count 1 --timeout 15s` returns an event.
2. Owner GPS fresh (age under 60 s is good, under 120 s acceptable).
3. `takserver-messaging` and `takserver-plugins` pods Ready.
4. The owner's client shows connected in `contacts/all` or a recent event.
5. If voice replies are expected, `voice.auto_tts` and the gateway voice mode
   are on (both were seen reverting after restarts on the predecessor).

## Style

Maximum two sentences per message, no greetings or sign-offs, answer first.
Mission acknowledgements are one line. If replies drift verbose, restate this in
the runtime system prompt via `hermes config` (read the current value first).

## Pitfalls

- Skills do not fire on a keyword alone on a chat channel; the instruction to
  honour phase words must also be in the soul or system prompt.
- Never put the owner's live coordinates in a skill, memory entry or Git.
- Real coordinates and device UIDs belong to the mission, not the procedure.
