---
id: openclaw-voice-research
title: "OpenClaw Voice — Deep Research & OneMind Integration Plan"
type: reference
status: active
created: 2026-04-23
owner: zeus
tags: [type/reference, domain/so, project/flutter-app, voice, realtime]
ai_summary: "Comprehensive research on OpenClaw's native voice stack (Talk mode, TTS, Voice Wake), the broader OpenClaw voice ecosystem (VisionClaw, Nerve, community projects), and a concrete plan to build universal real-time voice into the OneMind Flutter app."
---

# OpenClaw Voice — Deep Research & Integration Plan

> Research date: April 23, 2026
> Scope: OpenClaw voice capabilities, ecosystem projects, and building universal real-time voice for OneMind across all devices.

---

## 1. What OpenClaw Ships Natively

### 1.1 Talk Mode — The Core Voice Loop

OpenClaw has a built-in **Talk Mode** — a continuous voice conversation loop that runs on macOS/iOS/Android nodes. It is NOT just TTS bolted on. It's a full bidirectional voice pipeline.

**The loop:**
```
Listen for speech (mic + VAD)
    → Silence detected (700ms macOS/Android, 900ms iOS default)
    → Send transcript to: chat.send (main session)
    → Agent responds
    → Speak via talk.speak (configured provider)
    → Interrupt-aware: if user talks during playback → stop + note timestamp
```

**Config in `~/.openclaw/openclaw.json`:**
```json
{
  "talk": {
    "voiceId": "elevenlabs_voice_id",
    "modelId": "eleven_v3",
    "outputFormat": "mp3_44100_128",
    "apiKey": "elevenlabs_api_key",
    "silenceTimeoutMs": 1500,
    "interruptOnSpeech": true
  }
}
```

**Zeus's current config already has:**
```json
"talk": {
  "interruptOnSpeech": true,
  "provider": "google",
  "providers": {
    "google": {
      "apiKey": "[configured]",
      "modelId": "gemini-2.5-flash-preview-tts",
      "voiceId": "Charon"
    }
  }
}
```
**This is already active.** The gateway is configured for real-time voice with Google Gemini TTS.

**Voice directives** (model can control its own voice per-reply):
```
{ "voice": "pMsXgVXv3BLzUgSXRplE", "once": true }
[[tts:speed=1.2 voiceId=Charon]]
[[tts:text]](laughs) Here's what I found.[[/tts:text]]
```

**RPC interface:**
- `talk.speak` — play audio via the gateway's talk provider
- `chat.send` — send transcript as text
- Talk mode replies go to WebChat (same as typing)

### 1.2 TTS Providers — Full Comparison

| Provider | Model | Best For | Key Strength |
|----------|-------|----------|-------------|
| **Google Gemini** | `gemini-2.5-flash-preview-tts` | Real-time, low cost | Fast, we already have key |
| **ElevenLabs** | `eleven_v3` / `eleven_multilingual_v2` | Richest voices, expressive | Best quality, voice cloning |
| **OpenAI** | `gpt-4o-mini-tts` | Reliable, fast | Widely supported |
| **Microsoft** | Edge Neural TTS | Free, no API key | Good for fallback |
| **MiniMax** | `speech-2.8-hd` | High quality, affordable | 32kHz |
| **xAI** | `eve`, `ara`, `leo`, etc. | Expressive | 6 voices |

**Output formats by channel:**
- Telegram / WhatsApp / Matrix: **Opus voice message** (48kHz/64kbps)
- Other channels: **MP3** (44.1kHz/128kbps)
- Google Gemini: returns raw 24kHz PCM → wrapped as WAV
- Android: `pcm_16000`, `pcm_22050`, `pcm_24000`, `pcm_44100` for low-latency AudioTrack

**Gateway RPC:**
```
tts.status  /  tts.enable  /  tts.disable
tts.convert  /  tts.setProvider  /  tts.providers
```

### 1.3 Voice Wake — Always-On Trigger

OpenClaw has a **global Voice Wake** system — shared list across all nodes.

```json
// ~/.openclaw/settings/voicewake.json
{ "triggers": ["openclaw", "legacy", "jarvis"], "updatedAtMs": 1730000000000 }
```

