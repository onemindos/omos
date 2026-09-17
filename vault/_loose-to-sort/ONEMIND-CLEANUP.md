# One Mind Cleanup — Session Handoff Doc

> Last updated: 2026-06-16
> Purpose: Continue the cleanup/merge/organize operation in a new session without losing context.

---

## The Goal

Unify ALL of Zeus's data into ONE folder: `/Users/onemind/onemind/`

Every file goes into the pillar and domain it belongs to. No separate assets folder. No separate media folder. Credit report goes in `25-49 HP/30 Finance/`. Family photo goes in `50-74 LE/51 Covenant/`. Everything filed where it belongs.

### The Design (LOCKED)

```
/Users/onemind/onemind/
├── 00-24 SO (Sovereign Operations)/
├── 25-49 HP (Holistic Performance)/
├── 50-74 LE (Legacy Evolution)/
├── 75-99 GE (Generational Empire)/
├── _system/
└── .gitignore  ← git tracks .md + code only. Assets cloud-synced only.
```

- **Obsidian** opens `/Users/onemind/onemind/` as vault ✅
- **Git** tracks `.md`, `.sh`, `.py`, `.ts`, `.js`, `.json`, `.yaml` only ✅
- **Cloud sync** handles everything including assets ✅
- **RAG** (turbovec) points to `/Users/onemind/onemind/` as single brain ✅
- **Assets** (PDF, JPG, MP4, MP3) live IN the domain folder they belong to ✅

---

## Rules for This Cleanup

1. **Go one folder at a time** — no skipping ahead
2. **Explain WHY before any delete** — not just "it's safe"
3. **Verify with file counts + diff before deleting** — never assume
4. **Rescue unique files before deleting old versions** — always
5. **Zeus says "delete" — then delete** — never delete on assumption
6. **Vault is source of truth** — domain numbering follows current vault, not old folders

---

## Source Folders to Merge Into onemind/

| Folder | Size | Status |
|--------|------|--------|
| `onemind-assets/` | 241GB | ⏳ Not started |
| `om-primary/onemind-media/` | 1.4TB | ⏳ Not started |
| `om-primary/onemind-edu/` | 1.9TB | ⏳ Not started — Zeus wants to delete most of it |
| `Inbox/` | 10GB | ⏳ Not started |
| `Credentials/` | 12MB | ⚠️ Security risk — must go to 1Password first |
| `onemind-video/` | 665MB | ⏳ Separate repo `omos-video` — do not merge |

---

## Disk Space — Current State

| Volume | Total | Used | Free |
|--------|-------|------|------|
| Mac SSD | 926GB | ~788GB | ~138GB |
| OM PRIMARY (external) | 4.5TB | 3.0TB | 1.5TB |
| Time Machine | 931GB | 756GB | 175GB |

### Already Freed This Session (~164GB)
- `Library/Caches` cleared → ~24GB
- `Library/pnpm` cleared → ~6GB
- TimeMachine local snapshots (8) deleted → ~30GB
- `iPhone Videos - june-3-2026 2` (duplicate) deleted → 80GB
- `HP/LE/SO - old2` deleted (after rescue) → ~20GB
- `codex-media-*-_to-delete` deleted → ~4GB

---

## Inside `onemind/` — Cleanup Status

### ✅ Done
| Folder | Action |
|--------|--------|
| `HP - old2` | Rescued → deleted |
| `LE - old2` | Rescued → deleted |
| `SO - old2` | Rescued → deleted |
| `codex-media-00-24-_to-delete` | Deleted |
| `codex-media-50-74-_to-delete` | Deleted |
| `iPhone Videos - june-3-2026 2` | Confirmed duplicate → deleted |

### ⏳ In Progress — Reviewing One by One
Currently inside: `/Users/onemind/onemind/archive/`

