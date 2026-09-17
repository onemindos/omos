# NODES.md — Client Device Setup

Instructions for connecting client devices (nodes) to the OneMind OS gateway. Each connected device becomes an OpenClaw **node** — a trusted endpoint that can chat with Legacy, run tools, and receive canvases.

---

## Gateway Connection Info

| Field | Value |
|-------|-------|
| Gateway URL | `https://gateway.onemindos.dev` |
| WebSocket URL | `wss://gateway.onemindos.dev` |
| Tailscale URL | `https://onemind-main.tail717ba0.ts.net` |
| Gateway Token | Stored in 1Password: `00-24 SO / OpenClaw Gateway / gateway-token` |
| QR Pairing | From any paired device: run `/pair qr` in chat |

> **Security**: Never share the gateway token publicly. Use QR pairing for new devices when possible.

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

### Connect to Remote Gateway

1. Open **OpenClaw.app**
2. On first launch, the onboarding wizard starts
3. Choose **Remote Gateway** mode
4. Enter the gateway URL: `https://gateway.onemindos.dev`
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

## Troubleshooting

### "Could not connect to gateway"

1. Verify the gateway is running: `ssh root@104.248.4.212 "systemctl status openclaw"`
2. Check the URL is correct: `https://gateway.onemindos.dev`
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
