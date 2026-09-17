# OneMind Topology — The Four-Plane Infrastructure Map

**Decided:** 2026-07-04 (after five turns of pressure-testing: K8s drawbacks → K8s vs NEX → omnipresence/physics → KubeSolo discovery → final shape)
**Status:** Standing infra doctrine. Build to this. Do not re-litigate.
**Companion:** memory `reference/omos-infra-doctrine.md`, `skills/onemind-stack/SKILL.md`

---

## The one-line bottom line

```
Portainer  = one screen over all of it
Cloud K8s  = always-on hub brain (Hermes, agents, control plane, router)
NEX        = omnipresent reach across the bus (the moat)
KubeSolo   = self-heal on edge boxes, offline-safe
NATS       = the bus under all of it
Mac        = cockpit, not server
```

---

## The core reframe

The real problem was never "K8s yes/no." It was three things:

1. **"I can't see what's running or what's broken"** — the Mac is opaque.
2. **"The Mac↔Docker network is an issue"** — Docker Desktop-on-Mac has a real networking pathology (already hit with the Mumble bridge port-forward bug).
3. **Self-heal is appealing.**

**Root cause:** the Mac is being used as an always-on server, and it is bad at that job.

**The fix is architectural:** the Mac becomes the **cockpit**, not the server. Move always-on load to real Linux hosts (cloud hub + edge boxes). That alone kills the Docker Desktop networking pathology.

---

## The whole system, one diagram

```
╔═══════════════════════════════════════════════════════════════════════╗
║  MANAGEMENT PLANE                                                      ║
║  ┌─────────────────────────────────────────────────────────────────┐  ║
║  │  PORTAINER  — one screen over everything                        │  ║
║  │  sees: cloud K8s pods · edge KubeSolo nodes · Docker · agents   │  ║
║  │  does: visibility · health · deploy · rollback · RBAC · offline │  ║
║  └─────────────────────────────────────────────────────────────────┘  ║
╚═══════════════════════════════════════════════════════════════════════╝
                                   │ observes ▼
╔═══════════════════════════════════════════════════════════════════════╗
║  CLOUD HUB          runtime: MANAGED K8s (DOKS)      always-on, self-heal║
║                                                                        ║
║   ┌── control/brain ──────────┐   ┌── data/agents ─────────────┐       ║
║   │ om-control-plane          │   │ Hermes/Legacy gateway      │       ║
║   │ om-router  (event rules)  │   │ agent workers (Composio…)  │       ║
║   │ om-command-gateway        │   │ Composio connector         │       ║
║   │ om-feed                   │   │ HomeAssistant CONNECTOR    │       ║
║   │ om-mission                │   │ simulation services        │       ║
║   │ om-scheduler              │   │ TAK bridge (primary)       │       ║
║   └───────────────────────────┘   └────────────────────────────┘       ║
║                                                                        ║
║   ┌── stateful (managed first) ──┐  ┌── observability ─────────┐       ║
║   │ Postgres (managed)           │  │ Prometheus/Grafana/Loki  │       ║
║   │ ClickHouse (omos-observ)     │  │ Tempo · alerts           │       ║
║   │ Redis (cache only)           │  │ Synadia Insights/Protect │       ║
║   │ media registry → Spaces/S3   │  └──────────────────────────┘       ║
║   └──────────────────────────────┘                                     ║
║   web console: omos-atoc UI (hosted for remote reach)                  ║
╚═══════════════════════════════════════════════════════════════════════╝
                                   │  NATS leaf ▲▼
╔═══════════════════════════════════════════════════════════════════════╗
║  EDGE / OMNIPRESENT                                                    ║
║                                                                        ║
║   runtime NEX (bus-native organs) ── the MOAT ──────────────────┐      ║
║   │ om-overwatch-{state,relay,tak,video,interceptor,api}        │      ║
║   │ om-overwatch-sentinel · om-overwatch-vision                 │      ║
║   │ om-unity (ORP entity graph, 5,938)                          │      ║
║   │ om-bridge-ros2 (Unitree/Go2)  · mavlink-nats-bridge         │      ║
║   │ omos-voice organs · TAK bridge (local fallback)             │      ║
║   │ hardware/sensor organs (ENV, ZEUS telemetry)                │      ║
║   └─────────────────────────────────────────────────────────────┘      ║
║                                                                        ║
║   runtime KubeSolo (single-node k8s self-heal, offline-first) ──┐      ║
║   │ vendor/off-the-shelf containers that expect Kubernetes      │      ║
║   │ Jetson inference containers · industrial Helm packages      │      ║
║   │ edge boxes: OMOS Dell · future gateways                     │      ║
║   └─────────────────────────────────────────────────────────────┘      ║
║                                                                        ║
║   runtime Docker/VM (local, known-good) ────────────────────────┐      ║
║   │ TAK Server 5.7 (+ Postgres, GeoServer, Murmur)              │      ║
║   │ Home Assistant itself (needs LAN/Zigbee/BLE/mDNS)           │      ║
║   └─────────────────────────────────────────────────────────────┘      ║
║                                                                        ║
║   physical: drones · radios · SDR · cameras · robots · Meshtastic      ║
║             → publish INTO the bus, never orchestrated by cloud        ║
╚═══════════════════════════════════════════════════════════════════════╝
                                   │
╔═══════════════════════════════════════════════════════════════════════╗
║  NATS  — the bus under all planes (leaf nodes everywhere)              ║
║  Envelope v1 contract · JetStream streams · KV · req/reply            ║
║  Tailscale for secure links · one logical fabric across the planet    ║
╚═══════════════════════════════════════════════════════════════════════╝

  MAC = cockpit (dev + command), NOT a server. Runs atoc dev, tooling,
        you. Always-on load lives in the hub + edge boxes.
```

