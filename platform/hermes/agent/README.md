# Hermes / Legacy on DOKS — Legacy Agent Fabric, Step 1

**Goal:** run a cloud Hermes as a **second body of the same Legacy identity**, on the same
NATS Agent-Protocol subject as the Mac Hermes (`agents.prompt.legacy.zeus`), so NATS
queue-group HA auto-routes prompts to whichever body is alive/nearest. Cloud = primary
(always-on, self-healing via K8s); Mac Hermes stays as a warm edge body. They talk over
the bus (bridged Mac↔DOKS NATS) via the Synadia Agent SDK.

See the full design: agent memory `reference/legacy-agent-fabric.md`.

## Why this is safe/movable
- Hermes is its own git repo (`~/.hermes/hermes-agent/`), production Dockerfile (multi-arch,
  s6-overlay supervision). Models = AWS Bedrock (no local GPU). Identity via env:
  `HERMES_NATS_SESSION_NAME=legacy`, `HERMES_NATS_OWNER=zeus`.
- In-cluster it connects DOKS NATS directly (`nats.nats.svc.cluster.local:4222`, ONEMIND
  account) and shares subjects with Mac Hermes over the leaf bridge.

## Prereqs (gates)
1. 🔴 **AWS creds rotation** — Bedrock IAM keys are currently plaintext in the Mac launchd
   plist. Rotate in AWS console → inject as K8s secret `hermes-aws` (NOT env-in-plist).
2. ✅ **Bus bridge** — done (Mac↔DOKS leaf, bidirectional).
3. **HERMES_HOME state** — `~/.hermes/` `.env` (~24KB) + `SOUL.md` etc. → K8s secret
   `hermes-env` + PVC for writable `/opt/data`.
4. **Image** — build `~/.hermes/hermes-agent/Dockerfile` for `linux/amd64`, push to DOCR
   `registry.digitalocean.com/onemind-hub/hermes:<tag>`.

## Planned manifests (to add here — GitOps source of truth)
- `namespace.yaml` — ns `hermes` (or `agents`).
- `secret-hermes-env.yaml` — sealed/out-of-band (`.env`), gitignored real values.
- `secret-hermes-aws.yaml` — rotated Bedrock creds, out-of-band.
- `pvc.yaml` — writable state for `/opt/data`.
- `deployment.yaml` — Hermes, `NATS_URL=nats://nats.nats.svc.cluster.local:4222`,
  ONEMIND creds, `HERMES_NATS_SESSION_NAME=legacy`, `HERMES_NATS_OWNER=zeus`,
  liveness probe on the agent heartbeat.
- (later) `self-maintenance-cronjob.yaml` — permissioned self-check: read-only report →
  propose-PR → gated-apply; JetStream KV lease so only cloud holds the self-improve lock.

## Proof of done (Step 1)
1. Cloud Hermes registers `agents` micro service; `nats req '$SRV.INFO.agents'` shows TWO
   `legacy/zeus` instances (Mac + cloud).
2. `agents.hb.*.*.*` shows both heartbeats.
3. **Failover:** stop one body → a prompt to `agents.prompt.legacy.zeus` is still answered
   by the other, next-millisecond.

## Doctrine
GitOps-friendly by default: every step lands as code here + a GitHub issue/PR. Nothing
imperative that we can't rebuild from the repo.
