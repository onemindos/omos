# IRON — Industrial Automation Platform

**Repo:** (local only — no public GitHub URL found in source)  
**Language:** Rust (edge agent) + Elixir/Phoenix (dashboard) + TimescaleDB  
**Role in OneMind:** STAGED — Phase 3+, industrial/environmental sensor integration

Open-source SCADA platform designed for factories, water treatment, greenhouses, and workshops. Currently in architecture phase — complete engineering blueprint (15 normative specs, 7 ADRs) but no runnable code yet.

## Stack

- Rust edge agent (Modbus, OPC-UA, S7 protocol support)
- Phoenix LiveView dashboard
- TimescaleDB historian
- CLI: `iron new`, `iron dev`, `iron deploy`, `iron validate`

## Why it matters for OneMind

Connects industrial/environmental hardware to the data plane. When PLCs, Modbus devices, greenhouse sensors, or water system controllers arrive, IRON is the bridge.

## Activation condition

No runnable code yet — architecture phase only. Activate when industrial/environmental hardware arrives. The edge agent publishes sensor data into NATS JetStream → ClickHouse time machine.

## Integration point

Edge agent publishes to `tel.{node}.industrial.{metric}` on the FABRIC account.
