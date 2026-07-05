# TAK Server 5.7 — OneMind cluster on DOKS

Horizontally-scalable TAK Server (messaging / api / plugins / config split) on
DigitalOcean Kubernetes, adapted from the official TAK Product Center cluster
blueprint, with **Zeus's NATS bridge plugin baked in** so the tactical fabric
projects onto the OneMind bus.

**Status:** scaffold + runbook. Not yet deployed. Two decisions block first apply
(see _Open decisions_).

---

## What this is (and where it came from)

- `upstream/` — the **authoritative blueprint**, vendored verbatim from
  `github.com/TAK-Product-Center/Server` → `src/takserver-cluster` (public mirror,
  chart appVersion 4.5.0 / ~TAK 5.2). This is the real Helm chart + Dockerfiles +
  Ignite/NATS wiring + the AWS/kops `ManualClusterGuide.md`. We adapt it; we don't
  fork the logic.
- `values-onemind.yaml` — our DOKS overlay values (sized down, DOCR images, namespace
  `takserver`, in-cluster PostGIS to start, internal NATS at 1 replica).
- `docker/Dockerfile.takserver-plugins.onemind` — the plugins image, extended to bake
  in `takserver-nats-plugin-1.0.0.jar` + `NatsBridgePlugin.yaml`.
- `plugin/` — the plugin's config (tracked) + its JAR (copied at build time, gitignored).

## Architecture (what actually gets deployed)

| Tier | Pod(s) | Ports | Scale knob |
|------|--------|-------|-----------|
| **config** | `takserver-config` | 10800 (ignite) | 1 (singleton) |
| **messaging** | `takserver-messaging` | 8089 CoT/TLS, 9000/9001 fed | **scale with EUD count** |
| **api** | `takserver-api` | 8443 WebTAK, 8444 fed, 8446 enroll | scale with UI/API load |
| **plugins** | `takserver-plugins` | — (egress to OneMind NATS) | 1 |
| **ignite** | statefulset | 10800 | cache/discovery backbone |
| **nats (internal)** | statefulset | 4222 | TAK-internal only |
| **postgres** | PostGIS | 5432 | in-cluster now → DO Managed later |

Config is delivered by **three ConfigMaps** the pods mount:
- `cert-migration` → `/certs-configmap/` (CA + X.509, generated via `Dockerfile.ca`)
- `core-config` → `/CoreConfig.xml`
- `tak-ignite-config` → `/TAKIgniteConfig.xml`

### Two different NATS — do not conflate
1. **Internal NATS** (`nats.*` in values) — the cluster's own inter-pod JetStream bus.
   Private plumbing. Never touches OneMind.
2. **OneMind NATS** — your sovereign bus. Reached **only** by Zeus's bridge plugin in
   the `plugins` pod, via `plugin/tak.server.plugins.NatsBridgePlugin.yaml`
   (`tak.cot.inject` in, `tak.cot.out` out). This is the seam.

## Hard constraints (learned from the blueprint + upstream issues)

