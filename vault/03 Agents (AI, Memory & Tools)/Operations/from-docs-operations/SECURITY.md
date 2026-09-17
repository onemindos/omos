# 1. Security Practices

**Protect your OneMind OS system from threats**

---

## Overview

Security is the **foundation** of production infrastructure. This practice covers:

- ✅ Automatic security updates
- ✅ Firewall configuration
- ✅ SSH hardening
- ✅ Brute-force protection
- ✅ Docker security

**Time to implement:** 30 minutes
**Maintenance:** Zero (automated)

---

## Quick Setup

```bash
# Run automated setup
cd /opt/onemind
./production/1-security/setup-security.sh

# Follow prompts to configure:
# - Unattended upgrades
# - UFW firewall
# - SSH hardening
# - Fail2Ban
```

---

## Manual Setup

### 1. Automatic Security Updates

**Install unattended-upgrades:**

```bash
sudo apt update
sudo apt install -y unattended-upgrades apt-listchanges

# Configure interactively
sudo dpkg-reconfigure -plow unattended-upgrades
```

**Or configure manually:**

```bash
# Edit: /etc/apt/apt.conf.d/50unattended-upgrades
sudo nano /etc/apt/apt.conf.d/50unattended-upgrades
```

Add:
```
Unattended-Upgrade::Allowed-Origins {
    "${distro_id}:${distro_codename}-security";
    "${distro_id}ESMApps:${distro_codename}-apps-security";
    "${distro_id}ESM:${distro_codename}-infra-security";
};

Unattended-Upgrade::Automatic-Reboot "true";
Unattended-Upgrade::Automatic-Reboot-Time "03:00";
Unattended-Upgrade::Remove-Unused-Dependencies "true";
Unattended-Upgrade::Remove-Unused-Kernel-Packages "true";
```

**Enable automatic updates:**

```bash
# Edit: /etc/apt/apt.conf.d/20auto-upgrades
sudo nano /etc/apt/apt.conf.d/20auto-upgrades
```

Add:
```
APT::Periodic::Update-Package-Lists "1";
APT::Periodic::Download-Upgradeable-Packages "1";
APT::Periodic::AutocleanInterval "7";
APT::Periodic::Unattended-Upgrade "1";
```

**Verify:**
```bash
sudo systemctl status unattended-upgrades
```

---

### 2. Firewall Configuration (UFW)

**Install UFW:**

```bash
sudo apt install -y ufw
```

**Configure firewall rules:**

```bash
# Default policies
sudo ufw default deny incoming
sudo ufw default allow outgoing

# CRITICAL: Allow SSH first!
sudo ufw allow 22/tcp comment 'SSH'

# OneMind OS ports
sudo ufw allow 8080/tcp comment 'Agno API'
sudo ufw allow 8000/tcp comment 'API Gateway'
sudo ufw allow 4222/tcp comment 'NATS Client'
sudo ufw allow 8222/tcp comment 'NATS Monitoring'
sudo ufw allow 5678/tcp comment 'N8N'
sudo ufw allow 3001/tcp comment 'Uptime Kuma'
sudo ufw allow 9000/tcp comment 'Portainer'
sudo ufw allow 80/tcp comment 'HTTP'
sudo ufw allow 443/tcp comment 'HTTPS'

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status verbose
```

**Port reference:**

| Port | Service | Required |
|------|---------|----------|
| 22 | SSH | ✅ Critical |
| 80/443 | HTTP/HTTPS | ✅ If using domain |
| 8080 | Agno API | ✅ Critical |
| 8000 | API Gateway | ✅ Critical |
| 4222 | NATS Client | ✅ Critical |
| 8222 | NATS Monitoring | ⚪ Optional |
| 5678 | N8N | ⚪ Optional (internal only) |
| 3001 | Uptime Kuma | ⚪ Optional (internal only) |
| 9000 | Portainer | ⚪ Optional (internal only) |

---

### 3. SSH Hardening

**Edit SSH config:**

```bash
sudo nano /etc/ssh/sshd_config
```

**Recommended settings:**

```bash
# Disable root login
PermitRootLogin no

# Use key-based authentication only
PasswordAuthentication no
PubkeyAuthentication yes

# Limit login attempts
MaxAuthTries 3

# Connection timeouts
ClientAliveInterval 300
ClientAliveCountMax 2

# Disable unused features
X11Forwarding no
PermitEmptyPasswords no
```

**Restart SSH:**

```bash
sudo systemctl restart sshd
```

**⚠️ IMPORTANT:** Test SSH in a new terminal before closing your current session!

---

### 4. Fail2Ban (Brute Force Protection)

**Install:**

```bash
sudo apt install -y fail2ban
```

**Configure:**

```bash
# Create local config
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
sudo nano /etc/fail2ban/jail.local
```

**Add/modify:**

```ini
[DEFAULT]
bantime = 3600          # 1 hour ban
findtime = 600          # 10 minute window
maxretry = 5            # 5 failed attempts

[sshd]
enabled = true
port = 22
logpath = /var/log/auth.log
```

