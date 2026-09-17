# OneMind OS - Observability Stack

## Overview

Complete real-time monitoring with millisecond-level metrics, long-term storage, custom dashboards, and API access for Legacy.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         OBSERVABILITY ARCHITECTURE                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│   DATA SOURCES              COLLECTION           STORAGE         VISUALIZATION  │
│   ───────────               ──────────           ───────         ─────────────  │
│                                                                                  │
│   ┌──────────┐                                                                   │
│   │Containers│──┐                                                                │
│   │ (Docker) │  │                                                                │
│   └──────────┘  │          ┌─────────────────┐                                   │
│                 │          │                 │                                   │
│   ┌──────────┐  │          │    NETDATA      │                                   │
│   │  System  │──┼─────────►│                 │                                   │
│   │(CPU/RAM) │  │          │ • 1-second data │      ┌─────────────────┐          │
│   └──────────┘  │          │ • Auto-discover │      │                 │          │
│                 │          │ • Anomaly detect│─────►│ VICTORIAMETRICS │          │
│   ┌──────────┐  │          │ • 300+ charts   │      │                 │          │
│   │ Network  │──┤          │                 │      │ • 90-day retain │          │
│   │(packets) │  │          │ GUI: :19999     │      │ • PromQL query  │          │
│   └──────────┘  │          │ API: REST       │      │ • 10x compress  │          │
│                 │          └────────┬────────┘      │                 │          │
│   ┌──────────┐  │                   │               │ GUI: :8428/vmui │          │
│   │  Disks   │──┤                   │               │ API: PromQL     │          │
│   │  (I/O)   │  │                   │               └────────┬────────┘          │
│   └──────────┘  │                   │                        │                   │
│                 │                   │                        │                   │
│   ┌──────────┐  │                   │                        │                   │
│   │Processes │──┘                   │                        │                   │
│   │(per-app) │                      │                        │                   │
│   └──────────┘                      │                        │                   │
│                                     │                        │                   │
│                                     ▼                        ▼                   │
│                            ┌─────────────────────────────────────────┐           │
│                            │              GRAFANA                    │           │
│                            │                                         │           │
│                            │  • Custom dashboards                    │           │
│                            │  • Alerting rules                       │           │
│                            │  • Team sharing                         │           │
│                            │  • Annotations                          │           │
│                            │                                         │           │
│                            │  GUI: :3001                             │           │
│                            │  API: Full REST                         │           │
│                            └─────────────────────────────────────────┘           │
│                                                                                  │
│                                                                                  │
│   ┌──────────────────────────────────────────────────────────────────────────┐  │
│   │                         Legacy API ACCESS                                │  │
│   │                                                                          │  │
│   │   Legacy can query any metric via:                                       │  │
│   │                                                                          │  │
│   │   • Netdata REST API  → Real-time (last 60 seconds, 1s granularity)     │  │
│   │   • VictoriaMetrics   → Historical (PromQL, 90 days)                    │  │
│   │   • Grafana API       → Dashboards, annotations, alerts                  │  │
│   │                                                                          │  │
│   └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Component Details

### Netdata (Real-time Metrics)

| Feature | Details |
|---------|---------|
| **Purpose** | Real-time metrics collection and visualization |
| **Granularity** | 1 second (configurable to 0.5s) |
| **Retention** | ~1 hour in memory |
| **Charts** | 300+ auto-discovered |
| **GUI** | ✅ Beautiful, real-time |
| **Port** | 19999 |
| **Resources** | ~150MB RAM |

**What it monitors automatically:**
- CPU per core, per process
- Memory (RAM, swap, cgroups)
- Disk I/O, space, latency
- Network per interface, per connection
- Docker containers (CPU, memory, network, disk)
- System services
- Application-specific (Redis, Postgres, NATS, etc.)

### VictoriaMetrics (Long-term Storage)

| Feature | Details |
|---------|---------|
| **Purpose** | Long-term metrics storage |
| **Retention** | 90 days (configurable) |
| **Query Language** | PromQL (Prometheus-compatible) |
| **Compression** | 10x better than Prometheus |
| **GUI** | Basic (vmui) |
| **Port** | 8428 |
| **Resources** | ~100MB RAM + disk |

**Why VictoriaMetrics over Prometheus:**
- Same PromQL, better performance
- 10x less storage
- Single binary, simpler
- Drop-in replacement

### Grafana (Dashboards & Alerts)

