# InferMesh

**Repo:** github.com/kamalgs/infermesh  
**Language:** Go  
**Role in OneMind:** STAGED — Phase 3, LLM inference mesh router

Dynamic geographic load shifting for LLM inference, powered by NATS global clustering. GPU nodes, cloud API adapters, and clients connect to a NATS mesh. Routing shifts inference load across regions based on capacity, latency, cost, and availability.

## Why it matters for OneMind

- OpenAI-compatible HTTP proxy — swap `baseURL`, zero code changes in agent callers
- Elastic: adding GPU capacity = starting a NATS subscriber, no reconfiguration
- NATS leaf nodes for on-prem GPU clusters — no public IPs, no VPNs
- Multi-tenancy via NATS accounts

## Activation condition

Staged until a second inference node is available (VPS GPU, second Mac, etc.). Running load balancing across a single LiteLLM node has no value. Activate at Phase 3.

## Integration point

Sits between NATS INTERFACE account and any agent that calls an LLM. Agents publish to `infer.request.*` and InferMesh routes to the best available node.
