# GitHub Workflows Improvement Plan

## Executive Summary

This document outlines the comprehensive CI/CD strategy for OneMind's three repositories. The goal is to optimize build times, enhance security, enable faster feedback loops, and automate routine operations.

---

## Current State Analysis

### onemind-os (Main Application)

**Existing Workflows:**
| Workflow | Triggers | Purpose |
|----------|----------|---------|
| `ci.yml` | push (main/develop), PR (main) | TypeScript/Python lint, Docker build |
| `deploy.yml` | push (main), manual | Build images → GHCR → Deploy to Azure VM |

**Gaps Identified:**
- ❌ No test execution (unit/integration/e2e)
- ❌ No security scanning (dependencies, containers)
- ❌ No preview deployments for PRs
- ❌ No semantic versioning/changelog
- ❌ No caching optimization for npm/pip
- ❌ No Lighthouse/performance checks
- ❌ No staging environment gate

### onemind-infra (Infrastructure)

**Existing Workflows:**
| Workflow | Triggers | Purpose |
|----------|----------|---------|
| `backup-postgres.yml` | schedule (2AM UTC) | PostgreSQL backup to S3 |
| `deploy-cloud.yml` | manual | Infrastructure config deployment |

**Gaps Identified:**
- ❌ No Terraform plan on PRs
- ❌ No infrastructure drift detection
- ❌ No cost estimation for changes
- ❌ No security scanning for IaC (tfsec/checkov)
- ❌ No backup restoration testing

### OneMind-Codex (Documentation)

**Existing Workflows:** None

**Gaps Identified:**
- ❌ No markdown linting
- ❌ No link checking
- ❌ No automated sync capabilities

---

## Proposed Workflow Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         OneMind CI/CD Pipeline                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  onemind-os                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │   CI     │→ │ Security │→ │  Build   │→ │  Stage   │→ │  Prod    │  │
│  │ (lint,   │  │ (trivy,  │  │ (docker, │  │ (preview │  │ (deploy, │  │
│  │  test)   │  │  codeql) │  │  cache)  │  │  vercel) │  │  health) │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                                         │
│  onemind-infra                                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐                │
│  │  Validate │→ │ Security │→ │  Plan    │→ │  Apply   │                │
│  │ (format,  │  │ (tfsec,  │  │ (cost,   │  │ (manual  │                │
│  │  lint)    │  │ checkov) │  │  diff)   │  │  gate)   │                │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘                │
│                                                                         │
│  OneMind-Codex                                                          │
│  ┌──────────┐  ┌──────────┐                                            │
│  │  Quality  │→ │  Sync    │                                            │
│  │ (lint,    │  │ (backup, │                                            │
│  │  links)   │  │  export) │                                            │
│  └──────────┘  └──────────┘                                            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Detailed Workflow Specifications

### 1. onemind-os Workflows

#### 1.1 ci.yml (Enhanced)
**Improvements:**
- Add npm/pip caching
- Add test execution step
- Add code coverage reporting
- Parallel job execution optimization

#### 1.2 security.yml (NEW)
**Purpose:** Vulnerability scanning for dependencies and containers
**Triggers:** push, PR, weekly schedule
**Jobs:**
- CodeQL analysis (TypeScript/Python)
- Dependency scanning (npm audit, pip-audit)
- Container scanning (Trivy)
- Secret scanning (Gitleaks)

#### 1.3 preview.yml (NEW)
**Purpose:** Deploy preview environments for PRs
**Triggers:** PR (open, sync)
**Jobs:**
- Build preview image
- Deploy to ephemeral environment
- Run smoke tests
- Comment PR with preview URL
- Cleanup on PR close

#### 1.4 release.yml (NEW)
**Purpose:** Semantic versioning and changelog
**Triggers:** PR merge to main (conventional commits)
**Jobs:**
- Determine version bump (major/minor/patch)
- Generate changelog
- Create GitHub release
- Tag Docker images with version

#### 1.5 deploy.yml (Enhanced)
**Improvements:**
- Add staging deployment step
- Add smoke test verification
- Add rollback capability
- Add deployment notifications (Discord)

### 2. onemind-infra Workflows

