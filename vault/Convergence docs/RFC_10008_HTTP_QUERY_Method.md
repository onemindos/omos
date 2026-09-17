# RFC 10008: The HTTP QUERY Method

**Official Link:** https://www.rfc-editor.org/info/rfc10008/

**Status:** PROPOSED STANDARD (IETF Internet Standards Track)

**Authors:** 
- Julian Reschke
- James M. Snell
- Mike Bishop

**Published:** 2026 (IETF)

**RFC Editor Home:** https://www.rfc-editor.org/

**Related RFCs:**
- RFC 7841 (Internet Standards Process)
- RFC 9110 (HTTP Semantics)
- RFC 3986 (URI Generic Syntax)

---

## What Is RFC 10008?

**RFC 10008** defines the **HTTP QUERY method** — a new HTTP request method designed for **safe, idempotent queries** with **request body content** (unlike GET, which uses URI parameters).

### Core Problem Solved

**The Challenge:** 
When query data is too large for URI parameters, developers resort to POST requests. But POST is semantically ambiguous — it looks like state-changing, non-idempotent operation, even when performing safe queries.

**The Solution:**
QUERY method explicitly indicates:
- ✅ Safe operation (read-only, no side effects)
- ✅ Idempotent (can be safely repeated/restarted)
- ✅ Request body contains query parameters
- ✅ Can be automatically retried without concern for partial state

---

## Problem Context

### Issue 1: URI Size Limits

```http
GET /feed?q=foo&limit=10&sort=-published HTTP/1.1
Host: example.org
```

**Problem:** 
- Size limits unknown (passes through many systems)
- HTTP recommends minimum 8000 octets, but not guaranteed
- Large complex queries can exceed limits

### Issue 2: URI Encoding Inefficiency

**Problem:**
- Encoding complex data into URI is inefficient
- Overhead of URI-safe encoding (special characters → %XX)
- Many data types don't naturally fit URI format

### Issue 3: Privacy & Logging

**Problem:**
- Request URIs logged more than request content
- Query parameters appear in bookmarks, browser history
- Sensitive data exposed in logs/analytics

### Issue 4: Resource Semantics

**Problem:**
- Every query combination becomes distinct URI resource
- Violates REST principles (collapsing query space into resource space)
- Creates unnecessary URL proliferation

---

## Current Workaround: POST for Queries

```http
POST /feed HTTP/1.1
Host: example.org
Content-Type: application/x-www-form-urlencoded

q=foo&limit=10&sort=-published
```

**Problem:** Semantically unclear
- POST typically means "create" or "change state"
- Server needs custom logic to recognize safe queries
- Clients can't reliably retry (POST is not idempotent)
- Caching systems treat it as non-cacheable

---

## HTTP QUERY Method Specification

### Semantics

**QUERY is:**
- ✅ **Safe** — does not alter server state
- ✅ **Idempotent** — can be repeated/restarted without additional side effects
- ✅ **Cacheable** — responses can be cached (with appropriate headers)
- ✅ **Request body** — includes query parameters/content

### HTTP Semantics

```http
QUERY /feed HTTP/1.1
Host: example.org
Content-Type: application/x-www-form-urlencoded
Content-Length: 35

q=foo&limit=10&sort=-published
```

**Key Differences from POST:**
| Aspect | POST | QUERY |
|--------|------|-------|
| **Semantic** | State-changing | Safe query |
| **Idempotent** | No (not safe to retry) | Yes (safe to retry) |
| **Cacheable** | No (by default) | Yes (with headers) |
| **Automatic Retry** | Not recommended | Safe & recommended |
| **Request Target** | Creates/modifies resource | Queries resource |

---

## Use Cases

### 1. Complex Search Queries

**Problem:** Search filters too complex for URI

