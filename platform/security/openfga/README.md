# OpenFGA actions model (proposal, FEATURES I4)

**Not applied.** `platform-security` reads `platform/security/` without recursion, so nothing in
this directory is synced. The live store "OneMind Operations" keeps the model written by
`platform/ecosystem/resources/authorization-context.yaml` (types `user`, `group#member`,
`capability#runner`, used by Backstage Missions only) until the owner adopts this one.

| File | What it is |
|---|---|
| `onemind-actions.fga` | The model in the OpenFGA DSL (schema 1.1), for review |
| `onemind-actions.json` | The same model in the API form (`POST /stores/{id}/authorization-models` body). `services/policy` loads it for its in-process evaluator and its tests, and `TestModelDSLMatchesJSON` fails if the two files disagree |

## What it adds

The live types stay unchanged (`user`, `group`, `capability#runner`), so the Backstage Missions
checks keep working. New types, from onemind-xr `docs/design/ONTOLOGY.md` section 7.2:

| Type | Relations | Why |
|---|---|---|
| `agent`, `service` | — | agents and services become principals (`agent:hermes/zeus/legacy`, `service:watch`) |
| `site` | `operator` | operate every asset of a site |
| `asset` | `site`, `operator` (user, group member, or user with `within_window`; or operator of its site), `viewer` | asset ids are `sim-*` in this run |
| `action_type` | `proposer` (user, group member, agent, service), `approver` and `second_approver` (users and group members only), `executor` (service) | **an agent cannot be written as an approver**: OpenFGA rejects the tuple at write time (type restriction), tested in `TestModelTypeRestrictions` |
| `mission`, `team` | owner/member/viewer, leader/member | for A2 mission and team actions |

Condition `within_window(current_time, grant_time, grant_duration)` time-boxes an operator
grant; the policy service passes `current_time` on every Check.

## Tuples the policy service needs (example, for the owner to write)

```
user:zeus                 operator         asset:sim-x1
user:zeus                 proposer         action_type:asset.goto      (and every catalogue action)
user:zeus                 approver         action_type:asset.goto
user:zeus                 second_approver  action_type:asset.arm       (A4 verbs; a second person too)
agent:hermes/zeus/legacy  proposer         action_type:asset.goto
service:watch             proposer         action_type:task.assign
```

## Adopting it (owner)

1. Review the model. Then change the bootstrap Job's `definitions` in
   `authorization-context.yaml` to this JSON (or add a second Job), bump its name
   (`onemind-authorization-bootstrap-v2`) so Argo re-runs it, and write the tuples above.
2. Put the new model id (and the store id) into `platform/policy/configmap.yaml`
   (`OPENFGA_MODEL_ID`, `OPENFGA_STORE_ID`).
3. OpenFGA runs one replica. The policy service fails closed when it is unreachable: every
   proposal is denied `policy_unavailable`.
