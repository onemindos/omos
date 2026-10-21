# `schemas/`

**The Contract.** Versioned JSON Schemas for the OneMind OS event envelope and all typed payloads. Everything that flows on the NATS bus conforms to these schemas.

> The Flutter app (`onemindos-command`) depends directly on this folder. Do not change envelope field names or types without a version bump and a migration plan.

---

## Structure

### `envelope/`
The core P0-E event envelope — the wrapper around every message on the bus.

```json
{
  "v":       1,           // envelope version
  "id":      "uuid-v4",   // unique message id
  "ts":      "RFC3339",   // authoritative timestamp (ms precision)
  "type":    "evt.task.issue.created",  // event type = NATS subject
  "source":  "actor-id",  // originating actor
  "subject": "evt.task.issue.created",  // NATS subject it arrived on
  "data":    { ... }      // typed payload (schema per event type)
}
```

Files:
- `envelope.schema.json` — JSON Schema definition
- `examples/` — example envelope payloads
- `README.md` — envelope field definitions and rules

### `codegen/`
Code generators that produce typed classes from the JSON Schemas:

| Script | Output |
|--------|--------|
| `gen_dart.py` | Dart classes for the Flutter app |
| `gen_python.py` | Python dataclasses for services + subscribers |

Generated output lives in `codegen/out/` — do not hand-edit, regenerate instead.

---

## Adding a new event type

1. Add a JSON Schema to `envelope/` (or a new subfolder per domain)
2. Run codegen: `python3 services/schemas/codegen/gen_dart.py && python3 services/schemas/codegen/gen_python.py`
3. Commit both the schema and the generated output
4. Update `docs/NATS-DOCTRINE.md §subject-grammar` if a new subject prefix is introduced