**RPC:**
- `voicewake.get` → current triggers
- `voicewake.set` → update globally (broadcasts to all nodes)
- `voicewake.changed` event → pushed to all WebSocket clients on change

**Platform support:**
- **macOS**: fully supported, menu bar toggle
- **iOS node**: fully supported, `VoiceWakeManager`
- **Android node**: disabled — uses manual mic in Voice tab instead

### 1.4 TTS Auto-Modes

```
auto: "off"       — no automatic TTS
auto: "always"    — TTS every reply
auto: "inbound"   — TTS only after a voice message comes in
auto: "tagged"    — TTS only when model emits [[tts:...]] directive
```

**Flow:**
```
Reply → TTS enabled?
  no → send text
  yes → has media / short reply?
         yes → send text
         no → length > limit?
                no → TTS → attach audio
                yes → summary enabled?
                         yes → summarize → TTS → attach audio
                         no → send text
```

---

## 2. The OpenClaw Voice Ecosystem (Community Projects)

Six major community projects found via GitHub research. All relevant to OneMind.

### 2.1 VisionClaw — The Gold Standard Reference
**Repo:** `Intent-Lab/VisionClaw` — 2,200+ stars, active, 394 forks

This is the most important project for us. It's exactly what we want to build, already proven.

**What it does:**
- Real-time voice + vision AI via **Gemini Live API** (WebSocket, native audio — NOT STT-first)
- Routes tool calls through **OpenClaw gateway** (56+ skills)
- Works on iPhone (glasses POV) OR just phone camera
- WebRTC live streaming to browser viewer

**Architecture:**
```
Meta Ray-Ban / Phone Camera
       ↓
iOS / Android App
       ↓ JPEG frames (~1fps) + PCM audio (16kHz)
Gemini Live API (WebSocket)
       ↓ Audio response (PCM 24kHz) → Speaker
       ↓ Tool calls → App → OpenClaw Gateway → 56+ skills
       ↓ Tool response → Gemini speaks result
```

**Audio Pipeline (proven, production):**
- Input: PCM Int16, 16kHz mono, 100ms chunks → Gemini WebSocket
- Output: Gemini WebSocket → AudioManager playback queue → Speaker
- iOS: `.voiceChat` audio session (echo cancellation + mic gating during AI speech)
- Android: `VOICE_COMMUNICATION` audio source (built-in AEC)

**Key source files to port:**
| File | Language | What to extract |
|------|----------|-----------------|
| `GeminiLiveService.swift` | Swift | WebSocket client for Gemini Live API |
| `AudioManager.swift` | Swift | PCM capture + playback pipeline |
| `GeminiLiveService.kt` | Kotlin | Same but Android (OkHttp WebSocket) |
| `AudioManager.kt` | Kotlin | AudioRecord(16kHz) + AudioTrack(24kHz) |
| `OpenClawBridge.swift/.kt` | Swift/Kotlin | HTTP client to OpenClaw |
| `ToolCallRouter.swift/.kt` | Swift/Kotlin | Routes Gemini tool calls to OpenClaw |

**Flutter porting strategy:**
- VisionClaw uses Gemini Live's native audio (bidirectional PCM over WebSocket)
- For our app: we can use either Gemini Live (ultra-low latency, bidirectional) OR OpenClaw's `talk.speak` (simpler, uses existing config)
- VisionClaw approach gives us: voice + vision + tool calling in one pipeline

### 2.2 openclaw-nerve — The Web Cockpit
**Repo:** `daggerhashimoto/openclaw-nerve` — 752 stars, actively maintained (daily commits), v1.5.3

The best open-source web UI for OpenClaw. Directly relevant to our setup.

**Architecture:**
```
Browser ── Nerve (:3080) ── OpenClaw Gateway (:18789)
  ├─ WS ───────────────── proxied to gateway
  ├─ SSE ──────────────── file watchers, real-time sync
  └─ REST ─────────────── files, memories, TTS, models
```

**Voice capabilities built in:**
- Push-to-talk
- Wake word flows
- Explicit language selection
- **Local Whisper transcription** (offline STT)
- Multilingual stop/cancel phrases
- Multiple TTS providers

**Capability snapshot:**
- Agent fleet control (multiple agents, subagents)
- Session trees, kanban task board, review flows
- Per-agent workspace editor
- Token usage + cost tracking
- 14 themes, responsive mobile UI
- Cron scheduling

