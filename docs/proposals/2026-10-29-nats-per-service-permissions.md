# Proposal: per-service NATS users and permissions (stage 1)

Status: **proposal for the owner**, 2026-10-29. Nothing here is live. No credentials are in this
file; every password is a placeholder filled from a SealedSecret or the gitignored NATS Helm
`values.secret.yaml`. Source design: onemind-xr `docs/design/BUS-DESIGN.md` sections 7.1–7.3
(branch `run/m0`); filed with FEATURES I4 (policy service, `services/policy`).

## Why

Today all application traffic runs in one account `ONEMIND` on shared passwords (onemind-xr
`docs/research/nats-platform.md` §1). Any holder of a shared login can publish `cmd.fleet.>`
and bypass policy, and any reply inbox can be answered by anyone. The policy service makes the
policy decision, but the server must make the bypass impossible: the policy service is the only
publisher of `cmd.>` (except `cmd.map.>`), and only a verified human can reach
`svc.policy.*.confirm`.

## Accounts

| Account | Members | Exports / imports |
|---|---|---|
| `SYS` | ops, monitor | — |
| `ONEMIND` | every service, the TAK leaf, site leafs (keeps every stream; no migration) | exports service `svc.policy.>` (and later `svc.entity.>`, `svc.layer.>`, `svc.cot.send`) |
| `CLIENT` | humans: Unreal, L0 web, Quest, ATAK plugin; browsers later through Authentik and auth callout | imports `svc.policy.>` from `ONEMIND` **with `share: true`**; no JetStream; no data publish |
| `AGENTS` | Hermes and other harnesses | imports only `svc.policy.command.propose`, `svc.policy.task.propose`, `svc.policy.action.get` with `share: true`; `agents.>`; `cmd.map.>` |

### How the policy service knows who is asking

On a service import with `share: true`, nats-server adds the header `Nats-Request-Info` with the
requester's account and user, and **replaces** any value the client sent (nats-server v2.14.8
`server/client.go`, service import processing; `server/accounts.go` `ClientInfoHdr`). The policy
service maps `{account, user}` to a principal (`platform/policy/configmap.yaml`
`principals.json`) and treats every identity from `AGENTS` as an agent, whatever the map says.
Requests from inside `ONEMIND` carry no server-set identity: they are recorded
`actor_verified: false` and can never confirm.

Proved in `services/policy/internal/policy/service_test.go` against an in-process server with
exactly this layout: a CLIENT human is verified; a forged header from an AGENTS user is
replaced; an AGENTS user without the confirm import is refused by the server; a misconfigured
agent that reaches confirm is refused by the policy (`agent_cannot_confirm`); CLIENT and AGENTS
users get a permission violation on `cmd.fleet.>`; an ONEMIND service cannot publish to
`svc.policy.>`.

**Residual risk:** a user inside `ONEMIND` that may publish `svc.policy.*.propose` could set
its own `Nats-Request-Info` header and appear as a CLIENT user on a *proposal*. It still cannot
confirm (no `ONEMIND` user may publish to `svc.policy.*.confirm`). Keep the list of ONEMIND
proposers short (watch jobs, the task service) and move agents to `AGENTS` from day one.
Stage 2 (operator/JWT and auth callout, doctrine §5) removes it.

## Per-identity allow-lists

`_INBOX.>` and each identity's own JetStream API calls are granted per identity; responders get
`allow_responses: true` (needed to answer requests that arrive through an import, whose reply
subject is rewritten by the server).

