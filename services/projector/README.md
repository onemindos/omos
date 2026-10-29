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
| serves | `svc.cot.send` (queue group `omos-tak-projector`) | client-originated CoT, below |
| writes | `tak.cot.inject` (`PROJECTOR_TARGET=server`) or `tak.cot.mesh.tx` (`mesh`) | core publish, no reply subject; headers `Nats-Msg-Id: {uid}:{time_ms}`, `OneMind-Marking`, `OneMind-Releasable-To`, `OneMind-Projector`; client sends add `OneMind-Origin: svc.cot.send`, `OneMind-Sender: {account}/{user}` and, for a directed send, `Tak-Dest-Uid` |

## svc.cot.send (client-originated CoT)

The projector is the single writer of its output subject, so client CoT (GeoChat, drawings,
routes, markers, SPI) goes through it and through the same publisher and ceiling as entities
(cot-projection.md §1 and §9; register entry INT-R4 in `docs/proposals/2026-10-29-grammar-register/`,
PR #556). Two bodies are accepted:

| Form | Request | Reply |
|---|---|---|
| Structured GeoChat (L0 web, WEB-P4) | `{v: 1, kind: "geochat", target: server\|mesh\|both, message_uid, room {id, name, kind: all\|group\|direct}, text (≤ 1000), sender {uid, callsign}, marking?}`; the projector builds the `b-t-f` CoT | `{v, ok: true, uid, targets[, skipped]}` or `{v, ok: false, error: {code, message}}` |
| CoT (Unreal client, B-P4 and UE-P1) | `{v: 1, kind: marker\|spi\|route\|geofence\|shape\|chat, targets[], client_id, client_ref, marking {tlp, releasable_to?}, dest_uid?, cot}` | `{v, ok: true, uid, targets[, skipped]}` or `{v, ok: false, error: "<code>", detail}` |

Codes: `invalid`, `forbidden`, `rate_limited`, `target_failed`. Checks, in order:

1. **Caller.** Only a user of an account in `COT_SEND_ACCOUNTS` (`CLIENT`), read from the
   `Nats-Request-Info` header that nats-server sets on the `share: true` import (omos #562).
   A caller in the projector's own account with no header, or with a header naming the own
   account, is refused (`forbidden`). **That header is not proof of identity by itself:**
   nats-server passes a header that an own-account publisher wrote by hand unchanged, so any
   other `ONEMIND` user (the shared `onemind` and `leaf_tak` logins) can claim
   `{"acc":"CLIENT","user":"zeus"}` and the projector cannot tell
   (`TestSendForgedRequestInfo`, "without bus control"). The identity is only trustworthy when
   no own-account user except the projector may publish to `svc.cot.>` (owner proposal A6-P13).
   Until that deny is applied the endpoint is **off**: `COT_SEND_ENABLED` defaults to `false`.
2. **Replay.** A repeated `client_ref` / `message_uid` from the same caller gets the stored reply
   and is not sent again (`COT_SEND` dedup window 10 min). The key is reserved before the send
   is queued, so a retry while the first is in flight gets `target_failed` ("in progress") and
   never publishes twice; once queued, the request waits for the loop's answer (no timeout).
3. **Rate.** Token bucket per caller (`COT_SEND_RATE_PER_MIN`, `COT_SEND_BURST`).
4. **Content.** Reserved uids: any attribute naming `OM.*` or `INTEL.*` (also after a
   `BAO.F.*.` source prefix) is `forbidden`. The uid is bound to the sender: `onemind-{client_id}-…`
   (`GeoChat.onemind-{client_id}.…` for chat), so a client cannot take over a TAK user's marker;
   a GeoChat `sender.uid` must be the calling NATS user or end in `-{user}` (web `OM-WEB-{user}`).
   The kind fixes the CoT type (marker `a-*` or `b-m-p-s-m`, spi `b-m-p-s-p-i`, route `b-m-r`,
   geofence `u-d-f|u-d-c-c|u-d-r`, shape adds `b-m-p-s-m|u-r-b-c-c|u-rb-a`, chat `b-t-f`);
   civilian scope (no hostile, suspect, weapon or tasking type); `how` must be human (`h-*`).
   The XML is parsed strictly (one `<event>`, no DOCTYPE, comments or namespaces, depth ≤ 4) and
   `<detail>` children are allowlisted. Only ephemeral kinds (`spi`, `chat`, structured
   GeoChat) are served by default; the persistent kinds `marker`, `route`, `geofence` and
   `shape` belong to `entity.create` (cot-projection.md section 9) and need
   `COT_SEND_PERSISTENT_KINDS=true` (A6-P6). **Person-track data** is `forbidden`: an EUD's own
   position report (`<takv>`, `<__group>`, `<uid>`, `<status>`, `<track>`, `<contact endpoint|phone|…>`)
   would become `ent.person.*` on the reverse path; no attribute value (both forms: contact
   callsign, link and label attributes, GeoChat sender, room and uid) and no non-chat remarks
   text may reference `ent.person.*` / `ent.station.*` or device identifiers (IMSI, IMEI, MAC
   address). The free text of a chat message is not screened (A6-P14). Projector-only detail (`<onemind>`,
   `<_flow-tags_>`, `<marti>`) and the `access` / `caveat` / `releasableTo` attributes are refused.
   Time within `±10 min`, stale ≤ `COT_SEND_MAX_STALE_SEC`, size ≤ `COT_SEND_MAX_BYTES`
   (mesh: `MESH_MAX_BYTES`), vertices ≤ `MAX_VERTICES`.
5. **Marking gate G7**, the same rule as for entities. A GeoChat without `marking` takes
   `COT_SEND_CHAT_DEFAULT_TLP` (`CLEAR`; empty refuses); CoT-form sends must carry one.
6. **Target.** Each instance publishes only its own target; other requested targets are listed in
   `skipped`, and a request for none of this instance's targets is `target_failed`.
7. **Kill switch G1** (`cmd.projector.disable` refuses sends) and the **ceilings**: at most
   `COT_SEND_MAX_PER_SEC` client publishes per second, counted in the shared
   `MAX_RATE_PER_SEC` window, so the target never sees more than the global ceiling.

Client CoT is published unchanged after these checks (the Unreal client's bytes); GeoChat XML is
built by the projector.

## Configuration (contract section 10)

`NATS_URL`, `NATS_USER`, `NATS_PASSWORD`; `PROJECTOR_TARGET` (`server`), `PROJECTOR_INSTANCE`
(`hub`), `TAK_INJECT_SUBJECT`, `MESH_TX_SUBJECT`, `MIN_INTERVAL_SEC` (5 / 15 mesh),
`OWN_MIN_INTERVAL_SEC` (1 / 5), `MAX_RATE_PER_SEC` (50 / 10), `MAX_PROJECTED` (500 / 100),
`TAK_RELEASE_GROUPS` (empty in code; the deployment sets `homestead`, A6-P15), `AOI_EXIT_GRACE_S` (60), `MESH_WORLD_FEEDS` (false),
`MESH_MAX_BYTES` (1200), `MAX_VERTICES` (200), `ENT_STREAM` (`ENT`), `ENT_MODE=core`,
`HTTP_ADDR` (`:8080`: `/healthz`, `/readyz`, `/status` with counters), `LOG_LEVEL`.
`TAK_GROUP` is deprecated and ignored (`<__group>` is no longer sent, section 5.4).

svc.cot.send: `COT_SEND_ENABLED` (`false`; set `true` only after A6-P13),
`COT_SEND_PERSISTENT_KINDS` (`false`), `COT_SEND_SUBJECT` (`svc.cot.send`),
`COT_SEND_ACCOUNTS` (`CLIENT`), `COT_SEND_OWN_ACCOUNT` (`ONEMIND`), `COT_SEND_RATE_PER_MIN` (30),
`COT_SEND_BURST` (10), `COT_SEND_MAX_PER_SEC` (10 / 2 mesh), `COT_SEND_MAX_BYTES` (65536 / mesh
`MESH_MAX_BYTES`), `COT_SEND_MAX_STALE_SEC` (86400), `COT_SEND_CHAT_STALE_SEC` (600),
`COT_SEND_CHAT_DEFAULT_TLP` (`CLEAR`), `COT_SEND_ALLOW_UNVERIFIED` (`false`; local tests only).

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

`internal/projector/cotsend_test.go` covers svc.cot.send with the clients' own shapes (Unreal
`CotSendRequest` and builders, M4 drawings, web `geochatRequest`, an ATAK EUD report) and, on a
disposable `nats-server` with the omos #562 account layout (`CLIENT` imports `svc.cot.send` with
`share: true`), verified identity, impersonation (including an own-account forged
`Nats-Request-Info`, with and without the A6-P13 deny), person data, marking, replay, per-caller and
global rate limits, the kill switch and the single-writer count.