**Why this matters for us:**
- We could run Nerve on VPS and access it as a web app from Safari/Chrome on iPhone
- It already has Tailscale integration (local, hybrid, cloud deploy guides)
- It proxies the gateway, so it speaks the same WebSocket protocol our Flutter app uses
- Could serve as our "web fallback" or desktop power-user interface

**Deploy on VPS (our setup):**
```bash
curl -fsSL https://raw.githubusercontent.com/daggerhashimoto/openclaw-nerve/master/install.sh | bash
# Then configure for Tailscale Serve (same as our gateway setup)
```

### 2.3 openclaw-assistant (Android)
**Repo:** `yuga-hashimoto/openclaw-assistant` — 250 stars, Kotlin

- Wake word activation
- System assistant integration
- TTS responses
- Runs as a system-level Android service

**Interesting for us:** Shows how to make OpenClaw the system voice assistant on Android — replaces Google Assistant.

### 2.4 openclaw-assistant-mvp (Desktop)
**Repo:** `andyhuo520/openclaw-assistant-mvp` — 410 stars, Electron/JavaScript

- Live2D character animations
- Real-time speech recognition
- TTS responses
- Desktop UI companion

**Interesting for us:** The Live2D companion model idea could be a future addition to OneMind app.

### 2.5 openclaw-command-center
**Repo:** `mayukh4/openclaw-command-center` — 64 stars, JavaScript

- 24/7 interactive (voice + touch) mission control
- Live mission control dashboard for agents

### 2.6 openclaw-homeassistant
**Repo:** `ddrayne/openclaw-homeassistant` — 45 stars, Python

- OpenClaw as Home Assistant voice provider
- Replaces Assist pipeline in HA with OpenClaw agents

---

## 3. Real-Time Voice — Technology Options

Three architectural approaches. Each with different tradeoffs.

### 3.1 Approach A: OpenClaw Native Talk Mode (Simplest)

**Use the gateway's `talk.speak` RPC + `chat.send`**

```
Flutter app
  ↓ mic → on-device STT (or Whisper) → text
  ↓ chat.send via WebSocket
OpenClaw Gateway
  ↓ agent processes, responds
  ↓ talk.speak → TTS audio bytes
Flutter app
  ↓ play audio
```

**Pros:** Already configured (Zeus has Google Gemini TTS active), uses existing gateway stack, works with all 9 agents, tool calls happen server-side
**Cons:** Not truly "real-time" — full round trip (STT → process → TTS), 1-3s latency typical
**Best for:** Standard voice conversation mode in the app (good enough for 95% of use cases)

**Flutter implementation:**
```dart
// 1. Capture audio with `record` package
final recorder = AudioRecorder();
await recorder.start(const RecordConfig(encoder: AudioEncoder.pcm16bits, sampleRate: 16000));

// 2. When silence detected, stop and get transcription
// Option A: On-device STT (speech_to_text package)
// Option B: Send PCM to gateway's transcription endpoint

// 3. Send transcript via WebSocket
ws.sink.add(jsonEncode({
  'jsonrpc': '2.0',
  'id': uuid.v4(),
  'method': 'chat.send',
  'params': {'sessionKey': 'agent:main:flutter', 'message': transcript},
}));

// 4. Listen for agent response + talk.speak event
// Gateway sends audio bytes or URL back

// 5. Play audio with just_audio
final player = AudioPlayer();
await player.setAudioSource(AudioSource.uri(Uri.parse(audioUrl)));
await player.play();
```

### 3.2 Approach B: Gemini Live (VisionClaw Pattern — Ultra Low Latency)

**Direct Gemini Live WebSocket with OpenClaw as tool executor**

```
Flutter app
  ↓ PCM 16kHz raw audio stream → Gemini Live WebSocket
Gemini Live API
  ↓ PCM 24kHz audio response → Flutter → speaker
  ↓ Tool calls → Flutter → OpenClaw gateway → 56+ skills
  ↓ Tool response → Gemini continues speaking
```

