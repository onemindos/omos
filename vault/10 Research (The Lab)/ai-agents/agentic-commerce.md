---
title: "Agentic Commerce — Landscape Research"
type: research
status: monitoring
created: 2026-06-08
updated: 2026-06-08
owner: zeus
category: ai-commerce
tags: [research, ai, agents, commerce, payments, stripe, visa, mastercard, mcp, a2a]
ai_summary: "AI agents acting autonomously on behalf of consumers/businesses to browse, compare, negotiate, and complete purchases. Major payment networks (Visa, Mastercard, Stripe) building agent-native infrastructure in 2025."
---

# Agentic Commerce

> AI agents acting autonomously on behalf of consumers or businesses to perform commerce tasks — browsing, comparing, negotiating, purchasing, managing returns — with minimal human involvement per transaction.

## What It Is

Traditional e-commerce is built for humans: search, click, checkout. Agentic commerce shifts that active role to AI agents. Humans set goals, budgets, and rules — agents execute end-to-end.

**Three layers (McKinsey):**

| Layer | Description |
|---|---|
| **Consumer Agents** | Personal shoppers — browse, compare, buy on a person's behalf |
| **Business Agents** | Manage inventory, pricing, supplier negotiation, customer service |
| **Agent-to-Agent (A2A)** | Automated B2B commerce — buyer agent negotiates directly with seller agent |

---

## Payment Infrastructure

### Visa — Intelligent Commerce (May 2025)
First major payments network to build infrastructure for AI agents. Authenticates agents the same way it authenticates human cardholders using existing global trust network.

- **Authentication** — AI agents authenticated like human cardholders
- **Controls** — consumers set spending limits and rules per agent
- **Embedded Credentials** — Visa credentials embedded directly in AI agent platforms
- **Programmatic Payments** — agents request, issue, and pay with Visa credentials dynamically
- **Partners:** Anthropic, Microsoft, OpenAI, Samsung, Perplexity, Stripe
- **Commercial pilots:** Summer 2025

Links:
- https://investor.visa.com/news/news-details/2025/Visa-Introduces-New-Capabilities-to-Advance-Agentic-Commerce/default.aspx
- https://www.visa.com/visa-intelligent-commerce/index.html
- https://developer.visa.com/visa-intelligent-commerce

---

### Mastercard — Agent Pay (May 2025)
Tokenized agent identities — each AI agent gets a unique verifiable digital ID on the Mastercard network.

- Tokenized agent identities
- Agent-to-agent payment flows
- Consumer-controlled spending rules
- **Launch partner:** Microsoft Azure AI Foundry

Links:
- https://www.mastercard.com/news/press/2025/may/mastercard-introduces-agent-pay-to-power-next-generation-of-agentic-payments/
- https://investor.mastercard.com/news/news-details/2025/Mastercard-Introduces-Agent-Pay-to-Power-Next-Generation-of-Agentic-Payments/default.aspx

---

### Stripe — Open Payments Protocol (April 2025)
Most significant **open protocol** move in the space. Built on MCP — designed as an open standard, not Stripe-proprietary.

- Merchants publish machine-readable product catalogs via MCP
- AI agents discover products and complete purchases without human login flows
- Open specification — all payment processors invited to adopt
- Collaborating with OpenAI, Anthropic to finalize spec

Links:
- https://techcrunch.com/2025/04/25/stripe-launches-open-payments-protocol-for-ai-agents/
- https://stripe.com/blog/from-chatbots-to-checkout-how-stripe-is-remaking-payments-for-the-agent-era
- https://docs.stripe.com/stripe-apps/mcp
- https://a16z.com/stripes-open-payments-protocol-for-ai-agents-and-what-it-might-mean/

---

### PayPal
Agentic checkout flows using stored PayPal credentials, Smart Receipts, Intelligent Offers Engine.

- https://newsroom.paypal.com/news-details/2025/PayPal-AI-Agents-Agentic-Commerce/default.aspx

---

## AI Agent Platforms (Consumer-Facing)

