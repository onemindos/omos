---
id: ui-08-openclaw-deployment
title: "OpenClaw Gateway — Live Deployment Reference"
type: operations
status: active
created: 2026-03-05
modified: 2026-03-05
owner: zeus
tags: [type/operations, status/active, domain/ui, platform/openclaw]
ai_summary: "Complete reference for the live OpenClaw gateway deployment on DigitalOcean VPS. Covers server, agents, voice, channels, devices, networking, and environment."
---

# OpenClaw Gateway — Live Deployment Reference

> **Last verified**: 2026-03-05
> **Status**: OPERATIONAL — All systems green
> **Version**: OpenClaw `2026.3.2`

---

## 1. Server Infrastructure

| Property | Value |
|----------|-------|
| **Provider** | DigitalOcean Droplet |
| **Public IP** | `134.122.8.135` |
| **Tailscale IP** | `100.104.61.29` |
| **Tailscale hostname** | `onemind-gateway` |
| **OS** | Ubuntu (systemd) |
| **SSH Access** | `ssh root@134.122.8.135` or `ssh onemind@100.104.61.29` |
| **Service user** | `onemind` |
| **Service manager** | `systemctl {start|stop|restart|status} openclaw` |
| **Logs** | `journalctl -u openclaw -f --no-pager` |

### Key Paths (on VPS)

| Resource | Path |
|----------|------|
| Config | `/home/onemind/.openclaw/openclaw.json` |
| Environment | `/opt/onemind/.env` |
| Legacy workspace | `~/.openclaw/workspace-legacy` |
| HP workspace | `~/.openclaw/workspace-hp` |
| LE workspace | `~/.openclaw/workspace-le` |
| GE workspace | `~/.openclaw/workspace-ge` |
| Agent dirs | `~/.openclaw/agents/{id}/agent` |

---

## 2. Networking & Access

### Gateway Listener

```
ws://0.0.0.0:18789  (bind: lan)
```

| Setting | Value | Notes |
|---------|-------|-------|
| **Port** | `18789` | WebSocket gateway |
| **Bind** | `lan` (`0.0.0.0`) | Accepts from any interface |
| **Mode** | `local` | No cloud relay |
| **Auth** | Token-based | `${OPENCLAW_GATEWAY_TOKEN}` |

### Control UI (Studio)

| Setting | Value |
|---------|-------|
| **Listen** | `http://127.0.0.1:18791/` (loopback only) |
| **Auth** | Token (`STUDIO_ACCESS_TOKEN`) |
| **Allowed Origins** | `*` (wide open — temporary) |
| **Insecure Auth** | `true` |
| **Device Auth** | Disabled (`dangerouslyDisableDeviceAuth: true`) |

> ⚠️ **Security Note**: Control UI is currently wide open for development. Before production, tighten: remove wildcard origins, enable device auth, disable insecure auth.

### How Devices Connect

```
iPhone / Mac / Browser
    ↓  Tailscale VPN (100.x.x.x)
    ↓  ws://100.104.61.29:18789
    ↓  Token auth → Device pairing → Connected
Gateway (0.0.0.0:18789)
```

---

## 3. Agents (Hive Mind)

Four agents run in the gateway, each with dedicated workspaces and Discord channel bindings:

| Agent | ID | Model | Heartbeat | Discord Channel |
|-------|----|-------|-----------|-----------------|
| **LEGACY** (default) | `legacy` | `azure-openai/gpt-5.2` | 1h → heartbeat channel | Default (all DMs + unbound) |
| **HP** (Holistic Performance) | `hp` | `azure-openai/gpt-4o` | 1h (08:00–22:00 ET) | `1478604249813221528` |
| **LE** (Legacy Evolution) | `le` | `azure-openai/gpt-4o` | 1h (09:00–20:00 ET) | `1478604259246215361` |
| **GE** (Generational Entrepreneurship) | `ge` | `azure-openai/gpt-5.2` | 1h (07:00–23:00 ET) | `1478604264942207017` |

### Agent Identity

