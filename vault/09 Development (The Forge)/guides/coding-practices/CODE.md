# CODE.md — Codebase Inventory

> Last updated: March 16, 2026
> Excludes: `node_modules`, `.git`, `codex/.obsidian` plugin bundles (vendored)
> Grand total: **~26,860 lines of original code**

---

## Active / Deployed Code

### `openclaw/workspaces/legacy/skills/` — Node.js (21 skills, ~6,900 lines)

Each skill is a CommonJS package (`index.js` + `package.json` + `SKILL.md`) that becomes a live tool for Legacy.

| Skill | Lines | Purpose |
|-------|-------|---------|
| `skill-forge-1.0.0` | 839 | Autonomous skill builder |
| `clickup-1.0.0` | 566 | ClickUp task management |
| `credit-repair-1.0.0` | 461 | Credit management tools |
| `zoho-mail-1.0.0` | 484 | Zoho Mail API |
| `nats-manager-1.0.0` | 445 | NATS monitoring & management |
| `vision-1.0.0` | 393 | Computer vision integration |
| `homestead-control` | 309 | Home/estate control |
| `onepassword-1.0.0` | 292 | 1Password credential queries |
| `rss-digest-1.0.0` | 289 | RSS feed aggregation |
| `uber-1.0.0` | 281 | Uber ride integration |
| `lob-1.0.0` | 281 | Physical mail API |
| `recipe-scraper-1.0.0` | 247 | Recipe data extraction |
| `transit-1.0.0` | 273 | Public transit data |
| `cloudflare-api-1.0.0` | 215 | Cloudflare API control |
| `nexus-bridge` | 206 | NATS dispatch & query |
| `weather-1.0.0` | 184 | Weather API |
| `spotify-1.0.0` | 178 | Spotify playback control |
| `hitl-escalate-1.0.0` | 139 | Human-in-the-loop escalation |
| `resend-1.0.0` | 124 | Email via Resend API |
| `twilio-1.0.0` | 118 | SMS via Twilio |
| `home-assistant-1.0.0` | 189 | Smart home control (6 tools) |

**14 additional skill folders** (`github-api-1.0.3`, `obsidian-1.0.0`, `browser-1.0.0`, `todoist`, `todoist-api`, `zoho-cli`, `api-gateway-1.0.48`, `proactive-agent-3.1.0`, `agent-browser-0.2.0`, `find-skills-0.1.0`, etc.) are **doc/config-only** — they contain `SKILL.md` and `_meta.json` but no runnable code. The skill logic is built into the OpenClaw runtime for those.

**Skills also include ~2,300 lines of Python:**
- `desktop-control-1.0.0/` — `ai_agent.py`, `demo.py`, `__init__.py` (desktop automation)
- `legacy-eyes/legacy-eyes.py` — Computer vision bridge
- `legacy-pulse/legacy-pulse.py` + `legacy-session-start.py` — Session/heartbeat tracking
- `skill-creator-0.1.0/scripts/` — `init_skill.py`, `package_skill.py`, `quick_validate.py`

---

### `openclaw/workspaces/legacy/` (root files) — Python + Bash (~620 lines)

| File | Lines | Purpose |
|------|-------|---------|
| `bridge.py` | 267 | NATS ↔ OpenClaw webhook relay (FastAPI) |
| `project_manager.py` | 351 | Project planning & tracking logic |
| `launch.sh` / `stop.sh` | — | Gateway startup/shutdown scripts |
| `scripts/todoist_api.sh` | — | Todoist REST API helpers |
| `scripts/add_comment.sh` | — | Add comment to Todoist task |
| `scripts/sync_task.sh` | — | Sync task state |

---

### `infra/nats/scripts/` — Python (7 files, ~2,040 lines)

NATS bridge and event ingestion services.

| File | Lines | Purpose |
|------|-------|---------|
| `message_protocol.py` | 461 | Protocol definitions & data models |
| `webhook_receiver.py` | 352 | FastAPI webhook receiver → NATS publisher |
| `nats_gesture_bridge.py` | 349 | Gesture events → NATS normalization |
| `nats_tools.py` | 249 | NATS utility functions |
| `nats_bridge.py` | 225 | Generic NATS bridge |
| `bus.py` | 203 | Event bus abstraction |
| `nats_ingester.py` | 198 | Event ingestion pipeline |

---

### `infra/scripts/` — Bash (4 scripts)

| File | Purpose |
|------|---------|
| `health-check.sh` | VPS health monitoring |
| `install-openclaw.sh` | OpenClaw installation |
| `setup-droplet.sh` | DigitalOcean droplet provisioning |
| `setup-firewall.sh` | UFW firewall configuration |

---

### `openclaw/deploy/scripts/` — Bash (8 scripts, ~1,600 lines combined)