**Pros:** Lowest latency (bidirectional native audio, no STT step), voice + vision possible, Gemini can see phone camera, interrupt support baked in
**Cons:** Requires Gemini API directly (already have key in config), tool calls routed through extra hop, more complex WebSocket binary frame handling
**Best for:** "Talk to Legacy" native voice mode, future camera integration

**Flutter packages needed:**
- `web_socket_channel` — for Gemini Live WebSocket binary frames
- `record` — PCM 16kHz capture
- `just_audio` / `audioplayers` — PCM 24kHz playback

```dart
// Gemini Live WebSocket connection
final wsUri = Uri.parse(
  'wss://generativelanguage.googleapis.com/ws/google.ai.generativelanguage.v1alpha.GenerativeService.BidiGenerateContent'
  '?key=$geminiApiKey'
);

// Send setup message
ws.sink.add(jsonEncode({
  'setup': {
    'model': 'models/gemini-2.0-flash-live-001',
    'generation_config': {
      'response_modalities': ['AUDIO'],
    },
    'tools': [{'function_declarations': openClawToolDeclarations}],
  }
}));

// Stream PCM audio
// Handle binary AudioData frames from Gemini
// Route tool_call frames to OpenClaw gateway
```

### 3.3 Approach C: Hybrid — OpenClaw Talk via Flutter WebSocket

**Flutter triggers Talk mode on the gateway, gateway handles the audio**

This is how the macOS app and iOS node work with Talk mode. The Flutter app could implement the same `talk.speak` subscription:

```
Flutter app
  ↓ `voicewake.get` on connect → get trigger words
  ↓ mic → on-device VAD → silence → `chat.send`
OpenClaw Gateway
  ↓ agent responds with text
  ↓ gateway calls `talk.speak` → audio bytes returned as binary WS frame
Flutter app
  ↓ receives binary audio frame → play via AudioTrack/AVAudioEngine
```

**Key insight from docs:** Android Talk mode returns PCM at `pcm_24000` by default for low-latency AudioTrack streaming. Flutter can match this.

---

## 4. Flutter Integration Blueprint

### 4.1 Package Stack

```yaml
# Voice-specific additions needed in pubspec.yaml
dependencies:
  # Audio I/O
  record: ^5.1.0              # PCM 16kHz mic capture (already in FLUTTER-APP-DEV.md)
  just_audio: ^0.9.36         # Playback (WAV, MP3, PCM streaming)
  flutter_pcm_sound: ^0.4.0   # Real-time PCM playback (AudioTrack/AudioQueue)

  # STT (on-device fallback when not using Gemini Live)
  # NOTE: Previously removed due to iOS 26 JIT crash — re-add when ready
  # speech_to_text: ^7.0.0   # REMOVED — re-add after iOS 26 stability confirmed

  # WebRTC (future — for LiveKit/peer-to-peer voice)
  # flutter_webrtc: ^0.9.47  # Full WebRTC (STUN/TURN) for VisionClaw-style live streaming

  # WebSocket (already have web_socket_channel)
  # Binary frame support is native to web_socket_channel
```

### 4.2 Voice Service Architecture

```dart
// lib/core/voice/
//
// voice_service.dart         — coordinator (Talk mode state machine)
// mic_recorder.dart          — PCM capture from mic
// vad_detector.dart          — Voice Activity Detection (silence threshold)
// tts_player.dart            — audio playback (PCM/MP3/WAV)
// gemini_live_client.dart    — Gemini Live WebSocket (Approach B)
// openclaw_talk_client.dart  — OpenClaw Talk RPC (Approach A/C)

enum VoiceMode { off, listening, thinking, speaking }

class VoiceService {
  // State
  VoiceMode mode = VoiceMode.off;
  bool interruptOnSpeech = true;

  // Providers
  final MicRecorder _mic;
  final VadDetector _vad;
  final TtsPlayer _player;
  final GatewayClient _gateway;

  Future<void> startTalkMode() async {
    // 1. Get wake words from gateway
    final triggers = await _gateway.requestMethod('voicewake.get', {});

    // 2. Start listening
    mode = VoiceMode.listening;
    await _mic.startCapture(sampleRate: 16000);

    // 3. VAD loop
    _vad.onSilence = (transcript) async {
      mode = VoiceMode.thinking;
      await _stopMic();

      // 4. Send to gateway
      await _gateway.sendMessage(sessionKey: 'agent:main:flutter', text: transcript);
    };

    // 5. Listen for gateway TTS response
    _gateway.onTalkSpeak = (audioBytes) async {
      mode = VoiceMode.speaking;
      await _player.playPcm(audioBytes);

      // After playback, back to listening
      mode = VoiceMode.listening;
      await _mic.startCapture(sampleRate: 16000);
    };
  }

  Future<void> interrupt() async {
    if (mode == VoiceMode.speaking) {
      await _player.stop(); // Stop current playback
      mode = VoiceMode.listening;
      // Gateway notes interrupt via interruptOnSpeech config
    }
  }
}
```

