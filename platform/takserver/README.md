# TAK Server 5.8 — OneMind cluster on DOKS

> **2026-10-26 partial Argo ownership:** `platform-takserver` owns the existing
> `takserver-postgresql`, `takserver-nats`, `takserver-ignite` StatefulSets,
> `takserver-api`, `takserver-config` Deployments and `cot-retention` CronJob
> through `deploy/`, in addition to earlier messaging/plugins and overlays. Freeze whole-release
> `helm upgrade`/`helm rollback`/`helm uninstall`: these would create a second writer.
> Helm revision34 remains historical; the rest of the release is not adopted here.
> Follow [the bounded preflight](../../docs/operations/2026-10-26-tak-policy-preflight.md).

**Status (verified live 2026-10-27, read-only):** deployed and healthy in namespace
`takserver` on `omos-cloud`; TAK Server **5.8** on OneMind images; every component runs a
single replica. TAK is **TIER-1**: a change to any TAK pod template restarts TAK for all EUDs.

This folder holds the deployment half only. Image builds, the NatsBridgePlugin source and JAR,
build scripts and `cot-tap` (not deployed) are in `services/takserver/`; the edge TAK deployment
and device enrollment are in `infra/edge/takserver/`. Deploy and upgrade steps:
[`DEPLOY-RUNBOOK.md`](DEPLOY-RUNBOOK.md).

## Layout and ownership

| Path | What | Applied by |
|------|------|-----------|
| `chart/` + `values-onemind.yaml` | TAK cluster Helm chart (adapted from the TAK Product Center `takserver-cluster` blueprint) | Helm release `takserver`, revision 34 (2026-09-29), **frozen**: no whole-release Helm operations |
| `deploy/` | Adopted core workloads, Services, ConfigMaps, LB proxy, NATS leaf, LDAP, Mumble, `cot-retention`, SealedSecrets | Argo `platform-takserver` (plain directory; `CoreConfig.xml` and `TAKIgniteConfig.xml` excluded) |
| `fedhub/` | TAK Federation Hub + MongoDB + tailnet UI | Argo `platform-fedhub` |

## What runs (namespace `takserver`)

| Component | Kind | Image | Ports |
|-----------|------|-------|-------|
| `takserver-messaging` | Deployment | `ghcr.io/onemindos/takserver-base@sha256:…` (5.8, digest-pinned) | 8089 CoT/TLS, 9000/9001 |
| `takserver-api` | Deployment | `ghcr.io/onemindos/takserver-base:api-5.8-onemind` | 8443 Marti/WebTAK, 8444 fed_https, 8446 enrollment/OAuth |
| `takserver-config` | Deployment | `ghcr.io/onemindos/takserver-base:config-5.8-onemind` | internal |
| `takserver-plugins` | Deployment | `ghcr.io/onemindos/takserver-base@sha256:…` (5.8 + NatsBridgePlugin) | none inbound |
| `takserver-ignite` | StatefulSet | `ghcr.io/onemindos/takserver-base:ignite-5.8-onemind` | 47100/47500/10800/… (cache and discovery) |
| `takserver-postgresql` | StatefulSet | `postgis/postgis:15-3.4` | 5432 (database `cot`, PVC 20Gi) |
| `takserver-nats` (+ `takserver-nats-box`) | StatefulSet | `nats:2.9.21-alpine` | 4222, TAK-internal only |
| `tak-ldap` | StatefulSet | `ghcr.io/onemindos/tak-ldap:5.8` | 389/636 (CoreConfig LDAP auth) |
| `tak-lb-proxy` | Deployment | `nginx:1.27-alpine` (L4 stream) | 8089, 8443, 8446, 64738 |
| `tak-nats-leaf` | Deployment | `nats:2.10-alpine` | 4222 (from the plugin) |
| `murmur`, `mumble-bridge` | Deployment | Mumble server; Python STT bridge | 64738 |
| `fedhub`, `fedhub-mongodb` | Deployment, StatefulSet | `ghcr.io/onemindos/fedhub:5.8`, `mongo:6.0` | 9100/9443 UI, 9101/9102 federation (ClusterIP only) |

Images are built by `.github/workflows/build-takserver.yml` on the self-hosted runner and pushed
to GHCR (DOCR is not used).

