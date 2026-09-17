# Access Methods Reference

> Comprehensive guide to accessing apps and services via AI, CLI, API, SDK, and MCP

This document catalogs all access methods available for each service/app in the OneMind OS ecosystem. Use this to determine the best integration approach for AI agents and automation.

---

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Available and confirmed |
| ⚠️ | Limited/Beta/Partial support |
| ❌ | Not available |
| 🔗 | Has official MCP server |

---

## Quick Reference Table

| Service | API | CLI | SDK | MCP | Best AI Access |
|---------|-----|-----|-----|-----|----------------|
| **Letta** | ✅ | ✅ | ✅ | ✅ | MCP/SDK |
| **Twilio** | ✅ | ✅ | ✅ | ✅ | MCP/API |
| **Vapi** | ✅ | ✅ | ✅ | ✅ | MCP/CLI |
| **LiveKit** | ✅ | ✅ | ✅ | ✅ | MCP/SDK |
| **Cloudflare** | ✅ | ✅ | ✅ | ✅ | CLI (wrangler)/MCP |
| **GitHub** | ✅ | ✅ | ✅ | ✅ | MCP/CLI (gh) |
| **GitHub Copilot** | ✅ | ✅ | ⚠️ | ✅ | CLI/MCP |
| **1Password** | ✅ | ✅ | ✅ | ✅ | CLI/MCP |
| **AWS** | ✅ | ✅ | ✅ | ✅ | CLI/MCP |
| **Azure** | ✅ | ✅ | ✅ | ✅ | CLI/MCP |
| **Circle.so** | ✅ | ❌ | ✅ | ⚠️ | API/SDK |
| **Whop** | ✅ | ❌ | ✅ | ✅ | API/SDK/MCP |
| **Stripe** | ✅ | ✅ | ✅ | ✅ | MCP/API |
| **n8n** | ✅ | ❌ | ⚠️ | ✅ | MCP |
| **Webflow** | ✅ | ❌ | ✅ | ✅ | MCP |
| **Raindrop.io** | ✅ | ❌ | ⚠️ | ✅ | MCP/API |
| **Bouncie** | ✅ | ❌ | ❌ | ⚠️ | API |
| **NATS** | ✅ | ✅ | ✅ | ❌ | CLI/SDK |
| **OpenAI** | ✅ | ⚠️ | ✅ | ✅ | SDK/MCP |
| **Supabase** | ✅ | ✅ | ✅ | ✅ | MCP/CLI |
| **Vercel** | ✅ | ✅ | ✅ | ✅ | CLI/MCP |
| **Linear** | ✅ | ⚠️ | ✅ | ✅ | MCP/API |
| **Notion** | ✅ | ❌ | ✅ | ✅ | MCP |
| **Slack** | ✅ | ❌ | ✅ | ✅ | MCP/SDK |
| **Discord** | ✅ | ❌ | ✅ | ✅ | MCP/SDK |
| **Firecrawl** | ✅ | ❌ | ✅ | ✅ | MCP |
| **Resend** | ✅ | ❌ | ✅ | ✅ | MCP/SDK |

---

## Detailed Service Breakdown

### AI & Agent Platforms

#### Letta 🔗
Memory-first AI agent framework with persistent context.

| Method | Details |
|--------|---------|
| **API** | REST API for agent management, memory, MCP routes under `/mcp-servers` |
| **CLI** | `letta` CLI for local development and management |
| **SDK** | Python SDK (`letta`), TypeScript SDK, Vercel AI SDK provider |
| **MCP** | Native MCP support - can connect to external MCP servers |

```bash
# Install
pip install letta

# MCP Config (Claude Desktop)
{
  "mcpServers": {
    "letta": {
      "command": "letta",
      "args": ["server", "--mcp"]
    }
  }
}
```

**Best for AI**: SDK for programmatic access, MCP for AI assistants

---

#### Vapi 🔗
Voice AI platform for building conversational agents.

| Method | Details |
|--------|---------|
| **API** | Full REST API for assistants, calls, phone numbers, workflows |
| **CLI** | `vapi` CLI - init projects, manage assistants, local webhook debugging |
| **SDK** | `@vapi-ai/web` (Web), Python SDK, server SDKs |
| **MCP** | Official MCP server with docs integration |

