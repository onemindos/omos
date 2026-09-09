# OneMind Control Plane

## Purpose

OneMind needs one backend control plane that makes the bus understandable to operators and clients.

The bus remains the source of truth for real-time communication. The control plane does not replace NATS, TAK, NEX, om-unity, Legacy, or JetStream. It observes them, normalizes their state, maintains registries, and publishes bounded UI projections for clients.

## Doctrine

- NATS is the backbone.
- TAK is the primary field/tactical map.
- `omos-atoc` is the control-plane/workflow/AI console.
- Legacy reasons over selected context and acts through commands.
- Flutter must not consume the raw firehose for routine UI.
- Backend services publish clean `ui.*` projections for app views.
- Every action is a command; every result is an event.
- Every connected thing is an entity, source, service, integration, workflow, or capability.

## Core services

### `om-control-plane`

Owns the system registries and UI projections.

Inputs:

- NATS monitor API (`:8222`): `varz`, `connz`, `jsz`, `accountz`, `healthz`
- NEX workload status
- JetStream streams and KV buckets
- `agents.hb.>` and `agents.status.>`
- `ent.*.*.state`
- `evt.>`, `cmd.>`, `voice.>`, `task.>`, `evt.mission.>`
- om-unity API (`:9090`) entity graph
- TAK bridge/server health
- Hermes/Legacy dashboard/heartbeat
- connector/bridge status

Outputs:

- `ui.now.summary`
- `ui.assets.index`
- `ui.assets.{id}.state`
- `ui.sources.index`
- `ui.services.index`
- `ui.integrations.index`
- `ui.feed.recent`
- `ui.feed.important`
- `ui.work.tasks.board`
- `ui.work.missions.active`
- `ui.systems.index`
- `ui.insights.summary`
- `ui.protect.summary`
- `ui.legacy.status`

### `om-tak-projector`

Controls what is projected into TAK.

Inputs:

- `registry.entity.{id}` / KV entity registry
- `ent.*.*.state`
- visibility commands: `cmd.entity.{id}.tak.*`, `cmd.layer.{layer}.tak.*`

Outputs:

- `tak.cot.inject`
- `evt.tak.projected`
- `evt.tak.hidden`
- `ui.layers.tak`

Rule: an entity appears in TAK only when `visibility.tak == true` and policy allows it.

### `om-router`

Workflow/routing engine.

Inputs:

- Event streams (`evt.>`, `ent.*.*.state`, `cmd.>`, `voice.>`, `agents.hb.>`)
- Workflow registry

Outputs:

- Commands (`cmd.*`)
- alerts (`ui.alerts.*`, `evt.alert.*`)
- tasks/missions (`cmd.task.*`, `cmd.mission.*`)
- Legacy prompts (`agents.prompt.hermes.*`) only when a rule requires AI.

### `om-connector-manager`

Manages integrations and connector lifecycle.

Responsibilities:

- credential collection and validation
- entity discovery
- bridge start/stop/restart
- connector health
- source registry updates
- integration UI projections

Examples: Home Assistant, Todoist, Google Calendar, Slack, Discord, GitHub, Tailscale, WorldWideView, WorldMonitor, Osiris, om-unity feeds.

## Registries

Canonical registries should live in JetStream KV and emit updates.

Suggested KV buckets:

- `KV_ENTITY_REGISTRY`
- `KV_SOURCE_REGISTRY`
- `KV_SERVICE_REGISTRY`
- `KV_INTEGRATION_REGISTRY`
- `KV_WORKFLOW_REGISTRY`
- `KV_CAPABILITY_REGISTRY`
- `KV_COMMAND_POLICY`

Suggested subjects:

- `registry.entity.{id}`
- `registry.source.{id}`
- `registry.service.{id}`
- `registry.integration.{id}`
- `registry.workflow.{id}`
- `registry.capability.{id}`

## Entity model

Every connected/operational object is an entity.

Examples:

- `robot.go2-01`
- `drone.raven-01`
- `agent.legacy`
- `node.zeus-mac`
- `node.barn-pi-01`
- `sensor.greenhouse-temp-01`
- `camera.gate-rtsp-01`
- `ship.mmsi-636092200`
- `tak.user.omos-iphone17`

Minimum fields:

```json
{
  "id": "node.zeus-mac",
  "kind": "machine",
  "name": "zeus-mac",
  "status": "online",
  "source": "omd",
  "subjects": {
    "telemetry": "tel.node.zeus-mac.>",
    "commands": "cmd.node.zeus-mac.>"
  },
  "visibility": {
    "tak": false,
    "atoc": true,
    "feed": true,
    "legacy_watch": false
  },
  "tak": {
    "uid": "OM.node.zeus-mac",
    "cot_type": "a-f-G-U-i-I",
    "layer": "infrastructure",
    "label": "ZEUS-MAC"
  },
  "capabilities": ["health", "logs", "restart_workload"],
  "last_seen": "2026-07-03T00:00:00Z"
}
```

## Command model

UI buttons publish commands; backend validates and emits events.

Examples:

- `cmd.entity.{id}.tak.show`
- `cmd.entity.{id}.tak.hide`
- `cmd.entity.{id}.legacy.watch`
- `cmd.entity.{id}.feed.mute`
- `cmd.integration.homeassistant.test`
- `cmd.integration.homeassistant.discover`
- `cmd.integration.homeassistant.start`
- `cmd.workflow.create`
- `cmd.workflow.enable`
- `cmd.task.create`
- `cmd.mission.create`
- `cmd.nex.workload.restart.{id}`
- `cmd.legacy.ask`

Results/events:

- `evt.command.accepted`
- `evt.command.rejected`
- `evt.command.executed`
- `evt.command.failed`
- `evt.entity.tak.visible`
- `evt.entity.tak.hidden`
- `evt.integration.connected`
- `evt.integration.failed`
- `evt.workflow.triggered`

## UI projection rule

Flutter and other clients consume `ui.*` projections for normal operation.

Allowed routine UI subjects:

- `ui.now.summary`
- `ui.feed.recent`
- `ui.assets.index`
- `ui.assets.{id}.state`
- `ui.services.index`
- `ui.integrations.index`
- `ui.work.*`
- `ui.systems.*`
- `ui.insights.*`
- `ui.protect.*`

Raw firehose subjects (`evt.>`, `tel.>`, `tak.cot.>`, `$SYS.>`) are for backend processing and advanced inspector views only.

## Phase 1 MVP

Build one thin vertical slice:

1. `om-control-plane` publishes:
   - `ui.now.summary`
   - `ui.assets.index`
   - `ui.systems.index`
   - `ui.feed.recent`
2. `om-control-plane` discovers:
   - NATS server health
   - NEX workloads
   - current NATS connections
   - om-unity entity count
   - TAK bridge/server health
   - Hermes/Legacy health
3. `om-tak-projector` supports:
   - `cmd.entity.{id}.tak.show`
   - `cmd.entity.{id}.tak.hide`
4. `omos-atoc` consumes only the `ui.*` projection subjects for the normal operator views.

