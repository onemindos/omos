---
title: "Tools & APIs Registry"
type: reference
status: active
created: 2026-04-09
owner: zeus
tags: [type/reference, domain/so-03, tools, apis]
ai_summary: "Master index of all tools and APIs available to One Mind agents — status, credentials, skill availability, and domain routing."
---

# Tools & APIs Registry

> Single source of truth for every tool and API in the One Mind ecosystem.
> Agents: check here to know what you can use and how.

---

## Status Key

| Status | Meaning |
|--------|---------|
| ✅ Active | Skill exists, credentials configured, ready to use |
| 🔧 Configured | API access exists but no skill yet |
| 📋 Planned | On the roadmap, not yet set up |
| ❌ Deprecated | Was used, no longer active |

---

## Active Skills (`.claude/skills/`)

These have full SKILL.md implementations with API patterns, auth, and usage docs.

| Tool | Skill Name | Domain | Credentials |
|------|-----------|--------|-------------|
| Home Assistant | `home-assistant` | SO 03 / LE 51 | `op://SO/Home Assistant/token` |
| Todoist | `todoist` | SO 09 | `op://SO/Todoist/api-token` |
| Telegram | `telegram` | SO 03 | `op://SO/Telegram Bot/token` |
| Discord | `discord` | SO 03 | `op://SO/Discord Webhook/url` |
| GitHub | `github` | SO 10 | `gh` CLI auth |
| Letta | `letta` | SO 03 | `op://SO/Letta/api-key` |
| N8N | `n8n` | SO 09 | `op://SO/N8N/api-key` |
| 1Password | `1password` | SO 11 | `op` CLI auth |
| Zoho (CRM/Desk/Calendar) | `zoho` | GE 75+ | `op://SO/Zoho/client-id` |
| Vapi (Voice AI) | `vapi` | SO 03 | `op://SO/Vapi/api-key` |
| Circle.so | `circle` | GE 78 | `op://SO/Circle/api-key` |
| LikeKit | `likekit` | GE 78 | `op://SO/LikeKit/api-key` |
| Legacy Library | `legacy-library` | HP 28 | — (local processing) |
| Video Production | `video-production` | GE 78 | — (local tools) |
| Capture Pipeline | `capture` | SO 24 | — (uses other skills) |

---

## Configured / Planned Tools

Tools identified for integration. Each becomes a skill when implemented.

### Browser Automation
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| Puppeteer | 🔧 | Headless browser automation, scraping | Yes |
| Selenium | 📋 | Browser testing, web automation | Maybe |
| Bytebot | 📋 | AI-powered browser agent | Maybe |

### Communication
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| Twilio SMS | 📋 | Outbound SMS notifications, 2FA | Yes |
| Twilio Voice | 📋 | Voice calls, IVR, phone integration | Yes |
| Resend | 📋 | Transactional email API | Yes |
| Zoho Mail | 🔧 | Email (covered by zoho skill) | — |

### Search & Research
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| Exa | 📋 | AI-native web search API | Yes |
| DuckDuckGo | 📋 | Privacy-first web search | Yes |
| Google Maps | 📋 | Location, directions, places | Yes |
| Arxiv | 📋 | Academic paper search | Yes |
| Wikipedia | 📋 | Knowledge lookup | Yes |

### Data & Finance
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| YFinance | 📋 | Stock/crypto market data | Yes |

### Transportation & IoT
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| Bouncie | 📋 | Vehicle GPS tracking, trip logs | Yes |
| Citibike | 📋 | Bike availability near location | Maybe |
| Lyft | 📋 | Ride pricing/booking | Maybe |
| Uber | 📋 | Ride pricing/booking | Maybe |

### Productivity
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| ClickUp | 📋 | Project management (deep work tracking) | Maybe |

### Agent Patterns (Internal)
| Pattern | Status | Purpose |
|---------|--------|---------|
| Run Team | 🔧 | Multi-agent parallel execution |
| Ask Agent | 🔧 | Agent-to-agent delegation |
| Sequential Thinking | 🔧 | Chain-of-thought reasoning pattern |
| Shell | ✅ | Command execution (built-in) |
| Calculator | ✅ | Math operations (built-in) |

### Mail & Print
| Tool | Status | Purpose | Future Skill? |
|------|--------|---------|:---:|
| Lob | 📋 | Physical mail/postcard API | Maybe |

---

## Domain Routing Guide

When a new tool is added, route it to the correct domain:

| Tool Category | Primary Domain | Why |
|--------------|---------------|-----|
| Agent tools, APIs, integrations | **SO 03 Agents** | Core agent capability |
| Automation, workflow | **SO 09 Automation** | Orchestration layer |
| Security, credentials | **SO 11 Security** | Access control |
| Research tools | **SO 12 Research** | Knowledge gathering |
| Finance APIs | **HP 30 Finance** | Personal finance data |
| Health/fitness APIs | **HP 27 Body** | Body metrics |
| Transportation | **HP 25 Mind** or **LE 51** | Daily life logistics |
| Home/IoT | **LE 51 Home** | Smart home |
| Marketing/email | **GE 78 Marketing** | Growth & outreach |
| CRM/sales | **GE 75 Brand** | Business relationships |

---

## Adding a New Tool

1. Add entry to this registry with status 📋
2. Create Todoist task under `06 Tools & APIs` → "Set Up Tools & APIs"
3. When implementing: create `.claude/skills/{tool-name}/SKILL.md`
4. Update status to ✅ Active in this registry
5. Add credentials to 1Password under appropriate SO vault

---

*Part of the One Mind framework. Domain: SO 03 Agents.*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
