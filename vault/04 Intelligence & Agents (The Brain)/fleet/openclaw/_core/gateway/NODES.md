# NODES.md — Client Device Setup

Instructions for connecting client devices (nodes) to the OneMind OS gateway. Each connected device becomes an OpenClaw **node** — a trusted endpoint that can chat with Legacy, run tools, and receive canvases.

---

## Gateway Connection Info

| Field | Value |
|-------|-------|
| VPS Hostname | `onemindos-cloud-1` |
| Gateway URL | `https://omos-core.tail717ba0.ts.net` |
| WebSocket URL | `wss://omos-core.tail717ba0.ts.net` |
| Tailscale IP | `100.98.25.10` |
| Public IP | `64.227.8.23` |
| Gateway Port | `18789` (loopback only — Tailscale Serve exposes via HTTPS) |
| Gateway Token | Stored in 1Password: `00-24 SO / OpenClaw Gateway / gateway-token` |
| Gateway Version | OpenClaw 2026.4.14 |
| QR Pairing | From any paired device: run `/pair qr` in chat |

> **Access is Tailscale-only.** The gateway binds to loopback. Tailscale Serve proxies HTTPS → localhost:18789.
> Use the Tailscale URL for all connections.

### Paired Nodes (April 2026)

| Node | Client | Tailscale IP | Caps | Status |
|------|--------|-------------|------|--------|
| Zeus's MacBook Pro (app) | `openclaw-macos/node` | `100.114.142.55` | browser, camera, canvas, screen | paired · connected |
| Zeus's MacBook Pro (headless) | `node-host/node` | `100.114.142.55` | browser, system | paired · connected |

**Two nodes, one Mac.** The macOS app registers as a node with full GUI capabilities (camera, canvas, screen recording). The headless node host (LaunchAgent) registers separately with system/browser capabilities. Both run simultaneously.

---

## macOS — Official OpenClaw App

The official OpenClaw macOS app auto-updates via Sparkle. No developer account needed.

### Install

```bash
# Option A: Homebrew (recommended)
brew install --cask openclaw

# Option B: Download DMG from GitHub releases
# https://github.com/openclaw/openclaw/releases → download the .dmg asset
```

### Current Mac Config (April 2026)

**Transport: Direct (ws/wss)** — no SSH tunnel. The app connects straight to the gateway via Tailscale Serve HTTPS.

Key config in `~/.openclaw/openclaw.json`:
```json
{
  "gateway": {
    "mode": "remote",
    "remote": {
      "url": "wss://omos-core.tail717ba0.ts.net",
      "transport": "direct",
      "token": "<gateway-token>"
    }
  }
}
```

NSUserDefaults (`ai.openclaw.mac.debug`):
- `openclaw.connectionMode` = `remote`
- `openclaw.remoteTransport` = `direct`
- `openclaw.remoteUrl` = `wss://omos-core.tail717ba0.ts.net`
- `openclaw.remoteTarget` = `root@omos-core.tail717ba0.ts.net`

### Headless Node Host (LaunchAgent)

The Mac also runs a headless node service that stays connected to the gateway even when the app is closed.

```bash
# Install (one-time)
openclaw node install --host omos-core.tail717ba0.ts.net --port 443 --tls

# Status
openclaw node status

# Restart
openclaw node restart
```

LaunchAgent plist: `~/Library/LaunchAgents/ai.openclaw.node.plist`
Logs: `~/.openclaw/logs/node.log` and `node.err.log`

### Connect to Remote Gateway

1. Open **OpenClaw.app**
2. On first launch, the onboarding wizard starts
3. Choose **Remote Gateway** mode
4. Enter the gateway URL: `https://onemindos-cloud.tail717ba0.ts.net`
5. Enter the gateway token (from 1Password)
6. The app will verify the connection and pair as a trusted device

### Update

Updates are automatic via Sparkle. When a new version is available, the app will prompt you to install it. You can also update manually:

```bash
brew upgrade openclaw
```

### Verify

```bash
# Check installed version
openclaw --version

# Check gateway connection status
openclaw gateway status
```

---

## iPhone — Chowder (Sideloaded)

