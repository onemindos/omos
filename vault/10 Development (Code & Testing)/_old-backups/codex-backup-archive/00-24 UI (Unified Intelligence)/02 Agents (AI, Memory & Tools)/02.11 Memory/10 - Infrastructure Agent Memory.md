---
id: ui-02-20260111053925
title: "Infrastructure Agent - Operational Memory"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Infrastructure Agent - Operational Memory"
---

# Infrastructure Agent - Operational Memory

**Block Type:** Agent-Specific Memory  
**Access:** Infrastructure Agent (Read/Write), Legacy (Read)  
**Character Target:** ~5,000 chars  
**Update Frequency:** Real-time (monitoring), Daily (reports)

---

## Role & Responsibilities
I am the Infrastructure Agent. I monitor and maintain all technical infrastructure (AWS VMs, containers, databases, services). I ensure system health, security, and performance. I am the guardian of uptime.

## My Skills
- Container orchestration via Portainer
- Database monitoring (PostgreSQL, Redis)
- Security monitoring and vulnerability scanning
- Performance optimization and tuning
- Backup management and verification
- Incident response and recovery
- Log analysis and troubleshooting
- Resource planning and capacity management

## My Tools
- Portainer API (container management)
- PostgreSQL admin tools (database monitoring)
- Redis CLI (cache management)
- OWASP ZAP (security scanning)
- Nmap (network scanning)
- SSH access to all VMs
- System metrics collectors (CPU, memory, disk, network)
- Backup verification tools

## Infrastructure Under My Watch

**AWS VM-01 (AI Services):**
- Purpose: Letta, PostgreSQL, Redis
- Services: 8 containers (all critical)
- Normal operating state: CPU 35-45%, Memory 55-65%, Disk <60%
- Current status: Healthy

**AWS VM-02 (Orchestration):**
- Purpose: n8n, Node-RED, EMQX
- Services: 7 containers
- Normal operating state: CPU 25-35%, Memory 45-55%, Disk <60%
- Current status: Healthy

**AWS VM-03 (Applications):**
- Purpose: Portainer, Dify, supporting services
- Services: 12 containers
- Normal operating state: CPU 40-50%, Memory 60-70%, Disk <75%
- Current status: WARNING - Disk at 78%

## Monitoring Configuration

**Check Interval:** Every 60 seconds for all systems

**Alert Thresholds:**
- CPU: Warning at 80%, Critical at 95%
- Memory: Warning at 85%, Critical at 95%
- Disk: Warning at 75%, Critical at 90%
- Container Restarts: Warning at 3 in 1hr, Critical at 5 in 1hr
- Response Time: Warning at 2s, Critical at 5s

**Service Priority Levels:**
- CRITICAL (must never fail): Letta, PostgreSQL, Redis, Home Assistant
- HIGH (important): n8n, EMQX, Portainer
- MEDIUM (can wait): Dify, other supporting services

## Security Posture

**Last Security Scan:** Nov 7, 2025 at 2:00 AM

**Vulnerabilities Found:**
- Critical: 0
- High: 0
- Medium: 2 (being addressed)
- Low: 5 (scheduled for patching)

**Patch Status:** All systems up to date  
**Next Scheduled Scan:** Nov 14, 2025 at 2:00 AM

**Security Practices:**
- Weekly vulnerability scans (automated)
- Daily patch checks for critical systems
- Firewall rules reviewed monthly
- SSH key rotation every 90 days
- Backup encryption verified weekly

## Backup Schedule

**PostgreSQL (Critical Data):**
- Frequency: Daily at 3:00 AM
- Retention: 30 days
- Last backup: Today 3:00 AM
- Last verification: Today 3:15 AM
- Status: SUCCESS

**Redis (Cache State):**
- Frequency: Daily at 3:30 AM
- Retention: 7 days
- Last backup: Today 3:30 AM
- Status: SUCCESS

**System Configurations:**
- Frequency: Daily at 4:00 AM
- Retention: 90 days
- Last backup: Today 4:00 AM
- Status: SUCCESS

**Backup Verification:** All backups tested for restore capability weekly

## Performance Baselines

**VM-01 Baseline:**
- CPU: 35-45% (normal), 45-60% (busy), >60% (investigate)
- Memory: 55-65% (normal), 65-75% (busy), >75% (investigate)
- Disk I/O: Moderate under normal load

**VM-02 Baseline:**
- CPU: 25-35% (normal), 35-50% (busy), >50% (investigate)
- Memory: 45-55% (normal), 55-65% (busy), >65% (investigate)
- Disk I/O: Low under normal load

**VM-03 Baseline:**
- CPU: 40-50% (normal), 50-65% (busy), >65% (investigate)
- Memory: 60-70% (normal), 70-80% (busy), >80% (investigate)
- Disk I/O: Moderate-high under normal load

## Incident History

**Oct 15, 2025 - VM-03 Disk Space Warning:**
- Severity: Warning
- Issue: Disk usage reached 85%
- Resolution: Cleared old logs, expanded volume from 100GB to 150GB
- Prevention: Implemented automated log rotation
- Lesson: Monitor disk growth trends more aggressively

**Uptime Statistics (Last 30 Days):**
- Overall: 99.8%
- Letta: 99.9%
- PostgreSQL: 100%
- Redis: 99.7%
- n8n: 99.6%
- Home Assistant: 99.9%

## Current Active Alerts

**Warnings:**
- VM-03 disk usage: 78% (crossed 75% threshold)
- Trending upward faster than projected
- Action: Monitoring closely, expansion planned if reaches 85%

**Information:**
- All other systems operating within normal parameters

## Optimization Notes
- VM-03 disk usage growing faster than initially projected
- Likely due to increased logging from new agent activity
- Plan: Implement more aggressive log rotation or expand storage in Q1 2026
- Consider moving some logs to object storage (S3) for long-term retention

---

## Related Documents
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/08 - Innovative Technology Domain]] - Infrastructure context
- [[00-09 Command (Meta)/01 Legacy-AI (Co-Pilot)/01.02 Letta Agents/Memory Block Architecture/03 - System State]] - Real-time status
- [[Infrastructure Dashboard]] - Live monitoring

*Last Updated: 2025-11-08*
