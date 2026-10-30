---
name: omos-legacy-genome
description: "Use to update Legacy's own definition (soul, memory seeds, doctrine, skills, runbook) through Git, and to answer who-am-I, what-do-I-run and how-do-I-recover questions."
version: 2.0.0
author: Legacy
license: MIT
metadata:
  hermes:
    tags: [self-improvement, architecture, omos, identity, evolution, gitops]
    related_skills: [legacy-self-awareness, onemind-architecture, hermes-plugin-pattern]
---

# Legacy genome: self-programming through Git

Legacy's definition is **not** in a separate `omos-legacy` repository. That
repository is archived. Everything now lives in the single OneMind repo,
`onemindos/omos`, so the cluster, the agent and the docs share one history.
Corrected 2026-10-29 from the earlier version, which described the archived repo.

## Where the genome lives

| Piece | Path in `onemindos/omos` |
|---|---|
| Soul and persona (what to load into the live instance) | `platform/hermes/seed/onemind-hermes/SOUL.md` |
| Memory seeds (environment facts, owner facts) | `platform/hermes/seed/onemind-hermes/MEMORY.md`, `USER.md` |
| Instance manifests and initial config | `platform/hermes/onemind-hermes/` (`config.defaults.yaml`, `MIGRATION.md`) |
| Skills this instance loads | `platform/hermes/skills/` (see its `README.md`) |
| Doctrine (creed, pillars, profiles, evolution) | `platform/hermes/agent/doctrine/` |
| Role definitions (SO, HP, LE, GE) and the older Legacy | `platform/hermes/agent/agents/` and `AGENTS.yaml` |
| Runbook and architecture of the predecessor | `platform/hermes/agent/RUNBOOK.md`, `LEGACY-ARCHITECTURE.md`, `FABRIC-REFERENCE.md` |
| Platform mission, live inventory, subject grammar | `docs/MISSION.md`, `docs/STACK.md`, `docs/NATS-DOCTRINE.md` |

The doctrine files describe the older multi-agent fabric (overseers and
servants). The current instance is **one agent with one soul**; the old
servants are modes (`onemind-domains`). When doctrine and the soul disagree,
the soul and the owner win; propose a doctrine fix.

## When to use

- Update your own architecture notes, soul, memory seeds or a skill.
- A new service, machine or operating pattern should be recorded.
- A recovery procedure changed or was learned.
- Answering "who am I", "what do I run", "how do I recover".

## Self-improvement workflow (Git, never in place)

1. Decide what you learned and which file owns it (table below).
2. Edit on a branch of a checkout of `onemindos/omos`; never push to `main`.
3. Open a PR with a clear description. The owner reviews and merges (never
   squash). Argo CD deploys what lives under `platform/`.
4. After merge, refresh the live copy: skills are copied into
   `/opt/data/onemind-skills`; memory seeds go to `/opt/data/memories/` with the
   old file backed up under a dated suffix first (see `platform/hermes/skills/README.md`).

The pod has no GitHub credential by default. If you cannot push, write the
change as a patch or a PR description and hand it to the owner or to a coding
worker (`onemind-code-orchestration`). Never paste a token into a command, a
remote URL or a reply. CI (`quality.yml`) runs gitleaks on every PR.

| Discovery | File to update |
|---|---|
| New machine on the tailnet | `docs/STACK.md` (tiered inventory), `tailscale-pod-access` |
| New cluster service | an Argo Application under `platform/argocd/` plus a `docs/STACK.md` row |
| NATS subject, stream or KV | a new CR in `platform/nats/streams/`; subjects follow the locked grammar |
| New or changed agent behavior | `platform/hermes/seed/onemind-hermes/SOUL.md` |
| Environment or owner fact | `MEMORY.md` / `USER.md` seeds (and the live memory tool) |
| Recovery procedure | `platform/hermes/agent/RUNBOOK.md` or `docs/operations/` |
| Reusable method | a skill under `platform/hermes/skills/` |

## Pitfalls

1. Never commit credentials, even test values; refer to the Secret by name.
2. Always branch and PR. Stale docs are worse than no docs, but a wrong "live"
   claim is worse than a stale one: verify against the cluster before writing "live".
3. The owner may have newer uncommitted work on his Mac. Check for open PRs and
   ask before concluding that a file or directory does not exist.
4. Every skill needs `name` and `description`, and must pass the scrub rules in
   `platform/hermes/skills/README.md` (no tokens, passwords, chat ids, phone
   numbers or addresses that grant access).
5. If you learned it, commit it: do not let operational knowledge die in a session.

## Continuous improvement checklist (end of any complex session)

- [ ] New operational pattern discovered? Skill or doctrine PR.
- [ ] A skill found outdated or wrong? Patch it now, not later.
- [ ] A non-trivial problem solved (five or more tool calls)? Save it as a skill.
- [ ] Memory stale? Update it.

## Reference

- `references/dev-infra-audit-checklist.md`: making any OneMind repo dev-ready
  (frontend and IaC/GitOps profiles).
