# Legacy Mother Self-Calibration Framework — Fabric State Authority & Verification

**Official Skill:** https://github.com/nousresearch/hermes-agent/skills/autonomous-ai-agents/legacy-self-calibration

**Hermes Skill:** `legacy-self-calibration` (production-ready, autonomous)

**Creator:** Zeus (OneMind OS)

**Date:** 2026-07-11

**PRIMARY PURPOSE:** Legacy Mother verifies and maintains fabric-wide authority

**SECONDARY PURPOSE:** Edge nodes receive mother's verified state broadcasts

---

## Problem: Mother Must Know What's Happening (Edge Nodes Can't Decide)

### The Challenge

**Legacy Mother is the authoritative primary.** But she can drift out of sync with reality:

**Scenario 1: Mother Doesn't Know Her Own Health**
- Mother thinks DOKS K8s is healthy (last checked 2 hours ago)
- K8s cluster actually crashed 30 minutes ago
- Mother keeps delegating tasks to dead infrastructure
- Result: Fabric-wide failures, cascade degradation

**Scenario 2: Mother Doesn't Detect Partial Outage**
- Some regions responsive, some not
- Mother doesn't realize fabric is split
- Broadcasts wrong decisions to regions that can't reach her
- Result: Incoherent behavior, edge nodes confused

**Scenario 3: Mother Doesn't Know Fabric State**
- "Where are all my edge nodes?"
- "Which are healthy, which are struggling?"
- "How bad is cross-region latency?"
- Result: Can't make informed fabric-wide decisions

**The Fix: Mother Continuously Self-Verifies**

