---
title: "Graphite Bridge"
type: reference
status: active
tags: [platform, infrastructure, monitoring, metrics, netdata, clickhouse]
---
# Graphite Bridge

> Pipes Netdata system metrics into ClickHouse for permanent storage.

## Role in OneMind OS
- Receives metrics from Netdata via Graphite plaintext protocol (TCP :2003)
- Batches and inserts into ClickHouse HTTP API every 10 seconds
- Enables SQL queries over years of system metrics history

## Access
- TCP listener: `127.0.0.1:2003`
- Source: `/Users/onemind/onemind/subscribers/graphite_bridge.py`
- Logs: `/Users/onemind/onemind/subscribers/logs/graphite-bridge.log`
- LaunchAgent: `com.onemind.graphite-bridge`

## Data Flow

```
Netdata (collector, 1s resolution)
  → Graphite exporter (every 60s, configured in /opt/homebrew/etc/netdata/exporting.conf)
    → Graphite Bridge (TCP :2003)
      → Batches every 10s
        → ClickHouse HTTP :8124 → onemind.metrics table
```

## Volume
- ~3,800 metrics per minute (all Netdata charts)
- ~5.5M metrics/day
- With ClickHouse compression: ~50-100MB/month

## Configuration
- Netdata exporting: `/opt/homebrew/etc/netdata/exporting.conf`
- ClickHouse target: `http://127.0.0.1:8124` (onemind.metrics table)

## Dependencies
- Netdata (metric collector)
- ClickHouse (storage on :8124)
- Python httpx library
