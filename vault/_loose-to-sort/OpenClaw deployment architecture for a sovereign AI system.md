# OpenClaw deployment architecture for a sovereign AI system

**OpenClaw's single-gateway design means true active/active clustering isn't natively supported — but a resilient, self-healing multi-node architecture is absolutely achievable.** The recommended path for One Mind OS combines a warm-standby Gateway pair with Cloudflare DNS failover, NATS JetStream for state replication, Tailscale mesh networking, and OpenClaw's native `sessions_spawn` for agent duplication. This architecture delivers near-zero downtime, automatic failover, and the agent-spawning capabilities needed for "Legacy" to distribute work across nodes — all without compromising security or feature access.

OpenClaw (356k+ GitHub stars, MIT license) is a hub-and-spoke AI orchestration runtime built around a single Gateway process. Understanding this core constraint is essential: **the Gateway is a single-process, local-first Node.js application** that manages sessions as append-only JSONL files, runs Lobster workflows in-process, and coordinates channel adapters, tool execution, and LLM calls from one control plane. Every architecture decision flows from this reality.

---

## Every deployment method ranked for this use case

OpenClaw supports six deployment methods. For a solo operator running a sovereign AI system on bare metal DigitalOcean, the choice matters more than usual because it determines your self-healing options, security posture, and operational overhead.

**Native npm/pnpm on bare metal with systemd is the clear winner.** It's OpenClaw's primary install path (`npm install -g openclaw@latest && openclaw onboard --install-daemon`), requires Node.js 24 (recommended) or 22.16+, and runs as a systemd user service on Linux. The Gateway is remarkably lightweight — roughly **1 vCPU and under 1 GB RAM** — leaving the vast majority of your 4–8 vCPU / 16–32 GB for LLM inference, NATS, TimescaleDB, and LiveKit. Fewer abstraction layers means easier debugging, direct filesystem access for session data, and native systemd integration for process management. The OpenClaw community's "Clawculus" deployment guide explicitly states: "fewer abstractions, easier debugging, works on smaller VPS sizes. Use systemd for auto-restart and you're golden."

**Docker** (official GHCR image at `ghcr.io/openclaw/openclaw:latest`) is the second-best option. It adds isolation and built-in `HEALTHCHECK` (`node dist/index.js health`), and the repo ships `docker-compose.yml` plus a `docker-setup.sh` one-step script. Docker Sandbox mode provides tool execution isolation. The tradeoff: an extra abstraction layer, volume mount management for `~/.openclaw`, and potential networking complexity when integrating with Tailscale and NATS on the host.

**Podman** has explicit support via `setup-podman.sh` with rootless containers and systemd Quadlet integration (`./setup-podman.sh --quadlet`). It creates a dedicated `openclaw` system user with `nologin` shell. Equivalent to Docker with better security defaults, but smaller community ecosystem.

**NixOS** has first-class support via `openclaw/nix-openclaw` and community modules providing `services.openclaw.*` declarations with built-in systemd sandboxing and tool allowlists. Excellent for reproducibility but adds declarative-config complexity.

**Kubernetes** has a community operator (`openclaw-rocks/k8s-operator`) supporting multi-replica deployments with Prometheus metrics, but requires RWX storage (NFS/CephFS/Longhorn) and is overkill for a solo operator with two nodes.

