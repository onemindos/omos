---
id: root-00-20260114235900
title: "EMAIL PHASE 1 FILTER SPECIFICATIONS"
type: reference
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/reference, status/active, email, filters, implementation]
ai_summary: "Step-by-step specifications for creating 4 Phase 1 email category filters in Zoho Mail"
---
# Email Phase 1 Filter Specifications

**Created:** 2026-01-14
**Purpose:** Exact specifications to manually create 4 Phase 1 category filters in Zoho Mail
**Based on:** Analysis of 173 real emails from existing folders
**Priority:** Execute after domain routing filters (1-5) are complete

---

## How to Create Filters in Zoho Mail

1. **Log into Zoho Mail** → https://mail.zoho.com
2. **Go to Settings** (gear icon) → **Filters**
3. **Click "Create Filter"** for each filter below
4. **Set criteria exactly as specified**
5. **Set actions exactly as specified**
6. **Save and activate**

---

## Filter #6: Auto-Tag Mind (Learning & Education)

### Basic Info
- **Filter Name:** `Auto-Tag Mind`
- **Priority:** 6 (execute after filters 1-5)
- **Status:** Active
- **Apply to:** All incoming emails

### Criteria (Match ANY)
**FROM domain contains ANY of:**
- `substack.com`
- `skool.com`
- `deeplearning.ai`
- `joinsuperhuman.ai`
- `mail.beehiiv.com`
- `coursera.org`
- `udemy.com`
- `forwardfuture.com`
- `bytebytego.com`

**OR**

**SUBJECT contains ANY of:**
- `course`
- `lesson`
- `webinar`
- `newsletter`
- `learning path`

### Actions
- **Apply tag:** `#26-mind`
- **Do NOT move folder** (already routed by pillar filters)
- **Do NOT mark as read**
- **Do NOT forward**

### Zoho Mail Setup
```
Criteria Section:
├─ Match: Any condition
├─ Condition 1: From → contains → substack.com
├─ Condition 2: From → contains → skool.com
├─ Condition 3: From → contains → deeplearning.ai
├─ Condition 4: From → contains → joinsuperhuman.ai
├─ Condition 5: From → contains → mail.beehiiv.com
├─ Condition 6: From → contains → coursera.org
├─ Condition 7: From → contains → udemy.com
├─ Condition 8: From → contains → forwardfuture.com
├─ Condition 9: From → contains → bytebytego.com
├─ Condition 10: Subject → contains → course
├─ Condition 11: Subject → contains → lesson
├─ Condition 12: Subject → contains → webinar
├─ Condition 13: Subject → contains → newsletter
└─ Condition 14: Subject → contains → learning path

Actions Section:
└─ Apply tag → #26-mind
```

### Expected Performance
- **Accuracy:** 95%
- **Email volume:** ~50 existing emails match
- **Common matches:** Substack newsletters, Skool courses, DeepLearning.AI lessons

---

## Filter #7: Auto-Tag Body (Health & Wellness)

### Basic Info
- **Filter Name:** `Auto-Tag Body`
- **Priority:** 7
- **Status:** Active
- **Apply to:** All incoming emails

### Criteria (Match ANY)
**FROM domain contains ANY of:**
- `pausestudio.com`
- `riseabovefloatation.com`
- `email.alltrails.com`
- `mychart.ynhh.org`
- `equinox.com`
- `enagic.co.jp`
- `accesshealthct.com`
- `fitbit.com`
- `myfitnesspal.com`
- `peloton.com`

**OR**

**SUBJECT contains ANY of:**
- `appointment`
- `health record`
- `wellness`
- `workout`
- `membership`

### Actions
- **Apply tag:** `#27-body`
- **Do NOT move folder**
- **Do NOT mark as read**

