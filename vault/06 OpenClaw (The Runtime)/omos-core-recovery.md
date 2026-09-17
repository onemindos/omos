---
domain: 02
name: omos-core Disaster Recovery
agent: Sovereignty
status: active
updated: 2026-05-14
---

# omos-core — Disaster Recovery

> Full recovery playbook for the primary OpenClaw node.  
> Canonical source: **[OneMind-OS/onemind-recovery](https://github.com/OneMind-OS/onemind-recovery)** (private)

---

## Node Specs

| Field | Value |
|-------|-------|
| Provider | AWS Lightsail |
| Instance | `omos-core` |
| vCPU / RAM | 2 vCPU / 4 GB |
| OS | Ubuntu 24.04 |
| SSH | `ssh onemindos-core` |
| Public IP | `64.227.8.23` |
| Tailscale | `omos-core.tail717ba0.ts.net` / `100.98.25.10` |
| Gateway URL | `https://omos-core.tail717ba0.ts.net` (port 18789) |
| User | `ubuntu` |

---

## Automated Backup System (Live — 2026-05-14)

### What's backed up

| Backup | Schedule | Location |
|--------|----------|----------|
| OpenClaw full archive (`~/.openclaw/`) | Daily 3:00 AM UTC | `~/backups/*.tar.gz` (last 7 kept) |
| Sanitized `openclaw.json` template | Weekly Sunday 4:00 AM UTC | `onemind-recovery` GitHub repo |

### Cron (on omos-core)

```
0 3 * * *   ~/.openclaw/recovery-scripts/backup-cron.sh   → ~/backups/backup.log
0 4 * * 0   ~/.openclaw/recovery-scripts/push-templates.sh → ~/backups/push-templates.log
```

### Manual backup

```bash
ssh onemindos-core 'openclaw backup create --output ~/backups/ --verify'
```

---

## Emergency Triage (30 seconds)

```bash
# 1. Can you SSH?
ssh onemindos-core 'systemctl status openclaw-gateway.service'

# 2. Is gateway responding?
curl -sf https://omos-core.tail717ba0.ts.net/health

# 3. Restart gateway
ssh onemindos-core 'sudo systemctl restart openclaw-gateway.service'

# 4. View logs
ssh onemindos-core 'journalctl -u openclaw-gateway.service -n 50'
```

---

## Restore from Backup

```bash
# List available backups
ssh onemindos-core 'ls -lh ~/backups/*.tar.gz'

# Restore (replaces current ~/.openclaw)
ssh onemindos-core '
BACKUP=$(ls ~/backups/*.tar.gz | tail -1)
echo "Restoring: $BACKUP"
sudo systemctl stop openclaw-gateway.service
mv ~/.openclaw ~/.openclaw.bak.$(date +%s)
tar -xzf "$BACKUP" -C ~/
sudo systemctl start openclaw-gateway.service
'
```

---

## Full Rebuild (Node Dead)

See `scripts/rebuild-omos-core.sh` in [onemind-recovery](https://github.com/OneMind-OS/onemind-recovery).

Steps:
1. Provision new Lightsail Ubuntu 24.04 instance
2. Add SSH key, update `~/.ssh/config` entry `onemindos-core`
3. Run `rebuild-omos-core.sh` — installs Node 22, OpenClaw, deploys secrets via `op inject`, starts gateway
4. Apply voice patch: `~/.openclaw/patches/patch-universal-voice.sh`
5. Verify all 5 agents online at the dashboard

Secrets injection requires `op` CLI + 1Password vault `aaeyris5z44pdqspk277u2oppi`.

---

## onemind-recovery Repo Structure

```
onemind-recovery/
├── README.md                        ← Full runbook (read this first)
├── configs/
│   ├── gateway.env.template         ← All 23 env vars (op:// refs)
│   ├── openclaw.json.template       ← Sanitized config (auto-updated weekly)
│   └── systemd/
│       ├── openclaw-gateway.service
│       └── drop-ins/                ← 6 API key drop-in configs
├── patches/
│   └── patch-universal-voice.sh    ← Voice patch (auto-detects hash)
└── scripts/
    ├── backup-cron.sh              ← Daily backup + rotation
    ├── push-templates.sh           ← Weekly config sync to GitHub
    ├── install-cron.sh             ← One-time cron installer
    └── rebuild-omos-core.sh        ← Full automated rebuild
```

---

## Related

- [VPS-MAIN.md](VPS-MAIN.md) — Legacy node reference (stale — `onemind-main` deleted 2026-04-15)
- `00-24 SO/06 OpenClaw (The Runtime)/` — Agent config and workspace docs
- `00-24 SO/03 Security (The Shield)/` — Hardening notes (2026-04-15)