**NemoClaw** (NVIDIA's hardened reference stack) bundles OpenClaw inside Docker with k3s and the NVIDIA OpenShell runtime. It's early preview, ~2.4 GB compressed, and designed for security-first deployments. Worth monitoring but not production-ready.

**Recommendation: bare metal npm + systemd on both nodes.** This gives you the thinnest possible runtime, maximum resource availability, direct integration with your existing NATS/Tailscale/LiveKit stack, and native systemd watchdog capabilities. Skip Docker unless you specifically need tool execution sandboxing, which can be added later via Docker Sandbox mode even on a bare-metal install.

---

## The active/active question and what OpenClaw actually supports

OpenClaw's architecture is fundamentally **single-gateway**. The official documentation and community guides are unambiguous: "OpenClaw's current architecture is designed around a single-gateway model." Everything — routing, authentication, session management, channel adapters, cron scheduling, Lobster workflows — flows through one Gateway process. Sessions are JSONL files on the local filesystem. Cron state lives in a local `jobs.json`. Memory files are local.

**True active/active with two Gateways processing requests simultaneously introduces three hard problems.** First, session state conflicts: two Gateways writing to the same session JSONL files create race conditions. Second, cron deduplication: both Gateways would fire the same scheduled jobs. The Kubernetes documentation addresses this with `OPENCLAW_SKIP_CRON=1` on replica pods, but that means only one node runs cron — breaking the active/active model. Third, channel binding: messaging adapters (WhatsApp via Baileys, Telegram via grammY, etc.) maintain stateful connections. Two Gateways binding the same WhatsApp account would conflict.

**The pragmatic architecture is active/warm-standby with fast automated failover**, not active/active. Here's the concrete design:

The **primary node** (initially DigitalOcean, eventually home) runs the active Gateway — all traffic, all channels, all cron, all Lobster workflows. The **standby node** runs an identical OpenClaw installation with the same configuration, SOUL.md, skills, and credentials, but with the Gateway **stopped or in a receive-only sync mode**. NATS JetStream continuously replicates session data, memory files, and configuration between nodes. When the primary fails, Cloudflare detects it via health check, routes traffic to the standby, and a systemd-triggered script starts the standby Gateway. Recovery time: **60–120 seconds** (Cloudflare health check interval + Gateway startup).

To achieve the closest thing to active/active, you can run both Gateways simultaneously with **role partitioning**: one Gateway handles interactive messaging channels while the other handles background agent work (sub-agent execution, Lobster workflows, cron). This avoids state conflicts because each Gateway owns distinct workloads. NATS JetStream bridges them for coordination. This is not OpenClaw-native but is architecturally sound.

---

## NATS JetStream as the state sync layer between nodes

OpenClaw **does not use NATS internally** — its Gateway communicates via WebSocket and manages state through local filesystem operations. However, NATS JetStream is an excellent choice for bridging two OpenClaw instances externally, and a community project (`PeachForge/openclaw-nats-cluster`) demonstrates this pattern.

**The recommended NATS JetStream architecture for One Mind OS** uses three sync streams. First, a **session replication stream** that watches `~/.openclaw/agents/*/sessions/*.jsonl` via `inotifywait` or a custom file watcher, publishing new JSONL entries to a NATS subject like `openclaw.sessions.>`. The standby node subscribes and appends entries to its local session files. Because sessions are append-only JSONL, this is conflict-free — you're only appending new lines. Second, a **config/memory sync stream** that replicates `~/.openclaw/openclaw.json`, SOUL.md, skills, and memory files via NATS KV store, using file modification timestamps for conflict resolution. Third, a **health/coordination stream** for leader election and failover signaling — each node publishes heartbeats to `openclaw.health.<node-id>`, and missing heartbeats trigger failover.

NATS JetStream's **at-least-once delivery with message deduplication** ensures no session entries are lost during replication. Its **built-in clustering** across your Tailscale mesh means the message bus itself has no single point of failure. Since you're already running NATS JetStream in your stack, this adds no new infrastructure — just new streams and subjects.

**Key limitation**: this sync is eventually consistent. During a failover, the standby node may be a few seconds behind. For a solo operator, this gap is negligible — you'll lose at most the last few messages in an active conversation, which can be re-sent.

---

## Self-healing: what's built in and what to add

OpenClaw ships with basic health infrastructure: a **`/health` HTTP endpoint** for liveness probes, the **`openclaw doctor`** CLI command (with `--fix` for auto-repair and `--deep` for diagnostics), and **`openclaw security audit`** for configuration review. On bare metal with systemd, the daemon installs with `Restart=always` and `RestartSec=1`, so the process auto-restarts on crash.

For a sovereign AI system, you need more. The community project **`Ramsbaby/openclaw-self-healing`** provides a battle-tested 4-tier autonomous recovery system. Tier 1 is KeepAlive via systemd restart. Tier 2 is a watchdog script monitoring every 15 seconds, triggering recovery after 3 consecutive failures. Tier 3 invokes `openclaw doctor --fix` automatically when crash count exceeds 2. Tier 4 sends alerts via Discord or Telegram. This system reduces average recovery time from **36 minutes to 7 minutes**. Install with a single command: `curl -fsSL https://raw.githubusercontent.com/Ramsbaby/openclaw-self-healing/main/install.sh | bash`.

**Recommended self-healing stack for One Mind OS:**

- **systemd** with `Restart=always`, `RestartSec=1`, `WatchdogSec=30` for process-level recovery
- **Ramsbaby/openclaw-self-healing** watchdog for application-level recovery and `doctor --fix` automation
- **NATS heartbeat** stream for cross-node health awareness
- **Cloudflare health checks** hitting `/health` endpoint for DNS-level failover
- **systemd sandboxing** for bare metal hardening: `ProtectSystem=strict`, `ProtectHome=read-only`, `PrivateTmp=true`, `NoNewPrivileges=true`, `ReadWritePaths=/home/openclaw/.openclaw`

This creates four independent failure detection layers. A crashed process restarts in 1 second. A hung Gateway triggers watchdog recovery in under a minute. A dead node triggers Cloudflare failover in 60–120 seconds.

---

## Lobster workflows and agent spawning across nodes

**Lobster is explicitly local-first and cannot distribute workflows across nodes.** The documentation states: "Install the Lobster CLI on the same host that runs the OpenClaw Gateway." Workflows execute in-process inside the Gateway. Resume tokens are stored locally under Lobster's state directory. An open GitHub issue (#26872) requesting cross-restart continuation confirms this limitation — if the Gateway restarts, workflows don't automatically resume; the user must manually prompt continuation. **Cross-node workflow migration is not supported.**

For One Mind OS, this means Lobster workflows must be treated as node-local. If the active node fails mid-workflow, the workflow state is lost. The mitigation: design Lobster workflows to be **idempotent and re-triggerable**. Use approval gates as checkpoints. After failover, Legacy can re-invoke the workflow on the new node from the last approved step.

**Agent spawning is where OpenClaw shines for your "Legacy duplicates herself" requirement.** The `sessions_spawn` tool creates child agent sessions dynamically and non-blockingly. Each sub-agent runs in its own isolated session with independent context and token usage. Key configuration for the orchestrator pattern:

```json
{
  "agents": {
    "defaults": {
      "subagents": {
        "maxSpawnDepth": 2,
        "maxChildrenPerAgent": 5,
        "maxConcurrent": 8,
        "runTimeoutSeconds": 900
      }
    }
  }
}
```

With `maxSpawnDepth: 2`, Legacy (the main agent) can spawn orchestrator sub-agents, which can in turn spawn worker sub-agents — creating a **Main → Orchestrator → Workers** hierarchy. Sub-agents can use different models (e.g., Opus for orchestration, Haiku for workers via `agents.defaults.subagents.model`) to optimize costs. Sub-agents announce results back to the parent's chat channel when complete.

**However, all sub-agents run on the same Gateway node.** There is no native mechanism to spawn a sub-agent on a remote node. To distribute agent work across nodes, you'd need to use NATS JetStream as a task queue: Legacy on Node A publishes a task to a NATS subject, and a custom listener on Node B invokes OpenClaw's API to spawn a local sub-agent with that task. This is custom integration work, not built-in, but it's architecturally clean and leverages your existing NATS infrastructure.

---

## Security hardening for bare metal without NemoClaw

Running OpenClaw bare metal as a solo operator with no Docker isolation demands rigorous hardening. The Gateway runs with your full user permissions — every file you can read, it reads; every command you can run, it runs.

**Network isolation is the single most important step.** Bind the Gateway to loopback only (`"gateway": {"bind": "loopback"}`) and use **Tailscale Serve** for all remote access. OpenClaw's native Tailscale integration (`"gateway": {"tailscale": {"mode": "serve"}}`) makes the Gateway accessible only within your tailnet while remaining invisible to port scanners. This is described in the security documentation as "the single highest-impact security change you can make." Configure UFW to deny all incoming except SSH: `ufw default deny incoming && ufw allow 22/tcp`.

**Exec approvals are your primary security boundary on bare metal.** Configure `~/.openclaw/exec-approvals.json` with `"security": "allowlist"` and `"ask": "always"` with `"askFallback": "deny"`. This means only explicitly allowlisted commands execute, human approval is required every time, and if no UI is available to prompt, the request is denied. Allowlist specific binary paths (`/usr/bin/git`, `/usr/bin/rg`) rather than basename matches.

**Additional hardening layers**: run OpenClaw under a dedicated non-root user with `nologin` shell; set `chmod 600` on `~/.openclaw/openclaw.json` and all credential files; use environment variables for API keys; add SOUL.md security rules forbidding destructive operations; run `openclaw security audit` after every configuration or skill change; and install only skills from trusted sources — the ClawHub registry has had documented incidents with malicious skills. The automated `openclaw-hardened-ansible` playbook from Next-Kick can bootstrap most of these controls in one shot.

---

## Cloudflare DNS failover configuration

**Cloudflare Load Balancing** at **$5/month** is purpose-built for this use case. It includes 2 origins, 60-second health check intervals, and automatic failover — exactly what two-node active/warm-standby requires.

Create two origin pools: Pool 1 (cloud node) and Pool 2 (home node). Attach an HTTP health monitor hitting OpenClaw's `/health` endpoint with expected response code 200. Set traffic steering to **"Off" (Failover mode)** — this routes all traffic to Pool 1 until it fails health checks, then automatically shifts to Pool 2. When Pool 1 recovers, traffic returns automatically. For proxied (orange-cloud) traffic, **failover is near-instant** because it doesn't depend on DNS TTL propagation.

The alternative — and arguably better — approach uses **Cloudflare Tunnel** (`cloudflared`) on each node instead of exposing origin IPs directly. Each node runs its own tunnel. Load Balancing supports tunnel origins natively, and this keeps both servers completely invisible to the public internet while Tailscale handles inter-node communication. Combined with your existing Tailscale mesh, this creates a zero-exposed-ports architecture where the only ingress path is through Cloudflare's edge network.

To eventually flip home as primary and cloud as standby, simply reorder the pool priority in Cloudflare's dashboard — Pool 2 (home) becomes primary, Pool 1 (cloud) becomes fallback. No DNS changes, no configuration pushes, no downtime.

---

## Concrete architecture recommendation for One Mind OS

The complete architecture layers as follows. **Tailscale mesh** connects all nodes (cloud, home, Mason Android devices, Unitree robots via ROSClaw) into a single encrypted network. **OpenClaw Gateway** runs bare metal on the primary node with systemd and the self-healing watchdog stack. **NATS JetStream** handles three roles: session state replication between nodes, cross-node task delegation for agent spawning, and health heartbeat coordination. **Cloudflare Load Balancing** provides the external failover trigger. **Lobster** runs locally on the active Gateway for deterministic workflows. **Sub-agents** spawn locally via `sessions_spawn` with `maxSpawnDepth: 2` for Legacy's self-duplication capability.

The migration path is straightforward. Phase 1: deploy OpenClaw on DigitalOcean bare metal, harden security, configure Cloudflare health checks. Phase 2: set up the home node as a warm standby with NATS JetStream session sync. Phase 3: flip primary to home, cloud becomes standby. Phase 4: build custom NATS-based cross-node agent delegation for distributed sub-agent spawning. Each phase is independently valuable and doesn't require the next.

**The honest constraint**: OpenClaw wasn't designed for multi-Gateway distributed operation. You're engineering around a single-process architecture. The session replication, failover automation, and cross-node agent delegation are all external integrations you'll build and maintain. This is tractable for a power user — the primitives (NATS, Tailscale, Cloudflare, systemd) are all proven infrastructure — but it's custom work, not turnkey. The OpenClaw community is actively discussing native clustering (K8s operator, NATS cluster project, feature requests for cross-node workflow resumption), so some of this may become built-in over the next 6–12 months.