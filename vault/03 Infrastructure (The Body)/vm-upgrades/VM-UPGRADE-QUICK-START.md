---
id: ui-12-20260111053926
title: "VM Upgrade - Quick Reference Card"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/12]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "VM Upgrade - Quick Reference Card"
---

# VM Upgrade - Quick Reference Card

## 🚀 ONE-COMMAND DEPLOYMENT

```bash
cd /Users/zeuslegacy/Documents/Legacy-AI
./scripts/deploy-vm-upgrade.sh
```

**That's it!** The script handles everything automatically.

---

## 📋 What Happens Automatically

1. ✅ Connects to VM via SSH
2. ✅ Creates full system backup
3. ✅ Updates OS packages
4. ✅ Upgrades Docker & Portainer
5. ✅ Updates all 20 containers
6. ✅ Fixes health checks
7. ✅ Optimizes Caddy SSL
8. ✅ Enables firewall
9. ✅ Cleans up resources
10. ✅ Generates report

**Duration**: 10-15 minutes  
**Downtime**: <2 minutes (during Docker restart)

---

## 🔍 Monitor Progress

The script shows real-time output. You'll see:

```
╔════════════════════════════════════════════════════════════╗
║    LEGACY AI - VM UPGRADE DEPLOYMENT (Remote Execution)   ║
╚════════════════════════════════════════════════════════════╝

═══ PREFLIGHT CHECKS ═══
🔍 Checking SSH connection...
✅ SSH connection successful

═══ UPLOADING UPGRADE SCRIPTS ═══
📤 Uploading upgrade scripts...
✅ Scripts uploaded

═══ EXECUTING UPGRADE ON VM ═══
🚀 Starting upgrade process...

[Progress indicators...]

✅ Main upgrade completed

═══ APPLYING HEALTH CHECK FIXES ═══
✅ Health checks updated

═══ UPDATING CADDY CONFIGURATION ═══
✅ Caddy configuration updated

═══ FINAL VALIDATION ═══
📊 Container status: [20/20 running]
🌐 Testing domains: [All ✅]

╔════════════════════════════════════════════════════════════╗
║              UPGRADE DEPLOYMENT COMPLETE                   ║
╚════════════════════════════════════════════════════════════╝
```

---

## 🛡️ Safety Net

### Automatic Backup
Before any changes, full backup created at:
```
/home/ubuntu/backups/pre-upgrade-YYYYMMDD_HHMMSS/
```

### Rollback (If Needed)
```bash
ssh -i "/path/to/key.pem" ubuntu@3.215.29.128
sudo bash rollback-upgrade.sh ~/backups/pre-upgrade-YYYYMMDD_HHMMSS
```

---

## 📊 Verify Success

After deployment, check:

### 1. Containers Running
```bash
ssh -i "/path/to/key.pem" ubuntu@3.215.29.128 "docker ps"
```
Expected: 20/20 running

### 2. All Healthy
```bash
ssh -i "/path/to/key.pem" ubuntu@3.215.29.128 "docker ps --filter 'health=unhealthy'"
```
Expected: Empty (no unhealthy containers)

### 3. Domains Working
- ✅ https://command.zeusops.dev (Portainer)
- ✅ https://n8n.zeusops.dev (n8n)
- ✅ https://dify.zeusops.dev (Dify)
- ✅ https://letta.zeusops.dev (Letta)

### 4. Resources Optimal
```bash
ssh -i "/path/to/key.pem" ubuntu@3.215.29.128 "free -h && df -h /"
```
Expected: <60% memory, <80% disk

---

## 🔧 Troubleshooting

### Script Won't Run
```bash
# Make executable
chmod +x /Users/zeuslegacy/Documents/Legacy-AI/scripts/deploy-vm-upgrade.sh

# Check SSH key
ls -l /Users/zeuslegacy/Documents/LightsailDefaultKey-us-east-1.pem
# Should be: -rw-------
```

### Upgrade Fails
```bash
# Check logs on VM
ssh -i "/path/to/key.pem" ubuntu@3.215.29.128
tail -100 /var/log/vm-upgrade-*.log
```

### Container Unhealthy
```bash
# Re-run health check fixer
ssh -i "/path/to/key.pem" ubuntu@3.215.29.128
bash ~/upgrade-scripts/fix-health-checks.sh
```

---

## 📁 Files Created

### Local (Repository)
- `scripts/vm-upgrade-complete.sh` - Main upgrade script
- `scripts/fix-health-checks.sh` - Health check fixer
- `scripts/deploy-vm-upgrade.sh` - Deployment orchestrator
- `scripts/rollback-upgrade.sh` - Rollback script
- `ENVIRONMENTS/production/infrastructure/caddy/Caddyfile.optimized` - Optimized Caddy config
- `KNOWLEDGE/docs/operations/VM-UPGRADE-GUIDE.md` - Full documentation

### On VM (After Deployment)
- `~/upgrade-scripts/` - All upgrade scripts
- `~/backups/pre-upgrade-*/` - System backup
- `/var/log/vm-upgrade-*.log` - Upgrade log

---

## ✅ Post-Upgrade Checklist

- [ ] All containers running (20/20)
- [ ] No unhealthy containers
- [ ] All domains accessible
- [ ] SSL certificates valid
- [ ] Firewall enabled
- [ ] Resources optimal
- [ ] Applications functional

---

## 🎯 Next Steps After Upgrade

1. **Test Critical Workflows**
   - n8n automations
   - Dify conversations
   - Letta agent responses

2. **Monitor for 24 Hours**
   - Container logs
   - Resource usage
   - Error patterns

3. **Document Changes**
   - Update VM-STACKS-INVENTORY.md
   - Note any version changes
   - Record any issues/fixes

4. **Plan Next Maintenance**
   - Schedule: Monthly
   - Focus: Security updates
   - Review: Performance metrics

---

## 📞 Support

**Full Documentation**: `KNOWLEDGE/docs/operations/VM-UPGRADE-GUIDE.md`

**Quick Commands**:
```bash
# Check status
ssh ubuntu@3.215.29.128 "docker ps"

# View logs
ssh ubuntu@3.215.29.128 "docker logs <container>"

# Restart container
ssh ubuntu@3.215.29.128 "docker restart <container>"

# View upgrade report
ssh ubuntu@3.215.29.128 "cat ~/backups/pre-upgrade-*/UPGRADE_REPORT.txt"
```

---

**Ready to Upgrade?**

```bash
cd /Users/zeuslegacy/Documents/Legacy-AI
./scripts/deploy-vm-upgrade.sh
```

Type "yes" when prompted and watch the magic happen! ✨


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
