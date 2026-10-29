# omos TAK projector (FEATURES A6)

Turns fused entities on `ent.{entity_type}.{entity_id}.state` into Cursor-on-Target for TAK,
only when every gate of onemind-xr `docs/contracts/cot-projection.md` section 3 passes, and sends
a removal when an entity stops passing. Proposal and owner decisions:
`docs/proposals/2026-10-29-projector-aoi.md`. Deployment: `platform/omos/projector/` (manual
sync, `platform-omos-projector`).

| Direction | Subject or bucket | How |
|---|---|---|
| reads | `ENT` (`ent.*.*.state`) | ordered ephemeral consumer: last message per subject on start, then every change, including `Nats-Marker-Reason` delete markers; core subscription when `ENT` is absent |
| reads | KV `geofences`, KV `mission-state` | watch all keys; every change re-runs the gates for every entity |
| reads | `tak.cot.out`, `tak.cot.mesh.rx` | tasking guard only: logs and counts `t-k-*` / `t-e-*` (refused, SL1) and other `t-*`; never acts |
| reads | `cmd.projector.disable`, `cmd.projector.enable` | live legacy control subjects (contract CP-1) |
| writes | `tak.cot.inject` (`PROJECTOR_TARGET=server`) or `tak.cot.mesh.tx` (`mesh`) | core publish, no reply subject; headers `Nats-Msg-Id: {uid}:{time_ms}`, `OneMind-Marking`, `OneMind-Releasable-To`, `OneMind-Projector` |

## Configuration (contract section 10)

`NATS_URL`, `NATS_USER`, `NATS_PASSWORD`; `PROJECTOR_TARGET` (`server`), `PROJECTOR_INSTANCE`
(`hub`), `TAK_INJECT_SUBJECT`, `MESH_TX_SUBJECT`, `MIN_INTERVAL_SEC` (5 / 15 mesh),
`OWN_MIN_INTERVAL_SEC` (1 / 5), `MAX_RATE_PER_SEC` (50 / 10), `MAX_PROJECTED` (500 / 100),
`TAK_RELEASE_GROUPS` (empty), `AOI_EXIT_GRACE_S` (60), `MESH_WORLD_FEEDS` (false),
`MESH_MAX_BYTES` (1200), `MAX_VERTICES` (200), `ENT_STREAM` (`ENT`), `ENT_MODE=core`,
`HTTP_ADDR` (`:8080`: `/healthz`, `/readyz`, `/status` with counters), `LOG_LEVEL`.
`TAK_GROUP` is deprecated and ignored (`<__group>` is no longer sent, section 5.4).

## Tests

```
go vet ./...
NATS_SERVER_BIN=$(which nats-server) go test -race -count=1 ./...
```

Unit tests cover the type table (every row of section 4.1; with `COTTYPES_XML` pointing at MITRE
`CoTtypes.xml` each type is also checked for membership), the civilian check, each gate, AOI
geometry, the cap order, intervals, the ceiling and removals. The acceptance tests in
`internal/projector/acceptance_test.go` run contract section 12 tests 2 to 6 and the projector
half of test 8, plus JADC2-MAP stop-line tests 8 and 10, on a disposable `nats-server`. The echo
and tasking test builds and runs the entity service from `../entity` (skipped if it does not
build; set `ENTITY_SERVICE_BIN` to use a prebuilt one). `A6_EVIDENCE_LOG` appends a one-line
result per check to a file.
