---
id: ui-08-filter1-mind-update
title: "Filter #1 Mind - Missing Domains Update"
type: guide
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/guide, status/active, platform/zoho, system/email]
ai_summary: "Update specifications for Filter #1 (Mind) to include missing domains discovered in Education folder analysis"
---

# Filter #1: Mind & Education - Update Guide

> **Issue:** 50 emails in "Zeus Delacruz/Education" folder not tagged with #26-mind
> **Cause:** Missing domains and Substack senders from Filter #1
> **Solution:** Add 15+ missing conditions to existing filter

---

## Analysis Summary

**Untagged Emails Found:** 50 in Education subfolder (folderId: 8698543000000494019)
**Current Filter #1:** 39 conditions (16 domains + 15 subject + 8 content)
**Missing Conditions:** 15+ domains and email senders

---

## Missing Domains to Add

### High Priority Domains (6)

**CRITICAL:** Add these as `FROM` → `contains` (NOT "domain contains")

```
1. forwardfuture.ai        - Forward Future (Matthew Berman) - AI news/education
2. liamottley.com          - Liam Ottley - AI automation education
3. gethalfbaked.com        - Half Baked - Startup ideas newsletter
4. simplilearnmailer.com   - Purdue/Simplilearn - Professional courses
5. fundrapp.com            - Business funding/credit education
6. fundingvault.io         - Money Broker Society - Business funding
```

**In Zoho Mail UI:**
- Field: `FROM`
- Operator: `contains`
- Value: `forwardfuture.ai` (domain) or `astro@forwardfuture.ai` (email)

### Substack Senders (8+)

**Problem:** Many newsletters route through substack.com (generic domain)
**Solution:** Use `FROM contains` with specific email addresses

Add these as `FROM contains`:

```
7.  astro@forwardfuture.ai           - Forward Future (AI news)
8.  thecapitalistmag@substack.com    - The Capitalist (Finance/economics)
9.  ainutpam@substack.com            - AI Monk (AI tools education)
10. linas@substack.com               - Linas Beliūnas (FinTech/AI education)
11. ideas@mail.gethalfbaked.com      - Half Baked (Startup ideas)
12. willywoo@substack.com            - Bitcoin Vector Lite (Crypto education)
13. natehagens@substack.com          - The Great Simplification (Systems thinking)
14. yungpueblo@substack.com          - Personal development/mindfulness
15. superhuman@mail.joinsuperhuman.ai - AI news digest
```

### Medium Priority (Optional - Add if seeing untagged emails)

```
16. news.militaryaerospace.com       - Military aerospace education
17. verkada.com                      - Security/camera tech education
18. pagerduty.com                    - Incident management education
19. billing.venturebeat.com          - Tech industry education (VentureBeat)
```

---

## Step-by-Step Update Instructions

### Option A: Edit Existing Filter (Recommended)

1. **Go to:** https://mail.zoho.com → Settings (⚙️) → Filters
2. **Find:** Filter "6 - Auto-Tag Mind"
3. **Click:** Edit (pencil icon)
4. **Scroll to:** Conditions section
5. **Add new conditions** (15 total):
   - Add 6 domain conditions (`FROM domain contains`)
   - Add 9 email sender conditions (`FROM contains`)
6. **Click:** Save
7. **Test:** Send email from one of the new domains

### Option B: Create Second Mind Filter (If condition limit hit)

If Zoho has a condition limit per filter:

**Filter 6a: Auto-Tag Mind (Domains)**
- Keep existing 16 domains + add 6 new = 22 domains
- Keep existing subject keywords
- Keep existing content keywords

**Filter 6b: Auto-Tag Mind (Substack)**
- Add 9 Substack sender conditions
- Same action: Apply label #26-mind

---

## Complete Updated Filter Specification

### Filter #1 Updated: Mind & Education (#26-mind)

**Total Conditions:** 54 (was 39)
- 22 domain conditions (was 16)
- 9 sender email conditions (NEW)
- 15 subject keywords (unchanged)
- 8 content keywords (unchanged)

#### Full Domain List (22):

**IMPORTANT:** Use `FROM` → `contains` for all conditions (NOT "domain contains")

```
FROM contains:
1.  deeplearning.ai
2.  coursera.org
3.  udemy.com
4.  linkedin.com
5.  skillshare.com
6.  masterclass.com
7.  pluralsight.com
8.  datacamp.com
9.  codecademy.com
10. khanacademy.org
11. edx.org
12. udacity.com
13. workweek.com
14. venturebeat.com
15. infinitusmail.com
16. relevanceai.com
17. forwardfuture.ai         ← NEW
18. liamottley.com           ← NEW
19. gethalfbaked.com         ← NEW
20. simplilearnmailer.com    ← NEW
21. fundrapp.com             ← NEW
22. fundingvault.io          ← NEW
```

**Why `FROM contains` not `FROM domain contains`:**
- `FROM contains: coursera.org` → Matches emails from @coursera.org (specific)
- `FROM domain contains: substack` → Matches ALL @substack.com (1000+ emails - TOO BROAD!)
- Always use `FROM contains` with the domain or full email address

#### Substack Senders (9 NEW):

