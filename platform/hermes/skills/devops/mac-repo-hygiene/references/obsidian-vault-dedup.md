# Obsidian Vault Deduplication — Zeus's Vault Facts & Audit Pattern

> Reference only. Deleting or removing any folder named here needs Zeus to write "execute" or "override" for that specific deletion (SOUL rule 4). Quarantine and report; never delete on your own.

## Zeus's Canonical Vault

- **Primary vault:** `~/onemind/` on the Mac (Obsidian-registered, GitHub-connected)
- **GitHub:** `github.com/Zeus-Delacruz/onemind` — clean, active (last commit July 22 2026)
- **Structure:** Four life pillars — `00-24 SO`, `25-49 HP`, `50-74 LE`, `75-99 GE`
- **`onemind-2 /`** (trailing space) — abandoned migration attempt. Broken git, not registered in Obsidian, nothing newer than `onemind/`. Being deduped/deleted folder by folder.

## Which Vault Is Active?

```bash
cat ~/Library/Application\ Support/obsidian/obsidian.json | python3 -c \
  "import sys,json; d=json.load(sys.stdin); [print(v.get('path'), '| open:',v.get('open',False)) for v in d.get('vaults',{}).values()]"
```

## Byte-by-Byte Folder Audit Script (Single Source)

Run folder by folder — never all at once. Zeus reviews each result before deletion.

```python
import os, hashlib, glob

def hash_file(path):
    h = hashlib.md5()
    try:
        with open(path, 'rb') as f:
            while chunk := f.read(65536):
                h.update(chunk)
        return h.hexdigest()
    except:
        return None

# Build canonical vault hash map once
v1_base = '/Users/zeus/onemind'
v1_hashes = {}
for root, dirs, files in os.walk(v1_base):
    dirs[:] = [d for d in dirs if d not in ['.git', '.Trash']]
    for fn in files:
        if fn == '.DS_Store': continue
        full = os.path.join(root, fn)
        h = hash_file(full)
        if h:
            v1_hashes[h] = os.path.relpath(full, v1_base)

# Check one folder — use glob for trailing-space folder names
matches = glob.glob('/Users/zeus/onemind-2 /FOLDER_NAME*')
folder = matches[0] if matches else None

unique = []
exists = 0
for root, dirs, files in os.walk(folder):
    dirs[:] = [d for d in dirs if d != '.DS_Store']
    for fn in sorted(files):
        if fn == '.DS_Store': continue
        full = os.path.join(root, fn)
        rel = os.path.relpath(full, folder)
        size = os.path.getsize(full)
        h = hash_file(full)
        match = v1_hashes.get(h)
        if match:
            exists += 1
        else:
            unique.append((rel, size))

print(f"--- SUMMARY ---")
print(f"Total exists in onemind/: {exists}")
print(f"Total UNIQUE (not in onemind/): {len(unique)}")
if unique:
    print("\nUNIQUE FILES:")
    for f, s in unique:
        print(f"  🆕 {f} ({s/1024:.1f}K)")
else:
    print("\n✅ Zero unique files — entire folder is safe to delete.")
```

## Dual-Source Audit Script (External Drive — code/md vs onemind/, assets vs om-assets/)

Use when auditing `onemind (07-09-26)` on OM PRIMARY. One scan, two hash maps.

