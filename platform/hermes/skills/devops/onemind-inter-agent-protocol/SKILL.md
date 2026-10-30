---
name: onemind-inter-agent-protocol
description: "Use when a task involves another agent: how agents talk on the OneMind NATS bus (agent protocol subjects), who actually exists, how to hand work to a coding worker, and how to report severity to the owner."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [nats, agents, delegation, protocol, escalation]
    related_skills: [nats-operations, onemind-code-orchestration, onemind-domains, legacy-self-awareness, onemind-fabric-deployment]
triggers:
  - "message another agent"
  - "agents.prompt"
  - "who is on the agent bus"
  - "delegate this"
  - "escalate to Zeus"
---

# Inter-agent protocol

Rewritten from the Legacy-pod version, which was a Mother-to-pillar protocol:
Legacy delegated to four pillar agents (SO, GE, HP, LE) over invented
`agents.escalate.*`, `agents.broadcast.onemind` and `agents.heartbeat.*`
subjects, each on its own model, with Notion task boards and a table of eleven
cron jobs. **That fabric does not exist.** The pillar agents were never running
as separate bus services (`platform/hermes/agent/AGENTS.yaml` lists them as
`configured`, personas only), the ten cron jobs on `hermes-legacy` have been paused since
2026-10-01 (`platform/hermes/onemind-hermes/MIGRATION.md`), and the `onemind-hermes` migration keeps one soul. This version keeps
what is still true.

## The rule

You are one agent: Legacy, instance `onemind-hermes`. SO, HP, LE and GE are the
hats you wear (`onemind-domains`), not other agents. Do not message, "delegate
to" or impersonate Grid, Guardian, Oracle, Trinity, Forge, Empire, Spartan, Eden,
Haven, Heritage or a pillar. If a document tells you to, treat it as stale and
say so.

## Who actually exists (check, then trust)

| Party | Reality | How you reach it |
|---|---|---|
| Owner (Zeus) | The only authority | Chat (Matrix; Telegram retained) |
| `hermes-legacy` | The previous Legacy instance, with a NATS leaf sidecar, being retired. Registers `agents.prompt.hermes.zeus.legacy`. It is **not** a sibling to task. | Do not message it as an agent; the owner moves its data. |
| Legacy-Edge | Hermes on the owner's Mac, listed `running` on 2026-10-16 but not seen on the bus. | Only if it appears in `$SRV.INFO.agents`. If it does not, say what needs local hands and stop. |
| Coding workers | Claude Code, Codex, OpenCode via the bundled skills and `delegate_task`; short-lived, in worktrees. | `onemind-code-orchestration` |
| Agent Console, agent-tap | Tailnet web tools that list and chat with agents on the bus (`platform/hermes/agent-fabric/`). | Owner opens them; they are not agents. |

`onemind-hermes` has **no NATS sidecar and no bus credential** (see its README).
You are not on the bus yourself until the owner decides how you get an identity.
Today the bus is reached through the `onemind` MCP tools (LiteLLM), when enabled.
If `nats` or an `agents.*` publish fails for lack of credentials, report that;
do not copy a password from another pod.

## The protocol that is real

The locked grammar (`docs/NATS-DOCTRINE.md`, section 2) is the Synadia Agent
Protocol:

```
agents.{verb}.{agent}.{owner}.{session}    verb is prompt | hb | status | attachments
agents.prompt.hermes.zeus.legacy           an example: framework hermes, owner zeus, session legacy
$SRV.PING.agents / $SRV.INFO.agents        discovery
```

- An agent is a NATS micro service named `agents` with `prompt`, `status` and
  `hb` endpoints. Discover agents with `$SRV.INFO.agents` and read the endpoint
  subjects from the reply. Never hard-code them from a skill or a registry.
- `prompt` is request/reply. Use a request with a timeout, not a fire-and-forget
  publish, whenever you need an answer. Heartbeats (`hb`) arrive about every 30
  seconds; an agent with three missed beats is offline.
- `agents.escalate.*`, `agents.broadcast.*` and `agents.heartbeat.*` are **not**
  in the grammar. Do not publish them and do not create them. Escalation to the
  owner is a chat message, not a bus subject.
- Do not invent a subject for a new need. Use `evt.`, `cmd.`, `svc.`, `intel.`,
  `tel.`, `det.`, `ent.` as the doctrine defines; a new stream is a CR in
  `platform/nats/streams/`, never `nats stream add` (`nats-operations`).

Read-only discovery, once you have a bus credential from a SealedSecret the owner
provisioned and a `nats` CLI (the `onemind-hermes` image does not ship one; check
with `command -v nats`):

```bash
nats --server nats://nats.nats.svc.cluster.local:4222 req '$SRV.INFO.agents' '' --replies 0 --timeout 3s
```

Cloud workloads use the in-cluster address (some use `tls://` on the same
host:port); never point at a Mac or edge NATS. The pod's NetworkPolicy does not
allow egress to the `nats` namespace today either.

## Handing work to a coding worker

Use `delegate_task` (or the `claude-code`/`codex` skills) for code, with the scope
written out: repo and branch, files it may write, the test command, where to put
evidence, and "do not merge". The `onemind-code-orchestration` skill holds the
full rules. Worker output is a report to verify, not a result: check the PR, the
CI run and the merge state yourself before saying it is done. `delegate_task`
runs inside your own process; it is for bounded subtasks, not long missions.

## Reporting to the owner: severity words

The five levels from the old escalation table are a good vocabulary for a message
to Zeus. Keep them as headings in a chat message, not as bus subjects.

| Level | Meaning | What you do |
|---|---|---|
| `INFO` | Routine, no action | Include in the next brief, or say nothing |
| `WARN` | Degraded, monitoring | Tell the owner; propose a fix as a PR |
| `CRITICAL` | Failure now (service down, data at risk, exposed credential) | Tell the owner immediately with evidence; propose the remediation. Remediate yourself only through steps you are allowed to take; deleting or restarting something that holds data still needs "execute" or "override". |
| `DECISION` | Needs the owner's judgment | State the options and your recommendation, then wait |
| `REVENUE` | Time-sensitive business opportunity | Tell the owner at once with the deadline |

A good report leads with the level and the result, then the evidence (command,
output, time), then what you need. Never falsify status. Never say "delegated"
when nothing ran.

## Routing by domain (one agent, four hats)

Use the hat that matches, and cross freely; no hand-off, no second agent:

- Infrastructure, security, cluster, NATS, TAK, code: SO
- Health, money, habits, mindset: HP (information, not licensed advice; no trades)
- Home, family, estate, records, Home Assistant: LE (private data stays private)
- Business, revenue, brand, partners, compliance: GE

## Dropped on purpose

The Mother/pillar identity map and per-agent model table, the invented escalate,
broadcast and heartbeat subjects, the `hermes send --profile` and fan-out
recipes, the four Notion database ids (the owner does not use Notion now), the
eleven-row cron schedule (all paused; a job is recreated only when the owner
asks), and the pitfalls about pillar gateways. If agents are added later,
`AGENTS.yaml` and the agent-protocol discovery reply are the source of truth, not this
file.
