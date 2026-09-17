---
id: root-00-20260113225954
title: "SECRETS MANAGEMENT"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "SECRETS MANAGEMENT"
---
# OneMind - Secrets Management

## Quick Reference

| What | Location |
|------|----------|
| Your secrets file | `/Users/zeuslegacy/.secrets/onemind.env` |
| Backup | 1Password → "OneMind Secrets" |
| Template (no real keys) | `.env.example` in repo |

---

## How It Works

```
1Password (backup)
     │
     ▼
~/.secrets/onemind.env  ←── Your actual keys live here
     │
     ▼ (loaded automatically when terminal opens)
     │
Terminal / Claude Code  ←── Keys available as $LETTA_API_KEY, etc.
```

---

## Day-to-Day Usage

**You don't need to do anything.** Keys load automatically when you open a terminal.

To verify:
```bash
echo $LETTA_API_KEY
```

---

## When You Get a New Key

### Step 1: Edit the file
```bash
code ~/.secrets/onemind.env
```

### Step 2: Add the new key
```bash
export MY_NEW_KEY="the-key-value"
```

### Step 3: Reload
Open a new terminal, OR run:
```bash
source ~/.secrets/onemind.env
```

### Step 4: Update 1Password backup
Copy the updated file contents to your 1Password "OneMind Secrets" note.

---

## Finding the Secrets File

The file is hidden (starts with `.`). To access:

**Terminal:**
```bash
code ~/.secrets/onemind.env
```

**Finder:**
1. Press `Cmd + Shift + G`
2. Type: `~/.secrets/`
3. Press Enter

---

## Backing Up to 1Password

1. Open 1Password
2. Create new **Secure Note** called "OneMind Secrets"
3. Copy entire contents of `~/.secrets/onemind.env`
4. Paste into the note
5. Save

---

## If You Lose the File

1. Open 1Password
2. Find "OneMind Secrets" note
3. Create new file: `~/.secrets/onemind.env`
4. Paste contents from 1Password
5. Set permissions: `chmod 600 ~/.secrets/onemind.env`
6. Restart terminal

---

## Current Environment Variables

```bash
# Letta AI
LETTA_API_KEY
LEGACY_AI_AGENT_ID
HP_CONDUCTOR_AGENT_ID
LE_CONDUCTOR_AGENT_ID
GE_CONDUCTOR_AGENT_ID
IT_CONDUCTOR_AGENT_ID

# Discord
DISCORD_BOT_TOKEN
DISCORD_GUILD_ID
DISCORD_OWNER_ID

# Todoist
TODOIST_API_TOKEN

# GitHub
GITHUB_TOKEN
GITHUB_REPO

# Telegram
TELEGRAM_BOT_TOKEN
TELEGRAM_CHAT_ID

# Zoho
ZOHO_CLIENT_ID
ZOHO_CLIENT_SECRET
ZOHO_CRM_REFRESH_TOKEN
ZOHO_DESK_REFRESH_TOKEN
ZOHO_DESK_ORG_ID
ZOHO_CALENDAR_REFRESH_TOKEN

# Circle.so
CIRCLE_API_TOKEN
CIRCLE_COMMUNITY_ID

# Home Assistant
HA_URL
HA_TOKEN
MQTT_PASSWORD
HA_SECRET_KEY
```
