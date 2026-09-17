---
id: ui-02-20260111053925
title: "Agent File Cleanup Guide"
type: reference
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/02]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Agent File Cleanup Guide"
---

# Agent File Cleanup Guide

**Created:** 2025-11-08  
**Purpose:** Guide for cleaning up duplicate agent files after renaming

---

## Current State

We have TWO sets of agent files:
1. **OLD files** (without domain prefixes): `01 - Intelligence Agent.md`, etc.
2. **NEW files** (with domain prefixes): `01 - HP - Intelligence Agent.md`, etc.

---

## Files to DELETE (Old Files Without Prefixes)

Delete these 13 files from `Legacy AI/Infrastructure - Agents/`:

- [ ] `00 - Legacy Agent.md` (replaced by `00 - LEGACY - Master Orchestrator.md`)
- [ ] `01 - Intelligence Agent.md` (replaced by `01 - HP - Intelligence Agent.md`)
- [ ] `02 - Performance Agent.md` (need to create `02 - HP - Performance Agent.md`)
- [ ] `03 - Finance Agent.md` (need to create `03 - HP - Finance Agent.md`)
- [ ] `04 - Business Operations Agent.md` (need to create `04 - GE - Business Operations Agent.md`)
- [ ] `05 - Founder Assistant Agent.md` (need to create `05 - GE - Founder Assistant Agent.md`)
- [ ] `06 - Capital Systems Agent.md` (need to create `06 - GE - Capital Systems Agent.md`)
- [ ] `07 - Home Command Agent.md` (need to create `07 - LE - Home Command Agent.md`)
- [ ] `08 - Homestead Operations Agent.md` (need to create `08 - LE - Homestead Operations Agent.md`)
- [ ] `09 - Family Wealth Agent.md` (need to create `09 - LE - Family Wealth Agent.md`)
- [ ] `10 - Infrastructure Agent.md` (need to create `10 - IT - Infrastructure Agent.md`)
- [ ] `11 - Robotics & Automation Agent.md` (need to create `11 - IT - Robotics & Automation Agent.md`)
- [ ] `12 - Communications Agent.md` (need to create `12 - EI - Communications Agent.md`)

---

## Files to KEEP (New Files With Prefixes)

Keep these files:

✅ `00 - Agent Roster Overview.md` - Overview document  
✅ `00 - LEGACY - Master Orchestrator.md` - Renamed  
✅ `01 - HP - Intelligence Agent.md` - Renamed  
✅ `Agent Naming Convention.md` - Reference guide

---

## Action Required

**Option 1: Manual rename in Obsidian (RECOMMENDED)**
1. Open Obsidian
2. Navigate to `Legacy AI/Infrastructure - Agents/`
3. For each file 02-12, rename by adding the domain prefix:
   - `02 - Performance Agent` → `02 - HP - Performance Agent`
   - `03 - Finance Agent` → `03 - HP - Finance Agent`
   - `04 - Business Operations Agent` → `04 - GE - Business Operations Agent`
   - `05 - Founder Assistant Agent` → `05 - GE - Founder Assistant Agent`
   - `06 - Capital Systems Agent` → `06 - GE - Capital Systems Agent`
   - `07 - Home Command Agent` → `07 - LE - Home Command Agent`
   - `08 - Homestead Operations Agent` → `08 - LE - Homestead Operations Agent`
   - `09 - Family Wealth Agent` → `09 - LE - Family Wealth Agent`
   - `10 - Infrastructure Agent` → `10 - IT - Infrastructure Agent`
   - `11 - Robotics & Automation Agent` → `11 - IT - Robotics & Automation Agent`
   - `12 - Communications Agent` → `12 - EI - Communications Agent`
4. Delete the old `00 - Legacy Agent.md` file

**This takes about 2 minutes and preserves all content perfectly.**

---

## Domain Prefix Reference

- **LEGACY** = Master Orchestrator (red)
- **HP** = Holistic Performance (blue)
- **GE** = Generational Entrepreneurship (green)
- **LE** = Legacy Evolution (purple)
- **IT** = Innovative Technology (orange)
- **EI** = External Interface

---

*Delete this file after cleanup is complete*