```http
QUERY /search HTTP/1.1
Host: api.example.org
Content-Type: application/json

{
  "q": "machine learning",
  "filters": {
    "date_from": "2024-01-01",
    "date_to": "2026-07-11",
    "language": "en",
    "source": ["github", "arxiv", "blogs"]
  },
  "limit": 50,
  "sort": "-relevance"
}
```

### 2. Large Form Submissions for Filtering

```http
QUERY /products HTTP/1.1
Host: ecommerce.example.org
Content-Type: application/x-www-form-urlencoded

category=electronics&
price_min=100&price_max=5000&
brand=Apple|Samsung|Google&
in_stock=true&
ratings_min=4.0&
sort=-popularity
```

### 3. Document/Database Queries

```http
QUERY /documents HTTP/1.1
Host: db.example.org
Content-Type: application/json

{
  "collection": "users",
  "filter": {
    "age": {"$gte": 18, "$lte": 65},
    "country": "US",
    "active": true
  },
  "projection": ["name", "email", "created_at"],
  "sort": {"created_at": -1},
  "limit": 100,
  "skip": 0
}
```

### 4. API Rate Limiting with Query Details

```http
QUERY /api/process HTTP/1.1
Host: processor.example.org
Content-Type: application/json

{
  "operation": "analyze_data",
  "data": [/* large array of objects */],
  "options": {
    "parallel": true,
    "timeout": 30000,
    "format": "json"
  }
}
```

---

## Advantages Over POST

### 1. Semantic Clarity

| Method | Meaning |
|--------|---------|
| GET | Retrieve resource (query in URI) |
| POST | Create/modify resource |
| **QUERY** | **Query resource (query in body)** |

### 2. Safe Retry

**Automatic retries are safe:**
- Network failure → retry without concern
- Proxy can automatically retry on timeout
- Client doesn't need custom retry logic
- Matches idempotent semantics of GET

### 3. Caching

**Responses are cacheable (with appropriate headers):**
- Content-Delivery Networks (CDNs) can cache QUERY responses
- Proxy servers can cache results
- Conditional requests (If-Modified-Since, ETag) work
- Reduces server load for repeated queries

### 4. Logging & Privacy

**Query data NOT in URI:**
- Not exposed in browser history
- Not logged in access logs (by default)
- Not visible in bookmarks
- Better privacy compliance

### 5. Intermediary Handling

**Proxies/gateways understand intent:**
- Can apply query-specific caching rules
- Can enforce rate limiting per query type
- Can route based on query complexity
- Can prioritize based on idempotence

---

## Comparison to Existing Methods

| Aspect | GET | POST | PUT | PATCH | **QUERY** |
|--------|-----|------|-----|-------|-----------|
| **Safe** | ✅ | ❌ | ❌ | ❌ | ✅ |
| **Idempotent** | ✅ | ❌ | ✅ | ❌ | ✅ |
| **Cacheable** | ✅ | ❌* | ❌ | ❌ | ✅ |
| **Request Body** | ❌ | ✅ | ✅ | ✅ | ✅ |
| **Query in URI** | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Query in Body** | ❌ | ✅ (misused) | ❌ | ❌ | ✅ |

*POST can be cached but rarely is due to semantic ambiguity

---

## Implementation Considerations

### Server-Side

**Servers must:**
1. Recognize QUERY method in request routing
2. Process request body as query parameters
3. Return appropriate status codes:
   - **200 OK** — query succeeded, results included
   - **204 No Content** — query succeeded, no results
   - **400 Bad Request** — malformed query
   - **404 Not Found** — target resource not found
   - **415 Unsupported Media Type** — unsupported Content-Type

### Client-Side

**Clients should:**
1. Use QUERY for safe, idempotent queries with body content
2. Set appropriate Content-Type header
3. Can automatically retry on network failure
4. Can include caching headers (Cache-Control, ETag, etc.)

### Content Types

**Supported (recommended):**
- `application/x-www-form-urlencoded` (form data)
- `application/json` (JSON queries)
- `application/xml` (XML queries)
- Custom content types specific to API

