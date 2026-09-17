---
title: "Matrix"
type: reference
status: planned
tags: [platform, communication, decentralized, encrypted, self-hosted]
---
# Matrix

> Decentralized, end-to-end encrypted communication protocol. Self-hostable, federated, open standard.

## Use Case in OneMind OS
- Sovereign comms — no dependency on corporate platforms
- Agent-to-agent messaging (machine-readable, structured)
- Fleet communication backbone (alternative/complement to NATS for human-facing)
- Bridging to other platforms (Telegram, Discord, SMS)

## Key Concepts
| Concept | Description |
|---------|-------------|
| Homeserver | Self-hosted instance (Synapse, Dendrite, Conduit) |
| Rooms | Encrypted group chats / channels |
| Spaces | Hierarchical org of rooms (like Discord servers) |
| Federation | Homeservers talk to each other (optional) |
| Bridges | Connect to Telegram, Discord, SMS, etc. |

## Clients
- **Element** — Primary GUI client (web, desktop, mobile)
- **matrix-commander** — CLI for automation/scripts
- **matrix-nio** — Python SDK for bots/agents

## Self-Hosting Options
| Server | Language | Notes |
|--------|----------|-------|
| Synapse | Python | Reference impl, feature-complete, heavy |
| Dendrite | Go | Lighter, newer, still maturing |
| Conduit | Rust | Minimal, fast, single-binary |

## Integration Points
- Hermes has native Matrix support (`send_message target='matrix:...'`)
- Bridges via mautrix (Telegram, Discord, Signal, SMS)
- Webhooks for alerting
