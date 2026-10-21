# LibreChat → Legacy: staged, operator-only pilot

**Not deployed. Do not sync this Application yet.** This PR prepares a private
conversation surface for the existing Hermes **agent**, not a direct model chat
and not an agent-framework replacement. The runtime inference test currently
fails with DigitalOcean HTTP 402 (account blocked). See the timestamped
[preflight evidence](../../docs/operations/legacy-librechat-preflight-2026-10-20.md).

## Ownership and route

```text
Zeus → Tailscale HTTPS → LibreChat → existing Hermes API :8642/v1
                          │                   │
                     dedicated MongoDB    existing tools / NATS fabric
```

- Argo Application `platform-librechat-pilot` owns the pinned upstream chart and
  the support Kustomization in namespace `librechat-pilot`. **Manual sync only**;
  no automated sync, prune, or deletion finalizer. After an eventual reviewed
  merge, the root Application would create the child Application, not deploy its
  workloads. Do not use Helm CLI installs or imperative workload applies.
- Chart `2.0.16`, app `v0.8.8`, multi-architecture image pinned by digest. Chart
  digest observed during inspection:
  `sha256:09f57a71294bcde827573600c3fe93226fedf8ab3b59c5e4c2e5ea9f76fbe6de`.
  Argo pins the chart version; this is not an OCI chart digest-enforcement claim.
- Backstage remains the workspace/catalog; NATS Agent Console remains the fabric
  inspection UI. This endpoint reaches Legacy only. It does **not** make every
  NATS agent appear automatically in LibreChat or implement a new NATS bridge.
- No new Kubernetes RBAC, service-account token, NATS credential, subject, model
  provider, or tool permission. NetworkPolicy permits only DNS, the existing
  mesh OIDC proxy, the actual Hermes API, and a labelled same-namespace database.

## Security boundary

This is a **Zeus-only pilot with Legacy's existing powers**, not a read-only agent.
A shared Hermes bearer credential does not implement per-tool authorization.
LibreChat is another caller of the existing harness and inherits its actual
permissions. Do not extend this to family accounts or workers until distinct
identities, approvals, revocation, and tool scopes have been acceptance-tested.
Existing read-only Kubernetes/Backstage MCP roles are unchanged.

Create a dedicated Authentik application/provider `onemind-librechat` through a
reviewed identity change, restricted to Zeus's `onemind-librechat-pilot` group.
Redirect URI must be exactly
`https://librechat.tail717ba0.ts.net/oauth/openid/callback`.
Include the group in the **signed ID token's `groups` claim**. LibreChat checks
the required group; local email login/registration and public sharing are off.
Social registration is enabled only to provision an authorized OIDC user on first
login; it is not permission for other Authentik users. Test denial explicitly.
No automatic admin-role mapping or ID-token forwarding to Hermes is configured.

Issuer uses the existing private `onemind-auth` route, including the CoreDNS
rewrite and `backstage/authentik-mesh` egress service owned by `platform-backstage`.
NetworkPolicy allows the existing `tailscale` egress ProxyGroup on TCP 443, not
unrestricted internet HTTPS. Verify DNS, issuer equality, JWKS and TLS from the
pilot namespace; a successful browser login alone does not prove server routing.
If that shared route changes, update the dependency deliberately, not by allowing
all outbound traffic. No global DNS change is included here.

## Required credentials — reference contract, no values

Every key below is an explicit non-optional `secretKeyRef` to
`librechat-pilot/librechat-runtime`. Missing keys prevent container creation; the
chart is not allowed to generate temporary cryptographic secrets. **This Secret
and its SealedSecret do not exist in this PR.** Never paste keys into chat or
commit a plaintext Secret. Generate/seal the approved values in a separate step.

