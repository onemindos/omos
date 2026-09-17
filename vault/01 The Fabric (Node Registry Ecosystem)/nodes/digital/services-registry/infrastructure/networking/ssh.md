---
title: "CLI — ssh"
type: reference
status: active
created: 2026-04-10
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/ssh, topic/infra]
ai_summary: "SSH reference for onemind-main VPS — host config, key management, tunnels, SCP, and common server admin patterns."
---

# ssh — Secure Shell

> Primary access method for `onemind-main`. Prefer Tailscale IP (`100.77.52.93`) over public IP for security.

---

## Quick Connect

```bash
# Via Tailscale (preferred)
ssh root@100.77.52.93

# Via public IP (fallback)
ssh root@104.248.4.212

# Short alias (after ~/.ssh/config setup)
ssh onemind
```

---

## ~/.ssh/config

Add these to `~/.ssh/config` on Mac:

```
Host onemind
  HostName 100.77.52.93
  User root
  IdentityFile ~/.ssh/id_ed25519
  ServerAliveInterval 60
  ServerAliveCountMax 3

Host onemind-pub
  HostName 104.248.4.212
  User root
  IdentityFile ~/.ssh/id_ed25519
```

---

## Key Management

```bash
# Generate new key
ssh-keygen -t ed25519 -C "zeus@onemindos.com"

# Copy public key to server
ssh-copy-id -i ~/.ssh/id_ed25519.pub root@104.248.4.212

# List authorized keys on VPS
ssh root@100.77.52.93 'cat ~/.ssh/authorized_keys'

# Add key manually (VPS)
echo "ssh-ed25519 AAAA... zeus@mac" >> ~/.ssh/authorized_keys
```

---

## Port Forwarding & Tunnels

```bash
# Local port forward — access VPS service locally
# Access onemind-web (port 3456 on VPS) at localhost:3456
ssh -L 3456:localhost:3456 root@100.77.52.93 -N

# Access OmniDB postgres
ssh -L 5432:localhost:5432 root@100.77.52.93 -N

# Access NATS
ssh -L 4222:localhost:4222 root@100.77.52.93 -N

# Keep in background
ssh -L 3456:localhost:3456 root@100.77.52.93 -fN

# Close background tunnel
kill $(pgrep -f "ssh -L 3456")
```

---

## File Transfer (SCP / rsync)

```bash
# Copy file to VPS
scp ./local-file.txt root@100.77.52.93:/opt/onemind/

# Copy file from VPS
scp root@100.77.52.93:/var/log/app.log ./app.log

# Sync directory (rsync — faster, handles diffs)
rsync -avz ./local-dir/ root@100.77.52.93:/opt/onemind/remote-dir/

# Sync with delete (mirror)
rsync -avz --delete ./local-dir/ root@100.77.52.93:/opt/onemind/remote-dir/
```

---

## Remote Commands

```bash
# Run single command
ssh root@100.77.52.93 'docker ps --format "table {{.Names}}\t{{.Status}}"'

# Run multiple commands
ssh root@100.77.52.93 '
  systemctl status openclaw --no-pager -l
  systemctl status n8n --no-pager -l
'

# Run script from Mac on VPS
cat ./setup.sh | ssh root@100.77.52.93 'bash'
```

---

## Troubleshoot

```bash
# Verbose mode (debug connection issues)
ssh -vvv root@100.77.52.93

# Check VPS SSH service
ssh root@100.77.52.93 'systemctl status sshd'

# View SSH logs on VPS
ssh root@100.77.52.93 'journalctl -u sshd -n 50'
```
