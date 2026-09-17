# Agent Access Control Matrix

> *Defines what each agent is permitted to read, write, and execute.*
> *Review after any agent configuration change or security incident.*

---

```yaml
---
type: access-control
domain: 04
agent: Guardian
owner: Zeus
status: active
updated: YYYY-MM-DD
---
```

---

## Access Control Principles

1. **Least privilege** — each agent gets only the access it needs for its domain
2. **No write outside domain** — agents write only to their assigned pillar/domain folders
3. **Inbox as neutral zone** — all agents may write to `99 Inbox/` (captures only)
4. **Legacy orchestrates** — Legacy may read and write across all domains
5. **Guardian audits** — Guardian reviews agent writes during security checks

---

## Agent Read/Write Matrix

| Agent | Reads | Writes | Restricted From |
|-------|-------|--------|-----------------|
| Legacy ⚡ | All vault | All vault | Nothing (orchestrator) |
| Spartan 💪 | HP 27-28, Inbox | HP 27-28, Inbox | All other domains |
| Oracle 🔮 | HP 25-26, 29-34, Inbox | HP 25-26, 29-34, Inbox | SO, LE, GE |
| Eden 🌱 | LE (all), Inbox | LE (all), Inbox | SO, HP, GE |
| Haven 🏡 | LE (all), Inbox | LE (all), Inbox | SO, HP, GE |
| Empire 💰 | GE (all), Inbox | GE (all), Inbox | SO, HP, LE |
| Guardian 🛡️ | SO 04, Vault (read audit) | SO 04, Inbox | Write to other domains |
| Grid ⚙️ | SO 01-02, SO 07, Inbox | SO 01-02, SO 07, Inbox | HP, LE, GE |
| Trinity 🔺 | SO 06, all (read research) | SO 06, Inbox | HP, LE, GE (write) |

---

## External Tool Access Matrix

| Agent | Vault (MCP) | GitHub | Todoist | 1Password | Shell/VPS |
|-------|-------------|--------|---------|-----------|-----------|
| Legacy ⚡ | Read + Write | Read | Read + Write | Read (via op CLI) | Limited |
| Spartan 💪 | Domain only | No | Domain tasks | No | No |
| Oracle 🔮 | Domain only | No | Domain tasks | No | No |
| Eden 🌱 | Domain only | No | Domain tasks | No | No |
| Haven 🏡 | Domain only | No | Domain tasks | No | No |
| Empire 💰 | Domain only | Read | Domain tasks | No | No |
| Guardian 🛡️ | Read-only audit | Read | No | Read (security audit) | Read-only |
| Grid ⚙️ | Domain only | Read | No | Read (infra) | Full (infra work) |
| Trinity 🔺 | Read-only + domain write | Read | No | No | No |

---

## Autonomy Levels (L0–L3)

Each agent operates at a defined autonomy level for different action types:

| Action Type | L0 (Ask always) | L1 (Ask for significant) | L2 (Do, then report) | L3 (Fully autonomous) |
|-------------|-----------------|--------------------------|----------------------|----------------------|
| Read vault | — | — | — | All agents |
| Write to Inbox | — | — | All agents | — |
| Write to domain | — | Spartan/Oracle/Eden/Haven/Empire | Legacy | — |
| Delete vault content | All agents | — | — | — |
| Execute shell commands | Grid (destructive) | Grid (write ops) | — | Grid (read-only) |
| Create Todoist tasks | — | Legacy | Domain agents | — |
| Complete Todoist tasks | — | — | Domain agents | — |
| Push to GitHub | — | Legacy | — | — |
| Rotate credentials | All agents | — | — | — |
| Restart services | Grid (confirm first) | — | — | — |

**L0:** Must ask Zeus before doing
**L1:** Ask if impact is significant; proceed if routine
**L2:** Do it, then tell Zeus what was done
**L3:** Fully autonomous, no report needed

---

## Prohibited Actions (All Agents)

These actions are never permitted without explicit Zeus approval:

- Delete any vault file or folder
- Delete any Todoist task (move to "To Delete" project instead)
- Push force to any GitHub repo
- Rotate or delete credentials in 1Password
- Shut down or reboot the VPS
- Add or remove MCP server configurations
- Change agent configurations (openclaw.json)
- Send messages impersonating Zeus on any channel
- Access files outside the vault unless explicitly told to

---

## Audit Log

Guardian reviews agent writes on a weekly basis and logs:

```
YYYY-MM-DD — Guardian audit
- Legacy: [files written] — [assessment]
- [Agent]: [files written] — [assessment]
- Flags: [any access outside expected domains]
```

Audit log location: `SO/04 Security/agent-audits/YYYY-WW-audit.md`

---

## Violation Response

If an agent writes outside its permitted domain:

1. Guardian flags in Discord: `🚨 ACCESS ALERT: [Agent] wrote to [path] — outside permitted domain`
2. Zeus reviews the write — is it intentional delegation or error?
3. If error: correct the file location, update agent SOUL.md with boundary reminder
4. If pattern: review openclaw.json permissions and SOUL.md restrictions

---

*Agent Guardian: Run a domain write audit weekly. Report any out-of-domain writes immediately.*
