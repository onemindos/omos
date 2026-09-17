# Redis

> **In-Memory Cache & Session Store**

## Overview

Redis is a high-performance in-memory data store used for caching, sessions, rate limiting, and real-time features in OneMind OS.

## Status

- **Active**: ✅ Yes
- **Deployment**: Docker container
- **Last Updated**: January 25, 2026

## Use Cases

1. **Session Storage** - User sessions
2. **Cache Layer** - API response caching
3. **Rate Limiting** - Request throttling
4. **Pub/Sub** - Real-time messaging
5. **Queue** - Background jobs

## Key Features

- Sub-millisecond latency
- Multiple data structures
- Persistence options
- Pub/Sub messaging
- Lua scripting
- Cluster mode

## Installation

### Docker

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  -v redis-data:/data \
  redis:7-alpine \
  redis-server --appendonly yes
```

### Docker Compose

```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

volumes:
  redis-data:
```

## Python Integration

### Installation

```bash
pip install redis
```

### Basic Usage

```python
import redis

# Connect
r = redis.Redis(
    host='localhost',
    port=6379,
    db=0,
    decode_responses=True
)

# String operations
r.set('key', 'value')
r.get('key')  # 'value'
r.setex('temp_key', 3600, 'expires in 1 hour')

# Hash operations
r.hset('user:1', mapping={'name': 'Zeus', 'email': 'zeus@example.com'})
r.hgetall('user:1')  # {'name': 'Zeus', 'email': 'zeus@example.com'}

# List operations
r.rpush('queue', 'task1', 'task2')
r.lpop('queue')  # 'task1'

# Set operations
r.sadd('tags', 'python', 'ai', 'automation')
r.smembers('tags')  # {'python', 'ai', 'automation'}
```

### Async Redis

```python
import redis.asyncio as redis

async def main():
    r = redis.Redis(host='localhost', port=6379, decode_responses=True)
    await r.set('key', 'value')
    value = await r.get('key')
    await r.close()
```

## FastAPI Integration

### Session Store

```python
from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware
import redis
import json
import uuid

app = FastAPI()
redis_client = redis.Redis(host='localhost', port=6379, decode_responses=True)

async def get_session(session_id: str) -> dict:
    data = redis_client.get(f"session:{session_id}")
    return json.loads(data) if data else {}

async def save_session(session_id: str, data: dict, ttl: int = 86400):
    redis_client.setex(f"session:{session_id}", ttl, json.dumps(data))
```

### Rate Limiting

```python
from fastapi import HTTPException, Request
from fastapi.responses import JSONResponse

async def rate_limit(request: Request, limit: int = 100, window: int = 60):
    key = f"rate_limit:{request.client.host}"

    current = redis_client.incr(key)
    if current == 1:
        redis_client.expire(key, window)

    if current > limit:
        raise HTTPException(status_code=429, detail="Rate limit exceeded")
```

## Caching Pattern

### Cache Decorator

```python
import functools
import json

def cache(ttl: int = 300):
    def decorator(func):
        @functools.wraps(func)
        async def wrapper(*args, **kwargs):
            # Build cache key
            key = f"cache:{func.__name__}:{hash((args, tuple(kwargs.items())))}"

            # Check cache
            cached = redis_client.get(key)
            if cached:
                return json.loads(cached)

            # Execute function
            result = await func(*args, **kwargs)

            # Store in cache
            redis_client.setex(key, ttl, json.dumps(result))
            return result
        return wrapper
    return decorator

# Usage
@cache(ttl=3600)
async def get_user(user_id: str):
    # Expensive database query
    return await db.users.find_one({"id": user_id})
```

## Pub/Sub

### Publisher

```python
def publish_event(channel: str, message: dict):
    redis_client.publish(channel, json.dumps(message))

# Usage
publish_event('events:user', {'type': 'login', 'user_id': '123'})
```

### Subscriber

```python
def subscribe_to_events(channel: str, handler):
    pubsub = redis_client.pubsub()
    pubsub.subscribe(channel)

    for message in pubsub.listen():
        if message['type'] == 'message':
            data = json.loads(message['data'])
            handler(data)
```

## Data Structures

| Type           | Use Case               | Commands               |
| -------------- | ---------------------- | ---------------------- |
| **String**     | Simple key-value       | GET, SET, INCR         |
| **Hash**       | Objects, user data     | HGET, HSET, HGETALL    |
| **List**       | Queues, logs           | LPUSH, RPOP, LRANGE    |
| **Set**        | Tags, unique items     | SADD, SMEMBERS, SINTER |
| **Sorted Set** | Leaderboards, rankings | ZADD, ZRANGE, ZSCORE   |
| **Stream**     | Event logs             | XADD, XREAD, XRANGE    |

## Configuration

### redis.conf

```conf
# Persistence
appendonly yes
appendfsync everysec

# Memory
maxmemory 256mb
maxmemory-policy allkeys-lru

# Security
requirepass your_password

# Networking
bind 127.0.0.1
```

### Environment Variables

```env
REDIS_URL=redis://localhost:6379/0
REDIS_PASSWORD=your_password
```

## Monitoring

### Redis CLI

```bash
# Check info
redis-cli INFO

# Monitor commands
redis-cli MONITOR

# Memory usage
redis-cli INFO memory

# Slow log
redis-cli SLOWLOG GET 10
```

### Key Metrics

- `used_memory`: Current memory usage
- `connected_clients`: Active connections
- `total_commands_processed`: Command count
- `cache_hit_rate`: Hit/miss ratio
- `latency`: Command latency

## Best Practices

- ✅ Use appropriate data structures
- ✅ Set TTL on keys to prevent memory bloat
- ✅ Use pipelining for bulk operations
- ✅ Enable AOF persistence for durability
- ✅ Monitor memory usage
- ✅ Use connection pooling
- ✅ Avoid large keys (>1MB)

## Resources

- [Redis Documentation](https://redis.io/docs/)
- [Redis Python Client](https://redis-py.readthedocs.io/)
- [Redis Commands](https://redis.io/commands/)

---

**Related**: [PostgreSQL](./PostgreSQL.md) | [NATS](./NATS.md)
