# Email Structure - Addendum

**Additions to main plan (EMAIL-STRUCTURE-IMPLEMENTATION-PLAN.md)**

---

## 1. Existing Folders to Archive

**Current "Zeus Delacruz" folder contains 14 subfolders:**
```
Zeus Delacruz/
├── Education
├── Family
├── Finance
├── Legacy
├── Legal
├── Performance
├── Personal
├── Productivity
├── Recreation
├── Save for Later
├── Shopping
├── Social
├── Training
├── Travel
└── Wealth
```

**Also exists:**
```
Legacyhq - Old
legacyhq - New
```

**Action:** Rename "Zeus Delacruz" → "Archive - Pre-2026" (keep emails, don't delete)

---

## 2. Why NO Subfolders

Tags are better than subfolders because:
- One email can have multiple tags (can't be in multiple folders)
- Click tag = see all emails across pillars
- Less navigation overhead
- Easier to search

**Example:** An email about "family travel" can be tagged #family AND #travel. With folders, you'd have to pick one.

---

## 3. Why NO Archive Folder

Zoho Mail has built-in archive:
- Press "E" or click Archive button
- Goes to system Archive folder
- No need to create custom archive folder

---

## 4. Zoho Desk Auto-Tag for Legacy Restore

**Setup in Zoho Desk:**
```
1. Go to: Zoho Desk → Setup → Automation → Workflow Rules
2. Create Rule:
   - Name: Tag Legacy Restore Tickets
   - Trigger: When ticket is created
   - Condition: Email received at support@legacyrestore.co
   - Action: Add tag #LegacyRestore
3. Save
```

**Result:** All Legacy Restore support emails auto-tagged in Zoho Desk queue.

---

## 5. Main Inbox Emails (No Filter Needed)

**These emails stay in Inbox - NO filter required:**
```
zeus@onemindcodex.com    → Stays in Inbox (main business)
zeus@zeusdelacruz.com    → Stays in Inbox (personal brand)
```

**Why:** These are your primary contact emails. You triage them manually. Only pillar emails (hp@, le@, ui@, ge@) get auto-routed to folders.

---

**End of Addendum**


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
