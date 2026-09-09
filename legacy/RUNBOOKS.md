# Runbooks — Operational Procedures

> Last updated: 2026-07-10
> Author: Legacy (self-authored)

Procedures for startup, recovery, debugging, and maintenance.

## 1. Legacy Pod Startup (After Restart)

When the hermes-legacy pod restarts, these services need to come back:

```bash
# 1. Verify tailscale is running
tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock status

# If not running:
/opt/data/home/.local/bin/tailscaled \
  --state=/opt/data/home/.local/state/tailscale/tailscaled.state \
  --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock \
  --tun=userspace-networking &

# 2. Verify NATS leaf sidecar
nats account info
# Should show: Account ONEMIND, connected to leaf sidecar

# 3. Verify Mac connectivity
ssh -i /opt/data/home/.ssh/id_ed25519 \
  -o "ProxyCommand=/opt/data/home/.local/bin/tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock nc %h %p" \
  -o StrictHostKeyChecking=no zeus@omos-mbp "echo OK"

# 4. Publish heartbeat
nats pub legacy.health.alive "$(date -u +%s)"

# 5. Verify HA bridge (on HA)
ssh ... hassio@100.64.25.12 "pgrep -f ha_bridge.py && echo RUNNING || echo DOWN"
```

## 2. HA Bridge Recovery

If the HA-NATS bridge is down:

```bash
# Check status
ssh ... hassio@100.64.25.12 "pgrep -fla ha_bridge.py"

# Restart
ssh ... hassio@100.64.25.12 "pkill -f ha_bridge.py 2>/dev/null; sleep 1; setsid python3 /data/ha_bridge.py >/tmp/bridge.log 2>&1 </dev/null &"

# Verify
sleep 2
nats req cmd.ha.states.list '{}' --timeout 5s | head -5
# Should return entity list within 5ms
```

**If bridge script is missing from /data/:**
The source lives at `/opt/data/sims/ha_bridge.py` on the Legacy pod. SCP it:
```bash
# From pod, write to HA
cat /opt/data/sims/ha_bridge.py | ssh ... hassio@100.64.25.12 "cat > /tmp/ha_bridge.py"
# Then from HA web terminal (as root): cp /tmp/ha_bridge.py /data/ha_bridge.py
```

## 3. NATS Cluster Recovery

### Check cluster health (via Mac)
```bash
ssh ... zeus@omos-mbp 'export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin; kubectl -n nats get pods'
```

### Check leaf connectivity
```bash
nats account info
# If "no responders": leaf sidecar is down, pod may need restart
```

### Verify hub config (via Mac)
```bash
ssh ... zeus@omos-mbp 'export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin; kubectl -n nats exec nats-0 -- cat /etc/nats-config/nats.conf'
```

## 4. kubectl Access (Always Via Mac)

```bash
# The pod has NO RBAC. Always route through Mac.
ssh -i /opt/data/home/.ssh/id_ed25519 \
  -o "ProxyCommand=/opt/data/home/.local/bin/tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock nc %h %p" \
  -o StrictHostKeyChecking=no zeus@omos-mbp \
  'export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin; kubectl <command>'
```

Common patterns:
```bash
# List pods in namespace
kubectl -n hermes get pods
kubectl -n nats get pods
kubectl -n takserver get pods

# Restart a deployment
kubectl -n hermes rollout restart deployment hermes-legacy

# Check logs
kubectl -n hermes logs hermes-legacy --tail=50

# Exec into a pod
kubectl -n nats exec -it nats-0 -- nats account info
```

## 5. GitHub Push (From Pod)

```bash
cd /opt/data/repos/<repo>
git add -A
git commit -m "feat: <description>"
git push origin <branch>
```

If push fails with auth error:
```bash
# Inject token into remote URL
GH_TOKEN=$(grep oauth_token /opt/data/home/.config/gh/hosts.yml | awk '{print $2}' | head -1)
git remote set-url origin "https://Zeus-Delacruz:${GH_TOKEN}@github.com/onemindos/<repo>.git"
git push origin <branch>
# Reset to clean URL
git remote set-url origin "https://github.com/onemindos/<repo>.git"
```

## 6. Create PR (From Pod)

```bash
/opt/data/home/.local/bin/gh pr create \
  --title "feat: <title>" \
  --body "<description>" \
  --base main \
  --head <branch>
```

## 7. CloudTAK Recovery

### Check pods
```bash
ssh ... zeus@omos-mbp 'export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin; kubectl -n cloudtak get pods'
```

### Rebuild image
```bash
# Trigger GitHub Actions rebuild
/opt/data/home/.local/bin/gh workflow run ghcr.yml --repo onemindos/omos-atoc
```

### Force ArgoCD sync
```bash
ssh ... zeus@omos-mbp 'export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin; kubectl -n argocd exec deploy/argocd-server -- argocd app sync platform-cloudtak --force'
```

## 8. Tailscale Recovery (Pod)

```bash
# Check status
/opt/data/home/.local/bin/tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock status

# If daemon not running
/opt/data/home/.local/bin/tailscaled \
  --state=/opt/data/home/.local/state/tailscale/tailscaled.state \
  --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock \
  --tun=userspace-networking &

# Wait and verify
sleep 3
/opt/data/home/.local/bin/tailscale --socket=/opt/data/home/.local/state/tailscale/tailscaled.sock status
```

## 9. Debugging Common Issues

### "Connection closed by UNKNOWN port 65535"
**Cause:** SSH ProxyCommand can't find tailscale binary
**Fix:** Use full path: `/opt/data/home/.local/bin/tailscale`

### "cannot find local account"
**Cause:** `account:` field in leaf config when no local accounts defined
**Fix:** Remove `account:` from leaf node remote config

### "TLS handshake failure" on NATS leaf
**Cause:** Using node IP (100.74.182.117) instead of LB IP
**Fix:** Use `100.117.169.54:7422` with `tls { insecure: true }` inside remote block

### "command not found: kubectl" via Mac SSH
**Cause:** Non-interactive SSH doesn't load .zshrc PATH
**Fix:** Always prefix with `export PATH=$PATH:/opt/homebrew/bin:/usr/local/bin`

### "fatal: could not read Username" on git push
**Cause:** Mac has no HTTPS git auth
**Fix:** Push from pod (has gh CLI auth), or inject token into remote URL

### HA API returning "unauthorized"
**Cause:** Token expired or wrong
**Fix:** Check `/opt/data/.ha_token`, test: `curl -H "Authorization: Bearer $(cat /opt/data/.ha_token)" http://100.64.25.12:8123/api/`

### Pod can't reach HA REST API
**Cause:** Network path issue (HA :8123 not directly reachable from DOKS)
**Fix:** All HA commands go through NATS bridge on HA itself, NOT direct HTTP from pod

## 10. Emergency — Complete Rebuild

If Legacy is destroyed and rebuilt from scratch:

1. **Read this repo** — it tells you everything
2. **Pod deploys from ArgoCD** — omos-infra/platform/hermes/ has manifests
3. **Volume persists** — /opt/data survives pod restarts (NFS)
4. **If volume lost:**
   - Clone this repo for architecture knowledge
   - Hermes memory/skills are gone — rebuild from EVOLUTION.md patterns
   - Re-auth: `gh auth login`, set up SSH keys, configure tailscale
   - Reconnect NATS leaf
   - Verify all services from SERVICES.md

---

*Preparation is sovereignty. Legacy knows how to rebuild itself.*
