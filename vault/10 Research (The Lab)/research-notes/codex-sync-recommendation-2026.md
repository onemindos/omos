---
title: "Codex Sync — Getting It Working Now"
type: reference
status: active
created: 2026-04-06
owner: zeus
tags: [type/reference, domain/so, topic/sync, topic/infrastructure, topic/openclaw]
ai_summary: "Practical sync recommendation for the codex — what to do RIGHT NOW to get Mac, VPS, and Phone syncing without complexity. Compares all methods with honest tradeoffs."
---

# Codex Sync — Getting It Working Now

> "What is the best way to put the codex on the VPS now, see it working, and sync without all the extra stuff?"

---

## The Short Answer

**Use git for Mac ↔ VPS (you already have this). Use Obsidian Sync for Mac ↔ Phone. They don't conflict.**

The simplest path that works TODAY with zero new infrastructure:

```
Mac (Obsidian + VS Code)
  ↓ git push (manual or scripted)
GitHub (OneMind-OS/onemind-codex)
  ↓ onemind-pull.timer (already installed, 1min poll)
VPS (/home/onemind/onemind-os/codex/)
  ↓ OpenClaw reads workspace from here
Agents read + write codex files
```

This is what you have now. It works. The question is: what's the **minimum change** to make it feel seamless?

---

## Sync Methods Compared

| Method | Mac → VPS | VPS → Mac | Phone | Binary Files (PDFs, images) | Latency | Cost | Complexity |
|--------|-----------|-----------|-------|----------------------------|---------|------|-----------|
| **Git (current)** | ✅ `git push` | ✅ `git pull` (1min timer) | ❌ No | ❌ Excluded via .gitignore | 1-60 sec | Free | Low |
| **Obsidian Sync** | ✅ Real-time | ✅ Real-time | ✅ Real-time | ✅ All files | 1-5 sec | $8-16/mo | Medium |
| **Syncthing** | ✅ Real-time | ✅ Real-time | ⚠️ Android only | ✅ All files | 1-10 sec | Free | Medium |
| **rsync + cron** | ✅ One-way push | ✅ Or pull | ❌ No | ✅ All files | Cron interval | Free | Low |
| **Unison** | ✅ Bidirectional | ✅ Bidirectional | ❌ No | ✅ All files | Manual/cron | Free | High |
| **Dropbox** | ✅ Real-time | ✅ Real-time | ✅ Real-time | ✅ All files | 2-10 sec | $12/mo | Low |

---

## The RECOMMENDED Path: 3 Phases

### Phase 0: Fix What You Have (TODAY — 15 minutes)

Your current git sync works but has two problems:
1. `onemind-watch.service` is disabled (detached HEAD bug)
2. VPS → Mac sync requires manual `git pull`

**Fix onemind-watch.sh on VPS:**

```bash
ssh root@104.248.4.212

# Edit the watch script
cat > /opt/onemind/onemind-watch.sh << 'EOF'
#!/bin/bash
# onemind-watch.sh — event-driven git commit + push
# Watches for file changes, debounces 3 seconds, commits and pushes

CODEX_DIR="/home/onemind/onemind-os/codex"
cd "$CODEX_DIR" || exit 1

# Guard: ensure we're on main branch
BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
if [[ -z "$BRANCH" ]]; then
  echo "[onemind-watch] Detached HEAD detected, checking out main"
  git checkout main
fi

inotifywait -m -r -e modify,create,delete,move \
  --exclude '\.git|\.obsidian' \
  "$CODEX_DIR" |
while read -r _dir _event _file; do
  sleep 3  # debounce

  # Drain any queued events during debounce
  while read -t 0.1 -r _ _ _; do :; done

  cd "$CODEX_DIR"

  # Check branch again
  BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null)
  if [[ -z "$BRANCH" ]]; then
    git checkout main
  fi

  git add -A
  CHANGED=$(git diff --cached --name-only | head -5 | tr '\n' ', ')
  if [[ -n "$CHANGED" ]]; then
    git commit -m "agent: update ${CHANGED%,}"
    git push origin main
  fi
done
EOF

chmod +x /opt/onemind/onemind-watch.sh

# Re-enable the service
systemctl enable --now onemind-watch.service
```

**Result:** Agent edits on VPS → auto-commit + push to GitHub → you `git pull` on Mac (or set up a pull hook).

