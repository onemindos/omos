---
id: root-00-20260113225954
title: "FRAMEWORK CONSOLIDATION PLAN"
type: stub
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/stub, status/stub]
ai_summary: "FRAMEWORK CONSOLIDATION PLAN"
---
# Framework Consolidation Plan

## Goal
Create a single source of truth that both humans and AI agents can navigate without confusion.

---

## The New Structure

### ROOT LEVEL (Critical Files Only)

| File | Purpose | Status |
|------|---------|--------|
| **ONEMIND-CODEX.md** | THE CANONICAL SOURCE OF TRUTH | ✅ KEEP (Master) |
| **README.md** | GitHub landing page | ✅ KEEP (Points to CODEX) |
| **CLAUDE.md** | AI assistant instructions | ✅ KEEP (Technical) |
| **KNOWLEDGE-BASE.md** | Quick navigation index | 🔄 UPDATE (Index to CODEX) |
| All other docs | Move to appropriate locations | 🗑️ RELOCATE |

---

## 00 Framework Folder - NEW STRUCTURE

```
00-24 SO/00 Framework (Doctrine)/
├── INDEX.md                          ← Navigation hub (NEW)
├── QUICK-START.md                    ← 5-minute onboarding (NEW)
│
├── 00.1 Core-Doctrine/
│   ├── README.md                     ← Points to root ONEMIND-CODEX.md
│   └── (DELETE all other files - they duplicate CODEX)
│
├── 00.2 Architecture/
│   ├── Agent-Architecture.md         ← Specific: How agents work
│   ├── Vault-Structure.md            ← Specific: Folder organization
│   └── Integration-Map.md            ← Specific: Tool connections
│
├── 00.3 Operating-System/
│   ├── Daily-Rhythm.md               ← Extracted from CODEX (templates)
│   ├── Weekly-Review.md              ← Extracted from CODEX (templates)
│   └── Quarterly-Planning.md         ← Extracted from CODEX (templates)
│
├── 00.4 AI-Agent-Guides/
│   ├── README.md                     ← How agents use this vault
│   ├── Navigation-Protocol.md        ← Where agents look first
│   ├── Memory-Sync-Protocol.md       ← How agents stay current
│   └── Routing-Rules.md              ← How to route requests
│
└── 00.5 Reference/
    ├── Vocabulary.md                 ← Glossary (from CODEX Appendix A)
    ├── Naming-Conventions.md         ← File/folder naming rules
    └── Decision-Trees.md             ← Quick routing decisions
```

---

## What to DO with Each Current File

### 00.1 Core-Doctrine/

| Current File | Action | Reason |
|--------------|--------|--------|
| `ONE-MIND-FRAMEWORK-v1.0.md` | 🗑️ DELETE | 80% duplicate of ONEMIND-CODEX.md |
| `Zeus-LIFE-OS-ARCHITECTURE.md` | 🗑️ DELETE | Implementation details in CODEX Part 6 |

**Replace with:**
- `README.md` → Single file that says: "See root ONEMIND-CODEX.md - THE canonical source"

---

### 00.2 Architecture/

| Current File | Action | Reason |
|--------------|--------|--------|
| `Vault Architecture.md` | ✅ KEEP | Folder structure reference |
| `OneMind-OS-Brand-Architecture.md` | 🗑️ DELETE | Business content → move to 75-99 GE |
| `Cross-Pillar Linking Strategy.md` | ✅ KEEP | Implementation guide |
| `The Area Template.md` | ✅ KEEP | Practical template |

**Rename to:**
- `Vault-Structure.md` (clearer name)
- `Cross-Pillar-Routing.md` (clearer name)
- `Area-Template.md` (clearer name)

---

### 00.3 Development/

