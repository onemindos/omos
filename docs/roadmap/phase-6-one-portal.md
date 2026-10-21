# Phase 6: one portal — completed

Completed 2026-10-19 after live gates, not merely compilation.
Tracker: [PLATFORM-ROADMAP.md](../../PLATFORM-ROADMAP.md).
Detailed evidence/rollback: [integration runbook](../operations/phase6-portal-integration.md).

## Delivered

- Existing OneMind Backstage **1.55.3** remains the host, preserving branding,
  home/Today, catalog, user settings, runbooks, integrations and PostgreSQL data.
- Pinned OpenChoreo **1.3.1** core/catalog/scaffolder and observability plugins add
  real Definition, Deploy, Logs and Metrics views. No workflow plane or assistant
  pack was silently enabled.
- Authentik `openchoreo-auth` retains verified signed-subject catalog mapping and
  existing operator membership, with one intersected host/OpenChoreo policy.
- Shared capability caching in the managed permission DB fixes cross-replica
  requests. Keys are SHA-256, values are authorization results, TTL follows the
  upstream token expiration, DB failures deny closed; raw tokens are not stored.
- Temporary built-in console, its private ingress/Application and the memory
  candidate are retired through GitOps. Reused credentials and bindings remain;
  no PVC, database, backup or identity was deleted.

## Release and live gate

Source `3dd34900a5b455696387394398b72f80ce3376ba`; verified GHCR digest
`sha256:b2fc101eacc504fadff99b638b342ce9bf56782a4901ea49dffe2fdd09f1967a`.
Main portal: https://onemind-preview.tail717ba0.ts.net/ — **2/2 ready**, Argo
Synced/Healthy. Control plane, observability, identity and app-of-apps also healthy.

Real Authentik login maps to `user:default/akadmin` and existing administrators.
Home renders 130 launcher entries; catalog has 606 entities. Today, notifications,
settings and packaged runbooks remain available. Pilot Deploy shows active
`stac-browser-v5`; Definition is real; CPU/memory charts and runtime logs work,
including a fresh HTTP 200 observed at 18:27:08 UTC.

Both new pods passed anonymous 401, guest 404, machine catalog read 200, catalog
writes/platform/cache/context 403, and six explicitly read-only MCP tools. Actual
MCP health reads work. Non-admin human policy has regression coverage; no second
human account was created or represented as live-tested. No live mutation test.

## Verification and deliberate boundaries

49 tests, full TypeScript, immutable Linux release, TechDocs/backend bundle, JSON
exports, native AMD64 SQLite and remote digest verification passed. An isolated
PostgreSQL 17 exercise verified concurrent initialization/upsert and cross-client
cache visibility/expiration. Changed YAML/schema/Helm/kustomize and new-secret
checks passed. Existing encrypted backup restore verification passed today.

The browser/live gates found and corrected request-body reuse, exact-endpoint
CSP/CORS, owner-name prefixing and the multi-replica cache gap before retirement.
Upstream's published compatibility matrix still tests 1.51; this is empirical
OneMind verification on 1.55.3. Keep the frontend cohort and private ABAC-router
adapter pinned and reverify future updates.

Catalog sync is deliberately poll-only until authenticated event ingress is
reviewed. The unused OpenChoreo webhook forwarder is disabled; native Backstage
PostgreSQL Events and NATS are unchanged. The control-plane app retains automated
prune=false; a selected Argo operation pruned exactly the 10 retired UI/forwarder
resources. The identity hook removed retired exact callbacks and succeeded.

GitHub Actions cannot start due account payment/spending limits. Whole-tree
`make validate` still fails unrelated old YAML debt. Neither was reported green
or bypassed; local and live evidence is recorded separately.

Implementation/rollout: PRs #311–#319. This PR (#282) records final acceptance.
Phases **7 and 8 remain Not started**; no Layer 3 migration is bundled here.
