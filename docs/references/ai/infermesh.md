# InferMesh — NATS-native LLM Inference Router

**Repo:** https://github.com/kamalgs/infermesh
**License:** MIT
**Language:** Go + TypeScript SDK

## What it is

Dynamic geographic load-shifting for LLM inference powered by NATS global clustering. GPU nodes, cloud API adapters, and clients all connect to a NATS mesh. The routing layer shifts inference load across regions based on capacity, latency, cost, and availability — in real time.

## Architecture

Two application components — everything else is NATS infrastructure:

```
Proxy (cmd/proxy)          Provider (cmd/provider)
HTTP → NATS                NATS → HTTP LLM
POST /v1/chat/completions  Anthropic / OpenAI / Ollama adapters
   │                            │
   └──── llm.chat.{model} ──────┘
              (NATS subject)
```

No router, orchestrator, or gateway code. The model name IS the NATS subject — the proxy publishes directly, NATS routes to the nearest available provider.

## Why it's relevant to OneMind

| Pattern | OneMind application |
|---------|---------------------|
| LLM calls over NATS | Route Hermes/Legacy inference requests to local GPU (Mac) or cloud (Bedrock/OpenAI) based on availability |
| Leaf node GPU clusters | On-prem edge GPU connects via outbound NATS leaf — no public IP, no VPN |
| Multi-tenant NATS accounts | Hard isolation between different agent tenants on the same bus |
| OpenAI-compatible HTTP proxy | Drop-in replacement — change `baseURL` only, zero code changes in callers |
| Provider adapters | Anthropic, OpenAI, Ollama — all three in use in OneMind already |

## Subject grammar

```
llm.chat.{model}        → route to provider serving that model
llm.chat.gpt-4o         → OpenAI
llm.chat.claude-opus-5  → Anthropic
llm.chat.llama3         → Ollama (local GPU)
```

## When to use this

If OneMind ever needs to route LLM inference across multiple providers or GPU nodes (e.g. Mac GPU + K8s GPU + Bedrock), InferMesh is the reference implementation. The NATS leaf node pattern for GPU clusters is directly applicable.

## Key files

| File | Purpose |
|------|---------|
| `internal/proxy/proxy.go` | HTTP → NATS bridge |
| `internal/provider/` | Anthropic / OpenAI / Ollama adapters |
| `internal/session/store.go` | JetStream KV session store |
| `demo/docker-compose.yaml` | Full local demo stack |
| `docs/TRANSPORT-NATS.md` | NATS transport design |
| `sdk/js/` | TypeScript client SDK |
