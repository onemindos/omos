---
id: root-00-20260114140000
title: "EMAIL FINANCE FLOW EXAMPLE"
type: reference
status: active
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [type/reference, status/active, email, example, finance]
ai_summary: "Concrete examples showing how finance emails flow through the OneMind email system"
---
# Email Finance Flow Example

**Created:** 2026-01-14
**Purpose:** Demonstrate how finance emails route through pillar + category tagging system
**Scenario:** Understanding the difference between personal finance (#30-finance) and business finance (#80-treasury)

---

## The Confusion: What's the Difference?

You have TWO finance-related categories in OneMind Codex:

| Tag | Location | Purpose |
|-----|----------|---------|
| `#30-finance` | HP (25-49) - Holistic Performance | **Personal money** - Your salary, credit cards, investments, personal spending |
| `#80-treasury` | GE (75-99) - Generational Endeavor | **Business money** - Company revenue, business accounts, QuickBooks, business credit |

**The rule:** Same bank, different tag depending on WHERE the email is sent!

---

## Example 1: Chase Personal Credit Card

### Scenario
Chase sends you a statement for your **personal credit card** ending in 1234.

### Email Details
```
From: alerts@chase.com
To: hp@zeusdelacruz.me  ← YOU decide this by giving Chase this address
Subject: "Your Chase Sapphire credit card statement is ready"
```

### What Happens (Filter Execution Order)

**Step 1: Pillar Routing Filter (Priority 1)**
```
Filter name: Route HP
Criteria: To contains "hp@zeusdelacruz.me"
Action:
  - Move to folder: Holistic Performance
  - Apply tag: #hp
Result: ✅ Email lands in HP folder with #hp tag
```

**Step 2: Category Tagging Filter (Priority 10)**
```
Filter name: Auto-Tag Finance
Criteria: FROM contains "@chase.com"
Action:
  - Apply tag: #30-finance
  - Do NOT move (already in HP folder)
Result: ✅ Email gets second tag #30-finance
```

### Final Result
```
📧 Chase credit card statement
Location: Holistic Performance folder
Tags: #hp, #30-finance
Why: This is YOUR PERSONAL credit card
```

---

## Example 2: Chase Business Credit Card

### Scenario
Chase sends you a statement for your **business credit card** ending in 5678 (Legacy Restore LLC).

### Email Details
```
From: alerts@chase.com  ← SAME SENDER as Example 1
To: ge@zeusdelacruz.me  ← YOU decide this by giving Chase this address
Subject: "Your Chase Ink business card statement is ready"
```

### What Happens (Filter Execution Order)

**Step 1: Pillar Routing Filter (Priority 4)**
```
Filter name: Route GE
Criteria: To contains "ge@zeusdelacruz.me"
Action:
  - Move to folder: Generational Endeavor
  - Apply tag: #ge
Result: ✅ Email lands in GE folder with #ge tag
```

**Step 2: Category Tagging Filter (Priority 10)**
```
Filter name: Auto-Tag Finance
Criteria: FROM contains "@chase.com"
Action:
  - Apply tag: #30-finance
  - Do NOT move (already in GE folder)
Result: ✅ Email gets second tag #30-finance
```

**Wait... that's not right! We want #80-treasury, not #30-finance!**

### The Problem
**Both emails get #30-finance tag because they're from the same sender!**

This is where we need **subject line detection** to differentiate.

---

## Solution: Smart Finance Filter with Subject Detection

### Updated Filter #10: Auto-Tag Finance (Smart Version)

**Priority:** 10
**Name:** Auto-Tag Personal Finance
**Target tag:** `#30-finance`

```
IF (FROM contains ANY OF):
  - @chase.com
  - @bankofamerica.com
  - @paypal.com
  - @stripe.com
  - @schwab.com
  - @fidelity.com
  - @venmo.com

AND (SUBJECT does NOT contain):
  - "business"
  - "LLC"
  - "QuickBooks"
  - "business account"

THEN:
  - Apply tag: #30-finance
```

### New Filter #11: Auto-Tag Business Finance

**Priority:** 11
**Name:** Auto-Tag Business Finance
**Target tag:** `#80-treasury`

```
IF (FROM contains ANY OF):
  - @chase.com
  - @bankofamerica.com
  - @paypal.com
  - @stripe.com
  - @quickbooks.intuit.com

AND (SUBJECT contains ANY OF):
  - "business"
  - "LLC"
  - "QuickBooks"
  - "business account"
  - "business credit"

THEN:
  - Apply tag: #80-treasury
```

---

## Example 2 REVISED: Chase Business Credit Card (With Smart Filter)

### Email Details
```
From: alerts@chase.com
To: ge@zeusdelacruz.me
Subject: "Your Chase Ink business card statement is ready"
```

### What Happens (Fixed)

**Step 1: Pillar Routing**
```
Filter: Route GE
Result: Generational Endeavor folder + #ge tag
```

**Step 2: Finance Filter (Priority 10)**
```
Filter: Auto-Tag Personal Finance
Check: FROM contains "@chase.com" ✅
Check: SUBJECT contains "business" ❌ (EXCLUDED)
Result: SKIP - Not tagged #30-finance
```

**Step 3: Business Finance Filter (Priority 11)**
```
Filter: Auto-Tag Business Finance
Check: FROM contains "@chase.com" ✅
Check: SUBJECT contains "business" ✅ (MATCHED)
Result: ✅ Tagged #80-treasury
```

### Final Result (Correct!)
```
📧 Chase business credit card statement
Location: Generational Endeavor folder
Tags: #ge, #80-treasury
Why: Subject says "business" - this is BUSINESS money
```

---

## Example 3: Personal PayPal Payment

### Email Details
```
From: service@paypal.com
To: hp@zeusdelacruz.me
Subject: "You sent $50.00 to John Smith"
```

### What Happens

**Step 1: Pillar Routing**
```
Filter: Route HP
Result: Holistic Performance folder + #hp tag
```

**Step 2: Finance Filter**
```
Filter: Auto-Tag Personal Finance
Check: FROM contains "@paypal.com" ✅
Check: SUBJECT contains "business" ❌
Result: ✅ Tagged #30-finance
```

**Step 3: Business Finance Filter**
```
Filter: Auto-Tag Business Finance
Check: SUBJECT contains "business" ❌
Result: SKIP
```

### Final Result
```
📧 Personal PayPal payment
Location: Holistic Performance folder
Tags: #hp, #30-finance
Why: No "business" keyword - this is PERSONAL money
```

---

## Example 4: Business PayPal Payment

### Email Details
```
From: service@paypal.com
To: ge@zeusdelacruz.me
Subject: "Your business account received $500.00"
```

### What Happens

**Step 1: Pillar Routing**
```
Filter: Route GE
Result: Generational Endeavor folder + #ge tag
```

**Step 2: Finance Filter**
```
Filter: Auto-Tag Personal Finance
Check: FROM contains "@paypal.com" ✅
Check: SUBJECT contains "business" ❌ (EXCLUDED)
Result: SKIP
```

**Step 3: Business Finance Filter**
```
Filter: Auto-Tag Business Finance
Check: FROM contains "@paypal.com" ✅
Check: SUBJECT contains "business" ✅
Result: ✅ Tagged #80-treasury
```

### Final Result
```
📧 Business PayPal payment
Location: Generational Endeavor folder
Tags: #ge, #80-treasury
Why: Subject says "business" - this is BUSINESS money
```

---

## Example 5: QuickBooks Email (Always Business)

### Email Details
```
From: quickbooks@intuit.com
To: ge@zeusdelacruz.me
Subject: "Your monthly summary for Legacy Restore LLC"
```

### What Happens

**Step 1: Pillar Routing**
```
Filter: Route GE
Result: Generational Endeavor folder + #ge tag
```

**Step 2: Finance Filter**
```
Filter: Auto-Tag Personal Finance
Check: FROM contains "@quickbooks" ❌
Result: SKIP
```

**Step 3: Business Finance Filter**
```
Filter: Auto-Tag Business Finance
Check: FROM contains "@quickbooks.intuit.com" ✅
Check: SUBJECT contains "LLC" ✅
Result: ✅ Tagged #80-treasury
```

### Final Result
```
📧 QuickBooks monthly summary
Location: Generational Endeavor folder
Tags: #ge, #80-treasury
Why: QuickBooks is ALWAYS business
```

---

## Search Examples: The Power of Tags

### Search 1: All Finance Emails (Personal + Business)
```
Search: tag:#30-finance OR tag:#80-treasury
Result: Every financial email across ALL pillars and folders
Use case: Tax preparation, financial audit
```

### Search 2: All Personal Finance Only
```
Search: tag:#30-finance
Result: Personal credit cards, personal PayPal, personal investments
Use case: Personal budget review
```

### Search 3: All Business Finance Only
```
Search: tag:#80-treasury
Result: Business accounts, QuickBooks, business PayPal
Use case: Business accounting
```

### Search 4: All Chase Emails (Personal + Business)
```
Search: from:@chase.com
Result: Every Chase email regardless of tag
Use case: Looking for a specific Chase transaction
```

### Search 5: Personal Finance in HP Folder Only
```
Search: tag:#hp AND tag:#30-finance
Result: Only personal finance emails in HP folder
Use case: Review HP pillar financial health
```

### Search 6: Business Finance in GE Folder Only
```
Search: tag:#ge AND tag:#80-treasury
Result: Only business finance emails in GE folder
Use case: Review GE pillar financial health
```

---

## The Key Insight

**YOU control routing by which email address you give to services:**

| You Give Them | Email Routes To | Pillar Tag | Category Tag |
|---------------|----------------|------------|--------------|
| hp@zeusdelacruz.me | HP folder | #hp | #30-finance (personal) |
| ge@zeusdelacruz.me | GE folder | #ge | #80-treasury (business)* |

*Subject detection ensures correct business tagging

**Filter just adds the category tag based on sender + subject keywords!**

---

## Why This System is Powerful

### Before: One Big Mess
```
zeus@onemindcodex.com → ALL emails in one inbox
- Chase personal card statement
- Chase business card statement
- PayPal personal payment
- PayPal business payment
- QuickBooks summary
→ You manually sort everything 😵
```

### After: Automatic Organization
```
hp@zeusdelacruz.me → HP folder + #hp + #30-finance (personal money)
ge@zeusdelacruz.me → GE folder + #ge + #80-treasury (business money)
→ Zero manual sorting 🎉
```

---

## Summary: How Finance Tagging Works

1. **You choose the pillar** by which email address you give (hp@ or ge@)
2. **Filter applies category tag** based on sender domain + subject keywords
3. **Result:** Proper folder + pillar tag + category tag automatically

**Personal finance example:**
- Give Chase your `hp@zeusdelacruz.me` address
- Email arrives → HP folder + #hp + #30-finance
- Search `tag:#30-finance` to see ALL personal finance

**Business finance example:**
- Give QuickBooks your `ge@zeusdelacruz.me` address
- Email arrives → GE folder + #ge + #80-treasury
- Search `tag:#80-treasury` to see ALL business finance

**Mixed example (same bank, two accounts):**
- Personal Chase → `hp@zeusdelacruz.me` → #30-finance
- Business Chase → `ge@zeusdelacruz.me` → #80-treasury (subject detection)

---

## Does This Clear Up the Confusion?

The **critical understanding:**

1. **Email address (hp@ vs ge@)** = Which FOLDER it lands in
2. **Category tag (#30-finance vs #80-treasury)** = What TYPE of email it is
3. **Subject keywords** = Disambiguate when same sender has both personal and business

**Next step:** Apply this same logic to ALL categories (health, shopping, learning, etc.)

---

*Document created: 2026-01-14*
*Demonstrates OneMind email system using finance as concrete example*
*Ready to apply this pattern to all other categories*
