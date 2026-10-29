---
name: ha-nats-bridge
description: "Use to control or observe Home Assistant through the OneMind NATS bus (cmd.ha.*, evt.ha.*) and to check the bridge's health."
version: 2.0.0
triggers:
  - "HA bridge"
  - "home assistant NATS"
  - "cmd.ha"
  - "evt.ha.state_changed"
  - "turn on the light"
  - "switch the plug"
---

# HA to NATS Bridge

Home Assistant is controlled and observed only through the bus. The bridge runs
at the edge next to Home Assistant, on a NATS leaf node that forwards to the
cluster through the `omos-nats` Tailscale load balancer. Direct REST from the
cloud to Home Assistant is not exposed. The edge being offline means the bridge
is down; the cloud bus keeps running.

## Subject contract

| Direction | Subject | Payload |
|-----------|---------|---------|
| command | `cmd.ha.service.call` | `{"domain":"switch","service":"turn_on","entity_id":"<entity>"}` |
| query | `cmd.ha.state.get` | `{"entity_id":"<entity>"}` |
| live stream | `evt.ha.state_changed.<entity_id>` | `{"entity_id","state","attributes","last_changed"}` |
| command result | `evt.ha.cmd.result.<entity_id>` | `{"status":200,"entity_id":"...","confirmed_state":"on","elapsed_ms":600}` |

## Use

```bash
# Fire a command (publish, then watch for the confirmation)
nats pub cmd.ha.service.call '{"domain":"switch","service":"turn_on","entity_id":"<entity>"}'
timeout 20 nats sub 'evt.ha.cmd.result.<entity_id>' --count 1
```

- Use `nats pub` plus a subscription. Do not use `nats req` for state queries:
  a JetStream stream captures `cmd.>` and the reply is the stream ack, not the state.
- Discover entities from the live stream (`timeout 10 nats sub 'evt.ha.state_changed.>'`)
  rather than from a remembered list; devices come and go. Controllable devices
  have mostly been `switch.*` entities.
- Physical latency (hundreds of milliseconds for a WiFi plug) is the device,
  not the bus.
- Acting on locks, alarms or anything that affects physical safety needs the
  owner's explicit approval for that action.

## Health

- No `evt.ha.state_changed.>` traffic for several minutes means the bridge, the
  edge leaf or the edge machine is down. Confirm from the leaf link on the
  cluster side first (`nats` monitoring `leafz` on the hub), then ask the owner
  to check the edge host. This agent does not hold a shell on the edge host.
- Bridge behavior to expect: infinite NATS reconnect, a 30 s Home Assistant
  WebSocket heartbeat with exponential backoff (2 s to 60 s), a 10 s timeout on
  service calls and a 5 s timeout on state reads.

## Credentials

The Home Assistant access token and the leaf credentials are secrets held at the
edge and in SealedSecrets. They are not in this skill; never paste them into a
command, a file or a reply. Re-deploying the bridge is an owner-led edge task
(see `infra/edge/` and `docs/` for the edge leaf configuration).
