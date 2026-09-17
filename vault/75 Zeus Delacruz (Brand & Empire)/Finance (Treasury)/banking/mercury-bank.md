# Mercury Bank

**Website:** https://mercury.com  
**Type:** Business Banking Platform for Startups & Tech Companies  
**Founded:** 2017  
**Founders:** Immad Akhund (CEO), Max Tagher (CTO), Jason Zhang (COO)  
**Status:** 🟡 **Researched** | Not yet applied  

---

## Overview

Mercury is an all-in-one business banking and money management platform designed for startups, small businesses, and scaling companies. Known for its modern UX, developer-friendly API, and comprehensive financial stack.

**Tagline:** "Everything you do with money. All in one place."

---

## Key Features

### 💳 **Banking Products**
- **Business Checking** — $0 monthly fee, $0 minimum balance
- **Business Savings** — $0 monthly fee, $0 minimum balance
- **Treasury** — Up to **3.66% yield** via Mercury Advisory (SEC-registered investment advisor)
- **Credit Card** — **1.5% cashback** on all business spend
- **Loans** — Venture debt & working capital loans for growth-stage companies

### 💸 **Payments & Invoicing**
- **Bill Pay** — Free USD domestic payments, international wire support
- **Invoicing** — Create free invoices, get paid directly to Mercury account
- **Business Payments** — Send money around the world

### 👥 **Team & Expense Management**
- **Unlimited Virtual Cards** — Issue cards for team members
- **Physical Cards** — Metal debit cards for founders
- **Spending Limits** — Set per-card limits and approval workflows
- **Expense Categories** — Automatic categorization for accounting

### 📊 **Accounting & Integrations**
- **QuickBooks Sync** — Real-time transaction sync
- **Xero Sync** — Two-way accounting integration
- **NetSuite Sync** — Enterprise-grade ERP integration
- **AI Automations** — Auto-categorization, invoice matching, reconciliation

### 🔧 **Developer Tools**
- **REST API** — Full programmatic access to banking features
- **Webhooks** — Real-time event notifications
- **Developer Docs** — Comprehensive API documentation
- **Sandbox Environment** — Test integrations before production

### 📱 **Mobile & Web**
- **iOS App** — 4.9/5 rating (10,730 reviews)
- **Android App** — Available
- **Web Dashboard** — Full-featured desktop interface

---

## Banking Partners

- **Choice Financial Group** (FDIC insured)
- **Evolve Bank & Trust** (FDIC insured)
- All deposits insured up to $250,000 per depositor, per bank

---

## Target Market

- ✅ VC-backed startups
- ✅ Tech companies & SaaS businesses
- ✅ Scaling companies ($1M-$100M+ ARR)
- ✅ Businesses needing API integrations
- ✅ Companies raising or managing investor capital

**Customer Base:** 300,000+ businesses including Supabase, Y Combinator startups, and other high-growth companies

---

## Pricing

| Product | Cost |
|---------|------|
| **Business Checking** | $0/month |
| **Business Savings** | $0/month |
| **Treasury** | $0/month (yield varies) |
| **Credit Card** | $0 annual fee |
| **Domestic Wire (Outgoing)** | $0 |
| **International Wire (Outgoing)** | Varies by currency |
| **Physical Debit Card** | $0 |
| **Virtual Cards** | $0 (unlimited) |

**Note:** No monthly fees, no minimum balances, no hidden charges on standard products.

---

## Strengths ✅

- **Best-in-class UX** — Modern, intuitive interface (better than traditional banks)
- **API-first architecture** — Full developer access for integrations
- **High treasury yields** — 3.66% competitive with Brex, Ramp
- **1.5% cashback credit card** — Strong rewards for business spend
- **Fast onboarding** — Apply in 10 minutes, approval within 1-3 days
- **Comprehensive stack** — Banking + payments + invoicing + expense management in one platform
- **No personal guarantee** — Some products don't require personal credit check
- **Startup-friendly** — Understands venture-backed business models

---

## Weaknesses ❌