```bash
# Install CLI
curl -sSL https://vapi.ai/install.sh | bash

# Add MCP to IDE
vapi mcp setup

# MCP Config
{
  "mcpServers": {
    "vapi": {
      "command": "npx",
      "args": ["-y", "@vapi-ai/mcp-server"]
    }
  }
}
```

**Best for AI**: CLI for development, MCP for IDE integration

---

#### LiveKit 🔗
Real-time audio/video infrastructure.

| Method | Details |
|--------|---------|
| **API** | REST API + WebSocket for rooms, participants, tracks |
| **CLI** | `lk` CLI for room management, token generation, testing |
| **SDK** | SDKs: Python, JS/TS, Go, Swift, Kotlin, Rust, Unity, Flutter |
| **MCP** | Official docs MCP server at `https://docs.livekit.io/mcp` |

```bash
# Install CLI
brew install livekit-cli  # or download from releases

# Add MCP (Claude Code)
claude mcp add --transport http livekit-docs https://docs.livekit.io/mcp

# Python with MCP
uv add livekit-agents[mcp]~=1.3
```

**Best for AI**: MCP for documentation, SDK for agent development

---

#### OpenAI 🔗
LLM and AI services.

| Method | Details |
|--------|---------|
| **API** | Responses API, Assistants API, Chat Completions |
| **CLI** | Limited - via Codex CLI or third-party tools |
| **SDK** | Official: Python (`openai`), Node.js, .NET. Agents SDK for MCP |
| **MCP** | Remote MCP server support in Responses API |

```python
# Using MCP with OpenAI Agents SDK
from agents.mcp import MCPServerStdio

mcp_server = MCPServerStdio(
    params={"command": "npx", "args": ["-y", "@modelcontextprotocol/server-filesystem", "."]}
)

agent = Agent(
    name="Assistant",
    instructions="Use tools to access files.",
    mcp_servers=[mcp_server]
)
```

**Best for AI**: Agents SDK with MCP integration

---

### Communication Services

#### Twilio 🔗
Programmable SMS, Voice, Video, and more.

| Method | Details |
|--------|---------|
| **API** | Full REST API for messaging, voice, video, serverless |
| **CLI** | `twilio` CLI - send messages, make calls, manage resources |
| **SDK** | Python, Node.js, C#, Java, PHP, Ruby, Go |
| **MCP** | Official Alpha MCP server from twilio-labs |

```bash
# Install CLI
npm install -g twilio-cli

# MCP Config
{
  "mcpServers": {
    "twilio": {
      "command": "npx",
      "args": ["-y", "@anthropic/mcp-server-twilio"],
      "env": {
        "TWILIO_ACCOUNT_SID": "your_sid",
        "TWILIO_AUTH_TOKEN": "your_token"
      }
    }
  }
}
```

**Best for AI**: MCP for conversational access, API for programmatic

---

#### Slack 🔗
Team messaging and collaboration.

| Method | Details |
|--------|---------|
| **API** | Web API, Events API, RTM API (legacy) |
| **CLI** | No official CLI |
| **SDK** | Bolt SDK (Python, JS), `slack_sdk` Python |
| **MCP** | Official Slack MCP server (select partners) |

```python
# Slack SDK
from slack_sdk import WebClient
client = WebClient(token="xoxb-your-token")
client.chat_postMessage(channel="#general", text="Hello!")
```

**Best for AI**: MCP (if available), SDK for bot development

---

#### Discord 🔗
Gaming/community platform.

| Method | Details |
|--------|---------|
| **API** | REST API + Gateway (WebSocket) |
| **CLI** | No official CLI |
| **SDK** | `discord.py`, `discord.js`, `JDA` (Java), `discordgo` |
| **MCP** | Community MCP servers (hanweg/mcp-discord, SaseQ/discord-mcp) |

```bash
# MCP Config
{
  "mcpServers": {
    "discord": {
      "command": "npx",
      "args": ["-y", "@hanweg/mcp-discord"],
      "env": {
        "DISCORD_TOKEN": "your_bot_token"
      }
    }
  }
}
```

**Best for AI**: MCP for chat integration, SDK for full bot development

---

#### Resend 🔗
Email API for developers.

| Method | Details |
|--------|---------|
| **API** | REST API for sending emails, domains, audiences, webhooks |
| **CLI** | No official CLI |
| **SDK** | Node.js, Python, Ruby, PHP, Go, Rust, Elixir, Java, .NET |
| **MCP** | Official Resend MCP server for sending emails |

