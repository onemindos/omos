---
id: root-00-20260113225954
title: "1PASSWORD VAULT RESTRUCTURING"
type: reference
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [type/reference, status/active]
ai_summary: "1PASSWORD VAULT RESTRUCTURING"
---
# 1Password Vault Restructuring & OneMind Integration
**Project Status: COMPLETE (94% Tagged)**
**Last Updated:** January 13, 2026

---

## Executive Summary

Successfully restructured 1Password from 30+ isolated vaults into 4 OneMind-aligned domain vaults with comprehensive tag-based organization. This enables credential access via tags, supports multi-agent queries, and aligns with the OneMind framework (UI/HP/LE/GE domains).

**Final Status:**
- ✅ **245 total items** across 6 vaults
- ✅ **231 items tagged with path:** (94.3%)
- ✅ **229 items tagged with category:** (93.5%)
- ✅ **225+ items fully tagged** (both path + category)
- ⚠️ **14 items in Z Delete** vault remain untagged (archive candidates)

---

## Architecture

### OneMind Domain Vaults

```
1Password Vaults (New Structure)
├── OneMind - Unified Intelligence (UI) → 63 items
│   └── Categories: it, infrastructure, communications, Intelligence
├── OneMind - Holistic Performance (HP) → 103 items
│   └── Categories: finance, health, shopping, travel, learning, family
├── OneMind - Legacy Evolution (LE) → 14 items
│   └── Categories: legal, family, finance
├── OneMind - Generational Entrepreneurship (GE) → 48 items
│   └── Categories: finance, marketing, content, mlm, partnernetwork, it
├── Z Delete (Archive) → 15 items
│   └── Categories: untagged (candidates for deletion)
└── Personal → 2 items
    └── For items you want outside the system
```

### Tag Structure

All items use a **2-tag system**: `path:` + `category:`

#### Path Tags (Domain Classification)
```
path:ui          → Unified Intelligence (00-24) - AI, tech, infrastructure
path:hp          → Holistic Performance (25-49) - Personal, health, wealth
path:le          → Legacy Evolution (50-74) - Family, legal, estate
path:ge          → Generational Entrepreneurship (75-99) - Business, revenue
path:archive     → Archive items (in Z Delete vault - candidates for deletion)
```

#### Category Tags (Functional Classification)
```
category:finance          → Financial accounts, banking, investments (48 items)
category:it              → Software tools, dev platforms, services (43 items)
category:infrastructure  → Core systems, hosting, domain management (31 items)
category:Intelligence    → AI services, LLM APIs, knowledge tools (26 items)
category:shopping        → Commerce, retail, marketplaces (24 items)
category:travel          → Airlines, hotels, travel services (11 items)
category:health          → Medical, wellness, fitness (10 items)
category:partnernetwork  → Partners, affiliates, integrations (8 items)
category:legal           → Legal services, compliance, documents (7 items)
category:marketing       → Marketing tools, analytics, campaigns (6 items)
category:communications  → Email, messaging, telecom (5 items)
category:family          → Family accounts, household management (4 items)
category:mlm             → Multi-level marketing platforms (2 items)
category:id              → Identity documents, credentials (2 items)
category:content         → Content management, publishing (2 items)
```

---

## Current Status Details

### Vault Distribution
| Vault | Count | % of Total | Status |
|-------|-------|-----------|--------|
| OneMind - Holistic Performance | 103 | 42.0% | ✅ Complete |
| OneMind - Unified Intelligence | 63 | 25.7% | ✅ Complete |
| OneMind - Generational Entrepreneurship | 48 | 19.6% | ✅ Complete |
| Z Delete | 15 | 6.1% | ⚠️ Untagged |
| OneMind - Legacy Evolution | 14 | 5.7% | ✅ Complete |
| Personal | 2 | 0.8% | ✅ Manual |
| **TOTAL** | **245** | **100%** | — |

### Tagging Completion
- **Path Tags Applied:** 231/245 (94.3%)
- **Category Tags Applied:** 229/245 (93.5%)
- **Fully Tagged (both):** 225+/245 (91.8%)
- **Untagged:** 14 items (5.7%) - all in Z Delete vault

