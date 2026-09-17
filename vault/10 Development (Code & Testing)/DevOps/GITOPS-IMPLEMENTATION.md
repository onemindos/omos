---
id: ui-12-20260111053926
title: "🚀 GitOps Implementation Guide"
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
ai_summary: "🚀 GitOps Implementation Guide"
---

# 🚀 GitOps Implementation Guide

## Overview

This document outlines the complete GitOps implementation for Legacy AI development, providing automated deployment, rollback, and environment management through GitHub Actions.

## 📋 Workflow Inventory

### 1. **🚀 Automated Deployment** (`deploy-production.yml`)
- **Trigger**: Push to main branch
- **Purpose**: Automated deployment of changes to production VM
- **Features**:
  - Path-based change detection
  - Multi-service deployment (Core → Automation → IoT → Infrastructure)
  - Health checks with retry logic
  - Automatic notifications via n8n webhooks

### 2. **🔧 Manual Service Control** (`manual-control.yml`)
- **Trigger**: Manual workflow dispatch
- **Purpose**: Administrative control over services
- **Actions**:
  - Restart services
  - View logs
  - System health checks
  - Emergency cleanup

### 3. **🔄 Emergency Rollback** (`rollback.yml`)
- **Trigger**: Manual workflow dispatch (emergency use)
- **Purpose**: Rapid rollback to previous stable state
- **Features**:
  - Automatic pre-rollback backup
  - Git-based rollback to specific commit
  - Service restart with rolled-back configuration
  - Post-rollback health verification

### 4. **🌍 Environment Management** (`environment.yml`)
- **Trigger**: Manual workflow dispatch
- **Purpose**: System maintenance and monitoring
- **Operations**:
  - Full system backups
  - Configuration-only backups
  - System cleanup
  - Resource reporting
  - Maintenance mode toggle

## 🔧 Setup Requirements

### GitHub Secrets Configuration

The following secrets must be configured in your GitHub repository:

```bash
# Required for VM access
VM_PRIVATE_KEY=<your-ssh-private-key-content>

# Optional but recommended for enhanced features
PORTAINER_TOKEN=<portainer-api-token>
SLACK_WEBHOOK=<slack-webhook-url>
```

### VM Prerequisites

1. **SSH Key Access**: GitHub Actions must have SSH access to your VM
2. **Docker Environment**: All services must be containerized
3. **n8n Webhook Endpoints**: For deployment notifications

## 📊 Change Detection Logic

The deployment workflow uses intelligent path-based detection:

```yaml
Core Services (stacks/legacy-ai/):
- letta/ → Triggers letta service deployment
- memory/ → Triggers memory service deployment
- llm/ → Triggers LLM service deployment

Automation (stacks/apps/):
- n8n-queue-mode/ → Triggers n8n deployment
- n8n/ → Triggers n8n deployment

IoT Stack (stacks/iot/):
- home-assistant/ → Triggers Home Assistant deployment
- node-red/ → Triggers Node-RED deployment
- emqx/ → Triggers EMQX deployment

Infrastructure:
- scripts/ → Triggers infrastructure updates
- .github/workflows/ → Triggers workflow updates
```

## 🚀 Deployment Process

### Automatic Deployment Flow

1. **Code Push** → Repository receives changes
2. **Change Detection** → Workflow analyzes modified paths
3. **Service Selection** → Determines which services need updates
4. **Sequential Deployment**:
   - Core services (critical infrastructure)
   - Automation services (n8n)
   - IoT services (Home Assistant, Node-RED, EMQX)
   - Infrastructure updates
5. **Health Verification** → Confirms all services are healthy
6. **Notification** → Reports deployment status

### Manual Override Capabilities

```bash
# Restart specific service
GitHub Actions → Manual Control → restart-service → Select service

# Emergency rollback
GitHub Actions → Rollback → Specify commit/reason → CONFIRM

# System maintenance
GitHub Actions → Environment → backup-full/cleanup-system
```

## 🔄 Rollback Strategy

### When to Use Rollback

- Service failures after deployment
- Configuration errors causing instability
- Performance degradation
- Security issues requiring immediate reversion

### Rollback Process

1. **Confirmation Required**: Type "CONFIRM" to proceed
2. **Pre-Rollback Backup**: Automatic backup of current state
3. **Service Shutdown**: Graceful stop of all services
4. **Git Rollback**: Repository reset to target commit
5. **Service Restart**: Start services with rolled-back config
6. **Health Check**: Verify rollback success

