---
title: "MindsDB — Self-Hosted AI Agent + Data Query Engine"
type: reference
domain: SO
status: active
created: 2026-05-26
updated: 2026-05-27
tags: [ai-agents, open-source, database, mindsdb, self-hosted, tools]
---

# MindsDB

**Docs:** https://docs.mindsdb.com
**GitHub:** https://github.com/mindsdb/mindsdb
**License:** Open source (free, self-hosted)
**Cloud/Enterprise version:** Minds (mindshub.ai) — paid

---

## Two Products — Know the Difference

| | MindsDB | Minds (MindsHub) |
|-|---------|-----------------|
| Cost | Free / open source | Paid / hosted |
| Hosting | Self-hosted (Docker, PyPI) | Managed cloud |
| Agents | Yes — basic NL agents over data | Yes — production-grade, more capable |
| Use case | Dev / self-hosted / OneMind OS | Enterprise / teams |

**For Zeus: MindsDB (self-hosted) is the play.**

---

## What It Does

Three-step model:

1. **Connect** — plug in any data source (Postgres, MySQL, Snowflake, BigQuery, CSVs, APIs, vector DBs, Salesforce, 100+ integrations)
2. **Unify** — create knowledge bases, views, and jobs across those sources
3. **Respond** — run AI agents or expose an MCP API that answers NL questions over that unified data

In plain terms: you connect your databases, and then ask questions in plain English. MindsDB figures out which tables to query, writes and runs the SQL, and returns a synthesized answer.

---

## MindsDB Agents (Open Source)

Built on Pydantic. Free, included in the self-hosted install. How they work:

1. **Input Processing** — builds a real-time data catalog (5-row sample from each connected table), extracts the query
2. **Planning** — determines relevant sources, plans SQL queries
3. **Exploration Loop** — iterative query execution, up to 20 queries per request, adjusts as needed
4. **Error Handling** — catches SQL errors, retries with corrections (up to 3 accumulated errors retained as context)
5. **Output Synthesis** — aggregates results, returns NL or structured response

---

## What's Needed to Install

**Recommended: Docker Desktop** (easiest, works on Mac M2)

### Option A — Docker Desktop Extension (recommended)
```bash
# Prerequisites: Docker Desktop installed
# 1. Go to Extensions in Docker Desktop → search MindsDB → Install
# 2. MindsDB runs at http://127.0.0.1:47334
# 3. Add data integrations via Settings → Manage Integrations
```
Image used: `mindsdb/mindsdb:latest`
Comes pre-bundled with: MySQL, PostgreSQL, Snowflake, BigQuery, MSSQL, Salesforce

### Option B — PyPI / Source (for dev/contribution)
```bash
# Requirements: Python 3.10–3.11 (64-bit), pip >= 20.3, ~6 GB free disk

git clone https://github.com/mindsdb/mindsdb.git
python -m venv mindsdb-venv
source mindsdb-venv/bin/activate
cd mindsdb
pip install -e .
pip install -r requirements/requirements-dev.txt
python -m mindsdb
# MindsDB Studio → http://127.0.0.1:47334
```

### Optional: specific integrations
```bash
pip install mindsdb[handler_name]
# e.g. pip install mindsdb[postgres]
```

---

## OneMind OS Integration Path

MindsDB exposes an MCP API — meaning Hermes/Legacy can connect to it as an MCP server and query any connected database via natural language.

Potential use:
- Connect OneMind vault data (SQLite, Postgres) → ask questions across all life domains
- Hook in Google Sheets, Airtable, Notion as data sources
- Run agents that answer "what's my financial status" or "summarize last 30 days of tasks" by querying live data

---

## Live Instance on Zeus's Mac (as of 2026-05-27)

**Status: RUNNING**

| Item | Value |
|------|-------|
| Container | `mindsdb` (Docker) |
| Studio UI | http://127.0.0.1:47334 |
| MySQL API | port 47335 |
| Storage volume | `mindsdb_data` |
| Version | 26.1.0 |
| ML Engine | `openai_azure` → Azure OpenAI `onemind-openai` endpoint |
| Agent | `onemind_agent` using `gpt-4o` deployment |

### Start/Stop

```bash
docker start mindsdb      # start
docker stop mindsdb       # stop
docker logs mindsdb -f    # view logs
```

### Query the agent

```bash
curl -s -X POST http://127.0.0.1:47334/api/sql/query \
  -H "Content-Type: application/json" \
  -d '{"query": "SELECT * FROM mindsdb.onemind_agent WHERE question = '\''your question here'\'';"}' \
  | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data'][0][0])"
```

### First test result

Query: "What is OneMind OS and what is your role?"
Response: Answered as a MindsDB SQL analyst — agent is live and connected to Azure OpenAI gpt-4o. Next step: connect data sources (vault, Postgres, Airtable) and retrain with OneMind context.

---

## Links

- Docs home: https://docs.mindsdb.com/homepage
- Agents docs: https://docs.mindsdb.com/agents/agent
- Docker Desktop install: https://docs.mindsdb.com/setup/self-hosted/docker-desktop
- PyPI/source install: https://docs.mindsdb.com/setup/self-hosted/pip/source
- GitHub: https://github.com/mindsdb/mindsdb
- MCP API: https://docs.mindsdb.com/mcp-api/overview