- **LEGACY** has `IDENTITY.md` (🧭 Legacy), `SOUL.md`, `HEARTBEAT.md`, `BOOT.md` in workspace
- Each agent has isolated workspace + agent dir
- Subagent model: `azure-openai/gpt-4o-mini`
- Thinking default: `medium`
- Context pruning: cache-ttl, 1h TTL
- Compaction: safeguard mode

### Agent Routing (Discord Bindings)

```
Discord channel 1478604249813221528 → HP agent
Discord channel 1478604259246215361 → LE agent
Discord channel 1478604264942207017 → GE agent
Everything else (DMs, other channels) → LEGACY (default)
```

### Inter-Agent Communication

```json
"agentToAgent": {
  "enabled": true,
  "allow": ["legacy", "hp", "le", "ge"]
}
```

All four agents can talk to each other. LEGACY can delegate to HP/LE/GE and vice versa.

---

## 4. Voice & TTS — Dual Pipeline

### Pipeline 1: iOS/Device Talk Mode (ElevenLabs)

This is the `talk` section — used by the iOS app and any device node connecting via WebSocket.

```
You speak → iPhone mic → Apple SFSpeechRecognizer (on-device STT)
    ↓
Transcription → sent to gateway → LLM processes
    ↓
Response text → ElevenLabs API (server-side TTS)
    ↓
Audio streamed back → iPhone speaker (PCM 24kHz)
```

| Setting | Value |
|---------|-------|
| **Provider** | `elevenlabs` |
| **Model** | `eleven_v3` (latest, most natural) |
| **Voice** | `pNInz6obpgDQGcFmaJgB` (Adam — warm, deep male) |
| **API Key** | `${ELEVENLABS_API_KEY}` (from env) |
| **Output Format** | `pcm_24000` (low-latency streaming) |
| **Interrupt on Speech** | `true` (stops speaking when you start) |

#### iOS Talk Mode Details

- **STT**: Apple `SFSpeechRecognizer` (on-device, no API cost)
- **Audio Capture**: `AVAudioEngine` (real-time)
- **TTS Client**: `ElevenLabsTTSClient` (streams PCM, falls back to MP3)
- **Fallback**: Apple `AVSpeechSynthesizer` (system voice if ElevenLabs unavailable)
- **Wake Words**: `SFSpeechRecognizer` continuous listening (configurable trigger words)

#### Config Block

```json
"talk": {
  "provider": "elevenlabs",
  "providers": {
    "elevenlabs": {
      "apiKey": "${ELEVENLABS_API_KEY}",
      "voiceId": "pNInz6obpgDQGcFmaJgB",
      "modelId": "eleven_v3"
    }
  },
  "interruptOnSpeech": true,
  "outputFormat": "pcm_24000"
}
```

### Pipeline 2: Discord Voice Channel (OpenAI TTS)

This is the `channels.discord.voice.tts` section — used when Legacy speaks in a Discord voice channel.

```
You speak in Discord voice channel
    ↓
Discord bot receives Opus audio stream
    ↓
Gateway decodes → STT (server-side)
    ↓
Transcription → LLM processes
    ↓
Response text → OpenAI gpt-4o-mini-tts (server-side)
    ↓
Audio → Discord voice channel (Opus)
```

| Setting | Value |
|---------|-------|
| **Provider** | `openai` |
| **Model** | `gpt-4o-mini-tts` |
| **Voice** | `alloy` |
| **API Key** | `${AZURE_OPENAI_API_KEY}` (Azure endpoint) |
| **DAVE Encryption** | Enabled |

#### Config Block

```json
"channels": {
  "discord": {
    "voice": {
      "enabled": true,
      "autoJoin": [{
        "guildId": "1453411552181883027",
        "channelId": "1478851954791678092"
      }],
      "daveEncryption": true,
      "tts": {
        "provider": "openai",
        "openai": { "voice": "alloy" }
      }
    }
  }
}
```

### Pipeline 3: Discord Text-to-Speech (Message TTS)

When Legacy reads messages aloud (text message TTS, separate from voice channel):

