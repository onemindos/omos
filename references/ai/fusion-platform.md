# fusion-platform — NATS JetStream Sensor Fusion Platform

**Repo:** https://github.com/bragind/fusion-platform
**Language:** Python 3.11 + C++17
**License:** MIT

## What it is

A production-ready multi-sensor fusion platform built on NATS JetStream as the message bus. Combines IMU, GPS, LiDAR, camera, and telemetry data through an Extended Kalman Filter (EKF) in real time. Designed for drones, AGVs, underwater vehicles, and industrial diagnostics.

## Architecture

```
SENSORS (sim or real)
  IMU 100Hz │ GPS 10Hz │ LiDAR │ Telemetry
            ▼
  NATS JETSTREAM (message bus)
  Persistence · Replay · Guaranteed Delivery
            ▼
  PROCESSING PIPELINE
  Time Sync → Noise Filter (3σ) → Kalman EKF (C++/Eigen)
            ▼
  SERVING LAYER
  REST API │ WebSocket │ Live Dashboard
            ▼
  OBSERVABILITY
  Prometheus │ Grafana │ Alerts
```

## Tech Stack

| Component | Tech |
|-----------|------|
| Kalman core | C++17, Eigen 3.4, pybind11 |
| Services | Python 3.11, FastAPI, NumPy, asyncio |
| Message bus | NATS JetStream |
| Replay | Apache Parquet, PyArrow |
| Containers | Docker, Docker Compose |
| Monitoring | Prometheus, Grafana |

## Why it's relevant to OneMind

| Pattern | OneMind application |
|---------|---------------------|
| Sensors publish to NATS JetStream | Same as OneMind `tel.{node}.{class}.{metric}` — drop-in pattern for drone/robot telemetry feeds |
| Extended Kalman Filter | Fuse MAVLink IMU + GPS from the drone bridges into a clean position estimate — critical for constellation tracking |
| Replay via Apache Parquet | Records raw sensor streams, replays at configurable speed — directly analogous to ClickHouse Time Machine time-travel |
| NATS subject per sensor type | `tel.imu`, `tel.gps`, `tel.lidar` — mirrors OneMind subject grammar |
| Prometheus + Grafana metrics | `fusion_latency_seconds`, `sensor_drop_total`, `kalman_innovation_magnitude` — ready to wire into OneMind's Grafana stack |
| C++/Python bridge (pybind11) | Pattern for embedding high-performance C++ compute inside Python NATS services |

## NATS subjects used

```
tel.imu       IMU data (100 Hz)
tel.gps       GPS data (10 Hz)
tel.lidar     LiDAR data
tel.telemetry General telemetry
```

## Key files

| File | Purpose |
|------|---------|
| `processing/kalman/kalman_core.cpp` | C++ Extended Kalman Filter (Eigen) |
| `processing/kalman/fusion_service.py` | Python wrapper — subscribes to NATS, runs EKF, publishes fused state |
| `processing/time_sync/sync_service.py` | Aligns async sensor streams by timestamp |
| `processing/noise_filter/outlier_filter.py` | Median + 3σ outlier rejection before Kalman |
| `replay/recorder.py` | Writes raw NATS messages to Parquet |
| `replay/controller.py` | HTTP API — replay a time window at N× speed |
| `serving/api.py` | REST + WebSocket serving fused state + covariance |
| `docker/docker-compose.dev.yml` | Full stack: NATS + Prometheus + Grafana |

## REST API

| Method | URL | Returns |
|--------|-----|---------|
| `GET` | `/state` | Current fused state `[x, y, vx, vy]` |
| `GET` | `/covariance` | Covariance matrix diagonal |
| `GET` | `/health` | Service status |
| `GET` | `/metrics` | Prometheus metrics |

## Replay API

```bash
curl -X POST http://localhost:8001/replay \
  -H "Content-Type: application/json" \
  -d '{"start_time": 1700000000.0, "end_time": 1700000060.0, "speed": 2.0}'
```

Recordings saved to `data/recordings/` as Parquet files.

## Prometheus metrics

| Metric | Type | What it measures |
|--------|------|-----------------|
| `fusion_latency_seconds` | Histogram | EKF predict/update cycle latency |
| `sensor_drop_total` | Counter | Lost sensor messages |
| `kalman_innovation_magnitude` | Gauge | Filter innovation — spikes = sensor anomaly |

## Quick start

```bash
git clone https://github.com/bragind/fusion-platform.git
cd fusion-platform
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt

# Set NATS URL in config.json
echo '{"nats_url": "nats://localhost:4222"}' > config.json

# Run full stack
python sensors/imu_sim.py &
python sensors/gps_sim.py &
python processing/time_sync/sync_service.py &
python processing/noise_filter/outlier_filter.py &
python serving/api.py &
```

## Notes

- Cloned locally at `ignore /fusion-platform/` in omos-core local repo (now deleted)
- README is in Russian but code and comments are standard
- The Parquet replay pattern is the most immediately useful piece for OneMind — same concept as ClickHouse time-travel but at the raw sensor level
- C++ Kalman core requires CMake + Eigen to build; Python services work standalone