**Enable and start:**

```bash
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

**Check status:**

```bash
# Overall status
sudo systemctl status fail2ban

# SSH protection status
sudo fail2ban-client status sshd

# View banned IPs
sudo fail2ban-client get sshd banip
```

---

### 5. Docker Security

**Configure Docker daemon:**

```bash
sudo nano /etc/docker/daemon.json
```

**Add security settings:**

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "userns-remap": "default",
  "no-new-privileges": true
}
```

**Restart Docker:**

```bash
sudo systemctl restart docker
```

**Verify services restart:**

```bash
cd /opt/onemind
docker compose -f containers/cloud/docker-compose.yml ps
```

---

## Verification Checklist

After setup, verify everything works:

### ✅ Auto-Updates
```bash
sudo systemctl status unattended-upgrades
# Should show: active (running)
```

### ✅ Firewall
```bash
sudo ufw status verbose
# Should show: Status: active
# Should list all OneMind ports as ALLOW
```

### ✅ SSH Config
```bash
# Test SSH connection in NEW terminal
ssh ubuntu@YOUR_SERVER_IP

# Should connect without password (key-based only)
```

### ✅ Fail2Ban
```bash
sudo fail2ban-client status sshd
# Should show: Currently banned: 0 (hopefully!)
```

### ✅ Docker
```bash
docker info | grep -i security
# Should show security options enabled
```

---

## Monitoring Security

### Check for Security Updates

```bash
# List available security updates
sudo unattended-upgrade --dry-run --debug

# View update history
grep "Unattended-Upgrade" /var/log/dpkg.log | tail -20
```

### Monitor Failed Login Attempts

```bash
# View recent failed logins
sudo tail -50 /var/log/auth.log | grep "Failed password"

# Check Fail2Ban bans
sudo fail2ban-client status sshd
```

### Review Firewall Logs

```bash
# View denied connections
sudo grep "UFW BLOCK" /var/log/syslog | tail -20
```

---

## Troubleshooting

### Can't SSH After Changes

**If you get locked out:**

1. Use cloud provider console to access server
2. Revert SSH config changes:
   ```bash
   sudo nano /etc/ssh/sshd_config
   # Set PasswordAuthentication yes (temporarily)
   sudo systemctl restart sshd
   ```
3. Test SSH connection
4. Re-apply security settings carefully

### Firewall Blocking Services

**If services can't be accessed:**

```bash
# Check firewall rules
sudo ufw status numbered

# Allow specific port
sudo ufw allow XXXX/tcp

# Or temporarily disable (for debugging)
sudo ufw disable
```

### Fail2Ban Banning Legitimate IPs

**Unban an IP:**

```bash
sudo fail2ban-client set sshd unbanip YOUR_IP
```

**Whitelist an IP permanently:**

```bash
sudo nano /etc/fail2ban/jail.local

# Add under [DEFAULT]:
ignoreip = 127.0.0.1/8 ::1 YOUR_TRUSTED_IP
```

---

## Maintenance

### Weekly Review (5 minutes)

```bash
# Check security updates applied
tail -50 /var/log/dpkg.log | grep security

# Check failed login attempts
sudo fail2ban-client status sshd

# Check firewall blocks
sudo grep "UFW BLOCK" /var/log/syslog | tail -20
```

### Monthly Review (15 minutes)

```bash
# Review SSH log
sudo tail -100 /var/log/auth.log

# Check for unusual activity
last -20

# Review system accounts
cat /etc/passwd | grep -v nologin | grep -v /bin/false

# Update password (if using any)
passwd
```

### Quarterly Tasks

- [ ] Rotate SSH keys
- [ ] Review firewall rules (remove unused)
- [ ] Security audit
- [ ] Update documentation

---

## Security Best Practices

### ✅ DO

- ✅ Keep system updated automatically
- ✅ Use SSH keys (never passwords)
- ✅ Minimize open ports
- ✅ Enable Fail2Ban
- ✅ Monitor logs regularly
- ✅ Use strong passwords (for sudo)
- ✅ Limit sudo access
- ✅ Use Tailscale VPN for internal services

### ❌ DON'T

- ❌ Allow root SSH login
- ❌ Use weak passwords
- ❌ Expose internal services publicly
- ❌ Disable firewall "temporarily"
- ❌ Ignore security updates
- ❌ Run containers as root
- ❌ Store secrets in git

---

## Next Steps

After securing your system:

1. **Setup Backups** → [../2-backups/](../2-backups/README.md)
2. **Configure Monitoring** → [../3-monitoring/](../3-monitoring/README.md)
3. **Enable Automation** → [../5-automation/](../5-automation/README.md)

---

## Resources

- [Ubuntu Security Guide](https://ubuntu.com/security)
- [UFW Documentation](https://help.ubuntu.com/community/UFW)
- [Fail2Ban Documentation](https://www.fail2ban.org/)
- [Docker Security](https://docs.docker.com/engine/security/)

---

**🔒 Your system is now hardened and secure!**
