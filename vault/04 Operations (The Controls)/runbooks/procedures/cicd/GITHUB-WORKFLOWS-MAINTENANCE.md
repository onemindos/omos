# GitHub Workflows - Maintenance Guide

This document provides a comprehensive reference for all GitHub Actions workflows across the OneMind repositories, including configuration, troubleshooting, and maintenance procedures.

---

## Table of Contents

1. [Repository Overview](#repository-overview)
2. [onemind-os Workflows](#onemind-os-workflows)
3. [onemind-infra Workflows](#onemind-infra-workflows)
4. [OneMind-Codex Workflows](#onemind-codex-workflows)
5. [Required Secrets](#required-secrets)
6. [Troubleshooting](#troubleshooting)
7. [Maintenance Schedule](#maintenance-schedule)
8. [Cost Optimization](#cost-optimization)

---

## Repository Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     OneMind CI/CD Architecture                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   onemind-os (Application)                                          │
│   ├── ci.yml           → Lint, type check, build validation        │
│   ├── security.yml     → CodeQL, Trivy, dependency scanning        │
│   ├── release.yml      → Semantic versioning, changelog            │
│   └── deploy.yml       → Build images, deploy to Azure VM     │
│                                                                     │
│   onemind-infra (Infrastructure)                                    │
│   ├── terraform-validate.yml → Format, lint, security scan         │
│   ├── terraform-plan.yml     → Plan with cost estimation           │
│   ├── drift-detection.yml    → Weekly drift monitoring             │
│   └── backup-postgres.yml    → Daily database backups              │
│                                                                     │
│   OneMind-Codex (Documentation)                                     │
│   ├── docs-quality.yml → Markdown lint, links, spell check         │
│   └── sync.yml         → Daily backup, health checks               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## onemind-os Workflows

### ci.yml - Continuous Integration

**Purpose:** Fast feedback on code quality for every push and PR.

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main`
- Manual dispatch

**Jobs:**
| Job | Duration | Purpose |
|-----|----------|---------|
| `changes` | ~10s | Detect which files changed |
| `typescript-check` | ~2m | TypeScript lint, type check, build |
| `python-check` | ~1m | Ruff, Black, mypy |
| `docker-build` | ~5m | Build all Docker images |
| `ci-success` | ~10s | Summary and PR comment |

**Optimization Features:**
- Path filtering (only runs relevant checks)
- GitHub Actions cache for npm/pip
- Parallel Docker builds

**Configuration:**
```yaml
env:
  NODE_VERSION: '20'
  PYTHON_VERSION: '3.12'
```

---

### security.yml - Security Scanning

**Purpose:** Identify vulnerabilities in code and containers.

**Triggers:**
- Push to `main` or `develop`
- Pull requests to `main`
- Weekly schedule (Sundays 3 AM UTC)
- Manual dispatch

**Jobs:**
| Job | Purpose |
|-----|---------|
| `codeql-analysis` | Static analysis for TypeScript/Python |
| `dependency-scan` | npm audit, pip-audit |
| `container-scan` | Trivy vulnerability scanner |
| `secret-scan` | Gitleaks for exposed secrets |

**Results:** Uploads SARIF files to GitHub Security tab.

---

### release.yml - Automated Releases

**Purpose:** Semantic versioning and changelog generation.

**Triggers:**
- Push to `main` with conventional commit prefixes (`feat`, `fix`)
- Manual dispatch with version selection

**Version Bumping:**
- `feat:` → Minor version (0.X.0)
- `fix:` → Patch version (0.0.X)
- `BREAKING:` or `!:` → Major version (X.0.0)

**Outputs:**
- GitHub Release with changelog
- Docker images tagged with version
- Discord notification (if configured)

---

### deploy.yml - Production Deployment

**Purpose:** Build and deploy to Azure VM.

**Triggers:**
- Push to `main`
- Manual dispatch

**Deployment Flow:**
```
1. Notify start (Discord)
2. Build Docker images (6 services, parallel)
3. Push to GHCR
4. SSH to Azure VM
5. Pull latest code
6. docker compose up -d
7. Health check (6 attempts, 60s timeout)
8. Notify complete
```

**Rollback:** Automatic on health check failure.

**Configuration Options:**
- `skip_build`: Use existing images (faster redeploy)
- `environment`: Target environment (currently production only)

---

## onemind-infra Workflows

> **Note:** These workflows are in `docs/workflows-infra/` and need to be copied to the onemind-infra repository.

### terraform-validate.yml

**Purpose:** Validate Terraform configuration on PRs.

**Jobs:**
| Job | Tools |
|-----|-------|
| `validate` | terraform fmt, init, validate, tflint |
| `security-scan` | tfsec, checkov |

---

### terraform-plan.yml

**Purpose:** Show infrastructure changes with cost estimation.

**Features:**
- Terraform plan output in PR comments
- Cost estimation via Infracost (optional)
- Destructive change warnings
- Approval gate for apply

---

### drift-detection.yml

**Purpose:** Detect when cloud resources drift from Terraform state.

**Schedule:** Mondays at 6 AM UTC

**Actions on drift:**
- Creates GitHub issue
- Sends Discord alert
- Uploads drift report artifact

---

### backup-postgres-enhanced.yml

**Purpose:** Daily PostgreSQL backups with verification.

**Schedule:** Daily at 2 AM UTC

**Features:**
- Backup to S3
- Integrity verification
- Optional restore testing
- 30-day retention (configurable)
- Discord notifications

---

## OneMind-Codex Workflows

### docs-quality.yml

**Purpose:** Ensure documentation quality.

**Triggers:**
- Push/PR with `.md` file changes

**Checks:**
| Check | Tool | Severity |
|-------|------|----------|
| Markdown lint | markdownlint | Warning |
| Link check | lychee | Warning |
| Spell check | cspell | Warning |
| YAML validation | yamllint | Warning |

---

### sync.yml

**Purpose:** Daily vault backup and health monitoring.

**Schedule:** Daily at 4 AM UTC

**Jobs:**
- `backup`: Create archive, upload to S3
- `health-check`: Verify vault structure
- `export-docs`: Optional public doc export

---

## Required Secrets

### onemind-os

| Secret | Required | Purpose |
|--------|----------|---------|
| `AZURE_SSH_PRIVATE_KEY` | Yes | SSH to Azure VM |
| `AZURE_HOST` | Yes | Azure VM IP/hostname |
| `AWS_USER` | Yes | SSH username |
| `DISCORD_WEBHOOK_URL` | No | Notifications |
| `CODECOV_TOKEN` | No | Coverage reporting |

### onemind-infra

| Secret | Required | Purpose |
|--------|----------|---------|
| `AWS_ACCESS_KEY_ID` | Yes | Terraform AWS access |
| `AWS_SECRET_ACCESS_KEY` | Yes | Terraform AWS access |
| `AWS_SSH_PRIVATE_KEY` | Yes | SSH to instances |
| `INFRACOST_API_KEY` | No | Cost estimation |
| `DISCORD_WEBHOOK_URL` | No | Alert notifications |

### OneMind-Codex

| Secret | Required | Purpose |
|--------|----------|---------|
| `AWS_ACCESS_KEY_ID` | No | S3 backup upload |
| `AWS_SECRET_ACCESS_KEY` | No | S3 backup upload |
| `DISCORD_WEBHOOK_URL` | No | Sync notifications |

---

## Troubleshooting

### Common Issues

#### "Unable to resolve action"
The IDE may show false positives. Verify the action exists on GitHub:
```bash
gh api repos/owner/action --jq '.name'
```

#### Docker build failures
1. Check Dockerfile syntax
2. Verify base image availability
3. Review build logs for dependency issues

#### SSH connection failures
1. Verify `AWS_HOST` is correct
2. Check `AWS_SSH_PRIVATE_KEY` format (should include newlines)
3. Ensure security group allows port 22

#### Terraform plan failures
1. Run `terraform init` locally
2. Check AWS credentials
3. Verify state backend is accessible

### Debug Mode

Add to workflow for verbose output:
```yaml
env:
  ACTIONS_STEP_DEBUG: true
```

---

## Maintenance Schedule

### Daily (Automated)
- PostgreSQL backups (2 AM UTC)
- Vault sync (4 AM UTC)

### Weekly (Automated)
- Security scans (Sunday 3 AM UTC)
- Drift detection (Monday 6 AM UTC)

### Monthly (Manual)
- [ ] Review security scan results
- [ ] Update GitHub Actions versions
- [ ] Prune old Docker images
- [ ] Review backup retention

### Quarterly (Manual)
- [ ] Audit secrets and rotate if needed
- [ ] Review workflow costs
- [ ] Update dependencies (Node, Python versions)
- [ ] Test disaster recovery procedures

---

## Cost Optimization

### GitHub Actions Minutes

**Free tier:** 2,000 minutes/month (public repos) or 500 minutes (private)

**Optimization strategies:**
1. Path filtering reduces unnecessary runs
2. Caching speeds up builds (saves ~50% time)
3. Matrix jobs run in parallel (wall-clock faster)
4. Skip CI with `[skip ci]` in commit message

### Estimated Monthly Usage

| Workflow | Runs/Month | Minutes/Run | Total |
|----------|------------|-------------|-------|
| CI | 100 | 8 | 800 |
| Deploy | 20 | 12 | 240 |
| Security | 4 | 15 | 60 |
| Terraform | 30 | 5 | 150 |
| Backup | 30 | 3 | 90 |
| Docs | 50 | 2 | 100 |
| **Total** | | | **~1,440** |

---

## Quick Reference

### Manual Workflow Triggers

```bash
# Trigger workflow via gh CLI
gh workflow run ci.yml --repo Zeus-Delacruz/onemind-os

# Trigger with inputs
gh workflow run deploy.yml --repo Zeus-Delacruz/onemind-os -f skip_build=true

# List recent runs
gh run list --repo Zeus-Delacruz/onemind-os --limit 10

# View run logs
gh run view <run-id> --repo Zeus-Delacruz/onemind-os --log
```

### Workflow Status Badges

Add to README:
```markdown
![CI](https://github.com/Zeus-Delacruz/onemind-os/actions/workflows/ci.yml/badge.svg)
![Deploy](https://github.com/Zeus-Delacruz/onemind-os/actions/workflows/deploy.yml/badge.svg)
![Security](https://github.com/Zeus-Delacruz/onemind-os/actions/workflows/security.yml/badge.svg)
```

---

## Changelog

| Date | Change |
|------|--------|
| 2026-01-18 | Initial workflow suite created |
| | - Added ci.yml enhancements (path filtering, PR comments) |
| | - Added security.yml (CodeQL, Trivy, Gitleaks) |
| | - Added release.yml (semantic versioning) |
| | - Enhanced deploy.yml (notifications, rollback) |
| | - Created Terraform workflows for onemind-infra |
| | - Created docs-quality.yml for OneMind-Codex |
| | - Created sync.yml for vault backup |

---

*Last updated: 2026-01-18*


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