| Key | Requirement |
|---|---|
| `MONGO_URI` | Authenticated, dedicated LibreChat database; approved same-namespace service on port 27017, labelled `onemind.io/librechat-database: 'true'`. No shared application's database/user. |
| `CREDS_KEY` | Persistent 32-byte cryptographic key represented as 64 hex characters. |
| `CREDS_IV` | Persistent 16-byte IV represented as 32 hex characters. |
| `JWT_SECRET`, `JWT_REFRESH_SECRET` | Distinct, persistent high-entropy secrets. |
| `OPENID_CLIENT_ID`, `OPENID_CLIENT_SECRET` | Dedicated confidential Authentik client, not Backstage's client. |
| `OPENID_SESSION_SECRET` | Persistent high-entropy OIDC session secret. |
| `HERMES_API_KEY` | Approved credential for the existing Legacy API. Do not copy a provider's model key here. Confirm the scope/risk of sharing Legacy's credential before activation. |

The chart's bundled MongoDB, Meilisearch, Redis, RAG and fanout services are
disabled. **This PR does not provision MongoDB or backups.** It deliberately
does not deploy the chart's unauthenticated legacy MongoDB default. A reviewed
database/backup manifest and sealed credentials are prerequisites, not implicit
external services. Database ownership and backup restoration must be in GitOps.
Do not sync an incomplete frontend and describe history as persistent.

## Conversation, cost and file limits

The configured header combines the authenticated LibreChat user ID and the
conversation ID into `X-Hermes-Session-Key`. Verify first-message generation,
conversation reload, a second conversation, and a second isolated test identity
before relying on it. Never use a single static session key. Static configuration
tests do not prove runtime isolation. Do not persist a returned Hermes session ID
globally across users/conversations. Check for duplicated history when both ends
maintain context.

Automatic title and summarization calls are disabled: these would otherwise call
the privileged harness and generate extra work/spend. File uploads for this
endpoint are disabled pending attachment compatibility and backup testing. The
2 GiB image PVC is **not** the conversation database and does not cover every
possible upload directory. Single replica plus `Recreate` avoids a RWO rollout
deadlock; it is not high availability. Readiness is `/readyz`, liveness `/livez`.

LibreChat does not know the true cost of all Legacy's delegated/tool/model calls.
Do not present an empty/zero cost in either UI as free usage. The later fabric
accounting contract must distinguish provider-reported, estimated and unknown
spend, include worker/delegation costs, and correlate agent/session/request IDs.
That integration is not implemented by this pilot.

## Activation and recovery gates

Keep the PR draft and Application unsynced until these are recorded as passed:

- [ ] Resolve Legacy's configured inference provider's 402 without unauthorized
      billing, credential, or provider changes. Repeat one authorized harmless
      streaming test: HTTP 200 alone is not success; require a real answer and
      successful terminal event. Do not automatically retry paid probes.
- [ ] Approve pilot capacity; provision a dedicated authenticated MongoDB and PVC
      through GitOps, choose retention/storage and demonstrate denied anonymous
      database access. Decide the pilot's RPO/RTO before accepting real history.
- [ ] Add encrypted off-host database backups and a scheduled success/freshness
      check. Prove an application-consistent backup (e.g. a supported replica-set
      backup strategy or quiesced disposable pilot), not just a live PVC snapshot.
- [ ] Restore to an **isolated** database/namespace, verify a known conversation,
      retain image data and the encryption/JWT/OIDC secrets, and record recovery
      time. Off-host must not mean only the same physical Mac disk. Do not restore
      over production or delete existing volumes. Never assume existing PG backup
      jobs also back up MongoDB.
- [ ] Create the dedicated Authentik provider/group/binding and sealed runtime
      credentials. Verify issuer/JWKS, exact redirect, authorized login, denied
      nonmember/anonymous access, logout, and session revocation behavior.
- [ ] Verify restricted PodSecurity, storage permissions, probes, DNS/OIDC and
      explicit network-denial cases in the target environment.
- [ ] Verify two-conversation/user separation, reconnect, persisted history after
      restart, visible inference failure, and no hidden title/model requests.
- [ ] Test an approval-required, non-destructive operation. Cancellation or closing
      a browser must not be represented as stopping remote execution. Verify the
      actual Hermes/NATS approval and cancellation behavior before agent actions.
- [ ] Add a working private launch link to Backstage only after those gates; do
      not advertise a dead route as an installed/healthy service.
- [ ] Human review approves activation; sync the existing Argo Application. No
      foundation migration, TAK migration or Phase 7 completion is implied.