```json
"messages": {
  "tts": {
    "provider": "openai",
    "openai": {
      "apiKey": "${AZURE_OPENAI_API_KEY}",
      "model": "gpt-4o-mini-tts",
      "voice": "alloy"
    }
  }
}
```

### Voice Summary Table

| Context | Provider | Model/Voice | API Key Source | Latency |
|---------|----------|-------------|---------------|---------|
| **iOS Talk Mode** | ElevenLabs | `eleven_v3` / Adam | `ELEVENLABS_API_KEY` | ~500ms (streaming PCM) |
| **Discord Voice** | OpenAI | `gpt-4o-mini-tts` / alloy | `AZURE_OPENAI_API_KEY` | ~800ms |
| **Discord Msg TTS** | OpenAI | `gpt-4o-mini-tts` / alloy | `AZURE_OPENAI_API_KEY` | ~800ms |

---

## 5. Channels

### Discord

| Setting | Value |
|---------|-------|
| **Enabled** | `true` |
| **Bot Token** | `${DISCORD_BOT_TOKEN}` |
| **Guild** | `1453411552181883027` (OneMind server) |
| **Require Mention** | `false` (responds to all messages in guild) |
| **DM Policy** | `open` (responds to all DMs) |
| **Group Policy** | `open` |
| **Streaming** | `progress` (shows typing indicator) |
| **Thread Bindings** | Enabled, 24h idle timeout, spawn subagent sessions |
| **Voice Auto-Join** | Channel `1478851954791678092` in guild |

---

## 6. Paired Devices

