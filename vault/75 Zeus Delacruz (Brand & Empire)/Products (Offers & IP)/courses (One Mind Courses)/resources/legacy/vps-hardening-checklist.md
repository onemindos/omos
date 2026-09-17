# VPS Hardening Checklist

> *Run this checklist when setting up a new VPS or auditing an existing one.*
> *Each item maps to a specific command. Check and verify, don't just trust memory.*

---

```yaml
---
type: checklist
domain: 04
name: VPS Hardening Checklist
agent: Guardian
status: active
updated: YYYY-MM-DD
---
```

---

## 1. Initial Server Access

- [ ] SSH in as root to confirm access
- [ ] Check Ubuntu/Debian version: `lsb_release -a`
- [ ] Update all packages: `apt update && apt upgrade -y`
- [ ] Set hostname: `hostnamectl set-hostname [your-hostname]`
- [ ] Set timezone: `timedatectl set-timezone [your/timezone]`

---

## 2. Create Non-Root User

- [ ] Create sudo user: `adduser zeus`
- [ ] Add to sudo group: `usermod -aG sudo zeus`
- [ ] Test sudo access: `su - zeus && sudo whoami` (should return `root`)
- [ ] Copy SSH keys: `rsync --archive --chown=zeus:zeus ~/.ssh /home/zeus/`
- [ ] Test SSH login as new user: `ssh zeus@[server-ip]`

---

## 3. SSH Hardening

Edit `/etc/ssh/sshd_config`:

- [ ] `PermitRootLogin no`
- [ ] `PasswordAuthentication no`
- [ ] `PubkeyAuthentication yes`
- [ ] `AuthorizedKeysFile .ssh/authorized_keys`
- [ ] `X11Forwarding no`
- [ ] `AllowTcpForwarding no` (unless you need port forwarding)
- [ ] Restart SSH: `systemctl restart sshd`
- [ ] Verify SSH still works before closing current session

---

## 4. Firewall (UFW)

- [ ] Install: `apt install ufw -y`
- [ ] Default deny inbound: `ufw default deny incoming`
- [ ] Default allow outbound: `ufw default allow outgoing`
- [ ] Allow SSH: `ufw allow ssh`
- [ ] Allow only needed ports: `ufw allow [port]/tcp`
- [ ] Enable: `ufw enable`
- [ ] Verify: `ufw status verbose`

**Ports to allow by service:**
| Service | Port | Protocol |
|---------|------|----------|
| SSH | 22 | TCP |
| OpenClaw gateway | 18789 | TCP (Tailscale only) |
| Tailscale | (managed) | — |

---

## 5. Fail2ban

- [ ] Install: `apt install fail2ban -y`
- [ ] Copy default config: `cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local`
- [ ] Edit `/etc/fail2ban/jail.local` — set `bantime = 1d`, `maxretry = 5`
- [ ] Enable SSH jail: confirm `[sshd]` section has `enabled = true`
- [ ] Start and enable: `systemctl enable --now fail2ban`
- [ ] Verify: `fail2ban-client status sshd`

---

## 6. Automatic Security Updates

- [ ] Install: `apt install unattended-upgrades -y`
- [ ] Enable: `dpkg-reconfigure --priority=low unattended-upgrades`
- [ ] Verify config: `cat /etc/apt/apt.conf.d/20auto-upgrades`

---

## 7. Tailscale (Private Network)

- [ ] Install: `curl -fsSL https://tailscale.com/install.sh | sh`
- [ ] Connect: `tailscale up --ssh`
- [ ] Verify IP: `tailscale ip -4`
- [ ] In Tailscale admin: enable SSH access via Tailscale
- [ ] Lock down UFW: remove direct SSH (22) from public, allow only Tailscale subnet

---

## 8. Secrets Management

- [ ] Install 1Password CLI: `op --version`
- [ ] Authenticate: `op signin`
- [ ] Test secret retrieval: `op read "op://vault-id/test-item/field"`
- [ ] No hardcoded secrets in any config file — verify: `grep -r "password\|secret\|token\|key" /etc/` (investigate findings)
- [ ] Environment variables only injected at runtime, not persisted in `.bashrc`

---

## 9. Monitoring

- [ ] Install basic monitoring: `apt install htop ncdu -y`
- [ ] Check listening ports: `ss -tlnp`
- [ ] Check running services: `systemctl list-units --type=service --state=running`
- [ ] Check failed services: `systemctl --failed`
- [ ] Set up DigitalOcean/Hetzner monitoring alerts (CPU, disk, memory) via provider dashboard

---

## 10. Backups

- [ ] DigitalOcean/provider backups enabled from dashboard
- [ ] Automated local snapshot configured in `_system/backup/`
- [ ] Test restore: confirm at least one backup is restorable
- [ ] Vault backup: `onemind-watch.service` active — vault auto-commits to GitHub
- [ ] Verify last auto-commit: `cd /root/onemind-codex && git log --oneline -5`

---

## 11. Final Verification

- [ ] SSH in as `zeus` user (not root) — works
- [ ] Root login blocked: `ssh root@[ip]` — fails with "Permission denied"
- [ ] Password auth blocked: confirm no password prompt on SSH
- [ ] UFW enabled with correct rules: `ufw status verbose`
- [ ] Fail2ban running: `systemctl status fail2ban`
- [ ] All critical services running: check `systemctl status openclaw-gateway onemind-watch onemind-pull`
- [ ] No unexpected listening ports: `ss -tlnp`

---

## Hardening Score

Run after completing:

| Category | Checks | Complete | Score |
|----------|--------|----------|-------|
| Initial setup | 5 | ___ | ___% |
| Non-root user | 4 | ___ | ___% |
| SSH hardening | 7 | ___ | ___% |
| Firewall | 6 | ___ | ___% |
| Fail2ban | 5 | ___ | ___% |
| Auto updates | 3 | ___ | ___% |
| Tailscale | 5 | ___ | ___% |
| Secrets | 5 | ___ | ___% |
| Monitoring | 5 | ___ | ___% |
| Backups | 5 | ___ | ___% |
| Final verification | 7 | ___ | ___% |

100% completion = hardened server.

---

*Agent Guardian: Run this checklist and report results to Zeus. Flag any incomplete items with severity.*
