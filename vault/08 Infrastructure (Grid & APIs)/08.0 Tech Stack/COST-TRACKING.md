# Cost Tracking - Monthly Expenses

> **Monitor recurring costs for OneMind OS infrastructure**

## Current Monthly Costs

### Total: ~$200/month

---

## AI & Language Models ($80/month)

| Service         | Purpose            | Monthly Cost | Billing Type |
| --------------- | ------------------ | ------------ | ------------ |
| **AWS Bedrock** | Claude Opus 4.5    | ~$50         | Pay-per-use  |
| **OpenAI API**  | GPT-4, embeddings  | ~$30         | Pay-per-use  |
| **Composio**    | OAuth integrations | $0           | Free tier    |

**Notes:**

- Claude cost varies by usage (~500K tokens/month)
- OpenAI backup for when Claude unavailable

---

## Development Tools ($20/month)

| Service            | Purpose         | Monthly Cost | Billing Cycle      |
| ------------------ | --------------- | ------------ | ------------------ |
| **Cursor**         | AI code editor  | $20          | Monthly            |
| **GitHub Copilot** | Code completion | $0           | Included in GitHub |
| **VS Code**        | Code editor     | $0           | Free               |

---

## Infrastructure ($50/month)

| Service        | Purpose            | Monthly Cost | Notes         |
| -------------- | ------------------ | ------------ | ------------- |
| **AWS EC2**    | Compute            | ~$30         | t3.medium     |
| **AWS RDS**    | PostgreSQL         | ~$15         | db.t3.micro   |
| **AWS S3**     | Storage            | ~$2          | Minimal usage |
| **Vercel**     | Web hosting        | $0           | Free tier     |
| **Docker Hub** | Container registry | $0           | Free tier     |

**Notes:**

- Can optimize by shutting down dev instances when not in use
- RDS might be replaced with self-hosted PostgreSQL

---

## Productivity & Tools ($10/month)

| Service       | Purpose         | Monthly Cost | Billing Cycle |
| ------------- | --------------- | ------------ | ------------- |
| **1Password** | Secrets manager | $5           | Monthly       |
| **Todoist**   | Task management | $0           | Free tier     |
| **Linear**    | Issue tracking  | $0           | Free tier     |
| **Obsidian**  | Knowledge base  | $0           | Local-first   |
| **Tailscale** | VPN             | $0           | Free tier     |

---

## Communication ($0/month - Self-Hosted)

| Service     | Purpose        | Monthly Cost | Notes       |
| ----------- | -------------- | ------------ | ----------- |
| **LiveKit** | Voice/video    | $0           | Self-hosted |
| **NATS**    | Message broker | $0           | Self-hosted |
| **MQTT**    | IoT messaging  | $0           | Self-hosted |

---

## Business & CRM ($40/month)

| Service       | Purpose            | Monthly Cost | Notes          |
| ------------- | ------------------ | ------------ | -------------- |
| **Circle.so** | Community platform | ~$40         | Varies by plan |
| **Zoho CRM**  | Customer mgmt      | $0           | Free tier      |

---

## Planned Additions (Not Yet Active)

### Transportation APIs ($0/month)

| Service         | Purpose       | Estimated Cost | Status                         |
| --------------- | ------------- | -------------- | ------------------------------ |
| **Uber API**    | Ride ordering | $0             | 🔄 API free, rides cost normal |
| **Lyft API**    | Ride ordering | $0             | 🔄 API free, rides cost normal |
| **Google Maps** | Navigation    | $0             | 🔄 $200 credit/month           |
| **MTA API**     | NYC transit   | $0             | 🔄 Free                        |

### Research & Search APIs ($25/month)

| Service            | Purpose         | Estimated Cost | Status             |
| ------------------ | --------------- | -------------- | ------------------ |
| **Arxiv**          | Academic papers | $0             | 🔄 Free            |
| **Brave Search**   | Web search      | $0             | 🔄 2000 free/month |
| **Exa (Metaphor)** | Semantic search | $10            | 🔄 Planned         |
| **Tavily**         | AI search       | $10            | 🔄 Planned         |
| **Wolfram Alpha**  | Computation     | $5             | 🔄 Planned         |

### Web Scraping ($20/month)

| Service         | Purpose        | Estimated Cost | Status         |
| --------------- | -------------- | -------------- | -------------- |
| **Firecrawl**   | Web scraping   | $20            | 🔄 Planned     |
| **Jina Reader** | URL extraction | $0             | 🔄 Free tier   |
| **Crawl4AI**    | Web crawler    | $0             | 🔄 Self-hosted |
| **SearXNG**     | Meta-search    | $0             | 🔄 Self-hosted |

---

## Cost Breakdown by Category

| Category           | Current  | Planned | Future Total |
| ------------------ | -------- | ------- | ------------ |
| **AI/ML**          | $80      | $0      | $80          |
| **Infrastructure** | $50      | $0      | $50          |
| **Development**    | $20      | $0      | $20          |
| **Business/CRM**   | $40      | $0      | $40          |
| **Productivity**   | $10      | $0      | $10          |
| **Transportation** | $0       | $0      | $0           |
| **Research**       | $0       | $25     | $25          |
| **Web Scraping**   | $0       | $20     | $20          |
| **Communication**  | $0       | $0      | $0           |
| **TOTAL**          | **$200** | **$45** | **$245**     |

---

## Annual Costs

**Current:** $2,400/year  
**With Planned Additions:** $2,940/year (~$245/month)

---

## Cost Optimization Opportunities

### Immediate Savings

1. **Self-host PostgreSQL** - Save $15/month (AWS RDS → local)
2. **Optimize EC2** - Use spot instances or smaller instance types
3. **Use AWS free tier** - Student credits if eligible

### Potential Savings: $30-50/month

### Future Considerations

1. **Consolidate AI models** - Pick Claude OR OpenAI (not both)
2. **Free alternatives** - Use SearXNG instead of paid search APIs
3. **Self-host more** - Reduce cloud dependencies

---

## Budget Alerts

| Threshold      | Action                        |
| -------------- | ----------------------------- |
| **$250/month** | Review newest additions       |
| **$300/month** | Pause non-essential services  |
| **$350/month** | Emergency optimization needed |

---

## Historical Costs

| Month    | Total Cost | Notes                      |
| -------- | ---------- | -------------------------- |
| Jan 2026 | $200       | Baseline established       |
| Dec 2025 | $180       | Before Cursor subscription |
| Nov 2025 | $165       | Lower AI usage             |

---

## Payment Methods

- **AWS**: Personal credit card (autopay)
- **OpenAI**: Personal credit card
- **Cursor**: Personal credit card
- **1Password**: Personal credit card
- **Circle.so**: Business credit card

---

## Tax Deductions (Business Expenses)

Most of these costs qualify as business expenses for freelance/consulting work:

- ✅ AI API costs (tools of trade)
- ✅ Cloud infrastructure (hosting)
- ✅ Development tools (professional development)
- ✅ CRM/business tools (customer management)

**Estimated Annual Deduction:** ~$2,500-3,000

---

Last Updated: January 25, 2026  
Next Review: February 1, 2026
