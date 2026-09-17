–---
id: library-readme
title: "OneMind Book Library — System Guide"
type: sop
status: active
created: 2026-03-29
updated: 2026-03-29
owner: zeus
tags: [type/sop, domain/so, status/active, system/library]
ai_summary: "Complete operating guide for the OneMind book library inventory system — SKU format, folder structure, processing workflow, frontmatter spec, and location tracking"
---

# OneMind Book Library

A structured, AI-native inventory system for tracking physical and digital books. Built for Obsidian, designed for batch processing with vision AI models.

---

## Quick Start

1. **To add a book manually:** Copy `_templates/book.md` → fill in frontmatter → save to the correct genre folder
2. **To add books from photos:** Upload image → AI identifies book → generates MD file → place in correct folder
3. **To find a book:** Use Obsidian search, Dataview queries, or browse by genre folder
4. **To check a location:** Open the location file (e.g., `locations/local.md`) for a Dataview-powered inventory

---

## Directory Structure

```
_library/
├── README.md                              ← This file
├── _templates/
│   └── book.md                            ← Book entry template
├── locations/
│   ├── local.md                           ← In-room / bedroom shelf
│   ├── attic.md                           ← Attic storage
│   └── desk.md                            ← Desk / workspace
├── non-fiction/
│   ├── business-finance/
│   │   └── becoming-your-own-banker/      ← One folder per book
│   │       ├── becoming-your-own-banker.md
│   │       ├── cover-front.jpg
│   │       └── cover-back.jpg
│   ├── self-help-personal-development/
│   ├── technology-science/
│   ├── philosophy-spirituality/
│   ├── health-wellness/
│   ├── biography-memoir/
│   ├── history-politics/
│   └── education-reference/
├── fiction/
│   ├── fantasy/
│   ├── science-fiction/
│   ├── thriller-mystery/
│   ├── literary/
│   └── romance/
├── reference/
│   ├── textbooks/
│   ├── manuals-guides/
│   └── encyclopedias-dictionaries/
└── children/
    ├── picture-books/
    ├── early-readers/
    └── young-adult/
```

### Folder-Per-Book Convention

Each book lives in its own subfolder under the genre directory:

```
genre-folder/
└── book-slug/
    ├── book-slug.md           ← Main entry file
    ├── cover-front.jpg        ← Front cover photo
    ├── cover-back.jpg         ← Back cover photo
    └── book-title.pdf         ← PDF version (if available)
```

- Folder name = slugified book title (lowercase-hyphens)
- `.md` filename matches the folder name
- Cover images are always `cover-front.jpg` and `cover-back.jpg`
- PDFs use a descriptive filename (title or source name)
- Additional photos or scans can be added to the folder as needed

**Adding new genre folders:** Create a new folder under the appropriate top-level category. Use `lowercase-hyphens`. No spaces.

---

## SKU System

Every book gets a unique internal SKU that encodes format and sequence.

### Format

```
BK-{FORMAT}-{SEQUENCE}
```

| Segment | Values | Meaning |
|---------|--------|---------|
| `BK` | Always `BK` | Book identifier prefix |
| `FORMAT` | `PHY` / `DIG` / `BTH` | Physical, Digital, or Both |
| `SEQUENCE` | `0001`–`9999` | Auto-incrementing sequence number |

### Examples

| SKU | Meaning |
|-----|---------|
| `BK-PHY-0001` | First physical book cataloged |
| `BK-DIG-0042` | 42nd digital book cataloged |
| `BK-BTH-0007` | 7th book with both physical and digital copies |

### SKU Rules

- SKUs are **permanent** — never reassign a retired SKU
- Sequence is **global** across all formats (not per-format)
- The SKU also serves as the `id` field in frontmatter
- To find the next available SKU: sort all existing books by sequence, take the highest + 1

### Finding Next SKU (Dataview)

```dataview
TABLE sku, title, format
FROM "_library"
WHERE type = "book"
SORT sku DESC
LIMIT 1
```

---

## Frontmatter Specification

Every book file must have this YAML frontmatter. Required fields are marked with `*`.

### Core Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | * | Same as SKU (e.g., `BK-PHY-0001`) |
| `title` | string | * | Full book title |
| `type` | string | * | Always `book` |
| `status` | string | * | `owned` / `wishlist` / `lent` / `lost` / `sold` / `donated` |
| `format` | string | * | `physical` / `digital` / `both` |
| `sku` | string | * | Internal SKU (same as `id`) |

### Book Details

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `author` | list | * | Author name(s) as array |
| `publisher` | string | | Publisher name |
| `published` | string | | Publication year or full date |
| `edition` | string | | Edition info (e.g., "2nd", "Revised") |
| `pages` | number | | Page count |
| `language` | string | | ISO 639-1 code (default: `en`) |
| `cover_type` | string | | `hardcover` / `paperback` / `spiral` / `board` |
| `isbn` | string | | ISBN-10 |
| `isbn13` | string | | ISBN-13 |
| `asin` | string | | Amazon ASIN |

