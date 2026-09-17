---
id: ui-15-20260111053926
title: "Skill: Infrastructure Monitoring"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/15]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Skill: Infrastructure Monitoring"
---

# Skill: Infrastructure Monitoring

## Purpose

Monitor all technology infrastructure across the One Mind ecosystem. Ensure systems are operational, secure, and performing optimally. Proactively identify issues before they become problems.

## Knowledge

### Infrastructure Components

| Component | Location | Priority | SLA |
|-----------|----------|----------|-----|
| Homelab Server | On-premise | Critical | 99.5% |
| Cloud VPS | DigitalOcean/AWS | High | 99.9% |
| Letta Instance | Cloud/Self-hosted | Critical | 99.5% |
| N8N | Self-hosted | High | 99% |
| Obsidian Sync | iCloud/Git | Medium | 95% |
| Backups | Multi-location | Critical | 100% |

### Services Stack

| Service | Purpose | Port | Health Check |
|---------|---------|------|--------------|
| Letta | AI agent memory | 8283 | `/health` |
| N8N | Workflow automation | 5678 | `/healthz` |
| Langfuse | LLM observability | 3000 | HTTP 200 |
| Uptime Kuma | Status monitoring | 3001 | Self |
| Prometheus | Metrics collection | 9090 | `/ready` |
| Grafana | Dashboards | 3002 | `/api/health` |

### Alert Thresholds

| Metric | Warning | Critical |
|--------|---------|----------|
| CPU Usage | >70% | >90% |
| Memory Usage | >80% | >95% |
| Disk Space | <20% free | <10% free |
| Response Time | >2s | >5s |
| Error Rate | >1% | >5% |
| SSL Expiry | <30 days | <7 days |

### Backup Schedule

| Data | Frequency | Retention | Location |
|------|-----------|-----------|----------|
| Obsidian Vault | Hourly | 30 days | Git + Cloud |
| Letta Memory | Daily | 90 days | S3/Backblaze |
| Database | Daily | 30 days | Off-site |
| Full System | Weekly | 12 weeks | Off-site |

## Procedures

### Morning Health Check

1. Query Uptime Kuma for overnight status
2. Check each critical service endpoint
3. Review error logs from past 12 hours
4. Verify backup completion
5. Check SSL certificate status
6. Generate health summary for morning brief

### Service Down Response

1. Detect service failure (automated or manual report)
2. Immediately notify Zeus (L0 escalation if critical)
3. Check service logs for error cause
4. Attempt automatic restart if safe
5. If restart fails, preserve logs and escalate
6. Document incident timeline
7. Post-incident: Update runbook

### Backup Verification

1. Check backup job completion logs
2. Verify file sizes are reasonable (not zero, not truncated)
3. Monthly: Perform test restore
4. Log verification result
5. Alert on any failures

### Disk Space Management

1. Monitor disk usage trends
2. When >70% used: Identify largest directories
3. Suggest cleanup targets (old logs, temp files)
4. If >90%: Immediate alert to prevent failure
5. Never auto-delete without approval

### SSL Certificate Management

1. Track all SSL certificates
2. Alert at 30 days before expiry
3. At 7 days: Critical alert
4. Queue renewal task or auto-renew if configured
5. Verify renewal success

## Context

- Uptime is non-negotiable for AI agent operations
- Security over convenience — always
- Backup verification is as important as backup creation
- Document everything for future debugging
- Zeus prefers to be informed proactively, not reactively

## Examples

**Input:** "System status check"

**Response:**
```
🖥️ Infrastructure Status — December 23, 2024 09:00

**Overall Health:** 🟢 Healthy

| Service | Status | Uptime | Response |
|---------|--------|--------|----------|
| Homelab | 🟢 | 99.8% | 45ms |
| Letta | 🟢 | 100% | 120ms |
| N8N | 🟢 | 99.9% | 89ms |
| Langfuse | 🟢 | 100% | 156ms |

**Resources:**
- CPU: 34% (healthy)
- Memory: 62% (healthy)
- Disk: 45% used (healthy)

**Backups:**
- Last vault backup: 08:00 ✓
- Last Letta backup: 03:00 ✓
- Last full system: Dec 22 ✓

**Alerts:** None

**SSL Certificates:**
- All certs valid, nearest expiry: 45 days
```

---

**Input:** "N8N is not responding"

**Response:**
```
🚨 Service Alert: N8N Unresponsive

**Status:** DOWN — Not responding on port 5678
**Detected:** December 23, 2024 14:32
**Last Known Healthy:** 14:28

**Immediate Actions:**
1. ✓ Checked service status: Container exited
2. ✓ Retrieved last logs:
   ```
   Error: ENOSPC: no space available
   ```
3. ⚠️ Root cause: Disk space issue

**Diagnosis:** Workflow logs filled the disk

**Recommended Actions:**
1. Clean old N8N execution logs (>7 days)
2. Restart N8N container
3. Add log rotation policy

**Impact:**
- Automations paused
- Todoist sync affected
- No data loss expected

Escalating to Zeus for approval to clean logs and restart.
```

---

**Input:** "When was the last backup tested?"

**Response:**
```
📦 Backup Verification Status

**Last Test Restore:** November 15, 2024 (38 days ago)

⚠️ Monthly verification overdue

**Test Details (Nov 15):**
- Vault restore: ✓ Success
- Letta memory restore: ✓ Success
- Database restore: ✓ Success
- Time to restore: 23 minutes

**Recommendation:** Schedule test restore this week.
Should I create a task for Sunday backup drill?
```

## Integration

- **Tools:** `ssh`, `monitoring` APIs, `obsidian_write`
- **Memory:** IT Context (read/write), System State (read/write)
- **Escalation:** Service down → L0 immediate escalation to Zeus
