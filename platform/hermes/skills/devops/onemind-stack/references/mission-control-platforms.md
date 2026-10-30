# Mission operations platforms: verify execution, not dashboard labels

Use when comparing Paperclip, builderz-labs Mission Control, AionUi, or similar consoles for TAK/robotics/live project work.

## Decision rules

- Treat the requirement as a control plane: dispatch work, observe a correlated live run, pause for an attributable scoped decision, resume only the approved action, and attach execution evidence. Do not substitute a task list, messenger, or cron-prompt UI.
- Separate canonical work/run state, human presentation, execution engine, and device/event adapters. Keep one authoritative job and approval decision; avoid a Paperclip task and Mission Control task becoming competing owners. Keeping Paperclip during a pilot is not proof that a second permanent control plane is necessary.
- Inspect runtime-specific code before crediting a console with Hermes support. Mission Control's current Hermes tasks route documents a read-only bridge to cron files, while its events route receives lifecycle hooks. Its exec-approvals route targets an OpenClaw gateway/home. Those facts do not establish an end-to-end Hermes execution/approval bridge. Recheck source and exercise the exact release before deployment.
- Paperclip's current Hermes adapter supports local CLI and an existing HTTP/SSE gateway. Its gateway path creates runs, streams events, polls on stream interruption, and can stop timed-out runs. Its reverse task bridge is separate, scoped access. Validate deployed adapter version and permission/approval behavior rather than assuming current source is installed.
- Inspect Paperclip's current Plugin SDK before recommending an external full dashboard. It documents pages, dashboard widgets, detail tabs, webhooks, managed routines, streams, approval events, and agent sessions. These are extension points, not a shipped TAK/NATS/robotics integration or proof of mobile fitness. A platform plugin extending an owned platform differs from asking a standalone agent plugin to be the platform.
- AionUi is a multi-agent workspace with remote WebUI and scheduled prompts. Its scheduler requires the process to remain running. Generic ACP support is not verified Hermes integration. Its source includes mobile code; do not infer publicly supported mobile delivery, reliable push, or unified mission management from source paths.
- Design proposed OneMind integration around one mission/run ID, authenticated event ingestion, idempotency, state reconciliation, and action-bound expiring approvals. High-rate telemetry stays on the fabric; send useful state changes to the work record. Physical-device safety/interlocks belong at the executor/controller as well as the UI.
- Use a physical-phone end-to-end demonstration as the acceptance gate: event -> visible work -> proposal -> approve/deny -> authorized execution -> live progress -> durable evidence, with desktop closed and after reconnect. Do not deploy or migrate from a research-only request.

## Control-plane architecture review

- Distinguish a suitable UI/framework stack from an implemented mission execution architecture. A diagram and test-tool list are not proof of end-to-end enforcement.
- Put validation, resource-scoped authorization, policy, and required approval before dispatch; recheck applicable authority and approval freshness at execution. Bind approval to action/target/parameters and actor, with expiry and revocation behavior.
- Route Web, CLI, MCP, REST, and webhook commands through one command service. Persist admission/outbox before dispatch and use durable workers rather than keeping mission lifetime inside an HTTP request or browser session.
- Propose JetStream-capable fabric subjects for durable commands/events; verify actual account and stream availability before implementation. Core NATS has no replay. Design for redelivery, application idempotency, reconciliation, and explicit distinction between bus acknowledgment and verified real-world completion.
- SQLite WAL permits concurrent readers with a single writer and requires same-host access, not a database file shared over a network filesystem. Keep it initially for single-writer control-plane state; reassess storage/HA when multiple active writers are required. Keep high-rate raw telemetry out of task/audit tables, use proper backups and test restores.
- Declare ownership of project metadata, mission intent, execution facts, approvals, and receipts. Treat the bus as the execution-event fabric and UI caches/projections as derived views, not independent competing masters.
- Mobile requires explicit PWA/push/reconnect behavior, stale-data indicators, authenticated deep links, and snapshot plus event resumption; WebSockets/SSE do not by themselves establish phone background delivery. Do not auto-replay expired physical commands after reconnect.
- Require fault-path tests for duplicate delivery, crash/restart, expired or revoked approval, denied alternate API paths, lost connectivity, and restore. Keep physical-device safety/interlocks independent of UI uptime. Keep Google login optional if independent self-hosted identity is required; Zod validates shape, not authorization.

Primary architecture references:
- https://www.sqlite.org/wal.html
- https://docs.nats.io/concepts/jetstream
- https://docs.nats.io/learn/jetstream/delivery-and-acknowledgment
- https://nextjs.org/docs/app/guides/self-hosting
- https://nextjs.org/docs/app/guides/progressive-web-apps
- https://cheatsheetseries.owasp.org/cheatsheets/Authorization_Cheat_Sheet.html

## Fork evaluation

- Evaluate an owned fork differently from adopting a product unchanged. Prefer the base whose authority model already matches the intended Hermes/NATS architecture, not simply the one with more panels. Keep one task/mission owner and preserve current work until migration acceptance passes.
- SwarmClaw does support Hermes through an API endpoint. Its provider sets explicitly mark Hermes as runtime-managed (own tool loop), native-capability, and worker-only (no coordinator/heartbeat/advanced settings in the stock model). Do not claim it necessarily replaces Hermes, but inspect those policy boundaries before proposing it as the Mother/control-plane shell. Source: `src/lib/provider-sets.ts` and `src/lib/providers/index.ts`.
- SwarmClaw also supplies its own orchestration, schedules, memory, tools, and UI extensions. These are useful if intentionally selecting that orchestrator; otherwise resolve duplicated authority instead of silently running two brains.
- Mission Control describes itself as alpha and as a control plane rather than an agent framework. Its Hermes execution/approval gaps still need a real adapter. Alpha APIs/schema churn and mobile behavior are fork-maintenance risks, not solved by rebranding.
- Both repositories currently use MIT licenses. Verify each fork's pinned revision and dependency/asset terms, preserve original notices, isolate OneMind domain adapters and mobile views, track upstream, and exercise upgrades plus restores. An opinion request is not authorization to fork repositories, publish code, or deploy.

## Primary sources

- https://github.com/swarmclawai/swarmclaw
- https://github.com/swarmclawai/swarmclaw/blob/main/LICENSE
- https://github.com/swarmclawai/swarmclaw/blob/main/src/lib/provider-sets.ts
- https://github.com/swarmclawai/swarmclaw/blob/main/src/lib/providers/index.ts

- https://github.com/builderz-labs/mission-control
- https://github.com/builderz-labs/mission-control/blob/main/src/app/api/hermes/tasks/route.ts
- https://github.com/builderz-labs/mission-control/blob/main/src/app/api/hermes/events/route.ts
- https://github.com/builderz-labs/mission-control/blob/main/src/app/api/exec-approvals/route.ts
- https://github.com/paperclipai/paperclip/blob/master/packages/adapters/hermes/README.md
- https://github.com/paperclipai/paperclip/blob/master/packages/plugins/sdk/README.md
- https://github.com/paperclipai/paperclip/blob/master/doc/plugins/PLUGIN_AUTHORING_GUIDE.md
- https://www.aionui.com/en/#automation
- https://github.com/iOfficeAI/AionUi/blob/main/docs/guides/webui.md
- https://github.com/iOfficeAI/AionUi/wiki/Scheduled-Tasks-Guide