| Feature | Details |
|---------|---------|
| **Purpose** | Custom dashboards, alerting, visualization |
| **Data Sources** | VictoriaMetrics, Netdata, Postgres, etc. |
| **Dashboards** | Unlimited custom |
| **Alerting** | Built-in with routing |
| **GUI** | ✅ Beautiful, customizable |
| **Port** | 3001 |
| **Resources** | ~200MB RAM |

## Docker Compose Configuration

```yaml
# Add to infra/compose/cloud/docker-compose.yml

services:
  # ===== OBSERVABILITY PROFILE =====

  # Real-time metrics (1-second granularity)
  netdata:
    image: netdata/netdata
    container_name: netdata
    profiles: ["observability"]
    hostname: ${HOSTNAME:-onemind-cloud}
    pid: host
    network_mode: host
    cap_add:
      - SYS_PTRACE
      - SYS_ADMIN
    security_opt:
      - apparmor:unconfined
    volumes:
      - netdata_config:/etc/netdata
      - netdata_lib:/var/lib/netdata
      - netdata_cache:/var/cache/netdata
      - /etc/passwd:/host/etc/passwd:ro
      - /etc/group:/host/etc/group:ro
      - /etc/localtime:/etc/localtime:ro
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /etc/os-release:/host/etc/os-release:ro
      - /var/log:/host/var/log:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - NETDATA_EXTRA_APK_PACKAGES=lm-sensors
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:19999/api/v1/info"]
      interval: 30s
      timeout: 10s
      retries: 3
    # GUI: http://localhost:19999 (via Tailscale)

  # Long-term storage (replaces Prometheus)
  victoriametrics:
    image: victoriametrics/victoria-metrics:stable
    container_name: victoriametrics
    profiles: ["observability"]
    ports:
      - "8428:8428"
    volumes:
      - vm_data:/victoria-metrics-data
    command:
      - "-retentionPeriod=90d"
      - "-storageDataPath=/victoria-metrics-data"
      - "-httpListenAddr=:8428"
      - "-search.latencyOffset=0s"
      - "-search.maxConcurrentRequests=16"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8428/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    # GUI: http://localhost:8428/vmui
    # API: http://localhost:8428/api/v1/query

  # Custom dashboards and alerting
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    profiles: ["observability"]
    ports:
      - "3001:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD:-changeme}
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_ROOT_URL=http://localhost:3001
      - GF_INSTALL_PLUGINS=grafana-clock-panel,grafana-piechart-panel
    restart: unless-stopped
    depends_on:
      - victoriametrics
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    # GUI: http://localhost:3001

  # Uptime monitoring
  uptime-kuma:
    image: louislam/uptime-kuma:latest
    container_name: uptime-kuma
    profiles: ["observability"]
    ports:
      - "3002:3001"
    volumes:
      - uptime_data:/app/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3001"]
      interval: 30s
      timeout: 10s
      retries: 3
    # GUI: http://localhost:3002

volumes:
  netdata_config:
  netdata_lib:
  netdata_cache:
  vm_data:
  grafana_data:
  uptime_data:
```

## Grafana Data Source Configuration

Create `infra/compose/cloud/grafana/provisioning/datasources/datasources.yml`:

```yaml
apiVersion: 1

datasources:
  # VictoriaMetrics (primary for historical)
  - name: VictoriaMetrics
    type: prometheus
    access: proxy
    url: http://victoriametrics:8428
    isDefault: true
    editable: false

  # Netdata (real-time)
  - name: Netdata
    type: prometheus
    access: proxy
    url: http://host.docker.internal:19999/api/v1
    editable: false
```

## Netdata Export to VictoriaMetrics

Configure Netdata to export to VictoriaMetrics:

```yaml
# netdata.conf (or via environment)
[exporting:global]
    enabled = yes

[prometheus_remote_write:victoriametrics]
    enabled = yes
    remote write URL = http://victoriametrics:8428/api/v1/write
    data source = average
    send charts matching = *
```

## API Access for Legacy

### Full Integration Code