```bash
# Install SDK (Node.js)
npm install resend

# Install SDK (Python)
pip install resend

# MCP Config
{
  "mcpServers": {
    "resend": {
      "command": "node",
      "args": ["path/to/mcp-send-email/build/index.js"],
      "env": {
        "RESEND_API_KEY": "re_xxxxx"
      }
    }
  }
}
```

```typescript
// SDK Usage
import { Resend } from 'resend';

const resend = new Resend('re_xxxxx');

await resend.emails.send({
  from: 'you@example.com',
  to: 'user@example.com',
  subject: 'Hello',
  html: '<p>Hello World</p>'
});
```

**Best for AI**: MCP for composing/sending emails, SDK for programmatic

---

### Cloud Providers

#### AWS 🔗
Amazon Web Services.

| Method | Details |
|--------|---------|
| **API** | All services have REST/HTTP APIs |
| **CLI** | `aws` CLI - full service coverage |
| **SDK** | Boto3 (Python), AWS SDK for JS, Java, .NET, Go, etc. |
| **MCP** | AWS API MCP Server on AWS Marketplace, Cloud Control API MCP |

```bash
# Install CLI
pip install awscli

# MCP Config
{
  "mcpServers": {
    "aws": {
      "command": "npx",
      "args": ["-y", "@aws/mcp-server"],
      "env": {
        "AWS_PROFILE": "default"
      }
    }
  }
}
```

**Best for AI**: CLI for automation, MCP for natural language operations

---

#### Azure 🔗
Microsoft Azure cloud.

| Method | Details |
|--------|---------|
| **API** | Azure Resource Manager REST APIs |
| **CLI** | `az` CLI - comprehensive Azure management |
| **SDK** | Azure SDK for Python, JS, .NET, Java, Go |
| **MCP** | Official Azure MCP Server from Microsoft |

```bash
# Install CLI
curl -sL https://aka.ms/InstallAzureCLIDeb | sudo bash

# MCP Server (Java-based)
git clone https://github.com/jdubois/azure-cli-mcp.git
cd azure-cli-mcp
./mvnw clean package -DskipTests
java -jar target/azure-cli-mcp-*.jar
```

**Best for AI**: CLI through MCP wrapper, native MCP server

---

#### Cloudflare 🔗
Edge computing, CDN, DNS, security.

| Method | Details |
|--------|---------|
| **API** | Full REST API for all services |
| **CLI** | `wrangler` CLI - Workers, Pages, R2, D1, etc. |
| **SDK** | `cloudflare` Python SDK, JS SDK |
| **MCP** | Build MCP servers ON Cloudflare Workers, Agents SDK |

```bash
# Install Wrangler
npm install -g wrangler

# Deploy MCP server to Cloudflare
npm create cloudflare@latest -- my-mcp-server \
  --template=cloudflare/ai/demos/remote-mcp-github-oauth
```

**Best for AI**: Wrangler CLI, deploy MCP servers to Workers

---

#### Vercel 🔗
Frontend cloud platform.

| Method | Details |
|--------|---------|
| **API** | REST API for deployments, projects, domains |
| **CLI** | `vercel` CLI - deploy, dev, env management |
| **SDK** | AI SDK with MCP adapter |
| **MCP** | MCP adapter for building servers, official Vercel MCP server |

```bash
# Install CLI
npm i -g vercel

# Using Vercel MCP Adapter
import { McpServer } from "@vercel/mcp-adapter";
```

**Best for AI**: CLI for deployments, AI SDK for MCP integration

---

#### Supabase 🔗
Open source Firebase alternative.

| Method | Details |
|--------|---------|
| **API** | REST API (PostgREST), Realtime API |
| **CLI** | `supabase` CLI - local dev, migrations, functions |
| **SDK** | `@supabase/supabase-js`, Python client |
| **MCP** | Official Supabase MCP server with 20+ tools |

```bash
# Install CLI
npm install -g supabase

# MCP Config
{
  "mcpServers": {
    "supabase": {
      "command": "npx",
      "args": ["-y", "@supabase/mcp-server"],
      "env": {
        "SUPABASE_ACCESS_TOKEN": "your_token"
      }
    }
  }
}
```

**Best for AI**: MCP for database operations, CLI for migrations

---

### Developer Tools

