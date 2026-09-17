---
title: "Reference — Section Index"
section: "14-REFERENCE"
type: index
tags: [nats, reference, glossary, cheatsheet, cli, config, limits, subjects]
created: 2026-06-14
---

# 14 — Reference

The lookup layer. When you don't need a tutorial — you need the **exact flag**, the **default value**, the **monitoring endpoint**, or the **term definition** — it's here. This section consolidates and cross-links the practical details scattered across Sections 00–13.

See also: [[NATS/README|KB Home]] · [[NATS/09-CLI-TOOLING/README|CLI & Tooling]] · [[NATS/07-DEPLOYMENT-OPS/README|Deployment & Ops]]

## Notes
- [[NATS/14-REFERENCE/glossary|Glossary]] — every NATS term, one line each (subject, JetStream, consumer, leaf, gateway, NKEY, RTT…)
- [[NATS/14-REFERENCE/config-cheatsheet|Config Cheatsheet]] — server `.conf` blocks at a glance (jetstream, cluster, leafnodes, tls, auth, mqtt, websocket)
- [[NATS/14-REFERENCE/cli-cheatsheet|CLI Cheatsheet]] — the `nats` / `nsc` / `nk` commands you actually type
- [[NATS/14-REFERENCE/subject-design-patterns|Subject Design Patterns]] — naming conventions, wildcards, tokens, anti-patterns
- [[NATS/14-REFERENCE/monitoring-endpoints|Monitoring Endpoints]] — every `:8222` HTTP endpoint and what it returns
- [[NATS/14-REFERENCE/limits-and-defaults|Limits & Defaults]] — ports, payload caps, timeouts, JetStream limits, every default value

## How to use this section
| You need… | Go to |
|-----------|-------|
| "What does *X* mean?" | [[NATS/14-REFERENCE/glossary\|Glossary]] |
| "What's the config block for *Y*?" | [[NATS/14-REFERENCE/config-cheatsheet\|Config Cheatsheet]] |
| "What's the `nats` command to…?" | [[NATS/14-REFERENCE/cli-cheatsheet\|CLI Cheatsheet]] |
| "How should I name subjects?" | [[NATS/14-REFERENCE/subject-design-patterns\|Subject Patterns]] |
| "Which endpoint shows JetStream health?" | [[NATS/14-REFERENCE/monitoring-endpoints\|Monitoring Endpoints]] |
| "What's the default *Z* / max *W*?" | [[NATS/14-REFERENCE/limits-and-defaults\|Limits & Defaults]] |

## Version note
Built against **NATS server v2.14.x** line, `nats` CLI 0.2.x, 2026-06-14. Defaults and flags can change across releases — when this KB and `nats-server --help` / `docs.nats.io` disagree, **the live tool/docs win**. Always confirm version-sensitive values with `nats-server --version` and `nats --version`.

## OneMind Relevance
- This is the section to keep open in a second pane during fleet ops — the muscle-memory layer for [[NATS/09-CLI-TOOLING/nats-cli|nats CLI]], config edits, and health checks.
- The [[NATS/14-REFERENCE/subject-design-patterns|Subject Design Patterns]] doc encodes the conventions our fleet subject law should follow (`omos.<node>.<domain>.<event>`), so it's a design reference, not just a lookup.

## Sources
- Server config — https://docs.nats.io/running-a-nats-service/configuration
- nats CLI — https://github.com/nats-io/natscli
- Monitoring — https://docs.nats.io/running-a-nats-service/nats_admin/monitoring