---

## The four planes explained

### 1. Management plane — Portainer
Single pane of glass over **all** of it: cloud K8s pods, edge KubeSolo nodes, Docker containers, agents. Visibility · health · deploy · rollback · RBAC · offline-safe async agents. **This directly fixes the #1 pain: "I can't see what's running or broken."** Three agent types: LAN (trusted nets), Remote (outbound-only, no inbound ports), Async Edge (offline-safe, queues until reconnect).

### 2. Cloud hub — Managed Kubernetes (DigitalOcean DOKS)
The always-on brain. Dense services that must run 24/7, self-heal, and stay remotely reachable without babysitting. This is where K8s earns its keep — restart, deploy, scale, monitor, secure, roll back, automatically, forever, without you touching a terminal.

### 3. Edge / omnipresent — NEX + KubeSolo + Docker/VM
Everywhere the cloud can't reach: the Mac, the Dell, vehicles, drones, offline sites. Three runtimes coexist here, each for a different job (see the seams below).

### 4. NATS — the bus under all planes
Leaf nodes everywhere make one logical fabric across the planet. The Envelope v1 contract (`schemas/envelope/v1`) is the shared language every organ speaks. Tailscale secures the links.

---

## Why NEX is the omnipresent one — physics, not preference

K8s assumes a stable cluster network + etcd quorum + low-latency node membership. That survives a datacenter or a home rack. It does **not** survive a drone, a Meshtastic node, a vehicle on flaky LTE, or an offline NAT'd edge box. That terrain **is** the omnipresent system — so multi-node clustered K8s **cannot** be omnipresent. Architectural, not tunable.

NEX rides NATS leaf nodes that already span the planet, tolerate intermittent links, and traverse NAT through the bus. A NEX workload is **bus-native from birth** — it speaks by NATS subject the instant it starts, no networking to wire. NEX can place an organ on the Mac, a truck, or a cloud node through the same fabric. **For omnipresence, NEX is the only one of the two that can actually be omnipresent.**

---

## What "organs" means

An **organ** is a single bus-native service — one small program that does one job and speaks the bus. The body metaphor fits:

- **NATS** = nervous system (the bus everything signals across)
- **Organs** = functional units — each listens on some subjects, publishes on others
- **The Envelope** = the blood — the shared message contract every organ speaks

Live organs today: `om-overwatch-state` (fusion spine), `om-overwatch-tak` (pushes to TAK), `om-overwatch-interceptor` (the kill), `om-unity` (ORP graph, 5,938 entities), `om-control-plane` (ui.* projections). You grow the body one organ at a time.

---

## NEX vs KubeSolo — different axes, not competitors

| | **NEX** | **KubeSolo** |
|---|---|---|
| **Answers** | "Place this organ *anywhere on the bus*" | "Keep this *one box's* containers alive, even offline" |
| **Native to** | NATS bus (bus-native from birth) | Kubernetes (single-node, no etcd) |
| **Reach** | The whole fabric — Mac, vehicle, cloud, drone | Just the single device it runs on |
| **Networking** | None to wire — speaks by subject | Standard container networking on that box |
| **Best for** | Bus organs: overwatch, bridges, agent workloads on NATS | Off-the-shelf containers/Helm charts on an edge box |
| **Ecosystem** | Sovereign, rare — **your moat** | Commodity k8s tooling (Helm, kubectl, Argo, K9s) |
| **Maturity** | Young (Synadia) | Hardened k3s fork; not CNCF-certified |

