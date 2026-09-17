# PRODUCTION ROADMAP — One Mind OS v0.5 Launch

> **The master plan**: from where we are today to a production-grade agent system + community launch.
> Generated: March 18, 2026 | Source: VPS audit, STATUS.md, MEMORY-ARCHITECTURE.md, ACADEMY-CATALOG.md

---

## CURRENT STATE — Honest Assessment (March 18, 2026)

### What's LIVE and Working

| Component | Status | Details |
|-----------|--------|---------|
| OpenClaw Gateway | **LIVE** | v2026.3.13, 9 agents, 22 models, 5 providers |
| NATS JetStream | **LIVE** | 6 streams, MQTT bridge, Docker container, 4 days uptime |
| OmniDB (PostgreSQL) | **LIVE** | TimescaleDB + pgvector, Docker, healthy |
| legacy-bridge | **LIVE** | NATS → OpenClaw webhook relay |
| nats-ingester | **LIVE** | NATS → OmniDB archival pipeline |
| Discord (9 channels) | **LIVE** | All 9 agents bound to channels |
| iMessage (BlueBubbles) | **LIVE** | Connected via Tailscale to MacBook |
| Lossless Claw plugin | **LIVE** | v0.3.0, DAG-based lossless summarization |
| Memory search | **LIVE** | Hybrid BM25+vector, indexes codex+memory+sessions |
| SwarmClaw dashboard | **LIVE** | v0.7.7 on port 3456 |
| OpenClaw Office | **LIVE** | "The Homestead" on port 5180 |
| Caddy reverse proxy | **LIVE** | 4 custom domains with HTTPS |
| Tailscale mesh | **LIVE** | Mac + VPS + iPhone connected |
| Cron jobs | **LIVE** | morning-brief (7am), health-check (6h), evening-digest (9pm) |
| legacy-pulse | **LIVE** | NATS ambient state → AMBIENT_STATE.md |
| OpenViking Context DB | **LIVE** | v0.2.6, port 1933 |
| NATS dashboard | **LIVE** | port 18300 |

### What's BROKEN (Must Fix Immediately)

| Component | Status | Impact |
|-----------|--------|--------|
| **codex-watch.service** | **NOT FOUND** | Codex changes on VPS are NOT auto-committing to GitHub |
| **codex-pull.service/timer** | **NOT FOUND** | GitHub codex changes are NOT auto-pulling to VPS |
| **openclaw-reload.service** | **DEAD since Mar 16** | Skill file changes NOT triggering hot-reload |

### What's NOT Running (Should Be)

| Component | Status | Priority |
|-----------|--------|----------|
| n8n workflow engine | Not started | P0 — needed for webhook ingestion |
| Telegram channel | Not configured | P1 — needed for omnipresent mode |
| WhatsApp channel | Not configured | P2 |
| legacy-eyes (vision) | Not deployed | P2 |
| ClawVault (structured memory) | Not evaluated | P3 |

---

## PHASE 0: EMERGENCY FIXES (Today — March 18)

Fix the broken services that are preventing the system from self-healing.

### 0.1 — Restore codex-watch.service

The codex auto-commit service is missing on VPS. Legacy writes to the codex (memory logs, vault edits) but they aren't pushing to GitHub.

```bash
# SSH to VPS and check what happened
ssh root@104.248.4.212

# Check if the unit file exists
ls -la /etc/systemd/system/codex-watch.*

# If missing, recreate from the monorepo
cat > /etc/systemd/system/codex-watch.service << 'EOF'
[Unit]
Description=Codex Watch — auto-commit codex + workspaces to GitHub
After=network.target

[Service]
Type=simple
User=onemind
WorkingDirectory=/home/onemind/onemind-os
ExecStart=/home/onemind/onemind-os/openclaw/deploy/scripts/codex-watch.sh
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now codex-watch.service
systemctl status codex-watch.service
```

### 0.2 — Restore codex-pull.service + timer

