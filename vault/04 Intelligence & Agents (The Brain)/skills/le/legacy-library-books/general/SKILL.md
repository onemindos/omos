---
name: legacy-library
description: Process book cover photos (HEIC/JPG from iPhone) into structured library entries in the One Mind Codex vault. Use when you need to catalog books, convert cover images, create book markdown entries, find next SKU, audit library entries, or batch-process a shelf of books.
allowed-tools:
  - Bash(sips:*)
  - Bash(cp:*)
  - Bash(mv:*)
  - Bash(mkdir:*)
  - Bash(grep:*)
  - Bash(find:*)
  - Bash(ls:*)
  - Bash(cat:*)
  - Bash(wc:*)
  - Bash(sort:*)
  - Bash(bash:*)
  - Read
  - Write
requires: [read_file, list_directory, write_file]
offline: true
agent: any
---

# Library Skill — Book Cataloging from Photos

Process iPhone photos of book covers into structured, searchable library entries.

## Overview

```
iPhone HEIC photos → Convert to JPG → Vision AI identifies book →
Create folder + markdown entry → Place covers → Validate
```

## Paths

| Path                          | Purpose                                         |
| ----------------------------- | ----------------------------------------------- |
| `_library/`                   | Library root                                    |
| `_library/_staging/`          | Drop zone for raw images (HEIC/JPG from iPhone) |
| `_library/_templates/book.md` | Book entry template                             |
| `_library/_images/covers/`    | Cover image archive                             |
| `_library/non-fiction/`       | Non-fiction genre folders                       |
| `_library/fiction/`           | Fiction genre folders                           |
| `_library/reference/`         | Reference material                              |
| `_library/children/`          | Children's books                                |
| `_library/locations/`         | Physical location files (Dataview-powered)      |

## Scripts

Helper scripts live in `.claude/skills/library/scripts/`:

| Script             | Purpose                                 |
| ------------------ | --------------------------------------- |
| `heic-convert.sh`  | Convert HEIC → JPG (macOS sips)         |
| `process-book.sh`  | End-to-end: mkdir → place covers        |
| `next-sku.sh`      | Find next available SKU                 |
| `audit-library.sh` | Audit for missing fields, orphan images |

---

## Step-by-Step Processing Workflow

### Step 1: Receive Images

User provides HEIC or JPG photos from iPhone. Usually in pairs (front cover + back cover).

**If HEIC:** Convert to JPG first using macOS `sips`:

```bash
# Single file
sips -s format jpeg input.HEIC --out output.jpg

# Batch convert all HEIC in staging
for f in _library/_staging/*.HEIC; do
  sips -s format jpeg "$f" --out "${f%.HEIC}.jpg"
done
```

**If already JPG:** No conversion needed — proceed to identification.

### Step 2: Identify the Book (Vision AI)

View the cover image(s) and extract:

| Field          | Source                  | Required |
| -------------- | ----------------------- | -------- |
| `title`        | Front cover             | Yes      |
| `author`       | Front/back cover        | Yes      |
| `publisher`    | Spine or copyright page | No       |
| `edition`      | Title page or copyright | No       |
| `isbn13`       | Back cover barcode      | No       |
| `isbn`         | Back cover              | No       |
| `pages`        | Back cover or look up   | No       |
| `cover_type`   | Visual inspection       | Yes      |
| `condition`    | Visual inspection       | Yes      |
| `description`  | Back cover blurb        | Yes      |
| `subject_tags` | Back cover / knowledge  | Yes      |

**Fallback:** If cover is unclear, search Google Books or Open Library by visible title/author.

### Step 3: Determine Classification

| Field      | How to Decide                                          |
| ---------- | ------------------------------------------------------ |
| `category` | `non-fiction` / `fiction` / `reference` / `children`   |
| `genre`    | Map to existing genre folder (see Genre Folders below) |
| `subgenre` | More specific tags                                     |

### Step 4: Get Next SKU

```bash
# Find highest existing SKU
grep -rh "^sku:" _library/ 2>/dev/null | sort -t'-' -k3 -n | tail -1
# → Next SKU: increment the number by 1
```

SKU format: `BK-{FORMAT}-{SEQUENCE}`

| Format Code | Meaning                   |
| ----------- | ------------------------- |
| `PHY`       | Physical book             |
| `DIG`       | Digital only              |
| `BTH`       | Both physical and digital |

Sequence is **global** — always increment from the last highest, regardless of format.

### Step 5: Create Book Folder

```bash
# Pattern: _library/{category}/{genre}/{book-slug}/
mkdir -p "_library/non-fiction/business-finance/the-book-title/"
```

