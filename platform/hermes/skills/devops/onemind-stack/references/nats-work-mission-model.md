# NATS-centered work and mission model

Use when designing OneMind's owned, phone-first control plane to replace Todoist and route work to people, agents, automation, and robots.

## SDK versus application

- Recheck the official Synadia Agent Protocol and SDK sources before choosing packages or assuming protocol compatibility. Caller and host SDKs exist in TypeScript/Python; a standalone Hermes NATS gateway is documented. Verify installed Hermes, gateway, SDK and protocol versions through an actual safe round trip before claiming interoperability.
- Use the SDK for agent discovery, prompting, streamed activity, liveness, and supported mid-stream queries. Do not treat a prompt stream as a durable task/mission ledger. The inspected protocol v0.3 draft explicitly excludes JetStream-backed at-least-once streaming.
- The inspected standalone Hermes NATS gateway README documents MVP omissions: no JetStream at-least-once delivery, cron proactive delivery or send_message routing to NATS, running-agent /stop interruption, or same-process multi-session multiplexing. It also warns that simultaneous subagent approvals route FIFO/oldest-first and can mismatch commands. Treat these as current documented upstream limits, not verified facts about the installed gateway. Do not rely on that query path for physical mission authorization; implement command-bound domain approvals and cancellation in the mission/device executor and verify exact identity/correlation on real runs. Audit its documented NATS_CONTEXT credential/TLS propagation caveat before configuring a production connection.
- The official agent-web-ui is a Bun/Vue test client and useful reference for discovery, tool cards and mid-stream allow/deny. Its server holds the NATS connection; it is not evidence of a production-ready phone task application.

## Domain model

- Model areas of responsibility, optional projects, independent tasks, reusable routines, generated occurrences, execution runs, approvals and evidence. A task can have zero or multiple runs. Preserve Todoist-style inbox, quick capture, due dates, reminders, priorities, labels, subtasks, recurrence, sharing, and personal views.
- Separate due dates from execution schedules. Define timezone/DST, calendar-based versus completion-based recurrence, missed-run behavior, and deduplication of occurrences. Ordinary human tasks must not require an AI run or paid model call.
- Assign work to a human, Hermes/other agent, deterministic automation, or typed device/robot adapter. A robot need not run an LLM or implement a prompt endpoint.
- Route robotics through bounded typed goals and explicit capability/policy checks. Keep servo/control loops, local interlocks and emergency handling at the device. Bind approval to action, target, parameters, actor, expiry and command identity; do not replay expired physical commands on reconnect.

## Authority and delivery

- Center command/event integration on the NATS fabric; keep the UI/fork replaceable. Publish lifecycle and receipt events with stable task, occurrence, mission and run identities. Define one authoritative work/mission state model rather than letting Paperclip, Hermes and a fork independently own copies.
- Use JetStream-capable subjects/accounts for durable work/event delivery, with application idempotency, deduplication and reconciliation. Verify real account capabilities; Core NATS cannot stand in for a durable queue. Do not assume Synadia prompt traffic is persisted merely because JetStream exists elsewhere.
- Use SQL for domain/query needs and projections with an explicit admission/outbox or event-log consistency design; do not make every task-list query a scan of broker messages. Store bulky evidence separately and carry references. Keep high-rate telemetry distinct from workflow transition events.
- Scope access per human/device/agent and resource; areas are organizing categories, not automatically isolated tenants. NATS subject permissions do not inspect command parameters or establish application-level approval authority. Keep broad bus credentials out of phone clients.
- Start as logical modules in one backend plus durable workers, not a separate microservice for every noun. A Mission Control fork is a client/control surface; it should use the work/mission API rather than own the fabric's runtime lifecycle.

## Acceptance

Prove both an ordinary recurring human responsibility and one safe, bounded field workflow: capture/trigger -> assignment -> plan -> required approval -> execution -> map/live progress -> evidence -> next occurrence. Test phone background/reconnect, service restart, redelivery, revocation, expiry, backup/restore and private-user visibility before migrating existing work. A receipt of delivery is not completion or verification.

## Primary sources

- https://github.com/synadia-ai/synadia-agents
- https://github.com/synadia-ai/synadia-agent-sdk-docs/blob/main/core-protocol.md
- https://github.com/synadia-ai/hermes-nats-gateway
- https://github.com/synadia-ai/synadia-agents/tree/main/examples/agent-web-ui
- https://docs.nats.io/concepts/jetstream
- https://docs.nats.io/learn/security/authorization