GitHub → VPS sync is missing. Changes pushed from local Obsidian aren't reaching Legacy.

```bash
cat > /etc/systemd/system/codex-pull.service << 'EOF'
[Unit]
Description=Codex Pull — sync GitHub to VPS monorepo
After=network.target

[Service]
Type=oneshot
User=onemind
WorkingDirectory=/home/onemind/onemind-os
ExecStart=/usr/bin/git pull --rebase --autostash origin main
StandardOutput=journal
StandardError=journal
EOF

cat > /etc/systemd/system/codex-pull.timer << 'EOF'
[Unit]
Description=Codex Pull Timer — every 1 minute

[Timer]
OnCalendar=*:*:00
Persistent=true

[Install]
WantedBy=timers.target
EOF

systemctl daemon-reload
systemctl enable --now codex-pull.timer
systemctl status codex-pull.timer
systemctl start codex-pull.service
journalctl -u codex-pull -n 5
```

### 0.3 — Fix openclaw-reload.service

Skill changes need to trigger a gateway hot-reload. Dead since March 16.

```bash
# Check why it died
journalctl -u openclaw-reload -n 20

# Restart it
systemctl restart openclaw-reload.service
systemctl status openclaw-reload.service

# If it keeps dying, check the script it runs
cat /home/onemind/onemind-os/openclaw/deploy/scripts/openclaw-reload.sh
```

### 0.4 — Verify all services after fixes

```bash
systemctl status openclaw codex-watch openclaw-reload codex-pull.timer
# All should be active/running (or active/waiting for timer)
```

---

## PHASE 1: WEBHOOK INGESTION PIPELINE (Days 1-2)

### 1.1 — Start n8n Container

n8n is already in docker-compose.yml under the `full` profile.

```bash
ssh root@104.248.4.212

# Start n8n alongside existing services
cd /home/onemind/onemind-os/infra
docker-compose --profile full up -d n8n

# Verify n8n is running
docker ps | grep n8n
curl -s http://localhost:5678/healthz

# Add Caddy route for n8n
# In Caddyfile, add:
# n8n.onemindos.dev {
#   reverse_proxy localhost:5678
# }
```

### 1.2 — Install Synadia NATS Node for n8n

```bash
docker exec -it onemind-n8n sh -c "cd /home/node && npm install n8n-nodes-synadia"
docker restart onemind-n8n
```

### 1.3 — Create Zoho CRM Webhook Pipeline

n8n workflow:
1. **Webhook Trigger** → `POST /webhook/zoho-crm`
2. **Function Node** → Normalize to CloudEvent envelope
3. **NATS Publish** → subject `saas.zoho.crm.{resource}.{action}`

Subject mapping:
```
saas.zoho.crm.deal.created
saas.zoho.crm.deal.updated
saas.zoho.crm.deal.closed
saas.zoho.crm.lead.created
saas.zoho.crm.contact.updated
saas.zoho.desk.ticket.created
saas.zoho.desk.ticket.resolved
saas.zoho.calendar.event.created
saas.zoho.calendar.event.updated
```

### 1.4 — Create Todoist Webhook Pipeline

n8n workflow:
1. **Webhook Trigger** → `POST /webhook/todoist`
2. **Function Node** → CloudEvent envelope
3. **NATS Publish** → subject `saas.todoist.{resource}.{action}`

Subject mapping:
```
saas.todoist.task.created
saas.todoist.task.completed
saas.todoist.task.updated
saas.todoist.task.deleted
saas.todoist.project.created
saas.todoist.project.updated
```

### 1.5 — Create Circle.so Webhook Pipeline (Future)

n8n workflow:
1. **Webhook Trigger** → `POST /webhook/circle`
2. **Function Node** → CloudEvent envelope
3. **NATS Publish** → subject `saas.circle.{resource}.{action}`

Subject mapping:
```
saas.circle.member.joined
saas.circle.member.left
saas.circle.post.created
saas.circle.comment.created
saas.circle.event.rsvp
```