## How clients reach TAK (tailnet only)

```
EUD (ATAK/iTAK/WinTAK/WebTAK) ── tailnet ──► tak.tail717ba0.ts.net
   Service tak-tailnet (LoadBalancer, class tailscale; proxy pod in ns tailscale)
     └─► tak-lb-proxy (nginx stream) ──► takserver-messaging:8089
                                     ├─► takserver-api:8443 / 8446
                                     └─► murmur:64738
```

- There is no public listener: the DigitalOcean LB `tak-external` was removed on 2026-10-23
  and `tak.onemindos.dev` is retired. No Cloudflare Tunnel or ingress-nginx path exists for TAK.
- Federation ports 9101/9102 are not exposed outside the cluster; the FedHub UI is at
  `fedhub.tail717ba0.ts.net`.
- In-cluster clients: CloudTAK (`platform/cloudtak/takserver.yaml`, ExternalName to `tak-tailnet`),
  TAK-ML (`takserver-api-service:8443`). The TAK portal (`onetak/tak-portal`,
  `enroll.tail717ba0.ts.net`) writes `tak.tail717ba0.ts.net:8089` into enrollment packages.

## Two different NATS — do not conflate

1. **Internal NATS** (`takserver-nats`): TAK's own inter-pod plumbing. Never touches OneMind.
2. **OneMind authority bus** (`nats.nats`): reached only by the NatsBridgePlugin in
   `takserver-plugins`, through the in-namespace leaf:
   `plugin → tak-nats-leaf:4222 (plaintext, in-namespace) → leaf remote → nats.nats:7422 (TLS, leaf credentials)`.

Live subjects (`services/takserver/plugin/tak.server.plugins.NatsBridgePlugin.yaml`):

| Direction | Subject | Stream |
|-----------|---------|--------|
| NATS → TAK (inject) | `tak.cot.inject` | — |
| TAK → NATS (all EUD CoT) | `tak.cot.out` | `COT` (`tak.cot.>`) |

[NATS-DOCTRINE §9](../../docs/NATS-DOCTRINE.md) specifies `evt.tak.cot.{uid}` for TAK CoT; nothing
publishes there today (the `GEO` stream is empty). This is an **open decision**, recorded in
NATS-DOCTRINE §9.8. Do not change the runtime subjects without it.

## Data, retention, backup, monitoring

- **CoT history:** table `cot_router` in database `cot` (~11.5M rows / 9.6 GB on 2026-10-27).
  `cot-retention` (every 4 h) is meant to keep 24 h. Until `fix/tak-cot-retention` merges, it
  deletes nothing: Kubernetes `$$` expansion breaks its inline SQL and the error is swallowed.
- **Backups:** `backup-takserver-postgres` dumps the database daily at 02:00 UTC (~2 h 20 min);
  VolSync backs up the PVCs (`platform/backups/volsync/`).
- **Monitoring:** Prometheus scrapes no TAK target. The only TAK alert on `main` is
  `TAKCoTFlowStopped` (COT stream rate); workload, volume and job alerts are proposed in
  `feat/tak-alerts`.
- **Network:** namespace `takserver` has no NetworkPolicy on `main`; ingress policies are proposed
  in `feat/takserver-netpol`.

## Hard constraints

- **Namespace must be `takserver`.** Ignite's `TcpDiscoveryKubernetesIpFinder` hardcodes service
  `takserver-ignite` in namespace `takserver`.
- TAK terminates its own mTLS. Keep every path L4 (the nginx `stream` proxy and Tailscale); never
  put a TLS-terminating ingress in front of 8089/8443/8446.
- Config is mounted from ConfigMaps: `cert-migration` (`/certs-configmap/`), `core-config`
  (`CoreConfig.xml`) and `tak-ignite-config`. They contain keys and passwords; do not print them.
  `deploy/CoreConfig.xml` still carries the 5.7 `version` attribute from the original distro; the
  running images are 5.8.

## Provenance

Blueprint: TAK-Product-Center/Server `src/takserver-cluster` (public mirror). Adapted for OneMind
on DOKS from July 2026; upgraded to 5.8 OneMind images. Plugin: NatsBridgePlugin
(`services/takserver/plugin-src/`).
