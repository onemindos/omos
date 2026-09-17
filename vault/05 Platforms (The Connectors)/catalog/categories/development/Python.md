# Python

> **Primary scripting and AI/ML language across all One Mind infrastructure**

## Overview

Python is the dominant language in the One Mind stack. Used for agent code, automation scripts, data processing, ML/AI pipelines, and VPS tooling. Combined with FastAPI for APIs and with the OpenAI/Anthropic SDK for agent integration.

## Status

| Attribute   | Value                                       |
| ----------- | ------------------------------------------- |
| Status      | ✅ Active                                    |
| Category    | Development / Language                      |
| Type        | Programming Language                        |
| Integration | pip, poetry, venv, pyenv, conda             |

## Use Cases in OneMind

1. **Agent Scripts** — VPS automation, cron jobs, vault management scripts
2. **AI Pipelines** — OpenAI, Anthropic SDK calls from Python
3. **Data Processing** — Dedup scripts, file migration, vault cleanup tools
4. **FastAPI Backend** — API layer for internal services
5. **OpenClaw Skills** — Python-based agent tool implementations

## Core Libraries in Stack

| Library       | Use                              |
| ------------- | -------------------------------- |
| `anthropic`   | Claude API client                |
| `openai`      | OpenAI API + embeddings          |
| `fastapi`     | REST API backend                 |
| `uvicorn`     | ASGI server for FastAPI          |
| `psycopg2`    | PostgreSQL adapter               |
| `redis`       | Redis client                     |
| `requests`    | HTTP client                      |
| `composio`    | Composio tool integration        |

## VPS Python

```bash
python3 --version          # Python 3.12 on Ubuntu 24.04
pip install -r requirements.txt
# Recommended: use venv
python3 -m venv .venv && source .venv/bin/activate
```

## Links

- Website: [python.org](https://python.org)
- PyPI: [pypi.org](https://pypi.org)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