### 1.6 — Add NATS JetStream Stream for SaaS Events

```bash
# Add a new SAAS stream for third-party events
nats stream add SAAS \
  --subjects="saas.>" \
  --retention=limits \
  --max-age=90d \
  --storage=file \
  --replicas=1 \
  --discard=old \
  --max-msgs=-1 \
  --max-bytes=1073741824 \
  --server=nats://localhost:4222
```

### 1.7 — Update OpenClaw Hooks for SaaS Events

Add to `openclaw.json` hooks mappings:
```json
{
  "zoho-crm-deal": {
    "action": "agent",
    "agentId": "empire",
    "match": { "source": "zoho.crm" },
    "deliver": true,
    "channel": "discord"
  },
  "todoist-task": {
    "action": "agent",
    "agentId": "legacy",
    "match": { "source": "todoist" },
    "deliver": false
  },
  "circle-member": {
    "action": "agent",
    "agentId": "trinity",
    "match": { "source": "circle" },
    "deliver": true,
    "channel": "discord"
  }
}
```

---

## PHASE 2: AGENT ROSTER + SHARED MEMORY (Days 2-3)

### 2.1 — The Approved 9-Agent Roster (Confirmed)

**Already configured** — all 9 agents are live in `openclaw.json`. Verify current state:

| Agent | Model | Tools | Memory Scope |
|-------|-------|-------|-------------|
| ⚡ Legacy | gpt-5.2 | `full` + lobster + voice_call | All codex |
| ⚔️ Spartan | gpt-4o | messaging + weather + todoist + obsidian | HP (25-49) |
| 🔮 Oracle | gpt-4o | messaging + weather + todoist + obsidian | HP (25-49) |
| 🌳 Eden | gpt-4o | messaging + todoist + obsidian | LE (50-74) |
| 👑 Empire | gpt-5.2 | messaging + todoist + obsidian + zoho-crm | GE (75-99) |
| 🌿 Haven | gpt-4o | messaging + home-assistant + weather + todoist + obsidian | LE (50-74) |
| 🛡️ Guardian | gpt-4o | messaging + onepassword + todoist + obsidian | UI (00-24) |
| ⚙️ Grid | gpt-4o | messaging + github-api + obsidian + exec | UI (00-24) |
| 🔺 Trinity | gpt-4o | messaging + discord + todoist + obsidian | GE (75-99) |

### 2.2 — Shared Memory Architecture (Implement)

Current state: All agents search the full codex via `memorySearch.extraPaths`. Need to add domain-scoped priority paths.

**Config changes for `openclaw.json`:**

Each agent gets its domain codex path as the FIRST entry in `extraPaths`:

```json
// Spartan — Body & Physical Mastery
{
  "id": "spartan",
  "memorySearch": {
    "extraPaths": [
      "/home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)",
      "/home/onemind/onemind-os/codex"
    ]
  }
}

// Empire — Business & Wealth
{
  "id": "empire",
  "memorySearch": {
    "extraPaths": [
      "/home/onemind/onemind-os/codex/75-99 GE (Generational Endeavor)",
      "/home/onemind/onemind-os/codex"
    ]
  }
}

// Haven — Home, Family & Estate
{
  "id": "haven",
  "memorySearch": {
    "extraPaths": [
      "/home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)",
      "/home/onemind/onemind-os/codex"
    ]
  }
}

// Grid — Infrastructure Ops
{
  "id": "grid",
  "memorySearch": {
    "extraPaths": [
      "/home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)",
      "/home/onemind/onemind-os/codex"
    ]
  }
}
```

### 2.3 — Create Missing Agent Workspace Files

Each agent workspace needs at minimum: `SOUL.md`, `AGENTS.md`, `HEARTBEAT.md`

