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
| Mesh → bus | SA group, GeoChat group, directed TCP and optional UDP → CoT XML on `tak.cot.mesh.rx` |
| Bus → mesh | CoT XML on `tak.cot.mesh.tx` → multicast in the negotiated protocol version (GeoChat `b-t-f*` to the chat group), or directed TCP with header `Tak-Dest-Uid` |
| Mesh → TAK Server (optional) | `MESH_TO_TAK_SERVER=true` (default `false`) also publishes to `tak.cot.inject`, except OneMind's own CoT |
| Data packages | Directed `b-f-t-r` requests are fetched into `DATAPKG_DIR/inbox` and acked; `DATAPKG_DIR/outbox` is served at `/getfile` |
| Health | `tel.fleet.edge.{GATEWAY_ID}.status` every 30 s |

Protocol handling follows `takproto/README.txt` and the Commo sources in
`TAK-Product-Center/atak-civ`:

- Receives both v0 XML and v1 protobuf; always publishes XML to NATS, so
  consumers see the same format as `tak.cot.out` from TAK Server.
- Broadcasts TakControl every 60 s and tracks each contact's versions; sends v1
  only when every contact it hears can decode it, otherwise XML. A TakControl
  with unset versions counts as version 1, and one without `contactUid` applies
  to the paired event's uid (`takcontrol.proto`).
- On encode, typed detail fields (`contact`, `__group`, `precisionlocation`,
  `status`, `takv`, `track`) are used only when the whole element maps;
  everything else stays in `xmlDetail`.
- Drops its own multicast echo, and sets `Nats-Msg-Id` so two gateways hearing
  the same report store it once in stream `COT` (2 min duplicate window).
- Registered detail extensions it cannot decode are counted in the
  `Tak-Undecoded-Extensions` header.

The codec and the mesh encryption are written from the published wire spec and
checked against OpenSSL and the `takproto` Python package; no GPL code is
included.

### Interfaces

`MESH_INTERFACE` applies to both directions: the gateway joins the groups on
that interface, and its transmit socket sends multicast there
(`IP_MULTICAST_IF`). Leave it unset only on a host with one network. On a Mac
with Wi-Fi, Ethernet and Tailscale, set it to the interface on the ATAK LAN
(`route -n get 239.2.3.1` shows which one the system would pick otherwise).
The status message reports the interface in `meta.tx_interface`.

**Run one gateway per mesh interface and NATS subject.** Every gateway
subscribed to `tak.cot.mesh.tx` transmits every message it gets, so two gateways
on the same LAN and bus send each marker twice. Two gateways on the same LAN
that both receive are fine: `Nats-Msg-Id` deduplicates what they publish.

### GeoChat

ATAK's "All Chat Rooms" uses multicast group `224.10.10.1:17012`. The gateway
joins it (`MESH_CHAT_GROUP`, `off` disables). Chat it hears is published on
`tak.cot.mesh.rx` with `Tak-Delivery: chat`. CoT of type `b-t-f` or `b-t-f-*`
on `tak.cot.mesh.tx` is sent to the chat group; everything else goes to the SA
group.

### Mesh encryption

With `MESH_AUTH_KEY` and `MESH_CIPHER_KEY` set (two different 256-bit keys, 64
hex characters each), the gateway does what Commo does after `setCryptoKeys`:
every multicast and directed message is AES-256-CBC encrypted with a random IV,
followed by an HMAC-SHA256 under the auth key over IV and ciphertext (wire
layout `IV | ciphertext | HMAC`). Inbound mesh, directed and UDP messages that
do not authenticate are dropped and counted as errors, including plaintext.

The keys are secrets. Put them in the host's secret store and export them into
the gateway's environment; never write them into `.env.example`, `bridge.yaml`
or anything committed. Every device on the mesh must use the same pair.

### Data packages

With `DATAPKG_DIR` set, the gateway takes part in Commo's peer-hosted transfer
(`commoncommo/core/missionpackage-rules.txt`):

- **Receive.** A file transfer request (`b-f-t-r` with `<fileshare>`) sent
  **directly to the gateway** (directed TCP) is fetched over http from the
  request's source address (the URL's port and path are kept; the host is
  replaced, as Commo's own fallback does), up to `DATAPKG_MAX_BYTES`. The file
  is kept only if its size and SHA-256 match the request, as
  `inbox/<sha256>-<filename>`. If the sender asked for an ack, the gateway
  sends a `b-f-t-a` `<ackresponse>` with `success` true or false to the
  sender's endpoint. Requests heard on multicast are not fetched. Peer-hosted
  https (`peerHosted`, `httpsPort`) and server-hosted (`/Marti/`) URLs are
  refused.
- **Serve.** With `DATAPKG_HTTP_PORT` set, files placed in `outbox/` are served
  at `http://<host>:<port>/getfile?file=<name>`. To send one to an ATAK user,
  publish a `b-f-t-r` naming that URL on `tak.cot.mesh.tx` with
  `Tak-Dest-Uid: <their uid>`.

Received files stay on the edge host; nothing about them is published to NATS
yet (see the owner proposals in
`docs/roadmap/2026-10-29-tak-mesh-gateway-owner-proposals.md`).

## Quick Start

On the Mac, next to the local leaf (`infra/edge/mac-leaf`):

```bash
cd services/bridge/bridges/tak
cp .env.example .env          # edit; secrets come from the environment, not this file
MESH_INTERFACE=en0 ./run-mac.sh
```

`run-mac.sh` builds the gateway, checks that the leaf answers on `NATS_URL`,
and prints the transmit interface before starting. Anywhere else:

```bash
go build -o tak . && set -a && . ./.env && set +a && ./tak
```

