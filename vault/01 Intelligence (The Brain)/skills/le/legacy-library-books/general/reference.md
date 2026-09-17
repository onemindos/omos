# Library Skill — Quick Reference

## SKU Format
```
BK-PHY-0001   (physical)
BK-DIG-0042   (digital)
BK-BTH-0007   (both)
```

## Find Next SKU
```bash
grep -rh "^sku:" _library/ 2>/dev/null | sort -t'-' -k3 -n | tail -1
```

## Convert HEIC → JPG
```bash
sips -s format jpeg input.HEIC --out output.jpg
```

## Folder Structure
```
_library/{category}/{genre}/{book-slug}/
  ├── book-slug.md
  ├── cover-front.jpg
  └── cover-back.jpg
```

## Categories
- `non-fiction` — 26 genre folders
- `fiction` — 8 genre folders
- `children` — 3 genre folders
- `reference` — 3 genre folders

## Required Frontmatter Fields
`id` · `title` · `type` · `status` · `format` · `sku` · `author` · `genre` · `category` · `location` · `tags` · `ai_summary` · `created` · `updated`

## Required Tags
`type/book` · `status/{status}` · `format/{format}` · `category/{category}` · `genre/{genre}` · `location/{location}` · `read/{read_status}`

## Scripts
```bash
# Convert HEIC images
bash .claude/skills/library/scripts/heic-convert.sh [file|directory]

# Create book folder + place covers
bash .claude/skills/library/scripts/process-book.sh <category> <genre> <slug> <front> [back]

# Get next SKU
bash .claude/skills/library/scripts/next-sku.sh [PHY|DIG|BTH]

# Full audit
bash .claude/skills/library/scripts/audit-library.sh
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