### Zoho Mail Setup
```
Criteria Section:
├─ Match: Any condition
├─ Condition 1: From → contains → pausestudio.com
├─ Condition 2: From → contains → riseabovefloatation.com
├─ Condition 3: From → contains → email.alltrails.com
├─ Condition 4: From → contains → mychart.ynhh.org
├─ Condition 5: From → contains → equinox.com
├─ Condition 6: From → contains → enagic.co.jp
├─ Condition 7: From → contains → accesshealthct.com
├─ Condition 8: From → contains → fitbit.com
├─ Condition 9: From → contains → myfitnesspal.com
├─ Condition 10: From → contains → peloton.com
├─ Condition 11: Subject → contains → appointment
├─ Condition 12: Subject → contains → health record
├─ Condition 13: Subject → contains → wellness
├─ Condition 14: Subject → contains → workout
└─ Condition 15: Subject → contains → membership

Actions Section:
└─ Apply tag → #27-body
```

### Expected Performance
- **Accuracy:** 90%
- **Email volume:** ~33 existing emails match
- **Common matches:** Medical appointments, gym memberships, wellness services

---

## Filter #8: Auto-Tag Shopping (Personal Shopping)

### Basic Info
- **Filter Name:** `Auto-Tag Shopping`
- **Priority:** 8
- **Status:** Active
- **Apply to:** All incoming emails

### Criteria (Match ANY)
**FROM domain contains ANY of:**
- `uber.com`
- `lyftmail.com`
- `airsoftgi.com`
- `bemightykind.com`
- `gardencatering.com`
- `walmart.com`
- `joinsecret.com`
- `doordash.com`
- `grubhub.com`
- `instacart.com`

**OR**

**SUBJECT contains ANY of:**
- `order confirmation`
- `delivery scheduled`
- `shipment`
- `tracking number`

### Actions
- **Apply tag:** `#hp-shopping`
- **Do NOT move folder**
- **Do NOT mark as read**

### Zoho Mail Setup
```
Criteria Section:
├─ Match: Any condition
├─ Condition 1: From → contains → uber.com
├─ Condition 2: From → contains → lyftmail.com
├─ Condition 3: From → contains → airsoftgi.com
├─ Condition 4: From → contains → bemightykind.com
├─ Condition 5: From → contains → gardencatering.com
├─ Condition 6: From → contains → walmart.com
├─ Condition 7: From → contains → joinsecret.com
├─ Condition 8: From → contains → doordash.com
├─ Condition 9: From → contains → grubhub.com
├─ Condition 10: From → contains → instacart.com
├─ Condition 11: Subject → contains → order confirmation
├─ Condition 12: Subject → contains → delivery scheduled
├─ Condition 13: Subject → contains → shipment
└─ Condition 14: Subject → contains → tracking number

Actions Section:
└─ Apply tag → #hp-shopping
```

### Expected Performance
- **Accuracy:** 85%
- **Email volume:** ~40 existing emails match
- **Common matches:** Food delivery, online shopping, ride-sharing receipts

---

## Filter #9: Auto-Tag Infrastructure (Tech & Hosting)

### Basic Info
- **Filter Name:** `Auto-Tag Infrastructure`
- **Priority:** 9
- **Status:** Active
- **Apply to:** All incoming emails

### Criteria (Match ANY)
**FROM domain contains ANY of:**
- `amazon.com`
- `aws.amazon.com`
- `mercury.com`
- `github.com`
- `microsoft.com`
- `azure.com`
- `nvidia.com`
- `verizon.com`
- `cloudflare.com`
- `digitalocean.com`
- `heroku.com`
- `vercel.com`

**OR**

**SUBJECT contains ANY of:**
- `support case`
- `billing`
- `invoice`
- `service update`
- `security alert`

### Actions
- **Apply tag:** `#08-infrastructure`
- **Do NOT move folder**
- **Do NOT mark as read**

### Zoho Mail Setup
```
Criteria Section:
├─ Match: Any condition
├─ Condition 1: From → contains → amazon.com
├─ Condition 2: From → contains → aws.amazon.com
├─ Condition 3: From → contains → mercury.com
├─ Condition 4: From → contains → github.com
├─ Condition 5: From → contains → microsoft.com
├─ Condition 6: From → contains → azure.com
├─ Condition 7: From → contains → nvidia.com
├─ Condition 8: From → contains → verizon.com
├─ Condition 9: From → contains → cloudflare.com
├─ Condition 10: From → contains → digitalocean.com
├─ Condition 11: From → contains → heroku.com
├─ Condition 12: From → contains → vercel.com
├─ Condition 13: Subject → contains → support case
├─ Condition 14: Subject → contains → billing
├─ Condition 15: Subject → contains → invoice
├─ Condition 16: Subject → contains → service update
└─ Condition 17: Subject → contains → security alert

Actions Section:
└─ Apply tag → #08-infrastructure
```

