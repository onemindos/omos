---
name: onemind-architecture
description: "Use when deciding how OneMind knowledge is organized: which areas and record types exist, whether something is a document or a structured record, and where each kind of content belongs."
version: 2.0.0
author: Legacy
license: private
platforms: [linux]
metadata:
  hermes:
    tags: [onemind, architecture, knowledge, taxonomy, routing]
    related_skills: [onemind-domains, doks-cluster-ops, nats-operations]
---

# OneMind Architecture (knowledge taxonomy)

This skill keeps the OneMind information architecture: the four areas, the record
types each one needs, and rules for routing content. It is platform-neutral.

**Status:** Zeus does not use Notion now, and no replacement is chosen. A
knowledge platform will be chosen later. Until then, do not create or sync pages or
databases in any external tool, and do not invent identifiers for a store that does
not exist. Durable knowledge lives in Git (`docs/`, `platform/hermes/skills/`, the
repo's `docs/STACK.md`) and in Hermes memory; unsorted capture goes to the inbox
(the chat thread with Zeus or a note in the session) for later triage. When the
platform is picked, this taxonomy is the schema to build from.

How to act in each area (rules, safety, procedures) is in `onemind-domains`. This
skill is only about what is recorded and where it goes; do not duplicate those
procedures here.

## Document or structured record

| Signal | Document | Structured records |
|---|---|---|
| One reference doc, research note, ADR, SOP | yes | |
| Five or more entries of the same type | | yes |
| Needs filtering, sorting or grouping | | yes |
| Needs relations to other records | | yes |
| Inventory, catalog, registry or log | | yes |

Rule: more than five of the same thing that you will query means a structured
record set. One deep write-up means a document. Do not build an empty record set;
build it when content is ready to fill it.

## Areas

Four life and business areas (the same four named in the soul): SO sovereign
operations, GE generational empire (business), HP holistic performance, LE legacy
evolution (land, home, estate). Each has actionable items (tasks with deadlines)
and reference knowledge (wiki). The inbox holds anything unsorted.

## Record types by area

**SO, infrastructure:** devices (phones, laptops, routers, edge boards, drones),
services (namespace, port, version, status), AI agents on the bus, NATS subjects
(registry of the locked grammar patterns in `docs/NATS-DOCTRINE.md`), code repos,
deployments (Argo Applications and Helm releases), ADRs, runbooks, incidents,
certificates (metadata only: type, issuer, expiry; never the secret), network
topology, research pipeline, vendors.

**GE, business:** contacts/CRM, companies, deals, products, features, content plans,
revenue streams, contracts (NDAs, MSAs), IP, grants and programs, competitive
intel, events.

**HP, performance:** journal, accounts and subscriptions (metadata, not numbers
or credentials), health metrics and workouts, learning, goals, habits,
relationships, values, media.

**LE, legacy:** properties, inventory of gear and appliances, documents (type,
status, expiry, where the original is kept; never the content), vehicles,
maintenance, homestead (land, garden, animals), education, estate plan, home
expenses.

Never store secrets, credentials, identity numbers or account numbers in any
record. Record where the credential lives (for example "SealedSecret X in
namespace Y", or "Zeus holds it").

## Routing content

```
Actionable with a deadline?        -> tasks for that area
Reference or knowledge?            -> wiki/document for that area
One of five or more like things?   -> the structured record type above
Unsure?                            -> inbox, triage later
```

| Content | Now | When record sets exist |
|---|---|---|
| Architecture doc, research note | document in Git `docs/` | stays a document |
| Device spec or profile | SO document | devices |
| Running service or deployment | `docs/STACK.md` (generated from the cluster) | services, deployments |
| Agent profile | `platform/hermes/agent/` | AI agents |
| NATS subject definition | `docs/NATS-DOCTRINE.md`, `services/schemas/` | NATS subjects |
| Architecture decision | ADR in `docs/` | ADRs |
| SOP or runbook | `docs/operations/` or a skill | runbooks |
| Incident | a dated note under `docs/operations/` | incidents |
| Cert or secret metadata | document, metadata only | certificates |
| Partner, deal, competitor | GE document | CRM, deals, competitive intel |
| Financial record | out of scope for the agent; Zeus keeps it | accounts, transactions |
| Property, vehicle, legal or estate document | LE document, metadata only | properties, vehicles, documents |

## Relations worth modelling

Device runs service; service has deployments and incidents; incident links to a
runbook; ADR links to the research that led to it; AI agent links to its NATS
subjects; repo powers service; product has features; feature links to tasks; deal
links to contacts; vehicle links to maintenance; property links to documents.
Create both record sets before adding a relation.

## Suggested build order (when a platform exists)

Devices, services, AI agents, contacts, products, competitive intel, ADRs,
subscriptions, gear inventory; the rest as content arrives.

## Pitfalls

- Do not route content to a tool that is not chosen yet; use Git and memory.
- Check `docs/STACK.md` and the live cluster before recording a service as live
  (SOUL rule 6). Records of pod or node counts go stale fast.
- When a record type is created later, update this skill in the same change.
- Never copy this skill's contents into `onemind-domains`; link to it instead.
