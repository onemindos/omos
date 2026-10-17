# OneMind Event Envelope — Specification

> **Status:** v1, canonical. Every service, agent, bridge, and the Flutter UI speaks this envelope.
> Defined by card P0-E. Stored live in `om-schema` KV bucket under key `envelope.v1`.

---

## What Is The Envelope

Every message on the OneMind NATS bus is wrapped in this versioned envelope. The envelope is
the outer structure; `data` is the typed payload that varies by event type. This separation
lets all consumers share the same framing logic and parse the payload according to `type`.

---

## Field Table

| Field       | Type    | Required | Description |
|-------------|---------|----------|-------------|
| `v`         | integer | **yes**  | Schema version. Starts at 1. Consumers pin a minimum `v` and tolerate unknowns. Breaking changes bump `v`. |
| `id`        | string  | **yes**  | Unique event ID — UUID v4. |
| `ts`        | string  | **yes**  | Event timestamp — RFC3339 / ISO8601 with millisecond precision (e.g. `2026-06-10T14:30:00.123Z`). |
| `type`      | string  | **yes**  | Dotted event type mirroring the NATS subject grammar (e.g. `evt.task.issue.created`). |
| `source`    | string  | **yes**  | Actor ID that emitted this event (e.g. `agent.legacy`, `node.barn-pi-1`). |
| `subject`   | string  | **yes**  | The NATS subject this envelope was published on. May include dynamic tokens. |
| `data`      | object  | **yes**  | Typed payload. Schema varies by `type`. Parse according to `type`. |
| `corr`      | string  | no       | Correlation ID (UUID v4). Groups related events across a workflow or request chain. |
| `causation` | string  | no       | Causation ID (UUID v4). The `id` of the parent event that caused this one. |
| `trace`     | string  | no       | Distributed trace ID for OpenTelemetry / Jaeger integration. |

**Forward-compatibility rule:** `additionalProperties: true` at the envelope root. Consumers
**MUST** tolerate unknown fields — never reject an envelope because it has extra fields.

---

## Versioning Policy

- **Consumers pin a minimum `v`**, not an exact version. A consumer that requires `v >= 1`
  accepts any envelope with `v` equal to or greater than 1.
- **Tolerate unknowns.** The envelope deliberately allows additional properties. A consumer
  that rejects unknown fields is broken by design.
- **Backward-compatible changes** (adding optional fields, relaxing constraints) do NOT bump
  `v`. New optional fields just appear — consumers that don't know them ignore them.
- **Breaking changes** (renaming/removing required fields, tightening constraints, changing
  field semantics) MUST bump `v` and ship a migration note in this README and in the
  `om-schema` KV bucket entry's description.

### Migration note template

When `v` bumps from N to N+1, add a section here:

```
## v1 → v2 migration (TEMPLATE)
- Changed: `source` renamed to `origin` (was ambiguous with NATS connection source).
- Action: consumers reading `source` must also check `origin` for v2+ envelopes.
- Grace period: both fields present in v2 until 2027-01-01, then `source` removed in v3.
```

---

## Worked Examples

### 1. Task Event — `evt.task.issue.created`

A task event emitted when a new OneMind issue is created (e.g. by Legacy after receiving a
request from Zeus). This event is recorded forever in ClickHouse and triggers any autopilot
rules watching for new issues.

```json
{
  "v": 1,
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "ts": "2026-06-10T14:30:00.123Z",
  "type": "evt.task.issue.created",
  "source": "agent.legacy",
  "subject": "evt.task.issue.created",
  "corr": "a1b2c3d4-e5f6-4a7b-8c9d-000000000001",
  "data": {
    "issue_id": "OM-001",
    "title": "Re-roof the barn section",
    "assignee": "node.barn-pi-1",
    "priority": "high",
    "status": "open",
    "created_by": "agent.legacy",
    "tags": ["infrastructure", "barn"]
  }
}
```

**Notes:**
- `corr` groups this event with other events in the same user request chain.
- `data.assignee` is an actor ID (matches the polymorphic actor model).
- ClickHouse ingests this on `evt.>` subscription.

---

### 2. Telemetry Event — `tel.barn-pi-1.soil.moisture`

A telemetry event from a barn sensor. Below the `threshold_min`, an autopilot would
create an issue assigned to the irrigation device actor.

```json
{
  "v": 1,
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "ts": "2026-06-10T14:32:15.456Z",
  "type": "tel.barn-pi-1.soil.moisture",
  "source": "node.barn-pi-1",
  "subject": "tel.barn-pi-1.soil.moisture",
  "trace": "abc123def456ghi789",
  "data": {
    "sensor_id": "soil-1",
    "value": 18.3,
    "unit": "%",
    "threshold_min": 20,
    "alert": true,
    "location": "greenhouse-row-3"
  }
}
```

**Notes:**
- `trace` enables correlation with any distributed trace (e.g. from an OpenTelemetry-wired node).
- `type` mirrors the NATS subject grammar: `tel.{node}.{class}.{metric}`.
- FABRIC account publishes; CORE and INTERFACE subscribe via import.

---

### 3. Agent Event — `agents.status.legacy.zeus.sess-001`

An agent status event per Synadia Agent Protocol v0.3 subject grammar. The `causation` field
links back to the task event that triggered this agent session.

```json
{
  "v": 1,
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "ts": "2026-06-10T14:35:00.789Z",
  "type": "agents.status.legacy.zeus.sess-001",
  "source": "agent.legacy",
  "subject": "agents.status.legacy.zeus.sess-001",
  "corr": "a1b2c3d4-e5f6-4a7b-8c9d-000000000002",
  "causation": "550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "agent": "legacy",
    "status": "idle",
    "session": "sess-001",
    "capabilities": ["code", "plan", "review", "kanban"],
    "owner": "zeus",
    "model": "claude-sonnet-4.6"
  }
}
```

**Notes:**
- `causation` traces this status back to the issue-created event (example 1) that spawned the session.
- `type` uses the Synadia Agent Protocol v0.3 subject grammar: `agents.{verb}.{agent}.{owner}.{session}`.
- The Flutter app subscribes on `agents.status.>` to show live agent state.

---

## Schema File

The JSON Schema (draft 2020-12) lives at `services/schemas/envelope/envelope.schema.json`.
The live copy is seeded into the NATS KV registry: `nats kv get om-schema envelope.v1`.

## Codegen

Typed models are generated from the schema. See `services/schemas/codegen/README.md` for instructions.

- **Dart** → `services/schemas/codegen/out/envelope.dart` (consumed by `onemindos-command/contract/`)
- **Python** → `services/schemas/codegen/out/envelope.py` (consumed by core services)
