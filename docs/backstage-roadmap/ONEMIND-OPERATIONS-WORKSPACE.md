# OneMind operations workspace integration

This branch integrates eight independently reviewable feature PRs. It reconciles shared
frontend/backend registration, home layout, documentation navigation and the offline/MCP
runbook allowlists. It does not merge feature PRs into main or change the deployed image.

| Feature PR | Delivered capability | Activation requirement |
| --- | --- | --- |
| #55 | Reviewed scheduled inventory refresh; retained missing identities and freshness | Dedicated observational runner, scoped org metadata token, full visibility confirmation and explicit enablement |
| #57 | Activity inbox and bounded authenticated event intake | Restricted producer identity/configuration; recipient and delivery readiness checks |
| #58 | Systems dependency and impact workspace | Catalog relationships determine coverage; no automatic inference of missing dependencies |
| #59 | Paginated actionable Zoho task queue and provider contract | Read-only US Zoho credentials, numeric portal ID, account API entitlement |
| #60 | Evidence-grounded assistant and source explorer | Optional operator-controlled Ollama endpoint/model; otherwise explicitly no connected model |
| #61 | Service Operations tab with Argo, namespace metrics and Grafana links | Existing read-only cluster observations; per-service dashboard mappings beyond NATS |
| #62 | App onboarding, deployment change and recovery preparation draft-PR workflows | Administrator user and restricted short-lived GitHub token for a real workflow run |
| #63 | Explicit private offline documentation pack plus local exporter | Download while connected, or export from a trusted repository checkout |

## Release boundaries

- Authentication, administrator membership, machine catalog write protection and existing
  read-only cluster permissions remain enforced. No guest login or broad execution tokens.
- All runtime additions are useful without fabricated provider data: unavailable,
  configuration-required, partial and stale states are distinct from healthy observations.
- New event/model/cloud credentials are not installed by this branch. Inventory automation
  remains disabled unless its explicit prerequisites and repository variables are satisfied.
- Workflows create draft proposal PRs outside active GitOps paths. They do not deploy or restore.
- The offline pack is static knowledge, not a disconnected copy of live cloud applications.
- Existing portal URLs, backgrounds, logos, database volumes and backup keys are preserved.

## Evidence and promotion

The independent UI/backend work and this assembled branch compile with the supported Node
22 toolchain. Inventory scripts passed syntax compilation; new workflow artifacts were
rendered with installed Backstage template tooling. No new tests were added or run.

A combined image build must pass before promotion. Record its digest and final branch SHA
in the integration PR. Then review frontend behavior with administrator identity before
pinning that digest in a separate deployment PR. Do not deploy from a moving image tag.

Verify runtime read/deny behavior, the event recipient/inbox path, documentation navigation,
entity associations, provider freshness, each real template run, model evidence handling,
and offline download on the release candidate. Compilation and container packaging do not
prove account integrations or live end-to-end workflows.

Recommended activation sequence: workspace and existing observations first; restricted
event producers next; inventory automation and Zoho after credentials are configured;
model connection and home-cluster cutover only after their compute/routing plans exist.
