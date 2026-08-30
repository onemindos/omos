# om-overwatch-tak

**Organ 3/7 of om-overwatch** — the C2 bridge between the fusion spine and TAK/ATAK.

The TAK Server's Java **NatsBridgePlugin** already owns all TAK transport (TLS :8089,
certs, the TAK protocol) and exposes CoT as plain XML on NATS. This organ is a pure
**translation + fusion layer** — no TAK server comms, no certs, no protobuf.

```
INBOUND   tak.cot.out ──parse CoT──► det.cot.{id} ──► spine fuses
          Human-tagged affiliation (a-f/h/n/u) becomes AUTHORITATIVE in the spine.

OUTBOUND  ent.air.{id}.state ──render CoT──► tak.cot.inject ──► every ATAK client
          Operators finally see the fused sensor picture (drone detections,
          friendly UAV telemetry, etc.) on their tablets.
```

## Affiliation authority (human-in-the-loop)

When an operator tags a track friendly/hostile/neutral in ATAK, the CoT flows in on
`det.cot.{id}` and the spine **locks** that affiliation (`affiliation_locked`). Machine
sensors (RF/acoustic/visual/telemetry) may still update position, confidence, and add
their source, but they **cannot flip the operator's friend/foe call**. Proven:

```
CoT a-h "Hostile-Recon"  → affiliation=hostile locked=true sources=[cot]
visual sensor claims friendly → affiliation STILL hostile, sources=[cot,visual]
```

## Loop prevention

- Outbound never echoes a track whose sources include `cot` (ATAK already has it).
- Inbound ignores CoT with uid prefix `OW.` or an `<__overwatch fused="true"/>` marker
  (our own re-broadcasts).

## Subjects

| Direction | Subject | Notes |
|---|---|---|
| IN  | `tak.cot.out` | all CoT the TAK server sees (devices, operators) |
| IN→spine | `det.cot.{id}` | contact report with human affiliation, `source:"cot"` |
| spine→OUT | `ent.air.*.state` | fused entity state |
| OUT | `tak.cot.inject` | CoT pushed to all ATAK clients (uid `OW.{id}`) |

## CoT ↔ entity mapping

- **Affiliation:** CoT 2nd token `f/h/n/u` ↔ friendly/hostile/neutral/unknown
  (j/k/s → hostile; p/a → unknown).
- **Type:** CoT 3rd token `A/G/S` ↔ air/ground/sea; ground unit (`-U-`) → operator_station,
  ground equipment sensor (`-E-S`) → sensor_platform. Reverse map builds MIL-STD-2525
  CoT types (e.g. friendly fixed-wing → `a-f-A-M-F`, rotary → `a-{aff}-A-M-H-Q`).

## Build & run

```bash
go build -o om-overwatch-tak .
NATS_URL=nats://127.0.0.1:4222 ./om-overwatch-tak
```

Deployed under NEX via `Nexfile`; see `nex/nex-boot-deploy.sh`.

## Live ATAK test — plug-and-play

Everything server-side is wired (TAK Server :8089, NatsBridgePlugin connected,
this organ running). To confirm on a real device:

1. **Connect** an ATAK (Android) / iTAK (iPhone) to `192.168.1.183:8089` (LAN) or
   `100.74.182.117:8089` (Tailscale), SSL/TLS, cert `admin.p12` (pw `atakatak`).
2. **Outbound test** — run the live feed; moving tracks appear on the map:
   ```bash
   python3 cmd/demo-feed.py          # DEMO-PATROL (friendly, blue) + DEMO-THREAT (hostile, red)
   # CENTER_LAT / CENTER_LON to recenter on your AO. Ctrl-C to stop.
   ```
   Drives REAL fusion (spine → ent.air → this organ → tak.cot.inject → ATAK).
3. **Inbound test** — drop/tag a hostile marker on the device, then:
   ```bash
   nats kv get OVERWATCH_STATE <uid>  # expect affiliation=hostile, affiliation_locked=true, sources=[cot]
   ```
Cleanup: `nats kv del OVERWATCH_STATE DEMO-PATROL DEMO-THREAT`.

## Provenance

`bridge.py.ref` — the original `cot2constellation` Python bridge, kept as the reference
spec for the CoT parser and affiliation mapping. This Go organ supersedes it: adds the
outbound direction, drops the legacy :4299/org_id/REST transport, and integrates with the
NatsBridgePlugin.
