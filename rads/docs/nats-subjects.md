# RADS NATS Subject Reference

All subjects follow the pattern: `rad.<homestead_id>.<subsystem>.<...>`

Default homestead_id: `greenwich`

## Land / Soil

| Subject | Direction | Schema | Description |
|---------|-----------|--------|-------------|
| `rad.greenwich.land.field.1.soil` | PUB | `schemas/soil_reading.json` | Field 1 soil reading |
| `rad.greenwich.land.field.*.soil` | SUB wildcard | — | All fields |
| `rad.greenwich.land.field.1.irrigate` | PUB | `schemas/land_decision.json` | Irrigate field 1 |
| `rad.greenwich.land.decisions` | PUB | `schemas/land_decision.json` | Decision audit log |

## Robots / Drones

| Subject | Direction | Schema | Description |
|---------|-----------|--------|-------------|
| `rad.greenwich.robots.drone1.telemetry` | PUB | `schemas/drone_telemetry.json` | Drone 1 telemetry |
| `rad.greenwich.robots.drone1.cmd` | SUB | — | Drone 1 commands |

## Legacy Override

| Subject | Direction | Description |
|---------|-----------|-------------|
| `rad.greenwich.legacy.override` | SUB | Legacy AI override channel for land agent |

## JetStream Streams

| Stream | Subjects | Retention |
|--------|----------|-----------|
| `RADS_LAND` | `rad.greenwich.land.>` | 10 years |
| `RADS_ROBOTS` | `rad.greenwich.robots.>` | 10 years |
