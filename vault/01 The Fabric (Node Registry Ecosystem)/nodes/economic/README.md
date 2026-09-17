---
title: "Economic — The Fabric (Node Registry)"
category: economic
type: index
status: scaffold
created: 2026-06-05
owner: zeus
tags: [type/index, domain/so, category/economic, status/scaffold]
ai_summary: "The economic category of The Fabric — the 6th kingdom. Holds the non-physical, non-digital nouns of Zeus's world that carry value: legal entities (LLCs, holdings, brands, community orgs), financial accounts (bank, brokerage, crypto, credit), and holdings (investments, real estate as an asset, IP, domains, licenses). SCAFFOLD ONLY — Zeus populates real entities."
---

# Economic — The 6th Kingdom

> The value layer of The Fabric. Everything Zeus **owns, controls, or transacts through** that is neither a physical-world object nor a piece of running software.

> **STATUS: SCAFFOLD.** Folders + doctrine only. No real entities populated yet — Zeus fills these in. Legacy will NOT fabricate accounts, LLCs, or holdings.

---

## Why this category exists

The original Node System had 5 categories — biological, ecological, physical, mechanical, digital. None of them hold an LLC, a brokerage account, a brand, a domain, or a license. Those are real nodes in Zeus's world (they have telemetry — balances, valuations, renewals — and they receive commands — transfer, renew, file). Forcing them into `digital` was a category error. `economic` closes the full-scope gap.

## The 3 classes

| Class | Folder | What lives here | Example subjects |
|-------|--------|-----------------|------------------|
| **org** | `orgs-registry/` | Legal + brand entities: LLCs, holding companies, brands, the community org | `*.org.onemind.>`, `*.org.weareonemind.>`, `cmd.org.<llc>.file_annual_report` |
| **account** | `accounts-registry/` | Financial accounts: bank, brokerage, crypto wallets, credit lines | `*.account.mercury.>`, `tel.account.<brokerage>.balance`, `qry.account.<crypto>.holdings` |
| **holding** | `holdings-registry/` | Owned value: investments, real-estate-as-asset, IP, domains, licenses, subscriptions | `*.holding.<domain>.>`, `tel.holding.<property>.valuation`, `cmd.holding.<license>.renew` |

## Modeling rules (inherited from The Fabric)

1. **Identity in the subject, properties in the registry.** `id: mercury` → `*.account.mercury.>`. The bank's address, login portal, account numbers (1Password ref only) live as frontmatter properties, never as subject tokens.
2. **A node is the thing, not the work.** A wire transfer is an EVENT (`evt.account.mercury.transfer_completed`), not a class. The account is the node; the transfer happened to it.
3. **Value readings are telemetry.** Balances, valuations, credit scores publish under `tel.<class>.<id>.<metric>` — ephemeral, measured, never the source of truth (the institution is).

## Bridge frontmatter (makes the file a bus node)

```yaml
id: mercury
class: account
category: economic
subject: "*.account.mercury.>"
access: L0            # money = owner-only by default
institution: "Mercury"   # registry PROPERTY — never a subject token
pillar: GE            # ClickHouse column, never a routing token
secret_ref: "op://OneMind/mercury"   # 1Password only — NEVER the credential itself
```

## What to populate (Zeus's call)

- `orgs-registry/` — the LLC(s), the holding co, `onemind`, `weareonemind`, `zeusdelacruz` brand, the community org
- `accounts-registry/` — Mercury + personal banking, any brokerage, crypto wallets, credit lines (Credit Karma / Equifax / Experian / TransUnion are *services* that report ON these — they live in `services-registry`, the accounts themselves live here)
- `holdings-registry/` — real estate held as asset, domains, software licenses, IP, recurring subscriptions of record

> **No secrets in markdown.** Every credential is a `op://` 1Password reference. Balances and account numbers are pulled live, never hard-coded.
