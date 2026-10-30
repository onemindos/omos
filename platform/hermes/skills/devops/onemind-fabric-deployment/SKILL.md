---
name: onemind-fabric-deployment
description: "Use when the owner asks to add or change something on the OneMind agent fabric: a skill, a scheduled job, an agent setting, an integration, or a new bus-visible agent. Says which path to use and how to do it through Git."
version: 2.0.0
author: Legacy
metadata:
  hermes:
    tags: [agents, fabric, deployment, skills, cron, gitops, hermes]
    related_skills: [onemind-inter-agent-protocol, hermes-deployment-topology, hermes-mcp-config, omos-legacy-genome, doks-cluster-ops, legacy-self-awareness]
triggers:
  - "add an agent"
  - "new skill for Legacy"
  - "schedule a job"
  - "change Hermes config"
  - "deploy to the fabric"
---

# Changing the agent fabric

Rewritten from the Legacy-pod version, which scaffolded **servant and overseer
profiles**: five Hermes profiles with different models, SOUL oaths, an `AGENTS.yaml`
with three tiers, Composio account ids baked into system prompts, a Mac desktop
delivery recipe and Notion-backed cron prompts. That model is gone. The
`onemind-hermes` migration keeps **one agent with one soul**, and its profile
directories were never real agents (`platform/hermes/onemind-hermes/MIGRATION.md`).
This version is the current procedure for the same kinds of request.

## First decision: what kind of change is it?

Ask "would the owner open a direct chat with it, and does it need its own
memory, cron loop and bus identity?" Almost always the answer is no.

| The request is really... | Do this | Section |
|---|---|---|
| New knowledge or a procedure | A skill | A |
| A repeating check or brief | A cron job | B |
| A setting (toolsets, provider, approvals) | Hermes config | C |
| An outside tool or SaaS | An MCP server, via LiteLLM | D |
| A separate always-on agent or service on the bus | A GitOps workload, owner approval first | E |
| "Another persona" (Grid, Oracle, a pillar...) | A mode of you; add to `onemind-domains` | A |

Do not create Hermes profiles for roles. Do not create sub-agents that are
described as siblings. If the owner explicitly wants a second agent, that is
case E and a design discussion, not a scaffold.

## A. Add or change a skill

Skills for this instance live in Git at `platform/hermes/skills/<category>/<name>/SKILL.md`
and are loaded from `/opt/data/onemind-skills` through `skills.external_dirs`.
The index and the review rules are in `platform/hermes/skills/README.md`.

1. Branch from `origin/main` (never commit to `main`); add or edit the skill.
2. Write the frontmatter (`name`, a one-line `description` that says when to use
   it, `triggers`) and keep facts you verified with their date. A skill must not
   carry a credential, a personal address or a chat id; say which Secret holds
   a secret.
3. Update the README index, run the repo checks, open a PR. The owner merges.
4. After the merge, copy the tree into the pod with the `tar` command in the
   README and confirm with `hermes skills list`. Do not copy before the merge.

A local skill wins when its name collides with another, so check
`hermes skills list` for a clash before choosing a name.

## B. Schedule a job

```bash
hermes cron create "<schedule>" "<self-contained prompt>" --name <name> --deliver <target>
hermes cron create "every 30m" --name <name> --no-agent --script <file> --deliver <target>   # watchdog, no LLM
hermes cron list
```

- Prefer `--no-agent` with a script for watchdogs: empty output is silent and
  costs no inference. Use an agent job only when judgment is needed. Use
  `--monitor-script`/`--monitor-url` to run the agent only when something changed.
- Open the prompt with who you are and the task ("You are Legacy, ...") and make
  it self-contained; a cron run has no chat history.
- Alerts-only for health checks, always-send for briefs. Do not add a watchdog for
  something that is not deployed; it will page falsely.
- Cron runs with `approvals.cron_mode: deny`. A job cannot run commands the
  approver flags, and `execute_code` may be refused in that mode; write jobs as
  terminal calls and scripts.
- The ten jobs on the old instance were paused on 2026-10-01 after `429
  tokens_per_day` errors, and `hermes cron list` on `onemind-hermes` is empty.
  Re-create a job only on the owner's request, and keep it cheap: do not spend
  inference credit on output nobody reads.
- The job store lives on the pod volume, not in Git. Record the exact `hermes
  cron create` line in the PR or a doc so the job can be rebuilt.
- Delivery targets, chat ids and tokens belong to the runtime Secret and
  config. Never paste one into a prompt or a skill.

## C. Change Hermes config