| Identity (account) | Publish | Subscribe |
|---|---|---|
| `policy` (ONEMIND) | `cmd.>`, `evt.so.action.>`, `svc.entity.get`, `$JS.API.>`, `_INBOX.>`; `allow_responses` | `svc.policy.>`, `evt.drone.>`, `evt.robot.>`, `evt.vehicle.>`, `evt.mission.>`, `evt.entity.>`, `_INBOX.>` |
| `entity` (ONEMIND) | `ent.*.*.state`, `evt.entity.>`; `allow_responses` | `svc.entity.>`, `cmd.entity.>`, its inputs (`tak.cot.out`, `tel.fleet.>`, `intel.>`, `det.>`, …) |
| `task` (ONEMIND) | `evt.mission.>`, KV `mission-state`; `svc.policy.task.propose` (decomposed children); `allow_responses` | `cmd.mission.>`, `svc.mission.>`, `evt.so.action.>`, `evt.drone.>`, `evt.robot.>` |
| `watch` (ONEMIND) | `svc.policy.command.propose`, `svc.policy.task.propose` | its inputs, `_INBOX.>` |
| `bridge-mavlink-<site>` (ONEMIND, on the site leaf) | `tel.fleet.drone.<id>.>`, `evt.drone.<id>.>` for **each** of its ids (NATS wildcards cannot match `sim-*` partially) | `cmd.fleet.drone.<id>.>` per id, no wildcard over ids |
| `bridge-ros-<site>` | same with `robot` | same |
| `projector` (ONEMIND) | `tak.cot.inject` | `ent.*.*.state`, KV `geofences`, `mission-state`; serves `svc.cot.send` |
| `tak-leaf` (ONEMIND) | `tak.cot.out` | `tak.cot.inject` |
| collectors (ONEMIND, one user per group) | `intel.>`, `evt.world.>`, `tel.feed.>`, `evt.sys.*.rejected` | KV `ingest-locks`, `feed-config` |
| every other ONEMIND service | its own subjects; **deny** `svc.policy.>`, `cmd.>`, `evt.so.action.>` | its own subjects |
| human users (CLIENT) | `svc.policy.>`, `svc.entity.>`, `svc.layer.>`, `cmd.map.>` | `_INBOX.>`, `cmd.map.>` |
| `hermes` (AGENTS) | `svc.policy.command.propose`, `svc.policy.task.propose`, `svc.policy.action.get`, `agents.>`, `cmd.map.>` | `_INBOX.>`, `agents.prompt.hermes.>` |

## Config skeleton (stage 1, config mode)

Placeholders `<...>` come from the gitignored `values.secret.yaml` or a SealedSecret. This is the
shape the tests run (`natsConfig` in `harness_test.go`), not a file to apply.

```
accounts {
  ONEMIND {
    jetstream: enabled
    users: [
      { user: policy, password: "<from secret>", permissions: {
          publish:   { allow: ["cmd.>", "evt.so.action.>", "svc.entity.get", "$JS.API.>", "_INBOX.>"] }
          subscribe: { allow: ["svc.policy.>", "evt.drone.>", "evt.robot.>", "evt.vehicle.>", "evt.mission.>", "evt.entity.>", "_INBOX.>"] }
          allow_responses: true } }
      { user: bridge-mavlink-home, password: "<from secret>", permissions: {
          publish:   { allow: ["evt.drone.sim-x1.>", "tel.fleet.drone.sim-x1.>"] }
          subscribe: { allow: ["cmd.fleet.drone.sim-x1.>"] } } }
      # ... one entry per service, as the table above
    ]
    exports: [ { service: "svc.policy.>" } ]
  }
  CLIENT {
    users: [ { user: zeus, password: "<from secret>", permissions: {
        publish: { allow: ["svc.policy.>", "svc.entity.>", "svc.layer.>", "cmd.map.>"] }
        subscribe: { allow: ["_INBOX.>", "cmd.map.>"] } } } ]
    imports: [ { service: { account: ONEMIND, subject: "svc.policy.>" }, share: true } ]
  }
  AGENTS {
    users: [ { user: hermes, password: "<from secret>", permissions: {
        publish: { allow: ["svc.policy.command.propose", "svc.policy.task.propose", "svc.policy.action.get", "agents.>", "cmd.map.>"] }
        subscribe: { allow: ["_INBOX.>", "agents.prompt.hermes.>"] } } } ]
    imports: [
      { service: { account: ONEMIND, subject: "svc.policy.command.propose" }, share: true }
      { service: { account: ONEMIND, subject: "svc.policy.task.propose" }, share: true }
      { service: { account: ONEMIND, subject: "svc.policy.action.get" }, share: true }
    ]
  }
}
```

## Migration order (owner)

1. Add the `policy`, `entity`, `task` and bridge users inside `ONEMIND` next to the existing
   shared users (nothing breaks; nothing is removed).
2. Add `CLIENT` and `AGENTS` with their imports; move Hermes to `AGENTS` (its leaf sidecar
   credential changes; `platform/hermes/agent`).
3. Deny `cmd.>` (except `cmd.map.>`) and `svc.policy.*.confirm` for every remaining shared
   ONEMIND user. From this step the server enforces "only the policy service publishes
   commands".
4. Per-account JetStream limits (today all `-1`).

Each step is a NATS Helm values change (the release is still Helm CLI, omos `CLAUDE.md`), so it
needs the owner. The policy service runs before step 3, but until then the bypass is closed only
by convention.