### 4.3 UI — Voice Mode Screen

The `voice_input_sheet.dart` is currently a stub. Here's the real implementation blueprint:

```
┌─────────────────────────────────┐
│         VOICE MODE              │
│                                 │
│    ┌─────────────────────┐      │
│    │    ◉  LEGACY  ⚡    │      │  ← Active agent
│    └─────────────────────┘      │
│                                 │
│         ╭─────────╮             │
│         │  🎙️     │  ← Animated │  Listening: pulse
│         │  ~~~~   │    phase    │  Thinking: spinner
│         │  ····   │    clouds   │  Speaking: rings
│         ╰─────────╯             │
│                                 │
│   "What's on my schedule?"      │  ← Transcript preview
│                                 │
│   [──────────────────────────]  │  ← Audio level bar
│                                 │
│   Spartan  Oracle  Eden  More   │  ← Agent switcher
│                                 │
│        [ END VOICE MODE ]       │
└─────────────────────────────────┘
```

**Interaction states:**
- **Listening**: cloud pulses with mic level (matches macOS app behavior)
- **Thinking**: sinking/spinner animation
- **Speaking**: radiating rings
- **Tap animation**: interrupt playback
- **Swipe up**: expand to full chat view
- **Long press agent**: switch active agent for voice session

### 4.4 Voice Wake Integration

Once `speech_to_text` is stable on iOS 26 or we use Gemini Live (no need for STT package):

```dart
// Subscribe to wake word changes on connect
void _subscribeVoiceWake() {
  gateway.on('voicewake.changed', (data) {
    _wakeWords = List<String>.from(data['triggers']);
  });
  gateway.requestMethod('voicewake.get', {}).then((result) {
    _wakeWords = List<String>.from(result['triggers']);
  });
}

// Custom trigger words from settings: "Legacy", "Hey Legacy", "OneMind"
// On match → startTalkMode()
```

---

## 5. Universal Voice — All Devices

This is the key strategic point. The goal: **one voice backend (OpenClaw gateway on VPS), multiple frontends**.

### 5.1 Device Matrix

| Device | Voice Method | Status |
|--------|-------------|--------|
| **iPhone** (OneMind app) | Flutter Talk mode + Gemini Live | Build this |
| **Android** | Flutter Talk mode + manual mic | Build this |
| **Mac** | OpenClaw native macOS app | Already works |
| **Telegram** | Voice messages + auto-TTS | Already works |
| **Browser (Safari/Chrome)** | Nerve web UI with PTT | Install Nerve on VPS |
| **Smart Speaker** (future) | OpenClaw Raspberry Pi node | Future |
| **Meta Ray-Ban** | VisionClaw app | Exists, can use on iPhone |
| **Apple Watch** | Watchkit → WatchConnectivity → iPhone → gateway | Future |

### 5.2 The Universal Pattern

All these surface the same gateway: `wss://omos-core.tail717ba0.ts.net`

All use the same session: `agent:main:flutter` (or `agent:main:mobile`)

All use the same agent: `main` (Legacy) with specialist routing

The shared config (`talk:`, `messages.tts:`) applies everywhere.

**Changing voice globally:**
```bash
ssh root@omos-core.tail717ba0.ts.net 'openclaw config set talk.voiceId "new_voice_id"'
# → applies to ALL clients simultaneously
```

### 5.3 OneMind Voice Identity

With ElevenLabs, we can clone a custom voice for Legacy:
- Record 1-5 minutes of Zeus's own voice or a chosen style
- Create a custom voice ID in ElevenLabs dashboard
- Set as default: `openclaw config set talk.voiceId "your_clone_id"`
- Every device speaks with the same identity

