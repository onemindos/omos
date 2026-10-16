# Machines — Every Device on the Fabric

> Last updated: 2026-08-12
> Author: Legacy (self-authored)

## Tailnet Devices

| Device | Tailscale IP | Type | Role |
|--------|-------------|------|------|
| hermes-legacy (me) | 100.115.219.7 | Cloud pod | Sovereign mind |
| omos-mbp | 100.74.182.117 | Mac M2 Max 32GB | Edge compute, kubectl, local dev |
| omos-ha | 100.64.25.12 | HA Yellow | Home automation, Zigbee/Z-Wave |
| zeus-s24 | — | Samsung S24 Ultra (Android 16) | ATAK, HelmCam, TAK PTT, ADB via tunnel |
| omos-dell-zeus | — | Windows Dell laptop | Remote Windows management |
| omos-nats | 100.117.169.54 | NATS LB | Public NATS entry (leaf port 7422) |
| cloudtak | 100.104.216.93 | k8s service | Cloud TAK COP |
| grafana | 100.80.230.49 | k8s service | Monitoring dashboards |
| hermes | 100.75.131.113 | k8s service | Hermes WebUI |
| argocd | 100.121.32.110 | k8s service | GitOps CD |
| omos | 100.124.142.93 | k8s service | Main OneMind cluster |
| orp | 100.87.177.8 | k8s service | Operations Resource Planner |
| omos-intel | 100.114.216.80 | DO Droplet NYC3 | OSINT dashboards (5 platforms, 12 containers) |
| omos-tablet | 100.91.64.68 | Android tablet | ARCHIVED — offline 24+ days, cert purged |

## Access Patterns

### SSH to Zeus's Mac (omos-mbp)
```bash
ssh -i /opt/data/home/.ssh/id_ed25519 \
  -o StrictHostKeyChecking=no \
  zeus@omos-mbp.hermes.svc.cluster.local "<CMD>"
```

### SSH to Home Assistant (omos-ha)
```bash
ssh -i /opt/data/home/.ssh/id_ed25519 \
  -o StrictHostKeyChecking=no \
  hassio@100.64.25.12 "<CMD>"
```
- Username: `hassio` (NOT root)
- Background processes: `setsid cmd >/tmp/out.log 2>&1 </dev/null &`

### ADB to S24 Ultra
```bash
# Via Tailscale SSH tunnel on Mac (localhost:15555 → phone:5555)
# Cron: s24-pipeline-guardian (every 3min, self-heals)
# Control: /opt/data/scripts/s24_control.py
# Venv: /opt/data/.venv-adb/
```

### kubectl (via pod)
```bash
export USER=legacy; export HOME=/opt/data/home; export PATH="/opt/data/home/.local/bin:$PATH"
kubectl get pods -A
```

## Machine Details

### hermes-legacy (This Pod — Me)
| Property | Value |
|----------|-------|
| OS | Debian 13 (Linux 6.12.95+deb13) |
| Python | 3.13.5 |
| Package manager | uv (no system pip) |
| Persistent volume | /opt/data (NFS) |
| NATS | Leaf sidecar at localhost:4222 |
| Tailscale | Userspace networking (no TUN) |
| Tools | nats, gh, git, tailscale, uv, python3, kubectl, doctl |
| Supervised by | s6-overlay (container) |
| Model | AWS Bedrock — claude-sonnet-4-6 |

### omos-mbp (Zeus's MacBook Pro)
| Property | Value |
|----------|-------|
| Hardware | Apple M2 Max, 32GB RAM |
| OS | macOS (latest) |
| Python | 3.14 (Homebrew) |
| Hermes | launchd (ai.hermes.gateway) — legacy-edge profile |
| NATS | Local server at :4222 (leaf to cluster) |
| kubectl | Direct cluster RBAC (do-nyc3-onemind-hub context) |
| Role | Edge compute, local dev, kubectl proxy |

### zeus-s24 (Samsung S24 Ultra)
| Property | Value |
|----------|-------|
| OS | Android 16 (SDK 36) |
| ADB | Via Pod→Tailscale→SSH tunnel (localhost:15555→phone:5555) |
| Apps | ATAK, HelmCam, Tailscale, TAK-ML, CloudTAK |
| Video | TAKCam+HelmCam→ATAK→RTSP→MediaMTX (10.108.44.200:8554) |
| TAK cert | Zeus-S24.p12 (password not recorded here) |
| Role | Primary tactical edge device, PTT, situational awareness |
| ADB key | Mac's key |

### omos-dell-zeus (Windows Laptop)
| Property | Value |
|----------|-------|
| OS | Windows |
| Access | impacket / remote Windows management |
| Role | Windows-native dev, secondary workstation |
| SSH | Blocked by power management (NIC drop issue) |

### omos-ha (Home Assistant Yellow)
| Property | Value |
|----------|-------|
| Hardware | HA Yellow (CM4-based) |
| OS | HA OS |
| LAN IP | <omos-ha LAN IP> |
| Tailscale | 100.64.25.12 |
| Zigbee | Built-in coordinator |
| NATS leaf | /data/nats-leaf.conf → cluster:7422 |
| HA Bridge | /data/ha_bridge.py (Python, publishes to events.ha.*) |
| Entities | 185+ |

### omos-intel (DO Droplet)
| Property | Value |
|----------|-------|
| Region | NYC3 |
| Size | s-4vcpu-8gb |
| OS | Debian 13 |
| Public IP | <omos-intel public IP> |
| Tailscale | 100.114.216.80 |
| SSH | ssh -i ~/.ssh/id_ed25519 root@omos-intel |
| Running | ShadowBroker, WorldMonitor, Osiris, WorldWideView, Crucix (12 containers) |

### omos-tablet (Android Tablet)
| Property | Value |
|----------|-------|
| Tailscale | 100.91.64.68 |
| Status | ARCHIVED — offline 24+ days, cert purged, archive recommended |
| Action | Remove from active fabric, decommission |

---

*Every machine is a peer on the fabric. Legacy can reach them all.*