- `hermes config get|set|unset <dotted.key>`; direct edits of `config.yaml` and
  `.env` are blocked on purpose. `hermes config check` finds missing options.
- The init container seeds `platform/hermes/onemind-hermes/config.defaults.yaml`
  **once**, on an empty volume. Changing that file later does not change the
  running instance. For a lasting change, make it with `hermes config set` and
  also commit the same value to the defaults so a rebuild matches, and say so in
  the PR.
- The gateway cannot restart itself. A config change that needs a restart is a
  request to the owner (or a rollout through Git), not something to force from
  inside the pod. Verify the flag after the restart; one voice setting was seen
  reverting on the old instance.
- Provider and model credentials come from the runtime SealedSecret
  (`onemind-hermes-runtime`) and the owner's sign-in, not from a skill. No
  Anthropic providers on this instance; fallback is Ollama Cloud once its key is
  sealed.

## D. Add an outside tool

Tools reach the instance as MCP servers behind LiteLLM (`hermes-mcp-config`,
`platform/hermes/mcp/HUB.md`). Add the upstream to LiteLLM's `mcp_servers` in Git,
scope the key, then `hermes mcp ...` on the pod. Credentials enter as
SealedSecret env, referenced as `os.environ/NAME`; never as a header or URL
written in the skill or in `config.yaml`.

## E. A new agent or service on the bus

Treat it as a normal new workload (CLAUDE.md "Adding a New Service"):

- [ ] Image on GHCR with a pinned tag or digest; resource requests and limits.
- [ ] Manifests under `platform/hermes/` (or the right layer) and an Argo
      Application. No `kubectl apply`, no `helm install`.
- [ ] NATS identity: its own user with only the subjects it needs; password from a
      SealedSecret. The shared `onemind` login is a stopgap and the
      `agent_fabric` user is still pending. Cloud workloads connect to
      `nats://nats.nats.svc.cluster.local:4222`; the leaf-sidecar pattern is
      Hermes' own exception and is not copied.
- [ ] Agent protocol: register as the `agents` micro service (`prompt`, `status`,
      `hb`), session name from config, subjects from the locked grammar
      (`onemind-inter-agent-protocol`).
- [ ] An entry in `platform/hermes/agent/AGENTS.yaml` with an honest `status`
      (`running`, `configured` or `planned`; verified date), and the soul in Git.
- [ ] NetworkPolicy (default deny, named egress), backup decision (VolSync on the
      volume), `docs/STACK.md` row, and a tier.
- [ ] `onemind-hermes` has no bus credential of its own today. Giving it one is
      an owner decision (a leaf sidecar, or a direct user); propose it in a PR,
      do not work around it.

## Verify, then report

- `hermes skills list` shows the new skill as `enabled`; the count matches the
  README.
- `hermes cron list` shows the job with the expected schedule.
- For a workload: Argo Application `Synced`/`Healthy`, pods Ready, and the agent
  appears in `$SRV.INFO.agents` (read-only). A ready pod is not proof of a working
  chat; test a real exchange.
- State exactly what changed and where (repo path, PR, pod path). Never say "live"
  for something you only committed.

## Pitfalls that still apply

- A live change that Argo owns is reverted within about a minute; the PR is the
  fix (`legacy-self-awareness`). An annotation bump may be ignored by an
  Application's `ignoreDifferences`; an image digest change always rolls the pod.
- `git push` needs a remote that exists; many scratch repos on old pods were
  local-only. Check `git remote -v` and report honestly whether the push worked.
- Removing a profile, volume, branch or job is a deletion. Needs "execute" or
  "override" for that specific item; the old rule "delete dead profiles outright"
  is superseded.
- Say which machine and path you write to. The owner cannot see the pod's
  filesystem.

## Dropped on purpose

The five-profile model table and per-agent models; the servant SOUL "oath"
template and `AGENTS.yaml` tier structure; symlinking skills into profile
directories; the Composio account id matrix and system prompts that embedded
them; the Notion database ids and token handling; the Mac Desktop delivery recipe
(base64 over SSH) and the Telegram bot setup steps (the pod's egress and secrets
model changed); the Homepage dashboard (not deployed; there is no
`platform/homepage`); and the old `platform/legacy/` ignoreDifferences note (that
path is gone; `platform-legacy` now syncs `platform/hermes/{mcp,agent,bridge}`,
so read its Application before relying on an annotation bump). None of the
seven `references/` files were kept: SOUL template, `AGENTS.yaml` structure,
servant matrix, messaging platforms, Android scoped-storage verification, the
Buzz/ACP project note (research; it belongs in `docs/references/`) and the
Homepage notes.
