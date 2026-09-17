# OneMind OS Scaling Guide

Guide for scaling OneMind OS infrastructure from current state to heavy concurrent workflows.

---

## Current State

**VM**: Standard_B2s (2 vCPU, 4GB RAM)
- Status: ❌ **Inadequate for production workflows**
- Issue: Burstable CPU will throttle under sustained load
- Cost: ~$30/month

---

## Scaling Path

### Phase 1: Production Readiness (NOW)

**Target**: Standard_D2s_v3
- **vCPU**: 2 (Intel Xeon Platinum 8272CL)
- **RAM**: 8GB DDR4
- **Cost**: ~$76/month (+$46/month from current)
- **Performance**: Sustained CPU, no throttling
- **Workload**: Moderate workflows, 2-3 concurrent agents

**Why This First**:
- Minimum viable size for production workflows
- Eliminates CPU throttling from B-series
- 2x memory (4GB → 8GB)
- Affordable baseline for current needs

**Execute**:
```bash
cd /Users/zeuslegacy/Documents/onemind-os/infra/deployment/azure
./resize-vm.sh Standard_D2s_v3
```

**Downtime**: 2-3 minutes (Docker containers restart automatically)

---

### Phase 2: Heavy Concurrent Workloads (GOAL)

**Target**: Standard_D4s_v3
- **vCPU**: 4 (2x compute power)
- **RAM**: 16GB (2x memory)
- **Cost**: ~$146/month (+$70/month from Phase 1)
- **Performance**: High sustained CPU, excellent concurrency
- **Workload**: Many agents running simultaneously, complex workflows

**Why This Eventually**:
- 2x more CPU cores for parallel agent execution
- 2x more memory for larger context windows and caching
- Handles 5-10+ concurrent agent runs
- Supports complex multi-agent teams
- No performance bottlenecks

**When to Upgrade**:
Monitor these metrics on Standard_D2s_v3:

1. **CPU Usage** (via `htop` or Azure Monitor):
   - Consistently above 70% → Need more CPU
   - Spikes to 100% during workflows → Need D4s_v3

2. **Memory Usage** (via `free -h`):
   - Above 6GB consistently → Need more RAM
   - OOM (Out of Memory) errors → Immediate upgrade needed

3. **Workflow Performance**:
   - Agents taking longer than expected
   - Timeouts or failures under load
   - Queue building up

4. **Concurrent Agent Count**:
   - Running 5+ agents regularly → Consider D4s_v3
   - Running 10+ agents → Definitely need D4s_v3

**Execute When Ready**:
```bash
cd /Users/zeuslegacy/Documents/onemind-os/infra/deployment/azure
./resize-vm.sh Standard_D4s_v3
```

---

## Monitoring Performance

### Check CPU Usage

**Via SSH**:
```bash
ssh zeus@20.121.38.186 "htop"
# Press F10 to quit
```

**Via Azure CLI**:
```bash
# Average CPU over last hour
az monitor metrics list --resource /subscriptions/acba54ab-f643-4db6-9f77-57c97aac244f/resourceGroups/onemind-production/providers/Microsoft.Compute/virtualMachines/onemind-vm \
  --metric "Percentage CPU" \
  --interval PT1H \
  --output table
```

### Check Memory Usage

**Via SSH**:
```bash
ssh zeus@20.121.38.186 "free -h"
```

**Check Docker Container Memory**:
```bash
ssh zeus@20.121.38.186 "docker stats --no-stream"
```

### Check Agent Performance

**Backend Logs** (look for slow responses):
```bash
ssh zeus@20.121.38.186 "cd ~/OneMind-OS && docker compose logs --tail=100 agno | grep -i 'took\|duration\|timeout'"
```

**Database Connections**:
```bash
ssh zeus@20.121.38.186 "cd ~/OneMind-OS && docker compose exec -T postgres psql -U onemind -d onemind -c 'SELECT count(*) FROM pg_stat_activity;'"
```

---

## Cost Comparison

| Phase | VM Size | vCPU | RAM | Cost/Month | vs Current | Use Case |
|-------|---------|------|-----|------------|------------|----------|
| Current | Standard_B2s | 2 | 4GB | $30 | - | ❌ Development only |
| **Phase 1** | **Standard_D2s_v3** | 2 | 8GB | **$76** | **+$46** | ✅ **Production workflows** |
| **Phase 2** | **Standard_D4s_v3** | 4 | 16GB | **$146** | **+$116** | ✅ **Heavy concurrent** |

**Cost Savings with Reserved Instances**:

| Commitment | Standard_D2s_v3 | Standard_D4s_v3 | Savings |
|------------|-----------------|-----------------|---------|
| Pay-as-you-go | $76/month | $146/month | 0% |
| 1-year RI | $53/month | $102/month | 30% |
| 3-year RI | $38/month | $73/month | 50% |

**Recommendation**: Use Pay-as-you-go for Phase 1 (testing period), then commit to 1-year RI in Phase 2 when stable.

---

## Performance Benchmarks

### Standard_D2s_v3 (Phase 1) - Expected Performance

**Single Agent Run**:
- Simple query: 2-5 seconds
- Complex reasoning: 10-30 seconds
- With tools: 15-60 seconds

**Concurrent Agents**:
- 2 agents: Smooth performance
- 3 agents: Good performance
- 4 agents: Possible slowdown
- 5+ agents: Consider Phase 2

**Workflow Examples**:
- Daily task sync (Todoist): < 10 seconds
- Knowledge base search: < 5 seconds
- Multi-step workflow (3 agents): 30-90 seconds

