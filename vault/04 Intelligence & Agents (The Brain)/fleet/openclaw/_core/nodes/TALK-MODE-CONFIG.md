# OpenClaw Talk Mode — Configuration & Architecture

**Updated:** 2026-05-11
**Provider:** Azure OpenAI Realtime (`gpt-realtime-1.5`) — ✅ deployed and active
**Azure Resource (voice):** `onemind-voice.openai.azure.com` (East US 2)
**Azure Resource (chat):** `onemind-openai.openai.azure.com` (East US)
**OpenClaw Version:** v2026.5.7 (eeef486)

---

## Active Config (omos-core)

| Setting | Value |
|---|---|
| `talk.provider` | `openai` (Azure OpenAI Realtime) |
| `talk.providers.openai.model` | `gpt-realtime-1.5` |
| `talk.providers.openai.voice` | `coral` |
| `talk.providers.openai.azureEndpoint` | `https://onemind-voice.openai.azure.com` |
| `talk.providers.openai.azureDeployment` | `gpt-realtime-1.5` |
| `talk.providers.openai.azureApiVersion` | `2025-04-01-preview` |
| API Key env | `AZURE_OPENAI_VOICE_KEY` (in `gateway.systemd.env`) |

### omos-core (100.123.126.34) — actual config

```json
{
  "provider": "openai",
  "providers": {
    "openai": {
      "model": "gpt-realtime-1.5",
      "voice": "coral",
      "apiKey": { "source": "env", "provider": "default", "id": "AZURE_OPENAI_VOICE_KEY" },
      "azureEndpoint": "https://onemind-voice.openai.azure.com",
      "azureDeployment": "gpt-realtime-1.5",
      "azureApiVersion": "2025-04-01-preview"
    }
  }
}
```

All env vars are injected via `/home/ubuntu/.openclaw/gateway.systemd.env`.

### Azure Resources

| Resource | Endpoint | Location | Deployments |
|---|---|---|---|
| `onemind-openai` | `onemind-openai.openai.azure.com` | East US | gpt-4o, gpt-4o-mini, gpt-5, gpt-5-mini, gpt-5-nano, gpt-4.1, gpt-4.1-mini, gpt-4.1-nano, text-embedding-3-small |
| `onemind-voice` | `onemind-voice.openai.azure.com` | East US 2 | **gpt-realtime-1.5** ✅, gpt-realtime-mini ✅, gpt-4o-mini-tts ✅, gpt-4o-transcribe ✅ |

### omos-sentinel (100.66.235.38)

Talk mode not yet confirmed active on sentinel. Assume same config shape needed if enabled.

---

## Architecture — How Talk Mode Works End-to-End

Talk mode has **two layers** that work together:

### Layer 1: The Voice Model (Gemini)

When you click Talk in the browser:

1. **Gateway** creates a Google Live API session via `BidiGenerateContentConstrained`
2. **Browser** connects directly to Google over WebSocket (`json-pcm-websocket` transport)
3. **Gemini** handles speech-to-text and text-to-speech natively (no separate ASR/TTS)
4. The voice model receives a **hardcoded system prompt**:

> "You are OpenClaw's realtime voice interface. Keep spoken replies concise. If the user asks for code, repository state, tools, files, current OpenClaw context, or deeper reasoning, call `openclaw_agent_consult` and then summarize the result naturally."

