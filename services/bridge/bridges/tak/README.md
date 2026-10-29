# tak-mesh — OneMind TAK mesh gateway

Part of the **One Mind OS (omos)** bridge ecosystem. Design and context:
[ADR 0003](../../../../docs/architecture/decisions/0003-tak-native-clients-and-mesh.md).

## What It Does

ATAK, iTAK, WinTAK and VR-TAK talk to each other through Common Commo, TAK's
comms library: SA over UDP multicast, direct messages over short TCP
connections, and TAK Protocol negotiation between XML (v0) and protobuf (v1).
This gateway does the same, so OneMind is a peer on the mesh rather than only a
TAK Server feed. It works with no server and no cloud link; the NATS leaf node
forwards to the cluster bus whenever it is connected.

| Direction | Path |
|---|---|
| Mesh → bus | SA group, directed TCP and optional UDP → CoT XML on `tak.cot.mesh.rx` |
| Bus → mesh | CoT XML on `tak.cot.mesh.tx` → multicast in the negotiated protocol version |
| Mesh → TAK Server (optional) | `MESH_TO_TAK_SERVER=true` also publishes to `tak.cot.inject` |
| Health | `tel.fleet.edge.{GATEWAY_ID}.status` every 30 s |

Protocol handling follows `takproto/README.txt` in `TAK-Product-Center/atak-civ`:

- Receives both v0 XML and v1 protobuf; always publishes XML to NATS, so
  consumers see the same format as `tak.cot.out` from TAK Server.
- Broadcasts TakControl every 60 s and tracks each contact's versions; sends v1
  only when every contact it hears can decode it, otherwise XML.
- On encode, typed detail fields (`contact`, `__group`, `precisionlocation`,
  `status`, `takv`, `track`) are used only when the whole element maps;
  everything else stays in `xmlDetail`.
- Drops its own multicast echo, and sets `Nats-Msg-Id` so two gateways hearing
  the same report store it once in stream `COT` (2 min duplicate window).
- Registered detail extensions it cannot decode are counted in the
  `Tak-Undecoded-Extensions` header.

The codec is written from the published wire spec; no GPL code is included.

## Quick Start

```bash
cd services/bridge/bridges/tak
cp .env.example .env
go build -o tak .
set -a; source .env; set +a; ./tak
```

Run it on an edge host on the same LAN or radio network as the TAK devices, next
to a NATS leaf. Multicast does not reach pods in the cloud cluster.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS (normally the local leaf) |
| `NATS_USER`, `NATS_PASSWORD` | unset | NATS credentials, if the leaf requires them |
| `GATEWAY_ID` | `tak-mesh` | Gateway name in subjects, headers and TakControl (`omos-gw-<id>`) |
| `MESH_GROUP` | `239.2.3.1:6969` | SA multicast group |
| `MESH_INTERFACE` | system default | Interface to join the group on |
| `MESH_PROTOCOL` | `auto` | `auto` negotiates; `xml` never sends protobuf |
| `DIRECT_TCP_PORT` | `4242` | Directed TCP messages; `0` disables |
| `UDP_LISTEN_PORT` | `0` | Unicast UDP CoT input; `0` disables |
| `MESH_TO_TAK_SERVER` | `false` | Also publish received CoT to `tak.cot.inject` |
| `GATEWAY_CALLSIGN` | unset | If set, the gateway sends its own SA every 30 s |
| `GATEWAY_ENDPOINT_HOST` | unset | LAN address advertised as `host:DIRECT_TCP_PORT:tcp`, so TAK users can message the gateway |
| `GATEWAY_COT_TYPE` | `a-f-G-E-S` | CoT type of the gateway's SA |
| `GATEWAY_LAT`, `GATEWAY_LON` | `0` | Position of the gateway's SA |

## NATS Subjects

| Subject | Direction | Payload |
|---|---|---|
| `tak.cot.mesh.rx` | Publish | CoT XML. Headers: `Nats-Msg-Id`, `Tak-Gateway`, `Tak-Delivery` (`mesh`, `direct`, `udp`), `Tak-Proto` (`0`, `1`), `Tak-From` |
| `tak.cot.mesh.tx` | Subscribe | CoT XML to broadcast. Request/reply answers `ok` or `error: …` |
| `tak.cot.inject` | Publish (optional) | CoT XML for TAK Server via NatsBridgePlugin |
| `tel.fleet.edge.{id}.status` | Publish | `{uid, type: "fleet.edge", t, meta: {proto_version, contacts, rx, tx, errors, last_error, …}}` |

Anyone allowed to publish `tak.cot.mesh.tx` can put markers on every device on
the mesh. Restrict it with NATS permissions.

```bash
nats sub tak.cot.mesh.rx --headers-only
nats req tak.cot.mesh.tx '<event version="2.0" uid="omos-test-1" type="b-m-p-s-m" how="h-g-i-g-o" time="2026-10-29T12:00:00Z" start="2026-10-29T12:00:00Z" stale="2026-10-29T13:00:00Z"><point lat="30.2672" lon="-97.7431" hae="0" ce="9999999" le="9999999"/><detail><contact callsign="OneMind test"/></detail></event>'
nats sub 'tel.fleet.edge.*.status'
```

## Tests

```bash
go test ./...
```

Covers the protobuf round trip of a full ATAK SA, the whole-element rule,
negotiation rules 3a–3c and 4, malformed input, real multicast on the host
(skipped where multicast is unavailable) and directed TCP.

## Changes from v0.1

v0.1 accepted XML on unicast UDP 8087, published JSON envelopes to
`omos.field.<uid>.telemetry` and sent `omos.field.*.command` to a TAK server as
CoT at 0,0. Those subjects were outside the platform grammar and had no
consumers; they are removed. Set `UDP_LISTEN_PORT=8087` to keep the old input.

## License

Apache-2.0.