### 5.4 Multi-Agent Voice Routing

Voice becomes a routing surface — different agents have different voices:

```json
{
  "talk": {
    "voiceId": "Charon",  // Legacy default (ominous, deep)
    "provider": "google"
  }
}
```

Future: per-agent voice profiles via voice directives in system prompts:
```
// In spartan's IDENTITY.md workspace file:
// Start each response with: {"voice": "energetic-male-voice-id", "once": true}
```

---

## 6. Implementation Plan (Phased)

### Phase 1: Basic TTS in Flutter App (1-2 days)

**What:** When Legacy responds in chat, auto-play TTS via gateway.

**How:**
1. Subscribe to `chat.delta` events — collect full response
2. On `chat.finish`, call `tts.convert` RPC with the text
3. Gateway returns audio URL or binary bytes
4. Play via `just_audio`

**Why first:** Lowest complexity, zero new plumbing. Uses existing `talk:` config (Google Gemini TTS already set). Validates the pipeline before adding mic.

**Flutter changes:**
- Add `just_audio` to pubspec
- Add TTS toggle in chat input widget
- Add `tts.convert` call to `ChatNotifier._onFinish()`

### Phase 2: Push-to-Talk Voice Input (2-3 days)

**What:** Hold mic button → speak → release → text appears in chat input → send.

**How:**
1. Add `record` package back (was removed — iOS 26 stable enough now in release mode)
2. `MicRecorder` captures PCM until button release
3. For STT: use `speech_to_text` for on-device transcription OR send PCM to Whisper via gateway skill
4. Result populates chat text input

**Why second:** Meets the "voice input" UX without full real-time loop complexity.

### Phase 3: Full Talk Mode (3-5 days)

**What:** Tap voice button → continuous conversation loop (Listening → Thinking → Speaking) matching OpenClaw macOS Talk mode behavior.

**How:**
1. `VoiceService` state machine
2. VAD (Voice Activity Detection) via amplitude monitoring
3. Silence threshold (700ms) triggers send
4. Listen for binary audio frames on WebSocket (`talk.speak` event)
5. Play PCM 24kHz via `flutter_pcm_sound` (low-latency AudioTrack)
6. Interrupt: tap during Speaking → stop playback, back to Listening

**Gateway side:** Already configured. `interruptOnSpeech: true` is live.

### Phase 4: Gemini Live Mode (5-7 days)

**What:** Upgrade Talk Mode to use Gemini Live directly — native bidirectional audio, sub-second latency, optional camera feed.

**How:**
1. `GeminiLiveClient` — WebSocket to Gemini Live API
2. Declare `execute` tool that routes to OpenClaw (VisionClaw pattern)
3. Stream PCM 16kHz mic → Gemini Live
4. Receive PCM 24kHz from Gemini → play
5. Route tool calls → `GatewayClient.toolInvoke(task: ...)` → gateway handles

**This is the VisionClaw architecture, built in Flutter.**

### Phase 5: Install Nerve on VPS (1 day)

**What:** Web browser voice UI accessible via Tailscale from any device.

```bash
ssh root@omos-core.tail717ba0.ts.net 'curl -fsSL https://raw.githubusercontent.com/daggerhashimoto/openclaw-nerve/master/install.sh | bash'
# Configure for Tailscale Serve on port 3080
# Access: https://omos-core.tail717ba0.ts.net:3080/
```

**Why:** Instant voice access from any browser (Mac, iPad, borrowed device) without the Flutter app. Also gives us the Kanban + agent fleet UI.

### Phase 6: Voice Wake on iPhone (3-5 days)

**What:** Say "Legacy" → iPhone activates Talk mode automatically.

**How:**
1. Re-add `speech_to_text` or use `SFSpeechRecognizer` via platform channel
2. Listen for wake word in background
3. On trigger → start Talk Mode

**iOS 26 caveat:** Background audio processing requires careful permission handling. May need Siri Shortcuts bridge as intermediate.

---

## 7. Quick Wins — Deploy Today

These require no app code changes — just gateway config.

### 7.1 Enable TTS Auto-Mode in Telegram

Every agent reply in Telegram becomes a voice message:
```bash
ssh root@omos-core.tail717ba0.ts.net 'openclaw config set messages.tts.auto always'
```

