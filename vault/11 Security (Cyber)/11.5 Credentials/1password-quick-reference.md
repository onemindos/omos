---
id: root-00-20260113225954
title: "1password quick reference"
type: reference
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/reference, status/active]
ai_summary: "1password quick reference"
---
# 1Password OneMind Integration - Quick Reference

## Current Status (as of 2026-01-13)
```
✅ 245 items organized
✅ 231 items with path tags (94.3%)
✅ 229 items with category tags (93.5%)
⚠️  14 items in Z Delete (untagged - archive candidates)
```

## Vault Structure
| Vault | Items | Purpose |
|-------|-------|---------|
| OneMind - Unified Intelligence | 63 | Tech, AI, infrastructure |
| OneMind - Holistic Performance | 103 | Personal, health, finance |
| OneMind - Generational Entrepreneurship | 48 | Business, revenue, partners |
| OneMind - Legacy Evolution | 14 | Family, legal, estate |
| Z Delete | 15 | Archive candidates |
| Personal | 2 | Private items |

## Tag Examples
```
✓ Correct: path:ui, category:it
✓ Correct: path:hp, category:finance
✓ Correct: path:le, category:legal
✓ Correct: path:ge, category:marketing

✗ Wrong: path:infrastructure (use category:infrastructure)
✗ Wrong: path:finance (use category:finance)
✗ Wrong: Only one tag (need both path: and category:)
```

## Common Commands

### Search/Query
```bash
# All UI domain items
op item list --tags "path:ui" --format json

# All finance credentials
op item list --tags "category:finance" --format json

# Untagged items
op item list --format json | jq '.[] | select(.tags == null or .tags == [])'
```

### Add/Update
```bash
# Update tags on existing item
op item edit [ITEM_ID] --tags "path:ui,category:it"

# Move item to different vault
op item move "Item Name" --current-vault "Old Vault" --destination-vault "OneMind - Unified Intelligence"
```

### Maintenance
```bash
# Find incomplete tags (only path, no category)
op item list --format json | jq '.[] | select(.tags | length == 1)'

# Count items per path
for tag in ui hp le ge; do
  echo "path:$tag: $(op item list --tags "path:$tag" --format json | jq length)"
done
```

## Path Tag Reference
| Tag | Domain | Purpose |
|-----|--------|---------|
| `path:ui` | Unified Intelligence (00-24) | AI agents, tech infrastructure, automations |
| `path:hp` | Holistic Performance (25-49) | Personal development, health, shopping |
| `path:le` | Legacy Evolution (50-74) | Family, legal, generational planning |
| `path:ge` | Generational Entrepreneurship (75-99) | Business, revenue, partner networks |

## Category Tags Quick List
Finance (48) • IT (43) • Infrastructure (31) • Intelligence (26) • Shopping (24) • Travel (11) • Health (10) • Partner Network (8) • Legal (7) • Marketing (6) • Communications (5) • Family (4) • MLM (2) • ID (2) • Content (2)

## Troubleshooting

**Tags not applying?**
- Use both tags in one command: `--tags "path:ui,category:it"`
- If in Z Delete vault, use web UI instead

**Item won't move?**
- Use `op item move` not `op item edit --vault`
- If fails, use 1Password.com web UI

**Can't find item?**
- Check if it's in Personal vault (not indexed by default)
- Use: `op item list --vault "Personal" --format json`

## Next Steps for Maintenance

- [ ] Delete/archive items in Z Delete vault after review
- [ ] Create Letta agent integration for credential queries
- [ ] Setup N8N workflows for auto-categorization
- [ ] Configure credential rotation for category:infrastructure items

---

For full documentation, see: `/00-24 UI (Unified Intelligence)/08 Infrastructure (Grid & APIs)/1PASSWORD-VAULT-RESTRUCTURING.md`