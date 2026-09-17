# LiveKit Self-Hosted: The Honest Truth

> Created: 2026-01-23
> Purpose: 100% honest assessment of self-hosted LiveKit maintenance

---

## The "Super Hard" Reputation - Where It Comes From

### What People ACTUALLY Mean

**Hard Scenario 1: Complex Enterprise Deployments**
- Multi-region with 1000+ concurrent users
- Kubernetes clusters across 5 continents
- Custom TURN infrastructure
- Load balancing, auto-scaling
- **Reality**: This IS hard - not your use case

**Hard Scenario 2: Zero DevOps Experience**
- Never touched Docker
- Don't understand ports/firewalls
- No server management experience
- **Reality**: This would be hard - but you already run Docker + compose + services

**Hard Scenario 3: WebRTC First-Timers**
- Don't understand NAT/STUN/TURN
- Debug ICE connection failures blindly
- Blame LiveKit for network issues
- **Reality**: WebRTC has a learning curve - but LiveKit handles 90% automatically

---

## Your Specific Use Case: Single-Server, <100 Users

### What's Actually Hard (The Real Challenges)

#### 1. **NAT/Firewall Issues** (Most Common Problem)

**What goes wrong**:
```bash
# Client logs: "ICE connection failed"
# Server logs: No errors
# Problem: UDP ports blocked by firewall
```

