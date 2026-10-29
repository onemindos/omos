---
name: tak-server
description: "Use for TAK Server on OneMind: topology, how data reaches the TAK map over NATS, the Marti REST API, CoT format and client enrollment."
version: 2.0.0
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
deployment and contained credentials. Authoritative sources, in order:
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

## REST API (Marti)

Base `https://takserver-api-service.takserver.svc.cluster.local:8443/Marti/api/`,
mutual TLS. Useful read endpoints: `version`, `contacts/all`, `groups/all`,
`missions`, `missions/{name}`, `missions/{name}/cot`, `devices`, `excheck/templates`;
file sync under `/Marti/sync/content`. Do not create users, groups or certificates
without the owner's approval. Never put passwords or certificate passphrases in
a command line, a file in the repo or a reply.

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