| Current File | Action | Reason |
|--------------|--------|--------|
| `Operating-Rhythm.md` | ✅ KEEP | Extract templates from CODEX |
| `HIVE-PROTOCOL.md` | ✅ KEEP | Capture methodology |
| `Content Flow Decision Tree.md` | ✅ KEEP | Routing logic |
| `FRAMEWORK-INTEGRATIONS.md` | ✅ KEEP | Tool integration guide |
| `7-PROBLEMS-SECOND-BRAINS.md` | 🗑️ DELETE | Marketing content → 75-99 GE |
| `Meta-Cog Governance.md` | 🗑️ DELETE | Duplicate of CODEX philosophy |

**Move to:**
- Rename folder to `00.3 Operating-System/`
- Keep only operational templates

---

## New Files to CREATE

### 1. `00 Framework/INDEX.md`

```markdown
# Framework Index

**Start here:** See [ONEMIND-CODEX.md](/ONEMIND-CODEX.md) at repository root for the complete system.

## Quick Navigation

### For Humans
- [Quick Start](QUICK-START.md) - 5-minute onboarding
- [Daily Rhythm](00.3 Operating-System/Daily-Rhythm.md)
- [Vault Structure](00.2 Architecture/Vault-Structure.md)

### For AI Agents
- [Agent Navigation Protocol](00.4 AI-Agent-Guides/Navigation-Protocol.md)
- [Memory Sync Protocol](00.4 AI-Agent-Guides/Memory-Sync-Protocol.md)
- [Routing Rules](00.4 AI-Agent-Guides/Routing-Rules.md)

### Reference
- [Vocabulary](00.5 Reference/Vocabulary.md)
- [Naming Conventions](00.5 Reference/Naming-Conventions.md)
```

---

### 2. `00 Framework/QUICK-START.md`

```markdown
# One Mind Quick Start

## 1. Read the Master Document (30 min)
[ONEMIND-CODEX.md](/ONEMIND-CODEX.md) - The complete system

## 2. Understand Your Three Paths (10 min)
- **HP (25-49)**: Personal optimization → PI Score
- **LE (50-74)**: Family & legacy → LI Score
- **GE (75-99)**: Business & wealth → WI Score
- **UI (00-24)**: The brain that coordinates all paths

## 3. Set Up Your Tools (1 hour)
See: [ONEMIND-CODEX.md Part 6: Implementation](/ONEMIND-CODEX.md#part-6-the-implementation)

## 4. Configure AI Agents (30 min)
See: [00.4 AI-Agent-Guides/](00.4 AI-Agent-Guides/)

## 5. Start Daily Rhythm (15 min/day)
See: [00.3 Operating-System/Daily-Rhythm.md](00.3 Operating-System/Daily-Rhythm.md)
```

---

### 3. `00.4 AI-Agent-Guides/Navigation-Protocol.md`

```markdown
# AI Agent Navigation Protocol

## First Read (Every Session)
1. [ONEMIND-CODEX.md](/ONEMIND-CODEX.md) - Read Parts 1-4 for context
2. Current user context from memory blocks
3. Today's priorities from Command Center

## When Routing Requests
1. Check [Routing Rules](Routing-Rules.md)
2. Use domain question filters:
   - "Does this improve ME?" → HP (25-49)
   - "Does this protect US?" → LE (50-74)
   - "Does this create VALUE?" → GE (75-99)
   - "How do we coordinate?" → UI (00-24)

## When Updating Knowledge
1. Read relevant section of ONEMIND-CODEX.md
2. Update memory blocks
3. Log to daily note

## File Locations
| Content Type | Location |
|--------------|----------|
| Master doctrine | `/ONEMIND-CODEX.md` |
| Agent configs | `00-24 SO/02 Agents/` |
| Protocols | `00-24 SO/04 Protocols/` |
| Daily operations | `00-24 SO/01 Command/` |
```

---

## Files to DELETE (Duplicates)

