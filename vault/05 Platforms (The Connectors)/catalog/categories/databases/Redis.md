# Redis

> In-memory data store for caching, sessions, rate limiting, and pub/sub.

---

## Overview

| Attribute         | Value                                  |
| ----------------- | -------------------------------------- |
| **Type**          | In-Memory Data Store                   |
| **Protocol**      | RESP (Redis Serialization Protocol)    |
| **Default Port**  | 6379                                   |
| **License**       | BSD-3-Clause                           |
| **Documentation** | [redis.io/docs](https://redis.io/docs) |

---

## What is Redis?

Redis is an open-source, in-memory data structure store used as:

- **Cache** — Fast data access layer
- **Session store** — User session management
- **Message broker** — Pub/sub messaging
- **Rate limiter** — API throttling
- **Queue** — Job/task queues
- **Leaderboard** — Sorted sets for rankings

---

## OneMind Usage

### Primary Use Cases

```
┌─────────────────────────────────────────────────────────────────┐
│                         OneMind Redis                            │
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │   Caching    │    │   Sessions   │    │ Rate Limits  │       │
│  │  LLM resp.   │    │  User state  │    │  API calls   │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐       │
│  │   Pub/Sub    │    │   Queues     │    │  Temp Data   │       │
│  │  Real-time   │    │  Background  │    │  TTL-based   │       │
│  └──────────────┘    └──────────────┘    └──────────────┘       │
└─────────────────────────────────────────────────────────────────┘
```

| Use Case                 | Description                                |
| ------------------------ | ------------------------------------------ |
| **LLM Response Caching** | Cache repeated queries to reduce API costs |
| **Session State**        | Store user/agent session data              |
| **Rate Limiting**        | Throttle API calls per user/endpoint       |
| **Real-time Pub/Sub**    | Lightweight event broadcasting             |
| **Background Queues**    | Task queue for async jobs                  |
| **Temporary Storage**    | Short-lived data with TTL                  |

---

## Docker Compose

```yaml
services:
  redis:
    image: redis:7-alpine
    container_name: onemind-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - onemind-network

volumes:
  redis_data:
```

---

## Configuration

### redis.conf (Production)

```conf
# Memory management
maxmemory 512mb
maxmemory-policy allkeys-lru

# Persistence
appendonly yes
appendfsync everysec

# Security
requirepass your_strong_password_here

# Network
bind 0.0.0.0
protected-mode yes

# Performance
tcp-keepalive 300
timeout 0
```

### Environment Variables

```bash
# Connection
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=your_password

# With auth
REDIS_URL=redis://:your_password@localhost:6379/0
```

---

## Python Integration

### Basic Connection

```python
import redis

# Simple connection
r = redis.Redis(host='localhost', port=6379, db=0)

# From URL
r = redis.from_url("redis://localhost:6379/0")

# With password
r = redis.Redis(
    host='localhost',
    port=6379,
    password='your_password',
    decode_responses=True  # Return strings instead of bytes
)
```

### Async Connection

```python
import redis.asyncio as aioredis

async def get_redis():
    return await aioredis.from_url(
        "redis://localhost:6379",
        encoding="utf-8",
        decode_responses=True
    )

async def main():
    redis = await get_redis()
    await redis.set("key", "value")
    value = await redis.get("key")
```

---

## Data Structures

### Strings

```python
# Basic key-value
r.set("user:1:name", "Zeus")
name = r.get("user:1:name")

# With expiration
r.setex("session:abc123", 3600, "session_data")  # 1 hour TTL

# Increment/Decrement
r.incr("api:calls:today")
r.incrby("points:user:1", 10)
```

### Hashes

```python
# Store object-like data
r.hset("user:1", mapping={
    "name": "Zeus",
    "email": "zeus@onemind.dev",
    "awareness_level": "engaged"
})

# Get single field
name = r.hget("user:1", "name")

# Get all fields
user = r.hgetall("user:1")
```

### Lists

```python
# Queue (FIFO)
r.lpush("tasks:pending", "task_1")
r.lpush("tasks:pending", "task_2")
task = r.rpop("tasks:pending")  # Returns "task_1"

# Recent items
r.lpush("recent:commands", "lights on")
r.ltrim("recent:commands", 0, 99)  # Keep only last 100
```

### Sets

```python
# Unique values
r.sadd("online:users", "user:1", "user:2")
r.sismember("online:users", "user:1")  # True
r.srem("online:users", "user:1")

# Set operations
r.sinter("tags:post:1", "tags:post:2")  # Intersection
```

### Sorted Sets

```python
# Leaderboard
r.zadd("scores", {"user:1": 100, "user:2": 85})
r.zincrby("scores", 10, "user:1")  # Add 10 to user:1

# Top 10
top = r.zrevrange("scores", 0, 9, withscores=True)
```

---

## OneMind Patterns

### LLM Response Caching

```python
import hashlib
import json

class LLMCache:
    def __init__(self, redis_client, ttl=3600):
        self.redis = redis_client
        self.ttl = ttl

    def _cache_key(self, prompt: str, model: str) -> str:
        content = f"{model}:{prompt}"
        return f"llm:cache:{hashlib.sha256(content.encode()).hexdigest()}"

    async def get(self, prompt: str, model: str) -> str | None:
        key = self._cache_key(prompt, model)
        cached = await self.redis.get(key)
        if cached:
            return json.loads(cached)
        return None

    async def set(self, prompt: str, model: str, response: str):
        key = self._cache_key(prompt, model)
        await self.redis.setex(key, self.ttl, json.dumps(response))

# Usage
cache = LLMCache(redis_client)

# Check cache first
cached = await cache.get(prompt, model)
if cached:
    return cached

# Call LLM
response = await agent.run(prompt)

# Cache response
await cache.set(prompt, model, response)
```

### Rate Limiting

```python
class RateLimiter:
    def __init__(self, redis_client):
        self.redis = redis_client

    async def is_allowed(
        self,
        key: str,
        limit: int = 100,
        window: int = 60
    ) -> bool:
        """Sliding window rate limiter."""
        current = await self.redis.incr(f"ratelimit:{key}")

        if current == 1:
            await self.redis.expire(f"ratelimit:{key}", window)

        return current <= limit

    async def get_remaining(self, key: str, limit: int = 100) -> int:
        current = await self.redis.get(f"ratelimit:{key}")
        current = int(current) if current else 0
        return max(0, limit - current)

# Usage in FastAPI
limiter = RateLimiter(redis_client)

@app.middleware("http")
async def rate_limit_middleware(request: Request, call_next):
    client_ip = request.client.host

    if not await limiter.is_allowed(f"ip:{client_ip}"):
        raise HTTPException(429, "Rate limit exceeded")

    return await call_next(request)
```

### Session Management

```python
import uuid
from datetime import timedelta

class SessionManager:
    def __init__(self, redis_client, ttl: int = 86400):
        self.redis = redis_client
        self.ttl = ttl

    async def create_session(self, user_id: str, data: dict) -> str:
        session_id = str(uuid.uuid4())
        session_key = f"session:{session_id}"

        await self.redis.hset(session_key, mapping={
            "user_id": user_id,
            **data
        })
        await self.redis.expire(session_key, self.ttl)

        return session_id

    async def get_session(self, session_id: str) -> dict | None:
        session_key = f"session:{session_id}"
        data = await self.redis.hgetall(session_key)

        if data:
            # Refresh TTL on access
            await self.redis.expire(session_key, self.ttl)

        return data if data else None

    async def destroy_session(self, session_id: str):
        await self.redis.delete(f"session:{session_id}")
```

### Pub/Sub for Real-time

```python
import asyncio

class EventBroadcaster:
    def __init__(self, redis_client):
        self.redis = redis_client
        self.pubsub = self.redis.pubsub()

    async def publish(self, channel: str, message: dict):
        await self.redis.publish(channel, json.dumps(message))

    async def subscribe(self, channel: str, callback):
        await self.pubsub.subscribe(channel)

        async for message in self.pubsub.listen():
            if message["type"] == "message":
                data = json.loads(message["data"])
                await callback(data)

# Usage
broadcaster = EventBroadcaster(redis_client)

# Publisher
await broadcaster.publish("awareness:changes", {
    "level": "engaged",
    "timestamp": datetime.now().isoformat()
})

# Subscriber
async def on_awareness_change(data):
    print(f"Awareness changed to {data['level']}")

asyncio.create_task(
    broadcaster.subscribe("awareness:changes", on_awareness_change)
)
```

### Background Task Queue

```python
import json
from datetime import datetime

class TaskQueue:
    def __init__(self, redis_client):
        self.redis = redis_client

    async def enqueue(self, queue: str, task: dict):
        task["enqueued_at"] = datetime.now().isoformat()
        await self.redis.lpush(f"queue:{queue}", json.dumps(task))

    async def dequeue(self, queue: str, timeout: int = 0):
        result = await self.redis.brpop(f"queue:{queue}", timeout)
        if result:
            return json.loads(result[1])
        return None

    async def length(self, queue: str) -> int:
        return await self.redis.llen(f"queue:{queue}")

# Producer
queue = TaskQueue(redis_client)
await queue.enqueue("notifications", {
    "type": "email",
    "to": "zeus@example.com",
    "subject": "Daily Summary"
})

# Worker
async def worker():
    while True:
        task = await queue.dequeue("notifications", timeout=5)
        if task:
            await process_notification(task)
```

---

## Monitoring

### Redis CLI

```bash
# Connect
redis-cli

# Check memory
INFO memory

# List keys
KEYS *

# Monitor commands
MONITOR

# Check slow log
SLOWLOG GET 10
```

### Key Metrics

```bash
# Memory usage
redis-cli INFO memory | grep used_memory_human

# Connected clients
redis-cli INFO clients | grep connected_clients

# Operations per second
redis-cli INFO stats | grep instantaneous_ops_per_sec

# Hit rate
redis-cli INFO stats | grep keyspace_hits
redis-cli INFO stats | grep keyspace_misses
```

### Docker Health Check

```yaml
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 10s
  timeout: 5s
  retries: 5
```

---

## Best Practices

### Key Naming

```python
# ✅ Good: Hierarchical, descriptive
"user:1:session"
"cache:llm:sha256hash"
"ratelimit:ip:192.168.1.1"
"queue:notifications"

# ❌ Bad: Flat, ambiguous
"session1"
"cache"
"limit"
```

### TTL Strategy

| Data Type  | TTL              |
| ---------- | ---------------- |
| LLM cache  | 1-24 hours       |
| Session    | 24 hours         |
| Rate limit | 1 minute         |
| Temp data  | 5-15 minutes     |
| Queues     | No TTL (consume) |

### Memory Management

```conf
# Set max memory
maxmemory 512mb

# Eviction policy
maxmemory-policy allkeys-lru  # Evict least recently used
# Options: volatile-lru, allkeys-random, volatile-ttl, noeviction
```

---

## Troubleshooting

### Connection Issues

```bash
# Test connectivity
redis-cli -h localhost -p 6379 ping

# Check if password required
redis-cli -h localhost -p 6379 -a your_password ping

# Check binding
netstat -tlnp | grep 6379
```

### Memory Issues

```bash
# Check memory
redis-cli INFO memory

# Find big keys
redis-cli --bigkeys

# Clear all data (careful!)
redis-cli FLUSHALL
```

### Slow Queries

```bash
# Enable slow log
redis-cli CONFIG SET slowlog-log-slower-than 10000  # 10ms

# View slow queries
redis-cli SLOWLOG GET 25
```

---

## Resources

- [Redis Documentation](https://redis.io/docs)
- [Redis Python Client](https://redis-py.readthedocs.io/)
- [Redis Best Practices](https://redis.io/docs/management/optimization/)

---

_Last updated: January 2026_


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