### Classification

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `genre` | list | * | Primary genre(s) — maps to folder placement |
| `subgenre` | list | | More specific classification |
| `category` | string | * | Top-level: `non-fiction` / `fiction` / `reference` / `children` |
| `subject_tags` | list | | Topic keywords for search |
| `dewey` | string | | Dewey Decimal number (if known) |
| `lcc` | string | | Library of Congress Classification |

### Physical Location

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `location` | string | * | `local` / `attic` / `desk` / `digital` / `lent` |
| `shelf` | string | | Shelf identifier within location (e.g., "top-left", "row-2") |
| `condition` | string | | `new` / `like-new` / `good` / `fair` / `poor` |

### Digital Details (if format is digital or both)

| Field | Type | Description |
|-------|------|-------------|
| `digital_format` | string | `pdf` / `epub` / `mobi` / `kindle` / `audiobook` |
| `digital_platform` | string | `kindle` / `apple-books` / `audible` / `local` / `other` |
| `file_path` | string | Path to file if stored locally (relative to vault or `_docs/`) |

### Reading Status

| Field | Type | Description |
|-------|------|-------------|
| `read_status` | string | `unread` / `reading` / `completed` / `abandoned` / `re-reading` |
| `date_started` | date | When reading began |
| `date_finished` | date | When reading completed |
| `rating` | number | 0–5 (0 = unrated) |
| `review` | string | One-line review or link to full review note |

### Acquisition

| Field | Type | Description |
|-------|------|-------------|
| `acquired` | date | Date acquired |
| `acquired_from` | string | Where purchased (Amazon, bookstore, gift, etc.) |
| `price_paid` | number | Purchase price |
| `estimated_value` | number | Current estimated value |

### Media

| Field | Type | Description |
|-------|------|-------------|
| `cover_image` | string | Path to cover image (in `_images/` or URL) |
| `photos` | list | Paths to photos of the physical book |

### External Links

| Field | Type | Description |
|-------|------|-------------|
| `goodreads_url` | string | Goodreads page |
| `amazon_url` | string | Amazon product page |
| `openlibrary_url` | string | Open Library page |
| `worldcat_url` | string | WorldCat listing |
| `source_urls` | list | Any other reference URLs |

### Metadata

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `tags` | list | * | Standard vault tags + `type/book` |
| `related` | list | | Wiki-links to related vault notes |
| `ai_summary` | string | * | One-line book description (≤20 words) |
| `created` | date | * | Entry creation date |
| `updated` | date | * | Last modification date |

---

## Tag System

Every book gets these base tags plus classification tags:

```yaml
tags:
  - type/book
  - status/owned              # or wishlist, lent, sold, etc.
  - format/physical            # or digital, both
  - category/non-fiction       # or fiction, reference, children
  - genre/business-finance     # matches folder name
  - location/local             # or attic, desk, digital
  - read/unread               # or reading, completed, abandoned
```

This enables powerful Dataview queries and Obsidian tag-based browsing.

---

## File Naming Convention

Each book gets its own folder under the genre directory:

```
{category}/{genre-folder}/{book-slug}/
├── {book-slug}.md
├── cover-front.jpg
├── cover-back.jpg
└── optional.pdf
```