Workspaces to populate (on VPS):
```
/home/onemind/onemind-os/openclaw/workspaces/spartan/
/home/onemind/onemind-os/openclaw/workspaces/oracle/
/home/onemind/onemind-os/openclaw/workspaces/eden/
/home/onemind/onemind-os/openclaw/workspaces/empire/
/home/onemind/onemind-os/openclaw/workspaces/haven/
/home/onemind/onemind-os/openclaw/workspaces/guardian/
/home/onemind/onemind-os/openclaw/workspaces/grid/
/home/onemind/onemind-os/openclaw/workspaces/trinity/
```

Each gets:
- `SOUL.md` — personality, tone, domain focus
- `AGENTS.md` — operating instructions, domain scope, escalation rules
- `HEARTBEAT.md` — what to check on each heartbeat cycle
- `memory/` — empty directory for daily logs

### 2.4 — Per-Agent Model Configuration (Already Live)

Current config confirmed working:
- Legacy: `azure-openai/gpt-5.2`
- Empire: `azure-openai/gpt-5.2`
- All others: `azure-openai/gpt-4o` (from defaults)
- Subagents: `azure-openai/gpt-4o-mini`

No changes needed — this is already correctly configured.

---

## PHASE 3: OMNIPRESENT MODE (Days 3-4)

### 3.1 — Enable Telegram Channel

```json
// Add to openclaw.json channels
"telegram": {
  "botToken": "${TELEGRAM_BOT_TOKEN}",
  "allowFrom": ["${ZEUS_TELEGRAM_ID}"]
}
```

Add Telegram bot token to 1Password:
```bash
op item create --category=login \
  --vault="00-24 SO (Sovereign Operations)" \
  --title="Telegram Bot" \
  "bot-token[password]=YOUR_TOKEN"
```

Add to env.tpl:
```
TELEGRAM_BOT_TOKEN=op://aaeyris5z44pdqspk277u2oppi/Telegram Bot/bot-token
ZEUS_TELEGRAM_ID=op://aaeyris5z44pdqspk277u2oppi/Telegram Bot/zeus-id
```

### 3.2 — Session Architecture for Omnipresent

Current config:
```json
{
  "session": { "dmScope": "per-channel-peer" },
  "agents": {
    "defaults": {
      "subagents": { "maxConcurrent": 8, "maxSpawnDepth": 3 }
    }
  }
}
```

**Omnipresent upgrade:**
- Main lane stays at 4 concurrent (OpenClaw default, sufficient for current load)
- Subagent lane at 8 (already configured)
- Session reset: `idleMinutes: 10080` (7 days, for Lossless Claw context retention)
- Per-channel-peer isolation: already correct

Sessions are automatically created per channel+peer. No changes needed — OpenClaw already supports simultaneous conversations across Discord, iMessage, and Telegram. Each gets its own session context.

### 3.3 — Heartbeat Updates for Omnipresent

Update heartbeat targets for multi-channel:
- Legacy: heartbeat → `target: "last"` (replies to most recent channel)
- Domain agents: heartbeat → `target: "telegram"` (consolidated updates)

### 3.4 — Document Session Architecture

The session model explained:

```
OMNIPRESENT MODE — How Sessions Work

Every conversation gets its own session, keyed by channel+peer:

  Discord #legacy-chat + Zeus    → Session A (main lane)
  iMessage + Zeus                → Session B (main lane)
  Telegram + Zeus                → Session C (main lane)
  Cron: morning-brief            → Session D (cron lane)
  NATS hook: zoho deal closed    → Session E (main lane)
  Subagent: research task        → Session F (subagent lane)

All sessions run in parallel. No blocking between channels.

Lane Limits (concurrent tasks per lane):
  main:     4 simultaneous conversations
  cron:     3 simultaneous scheduled jobs
  subagent: 8 simultaneous spawned workers

Session Lifecycle:
  1. Message arrives on channel+peer
  2. OpenClaw finds or creates session for that key
  3. Task enters the appropriate lane queue
  4. If lane has capacity → process immediately
  5. If lane full → FIFO queue (next available slot)
  6. Session stays alive for 7 days (idleMinutes: 10080)
  7. Lossless Claw retains full context within session
  8. After 7 days idle → session resets, memory persisted in daily logs
```

