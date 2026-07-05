# Legacy Agent Fabric

> NATS-native, self-healing, omnipresent design for Legacy/Hermes.
> Mirrors the agent-memory design doc (`reference/legacy-agent-fabric.md`); this repo copy
> is the shared/team-visible version. Deployment lives in `platform/hermes/`, tracked by
> issue #1.

## Thesis
Legacy stops being "a process on the Mac" and becomes **an identity on the NATS fabric**
that can wear many bodies at once — all reachable as one, all self-healing, all governed by
code. That is what makes Legacy omnipresent and millisecond-fast.

## Foundation: the Synadia Agent Protocol (open, NATS-native)
- A **contract, not a runtime/framework.** An agent = a NATS micro service named `agents`
  serving 3 endpoints: `prompt`, `status`, `hb` (heartbeat).
- Subjects: `agents.prompt.{agent}.{owner}.{name}` (verb-first). Legacy today =
  `agents.prompt.legacy.zeus`.
- **Discovery = standard NATS micro:** `$SRV.PING.agents` / `$SRV.INFO.agents`. No registry,
  no coordinator process to keep alive.
- **Liveness:** heartbeat on `agents.hb.*.*.*` (30s cadence); offline after 3 missed beats.
- **Two SDKs:** **host SDK** (`@synadia-ai/agent-service` / `synadia-ai-agent-service`) puts
  an agent on the fabric; **caller SDK** (`@synadia-ai/agents` / `synadia-ai-agents`) builds
  META AGENTS (discover, fan-out, merge, orchestrate).
- Spec: `github.com/synadia-ai/synadia-agent-sdk-docs` · SDKs/harnesses:
  `github.com/synadia-ai/synadia-agents` · Blog: nats.io/blog/nats-native-protocol-for-ai-agents.
- Hermes already speaks this; Legacy already lives at `agents.prompt.legacy.zeus`.

## Architecture
1. **One Legacy, many bodies (omnipresence + millisecond routing).** The `prompt` endpoint
   registers with NATS **queue group `agents`**. Run multiple Legacy instances on the SAME
   subject — cloud (primary, always-on) + Mac (warm edge). NATS auto-routes each prompt to
   whichever instance is alive/nearest. Kill one → the other answers next-millisecond.
   Automatic HA + failover, no load balancer, no failover code.
2. **Two-layer self-healing.**
   - *Infra (K8s):* Deployment restarts/reschedules a dead pod; liveness probe = heartbeat.
   - *Agent (self-cron, permissioned):* in-cluster Legacy runs a scheduled self-check —
     reads its own logs/metrics off the bus and (via `kubectl`-over-tailnet + GitOps) can
     **open a PR against this repo to change its own deployment.** Propose → approve → merge.
     Never blind self-mutation.
3. **The meta agent (conductor).** Built with the **caller SDK**. Discovers every agent,
   fans a prompt to all, merges responses, watches every heartbeat. Legacy's hands for
   commanding the whole swarm (GLM workers, sub-agents, TAK, cot-tap) as one org.
4. **Per-agent zero-trust identity (nauth JWT CRDs).** Every body (Legacy-cloud, Legacy-mac,
   meta-agent, each worker, cot-tap) gets its own NATS **Account/User as a CRD**, GitOps'd
   here. Scoped: who may prompt whom, who may publish where. Compromise one → revoke/scope
   at the bus instantly. Synadia Protect ideology **as code**.

## Guardrails
- **Split-brain on self-modification.** Many bodies is great for *answering*, dangerous for
  *self-editing*. Only ONE body (cloud) holds the "self-improve" lease — enforce via a
  **JetStream KV lock**. Mac stays a responder, never mutates itself.
- **Self-improvement is permission-gated + PR-reviewed.**
- **No paid Synadia product.** Protocol, both SDKs, and `nauth` are all OPEN. We own the
  cluster. Synadia "Deploy for Kubernetes" ($1,250–2,000/mo) is NOT needed.
- **AWS creds gate.** Cloud Hermes uses AWS Bedrock. Rotate the long-lived IAM keys
  (currently plaintext in the Mac launchd plist) and inject as K8s secrets before deploy.

## Build sequence
1. **Cloud Hermes = second Legacy body** on DOKS, same `agents.prompt.legacy.zeus` subject
   → prove failover. (See `platform/hermes/`, issue #1.)
2. **Self-maintenance cron** (permissioned): read-only report → propose-PR → gated-apply;
   JetStream KV lease so only cloud holds it.
3. **Meta agent** (caller SDK): discovery + fan-out + live heartbeat board.
4. **nauth layer:** per-agent Accounts/Users as CRDs committed here.

## Related
- `platform/hermes/README.md` — step 1 deployment plan
- `platform/nats/` — the bus + leaf bridge
- `platform/tailscale/` — mesh, kubectl-over-tailnet, GitOps ACLs
