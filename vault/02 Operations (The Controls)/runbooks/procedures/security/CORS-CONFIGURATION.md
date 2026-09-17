# CORS Configuration and Testing

## Overview

This document describes the Cross-Origin Resource Sharing (CORS) configuration for OneMind OS and the testing performed to verify UI-to-API communication.

## Current Configuration

The CORS middleware is configured in `agents/legacy/main.py` (lines 321-327):

```python
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # TODO: Restrict to TOC domain in production
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
```

### Configuration Details

| Setting | Value | Description |
|---------|-------|-------------|
| `allow_origins` | `["*"]` | **Currently allows ALL origins** - any domain can access the API |
| `allow_credentials` | `True` | Allows cookies and authorization headers to be sent |
| `allow_methods` | `["*"]` | Allows all HTTP methods (GET, POST, PUT, DELETE, OPTIONS, etc.) |
| `allow_headers` | `["*"]` | Allows all custom headers |

## Test Results

All 18 CORS tests pass successfully. Test coverage includes:

### ✅ Standard Endpoints
- `/health` - Health check endpoint
- `/api/agents` - Agent listing endpoint
- `/awareness` - Awareness state endpoint
- `/sessions` - Session management endpoint

### ✅ Preflight OPTIONS Requests
- Preflight requests return 200 status
- All required CORS headers present:
  - `access-control-allow-origin`
  - `access-control-allow-methods`
  - `access-control-allow-headers`

### ✅ Streaming Endpoints (SSE)
- `/api/agents/{id}/runs` with `stream=true`
- Proper headers for Server-Sent Events:
  - `content-type: text/event-stream`
  - `cache-control: no-cache`
  - `connection: keep-alive`
- CORS headers present on streaming responses

### ✅ Multiple Origins
Tested with:
- `https://agno.onemindos.dev` (production domain)
- `https://ui.onemindos.dev` (alternative domain)
- `http://localhost:3000` (development UI)
- `http://localhost:8080` (development API)

### ✅ Custom Headers
- Authorization headers work correctly
- Custom headers (`X-Custom-Header`, etc.) are allowed
- Content-Type headers work as expected

### ✅ Edge Cases
- CORS headers present on 404 responses
- CORS headers present on 500 error responses
- Same-origin requests work without Origin header

### ✅ Real-World Scenarios
- Full chat flow with CORS (health → agents → sessions → message)
- Streaming with CORS (preflight → streaming request)

## Production Recommendations

### 🔴 Critical: Restrict Origins

The current configuration (`allow_origins=["*"]`) is **NOT recommended for production**. 

**Recommended change:**

```python
app.add_middleware(
    CORSMiddleware,
    allow_origins=[
        "https://agno.onemindos.dev",
        "https://ui.onemindos.dev",
        "https://onemindos.dev",
        # Add other production domains as needed
    ],
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["Content-Type", "Authorization", "X-Requested-With"],
)
```

For development, use environment variables:

```python
import os

# Get allowed origins from environment, default to wildcard for development
allowed_origins = os.getenv("CORS_ALLOWED_ORIGINS", "*")
if allowed_origins != "*":
    allowed_origins = [origin.strip() for origin in allowed_origins.split(",")]

app.add_middleware(
    CORSMiddleware,
    allow_origins=allowed_origins,
    allow_credentials=True,
    allow_methods=["GET", "POST", "PUT", "DELETE", "OPTIONS"],
    allow_headers=["Content-Type", "Authorization", "X-Requested-With"],
)
```

Environment variable example:
```bash
# Production
CORS_ALLOWED_ORIGINS="https://agno.onemindos.dev,https://ui.onemindos.dev"

# Development
CORS_ALLOWED_ORIGINS="*"
```

### Security Considerations

1. **Origin Restriction**
   - Limit to known production domains
   - Use environment variables for flexibility
   - Never use `"*"` with `allow_credentials=True` in production

2. **Credentials**
   - Only enable if you need to send cookies/auth headers
   - When enabled, cannot use `allow_origins=["*"]`

3. **Methods**
   - Restrict to only needed HTTP methods
   - Remove methods not used by your API

