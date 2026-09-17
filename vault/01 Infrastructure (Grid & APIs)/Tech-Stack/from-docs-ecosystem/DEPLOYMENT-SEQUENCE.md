# OneMind OS — Deployment Sequence

> Step-by-step deployment order for all 42 ecosystem additions.
> Baseline OpenClaw first, then layer by tier.

**Philosophy**: Every tier builds on the previous one. Never skip ahead — dependencies compound.

---

## Prerequisites (Before Any Tier)

| Requirement | Status | Verify |
| --- | --- | --- |
| DO Droplet running | ✅ | `ssh openclaw@<ip>` |
| OpenClaw gateway running | ✅ | `curl http://localhost:18789/health` |
| Docker + Compose installed | ✅ | `docker compose version` |
| NATS JetStream running | ✅ | `docker compose up -d nats` |
| OmniDB (PG16+pgvector) running | ✅ | `docker compose up -d omnidb` |
| Tailscale mesh connected | ✅ | `tailscale status` |
| GitHub deploy key added | ⚠️ | `ssh -T git@github.com` |
| .env file configured | ✅ | `cat .env \| wc -l` |

---

## Tier 0 — Immediate (< 1 hour total)

Deploy in this exact order. Each builds on the previous.

### Step 0.1: Workspace Backup (#5)

**Why first**: If anything goes wrong in later steps, we can recover.

```bash
# Initialize git in workspace
cd ~/.openclaw/workspace
git init
git remote add origin git@github.com:Zeus-Delacruz/openclaw-workspace-backup.git
git add -A
git commit -m "initial backup"
git push -u origin main

# Add cron
crontab -e
# Add: */15 * * * * cd ~/.openclaw/workspace && git add -A && git commit -m "auto-backup $(date +\%Y-\%m-\%d_\%H:\%M)" --allow-empty && git push origin main 2>/dev/null
```

**Verify**: `git log --oneline -3` shows recent commits.

### Step 0.2: clawsec (#1)

**Why second**: Security before adding anything else.

```bash
openclaw skills install clawsec
openclaw gateway restart
```

**Verify**: Ask Legacy "run a security audit" → should list installed skill permissions.

### Step 0.3: ClawRouter (#2)

```bash
# Install router plugin
# Configure model tiers in openclaw.json
```

**Verify**: Check logs — simple queries should route to GPT-4o-mini, not Opus.

### Step 0.4: MCP Plugin (#4)

```bash
# Install MCP plugin for universal MCP server access
```

**Verify**: Legacy can list available MCP servers.

### Step 0.5: Native Sessions (#6)

No install needed — built into OpenClaw. Learn and test:

```bash
# Test sessions_spawn
# Test sessions_send
# Document results in workspace/SESSIONS-NOTES.md
```

**Verify**: Legacy can spawn a child agent and receive results.

### Step 0.6: Agent Protocol Skill (#3)

```bash
openclaw skills install agent-protocol
```

**Verify**: Legacy can chain two agents (e.g., research → summary).

**Tier 0 Complete Checklist**:
- [ ] Workspace backed up to private GitHub (auto-sync every 15 min)
- [ ] clawsec installed and scanning
- [ ] ClawRouter routing to cheap models for simple queries
- [ ] MCP plugin unlocking external MCP servers
- [ ] sessions_spawn/sessions_send tested and documented
- [ ] Agent Protocol enabling event-based chains

---

## Tier 1 — Same Day (4-8 hours)