| File | Purpose |
|------|---------|
| `setup.sh` | Full VPS setup from scratch |
| `openclaw-reload.sh` | Hot-reload workspace on VPS |
| `workspace-backup.sh` | Backup VPS workspace → repo |
| `setup-symlinks.sh` | Sync repo workspace → VPS paths |
| `onemind-watch.sh` | Auto-commit codex changes → GitHub |
| `create-1p-items.sh` | Migrate secrets → 1Password |
| `deploy-frontends.sh` | Deploy command center UIs |
| `sentinel-setup.sh` | Sentinel VPS setup |

---

## R&D / Prototype Code

### `roadmap/` — Python + Dart (~7,500 + 1,130 lines)

Not deployed — reference implementations and prototypes.

#### Consciousness (`roadmap/consciousness/reference/`) — 10 Python files

| File | Lines | Purpose |
|------|-------|---------|
| `bus.py` | 653 | Event bus for awareness system |
| `processor.py` | 570 | Event processing pipeline |
| `temporal_models.py` | 480 | Time-series data models |
| `correlator.py` | 465 | Event correlation engine |
| `deep_links.py` | 446 | Context linking |
| `temporal_context.py` | 401 | Temporal reasoning |
| `triage.py` | 374 | Event triage & prioritization |
| `activity_models.py` | 343 | Activity classification |
| `awareness_provider.py` | 271 | Awareness state provider |
| `aware_runner.py` | 213 | Awareness loop runner |

#### Hardware (`roadmap/hardware/`) — 5 Python files

| File | Lines | Purpose |
|------|-------|---------|
| `3d-printing/printer.py` | 962 | 3D printer control |
| `3d-printing/cad.py` | 539 | CAD file generation |
| `biometrics/face_auth.py` | 483 | Facial recognition / face auth |
| `gesture-systems/gesture_service.py` | 233 | Gesture recognition service |
| `edge-nodes/remote_agent.py` | 195 | Remote device agent |

#### OpenMind (`roadmap/openmind/event-filter/`) — 2 Python files

| File | Lines | Purpose |
|------|-------|---------|
| `digest.py` | 350 | NATS event digest & batching |
| `throttle.py` | 196 | Event rate throttling |

#### Voice (`roadmap/voice/realtime-audio/`)

| File | Lines | Purpose |
|------|-------|---------|
| `s2s.py` | 335 | OpenAI Realtime API full-duplex audio |

#### OneIntel widgets (`roadmap/oneintel/widgets/`) — Dart (archived Flutter)

| File | Lines | Purpose |
|------|-------|---------|
| `widgets.dart` | — | Widget barrel export |
| `mini_heatmap.dart` | — | Heatmap calendar widget |
| `progress_ring.dart` | — | Circular progress ring |
| `streak_indicator.dart` | — | Habit streak indicator |
| `trend_chart.dart` | — | Trend line chart |
| **Total** | ~1,128 | Reference designs for TS rebuild |

---

### `scripts/zoho_mail_cleanup/` — Python (20 files, ~4,280 lines)

Full CLI suite for auditing and bulk-managing Zoho Mail.

| File | Lines |
|------|-------|
| `move_to_delete.py` | 454 |
| `scan_deep.py` | 392 |
| `scan_multi.py` | 357 |
| `bloat_audit.py` | 310 |
| `scan_review.py` | 283 |
| `cleanup_unread.py` | 274 |
| `audit.py` | 256 |
| `cleanup_senders.py` | 249 |
| `move_trash_to_review.py` | 244 |
| `move_trash_multi.py` | 213 |
| `cleanup_old.py` | 197 |
| `cleanup_folders.py` | 189 |
| `cleanup_review.py` | 162 |
| `move_evergreen_to_inbox.py` | 138 |
| `config.py` | 112 |
| `zoho_mcp_audit.py` | 65 |
| `menu.py` | 60 |
| `zoho_mcp_policies.py` | 59 |
| `move_2026_to_ceo.py` | 43 |
| `count_by_year.py` | 29 |

---

### `scripts/migrate-config.py` — Python

Config migration utility (~50 lines).

---

## Grand Total

| Area | Language | Lines |
|------|----------|-------|
| Skills — 21 active (`index.js`) | Node.js | ~6,900 |
| Skills — Python components | Python | ~2,300 |
| `bridge.py` + `project_manager.py` | Python | ~620 |
| `infra/nats/scripts/` | Python | ~2,040 |
| Deploy + infra ops scripts | Bash | ~2,090 |
| `roadmap/` — prototype code | Python | ~7,500 |
| `roadmap/oneintel/widgets/` | Dart | ~1,130 |
| `scripts/zoho_mail_cleanup/` | Python | ~4,280 |
| **Total** | | **~26,860 lines** |


---

## Related

- [[00-24 SO (Sovereign Operations)/05 Development (Code & Testing) (The Forge)/05 Development|05 Development]]
