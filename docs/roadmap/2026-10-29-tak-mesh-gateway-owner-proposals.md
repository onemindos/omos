# TAK mesh gateway: owner proposals (2026-10-29)

Context: the mesh gateway in `services/bridge/bridges/tak` (v0.3.0, ADR 0003,
feature A7). It uses only subjects that already exist under the locked grammar
(`tak.cot.mesh.rx`, `tak.cot.mesh.tx`, `tak.cot.inject`,
`tel.fleet.edge.{id}.status`). The items below are the gaps the owner needs to
decide. Nothing here is implemented beyond what the gateway README describes.

| # | Gap | Proposal | Until decided |
|---|---|---|---|
| P1 | The grammar line for `tak.cot.>` reads "CoT from TAK (cot-tap output)", but `tak.cot.mesh.tx` carries CoT **toward** TAK devices, and stream `COT` (`tak.cot.>`) stores both directions | Ratify `tak.cot.mesh.rx` (mesh to bus) and `tak.cot.mesh.tx` (bus to mesh) as the mesh pair, and change the grammar line in `CLAUDE.md` to "CoT to and from TAK". Same as CP-4 in the OneMind cot-projection contract | The gateway keeps the two subjects as they are on this branch |
| P2 | Received data packages stay on the edge host. There is no Object Store CR in `platform/nats/streams/` and no subject for "a data package arrived" | Add an Object Store `datapkg` (R3 at the hub, R1 at a site, size cap) and an event `evt.tak.datapkg.received` (pattern `evt.{domain}.{entity}.{event}`) with `{sha256, name, size, sender_uid, gateway}`. The gateway would put the file and then publish the event | Files are in `DATAPKG_DIR/inbox`; `meta.datapkg_rx` on the status counts them |
| P3 | Two message headers are new: `Tak-Delivery: chat` (on `tak.cot.mesh.rx`) and `Tak-Dest-Uid` (on `tak.cot.mesh.tx`, directed send to one contact) | Ratify both as part of the gateway's header set (with `Tak-Gateway`, `Tak-Proto`, `Tak-From`). They are headers, not subjects, so the grammar does not change | In use as documented in the README and `docs/NATS-DOCTRINE.md` §9 |
| P4 | Anyone who can publish `tak.cot.mesh.tx` or `tak.cot.inject` can put markers on every TAK device | Give the edge leaf account publish permission on `tak.cot.mesh.tx` only for the site projector identity, and on `tak.cot.inject` only for the hub projector (and the gateway, when `MESH_TO_TAK_SERVER=true` is approved for a site) | Default deny is not in place; the README warns |
| P5 | Mesh encryption keys must reach every ATAK device and the gateway. How ATAK 5.6 imports the key pair was not verified on the build machine (no ATAK here) | Keep the pair in the Mac's Keychain (or the site host's secret store) and export it to the gateway at start; the gateway runs outside the cluster, so a SealedSecret does not apply. The owner records the ATAK import steps after the LAN test | Encryption is off unless both variables are set |
| P6 | `docs/STACK.md` lists live cluster workloads only, and the gateway is not a cluster workload | When a gateway runs on the Mac or a site host, add an "Edge hosts (not in the cluster)" section to `docs/STACK.md` with its host, version and interface | No STACK row; the gateway is not deployed |
| P7 | `MESH_TO_TAK_SERVER` | Keep `false` everywhere (decision 0004 #6); `true` only in the owner's local test config | Default `false`; with `true`, `OM.*`, `INTEL.*` and `omos-gw-*` uids are never copied (CP-2) |
| P8 | Peer-hosted https (`peerHosted="true"`, `httpsPort`) and server-hosted (`/Marti/`) data package URLs are refused | https peer hosting needs a certificate for the gateway; server-hosted packages belong to the TAK Server client (feature A21), not the mesh gateway | Refused with a nack and a logged reason |
