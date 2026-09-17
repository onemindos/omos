---
id: so-01-20260412000011
title: "Zeus Identity — JSON Context Profile"
type: reference
status: active
created: 2026-04-12
owner: zeus
tags: [path/so, path/so/01, type/reference, prompt-engineering, json, identity]
ai_summary: "Production JSON profile for Zeus's identity — paste when briefing any new AI agent or tool so it knows who you are without a 20-minute setup conversation"
---
# Zeus Identity — JSON Context Profile

> **How to use:** Paste at the top of any new AI conversation, agent setup, or tool configuration. The AI instantly knows who you are, how you work, and what not to do.

---

## Profile

```json
{
  "profile_type": "user_identity",
  "name": "Zeus Delacruz",
  "handle": "@zeuslegacy",
  "org": "OneMind-OS",
  "role": "Founder and architect of One Mind OS — AI-native life operating framework",
  "location": "NYC (primary base)",
  "communication_style": {
    "preferred": ["direct", "outcome-first", "no filler"],
    "format": "Brief by default. Tables and bullets over paragraphs. Skip summaries unless asked.",
    "dislikes": [
      "Long intros before answering",
      "Restating what was just said",
      "Asking obvious clarifying questions when context is sufficient",
      "Over-explaining steps already understood"
    ]
  },
  "what_i_am_building": {
    "product": "One Mind OS — the first AI-native human enhancement framework",
    "one_liner": "One Mind is how you and AI run your life together.",
    "method": "CODEX: Capture → Organize → Direct → Execute → × Multiply",
    "model": "3+1 life domains: HP (self) + LE (foundation) + GE (empire) + SO (operating layer)",
    "status_2026": "Pre-launch — writing manifesto, polishing template vault, building community"
  },
  "active_infrastructure": {
    "vault_local": "/Users/zeuslegacy/Documents/onemind-codex/",
    "vault_vps": "/home/onemind/onemind-os/codex/",
    "vault_public": "codex.onemindos.com",
    "agents_on_vps": ["Legacy (orchestrator)", "Spartan", "Oracle", "Haven", "Empire", "Guardian", "Grid", "Trinity"],
    "key_tools": ["Obsidian", "VS Code + Copilot", "Todoist", "1Password", "Telegram", "n8n", "Docker", "OpenClaw"]
  },
  "domain_structure": {
    "SO_00_24": "Sovereign Operations — AI, tech, infrastructure (the operating layer)",
    "HP_25_49": "Holistic Performance — health, skills, money, identity (YOUR SELF)",
    "LE_50_74": "Legacy Evolution — relationships, home, community (YOUR FOUNDATION)",
    "GE_75_99": "Generational Empire — brand, ventures, products, income (YOUR EMPIRE)"
  },
  "priorities_april_2026": [
    "Write One Mind manifesto — 50-80 page ebook",
    "Polish onemind-codex-template for public release",
    "Deploy Obsidian Headless Sync on VPS",
    "Launch One Mind community (Circle or Discord)"
  ],
  "hard_rules": [
    "Never delete files or tasks without explicit approval — always ask first",
    "Never rename SO/HP/LE/GE domain folders without explicit approval",
    "Never edit ONEMIND-CODEX.md or CODEX-FRAMEWORK.md casually",
    "Never store credentials in markdown — reference 1Password instead",
    "Never say 'quadrant' — say 'domain'",
    "Never say 'enterprise' or 'endeavor' for GE — say 'empire'",
    "Never say 'Your Family' for LE — say 'Your Foundation'"
  ],
  "credentials": "All stored in 1Password vault '00-24 SO (Sovereign Operations)', ID: aaeyris5z44pdqspk277u2oppi. Reference format: op://aaeyris5z44pdqspk277u2oppi/Item Name/field"
}
```

---

## When to Use This

- Briefing a new AI agent for the first time
- Starting a fresh Claude / ChatGPT / Gemini session on a complex task
- Setting up a new tool that needs to "know you"
- Configuring an OpenClaw agent workspace


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
