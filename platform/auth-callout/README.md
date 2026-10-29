# Role modes with Authentik: NATS auth callout (FEATURES J2)

Status, 2026-10-29: **not live**. Nothing here changes the bus until the owner runs the steps
below. Service code and tests: [`services/auth-callout`](../../services/auth-callout/README.md).

## The five roles

Source: [`platform/security/roles/onemind-roles.json`](../security/roles/onemind-roles.json)
([README](../security/roles/README.md)). One Authentik group per role; a user gets the highest
role among their groups.

| Role | Authentik group | Client modes | NATS (CLIENT account, request/reply only) | OpenFGA (actions model, #554) |
|---|---|---|---|---|
| newcomer | `onemind-newcomers` | explore | `svc.layer.*`, `svc.entity.get/search`. **No `svc.policy.*`, no `cmd.*`, no `agents.*`**: cannot see or propose asset commands | nothing |
| viewer | `onemind-viewers` | + watch | + `svc.entity.query/history`, `svc.policy.action.get`; receives `cmd.map.>` | nothing beyond newcomer |
| operator | `onemind-operators` | + operate | + `svc.policy.{command,task}.{propose,confirm}`, `svc.policy.check`, `svc.cot.send`, `agents.prompt.hermes.>`, `cmd.map.>` | proposer and approver of A1–A3 actions; operator of site `sim` |
| commander | `onemind-commanders` | + command | as operator | + proposer, approver and second approver of A4 |
| admin | `onemind-admins` | + admin | as commander | as commander (granting roles happens in Authentik, not on the bus) |

Two independent layers enforce it: the **server** (the user JWT's allow lists; a newcomer's
`svc.policy.*` request is a `Permissions Violation` and never reaches the policy service) and
the **policy service** (OpenFGA Check with the user's role groups as contextual tuples; a
newcomer or viewer who somehow reached it is denied `not_authorized`). Agents can never hold a
role: OpenFGA refuses an agent as a group member.

## Files

| File | Applied by | What |
|---|---|---|
| `deployment.yaml`, `configmap.yaml`, `network-policy.yaml` | `platform-auth-callout` (manual sync) | The callout, 2 replicas, ecosystem pool, namespace `onemind-core` |
| `authentik-blueprint.yaml` | `platform-auth-callout` (ConfigMap in namespace `auth`) | Groups, `groups` claim, public OIDC client `onemind-xr`; inert until step 5 |
| `nats-callout.conf` | owner, by hand (excluded from Argo) | `AUTH` account, callout user, `authorization.auth_callout` block |
| `identities.fragment.yaml` | nobody yet (excluded) | The same identity in the A19 spec format, to move into `platform/nats/identities/identities.yaml` once #562 merges |
| [`../security/openfga/onemind-roles.tuples.json`](../security/openfga/onemind-roles.tuples.json) | owner, with the model | Role grants for OpenFGA (generated) |

## What the callout touches

Verified against nats-server v2.14.2 source (`server/auth.go`, `server/auth_callout.go`) and in
the integration tests on an in-process 2.14.2 server:

- **Service users (password, any account):** never sent to the callout, because
  `allowed_accounts: [AUTH]` (`TestPasswordUsersAreNotCalledOut`). An empty `allowed_accounts`
  would send every password user to the callout and lock every service out.
- **Unknown logins** (an OIDC token, or a wrong password) go to the callout. A wrong password is
  then refused by the callout instead of by the server; the result is the same.
- **Leafs through `leafnodes { authorization { users } }`** (today's `leaf_tak`, A19 stage A1):
  never sent to the callout.
- **Leafs that log in as account users** (A19 stage A2): always sent to the callout, which
  refuses them. **The callout and A19 stage A2 cannot both be on** (proposal J2-P1).
- **Callout down:** token logins fail closed after the server's auth timeout; everything else
  is unaffected (`TestCalloutDownFailsClosed`).

## Owner steps

Prerequisites: #554 (policy service and actions model) and #562 (accounts, stage A1) merged and
stage A1 live; `kubeseal`; Go (for `nk`); a private directory outside the repo
(`W=~/onemind-private/nats`, as in `platform/nats/IDENTITIES.md`).

1. **Keys** (outside the repo):
   ```bash
   cd "$W"
   go run github.com/nats-io/nkeys/nk@v0.4.16 -gen account -pubout > callout-issuer.nk   # line 1 seed SA..., line 2 public A...
   go run github.com/nats-io/nkeys/nk@v0.4.16 -gen curve -pubout   > callout-xkey.nk     # SX... / X...
   python3 -c 'import secrets;print(secrets.token_urlsafe(32))' > auth-callout.pw
   ```
2. **SealedSecrets** (namespace-wide scope, as the repo `CLAUDE.md`), committed beside this file:
   ```bash
   kubectl create secret generic auth-callout-runtime -n onemind-core \
     --from-literal=NATS_USER=auth-callout --from-file=NATS_PASSWORD="$W/auth-callout.pw" \
     --dry-run=client -o json | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/auth-callout/auth-callout-runtime.sealedsecret.yaml
   kubectl create secret generic auth-callout-keys -n onemind-core \
     --from-literal=issuer.seed="$(sed -n 1p "$W/callout-issuer.nk")" \
     --from-literal=xkey.seed="$(sed -n 1p "$W/callout-xkey.nk")" \
     --dry-run=client -o json | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/auth-callout/auth-callout-keys.sealedsecret.yaml
   ```
3. **NATS: add the `AUTH` account and user only** (part 1 of `nats-callout.conf`, the password
   as a bcrypt hash: `nats server passwd -p "$(cat "$W/auth-callout.pw")"`). Apply it like an
   A19 stage (`platform/nats/IDENTITIES.md` step 3: edit the saved `live.conf`, `kubectl apply`
   the ConfigMap, wait for `Reloaded server configuration` on nats-0..2). Nothing else changes.
4. **Deploy the callout:** build on main (`build-auth-callout.yml`), pin the digest in
   `deployment.yaml`, then sync `platform-auth-callout` by hand. Expected:
   `kubectl -n onemind-core logs deploy/auth-callout` shows `auth-callout ready ... "roles":5,"xkey":true`.
5. **Authentik:** add to `platform/argocd/platform-auth-helm.yaml` values
   `blueprints: {configMaps: [onemind-xr-blueprint]}` (check the key against the goauthentik
   chart 2026.8.2 values first) in a PR. After it syncs, Authentik → System → Blueprints shows
   "OneMind XR role modes" as successful, and Applications lists `onemind-xr`. Put yourself in
   `onemind-admins`. Check the token:
   `curl -s https://auth.onemindos.dev/application/o/onemind-xr/.well-known/openid-configuration | jq .issuer`
   prints the `OIDC_ISSUER` in `configmap.yaml`, and the JWKS lists an RSA key.
6. **OpenFGA:** with the actions model (#554 `platform/security/openfga/README.md`), also write
   `onemind-roles.tuples.json` and the asset-site tuples of the simulated assets
   (`site:sim site asset:sim-x1`). `platform/policy/configmap.yaml` already lists the callout
   account.
7. **Turn the callout on:** add part 2 of `nats-callout.conf` (`authorization.auth_callout`)
   with the two public keys, apply, wait for the reload on all three pods. Check that nothing
   else changed:
   ```bash
   kubectl -n onemind-core logs deploy/auth-callout --since=10m | grep -c '"msg":"denied"'   # expect 0 or only your tests
   for p in nats-0 nats-1 nats-2; do kubectl -n nats logs $p -c nats --since=10m | grep -E 'Authorization Violation|callout'; done
   ```
   Every service keeps running (none is called out). Then the J2 acceptance test from a browser
   or with the nats CLI and a token (`OWNER-STEPS.md` in onemind-xr `docs/evidence/M4/J2/`).

**Rollback:** remove the `authorization { auth_callout }` block and apply (reload). Token logins
stop; nothing else is affected. Removing the `AUTH` account is a second, separate reload.

## Limits

- Config mode only: the role groups ride in the user name (J2-P2); tags and per-account
  callouts arrive with operator/JWT mode (BUS-DESIGN §7, stage 2).
- Not compatible with A19 stage A2 leaf users (J2-P1).
- An agent proposing on behalf of a callout user is checked without that user's groups
  (J2-P3).
- The shared `onemind` password still works on `ws.onemindos.dev` until A19 stage B/C closes
  it; the callout adds a safe path, it does not remove the unsafe one (J2-P6).