### Path Distribution
```
path:hp   = 104 items (42.4%)
path:ui   = 64 items (26.1%)
path:ge   = 49 items (20.0%)
path:le   = 14 items (5.7%)
TOTAL     = 231 items
```

### Untagged Items (14 - All in Z Delete vault)
These are archive candidates awaiting deletion:
1. Agora
2. Coda
3. Cursor
4. Databox Affiliate
5. Envato
6. GoDaddy (Need to Call)
7. Miro
8. Notion
9. PluginsforWP
10. Post Flow
11. Replit
12. Robinhood (Recover?)
13. Send Certified Mail
14. Trello -> Move to Clickup?

---

## How It Works

### For AI Agents (Letta & Automations)

1. **Query by Domain Path:**
   ```bash
   op item list --tags "path:ui" --format json
   # Returns all Unified Intelligence credentials
   ```

2. **Query by Category:**
   ```bash
   op item list --tags "category:finance" --format json
   # Returns all financial accounts
   ```

3. **Combined Queries:**
   ```bash
   op item list --tags "path:ge" --tags "category:finance" --format json
   # Returns GE domain financial accounts
   ```

### For Manual Access

1. **1Password.com Web UI** - Use tag filters in the search/filter interface
2. **1Password Apps** - Search by tag or browse by vault
3. **Tag Filtering** - Click "All Items" → Filter by path or category tag

---

## Operations Guide

### Adding New Items

When adding new credentials:

1. **Determine Domain (Path Tag):**
   - UI (Unified Intelligence) → `path:ui`
   - HP (Holistic Performance) → `path:hp`
   - LE (Legacy Evolution) → `path:le`
   - GE (Generational Entrepreneurship) → `path:ge`

2. **Determine Category:**
   - Use existing categories when possible
   - Create new category only if no suitable category exists
   - Naming: lowercase, hyphenated (e.g., `category:partner-network`)

3. **Apply Tags via CLI:**
   ```bash
   op item create \
     --vault "OneMind - [Domain]" \
     --tags "path:xx,category:yyy" \
     [other flags]
   ```

### Moving Items Between Vaults

```bash
op item move "Item Title" \
  --current-vault "Old Vault" \
  --destination-vault "OneMind - Unified Intelligence"
```

**Note:** This creates a copy in destination and deletes from source. The item gets a new ID.

### Updating Tags on Existing Items

```bash
op item edit [ITEM_ID] --tags "path:ui,category:it"
```

