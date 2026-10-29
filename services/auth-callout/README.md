# auth-callout (FEATURES J2)

NATS auth callout for browser and remote clients. A client connects to NATS (WebSocket or plain
NATS) with its **Authentik OIDC token** as the CONNECT `auth_token`; the server asks this service
(`$SYS.REQ.USER.AUTH`, account `AUTH`) and the service answers with a short-lived user that the
server places in the `CLIENT` account with the permissions of one OneMind role. No shared
password ever reaches a browser.

```
browser ──PKCE──► Authentik (onemind-xr)        role groups in the `groups` claim
   │  token
   └──NATS WebSocket CONNECT {auth_token}──► nats-server ──$SYS.REQ.USER.AUTH──► auth-callout
                                                  ▲                                   │ verify RS256/ES256 via JWKS,
                                                  └──── user JWT (aud CLIENT, ◄────────┘ map groups → role
                                                        allow lists, exp ≤ token, ≤ 1 h)
CLIENT user "ak:olga;onemind-operators" ──svc.policy.command.propose──► policy service
   (server-set Nats-Request-Info) → user:olga + contextual tuple member group:onemind-operators → OpenFGA
```

## What it decides

| Input | Rule |
|---|---|
| Token | Compact JWS, `alg` RS256 or ES256 only (no `none`, no HS*), `kid` in the provider JWKS (refetched on an unknown kid, at most every 30 s; cached keys survive an Authentik outage), exact `iss`, `aud` contains `onemind-xr`, `exp`/`nbf`/`iat` with 30 s leeway |
| Username | `preferred_username` must match `^[a-z0-9][a-z0-9._@-]{0,63}$` (it becomes the OpenFGA id `user:<username>`) |
| Role | The highest-ranked role whose Authentik group is in `groups` ([role map](../../platform/security/roles/onemind-roles.json)). No OneMind group: refused |
| User JWT | `aud` = `CLIENT` (config-mode placement); `name` = `ak:<username>;<role groups>`; publish = the role's allow list, deny `$JS.>`, `$KV.>`, `$SYS.>`, `_INBOX.>`; subscribe = `_INBOX.<sub>.>` plus the role's list; connection types WEBSOCKET and STANDARD; `exp` = token `exp`, at most `max_ttl_seconds` (1 h). The server disconnects the client when it expires |
| Refusal | The reply carries only `not authorized`; the reason is logged (`username`, `reason`, `client_type`, `host`) |
| Leaf, password, nkey logins | Refused if they ever arrive (they do not, see the platform README) |

Clients must use the inbox prefix `_INBOX.<sub>` (`sub` claim of their own token;
nats.ws `inboxPrefix`, nats.go `CustomInboxPrefix`). That stops one `CLIENT` user from reading
another's replies (A19-P5), which would otherwise let a newcomer read an operator's policy
replies.

### Why the role groups are in the user name

In config mode nats-server copies only the user name into `Nats-Request-Info`: the `tags` of a
callout user JWT are applied on the operator-mode JWT path only (nats-server v2.14.2
`server/auth.go`, `c.tags = juc.Tags`), so they never reach the policy service. The name is
set by the server from the callout's answer and a client cannot change it, so the groups travel
there (`PrincipalName`). Proposal J2-P2.

## Run and test

```bash
cd services/auth-callout
go vet ./...
go test -race -count=1 ./...      # in-process nats-server v2.14.2 (the live version)
```

The tests start a disposable server with the stage-1 account layout (A19), the `AUTH` account and
the `auth_callout` block from [`platform/auth-callout/nats-callout.conf`](../../platform/auth-callout/nats-callout.conf),
a stand-in Authentik JWKS, WebSocket and leaf listeners. Every key and password is generated per
run; nothing is committed. They cover the FEATURES J2 acceptance test (a newcomer reaches no
`svc.policy.*` endpoint and cannot subscribe to anything that shows a command), viewer and
operator rights, the verified identity at the policy responder (and that a forged
`Nats-Request-Info` is replaced), inbox isolation, eleven kinds of bad token, token expiry
disconnect, callout outage (fails closed; password users unaffected), JWKS outage and key
rotation, xkey-encrypted callouts, password and leaf logins (which never reach the callout,
except stage-A2 leafs, pinned in `TestLeafNodes`), and that the deployed ConfigMap and the
Authentik blueprint match the role map.

Configuration (environment): `NATS_URL`, `NATS_USER`, `NATS_PASSWORD`, `NATS_CA_FILE`,
`CALLOUT_ISSUER_SEED_FILE` (account seed `SA...`), `CALLOUT_XKEY_SEED_FILE` (curve seed `SX...`,
optional), `OIDC_ISSUER`, `OIDC_AUDIENCE`, `OIDC_JWKS_URL`, `ROLES_FILE`, `HEALTH_ADDR`.
Deployment and owner steps: [`platform/auth-callout/README.md`](../../platform/auth-callout/README.md).
