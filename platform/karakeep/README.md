# Karakeep — private native Kubernetes deployment

The user authorized deployment on 2026-10-25. Argo Application `platform-karakeep` now automatically reconciles one replica each of web/worker, Chrome and Meilisearch. Web starts in sync wave 2 after the app-specific Authentik hook succeeds in wave 1. **First web start runs the official SQLite migrations.** The private address is **https://karakeep.tail717ba0.ts.net**. Source activation is not proof of a successful live rollout, owner sign-in or restore; verify the acceptance checks below. No public Ingress, LoadBalancer, NodePort or Funnel is configured.

## Native source and dependencies

Adapted from the official [Kubernetes/Kustomize installation](https://docs.karakeep.app/installation/kubernetes/) and [v0.34.1 native manifests](https://github.com/karakeep-app/karakeep/tree/ebd122bb45b38d9dd45f7c6db55b5587bf6b1265/kubernetes), not a Compose translation or community Helm chart. `upstream.lock.json` records exact upstream commit, source checksums and registry-verified image index digests. Updates must review those native resources and release notes before replacing the pins.

| Piece | Pin | Role / persistent state |
| --- | --- | --- |
| `web` | Karakeep 0.34.1 | Official combined web/worker image; SQLite and assets at `/data` on `data-pvc` (20Gi) |
| `chrome` | Karakeep Chrome 154.0.8037.57-r1 | Official browser image; ephemeral `/tmp`, no PVC |
| `meilisearch` | Meilisearch v1.41.0 (upstream native pin) | Supported full-text search; `meilisearch-pvc` (5Gi) |

The metadata preflight found no existing Karakeep, Hoarder, Meilisearch or shared browser workload. Existing shared OpenSearch cannot replace Meilisearch through the supported Karakeep setting. Karakeep's database is SQLite; no shared CNPG database or new PostgreSQL server is introduced. Search and browser dependencies are dedicated because they are the official supported native layout. OpenAI-compatible AI tagging and embeddings are optional; no inference credential, model server or gateway route is configured. Automatic tagging/summaries, embeddings, semantic search and OCR start disabled. Enabling them is a separate reviewed configuration/integration, not app acceptance.

The official native service defaults to LoadBalancer; here it is ClusterIP behind the Tailscale operator. Plaintext Kustomize secretGenerator/.secrets is replaced by the preserved strict SealedSecret `karakeep/karakeep-secrets`, keys `NEXTAUTH_SECRET` and `MEILI_MASTER_KEY`. The upstream sample's legacy `NEXT_PUBLIC_SECRET` is absent from v0.34.1 shared config and is not generated. A fresh dedicated OIDC client secret has strict sealed copies in `karakeep/karakeep-oidc` and `auth/karakeep-oidc-bootstrap`. Existing application/global identity credentials are not rotated.

## Owner sign-in

Authentication uses the native Karakeep OIDC integration, issuer discovery at `https://auth.onemindos.dev/application/o/onemind-karakeep/.well-known/openid-configuration` and strict callback `https://karakeep.tail717ba0.ts.net/api/auth/callback/custom`. Password login and local signup are disabled (`DISABLE_PASSWORD_AUTH=true`). Implicit linking by email is disabled. Authentik's existing login/enrollment and verification flows are reused unchanged.

The Argo Sync hook creates only the `onemind-karakeep` provider/application and an access binding for the existing active owner `akadmin`, stable UID `f73e10d872325ed63330a138782f12d3e8252e00a05466e2f0180bcc885a42c4` (live-verified against the existing Backstage registration on 2026-10-25). It refuses missing/replaced owners, duplicate providers or unexpected access bindings. The bootstrap token is used only by the short-lived hook in `auth`, never by the application. The hook has no Kubernetes API token or workload mutation permissions.

Pinned v0.34.1 sets its OIDC `disableSignUp` directly from `DISABLE_SIGNUPS`, so `DISABLE_SIGNUPS=false` is necessary for the first OIDC account creation. Only the bound Authentik owner can enroll; local/password registration remains disabled. After that owner has signed in successfully, set `DISABLE_SIGNUPS=true` through reviewed Git/Argo to prevent any new account creation. Do not enable password signup, add new users or weaken Authentik verification as a bootstrap shortcut. Native integration references: [Karakeep configuration](https://docs.karakeep.app/configuration/environment-variables/), [Authentik integration](https://integrations.goauthentik.io/documentation/karakeep/), and the [pinned authentication source](https://github.com/karakeep-app/karakeep/blob/ebd122bb45b38d9dd45f7c6db55b5587bf6b1265/apps/web/server/auth.ts).

## Initial budgets and isolation

App requests total **600m CPU / 1792Mi memory**, limits **4 CPU / 4608Mi memory**, plus the Tailscale proxy and transient bootstrap/backup overhead. Each of two backup movers requests 100m CPU / 256Mi memory with a 1Gi cache PVC. The bootstrap hook requests 50m CPU / 64Mi memory. These are initial estimates, not measured usage. Recheck capacity alongside the shared-database recovery and other app rollouts; keep a single web writer and Recreate strategy for both RWO stateful workloads.

Startup/readiness/liveness probes use upstream `/api/health`, Chrome `/json/version` and Meilisearch `/health`. The web health route proves the HTTP process responds; it does not certify search, crawling, workers, authentication or recovery. Namespace PSA is restricted, Kubernetes API tokens are disabled, and pods are non-root with dropped capabilities and RuntimeDefault seccomp. The web image's s6 init needs writable container runtime/cache paths, so only Chrome/Meilisearch have read-only root filesystems. These settings follow the upstream baked node-owned `/run` and browser UID65534; inspect bounded startup failures before changing them.

Default-deny ingress/egress allows DNS; only the named Karakeep Tailscale ingress proxy can reach web3000; only web can reach chrome9222 and meilisearch7700. Web and Chrome can fetch public IPv4 HTTP/HTTPS, excluding private, tailnet, loopback, link-local/metadata and reserved ranges. Meilisearch has no public egress. There is no route to CNPG, model gateways, internal services or private bookmark targets. Changing that is an intentional scope expansion. [Upstream security guidance](https://docs.karakeep.app/administration/security-considerations/) still treats crawler input as trusted-user access: the browser runs upstream `--no-sandbox`, and website crawling uses this server's public IP. Do not treat the basic SSRF checks as complete isolation.

## Rollout and acceptance

1. Review the exact-main Argo diff and current capacity. Merge/sync through Argo. Confirm the namespace/PVC protections, SealedSecret Synced status and expected secret key names without printing values. Record PVC UIDs/PVs. No manual workload apply, scaling, init SQL or claim recreation.
2. Confirm Chrome and Meilisearch ready, the Authentik hook successful, then web ready. A failed provider hook must leave wave-2 web unstarted. Inspect logs/events before changing security or resource settings. The official image performs SQLite migrations before web/workers.
3. Verify the tailnet route's TLS and login redirect. Sign in as the existing owner, confirm session/CSRF behavior, denied password signup and denied access for unbound Authentik users. Then close OIDC account creation with `DISABLE_SIGNUPS=true` through Git/Argo.
4. Store a harmless bookmark and attachment, exercise full-text search and a permitted public-page crawl/screenshot; confirm private/metadata targets remain blocked. A unique scoped user API key and official API/CLI or MCP access need separate authenticated enrollment; do not reuse app secrets as user API keys. No agent key or messages are created by this source.
5. Trigger both initial VolSync sources through reviewed Git and confirm encrypted repository completion. Follow [BACKUP.md](BACKUP.md) for SQLite integrity and isolated restore acceptance before normal import. Daily whole-volume snapshots are crash-consistent; successful upload alone does not prove recovery.
6. Confirm restart/persistence and an isolated restore with bookmark, attachment and search checks. Verify Android/iOS/extension compatibility, offline download/reading and later resync separately; no offline/home-autonomy claim follows from a tailnet URL.

## Backup, upgrade and rollback gates

See [BACKUP.md](BACKUP.md). Keep namespace/PVCs/sealed values on rollback. Before upgrades, take and verify a recoverable full data backup; read release/migration compatibility and Meilisearch dump/reindex requirements. Reverting an image commit does not reverse database/index migrations. Source-only failure can revert the Argo/source change while no app has started; after first start stop web via reviewed Git and restore into an isolated new claim first. Never delete/recreate claims, reset SQLite or rotate auth/search keys just to clear readiness.

Local verification: `kubectl kustomize platform/karakeep`, strict kubeconform, changed-YAML/manifest gates and `python3 -m unittest discover -s scripts/tests -p 'test_gitops_karakeep.py'`. Live rollout, owner enrollment, backup and restore results must be recorded from observed evidence separately.
