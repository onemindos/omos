# Weekly VM Performance Review Checklist

**Purpose**: Review VM performance weekly to determine if upgrade to Standard_D4s_v3 is needed.

**Current Configuration**: Standard_D2s_v3 (2 vCPU, 8GB RAM, ~$76/month)
**Target Configuration**: Standard_D4s_v3 (4 vCPU, 16GB RAM, ~$146/month)

---

## Review Schedule

**Start Date**: 2026-01-27 (VM resized to Standard_D2s_v3)
**Review Frequency**: Every Monday at 9:00 AM
**Next Review**: 2026-02-03 (Week 1)

---

## Weekly Performance Checklist

### Week of: _______________

Run these commands and record the results:

#### 1. CPU Usage Check

**Command**:
```bash
ssh zeus@20.121.38.186 "top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\([0-9.]*\)%* id.*/\1/' | awk '{print 100 - \$1}'"
```

**Average CPU Usage**: _______%

**Azure Monitor (last 7 days)**:
```bash
az monitor metrics list \
  --resource /subscriptions/acba54ab-f643-4db6-9f77-57c97aac244f/resourceGroups/onemind-production/providers/Microsoft.Compute/virtualMachines/onemind-vm \
  --metric "Percentage CPU" \
  --start-time $(date -u -v-7d +%Y-%m-%dT%H:%M:%SZ) \
  --end-time $(date -u +%Y-%m-%dT%H:%M:%SZ) \
  --interval PT1H \
  --aggregation Average \
  --output table
```

**Status**:
- [ ] ✅ **UPGRADE NEEDED**: CPU consistently above 70%
- [ ] ⚠️ **MONITOR**: CPU between 50-70%
- [ ] ✓ **OK**: CPU below 50%

---

#### 2. Memory Usage Check

**Command**:
```bash
ssh zeus@20.121.38.186 "free -h && echo '---' && free | grep Mem | awk '{print \"Used: \" \$3 \" / Total: \" \$2 \" (\" int(\$3/\$2*100) \"%)\" }'"
```

**Memory Usage**: _______GB / 8GB (______%)

**Docker Container Memory**:
```bash
ssh zeus@20.121.38.186 "docker stats --no-stream --format 'table {{.Name}}\t{{.MemUsage}}'"
```

**Status**:
- [ ] ✅ **UPGRADE NEEDED**: Memory consistently above 6GB (75%)
- [ ] ⚠️ **MONITOR**: Memory between 4-6GB (50-75%)
- [ ] ✓ **OK**: Memory below 4GB (50%)

---

#### 3. Concurrent Agent Activity

**Check active agent runs**:
```bash
ssh zeus@20.121.38.186 "curl -s http://localhost:8080/agents | jq '.[] | {id: .agent_id, name: .name, active_runs: .stats.active_runs}'"
```

**Peak concurrent agents this week**: _______

**Average concurrent agents**: _______

**Status**:
- [ ] ✅ **UPGRADE NEEDED**: Running 5+ agents concurrently regularly
- [ ] ⚠️ **MONITOR**: Running 3-4 agents concurrently
- [ ] ✓ **OK**: Running 1-2 agents concurrently

---

#### 4. Workflow Performance

**Check recent run times** (from backend logs):
```bash
ssh zeus@20.121.38.186 "docker logs onemind-agno --tail=100 2>&1 | grep -i 'duration\|took\|completed' | tail -20"
```

**Average response time**: _______seconds

**Slowest workflows**:
1. _______________________
2. _______________________
3. _______________________

**Status**:
- [ ] ✅ **UPGRADE NEEDED**: Workflows taking 2x longer than expected, frequent timeouts
- [ ] ⚠️ **MONITOR**: Occasional slow responses (>30s for complex queries)
- [ ] ✓ **OK**: Workflows completing in expected timeframe

---

#### 5. System Health

**Check for errors/warnings**:
```bash
ssh zeus@20.121.38.186 "docker logs onemind-agno --tail=200 2>&1 | grep -iE 'error|warning|timeout|failed' | wc -l"
```