```python
import os, hashlib, glob

CODE_EXTS = {
    '.md', '.txt', '.py', '.js', '.ts', '.tsx', '.jsx', '.json', '.yaml', '.yml',
    '.sh', '.bash', '.go', '.rs', '.toml', '.env', '.conf', '.cfg', '.ini',
    '.html', '.css', '.scss', '.xml', '.sql', '.tf', '.hcl', '.vue', '.svelte',
    '.rb', '.java', '.c', '.cpp', '.h', '.cs', '.swift', '.dart', '.kt',
    '.r', '.ipynb', '.csv', '.tsv', '.graphql', '.proto'
}

def hash_file(path):
    h = hashlib.md5()
    try:
        with open(path, 'rb') as f:
            while chunk := f.read(65536):
                h.update(chunk)
        return h.hexdigest()
    except:
        return None

# Build mac onemind hash map
v1_base = '/Users/zeus/onemind'
v1_hashes = {}
for root, dirs, files in os.walk(v1_base):
    dirs[:] = [d for d in dirs if d not in ['.git', '.Trash']]
    for fn in files:
        if fn == '.DS_Store': continue
        full = os.path.join(root, fn)
        h = hash_file(full)
        if h:
            v1_hashes[h] = os.path.relpath(full, v1_base)

# Build om-assets hash map
assets_base = '/Volumes/OM PRIMARY/om-assets (07-22-26)'
assets_hashes = {}
if os.path.exists(assets_base):
    for root, dirs, files in os.walk(assets_base):
        dirs[:] = [d for d in dirs if d not in ['.git', '.Trash', '.fseventsd']]
        for fn in files:
            if fn.startswith('.'): continue
            full = os.path.join(root, fn)
            h = hash_file(full)
            if h:
                assets_hashes[h] = os.path.relpath(full, assets_base)

# Scan one folder
folder = '/Volumes/OM PRIMARY/onemind (07-09-26)/FOLDER_NAME'
unique_code = []
unique_assets = []
ok_code = ok_assets = 0

for root, dirs, files in os.walk(folder):
    dirs[:] = [d for d in dirs if not d.startswith('.')]
    for fn in sorted(files):
        if fn.startswith('.'): continue
        full = os.path.join(root, fn)
        rel = os.path.relpath(full, folder)
        ext = os.path.splitext(fn)[1].lower()
        size = os.path.getsize(full)
        h = hash_file(full)
        if not h: continue
        is_code = ext in CODE_EXTS
        if is_code:
            match = v1_hashes.get(h)
            if match: ok_code += 1
            else: unique_code.append((rel, size))
        else:
            match = assets_hashes.get(h)
            if match: ok_assets += 1
            else: unique_assets.append((rel, size))

print(f"--- SUMMARY ---")
print(f"Code/MD:  ✅ {ok_code} confirmed  |  🆕 {len(unique_code)} unique")
print(f"Assets:   ✅ {ok_assets} confirmed  |  🆕 {len(unique_assets)} unique")
if unique_code:
    print(f"\nUNIQUE CODE/MD:")
    for f,s in unique_code: print(f"  🆕 {f} ({s/1024:.1f}K)")
if unique_assets:
    print(f"\nUNIQUE ASSETS:")
    for f,s in unique_assets: print(f"  🆕 {f} ({s/1024:.1f}K)")
if not unique_code and not unique_assets:
    print(f"\n✅ Zero unique files — safe to delete.")
```

## Zeus's Workflow Preference

1. Go **folder by folder** — Zeus calls the folder name to advance (do NOT auto-advance)
2. For each folder: report **summary + unique files only** — NEVER print every confirmed-duplicate line
3. Zeus decides what to move or delete — **never delete without explicit approval**
4. When unique files found: move them with `rsync --ignore-existing` (handles merging into existing dirs)
5. After all folders: commit rescued unique files to `onemind/` GitHub

## Moving Unique Files

```bash
# Single folder (no conflicts expected)
cp -r '/Users/zeus/onemind-2 /FOLDER' '/Users/zeus/onemind/FOLDER'

# Merging into existing dir (use rsync to avoid overwriting)
rsync -av --ignore-existing '/Users/zeus/onemind-2 /FOLDER/' '/Users/zeus/onemind/FOLDER/'

# Assets to OM PRIMARY
mkdir -p '/Volumes/OM PRIMARY/om-assets (07-22-26)/FOLDER'
rsync -av '/source/FOLDER/' '/Volumes/OM PRIMARY/om-assets (07-22-26)/FOLDER/'
```

## 2026-07-29 Audit Summary (onemind-2/ full scan)

| Metric | Value |
|--------|-------|
| Files in onemind/ | 68,525 |
| Files in onemind-2/ | 62,014 |
| Truly unique to onemind-2/ | **7,569 files — 47.3MB** |
| Same path, different content | **0** — zero conflicts |
| In onemind-2/ by path, content in onemind/ | 49,335 files |

Biggest false-unique: `omos-personal-ai-infrastructure/` (62M) — already lives at `~/omos-personal-ai-infrastructure/` as a separate repo.

