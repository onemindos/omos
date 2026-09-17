---
id: ui-12-20260111053926
title: "Bulk Operations Guide"
type: concept
pillar: UI
domain: 00-24
folder: "00-24 UI (Unified Intelligence)"
status: active
created: 2026-01-11
modified: 2026-01-11
owner: zeus
tags: [domain/ui, folder/12]
todoist_project: "00-24 UI (Unified Intelligence)"
ai_summary: "Bulk Operations Guide"
---

# Bulk Operations Guide

This document explains how to make system-wide changes to the One Mind infrastructure. These procedures enable you to update all 41 agents, redeploy workflows, rotate API keys, and perform other bulk modifications efficiently.

## Table of Contents

- [Overview](#overview)
- [Agent Bulk Operations](#agent-bulk-operations)
- [Workflow Operations](#workflow-operations)
- [API Key Rotation](#api-key-rotation)
- [Todoist Bulk Operations](#todoist-bulk-operations)
- [Model Upgrades](#model-upgrades)
- [Emergency Procedures](#emergency-procedures)

## Overview

### Why Bulk Operations?

One Mind manages **41 agents**, **6+ workflows**, and **dozens of integrations**. Making changes one-by-one is:
- ❌ Time-consuming (hours)
- ❌ Error-prone (easy to miss agents)
- ❌ Not scalable

With bulk operations, you can:
- ✅ Update all 41 agents in < 5 minutes
- ✅ Deploy all workflows atomically
- ✅ Rotate all API keys at once
- ✅ Apply changes consistently

### Safety Principles

1. **Always use dry-run first** - Preview changes before applying
2. **Work on a branch** - Test changes before merging to main
3. **Backup before bulk changes** - Run `./scripts/backup.sh`
4. **Validate after changes** - Run `python scripts/validate_system.py`
5. **Rollback plan** - Know how to undo changes

### Bulk Operation Workflow

```
┌─────────────────────┐
│ 1. Backup system    │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│ 2. Edit configs     │ (YAML files, .env, etc.)
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│ 3. Validate locally │ (syntax check, dry-run)
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│ 4. Apply changes    │ (sync scripts)
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│ 5. Verify results   │ (health checks, tests)
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│ 6. Commit to Git    │ (document changes)
└─────────────────────┘
```

---

## Agent Bulk Operations

### Operation 1: Update All Agent System Prompts

**Use case**: Refine agent instructions, add new capabilities, fix prompts

**Steps**:

1. **Backup current state**
   ```bash
   cd infrastructure
   ./scripts/backup.sh
   ```

2. **Edit agent YAML configs**
   ```bash
   # Edit specific pillar
   vim config/agents/command.yaml

   # Or all pillars
   vim config/agents/hp.yaml
   vim config/agents/le.yaml
   vim config/agents/ge.yaml
   vim config/agents/it.yaml
   ```

   Example change:
   ```yaml
   agents:
     - name: "cmd-mother-ai"
       system_prompt: |
         You are Mother AI, the central coordinator of the One Mind system.

         ## NEW CAPABILITY (2024-12-23)
         You can now access GitHub directly to review infrastructure changes.

         Your responsibilities include:
         - Orchestrating all 41 agents
         # ... rest of prompt
   ```

3. **Validate YAML syntax**
   ```bash
   python scripts/validate_system.py
   ```

4. **Dry-run to see what will change**
   ```bash
   python scripts/sync_agent_configs.py --dry-run
   ```

   Output:
   ```
   =========================================================
   COMMAND PILLAR
   =========================================================

   ↻ cmd-mother-ai: Config drift detected
     → System prompt differs
     → Would update (dry-run mode)

   ✓ cmd-human-interface: In sync
   ✓ cmd-todoist-sync: In sync
   ✓ cmd-system-monitor: In sync
   ```

5. **Apply changes to Letta Cloud**
   ```bash
   # Update all agents across all pillars
   python scripts/sync_agent_configs.py --update-all

   # Or update specific pillar
   python scripts/sync_agent_configs.py --update hp
   ```

6. **Verify changes applied**
   ```bash
   # Check for remaining drift (should be none)
   python scripts/sync_agent_configs.py --check
   ```

   Expected output:
   ```
   ✓ All agents in sync (0 drifted)
   ```

7. **Test agent functionality**
   ```bash
   # Send test message to Mother AI via Letta API
   curl -X POST https://api.letta.com/agents/{agent_id}/messages \
     -H "Authorization: Bearer $LETTA_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"messages": [{"role": "user", "text": "Hello, test your new capability"}]}'
   ```

8. **Commit changes**
   ```bash
   git add config/agents/
   git commit -m "Update all agent system prompts - add GitHub access capability"
   git push
   ```

**Time**: ~5 minutes (for all 41 agents)

---

### Operation 2: Add New Agent

**Use case**: Expand system with new agent (e.g., new HP sub-agent)

**Steps**:

1. **Add to YAML config**
   ```bash
   vim config/agents/hp.yaml
   ```

   Add new agent:
   ```yaml
   agents:
     # ... existing HP agents

     - name: "hp-sleep-optimizer"
       display_name: "Sleep Optimizer"
       model: "anthropic/claude-sonnet-4-5"
       enable_sleeptime: true
       tags:
         - "hp"
         - "health"
         - "sleep"
       system_prompt: |
         You are the Sleep Optimizer agent within the Holistic Performance pillar.

         Your responsibilities:
         - Track sleep metrics from Oura/Whoop
         - Analyze sleep quality trends
         - Provide sleep optimization recommendations
         - Coordinate with HP Coordinator on health insights
   ```

2. **Add to deployed state tracker**
   ```bash
   vim config/deployed_state.yaml
   ```

   Add entry:
   ```yaml
   deployed_agents:
     # ... existing agents

     hp-sleep-optimizer:
       letta_id: null
       created_at: null
       last_updated: null
       status: "not_deployed"
   ```

3. **Deploy new agent**
   ```bash
   # Using letta_agent_setup.py
   cd scripts
   python letta_agent_setup.py --agent hp-sleep-optimizer

   # Or deploy entire HP pillar (will skip existing, create new)
   python letta_agent_setup.py --pillar hp
   ```

4. **Update deployed state with agent ID**
   - Script should auto-update `deployed_state.yaml` with new Letta agent ID
   - Verify:
     ```bash
     grep -A 3 "hp-sleep-optimizer" config/deployed_state.yaml
     ```

5. **Verify agent is live**
   ```bash
   python scripts/validate_system.py
   ```

6. **Commit changes**
   ```bash
   git add config/agents/hp.yaml config/deployed_state.yaml
   git commit -m "Add hp-sleep-optimizer agent for sleep tracking"
   git push
   ```

**Time**: ~2 minutes per agent

---

### Operation 3: Delete Agent

**Use case**: Remove unused or redundant agent

**Steps**:

1. **Remove from YAML config**
   ```bash
   vim config/agents/[pillar].yaml
   # Delete the agent entry
   ```

2. **Optionally delete from Letta Cloud**
   ```bash
   # Get agent ID
   AGENT_ID=$(grep -A 1 "agent-name:" config/deployed_state.yaml | grep letta_id | cut -d'"' -f2)

   # Delete via API
   curl -X DELETE https://api.letta.com/agents/$AGENT_ID \
     -H "Authorization: Bearer $LETTA_API_KEY"
   ```

3. **Update deployed state**
   ```bash
   vim config/deployed_state.yaml
   # Remove the agent entry or mark as "deleted"
   ```

4. **Commit changes**
   ```bash
   git add config/
   git commit -m "Remove [agent-name] - no longer needed"
   git push
   ```

---

## Workflow Operations

### Operation 4: Redeploy All N8N Workflows

**Use case**: Update workflow logic, fix bugs, add new workflows

**Steps**:

1. **Backup current workflows**
   ```bash
   ./scripts/backup.sh
   ```

2. **Update workflow JSON files**

   Option A: Export from N8N UI
   - Edit workflow in N8N: https://n8n.zeusops.dev
   - Click **Download** → Save to `infrastructure/workflows/production/`

   Option B: Edit JSON directly
   ```bash
   vim workflows/production/one-mind-morning-briefing.json
   ```

3. **Deploy updated workflows**
   ```bash
   # Deploy all workflows (idempotent - deletes existing, imports new)
   ./scripts/deploy_workflows.sh
   ```

   Output:
   ```
   ========================================
   ONE MIND — N8N Workflow Deployment
   ========================================

   Found 6 workflow(s) to deploy

   Deploying: one-mind-morning-briefing
     → API integration pending (placeholder)
     ✓ Would deploy: workflows/production/one-mind-morning-briefing.json

   ... (5 more workflows)

   ========================================
   DEPLOYMENT SUMMARY
   ========================================
   ✓ Deployed: 6
   ```

   **Note**: Current script is a placeholder. Full N8N API integration is TODO.

4. **Activate workflows in N8N**
   - Open N8N UI
   - Verify workflows imported
   - Activate each workflow

5. **Test workflow execution**
   ```bash
   # Trigger morning briefing workflow manually
   curl -X POST https://n8n.zeusops.dev/webhook/morning-briefing-test
   ```

6. **Commit changes**
   ```bash
   git add workflows/production/
   git commit -m "Update N8N workflows - add GitHub activity to morning briefing"
   git push
   ```

**Time**: ~10 minutes (for 6 workflows)

---

## API Key Rotation

### Operation 5: Rotate All API Keys

**Use case**: Security best practice (every 90 days), suspected compromise

**Keys to Rotate**:
- `LETTA_API_KEY` - Letta Cloud
- `TODOIST_API_KEY` - Todoist
- `N8N_API_KEY` - N8N instance
- `GITHUB_TOKEN` - GitHub access
- `LANGFUSE_PUBLIC_KEY` / `LANGFUSE_SECRET_KEY` - Observability

**Steps**:

1. **Generate new API keys**

   For each service:

   **Letta Cloud**:
   - Go to: https://app.letta.com/settings/api-keys
   - Generate new key
   - Copy key

   **Todoist**:
   - Go to: https://todoist.com/app/settings/integrations/developer
   - Revoke old token
   - Generate new token

   **N8N**:
   - Go to: https://n8n.zeusops.dev/settings/api
   - Create new API key

   **GitHub**:
   - Go to: https://github.com/settings/tokens
   - Generate new token with required scopes

2. **Update `.env` file**
   ```bash
   cd infrastructure
   vim .env
   ```

   Update keys:
   ```bash
   # Old keys (revoke after successful rotation)
   # LETTA_API_KEY=letta_key_old_xxx
   # TODOIST_API_KEY=todoist_old_xxx

   # New keys (active)
   LETTA_API_KEY=letta_key_new_xxx
   TODOIST_API_KEY=todoist_new_xxx
   N8N_API_KEY=new_n8n_key_xxx
   GITHUB_TOKEN=ghp_new_xxx
   ```

3. **Reload environment variables**
   ```bash
   source .env

   # Verify
   echo $LETTA_API_KEY | tail -c 5  # Should show last 4 chars of new key
   ```

4. **Test all integrations**
   ```bash
   python scripts/validate_system.py
   ```

   Expected output:
   ```
   ============================================================
   ENVIRONMENT VARIABLES
   ============================================================

   ✓ LETTA_API_KEY
   ✓ TODOIST_API_KEY

   ============================================================
   API ENDPOINTS
   ============================================================

   ✓ Letta Cloud API
   ✓ Todoist API
   ✓ GitHub API
   ```

5. **Update N8N workflow credentials**
   - Open N8N: https://n8n.zeusops.dev
   - Go to **Credentials**
   - Update each credential with new API keys:
     - Letta API
     - Todoist API
     - GitHub API

6. **Test workflows**
   ```bash
   # Manually trigger test workflow
   # Verify it can still call all APIs
   ```

7. **Revoke old API keys**
   - Once confirmed new keys work, revoke old keys from each service
   - This ensures old keys can't be used if compromised

8. **Update rotation log**
   ```bash
   vim config/key_rotation.yaml
   ```

   ```yaml
   rotations:
     - date: "2024-12-23"
       keys_rotated:
         - LETTA_API_KEY
         - TODOIST_API_KEY
         - N8N_API_KEY
         - GITHUB_TOKEN
       next_rotation_due: "2025-03-23"  # 90 days
   ```

9. **Commit (only the log, NOT the .env)**
   ```bash
   git add config/key_rotation.yaml
   git commit -m "Rotate all API keys - quarterly security maintenance"
   git push
   ```

   **CRITICAL**: `.env` is in `.gitignore` - never commit API keys!

**Time**: ~20 minutes (for all services)

---

## Todoist Bulk Operations

### Operation 6: Update All Todoist Projects

**Use case**: Restructure project hierarchy, update project templates

**Steps**:

1. **Get current project structure**
   ```bash
   cd scripts
   python todoist_setup.py --list
   ```

2. **Backup Todoist data**
   - Todoist has built-in backups
   - Or export via Settings → Backups

3. **Update project configs**
   ```bash
   vim todoist_setup.py
   ```

   Modify project structure:
   ```python
   PROJECTS = {
       "HP": {
           "color": "blue",
           "sections": ["Daily", "Weekly", "Monthly", "Habits"],
           # ... update as needed
       },
       # ... other pillars
   }
   ```

4. **Apply changes**
   ```bash
   python todoist_setup.py --update-all
   ```

   Script will:
   - Create missing projects
   - Update project colors/names
   - Add missing sections
   - Skip existing items

5. **Verify in Todoist**
   - Open: https://todoist.com
   - Check all 4 pillar projects
   - Verify sections exist

6. **Commit changes**
   ```bash
   git add todoist_setup.py
   git commit -m "Update Todoist project structure - add Habits section"
   git push
   ```

**Time**: ~5 minutes

---

### Operation 7: Bulk Task Creation

**Use case**: Create recurring tasks for all pillars, initialize new workflow

**Steps**:

1. **Create task creation script**
   ```bash
   cd scripts
   vim bulk_create_tasks.py
   ```

   ```python
   import os
   import requests

   TODOIST_API_KEY = os.getenv("TODOIST_API_KEY")
   TODOIST_BASE_URL = "https://api.todoist.com/rest/v2"

   headers = {
       "Authorization": f"Bearer {TODOIST_API_KEY}",
       "Content-Type": "application/json"
   }

   # Get project IDs
   projects = requests.get(f"{TODOIST_BASE_URL}/projects", headers=headers).json()
   project_map = {p['name']: p['id'] for p in projects}

   # Define tasks to create
   tasks = [
       {
           "content": "Daily Review - HP",
           "project_id": project_map["HP - Holistic Performance"],
           "due_string": "every day at 9pm",
           "priority": 2,
           "labels": ["daily-review", "automated"]
       },
       {
           "content": "Daily Review - LE",
           "project_id": project_map["LE - Legacy Evolution"],
           "due_string": "every day at 9pm",
           "priority": 2,
           "labels": ["daily-review", "automated"]
       },
       # ... more tasks
   ]

   # Create tasks
   for task in tasks:
       response = requests.post(
           f"{TODOIST_BASE_URL}/tasks",
           headers=headers,
           json=task
       )
       if response.status_code == 200:
           print(f"✓ Created: {task['content']}")
       else:
           print(f"✗ Failed: {task['content']} - {response.text}")
   ```

2. **Run script**
   ```bash
   python bulk_create_tasks.py
   ```

3. **Verify in Todoist**
   - Check tasks were created
   - Verify recurring schedule

**Time**: ~2 minutes (for dozens of tasks)

---

## Model Upgrades

### Operation 8: Upgrade All Agents to New Model

**Use case**: New Claude version released (e.g., Sonnet 4.5 → Sonnet 5.0)

**Steps**:

1. **Backup configs**
   ```bash
   ./scripts/backup.sh
   ```

2. **Update all YAML configs**
   ```bash
   cd config/agents

   # Option A: Manual find-and-replace
   # vim all files, replace: anthropic/claude-sonnet-4-5 → anthropic/claude-sonnet-5-0

   # Option B: Automated with sed
   find . -name "*.yaml" -exec sed -i '' 's/anthropic\/claude-sonnet-4-5/anthropic\/claude-sonnet-5-0/g' {} \;
   ```

3. **Verify changes**
   ```bash
   # Check that all configs now use new model
   grep -r "model:" . | grep -v "sonnet-5-0"
   # Should return no results (or only haiku agents)
   ```

4. **Dry-run sync**
   ```bash
   python scripts/sync_agent_configs.py --dry-run
   ```

   Output:
   ```
   ↻ cmd-mother-ai: Config drift detected
     → Model differs: anthropic/claude-sonnet-4-5 → anthropic/claude-sonnet-5-0
     → Would update (dry-run mode)

   ↻ hp-coordinator: Config drift detected
     → Model differs: anthropic/claude-sonnet-4-5 → anthropic/claude-sonnet-5-0
     → Would update (dry-run mode)

   ... (39 more agents)
   ```

5. **Apply model upgrade**
   ```bash
   python scripts/sync_agent_configs.py --update-all
   ```

   Output:
   ```
   ✓ Updated: 41 agents
   ```

6. **Test agents with new model**
   ```bash
   # Send test message to Mother AI
   # Verify responses are coherent and model upgrade worked
   ```

7. **Monitor costs**
   - Check Letta billing: https://app.letta.com/billing
   - New model may have different pricing

8. **Commit changes**
   ```bash
   git add config/agents/
   git commit -m "Upgrade all agents to Claude Sonnet 5.0"
   git push
   ```

**Time**: ~10 minutes (for all 41 agents)

**Cost Impact**: May increase or decrease depending on new model pricing

---

## Emergency Procedures

### Operation 9: Emergency Rollback

**Use case**: Bulk change broke system, need to restore previous state

**Steps**:

1. **Identify last good backup**
   ```bash
   ls -lh infrastructure/backups/
   ```

   Output:
   ```
   one-mind-2024-12-22.tar.gz  (last good state)
   one-mind-2024-12-23.tar.gz  (current broken state)
   ```

2. **Stop all automated workflows**
   - Open N8N: https://n8n.zeusops.dev
   - Deactivate all workflows
   - This prevents system from making more changes during rollback

3. **Restore from backup**
   ```bash
   cd infrastructure
   ./scripts/restore.sh 2024-12-22
   ```

   Follow prompts:
   ```
   WARNING: This will overwrite current configs
   Continue? [y/N]: y

   Restoring from backup: backups/one-mind-2024-12-22.tar.gz
   ✓ Extracted backup
   ✓ Restored config files
   ✓ Restored deployed state
   ```

4. **Sync agents back to previous configs**
   ```bash
   python scripts/sync_agent_configs.py --update-all
   ```

5. **Verify system health**
   ```bash
   python scripts/validate_system.py
   ```

6. **Re-enable workflows**
   - Go back to N8N
   - Activate workflows

7. **Document what went wrong**
   ```bash
   vim infrastructure/logs/incidents/2024-12-23-rollback.md
   ```

   ```markdown
   # Incident: Emergency Rollback 2024-12-23

   ## What Happened
   - Bulk update to agent prompts caused agents to become unresponsive
   - Root cause: Syntax error in system_prompt YAML

   ## Impact
   - All 41 agents offline for 15 minutes
   - Morning briefing workflow failed

   ## Resolution
   - Rolled back to 2024-12-22 backup
   - Fixed YAML syntax error
   - Re-applied update with validation

   ## Prevention
   - Always run validate_system.py before bulk updates
   - Add YAML syntax check to pre-commit hook
   ```

8. **Fix the issue**
   - Identify what caused the problem
   - Apply fix carefully
   - Test thoroughly before re-deploying

**Time**: ~10 minutes

---

### Operation 10: Emergency Agent Shutdown

**Use case**: Agent behaving unexpectedly, need to stop it immediately

**Steps**:

1. **Identify agent ID**
   ```bash
   grep -A 1 "agent-name:" config/deployed_state.yaml | grep letta_id
   ```

2. **Disable agent via API**
   ```bash
   AGENT_ID="agent-xxx"

   curl -X PATCH https://api.letta.com/agents/$AGENT_ID \
     -H "Authorization: Bearer $LETTA_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"active": false}'
   ```

3. **Verify agent is inactive**
   ```bash
   curl https://api.letta.com/agents/$AGENT_ID \
     -H "Authorization: Bearer $LETTA_API_KEY" | jq '.active'
   # Should return: false
   ```

4. **Investigate issue**
   - Check agent conversation logs in Letta UI
   - Review recent system prompt changes
   - Check for API errors

5. **Fix and re-enable**
   ```bash
   # After fixing issue
   curl -X PATCH https://api.letta.com/agents/$AGENT_ID \
     -H "Authorization: Bearer $LETTA_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"active": true}'
   ```

**Time**: ~5 minutes

---

## Best Practices for Bulk Operations

### Pre-Flight Checklist

Before any bulk operation:
- [ ] Backup current state (`./scripts/backup.sh`)
- [ ] Run validation (`python scripts/validate_system.py`)
- [ ] Use dry-run mode first (`--dry-run`)
- [ ] Test on single agent/workflow before bulk
- [ ] Have rollback plan ready
- [ ] Schedule during low-usage time (early morning)

### Post-Operation Checklist

After any bulk operation:
- [ ] Verify changes applied (`--check` mode)
- [ ] Run full validation (`validate_system.py`)
- [ ] Test critical workflows manually
- [ ] Check for error logs
- [ ] Monitor for 24 hours
- [ ] Commit changes to Git with descriptive message

### Common Mistakes to Avoid

1. ❌ **No backup before changes** → Can't rollback if something breaks
2. ❌ **Skipping dry-run** → Unexpected changes applied
3. ❌ **Not validating YAML** → Syntax errors break sync
4. ❌ **Committing .env file** → API keys exposed in Git
5. ❌ **No testing** → Broken changes deployed to production
6. ❌ **Changing too much at once** → Hard to identify root cause if issue
7. ❌ **Not documenting changes** → Future confusion about what changed
8. ❌ **Ignoring error messages** → Silent failures

### Bulk Operation Templates

**Quick Reference**:

```bash
# Update all agents
python scripts/sync_agent_configs.py --dry-run  # Preview
python scripts/sync_agent_configs.py --update-all  # Apply

# Deploy all workflows
./scripts/deploy_workflows.sh

# Validate system
python scripts/validate_system.py

# Backup system
./scripts/backup.sh

# Rollback system
./scripts/restore.sh [date]
```

---

## Next Steps

- See [API-ACCESS.md](00-24%20SO%20(Sovereign%20Operations)/03%20Command%20(Operations%20&%20Navigation)%20(The%20Controls)/ops-runbooks/API-ACCESS.md) for API integration details
- See [CONTINUOUS-MAINTENANCE.md](00-24%20SO%20(Sovereign%20Operations)/03%20Command%20(Operations%20&%20Navigation)%20(The%20Controls)/ops-runbooks/CONTINUOUS-MAINTENANCE.md) for automated maintenance
- See [DEPLOYMENT.md](00-24%20SO%20(Sovereign%20Operations)/03%20Command%20(Operations%20&%20Navigation)%20(The%20Controls)/ops-runbooks/DEPLOYMENT.md) for full deployment procedures


---

## Related

- [[00-24 SO (Sovereign Operations)/03 Command (Operations & Navigation) (The Controls)/03 Command|03 Command]]