Folder name = slugified book title (lowercase, hyphens, no special chars).

### Step 6: Place Cover Images

```bash
# Copy/move covers into the book folder
cp staging/IMG_XXXX.jpg "_library/non-fiction/genre/book-slug/cover-front.jpg"
cp staging/IMG_XXXY.jpg "_library/non-fiction/genre/book-slug/cover-back.jpg"
```

- Front cover is always `cover-front.jpg`
- Back cover is always `cover-back.jpg`
- If only one image, it's `cover-front.jpg`

### Step 7: Create Markdown Entry

Create `book-slug.md` inside the book folder with full frontmatter.

**Frontmatter template:**

```yaml
---
id: BK-PHY-XXXX
title: "Full Book Title"
type: book
status: owned
format: physical
isbn: ""
isbn13: ""
asin: ""
sku: BK-PHY-XXXX

author:
  - Author Name
publisher: Publisher Name
published: "2020"
edition: "1st"
pages: 0
language: en
cover_type: softcover

genre:
  - genre-folder-name
subgenre:
  - specific-topic
category: non-fiction
subject_tags:
  - tag1
  - tag2
dewey: ""
lcc: ""

location: local
shelf: ""
condition: good

digital_format: ""
digital_platform: ""
file_path: ""

read_status: unread
date_started: ""
date_finished: ""
rating: 0
review: ""

acquired: "2026-04-01"
acquired_from: ""
price_paid: 0.00
estimated_value: 0.00

cover_image: "./cover-front.jpg"
photos:
  - "./cover-front.jpg"
  - "./cover-back.jpg"

goodreads_url: ""
amazon_url: ""
openlibrary_url: ""
worldcat_url: ""
source_urls: []

tags:
  - type/book
  - status/owned
  - format/physical
  - category/non-fiction
  - genre/genre-folder-name
  - location/local
  - read/unread
related: []
ai_summary: "One-line book description — max 20 words."
created: 2026-04-01
updated: 2026-04-01
---
```

**Body template:**

```markdown
# Book Title

## Covers

| Front                             | Back                            |
| --------------------------------- | ------------------------------- |
| ![Front Cover](./cover-front.jpg) | ![Back Cover](./cover-back.jpg) |

> [!info] Metadata
> **Author**: Author Name
> **Edition**: 1st Edition (2020)
> **Publisher**: Publisher Name
> **Pages**: N/A | **Cover**: Softcover | **Condition**: Good
> **SKU**: BK-PHY-XXXX | **Location**: local
> **Links**: [Google Books](https://www.google.com/search?tbm=bks&q=isbn:XXXXX) | [Amazon](https://www.amazon.com/s?k=XXXXX)

## Description

Back cover blurb or AI-generated summary.

## Key Topics

- Topic 1
- Topic 2
- Topic 3

## Details

- **ISBN-13**: XXXXX
- **Publisher**: Publisher Name

## Reading Status

- [x] Cataloged
- [ ] Reading
- [ ] Completed
- **Rating**: --/5

## Notes

<!-- Personal notes, highlights, or why you own this book -->

---

_Part of the One Mind Library. See [[_library/DASHBOARD|Library Dashboard]] for full collection._

## Source Information

- **Data Source:** Image scan (IMG_XXXX, IMG_XXXY)
- **Verified:** false
- **Last Checked:** 2026-04-01
```

### Step 8: Validate Entry

Check that:

- [ ] Frontmatter has all required fields (`id`, `title`, `type`, `status`, `format`, `sku`, `author`, `genre`, `category`, `location`, `tags`, `ai_summary`, `created`, `updated`)
- [ ] `id` matches `sku`
- [ ] Cover images exist and are referenced correctly
- [ ] Folder name matches the `.md` filename
- [ ] Genre folder exists under the correct category
- [ ] Tags follow `namespace/value` pattern
- [ ] `ai_summary` is present and ≤20 words

---

## Genre Folders (Current)

### Non-Fiction

`biography` · `biography-memoir` · `business` · `business-finance` · `economics` · `education-reference` · `health` · `health-wellness` · `history` · `history-politics` · `leadership` · `memoir` · `military` · `philosophy` · `philosophy-spirituality` · `politics-history` · `politics-media` · `psychology-relationships` · `relationships` · `science-nature` · `self-help` · `self-help-personal-development` · `social-sciences` · `spirituality` · `sports` · `technology-science`

### Fiction

`fantasy` · `historical-fiction` · `literary` · `literary-fiction` · `romance` · `science-fiction` · `thriller` · `thriller-mystery`

### Children

`picture-books` · `early-readers` · `young-adult`

### Reference