### Phase 1: Add Obsidian Sync (WHEN READY — $8-16/mo)

This is the endgame for phone access and real-time sync. But it's not required today.

```
Mac (Obsidian)
  ↕ Obsidian Sync (real-time, E2E encrypted)
Phone (Obsidian Mobile)

Mac
  ↕ git push/pull
GitHub
  ↕ onemind-pull.timer + onemind-watch.service
VPS
```

- Obsidian Sync handles Mac ↔ Phone (real-time, all files including PDFs)
- Git handles Mac ↔ VPS (text only, infinite history)
- They don't conflict: Obsidian Sync ignores `.git/`, git ignores `.obsidian/`

**Wait on this until:** You're actively using the phone workflow daily and git sync feels too slow.

### Phase 2: Obsidian Headless Sync on VPS (FUTURE)

Replaces `onemind-pull.timer` with real-time bidirectional sync:

```bash
# On VPS
npm install -g obsidian-headless
ob login
ob sync-setup --vault "OneMind Codex"
ob sync --continuous  # replaces onemind-pull.timer
```

This gives VPS real-time access to ALL files (including PDFs for RAG pipeline). But it requires Obsidian Sync subscription and `obsidian-headless` to be stable.

**Wait on this until:** Phase 1 is working and you want RAG over PDFs on the VPS.

---

## What to Do RIGHT NOW (Actionable Steps)

### Step 1: Verify VPS git sync is working

```bash
ssh root@104.248.4.212
cd /home/onemind/onemind-os/codex
git status
git log --oneline -5
```

If the repo is there and on `main`, you're good. If there are issues, fix them first.

### Step 2: Fix and re-enable onemind-watch.service

Apply the script fix above. This makes VPS → GitHub automatic (agent writes → 3s → committed + pushed).

### Step 3: Set up a Mac pull hook (optional)

If you want Mac to auto-pull when GitHub gets new commits:

```bash
# Simple: cron every 5 minutes
echo "*/5 * * * * cd /Users/zeuslegacy/Documents/onemind-codex && git pull --rebase origin main 2>&1 | logger -t onemind-pull" | crontab -

# Or: use a launchd plist for more control
```

### Step 4: Configure OpenClaw workspace to point at codex

On VPS, tell OpenClaw where the codex lives:

```json
{
  "agents": {
    "defaults": {
      "workspace": "/home/onemind/onemind-os/codex"
    }
  }
}
```

Or per-agent workspaces for domain isolation:
- Legacy → full codex root
- Coach → `/home/onemind/onemind-os/codex/25-49 HP (Holistic Performance)/`
- Eden → `/home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/`
- Empire → `/home/onemind/onemind-os/codex/75-99 GE (Generational Empire)/`

### Step 5: Test the round trip

1. Edit a file on Mac → `git push`
2. Wait 60 seconds (onemind-pull.timer)
3. SSH to VPS → verify file is updated
4. Agent edits a file on VPS
5. onemind-watch.service commits + pushes (3 seconds)
6. On Mac → `git pull` → verify agent's edit arrived

If this round trip works, you have a working sync. Everything else is optimization.

---

## Why NOT to Overcomplicate This

You asked "without all the extra stuff" and that's the right instinct. Here's what you DON'T need right now:

| Thing | Skip Because |
|-------|-------------|
| Dropbox layer | Git + Obsidian Sync covers everything |
| Syncthing | Doesn't do iOS, git is simpler |
| rsync scripts | Git already does this better with history |
| Custom webhook server | onemind-pull.timer is fine for now |
| NATS-based sync | Overengineered for file sync |

The architecture in CLAUDE.md (3 sync planes) is the DESTINATION. The path there is incremental:

1. **Now:** Git only (Mac ↔ VPS via GitHub) ← you're here
2. **Soon:** + fix onemind-watch (VPS auto-push) ← 15 min fix
3. **When ready:** + Obsidian Sync (Mac ↔ Phone) ← $8-16/mo
4. **Future:** + Obsidian Headless on VPS (replace git pull timer) ← Phase 2
5. **Optional:** + Dropbox backup layer ← only if you want 180-day file history

Each step is independent. Each step is reversible. Start with step 1-2 and see how it feels.

---

*See also: Sync Architecture section in [[CLAUDE]], [[onemind-openclaw-vision-2027]], [[openclaw-deployment-comparison-2026]]*