| Device | Role | Tailscale IP | Status |
|--------|------|-------------|--------|
| **iPhone** (Zeus's iPhone 17 Pro Max) | node, operator | `100.122.203.68` | ✅ Paired |
| **Zeus-Mac** | node | — | ✅ Paired |
| **Zeus's MacBook Pro** | node, operator | — | ✅ Paired |
| + 5 additional operator tokens | operator | — | ✅ Paired |

**Total**: 8 paired devices/tokens

### iOS App Build Info

| Property | Value |
|----------|-------|
| **Bundle ID** | `ai.openclaw.ios.zeus` |
| **Signing Team** | `Q8J6PP8YTF` (Zeus Delacruz, Personal) |
| **Deployment Target** | iOS 18.0 |
| **Xcode** | 26.3 (Swift 6.2.4) |
| **CoreDevice ID** | `C9D6A570-1E93-5DE3-A3BA-5E5CEF7471F7` |

### iOS Source Modifications (from upstream)

Three bugs were fixed to enable Tailscale/LAN connectivity:

1. **TLS Override** — `GatewayConnectionController.swift`: Added `isPrivateOrTailscaleIP()` to allow `ws://` for RFC-1918 + Tailscale CGNAT (100.64-127.x) IPs. Without this, app forced `wss://` for all non-loopback hosts.

2. **ATS Conflict** — `Info.plist` + `project.yml`: Removed `NSAllowsArbitraryLoadsInWebContent` which caused iOS to *ignore* `NSAllowsArbitraryLoads: true`. Now only `NSAllowsArbitraryLoads: true` remains.

3. **Build Compatibility** — `project.yml`: Removed WatchApp/ActivityWidget dependencies, disabled push/live activities, disabled SwiftFormat/SwiftLint build scripts (not installed locally).

---

## 7. LLM Provider Configuration

### Azure OpenAI (Primary)

| Setting | Value |
|---------|-------|
| **Base URL** | `https://onemind-openai.openai.azure.com/openai/v1` |
| **API Format** | `openai-responses` |
| **Auth** | API key in `api-key` header (not Bearer) |

### Available Models

| Model ID | Name | Capabilities |
|----------|------|-------------|
| `gpt-5.2` | GPT-5.2 | text, image, reasoning |
| `gpt-5.2-codex` | GPT-5.2 Codex | text, image, reasoning |
| `gpt-5.1-codex-mini` | GPT-5.1 Codex Mini | text, image, reasoning |
| `gpt-5-mini` | GPT-5 Mini | text, image |
| `gpt-5-nano` | GPT-5 Nano | text, image |
| `gpt-4o` | GPT-4o | text, image |
| `gpt-4o-mini` | GPT-4o Mini | text, image |
| `o1` | O1 | text, image, reasoning |
| `o3-mini` | O3 Mini | text, image, reasoning |

### Model Aliases (shortcuts)

```
gpt52     → azure-openai/gpt-5.2
gpt5mini  → azure-openai/gpt-5-mini
gpt5nano  → azure-openai/gpt-5-nano
gpt4o     → azure-openai/gpt-4o
gpt4omini → azure-openai/gpt-4o-mini
o1        → azure-openai/o1
o3mini    → azure-openai/o3-mini
llama70b  → gradient/llama3.3-70b-instruct
```

---

## 8. Tools & Capabilities

| Tool | Status | Notes |
|------|--------|-------|
| **Web Search** | ✅ Enabled | Max 8 results |
| **Web Fetch** | ✅ Enabled | 50K char cap |
| **Exec** | ✅ Enabled | Host: gateway, allowlist security |
| **Apply Patch** | ✅ Enabled | File editing via patches |
| **Loop Detection** | ✅ Enabled | Warning: 10, Critical: 20 |
| **Lobster** (browser) | ✅ Plugin enabled | LEGACY + HP + GE |
| **Agent-to-Agent** | ✅ Enabled | All 4 agents can communicate |

---

## 9. Hooks & Webhooks

| Hook ID | Trigger | Action | Target Agent |
|---------|---------|--------|-------------|
| `nats-event` | NATS source | Wake agent | (default) |
| `agno-result` | Agno worker result | Deliver to agent | `legacy` |
| `ha-alert` | Home Assistant alert | Deliver to Discord | `legacy` |
| `circle-member` | Circle.so member event | Process (no deliver) | `community` |
| `gmail` | Gmail preset | (built-in) | — |

**Hooks endpoint**: `/hooks` (token: `${OPENCLAW_HOOKS_TOKEN}`)

---

## 10. Scheduled Tasks

| Feature | Setting |
|---------|---------|
| **Cron** | Enabled, max 3 concurrent |
| **Heartbeats** | Every 1h per agent (active hours vary) |
| **Broadcast** | Parallel strategy |

---

## 11. Environment Variables (on VPS)

All secrets live in `/opt/onemind/.env`. Names only (values redacted):

| Category | Variables |
|----------|-----------|
| **Gateway Auth** | `OPENCLAW_GATEWAY_TOKEN`, `OPENCLAW_HOOKS_TOKEN`, `STUDIO_ACCESS_TOKEN` |
| **LLM Keys** | `ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `AZURE_OPENAI_API_KEY`, `AZURE_OPENAI_ENDPOINT`, `GRADIENT_API_KEY` |
| **Voice** | `ELEVENLABS_API_KEY`, `ELEVENLABS_VOICE_ID`, `LIVEKIT_API_KEY`, `LIVEKIT_API_SECRET`, `LIVEKIT_URL` |
| **Discord** | `DISCORD_BOT_TOKEN`, `HEARTBEAT_CHANNEL_ID`, `ZEUS_DISCORD_USER_ID`, `DISCORD_COMMUNITY_CHANNEL_ID` |
| **Database** | `OMNIDB_URL`, `OMNIDB_DB`, `OMNIDB_USER`, `OMNIDB_PASSWORD` |
| **External** | `NATS_URL`, `N8N_WEBHOOK_URL`, `N8N_ENCRYPTION_KEY` |
| **Home Assistant** | `EDGE_FARM_PASSWORD`, `EDGE_HOME_PASSWORD`, `EDGE_OFFICE_PASSWORD` |
| **WhatsApp** | `FAMILY_WHATSAPP_NUMBERS` |
| **Misc** | `BRIDGE_AGENT_ID`, `LEGACY_MODEL`, `LOG_LEVEL`, `MCP_CLIENT_SECRET`, `OPENAI_TTS_BASE_URL`, `OPENCLAW_ALLOW_INSECURE_PRIVATE_WS`, `OPENCLAW_MCP_GATEWAY_URL`, `OPENCLAW_NO_RESPAWN`, `OPENCLAW_WORKSPACE`, `NODE_COMPILE_CACHE` |

---

## 12. Quick Commands Reference

```bash
# ── Server Access ──
ssh root@134.122.8.135                      # Direct SSH
ssh onemind@100.104.61.29                   # Via Tailscale

# ── Service Management ──
systemctl restart openclaw                  # Restart gateway
systemctl status openclaw                   # Check status
journalctl -u openclaw -f --no-pager        # Live logs

# ── Device Management ──
openclaw devices list                       # Show all paired devices
openclaw devices approve <request-id>       # Approve pending device

# ── Agent Commands ──
openclaw agents                             # List all agents
openclaw doctor                             # Health check

# ── Config ──
cat /home/onemind/.openclaw/openclaw.json   # View config
# Edit config → restart gateway for changes to take effect

# ── Logs (filtered) ──
journalctl -u openclaw --since "5 min ago" --no-pager | grep -i "error\|warn"
journalctl -u openclaw --since "5 min ago" --no-pager | grep -i "talk\|tts\|voice"
journalctl -u openclaw --since "5 min ago" --no-pager | grep -i "discord"
```

---

## 13. Known Issues & TODOs

| Issue | Priority | Status |
|-------|----------|--------|
| Mac Studio "device identity required" in Chrome | Medium | 🔶 Unresolved |
| Gateway security wide open (`dangerouslyDisable*`) | High | ⏳ Tighten after all devices confirmed |
| Voice ID hardcoded (Adam) — may want custom clone | Low | 💡 Future |
| No local TTS fallback on server (if ElevenLabs down) | Medium | 💡 Jupiter Voice deferred |
| Heartbeat channel IDs are env vars | Info | ✅ Working |

---

## 14. Architecture Diagram

```
┌──────────────────────────────────────────────────────────────────────┐
│                        OneMind Gateway                               │
│                    DigitalOcean VPS (134.122.8.135)                   │
│                    Tailscale: 100.104.61.29                          │
│                    OpenClaw v2026.3.2                                 │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐    │
│  │  LEGACY    │  │     HP     │  │     LE     │  │     GE     │    │
│  │  GPT-5.2   │  │  GPT-4o    │  │  GPT-4o    │  │  GPT-5.2   │    │
│  │  (default) │  │  health    │  │  family    │  │  business  │    │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘    │
│        │               │               │               │            │
│        └───────────┬────┴───────────────┴───────┬───────┘            │
│                    │                            │                    │
│              Agent-to-Agent                  Bindings               │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                     CHANNELS                                 │    │
│  │                                                             │    │
│  │  Discord (text + voice)          iOS Talk Mode              │    │
│  │  ├─ Bot: all guilds/DMs          ├─ WebSocket :18789        │    │
│  │  ├─ Voice: auto-join             ├─ ElevenLabs eleven_v3    │    │
│  │  ├─ TTS: OpenAI alloy            ├─ Apple STT (on-device)   │    │
│  │  └─ Thread bindings              └─ PCM 24kHz streaming     │    │
│  │                                                             │    │
│  │  Hooks (webhooks)                                           │    │
│  │  ├─ NATS → wake                                             │    │
│  │  ├─ Agno → legacy                                           │    │
│  │  ├─ Home Assistant → discord                                │    │
│  │  └─ Circle.so → community                                   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                     TOOLS                                    │    │
│  │  web_search │ web_fetch │ exec │ lobster │ memory │ patch   │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    LLM PROVIDERS                             │    │
│  │  Azure OpenAI: GPT-5.2, 5-mini, 5-nano, 4o, 4o-mini,      │    │
│  │                o1, o3-mini, 5.2-codex, 5.1-codex-mini      │    │
│  │  Gradient: llama3.3-70b-instruct                            │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│  DEVICES (8 paired)                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │   iPhone 17  │  │  Zeus-Mac    │  │  MacBook Pro │              │
│  │  Pro Max     │  │  (Studio)    │  │              │              │
│  │ 100.122.     │  │              │  │              │              │
│  │ 203.68       │  │              │  │              │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└──────────────────────────────────────────────────────────────────────┘
```