Rollback: leave automated sync off, stop the pilot through a reviewed Git change,
remove its ingress only if approved, and retain the database/PVC/secrets/backups.
Do not cascade-delete the namespace. Revoking pilot access must not rotate shared
Legacy credentials without a separate approval and dependency review.

## Reproduce source checks

```bash
python3 -m unittest discover -s scripts/tests -p 'test_librechat*.py' -v
helm template librechat-pilot oci://ghcr.io/librechat-ai/librechat-chart/librechat \
  --version 2.0.16 --skip-tests --namespace librechat-pilot \
  -f platform/librechat/pilot/values.yaml > /tmp/librechat-pilot-rendered.yaml
python3 scripts/tests/test_librechat_pilot.py --rendered /tmp/librechat-pilot-rendered.yaml
kubectl kustomize platform/librechat/pilot > /tmp/librechat-pilot-support.yaml
```

These commands render locally, not deploy. The API preflight helper at
`scripts/ai-operator/legacy_api_preflight.py` defaults to metadata only; an optional
`--chat-probe` is a **real inference call** and must be authorized separately.
It reads an existing key from an environment variable, never prints it, does not
follow redirects with credentials, and does not change the model/provider.

### Primary references

- [Upstream Helm deployment](https://www.librechat.ai/docs/local/helm_chart)
- [Authentik integration](https://www.librechat.ai/docs/configuration/authentication/OAuth2-OIDC/authentik)
- [Custom endpoint configuration](https://www.librechat.ai/docs/configuration/librechat_yaml/object_structure/custom_endpoint)
- [Pinned LibreChat source](https://github.com/LibreChat-AI/LibreChat/tree/e8f3be08623663d4ad7f7241e693c94469b63bb0)
- [Hermes API server](https://hermes-agent.nousresearch.com/docs/user-guide/features/api-server)

Docs may describe newer features; pinned chart/source and actual runtime tests
take precedence over claims based on latest documentation alone.

## Executable acceptance matrix — preparation 2026-10-20

These are operator procedures for a later authorized test window, not results.
Preparation used the recorded 2026-10-20 preflight only; no new network or
inference test was performed. The 16:17:38–39 UTC provider402 remains the last
recorded inference result. Keep the draft unsynced. Source tests can pass while
runtime acceptance remains blocked.

Before execution, the coordinator must record exact source/image/chart, target
namespace, operator, UTC window, approved request/spend ceiling and stop conditions.
Identity, database, backup and restricted disposable test deployment need their
own approved preparation; this matrix does not authorize them. Resolve provider402
by separately approved means before live inference. Start with offline fixtures;
never generate provider errors by changing billing or keys. A no-tools prompt is
not a permission control.

Use synthetic markers `PILOT_A1`, `PILOT_A2`, `PILOT_B1` without personal content.
A is Zeus's pilot identity; B is an explicitly approved isolated test identity,
not permission to onboard family or workers. Without B, mark isolation blocked;
two tabs signed in as A do not substitute for two identities. Capture only
allowlisted, redacted correlation metadata: pseudonymous user/conversation/session
and request/run IDs, UTC timestamps, terminal outcomes and call counts. Never
retain bearer headers, cookies, tokens or real conversation bodies in PR evidence.
Keep any necessary sensitive evidence in approved storage.

| ID / current status | Procedure in the approved test environment | Required evidence and pass condition |
|---|---|---|
| S1 / NOT RUN | Sign in as A; open conversation C1 and send one harmless `PILOT_A1` request. Inspect the first outbound request through an approved metadata-only trace. | User/conversation IDs exist before the first Hermes call; `X-Hermes-Session-Key` resolves to A/C1, never blank, literal template text or a global constant. Record returned session/run correlation. Require a real answer and successful stream termination; HTTP200 alone fails. |
| S2 / NOT RUN | Reload C1, sign out/in as A, reopen it and send one approved follow-up. During a separately approved request, interrupt only the browser connection and reconnect. | Same A/C1 mapping, one persisted turn per request and no doubled context between LibreChat history and Hermes state. No blind resubmit or duplicate remote run. Browser disconnect must not be labelled remote cancellation. |
| S3 / NOT RUN | As A, open C2 and send `PILOT_A2`; compare C1/C2 traces and reopen both histories. | Distinct session keys and correlated Hermes sessions; each history contains only its own markers. Separate identifiers alone do not prove authorization. |
| I1 / BLOCKED: test identity | In a separate browser profile, sign in as approved B and create C3 with `PILOT_B1`. Using supported application history/read routes, attempt to open A's known C1/C2 IDs as B; repeat B-to-A. | Cross-user reads denied without content/session disclosure in UI and backend. A/C1, A/C2 and B/C3 mappings differ. Do not guess unrelated users' IDs or expose the shared Hermes credential to a browser. |
| I2 / BLOCKED: dedicated OIDC | Check A login, nonmember login denial and anonymous history denial. Sign A out, then exercise a retained test session after the approved revocation action. | Exact redirect, signed group enforcement and issuer/JWKS/TLS pass. Record actual logout/revocation timing and remaining session usability; unexpected continued access blocks acceptance. Never record raw tokens. |
| E1 / LIVE FAILURE RECORDED | Offline, replay redacted provider402 and HTTP200 SSE `finish_reason: error` fixtures, plus interrupted/malformed streams. After separately approved provider resolution, make one bounded harmless live request. | Fixtures show visible failure, no fabricated success and no automatic retry/title/summary calls. Retain the original 2026-10-20 failure receipt. Only a new dated real answer plus successful terminal event establishes recovery; fixtures do not. |
| P1 / BLOCKED: storage | After approved authenticated MongoDB/PVC provisioning, retain C1/C2. Perform the reviewed disposable pilot restart, reopen histories and inspect session correlations. | Histories/ownership survive with no duplicated turns; anonymous database access denied. Record database/PVC ownership, revision and restart evidence. Image PVC existence or an unchanged pod proves no conversation persistence. |
| R1 / BLOCKED: backup/restore | With approved RPO/RTO, capture an application-consistent encrypted off-host backup. Download/decrypt into an approved isolated recovery target; restore C1/C2 and necessary crypto/session material through the reviewed procedure. | Known markers, ownership/history and identity separation verified without inference or changing production. Record backup timestamp, independent destination, consistency method, integrity, measured data loss/recovery time and secret recovery references without values. Define revocation/relogin after restore. Restored MongoDB alone proves no Hermes session recovery; verify that state separately or block continuity acceptance. |
| A1 / BLOCKED: action controls | With an approved non-destructive fixture operation, exercise approval denial and expiry before one explicit approval. Test actual stop/cancel and reconnect correlation separately. | No action before approval; trace ties approval/outcome to the correct user/session/run. Denial/expiry cause no action. Confirm remote terminal state; UI close or JetStream acknowledgement is insufficient. Existing Legacy remains privileged until enforced restrictions are proven. |
| C1 / BLOCKED: spend controls | Before live calls, agree a per-window budget/request cap covering parent, worker, tool and delegated model calls. In an isolated fixture harness, exercise denial of unapproved delegation and cap exhaustion; reconcile one later approved request against available provider/tool records. | Name the enforcing component and prove blocking before excess/unauthorized dispatch; prompts/UI estimates are insufficient. Separate provider-reported, estimated and unknown cost with parent/child run IDs and call counts. Missing downstream accounting/enforcement leaves delegated/autonomous spend blocked; never infer zero cost. |

Execution order: offline E1 fixtures, prerequisite review, approved bounded live
S1, then S2/S3/I1/I2/P1/R1. A1/C1 need separately scoped operation/spend approval;
the request/spend ceiling applies before any live row. Stop on isolation leaks,
duplicate dispatch, unexpected tools, unapproved charges, recovery inconsistency
or renewed provider error. Do not retry paid requests automatically. Preserve
failed-row evidence before continuing. No real user history or autonomous/delegated
work is accepted while its identity, persistence, recovery or spend gate is open.

For each row the coordinator's receipt must contain: ID, PASS/FAIL/BLOCKED/NOT RUN,
actual UTC date/time, tested revisions, prerequisite/authorization reference,
redacted evidence location, observed result and remaining limitation. Keep offline
and live receipts separate. The coordinator owns tests, tracker reconciliation
and updates to existing #332; this appendix certifies no deployment, service
acceptance or Phase 6 completion.
