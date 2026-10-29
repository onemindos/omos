---
status: selected-first-slice-implemented
user-confirmation: owner direction, 2026-10-29 ("commo is the key"; "make vrtak aka onemind")
---

# OneMind speaks native TAK: Commo-level mesh at the edge, TAK-kernel clients on top

## Decision

OneMind builds on the shared core that native TAK clients use, not only on
TAK Server feeds:

1. **Edge nodes join the TAK mesh as peers.** Robots, drones, sensors and Hermes
   reach the map through a mesh gateway that behaves like Common Commo: UDP
   multicast SA, TCP directed messages, and TAK Protocol negotiation (XML v0,
   protobuf v1). Mesh traffic keeps working when the cloud and TAK Server are
   unreachable.
2. **The bus stays the authority.** Everything the gateway hears goes to the
   in-cluster NATS through the edge leaf node. Hermes and the entity layer read
   it there, and send to the mesh through one subject.
3. **The OneMind clients are the TAK-kernel clients.** ATAK and iTAK on mobile,
   TAKX on desktop (the Java build of the same TAK Kernel), and a OneMind build
   of VR-TAK for XR. CloudTAK remains the browser view.

## Evidence (read from the public repositories, 2026-10-29)

| Component | Source | What it is |
|---|---|---|
| Common Commo | `TAK-Product-Center/atak-civ/commoncommo` | C++ comms kernel with JNI, Objective-C and Win32 bindings: mesh SA, streaming, protocol negotiation, mission packages, enrollment, QUIC |
| TAK Kernel / Engine | `atak-civ/takkernel` | C++ map engine plus `gov.tak.api` Java models; build target `android` (ATAK) or `java` (TAKX); renders Cesium 3D Tiles through vendored cesium-native |
| VR-TAK | `TAK-Product-Center/VR` | Unreal 5.3.2, Windows; ships `commoncommo.dll` in its `VH_COT` plugin and uses its own `VH_Globe`, not the TAK engine |
| Wire spec | `atak-civ/takproto`, `atak-civ/takcot` | TAK Protocol framing and `.proto` field numbers; CoT XSDs |
| CloudTAK | `dfpc-coe/CloudTAK` | TypeScript + MapLibre, `@tak-ps/node-cot`; a TAK Server streaming client, not a mesh peer |

Licences: ATAK-CIV and the VR SDK are GPL-3.0 (federal-employee portions public
domain). Building `takkernel` needs TAK artifact-repository credentials
(`takRepoMavenUrl`, `takRepoConanUrl`), so in practice tak.gov access.

## First slice (this change)

`services/bridge/bridges/tak` becomes the **TAK mesh gateway** (v0.2.0):

- joins the SA group (default `239.2.3.1:6969`), accepts directed TCP on 4242,
  and optionally unicast UDP;
- decodes and encodes TAK Protocol v1 (`0xbf 0x01 0xbf` + `TakMessage`), applies
  the detail.proto whole-element rule, and follows mesh negotiation rules 1–5
  (TakControl every 60 s, per-contact versions, fallback to XML when any contact
  is XML-only, immediate TakControl on change);
- publishes CoT XML to `tak.cot.mesh.rx` with `Nats-Msg-Id`, so two gateways
  hearing the same report store it once in stream `COT`;
- transmits CoT XML from `tak.cot.mesh.tx` (request/reply returns `ok`);
- optionally copies mesh CoT to `tak.cot.inject` (`MESH_TO_TAK_SERVER=true`),
  so mesh devices also appear on TAK Server, CloudTAK and server-connected VR-TAK;
- optionally announces itself as a mesh contact with a TCP endpoint, so a TAK
  user can message the gateway (and through it, Hermes);
- reports health on `tel.fleet.edge.{id}.status`.

The codec is written from the wire spec with `protowire`; no GPL source is
copied, so the bridge stays Apache-2.0. It was cross-checked both ways against
the `takproto` Python package (Google protobuf runtime).

The old `omos.field.*` and `omos.bridge.*` subjects are removed from this bridge:
they were outside the locked grammar and nothing in the platform consumed them.
All subjects used are under the `tak.cot.>` and `tel.` patterns already in the
grammar, and land in the existing `COT` and `TEL` streams.

### Second slice (v0.3.0, feature A7)

- `MESH_INTERFACE` applies to transmit as well as receive (`IP_MULTICAST_IF`);
  before, SA left by the default route, which on a Mac is often not the ATAK LAN.
- GeoChat: joins `224.10.10.1:17012`; chat travels `tak.cot.mesh.rx` and
  `tak.cot.mesh.tx` like any CoT (`Tak-Delivery: chat` on receive).
- Commo-compatible mesh encryption (AES-256-CBC + HMAC-SHA256 over IV and
  ciphertext, keys from the environment), checked against OpenSSL.
- Peer-hosted data packages: directed `b-f-t-r` requests are fetched,
  verified and acked; an outbox is served at `/getfile`.
- Directed transmit to one contact with the `Tak-Dest-Uid` header.
- With `MESH_TO_TAK_SERVER=true`, OneMind's own CoT (`OM.*`, `INTEL.*`) and
  gateway presence are not copied to `tak.cot.inject` (contract proposal CP-2).
- A TakControl with unset versions counts as version 1, as `takcontrol.proto`
  says, so GeoChat from a `takproto` sender no longer drops the mesh to XML.

Open questions for the owner are in
`docs/roadmap/2026-10-29-tak-mesh-gateway-owner-proposals.md`.

Not deployed. The gateway runs on an edge host (Mac, home server, Jetson) next
to a NATS leaf, not in the cloud cluster: multicast does not cross the DOKS
network.

## Next slices (not started)

1. **Make Sense.** An entity service that reads `tak.cot.mesh.rx`, `tak.cot.out`
   and `tel.fleet.>`, and publishes `ent.{type}.{id}.state`, which
   `om-tak-projector` already consumes. This is the layer Hermes reasons over.
2. **Fleet in CoT.** MAVLink and other bridges publish `tel.fleet.{class}.{id}.position`;
   a projector turns fleet positions into CoT on `tak.cot.mesh.tx` (local) and
   `tak.cot.inject` (server), so drones and robots show on every client.
3. **OneMind XR (VR-TAK build).** Fork `TAK-Product-Center/VR` into a OneMind
   repository under GPL-3.0. It already speaks Commo, so it sees the mesh and
   TAK Server unchanged. OneMind work is Unreal plugins beside `VH_COT`: a NATS
   panel for Hermes, an entity layer, and Cesium for Unreal for the homestead
   3D tiles from WebODM. Needs a Windows build host with Unreal 5.3.2.
4. **Desktop and web.** TAKX when the civilian build ships (same kernel as ATAK,
   so the same plugin model); CloudTAK for browsers. Unreal Pixel Streaming of
   the XR client to a browser is a possible later option; it is untested for VR-TAK.

## Open decisions for the owner

- Whether `MESH_TO_TAK_SERVER` should be on by default at each site.
- Whether to track `evt.tak.cot.{uid}` per NATS-DOCTRINE §9.8 option 1, which the
  entity service would then publish.
- Whether the OneMind XR fork is public (GPL-3.0 requires source for distributed
  binaries either way).