> **Decisions applied**: supermemory (#7) DEFERRED (memory stays basic), Composio (#9) DEFERRED (use individual skills/MCP), openclaw-foundry (#8) DEFERRED (install after 2 weeks of manual patterns). Subtrees (#13, #14) already done.

### Step 1.1: Lobster Workflows (#10) — ✅ WRITTEN

**Why first**: These are the core automation that makes Legacy proactive.

Three `.lob` files already exist in `workspace/lobster/`:

```bash
# Files already in repo — copy to workspace on deploy
cp -r workspace/lobster/ ~/.openclaw/workspace/lobster/
openclaw gateway restart
```

- `morning-brief.lob` — daily 7AM briefing (calendar → weather → system → compose → WhatsApp)
- `health-check.lob` — 5-min system health monitor with alert/recovery
- `git-sync.lob` — 15-min Codex backup with change detection + summarization

**Verify**: Morning brief arrives via WhatsApp at 7 AM next day.

### Step 1.2: Starter Kit Patterns (#15)

No install — read and extract:

```bash
# Clone for reference
git clone https://github.com/zach-highley/openclaw-starter-kit /tmp/starter-ref
# Extract: 355-item security checklist
# Extract: spend cap pattern ($50/day)
# Extract: multi-model failover pattern
# Apply to openclaw.json
```

**Verify**: `openclaw.json` has `max_spend_per_day` set.

### Step 1.3: openclaw-ansible (#11) — REFERENCE ONLY

Do NOT run ansible playbooks — our `.do/` scripts already cover Droplet hardening. Study for patterns:

```bash
git clone https://github.com/openclaw/openclaw-ansible.git /tmp/oc-ansible
# Read: systemd service patterns, backup strategies, log rotation
# Compare against our .do/install-openclaw.sh
```

### Step 1.4: hardened-ansible (#12) — EVALUATE

```bash
# Evaluate for rootless Podman and LiteLLM proxy patterns
# Do NOT run yet — just study the approach
git clone https://github.com/openclaw/hardened-ansible.git /tmp/hardened-ref
```

### Step 1.5: elevenlabs-voices (#16) — DEFERRED

Install when voice channels are actually configured and tested:

```bash
# Not yet — run this after voice channel testing
openclaw skills install elevenlabs-voices
```

**Tier 1 Complete Checklist**:
- [x] awesome-openclaw subtree synced (#13)
- [x] awesome-openclaw-skills subtree synced (#14)
- [ ] 3 Lobster workflows running (morning, git-sync, health)
- [ ] Spend caps and failover configured from starter kit patterns
- [ ] Ansible playbooks reviewed for additional hardening ideas
- [ ] hardened-ansible evaluated for rootless Podman approach

---

## Tier 2 — Docker Compose Additions (1-3 days)

> **Decisions applied**: n8n (#17) SHIPS with profile:full. LiveKit (#18) uses CLOUD (self-hosted is fallback only). Mem0 (#20) DEFERRED (memory stays basic).

### Step 2.1: n8n (#17) — ✅ SHIP

n8n handles complex multi-API orchestration alongside Lobster.

```bash
# Already in docker-compose.yml behind profile:full
docker compose --profile full up -d n8n

# Access at http://localhost:5678 (Tailscale only)
# Create first workflow: OpenClaw webhook → email notification
```

**Division of labor**: Lobster for OpenClaw-native workflows. n8n for cross-platform chains (Todoist → Zoho → Telegram → Discord).

**Verify**: Trigger webhook from Legacy → n8n workflow executes → email arrives.

### Step 2.2: LiveKit Cloud (#18) — ✅ SHIP (cloud)

Primary: LiveKit Cloud (free tier: 50 participant-hours/mo). Self-hosted in docker-compose is fallback only.

```bash
# 1. Sign up at https://cloud.livekit.io
# 2. Get API key + secret
# 3. Add to .env:
#    LIVEKIT_CLOUD_URL=wss://your-project.livekit.cloud
#    LIVEKIT_API_KEY=your-key
#    LIVEKIT_API_SECRET=your-secret

# Test room creation via SDK or curl against cloud endpoint
```

**When to fall back to self-hosted**: >100 hrs/mo usage or hard privacy requirement.

**Verify**: WebRTC connection test passes against cloud endpoint.

### Step 2.3: openclaw-studio (#19) — 🔬 EVALUATE

```bash
# Evaluate — useful for debugging but not critical for launch
# Deploy studio (Docker or standalone)
# Point to OpenClaw gateway
```

**Verify**: Dashboard shows Legacy agent, recent tool calls, cron jobs.

### Step 2.4: acpx (#21) — ⏸ DEFERRED

Install when multi-agent dispatch across Claude/Gemini/Copilot is needed.

### Step 2.5: openclaw-mcp bridge (#22) — ⏸ DEFERRED

Deploy when Claude Desktop → Legacy bridge is useful.

**Tier 2 Complete Checklist**:
- [ ] n8n running at :5678 with first webhook workflow
- [ ] LiveKit Cloud configured with API credentials
- [ ] WebRTC test call successful
- [ ] openclaw-studio evaluated for debugging value

---

## Tier 3 — Hardware & Voice (2-4 weeks)

### Phase 3A: Voice (Week 1)

1. **Jupiter Voice (#30)** — local failover STT/TTS
2. **openclaw-voice (#31)** — browser-based voice chat
3. **voice-avatar (#32)** — lip-synced face (requires LiveKit from Tier 2)

### Phase 3B: Smart Home (Week 2)

1. **HA add-on (#27)** OR **HA conversation (#28)** — pick one based on HA setup
2. **HomeClaw (#29)** — Apple HomeKit devices

### Phase 3C: IoT Nodes (Week 3)

1. **WireClaw (#24)** — first ESP32 in office
2. **HeyClawy (#25)** — voice in kitchen
3. **MimiClaw (#26)** — evaluate for edge reasoning

### Phase 3D: Wearables (Week 4)

1. **VisionClaw (#23)** — Meta Ray-Ban integration (requires LiveKit)

**Tier 3 Complete Checklist**:
- [ ] Local voice failover working
- [ ] Browser-based voice chat accessible
- [ ] Voice avatar with lip sync on display
- [ ] Home Assistant integrated (add-on or conversation)
- [ ] HomeKit devices controllable
- [ ] First WireClaw sensor publishing to NATS
- [ ] HeyClawy voice in at least one room
- [ ] VisionClaw streaming from Meta glasses

---

## Tier 4 — Strategic (1-3 months)

Deploy based on business priorities. No strict ordering.

| Item | When | Trigger |
| --- | --- | --- |
| #33 Virtual Protocol ACP | When GE pillar launches | Revenue readiness |
| #34 Network-AI | When running 3+ frameworks and dispatch is messy | 3+ framework wall |
| #35 Langclaw | When Lobster workflows need branching logic | Complexity limit hit |
| #36 ExoClaw | When running untrusted community skills | Security requirement |
| #37 BankrBot | When crypto/DeFi integration needed | GE pillar decision |
| #38 ClawWork | Before pricing GE services | Revenue research |
| #39 CrewClaw | When managing 5+ agent fleet | Fleet scale |
| #40 OM1 sub-repos | When physical robot integration begins | Hardware available |
| #41 Clawdstrike | When running multiple agent instances | Fleet security |
| #42 Community | Immediately (low effort, continuous) | Always |

---

## Rollback Procedures

### Skill rollback
```bash
openclaw skills uninstall <skill-name>
openclaw gateway restart
```

### Docker service rollback
```bash
docker compose stop <service>
docker compose rm <service>
# Remove from docker-compose.yml
```

### Full workspace rollback
```bash
cd ~/.openclaw/workspace
git log --oneline -10
git checkout <commit-hash> -- .
openclaw gateway restart
```

### Nuclear option
```bash
# Restore from backup
cd ~/.openclaw/workspace
git fetch origin
git reset --hard origin/main
openclaw gateway restart
```

---

## Dependency Graph

```text
Nothing depends on ──→ #1 clawsec (do first always)
                   ──→ #5 workspace backup (do first always)
                   ──→ #42 community (join anytime)

LiveKit Cloud (#18) ←── required by ──→ VisionClaw (#23)
               ←── required by ──→ voice-avatar (#32)

OmniDB         ←── used by ──→ nats-ingester (metrics pipeline)

NATS           ←── required by ──→ WireClaw (#24)
               ←── required by ──→ nexus-bridge (existing)

Home Assistant ←── required by ──→ HA add-on (#27)
               ←── required by ──→ HA conversation (#28)
               ←── required by ──→ HomeClaw (#29)

Lobster        ←── required by ──→ Lobster workflows (#10)
               ←── enhanced by ──→ foundry (#8)

n8n (#17)      ←── for ──→ cross-platform API orchestration
               ←── enhanced by ──→ OpenClaw webhook triggers

#27 HA add-on   ←── competes with ──→ #28 HA conversation (pick one)
```
