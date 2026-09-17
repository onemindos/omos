---
id: ui-02-20260111053925
title: "Tool Examples — Real World Flows"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/02]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Tool Examples — Real World Flows"
---

# Tool Examples — Real World Flows

> **How Legacy-AI uses custom tools to automate your life**

---

## Overview

This document shows **real examples** of how each tool category works in practice — what you say, what happens behind the scenes, and what you receive.

All tools are Python functions that Letta agents call autonomously.

---

## Architecture Recap

```
YOU → Legacy-AI → Tools → External Services → Response
         ↓
    Worker Agents (HP, LE, GE, IT)
         ↓
    Shared Memory (instant sync)
```

---

# 1. FINANCIAL TOOLS

## `get_bank_balance`

**You say:**
> "@Legacy what's my cash position right now?"

**Legacy-AI processes:**
1. Recognizes: financial query → GE domain
2. Selects tool: `get_bank_balance`
3. Connects to: Plaid API (aggregates all bank accounts)
4. Retrieves: Checking, Savings, Business accounts

**You receive:**
```
💰 Cash Position (as of 2:34 PM)

Personal
├── Checking (Chase): $4,287.32
├── Savings (Marcus): $12,450.00
└── Emergency Fund: $8,000.00

Business (Legacy AI LLC)
├── Operating (Mercury): $23,891.45
├── Tax Reserve: $6,200.00
└── Runway: 4.2 months

Total Liquid: $54,828.77
```

**Background action:**
- Logs balance check to `30-39 GE/36 Treasury/balance-log.md`
- Updates shared memory block `ge_metrics` with current runway

---

## `categorize_expense`

**You say:**
> "I just spent $847 at Home Depot"

**Legacy-AI processes:**
1. Recognizes: expense entry
2. Analyzes: "Home Depot" → home improvement
3. Determines: Could be LE (home) or GE (business expense)
4. Asks for clarification OR auto-categorizes based on patterns

**You receive:**
```
🏠 Expense Logged: $847.00 at Home Depot

Categorized as: LE → Sanctuary → Home Improvement
(Based on: Previous Home Depot purchases were personal)

Is this correct?
├── ✅ Yes, personal home expense
├── 🏢 No, business expense (deductible)
└── 📝 Split between both
```

**If you confirm personal:**
- Logs to `20-29 LE/21 Sanctuary/expenses/2025-12.md`
- Updates monthly home spending tracker

**If you say business:**
- Logs to `30-39 GE/36 Treasury/expenses/2025-12.md`
- Tags as tax-deductible
- Adds to quarterly expense report

---

## `revenue_report`

**You say:**
> "How did we do this month revenue-wise?"

**Legacy-AI processes:**
1. Delegates to GE Worker (specialist in business data)
2. GE Worker calls: Stripe API, PayPal API, invoice tracker
3. Aggregates all revenue streams
4. Compares to previous month and targets

**You receive:**
```
📊 December 2025 Revenue Report

Revenue Streams
├── Course Sales: $3,240 (12 sales)
├── Consulting: $5,500 (2 clients)
├── Affiliate: $892
└── Sponsorship: $1,500
────────────────────────
Total: $11,132

vs Last Month: +18% ($9,432)
vs Target: 89% of $12,500 goal

🎯 Gap to Target: $1,368

Top Performers:
1. "One Mind Academy" course — $2,400
2. Consulting: Zeus Corp — $3,500

⚠️ Note: 2 invoices outstanding ($2,100)
```

**Background action:**
- Creates monthly report in `30-39 GE/35 Analytics/monthly/2025-12.md`
- Updates `ge_metrics` shared memory
- If below target, suggests actions

---

# 2. HOME AUTOMATION TOOLS

## `set_thermostat`

**You say:**
> "It's freezing in here, warm it up"

**Legacy-AI processes:**
1. Recognizes: home comfort request → LE domain
2. Checks current temp via Home Assistant API
3. Determines appropriate adjustment
4. Executes change