Run it on an edge host on the same LAN or radio network as the TAK devices, next
to a NATS leaf. Multicast does not reach pods in the cloud cluster, so the
gateway has no container image or Argo Application.

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `NATS_URL` | `nats://localhost:4222` | NATS (normally the local leaf) |
| `NATS_USER`, `NATS_PASSWORD` | unset | NATS credentials, if the leaf requires them |
| `GATEWAY_ID` | `tak-mesh` | Gateway name in subjects, headers and TakControl (`omos-gw-<id>`) |
| `MESH_GROUP` | `239.2.3.1:6969` | SA multicast group |
| `MESH_CHAT_GROUP` | `224.10.10.1:17012` | GeoChat group; `off` disables |
| `MESH_INTERFACE` | system default | Interface to join the groups on **and transmit on** |
| `MESH_TTL` | `1` | Multicast TTL (Commo's default is 1) |
| `MESH_PROTOCOL` | `auto` | `auto` negotiates; `xml` never sends protobuf |
| `MESH_AUTH_KEY`, `MESH_CIPHER_KEY` | unset | Mesh encryption keys, 64 hex characters each, different; both or neither. Secrets |
| `DIRECT_TCP_PORT` | `4242` | Directed TCP messages; `0` disables |
| `UDP_LISTEN_PORT` | `0` | Unicast UDP CoT input; `0` disables |
| `MESH_TO_TAK_SERVER` | `false` | Also publish received CoT to `tak.cot.inject`; skips uids `OM.*`, `INTEL.*` and `omos-gw-*` |
| `GATEWAY_CALLSIGN` | unset | If set, the gateway sends its own SA every 30 s |
| `GATEWAY_ENDPOINT_HOST` | unset | LAN address advertised as `host:DIRECT_TCP_PORT:tcp`, so TAK users can message the gateway |
| `GATEWAY_COT_TYPE` | `a-f-G-E-S` | CoT type of the gateway's SA |
| `GATEWAY_LAT`, `GATEWAY_LON` | `0` | Position of the gateway's SA |
| `DATAPKG_DIR` | unset | Enables data packages; holds `inbox/` and `outbox/` |
| `DATAPKG_HTTP_PORT` | `0` | Serves `outbox/` at `/getfile`; `0` disables |
| `DATAPKG_MAX_BYTES` | `52428800` | Largest data package accepted |

## NATS Subjects

| Subject | Direction | Payload |
|---|---|---|
| `tak.cot.mesh.rx` | Publish | CoT XML. Headers: `Nats-Msg-Id`, `Tak-Gateway`, `Tak-Delivery` (`mesh`, `chat`, `direct`, `udp`), `Tak-Proto` (`0`, `1`), `Tak-From` |
| `tak.cot.mesh.tx` | Subscribe | CoT XML to send. Optional header `Tak-Dest-Uid` sends it directed to that contact. Request/reply answers `ok` or `error: …`; senders that would use a JetStream publish should use a core publish instead (the reply would race the PubAck) and read `meta.tx` from the status |
| `tak.cot.inject` | Publish (optional) | CoT XML for TAK Server via NatsBridgePlugin |
| `tel.fleet.edge.{id}.status` | Publish | `{uid, type: "fleet.edge", t, meta: {proto_version, contacts, rx, tx, errors, last_error, tx_interface, chat_group, encrypted, to_tak_server, datapkg_rx, datapkg_failed, …}}` |

Anyone allowed to publish `tak.cot.mesh.tx` can put markers on every device on
the mesh, and anyone allowed to publish `tak.cot.inject` can do the same through
TAK Server. Restrict both with NATS permissions.

```bash
nats sub tak.cot.mesh.rx --headers-only
nats req tak.cot.mesh.tx '<event version="2.0" uid="omos-test-1" type="b-m-p-s-m" how="h-g-i-g-o" time="2026-10-29T12:00:00Z" start="2026-10-29T12:00:00Z" stale="2026-10-29T13:00:00Z"><point lat="30.2672" lon="-97.7431" hae="0" ce="9999999" le="9999999"/><detail><contact callsign="OneMind test"/></detail></event>'
nats sub 'tel.fleet.edge.*.status'
```

## Tests

```bash
go vet ./... && go test -race ./...
```

Covers the protobuf round trip of a full ATAK SA, the whole-element rule,
negotiation rules 3a–3c and 4 (including TakControl with unset versions),
malformed input, real multicast on the host, transmit on a named interface
(the loopback interface: a peer joined there must hear the gateway), GeoChat
routing, mesh encryption against an OpenSSL known-answer vector, plaintext
rejection on an encrypted mesh, the `MESH_TO_TAK_SERVER` filter, directed TCP
both ways, and data package fetch, verify, ack and `/getfile` serving.

## Changes

**v0.3.0.** `MESH_INTERFACE` now applies to transmit (it was receive only, so
SA could leave by the default route instead of the ATAK LAN). GeoChat group,
Commo-compatible mesh encryption, peer-hosted data packages, directed transmit
(`Tak-Dest-Uid`), `MESH_TTL`. `MESH_TO_TAK_SERVER=true` no longer copies
OneMind's own CoT (`OM.*`, `INTEL.*`) or gateway presence to `tak.cot.inject`.
A TakControl with unset versions (as `takproto` attaches to GeoChat) no longer
drops the whole mesh to XML.

**v0.2.0.** v0.1 accepted XML on unicast UDP 8087, published JSON envelopes to
`omos.field.<uid>.telemetry` and sent `omos.field.*.command` to a TAK server as
CoT at 0,0. Those subjects were outside the platform grammar and had no
consumers; they are removed. Set `UDP_LISTEN_PORT=8087` to keep the old input.

## License

Apache-2.0.
