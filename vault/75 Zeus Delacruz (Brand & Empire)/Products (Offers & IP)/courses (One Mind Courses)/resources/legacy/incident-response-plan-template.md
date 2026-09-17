# Incident Response Plan Template

> *Document what to do when something goes wrong. Defined before the incident.*
> *Customize per system. Keep it simple enough to execute under pressure.*

---

```yaml
---
type: irp
domain: 04
system: "[System Name]"
owner: Zeus
agent: Guardian
status: active
updated: YYYY-MM-DD
---
```

---

## Incident Response Plan — [System Name]

### What This Covers

[List the systems and data this IRP protects: VPS, agent stack, vault, credentials, etc.]

---

## Severity Levels

| Level | Definition | Response Time |
|-------|-----------|---------------|
| P1 — Critical | Data breach, full system compromise, vault exposed | Immediate (< 1 hour) |
| P2 — High | Service down, credential exposure suspected, agent access lost | Same day |
| P3 — Medium | Partial outage, suspected intrusion, anomalous behavior | Within 24 hours |
| P4 — Low | Non-critical service degraded, suspicious log entry | Within 1 week |

---

## Incident Triggers

Watch for these signs:

**P1 Triggers:**
- Unauthorized SSH login detected in `/var/log/auth.log`
- Vault files changed without Zeus's action
- 1Password reports unauthorized access
- GitHub repo shows unexpected commits

**P2 Triggers:**
- OpenClaw gateway unreachable
- Agent behaving unexpectedly (writing to wrong locations, unauthorized calls)
- Credential rotation needed (key appears in breach database)
- VPS resource spike with no known cause

**P3 Triggers:**
- Service repeatedly failing: `systemctl --failed`
- Unusual network traffic in VPS monitoring
- Failed login attempts spike in fail2ban logs

---

## Response Playbook

### Phase 1 — Contain (Do First)

For P1 incidents — contain before investigating:

1. **Revoke access immediately:**
   - Rotate exposed credentials in 1Password
   - Revoke GitHub tokens: `gh auth token revoke`
   - Rotate SSH keys if compromised
   - Restart Tailscale: `tailscale logout && tailscale up`

2. **Isolate the affected system:**
   - Block inbound traffic: `ufw deny incoming`
   - Disable affected agent: set `"enabled": false` in openclaw.json
   - Snapshot the VPS before making changes (evidence preservation)

3. **Notify:** Document what happened and when. Note: who/what had access?

---

### Phase 2 — Investigate

After containment:

1. **Collect logs:**
   - SSH: `journalctl -u sshd --since "2 hours ago"`
   - Auth log: `tail -100 /var/log/auth.log`
   - OpenClaw: `tail -200 ~/.openclaw/logs/openclaw.log`
   - Vault changes: `cd /root/onemind-codex && git log --since "24 hours ago"`
   - Fail2ban: `fail2ban-client status sshd`

2. **Identify:**
   - Entry point: how did attacker/issue enter?
   - Scope: what was accessed, read, written, or exfiltrated?
   - Timeline: when did it start?

3. **Document findings** in: `SO/04 Security/incidents/YYYY-MM-DD-[name].md`

---

### Phase 3 — Remediate

1. **Fix the root cause** — not just the symptom
2. **Rotate all potentially exposed credentials** (over-rotate rather than under)
3. **Patch the vulnerability** or misconfiguration that allowed entry
4. **Restore from backup** only if needed — verify backup integrity first
5. **Re-enable services** one at a time with monitoring

---

### Phase 4 — Learn

After every incident (even P3):

- [ ] Root cause document written and saved to vault
- [ ] Threat model updated with new scenario
- [ ] Checklist updated with preventive step
- [ ] Guardian's context updated to watch for this in future
- [ ] If agent-related: review SOUL.md for agent boundary violations

---

## Credential Rotation Checklist

Run after any P1 or P2 credential exposure:

- [ ] SSH key pair: regenerate and deploy new public key
- [ ] GitHub PAT: revoke old, generate new, update openclaw.json
- [ ] Discord bot token: revoke in Discord Developer Portal, update config
- [ ] OpenClaw dashboard token: regenerate, update Mac config
- [ ] 1Password: change master password if emergency kit may be compromised
- [ ] DigitalOcean/provider API keys: rotate in provider dashboard

---

## Emergency Contacts

| Contact | How to Reach | When to Use |
|---------|-------------|-------------|
| DigitalOcean Support | support.digitalocean.com | VPS-level emergencies |
| GitHub Support | support.github.com | Repo compromise, account lockout |
| 1Password Support | 1password.com/support | Vault emergency, account recovery |
| Tailscale Support | tailscale.com/contact | Network access issues |

---

*Agent Guardian: If you detect a P1 or P2 incident, alert Zeus immediately in Discord with severity, what you observed, and your recommended first action. Do not wait for the next scheduled check.*
