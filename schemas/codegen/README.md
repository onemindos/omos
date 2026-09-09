# OneMind Schema Codegen

Generates typed models from `schemas/envelope/envelope.schema.json`.

## Generators

| Script | Output | Target repo |
|--------|--------|-------------|
| `gen_python.py` | `out/envelope.py` | `onemindos-core` services |
| `gen_dart.py` | `out/envelope.dart` | `onemindos-command/contract/` |

## How to run

From the repo root:

```sh
# Python model (stdlib only — no deps needed)
python3 schemas/codegen/gen_python.py

# Dart model (stdlib only — no deps needed)
python3 schemas/codegen/gen_dart.py
```

If `datamodel-code-generator` is installed (`pip install datamodel-code-generator`), `gen_python.py`
will use it for a richer output with validators, aliases, and JSON serialization. Otherwise it
falls back to a hand-rolled dataclass that is functionally complete.

## Dependency policy

Generators are **stdlib only** by default. No heavy build-time deps. If a richer generator is
available it is used opportunistically; absence is not a failure. This keeps the `schemas/`
directory usable on fresh checkouts and in CI without a pre-install step.

## Consuming in onemindos-command

Copy `out/envelope.dart` to `contract/envelope.dart` in the `onemindos-command` repo. Pin the
envelope version you depend on in a companion `contract/VERSION` file (e.g. `envelope.v1`).

## Adding new schemas

1. Add `schemas/{domain}/{name}.schema.json`.
2. Add a `gen_{lang}.py` script in `schemas/codegen/` following the same pattern.
3. Add a row to this README's table.
4. Seed the live registry: `nats kv put om-schema {name}.v1 "@schemas/{domain}/{name}.schema.json"`.