- **Online-only** — No physical branches (not ideal if you need in-person banking)
- **Best for incorporated businesses** — Not optimized for freelancers or sole proprietors
- **Some features require thresholds** — Certain products need minimum revenue or balance
- **Limited cash deposit options** — No direct cash deposit (use third-party services)
- **Venture debt selective** — Loan products mainly for funded startups

---

## Application Requirements

### Documents Needed:
- ✅ **EIN (Employer Identification Number)**
- ✅ **Business formation documents** (Articles of Incorporation, Operating Agreement)
- ✅ **Beneficial owner information** (name, DOB, SSN, address for 25%+ owners)
- ✅ **Business address** (must be U.S.-based)
- ✅ **Business description** (what you do, revenue model)
- ✅ **Expected monthly transaction volume**

### Eligibility:
- ✅ U.S.-incorporated business (C-Corp, S-Corp, LLC, Partnership)
- ✅ Operating in legal industry (no cannabis, adult content, gambling)
- ✅ Valid business purpose
- ❌ Must not be in sanctioned countries

**Approval Timeline:** 1-3 business days for most applications

---

## OneMind OS Use Case

### Recommended Account Structure:

```
Mercury Account
├── Operating Checking (main runway)
│   └── For: Monthly expenses, vendor payments, software subscriptions
│
├── Payroll Checking (team payments)
│   └── For: Employee/contractor compensation, payroll taxes
│
├── Treasury Account (cash reserves)
│   └── For: 3-6 months operating expenses earning 3.66% yield
│
└── Credit Card (1.5% cashback)
    └── For: All operational expenses, AWS/infrastructure, tools
```

### Integration Points:
- **1Password** — Store credentials securely
- **Homepage Dashboard** — Display account balances
- **QuickBooks/Xero** — Accounting sync (if/when implemented)
- **NATS Events** — Future: Banking alerts via OneMind OS event bus
- **Mercury API** — Programmatic access for custom OneMind integrations

---

## Comparison to Alternatives

| Feature | Mercury | Relay | Lili |
|---------|---------|-------|------|
| **Best For** | Tech startups | Profit First SMBs | Freelancers |
| **Treasury Yield** | 3.66% | Standard | 4%+ |
| **Credit Card** | ✅ 1.5% cashback | ❌ | ❌ |
| **Loans** | ✅ Venture debt | ❌ | ❌ |
| **API Access** | ✅ Full REST API | Limited | ❌ |
| **Team Features** | ✅ Advanced | ✅ Good | ❌ |
| **Accounting Sync** | QB/Xero/NetSuite | QB/Xero | Basic |

**Winner for OneMind OS:** Mercury (most comprehensive feature set + API access)

---

## Next Steps

### Phase 1: Application ⏳
- [ ] Gather required documents (EIN, formation docs, beneficial owner info)
- [ ] Complete online application at mercury.com
- [ ] Wait for approval (1-3 business days)

### Phase 2: Account Setup ⏳
- [ ] Fund initial deposit (if required)
- [ ] Order physical debit card
- [ ] Set up Treasury account for cash reserves
- [ ] Apply for credit card (after checking account approved)

### Phase 3: Integration ⏳
- [ ] Store credentials in 1Password
- [ ] Add account to Homepage dashboard
- [ ] Connect to QuickBooks/Xero (if applicable)
- [ ] Set up API access for OneMind integrations

### Phase 4: Optimization ⏳
- [ ] Move operating cash to Treasury (earn 3.66% yield)
- [ ] Issue virtual cards for recurring subscriptions
- [ ] Set up expense categories and workflows
- [ ] Configure NATS event triggers for banking alerts

---

## Resources

- **Website:** https://mercury.com
- **Application:** https://mercury.com/start
- **API Docs:** https://docs.mercury.com
- **Support:** help@mercury.com
- **Status Page:** https://status.mercury.com

---

## Notes

*Created: 2026-05-23*  
*Last Updated: 2026-05-23*  
*Status: Researched, pending application*  
*Owner: Zeus Delacruz*  
*Entity: OneMind OS (primary business account)*

---

**🎯 Recommendation:** Mercury is the **top choice** for OneMind OS due to API access, treasury yields, credit card rewards, and startup-friendly features. Apply as primary business banking platform.
