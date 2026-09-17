---
id: ui-12-20260111053926
title: "Continuous Maintenance Workflows"
type: concept
domain: SO
domain: 00-24
folder: "00-24 SO (Sovereign Operations)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/so, folder/12]
todoist_project: "00-24 SO (Sovereign Operations)"
ai_summary: "Continuous Maintenance Workflows"
---

# Continuous Maintenance Workflows

This document defines automated maintenance workflows for the One Mind system. These workflows ensure the system stays healthy, synchronized, and up-to-date with minimal human intervention.

## Table of Contents

- [Overview](#overview)
- [Daily Workflows](#daily-workflows)
- [Weekly Workflows](#weekly-workflows)
- [Monthly Workflows](#monthly-workflows)
- [On-Demand Workflows](#on-demand-workflows)
- [Automated Recovery](#automated-recovery)
- [Monitoring and Alerts](#monitoring-and-alerts)

## Overview

### Maintenance Philosophy

One Mind is designed for **AI-driven continuous maintenance**:
- **Automated**: Workflows run on schedule without human intervention
- **Self-healing**: System detects and corrects drift automatically
- **Observable**: All actions logged and trackable
- **Human-in-the-loop**: Critical changes require approval
- **GitOps**: Configuration stored in Git, deployed via automation

### Maintenance Layers

```
┌─────────────────────────────────────────────┐
│  N8N Workflows (Orchestration)              │
│  - Scheduled triggers                       │
│  - Cross-service coordination               │
└─────────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────────┐
│  Python Scripts (Operations)                │
│  - validate_system.py                       │
│  - sync_agent_configs.py                    │
│  - backup automation                        │
└─────────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────────┐
│  External APIs (Services)                   │
│  - Letta Cloud (41 agents)                  │
│  - Todoist (task management)                │
│  - GitHub (version control)                 │
└─────────────────────────────────────────────┘
```

## Daily Workflows

### 1. Morning Health Check
**Schedule**: Every day at 5:00 AM
**N8N Workflow**: `one-mind-agent-health-check.json`
**Duration**: ~5 minutes

**What it does**:
1. Validates all system components are online
2. Checks agent availability in Letta Cloud
3. Verifies Todoist project structure
4. Tests N8N workflow execution
5. Checks disk space and Docker health
6. Sends summary to Morning Briefing

**Implementation**:
```python
# Can be run manually:
cd infrastructure
python scripts/validate_system.py

# Or via N8N workflow with HTTP Request nodes:
# 1. GET https://api.letta.com/agents → Count active agents
# 2. GET https://api.todoist.com/rest/v2/projects → Verify structure
# 3. Execute validate_system.py via SSH/Docker exec
# 4. POST results to notification channel
```

**Success Criteria**:
- ✅ All 41 agents respond to API calls
- ✅ All required projects exist in Todoist
- ✅ All critical N8N workflows are active
- ✅ Docker containers running
- ✅ Disk space > 10GB free

**Failure Actions**:
- Send alert to user via Todoist task (priority: urgent)
- Log error details to `infrastructure/logs/health-check.log`
- Attempt automated recovery (see [Automated Recovery](#automated-recovery))

---

### 2. Morning Briefing Generation
**Schedule**: Every day at 6:00 AM (after health check)
**N8N Workflow**: `one-mind-morning-briefing.json`
**Duration**: ~3 minutes

**What it does**:
1. Aggregates tasks from all 4 pillars (HP, LE, GE, IT)
2. Queries agents for overnight activities
3. Pulls GitHub activity (commits, PRs, issues)
4. Generates prioritized action list
5. Posts to daily briefing location

**Implementation**:
```javascript
// N8N workflow nodes:
1. Todoist API → Get today's tasks (filter by due_date)
2. For each pillar:
   - Query pillar coordinator agent via Letta API
   - Aggregate important updates
3. GitHub API → Get recent activity
4. Merge all data
5. POST to Obsidian via file write or webhook
```

**Output Format**:
```markdown
# Daily Briefing - 2024-12-23

## Health Status
✅ All systems operational

## Priority Tasks
- [ ] HP: Review sleep metrics from last week
- [ ] LE: Update weekly reflection notes
- [ ] GE: Review business metrics dashboard

## Agent Activity (overnight)
- Mother AI: Processed 3 maintenance tasks
- HP Coordinator: Updated habit tracking configs

## GitHub Activity
- 2 commits to infrastructure
- 1 PR opened: "Update agent system prompts"
```

---

### 3. Agent Config Drift Detection
**Schedule**: Every day at 12:00 PM
**Script**: `sync_agent_configs.py --check`
**Duration**: ~2 minutes

**What it does**:
1. Pulls latest agent configs from GitHub
2. Compares with deployed state in Letta Cloud
3. Detects configuration drift
4. Reports differences (no auto-fix without approval)

**Implementation**:
```bash
# Automated via N8N or cron:
cd infrastructure
python scripts/sync_agent_configs.py --check

# If drift detected (exit code 1):
# - Create Todoist task: "Agent config drift detected - review required"
# - Attach drift report
# - User reviews and runs: sync_agent_configs.py --update-all
```

**Drift Examples**:
- System prompt changed in GitHub but not in Letta
- Model version updated in config but agent still using old model
- New agent added to YAML but not deployed

---

### 4. Evening Review Preparation
**Schedule**: Every day at 9:00 PM
**N8N Workflow**: `one-mind-evening-review.json`
**Duration**: ~3 minutes

**What it does**:
1. Collects completed tasks from all pillars
2. Aggregates metrics (tasks completed, time tracked, etc.)
3. Queries agents for insights and reflections
4. Generates evening review document
5. Creates prompt for user reflection

**Implementation**:
```javascript
// N8N workflow:
1. Todoist API → Get completed tasks (today)
2. For each pillar:
   - Count completions
   - Calculate completion rate
3. Letta API → Ask each coordinator for daily insights
4. Format evening review
5. Create Todoist task: "Complete evening review"
```

**Output Format**:
```markdown
# Evening Review - 2024-12-23

## Completion Summary
- HP: 8/10 tasks (80%)
- LE: 5/7 tasks (71%)
- GE: 3/5 tasks (60%)
- IT: 6/8 tasks (75%)

**Total**: 22/30 tasks (73%)

## Key Wins
- Completed morning workout (HP)
- Wrote 500 words in journal (LE)
- Reviewed business metrics (GE)

## Agent Insights
**HP Coordinator**: "Strong consistency in morning routine this week"
**LE Coordinator**: "Journal writing streak: 5 days"

## Reflection Prompts
- What went well today?
- What could improve tomorrow?
- Any blockers to address?
```

---

## Weekly Workflows

### 5. Weekly System Review
**Schedule**: Every Sunday at 8:00 AM
**N8N Workflow**: `one-mind-weekly-review.json`
**Duration**: ~10 minutes

**What it does**:
1. Analyzes entire week's data across all pillars
2. Identifies trends and patterns
3. Suggests optimizations
4. Checks for stale tasks/projects
5. Reviews agent performance metrics

**Implementation**:
```python
# Weekly analysis script (called from N8N):
def weekly_review():
    # 1. Fetch all completed tasks (last 7 days)
    completed_tasks = todoist_api.get_completed(since="7 days ago")

    # 2. Calculate completion rates by pillar
    pillar_stats = analyze_by_pillar(completed_tasks)

    # 3. Query each coordinator agent for weekly insights
    agent_insights = {}
    for coordinator in ["hp-coordinator", "le-coordinator", "ge-coordinator", "it-coordinator"]:
        response = letta_api.send_message(
            coordinator,
            "Please provide your weekly review insights"
        )
        agent_insights[coordinator] = response

    # 4. Identify stale tasks (not updated in 14+ days)
    stale_tasks = todoist_api.get_tasks(filter="overdue & !@waiting")

    # 5. Generate report
    report = generate_weekly_report(pillar_stats, agent_insights, stale_tasks)

    # 6. Create review task in Todoist
    todoist_api.create_task(
        content="Complete Weekly Review",
        project="00-09 Command",
        priority=3,
        description=report
    )
```

**Report Includes**:
- Task completion trends (7-day graph)
- Top performing pillar
- Areas needing attention
- Stale task list
- Agent insights and recommendations

---

### 6. Backup Verification
**Schedule**: Every Sunday at 2:00 AM
**Script**: `infrastructure/scripts/backup.sh` + verification
**Duration**: ~15 minutes

**What it does**:
1. Runs backup script (Docker volumes, configs, etc.)
2. Verifies backup integrity
3. Tests restore process on sample data
4. Confirms backup size and completeness
5. Rotates old backups (keep last 4 weeks)

**Implementation**:
```bash
#!/bin/bash
# Weekly backup verification workflow

# 1. Run backup
cd infrastructure
./scripts/backup.sh

# 2. Verify backup exists and is not empty
BACKUP_FILE="backups/one-mind-$(date +%Y-%m-%d).tar.gz"
if [ ! -f "$BACKUP_FILE" ] || [ ! -s "$BACKUP_FILE" ]; then
    echo "❌ Backup failed or empty"
    # Create alert task
    exit 1
fi

# 3. Test backup integrity
if ! tar -tzf "$BACKUP_FILE" > /dev/null; then
    echo "❌ Backup corrupted"
    exit 1
fi

# 4. Rotate old backups (keep last 4 weeks)
find backups/ -name "one-mind-*.tar.gz" -mtime +28 -delete

# 5. Log success
echo "✅ Backup verified: $BACKUP_FILE ($(du -h $BACKUP_FILE | cut -f1))"
```

**Verification Checklist**:
- ✅ Backup file created
- ✅ Backup size reasonable (> 100MB, < 10GB)
- ✅ Archive is valid (tar test passes)
- ✅ Old backups rotated
- ✅ Backup log updated

---

### 7. Dependency Update Check
**Schedule**: Every Sunday at 10:00 AM
**Duration**: ~5 minutes

**What it does**:
1. Checks for Docker image updates
2. Checks for Python package updates
3. Reviews Letta Cloud model availability
4. Checks for N8N version updates
5. Creates update tasks if needed

**Implementation**:
```bash
# Check Docker images
docker images --format "{{.Repository}}:{{.Tag}}" | while read image; do
    docker pull $image --quiet
    if [ $? -eq 0 ]; then
        echo "✅ $image up to date"
    else
        echo "⚠️ $image has updates"
    fi
done

# Check Python packages
cd infrastructure/scripts
pip list --outdated | grep -E "(letta|todoist|requests|pyyaml)"

# If updates found:
# - Create Todoist task: "Review dependency updates"
# - Attach update list
# - User decides when to apply
```

---

## Monthly Workflows

### 8. Memory Consolidation
**Schedule**: First Sunday of each month at 9:00 AM
**N8N Workflow**: `one-mind-memory-consolidation.json`
**Duration**: ~20 minutes

**What it does**:
1. Reviews agent conversation histories
2. Identifies important patterns and learnings
3. Consolidates into agent core memories
4. Archives old conversation logs
5. Updates agent knowledge bases

**Implementation**:
```python
# Monthly memory consolidation
def consolidate_memories():
    for agent_id in all_agents:
        # 1. Fetch last 30 days of conversations
        messages = letta_api.get_messages(agent_id, since="30 days ago")

        # 2. Analyze for patterns (could use Claude API for summarization)
        summary = analyze_conversations(messages)

        # 3. Update agent's core memory block
        if summary['important_learnings']:
            letta_api.update_memory_block(
                agent_id,
                block_name="core_memory",
                content=summary['important_learnings']
            )

        # 4. Archive old messages (optional, if storage is concern)
        # letta_api.archive_messages(agent_id, older_than="90 days")

        print(f"✅ Consolidated memories for {agent_id}")
```

**Memory Categories**:
- User preferences and patterns
- Important decisions made
- Recurring issues/solutions
- System optimizations discovered

---

### 9. API Key Rotation Check
**Schedule**: First day of each month at 8:00 AM
**Duration**: ~2 minutes

**What it does**:
1. Checks API key expiration dates
2. Reviews API usage/quota consumption
3. Tests all API keys still work
4. Reminds user to rotate keys (best practice: every 90 days)

**Implementation**:
```python
# Check API keys
api_keys = {
    'LETTA_API_KEY': os.getenv('LETTA_API_KEY'),
    'TODOIST_API_KEY': os.getenv('TODOIST_API_KEY'),
    'N8N_API_KEY': os.getenv('N8N_API_KEY'),
    'GITHUB_TOKEN': os.getenv('GITHUB_TOKEN'),
}

for key_name, key_value in api_keys.items():
    if not key_value:
        print(f"❌ {key_name} not set")
        continue

    # Test key validity (simple API call)
    try:
        test_api_key(key_name, key_value)
        print(f"✅ {key_name} valid")
    except Exception as e:
        print(f"❌ {key_name} invalid: {e}")
        # Create urgent task to fix

# Check last rotation date (tracked in infrastructure/config/key_rotation.yaml)
# If > 90 days: Create task to rotate
```

---

### 10. System Audit
**Schedule**: First Sunday of each month at 11:00 AM
**Duration**: ~30 minutes

**What it does**:
1. Full system health audit
2. Review all agent configurations
3. Check for unused resources
4. Analyze cost/usage metrics
5. Generate comprehensive report

**Implementation**:
```python
def monthly_audit():
    audit_report = {
        'date': datetime.now().isoformat(),
        'agents': {},
        'workflows': {},
        'storage': {},
        'costs': {}
    }

    # 1. Agent audit
    agents = letta_api.list_agents()
    for agent in agents:
        details = letta_api.get_agent(agent['id'])
        audit_report['agents'][agent['name']] = {
            'active': details.get('active', True),
            'last_message': details.get('last_message_date'),
            'message_count': details.get('message_count', 0),
            'model': details['llm_config']['model']
        }

    # 2. Workflow audit
    workflows = n8n_api.list_workflows()
    for workflow in workflows:
        audit_report['workflows'][workflow['name']] = {
            'active': workflow['active'],
            'last_run': workflow.get('last_run'),
            'execution_count': workflow.get('execution_count', 0)
        }

    # 3. Storage audit
    audit_report['storage'] = {
        'backups_size_gb': get_directory_size('infrastructure/backups') / (1024**3),
        'docker_volumes_gb': get_docker_volume_sizes(),
        'obsidian_vault_gb': get_directory_size('.') / (1024**3)
    }

    # 4. Cost estimation (if available from APIs)
    audit_report['costs'] = {
        'letta_estimated': estimate_letta_costs(),
        'todoist': 0,  # Free tier
        'n8n': 0,       # Self-hosted
    }

    # 5. Save report
    with open(f'infrastructure/logs/audit-{datetime.now().strftime("%Y-%m")}.yaml', 'w') as f:
        yaml.dump(audit_report, f)

    # 6. Create review task
    create_todoist_task("Review monthly system audit", priority=2)

    return audit_report
```

**Audit Report Includes**:
- Agent activity levels (active vs idle)
- Workflow execution statistics
- Storage usage trends
- Cost projections
- Recommendations for optimization

---

## On-Demand Workflows

### 11. Agent Config Sync (Manual Trigger)
**Trigger**: After editing agent YAML configs
**Command**: `python sync_agent_configs.py --update-all`

**When to use**:
- After bulk editing agent system prompts
- When upgrading all agents to new model version
- After adding new agents to configs
- When fixing configuration drift

**Process**:
```bash
# 1. Edit configs in GitHub or locally
vim infrastructure/config/agents/command.yaml

# 2. Validate changes
python scripts/validate_system.py

# 3. Dry-run to see what would change
python scripts/sync_agent_configs.py --dry-run

# 4. Apply changes
python scripts/sync_agent_configs.py --update-all

# 5. Verify
python scripts/sync_agent_configs.py --check
```

---

### 12. Emergency Restore
**Trigger**: System failure or data loss
**Command**: `./scripts/restore.sh [backup-date]`

**When to use**:
- Accidental deletion of agents/workflows
- Corrupted configuration
- Need to rollback to previous state
- Disaster recovery

**Process**:
```bash
# 1. List available backups
ls -lh infrastructure/backups/

# 2. Choose backup to restore
BACKUP_DATE="2024-12-22"

# 3. Run restore (with confirmation)
cd infrastructure
./scripts/restore.sh $BACKUP_DATE

# 4. Verify restoration
python scripts/validate_system.py

# 5. Test critical workflows
# - Send test message to Mother AI
# - Check Todoist projects exist
# - Verify N8N workflows active
```

---

### 13. Workflow Deployment
**Trigger**: After updating N8N workflow JSON files
**Command**: `./scripts/deploy_workflows.sh`

**When to use**:
- Created new N8N workflow
- Updated existing workflow logic
- Need to redeploy all workflows

**Process**:
```bash
# 1. Update workflow JSON files
# - Edit in N8N UI, export to infrastructure/workflows/production/
# - Or edit JSON directly

# 2. Deploy with idempotency (deletes existing, imports new)
cd infrastructure
./scripts/deploy_workflows.sh

# 3. Activate workflows in N8N
# 4. Test execution
```

---

## Automated Recovery

### Self-Healing Procedures

**Problem**: Agent becomes unresponsive
**Detection**: Health check timeout
**Recovery**:
```python
# Automated recovery attempt
try:
    response = letta_api.get_agent(agent_id, timeout=5)
except Timeout:
    # Step 1: Wait and retry
    time.sleep(30)
    try:
        response = letta_api.get_agent(agent_id, timeout=10)
    except Timeout:
        # Step 2: Create alert task
        create_urgent_task(f"Agent {agent_id} unresponsive - manual investigation needed")
        # Step 3: Log for human review
        log_error(f"Agent {agent_id} timeout after retry")
```

**Problem**: Configuration drift detected
**Detection**: Daily drift check
**Recovery**:
```python
# Auto-sync if drift is minor (system prompt only)
if drift_type == 'system_prompt' and auto_sync_enabled:
    sync_agent_configs.py --update [agent-name]
else:
    # Create task for human review
    create_task("Review agent config drift", attach_drift_report())
```

**Problem**: Disk space low (< 5GB)
**Detection**: Health check
**Recovery**:
```bash
# Automated cleanup
# 1. Remove old Docker images
docker image prune -a --force --filter "until=720h"  # 30 days

# 2. Rotate old logs
find infrastructure/logs/ -name "*.log" -mtime +30 -delete

# 3. Remove old backups (keep last 4 weeks)
find infrastructure/backups/ -mtime +28 -delete

# 4. Check if space recovered
# 5. If still low: Create urgent task
```

---

## Monitoring and Alerts

### Alert Levels

| Level | Description | Action | Examples |
|-------|-------------|--------|----------|
| **INFO** | Normal operation | Log only | Daily health check passed |
| **WARNING** | Non-critical issue | Log + Todoist task (normal) | Config drift detected |
| **ERROR** | Service degradation | Todoist task (high priority) | Agent unresponsive |
| **CRITICAL** | System failure | Todoist task (urgent) + possible notification | All agents down, backup failed |

### Alert Destinations

1. **Todoist Tasks** (primary)
   - Project: `00-09 Command (Meta) / System Alerts`
   - Priority based on severity
   - Tagged with: `#alert`, `#one-mind`

2. **Log Files**
   - Location: `infrastructure/logs/`
   - Rotation: Daily, keep 30 days
   - Format: JSON for easy parsing

3. **Obsidian Daily Note** (info only)
   - Append to today's note
   - Section: `## System Events`

### Monitoring Dashboards

**Grafana Dashboard** (if configured):
- Agent response times
- API call rates
- Task completion rates by pillar
- Backup success/failure
- Disk usage trends

**Access**: https://grafana.zeusops.dev

---

## Running Maintenance Manually

### Daily Checklist
```bash
cd /Users/zeuslegacy/Documents/One Mind (Hive Mind)/infrastructure

# Morning
python scripts/validate_system.py           # 2 min
python scripts/sync_agent_configs.py --check  # 1 min

# As needed
./scripts/deploy_workflows.sh              # If workflows changed
python scripts/sync_agent_configs.py --update-all  # If configs changed
```

### Weekly Checklist
```bash
# Sunday morning
./scripts/backup.sh                        # 5 min
python scripts/validate_system.py          # Full validation

# Review logs
tail -100 infrastructure/logs/health-check.log
tail -100 infrastructure/logs/sync.log
```

### Monthly Checklist
```bash
# First Sunday
# Review audit report (generated automatically)
cat infrastructure/logs/audit-$(date +%Y-%m).yaml

# Check for updates
docker images --format "{{.Repository}}:{{.Tag}}" | xargs -I {} docker pull {}

# Review API key rotation schedule
cat infrastructure/config/key_rotation.yaml
```

---

## Best Practices

1. **Let automation run** - Don't manually intervene unless alerted
2. **Review weekly reports** - Spot trends before they become problems
3. **Test changes in dry-run** - Before bulk updates
4. **Keep backups** - Verify monthly that restore works
5. **Monitor costs** - Check Letta usage monthly
6. **Rotate API keys** - Every 90 days
7. **Update dependencies** - Within 30 days of release
8. **Document changes** - Commit message explains why
9. **Version everything** - All configs in Git
10. **Trust but verify** - Automation should alert on anomalies

## Next Steps

- See [API-ACCESS.md](00-24%20SO%20(Sovereign%20Operations)/05%20Development%20(Code%20&%20Testing)/Cloud/API-ACCESS.md) for API integration details
- See [BULK-OPERATIONS.md](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/BULK-OPERATIONS.md) for system-wide change procedures
- See [DEPLOYMENT.md](00-24%20SO%20(Sovereign%20Operations)/05%20Development%20(Code%20&%20Testing)/DevOps/DEPLOYMENT.md) for full deployment guide
