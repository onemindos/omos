---
name: mission-replay
description: "Use to replay a recorded TAK mission at compressed speed: read the track and markers back from the COT stream and re-publish them as REPLAY-* CoT on tak.cot.inject for live review on the map."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [tak, replay, mission, review, cot, nats]
    related_skills: [tak-mission-ops, tak-mission-runner, tak-server, nats-operations]
triggers:
  - "replay mission"
  - "replay the mission"
  - "mission replay"
  - "play back"
  - "rewind mission"
---

# Mission replay

Re-inject recorded CoT at compressed speed so the owner can watch a mission again
on the map. Rewritten 2026-10-29: the earlier version read the TAK database with
an embedded login and wrote to port 8089 with an admin certificate and a
hard-coded personal UID. Both are replaced by the bus: **read history from the
JetStream stream `COT`, publish the replay to `tak.cot.inject`**, with every UID
prefixed `REPLAY-` so a replayed item is never mistaken for a live one.

## Before you start

Ask in one message: which mission or time window, which UIDs (the owner's device
and any `MISSION-<name>-*` markers; list candidates from the stream first), and
the speed. Default 30x; offer 10x for detail, 60x for a quick scan. Do not use
100x for review: a 30-second track becomes 0.3 s jumps.

| Speed | 4 h mission replays in | Use |
|---|---|---|
| 10x | about 24 min | detailed review |
| 30x | about 8 min | standard (default) |
| 60x | about 4 min | quick scan |
| 100x | about 2.4 min | confirmation only |

The stream keeps 30 days (`maxAge: 720h`), so the window is limited to that.
Older material is not available from the bus.

## Read the history (no credentials needed beyond bus access)

```bash
# What is in the window? (adjust the times; UTC)
nats stream info COT
# Ephemeral pull consumer: it disappears by itself, so nothing is left to delete.
nats consumer add COT --ephemeral --pull --filter 'tak.cot.out' --deliver time \
  --start-time "2026-10-29T12:00:00Z" --ack none --defaults
# Then fetch with `nats consumer next COT <name> --count N` and save to a scratch file.
# Flags vary by CLI version: check `nats consumer add --help` first.
```

Filter events by `uid` (the owner's device UID from the live contact list, and
the `MISSION-<name>-*` / `LEGACY-*` markers). Drop events at latitude 0. If the
bus identity cannot create a consumer, ask the owner for a read-only consumer or
for a one-off export; a read-only database login from a SealedSecret is the
fallback, never an admin login.

## Replay

For each kept event in time order:

1. Rewrite `uid` to `REPLAY-<first 20 chars of the original uid>`.
2. Set `time` and `start` to the original timestamp shifted so the first event
   is "now", and `stale` to 5 minutes later (so replayed tracks linger briefly
   after the end and then clear by themselves; shorten to 30 s to clear sooner).
3. Keep `type`, point and `detail`; set `how="m-g"`.
4. `nats pub tak.cot.inject "$(cat event.xml)"` (never `< file`), then sleep
   `(t[i+1]-t[i]) / speed` seconds.

A short Python loop is fine; keep it in memory or a scratch file, not in the
skill. Publish one event per message. Do not hold an extra TLS connection to
TAK; the plugin on the server side does the injecting.

Stage proofs: the first replayed event appears on `tak.cot.out` (TAK echoes it
to subscribers) and the owner sees a `REPLAY-*` item. Tell the owner to open
the map before you start.

## What to expect on the map

- The track animates as a moving marker; phase markers (SP, MOVING, EYES ON, AT
  OBJ, PACKAGE, RTB, HOME) appear at their original positions.
- Injection is server-side: TAK broadcasts to all connected clients in the same
  group. If one device does not show it, that client is not connected or is in a
  different group; the injection path is not at fault.
- Replay items expire by themselves after their stale time.

## Pitfalls

- Events must be valid CoT with a future `stale`; an expired `stale` is dropped.
- Use the original owner of the data in mind: a replay of another person's track
  is sensitive. Only replay what the owner asked for.
- Do not edit the original timestamps in the stream; replay builds new events.
- To replay beyond the stream's 30 days, a database export is needed; ask.
