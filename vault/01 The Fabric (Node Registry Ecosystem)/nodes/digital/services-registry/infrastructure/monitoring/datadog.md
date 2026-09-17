---
title: "Datadog"
type: reference
status: active
tags: [platform, infrastructure, monitoring, observability, apm, logs]
---
# Datadog

> Cloud monitoring and observability platform. Metrics, traces, logs, APM, synthetics.

## Key Features
| Feature | Notes |
|---------|-------|
| Infrastructure | Host metrics, containers, cloud |
| APM | Distributed tracing |
| Logs | Centralized log management |
| Dashboards | Custom visualizations |
| Alerts | Threshold, anomaly, forecast |
| Synthetics | API + browser tests |
| RUM | Real user monitoring |
| Security | SIEM, threat detection |

## Agent Install
```bash
# macOS
DD_API_KEY=<key> DD_SITE="datadoghq.com"   bash -c "$(curl -L https://install.datadoghq.com/scripts/install_mac_os.sh)"

# Docker
docker run -d --name datadog-agent   -e DD_API_KEY=<key>   -e DD_SITE="datadoghq.com"   -v /var/run/docker.sock:/var/run/docker.sock:ro   datadog/agent:latest
```

## Integration
- Monitor NATS, Docker, cloud infra
- APM for Hermes/OpenClaw services
- Logs → centralized search
- Alerts → NATS → Hermes notification