4. **Headers**
   - Specify exact headers instead of `"*"`
   - Add headers as needed for functionality

## UI Integration

### Frontend Configuration

The UI automatically detects the API endpoint based on environment:

```typescript
// ui/src/lib/api.ts
const API_URLS = {
  agent: process.env.NEXT_PUBLIC_AGENT_API_URL || 'https://agno.onemindos.dev',
  gateway: process.env.NEXT_PUBLIC_GATEWAY_API_URL || 'https://api.onemindos.dev'
}
```

### Streaming Support

The UI uses Server-Sent Events (SSE) for streaming responses:

```typescript
// ui/src/hooks/useAIResponseStream.tsx
const response = await fetch(apiUrl, {
  method: 'POST',
  headers: {
    ...(!(requestBody instanceof FormData) && {
      'Content-Type': 'application/json'
    }),
    ...headers
  },
  body: requestBody instanceof FormData ? requestBody : JSON.stringify(requestBody)
})

const reader = response.body.getReader()
const decoder = new TextDecoder()
```

### CORS Headers Required by Browser

For streaming to work, these headers must be present:

1. **Response headers** (from server):
   - `access-control-allow-origin: <origin>` (or `*` for dev)
   - `access-control-allow-credentials: true`
   - `content-type: text/event-stream`
   - `cache-control: no-cache`
   - `connection: keep-alive`

2. **Preflight response** (for POST/PUT/DELETE):
   - `access-control-allow-origin: <origin>`
   - `access-control-allow-methods: GET, POST, OPTIONS`
   - `access-control-allow-headers: content-type, authorization`

## Testing

### Running Tests

```bash
# Run all CORS tests
pytest tests/test_cors.py -v

# Run specific test class
pytest tests/test_cors.py::TestCORSHeaders -v

# Run integration tests
pytest tests/test_cors.py -m integration -v
```

### Test Coverage

- ✅ 18/18 tests passing
- ✅ Standard endpoints tested
- ✅ Preflight requests tested
- ✅ Streaming (SSE) tested
- ✅ Multiple origins tested
- ✅ Custom headers tested
- ✅ Edge cases tested
- ✅ Real-world scenarios tested

## Troubleshooting

### Common CORS Errors

1. **"No 'Access-Control-Allow-Origin' header is present"**
   - CORS middleware not configured
   - Origin not in allowed list
   - Server error preventing CORS headers

2. **"Credentials flag is 'true', but the 'Access-Control-Allow-Credentials' header is ''"**
   - Set `allow_credentials=True` in middleware
   - Or remove credentials from request

3. **"Request header field X is not allowed by Access-Control-Allow-Headers"**
   - Add header to `allow_headers` list
   - Or use `allow_headers=["*"]` for development

4. **Streaming stops/disconnects**
   - Check `cache-control: no-cache` header
   - Check `connection: keep-alive` header
   - Verify SSE format is correct

### Debugging CORS

1. **Check browser console**
   ```javascript
   // Browser DevTools → Console
   // Look for CORS errors (shown in red)
   ```

2. **Check network tab**
   ```
   Browser DevTools → Network → Select request → Headers
   - Check Response Headers for CORS headers
   - Check Request Headers for Origin header
   ```

3. **Test with curl**
   ```bash
   # Test preflight
   curl -X OPTIONS http://localhost:8080/api/agents \
     -H "Origin: https://agno.onemindos.dev" \
     -H "Access-Control-Request-Method: GET" \
     -v

   # Test actual request
   curl http://localhost:8080/health \
     -H "Origin: https://agno.onemindos.dev" \
     -v
   ```

## Related Files

- `agents/legacy/main.py` - CORS middleware configuration
- `ui/src/lib/api.ts` - UI API client with CORS support
- `ui/src/hooks/useAIResponseStream.tsx` - Streaming handler
- `tests/test_cors.py` - CORS test suite
- `infra/compose/cloud/.env.example` - Environment variables

## References

- [FastAPI CORS Documentation](https://fastapi.tiangolo.com/tutorial/cors/)
- [MDN CORS Documentation](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [Server-Sent Events (SSE) Specification](https://html.spec.whatwg.org/multipage/server-sent-events.html)