#### GitHub 🔗
Code hosting and collaboration.

| Method | Details |
|--------|---------|
| **API** | REST API v3, GraphQL API v4 |
| **CLI** | `gh` CLI - issues, PRs, repos, actions, releases |
| **SDK** | Octokit (JS), PyGithub (Python) |
| **MCP** | Official GitHub MCP server, Copilot integration |

```bash
# Install CLI
brew install gh

# MCP Config (automatic in VS Code with Copilot)
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@github/mcp-server"]
    }
  }
}
```

**Best for AI**: `gh` CLI, GitHub MCP server through Copilot

---

#### GitHub Copilot 🔗
AI pair programmer.

| Method | Details |
|--------|---------|
| **API** | Extensions API (VS Code), Chat API |
| **CLI** | `copilot` CLI - interactive terminal AI |
| **SDK** | Extension SDK for VS Code |
| **MCP** | Full MCP support, GitHub MCP server integration |

```bash
# Install Copilot CLI
gh extension install github/gh-copilot

# Use in terminal
gh copilot explain "git rebase -i HEAD~5"
gh copilot suggest "deploy to production"
```

**Best for AI**: CLI for terminal, IDE extensions with MCP

---

#### Linear 🔗
Issue tracking for modern teams.

| Method | Details |
|--------|---------|
| **API** | GraphQL API |
| **CLI** | Limited - community tools |
| **SDK** | `@linear/sdk` (JS/TS) |
| **MCP** | Community MCP servers (linear-mcp, mcp-linear) |

```bash
# MCP Config
{
  "mcpServers": {
    "linear": {
      "command": "npx",
      "args": ["-y", "@linear/mcp-server"],
      "env": {
        "LINEAR_API_KEY": "your_key"
      }
    }
  }
}
```

**Best for AI**: MCP for issue management, SDK for integrations

---

### Productivity & Knowledge

#### Notion 🔗
All-in-one workspace.

| Method | Details |
|--------|---------|
| **API** | REST API for databases, pages, blocks, users |
| **CLI** | No official CLI |
| **SDK** | `@notionhq/client` (JS), community Python SDK |
| **MCP** | Official Notion MCP with OAuth, community servers |

```bash
# MCP Config (Official - uses OAuth)
{
  "mcpServers": {
    "notion": {
      "command": "npx",
      "args": ["-y", "@notionhq/mcp-server"]
    }
  }
}

# Community MCP (token-based)
{
  "mcpServers": {
    "notion": {
      "command": "npx",
      "args": ["-y", "mcp-notion-server"],
      "env": {
        "NOTION_API_TOKEN": "secret_xxx"
      }
    }
  }
}
```

**Best for AI**: MCP (official or community)

---

#### Raindrop.io 🔗
Bookmark manager.

| Method | Details |
|--------|---------|
| **API** | REST API for bookmarks, collections, tags |
| **CLI** | No official CLI |
| **SDK** | No official SDK |
| **MCP** | Community MCP server (hiromitsusasaki/raindrop-io-mcp-server) |

```bash
# MCP Config
{
  "mcpServers": {
    "raindrop": {
      "command": "npx",
      "args": ["-y", "raindrop-io-mcp-server"],
      "env": {
        "RAINDROP_TOKEN": "your_token"
      }
    }
  }
}
```

**Best for AI**: MCP for bookmark management

---

### Payments & Commerce

#### Stripe 🔗
Payment processing.

| Method | Details |
|--------|---------|
| **API** | Full REST API for payments, subscriptions, billing |
| **CLI** | `stripe` CLI - webhooks, testing, logs |
| **SDK** | `stripe` Python, Node.js, Ruby, PHP, Java, Go, .NET |
| **MCP** | Official Stripe MCP server with agent toolkit |

```bash
# Install CLI
brew install stripe/stripe-cli/stripe

# MCP Config
{
  "mcpServers": {
    "stripe": {
      "command": "npx",
      "args": ["-y", "@stripe/mcp-server"],
      "env": {
        "STRIPE_API_KEY": "sk_..."
      }
    }
  }
}
```

**Best for AI**: MCP for payment operations, CLI for testing

---

#### Whop 🔗
Digital product marketplace.

| Method | Details |
|--------|---------|
| **API** | REST API for products, payments, memberships |
| **CLI** | No official CLI |
| **SDK** | Python SDK, TypeScript SDK |
| **MCP** | Official Whop MCP server |