**You receive:**
```
🌡️ Thermostat Adjusted

Current: 66°F (too cold!)
Setting to: 72°F

Estimated time to reach: 18 minutes

💡 Tip: Your schedule shows you usually prefer 70°F at this time.
Want me to update your evening schedule to 72°F?
```

**Background action:**
- Logs temperature change to home automation history
- Learns your preference pattern
- Updates LE Worker's memory about comfort preferences

---

## `check_security`

**You say:**
> "Are all the doors locked? We're leaving for vacation tomorrow"

**Legacy-AI processes:**
1. Recognizes: security check → LE + IT domain
2. Queries Home Assistant for all lock/sensor states
3. Checks cameras
4. Provides comprehensive status

**You receive:**
```
🔒 Security Status Check

Doors
├── Front Door: 🔒 Locked
├── Back Door: 🔒 Locked
├── Garage Entry: 🔒 Locked
└── Garage Door: ⬇️ Closed

Windows
├── Living Room: ✅ Closed
├── Kitchen: ⚠️ OPEN (left window)
├── Bedroom: ✅ Closed
└── Office: ✅ Closed

Cameras
├── Front Porch: 🟢 Online, recording
├── Backyard: 🟢 Online, recording
└── Garage: 🟢 Online, recording

Motion Sensors: All armed

⚠️ ACTION NEEDED: Kitchen window is open!

Want me to:
├── 🔔 Send reminder to close it
├── 📋 Create vacation security checklist
└── 🏠 Activate "Away Mode" (lights on timer, extra alerts)
```

**Background action:**
- Creates vacation prep checklist in `20-29 LE/21 Sanctuary/checklists/`
- Logs security check timestamp
- If you activate Away Mode, schedules automated light patterns

---

## `water_plants`

**You say:**
> "Water the garden, it's been dry"

**Smart irrigation scenario:**
```
🌱 Irrigation Triggered

Zone 1 (Vegetable Garden): Running 15 min
Zone 2 (Herb Bed): Queued, starts in 15 min
Zone 3 (Fruit Trees): Queued, starts in 25 min

Total water usage estimate: 45 gallons

☀️ Weather Note: No rain expected for 5 days.
Recommend watering again Thursday.

Want me to set up automatic schedule based on weather?
```

**Manual scenario (no smart irrigation):**
```
🌱 Garden Watering Reminder Created

I don't see smart irrigation connected, so I've created a task:

📋 Task: Water the garden
📅 Due: Today, before 6 PM (avoid midday heat)
📍 Project: LE - Homestead

Checklist added:
- [ ] Vegetable beds
- [ ] Herb containers
- [ ] Fruit trees (deep water)
- [ ] Check soil moisture after

Want me to connect to a smart irrigation system?
```

---

# 3. CONTENT TOOLS

## `draft_blog_post`

**You say:**
> "Draft a blog post about why I built One Mind instead of using PARA"

**Legacy-AI processes:**
1. Recognizes: content creation → GE domain
2. Delegates to GE Worker (content specialist)
3. GE Worker:
   - Searches vault for existing notes on PARA, One Mind philosophy
   - Pulls from ONEMIND-CODEX.md
   - Generates structured draft
4. Returns draft for review

**You receive:**
```
📝 Blog Post Draft Ready

Title: "Why I Abandoned PARA and Built One Mind"
Word Count: 1,247
Reading Time: 5 min

---

OUTLINE:
1. The Problem with "Second Brains" (hook)
2. What PARA Gets Right
3. Where PARA Falls Short for Founders
4. The One Mind Difference: AI-Native Design
5. The 4 Pillars vs 4 Folders
6. Call to Action

---

DRAFT PREVIEW:

"Every productivity guru tells you to build a 'second brain.'
I tried. I really tried. PARA, Zettelkasten, Johnny Decimal —
I've used them all. But as a founder building a family legacy,
I kept hitting the same wall..."

[Full draft: 1,247 words]

---

📍 Saved to: 30-39 GE/34 Creation/blog/drafts/para-vs-onemind.md

Actions:
├── ✏️ Edit in Obsidian
├── 🔍 Review full draft
├── 📅 Schedule for publishing
└── 🗑️ Discard and start over
```

