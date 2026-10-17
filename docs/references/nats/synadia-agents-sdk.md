# Synadia Agent SDK — OneMind Agent Harness

**Repo:** https://github.com/synadia-io/agent-sdk (TypeScript + Python)
**Publisher:** Synadia
**Docs:** https://docs.synadia.com/agents

## What it is

The official SDK for building agents on the Synadia/NATS agent protocol. OneMind's agent harness (Hermes/Legacy) is already built on this — it's the foundation for all AI agents on the bus.

## Already in use in OneMind

- **Hermes/Legacy** runs as a Synadia agent on `agents.prompt.hermes.zeus.legacy`
- **Claude Code** is wired in as a Synadia agent
- All agents advertise capabilities, accept prompts, and return responses over NATS using this protocol

## SDKs available

| Language | Install |
|----------|---------|
| TypeScript | `npm install @synadia-io/agent-sdk` |
| Python | `pip install synadia-agent-sdk` |

## Key concepts

- **Agent session:** a stateful conversation context over NATS request/reply
- **Tool calling:** agents expose tools as NATS service endpoints
- **Persona/SOUL.md:** agent identity seeded from session name + SOUL.md file
- **Subject:** `agents.prompt.{harness}.{owner}.{name}` — standard addressing

## Notes

- This is the protocol Hermes, Legacy, and all future OneMind agents use
- `synadia-agents-labs` repo has tool-calling examples (TypeScript `responder.ts` pattern)
- Do not vendor — install via npm/pip from the official package
