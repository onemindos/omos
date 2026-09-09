# wlow — DAG Workflow Orchestrator on NATS

**Repo:** (third-party)  
**Language:** Rust  
**Role in OneMind:** Reference for NEX-style task execution with process/WASM/Firecracker tiers over NATS

DAG-based workflow orchestrator. Define a graph of tasks, submit it, wlow runs each task in the correct sandbox with dependency tracking, retries, and a final result.

## Execution tiers

| Runtime | How | KVM required |
|---------|-----|-------------|
| `process` | OS subprocess | No — `wlow start` on any NATS node |
| `wasm` | Wasmtime component | No — `wlow start --runtimes wasm` |
| `microvm` | Firecracker VM | Yes — KVM host required |

## Why it matters for OneMind

Directly parallels OneMind's NEX execution engine concept — workloads running on NATS nodes in different isolation tiers. wlow adds:
- DAG dependency graph (vs flat task queue)
- Firecracker microVM isolation (stronger than WASM sandbox)
- Retries + result tracking built-in

## Integration point

```
cmd.task.{id}.submit  → wlow job submitted (NATS)
evt.task.{id}.done    ← wlow result published (NATS)
```

DAG maps directly onto OneMind's `om-task` NEX workload. wlow's Firecracker tier is the path to full process isolation for untrusted agent workloads.

## See Also

- `om-task` in [[omos-services]] — OneMind's own NEX WASI task engine
- NEX (nats-io/nex) in [[nats-ecosystem-references]] — NATS execution engine