```python
# Python SDK
from whop import Whop
client = Whop(api_key="your_key")
```

**Best for AI**: SDK/API, MCP for AI assistants

---

### Security

#### 1Password 🔗
Password manager.

| Method | Details |
|--------|---------|
| **API** | Connect Server API for secrets |
| **CLI** | `op` CLI - vaults, items, secrets management |
| **SDK** | Python SDK, Go SDK |
| **MCP** | Official MCP Server for SaaS Manager, community POC servers |

```bash
# Install CLI
brew install 1password-cli

# Get secret in scripts
op read "op://vault/item/field"

# MCP Config (community)
{
  "mcpServers": {
    "1password": {
      "command": "uv",
      "args": ["run", "--with", "mcp[cli]", "--with", "onepassword-sdk", "mcp", "run", "server.py"],
      "env": {
        "OP_SERVICE_ACCOUNT_TOKEN": "your_token"
      }
    }
  }
}
```

**Best for AI**: CLI for scripts, MCP for AI credential access

---

### Automation & Workflows

#### n8n 🔗
Workflow automation.

| Method | Details |
|--------|---------|
| **API** | REST API for workflows, executions, credentials |
| **CLI** | No official CLI |
| **SDK** | Limited - node-based workflows |
| **MCP** | MCP Server Trigger + MCP Client Tool nodes |

```javascript
// n8n exposes workflows as MCP tools
// Use MCP Server Trigger node to expose workflows
// Use MCP Client Tool to connect to external MCP servers
```

**Best for AI**: MCP integration (both directions)

---

### Web Scraping & Data

#### Firecrawl 🔗
Web scraping API for LLMs - converts websites to clean markdown.

| Method | Details |
|--------|---------|
| **API** | REST API for scrape, crawl, map, search, extract, deep research |
| **CLI** | No official CLI |
| **SDK** | Python (`firecrawl-py`), Node.js (`@mendable/firecrawl-js`) |
| **MCP** | Official Firecrawl MCP server with 8+ tools |

```bash
# Install SDK (Python)
pip install firecrawl-py

# Install SDK (Node.js)
npm install @mendable/firecrawl-js

# MCP Config (Local)
{
  "mcpServers": {
    "firecrawl": {
      "command": "npx",
      "args": ["-y", "firecrawl-mcp"],
      "env": {
        "FIRECRAWL_API_KEY": "fc-xxxxx",
        "FIRECRAWL_RETRY_MAX_ATTEMPTS": "5",
        "FIRECRAWL_RETRY_INITIAL_DELAY": "2000"
      }
    }
  }
}

# MCP Config (Remote/Hosted)
{
  "mcpServers": {
    "firecrawl": {
      "url": "https://mcp.firecrawl.dev/{FIRECRAWL_API_KEY}/v2/sse"
    }
  }
}
```

```python
# SDK Usage
from firecrawl import FirecrawlApp

app = FirecrawlApp(api_key="fc-xxxxx")

# Scrape single page
result = app.scrape_url("https://example.com", params={"formats": ["markdown"]})

# Crawl entire site
crawl = app.crawl_url("https://example.com", params={"limit": 100})

# Search the web
search = app.search("AI agents frameworks 2024")
```

**MCP Tools Available:**
- `firecrawl_scrape` - Scrape single page to markdown
- `firecrawl_batch_scrape` - Batch scrape multiple URLs
- `firecrawl_crawl` - Crawl entire website
- `firecrawl_map` - Map website structure
- `firecrawl_search` - Web search with optional extraction
- `firecrawl_extract` - Extract structured data with LLM
- `firecrawl_deep_research` - Deep research on a topic
- `firecrawl_generate_llmstxt` - Generate llms.txt for a domain

**Best for AI**: MCP for web scraping/research tasks

---

### Web & CMS

#### Webflow 🔗
Visual web development.

| Method | Details |
|--------|---------|
| **API** | Data API, Designer API |
| **CLI** | No official CLI |
| **SDK** | `webflow` npm package |
| **MCP** | Official Webflow MCP server with OAuth |

```bash
# MCP Config (Remote - uses OAuth)
{
  "mcpServers": {
    "webflow": {
      "command": "npx",
      "args": ["mcp-remote", "https://mcp.webflow.com/sse"]
    }
  }
}

# Local MCP (token-based)
{
  "mcpServers": {
    "webflow": {
      "command": "npx",
      "args": ["-y", "webflow-mcp-server@latest"],
      "env": {
        "WEBFLOW_TOKEN": "your_token"
      }
    }
  }
}
```