- **Namespace must be `takserver`.** Ignite's `TcpDiscoveryKubernetesIpFinder` hardcodes
  service `takserver-ignite` + namespace `takserver` (issue #43). Deploy elsewhere → no discovery.
- **SSL passthrough** is required if you use an ingress controller — TAK terminates its own
  mTLS. We avoid that by fronting 8089 with a **DO TCP (L4) LoadBalancer** and tunneling
  8443/8446 through **Cloudflare Tunnel** (no TLS termination in the middle).
- Prod defaults in the blueprint are **60 messaging + 30 api** pods (~720Gi RAM). Ignore —
  `values-onemind.yaml` starts at 1/1 and scales up deliberately.

---

## Image build (the sovereign path) — and the honest fork

The blueprint's Dockerfiles (`upstream/docker-files/`) build the cluster images from a
TAK Server **source-build layout**:
`Dockerfile.takserver-base` expects `takserver-core/takserver-core*.war`,
`takserver-core/setenv-cluster.sh`, `takserver-usermanager/UserManager.jar`,
`takserver-schemamanager/SchemaManager.jar`, and the cluster launch/readiness scripts
(`takserver-config-cluster.sh`, `takserver-api-cluster.sh`, `takserver-messaging-cluster.sh`,
`*-readiness.sh`, `db-connection-configuration.sh`).

**The fork:** the single-node distro (`takserver-docker-5.7-RELEASE-43`) ships the WAR,
`setenv-cluster.sh`, `takserver-plugins-cluster.sh`, `UserManager.jar`, `SchemaManager.jar`
— but **may not contain all the config/api/messaging cluster scripts + readiness scripts**,
which live in the source tree. So:

- **Option A — build from source (recommended, needs your tak.gov account):** clone the
  full `git.tak.gov/core/takserver` at the **5.7 tag**, `./gradlew` the WARs, then build the
  cluster images with the vendored Dockerfiles → push to DOCR as `5.7-onemind`. Current + sovereign.
- **Option B — stage from the distro:** arrange the distro files into the Dockerfile layout and
  supply the missing cluster scripts from `upstream/`. Lighter, but we must verify every script exists.
- **Option C — pull provisioned images from tak.gov Artifactory** (`docker-devtest-local.artifacts.tak.gov/takserver-cluster/*`). Fastest, but **5.2 not 5.7**, account-gated, and historically flaky pulls (issue #28). Reference only.

Plugins image (all options): FROM the resulting `takserver-base`, then
`docker/Dockerfile.takserver-plugins.onemind` bakes in the NATS bridge. Copy the JAR first:
```
cp /Users/zeus/ATAK/plugins/nats-bridge/takserver-nats-plugin/build/libs/takserver-nats-plugin-1.0.0.jar plugin/
# (or the built jar already in the distro: .../tak/lib/takserver-nats-plugin-1.0.0.jar)
```

---

## Deploy runbook (DOKS-adapted from ManualClusterGuide.md)

1. **Cluster + registry:** `doctl kubernetes cluster create` (or existing), `doctl registry create onemind`, `kubectl create ns takserver`.
2. **DOCR pull secret:** `kubectl -n takserver create secret docker-registry docr-cred --docker-server=registry.digitalocean.com ...` (or `doctl registry kubernetes-manifest | kubectl apply`).
3. **Certs:** build `Dockerfile.ca` → produces CA + admin.p12 + server certs →
   `kubectl -n takserver create configmap cert-migration --from-file=<certs>`. Back up the `files/` dir.
4. **Database:** in-cluster PostGIS (default values) OR provision **DO Managed PostgreSQL**
   + set `global.postgresql.host` and the `CoreConfig.xml` Repository block.
5. **CoreConfig / Ignite configmaps:** `kubectl -n takserver create configmap core-config --from-file=CoreConfig.xml` and `tak-ignite-config --from-file=TAKIgniteConfig.xml`.
6. **Build + push images** (see Image build) → DOCR as `5.7-onemind`; set `<DOCR>` in `values-onemind.yaml`.
7. **Point the bridge** — edit `plugin/tak.server.plugins.NatsBridgePlugin.yaml` `natsUrl` at the
   real OneMind NATS (not `host.docker.internal`). If the OneMind bus needs TLS/creds, the plugin
   must support them — **review `nats-bridge` source before relying on it** (see Risks).
8. **Deploy:**
   ```
   helm upgrade --install takserver ./upstream/deployments/helm \
     -n takserver --create-namespace -f values-onemind.yaml
   ```
9. **Expose:** DO L4 LoadBalancer → `takserver-messaging-service:8089`; Cloudflare Tunnel → `takserver-api-service:8443/8446`.
10. **Verify:** load `admin.p12`, open `https://<host>:8443`; enroll a client on 8089; publish a CoT
    to `tak.cot.inject` on the OneMind bus and confirm it lands on the client; watch `tak.cot.out`.

## Risks / open verifications
- **Cluster scripts completeness** in the distro (drives Option A vs B).
- **Plugin scan dir** in cluster mode — confirm `takserver-pm.jar` loads the JAR from `lib/`.
- **Bridge auth/TLS** — does `nats-bridge` v1.0.0 support the secured OneMind bus? (review source).
- **Ignite RBAC** — `takserver-ignite` ServiceAccount needs endpoints read in ns `takserver`.

## Open decisions (need Zeus)
1. **Image path:** A (build from source, tak.gov) / B (stage from distro) / C (pull provisioned)?
2. **Database:** in-cluster PostGIS first, or straight to DO Managed PostgreSQL?
3. **8089 exposure:** DO L4 LoadBalancer (public) vs Tailscale-only ingress?

## Provenance
Blueprint: TAK-Product-Center/Server `src/takserver-cluster` (public mirror, ~5.2).
Adapted for OneMind / DOKS, July 2026. Plugin: Zeus's `nats-bridge` (takserver-nats-plugin v1.0.0).