**Error count (last 200 lines)**: _______

**Database connections**:
```bash
ssh zeus@20.121.38.186 "docker exec onemind-timescaledb psql -U onemind -d onemind -c 'SELECT count(*) as active_connections FROM pg_stat_activity;'"
```

**Active connections**: _______

**Status**:
- [ ] ✅ **UPGRADE NEEDED**: Frequent errors, connection limits reached
- [ ] ⚠️ **MONITOR**: Occasional errors
- [ ] ✓ **OK**: Minimal errors, stable

---

## Upgrade Decision Matrix

**Upgrade to Standard_D4s_v3 if ANY of these are true**:

- [ ] ✅ CPU consistently above 70%
- [ ] ✅ Memory above 6GB consistently
- [ ] ✅ Running 5+ agents concurrently regularly
- [ ] ✅ Workflows taking too long or timing out
- [ ] ✅ Ready to invest in concurrent operations (~$146/month)

**Current Score**: _____ / 5 indicators met

---

## Upgrade Command (When Ready)

```bash
cd /Users/zeuslegacy/Documents/onemind-os/infra/deployment/azure
./resize-vm.sh Standard_D4s_v3
```

**Estimated downtime**: 2-3 minutes
**Cost increase**: +$70/month ($76 → $146)

---

## Historical Performance Log

### Week 1: 2026-02-03

**CPU**: _____%
**Memory**: _____GB
**Concurrent Agents**: _____
**Decision**: [ ] Stay on D2s_v3 / [ ] Upgrade to D4s_v3
**Notes**: _______________________________________________

---

### Week 2: 2026-02-10

**CPU**: _____%
**Memory**: _____GB
**Concurrent Agents**: _____
**Decision**: [ ] Stay on D2s_v3 / [ ] Upgrade to D4s_v3
**Notes**: _______________________________________________

---

### Week 3: 2026-02-17

**CPU**: _____%
**Memory**: _____GB
**Concurrent Agents**: _____
**Decision**: [ ] Stay on D2s_v3 / [ ] Upgrade to D4s_v3
**Notes**: _______________________________________________

---

### Week 4: 2026-02-24

**CPU**: _____%
**Memory**: _____GB
**Concurrent Agents**: _____
**Decision**: [ ] Stay on D2s_v3 / [ ] Upgrade to D4s_v3
**Notes**: _______________________________________________

---

## Quick Monitoring Script

Save this as `~/check-vm-performance.sh` on your local machine:

```bash
#!/bin/bash

echo "========================================="
echo "OneMind OS VM Performance Check"
echo "========================================="
echo ""

echo "1. CPU Usage:"
ssh zeus@20.121.38.186 "top -bn1 | grep 'Cpu(s)'"
echo ""

echo "2. Memory Usage:"
ssh zeus@20.121.38.186 "free -h | grep -E 'Mem|Swap'"
echo ""

echo "3. Docker Containers:"
ssh zeus@20.121.38.186 "docker stats --no-stream"
echo ""

echo "4. Recent Backend Errors:"
ssh zeus@20.121.38.186 "docker logs onemind-agno --tail=50 2>&1 | grep -i error | tail -10"
echo ""

echo "========================================="
echo "Review complete - update WEEKLY_VM_REVIEW.md"
echo "========================================="
```

**Make executable**: `chmod +x ~/check-vm-performance.sh`
**Run weekly**: `~/check-vm-performance.sh`

---

## Contact & Support

**Maintained By**: Zeus
**VM IP**: 20.121.38.186
**Resource Group**: onemind-production
**VM Name**: onemind-vm

**Related Documentation**:
- [INFRASTRUCTURE.md](./INFRASTRUCTURE.md)
- [SCALING.md](./SCALING.md)
- [resize-vm.sh](./deployment/azure/resize-vm.sh)

---

**Last Updated**: 2026-01-27
**Current VM Size**: Standard_D2s_v3 (2 vCPU, 8GB RAM)