**Best for AI**: MCP for site management and Designer API

---

### Community & Membership

#### Circle.so
Community platform.

| Method | Details |
|--------|---------|
| **API** | Admin API v2, Headless Member API, Data API |
| **CLI** | No official CLI |
| **SDK** | Zapier, Make, Integrately integrations |
| **MCP** | No official MCP |

```javascript
// Circle API example
fetch('https://api.circle.so/v1/spaces', {
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY'
  }
})
```

**Best for AI**: API with automation platforms (Zapier/Make)

---

### Vehicle & IoT

#### Bouncie
Connected vehicle tracking.

| Method | Details |
|--------|---------|
| **API** | REST API with OAuth2, Webhook support |
| **CLI** | No official CLI |
| **SDK** | No official SDK |
| **MCP** | No official MCP (Home Assistant integration available) |

```javascript
// Bouncie API (via OAuth2)
// Setup at https://docs.bouncie.dev/
// Supports: trips, vehicle info, geofencing, diagnostics
```

**Best for AI**: API via custom integration

---

### Messaging Infrastructure

#### NATS
High-performance messaging.

| Method | Details |
|--------|---------|
| **API** | Client protocol (TCP), JetStream API |
| **CLI** | `nats` CLI - pub/sub, streams, consumers, benchmarking |
| **SDK** | Go, Python, Node.js, Rust, C, .NET, Java, Ruby |
| **MCP** | No official MCP |

```bash
# Install CLI
brew install nats-io/nats-tools/nats

# Basic commands
nats pub hello "world"
nats sub hello
nats stream ls
nats bench test --msgs 1000000 --pub 2 --sub 2
```

**Best for AI**: CLI for operations, SDK for integration

---

## MCP Server Registry

### Official MCP Servers

| Provider | URL/Package | Transport |
|----------|-------------|-----------|
| GitHub | `@github/mcp-server` | stdio |
| Notion | `@notionhq/mcp-server` | OAuth/Remote |
| Stripe | `@stripe/mcp-server` | stdio |
| Supabase | `@supabase/mcp-server` | stdio/Remote |
| Webflow | `https://mcp.webflow.com/sse` | SSE |
| Slack | Select partners only | Remote |
| LiveKit | `https://docs.livekit.io/mcp` | HTTP |
| Azure | `azure-cli-mcp` | stdio |
| AWS | AWS Marketplace | Various |
| Firecrawl | `firecrawl-mcp` / `https://mcp.firecrawl.dev` | stdio/SSE |
| Resend | `mcp-send-email` | stdio |

### Community MCP Servers

| Service | Repository |
|---------|------------|
| Discord | `hanweg/mcp-discord`, `SaseQ/discord-mcp` |
| Linear | `cosmix/linear-mcp`, `scoutos/mcp-linear` |
| Raindrop | `hiromitsusasaki/raindrop-io-mcp-server` |
| 1Password | `onepassword-mcp-server` (POC) |
| n8n | Built-in MCP nodes |

---

## Recommended Access Patterns

### For AI Agents (Claude, GPT, etc.)
1. **MCP First** - Use MCP servers when available for native integration
2. **CLI Second** - CLI tools work well with shell-capable agents
3. **API Third** - Direct API calls when no MCP/CLI available

### For Automation (n8n, Zapier)
1. **API/Webhooks** - Most reliable for workflow automation
2. **MCP** - n8n supports MCP natively
3. **CLI** - For script-based automation

### For Development
1. **CLI** - Best for local development and scripting
2. **SDK** - For production application integration
3. **MCP** - For AI-assisted development in IDEs

---

## Adding New Services

When integrating a new service, check in this order:

1. **Official MCP Server** - Check provider docs or https://mcpservers.org
2. **CLI** - Check if provider has official CLI tool
3. **SDK** - Look for official SDKs in your language
4. **API** - REST/GraphQL API documentation
5. **Community** - GitHub search for MCP servers or wrappers

---

## Maintenance

This document should be updated when:
- New MCP servers are released
- CLI tools are added/updated
- API versions change
- New services are integrated into OneMind OS

Last updated: 2026-02-03


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
