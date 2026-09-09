# Control-Plane Schemas

Canonical JSON Schema stubs for the OneMind control plane. These define the
contract between backend services and client apps (Flutter / omos-atoc).

## Ownership

**Backend owns these schemas.** The `omos-core` repo is the single source of
truth for all registry, UI projection, and command envelope definitions.

- **Backend services** (`om-control-plane`, `om-tak-projector`, `om-router`,
  `om-connector-manager`) implement against these schemas.
- **Clients** (`omos-atoc` / Flutter) generate or copy typed models from these
  schemas. Clients are **never** the source of truth for registry state or
  projection shapes.
- **Changes** require a schema version bump and a migration note. See the
  versioning policy in `schemas/envelope/README.md`.

## Structure

```
schemas/
  control-plane/
    README.md            ← this file
  registry/
    entity.schema.json       EntityRecord   — registry.entity.{id}
    source.schema.json       SourceRecord   — registry.source.{id}
    service.schema.json      ServiceRecord  — registry.service.{id}
    integration.schema.json  IntegrationRecord — registry.integration.{id}
    workflow.schema.json     WorkflowRecord — registry.workflow.{id}
  ui/
    now-summary.schema.json  UiNowSummary   — ui.now.summary
    asset.schema.json        UiAsset        — ui.assets.index / ui.assets.{id}.state
    system.schema.json       UiSystem       — ui.systems.index
    feed-item.schema.json    UiFeedItem     — ui.feed.recent / ui.feed.important
  command/
    envelope.schema.json     CommandEnvelope — cmd.* subjects
```

## Doctrine references

- `docs/CONTROL-PLANE.md` — control plane services, registries, entity model,
  command model, UI projection rules.
- `docs/TAK-PROJECTION.md` — TAK visibility, CoT mapping, show/hide commands.
- `docs/WORKFLOW-ROUTING.md` — workflow record shape, routing rules, protected
  command classes.
- `docs/NATS-DOCTRINE.md` — subject grammar, JetStream rules, account model.
- `docs/PR-ROADMAP.md` — scoped PR plan (this work is PR 1 / issue #6).

## Key rules

1. **No UI firehose.** Clients consume bounded `ui.*` projections for routine
   operation. Raw `evt.>`, `tel.>`, `tak.cot.>`, `$SYS.>` subjects are for
   backend processing and advanced inspector views only.
2. **Flutter is not source-of-truth.** Registry state lives in JetStream KV
   buckets managed by backend services. The app renders projections.
3. **Every action is a command; every result is an event.** UI buttons publish
   commands on `cmd.*` subjects. Backend validates, executes, and emits result
   events (`evt.command.accepted`, `evt.command.executed`, etc.).
4. **Forward compatibility.** All schemas use `additionalProperties: true` at
   the root. Consumers must tolerate unknown fields.

## Status

These are **stubs** (v1). Field sets are practical starting points derived from
the doctrine docs. They will be refined as the control-plane services are built
(PR 2+). No runtime behavior depends on them yet.

## TODOs

- [ ] Wire codegen (`schemas/codegen/`) to emit Dart + Python models for these
      schemas once the envelope codegen is extended to multi-schema.
- [ ] Seed live KV registry entries (`nats kv put om-schema ...`) when the
      control-plane services are deployed.
- [ ] Add `$ref` cross-references between schemas (e.g. UiAsset referencing
      EntityRecord) once a shared definitions file is needed.