### Rollback Safety

- Pre-rollback backups stored in `/opt/backups/pre-rollback-*`
- Docker volume backups included
- Container state preservation
- Automatic notification of rollback status

## 📊 Monitoring & Alerts

### Health Check Endpoints

The workflows monitor these key services:

- **Letta Agent Framework**: `http://localhost:8283`
- **n8n Automation**: `http://localhost:5678`
- **Dify AI Platform**: `http://localhost:3000`
- **Home Assistant**: `http://localhost:8123`
- **Node-RED**: `http://localhost:1880`

### Notification Integration

Automatic notifications sent to n8n webhooks:

```bash
# Deployment notifications
POST http://localhost:5678/webhook/deployment-status

# Rollback notifications  
POST http://localhost:5678/webhook/rollback-notification

# Health check alerts
POST http://localhost:5678/webhook/health-alerts
```

## 🛠 Troubleshooting

### Common Issues

**1. SSH Connection Failures**
```bash
Solution: Verify VM_PRIVATE_KEY secret is correctly configured
Check: SSH key format and permissions
```

**2. Service Start Failures**
```bash
Solution: Check Docker Compose files for syntax errors
Check: Resource availability (memory, disk space)
```

**3. Health Check Failures**
```bash
Solution: Increase health check timeout in workflow
Check: Service startup time requirements
```

### Debug Commands

```bash
# Check workflow run logs in GitHub Actions
# SSH into VM for manual inspection
ssh ubuntu@3.215.29.128

# Check service status
docker ps -a

# View service logs
docker logs <container-name>

# Check resource usage
df -h && free -h
```

## 📈 Performance Optimization

### Deployment Speed

- **Parallel Service Starts**: Where dependencies allow
- **Incremental Updates**: Only affected services are redeployed
- **Health Check Optimization**: Configurable timeouts and retries

### Resource Management

- **Memory Monitoring**: Automatic resource reports
- **Cleanup Automation**: Regular cleanup of old backups and unused images
- **Backup Rotation**: Intelligent backup retention policies

## 🔮 Future Enhancements

### Phase 2 Planned Features

1. **Blue-Green Deployments**: Zero-downtime deployments
2. **Canary Releases**: Gradual rollout with automatic rollback
3. **Multi-Environment Support**: Staging/Production pipelines
4. **Advanced Monitoring**: Prometheus/Grafana integration
5. **Security Scanning**: Automatic vulnerability assessment

### Integration Roadmap

- **Slack Integration**: Real-time deployment notifications
- **Portainer API**: Enhanced container management
- **Database Migrations**: Automated schema updates
- **Performance Testing**: Automated performance validation

## 📚 Best Practices

### Development Workflow

1. **Feature Branches**: Develop in feature branches
2. **Pull Requests**: All changes via PR with review
3. **Testing**: Local testing before merge to main
4. **Documentation**: Update this guide with changes

### Emergency Procedures

1. **Immediate Rollback**: Use rollback workflow for critical issues
2. **Service Isolation**: Use manual control to isolate problematic services
3. **Backup Verification**: Regular backup testing
4. **Incident Documentation**: Document all emergency responses

### Security Guidelines

1. **Secret Management**: Use GitHub Secrets for all sensitive data
2. **SSH Key Rotation**: Regular SSH key updates
3. **Access Control**: Limit workflow execution permissions
4. **Audit Logging**: Monitor all deployment activities

---

## 🎯 Quick Reference Commands

### GitHub Actions Triggers

```bash
# Automatic deployment
git push origin main

# Manual service restart
Actions → Manual Control → restart-service

# Emergency rollback
Actions → Rollback → [commit] → CONFIRM

# System backup
Actions → Environment → backup-full

# Resource report
Actions → Environment → resource-report
```

### Direct VM Commands (Emergency Use)

```bash
# Service status
cd /home/ubuntu/legacy-ai && docker ps

# Manual service restart
docker-compose -f stacks/legacy-ai/docker-compose.yml restart

# View deployment logs
tail -f /var/log/deployment.log

# Emergency service stop
docker stop $(docker ps -q)
```

This GitOps implementation provides a robust, automated deployment pipeline with comprehensive safety measures and monitoring capabilities for the Legacy AI hybrid architecture.