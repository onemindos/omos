---
id: ui-08-filter1-mind-guide
title: "Filter #1: Mind & Education - Step-by-Step Creation Guide"
type: guide
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/guide, status/active, platform/zoho, system/email]
ai_summary: "Complete step-by-step instructions for creating the Mind/Education smart filter in Zoho Mail UI"
---

# Filter #1: Mind & Education - Creation Guide

> **Goal:** Apply `#26-mind` tag to all education and learning-related emails in HP folder

---

## Quick Reference

**Filter Name:** `6 - Auto-Tag Mind`
**Category:** Holistic Performance → Mind (Education & Learning)
**Tag:** `#26-mind`
**Expected Matches:** 50-100 emails
**Estimated Setup Time:** 8-10 minutes

---

## Step-by-Step Instructions

### Step 1: Access Filter Settings
1. Open browser and go to: https://mail.zoho.com
2. Click **Settings** (⚙️ icon) in top-right corner
3. Click **Filters** in left sidebar
4. Click **Create Filter** button

---

### Step 2: Basic Filter Settings

**Filter Name:** `6 - Auto-Tag Mind`
**Description (optional):** `Apply #26-mind tag to education and learning emails`

---

### Step 3: Add Conditions - TIER 1 (Domain List)

**Operator:** Select **"Match ANY of the following"** (OR logic)

Add these conditions one by one:

#### Domain Conditions (16 total):

**CRITICAL:** Use `FROM` → `contains` (NOT "domain contains")

1. **From** → **contains** → `deeplearning.ai`
2. **From** → **contains** → `coursera.org`
3. **From** → **contains** → `udemy.com`
4. **From** → **contains** → `linkedin.com`
5. **From** → **contains** → `skillshare.com`
6. **From** → **contains** → `masterclass.com`
7. **From** → **contains** → `pluralsight.com`
8. **From** → **contains** → `datacamp.com`
9. **From** → **contains** → `codecademy.com`
10. **From** → **contains** → `khanacademy.org`
11. **From** → **contains** → `edx.org`
12. **From** → **contains** → `udacity.com`
13. **From** → **contains** → `workweek.com`
14. **From** → **contains** → `venturebeat.com`
15. **From** → **contains** → `infinitusmail.com`
16. **From** → **contains** → `relevanceai.com`

**Why NOT "domain contains":**
- `FROM contains: coursera.org` matches @coursera.org emails ✅
- `FROM domain contains: substack` matches ALL @substack.com (1000+) ❌

**Pro Tip:** Click "+ Add Condition" after each one. This section should take ~3 minutes.

---

### Step 4: Add Conditions - TIER 2 (Subject Keywords)

Continue adding conditions to the SAME filter:

#### Subject Conditions (15 total):
17. **Subject** → **contains** → `course`
18. **Subject** → **contains** → `learn`
19. **Subject** → **contains** → `lesson`
20. **Subject** → **contains** → `education`
21. **Subject** → **contains** → `training`
22. **Subject** → **contains** → `certificate`
23. **Subject** → **contains** → `certification`
24. **Subject** → **contains** → `tutorial`
25. **Subject** → **contains** → `workshop`
26. **Subject** → **contains** → `webinar`
27. **Subject** → **contains** → `study`
28. **Subject** → **contains** → `class`
29. **Subject** → **contains** → `lecture`
30. **Subject** → **contains** → `enrollment`
31. **Subject** → **contains** → `curriculum`

**Pro Tip:** These are case-insensitive by default. This section should take ~2-3 minutes.

---

### Step 5: Add Conditions - TIER 3 (Content/Body Keywords)

Continue adding conditions:

#### Content Conditions (8 total):
32. **Content** → **contains** → `complete the course`
33. **Content** → **contains** → `new lesson`
34. **Content** → **contains** → `learning path`
35. **Content** → **contains** → `skill development`
36. **Content** → **contains** → `professional development`
37. **Content** → **contains** → `continue your learning`
38. **Content** → **contains** → `course enrollment`
39. **Content** → **contains** → `certificate of completion`

**Pro Tip:** Content searching is slower but catches emails with generic subjects. This section should take ~1-2 minutes.

---

### Step 6: Set Actions

**Important:** We're ONLY applying a tag, NOT moving the email (already in HP folder)

1. Click **"Then perform these actions"** section
2. Add Action: **Apply Label**
3. Select label: `#26-mind`
   - If label doesn't exist, click "Create New Label" and name it `26-mind`
