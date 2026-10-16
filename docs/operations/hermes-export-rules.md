# Hermes export rules: skills and config.yaml

**Status:** proposal for review. Nothing here is built or executed. This is the rulebook for Step 2 of `hermes-recoverability-plan.md` (the nightly export workflow, decision D3), written before the workflow so the rules can be agreed first.
**Written:** 2026-10-16. Everything below was measured on the live `hermes-legacy` pod (read-only). Secret values are deliberately not recorded anywhere in this document.

## Summary

1. **A safe export is not possible yet.** Both `config.yaml` and the agent-written skills contain real credentials. The export must fail closed on them, and they must be rotated and moved out of those files first (section 4).
2. **Custom skills can be told apart from bundled ones** using Hermes' own metadata plus a comparison with the skills shipped inside the image. 71 skills qualify for export, 18 need a decision, 75 are skipped (section 1).
3. **`config.yaml` is exported as a redacted template**, using a fixed list of paths plus generic rules, and fails if it finds a sensitive-looking value it has no rule for (section 2).

## 1. Telling custom skills from bundled ones

### What identifies a skill

| Signal | Where | What it tells us |
|--------|-------|------------------|
| `created_by` | `skills/.usage.json`, per skill | Hermes' own provenance. Values seen: `agent` (64 entries), `learn` (12), `installed` (3), null (117, older entries) |
| Image copy | `/opt/hermes/skills` (58) and `/opt/hermes/optional-skills` (150), inside the image | If a live skill is byte-identical to the image's copy, it is upstream |
| Hub install | `skills/.hub/lock.json` | Skills installed from the Hermes hub (2: `llama-cpp`, `serving-llms-vllm`) |
| `author:` frontmatter | top of `SKILL.md` | Soft signal: `Legacy`, `Legacy (Mother)`, `OneMind OS / Zeus` mean ours; `Hermes Agent` or a person's name means upstream or community |
| `patch_count` | `.usage.json` | How many times the agent edited the skill. Above 0 on an upstream skill means a local change worth keeping |
| `.bundled_manifest` | `skills/` | **Not usable on its own.** It lists only 44 skills while 164 are live |

### The rules (applied in this order; first match wins)

| # | Rule | Result |
|---|------|--------|
| R1 | Byte-identical to the image's copy | Skip (upstream) |
| R2 | In `.hub/lock.json`, or `created_by` is `installed` | Skip (hub) |
| R3 | `created_by` is `agent` or `learn` | **Export** |
| R4 | `author:` is `Legacy`, `Legacy (Mother)` or `OneMind OS / Zeus` | **Export** |
| R5 | In the image but changed, and `patch_count` is above 0 | **Export** (a local patch) |
| R6 | In the image but changed, `patch_count` 0 | Skip (the live copy is just an older version of the upstream skill) |
| R7 | Not in the image and no `author:` | **Export** (hand-made locally) |
| R8 | Anything else | **Review** (listed in the PR, not exported) |

Archived skills (`state: archived`, and `skills/.archive/`) are never exported; the backups still hold them.

### Result on 2026-10-16 (164 live skills)

