---
name: tak-server
description: "Use for TAK Server on OneMind: topology, how data reaches the TAK map over NATS, the Marti REST API, CoT format and client enrollment."
version: 3.0.0
category: devops
triggers:
  - "TAK server API"
  - "send data to TAK"
  - "TAK map"
  - "CoT injection"
  - "tak.cot"
  - "Marti API"
  - "WebTAK"
  - "TAK enrollment"
  - "GeoChat"
---

# TAK Server on OneMind

Rewritten from the earlier long skill, which described a retired public
deployment and contained credentials, then merged (2026-10-29) with the
still-true parts of the old Legacy-pod version. Mission procedures are in
`tak-mission-ops`, `tak-mission-runner` and `mission-replay`. Authoritative
sources, in order:
`platform/takserver/DEPLOY-RUNBOOK.md`, `platform/takserver/README.md`,
`docs/TAK-PROJECTION.md`, `docs/NATS-DOCTRINE.md`, and the live cluster.

## Topology (verify with kubectl before relying on it)

- TAK Server 5.8 on OneMind hardened images (`ghcr.io/onemindos/takserver-base:*-5.8-onemind`),
  namespace `takserver`.
- The core (API, messaging, config, plugins, Ignite, PostgreSQL) is the Helm
  release `takserver` from `platform/takserver/chart` with `values-onemind.yaml`,
  applied by hand. An Argo commit to those files does not deploy by itself.
  Do not `helm upgrade` it without the owner's approval and the runbook.
- Argo manages only the overlay (`platform/takserver/deploy`: LB proxy, NATS leaf,
  Mumble, LDAP, secrets) and FedHub (`platform/takserver/fedhub`).
- There is no public listener. Clients (ATAK, iTAK, WinTAK, WebTAK) reach TAK over
  the tailnet at `tak.tail717ba0.ts.net` (ports 8089 CoT, 8443/8446 web/API, 64738 Mumble).
  The old public IP and `tak.onemindos.dev` are retired; federation ports are not exposed.
- In-cluster service names: `takserver-messaging-service` (8089),
  `takserver-api-service` (8443/8444/8446).

## Getting data onto the map

Preferred path: publish CoT XML to `tak.cot.inject`. The NatsBridgePlugin in the
plugin pod subscribes to it and injects into TAK; TAK-originated traffic comes
out on `tak.cot.out`. Publishing is only authorized on the subjects your bus
identity may use; see `nats-operations`.

```bash
nats pub tak.cot.inject '<event version="2.0" uid="demo-1" type="a-f-G-U-C" time="..." start="..." stale="..." how="m-g"><point lat="41.0" lon="-73.6" hae="0" ce="10" le="10"/><detail><contact callsign="DEMO"/></detail></event>'
```

- A publish is not proof of display. Prove each stage separately: bus delivery,
  TAK acceptance, client rendering (see `omos-feed-builder`).
- Feed collectors must not open their own TAK connection; use the bus and the
  shared projection boundary.
- Direct TLS to the messaging service on 8089 needs a client certificate, which
  is a credential. This agent has no standing client certificate in Git. Use it
  only if the owner provisions one as a SealedSecret.

## Bus direction, in one line

TAK to bus: `tak.cot.out` (all EUD traffic, also captured by stream `COT`
as `tak.cot.>`). Bus to TAK: `tak.cot.inject`. Read `tak.cot.out`; inject on
`tak.cot.inject`; never publish `evt.tak.cot.*` by hand (`docs/NATS-DOCTRINE.md` 9.8).
An older field note from the previous pod said nothing consumed `tak.cot.inject`.
If a published event never reaches a client, do not fall back to a raw TLS
connection to 8089 (that needs a client certificate this agent does not hold).
Check the `takserver-plugins` pod logs and the `COT` stream, report the gap, and
let the owner decide.

## REST API (Marti)

Base `https://takserver-api-service.takserver.svc.cluster.local:8443/Marti/api/`,
mutual TLS. Useful read endpoints: `version`, `contacts/all`, `groups/all`,
`missions`, `missions/{name}`, `missions/{name}/cot`, `devices`, `excheck/templates`;
file sync under `/Marti/sync/content`. Do not create users, groups or certificates
without the owner's approval. Never put passwords or certificate passphrases in
a command line, a file in the repo or a reply.

