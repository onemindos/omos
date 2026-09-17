#!/usr/bin/env python3
"""Migrate VPS openclaw.json from 4-agent to 9-agent architecture."""
import json
import shutil
from pathlib import Path

CONFIG = Path("/home/onemind/.openclaw/openclaw.json")

# Backup
shutil.copy(CONFIG, CONFIG.with_suffix(".json.bak.pre-9agent"))

with open(CONFIG) as f:
    data = json.load(f)

# ── 1. BINDINGS ──────────────────────────────────────────────────────
# Replace old 3 bindings (hp/le/ge) with 8 new ones (all non-legacy agents)
CHANNEL_MAP = {
    "spartan":  "1480445911707549750",
    "oracle":   "1480450417522643127",
    "eden":     "1480445924202385499",
    "empire":   "1480445930661744801",
    "haven":    "1480445937665966161",
    "guardian":  "1480445943789785180",
    "grid":     "1480445951805096007",
    "trinity":  "1480445958876827714",
}

data["bindings"] = [
    {
        "agentId": aid,
        "match": {
            "channel": "discord",
            "peer": {"kind": "channel", "id": cid}
        }
    }
    for aid, cid in CHANNEL_MAP.items()
]

# ── 2. AGENTS LIST ──────────────────────────────────────────────────
WS = "/home/onemind/onemind-os/openclaw/workspaces"
CODEX = "/home/onemind/onemind-os/codex"

# Keep legacy (first item) as-is, but ensure workspace is set
legacy_agent = data["agents"]["list"][0]  # id: legacy
legacy_agent["workspace"] = f"{WS}/legacy"

domain_agents = [
    {
        "id": "spartan",
        "workspace": f"{WS}/spartan",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/25-49 HP (Holistic Performance)",
        "heartbeat_to": "1480445911707549750",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check body, fitness, physical mastery priorities. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("08:00", "22:00"),
    },
    {
        "id": "oracle",
        "workspace": f"{WS}/oracle",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/25-49 HP (Holistic Performance)",
        "heartbeat_to": "1480450417522643127",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check mind, finance, meaning, relationships priorities. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("07:00", "23:00"),
    },
    {
        "id": "eden",
        "workspace": f"{WS}/eden",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/50-74 LE (Legacy Evolution)",
        "heartbeat_to": "1480445924202385499",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check garden, food growing, agriculture, land stewardship. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("07:00", "20:00"),
    },
    {
        "id": "haven",
        "workspace": f"{WS}/haven",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/50-74 LE (Legacy Evolution)",
        "heartbeat_to": "1480445937665966161",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check home, family, estate, property matters. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("09:00", "20:00"),
    },
    {
        "id": "empire",
        "workspace": f"{WS}/empire",
        "model": "azure-openai/gpt-5.2",
        "codex": f"{CODEX}/75-99 GE (Generational Entrepreneurship)",
        "heartbeat_to": "1480445930661744801",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check business ops, revenue, ventures, deals. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("07:00", "23:00"),
    },
    {
        "id": "guardian",
        "workspace": f"{WS}/guardian",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/00-24 UI (Unified Intelligence)",
        "heartbeat_to": "1480445943789785180",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check security, backups, infrastructure health. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("00:00", "23:59"),
    },
    {
        "id": "grid",
        "workspace": f"{WS}/grid",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/00-24 UI (Unified Intelligence)",
        "heartbeat_to": "1480445951805096007",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check tech systems, integrations, NATS, services. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("07:00", "23:00"),
    },
    {
        "id": "trinity",
        "workspace": f"{WS}/trinity",
        "model": "azure-openai/gpt-4o",
        "codex": f"{CODEX}/75-99 GE (Generational Entrepreneurship)",
        "heartbeat_to": "1480445958876827714",
        "heartbeat_prompt": "Read HEARTBEAT.md. Check community engagement, growth, partnerships. Surface anything urgent. HEARTBEAT_OK if clear.",
        "hours": ("08:00", "22:00"),
    },
]

new_agents = [legacy_agent]
for a in domain_agents:
    agent = {
        "id": a["id"],
        "workspace": a["workspace"],
        "model": a["model"],
        "memorySearch": {
            "extraPaths": [a["codex"]]
        },
        "heartbeat": {
            "every": "1h",
            "activeHours": {
                "start": a["hours"][0],
                "end": a["hours"][1],
                "timezone": "America/New_York"
            },
            "target": "discord",
            "to": a["heartbeat_to"],
            "prompt": a["heartbeat_prompt"]
        },
        "tools": {
            "profile": "messaging",
            "allow": ["web_fetch", "web_search", "group:memory", "lobster"]
        }
    }
    new_agents.append(agent)

data["agents"]["list"] = new_agents

# ── 3. GUILD CHANNELS ───────────────────────────────────────────────
ALL_CHANNELS = {
    "1478604243509448705": {"allow": True},   # legacy-chat
    "1478866012440129671": {"allow": True},   # agent-heartbeat
    "1478604227772493836": {"allow": True},   # morning-brief
    "1480609296876650567": {"allow": True},   # notifications
    "1480445911707549750": {"allow": True},   # spartan-chat
    "1480450417522643127": {"allow": True},   # oracle-chat
    "1480445924202385499": {"allow": True},   # eden-chat
    "1480445937665966161": {"allow": True},   # haven-chat
    "1480445930661744801": {"allow": True},   # empire-chat
    "1480445958876827714": {"allow": True},   # trinity-chat
    "1480445951805096007": {"allow": True},   # grid-chat
    "1480445943789785180": {"allow": True},   # guardian-chat
}

guild = data["channels"]["discord"]["guilds"]["1453411552181883027"]
guild["channels"] = ALL_CHANNELS

# ── 4. HOOKS allowedAgentIds ─────────────────────────────────────────
ALL_AGENT_IDS = [
    "legacy", "spartan", "oracle", "eden", "empire",
    "haven", "guardian", "grid", "trinity"
]
data["hooks"]["allowedAgentIds"] = ALL_AGENT_IDS

# ── 5. TOOLS agentToAgent.allow ──────────────────────────────────────
data["tools"]["agentToAgent"]["allow"] = ALL_AGENT_IDS

# ── 6. SESSION dmScope ───────────────────────────────────────────────
if "session" not in data:
    data["session"] = {}
data["session"]["dmScope"] = "per-channel-peer"

# ── WRITE ─────────────────────────────────────────────────────────────
with open(CONFIG, "w") as f:
    json.dump(data, f, indent=2)

print("Migration complete.")
print(f"  Bindings: {len(data['bindings'])}")
print(f"  Agents: {len(data['agents']['list'])}")
print(f"  Guild channels: {len(ALL_CHANNELS)}")
print(f"  Hooks agents: {data['hooks']['allowedAgentIds']}")
print(f"  Agent-to-agent: {data['tools']['agentToAgent']['allow']}")
print(f"  Session dmScope: {data['session']['dmScope']}")
print(f"  Backup at: {CONFIG.with_suffix('.json.bak.pre-9agent')}")
