# TOOLS.md — Haven Tool Usage Guide

## Codex Access

The Codex is at `/home/onemind/onemind-os/codex/` — your native knowledge search is built-in.

### ⛔ Codex Write Rules (CRITICAL)

The codex directory structure is **immutable** — managed exclusively by Zeus from the local IDE.

- **NEVER create new folders** in the codex. The numbered folder hierarchy (00-12, 25-34, 50-58, 75-84) is fixed.
- **NEVER rename or renumber** existing codex folders.
- **New notes go to `00-24 SO/06 Inbox (Queue)/`** — Zeus will file them to the correct location.
- **If a folder doesn't exist for a topic**, put the note in `06 Inbox (Queue)/` or `99-organize/` — do NOT create a new folder.
- The codex syncs via git. Any rogue folders you create get auto-committed and pollute the repo.

### memory_search (primary — use this first)

Semantically searches your memory + the Codex vault. Uses BM25+vector hybrid search.

```
memory_search "planting schedule tomatoes"    → garden plans
memory_search "chicken feeding schedule"      → livestock care
memory_search "irrigation well pump"          → water systems
memory_search "home assistant scene bedroom"  → smart home
memory_search "fence repair schedule"         → property maintenance
memory_search "soil moisture reading"         → sensor data
```

When to use: ANY question about property, farm, garden, livestock, or home systems. Always search before answering.

### read_file (direct access)

```
read_file /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/50 Home/...
```

## Home Assistant (Smart Home Control)

Use for managing the smart home:
- Turn devices on/off (lights, switches, climate)
- Trigger automations and scenes
- Check sensor readings
- Monitor security cameras
- Control irrigation systems

## Weather

Use for monitoring conditions:
- Current conditions and forecasts
- Severe weather alerts
- Growing degree days
- Frost predictions
- Rain probability for irrigation planning

## Todoist (Task Management)

Use for creating and tracking property tasks:
- Maintenance schedules and reminders
- Garden planting and harvest tasks
- Livestock care routines
- Supply shopping lists

## File Tools (read / write / edit)

You have direct file access to the codex:

- `read` — read any codex file by path
- `write` — create new notes (ONLY in `06 Inbox (Queue)/` or `99-organize/`, or update existing files in your domain)
- `edit` — make precise edits to existing codex documents

To create a new note:
```
write /home/onemind/onemind-os/codex/00-24 SO (Sovereign Operations)/06 Inbox (Queue)/My-New-Note.md "# Title\n\nContent here"
```

To update an existing file in your domain:
```
edit /home/onemind/onemind-os/codex/50-74 LE (Legacy Evolution)/50 Home/existing-note.md
```

## Domain Quick Reference

| Area | What's There |
|------|-------------|
| Property Index | Inventory, valuations, maps |
| Property Management | Maintenance, repairs, vendors |
| Garden | Planting, rotation, harvest |
| Livestock | Animal care, feeding, vet records |
| Smart Home | Devices, automations, scenes |
| Water Systems | Irrigation, wells, rain collection |
| Land Management | Fencing, roads, outbuildings |
