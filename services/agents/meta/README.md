# omos-meta — the Meta Agent (conductor)

The **meta agent** is Legacy's hands over the whole agent swarm. It speaks the
[Synadia Agent Protocol](https://nats.io/blog/nats-native-protocol-for-ai-agents/) as a
**caller**: it discovers every agent on the NATS fabric, watches their liveness, prompts any
of them, and can fan a single prompt out to all of them and merge the answers.

Part of the **Legacy Agent Fabric** — see `omos-infra/docs/legacy-agent-fabric.md`.

## What it does (protocol-native, no registry)
- **Discover** — `$SRV.INFO.agents` lists every live agent + its endpoints. No registry, no
  coordinator. (Agents register as the `agents` NATS micro service.)
- **Liveness board** — subscribes `agents.hb.*.*.*`; an instance is offline after 3 missed
  heartbeats (per its advertised cadence).
- **Prompt** — sends a prompt to a specific agent's discovered `prompt` endpoint subject
  (default shape `agents.prompt.{agent}.{owner}.{name}`).
- **Fan-out** — prompt every discovered agent at once and collect the responses.

## Why raw NATS + optional SDK
Discovery and heartbeats here use NATS micro directly (spec-accurate, dependency-light). For
richer streaming/session semantics you can drop in the official caller SDK
(`pip install synadia-ai-agents`) — the wire is identical, so this scaffold and the SDK
interoperate. See `github.com/synadia-ai/synadia-agents`.

## Run
```bash
pip install -r requirements.txt
export NATS_URL=nats://localhost:4222        # or tls://... for the cluster bus
python meta_agent.py list                    # discover agents
python meta_agent.py watch                    # live heartbeat board
python meta_agent.py prompt --agent legacy --owner zeus --name legacy "status?"
python meta_agent.py fanout "one-line health check please"
```

## Roadmap (fabric steps 2–4)
- Merge/summarize fan-out responses (Legacy can call this as a tool).
- Expose the meta agent itself as an `agents`-protocol service (host SDK) so Legacy prompts
  it like any other agent → recursive orchestration.
- Per-agent zero-trust identity via `nauth` JWT CRDs (who may prompt whom).
- JetStream KV lease for the single self-improve holder.