### From 00.1 Core-Doctrine/
- ❌ `ONE-MIND-FRAMEWORK-v1.0.md` (duplicate of CODEX)
- ❌ `Zeus-LIFE-OS-ARCHITECTURE.md` (duplicate of CODEX Part 6)

### From 00.2 Architecture/
- ❌ `OneMind-OS-Brand-Architecture.md` (move to 75-99 GE/Brand/)

### From 00.3 Development/
- ❌ `7-PROBLEMS-SECOND-BRAINS.md` (marketing content)
- ❌ `Meta-Cog Governance.md` (duplicate of CODEX philosophy)

**Total removed:** ~1,500 lines of duplicate content

---

## Migration Steps

### Step 1: Create New Structure (15 min)
```bash
cd "00-24 SO (Sovereign Operations)/00 Framework (Doctrine)"

# Create new folders
mkdir -p "00.3 Operating-System"
mkdir -p "00.4 AI-Agent-Guides"
mkdir -p "00.5 Reference"

# Create index files
touch INDEX.md
touch QUICK-START.md
touch "00.1 Core-Doctrine/README.md"
touch "00.4 AI-Agent-Guides/README.md"
touch "00.4 AI-Agent-Guides/Navigation-Protocol.md"
touch "00.4 AI-Agent-Guides/Memory-Sync-Protocol.md"
touch "00.4 AI-Agent-Guides/Routing-Rules.md"
```

### Step 2: Extract Templates from CODEX (30 min)
- Copy daily/weekly/quarterly templates from CODEX to `00.3 Operating-System/`
- Copy vocabulary from CODEX Appendix A to `00.5 Reference/Vocabulary.md`
- Copy naming conventions to `00.5 Reference/Naming-Conventions.md`

### Step 3: Rename & Reorganize (15 min)
```bash
# Rename for clarity
mv "00.2 Architecture/Vault Architecture.md" "00.2 Architecture/Vault-Structure.md"
mv "00.2 Architecture/Cross-Pillar Linking Strategy.md" "00.2 Architecture/Cross-Pillar-Routing.md"

# Move files
mv "00.3 Development/Operating-Rhythm.md" "00.3 Operating-System/"
mv "00.3 Development/HIVE-PROTOCOL.md" "00.3 Operating-System/"
```

### Step 4: Delete Duplicates (5 min)
```bash
# Remove duplicates
rm "00.1 Core-Doctrine/ONE-MIND-FRAMEWORK-v1.0.md"
rm "00.1 Core-Doctrine/Zeus-LIFE-OS-ARCHITECTURE.md"
rm "00.2 Architecture/OneMind-OS-Brand-Architecture.md"
rm "00.3 Development/7-PROBLEMS-SECOND-BRAINS.md"
rm "00.3 Development/Meta-Cog → One Mind Governance Layer.md"
```

### Step 5: Update Cross-References (30 min)
- Update all remaining docs to link to ONEMIND-CODEX.md
- Update KNOWLEDGE-BASE.md at root to point to new structure
- Update CLAUDE.md to reference new agent navigation protocol

---

## Verification Checklist

After consolidation:
- [ ] ONEMIND-CODEX.md is the ONLY place with complete philosophy/architecture
- [ ] All other docs either link to CODEX or contain implementation specifics
- [ ] AI agents have clear navigation protocol (00.4 AI-Agent-Guides/)
- [ ] Humans have quick start guide (QUICK-START.md)
- [ ] No duplicate content exists
- [ ] File count reduced by ~50%

---

## Result: Clean Information Architecture

**Before:**
- 15+ overlapping docs
- Unclear what's canonical
- Agents confused about which doc to trust

**After:**
- 1 master doc (ONEMIND-CODEX.md)
- Clear derivative hierarchy
- AI agents know exactly where to look
- Humans have fast navigation

**Total time investment:** ~2 hours
**Ongoing maintenance saved:** Massive (no more syncing multiple docs)

---

*This plan creates the single source of truth you need for both human and AI effectiveness.*
