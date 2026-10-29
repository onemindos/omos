# Owner proposals: role modes with Authentik and the NATS auth callout (FEATURES J2)

Status: **proposals for the owner**, 2026-10-29, filed under onemind-xr decision 0004 #8.
Nothing is live. **No new subject pattern**: the roles grant only subjects already in the locked
grammar and in #554/#562 (`svc.layer.*`, `svc.entity.*`, `svc.policy.*`, `svc.cot.send`,
`agents.prompt.hermes.>`, `cmd.map.>`). `$SYS.REQ.USER.AUTH` and `_INBOX.<sub>.>` are
server-internal subjects, not grammar subjects. So nothing is added to the subject register
(#556).

| ID | Decision needed | Evidence | Proposed ruling | Used until ruled |
|---|---|---|---|---|
| J2-P1 | **The callout and A19 stage A2 cannot both be on.** In config mode, a leaf that logs in as an account user is always sent to the callout (the `allowed_accounts` filter skips LEAF), and the callout cannot verify its password | nats-server v2.14.2 `server/auth.go` (`c.kind != LEAF && juc == nil && opts.AuthCallout != nil ...`); `TestLeafNodes` (A1 leaf: 0 callout requests, connected; A2 leaf: refused) | Choose one order: (a) turn the callout on at stage A1 and hold A2 until operator/JWT mode (stage 2), where callouts are per account; or (b) run A2 first and keep browsers on stage-1 passwords until stage 2. Recommended: (a), because A2 mainly adds leaf permissions while the callout removes the shared password from browsers | (a): runbook enables the callout after stage A1 only |
| J2-P2 | **Role groups ride in the NATS user name** `ak:<username>;<group>,...`, because config mode drops user-JWT tags from `Nats-Request-Info` | `server/auth.go` sets `c.tags` only on the operator-mode JWT path; `TestOperatorReachesPolicyWithVerifiedIdentity` | Ratify the name format; add a rule to `scripts/nats_identities.py check` (#562) refusing `;` in any configured username so no config user can look like a callout user | Format in use; checker rule documented in `platform/auth-callout/identities.fragment.yaml` |
| J2-P3 | An **agent proposing on behalf of a callout user** (`on_behalf_of: user:olga`) is checked without olga's groups: the request comes from the agent, so the policy service has no token for olga, and the proposal is denied unless static tuples exist | `services/policy` `groupTuples` (only for the verified caller) | Short term: the owner writes static `user:<name> member group:<group>` tuples for named people who use Hermes. Later: an Authentik-to-OpenFGA membership sync, or the client attaches a short-lived delegation the policy service can verify | Denied (fail closed) |
| J2-P4 | **Production URL of the L0 web client** (OIDC redirect URI) and whether browsers keep using the public `ws.onemindos.dev` (Cloudflare allowlist) or move to a tailnet-only WebSocket | Blueprint lists only `http://localhost:5173/auth/callback`; `ws.onemindos.dev` is on the `nginx` class with the Cloudflare allowlist (onemind-xr `security-and-trust.md` L5) | Pick the hostname with J1; add it to the blueprint's `redirect_uris`. Keep `ws.onemindos.dev` public only once the shared `onemind` password is gone from it (J2-P6) | Localhost only |
| J2-P5 | **Role definitions**: five roles (newcomer, viewer, operator, commander, admin). Viewer may call `svc.policy.action.get` (reads any action whose id it knows; #554 restricts only agents). Operator gets A1–A3 propose/approve and site `sim`; commander adds A4 and second approval; admin equals commander on the bus | `platform/security/roles/onemind-roles.json`; onemind-xr `ONTOLOGY.md` §7.1 (viewer, operator, second approver, owner) | Ratify, or narrow viewer to its own actions (a #554 change) | As listed |
| J2-P6 | The **shared `onemind` password** still works on `ws.onemindos.dev` after the callout is on, so a browser holding it bypasses every role | #562 IDENTITIES.md "Limits" (stage B/C) | After browsers use Authentik logins, give `onemind` `allowed_connection_types: [STANDARD]` (no WEBSOCKET) in the next A19 stage | Unchanged |
| J2-P7 | **Token lifetime**: callout users live until the token's `exp`, at most 1 h; the blueprint sets 15 min access tokens and 8 h refresh tokens. A role removed in Authentik takes effect at the next reconnect (at most 15 min). There is no push revocation in config mode | `max_ttl_seconds`; `TestTokenExpiryDisconnects` | Ratify; revocation lists come with operator/JWT mode | As listed |
| J2-P8 | Group names. The blueprint creates `onemind-newcomers`, `onemind-viewers`, `onemind-operators`, `onemind-commanders`, `onemind-admins`; who is in each is an owner decision in the Authentik UI (the blueprint adds nobody) | `platform/auth-callout/authentik-blueprint.yaml` | Ratify the names; owner adds himself to `onemind-admins` | — |

## Information only

- **Live server supports it.** The live bus runs nats-server 2.14.2 (#562, read-only check of the
  `nats-config` ConfigMap); auth callout has been in nats-server since 2.10 and the integration
  tests run on 2.14.2 itself. WebSocket logins with a token work (`TestNewcomerCannotSeeOrProposeAssetCommands` connects over WebSocket).
- **Password service users are untouched**: with `allowed_accounts: [AUTH]` no configured user
  is ever sent to the callout (`TestPasswordUsersAreNotCalledOut`).
- **OpenFGA model change** (in #554's model, not a fork): `group.member` also accepts
  `group#member`, so role groups nest. Type restrictions still keep agents out of every group.