---

## PHASE 4: REPO ORGANIZATION (Days 4-5)

### 4.1 — Move MEMORY-ARCHITECTURE.md to Codex

Currently sitting at repo root (was moved from codex research folder). Should live in the codex under the agents section.

```bash
mv /Users/zeuslegacy/Documents/onemindos-v0.4/MEMORY-ARCHITECTURE.md \
   "/Users/zeuslegacy/Documents/onemindos-v0.4/codex/00-24 SO (Sovereign Operations)/02 Agents (AI, Memory & Tools)/MEMORY-ARCHITECTURE.md"
```

### 4.2 — Repo Root Cleanup

Current root has accumulated docs. Organize:

| File | Action | Destination |
|------|--------|------------|
| `ACADEMY-CATALOG.md` | KEEP at root | Marketing/reference doc |
| `PRODUCTION-ROADMAP.md` | KEEP at root | Active work tracker |
| `STATUS.md` | KEEP at root | Live system status |
| `AGENTS.md` | KEEP at root | Repo instructions |
| `CLAUDE.md` | KEEP at root | Claude instructions |
| `BLUEPRINT.md` | Move to codex | `codex/00-24 SO/08 Infrastructure/` |
| `DEVELOPMENT.md` | Move to codex | `codex/00-24 SO/10 Development/` |
| `PRODUCTION-PLAN.md` | Archive | Replaced by PRODUCTION-ROADMAP.md |
| `MEMORY-ARCHITECTURE.md` | Move to codex | `codex/00-24 SO/02 Agents/` |
| `GITHUB-RESEARCH.md` | Move to codex | `codex/00-24 SO/12 Research/` |
| `GPU-ANALYSIS.md` | Move to codex | `codex/00-24 SO/12 Research/` |
| `FEATURES.md` | Move to roadmap | `roadmap/features/` |
| `SWARMCLAW.md` | Move to command-center | `command-center/` |
| `1PASSWORD-CLEANUP.md` | Archive | Done, no longer needed |
| `ACCESS.md` | Move to codex | `codex/00-24 SO/11 Security/` |
| `CODE.md` | Review | May be stale |
| `open-source-embedding-models-2026.md` | Move to codex | `codex/00-24 SO/12 Research/` |

### 4.3 — Standardize Agent Workspaces

Ensure all 8 domain agent workspaces match Legacy's structure:

```
workspaces/{agent}/
├── SOUL.md          ← personality, tone, values
├── AGENTS.md        ← operating instructions
├── HEARTBEAT.md     ← proactive checklist
├── IDENTITY.md      ← name, vibe, emoji
├── memory/          ← daily logs (YYYY-MM-DD.md)
└── lobster/         ← named routines (optional)
```

---

## PHASE 5: COMMUNITY LAUNCH PREPARATION (Week 2-3)

### 5.1 — Domain & Brand Finalization

| Asset | Status | Action |
|-------|--------|--------|
| `onemindos.dev` | **LIVE** | Keep — primary tech domain |
| `onemindclaw.com` | Available | **BUY** — OpenClaw-branded landing page |
| One Mind OS brand guide | Not created | Create: logo, colors, fonts, voice |
| Academy landing page | Not created | Build on `academy.onemindos.dev` |

### 5.2 — Circle.so Community Setup

| Space | Tier Access | Purpose |
|-------|------------|---------|
| Announcements | All | Updates, releases |
| Introductions | All | New member intros |
| General Discussion | Explorer+ | Framework discussion |
| Show & Tell | Operator+ | Share your vault, agent configs |
| Office Hours | Operator+ | Weekly live Q&A |
| Builder Lab | Builder | Code, configs, technical deep-dives |
| Templates | Operator+ | Vault templates, skills, configs |
| 1:1 Coaching | Builder | Direct access to Zeus |

