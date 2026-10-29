---
name: nats-operations
description: "Use for NATS bus operations: publishing, subscribing, health checks and diagnosing bridge or leaf connectivity on the OneMind bus."
version: 2.0.0
tags: [nats, messaging, bus, events, pub-sub]
triggers:
  - need to publish or subscribe on NATS
  - check bus health
  - monitor events
  - agent communication
---

# NATS Operations

The in-cluster NATS (3 nodes, JetStream, namespace `nats`) is the authority bus.
Everything publishes to it. Full doctrine: `docs/NATS-DOCTRINE.md` in the repo.

## Connection

- Hermes is the one documented exception to "connect to `nats.nats.svc`": its pod
  runs a `nats-leaf` sidecar and the agent uses `NATS_URL=nats://localhost:4222`,
  so it has its own bus identity. Check the variable before assuming:
  `echo $NATS_URL`.
- Credentials come from the pod's environment; never print them.
- What this identity may publish or subscribe to is set by account permissions.
  A permission error is a finding to report, not something to work around. Whether
  JetStream is available to this account is something to verify
  (`nats account info`), not assume.

## Locked subject grammar (do not invent patterns)

```
tel.{node}.{class}.{metric}     telemetry
det.{domain}.{type}.{id}        detections
evt.{domain}.{entity}.{event}   semantic events
cmd.{domain}.{entity}.{action}  commands
intel.{domain}.{source}.{type}  intelligence
ent.{entity_type}.{entity_id}.state  fused entity state
svc.{service}.{endpoint}        NATS RPC
agents.{verb}.{agent}.{owner}.* agent fabric (agents.prompt.hermes.<owner>.<session>)
tak.cot.inject                  CoT in to TAK (NatsBridgePlugin input)
tak.cot.>                       CoT out from TAK (cot-tap)
```

A new stream is a new CR file in `platform/nats/streams/` (applied by Argo, reconciled
by NACK), never `nats stream add`.

## Common operations

```bash
nats account info          # connection and account
nats rtt                   # round trip (iterations are positional)
nats server info

timeout 5 nats sub "tak.cot.>"          # TAK events
timeout 5 nats sub "agents.prompt.>"    # agent prompts
nats pub <subject> "<payload>"          # only on subjects you are authorized for
```

- Never subscribe to `>` on a busy bus except for a few seconds with `timeout`.
- `nats req` needs a unit on the timeout: `--timeout 15s`.
- Under `cmd.>` a JetStream stream may capture the subject, so `nats req` can
  return the stream ack instead of the service reply. Use `nats pub` plus a
  subscription to the result subject (see `ha-nats-bridge`).
- `nats bench` has subcommands (`pub`, `sub`, `service`, `js`, `kv`). Shell-measured
  publish latency of tens of milliseconds is process start-up, not the wire.

## TAK over NATS

The NatsBridgePlugin inside TAK Server subscribes to `tak.cot.inject` and injects
the CoT XML into TAK; TAK-originated events come out on `tak.cot.out` (the `COT` stream captures `tak.cot.>`).
A per-UID projection `evt.tak.cot.{uid}` is proposed in the doctrine but not live;
confirm direction with a short subscription before relying on it. Details: `tak-server`.

## Diagnosing a silent service or leaf

```bash
nats micro ls                      # services that respond
nats micro info <service>
nats req svc.<service>.health "" --timeout 3s   # only if the service defines it
```

The NATS monitoring port (8222) is reachable on localhost of the server or sidecar,
not on pod IPs. From the sidecar's own pod, `curl localhost:8222/connz?subs=1` and
`/leafz?detail=1` show connections and the leaf link to the cluster.

Leaf and bridge problems seen before:

- After a hub pod restart, leaf nodes can be handed an internal ClusterIP in the
  server's advertised endpoints and then fail to reconnect over Tailscale. The
  fix is on the hub side (do not advertise internal endpoints to leaves); report
  it rather than patching the hub by hand.
- Edge leaves reach the cluster through the `omos-nats` Tailscale load balancer
  (leaf port 7422). If the edge is offline the cloud bus keeps running.
- Never point a cloud workload at an edge or Mac NATS.