The slug (folder and file name) follows these rules:
- Lowercase with hyphens (vault convention)
- Drop leading "A" or "The" only if it makes the name clearer
- Include subtitle only if needed for disambiguation
- No author name in filename (it's in frontmatter)

Examples:
- `non-fiction/self-help-personal-development/atomic-habits/atomic-habits.md`
- `non-fiction/business-finance/becoming-your-own-banker/becoming-your-own-banker.md`
- `non-fiction/health-wellness/dr-sebi-the-cookbook/dr-sebi-the-cookbook.md`

---

## Processing Workflow: Image → Book Entry

### Step 1: Upload Image

Upload a photo of the book (cover, spine, or title page) to the AI agent.

### Step 2: AI Identification

The AI model will:
1. Read the title, author, and any visible details from the image
2. Cross-reference against book databases (Google Books, Open Library, Amazon)
3. Pull complete metadata: ISBN, publisher, page count, genre, description
4. Generate the full frontmatter + content sections

### Step 3: Agent Output

The AI generates a complete `.md` file with:
- Full frontmatter populated from database lookup
- AI-generated summary
- Genre classification and subject tags
- External links (Goodreads, Amazon, Open Library)
- Correct SKU assigned (next available in sequence)

### Step 4: File Placement

Create folder `_library/{category}/{genre-folder}/{book-slug}/` and save the `.md` file inside it. Add `cover-front.jpg` and `cover-back.jpg` to the same folder. If a PDF exists, copy it in too.

### Step 5: Verification

Mark `verified: true` in the source section once a human confirms accuracy.

---

## Batch Processing Protocol

When processing multiple book images at once:

1. Number each image or book clearly
2. AI processes each one and generates all MD files
3. Each file gets the next sequential SKU
4. Files are organized into correct genre folders
5. A summary report is generated showing: count, SKUs assigned, locations, any books that couldn't be identified

### Batch Naming

For photos of shelves with multiple books, name the image descriptively:
- `shelf-local-top-row.jpg`
- `attic-box-1.jpg`
- `desk-stack.jpg`

---

## Location System

### Current Locations

| Location ID | Name | Description |
|------------|------|-------------|
| `local` | In-Room / Bedroom | Main living space bookshelf |
| `attic` | Attic Storage | Long-term storage, less accessible |
| `desk` | Desk / Workspace | Active use, current reading |
| `digital` | Digital Only | No physical copy |
| `lent` | Lent Out | Currently with someone else |

### Adding Locations

1. Create a new `.md` file in `locations/` with the location ID as filename
2. Add a Dataview query that lists all books at that location
3. Update this README's location table

### Moving Books

When a book moves locations:
1. Update the `location` field in the book's frontmatter
2. Update the `location/` tag
3. Update `updated` date

---

## Obsidian Dataview Queries

### All Books

```dataview
TABLE author, format, location, read_status, rating
FROM "_library"
WHERE type = "book"
SORT title ASC
```

### Books by Location

```dataview
TABLE title, author, condition, shelf
FROM "_library"
WHERE type = "book" AND location = "local"
SORT shelf ASC, title ASC
```

### Unread Books

```dataview
TABLE title, author, genre, location
FROM "_library"
WHERE type = "book" AND read_status = "unread"
SORT acquired DESC
```

### Reading Stats

```dataview
TABLE length(rows) AS "Count"
FROM "_library"
WHERE type = "book"
GROUP BY read_status
```

### Books by Genre

```dataview
TABLE title, author, rating
FROM "_library"
WHERE type = "book" AND contains(genre, "business-finance")
SORT rating DESC
```

### Recently Added

```dataview
TABLE title, author, sku, format
FROM "_library"
WHERE type = "book"
SORT created DESC
LIMIT 20
```

### Collection Value

```dataview
TABLE sum(rows.price_paid) AS "Total Spent", sum(rows.estimated_value) AS "Total Value"
FROM "_library"
WHERE type = "book"
GROUP BY format
```

---

## AI Agent Instructions

### For Image Processing

When the user uploads a book image:

1. **Identify** the book from the image (title, author, edition)
2. **Look up** full metadata using the title/author:
   - Google Books API, Open Library API, Amazon product data
   - ISBN, publisher, page count, publication date, genre
3. **Determine** the next available SKU by checking existing entries
4. **Generate** the complete markdown file using `_templates/book.md`
5. **Classify** into the correct `category/genre-folder/`
6. **Include** all source URLs where data was found
7. **Set** `ai_summary` to a concise one-line description
8. **Set** `location` based on user context (ask if unclear)

### For Bulk Processing

- Process each book sequentially, assigning incrementing SKUs
- After all books are processed, give a summary table
- Flag any books that couldn't be confidently identified
- If a genre folder doesn't exist yet, create it

### Classification Guidelines

| If the book is about... | Category | Genre Folder |
|-------------------------|----------|-------------|
| Business, money, investing, economics | non-fiction | `business-finance` |
| Personal growth, habits, productivity | non-fiction | `self-help-personal-development` |
| Programming, AI, engineering, science | non-fiction | `technology-science` |
| Religion, philosophy, metaphysics | non-fiction | `philosophy-spirituality` |
| Fitness, nutrition, mental health | non-fiction | `health-wellness` |
| Life stories, memoirs | non-fiction | `biography-memoir` |
| History, politics, current events | non-fiction | `history-politics` |
| Academic, how-to, professional | non-fiction | `education-reference` |
| Fantasy, sci-fi, mystery, etc. | fiction | (appropriate subgenre) |
| Textbooks | reference | `textbooks` |
| Technical guides, cookbooks, manuals | reference | `manuals-guides` |
| Books for ages 0-17 | children | (appropriate subgenre) |

---

## Relationship to Other Vault Systems

| System | Relationship |
|--------|-------------|
| `_docs/` | If a book has a digital PDF/EPUB, the file lives in `_docs/` and `file_path` points to it |
| `_images/` | Cover photos and book images can be stored in `_images/` |
| `25-49 HP/26 Intelligence/` | Reading notes, takeaways, and reviews reference the book entry via wiki-link |
| `75-99 GE/` | Business books may have action items tracked as GE tasks |
| Dataview | Location files and dashboards use Dataview to query the library |

---

## Maintenance

### Monthly

- Verify location accuracy (spot-check 5 books)
- Update `read_status` for books started/finished
- Check for duplicate entries

### Quarterly

- Full location audit (count physical books vs entries)
- Update estimated values for valuable books
- Archive entries for books sold/donated (change status, don't delete)

---

*Part of the OneMind Codex vault. Created 2026-03-29.*


---

## Related

- [[25-49 HP (Holistic Performance)/27 Mind (Education & Learning)/27 Mind|27 Mind]]