### 5.3 — Content Pipeline

| Week | Content | Platform |
|------|---------|----------|
| W1 | "Why Your AI Doesn't Know You" blog post | Blog + Circle |
| W1 | Codex Starter Template (Level 1) | GitHub + Circle |
| W2 | "The CODEX Method" video (20min) | YouTube + Circle |
| W2 | Level 1 Module A1-A3 course content | Circle Courses |
| W3 | "Meet My AI Team" demo (live system walkthrough) | YouTube + Circle |
| W3 | Level 1 Module A4-A6 | Circle Courses |
| W4 | Level 2 opens for Operator tier | Circle |
| W4 | Agent configuration templates | GitHub + Circle |

### 5.4 — Trinity Agent for Community

Configure Trinity to manage Circle.so:
- Welcome new members (Circle webhook → NATS → Trinity)
- Monitor discussions (daily digest)
- Moderate content (flagged posts)
- Track engagement metrics

---

## PHASE 6: PRODUCTION HARDENING (Week 3-4)

### 6.1 — Monitoring & Alerting

| Component | Monitoring | Alert |
|-----------|-----------|-------|
| OpenClaw Gateway | Datadog APM + `journalctl` | Discord #grid-chat if down |
| NATS JetStream | Dashboard (18300) + health pub | Grid agent if streams unhealthy |
| OmniDB | pg_isready health check | Grid agent if connection fails |
| n8n | Execution history + error webhook | Legacy if workflow fails |
| Disk usage | cron `df -h` → NATS | Grid agent if >85% |
| Memory usage | cron `free -h` → NATS | Grid agent if >90% |

### 6.2 — Backup Strategy

| Data | Method | Frequency | Retention |
|------|--------|-----------|-----------|
| Codex | codex-watch → GitHub | Real-time | Git history |
| Workspaces | codex-watch → GitHub | Real-time | Git history |
| OmniDB | `pg_dump` → Backblaze B2 | Daily | 30 days |
| n8n workflows | n8n export → Git | Weekly | Git history |
| OpenClaw config | In monorepo | Real-time | Git history |
| NATS streams | JetStream file storage | Built-in | Per-stream config |

### 6.3 — Security Checklist

- [ ] All secrets in 1Password (no plaintext .env)
- [ ] Tailscale-only access (no public ports except 443 via Caddy)
- [ ] SSH key-only authentication
- [ ] fail2ban active
- [ ] UFW firewall rules reviewed
- [ ] n8n behind auth (self-hosted, access via Tailscale only)
- [ ] OpenClaw gateway token rotated
- [ ] Discord bot token in 1Password
- [ ] Read-only DB access for agents

---

## MASTER CHECKLIST — Everything in Order

### Immediate (Today — March 18)
- [ ] **P0** Fix codex-watch.service (missing on VPS)
- [ ] **P0** Fix codex-pull.service + timer (missing on VPS)
- [ ] **P0** Fix openclaw-reload.service (dead since Mar 16)
- [ ] **P0** Verify all 4 services running after fix

### This Week (March 18-22)
- [ ] **P0** Start n8n container (`docker-compose --profile full up -d`)
- [ ] **P0** Install Synadia NATS node in n8n
- [ ] **P0** Create Zoho CRM → NATS webhook pipeline in n8n
- [ ] **P0** Create Todoist → NATS webhook pipeline in n8n
- [ ] **P0** Add SAAS JetStream stream (`saas.>`, 90d retention)
- [ ] **P1** Update OpenClaw hooks for SaaS events
- [ ] **P1** Add domain-scoped `extraPaths` to each agent's memorySearch
- [ ] **P1** Create SOUL.md + AGENTS.md + HEARTBEAT.md for all 8 domain agents
- [ ] **P1** Enable Telegram channel in openclaw.json
- [ ] **P1** Add Telegram bot token to 1Password + env.tpl
- [ ] **P2** Move MEMORY-ARCHITECTURE.md to codex/02 Agents
- [ ] **P2** Move BLUEPRINT.md, GPU-ANALYSIS.md, etc. to codex
- [ ] **P2** Add Caddy route for n8n (`n8n.onemindos.dev`)