More endpoints that exist (read-only unless the owner approves): `contacts/all`,
`cot/xml/{uid}` (latest CoT for one UID), `cot/sa` (all friendly positions),
`subscriptions/all`, `missions/{name}/{log|cot|contents|changes|subscriptions}`,
`excheck/templates` (POST-only for create, XML body), `video` (JSON, not the XML
`/Marti/vcm`), `tls/profile/enrollment`. Mission names with spaces must be
URL-encoded. There is no OpenAPI endpoint. Port 8446 (WebTAK) does not need a
client certificate but many admin routes still return 403 there; admin work
goes to 8443 with a certificate, which the owner provisions.

## CoT format

`references/tak-cot-format.md` has the event structure, common types and detail
elements. Key points: `stale` should be short for live 1 Hz feeds so an icon disappears when the feed
drops, `track speed` units differ between the reference's two tables (knots in the
aircraft example, m/s in the generic table), so check the consumer, and use the
`uid` of the thing being reported, not of the sender.

## GeoChat

Chat arrives as CoT with `b-t-f` type. An agent filter that accepts only direct
messages silently drops "All Chat Rooms" traffic; accept both when the agent is
meant to hear the whole room. Diagnose in order: is the tap alive (`COT` stream
advancing), does TAK show the client connected, is the message on `tak.cot.out`.

## Enrollment and certificates

Enrollment, data packages and certificate issuance are owner-led and use the TAK
Portal (`enroll.tail717ba0.ts.net`, tailnet only) and the runbook. This agent
does not generate, decrypt or distribute client certificates or keys.

## Operations notes that still hold (verify live before acting)

- **Credentials live in SealedSecrets**, not here: `tak-server-stores` (keystore
  and truststore passwords) and `tak-ca-credentials` (CA), both in `takserver`.
  After the 2026 rotation the old default passphrases are invalid; never guess
  one and never ask for one in chat.
- **Retention.** `cot_router` has no built-in expiry. The Argo-managed CronJob
  `cot-retention` (`platform/takserver/deploy/cot-retention.yaml`) deletes rows
  older than 24 h in batches. Plain `VACUUM` frees space for reuse but does not
  shrink files; `VACUUM FULL` needs an exclusive lock and an owner window. A
  full Postgres volume crash-loops the database: check retention first.
- **WebTAK/API dead while ATAK stays connected** (8446/8443 down, 8089 up) is
  usually the `takserver-api` pod stuck `0/1 Running` with an Ignite
  `ClassNotFoundException ... Requesting mapping from grid failed` in its logs.
  A rollout restart of that Deployment clears it (about 90 s). It is a mutation:
  propose it with the log evidence, then act when the owner agrees. The core is
  a Helm release, so never `helm upgrade` it for this.
- **ConfigMaps mounted with `subPath` do not hot-update**; pods must be recreated.
  Container commands that override the image entrypoint must source `/setenv.sh`
  first or Java cannot find the TAK classes.
- **Groups.** Users in different groups cannot see each other's markers
  (fragmented picture). Put users who should share a picture in the same group.
- **ATAK client facts.** With file or LDAP auth ATAK always shows a
  username/password dialog; that is normal. A stale `truststore-root.p12` that
  predates the current CA makes ATAK fail silently with no server log line;
  rebuild the truststore from the current CA first. P12 files for macOS and
  Android need `openssl pkcs12 -export -legacy`. Android ATAK reads certificates
  from `/sdcard/atak/cert/`. ATAK server settings cannot be driven by ADB
  broadcasts or file writes, and tap automation is fragile: give the owner the
  few manual steps instead. The main activity is
  `com.atakmap.app.civ/com.atakmap.app.ATAKActivityCiv`.
- **`kubectl cp` from ConfigMap mounts** copies symlinks, not contents; read
  files with `kubectl exec ... -- cat` from the resolved path.
- **Voice and video.** TAK ships a Mumble *client* plugin, not a server: Murmur
  (`murmur`, tailnet port 64738) and `mumble-bridge` are the voice server side.
  TAK Server has no built-in video server; video feeds are registered through
  `/Marti/api/video` and served by a separate RTSP source.
- **Map sources and overlays.** ATAK tile source XML uses `{$x}`, `{$y}`, `{$z}`
  (not `{x}`). Distribute maps, KML/KMZ and data packages by uploading to
  Enterprise Sync (`/Marti/sync/upload`) and attaching to a mission, so every
  client receives them without per-device copying. Check each tile provider's
  terms before bulk caching.
- **CloudTAK** (namespace `cloudtak`) is a client and ETL layer on top of TAK
  Server, not a server. ATAK, iTAK and WinTAK connect to TAK Server, not to it.
  Its status changes; check `docs/STACK.md` before planning new work on it.
- **Federation** ports 9101/9102 are not exposed; FedHub is tailnet-only.
