---
id: ui-12-20260111053926
title: "One Mind Deployment Guide"
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
ai_summary: "One Mind Deployment Guide"
---

# One Mind Deployment Guide

This document provides complete deployment procedures for the One Mind system, from fresh installation to production deployment and ongoing updates.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Fresh Installation](#fresh-installation)
- [Component Deployment](#component-deployment)
- [Production Deployment](#production-deployment)
- [Update Procedures](#update-procedures)
- [Verification and Testing](#verification-and-testing)
- [Troubleshooting](#troubleshooting)

## Overview

### One Mind Architecture

One Mind consists of three deployment layers:

```
┌─────────────────────────────────────────────┐
│  Cloud Services (Letta, Todoist)            │
│  - 41 Letta agents in Letta Cloud           │
│  - Todoist task management                  │
│  - Managed via APIs (no deployment needed)  │
└─────────────────────────────────────────────┘
            ↓ (API calls)
┌─────────────────────────────────────────────┐
│  Self-Hosted Services (Docker on VM)        │
│  - N8N (workflow automation)                │
│  - Prometheus + Grafana (monitoring)        │
│  - Uptime Kuma (health checks)              │
│  - Langfuse (observability)                 │
└─────────────────────────────────────────────┘
            ↓ (local access)
┌─────────────────────────────────────────────┐
│  Local Environment (MacBook)                │
│  - Obsidian vault (knowledge base)          │
│  - Infrastructure configs (Git repo)        │
│  - Deployment scripts                       │
└─────────────────────────────────────────────┘
```

### Deployment Targets

1. **Letta Cloud** - Agent deployment
2. **Todoist Cloud** - Project/task structure
3. **VM (zeusops.dev)** - Docker services
4. **Local** - Configuration and orchestration

---

## Prerequisites

### Required Accounts

- [ ] **Letta Cloud** account (https://app.letta.com)
  - API key generated
  - Billing configured (if needed)
- [ ] **Todoist** account (https://todoist.com)
  - Premium subscription (for labels, filters)
  - API token generated
- [ ] **GitHub** account
  - Personal access token with repo access
- [ ] **VM or VPS** (optional, for self-hosted services)
  - Ubuntu 22.04 LTS recommended
  - 4GB RAM minimum
  - 50GB disk space
  - Docker installed

### Required Software

**On Local Machine**:
- [ ] Git
- [ ] Python 3.10+
- [ ] Docker Desktop (for local testing)
- [ ] Obsidian (for vault management)
- [ ] SSH client (for VM access)

**On VM (if deploying self-hosted services)**:
- [ ] Docker
- [ ] Docker Compose
- [ ] Git

### API Keys and Credentials

Create `.env` file in `infrastructure/`:

```bash
cd /Users/zeuslegacy/Documents/One Mind (Hive Mind)/infrastructure
cp .env.example .env
vim .env
```

Fill in all values:
```bash
# Required
LETTA_API_KEY=letta_key_xxx
TODOIST_API_KEY=todoist_token_xxx

# Optional (for self-hosted services)
N8N_API_KEY=your_n8n_api_key
GITHUB_TOKEN=ghp_xxx

# Observability (optional)
LANGFUSE_PUBLIC_KEY=pk-lf-xxx
LANGFUSE_SECRET_KEY=sk-lf-xxx

# N8N Configuration
N8N_ENCRYPTION_KEY=$(openssl rand -hex 32)
N8N_BASIC_AUTH_USER=admin
N8N_BASIC_AUTH_PASSWORD=secure_password_here
```

**Security**: Never commit `.env` to Git! It's in `.gitignore`.

---

## Fresh Installation

### Step 1: Clone Repository

```bash
# Clone the One Mind repository
git clone https://github.com/your-username/one-mind.git
cd one-mind/infrastructure
```

Or if already in Obsidian vault:
```bash
cd /Users/zeuslegacy/Documents/One Mind (Hive Mind)/infrastructure
```

### Step 2: Install Python Dependencies

```bash
# Create virtual environment
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt
```

Required packages:
```
pyyaml>=6.0
requests>=2.31.0
python-dotenv>=1.0.0
```

### Step 3: Configure Environment

```bash
# Copy example env file
cp .env.example .env

# Edit with your API keys
vim .env

# Load environment
source .env

# Verify
echo $LETTA_API_KEY | tail -c 5  # Should show last 4 chars
```

### Step 4: Validate Configuration

```bash
# Run pre-flight validation
python scripts/validate_system.py
```

Expected output:
```
============================================================
ENVIRONMENT VARIABLES
============================================================

✓ LETTA_API_KEY
✓ TODOIST_API_KEY

Optional:
✓ N8N_API_KEY
✓ GITHUB_TOKEN

============================================================
CRITICAL FILES
============================================================

✓ API endpoints config exists (1234 bytes)
✓ Deployed state tracker exists (5678 bytes)
✓ Letta agent setup script exists (28672 bytes)
✓ Todoist setup script exists (13456 bytes)

============================================================
VALIDATION SUMMARY
============================================================

✓ All checks passed! System is ready.
```

If validation fails, fix issues before proceeding.

---

## Component Deployment

### Deploy Letta Agents (Cloud)

**Option 1: Deploy All 41 Agents**

```bash
cd infrastructure/scripts
python letta_agent_setup.py --all
```

This will:
1. Create all 41 agents in Letta Cloud
2. Configure system prompts from YAML configs
3. Set models (Sonnet 4.5 or Haiku 3)
4. Enable background processing where configured
5. Update `deployed_state.yaml` with agent IDs

Time: ~15-20 minutes

**Option 2: Deploy Single Pillar**

```bash
# Deploy only Command layer (4 agents)
python letta_agent_setup.py --pillar command

# Deploy only HP pillar (8 agents)
python letta_agent_setup.py --pillar hp

# Similarly: le, ge, it
```

Time: ~3-5 minutes per pillar

**Option 3: Deploy Single Agent**

```bash
python letta_agent_setup.py --agent cmd-mother-ai
```

Time: ~1 minute

**Verification**:
```bash
# Check agent count in Letta Cloud
curl https://api.letta.com/agents \
  -H "Authorization: Bearer $LETTA_API_KEY" | jq '. | length'

# Should return: 41 (or however many deployed)

# Test Mother AI
curl -X POST https://api.letta.com/agents/{agent_id}/messages \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages": [{"role": "user", "text": "Hello, are you operational?"}]}'
```

### Deploy Todoist Structure

```bash
cd infrastructure/scripts
python todoist_setup.py --create
```

This will:
1. Create 4 pillar projects (HP, LE, GE, IT)
2. Create Command/Meta project
3. Set up sections within each project
4. Configure colors and favorites
5. Create initial task templates

**Verification**:
```bash
# List all projects
python todoist_setup.py --list

# Or check in Todoist UI
open https://todoist.com
```

Expected projects:
- `00-09 Command (Meta)`
- `10-19 HP - Holistic Performance`
- `20-29 LE - Legacy Evolution`
- `30-39 GE - Generational Empire`
- `40-49 IT - Innovative Technology`

### Deploy N8N Workflows

**Prerequisites**: N8N must be running (see [Deploy Self-Hosted Services](#deploy-self-hosted-services))

```bash
cd infrastructure
./scripts/deploy_workflows.sh
```

**Note**: Current script is a placeholder. Full deployment requires N8N API integration.

**Manual Alternative**:
1. Open N8N: https://n8n.zeusops.dev
2. For each workflow in `workflows/production/`:
   - Click **Import**
   - Upload JSON file
   - Configure credentials (Letta API, Todoist API, etc.)
   - Activate workflow

Workflows to deploy:
- `one-mind-morning-briefing.json`
- `one-mind-evening-review.json`
- `one-mind-agent-health-check.json`
- `one-mind-weekly-review.json`
- `one-mind-memory-consolidation.json`
- `one-mind-approval-webhook.json`

### Deploy Self-Hosted Services

**On VM or local Docker**:

```bash
# SSH to VM
ssh user@zeusops.dev

# Clone infrastructure
git clone https://github.com/your-username/one-mind.git
cd one-mind/infrastructure

# Create .env file
vim .env
# (Add same credentials as local)

# Deploy all services
docker-compose up -d

# Check status
docker-compose ps
```

**Services Deployed**:
- **N8N** - Workflow automation (port 5678)
- **Prometheus** - Metrics collection (port 9090)
- **Grafana** - Metrics visualization (port 3000)
- **Uptime Kuma** - Uptime monitoring (port 3001)
- **Langfuse** - LLM observability (port 3002)

**Access Services**:
- N8N: https://n8n.zeusops.dev
- Grafana: https://grafana.zeusops.dev
- Uptime Kuma: https://uptime.zeusops.dev
- Langfuse: https://langfuse.zeusops.dev

**Initial Setup for N8N**:
1. Navigate to https://n8n.zeusops.dev
2. Create admin account (credentials from `.env`)
3. Go to **Credentials** → Add credentials for:
   - Letta HTTP API
   - Todoist API
   - GitHub API
4. Import workflows (see above)

---

## Production Deployment

### Full System Deployment (Fresh Start)

For deploying the complete One Mind system from scratch:

```bash
# 1. Validate environment
cd infrastructure
python scripts/validate_system.py

# 2. Deploy Letta agents (cloud)
cd scripts
python letta_agent_setup.py --all

# 3. Deploy Todoist structure (cloud)
python todoist_setup.py --create

# 4. Deploy self-hosted services (VM)
ssh user@zeusops.dev
cd one-mind/infrastructure
docker-compose up -d

# 5. Deploy N8N workflows (manual import for now)
# Open N8N UI and import each workflow

# 6. Final validation (back on local)
cd infrastructure
python scripts/validate_system.py

# 7. Run health check
# N8N workflow: one-mind-agent-health-check

# 8. Enable continuous maintenance
# Activate all N8N scheduled workflows
```

**Total Time**: ~45-60 minutes

### GitOps Deployment (Updates)

For updating an existing deployment via Git:

```bash
# 1. Make changes locally
vim config/agents/command.yaml  # Edit agent configs

# 2. Commit to Git
git add config/
git commit -m "Update Mother AI system prompt"
git push

# 3. Pull on VM
ssh user@zeusops.dev
cd one-mind/infrastructure
git pull

# 4. Sync agents to updated configs
python scripts/sync_agent_configs.py --update-all

# 5. Restart services if needed
docker-compose restart

# 6. Verify
python scripts/validate_system.py
```

### Staged Deployment (Test → Production)

For testing changes before full deployment:

```bash
# 1. Create test branch
git checkout -b test/new-agent-prompts

# 2. Make changes
vim config/agents/hp.yaml

# 3. Deploy to test agent only
python scripts/sync_agent_configs.py --update hp-coordinator --dry-run
python scripts/sync_agent_configs.py --update hp-coordinator

# 4. Test thoroughly
# Send test messages, verify behavior

# 5. If successful, merge to main
git checkout main
git merge test/new-agent-prompts
git push

# 6. Deploy to all agents
python scripts/sync_agent_configs.py --update-all
```

---

## Update Procedures

### Update Agent Configurations

```bash
# 1. Edit YAML configs
vim config/agents/command.yaml

# 2. Validate
python scripts/validate_system.py

# 3. Dry-run to preview changes
python scripts/sync_agent_configs.py --dry-run

# 4. Apply updates
python scripts/sync_agent_configs.py --update-all

# 5. Verify
python scripts/sync_agent_configs.py --check
```

### Update Workflows

```bash
# 1. Edit workflow in N8N UI
# Or edit JSON: vim workflows/production/one-mind-morning-briefing.json

# 2. Export from N8N (if edited in UI)
# Download → Save to workflows/production/

# 3. Deploy updated workflow
./scripts/deploy_workflows.sh

# 4. Test execution
# Trigger workflow manually in N8N
```

### Update Self-Hosted Services

```bash
# SSH to VM
ssh user@zeusops.dev
cd one-mind/infrastructure

# Pull latest changes
git pull

# Update Docker images
docker-compose pull

# Recreate containers with new images
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f n8n
```

### Update Python Scripts

```bash
# 1. Pull latest code
git pull

# 2. Update dependencies (if requirements.txt changed)
pip install -r requirements.txt --upgrade

# 3. Test updated scripts
python scripts/validate_system.py
```

---

## Verification and Testing

### Post-Deployment Checklist

After any deployment:

- [ ] **Validation passes**
  ```bash
  python scripts/validate_system.py
  # Exit code: 0 (success)
  ```

- [ ] **All agents respond**
  ```bash
  # Check agent count
  curl https://api.letta.com/agents \
    -H "Authorization: Bearer $LETTA_API_KEY" | jq '. | length'
  # Should return: 41
  ```

- [ ] **Todoist projects exist**
  ```bash
  python scripts/todoist_setup.py --list
  # Should list all 5 projects
  ```

- [ ] **N8N workflows active**
  - Open N8N UI
  - Verify all workflows show as "Active"

- [ ] **Docker services running**
  ```bash
  ssh user@zeusops.dev
  docker-compose ps
  # All services: Up
  ```

- [ ] **Test critical workflows**
  - Manually trigger morning briefing
  - Verify health check runs
  - Test agent message sending

### Smoke Tests

**Test 1: Agent Communication**
```bash
# Send test message to Mother AI
MOTHER_AI_ID=$(grep -A 1 "cmd-mother-ai:" config/deployed_state.yaml | grep letta_id | cut -d'"' -f2)

curl -X POST https://api.letta.com/agents/$MOTHER_AI_ID/messages \
  -H "Authorization: Bearer $LETTA_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"messages": [{"role": "user", "text": "System test: Please confirm you are operational"}]}'

# Expected: Response with confirmation
```

**Test 2: Todoist Integration**
```bash
# Create test task
curl -X POST https://api.todoist.com/rest/v2/tasks \
  -H "Authorization: Bearer $TODOIST_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"content": "Deployment test task"}'

# Expected: HTTP 200, task created

# Delete test task
curl -X DELETE https://api.todoist.com/rest/v2/tasks/{task_id} \
  -H "Authorization: Bearer $TODOIST_API_KEY"
```

**Test 3: N8N Execution**
- Open N8N: https://n8n.zeusops.dev
- Select `one-mind-agent-health-check` workflow
- Click **Execute Workflow**
- Verify: Executes successfully, no errors

**Test 4: Config Sync**
```bash
# Check for drift (should be none after fresh deployment)
python scripts/sync_agent_configs.py --check

# Expected output:
# ✓ All agents in sync (0 drifted)
```

### Load Testing (Optional)

Test system under load:

```bash
# Send 10 messages to different agents in parallel
for i in {1..10}; do
  agent_id=$(get_random_agent_id)
  curl -X POST https://api.letta.com/agents/$agent_id/messages \
    -H "Authorization: Bearer $LETTA_API_KEY" \
    -H "Content-Type: application/json" \
    -d '{"messages": [{"role": "user", "text": "Load test message '$i'"}]}' &
done
wait

# Check all agents still responsive
python scripts/validate_system.py
```

---

## Troubleshooting

### Issue: Validation Fails

**Symptom**: `validate_system.py` exits with errors

**Common Causes**:
1. Missing API keys
2. API endpoints unreachable
3. Invalid YAML syntax

**Solution**:
```bash
# Check environment variables
env | grep -E "(LETTA|TODOIST|N8N|GITHUB)"

# Test API connectivity manually
curl https://api.letta.com/agents \
  -H "Authorization: Bearer $LETTA_API_KEY"

# Validate YAML files
python -c "import yaml; yaml.safe_load(open('config/agents/command.yaml'))"
```

### Issue: Agent Deployment Fails

**Symptom**: `letta_agent_setup.py` errors during agent creation

**Common Causes**:
1. Invalid Letta API key
2. Quota/billing issue
3. Malformed system prompt

**Solution**:
```bash
# Verify API key works
curl https://api.letta.com/agents \
  -H "Authorization: Bearer $LETTA_API_KEY"

# Check Letta account status
# Visit: https://app.letta.com/billing

# Test single agent creation
python letta_agent_setup.py --agent cmd-mother-ai

# Check logs for specific error
cat logs/letta_setup.log
```

### Issue: Config Sync Shows Drift

**Symptom**: `sync_agent_configs.py --check` reports drift immediately after deployment

**Common Causes**:
1. YAML config doesn't match deployed state
2. Agent was manually edited in Letta UI
3. Deployed state tracker out of sync

**Solution**:
```bash
# View exact differences
python scripts/sync_agent_configs.py --check

# If drift is expected (manual edits), sync configs:
python scripts/sync_agent_configs.py --update-all

# If drift is unexpected, investigate:
# - Check Letta UI for manual changes
# - Review deployed_state.yaml for correct agent IDs
```

### Issue: N8N Workflows Not Executing

**Symptom**: Workflows don't run on schedule or fail execution

**Common Causes**:
1. Workflow not activated
2. Invalid credentials
3. Missing environment variables

**Solution**:
```bash
# Check N8N logs
ssh user@zeusops.dev
docker-compose logs n8n | tail -100

# Verify workflow is active
# In N8N UI: Workflows should show green "Active" badge

# Test credentials
# In N8N UI: Credentials → Test each credential

# Check environment variables
docker-compose exec n8n env | grep -E "(LETTA|TODOIST)"
```

### Issue: Docker Services Won't Start

**Symptom**: `docker-compose up -d` fails or containers exit

**Common Causes**:
1. Port conflicts
2. Missing environment variables
3. Insufficient disk space

**Solution**:
```bash
# Check port conflicts
sudo lsof -i :5678  # N8N
sudo lsof -i :3000  # Grafana

# Check disk space
df -h

# Check Docker logs
docker-compose logs n8n
docker-compose logs prometheus

# Remove old containers
docker-compose down
docker system prune -a

# Recreate
docker-compose up -d
```

### Issue: API Rate Limits

**Symptom**: API calls return 429 (Too Many Requests)

**Solution**:
```bash
# Add delays between bulk operations
# Edit sync_agent_configs.py:
import time
time.sleep(1)  # Add 1 second delay between agents

# Check rate limit headers
curl -I https://api.letta.com/agents \
  -H "Authorization: Bearer $LETTA_API_KEY" | grep -i rate

# Wait for rate limit reset
# Header: X-RateLimit-Reset (Unix timestamp)
```

---

## Rollback Procedures

### Rollback Agent Configs

```bash
# 1. Restore from backup
./scripts/restore.sh 2024-12-22

# 2. Sync agents to previous configs
python scripts/sync_agent_configs.py --update-all

# 3. Verify
python scripts/validate_system.py
```

### Rollback Self-Hosted Services

```bash
# SSH to VM
ssh user@zeusops.dev
cd one-mind/infrastructure

# 1. Checkout previous version
git log --oneline  # Find commit hash
git checkout abc123  # Previous working commit

# 2. Recreate containers
docker-compose down
docker-compose up -d

# 3. Verify
docker-compose ps
```

---

## Deployment Automation (Future)

### CI/CD Pipeline (Planned)

```yaml
# .github/workflows/deploy.yml
name: Deploy One Mind

on:
  push:
    branches: [main]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate configs
        run: python scripts/validate_system.py

  deploy-agents:
    needs: validate
    runs-on: ubuntu-latest
    steps:
      - name: Sync agent configs
        run: python scripts/sync_agent_configs.py --update-all
        env:
          LETTA_API_KEY: ${{ secrets.LETTA_API_KEY }}

  deploy-workflows:
    needs: validate
    runs-on: ubuntu-latest
    steps:
      - name: Deploy N8N workflows
        run: ./scripts/deploy_workflows.sh
        env:
          N8N_API_KEY: ${{ secrets.N8N_API_KEY }}
```

---

## Best Practices

1. **Always validate before deployment**
   ```bash
   python scripts/validate_system.py
   ```

2. **Use dry-run for bulk changes**
   ```bash
   python scripts/sync_agent_configs.py --dry-run
   ```

3. **Backup before major changes**
   ```bash
   ./scripts/backup.sh
   ```

4. **Deploy during low-usage times**
   - Early morning (5-6 AM)
   - Minimize disruption

5. **Test in stages**
   - Single agent → Pillar → All agents

6. **Monitor for 24 hours after deployment**
   - Check health checks
   - Review error logs
   - Verify workflow execution

7. **Document all changes**
   - Clear commit messages
   - Update CHANGELOG.md
   - Note any manual steps

8. **Keep secrets secure**
   - Never commit `.env`
   - Rotate API keys quarterly
   - Use different keys for test/prod

---

## Next Steps

After successful deployment:

1. **Enable continuous maintenance**
   - See [CONTINUOUS-MAINTENANCE.md](00-24%20SO%20(Sovereign%20Operations)/05%20Development%20(Code%20&%20Testing)/DevOps/CONTINUOUS-MAINTENANCE.md)
   - Activate scheduled N8N workflows

2. **Set up monitoring**
   - Configure Grafana dashboards
   - Set up Uptime Kuma checks
   - Enable Langfuse tracing

3. **Test bulk operations**
   - See [BULK-OPERATIONS.md](00-24%20SO%20(Sovereign%20Operations)/02%20Infrastructure%20(VPS,%20Servers%20&%20Networking)%20(The%20Body)/Tech-Stack/BULK-OPERATIONS.md)
   - Practice emergency rollback

4. **Review API access patterns**
   - See [API-ACCESS.md](00-24%20SO%20(Sovereign%20Operations)/05%20Development%20(Code%20&%20Testing)/Cloud/API-ACCESS.md)
   - Optimize API usage

---

## Support

- **Issues**: https://github.com/your-username/one-mind/issues
- **Documentation**: `/infrastructure/docs/`
- **Logs**: `/infrastructure/logs/`
