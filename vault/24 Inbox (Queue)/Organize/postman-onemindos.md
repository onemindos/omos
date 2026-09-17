# Postman × One Mind OS
**How to leverage Postman across the One Mind OS stack**

---

## 1. API Testing Cockpit

Before wiring any integration into OpenClaw or your agent layer, verify it works first in Postman.

- Test OpenClaw endpoints before connecting to agent workflows
- Validate Perplexity Sonar API responses before building search tools
- Test Supabase database calls before committing to code
- Inspect NATS JetStream HTTP API responses
- Debug any webhook, REST API, or integration in isolation

**Rule:** If it's an API, test it in Postman before you build around it.

---

## 2. Collections — Your Stack API Library

Organize every endpoint in One Mind OS into structured Postman Collections.

**Recommended Collection Structure:**
```
One Mind OS/
├── OpenClaw/
│   ├── Gateway Endpoints
│   ├── Agent Loop Triggers
│   └── Tool Calls
├── Supabase/
│   ├── Auth
│   ├── Database Queries
│   └── Realtime Subscriptions
├── NATS JetStream/
│   ├── Publish
│   ├── Subscribe
│   └── Stream Management
├── Perplexity Sonar/
│   ├── Search Queries
│   └── Deep Research
└── External Integrations/
    ├── Zoho
    ├── ClickUp
    ├── GoDaddy
    └── Cloudflare
```

This becomes the **single source of truth** for every API connection in your stack.

---

## 3. Environments — Dev / Staging / Production

Switch your entire stack context with one click.

| Variable | Dev | Staging | Production |
|---|---|---|---|
| `OPENCLAW_URL` | localhost:8080 | staging.legacyhq.ai | legacyhq.ai |
| `SUPABASE_URL` | local instance | staging project | prod project |
| `API_KEY` | dev key | staging key | prod key |
| `NATS_URL` | localhost:4222 | staging NATS | prod NATS |

No more manually swapping URLs or keys. One environment toggle switches everything.

---

## 4. Automated Testing — Catch Breaks Early

Write test scripts that run after every API call and assert expected behavior.

**Example use cases:**
- Assert OpenClaw gateway returns correct agent response schema
- Verify Supabase auth token is valid and not expired
- Confirm NATS message publish returns success acknowledgment
- Validate Perplexity Sonar returns citations with every response
- Check OMCS coordinate lookups return correct pillar mappings (HP, LE, GE, IT)

Run entire test suites before any deployment. Know what broke before it hits production.

---

## 5. Monitors — Lightweight Uptime Alerts

Schedule Postman to automatically ping your critical endpoints on a timer.

**What to monitor:**
- OpenClaw Gateway health check
- Supabase connection status
- NATS JetStream availability
- Tailscale mesh node connectivity checks
- Home Assistant API responsiveness

Get alerted the moment something goes down — before your agents do.

---

## 6. Mock Servers — Build Before It's Live

Simulate APIs that don't exist yet so you can build against them immediately.

**One Mind OS use cases:**
- Mock OpenClaw endpoints while the actual gateway is still being configured
- Simulate OMCS coordinate API responses while the system is in development
- Build Legacy community frontend against a mocked backend API
- Test agent workflows against simulated tool responses

Build in parallel. Don't wait for infrastructure to catch up to your vision.

---

## 7. Flows — Visual API Orchestration

Chain API calls together visually — similar to n8n but purely for API workflows.

**Example flows:**
- Trigger OpenClaw agent → parse response → write result to Supabase
- Query Perplexity Sonar → format output → push to NATS topic
- ClickUp task created → enrich with AI → update Zoho CRM

Use Flows for lightweight integrations that don't need the full n8n treatment.

---

## 8. API Documentation — Onboarding Without the Headache

Postman auto-generates documentation from your Collections.

**When this matters for Legacy Platform:**
- Onboarding contractors or developers — hand them a shared Collection
- Documenting the One Mind OS API surface for community builders
- Creating internal API references for the OMCS coordinate system
- Certifying Legacy Platform partners who need to integrate with your stack

Write the API once in Postman. Documentation is automatic.

---

## 9. Request Chaining — Dynamic Workflows

Pass output from one API call as input to the next automatically.

**Example chain:**
1. Authenticate with Supabase → capture token
2. Use token to query OpenClaw gateway → capture agent session ID
3. Use session ID to trigger specific NATS message → capture acknowledgment
4. Log full chain result to Supabase

Test complex multi-step workflows end-to-end without writing a single line of code.

---

## 10. Version Control for APIs

Track changes to your API definitions over time.

- Know exactly when an OpenClaw endpoint changed behavior
- Roll back to a working configuration if an integration breaks
- Document the evolution of One Mind OS API surface as the stack grows

---

## Priority Actions

| Priority | Action |
|---|---|
| **Immediate** | Create One Mind OS Collection structure |
| **Immediate** | Set up Dev / Production environments |
| **Week 1** | Add all OpenClaw endpoints to Collection |
| **Week 1** | Add Supabase endpoints and auth flows |
| **Week 2** | Write automated tests for critical endpoints |
| **Week 2** | Set up monitors for OpenClaw and Supabase |
| **Month 1** | Build mock servers for OMCS coordinate API |
| **Ongoing** | Auto-generate docs as stack grows |

---

*One Mind OS Stack: OpenClaw → NATS JetStream → Supabase → Tailscale Mesh*
*Legacy Platform Holdings LLC — Confidential*
