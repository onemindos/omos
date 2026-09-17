---
title: "Domain Health: {DOMAIN_NAME}"
type: review
status: active
created: "{YYYY-MM-DD}"
owner: "{agent}"
tags: [type/review, cadence/monthly]
---

# Domain Health: {XX} {DOMAIN_NAME}

## Status: {🟢 Healthy / 🟡 Moderate / 🔴 Neglected}

| Metric | Value |
|--------|-------|
| Total files | {count} |
| Last updated | {date} |
| Active projects | {count} |
| Knowledge gaps | {count of status/gap files} |

## Recent Activity (Last 30 Days)

{List of files created or modified, from git log}

- {file} — {action} — {date}

## Knowledge Gaps

{Files tagged status/gap or status/stale}

- [ ] {gap 1 — what's missing}
- [ ] {gap 2 — what's outdated}

## Cross-Domain Connections

{Most-linked files, backlinks to/from other domains}

| Connected Domain | Type | Through |
|-----------------|------|---------|
| {domain} | {backlink/tag} | {file} |

## Recommendations

1. {AI-suggested improvement}
2. {AI-suggested improvement}
3. {AI-suggested improvement}

---

*Generated: {DATE} by {AGENT}*