**Why it's hard**:
- Corporate firewalls block UDP (home/office usually fine)
- Cloud providers (AWS, GCP) need security groups configured
- Testing requires EXTERNAL clients (can't test from same network)

**Your situation**:
- ✅ Already have server (100.102.21.44)
- ✅ Already managing firewall (NATS ports open)
- ✅ Tailscale mesh (simplifies NAT for most devices)
- ⚠️ **One-time setup**: Open UDP ports 50000-60000, 443

**Mitigation**:
```yaml
# Built-in TURN server handles NAT automatically
rtc:
  turn:
    enabled: true
    udp_port: 443  # Works through most firewalls
```

**Time to fix if broken**: 15-30 minutes (open ports, restart)
**Frequency**: Once during setup, rarely after

---

#### 2. **TURN Server Configuration** (Second Most Common)

**What goes wrong**:
```bash
# Clients behind restrictive networks can't connect
# LiveKit working, but some users can't join
```

**Why it's hard**:
- Historically required separate Coturn server
- TLS certificate setup for TURN
- Understanding TURN vs STUN

**Your situation**:
- ✅ LiveKit v1.5+ has **built-in TURN** (no Coturn needed!)
- ✅ Works on Tailscale IPs (no TLS cert needed)
- ⚠️ If exposing publicly: Need TLS cert (Let's Encrypt = free + automated)

**Mitigation**:
```yaml
# TURN works out-of-the-box with this config
rtc:
  turn:
    enabled: true
    udp_port: 443
  # That's it! No Coturn, no extra config
```

**Time to fix if broken**: 5-10 minutes (enable TURN, restart)
**Frequency**: Set once, forget forever

---

#### 3. **Redis Connection Issues** (Occasional)

**What goes wrong**:
```bash
# LiveKit can't connect to Redis
# Server logs: "Failed to connect to Redis"
# Symptom: Rooms don't persist across restarts
```

**Why it's hard**:
- Need separate Redis instance
- Network config between services

**Your situation**:
- ✅ Already running Redis (for Agno)
- ✅ Docker compose networking (trivial)

**Mitigation**:
```yaml
# In docker-compose
services:
  livekit:
    depends_on:
      - redis
  redis:
    image: redis:7-alpine
```

**Time to fix if broken**: 2 minutes (restart containers)
**Frequency**: Extremely rare (only if Redis crashes)

---

### What's Actually EASY (Surprisingly)

#### 1. **Installation** ✅
```bash
# That's it - one command
docker compose up -d livekit
```
**Time**: 30 seconds

---

#### 2. **Updates** ✅
```bash
# Update to new version
docker compose pull livekit
docker compose up -d livekit
```
**Time**: 2 minutes
**Breaking changes**: Rare (LiveKit has stable API)
**Rollback**: 30 seconds (use previous image tag)

---

#### 3. **Monitoring** ✅
```bash
# Check if running
docker logs onemind-livekit --tail 50

# Health check (built-in)
curl http://100.102.21.44:7880/
```
**Time**: 10 seconds
**Alerting**: Optional Prometheus metrics on port 6789

---

#### 4. **Scaling** ✅
- **Single server handles**: 50-100 concurrent participants
- **Your use case**: 1-5 concurrent users maximum
- **Headroom**: 10-20x capacity
- **Scaling needed**: Never (for your use case)

---

#### 5. **Debugging** ✅
```bash
# Enable debug logs
log_level: debug  # in config.yaml

# See exactly what's happening
docker logs onemind-livekit -f | grep ERROR
```
**Clarity**: LiveKit logs are actually readable
**Community**: Active Discord with maintainers

---

## Real Maintenance Schedule (Your Use Case)

### Daily: ZERO
- Runs indefinitely
- Auto-handles connections
- Self-healing (reconnects)

### Weekly: ZERO
- No manual intervention needed

### Monthly: 5 minutes
```bash
# Optional: Update to latest version
docker compose pull livekit
docker compose up -d livekit
```

### Quarterly: 10 minutes
- Review logs for patterns (optional)
- Check disk usage if recording (optional)

### Yearly: 30 minutes
- Major version upgrade (if any)
- Test from multiple devices

**Total yearly maintenance**: ~1-2 hours

---

## What ACTUALLY Breaks (Real Scenarios)

### Scenario 1: After Server Reboot
**Problem**: LiveKit doesn't start
**Cause**: `restart: unless-stopped` missing from compose
**Fix**: Add restart policy, `docker compose up -d`
**Frequency**: Once (during setup)
**Time to fix**: 1 minute

---

### Scenario 2: UDP Ports Blocked
**Problem**: Voice connects but drops immediately
**Cause**: Firewall changed, UDP ports blocked
**Fix**: Re-open UDP ports 50000-60000
**Frequency**: Rare (only if firewall changes)
**Time to fix**: 5 minutes

---

### Scenario 3: Redis Connection Lost
**Problem**: Rooms disappear after restart
**Cause**: Redis crashed or config changed
**Fix**: Restart Redis, check network config
**Frequency**: Very rare
**Time to fix**: 2 minutes

---

### Scenario 4: High CPU/Memory
**Problem**: Server slow, voice choppy
**Cause**: 50+ concurrent participants (unlikely for you)
**Fix**: Restart LiveKit, check for zombie sessions
**Frequency**: Extremely rare with <10 users
**Time to fix**: 2 minutes

---

### Scenario 5: Client Can't Connect
**Problem**: Some devices work, others don't
**Cause**: Corporate firewall blocking WebRTC
**Fix**: Nothing on server side - user network issue
**Note**: TURN helps but can't bypass all corporate firewalls
**Workaround**: User switches to non-corporate network

---

## Comparison with Managed Cloud

| Factor | Self-Hosted | LiveKit Cloud |
|--------|-------------|---------------|
| **Setup time** | 30 min (one-time) | 5 min (instant) |
| **Maintenance time** | ~1-2 hours/year | 0 hours/year |
| **Cost** | ~$0 (server exists) | ~$50-100/month |
| **NAT issues** | You debug | LiveKit debugs |
| **Scaling** | Manual (rarely needed) | Automatic |
| **Uptime SLA** | Your responsibility | 99.9% guaranteed |
| **Latency** | 10-50ms (same DC) | 50-150ms (edge) |
| **Egress limits** | Unlimited | Pay per minute |

**Trade-off**: ~2 hours/year maintenance → Save $600-1200/year

---

## The Honest Verdict

### Self-Hosted IS HARD If You:
- ❌ Never used Docker
- ❌ Don't understand networking (ports, firewalls)
- ❌ Need multi-region, high availability
- ❌ Have 1000+ concurrent users
- ❌ Don't have a dedicated server
- ❌ Can't afford any downtime

### Self-Hosted IS EASY If You:
- ✅ Already run Docker (you do - 7 containers running)
- ✅ Already manage services (postgres, redis, nats)
- ✅ Have a server with capacity (Core VM has room)
- ✅ Single location deployment (one datacenter)
- ✅ <100 concurrent users (you: <10 users)
- ✅ OK with 99% uptime vs 99.9% (acceptable for personal)

**Your situation**: ✅ All green checks - self-hosted makes sense

---

## What You DON'T Need to Worry About

### ❌ Complex TURN Setup
- LiveKit has built-in TURN (no Coturn)
- Works with UDP port 443
- No TLS certs needed on Tailscale IPs

### ❌ Load Balancing
- Single server handles 50-100 concurrent
- You: <10 concurrent users maximum
- Overkill for your scale

### ❌ Multi-Region Routing
- All your locations connect to one server
- Latency acceptable (50-150ms)
- No geo-routing needed

### ❌ Auto-Scaling
- Server resources fixed
- Capacity: 10-20x your needs
- No k8s, no horizontal scaling

### ❌ Advanced Monitoring
- Docker logs sufficient
- Optional: Prometheus metrics
- No need for Grafana dashboards (already have for other services)

---

## Decision Framework

### Go Self-Hosted If:
- Want unlimited usage (no per-minute costs)
- OK with ~2 hours/year maintenance
- Already comfortable with Docker
- Have server capacity

### Stay on Cloud If:
- Want absolute zero maintenance
- Usage <100 minutes/month (<$90/mo acceptable)
- Need guaranteed 99.9% uptime
- Don't want to debug NAT issues

---

## My Recommendation for OneMind OS

### ✅ GO SELF-HOSTED

**Reasons**:
1. You already run complex services (Postgres, NATS, Redis, Agno)
2. LiveKit is SIMPLER than what you already manage
3. Server has capacity (Azure VM can handle it)
4. Unlimited Gemini Live voice (no per-minute costs)
5. ~1-2 hours/year maintenance is acceptable
6. Save $600-1200/year

**Reality check**:
- Initial setup: 30 minutes
- First 3 months: Maybe 1-2 hours total (learning, tweaking)
- After that: 5 min/month (optional updates)

**Complexity comparison**:
- Agno + Teams + Workflows: 8/10 complexity
- NATS JetStream + MQTT: 7/10 complexity
- Frigate NVR + AI: 7/10 complexity
- **LiveKit self-hosted: 4/10 complexity** ← Easier than what you already do

---

## The "Super Hard" Myth - Debunked

**Why people say it's hard**:
1. They tried in 2020 (needed Coturn, complex setup)
2. They're running multi-region Kubernetes (enterprise scale)
3. They don't understand WebRTC (blame LiveKit for network issues)
4. They're comparing to zero-setup cloud services

**Reality for your use case**:
- LiveKit in 2024 is WAY easier (built-in TURN, better docs)
- Single server = simple (no k8s needed)
- Docker compose = 20 lines of YAML
- Maintenance = mostly "do nothing"

---

## Test Plan (Prove It's Easy)

### Week 1: Deploy and Test
```bash
# Day 1: Deploy (30 min)
docker compose up -d livekit

# Day 2-7: Use it normally
# - Voice conversations with Legacy
# - Multi-device testing
# - Monitor logs

# Result: Should "just work"
```

### Week 2-4: Burn-In Period
- Use voice daily
- Test from different networks
- Check logs for issues
- Measure actual maintenance time

### After 1 Month: Decision Point
**If maintenance > 30 min/month**: Switch to cloud
**If working well**: Keep self-hosted

---

## Emergency Rollback Plan

**If self-hosted fails badly**:
```bash
# 1. Switch back to cloud (5 min)
- Change 1 line in Agno: LiveKit Cloud URL
- Deploy

# 2. Or run both (hedging)
- Keep self-hosted for regular use
- Fall back to cloud if self-hosted down
- Feature flag in code
```

**Worst case**: You lose 2-3 hours of setup time
**Best case**: Save $600-1200/year forever

---

## Conclusion

**Is self-hosted LiveKit hard?**
- For enterprises: Yes (multi-region, k8s, scale)
- For beginners: Yes (no Docker/networking experience)
- **For you**: No (already manage more complex services)

**Maintenance burden?**
- Realistic: ~1-2 hours/year
- Most time: Reading logs (optional)
- Actual work: Mostly `docker compose restart`

**Worth it?**
- Save: $600-1200/year
- Cost: ~2 hours/year
- **ROI**: $300-600/hour of maintenance time 🎉

**My verdict**: The "super hard" reputation is outdated/overblown for your use case. Go self-hosted.