`textbooks` · `manuals-guides` · `encyclopedias-dictionaries`

**New genre?** Create a new folder under the correct category. Use `lowercase-hyphens`.

---

## Batch Processing Pattern

When processing multiple books from a photo session:

1. **Convert all HEIC → JPG** in `_staging/`
2. **View images sequentially** — identify each book
3. **Create all folders + markdown** files
4. **Generate a copy script** to move covers from staging to book folders
5. **Run the copy script**
6. **Verify** all entries

### Copy Script Pattern

```bash
#!/bin/bash
LIB=/Users/zeuslegacy/Documents/onemind-codex/_library
STG=$LIB/_staging

cp "$STG/IMG_3401.jpg" "$LIB/non-fiction/self-help/book-slug/cover-front.jpg"
cp "$STG/IMG_3402.jpg" "$LIB/non-fiction/self-help/book-slug/cover-back.jpg"
# ... repeat for each book
echo "COPIES_DONE"
```

---

## Location Values

| Value     | Meaning                    |
| --------- | -------------------------- |
| `local`   | In-room / bedroom shelf    |
| `attic`   | Attic storage              |
| `desk`    | Desk / workspace           |
| `digital` | Digital-only (no physical) |
| `lent`    | Currently lent to someone  |

---

## Useful Commands

```bash
# Find next SKU
grep -rh "^sku:" _library/ 2>/dev/null | sort -t'-' -k3 -n | tail -1

# Count all books
grep -rl "^type: book" _library/ 2>/dev/null | wc -l

# List books missing ai_summary
grep -rL "^ai_summary:" _library/**/*.md 2>/dev/null

# List books by genre
grep -rh "^genre:" _library/ 2>/dev/null | sort | uniq -c | sort -rn

# Convert single HEIC
sips -s format jpeg input.HEIC --out output.jpg

# Batch convert HEIC in staging
for f in _library/_staging/*.HEIC; do sips -s format jpeg "$f" --out "${f%.HEIC}.jpg"; done
```

---

## Installation Checklist

### System Check

- [ ] `_library/` directory exists at vault root
- [ ] `_library/_templates/book.md` template accessible
- [ ] `_library/_staging/` directory exists for image drops
- [ ] `sips` command available (macOS built-in)
- [ ] Vision AI model available (for reading book covers)

### Context Loading

1. Read `_library/README.md` — full library system guide
2. Read `_library/_templates/book.md` — frontmatter template
3. Read `_library/DASHBOARD.md` — Dataview queries for collection stats
4. Check current highest SKU: `grep -rh "^sku:" _library/ | sort -t'-' -k3 -n | tail -1`

### Skill Activation

1. Read this `SKILL.md`
2. Internalize the 8-step processing workflow
3. Note current genre folders (don't memorize — check when needed)
4. Confirm understanding of SKU format: `BK-{FORMAT}-{SEQUENCE}`

### Script Check

```bash
ls .claude/skills/library/scripts/
# Expected: heic-convert.sh, process-book.sh, next-sku.sh, audit-library.sh
```

---

## Verification

### Required

- [ ] SKILL.md readable and loaded
- [ ] Agent can describe the 8-step workflow in one sentence
- [ ] Agent knows the SKU format: `BK-{FORMAT}-{SEQUENCE}`
- [ ] Agent knows the folder-per-book convention
- [ ] Agent can find the next available SKU
- [ ] Agent can list the four top-level categories (non-fiction, fiction, reference, children)

### Functional Tests

**Test 1: SKU Lookup**

- [ ] Run: `grep -rh "^sku:" _library/ | sort -t'-' -k3 -n | tail -1`
- [ ] Agent correctly identifies the next available SKU number

**Test 2: HEIC Conversion**

- [ ] Agent can produce the correct `sips` command for HEIC → JPG conversion
- [ ] Test: `sips -s format jpeg test.HEIC --out test.jpg`

**Test 3: Folder Creation**

- [ ] Given a book "The Art of War" (non-fiction, military), agent would create:
  - Folder: `_library/non-fiction/military/the-art-of-war/`
  - Files: `the-art-of-war.md`, `cover-front.jpg`, `cover-back.jpg`

**Test 4: Frontmatter Completeness**

- [ ] Agent can produce a complete frontmatter block with all required fields
- [ ] Required fields: `id`, `title`, `type`, `status`, `format`, `sku`, `author`, `genre`, `category`, `location`, `tags`, `ai_summary`, `created`, `updated`

### Result

- ✅ All checks pass = skill operational
- ⚠️ Optional checks fail = reduced capability (e.g., sips not available on non-macOS)
- ❌ Required check fails = not ready
