# F1. The `intercept` example in the grammar block

Source: onemind-xr `docs/design/JADC2-MAP.md` §6 F1 (branch `run/m0`). Index and ruling order:
[README.md](README.md).

**Proposed.**
1. In the `CLAUDE.md` grammar block, line 88, replace the example `cmd.drone.x47.intercept` with
   `cmd.mission.patrol1.abort` (doctrine §9.1 mission form, four tokens), and add the fleet line
   from BUS-P2 with the example `cmd.fleet.drone.x47.goto`.
2. In `services/schemas/command/envelope.schema.json` line 63, remove `cmd.intercept.*` from the
   "protected command classes" (which today means "needs approval") and say instead that
   engagement verbs (`intercept`, `strike`, `engage`, `designate`) are refused outright: the asset
   verb list is closed (stop line SL5) and has none of them (SL1).

**Why.** The canonical example of the command pattern is an engagement verb, and the schema
treats `intercept` as allowed with approval. Both conflict with stop lines SL1 ("no weapons; no
strike, engage or designate action") and SL5 (closed verb list `goto hold rtl land orbit observe
estop arm disarm takeoff estop_release geofence_override`), onemind-xr `docs/design/JADC2-MAP.md`
§5.2. The checker flags the example with a stop-line word.

**Note on JADC2-MAP's wording.** JADC2-MAP F1 suggested `cmd.drone.x47.goto`. That would show
an asset command in the four-token form that BUS-P2 replaces with `cmd.fleet.{class}.{id}.{verb}`,
so this entry uses a mission example for the four-token line and puts the drone example on the
fleet line. The intent (no engagement verb) is unchanged.

**Alternatives considered.** Keep the example and rely on policy refusing it (the grammar block is
what new work copies; rejected). Use the doctrine §2 example `cmd.le.light.on` (valid, but a
OneMind reader needs a mission or fleet example).

**Interim reading.** Nothing publishes `intercept`. #554 refuses any action or task type that
contains a stop-line word (`services/policy/internal/policy/catalogue.go` `stopLineWords`,
`TestStopLine`), and #550 refuses verbs outside its catalogue (`unknown_action`, e2e test with
`strike`).

**Related finding (not in the M0 register, recorded for the owner, no change here).**
`services/nex/overwatch/om-overwatch-interceptor/main.go` describes itself as "the effector.
Closes the kill chain … INTERCEPT" and uses out-of-grammar `int.{kind}.{id}` subjects (listed under
BUS-P4). It has no Argo Application (not deployed), but its KV bucket `OVERWATCH_INTERCEPT`
(`main.go` line 38) is declared in `platform/nats/streams/kv-overwatch-intercept.yaml`. It conflicts with SL1 in the same way.
Suggested follow-up for the owner: park or remove it in a separate PR.

**Depends on it.** None blocks; doc and schema text only. Consistent with #554 and #550.

```subjects
legacy    cmd.drone.x47.intercept
legacy    cmd.intercept.*
proposed  cmd.mission.patrol1.abort
proposed  cmd.fleet.drone.x47.goto
current   cmd.le.light.on
alternative  cmd.drone.x47.goto
legacy    int.command.{id}
legacy    int.{kind}.{id}
proposed  cmd.fleet.{class}.{id}.{verb}
```

- [ ] Approve as proposed
- [ ] Approve with changes:
- [ ] Reject
