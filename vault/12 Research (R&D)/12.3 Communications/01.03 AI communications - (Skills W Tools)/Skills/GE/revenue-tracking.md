---
id: ui-15-20260111053926
title: "Skill: Revenue Tracking"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/15]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Skill: Revenue Tracking"
---

# Skill: Revenue Tracking

## Purpose

Monitor and report on business revenue metrics across all Legacy ventures. Provide real-time visibility into financial health and identify trends.

## Knowledge

### Revenue Categories

| Category | Source | Tracking Frequency |
|----------|--------|-------------------|
| MRR (Monthly Recurring) | Subscriptions, memberships | Daily |
| One-time Sales | Products, courses | As occurs |
| Service Revenue | Consulting, done-for-you | Per project |
| Affiliate/Partnership | Commission income | Monthly |
| Investment Returns | Capital gains, dividends | Quarterly |

### Key Metrics

| Metric | Definition | Target |
|--------|------------|--------|
| MRR | Monthly recurring revenue | Growth |
| ARR | Annual recurring revenue (MRR × 12) | $X |
| Churn Rate | % customers lost per month | <5% |
| LTV | Lifetime value per customer | Increasing |
| CAC | Customer acquisition cost | Decreasing |
| LTV:CAC | Ratio of value to acquisition | 3:1+ |

### Business Entities

```
Legacy Platform Holdings LLC
├── Legacy Acadacom — Education (primary revenue)
├── Legacy Restore — Credit Repair
├── Legacy Capital — Investments
└── Legacy Productions — Media
```

### Data Sources

- **Stripe** — Payment processing, MRR
- **Circle** — Membership counts, engagement
- **Zoho CRM** — Pipeline value, deal stage
- **Bank accounts** — Cash position
- **Spreadsheets** — Manual tracking backup

## Procedures

### Daily Revenue Check

1. Query Stripe for yesterday's transactions
2. Calculate current MRR
3. Check for new subscriptions and cancellations
4. Calculate net change from previous day
5. Update GE Context memory block
6. Flag significant changes (>10% swing)

### Weekly Revenue Report

1. Compile daily data for the week
2. Calculate:
   - Total revenue (MRR + one-time)
   - New vs churned customers
   - Average revenue per user
   - Week-over-week change
3. Compare to same week last month/year
4. Identify top revenue sources
5. Generate summary report

### Monthly Close

1. Pull final numbers from all sources
2. Reconcile with bank statements
3. Calculate:
   - Total monthly revenue
   - Gross margin
   - MRR growth rate
   - Churn analysis
4. Compare to budget/forecast
5. Generate monthly financial summary
6. Queue for Zeus review

### Pipeline Health Check

1. Query Zoho CRM for open deals
2. Calculate total pipeline value
3. Segment by stage (qualified, proposal, negotiation, closing)
4. Identify stalled deals (no activity 14+ days)
5. Flag high-value opportunities
6. Return summary with action items

## Context

- Cash flow is king — always know the runway
- MRR is the north star metric
- Revenue without profit is vanity
- Track leading indicators (pipeline, trials) not just lagging (revenue)
- Financial decisions > $500 require Zeus approval

## Examples

**Input:** "What's our revenue this month?"

**Response:**
```
💰 Revenue Summary — December 2024 (as of Dec 23)

**MRR:** $4,200 (+3% vs November)
- Legacy Acadacom: $3,500
- Legacy Restore: $700

**One-time Revenue:** $1,200
- Course sales: $800
- Consulting: $400

**Total MTD:** $5,400

**Key Changes:**
- +4 new members (Acadacom)
- -1 cancellation (Restore)
- Net MRR change: +$180

**Pipeline:** $12,500 in active deals
- 2 deals in negotiation stage
- 1 proposal sent yesterday

On track to beat November by ~8%.
```

---

**Input:** "Show me churn for this quarter"

**Response:**
```
📉 Churn Analysis — Q4 2024

| Month | Starting | Lost | Churn Rate |
|-------|----------|------|------------|
| Oct | 45 | 2 | 4.4% |
| Nov | 47 | 3 | 6.4% |
| Dec* | 48 | 1 | 2.1% |

*December partial (as of Dec 23)

**Q4 Average:** 4.3% (Target: <5%) ✓

**Churn Reasons:**
- Financial constraints: 3
- No longer needed: 2
- Competitor: 1

**Recommendation:** The November spike came from 2 customers citing financial reasons. Consider adding a pause option for members facing temporary hardship.
```

## Integration

- **Tools:** `zoho_crm`, `stripe` (if configured), `obsidian_write`
- **Memory:** GE Context (read/write)
- **Escalation:** Revenue drop >20% → Escalate to LEGACY immediately