- **Export: 71** (about 1.8 MB). They are the Legacy-written skills: the `so-*`, `le-*`, `hp-*` and `ge-*` servant skills, the `composio*` family, `omos-*`, `onemind-*`, `tak-*`, `mac-*`, `legacy-*`, `notion-*`, `nats-operations`, `ha-nats-bridge`, `zeus-site`, `wordpress-*`, and one patched upstream skill (`google-workspace`).
- **Skipped: 75:** 45 identical to the image, 27 older versions of upstream skills, 3 hub-installed.
- **Review: 18**, none patched: the `github-*` skills (6), `jupyter-live-kernel`, `merge-reconciler`, `ocr-and-documents`, `petdex`, `plan`, `session-librarian` (author "Hermes Agent", but not in this image's folders, so probably from an older image), and `audiocraft`, `lm-evaluation-harness`, `segment-anything`, `vllm`, `nano-pdf`, `touchdesigner-mcp` (community authors). Recommendation: leave them out; the workflow lists them weekly so nothing is missed.

The workflow recomputes this each night, so a skill the agent creates tomorrow is picked up under R3.

## 2. Redacting `config.yaml`

### What is in it

621 lines, 323 values, 20 MCP servers. Most of it is harmless settings. These values are not:

| Path | What it is |
|------|------------|
| `mcp.servers.notion.args` and `mcp_servers.notion.args` | A Notion API token inside the `OPENAPI_MCP_HEADERS` argument (the second block is disabled) |
| `mcp_servers.github.headers.Authorization` | A GitHub Copilot MCP bearer token |
| `mcp_servers.wordpress.headers.Authorization` | A WordPress Basic-auth header (disabled; retired domain) |
| `mcp_servers.composio.headers.x-consumer-api-key` | A Composio API key |
| `platforms.matrix.token` | A Matrix access token (and `platforms.matrix.user_id`, an identifier) |
| `dashboard.basic_auth.password_hash`, `.secret` | The dashboard login hash and a session secret |
| `TELEGRAM_ALLOWED_USERS` | A personal numeric Telegram ID (an identifier, not a secret) |

Already correct and safe to keep: `model.key_env` and `providers.digitalocean.key_env` (an env-variable name), `secrets.onepassword.env.*` (`op://` references, the right pattern), and `mcp_servers.*.auth: oauth`.

### The rules

1. **Parse, redact, re-serialise.** Load with a YAML parser, replace values, dump with the original key order. Never line-grep the raw file.
2. **Redact by path (explicit list):** every path in the table above becomes a placeholder `${NAME}`, where NAME comes from the path (for example `${MCP_GITHUB_AUTHORIZATION}`, `${PLATFORMS_MATRIX_TOKEN}`).
3. **Redact by name (generic):** any key whose name matches `token|secret|password|passwd|api[_-]?key|authorization|credential|cookie|bearer`, unless it is on the safe list: an env-variable name (`*_env`), an `op://` reference, a boolean, an empty value, or `auth: oauth`.
4. **Redact by value:** any string matching a known secret shape (`ntn_`, `ghp_`, `gho_`, `github_pat_`, `ck_`, `sk-`, `AKIA`, `tskey-`, a JWT, `Bearer ...`, `Basic ...`, a Telegram bot token) anywhere in the file, including inside list items and embedded JSON such as the Notion `args`.
5. **Identifiers are redacted too** (`TELEGRAM_ALLOWED_USERS`, `platforms.matrix.user_id`), since they are personal.
6. **Fail closed.** If a key name looks sensitive and matches none of the safe-list forms, the export stops and the PR says which path needs a rule. It never skips silently and never guesses.
7. **Never export `config.yaml.bak-*`** (at least nine dated copies, plus `.bak-bedrock` and `.pre-matrix-dev-*`). Older copies may hold older secrets.
8. After redaction, run gitleaks on the result. Any finding stops the export.
9. Profile configs under `profiles/*/config.yaml` follow the same rules.

The same checks (rules 4 and 8, plus a scan for credential-shaped assignments such as `password=...`) run over every exported skill file.

## 3. Output layout

```
platform/hermes/agent/
  config/config.template.yaml        redacted, with ${PLACEHOLDERS}
  profiles/<name>/config.template.yaml
  skills/<category>/<skill>/...      exported custom skills only
  EXPORT-REPORT.md                   what was exported, skipped, sent to review, and why
```

Only these paths are ever written by the workflow. It cannot touch manifests, workflows or anything else in the repo.

## 4. Secrets already in the skills and the config

The scan found credentials in the export set, so the first run would be refused. What was found (values not recorded):

- **Config:** the six items in section 2's table (Notion, GitHub Copilot, WordPress, Composio, Matrix, dashboard).
- **Skills (33 matches across about 23 files):** the Notion token again (`notion-workspace`, in full, not truncated), Composio API keys (`composio`, `composio-mcp-patterns`, `legacy-integrations`), Tailscale OAuth client credentials (a `doks-cluster-ops` reference), and passwords or tokens in `tak-server`, `tak-mission-ops`, `tak-mission-runner`, `zeus-site`, `legacy-self-awareness`, `omos-feed-builder` and `google-workspace/scripts/gws_bridge.py`. A few matches are placeholders or environment-variable references.
- This also changes an earlier assessment: the Notion token in the old docs was truncated, but the one on the pod is not.

**Prerequisite to Step 2** (tracked in Todoist): rotate each credential at its provider, replace inline values with env or `op://` references, edit the live skills to read from the environment, and only then run the export. These are live changes to the agent and are done with a fresh backup first.

## 5. Decisions needed

| # | Decision | Recommendation |
|---|----------|----------------|
| D9 | The 18 "review" skills | Leave them out of the export; the workflow lists them each week |
| D10 | IP addresses inside skills (public service addresses, droplet addresses, tailnet and LAN addresses; more than 200 occurrences) | Do not redact. The skills need them to work, and the repo is private. The one-off docs scrub in Step 0 was a different case. Alternative: map each address to a hostname |
| D11 | Order of work | Rotate and clean first, then build the workflow, so its first run succeeds instead of failing |
