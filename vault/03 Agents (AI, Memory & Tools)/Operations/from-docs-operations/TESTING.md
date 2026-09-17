# OneMind OS - Testing Guide

> Last Updated: 2026-01-21
> Status: Setting up automated testing

---

## Quick Start

```bash
# Run all tests
pytest tests/ -v

# Run specific test file
pytest tests/agents/legacy/test_config.py -v

# Run with coverage report
pytest tests/ --cov=agents --cov-report=html

# Run only fast tests (skip integration tests)
pytest tests/ -v -m "not integration"
```

---

## Test Status Dashboard

| Category | Tests | Passing | Failing | Notes |
|----------|-------|---------|---------|-------|
| **Config** | 9 | 9 | 0 | All passing |
| **Integrations API** | 4 | 4 | 0 | All passing |
| **Awareness** | 19 | 19 | 0 | All passing (fixed) |
| **Linear Tool** | 21 | 21 | 0 | All passing |
| **Todoist Tool** | 32 | 32 | 0 | All passing |
| **CORS** | 20 | - | - | Needs `openai` dep |
| **LiveKit** | 3 | - | - | Needs `livekit` dep |
| **UI Hooks** | 0 | - | - | Not yet written |

**Last Run:** 2026-01-21
**Total:** 86 tests passing (core tests), 23 skipped (missing deps)

---

## Test Categories

### 1. Unit Tests (Fast, No Dependencies)

Test individual functions in isolation.

```
tests/
├── agents/legacy/
│   ├── test_config.py       # ✅ Config loading
│   ├── test_awareness.py    # ⚠️ Awareness modes (needs fixes)
│   └── test_integrations_api.py  # ✅ Integration endpoints
```

**Example:**
```python
def test_config_imports():
    """Test that config module can be imported"""
    from agents.legacy.config import NATSConfig
    assert NATSConfig is not None
```

### 2. Integration Tests (Require Services)

Test components working together. Require database, NATS, etc.

```
tests/
├── test_cors.py           # Needs FastAPI + openai
├── test_livekit_connection.py  # Needs livekit
└── test_gemini_live.py    # Needs google-cloud
```

### 3. Tool Tests (Test API Integrations)

Test Legacy's tools against external APIs.

```
tests/
├── tools/
│   ├── test_linear.py     # Linear API (TODO)
│   ├── test_todoist.py    # Todoist API (TODO)
│   ├── test_github.py     # GitHub API (TODO)
│   └── test_comms.py      # Discord/Slack (TODO)
```

### 4. UI Hook Tests (TypeScript)

Test React hooks in isolation.

```
ui/src/hooks/__tests__/
├── useMemories.test.ts    # Memory CRUD (TODO)
├── useKnowledge.test.ts   # Knowledge CRUD (TODO)
├── useTeams.test.ts       # Teams fetch (TODO)
└── useAIStreamHandler.test.ts  # Streaming (TODO)
```

---

## Writing Tests

### Python Tests (pytest)

```python
# tests/tools/test_linear.py
import pytest
from unittest.mock import patch, MagicMock

def test_linear_list_issues():
    """Test listing Linear issues"""
    from agents.legacy.tools.linear import linear

    # Mock the HTTP request
    with patch('httpx.Client') as mock_client:
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "data": {
                "issues": {
                    "nodes": [
                        {"identifier": "OMOS-1", "title": "Test Issue"}
                    ]
                }
            }
        }
        mock_client.return_value.__enter__.return_value.post.return_value = mock_response

        result = linear(action="list", team="OMOS")

        assert "OMOS-1" in result
        assert "Test Issue" in result
```

### TypeScript Tests (Vitest)

```typescript
// ui/src/hooks/__tests__/useMemories.test.ts
import { renderHook, act } from '@testing-library/react'
import { useMemories } from '../useMemories'
import { vi } from 'vitest'

describe('useMemories', () => {
  it('fetches memories on mount', async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ memories: [] })
    })

    const { result } = renderHook(() => useMemories())

    await act(async () => {
      await result.current.fetchMemories()
    })

    expect(result.current.memories).toEqual([])
  })
})
```

---

## Test Fixtures

### Mocking External APIs

```python
# tests/conftest.py
import pytest
from unittest.mock import MagicMock

@pytest.fixture
def mock_linear_api():
    """Mock Linear GraphQL API"""
    return MagicMock(return_value={
        "data": {
            "issues": {"nodes": []}
        }
    })

@pytest.fixture
def mock_todoist_api():
    """Mock Todoist REST API"""
    return MagicMock(return_value={"tasks": []})
```

