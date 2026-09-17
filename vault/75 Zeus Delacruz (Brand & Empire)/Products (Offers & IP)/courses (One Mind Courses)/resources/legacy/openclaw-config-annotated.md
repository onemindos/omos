# openclaw.json — Annotated Reference

> *Complete annotated `openclaw.json` configuration. Every field explained.*
> *Your config lives at: `~/.openclaw/openclaw.json`*

---

```json
{
  "version": "1",

  "gateway": {
    // The gateway is the OpenClaw server — agents connect through it
    "port": 18789,
    // trustedProxies: Add your reverse proxy IP if using Tailscale serve or nginx
    "trustedProxies": ["127.0.0.1"],
    // maxConnections: how many simultaneous client connections
    "maxConnections": 10
  },

  "models": {
    // Define named model references used by agents
    // Format: "alias": { "provider": "...", "model": "..." }
    "default": {
      "provider": "github-copilot",
      "model": "claude-sonnet-4.6"
    },
    "fast": {
      "provider": "github-copilot",
      "model": "claude-haiku"
      // Use "fast" for lightweight tasks: summaries, routing, brief generation
    }
  },

  "agents": [
    {
      "id": "main",
      // id is used in CLI: openclaw agent chat main
      "name": "Legacy ⚡",
      "model": "default",
      // workspace: path to the directory containing IDENTITY.md, CODEX.md, USER.md
      "workspace": "/root/.openclaw/workspace",
      "enabled": true,
      // memory: enable persistent memory compaction
      "memory": {
        "enabled": true,
        "compactionThreshold": 80000,
        // maxTokens before compaction triggers
        "summaryModel": "fast"
      }
    },
    {
      "id": "spartan",
      "name": "Spartan 💪",
      "model": "default",
      "workspace": "/root/.openclaw/workspace-spartan",
      "enabled": true,
      "memory": { "enabled": true, "compactionThreshold": 80000, "summaryModel": "fast" }
    }
    // Repeat pattern for: oracle, eden, haven, empire, guardian, grid, trinity
  ],

  "channels": {
    "discord": {
      "enabled": true,
      "accounts": {
        "sentinel": {
          // token: Discord bot token — store in 1Password, inject at runtime
          "token": "op://aaeyris5z44pdqspk277u2oppi/discord-bot-token/credential",
          "defaultAgent": "main",
          // routingRules: route messages by channel to the right agent
          "routingRules": [
            { "channel": "spartan", "agent": "spartan" },
            { "channel": "oracle", "agent": "oracle" },
            { "channel": "empire", "agent": "empire" },
            { "channel": "guardian", "agent": "guardian" },
            { "channel": "grid", "agent": "grid" },
            { "channel": "trinity", "agent": "trinity" }
          ]
        }
      }
    },
    "telegram": {
      "enabled": false,
      // Enable telegram for mobile access — set token to activate
      "accounts": {
        "default": {
          "token": "",
          "defaultAgent": "main",
          // dmPolicy: "allowlist" restricts to specific user IDs (recommended)
          "dmPolicy": "allowlist",
          "allowedUserIds": []
        }
      }
    }
  },

  "cron": [
    {
      "id": "morning-briefing",
      "schedule": "0 6 * * *",
      // Cron syntax: min hour day month weekday
      // "0 6 * * *" = 6:00 AM every day
      "agent": "main",
      "channel": "discord",
      "account": "sentinel",
      "prompt": "Generate my morning brief. Pull active goals from USER.md, flag today's priorities, check system status.",
      "enabled": true
    },
    {
      "id": "eod-review",
      "schedule": "0 20 * * *",
      "agent": "main",
      "channel": "discord",
      "account": "sentinel",
      "prompt": "Generate EOD review. Summarize today, flag carries, set tomorrow's primary.",
      "enabled": true
    },
    {
      "id": "weekly-report",
      "schedule": "0 17 * * 0",
      // "0 17 * * 0" = 5:00 PM every Sunday (day 0)
      "agent": "main",
      "channel": "discord",
      "account": "sentinel",
      "prompt": "Generate full weekly report for Zeus's Sunday review.",
      "enabled": true
    }
  ],

  "mcp": {
    // MCP servers extend agent capabilities with external tools
    "servers": {
      "codex-vault": {
        // Gives agents read/write access to the vault filesystem
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-filesystem", "/root/onemind-codex"],
        "enabled": true
      },
      "github": {
        // Gives agents GitHub access — read repos, issues, PRs
        "command": "npx",
        "args": ["-y", "@modelcontextprotocol/server-github"],
        "env": {
          "GITHUB_PERSONAL_ACCESS_TOKEN": "op://aaeyris5z44pdqspk277u2oppi/github-pat/credential"
        },
        "enabled": true
      },
      "context7": {
        // Gives agents access to up-to-date library documentation
        "command": "npx",
        "args": ["-y", "@upstash/context7-mcp@latest"],
        "enabled": true
      }
    }
  },

  "memory": {
    // Global memory search path — vault is semantically searchable by all agents
    "extraPaths": ["/root/onemind-codex"],
    "indexing": {
      "enabled": true,
      "recursive": true
    }
  },

  "logging": {
    "level": "info",
    // "debug" for verbose output during troubleshooting
    "file": "/root/.openclaw/logs/openclaw.log",
    "maxSize": "10mb",
    "maxFiles": 5
  }
}
```

---

## Key Fields Reference

| Field | What It Does |
|-------|-------------|
| `agents[].id` | Used in CLI: `openclaw agent chat [id]` |
| `agents[].workspace` | Directory with IDENTITY.md, CODEX.md, USER.md |
| `channels.discord.routingRules` | Routes Discord channel messages to specific agents |
| `cron[].schedule` | Standard cron syntax (min hr day mo weekday) |
| `mcp.servers` | External tools available to all agents |
| `memory.extraPaths` | Directories included in semantic memory search |

---

## Common Edits

**Add a new agent:** Copy the agent block in `agents[]`, change `id`, `name`, `workspace`. Create the workspace directory and add IDENTITY.md.

**Add a cron job:** Add a new object to `cron[]`. Set `enabled: false` while testing.

**Disable a channel:** Set `"enabled": false` — does not delete config, just deactivates.

**Change model for all agents:** Update `models.default` — all agents using `"model": "default"` inherit it automatically.

---

*Part of the One Mind OS OpenClaw Deep Dive — Course 3*