**Important:** The `--tags` flag **replaces** all tags (doesn't accumulate if you apply both tags in one command).

### Batch Tag Updates

For multiple items in a domain:

```bash
# Get all UI items and add/update a category
op item list --tags "path:ui" --format json | jq -r '.[].id' | while read id; do
  op item edit "$id" --tags "path:ui,category:it"
done
```

### Finding Items

```bash
# By path domain
op item list --tags "path:ui" --format json

# By category
op item list --tags "category:finance" --format json

# By vault
op item list --vault "OneMind - Holistic Performance" --format json

# Untagged items
op item list --format json | jq '.[] | select(.tags == null or .tags == [])'

# Items with incomplete tags
op item list --format json | jq '.[] | select(.tags | length) == 1'
```

---

## Known Limitations & Workarounds

### 1. CLI Issues with Specific Vaults

**Problem:** Some vaults (especially "Z Delete") consistently fail with CLI tag/move operations.

**Workaround:** Use 1Password.com web UI for manual tagging:
1. Go to item → Click "Edit"
2. Scroll to Tags section
3. Add tags separated by commas

### 2. Item Move Creates New ID

**Problem:** `op item move` creates a copy with a new ID and deletes the original.

**Workaround:** This is expected behavior. IDs stored in other systems will need updating if referenced.

### 3. Tag Accumulation Issue (RESOLVED)

**Previous Issue:** `op item edit --tags` used to accumulate tags instead of replacing them.

**Resolution:** Always apply both tags (`path:` and `category:`) in a single command to avoid duplication.

### 4. Vault Rename Limitation

**Problem:** 1Password CLI doesn't support vault renaming.

**Workaround:** Rename vaults via 1Password.com web UI:
1. Settings → Vaults
2. Click vault menu → Rename

---

## Future Enhancements

### Phase 2 - AI Integration
- [ ] Configure Letta agents to query credentials by tag
- [ ] Create approval workflows for sensitive infrastructure access
- [ ] Set up credential rotation reminders for critical items

### Phase 3 - Automation
- [ ] N8N workflows to automatically categorize new items
- [ ] Discord notifications for new credentials added to specific tags
- [ ] Automated sync of credential metadata to vault

### Phase 4 - Governance
- [ ] Document access policies by path domain
- [ ] Set up encryption for highest-security categories
- [ ] Create rotation schedule for category:infrastructure items

---

## Maintenance Checklist

### Monthly
- [ ] Review items in Z Delete vault - permanently delete if confirmed
- [ ] Check for miscategorized items (review by category)
- [ ] Verify all new items have both tags

### Quarterly
- [ ] Audit category usage - consolidate if categories underused
- [ ] Review path distribution (should be roughly 40/25/20/5 for HP/UI/GE/LE)
- [ ] Cleanup old/duplicate vaults

### Annually
- [ ] Full audit of all 245+ items
- [ ] Rotate passwords for critical infrastructure
- [ ] Update documentation

---

## Scripts & Commands Reference

### Verification Script
```bash
#!/bin/bash
echo "=== 1PASSWORD TAG VERIFICATION ==="
total=$(op item list --format json | jq length)
echo "Total items: $total"

tagged=$(op item list --format json | jq '[.[] | select(.tags and .tags != [])] | length')
echo "Items with tags: $tagged"

echo ""
echo "PATH TAGS:"
for tag in ui hp le ge; do
  count=$(op item list --tags "path:$tag" --format json 2>/dev/null | jq length)
  echo "  path:$tag = $count items"
done

echo ""
echo "SAMPLE CATEGORIES:"
op item list --format json | jq '[.[] | .tags[]? | select(startswith("category:"))] | group_by(.) | map({tag: .[0], count: length}) | sort_by(.count) | reverse | .[0:5] | .[]'
```

### Find Untagged Items
```bash
op item list --format json | jq -r '.[] | select(.tags == null or .tags == []) | "\(.title) (Vault: \(.vault.name))"' | sort
```

### Find Items with Incomplete Tags
```bash
op item list --format json | jq -r '.[] | select(.tags and (.tags | length) == 1) | "\(.title) → \(.tags[0])"' | sort
```

### List Items by Path
```bash
for path in ui hp le ge; do
  echo "=== path:$path ==="
  op item list --tags "path:$path" --format json | jq -r '.[] | "\(.title) (Category: \(.tags[] | select(startswith("category:"))))"'
done
```

---

## Troubleshooting

### Issue: Tag not applied
**Solution:**
- Verify both tags were in the command: `op item edit [ID] --tags "path:ui,category:it"`
- Check if item is in a restricted vault (try web UI)

### Issue: Item not moving to vault
**Solution:**
- Use `op item move` instead of `op item edit --vault`
- If still fails, move via web UI

### Issue: Old vault showing items after move
**Solution:**
- The old vault takes time to sync - wait 30 seconds and refresh
- If persists, check vault cache settings in 1Password app

### Issue: Unrecognized tag category
**Solution:**
- Review category list above
- Create new category if necessary (lowercase, hyphenated)
- Apply to similar items for consistency

---

## Document Revision History

| Date | Version | Changes |
|------|---------|---------|
| 2026-01-13 | 1.0 | Initial completion - 94% tagged, all vaults reorganized |

---

## Contact & Support

For issues or questions about this structure:
- Check the scripts in `/tmp/` for automation examples
- Review recent commits in GitHub for migration details
- Refer to CLAUDE.md for 1Password CLI troubleshooting