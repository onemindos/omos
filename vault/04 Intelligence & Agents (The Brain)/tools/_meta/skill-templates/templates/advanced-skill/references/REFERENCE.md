# Advanced Skill — Full Reference

This document contains the complete technical reference for the advanced-skill.
It is loaded on demand by Legacy only when the main SKILL.md instructions
point to it. Keep this file focused and factual.

---

## API Reference

### Base URL

`ADVANCED_SKILL_URL` env var (e.g., `https://api.example.com`)

### Authentication

All requests use Bearer token authentication:
```
Authorization: Bearer <ADVANCED_SKILL_KEY>
```

---

## POST /process

Process input data.

**Request body:**
```json
{
  "input": "value or identifier",
  "mode": "standard | fast | thorough"
}
```

**Response:**
```json
{
  "id": "proc_abc123",
  "status": "completed",
  "input": "...",
  "output": "...",
  "duration_ms": 450
}
```

**Mode differences:**
| Mode | Speed | Accuracy | Use When |
|---|---|---|---|
| `fast` | ~100ms | 80% | Quick sanity check |
| `standard` | ~500ms | 95% | Normal use |
| `thorough` | ~2000ms | 99% | High-stakes decisions |

---

## GET /analyze/:target

Analyze a target.

**Parameters:**
- `target` (path) — URL-encoded target identifier

**Query parameters:**
- None for summary analysis

**Response:**
```json
{
  "target": "...",
  "summary": "...",
  "score": 0.85,
  "tags": ["..."]
}
```

---

## GET /report

Generate a report.

**Query parameters:**
- `format`: `markdown` | `json`
- `timeframe`: `1h` | `24h` | `7d` | `30d`

**Response:**
```json
{
  "generated_at": "2026-03-16T12:00:00Z",
  "timeframe": "24h",
  "content": "..."
}
```

---

## Status Codes

| Code | Meaning | Action |
|---|---|---|
| 200 | OK | Success |
| 204 | No Content | Success, no body |
| 400 | Bad Request | Check params |
| 401 | Unauthorized | Check API key |
| 404 | Not Found | Target/ID doesn't exist |
| 429 | Rate Limited | Wait and retry |
| 500 | Server Error | Retry once; report if persists |

---

## Rate Limits

- 60 requests/minute per API key
- 1000 requests/day per API key
- If rate limited, wait 60 seconds before retrying

---

## Pagination

List endpoints support:
```
GET /items?limit=20&offset=0
```

Response includes:
```json
{
  "items": [...],
  "total": 150,
  "limit": 20,
  "offset": 0
}
```


---

## Related

- [[00-24 SO (Sovereign Operations)/01 Intelligence (AI, Memory & Tools) (The Brain)/01 Intelligence|01 Intelligence]]