| Platform | Agent | What It Does |
|---|---|---|
| OpenAI | Operator | Browser agent — navigates web, fills forms, completes checkouts |
| Google | Project Mariner | Research agent — web browsing including shopping |
| Amazon | Rufus + Buy for Me | Rufus = assistant; Buy for Me = buys from third-party sites autonomously |
| Perplexity | Buy with Perplexity | In-chat purchasing (live Nov 2024, US) |
| Anthropic | Claude Computer Use | Navigates real UIs; fills checkout forms |
| Klarna | AI Shopping Assistant | BNPL-integrated shopping agent |

Links:
- https://openai.com/index/introducing-operator/
- https://deepmind.google/technologies/project-mariner/
- https://www.theverge.com/2025/3/5/24083917/amazon-buy-for-me-ai-agent
- https://techcrunch.com/2024/11/18/perplexity-introduces-an-in-app-shopping-and-checkout-feature-for-us-users/

---

## Commerce Platform Players

| Platform | Role |
|---|---|
| **Shopify** | Agent-accessible product catalogs + checkout; Sidekick AI for merchants |
| **Salesforce Agentforce** | Commerce Cloud agents — orders, returns, B2B quoting |
| **SAP Joule** | Enterprise procurement agents — POs, supplier negotiation, invoices |
| **IBM WatsonX** | Retail agents — recommendations, pricing, supply chain |

Links:
- https://www.theverge.com/2025/2/26/24083937/shopify-ai-agent-commerce-storefronts
- https://www.salesforce.com/news/stories/agentforce-commerce-retail-ai/

---

## Open Protocols

### Model Context Protocol (MCP) — Anthropic
The foundational open standard for connecting AI agents to data sources and tools. Commerce MCP servers: Stripe, Shopify, Amazon product search. Stripe's payments protocol is built on MCP.

- https://modelcontextprotocol.io/introduction

### Agent2Agent (A2A) — Google
Open protocol for agent-to-agent communication. Buyer agents communicate with merchant agents to negotiate, check inventory, transact. Agents publish "agent cards" advertising capabilities.

- https://developers.googleblog.com/en/a2a-a-new-era-of-agent-interoperability/
- https://github.com/google/A2A

> **MCP + A2A:** MCP = agent↔tool (connecting to Stripe, catalogs). A2A = agent↔agent (buyer agent ↔ seller agent).

---

## Agent Frameworks for Commerce

| Framework | Notes |
|---|---|
| LangChain | Python/JS, widely used, integrates with Stripe |
| LlamaIndex | Data retrieval + RAG, product search agents |
| AutoGen (Microsoft) | Multi-agent conversation + task automation |
| CrewAI | Role-based multi-agent orchestration |
| OpenAI Assistants API | Hosted, tool-calling |
| Amazon Bedrock Agents | AWS-hosted enterprise infrastructure |

---

## Key Challenges

| Challenge | Detail |
|---|---|
| Agent Authentication | How do merchants verify the agent has real authorization? |
| Authorization Scope | What is the agent allowed to buy? |
| Fraud Detection | Agents behave differently from humans |
| Liability | Who is responsible when agent buys wrong? |
| Brand Disruption | Agents optimize for specs, not brand loyalty |
| Privacy/GDPR | Explicit human authorization required before agents act |

---

## Relevance to OneMind OS

- **Stripe MCP + Visa Intelligent Commerce** = payment rails Legacy/agents could use for autonomous procurement (supplies, subscriptions, fleet hardware)
- **A2A protocol** = agent-to-agent commerce aligns with OneMind's multi-agent architecture
- **Education business (GE)** — agentic commerce will reshape how Wealth Warriors / Skool students shop for courses and tools; needs monitoring
- **Agent Experience Optimization (AEO)** — if OneMind ever sells digital products, AEO (making products machine-readable) becomes as important as SEO

---

## Reading List

- https://www.mckinsey.com/capabilities/mckinsey-digital/our-insights/welcome-to-the-era-of-agentic-commerce
- https://www.bloomreach.com/en/blog/agentic-commerce
- https://www.elastic.co/blog/agentic-commerce
- https://venturebeat.com/ai/visas-agentic-commerce-move-a-strategic-look-at-the-payments-giant-s-ai-push/
- https://www.ibm.com/think/insights/agentic-commerce
- https://www.infoq.com/articles/mcp-a2a-protocol-comparison/
