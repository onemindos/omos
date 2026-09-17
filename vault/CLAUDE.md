# OneMind __markdown Reorganization — Session State

## Source of Truth
`__onemind/` is the master vault. Four domain MAINs:
- `00-24 SO (Sovereign Operations) MAIN`
- `25-49 HP (Holistic Performance) MAIN`
- `50-74 LE (Legacy Evolution) MAIN`
- `75-99 GE (Generational Endeavor) MAIN`

## Current State of `__markdown/` Root

```
__onemind/                  ← SOURCE OF TRUTH (4 domain MAINs)
_byte-matched-dupes/        ← 60 files: byte-for-byte confirmed dupes of __onemind content. SAFE TO DELETE.
_loose-to-sort/             ← 258 files still needing work
  _byte-matched/            ← 42 files: byte-for-byte confirmed dupes of __onemind content. SAFE TO DELETE.
  (21 content-diff files)   ← same filename as __onemind but different content — need version check
  (237 truly unique files)  ← NOT in __onemind at all — MUST be routed to correct MAIN before any deletion
notion-import - MAIN/
obsidian-import - MAIN/
onemind-media-plr- MAIN/
```

## Pending Work (next session picks up here)

### Step 1 — Route 237 unique files in `_loose-to-sort/`
- Scan each file's content → assign to SO / HP / LE / GE domain
- Copy to proper MAIN subfolder
- After routing confirmed, mark originals for deletion

### Step 2 — Resolve 21 content-diff files in `_loose-to-sort/`
- Compare mod dates: `stat -f "%Sm" -t "%Y-%m-%d %H:%M" file`
- Keep newer version in __onemind, discard older root copy
- List: see `/tmp/loose_to_sort_audit.txt` (CONTENT_DIFF lines)

### Step 3 — Delete confirmed staging folders (need Zeus's go-ahead per file)
- `_byte-matched-dupes/` (60 files) — Zeus confirmed these are safe
- `_loose-to-sort/_byte-matched/` (42 files) — safe once Zeus confirms

### Step 4 — Process remaining __markdown MAINs
- `notion-import - MAIN/`
- `obsidian-import - MAIN/`
- `onemind-media-plr- MAIN/`

## Rules (DO NOT SKIP)
- **Never delete without explicit per-file/per-folder confirmation from Zeus**
- Always use `diff -rq` or `diff` for byte-level verification — no guessing
- When a filename appears multiple times in __onemind, check ALL matches (any identical = confirmed dupe)
- Use `rsync -a --ignore-existing` for safe merges into MAINs
- Domain routing: SO=ops/code/agents/infra, HP=health/performance/mind, LE=legacy AI/evolution, GE=ventures/generational

## Audit Files (in /tmp — regenerate if gone)
- `/tmp/loose_to_sort_audit.txt` — BYTE_MATCH | CONTENT_DIFF | UNIQUE tags for all 300 _loose-to-sort files
- `/tmp/onemind_mains_index.txt` — filename index of __onemind MAINs (29,582 entries)