---

## `schedule_social`

**You say:**
> "Schedule that blog post to go out across all platforms Tuesday at 9am"

**You receive:**
```
📱 Social Media Scheduled

Blog Post: "Why I Abandoned PARA and Built One Mind"
Date: Tuesday, Dec 31 at 9:00 AM EST

Platform Versions:

🐦 Twitter/X (Thread - 7 tweets)
├── Tweet 1: Hook + link
├── Tweet 2-6: Key points
└── Tweet 7: CTA
Status: ✅ Scheduled

💼 LinkedIn (Long-form)
├── Full article adaptation
├── 3 custom images attached
Status: ✅ Scheduled

📸 Instagram (Carousel)
├── 8 slides created
├── Caption with hashtags
Status: ⚠️ Needs approval (external post)

🎬 YouTube Community
├── Teaser post with link
Status: ✅ Scheduled

---

⚠️ APPROVAL NEEDED for Instagram
(External platform, requires confirmation)

[Approve Instagram Post] [Edit First] [Skip Instagram]
```

---

## `check_analytics`

**You say:**
> "How's my content performing this week?"

**You receive:**
```
📈 Content Performance (Dec 19-26)

YouTube
├── Views: 12,847 (+34% vs last week)
├── Watch Time: 892 hours
├── New Subs: 127
├── Top Video: "Building JARVIS for Real" — 4.2K views
└── 💡 Insight: Shorts driving 60% of new subs

Twitter/X
├── Impressions: 89,432
├── Engagements: 3,241 (3.6% rate)
├── New Followers: 234
├── Top Tweet: Thread on AI agents — 12K impressions
└── 💡 Insight: Threads outperform single tweets 4x

Newsletter
├── Subscribers: 2,847 (+89 this week)
├── Open Rate: 42% (industry avg: 21%)
├── Click Rate: 8.3%
└── 💡 Insight: Tuesday 9am sends perform best

Website
├── Visitors: 3,421
├── Avg Time: 4:32
├── Top Page: /one-mind-framework
└── 💡 Insight: 67% traffic from Twitter

---

🎯 RECOMMENDATIONS:
1. Double down on YouTube Shorts (highest growth)
2. Create more threads (4x engagement)
3. Repurpose "Building JARVIS" into a thread
```

---

# 4. RESEARCH TOOLS

## `web_search`

**You say:**
> "What are the best practices for Letta agent memory management?"

**You receive:**
```
🔍 Search Results: Letta Memory Management

Top Sources Found:

1. Letta Documentation (docs.letta.com)
   "Memory blocks should be kept under 5000 tokens for optimal
   performance. Use archival memory for long-term storage..."

2. Letta Blog: "Sleeptime Compute"
   "Agents can consolidate and optimize memory during idle periods.
   Enable sleeptime for background memory management..."

3. GitHub Discussion #847
   "Best practice: Separate core memory (active context) from
   archival memory (reference). Use labels consistently..."

---

📋 SUMMARY:

Key Best Practices:
1. **Core Memory**: Keep under 5K tokens, frequently accessed context
2. **Archival Memory**: Long-term storage, searchable via embeddings
3. **Memory Labels**: Use consistent naming (human, persona, context)
4. **Sleeptime**: Enable for automatic consolidation
5. **Shared Blocks**: Use for multi-agent state sharing

---

Want me to:
├── 📥 Save this to IT knowledge base
├── 📋 Create implementation checklist
└── 🔍 Search for more specific topics
```

---

## `summarize_url`

**You say:**
> "Summarize this: https://www.letta.com/blog/programmatic-tool-calling"

