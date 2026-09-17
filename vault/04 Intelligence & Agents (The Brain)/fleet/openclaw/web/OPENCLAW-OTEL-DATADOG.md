---
title: "OpenClaw → Datadog APM via OTEL"
type: sop
status: active
created: 2026-04-19
owner: zeus
tags: [type/sop, domain/SO, domain/01, status/active, system/openclaw, system/datadog, system/observability]
ai_summary: "Full OTEL→Datadog APM pipeline for OpenClaw. Traces every agent request and tool call to Datadog APM via OTLP gRPC. Enabled May 2026 on onemindos-cloud-1."
---

# OpenClaw → Datadog APM via OTEL

> **VPS**: `onemindos-cloud-1` (`64.227.8.23`, Tailscale `100.98.25.10`)
> **Implemented**: 2026-04-19
> **Status**: ✅ ACTIVE — both services running and pipeline confirmed

---

## What This Does

Every agent request, tool call, and LLM inference that passes through the OpenClaw gateway is automatically traced and sent to **Datadog APM**. This gives full observability into:

- Which agents are being called and how often
- Which tools are executing (tool call duration, success/fail)
- LLM provider latency per request
- Memory operations (active-memory, memory-core, memory-wiki reads/writes)
- Full distributed trace trees per conversation turn

This is the **operational nervous system** for One Mind — you can see exactly what Legacy, Spartan, Empire, and all agents are doing in production.

---

## Architecture

```
OpenClaw Gateway (VPS)
    ↓ OTEL SDK (built-in diagnostics-otel plugin)
    ↓ OTLP/gRPC → localhost:4317
Datadog Agent (VPS, systemd service)
    ↓ OTLP receiver (grpc:4317 + http:4318)
    ↓ forwards to Datadog cloud
Datadog APM
    → APM → Services → openclaw-gateway
    → Traces, spans, latency, errors
```

**Key insight**: OpenClaw's `diagnostics-otel` plugin instruments the gateway internally and exports spans via the standard OTEL SDK. The Datadog Agent acts as a local OTLP collector — no Datadog-specific SDK needed in OpenClaw.

---

## Configuration

### 1. Plugin Enabled

```bash
openclaw plugins enable diagnostics-otel
```

Verify on VPS:
```bash
openclaw plugins list | grep otel
```

### 2. OpenClaw Environment Variables

**File**: `/root/.config/systemd/user/openclaw-gateway.service.d/ollama-env.conf`

```ini
Environment="OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317"
Environment="OTEL_SERVICE_NAME=openclaw-gateway"
```

These tell the OTEL SDK to send spans to the local Datadog Agent OTLP receiver via gRPC.

> **1Password reference**: No secret needed for OTEL — no key required for local OTLP.

### 3. Datadog Agent OTLP Receiver

**File**: `/etc/datadog-agent/datadog.yaml` (bottom of file)

```yaml
otlp_config:
  receiver:
    protocols:
      grpc:
        endpoint: localhost:4317
      http:
        endpoint: localhost:4318
```

This enables the Datadog Agent to accept OTEL traces locally and forward them to Datadog cloud.

> **Datadog site**: `us5.datadoghq.com`
> **API key reference**: `op://aaeyris5z44pdqspk277u2oppi/Datadog/api-key`

---

## How to View in Datadog

1. Log in at `app.us5.datadoghq.com`
2. Navigate to **APM → Services**
3. Look for service: **`openclaw-gateway`**
4. Drill into **Traces** to see per-request span trees:
   - Root span: agent request (e.g. `legacy.run`)
   - Child spans: tool calls, LLM provider calls, memory reads

**What you'll see per agent turn:**
- Agent name + model used
- All tool calls with duration
- LLM provider (ollama/qwen2.5, anthropic, etc.) + latency
- Memory operations (search, read, write)
- Errors and timeouts

---

## Related: Tavily API Key (same session)

During this session, the **Tavily search plugin** was also enabled and its API key injected:

**File**: `/root/.config/systemd/user/openclaw-gateway.service.d/ollama-env.conf`

```ini
Environment="TAVILY_API_KEY=<redacted>"
```

> **1Password reference**: `op://aaeyris5z44pdqspk277u2oppi/Tavily ApI/credential`
> **Item ID**: `qvhwtmrhvkpxq4hstfb2gjdgpe`

```bash
openclaw plugins enable tavily
```

Tavily gives agents real-time web search with structured results — significantly better signal-to-noise than DuckDuckGo for research tasks.

---

## Service Management

### Restart services after config changes:

```bash
# Datadog Agent (root-level service)
systemctl restart datadog-agent

# OpenClaw Gateway (user-level service)
systemctl --user restart openclaw-gateway
```

### Check status:

```bash
systemctl status datadog-agent --no-pager
systemctl --user status openclaw-gateway --no-pager
```

### Check OTEL data is reaching Datadog Agent:

```bash
# Look for OTLP trace receiver messages
journalctl -u datadog-agent --since "5 minutes ago" | grep -i otlp
```

---

## Plugin State (as of 2026-04-19)

### Newly Enabled This Session

| Plugin | Purpose |
|--------|---------|
| `diagnostics-otel` | OTEL tracing → Datadog APM |
| `tavily` | Real-time web search (requires `TAVILY_API_KEY`) |
| `duckduckgo` | Free web search (no key) |
| `open-prose` | Writing mode + slash commands |
| `telegram` | Telegram bot channel (⏳ needs bot token) |
| `discord` | Discord channel (⏳ needs bot token) |

### Newly Disabled This Session

China/regional providers: `alibaba`, `byteplus`, `chutes`, `kimi`, `minimax`, `moonshot`, `qianfan`, `qwen`, `stepfun`, `volcengine`, `vydra`, `xiaomi`, `zai`
Routing: `amazon-bedrock-mantle`, `thread-ownership`

---

## Files Modified on VPS

| File | Change |
|------|--------|
| `/root/.openclaw/openclaw.json` | Plugin enable/disable state updated |
| `/root/.config/systemd/user/openclaw-gateway.service.d/ollama-env.conf` | Added `TAVILY_API_KEY`, `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_SERVICE_NAME` |
| `/etc/datadog-agent/datadog.yaml` | Appended `otlp_config` block (gRPC 4317, HTTP 4318) |

---

## Why This Matters

This is the first time One Mind has **full agent observability**. Previously, you had no visibility into what Legacy or any other agent was actually doing — tool calls were invisible, LLM latency was unknown, errors silently failed.

Now every agent turn is a traceable event. When an agent is slow, you can see exactly which tool is blocking. When an agent fails, you see the exact span that errored. This is the foundation for performance tuning, debugging, and eventually **SLO/SLA tracking per domain agent**.

> *"If it moves, measure it. If you measure it, you can improve it."*