#### 2.1 terraform-validate.yml (NEW)
**Purpose:** Validate Terraform on every PR
**Triggers:** PR (terraform/ directory)
**Jobs:**
- terraform fmt check
- terraform validate
- tflint
- tfsec security scan
- checkov compliance scan

#### 2.2 terraform-plan.yml (Enhanced)
**Purpose:** Show infrastructure changes before apply
**Triggers:** PR, manual
**Jobs:**
- terraform plan
- Cost estimation (Infracost)
- Post plan output to PR comment
- Require approval for destructive changes

#### 2.3 drift-detection.yml (NEW)
**Purpose:** Detect configuration drift
**Triggers:** Schedule (weekly)
**Jobs:**
- Compare actual vs desired state
- Alert on drift via Discord
- Create issue for significant drift

#### 2.4 backup-postgres.yml (Enhanced)
**Improvements:**
- Add restoration testing
- Add backup verification
- Enhanced monitoring/alerting

### 3. OneMind-Codex Workflows

#### 3.1 docs-quality.yml (NEW)
**Purpose:** Ensure documentation quality
**Triggers:** push, PR
**Jobs:**
- Markdown linting (markdownlint)
- Spell check (cspell)
- Link checking (lychee)
- YAML/JSON validation

#### 3.2 sync.yml (NEW)
**Purpose:** Automated backup and sync
**Triggers:** Schedule (daily), manual
**Jobs:**
- Sync to backup location
- Export specific docs for publishing

---

## Implementation Priority

### Phase 1: Security & Quality (Week 1)
1. ✅ security.yml for onemind-os
2. ✅ terraform-validate.yml for onemind-infra
3. ✅ docs-quality.yml for OneMind-Codex

### Phase 2: Automation (Week 2)
1. ✅ preview.yml for onemind-os
2. ✅ terraform-plan.yml (enhanced) for onemind-infra
3. ✅ sync.yml for OneMind-Codex

### Phase 3: Release Engineering (Week 3)
1. ✅ release.yml for onemind-os
2. ✅ drift-detection.yml for onemind-infra
3. ✅ deploy.yml enhancements

---

## Required Secrets & Variables

### onemind-os
| Secret | Purpose |
|--------|---------|
| `AWS_SSH_PRIVATE_KEY` | SSH to Azure VM (existing) |
| `AWS_HOST` | Azure VM hostname (existing) |
| `AWS_USER` | SSH user (existing) |
| `DISCORD_WEBHOOK_URL` | Deployment notifications |
| `VERCEL_TOKEN` | Preview deployments (optional) |
| `CODECOV_TOKEN` | Coverage reporting |

### onemind-infra
| Secret | Purpose |
|--------|---------|
| `AWS_ACCESS_KEY_ID` | Terraform AWS access |
| `AWS_SECRET_ACCESS_KEY` | Terraform AWS access |
| `AWS_SSH_PRIVATE_KEY` | SSH to instances (existing) |
| `INFRACOST_API_KEY` | Cost estimation |
| `DISCORD_WEBHOOK_URL` | Alert notifications |

### OneMind-Codex
| Secret | Purpose |
|--------|---------|
| `BACKUP_TOKEN` | Sync automation (if needed) |

---

## Estimated Time Savings

| Task | Before | After | Savings |
|------|--------|-------|---------|
| Security review | 30 min/PR | Automated | 30 min/PR |
| Preview testing | Manual deploy | Auto preview | 20 min/PR |
| Terraform review | Manual check | Auto plan | 15 min/PR |
| Doc quality check | Manual | Automated | 10 min/PR |
| Release management | Manual | Automated | 45 min/release |

**Total estimated savings:** 2+ hours per release cycle

---

## Success Metrics

1. **Build Time**: Reduce from current baseline by 30%
2. **Security**: 0 critical vulnerabilities shipped
3. **Preview Coverage**: 100% of PRs get preview environments
4. **Drift Detection**: Weekly automated checks
5. **Release Automation**: Fully automated changelogs

---

## Next Steps

1. Review and approve this plan
2. Create missing secrets in GitHub
3. Implement workflows in priority order
4. Test each workflow on a feature branch
5. Monitor and iterate