```python
# agents/legacy/tools/observability.py
"""
Observability tools for Legacy.

Provides real-time and historical access to all system metrics.
"""

import aiohttp
from datetime import datetime, timedelta
from typing import Optional, List, Dict, Any


class ObservabilityTools:
    """
    Complete observability access for Legacy.

    Three data sources:
    - Netdata: Real-time (1s granularity, last hour)
    - VictoriaMetrics: Historical (PromQL, 90 days)
    - Grafana: Dashboards, annotations, alerts
    """

    def __init__(
        self,
        netdata_url: str = "http://localhost:19999",
        vm_url: str = "http://localhost:8428",
        grafana_url: str = "http://localhost:3001",
        grafana_api_key: Optional[str] = None
    ):
        self.netdata = netdata_url
        self.vm = vm_url
        self.grafana = grafana_url
        self.grafana_key = grafana_api_key

    # =========================================================================
    # NETDATA - Real-time metrics (1-second granularity)
    # =========================================================================

    async def get_system_overview(self) -> Dict[str, Any]:
        """Get current system state - CPU, RAM, disk, network."""
        async with aiohttp.ClientSession() as s:
            results = {}

            # CPU
            async with s.get(f"{self.netdata}/api/v1/data", params={
                "chart": "system.cpu",
                "after": -1,
                "points": 1,
                "format": "json"
            }) as r:
                data = await r.json()
                if data.get("data"):
                    results["cpu_percent"] = sum(data["data"][0][1:])

            # RAM
            async with s.get(f"{self.netdata}/api/v1/data", params={
                "chart": "system.ram",
                "after": -1,
                "points": 1,
                "format": "json"
            }) as r:
                data = await r.json()
                if data.get("data"):
                    results["ram_used_mb"] = data["data"][0][1]
                    results["ram_free_mb"] = data["data"][0][2] if len(data["data"][0]) > 2 else 0

            # Network
            async with s.get(f"{self.netdata}/api/v1/data", params={
                "chart": "system.net",
                "after": -1,
                "points": 1,
                "format": "json"
            }) as r:
                data = await r.json()
                if data.get("data"):
                    results["net_in_kbps"] = abs(data["data"][0][1])
                    results["net_out_kbps"] = data["data"][0][2] if len(data["data"][0]) > 2 else 0

            return results

    async def get_realtime_cpu(self, seconds: int = 60) -> Dict[str, Any]:
        """Get CPU usage at 1-second granularity."""
        async with aiohttp.ClientSession() as s:
            async with s.get(f"{self.netdata}/api/v1/data", params={
                "chart": "system.cpu",
                "after": -seconds,
                "points": seconds,
                "format": "json"
            }) as r:
                return await r.json()

    async def get_container_metrics(self, container: str) -> Dict[str, Any]:
        """Get specific Docker container metrics."""
        async with aiohttp.ClientSession() as s:
            results = {}

            # Container CPU
            async with s.get(f"{self.netdata}/api/v1/data", params={
                "chart": f"cgroup_{container}.cpu_limit",
                "after": -60,
                "format": "json"
            }) as r:
                results["cpu"] = await r.json()

            # Container Memory
            async with s.get(f"{self.netdata}/api/v1/data", params={
                "chart": f"cgroup_{container}.mem_usage",
                "after": -60,
                "format": "json"
            }) as r:
                results["memory"] = await r.json()

            return results

    async def get_all_containers(self) -> List[Dict[str, Any]]:
        """Get all Docker containers and their status."""
        async with aiohttp.ClientSession() as s:
            async with s.get(f"{self.netdata}/api/v1/charts") as r:
                charts = await r.json()

            containers = []
            for chart_id, chart_info in charts.get("charts", {}).items():
                if chart_id.startswith("cgroup_") and ".cpu" in chart_id:
                    name = chart_id.replace("cgroup_", "").replace(".cpu_limit", "").replace(".cpu", "")
                    containers.append({
                        "name": name,
                        "chart_id": chart_id
                    })

            return containers

    async def get_anomalies(self) -> List[Dict[str, Any]]:
        """Get Netdata's built-in anomaly detection alerts."""
        async with aiohttp.ClientSession() as s:
            async with s.get(f"{self.netdata}/api/v1/alarms") as r:
                data = await r.json()

            return [
                {
                    "name": name,
                    "chart": alarm["chart"],
                    "status": alarm["status"],
                    "value": alarm["value"],
                    "units": alarm["units"],
                    "info": alarm.get("info", "")
                }
                for name, alarm in data.get("alarms", {}).items()
                if alarm["status"] not in ["CLEAR", "UNDEFINED"]
            ]

    async def get_charts_list(self) -> List[str]:
        """Get all available Netdata charts."""
        async with aiohttp.ClientSession() as s:
            async with s.get(f"{self.netdata}/api/v1/charts") as r:
                data = await r.json()
                return list(data.get("charts", {}).keys())

    # =========================================================================
    # VICTORIAMETRICS - Historical metrics (PromQL)
    # =========================================================================

    async def query_instant(self, promql: str) -> Dict[str, Any]:
        """Run instant PromQL query."""
        async with aiohttp.ClientSession() as s:
            async with s.get(f"{self.vm}/api/v1/query", params={
                "query": promql
            }) as r:
                return await r.json()

    async def query_range(
        self,
        promql: str,
        start: str = "-1h",
        end: str = "now",
        step: str = "1m"
    ) -> Dict[str, Any]:
        """Run range PromQL query."""
        async with aiohttp.ClientSession() as s:
            async with s.get(f"{self.vm}/api/v1/query_range", params={
                "query": promql,
                "start": start,
                "end": end,
                "step": step
            }) as r:
                return await r.json()

    async def get_cpu_history(self, hours: int = 24) -> Dict[str, Any]:
        """Get CPU usage over time."""
        return await self.query_range(
            "avg(netdata_system_cpu_percentage_average)",
            start=f"-{hours}h",
            step="5m"
        )

    async def get_memory_history(self, hours: int = 24) -> Dict[str, Any]:
        """Get memory usage over time."""
        return await self.query_range(
            "netdata_system_ram_MiB_average{dimension='used'}",
            start=f"-{hours}h",
            step="5m"
        )

    async def get_container_history(
        self,
        container: str,
        metric: str = "cpu",
        hours: int = 24
    ) -> Dict[str, Any]:
        """Get container metric history."""
        if metric == "cpu":
            query = f"netdata_cgroup_{container}_cpu_limit_percentage_average"
        else:
            query = f"netdata_cgroup_{container}_mem_usage_MiB_average"

        return await self.query_range(query, start=f"-{hours}h", step="5m")

    # =========================================================================
    # GRAFANA - Dashboards, annotations, alerts
    # =========================================================================

    def _grafana_headers(self) -> Dict[str, str]:
        """Get Grafana auth headers."""
        if self.grafana_key:
            return {"Authorization": f"Bearer {self.grafana_key}"}
        return {}

    async def get_dashboards(self) -> List[Dict[str, Any]]:
        """Get all Grafana dashboards."""
        async with aiohttp.ClientSession() as s:
            async with s.get(
                f"{self.grafana}/api/search",
                headers=self._grafana_headers()
            ) as r:
                return await r.json()

    async def get_dashboard(self, uid: str) -> Dict[str, Any]:
        """Get specific dashboard by UID."""
        async with aiohttp.ClientSession() as s:
            async with s.get(
                f"{self.grafana}/api/dashboards/uid/{uid}",
                headers=self._grafana_headers()
            ) as r:
                return await r.json()

    async def create_annotation(
        self,
        text: str,
        tags: List[str],
        dashboard_uid: Optional[str] = None
    ) -> Dict[str, Any]:
        """Create annotation on Grafana (mark events on graphs)."""
        payload = {
            "text": text,
            "tags": tags,
            "time": int(datetime.now().timestamp() * 1000)
        }
        if dashboard_uid:
            payload["dashboardUID"] = dashboard_uid

        async with aiohttp.ClientSession() as s:
            async with s.post(
                f"{self.grafana}/api/annotations",
                headers=self._grafana_headers(),
                json=payload
            ) as r:
                return await r.json()

    async def get_alerts(self) -> List[Dict[str, Any]]:
        """Get all Grafana alerts."""
        async with aiohttp.ClientSession() as s:
            async with s.get(
                f"{self.grafana}/api/alerts",
                headers=self._grafana_headers()
            ) as r:
                return await r.json()

    # =========================================================================
    # HIGH-LEVEL ANALYSIS METHODS FOR Legacy
    # =========================================================================

    async def get_health_summary(self) -> Dict[str, Any]:
        """
        Get overall system health summary for Legacy's context.
        Combines real-time and anomaly data.
        """
        overview = await self.get_system_overview()
        anomalies = await self.get_anomalies()
        containers = await self.get_all_containers()

        # Determine overall health
        if any(a["status"] == "CRITICAL" for a in anomalies):
            health = "critical"
        elif any(a["status"] == "WARNING" for a in anomalies):
            health = "warning"
        elif overview.get("cpu_percent", 0) > 80:
            health = "warning"
        elif overview.get("ram_used_mb", 0) > 3500:  # Assuming 4GB VM
            health = "warning"
        else:
            health = "healthy"

        return {
            "overall_health": health,
            "metrics": overview,
            "active_anomalies": len(anomalies),
            "anomalies": anomalies[:5],  # Top 5
            "containers_monitored": len(containers),
            "timestamp": datetime.now().isoformat()
        }

    async def investigate_issue(self, description: str) -> Dict[str, Any]:
        """
        Legacy calls this to investigate a reported issue.
        Returns relevant metrics and potential causes.
        """
        results = {
            "description": description,
            "findings": []
        }

        # Check current state
        overview = await self.get_system_overview()
        anomalies = await self.get_anomalies()

        # CPU issues
        if "slow" in description.lower() or "cpu" in description.lower():
            cpu_history = await self.get_cpu_history(hours=1)
            results["findings"].append({
                "metric": "cpu",
                "current": overview.get("cpu_percent"),
                "history": cpu_history
            })

        # Memory issues
        if "memory" in description.lower() or "ram" in description.lower() or "oom" in description.lower():
            mem_history = await self.get_memory_history(hours=1)
            results["findings"].append({
                "metric": "memory",
                "current": overview.get("ram_used_mb"),
                "history": mem_history
            })

        # Container issues
        if "container" in description.lower() or "docker" in description.lower():
            containers = await self.get_all_containers()
            results["findings"].append({
                "metric": "containers",
                "count": len(containers),
                "list": containers
            })

        results["anomalies"] = anomalies
        results["recommendation"] = self._generate_recommendation(overview, anomalies)

        return results

    def _generate_recommendation(
        self,
        overview: Dict[str, Any],
        anomalies: List[Dict[str, Any]]
    ) -> str:
        """Generate recommendation based on metrics."""
        recommendations = []

        if overview.get("cpu_percent", 0) > 80:
            recommendations.append("High CPU usage detected. Consider scaling or optimizing workloads.")

        if overview.get("ram_used_mb", 0) > 3500:
            recommendations.append("Memory usage is high. May need to increase RAM or reduce container memory limits.")

        for anomaly in anomalies:
            if anomaly["status"] == "CRITICAL":
                recommendations.append(f"Critical alert on {anomaly['chart']}: {anomaly.get('info', 'Check immediately')}")

        return " ".join(recommendations) if recommendations else "All systems operating normally."


# Register with Legacy agent
def register_observability_tools(agent):
    """Register observability tools with Legacy agent."""
    tools = ObservabilityTools()

    agent.register_tool(
        name="get_system_health",
        description="Get current system health summary including CPU, memory, and anomalies",
        func=tools.get_health_summary
    )

    agent.register_tool(
        name="investigate_issue",
        description="Investigate a system issue by analyzing metrics",
        func=tools.investigate_issue
    )

    agent.register_tool(
        name="query_metrics",
        description="Run custom PromQL query on historical metrics",
        func=tools.query_instant
    )

    return tools
```