4. **DO NOT add "Move to folder" action** - emails are already in correct pillar folder

---

### Step 7: Save and Test

1. Click **Save** button
2. **Send Test Email:**
   - From: zeus@zeusdelacruz.com
   - To: hp@zeusdelacruz.me
   - Subject: `TEST: Coursera new course available`
   - Body: `Complete the course on AI development. New lesson available.`

3. **Wait 10 seconds** for email to arrive
4. **Check HP folder** for the test email
5. **Verify DUAL TAGS applied:**
   - ✅ `#hp` (from pillar filter)
   - ✅ `#26-mind` (from this new filter)

---

## Verification Checklist

After creating the filter:

- [ ] Filter name: `6 - Auto-Tag Mind`
- [ ] Operator: "Match ANY of the following" (OR logic)
- [ ] 39 total conditions added (16 domains + 15 subject + 8 content)
- [ ] Action: Apply label `#26-mind`
- [ ] NO folder movement action
- [ ] Test email sent and received
- [ ] Dual tags visible on test email
- [ ] Filter is active/enabled

---

## Testing Existing Emails

Once filter is created and tested, manually apply it to existing emails:

### Method 1: Apply to Folder (Zoho Web UI)
1. Go to HP folder
2. Select multiple emails (Shift+Click)
3. Click **More** → **Apply Filters**
4. Zoho will retroactively apply the filter

### Method 2: Search and Tag (Manual)
1. Search for: `from:deeplearning.ai`
2. Select all results
3. Right-click → **Apply Label** → `#26-mind`

---

## Troubleshooting

### Issue: Test email doesn't have #26-mind tag
**Solution:**
- Check filter is enabled (not paused)
- Verify action is "Apply Label: 26-mind"
- Check filter priority (should be 6 or higher, after pillar filters)
- Wait 30 seconds and refresh mailbox

### Issue: Can't find #26-mind label
**Solution:**
- Create label manually: Settings → Labels → Create → Name: `26-mind`
- Return to filter and select the newly created label

### Issue: Too many conditions error
**Solution:**
- Zoho has condition limits per filter
- If hit limit, split into two filters:
  - `6a - Auto-Tag Mind (Domains)`
  - `6b - Auto-Tag Mind (Keywords)`

### Issue: Filter catching wrong emails
**Solution:**
- Review false positives
- Remove or adjust overly broad keywords
- Add exclusion conditions if needed

---

## Performance Monitoring

### Week 1: Daily Checks
- Count tagged emails: Search `tag:#26-mind`
- Review for false positives (wrongly tagged)
- Review for false negatives (missed emails)
- Add new domains discovered

### Week 2-4: Weekly Checks
- Assess filter performance
- Add 2-3 new domains as you discover them
- Adjust keywords if needed

### Monthly: Optimization
- Review most common senders
- Consolidate or split filters if needed
- Update domain list

---

## Quick Copy-Paste Reference

### Domains (Tier 1):
```
deeplearning.ai
coursera.org
udemy.com
linkedin.com
skillshare.com
masterclass.com
pluralsight.com
datacamp.com
codecademy.com
khanacademy.org
edx.org
udacity.com
workweek.com
venturebeat.com
infinitusmail.com
relevanceai.com
```

### Subject Keywords (Tier 2):
```
course
learn
lesson
education
training
certificate
certification
tutorial
workshop
webinar
study
class
lecture
enrollment
curriculum
```

### Content Keywords (Tier 3):
```
complete the course
new lesson
learning path
skill development
professional development
continue your learning
course enrollment
certificate of completion
```

---

## Success Metrics

**After 1 week, you should see:**
- ✅ 50+ emails tagged with `#26-mind`
- ✅ Zero false positives (non-education emails tagged)
- ✅ <5% false negatives (education emails not tagged)
- ✅ Clear category separation in HP folder

**Filter is working if:**
- Test email gets dual tags (#hp + #26-mind)
- Existing Coursera/Udemy emails now show #26-mind
- You can search `tag:#26-mind` and see only education emails

---

## Next Steps

Once Filter #1 is working perfectly:
1. ✅ Monitor for 2-3 days
2. ⏭️ Create Filter #2: Body & Fitness (#27-body)
3. ⏭️ Create Filter #3: Shopping (#hp-shopping)
4. ⏭️ Create Filter #4: Infrastructure (#08-infrastructure)

---

*Estimated Total Time: 8-10 minutes*
*Created: 2026-01-14*
*Status: Ready for implementation*