Chowder is our custom iOS client, forked to [`One-Mind-OS/chowder-iOS`](https://github.com/One-Mind-OS/chowder-iOS). It implements full OpenClaw Protocol v3 with streaming, Live Activity, and thinking step visualization.

### Prerequisites

- Mac with Xcode installed (latest stable)
- iPhone connected via USB (or same Wi-Fi for wireless deploy)
- Free Apple ID signed into Xcode (Settings → Accounts)
- USB cable for first deploy

### Initial Setup

```bash
# 1. Clone the fork
git clone https://github.com/One-Mind-OS/chowder-iOS.git
cd chowder-iOS

# 2. Open in Xcode
open Chowder.xcodeproj
```

3. In Xcode, select the **Chowder** target
4. Under **Signing & Capabilities**:
   - Check "Automatically manage signing"
   - Team: Select your personal Apple ID
   - Bundle Identifier: Change to something unique, e.g. `com.zeuslegacy.chowder`
5. Select your iPhone as the build destination (top toolbar)
6. Press **Cmd+R** to build and run

### First Launch on iPhone

1. If prompted "Untrusted Developer":
   - Go to **Settings → General → VPN & Device Management**
   - Tap your Apple ID under "Developer App"
   - Tap **Trust**
2. Open Chowder
3. Enter gateway URL: `https://gateway.onemindos.dev`
4. Enter gateway token
5. You're connected

### Sideload Signing Limits (Free Apple ID)

| Limitation | Free Apple ID | Paid Developer ($99/yr) |
|------------|---------------|------------------------|
| App validity | **7 days** then must re-sign | 1 year |
| Max sideloaded apps | 3 | 10 |
| TestFlight | No | Yes |
| App Store | No | Yes |

### Keeping the App Alive (Free Apple ID)

With a free Apple ID, the app expires every 7 days. Two options:

**Option A: Re-deploy from Xcode (manual)**

Every 7 days, connect iPhone to Mac and press Cmd+R in Xcode. Takes ~30 seconds.

**Option B: SideStore (automated, recommended)**

[SideStore](https://sidestore.io) auto-refreshes sideloaded apps so they never expire.

```bash
# 1. Install SideStore on iPhone via https://sidestore.io
# 2. Export Chowder as .ipa from Xcode:
#    Product → Archive → Distribute App → Custom → Development → Next → Export
# 3. AirDrop or transfer the .ipa to iPhone
# 4. Open in SideStore → Install
# SideStore will auto-refresh the signing before it expires
```

### Updating Chowder

```bash
cd chowder-iOS
git pull origin main    # Pull latest changes from our fork
# Open Xcode, Cmd+R to deploy updated build to iPhone
```

To pull upstream improvements from the original Chowder:

```bash
git remote add upstream https://github.com/newmaterialco/chowder-iOS.git
git fetch upstream
git merge upstream/main
# Resolve any conflicts, then Cmd+R to deploy
```

---

## iPad — Same as iPhone

Chowder works on iPad with the same sideload process. In Xcode, select your iPad as the build destination instead of iPhone. Same signing limits apply.

---

## Apple Watch — Future

The original OneMind iOS app had a WatchApp companion. Key modules to port into Chowder:

- `WatchApp/` — watchOS companion with chat + quick actions
- Watch complication for agent status

**Status**: Not yet implemented in Chowder fork. Tracked for future development.

---

## Android — Future

OpenClaw has an official Android app built from the monorepo, but it's not publicly released yet (per PR #23051). When it becomes available:

- APK sideload via GitHub releases
- Or build from `openclaw/openclaw` monorepo (Android source is in the main repo)

**Status**: Waiting for official public release.

---

## Android TV — Ambient Dashboard

Android TV can join the Tailscale mesh and display the OpenClaw Control UI or Home Assistant dashboards.

### Setup

1. Install **Tailscale** from the Google Play Store on the Android TV
2. Sign into the same tailnet (same Google/Tailscale account as all other devices)
3. The TV gets a Tailscale IP and can access any mesh service
4. Open the TV's web browser → navigate to `https://omos-core.tail717ba0.ts.net`
5. Bookmark it for quick access — use kiosk/fullscreen mode for a persistent dashboard

**Android TV is NOT a node.** It's a display-only device — no OpenClaw agent runs on it. It shows dashboards via browser.

**Use cases:**
- OpenClaw Control UI (agent status, conversations, system health)
- Home Assistant dashboard (weather, cameras, sensors, calendar)
- Chromecast target (cast from Mac or phone)

---

## Home Assistant Integration

Two integration directions — both solved with existing tools.

### Direction 1: OpenClaw → HA (agent controls devices)

`homeassistant-skill` on the VPS. Agents call HA REST API directly via Tailscale.

```bash
openclaw skills install homeassistant-skill
```

**Requirements:** Beelink on Tailscale, HA long-lived token (HA → Profile → Security → Long-Lived Access Tokens)

**Flow:** Legacy (VPS) → `homeassistant-skill` → HA REST API → physical devices

---

### Direction 2: Voice/Chat → OpenClaw (HACS Integration)

**GitHub:** `https://github.com/techartdev/OpenClawHomeAssistantIntegration`
**Version:** v0.1.62 (actively maintained, 70 stars)

Install on Home Assistant — connects HA to your existing VPS gateway via the `/v1/chat/completions` endpoint.

**What it adds:**

| Feature | Detail |
|---------|--------|
| Assist conversation agent | OpenClaw selectable in Settings → Voice Assistants |
| Lovelace chat card | `custom:openclaw-chat-card` — history, typing indicator, voice button |
| Service: `openclaw.send_message` | Automations can send messages to Legacy |
| Service: `openclaw.clear_history` | Clear session history |
| Service: `openclaw.invoke_tool` | Invoke any gateway tool |
| Event: `openclaw_message_received` | Trigger automations on Legacy responses |
| Event: `openclaw_tool_invoked` | Trigger on tool execution |
| Tool calls (optional) | Enables `execute_service` — Legacy controls HA devices directly |
| Voice input | Browser Web Speech or HA STT, wake word ("hey openclaw") |
| Model selector | Dashboard entity to switch models |

**Voice pipeline with this integration:**
```
Voice satellite → Wake word → HA STT → HACS conversation agent
  → POST /v1/chat/completions → Legacy (VPS) → response
  → TTS → speaker
```

**VPS prerequisite (already done — 2026-04-16):**
```bash
openclaw config set gateway.http.endpoints.chatCompletions.enabled true
# Verified: curl POST /v1/chat/completions → {"model":"openclaw/legacy",...} ✅
```

**Install steps (when Beelink HA is running):**
1. HACS → ⋮ → Custom repositories → `https://github.com/techartdev/OpenClawHomeAssistantIntegration` → Integration → Add
2. Install `OpenClaw` → Restart HA
3. Settings → Devices & Services → Add Integration → OpenClaw
4. Manual config (remote/standalone):
   - Gateway Host: `100.98.25.10` (VPS Tailscale IP)
   - Gateway Port: `18789`
   - Gateway Token: from 1Password `SO / OpenClaw Gateway / gateway-token`
   - Use SSL: ❌
5. Configure → **Enable tool calls** (bidirectional — Legacy controls HA devices)
6. Settings → Voice Assistants → set conversation agent to "OpenClaw"
7. Add card to dashboard: `type: custom:openclaw-chat-card`

> With tool calls enabled in the HACS integration, `homeassistant-skill` may be redundant — the integration provides both directions natively. Install `homeassistant-skill` only if you need agent-initiated HA calls outside of a voice/chat session.

---

### Voice Satellite Hardware

To use voice with the HACS integration you need a microphone + speaker connected to HA (or running ESPHome Wyoming protocol). Options:

| Device | Cost | Setup | Best For |
|--------|------|-------|----------|
| **HA Voice Preview Edition** | $59 | Plug in, HA auto-detects | Turn-key, no config |
| **ESP32 ATOM Echo** | $13 | Flash ESPHome, configure Wyoming | DIY, budget |
| **ESP32-S3-BOX** | ~$40 | ESPHome + Wyoming | More capable DIY |
| **Any old phone** (Android) | $0 | HA Companion app + mic | Quick test |

**Recommended:** HA Voice Preview Edition if you want it to just work. ESP32 ATOM Echo if you want cheap and are comfortable with ESPHome flashing.

**Wake word options** (configured in HA pipeline):
- `hey openclaw` (default from HACS integration)
- `hey jarvis`, `ok nabu`, `alexa` — built-in HA options
- Custom via openWakeWord

---

## Troubleshooting

### "Could not connect to gateway"

1. Verify the gateway is running: `ssh root@64.227.8.23 "systemctl --user status openclaw-gateway"`
2. Check the URL is correct: `https://omos-core.tail717ba0.ts.net`
3. Ensure the token hasn't been rotated (check 1Password)
4. If using Tailscale: ensure both devices are on the same tailnet

### iPhone: "Untrusted Developer"

Settings → General → VPN & Device Management → tap your Apple ID → Trust

### iPhone: App expired (7-day limit)

Re-deploy from Xcode (Cmd+R) or use SideStore for auto-refresh.

### macOS: App won't launch after update

```bash
# Clear cached state
rm -rf ~/Library/Application\ Support/OpenClaw/Cache
# Relaunch
open /Applications/OpenClaw.app
```

### QR Pairing (alternative to token entry)

From any already-connected device (Mac, Discord, etc.):

```
/pair qr
```

This generates a QR code. Scan it with the new device's camera during onboarding setup. The QR contains a short-lived bootstrap token — more secure than sharing the gateway token directly.

---

## Node Capabilities Reference

Each connected node can expose hardware capabilities to Legacy:

| Capability | Mac | iPhone | iPad | Watch |
|------------|-----|--------|------|-------|
| Chat | Yes | Yes | Yes | Yes |
| Voice (Talk Mode) | Yes | Yes | Yes | Yes |
| Camera | Yes | Yes | Yes | — |
| Screen capture | Yes | — | — | — |
| Location | — | Yes | Yes | Yes |
| Motion sensors | — | Yes | Yes | Yes |
| Desktop control | Yes | — | — | — |
| Notifications | Yes | Yes | Yes | Yes |
| Canvas display | Yes | Yes | Yes | — |
| File system access | Yes | — | — | — |
| Browser control | Yes | — | — | — |

---

## Adding a New Node Type

When adding support for a new device type:

1. Add a section to this document with setup instructions
2. Document the pairing/connection process
3. List the node capabilities it exposes
4. Note any signing/update limitations
5. Add troubleshooting entries for common issues


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