### Test Database

```python
@pytest.fixture
async def test_db():
    """Create a test database"""
    import asyncpg
    conn = await asyncpg.connect("postgresql://test:test@localhost/test")
    yield conn
    await conn.close()
```

---

## Automated Testing (CI/CD)

### GitHub Actions (Recommended)

Tests run automatically on every push/PR.

```yaml
# .github/workflows/test.yml
name: Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test
        ports:
          - 5432:5432

    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.11'

      - name: Install dependencies
        run: |
          pip install -r agents/legacy/requirements.txt
          pip install pytest pytest-asyncio pytest-cov

      - name: Run tests
        run: pytest tests/ -v --cov=agents --cov-report=xml

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: coverage.xml

  test-ui:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Node
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install dependencies
        run: cd ui && npm ci

      - name: Run tests
        run: cd ui && npm test

      - name: Type check
        run: cd ui && npm run type-check
```

### Coolify Integration

If using Coolify for CI/CD, add test step to build:

```dockerfile
# In Dockerfile
FROM python:3.11-slim AS test
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt pytest pytest-asyncio
COPY . .
RUN pytest tests/ -v --ignore=tests/test_livekit_connection.py

# Continue with build...
FROM python:3.11-slim AS production
# ...
```

---

## Running Tests Locally

### Prerequisites

```bash
# Install Python dependencies
pip install pytest pytest-asyncio pytest-cov httpx

# Install UI test dependencies
cd ui && npm install -D vitest @testing-library/react @testing-library/jest-dom
```

### Run All Tests

```bash
# Python
pytest tests/ -v

# UI (once set up)
cd ui && npm test
```

### Run with Coverage

```bash
# Python
pytest tests/ --cov=agents --cov-report=html
open htmlcov/index.html  # View report

# UI
cd ui && npm test -- --coverage
```

### Watch Mode (Development)

```bash
# Python
pytest tests/ -v --watch

# UI
cd ui && npm test -- --watch
```

---

## Test Markers

Use markers to categorize tests:

```python
# pytest.ini
[pytest]
markers =
    unit: Fast unit tests
    integration: Tests requiring external services
    slow: Slow tests (skip in CI)
    api: Tests for external API integrations
```

```python
# In test file
import pytest

@pytest.mark.unit
def test_config_loads():
    pass

@pytest.mark.integration
def test_database_connection():
    pass

@pytest.mark.api
def test_linear_api():
    pass
```

```bash
# Run only unit tests
pytest -m unit

# Run everything except slow
pytest -m "not slow"
```

---

## What Each Test Verifies

| Component | Test Verifies |
|-----------|---------------|
| `test_config.py` | Config classes load, env vars override defaults |
| `test_awareness.py` | Awareness modes exist, state serializes correctly |
| `test_integrations_api.py` | `/api/integrations/status` returns correct format |
| `test_cors.py` | CORS headers present on all endpoints |
| `test_linear.py` | Linear tool creates/lists/updates issues |
| `test_useMemories.ts` | Hook fetches, creates, deletes memories |

---

## Coverage Goals

| Module | Current | Target |
|--------|---------|--------|
| `agents/legacy/config.py` | 80% | 90% |
| `agents/legacy/tools/` | 0% | 70% |
| `agents/legacy/integrations/` | 30% | 60% |
| `ui/src/hooks/` | 0% | 70% |
| **Overall** | ~20% | 60% |

---

## Troubleshooting

### Missing Dependencies

```bash
# If tests fail with ImportError
pip install openai livekit google-cloud-aiplatform

# Or skip those tests
pytest tests/ --ignore=tests/test_livekit_connection.py
```

### Database Connection Errors

```bash
# Make sure Postgres is running
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=test postgres:15
```

### API Key Errors

```bash
# Set test API keys (or use mocks)
export LINEAR_API_KEY="test_key"
export TODOIST_API_KEY="test_key"
```

---

## Next Steps

1. **Add tool tests** - Linear, Todoist, GitHub
2. **Add UI hook tests** - useMemories, useKnowledge, useTeams
3. **Set up CI/CD** - GitHub Actions workflow
4. **Add coverage badges** - Show test status in README
5. **Pre-commit hooks** - Run tests before commit

---

## Resources

- [pytest documentation](https://docs.pytest.org/)
- [Vitest documentation](https://vitest.dev/)
- [Testing Library](https://testing-library.com/)
- [GitHub Actions](https://docs.github.com/en/actions)