**How to choose:**
- **Use NEX** when the thing is an *organ* — lives on the bus, talks by subject, must be placed across the fabric. This is most of OneMind's backend.
- **Use KubeSolo** when you have an *edge box* running normal containerized software (a vendor Helm chart, a Jetson inference container, an industrial package) and you want K8s-style self-heal + declarative config that survives going offline.

**They coexist.** One edge box can run KubeSolo (for a vendor container) *and* NEX organs (for the bus) at once — different jobs, same hardware. Never point both at the same workload.

**Honest simplification:** NEX is for *your* bus-native organs (the moat). KubeSolo is for *other people's* containers that expect Kubernetes, out at the edge, offline-safe. Cloud hub = full managed K8s. NATS under all of it.

---

## The three laws baked into the diagram

1. **One orchestrator per workload** — a service appears in exactly one runtime box.
2. **Everything speaks the Envelope** — so any box can move without a rewrite. This is the escape hatch: nothing is ever trapped in any runtime.
3. **The core loop is local** — Overwatch detect→fuse→C2→intercept lives in NEX on the edge. Cloud can vanish; OneMind still fights. Cloud is coordination + dense services, never the heartbeat.

---

## The seams (never violate)

- **Cloud K8s owns placement in the hub.** Dense always-on services.
- **NEX owns placement across the bus.** The omnipresent fabric / edge / vehicles / offline.
- **KubeSolo owns self-heal on a single edge box** where you want K8s semantics offline.
- **One orchestrator per workload — never two managing the same thing.** They only collide if you run NEX inside the cloud K8s cluster managing what K8s already manages. Don't.
- **Everything NATS-native** so nothing is trapped — falls back to NEX/Docker/systemd/K8s without a rewrite.
- **Core loop never depends on the cloud.**
- **TAK Server + Home Assistant stay Docker/VM/local.** Only their bridges/connectors move to hub/edge.

---

## The moat — do not lose it in the shine

K8s, KubeSolo, and Portainer are **commodity** — lots of people run them, low regret to learn. The **moat** is the NATS-native, bus-omnipresent, sovereign fabric: NEX + NATS security + the Envelope contract. Almost nobody can run a sovereign, bus-native, edge-omnipresent fabric. Adopt the commodity tools as operational muscle, but the law that makes OneMind *OneMind* is: **every organ speaks the bus.** Keep that and you can swap any runtime underneath.

---

## Honest caveats (standing)

- **Maturity risk lives on the NEX side** (young, Synadia). K8s/KubeSolo are more battle-tested. Mitigation = the scoping above.
- **KubeSolo is not CNCF-certified** (memory-optimized fork). Scope it to edge boxes only; cloud hub = standard managed K8s.
- **Portainer is a management layer, not magic.** Gives visibility + governance; you still design the services. Worth it — fixes the #1 visibility pain.
- **Don't stack all runtimes on the Mac.** Mac is the cockpit. Hub + edge boxes do the running.

---

## Live today vs where it's headed

- **Live now (edge/NEX + Docker on Mac):** Overwatch organs, om-unity, TAK Server, om-control-plane MVP, atoc dev. This is the working core loop — it stays where it is.
- **The migration = build the cloud hub**, then move the *always-on* services (Hermes, agents, router, feed, command-gateway, observability, connectors) off the Mac into K8s, and stand up Portainer over the whole thing.
- **Untouched:** the physical layer and the core kill chain stay on the edge/bus by design.

---

## Adoption stance

Zeus doesn't want "simple" — AI can manage complexity, and he's willing to learn. Goal is the **architecturally best design for an omnipresent, always-on, proactive AI agent system for TAK + all operations**, not the easy one. Go all in on learning cloud K8s as the commodity layer under the moat. Treat the cloud K8s buildout as a **bounded 30-day capability test**: keep it if it proves power; extract patterns and fall back to simpler infra if it slows us. No DigitalOcean lock-in (portable Helm/Argo/Terraform/manifests). Databases managed first, self-host later when backup/restore discipline exists.
