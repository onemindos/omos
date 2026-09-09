# NATS Ecosystem References

> Reference for Chunks 3, 4, 5 — upstream URLs preserved before fork/local deletion.

---

## CHUNK 3 — Synadia Official Tooling

Infrastructure you *run*, not code you *own*. Install via Homebrew, Docker, or binary.

| Repo | Upstream URL | Description |
|---|---|---|
| nats-server | https://github.com/nats-io/nats-server | High-performance NATS broker |
| natscli | https://github.com/nats-io/natscli | NATS CLI — pub/sub/stream management |
| nsc | https://github.com/nats-io/nsc | NKey/JWT config — operator/account/user provisioning |
| nats-surveyor | https://github.com/nats-io/nats-surveyor | NATS monitoring simplified |
| nats-top | https://github.com/nats-io/nats-top | top-like NATS server monitor |
| nats-operator | https://github.com/nats-io/nats-operator | NATS JetStream Kubernetes Operator |
| nack | https://github.com/nats-io/nack | NATS Controllers for Kubernetes |
| nex | https://github.com/nats-io/nex | NATS execution engine — run workloads on NATS |
| nats-memory-server | https://github.com/nats-io/nats-memory-server | In-memory NATS server for testing (Node.js) |
| nats-agent-edge-sensor | https://github.com/synadia-io/nats-agent-edge-sensor | LLM agent + ESP32 sensor demo |
| terraform-provider-jetstream | https://github.com/nats-io/terraform-provider-jetstream | Terraform provider for JetStream |
| client-tools | https://github.com/nats-io/client-tools | NATS client installers, build systems, CDN |
| nats-architecture-and-design | https://github.com/nats-io/nats-architecture-and-design | Official NATS architecture and design docs |
| nats-by-example | https://github.com/nats-io/nats-by-example | Runnable reference examples for NATS |

---

## CHUNK 4 — Synadia Agent Platform

The AI↔NATS wiring layer. **Integrated into omos-core/integrations/**.

| Repo | Upstream URL | Description | Status |
|---|---|---|---|
| synadia-agents | https://github.com/synadia-io/synadia-agents | SDKs + plugins — puts Hermes, Claude Code, OpenClaw on NATS as microservices | ✅ Integrated |
| synadia-agents-labs | https://github.com/synadia-io/synadia-agents-labs | Synadia Agent Protocol courseware + examples | ✅ Integrated |
| synadia-agent-sdk-docs | https://github.com/synadia-io/synadia-agent-sdk-docs | Wire protocol spec for Synadia Agent Protocol | ✅ Integrated |
| agenttap | https://github.com/M64GitHub/agenttap | Wireshark for AI agents — live NATS agent traffic visualizer | ✅ Integrated |
| omos-agenttap | https://github.com/onemindos/omos-agenttap | One Mind fork of agenttap — custom visualizer for omos fabric | Own repo |

---

## CHUNK 5 — Orbit SDK

NATS utility extensions. Import as a dependency — do not fork.

| Repo | Upstream URL | Description |
|---|---|---|
| orbit.go | https://github.com/nats-io/orbit.go | Go NATS utilities: JetStream ext, context, sys client |
| orbit.c | https://github.com/nats-io/orbit.c | C NATS utility extensions |
| orbit.js | https://github.com/nats-io/orbit.js | JavaScript NATS utility extensions |
| orbit.java | https://github.com/nats-io/orbit.java | Java NATS utility extensions |
| orbit.rs | https://github.com/nats-io/orbit.rs | Rust NATS utility extensions |
| orbit.net | https://github.com/nats-io/orbit.net | .NET NATS utilities: JetStream publisher, partitioned consumers, scheduled messages |