**You receive:**
```
📄 Article Summary

Title: "Programmatic Tool Calling with Any LLM"
Source: Letta Blog
Reading Time: 8 min → Summary: 1 min

---

🎯 TL;DR:
Letta agents can write and execute Python code that calls
multiple tools, enabling complex autonomous workflows without
multiple LLM invocations.

---

📌 KEY POINTS:

1. **What It Is**
   Agents dynamically write scripts that call available tools,
   rather than being limited to one tool at a time.

2. **How It Works**
   Agent receives task → writes Python code → code calls tools
   → results returned → agent processes results

3. **Example Use Case**
   "Query database, loop through results, web search each item,
   aggregate findings" — all in one agent action.

4. **Benefits**
   - Complex workflows without human intervention
   - Reduced token usage (one script vs many tool calls)
   - Model-agnostic (works with any LLM)

---

💡 RELEVANCE TO ONE MIND:
This enables Legacy-AI to orchestrate morning briefings,
cross-pillar reports, and complex automations autonomously.

---

Saved to: 40-49 IT/43 Research/articles/letta-programmatic-tools.md
```

---

## `save_bookmark`

**You say:**
> "Bookmark this for later, it's about homesteading solar setups"

**You receive:**
```
🔖 Bookmark Saved

Title: "Complete Guide to Off-Grid Solar for Homesteads"
URL: https://example.com/homestead-solar-guide
Source: Homestead Magazine

---

Categorized:
├── Primary: LE → Sanctuary → Infrastructure
├── Secondary: IT → Research → Renewable Energy
└── Tags: #solar #homestead #off-grid #sustainability

---

Saved to: 20-29 LE/21 Sanctuary/research/solar-setup.md

Added to Reading List:
├── Priority: Medium
├── Estimated Read Time: 15 min
└── Related: [[Homestead Energy Plan]], [[Solar Calculation Notes]]

---

📋 Would you like me to:
├── 📖 Summarize it now
├── 📅 Schedule reading time
└── 🔗 Find related resources
```

---

# 5. COMMUNICATION TOOLS

## `send_email`

**You say:**
> "Send an email to the accountant asking about Q4 tax estimates"

**Legacy-AI processes:**
1. Recognizes: external communication → needs approval
2. Drafts email based on context
3. Looks up accountant contact from vault
4. Requests approval before sending

**You receive:**
```
📧 Email Draft Ready

To: sarah@taxfirm.com (Sarah Chen, CPA)
Subject: Q4 2025 Estimated Tax Payment

---

Hi Sarah,

Hope you had a great holiday! I wanted to check in about my Q4
estimated tax payment for Legacy AI LLC.

Could you let me know:
1. The recommended payment amount
2. Deadline (I believe it's Jan 15?)
3. If I need to adjust anything based on this quarter's revenue

Revenue for Q4 was approximately $34,000. Let me know if you
need any additional documentation.

Thanks,
Zeus

---

⚠️ APPROVAL REQUIRED (External Communication)

This email will be sent to an external contact.

[✅ Send Now] [✏️ Edit First] [❌ Cancel]
```

**After approval:**
```
✅ Email Sent

To: sarah@taxfirm.com
Sent: 2:47 PM EST

Logged to: 30-39 GE/37 Compliance/correspondence/2025-12-26-accountant.md

📋 Follow-up task created:
"Check for accountant response re: Q4 taxes"
Due: Dec 30 (if no response by then)
```

---

## `schedule_call`

**You say:**
> "Schedule a call with mom for Sunday afternoon"

**You receive:**
```
📞 Scheduling Call: Mom

Your availability Sunday (Dec 29):
├── 1:00 PM - 3:00 PM ✅ Free
├── 3:00 PM - 4:00 PM ❌ "Review weekly plan"
└── 4:00 PM - 6:00 PM ✅ Free

Mom's preferences (from memory):
├── Prefers afternoon calls
├── Usually free after 2 PM
└── Last call: Dec 15 (11 days ago)

---

Suggested: Sunday 2:00 PM (1 hour)

---

[✅ Create Event] [📱 Send her a text to confirm] [🔄 Different time]
```

**If you confirm:**
```
✅ Call Scheduled

📅 Sunday, Dec 29 at 2:00 PM
📞 Call with Mom
⏱️ Duration: 1 hour

Added to calendar ✓
Reminder set: 30 min before ✓

Want me to text her to confirm the time?
```

---

## `send_sms`

**You say:**
> "Text mom that I'll call her Sunday at 2"