### Next Week (March 23-29)
- [ ] **P1** Configure Circle.so webhook → NATS pipeline (when Circle account ready)
- [ ] **P1** Create Circle.so community spaces matching tier structure
- [ ] **P1** Build Academy landing page on `academy.onemindos.dev`
- [ ] **P1** Publish Codex Starter Template on GitHub
- [ ] **P2** Write "Why Your AI Doesn't Know You" blog post
- [ ] **P2** Record "The CODEX Method" video (20min)
- [ ] **P2** Theme SwarmClaw → One Mind brand

### Month 1 (March-April)
- [ ] **P1** Launch Explorer tier (free) on Circle.so
- [ ] **P1** Launch Operator tier ($29/mo) with Level 1-2 content
- [ ] **P1** Create Level 1 Module A1-A6 course content
- [ ] **P2** Set up OmniDB daily backup to Backblaze B2
- [ ] **P2** Deploy security hardening checklist
- [ ] **P2** Buy `onemindclaw.com` domain
- [ ] **P3** Evaluate ClawVault for structured memory

### Month 2 (May)
- [ ] **P1** Launch Architect tier ($79/mo) with Level 3 content
- [ ] **P1** Create Level 2-3 course content
- [ ] **P2** Deploy legacy-eyes (vision from MacBook)
- [ ] **P2** Enable WhatsApp channel
- [ ] **P2** Self-hosted TTS (replace OpenAI TTS)

### Month 3+ (June+)
- [ ] **P1** Launch Builder tier ($199/mo) with Level 4 content
- [ ] **P2** Mac Studio M4 Ultra for local inference
- [ ] **P2** LiveKit persistent voice room
- [ ] **P3** Homestead IoT (when property acquired)

---

## REVENUE PROJECTION

| Month | Tier | Members | MRR |
|-------|------|---------|-----|
| M1 (April) | Explorer | 50 | $0 |
| M1 (April) | Operator | 10 | $290 |
| M2 (May) | Explorer | 150 | $0 |
| M2 (May) | Operator | 30 | $870 |
| M2 (May) | Architect | 5 | $395 |
| M3 (June) | Explorer | 300 | $0 |
| M3 (June) | Operator | 60 | $1,740 |
| M3 (June) | Architect | 15 | $1,185 |
| M3 (June) | Builder | 3 | $597 |
| **M3 Total** | | **378** | **$3,522/mo** |

The live system IS the marketing. Show, don't tell.

---

## VPS RESOURCE BUDGET

Current: 4 CPU, 16GB RAM, 193GB disk (21% used)

| Service | RAM | CPU | Disk |
|---------|-----|-----|------|
| OpenClaw Gateway (9 agents) | ~600MB | Variable | Minimal |
| NATS JetStream | ~100MB | Low | ~2GB (streams) |
| OmniDB (PostgreSQL) | ~500MB | Low | ~5GB (growing) |
| legacy-bridge | ~50MB | Low | Minimal |
| nats-ingester | ~50MB | Low | Minimal |
| n8n | ~300MB | Low | ~500MB |
| SwarmClaw (Next.js) | ~200MB | Low | ~200MB |
| OpenClaw Office (Vite) | ~100MB | Low | ~100MB |
| OpenViking | ~200MB | Low | ~1GB |
| legacy-pulse | ~50MB | Low | Minimal |
| Caddy | ~30MB | Low | Minimal |
| Datadog Agent | ~200MB | Low | ~500MB |
| **TOTAL** | **~2.4GB** | | **~10GB** |
| **Available** | **~12.6GB** | | **~183GB** |

Plenty of headroom. No upgrade needed yet.

---

> **Next action**: Fix the 3 broken services on VPS, then start n8n and build the first webhook pipeline.
