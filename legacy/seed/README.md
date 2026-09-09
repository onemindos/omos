# omos-agents — Hermes on the NATS bus (Legacy host)

Makes a Hermes instance a **first-class agent on the fabric** via the
[Synadia Agent Protocol](https://nats.io/blog/nats-native-protocol-for-ai-agents/).
This is the piece `hermes gateway run` does **not** provide (that's chat-platform
only — Slack/Matrix/iMessage). Part of the **Legacy Agent Fabric**
(`omos-infra/docs/legacy-agent-fabric.md`).

## Architecture (per Zeus)
- **Hub (cloud/DOKS) = the main Hermes** → registers `agents.prompt.hermes.zeus.legacy`.
- **Mac = a separate Hermes** → its own identity on the fabric.
- The fabric interconnects them (distinct agents, not one queue-grouped identity).

## `hermes_nats_host.py`
Runs **inside** the Hermes container. Connects to the in-pod NATS leaf
(`localhost:4222`), registers `AgentService(agent=hermes, owner=zeus,
session=legacy)`, and on each prompt drives the local Hermes one-shot
(`hermes -z <prompt>`), streaming the reply back. Self-heals its Bedrock model
config on a fresh `HERMES_HOME`. Heartbeats every 30s.

Key env (see the script header for all): `NATS_URL`, `HERMES_NATS_OWNER`,
`HERMES_NATS_SESSION_NAME`, `HERMES_MODEL`, `HERMES_PROMPT_TIMEOUT`.

## `Dockerfile`
Thin derived image `FROM ghcr.io/onemindos/hermes:legacy-v1` that bakes in the
host + its deps (`nats-py`, `synadia-ai-agent-service`, `synadia-ai-agents` —
NOT part of Hermes's own deps). Produces `ghcr.io/onemindos/hermes-legacy:v1`.
Deployed by `omos-infra/platform/hermes/`.

## Verify
```bash
# from nats-box with the onemind context:
nats micro ls                       # expect an 'agents' instance, session=legacy
nats req agents.prompt.hermes.zeus.legacy 'sitrep?'
```

## Latency note / roadmap
`hermes -z` cold-starts per prompt (loads skills + model each call → tens of
seconds). Fine for proving the loop; upgrade path is to drive the long-lived
`hermes serve` JSON-RPC backend for low-latency, session-persistent prompting.