| # | Folder | Size | Status |
|---|--------|------|--------|
| 1 | `archive/00-24 SO - media` | 208KB | 🔍 Reviewing — has .creds file + .lob files |
| 2 | `archive/00-24 SO - old3` | 20MB | ⏳ Not reviewed |
| 3 | `archive/00-24 SO - old2` | 688MB | ⏳ Not reviewed |
| 4 | `archive/25-49 HP - media` | 19MB | ⏳ Not reviewed |
| 5 | `archive/25-49 HP - old2` | 36MB | ⏳ Not reviewed |
| 6 | `archive/25-49 HP - old3` | 7MB | ⏳ Not reviewed |
| 7 | `archive/50-74 LE - media` | 12MB | ⏳ Not reviewed |
| 8 | `archive/50-74 LE - old2` | 13MB | ⏳ Not reviewed |
| 9 | `archive/50-74 LE - old3` | 1MB | ⏳ Not reviewed |
| 10 | `archive/75-99 GE - old3` | 5.1MB | ⏳ Not reviewed |
| 11 | `archive/75-99 GE - media` | 2.1GB | ⏳ Not reviewed |
| 12 | `archive/75-99 GE - old2` | 2.5GB | ⚠️ Has 13,019 files vs current GE 396 — needs rescue first |

### ⏳ Pending Review After archive/
| Folder | Size | Notes |
|--------|------|-------|
| `_legacy/` | 4.5GB | Old content folder |
| `move/` | 17GB | Staging folder with LE content |
| `GE - old2` (root level) | 19GB | Needs rescue — current GE nearly empty |
| `75-99 GE (Generational Endeavor) 7` | 25GB | Unknown — needs audit |
| `OneMind Academy (Education) 7` | 106GB | Zeus said delete most — needs review |
| `_docs` | 31GB | Unknown — needs audit |
| `33 Joy (Creativity & Play) 7` | 61GB | Needs routing to correct domain |
| `25 Identity 7` | 38GB | Needs routing |
| `26 Mind (Books & Courses) 7` | 34GB | Needs routing |
| `Zeus Delacruz Enterprise` | 33GB | Needs routing → GE |
| `-media` folders (SO/HP/LE/GE) | ~74GB | Need routing into correct pillar domains |
| `Inbox/` | 10GB | Staging — needs processing |
| `Review/` | 567MB | Needs audit |
| `_unsorted-pdfs/` | 412MB | Needs routing |

---

## Key Decisions Already Made

- Vault = `/Users/onemind/onemind/` → `github.com/Zeus-Delacruz/onemind` (private, Zeus-Delacruz personal)
- `.gitignore` is live — git tracks `.md` + code only
- `onemind-video/` → separate repo `omos-video`, NOT merged into vault
- `onemind-edu/` → delete most, keep only what's active
- `Credentials/` → 1Password first, then delete from disk
- Vault domain numbering is source of truth (not old folder numbers)
- RAG engine: turbovec (installed, ref: `/Users/onemind/onemind/turbovec.md`)

---

## Vault Domain Map (Current — Source of Truth)

### SO (00-24)
```
00 Framework (The Constitution)
01 The Fabric (Node Registry)
02 Intelligence & Agents (The Brain)
03 Infrastructure (The Body)
04 Security (The Shield)
05 Operations (The Controls)
05 Stack (Arsenal Ecosystem)
06 Architecture (The Blueprint)
07 Development (The Forge)
08 Research (The Lab)
```

### HP (25-49)
```
25 Identity (Self, Purpose & Spirit)
26 Spirituality (Metaphysics)
27 Mind (Education & Learning)
28 Body (Health & Wellness)
29 Relationships (Social)
30 Finance (Money)
31 Mastery (Skills & Training)
32 Rhythm (Routines & Productivity)
33 Sovereignty (Legal)
34 Joy (Recreation & Travel)
```

### LE (50-74)
```
50 Home (Sanctuary & Property)
51 Covenant (Relationship & Union)
52 Offspring (Parenting & Education)
53 Parents (Elder Care & Support)
54 Heritage (History & Traditions)
55 Health (Family Wellness)
56 Estate (Family Wealth & Legal)
57 Community (Extended Network)
```

### GE (75-99)
```
75 Zeus Delacruz (Brand & Empire)
76 Real Estate (Landlord & Property)
77 Acquisitions (Bought Businesses)
78 Holdings (Legal & Entity Structure)
79 Investments (Wealth & Assets)
80 Innovation (New Ventures Pipeline)
81 Strategy (Portfolio Intelligence)
82 Partnerships & Capital (Network & Deals)
83 Philanthropy (Giving)
```

---

## How to Continue in New Session

1. Open this file first — it's the full state of play
2. Tell Legacy: "Continue the onemind cleanup from ONEMIND-CLEANUP.md"
3. Pick up at Item 1 of archive/ — the `.creds` file decision is still open
4. Work one folder at a time, update this doc as items are completed

---

*Legacy — One Mind System*