---

## HTTP Status Codes

### Success Responses

```
200 OK           — Query processed, results returned
204 No Content   — Query processed, no results
```

### Client Error Responses

```
400 Bad Request           — Query syntax error
404 Not Found             — Resource not found
415 Unsupported Media Type — Content-Type not accepted
```

### Server Error Responses

```
500 Internal Server Error — Server-side error
501 Not Implemented       — QUERY not supported
503 Service Unavailable   — Server overloaded
```

---

## Caching with QUERY

### Caching Headers

**For cacheable QUERY responses:**

```http
QUERY /search HTTP/1.1
Host: api.example.org
Content-Type: application/json

{...query...}

---

HTTP/1.1 200 OK
Content-Type: application/json
Cache-Control: public, max-age=3600
ETag: "abc123"
Vary: Accept-Encoding

{...results...}
```

**Cache-Control Options:**
- `public` — any cache can store
- `private` — only client cache
- `max-age=3600` — cache for 1 hour
- `must-revalidate` — check freshness before using

### Conditional Requests

```http
QUERY /search HTTP/1.1
Host: api.example.org
If-None-Match: "abc123"
Content-Type: application/json

{...query...}

---

HTTP/1.1 304 Not Modified
```

---

## Security Considerations

### 1. Query Injection

**Risk:** Malicious query content could cause harm

**Mitigation:**
- Validate all query parameters
- Use parameterized queries (prepared statements)
- Sanitize input
- Apply principle of least privilege

### 2. DoS (Denial of Service)

**Risk:** Large queries or complex operations overload server

**Mitigation:**
- Enforce Content-Length limits
- Set timeouts on query execution
- Implement rate limiting per client/IP
- Cache common queries

### 3. Information Disclosure

**Risk:** Query results expose sensitive data

**Mitigation:**
- Enforce authentication/authorization
- Filter results based on user permissions
- Log access attempts
- Use HTTPS encryption

### 4. HTTPS Requirement

**Strongly recommended:**
- All QUERY operations should use HTTPS
- Encrypts request body content
- Prevents intermediary inspection
- Protects privacy

---

## Integration with OneMind/Convergence

### API Communication Pattern

**HTTP QUERY + OneMind NATS:**
- **QUERY** for complex API requests (safe, idempotent, with body)
- **NATS** for event streaming, pub/sub, real-time messaging
- **Together:** REST queries + event-driven architecture

### Applications

1. **Distributed System Queries**
   - Agent queries TAK server (large CoT filter criteria)
   - QUERY method signals safe, repeatable operation
   - Results can be cached at edge nodes

2. **Tactical Database Queries**
   - Complex threat intelligence queries
   - Large filter/sort criteria
   - Safe to retry on network failure
   - CDN cacheable at tactical edge

3. **Microservice Communication**
   - OneMind services query each other
   - QUERY for read-heavy operations
   - NATS for event notifications
   - Combined: efficient query + real-time updates

4. **API Gateway Pattern**
   - API gateway accepts QUERY requests
   - Routes to backend services
   - Caches responses
   - Implements rate limiting per query type

---

## Adoption Timeline

### Current Status (2026)

- **Status:** PROPOSED STANDARD
- **IETF Track:** Internet Standards Track
- **Expected:** Transition to STANDARD in 2027-2028
- **Implementation:** Early adoption by API platforms

### Predicted Adoption

| Timeline | Adoption |
|----------|----------|
| **2026** | IETF proposal, early adopters |
| **2027** | Major API platforms (GitHub, AWS, etc.) |
| **2028** | STANDARD status, mainstream adoption |
| **2029+** | Industry standard for body-based queries |

---

## Comparison to GraphQL & gRPC