### Standard_D4s_v3 (Phase 2) - Expected Performance

**Single Agent Run**:
- Simple query: 1-3 seconds (faster)
- Complex reasoning: 5-20 seconds (faster)
- With tools: 10-45 seconds (faster)

**Concurrent Agents**:
- 2-5 agents: Excellent performance
- 6-10 agents: Good performance
- 11-15 agents: Acceptable performance
- 16+ agents: Consider horizontal scaling (multiple VMs)

**Workflow Examples**:
- Complex multi-agent team (5 agents): 60-120 seconds
- Parallel workflow execution (10 agents): 90-180 seconds
- Heavy data processing: No bottlenecks

---

## Alternative Scaling Strategies

### Option 1: Vertical Scaling (Recommended)

**What**: Resize VM to larger size (D2s_v3 → D4s_v3)

**Pros**:
- Simple (one VM to manage)
- No code changes
- Fast to implement (2-3 min downtime)
- Lower operational complexity

**Cons**:
- Limited to single VM capacity
- Downtime during resize
- Single point of failure

**Best For**: Current needs (up to ~15 concurrent agents)

### Option 2: Horizontal Scaling

**What**: Add more VMs with load balancer

**Pros**:
- No upper limit on capacity
- No downtime during scaling
- High availability (redundancy)

**Cons**:
- More expensive (load balancer + multiple VMs)
- Complex setup (shared database, session management)
- Requires code changes (stateless agents)

**Best For**: Enterprise scale (50+ concurrent agents)

**Not Needed Yet** - Vertical scaling to D4s_v3 will handle your needs for the foreseeable future.

---

## Scaling Checklist

### Before Phase 1 (D2s_v3)

- [ ] Verify current VM is Standard_B2s
- [ ] Backup database (automatic in resize script)
- [ ] Test resize script with `--dry-run`
- [ ] Schedule resize during low-usage time
- [ ] Notify users of 2-3 min downtime

### After Phase 1 (D2s_v3)

- [ ] Verify all services restarted (`docker compose ps`)
- [ ] Check backend health (`curl http://20.121.38.186:8080/health`)
- [ ] Test agent runs (send test message)
- [ ] Monitor CPU/memory for 24 hours
- [ ] Verify no performance issues

### Before Phase 2 (D4s_v3)

- [ ] Monitor D2s_v3 for 2+ weeks
- [ ] CPU consistently above 70%? → Upgrade needed
- [ ] Memory above 6GB? → Upgrade needed
- [ ] Ready to invest $146/month? → Upgrade when ready
- [ ] Consider 1-year reserved instance for savings

### After Phase 2 (D4s_v3)

- [ ] Same verification as Phase 1
- [ ] Test heavy concurrent workloads (5-10 agents)
- [ ] Measure performance improvement
- [ ] Monitor for 1 week, then consider RI commitment
- [ ] Document performance benchmarks

---

## Rollback Plan

If Phase 1 or Phase 2 causes issues:

**Immediate Rollback**:
```bash
# Resize back to previous size
./resize-vm.sh Standard_B2s  # Roll back from Phase 1
./resize-vm.sh Standard_D2s_v3  # Roll back from Phase 2
```

**Common Issues**:

1. **Services don't restart**:
```bash
ssh zeus@20.121.38.186 "cd ~/OneMind-OS && docker compose up -d"
```

2. **Database connection errors**:
```bash
ssh zeus@20.121.38.186 "cd ~/OneMind-OS && docker compose restart postgres"
sleep 10
ssh zeus@20.121.38.186 "cd ~/OneMind-OS && docker compose restart agno"
```

3. **High costs**:
- Downgrade to smaller size
- Stop VM when not in use (saves compute, keeps disk)
- Consider burstable B2ms for non-production

---

## FAQ

**Q: Why not go straight to Standard_D4s_v3?**
A: Phase 1 (D2s_v3) at $76/month provides 80% of the performance for 52% of the cost. Start there, validate workflows, then upgrade when heavy concurrent workloads are actually needed. Avoid over-provisioning.

**Q: Can I use Standard_B2ms (8GB burstable) instead?**
A: Not recommended for production workflows. Burstable VMs use CPU credits that exhaust under sustained load, causing throttling. D-series has sustained performance always available.

**Q: What if I need more than Standard_D4s_v3?**
A: Next steps:
1. **Standard_D8s_v3** (8 vCPU, 32GB, ~$290/month) - Double capacity
2. **Horizontal scaling** (multiple VMs + load balancer) - Unlimited capacity
3. **Azure Kubernetes** (AKS) - Enterprise-grade orchestration

**Q: How long should I stay on each phase?**
A:
- Phase 1 (D2s_v3): Minimum 2-4 weeks to establish performance baseline
- Phase 2 (D4s_v3): When monitoring shows consistent need (CPU > 70%, 5+ concurrent agents)

**Q: Can I use Spot Instances to save money?**
A: Not for production workflows. Spot instances can be evicted with 30 seconds notice. Only use for non-critical dev/test environments.

---

## Summary

**Immediate Action** (Phase 1):
```bash
cd /Users/zeuslegacy/Documents/onemind-os/infra/deployment/azure
./resize-vm.sh Standard_D2s_v3
```

**Future Goal** (Phase 2):
- Monitor D2s_v3 performance
- Upgrade to D4s_v3 when ready for heavy concurrent workloads
- Cost: $146/month (or $102/month with 1-year RI)

**Current Path**: B2s → D2s_v3 → D4s_v3 (when ready)

---

**Last Updated**: 2025-01-27
**Document Version**: 1.0
