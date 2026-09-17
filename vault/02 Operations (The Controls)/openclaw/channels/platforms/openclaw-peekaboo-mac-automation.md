# OpenClaw × Peekaboo — macOS Desktop Automation

> Control any Mac app (including Bambu Studio) from OpenClaw using element-based UI targeting — not fragile pixel coordinates.

---

## What Is Peekaboo?

Peekaboo is the official first-party macOS UI automation skill for OpenClaw, published in the `openclaw/skills` repo. It uses macOS Accessibility APIs to:

- Capture and inspect screens
- Target UI elements by **ID or query** (not raw coordinates)
- Click buttons, type text, navigate menus
- Manage apps, windows, and the clipboard
- Output JSON for scripting and chaining commands

**Why it doesn't break:** Element ID targeting survives UI layout changes. Even if Bambu Studio redesigns their interface, Peekaboo re-identifies elements by label — not pixel position.

---

## Architecture

```
Your command (via Legacy / Claude)
        ↓
  OpenClaw on Mac mini
        ↓
  Peekaboo skill (clawhub install peekaboo)
        ↓
  Peekaboo CLI (brew install peekaboo)
        ↓
  Mac mini desktop → Bambu Studio
```

OpenClaw hosts a **Peekaboo Bridge** — a local UNIX socket server that brokers macOS TCC permissions so the CLI reuses the app's grants. You only need to approve permissions once.

---

## Install — 3 Commands on the Mac mini

```bash
# 1. Install Peekaboo CLI via Homebrew
brew install steipete/tap/peekaboo

# 2. Install the OpenClaw skill
clawhub install peekaboo

# 3. Verify permissions are set correctly
peekaboo permissions
```

---

## Mac Permissions (Required)

Go to **System Settings → Privacy & Security** and grant the following to OpenClaw (and Terminal if needed):

| Permission | Why |
|---|---|
| **Accessibility** | Clicking, typing, UI interaction |
| **Screen Recording** | Screenshots and UI inspection |

Enable the bridge inside OpenClaw: **Settings → Enable Peekaboo Bridge**

---

## Core Commands

### Check what's on screen
```bash
# Annotate all UI element IDs on screen
peekaboo see --app "BambuStudio" --annotate --path /tmp/bambu.png

# See the whole screen
peekaboo see --mode screen --screen-index 0 --analyze "What buttons are visible?"
```

### Click and interact
```bash
# Click by element ID (from annotated screenshot)
peekaboo click --on B1 --app "BambuStudio"

# Click by coordinates (fallback)
peekaboo click --coords 450,300 --app "BambuStudio"

# Type text into active field
peekaboo type "print settings" --app "BambuStudio"

# Press a key
peekaboo press return --app "BambuStudio"
```

### App and window control
```bash
# Launch an app
peekaboo app launch "BambuStudio"

# Focus a specific window
peekaboo window focus --app "BambuStudio" --window-title "Print"

# Navigate a menu path
peekaboo menu click --app "BambuStudio" --path "File > Export"
```

### Screenshots
```bash
# Retina screenshot of a specific app
peekaboo image --app "BambuStudio" --retina --path /tmp/screen.png

# Screenshot with AI analysis
peekaboo image --app "BambuStudio" --analyze "What is the current print status?"
```

---

## Example Workflow — Bambu Studio Print Job

```bash
# Step 1 — Launch and focus
peekaboo app launch "BambuStudio"
peekaboo window focus --app "BambuStudio"

# Step 2 — Capture and inspect the UI
peekaboo see --app "BambuStudio" --annotate --path /tmp/bambu-annotated.png

# Step 3 — Click the Print button (by element ID from annotation)
peekaboo click --on B4 --app "BambuStudio"

# Step 4 — Confirm via screenshot
peekaboo image --app "BambuStudio" --analyze "Did the print dialog open?"
```

---

## Peekaboo Bridge (Recommended)

OpenClaw can act as a Peekaboo Bridge host — this is the cleanest setup:

1. Open OpenClaw → **Settings → Enable Peekaboo Bridge**
2. OpenClaw starts a local UNIX socket server
3. Peekaboo CLI connects through it automatically
4. TCC permissions are shared — no duplicate permission grants

```bash
# Check which host Peekaboo is using
peekaboo bridge status --verbose
```

---

## Troubleshooting

| Issue | Fix |
|---|---|
| `Permission denied` | Run `peekaboo permissions` and grant Accessibility + Screen Recording |
| `Bridge client not authorized` | Ensure OpenClaw is signed, or set `PEEKABOO_ALLOW_UNSIGNED_SOCKET_CLIENTS=1` (debug only) |
| Element not found | Re-run `peekaboo see --annotate` — element IDs refresh each snapshot |
| Click lands in wrong spot | Check `peekaboo list screens` — confirm screen resolution matches |

---

## Resources

- **Peekaboo homepage:** https://peekaboo.boo
- **ClawHub skill page:** https://clawhub.com/skills/peekaboo
- **OpenClaw Peekaboo Bridge docs:** https://docs.openclaw.ai/platforms/mac/peekaboo
- **OpenClaw skills docs:** https://docs.openclaw.ai/tools/skills

---

> **Security note:** Always run `clawhub inspect peekaboo` or check the VirusTotal report on ClawHub before installing any community skill. Peekaboo is first-party (`openclaw/skills`) so it's safe — apply this rule to anything else you install.


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