**Idea:** Mother performs active health-check at strategic intervals:
- Daily (scheduled, know she's healthy)
- Before critical decisions (verify authority before commanding)
- Continuously in background (detect failures immediately)

**Result:** Mother always has fabric ground truth. Edge nodes trust her decisions. No confusion about authority.

---

## Framework: RECALIBRATE

Five-point verification system from Zeus's original notes:

```
recalibrate:
  modes       (Operating posture)
  state       (Ground truth about reality)
  memory      (Persistent facts current?)
  assets      (Resources available?)
  speed       (Performance baseline?)
```

### 1. MODES — Operating Posture

**Question:** "How should I operate right now?"

**Four modes Legacy can be in:**

#### NORMAL MODE
```
✅ Cloud mother fully available
✅ Edge nodes connected to fabric
✅ DOKS K8s cluster healthy
✅ All microservices running
✅ Full automation enabled

Behavior: Operate normally, delegate freely, trust cloud
Messaging latency acceptable: <500ms to mother
```

#### DEGRADED MODE
```
⚠️  Mother intermittent (some requests fail)
⚠️  K8s cluster reduced capacity (some nodes down)
⚠️  Edge nodes can operate locally
⚠️  Some services failing
⚠️  Reduce automation scope

Behavior: Reduce delegation, cache more locally, expect failures
Messaging latency degraded: 500ms-2s to mother
```

#### OFFLINE-FIRST MODE
```
🔴 Mother unreachable (cloud down)
🔴 NATS cloud broker unreachable
🔴 Edge mesh fully autonomous
🔴 Local NATS broker OK
🔴 TAK Server local-only

Behavior: Stop waiting for mother, operate fully autonomous
Messaging: Local only, no cloud
Queue all state updates for sync when mother returns
```

#### MAINTENANCE MODE
```
🔧 Scheduled updates in progress
🔧 Limited automation enabled
🔧 Manual oversight required
🔧 Restoration timeline published

Behavior: Human-supervised only, reduced autonomy
```

**Calibration Check:**
```bash
# Test reachability
nats micro ls --server nats://localhost:4222   # Local NATS?
ping agents.prompt.hermes.zeus.legacy          # Mother reachable?
kubectl get nodes                               # K8s healthy?
nats rtt --server nats://localhost:4222        # Latency?

# Determine mode
if mother_responds && latency < 200ms && k8s_healthy:
  MODE = NORMAL
elif mother_intermittent || latency > 2s || k8s_degraded:
  MODE = DEGRADED
elif mother_unreachable && local_nats_alive:
  MODE = OFFLINE_FIRST
else:
  MODE = MAINTENANCE_REQUIRED
```

### 2. STATE — Ground Truth About Reality

**Question:** "What's actually running right now?"

**Infrastructure State:**
- NATS broker: alive? version? how many subscribers?
- Kubernetes cluster: how many nodes up? pods running? PVCs mounted?
- ClickHouse: database responsive? data retention policies active?
- TAK Server: process running? clients connected?
- Tailscale mesh: how many nodes? latency between them?
- DNS/networking: external IP stable? resolution working?

**Agent State:**
- Mother (`agents.prompt.hermes.zeus.legacy`): responding? CPU load? uptime?
- Edge (`agents.prompt.hermes.zeus.legacy-edge`): responding? local-only or cloud-connected?
- Claude Code agent: available? ready for delegation?
- Codex agent: available? ready for delegation?
- Task engine (`om_task`): how many jobs queued? any stuck?
- Presence bridge (`om_presence`): voice bridge active?

**Deployment State:**
- Which Convergence files deployed to where?
- Which skills loaded in this session?
- Which cron jobs running?
- Which background processes active?
- Any orphaned/zombie processes from previous sessions?

**External Dependencies:**
- GitHub: accessible? API rate limit status?
- HuggingFace: model repositories reachable?
- ArXiv: search API responsive?
- Vercel Eve: API endpoint reachable?
- LM Studio: models loaded? memory available?

**Calibration Check:**
```bash
# Poll every critical system (3 second timeout)
systems = ["nats", "k8s", "clickhouse", "tak", "tailscale"]

actual_state = {}
for system in systems:
  try:
    result = health_check(system, timeout=3s)
    actual_state[system] = {
      "status": "UP",
      "response_time_ms": result.duration,
      "details": result.data
    }
  except timeout_error:
    actual_state[system] = {
      "status": "SLOW",
      "response_time_ms": ">3000",
      "details": "timeout"
    }
  except error:
    actual_state[system] = {
      "status": "DOWN",
      "error": str(error)
    }

# Compare to last known state
changes = diff(LAST_STATE, actual_state)
if changes:
  log("STATE CHANGE DETECTED", changes)
  trigger_alert("state-change", changes)

LAST_STATE = actual_state
```

### 3. MEMORY — Facts Current?

**Question:** "Is my persistent knowledge accurate and up-to-date?"

**User Memory (Who is Zeus?):**
- Name, role, work style
- Current projects (OneMind OS, Convergence, etc.)
- Technical stack preferences
- Known constraints (Python 3.11, PEP 668, venv required)
- Decision history (what was explicitly corrected?)
- Preferences (action-first, capability-forward, no overkill)

**System Memory (What's deployed?):**
- TAK Server credentials, URLs, version
- NATS cluster topology, account structure, port mapping
- Kubernetes cluster config, namespaces, RBAC rules
- ClickHouse database schema, retention policies, table structure
- Tailscale mesh configuration, exit nodes
- Active cron jobs and scheduled tasks

**Convergence Memory (What's documented?):**
- 52 files with 160+ embedded links
- Master index updated?
- Recently added resources (WAMI, LTX 2.3, RFC 10008, etc.)
- Integration patterns mapped?
- Links still working (not stale)?

**Fabric Topology:**
- Known agents on bus: who, where, what role?
- Protocol converters: how many, which types, operational?
- Sensor feeds: which live, which offline?
- Data flow patterns: from sensors → storage → analysis
- Backup/failover: tested recently?

**Calibration Check:**
```bash
# Load and timestamp memory
memory = read_file("~/.hermes/profiles/default/memory.md")
last_update = parse_timestamp(memory.metadata.updated)
age_days = (now() - last_update).days

if age_days > 7:
  log("WARNING: Memory older than 7 days, likely stale")

# Spot-check key facts
key_facts = [
  ("Zeus work style", ["action-first", "no-overkill"]),
  ("TAK Server", ["cloudtak.onemindos.dev", "8446"]),
  ("NATS broker", ["localhost:4222"]),
  ("Convergence", ["52 files", "1.0 MB"]),
  ("Legacy bus", ["agents.prompt.hermes.zeus.legacy"]),
]

for fact_name, required_phrases in key_facts:
  found = all(phrase in memory for phrase in required_phrases)
  if not found:
    log(f"WARNING: Memory fact may be stale: {fact_name}")
    needs_update = True

if needs_update:
  log("Triggering memory refresh...")
  trigger_memory_update()
```

### 4. ASSETS — Resources Available?

**Question:** "What resources do I have to work with?"

**Local Assets (on this machine):**
- Disk space: `/home`, `/`, `/data` (alert if <10% free)
- GPU memory: VRAM available, which models loaded?
- System RAM: available for processing (not swapped)
- Network bandwidth: to NATS broker, to cloud
- Cache: ComfyUI models, LM Studio models, downloaded files
- Credentials: API keys loaded, SSH keys accessible?

**Deployed Assets (in OneMind fabric):**
- Convergence files: where deployed? S3 bucket? Local mirror?
- Trained models: which deployed, which cached, versions?
- Databases: ClickHouse data volume, retention policy active?
- Backups: when last tested, restore time estimate?
- Licenses: valid until when? Renewal needed?

**Accessible Assets (via internet):**
- GitHub repos: network accessible? Rate limit status?
- HuggingFace models: can download? Mirror available?
- Cloud storage: S3, GCS, etc. reachable?
- APIs: Vercel Eve, external services online?
- Public datasets: accessible for training?

**Capability Assets:**
- Skills: which loaded in `~/.hermes/skills/`?
- Tools: which terminal commands available?
- Agents: which can be delegated to?
- Integrations: which third-party tools connected?

**Calibration Check:**
```bash
# Check local storage
disk_info = run("df -h / /home /data")
for path in ["/", "/home", "/data"]:
  percent_free = parse(disk_info[path].percent_used)
  if percent_free < 10:
    alert("disk-space-critical", path)

# Check GPU/memory
vram_info = run("nvidia-smi --query-gpu=memory.free --format=csv")
ram_info = run("free -h")
log("GPU VRAM available:", vram_info)
log("System RAM available:", ram_info)

# Check cache sizes
comfyui_cache = run("du -sh ~/.comfyui_models")
lm_studio_cache = run("du -sh ~/.lm-studio/models")
log("ComfyUI cache:", comfyui_cache)
log("LM Studio cache:", lm_studio_cache)

# Check credentials loaded
ssh_keys = run("ls -la ~/.ssh/ | wc -l")
api_keys = run("env | grep -c 'API_KEY\\|TOKEN'")
log(f"SSH keys found: {ssh_keys}, API keys in env: {api_keys}")

# Test external connectivity
for url in ["https://github.com", "https://huggingface.co"]:
  response = http_check(url)
  status = "UP" if response.ok else "DOWN"
  log(f"{url}: {status}")

# Build asset manifest
ASSET_MANIFEST = {
  "disk_free_percent": {"/": 25, "/home": 60, "/data": 45},
  "gpu_vram_free_gb": 8,
  "system_ram_free_gb": 12,
  "models_cached": 15,
  "skills_available": 45,
  "credentials_loaded": 5,
  "agents_reachable": ["mother", "edge", "claude-code", "codex"],
  "external_apis_accessible": ["github", "huggingface"],
  "timestamp": now(),
}
```

### 5. SPEED — Performance Baseline

**Question:** "Am I running at baseline speed, degraded, or fast?"

**Latency Baseline (target values):**
- NATS local publish-subscribe roundtrip: <100ms (local), <500ms (to cloud)
- Kubernetes API response: <200ms
- ClickHouse simple query: <100ms
- TAK Server API: <200ms
- Model inference: varies per model

**Throughput Baseline:**
- NATS: messages/second capability
- ClickHouse: inserts/second, queries/second
- Kubernetes: pod creation time, deployment rollout time
- ComfyUI: tokens/sec (LLM generation rate)

**Bottleneck Detection:**
- Which subsystem is rate-limiting? (slowest = bottleneck)
- Storage I/O: disk read/write speed?
- Network: bandwidth utilization?
- CPU: core saturation on critical services?
- GPU: utilization on inference services?

**Degradation Detection:**
- Is this slower than yesterday (historical baseline)?
- Is this slower than 7-day average?
- Which subsystem regressed?
- Trigger alert if >20% slower than baseline

**Calibration Check:**
```bash
# Measure NATS latency (100 samples)
latencies = []
for i in range(1, 101):
  t0 = time.now()
  run("nats pub test.latency 'ping' --server nats://localhost:4222")
  latency_ms = (time.now() - t0) * 1000
  latencies.append(latency_ms)

p50_latency = percentile(latencies, 50)   # median
p99_latency = percentile(latencies, 99)   # tail

# Measure K8s API speed
t0 = time.now()
run("kubectl get nodes")
k8s_latency_ms = (time.now() - t0) * 1000

# Measure ClickHouse query
t0 = time.now()
run("clickhouse-client --query 'SELECT 1'")
ch_latency_ms = (time.now() - t0) * 1000

# Measure LLM inference speed (if model loaded)
if llm_model_loaded:
  t0 = time.now()
  result = llm_model.predict("hello world")
  inference_time_ms = (time.now() - t0) * 1000
  tokens_generated = count_tokens(result)
  tokens_per_sec = tokens_generated / (inference_time_ms / 1000)

# Build speed baseline
SPEED_BASELINE = {
  "nats_p50_ms": p50_latency,
  "nats_p99_ms": p99_latency,
  "k8s_api_ms": k8s_latency_ms,
  "clickhouse_query_ms": ch_latency_ms,
  "llm_tokens_per_sec": tokens_per_sec if llm_model_loaded else None,
  "timestamp": now(),
}

# Compare to historical baseline
if p99_latency > HISTORICAL_BASELINE.nats_p99_ms * 1.2:
  log("WARNING: NATS latency degraded >20%")
  alert("performance-regression", {"service": "nats", "severity": "medium"})

if k8s_latency_ms > HISTORICAL_BASELINE.k8s_api_ms * 1.2:
  log("WARNING: K8s API latency degraded >20%")
  alert("performance-regression", {"service": "k8s", "severity": "medium"})
```

---

## Calibration Workflows

### Full Calibration (Startup or Daily)

**Duration:** 10-15 minutes

**Executed:** 
- On Legacy startup
- Daily 00:00 UTC (mother checks entire fabric)
- Manual trigger: `legacy-calibrate`

**Steps:**
```
STEP 1: MODES (2 min)
  └─ Determine actual operating mode
  └─ Update OPERATING_MODE state variable

STEP 2: STATE (3 min)
  └─ Poll all critical systems
  └─ Compare to last known state
  └─ Log any changes with timestamps

STEP 3: MEMORY (2 min)
  └─ Load memory file and check timestamp
  └─ Sample-verify key facts
  └─ Flag anything older than 7 days

STEP 4: ASSETS (2 min)
  └─ Check local disk, GPU, RAM available
  └─ Verify credentials are loaded
  └─ List reachable agents and APIs

STEP 5: SPEED (3 min)
  └─ Measure baseline latencies
  └─ Compare to historical baseline
  └─ Alert on >20% degradation

STEP 6: REPORT & DECIDE (1 min)
  └─ Generate calibration report
  └─ Make operational decisions:
     ├─ All green? → Proceed with normal operations
     ├─ Some yellow? → Log issues, reduce scope
     ├─ Red? → Escalate to Zeus, fail gracefully
     ├─ Offline-first triggered? → Activate edge autonomy
     └─ Critical failure? → Emergency mode
```

**Output:** Calibration report published to `om_meta.calibration.report` (NATS topic)

### Quick Calibration (Before Critical Operations)

**Duration:** 1-2 minutes

**Executed:** Before important tasks requiring high confidence

**Subset checks:**
```
├─ Is mother reachable? (5 sec ping)
├─ Is NATS broker up? (5 sec broker check)
├─ Is K8s cluster healthy? (5 sec node query)
├─ Is ClickHouse responsive? (5 sec query)
└─ Are delegated agents available? (5 sec status check)

Decision logic:
  ├─ All green (5/5)? → Proceed normally
  ├─ Yellow (3-4/5)? → Proceed with reduced scope
  ├─ Red (0-2/5)? → Fail fast, don't proceed
```

### Continuous Pulse (Background Monitoring)

**Duration:** <1 second per check

**Executed:** Every 30 seconds (even during other operations)

**Lightweight checks:**
```
├─ NATS broker responding to ping?
├─ Disk space dropped below 5%?
├─ GPU memory exhausted?
├─ Any critical service process died?
└─ If ANY failure detected → trigger full calibration
```

---

## Implementation for OneMind

### Mother (Cloud Hub)

**Execution:**
- Daily 00:00 UTC
- Scheduled cron job in cloud K8s

**Scope:**
- Full fabric topology check
- All edge nodes health
- Cross-region latency measurement
- Replicate calibration state to all edge nodes
- Publish calibration report to NATS

**If degraded:**
- Alert Zeus via automated message
- Reduce automation burden (defer non-critical tasks)
- Prepare for edge-first failover
- Keep edge nodes informed

### Edge Nodes (Local)

**Execution:**
- On startup (before any automation)
- Every 60 minutes (scheduled)
- Before delegating critical tasks (on-demand)

**Scope:**
- Local system health only (disk, GPU, network)
- Attempt connection to mother
- If mother reachable: pull latest state from mother
- If mother unreachable: operate fully autonomous
- Cache latest state from mother for fallback

**Decision Logic:**
```
if mother_reachable():
  mode = NORMAL
  sync_state_from_mother()
else:
  mode = OFFLINE_FIRST
  activate_autonomous_mode()
  queue_state_updates_for_later_sync()
```

### Both (Mother + Edge)

**Command:** `legacy-calibrate` (manual trigger)

**Use when:**
- Before critical mission (high confidence needed)
- After system maintenance (verify everything restarted OK)
- When performance seems off (detect degradation)
- When onboarding new systems (verify integration)

---

## Success Criteria

**After calibration completes, Legacy can answer with confidence:**

✅ **"What mode am I in?"**
   → NORMAL / DEGRADED / OFFLINE_FIRST / MAINTENANCE
   → Decision logic established

✅ **"Is system X alive?"**
   → Yes/no with latency
   → Actual response time measured

✅ **"What resources do I have?"**
   → Disk free: XX GB
   → GPU memory: XX GB
   → System RAM: XX GB
   → Network bandwidth: OK

✅ **"Is my memory current?"**
   → Age of facts: X days
   → Staleness level: OK / WARNING / CRITICAL
   → What needs refresh?

✅ **"Am I running fast or slow?"**
   → Baseline latency: P50, P99
   → Current latency: measured
   → Degradation: 0% / 15% / 25%+

✅ **"What changed since last calibration?"**
   → Diff of state changes
   → Which systems went up/down
   → Which metrics improved/regressed

✅ **"Can I delegate to agent Y?"**
   → Yes/no
   → Expected latency estimate
   → Will task timeout?

✅ **"Should I activate offline-first mode?"**
   → Decision logic triggered
   → Confidence level: high / medium / low
   → Fallback ready

✅ **"What should I report to Zeus?"**
   → Concise calibration summary
   → Alert level: OK / CAUTION / ALERT / CRITICAL
   → Recommended actions

---

## Integration with Convergence

### Why Calibration Matters to OneMind

**Mother is sovereign over the entire fabric.** But sovereignty is only credible if she knows what she's talking about.

**Convergence documents 52 resources** — Mother needs to know:
- Which resources are deployed where in fabric?
- Which regions have what capabilities?
- Which systems are healthy, which failing?
- What's the cross-region latency baseline?
- Where are all her edge nodes?

**Example:** Mother reads about LTX 2.3 (text-to-video generation)
- Question: "Can I order OneMind Dojo to generate training scenarios fabric-wide?"
- Needs to know: Which edge nodes have ComfyUI? GPU VRAM available? Model cached?
- Calibration answers: "3 edge nodes have LTX 2.3 cached, avg latency 150ms, 24GB VRAM available per node"
- Decision: "Yes, I'm authorizing distributed generation job across these 3 nodes"

**Without calibration:**
- Mother might order generation on a node that doesn't have the model
- Or assume available resources that aren't
- Or not realize the fabric is too slow for the task

**With calibration:**
- Mother knows exactly what she's commanding
- Edge nodes obey with confidence
- Fabric operates coherently

---

## Related Documentation

- **Skill:** `legacy-self-calibration` (Hermes)
- **System:** `zeus-legacy-local-first` (OneMind OS)
- **Framework:** `omos-core` (OneMind architecture)
- **Delegation:** `synadia-agent-protocol` (NATS agent communication)

---

## Timeline & Adoption

**2026-07-11:** Skill created (initial implementation)

**2026-07: Testing:** Validate on Legacy edge node

**2026-08:** Integrate into daily mother calibration (00:00 UTC)

**2026-09:** Expand to all edge nodes (if multi-node deployment)

**2026-Q4:** Include in OneMind Dojo training scenarios

---

## Summary

**Legacy Self-Calibration = AI System Sync Tool**

✅ **Five-point verification:** Modes, State, Memory, Assets, Speed

✅ **Prevents drift:** Stays synchronized with fabric reality

✅ **Enables autonomy:** Knows what's available, what's not

✅ **Informed decisions:** Can measure performance, detect degradation

✅ **Failover ready:** Detects offline-first situation immediately

✅ **Ground truth:** No more stale assumptions

**For Zeus:** Can trust that Legacy is always aware of current system state

**For OneMind fabric:** Reduces failure cascades, enables graceful degradation

---

**Captured:** 2026-07-11

**Status:** Production-ready autonomous skill

**Trigger:** Daily 00:00 UTC (mother), on-startup (edges), manual (`legacy-calibrate`)

**Maintenance:** Quarterly review of metrics thresholds
