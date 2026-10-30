---
name: ha-nats-bridge
description: "Use to control or observe Home Assistant through the OneMind NATS bus (cmd.ha.*, evt.ha.*) and to check the bridge's health."
version: 3.0.0
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

## Leaf troubleshooting (edge side; owner-led)

The edge leaf connects to the hub through the `omos-nats` Tailscale load balancer
on port 7422 (TLS). Error taxonomy from past sessions, so you can tell the owner
what to check:

| Symptom on the edge leaf | Meaning and fix |
|---|---|
| `first record does not look like a TLS handshake` | Plain NATS sent to a TLS port (or the reverse). Use a `tls://` URL; if the hub certificate cannot be verified, `tls { insecure: true }` goes at the `leafnodes` block level, not inside the remote entry. Prefer trusting the NATS CA over `insecure`. |
| `cannot find local account "..." specified in leafnode remote` | Do not set `account:` in the remote unless that local account is defined; the hub assigns the account from the credentials. |
| `Authorization Violation` after the TLS handshake | Credentials rejected or mapped to an account that does not exist on the hub. Compare with the current hub config in `platform/nats/`; do not paste the password anywhere. |
| leaf connects but lands in `$G` | It connected without credentials; events will not route to the agents. |
| `dial tcp 10.x.x.x:7422: i/o timeout` after a hub pod restart | The hub advertised its in-cluster IP to the leaf. `no_advertise: true` in the hub's leafnode block (already in `platform/nats/values.template.yaml`) stops this. |
| Home Assistant add-on container cannot resolve tailnet names | Its DNS is Home Assistant's own; use the load balancer's tailnet IP in the leaf URL instead of the MagicDNS name. |

To diagnose from the hub, the NATS monitoring endpoint (`leafz`, `accountz`) is
reachable from inside the `nats` namespace, not from this agent's pod.

On the Home Assistant OS host (reached only by the owner): processes started
over SSH die when the session closes unless started with `setsid`; `/data` is
root-writable only (use the web terminal); the `ha` CLI works only in the web
terminal because it needs the Supervisor token; and a watchdog should restart
the bridge, because `/tmp` is wiped when the add-on restarts.

## Entities

Do not rely on a remembered entity list. At last check the controllable devices
were Wi-Fi smart plugs (`switch.*`) and the readable ones included a thermometer
and hygrometer, a phone tracker and router WAN status; there were no `light.*`
entities. The Home Assistant box has also been offline for long periods
(`platform/hermes/agent/egress.yaml` notes September 2026), so silence may just
mean the host is down.

## Credentials

The Home Assistant access token and the leaf credentials are secrets held at the
edge and in SealedSecrets. They are not in this skill; never paste them into a
command, a file or a reply. Re-deploying the bridge is an owner-led edge task
(see `infra/edge/` and `docs/` for the edge leaf configuration).