| Aspect | REST+QUERY | GraphQL | gRPC |
|--------|-----------|---------|------|
| **Query Type** | HTTP QUERY (native) | Special language | RPC (binary) |
| **Request Body** | ✅ Yes | ✅ Yes | ✅ Yes (binary) |
| **Safe/Idempotent** | ✅ Yes | Depends | Depends |
| **Cacheable** | ✅ Yes | Difficult | Difficult |
| **HTTP Native** | ✅ Yes | Partial | No (protobuf) |
| **Learning Curve** | ✅ Low | High | Medium |
| **Tooling** | Standard HTTP | Custom tools | Custom tools |

**HTTP QUERY Advantage:** Native HTTP method, standard caching, minimal learning curve

---

## Examples

### Example 1: Search API

```http
QUERY /api/v1/articles/search HTTP/1.1
Host: content.example.org
Content-Type: application/json
Accept: application/json

{
  "text": "machine learning",
  "filters": {
    "published_after": "2024-01-01",
    "tags": ["AI", "neural-networks"],
    "min_views": 1000
  },
  "sort": "-published_date",
  "limit": 20
}

---

HTTP/1.1 200 OK
Content-Type: application/json
Cache-Control: public, max-age=300
Content-Length: 5432

{
  "total": 1247,
  "results": [
    {"id": "art_001", "title": "...", "published": "2026-07-10"},
    ...
  ]
}
```

### Example 2: Log Analysis Query

```http
QUERY /api/logs/query HTTP/1.1
Host: logging.example.org
Content-Type: application/json

{
  "time_range": {
    "start": "2026-07-11T00:00:00Z",
    "end": "2026-07-11T23:59:59Z"
  },
  "filters": {
    "level": "ERROR",
    "service": "auth-service",
    "region": "us-west"
  },
  "group_by": ["error_code"],
  "limit": 100
}
```

### Example 3: Database Query

```http
QUERY /api/database/query HTTP/1.1
Host: db.example.org
Content-Type: application/json

{
  "collection": "transactions",
  "filter": {
    "status": "pending",
    "amount": {"$gt": 1000, "$lt": 10000},
    "timestamp": {"$gte": "2026-07-01"}
  },
  "projection": ["id", "user_id", "amount", "created_at"],
  "sort": {"amount": -1},
  "limit": 50
}
```

---

## Resources

### Official

- **RFC 10008 Full Text:** https://www.rfc-editor.org/info/rfc10008/
- **RFC Editor Home:** https://www.rfc-editor.org/
- **IETF:** https://www.ietf.org/
- **HTTP Semantics (RFC 9110):** https://www.rfc-editor.org/info/rfc9110

### Related Standards

- **HTTP Status Codes:** https://httpwg.org/specs/rfc9110.html#status.codes
- **Content Negotiation:** https://httpwg.org/specs/rfc9110.html#content.negotiation
- **Caching:** https://www.rfc-editor.org/info/rfc9111
- **URI Syntax (RFC 3986):** https://www.rfc-editor.org/info/rfc3986

### Implementation Guidance

- **Web Framework Support:** Check your framework docs (Express, Django, Rails, etc.)
- **API Gateway:** Kong, AWS API Gateway, Nginx
- **HTTP Client Libraries:** httpx, axios, requests (Python)

---

## Summary

**RFC 10008: HTTP QUERY Method provides:**

✅ Native HTTP method for safe, idempotent queries with request body
✅ Semantic clarity (distinguishes safe queries from POST state-changes)
✅ Automatic retry safety (idempotent like GET, body like POST)
✅ Cacheable responses (reduces server load)
✅ Privacy (query in body, not URI)
✅ Better intermediary handling (proxies understand intent)

**Status:** PROPOSED STANDARD (2026), expected mainstream adoption 2027-2028

**For OneMind Convergence:** Standard protocol for distributed system queries + NATS event streaming = complete API communication layer

---

**Captured:** 2026-07-11

**Official Source:** https://www.rfc-editor.org/info/rfc10008/

**Status:** IETF Proposed Standard, Internet Standards Track