```
FROM contains:
1. astro@forwardfuture.ai
2. thecapitalistmag@substack.com
3. ainutpam@substack.com
4. linas@substack.com
5. ideas@mail.gethalfbaked.com
6. willywoo@substack.com
7. natehagens@substack.com
8. yungpueblo@substack.com
9. superhuman@mail.joinsuperhuman.ai
```

#### Subject Keywords (15 - Unchanged):

```
SUBJECT contains:
course, learn, lesson
education, training
certificate, certification
tutorial, workshop, webinar
study, class, lecture
enrollment, curriculum
```

#### Content Keywords (8 - Unchanged):

```
MESSAGE contains:
complete the course
new lesson
learning path
skill development
professional development
continue your learning
course enrollment
certificate of completion
```

#### Actions (Unchanged):
- Apply label: `#26-mind`
- Keep in current folder

---

## Testing After Update

### Test 1: Verify New Domains Work

Send test emails:

```bash
# Test 1: Forward Future domain
Subject: TEST: AI news from Forward Future
From: astro@forwardfuture.ai
To: hp@zeusdelacruz.me

# Test 2: Half Baked domain
Subject: TEST: Startup idea newsletter
From: ideas@mail.gethalfbaked.com
To: hp@zeusdelacruz.me

# Test 3: Liam Ottley domain
Subject: TEST: AI automation course
From: admin@liamottley.com
To: hp@zeusdelacruz.me
```

**Expected Result:** All 3 emails should have dual tags (#hp + #26-mind)

### Test 2: Apply to Old Education Folder

**Option A: Automatic (Zoho UI)**
1. Go to: Zeus Delacruz/Education folder
2. Select all 50 emails (Ctrl+A or Cmd+A)
3. Click: **More** → **Apply Filters**
4. Zoho will retroactively apply updated Filter #1

**Option B: Manual Move + Auto-Tag**
1. Select all 50 emails in Education folder
2. Move to: HP folder
3. Updated Filter #1 will automatically tag them with #26-mind

### Test 3: Search Verification

After update, search for:
```
tag:26-mind
```

**Expected Result:** Should show ~90-100 tagged emails (was ~40 before update)

---

## Expected Impact

### Before Update:
- 39 conditions
- ~40 emails tagged in HP folder
- 50 untagged emails in Education folder
- Missing newsletters from Substack

### After Update:
- 54 conditions (+15)
- ~90-100 emails tagged across folders
- Education folder emails tagged
- Substack newsletters caught

### Coverage Improvement:
- **Before:** ~45% of education emails tagged
- **After:** ~95% of education emails tagged
- **Gain:** +50% coverage

---

## Maintenance Notes

### Weekly Review (First Month)

Check for untagged education emails:
1. Search: `folder:HP NOT tag:26-mind`
2. Review senders of untagged emails
3. Add new domains/senders to filter

### Common Untagged Senders to Watch For:

- New Substack newsletters you subscribe to
- Podcast newsletters (Spotify, Apple Podcasts)
- YouTube education channels
- Medium writers
- Academic institutions (.edu domains)

### Adding New Domains (Quick Process)

1. Find untagged education email
2. Note sender domain or email
3. Edit Filter #1
4. Add new condition
5. Save
6. **Takes 30 seconds per domain**

---

## Troubleshooting

### Issue: Filter not catching Substack emails

**Problem:** Using `FROM domain contains: substack.com` is too broad
**Solution:** Use specific email addresses: `FROM contains: sender@substack.com`

### Issue: Still seeing untagged education emails

**Checklist:**
- [ ] Is filter enabled?
- [ ] Are new conditions saved?
- [ ] Is email in a folder that filters process? (Not in old Zeus Delacruz subfolders?)
- [ ] Does sender domain/email match any condition?

**Solution:** Move emails to HP folder, filters will auto-apply

### Issue: Filter catching non-education emails

**Rare with domain-based filtering**
- Review false positive
- If consistent, add exclusion: `SUBJECT does not contain: [keyword]`

---

## Alternative: Import via CSV (Advanced)

If manually adding 15 conditions is tedious, export/import:

### Export Current Filter
1. Note all current conditions
2. Create CSV with new conditions

### Import Format (Zoho Filter CSV)
```csv
Condition Type,Field,Operator,Value
FROM,domain,contains,forwardfuture.ai
FROM,domain,contains,liamottley.com
...
```

**Note:** Zoho may not support filter CSV import - manual entry recommended

---

## Next Steps After Update

1. ✅ Update Filter #1 with 15 new conditions
2. ✅ Test with 3 sample emails
3. ✅ Apply to Education folder (50 emails)
4. ✅ Verify ~90-100 total tagged emails
5. ⏭️ Monitor for 2-3 days
6. ⏭️ Create Filter #2 (Body & Fitness)

---

## Quick Copy-Paste List

**Domains (Add to filter):**
```
forwardfuture.ai
liamottley.com
gethalfbaked.com
simplilearnmailer.com
fundrapp.com
fundingvault.io
```

**Senders (Add to filter):**
```
astro@forwardfuture.ai
thecapitalistmag@substack.com
ainutpam@substack.com
linas@substack.com
ideas@mail.gethalfbaked.com
willywoo@substack.com
natehagens@substack.com
yungpueblo@substack.com
superhuman@mail.joinsuperhuman.ai
```

---

*Created: 2026-01-14*
*Status: Ready for implementation*
*Estimated Time: 5-7 minutes to add 15 conditions*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