### 7.2 Custom Wake Words

```bash
ssh root@omos-core.tail717ba0.ts.net 'openclaw config set voicewake.triggers ["legacy","hey legacy","one mind"]'
```
→ MacOS app and iOS node now respond to "Legacy" and "Hey Legacy"

### 7.3 Switch to ElevenLabs Voice (Zeus's custom voice)

1. Go to elevenlabs.io — create voice or use pre-built
2. Get voice ID and API key
3. Apply:
```bash
openclaw config set talk.provider elevenlabs
openclaw config set talk.voiceId "YOUR_ELEVENLABS_VOICE_ID"
openclaw config set talk.apiKey "YOUR_ELEVENLABS_API_KEY"
openclaw config set talk.modelId "eleven_v3"
```

### 7.4 Install Nerve (Voice Web UI)

```bash
ssh root@omos-core.tail717ba0.ts.net
curl -fsSL https://raw.githubusercontent.com/daggerhashimoto/openclaw-nerve/master/install.sh | bash
```
Access Nerve's push-to-talk voice UI from any browser on Tailscale.

---

## 8. Key Repos & References

| Project | URL | Stars | Why Important |
|---------|-----|-------|--------------|
| VisionClaw | `github.com/Intent-Lab/VisionClaw` | 2.2k | Audio pipeline to extract, Gemini Live pattern |
| openclaw-nerve | `github.com/daggerhashimoto/openclaw-nerve` | 752 | Install on VPS for web voice UI |
| openclaw-assistant | `github.com/yuga-hashimoto/openclaw-assistant` | 250 | Android wake word system service |
| openclaw-assistant-mvp | `github.com/andyhuo520/openclaw-assistant-mvp` | 410 | Desktop voice + Live2D |
| openclaw-homeassistant | `github.com/ddrayne/openclaw-homeassistant` | 45 | HA voice provider integration |
| OpenClaw Talk docs | `docs.openclaw.ai/nodes/talk` | — | Official Talk mode reference |
| OpenClaw TTS docs | `docs.openclaw.ai/tools/tts` | — | All TTS providers + config |
| OpenClaw Voice Wake | `docs.openclaw.ai/nodes/voicewake` | — | Wake word protocol |

---

## 9. Current Config State

Zeus's local Mac `~/.openclaw/openclaw.json` already has:
```json
"talk": {
  "interruptOnSpeech": true,
  "provider": "google",
  "providers": {
    "google": {
      "apiKey": "[configured]",
      "modelId": "gemini-2.5-flash-preview-tts",
      "voiceId": "Charon"
    }
  }
}
```

The **VPS gateway** (`omos-core`) also needs this config synced. Current VPS config may not have `talk:` configured. Verify:
```bash
ssh root@omos-core.tail717ba0.ts.net 'cat ~/.openclaw/openclaw.json | python3 -c "import sys,json; c=json.load(sys.stdin); print(json.dumps(c.get(\"talk\",{}),indent=2))"'
```

---

## 10. The Big Picture

```
OneMind Voice Schema (2026+)

                    OpenClaw VPS (omos-core)
                    ┌──────────────────────┐
                    │  9 Agents (Legacy,   │
                    │  Spartan, Oracle...) │
                    │                      │
                    │  TTS: Google Gemini  │
                    │  Voice: Charon       │
                    │  interruptOnSpeech   │
                    └──────────┬───────────┘
                               │ Tailscale mesh
              ┌────────────────┼────────────────────┐
              │                │                    │
    ┌─────────▼────┐  ┌────────▼──────┐  ┌─────────▼────┐
    │ OneMind App  │  │  Nerve (web)  │  │  Mac App     │
    │ (Flutter)    │  │  Browser UI   │  │  (native)    │
    │              │  │               │  │              │
    │ Talk Mode    │  │ Push-to-Talk  │  │ Talk Mode    │
    │ Phase 3+     │  │ Wake Word     │  │ Active now   │
    │              │  │ Whisper STT   │  │              │
    └──────────────┘  └───────────────┘  └──────────────┘
         iOS/Android       Any browser        macOS
```

**The gateway is the voice brain. Every surface speaks with the same Legacy voice.**