**You receive:**
```
📱 SMS Ready

To: Mom (+1 555-123-4567)

Message:
"Hey Mom! Can I call you Sunday around 2pm? Been too long! 🙂"

Character count: 62/160

---

⚠️ APPROVAL REQUIRED (External Communication)

[✅ Send] [✏️ Edit] [❌ Cancel]
```

**After approval:**
```
✅ SMS Sent

To: Mom
Sent: 2:52 PM

Logged to: 20-29 LE/22 Family/communication-log.md

💡 Reminder: You mentioned wanting to discuss holiday plans.
I've added that to your call notes for Sunday.
```

---

# 6. DATA TOOLS

## `query_notion`

**You say:**
> "What clients do I have in my Notion CRM that haven't been contacted in 30 days?"

**You receive:**
```
📊 CRM Query Results

Clients Not Contacted in 30+ Days:

1. Acme Corp
   ├── Last Contact: Nov 18 (38 days)
   ├── Status: Active client
   ├── Value: $5,000/mo retainer
   └── ⚠️ Risk: High value, needs attention

2. StartupXYZ
   ├── Last Contact: Nov 22 (34 days)
   ├── Status: Proposal sent
   ├── Value: $12,000 project
   └── ⚠️ Risk: Deal might go cold

3. Local Business Co
   ├── Last Contact: Nov 10 (46 days)
   ├── Status: Past client
   ├── Value: $2,400 (past)
   └── 💡 Opportunity: Upsell potential

---

📋 RECOMMENDED ACTIONS:

1. [Create task] Follow up with Acme Corp (high priority)
2. [Create task] Chase StartupXYZ proposal
3. [Create task] Reach out to Local Business for Q1

[Create all tasks] [Pick specific ones] [Export to CSV]
```

---

## `update_airtable`

**You say:**
> "Add a new lead to Airtable: John Smith from TechCorp, interested in consulting, met at conference"

**You receive:**
```
✅ Lead Added to Airtable

Table: Sales Pipeline
Record ID: rec123abc

---

Fields Populated:
├── Name: John Smith
├── Company: TechCorp
├── Interest: Consulting
├── Source: Conference
├── Status: New Lead
├── Created: Dec 26, 2025
└── Owner: Zeus

---

🔗 View in Airtable: [link]

---

📋 Suggested Next Steps:
1. [Create task] Send intro email to John
2. [Create task] Research TechCorp before outreach
3. [Add note] Conference details/context

[Do all] [Pick specific] [Done for now]
```

---

## `sync_spreadsheet`

**You say:**
> "Sync my monthly expenses from the spreadsheet to my vault"

**You receive:**
```
🔄 Spreadsheet Sync Complete

Source: "2025 Expense Tracker" (Google Sheets)
Destination: 30-39 GE/36 Treasury/expenses/

---

December 2025 Synced:
├── Total Expenses: $8,432.17
├── Categories: 12
├── Transactions: 47
└── New since last sync: 8 transactions

---

By Category:
├── Software/SaaS: $1,247.00
├── Contractors: $3,500.00
├── Marketing: $892.00
├── Office: $234.17
├── Travel: $1,459.00
└── Other: $1,100.00

---

Updated Files:
├── expenses/2025-12.md ✓
├── reports/monthly-summary.md ✓
└── Shared memory: ge_metrics ✓

---

⚠️ Anomaly Detected:
"Travel" is 73% higher than your 3-month average.
Want me to break down the travel expenses?
```

---

# 7. ORCHESTRATION (Chained Tool Calls)

The real power is when Legacy-AI **chains these together** autonomously.

## `morning_briefing`

**Triggered:** 5:30 AM daily (n8n schedule)

**Legacy-AI executes programmatically:**
```
1. calendar_get_events("today")           → 3 meetings found
2. todoist_get_tasks("today")             → 7 tasks due
3. get_bank_balance()                      → Cash position checked
4. check_analytics()                       → Content performance pulled
5. cross_pillar_status()                   → All workers report in
6. get_hp_metrics()                        → Sleep/energy data
7. get_pending_approvals()                 → Items needing attention
```