This prompt **cannot be customized** via config (it's in `buildRealtimeInstructions()` in `server-methods-DSL6KpJ1.js`).

### Layer 2: Agent Consult (The Full OneMind Agent)

The voice model gets one tool: `openclaw_agent_consult`. When called:

1. Gemini calls the tool with `question`, optional `context`, and optional `responseStyle`
2. The **browser UI** catches the tool call
3. Sends a `chat.send` to the **active OpenClaw session** (the full agent)
4. The agent processes the request with **all its tools**: memory, web search, file read, MCP, sandbox, etc.
5. The agent also loads its **workspace bootstrap files** as system context:
   - `IDENTITY.md` — agent name, role, vibe
   - `USER.md` — who Zeus is, preferences, projects
   - `SOUL.md` — personality, directives, voice mode behavior
   - `TOOLS.md` — infrastructure details, server map, key paths
6. Agent returns its response (text)
7. Browser sends response back to Gemini as tool result
8. Gemini **speaks the answer naturally** (120s timeout)

### What This Means Practically

- **Simple questions** (time, weather, quick facts) → Gemini answers directly
- **Anything requiring tools/context** → Gemini calls `openclaw_agent_consult` → full agent handles it → Gemini speaks the result
- The agent knows who Zeus is, knows the infrastructure, has full tool access
- Talk effectively gives you **voice access to the entire OneMind system**

---

## Workspace Bootstrap Files (Agent Identity)

These files at `~/.openclaw/workspace/` define how the agent behaves when consulted via Talk (and in all chat interactions):

### IDENTITY.md
- Name: **OneMind**
- Instance: omos-core (primary) / omos-sentinel (secondary)
- Creature: Sovereign AI operating system
- Vibe: Sharp, direct, competent

### USER.md
- Name: **Zeus**
- Timezone: America/New_York
- Context: OneMind OS founder, technical operator
- Preferences: Move don't copy, never delete without permission, Todoist for tasks

### SOUL.md
Core directives including:
- Be competent, not performative
- Know the system, use tools proactively
- Be direct, disagree when appropriate
- Voice mode specific: concise, no markdown, natural speech
- Never delete without permission
- Ask before external-facing actions

### TOOLS.md
Infrastructure map: servers (core/sentinel/titan), OpenClaw config, cloud services (AWS/GCP/Azure), key paths, integrations.

**To modify agent behavior**, edit these files on the server:
```bash
ssh ubuntu@100.123.126.34 'nano ~/.openclaw/workspace/SOUL.md'
# Changes take effect on next agent consult (no restart needed)
```

---

## Google Provider Config Options

All available via `talk.providers.google.<key>` (uses `additionalProperties`):

| Setting | Type | Description | Current |
|---|---|---|---|
| `apiKey` | string | Gemini API key | Set |
| `model` | string | Model ID | `gemini-2.5-flash-native-audio-preview-12-2025` |
| `voice` | string | Voice name | `Kore` |
| `temperature` | number | Response temperature | Unset (model default) |
| `enableAffectiveDialog` | boolean | Emotional expression in voice | `true` |
| `prefixPaddingMs` | number | Audio padding before speech detection | Unset |
| `silenceDurationMs` | number | Silence before end-of-speech (ms) | Unset (700ms default) |
| `startSensitivity` | `low` / `high` | Speech start detection sensitivity | Unset |
| `endSensitivity` | `low` / `high` | Speech end detection sensitivity | Unset |
| `activityHandling` | `interrupt` / `no-interruption` | Whether speech interrupts model | Unset |
| `turnCoverage` | `only-activity` / `all-input` / `audio-activity-and-all-video` | What counts as a turn | Unset |
| `automaticActivityDetectionDisabled` | boolean | Disable auto VAD | Unset |
| `thinkingLevel` | `minimal` / `low` / `medium` / `high` | Thinking budget level | Unset |
| `thinkingBudget` | number | Token budget for thinking | Unset |
| `apiVersion` | string | API version | Unset (`v1beta` default) |

### Available Google Voices

Default is `Kore`. Other options: `Puck`, `Charon`, `Fenrir`, `Aoede`, `Leda`, `Orus`, `Zephyr`.

To change voice:
```bash
openclaw config set talk.providers.google.voice "Puck"
openclaw gateway restart
```

---

## Talk Config Schema (Top-Level)

| Key | Type | Description | Current |
|---|---|---|---|
| `talk.provider` | string | Active talk provider ID | `google` |
| `talk.providers.<id>.*` | object | Provider-specific settings | See above |
| `talk.speechLocale` | string | BCP 47 locale for speech recognition | Unset (device default) |
| `talk.interruptOnSpeech` | boolean | Stop assistant speech when user talks | Unset (default: true) |
| `talk.silenceTimeoutMs` | integer | Ms of silence before transcript finalized | Unset (700ms macOS/Android, 900ms iOS) |

---

## The `openclaw_agent_consult` Tool

Injected into every Talk session automatically. Schema:

```json
{
  "name": "openclaw_agent_consult",
  "description": "Ask the full OpenClaw agent for deeper reasoning, current information, or tool-backed help before speaking.",
  "parameters": {
    "type": "object",
    "properties": {
      "question": { "type": "string", "description": "The concrete question or task the user asked." },
      "context": { "type": "string", "description": "Optional relevant context or transcript summary." },
      "responseStyle": { "type": "string", "description": "Optional style hint for the spoken answer." }
    },
    "required": ["question"]
  }
}
```

### Google-Specific Behavior

On Google Live, the agent consult tool is marked with `behavior: NON_BLOCKING` — Gemini can continue processing audio while waiting for the consult result. When the result arrives, it is scheduled with `FunctionResponseScheduling.WHEN_IDLE` so Gemini speaks the response at a natural break point.

### Agent Consult Tool Policies

These exist in the codebase for the **voice-call (telephony) plugin**, not for browser Talk. Browser Talk always gives the consult tool full agent access.

| Policy | Behavior | Used By |
|---|---|---|
| `safe-read-only` | Only: `read`, `web_search`, `web_fetch`, `x_search`, `memory_search`, `memory_get` | Voice-call plugin (telephony) |
| `owner` | Full access to all agent tools | Voice-call plugin (telephony) |
| `none` | Disables the consult tool entirely | Voice-call plugin (telephony) |

**For browser Talk mode**, the agent consult dispatches to the full session which has unrestricted access to all tools. No toolPolicy config is needed.

---

## Transport Types

| Transport | Provider | How It Works | Used By Talk? |
|---|---|---|---|
| `json-pcm-websocket` | Google | Browser ↔ Google Live API over WebSocket (PCM16 audio) | **Yes (primary)** |
| `gateway-relay` | Any | Browser → gateway → provider (PCM16 relayed through server) | Fallback |
| `webrtc-sdp` | OpenAI | Direct browser ↔ provider WebRTC | No (Google unsupported) |
| `managed-room` | (future) | Not yet in Control UI | No |

Google Talk sessions use `json-pcm-websocket`: the browser connects directly to `wss://generativelanguage.googleapis.com/ws/...`. The gateway only creates the session (generates an ephemeral token), then the browser handles the WebSocket directly.

If browser session creation fails, it falls back to `gateway-relay` where the gateway proxies all audio through the server.

---

## Session Lifecycle (Google Live)

1. Browser requests `talk.realtime.session` from gateway
2. Gateway calls `createGoogleRealtimeBrowserSession`:
   - Uses `v1alpha` API to create an auth token
   - Token is locked to: model, voice, system instructions, tools, config
   - Token expires in 30 minutes, new session window is 60 seconds
3. Gateway returns session descriptor to browser:
   - `transport: "json-pcm-websocket"`
   - `websocketUrl`, `clientSecret`, `initialMessage`, `audio` format spec
4. Browser validates URL (must be `wss://generativelanguage.googleapis.com`), opens WebSocket
5. Sends initial setup message (model + `responseModalities: [AUDIO]`)
6. Receives `setupComplete` → session is live
7. Audio streams bidirectionally as base64-encoded PCM16
8. Tool calls arrive as `toolCall` messages; results sent back as `sendToolResponse`

---

## Switching Voice Models

`gpt-realtime-1.5` and `gpt-realtime-mini` are both deployed on `onemind-voice.openai.azure.com`.

```bash
# Switch to cheaper fallback (gpt-realtime-mini)
ssh ubuntu@100.123.126.34 'sudo python3 -c "
import json
d=json.load(open(\"/home/ubuntu/.openclaw/openclaw.json\"))
d[\"talk\"][\"providers\"][\"openai\"][\"azureDeployment\"]=\"gpt-realtime-mini\"
d[\"talk\"][\"providers\"][\"openai\"][\"model\"]=\"gpt-realtime-mini\"
json.dump(d, open(\"/home/ubuntu/.openclaw/openclaw.json\",\"w\"), indent=2)
"'
ssh ubuntu@100.123.126.34 'XDG_RUNTIME_DIR=/run/user/1000 systemctl --user restart openclaw-gateway'
```

---

## GCP Infrastructure

| Resource | Detail |
|---|---|
| Project | OneMind OS (`legacy-aio-474306`) |
| Service Account | `openclaw-gemini@legacy-aio-474306.iam.gserviceaccount.com` (Vertex AI User) |
| SA Key | `/tmp/openclaw-gemini-key.json` (on Mac — move to secure location) |
| API Key | Restricted to `generativelanguage.googleapis.com` |
| APIs Enabled | `aiplatform.googleapis.com`, `generativelanguage.googleapis.com`, `apikeys.googleapis.com` |
| Credit | $1,000 GenAI App Builder trial |

---

## Cost Comparison (per minute of voice)

| Provider | Model | Total/min | Status |
|---|---|---|---|
| **Azure OpenAI Realtime** | gpt-realtime-1.5 | **~$0.21–0.40** | ✅ Active — `onemind-voice.openai.azure.com` |
| Azure OpenAI Realtime | gpt-realtime-mini | ~$0.09–0.21 | Available — deployed on `onemind-voice`, change `azureDeployment` |
| OpenAI Direct | gpt-realtime-1.5 | ~$0.21–0.40 | Available (needs direct OpenAI key) |

All voice runs on Azure. Switch to `gpt-realtime-mini` deployment to reduce cost ~2x if needed.

---

## Key Source Files (for debugging)

| File | Location | What It Does |
|---|---|---|
| `server-methods-DSL6KpJ1.js` | `/usr/lib/node_modules/openclaw/dist/` | Talk session handler, `buildRealtimeInstructions()`, `buildTalkRealtimeConfig()` |
| `session-runtime-pKRTzlWC.js` | Same | Agent consult tool schema, bridge session factory, tool policies |
| `realtime-voice-provider.js` | `extensions/google/` | Google Live API bridge, `buildGoogleLiveConnectConfig()`, browser session creation |
| `index-0RNoHvMd.js` | `control-ui/assets/` | Browser UI: transport handlers (GoogleLive, OpenAI WebRTC, relay), agent consult dispatcher |

**Note:** These filenames contain content hashes and will change on OpenClaw updates.

---

## Source Patches — Universal Agent-Aware Voice

**Applied:** 2026-05-01 on both omos-core and omos-sentinel
**Patch script:** `~/.openclaw/patches/patch-universal-voice.sh` (on both servers)
**Codex copy:** `openclaw/patches/patch-universal-voice.sh`

These patches make Talk mode **dynamically agent-aware** — the voice model automatically adopts whatever agent's identity it's connected to, and all references are generic (no hardcoded names).

### What Was Patched

| # | File | What Changed |
|---|---|---|
| 1 | `server-methods-DSL6KpJ1.js` | `buildRealtimeInstructions()` → accepts `agentName` param, uses it in system prompt |
| 2 | `server-methods-DSL6KpJ1.js` | `talk.realtime.session` handler → resolves agent identity from `sessionKey` via `resolveAssistantIdentity()`, passes name to instructions |
| 3 | `session-runtime-pKRTzlWC.js` | Consult tool description → "Query your own deeper capabilities..." (generic) |
| 4 | `session-runtime-pKRTzlWC.js` | Consult prompt → uses dynamic `assistantLabel` instead of hardcoded name |

### How It Works

1. Browser opens Talk with a session key (e.g. `main:session-abc`)
2. `talk.realtime.session` handler extracts agent ID from session key
3. `resolveAssistantIdentity({ cfg, agentId })` returns agent's configured `name` (e.g. "OneMind", "Sentinel AI")
4. `buildRealtimeInstructions(agentName)` generates a system prompt: "You are {name}. This is your realtime voice interface..."
5. The voice model IS that agent — the consult tool is "querying your own deeper capabilities"

### Re-applying After OpenClaw Updates

Source patches are overwritten when OpenClaw is updated (`npm install -g openclaw`). Re-apply:

```bash
# omos-core
ssh ubuntu@100.123.126.34 'sudo bash ~/.openclaw/patches/patch-universal-voice.sh'
openclaw gateway stop && systemctl --user start openclaw-gateway

# omos-sentinel
ssh ubuntu@100.66.235.38 'sudo bash ~/.openclaw/patches/patch-universal-voice.sh'
openclaw gateway stop && systemctl --user start openclaw-gateway
```

**Note:** After an update, filenames may change (they contain content hashes). If the script fails, update the `SM` and `SR` variables in the script to match the new filenames.

### Voice System Prompt (Current)

```
You are {agentName}. This is your realtime voice interface — you ARE this agent
speaking directly. Communicate with military brevity: acknowledge, execute, report.
Use Copy, Roger, Understood, Negative, Stand by. Keep spoken replies short and
decisive. For anything requiring tools, files, memory, infrastructure status, or
deeper reasoning, call agent_consult immediately without explaining what you are
doing. Just say Stand by and call it. When the result comes back, deliver it concisely.
```

---

## Limitations & Future Opportunities

### Current Limitations
- No conversation history carried between Talk sessions
- Agent consult adds latency (~2-10s for tool-backed queries)
- Google Live sessions expire after 30 minutes
- Per-agent voice (Kore/Puck/Aoede) must be selected in the Talk UI; no per-agent config field yet

### Workspace Files (Primary Customization)
The workspace bootstrap files (`IDENTITY.md`, `USER.md`, `SOUL.md`, `TOOLS.md`) at `~/.openclaw/workspace/` are loaded as system context for the agent that Talk consults. **Editing these files is the primary way to customize Talk behavior** — the agent consult returns responses shaped by these files. Changes take effect immediately on next consult (no gateway restart needed).

### Potential Upstream Improvements
- `talk.instructions` or `talk.systemPrompt` config field
- Custom tool injection via `talk.tools` for Talk-specific actions
- Per-agent `talkVoice` config field for automatic voice switching

---

## Quick Reference

### Restart Gateway (after config changes)
```bash
# Core
ssh ubuntu@100.123.126.34 'systemctl --user restart openclaw-gateway'

# Sentinel
ssh ubuntu@100.66.235.38 'systemctl --user restart openclaw-gateway'
```

### Verify Talk Config
```bash
openclaw config get talk
```

### Verify Workspace Files
```bash
ls ~/.openclaw/workspace/
cat ~/.openclaw/workspace/IDENTITY.md
```

### Change Voice
```bash
openclaw config set talk.providers.google.voice "Puck"
openclaw gateway restart
```

### Test Agent Consult (via chat)
Open OpenClaw web UI → start a chat → ask something that requires tools → confirm agent responds with OneMind identity and Zeus context.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