### Expected Performance
- **Accuracy:** 80%
- **Email volume:** ~50 existing emails match
- **Common matches:** AWS support, GitHub notifications, hosting bills

---

## Post-Creation Checklist

### After Creating All 4 Filters

**1. Verify Filter Order**
```
Priority 1: Route HP (domain routing)
Priority 2: Route LE (domain routing)
Priority 3: Route UI (domain routing)
Priority 4: Route GE (domain routing)
Priority 5: Catch-All zeusdelacruz.me (domain routing)
Priority 6: Auto-Tag Mind ← NEW
Priority 7: Auto-Tag Body ← NEW
Priority 8: Auto-Tag Shopping ← NEW
Priority 9: Auto-Tag Infrastructure ← NEW
```

**2. Test Each Filter**
Send test emails to verify tagging:
- Test email to hp@ with "course" in subject → Should get #hp + #26-mind
- Test email to hp@ with "appointment" in subject → Should get #hp + #27-body
- Test email to hp@ with "order confirmation" → Should get #hp + #hp-shopping
- Test email to ui@ from github.com → Should get #ui + #08-infrastructure

**3. Check Existing Emails**
Search for emails that should have been tagged:
- Search: `from:substack.com` → Should all have #26-mind tag
- Search: `from:mychart.ynhh.org` → Should all have #27-body tag
- Search: `from:doordash.com` → Should all have #hp-shopping tag
- Search: `from:github.com` → Should all have #08-infrastructure tag

**4. Retroactive Tagging (If Needed)**
If existing emails don't get tagged automatically:
- Manually select all matching emails
- Apply the appropriate tag
- Or wait for Zoho to process (may take 24 hours)

---

## Success Criteria

**You'll know it's working when:**
- ✅ New emails automatically get both pillar tag (#hp, #le, #ui, #ge) AND category tag (#26-mind, etc.)
- ✅ Search `tag:#26-mind` returns all learning emails across ALL folders
- ✅ Search `tag:#27-body` returns all health emails across ALL folders
- ✅ Search `tag:#hp-shopping` returns all shopping emails (mostly in HP)
- ✅ Search `tag:#08-infrastructure` returns all tech emails (mostly in UI)

---

## Troubleshooting

### Filter Not Working?
1. **Check filter priority** - Must be 6-9 (after domain routing)
2. **Check "Match ANY"** - Not "Match ALL"
3. **Check filter is active** - Toggle on if needed
4. **Test with new email** - Old emails may not retroactively apply

### Wrong Tag Applied?
1. **Check for typos** - Tag must be exact: `#26-mind` not `#26-Mind`
2. **Check duplicate filters** - Remove any conflicting filters
3. **Check criteria** - Verify sender domains are spelled correctly

### No Tags at All?
1. **Verify domain routing works first** - Email should land in HP/LE/SO/GE folder
2. **Check if tags exist** - Create tags manually if needed
3. **Wait 5 minutes** - Filters may take time to process

---

## Next Steps After Creation

1. ✅ Create all 4 filters (you're doing this now)
2. ⏳ Test with new emails (send test messages)
3. ⏳ Verify existing emails get tagged
4. ⏳ Monitor accuracy for 1 week
5. ⏳ Adjust keywords if needed
6. ⏳ Prepare Phase 2 filters (#30-finance, #78-legacy-restore)

---

## Phase 2 Preview (Coming Next)

After Phase 1 is stable, we'll add:
- **Filter #10:** Auto-Tag Finance (`#30-finance`)
- **Filter #11:** Auto-Tag Legacy Restore (`#78-legacy-restore`)

These will activate when you start account migration (Week 2).

---

*Document created: 2026-01-14*
*Ready for manual filter creation in Zoho Mail*
*Estimated creation time: 10-15 minutes total*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