**You receive in Discord:**
```
🌅 Good morning, Zeus
📅 Friday, December 27, 2025

---

🎯 TODAY'S TOP 3:
1. [GE] Client call with Acme Corp — 10 AM
2. [IT] Deploy n8n server to Coolify
3. [HP] Leg day workout

---

📆 CALENDAR:
• 10:00 AM - Acme Corp quarterly review
• 2:00 PM - Content planning session
• 4:00 PM - Family video call

---

💪 HP STATUS:
• Sleep: 7.2h (below 8h target)
• Recovery: 72%
• Workout streak: 5 days
💡 Suggestion: Moderate intensity today given recovery score

---

💰 GE STATUS:
• December Revenue: $11,132 (89% of goal)
• Outstanding invoices: $2,100
• Runway: 4.2 months

---

⏳ PENDING APPROVALS (2):
1. Send email to accountant re: Q4 taxes
2. Publish blog post to Instagram

---

⚠️ ALERTS:
• Acme Corp not contacted in 38 days (call today!)
• Blog post deadline tomorrow

---

Have a great day! 💪
```

---

## `prepare_for_monday`

**You say:**
> "Prepare for my Monday"

**Legacy-AI chains:**
```
1. calendar_get_events("monday")
2. todoist_get_tasks("monday")
3. query_notion("meetings this week")     → Client context
4. vault_read_file("weekly-goals.md")     → Goals reviewed
5. cross_pillar_status()                   → System health
6. Generates prep document
```

**You receive:**
```
📋 Monday Prep Complete

Your Monday at a Glance:
├── 3 meetings (first at 9 AM)
├── 7 tasks due
├── Revenue: On track
└── System health: 78%

Meeting Prep Notes:
1. 9 AM - Team standup
   • Review: Sprint progress
   • Discuss: Q1 planning

2. 11 AM - Client call (TechCorp)
   • Context: Proposal sent Nov 22, no response
   • Goal: Close the deal or get clear next steps
   • Prep: Review proposal, check competitor pricing

3. 3 PM - Content review
   • Pending: 2 blog drafts, 1 video script
   • Priority: "PARA vs One Mind" post

---

📄 Full prep doc saved to:
00-09 Command/01 Zeus/daily/2025-12-30-monday-prep.md

[View full prep] [Add tasks] [Anything else?]
```

---

# Tool Registry Summary

| Category | Tools | Pillar | API Dependencies |
|----------|-------|--------|------------------|
| **Financial** | `get_bank_balance`, `categorize_expense`, `revenue_report` | GE | Plaid, Stripe, PayPal |
| **Home Automation** | `set_thermostat`, `check_security`, `water_plants` | LE/IT | Home Assistant |
| **Content** | `draft_blog_post`, `schedule_social`, `check_analytics` | GE | Buffer, YouTube, Twitter API |
| **Research** | `web_search`, `summarize_url`, `save_bookmark` | All | Exa/Perplexity, Readability |
| **Communication** | `send_email`, `schedule_call`, `send_sms` | All | Gmail, Google Calendar, Twilio |
| **Data** | `query_notion`, `update_airtable`, `sync_spreadsheet` | All | Notion, Airtable, Google Sheets |
| **Multi-Agent** | `delegate_to_*`, `broadcast_to_all`, `cross_pillar_status` | Command | Letta internal |
| **Guardian** | `request_approval`, `check_spending_limit`, `validate_external` | Command | Internal |

---

# Letta Compatibility

All of this is **100% doable with Letta**:

| Capability | Letta Support |
|------------|---------------|
| Custom Python tools | ✅ Native |
| Multi-agent communication | ✅ Built-in |
| Shared memory blocks | ✅ Native |
| External API calls | ✅ Tools can call any API |
| Programmatic tool chaining | ✅ Agents write code that calls tools |
| MCP integration | ✅ Connect external MCP servers |
| Autonomous execution | ✅ Agents decide which tools to use |
| Background processing | ✅ Sleeptime compute |

---

*Part of the One Mind Operating System*
*Created: 2025-12-27*
