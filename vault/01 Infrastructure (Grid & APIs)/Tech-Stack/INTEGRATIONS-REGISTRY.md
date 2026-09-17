# Integrations Registry

> **Active APIs, webhooks, and external service integrations**

## Authentication & OAuth

| Service       | Type            | Integration Status | Auth Method |
| ------------- | --------------- | ------------------ | ----------- |
| **Composio**  | OAuth Platform  | ✅ Active          | API Key     |
| **1Password** | Secrets Manager | ✅ Active          | CLI (`op`)  |
| **Tailscale** | VPN             | ✅ Active          | OAuth       |

## AI & Language Models

| Service              | Purpose           | Integration Status | Cost        |
| -------------------- | ----------------- | ------------------ | ----------- |
| **AWS Bedrock**      | Claude Opus 4.5   | ✅ Active          | Pay-per-use |
| **OpenAI API**       | GPT-4, Embeddings | ✅ Active          | Pay-per-use |
| **Anthropic Direct** | Claude API        | ⚠️ Backup          | Pay-per-use |

## Productivity & Project Management

| Service             | Purpose         | Integration Status | Auth Method      |
| ------------------- | --------------- | ------------------ | ---------------- |
| **Todoist**         | Task management | ✅ Active          | OAuth + Webhooks |
| **Linear**          | Issue tracking  | ✅ Active          | API Key          |
| **Google Calendar** | Calendar sync   | ✅ Active          | OAuth            |
| **Obsidian**        | Knowledge base  | ✅ Active          | Local files      |

## Communication & Real-Time

| Service            | Purpose            | Integration Status | Cost        |
| ------------------ | ------------------ | ------------------ | ----------- |
| **LiveKit**        | Voice/video        | ✅ Active          | Self-hosted |
| **NATS JetStream** | Message broker     | ✅ Active          | Self-hosted |
| **MQTT**           | IoT messaging      | ✅ Active          | Self-hosted |
| **Firebase**       | Push notifications | 🔄 Planned         | Free tier   |

## Transportation & Location

| Service         | Purpose          | Integration Status | Cost                     |
| --------------- | ---------------- | ------------------ | ------------------------ |
| **Uber API**    | Ride ordering    | 🔄 Planned         | Free (rides cost normal) |
| **Lyft API**    | Ride ordering    | 🔄 Planned         | Free (rides cost normal) |
| **Google Maps** | Maps, navigation | 🔄 Planned         | $200 credit/month        |
| **MTA API**     | NYC transit data | 🔄 Planned         | Free                     |

## Research & Search

| Service            | Purpose         | Integration Status | Cost            |
| ------------------ | --------------- | ------------------ | --------------- |
| **Arxiv API**      | Academic papers | 🔄 Planned         | Free            |
| **Brave Search**   | Web search      | 🔄 Planned         | 2000 free/month |
| **Exa (Metaphor)** | Semantic search | 🔄 Planned         | ~$10/month      |
| **Tavily**         | AI search       | 🔄 Planned         | ~$10/month      |
| **Wolfram Alpha**  | Computation     | 🔄 Planned         | $5/month        |
| **Context7**       | Code docs       | ✅ Active (MCP)    | Free            |

## Web Scraping & Crawling

| Service         | Purpose         | Integration Status | Cost               |
| --------------- | --------------- | ------------------ | ------------------ |
| **Jina Reader** | URL to markdown | 🔄 Planned         | Free tier          |
| **Firecrawl**   | Web scraping    | 🔄 Planned         | ~$20/month         |
| **Crawl4AI**    | Web crawler     | 🔄 Planned         | Self-hosted (free) |
| **SearXNG**     | Meta-search     | 🔄 Planned         | Self-hosted (free) |

## Code & Development

| Service        | Purpose            | Integration Status | Auth Method           |
| -------------- | ------------------ | ------------------ | --------------------- |
| **GitHub API** | Repos, issues, PRs | ✅ Active (MCP)    | Personal Access Token |
| **Vercel**     | Web hosting        | ✅ Active          | OAuth                 |
| **Docker Hub** | Container registry | ✅ Active          | API Token             |

## Smart Home & IoT

| Service            | Purpose        | Integration Status | Cost          |
| ------------------ | -------------- | ------------------ | ------------- |
| **Home Assistant** | Smart home     | ✅ Active          | Self-hosted   |
| **Meshtastic**     | LoRa mesh      | 🔄 Planned         | Hardware only |
| **ESP32**          | Custom devices | 🔄 Planned         | Hardware only |

## Business & CRM

| Service       | Purpose             | Integration Status | Auth Method |
| ------------- | ------------------- | ------------------ | ----------- |
| **Zoho CRM**  | Customer management | ✅ Active          | OAuth       |
| **Circle.so** | Community platform  | ✅ Active          | API Key     |

## Cloud Infrastructure

| Service        | Purpose       | Integration Status | Cost        |
| -------------- | ------------- | ------------------ | ----------- |
| **AWS**        | Cloud hosting | ✅ Active          | Pay-per-use |
| **Render.com** | Alt hosting   | ⚠️ Backup          | Free tier   |
| **PostgreSQL** | Database      | ✅ Active          | Self-hosted |

## Integration Status Legend

- ✅ **Active** - Currently integrated and running
- 🔄 **Planned** - Linear issue created, in backlog
- ⚠️ **Backup** - Configured but not primary
- ❌ **Deprecated** - No longer used

## Webhook Endpoints

All webhooks route through NATS JetStream:

```
https://onemind.zeusdelacruz.com/webhooks/{integration}
↓
NATS subject: events.{integration}.{event_type}
↓
Legacy AI consciousness stream
```

**Active Webhooks:**

- Todoist: Task events → `events.todoist.*`
- GitHub: PR/issue events → `events.github.*`
- Zoho: CRM updates → `events.zoho.*`

## API Rate Limits

| Service         | Rate Limit        | Current Usage    |
| --------------- | ----------------- | ---------------- |
| **Todoist**     | 450 req/15 min    | ~50/day          |
| **Linear**      | 1000 req/hour     | ~100/day         |
| **Google Maps** | $200 credit/month | $0 (not started) |
| **Uber**        | 10 req/min        | $0 (not started) |

## Security

- **OAuth tokens** stored in 1Password
- **API keys** stored in 1Password
- **Webhooks** verified with HMAC signatures
- **HTTPS only** for all external communications

---

Last Updated: January 25, 2026