## onemind-2/ Folder-by-Folder Results (2026-07-29)

| Folder | Files | Unique | Action |
|--------|-------|--------|--------|
| Root files | 35 | 3 (`cleanup.md`, `openclaw-update-log-2026-05-08.md`, `openclaw-update-log-2026-05-13.md`) | Pending rescue |
| `Credentials/` | 12 | 0 | ✅ Safe to delete |
| `_RESET/` | 16 | 16 (ALL unique) | ✅ **MOVED** to `onemind/_RESET/` |
| `_system/` | 9 | 9 (ALL unique) | ✅ **MERGED** into `onemind/_system/` via rsync |
| `50-74 LE (Legacy Evolution)` | 100+ | 0 | ✅ Safe to delete |
| `25-49 HP (Holistic Performance)` | 1,410 | 0 | ✅ Safe to delete |
| `75-99 GE (Generational Empire)` | 662 | 0 | ✅ Safe to delete |
| `onemindcodex-v0.4` | 1,896 | 0 | ✅ Safe to delete |
| `00-24 SO (Sovereign Operations) 2` | 1,641 | 0 | ✅ Safe to delete |
| `00-24 SO (Sovereign Operations)` | 2,102 | 0 | ✅ Safe to delete |
| `99-INBOX/` | 22,756 | 1 (spam email, skip) | ✅ Safe to delete |
| `REVIEW/` | 22,555 | 1 (same spam email, skip) | ✅ Safe to delete |

## OM PRIMARY — onemind (07-09-26) Folder-by-Folder Results (2026-07-29)

**Drive:** `/Volumes/OM PRIMARY/` (space not dash)
**Total size:** 428GB — includes media, videos, PDFs not in Mac vault

| Folder | Size | Code/MD Unique | Asset Unique | Action |
|--------|------|----------------|--------------|--------|
| Root files | ~1MB | 3 (cleanup.md, openclaw logs) | 0 | Pending rescue to onemind/ |
| `ZLD - ID/` | 10MB | 0 | 11 (ALL — identity documents) | ✅ **MOVED** to `om-assets/ZLD - ID/` |
| `Credentials/` | 12MB | 2 (AWS keys .csv) + 3 assets (.pem, recovery PDF) | — | In progress |

**Remaining folders not yet audited in onemind (07-09-26):**
- `REVIEW/` (80GB), `75-99 GE/` (63GB), `25-49 HP/` (60GB), `Training Sessions (Long)/` (44GB)
- `References/` (28GB), `50-74 LE/` (10GB), `99-INBOX/` (3.6GB), `_rescued-inbox/` (2.4GB)
- `00-24 SO/` x2, `omos-personal-ai-infrastructure/`, `onemindcodex-v0.4/`
- `data-exports/`, `untitled folder/`, `_codex/`, `_RESET/`, `_system/`
- `jarvis-main 2/`, `mega-cycle-system/`, `onemind workflows/`

## Pitfalls

- `onemind-2 /` has a **trailing space** in the folder name — always quote carefully and use `glob` in Python
- `cp -r` into an existing directory fails silently on subdirs that already exist — use `rsync --ignore-existing` for merges
- The `_system/` folder already existed in onemind/ — `cp -rn` returned exit 1 even on success because some subdirs pre-existed; use rsync instead
- The omos-personal-ai-infrastructure subfolder in onemind-2/ shows 7,500 "unique" files but they're all from the PAI repo which lives separately at `~/omos-personal-ai-infrastructure/` — not truly missing from the canonical vault
- Always use the **summary-only** script pattern for large folders (1000+ files) — printing every file floods the context window
- External drive vault contains 428GB including media/videos not in Mac vault — do NOT compare sizes directly to Mac vault; they are not equivalent
- `OM PRIMARY` volume has a space: `/Volumes/OM PRIMARY/` — not `/Volumes/OM-PRIMARY/`
- Identity documents in `ZLD - ID/` are always unique (not in om-assets by default) — treat as irreplaceable, move before any deletion
- Building om-assets hash map on a large drive (07-22-26) takes 30-60s — do it once at top of script, reuse across all folder scans in the session