## Access Summary

| Service | URL (via Tailscale) | GUI | API |
|---------|---------------------|-----|-----|
| **Netdata** | `http://100.x.x.x:19999` | ✅ Beautiful | REST |
| **VictoriaMetrics** | `http://100.x.x.x:8428/vmui` | ⚠️ Basic | PromQL |
| **Grafana** | `http://100.x.x.x:3001` | ✅ Beautiful | Full REST |
| **Uptime Kuma** | `http://100.x.x.x:3002` | ✅ Good | REST |

## Starting the Stack

```bash
# Start observability profile
docker compose --profile core --profile observability up -d

# Check status
docker compose ps

# View Netdata logs
docker logs netdata -f

# Access (via Tailscale)
open http://100.x.x.x:19999  # Netdata
open http://100.x.x.x:3001   # Grafana
```

## Grafana Dashboard Setup

1. Open Grafana: `http://100.x.x.x:3001`
2. Login: admin / (your password)
3. Add Data Source: VictoriaMetrics
4. Import Dashboard: ID `1860` (Node Exporter Full)
5. Or create custom dashboards

## Resource Usage

| Service | RAM | CPU | Disk |
|---------|-----|-----|------|
| Netdata | ~150MB | 1-3% | Minimal |
| VictoriaMetrics | ~100MB | <1% | ~1GB/month |
| Grafana | ~200MB | <1% | ~100MB |
| Uptime Kuma | ~100MB | <1% | ~50MB |
| **Total** | **~550MB** | **~5%** | **~1.2GB/month** |
