# Hermes-first operations workspaces

Use this reference when comparing project systems, agent dashboards, communications apps, and mobile frontends. Reverify product capabilities and edition gates against current primary sources; the candidates below are evaluation paths, not deployment facts or approved selections.

## Discovery and comparison sequence

1. Inventory the active Hermes profile's enabled plugins and current task owner. Read the installed plugin README, security guidance, and supported host contracts; a catalog listing alone does not establish runtime readiness.
2. Build a capability table with columns for human projects/tasks, recurring routines, agent assignment/execution, approval/audit, mobile interaction, multi-user authorization, export, and operational burden.
3. Trace one task from capture to scheduling, execution, human review, and durable evidence. Name the authoritative store and task ID at each boundary. Use links or explicit mappings when integrating systems rather than silently creating independently editable copies.
4. Check current official product documentation for hosting, licensing, authentication, MCP scope, and edition restrictions. MCP availability does not establish that the free/self-hosted edition includes it, or that the exposed tools cover the intended workflow.
5. Verify candidate behavior in a bounded pilot before making a migration recommendation.

## Native and adjacent candidates

| Candidate | What to investigate | Boundary to preserve |
|---|---|---|
| Hermes Projects / Kanban / cron | Project/task organization, persistent workers, recurring execution and delivery | Agent work completion is not proof a human or physical-world task happened. |
| Project Stewardship | Goals, workflow progression, project memory, evidence, proposals, and human approval using native projects/boards | Verify evaluator support for non-code evidence; repository/test checks alone do not establish farm or household outcomes. Browser support may be a subset of desktop. |
| Hermes Workbench | Browser access to the desktop interface on the same backend | It is a community plugin, not evidence of an official native phone app or a multi-user workspace. Validate its current security boundary and phone ergonomics. |
| Hermes Mobile adapter | Direct gateway connection, media, streaming, clarification/approval cards, device delivery | An adapter and SDK are not a finished mobile client. Check identity/device binding, TLS, reconnect behavior, and actual push/background delivery. |
| LibreChat with Hermes API | General-purpose chat UX backed by Hermes's OpenAI-compatible API | Compatibility is not native project UI, shared history, or approval parity. Determine whether LibreChat or Hermes executes tools; avoid introducing a second agent brain by accident. |
| OpenProject with Nextcloud | Human project governance and file linkage | Verify exact integration scope and current MCP edition gate. File integration does not imply a unified calendar, task database, or Hermes session system. |
| Paperclip / agent control dashboards | Existing projects, tasks, schedules, delegation, and oversight worth preserving | Compare actual features rather than rejecting a tool solely for business terminology. Resolve ambiguous dashboard names to the exact repository/product before comparing runtimes. |
| Matrix / Mattermost / Lark / Signal | Communications, notifications, and bounded interaction with Hermes | Do not assume communications infrastructure supplies project orchestration or an operations UI. Separate self-hosting, encryption, mobile push dependencies, and interactive bot support. |

Primary documentation entry points:
- Hermes desktop: https://hermes-agent.nousresearch.com/docs/user-guide/desktop
- Hermes API server: https://hermes-agent.nousresearch.com/docs/user-guide/messaging/api-server
- Hermes Kanban: https://hermes-agent.nousresearch.com/docs/user-guide/features/kanban
- Project Stewardship catalog: https://hermes-agent.nousresearch.com/docs/plugins/hermes-project-stewardship
- Workbench catalog: https://hermes-agent.nousresearch.com/docs/plugins/hermes-workbench
- Mobile adapter catalog: https://hermes-agent.nousresearch.com/docs/plugins/hermes-mobile
- OpenProject MCP: https://www.openproject.org/docs/system-admin-guide/integrations/mcp-server/
- OpenProject Nextcloud integration: https://www.openproject.org/docs/system-admin-guide/integrations/nextcloud/
- LibreChat custom endpoints: https://www.librechat.ai/docs/quick_start/custom_endpoints

## Pilot acceptance matrix

- **Homestead:** recurring equipment or greenhouse check, assigned person, photo/sensor evidence, escalation, and an explicit distinction between reminder sent and inspection performed.
- **Family:** private parent planning, child-visible chores with restricted permissions, shared scheduling, and a parent-only approval path. Reject shared administrator credentials.
- **Personal:** an ongoing routine with recurrence, rescheduling, and review without generating unnecessary autonomous agent work.
- **Business:** a finite project with milestones, dependencies, accountable humans/agents, and a reviewable completion record.
- **Mobile:** capture a task and attachment, inspect a real tool run, answer an approval, disconnect/reconnect, and resume the intended session from another device. Test background notifications separately from foreground WebSocket delivery.
- **Recovery:** export records and attachments, restore them, retain identifiers/links, and demonstrate how the pilot can be abandoned without losing or duplicating work.

Report each result as verified, documented but untested, missing, or out of scope. Do not promote installed/enabled plugins or successful chat responses to an end-to-end readiness claim.
