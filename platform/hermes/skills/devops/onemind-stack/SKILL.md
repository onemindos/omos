---
name: onemind-stack
description: "Use for OneMind stack decisions: which tool owns tasks, files, mail, community, chat or the operations workspace, how to evaluate a replacement, and the routing rules that keep records portable."
version: 2.0.0
author: Legacy
license: MIT
metadata:
  hermes:
    tags: [onemind, stack, evaluation, routing, vault]
    related_skills: [onemind-architecture, onemind-domains, composio, tailscale-pod-access]
---

# OneMind stack: decisions and routing

Imported 2026-10-29 from the Legacy pod and rewritten. The original was a July
2026 decision log. It named Dropbox, ClickUp, Zoho, Circle, Notion, a WordPress
test site with an application password, cutover IPs and a Google Workspace
migration order. Those are **history, not the live stack**: the owner has dropped
Notion for now, the WordPress credential was a secret (rotate it, it is exposed on
the old volume) and the rest was never verified here. What remains is the method.

## What is live (ask `docs/STACK.md`, not this skill)

The tiered inventory is generated from the cluster. At the time of writing it
includes Matrix (tailnet only), Nextcloud, Karakeep, n8n, Backstage,
Paperclip (preparation), Authentik, Grafana and the NATS UIs, with TAK as the
tactical display and NATS as the bus. A tool in an old decision log is a
candidate until `docs/STACK.md` or the cluster says it runs.

Before claiming a tool is "the" system for tasks, files, mail or community,
check (1) the current request, (2) `docs/STACK.md`, (3) whether the connector
exists (`composio`), and (4) whether the owner said otherwise. A previous
preference is not a mandate; an explicit request to compare alternatives wins.

## Routing law (kept)

1. **Tasks organize, they never host.** A task holds a pointer to where the
   bytes live; one copy, one path.
2. **Wrapper-note pattern.** A heavy file (PDF, video, bill) gets a thin markdown
   note with frontmatter (`type`, `vendor`, `amount`, `due`, `status`, `file:`),
   so bytes become queryable and schedulable. Bytes on the drive, meaning in markdown.
3. **Public pages live on the owner's own site** (git plus a static host), not on
   a third-party workspace that may be cancelled.
4. **The tree is the asset; every tool is fungible.** The folder structure plus
   history outlives any vendor; apps are lenses over it.
5. **One authoritative project and task platform per work scope.** Do not let two
   tools each own an editable copy of the same task.
6. **Sovereign and local-first.** Prefer open-source, self-hosted capability; do
   not make the home depend on a remote dashboard or identity provider.

## Evaluating a replacement or a new platform

Use this when asked "should we adopt X" or "compare A and B".

0. **Audit what is already deployed first.** Inspect live workloads and Argo
   provenance, then the existing portal (Backstage), core, agent-fabric console,
   Paperclip and workflow services. Reuse a working backend; judge the frontend
   separately. Never substitute archived manifests or remembered topology for
   live evidence (`references/existing-platform-discovery.md`).
1. **Rank against outcomes, not product names.** The central need is a
   phone-first control plane for projects tied to live runs, map and asset
   context, events and webhooks, visible execution, scoped approvals, evidence
   and self-hosted records. A generic task manager is not a substitute.
2. **Separate the layers:** durable human work record, agent execution and policy,
   communications, operational UI. A chat client, an MCP endpoint and a project
   manager solve different problems. Keep Hermes as the execution engine unless
   the owner picks another, and name any duplicated task store, scheduler or
   approval authority before recommending an integration.
3. **Distinguish a platform from its plugins.** Separate: documented capability,
   dependency present, backend registered, enabled in config, live ingestion,
   and an exercised acceptance test. None proves the next.
4. **Test mobile as an operations surface:** navigation, editing, scheduling,
   attachments and voice, tool progress, approval cards, reconnect and
   notifications, not just chat. API compatibility is not session continuity.
5. **Apply household and sovereignty gates:** separate identities and permissions
   for adults, children, guests and agents; hosting, backup and export;
   authentication; notification dependencies; paid feature gates. Check licence
   editions component by component (source-available is not open source) and
   re-verify volatile details against primary sources before deployment
   (`references/openproject-edition-boundaries.md`).
6. **Recommend one coherent route with a bounded pilot:** what stays
   authoritative, what is only a lens, what is missing, and what wins if the
   pilot fails. When a recommendation changes, name the new verified fact and the
   trade-off it moves. Research and discussion authorize no install or migration.
   Do not retire a working system until its replacement passes the required
   workflows and has a rollback.

References for the above:

- `references/portal-vs-app-selection.md`: developer portals vs internal-app builders, headless frameworks, agent context graphs.
- `references/nats-work-mission-model.md`: a NATS-centered replacement for personal tasks plus agent and robot missions.
- `references/mission-control-platforms.md`: live TAK and robotics missions, unified work and approval consoles.
- `references/hermes-operations-workspaces.md`: capability and acceptance checklist.
- `references/openproject-edition-boundaries.md`: hosting vs paid-feature gates.
- `references/existing-platform-discovery.md`: inventory method (read-only, sanitized, never print Secret values).

## Pitfalls

- A cancelled background copy may already have completed and removed its source.
  Check the filesystem state (destination size against expectation) before
  reporting that nothing happened.
- Regenerable giants (documentation dumps, planet basemaps) are re-downloadable,
  but deleting anything still needs "execute" or "override" for that deletion.
  Quarantine first (move to a trash path), match byte counts second, report the
  path third. Never blind-delete; never delete a duplicate you have not
  byte-verified.
- Heavy file inventories change quickly. Treat sizes and counts from old notes as
  hints and re-measure.
- Do not record the owner's media, vault or drive layout in a skill; it is personal
  and goes stale.
