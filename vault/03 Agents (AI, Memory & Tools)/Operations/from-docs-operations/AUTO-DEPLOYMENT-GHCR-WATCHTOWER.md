# Auto-Deployment with GHCR + Watchtower

## Overview

OneMind OS uses an automated deployment pipeline that eliminates manual intervention. When you push code to `main`, GitHub Actions builds Docker images, pushes them to GitHub Container Registry (GHCR), and Watchtower automatically deploys them to your Azure VM.

**Total deployment time:** ~9 minutes from `git push` to live deployment

```
┌─────────────┐     ┌──────────────┐     ┌────────────┐     ┌─────────────┐
│  git push   │ --> │ GitHub Actions│ --> │    GHCR    │ --> │ Watchtower  │
│    main     │     │  (~4 minutes) │     │  (registry)│     │ (~5 minutes)│
└─────────────┘     └──────────────┘     └────────────┘     └─────────────┘
                           │                                        │
                           │ Build & Push                           │ Pull & Deploy
                           ▼                                        ▼
                    Docker Images                            Azure VM Services
```

## Architecture

### Components

1. **GitHub Actions** ([.github/workflows/deploy-ghcr.yml](../../../.github/workflows/deploy-ghcr.yml))
   - Triggers on every push to `main`
   - Builds two Docker images in parallel:
     - `ghcr.io/zeus-delacruz/onemind-agno:latest` (Backend)
     - `ghcr.io/zeus-delacruz/onemind-flutter:latest` (Frontend)
   - Pushes images to GHCR with automatic authentication

2. **GitHub Container Registry (GHCR)**
   - Private Docker registry hosted by GitHub
   - Stores versioned container images
   - Integrated with GitHub authentication

3. **Watchtower** (Running on Azure VM)
   - Monitors GHCR for new image versions
   - Polls every 5 minutes for updates
   - Automatically pulls new images and restarts containers
   - Performs rolling restarts (one service at a time)
   - Cleans up old images to save disk space

4. **Azure VM** (20.121.38.186)
   - Ubuntu 22.04 server
   - Runs Docker Compose services
   - Auto-updates via Watchtower

## How It Works

### 1. Developer Workflow

```bash
# Make your changes
vim backend/agno/some_file.py

# Commit and push
git add .
git commit -m "feat: add new feature"
git push origin main
```

That's it! No manual deployment commands needed.

### 2. GitHub Actions Build (4 minutes)

The [deploy-ghcr.yml](../../../.github/workflows/deploy-ghcr.yml) workflow automatically:

1. **Checks out code** from the repository
2. **Sets up Docker Buildx** for optimized multi-platform builds
3. **Authenticates to GHCR** using `GITHUB_TOKEN`
4. **Builds backend image**:
   - Context: Project root
   - Dockerfile: `infra/containers/agno/Dockerfile`
   - Tags: `ghcr.io/zeus-delacruz/onemind-agno:latest`
5. **Builds Flutter web app**:
   - Sets up Flutter 3.27.1
   - Runs `flutter pub get`
   - Builds web app with `flutter build web --release --no-tree-shake-icons`
6. **Builds frontend image**:
   - Context: Project root
   - Dockerfile: `infra/containers/flutter/Dockerfile`
   - Tags: `ghcr.io/zeus-delacruz/onemind-flutter:latest`
7. **Pushes both images to GHCR**

### 3. Watchtower Polling (5 minutes)

Watchtower runs on your Azure VM and:

1. **Checks GHCR every 5 minutes** for new image versions
2. **Compares digests** between local and remote images
3. **Pulls new images** when updates are detected
4. **Performs rolling restart**:
   - Stops container gracefully (SIGTERM)
   - Removes old container
   - Creates new container with updated image
   - Starts new container
   - Waits for health check before moving to next service
5. **Cleans up old images** to free disk space

### 4. Service Availability

- **Backend (Agno)**: http://20.121.38.186:8081
- **Frontend (Flutter)**: http://20.121.38.186:3000
- **Zero-downtime**: Rolling restart ensures services remain available

## Configuration Files

### GitHub Actions Workflow

Location: [.github/workflows/deploy-ghcr.yml](../../../.github/workflows/deploy-ghcr.yml)

Key configuration:
```yaml
name: Build and Push to GHCR

on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  packages: write  # Required for GHCR push

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    steps:
      # Build backend, frontend, push to GHCR
```

### Docker Compose with Watchtower

Location: [infra/deployment/docker/compose/docker-compose.core.ghcr.yml](../../../infra/deployment/docker/compose/docker-compose.core.ghcr.yml)

Key services:
```yaml
agno:
  image: ghcr.io/${GITHUB_REPOSITORY_OWNER:-zeus-delacruz}/onemind-agno:latest
  pull_policy: always
  labels:
    - "com.centurylinklabs.watchtower.enable=true"

flutter-web:
  image: ghcr.io/${GITHUB_REPOSITORY_OWNER:-zeus-delacruz}/onemind-flutter:latest
  pull_policy: always
  labels:
    - "com.centurylinklabs.watchtower.enable=true"

watchtower:
  image: containrrr/watchtower:latest
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
    - /home/zeus/.docker/config.json:/config.json:ro  # For GHCR auth
  environment:
    - WATCHTOWER_CLEANUP=true              # Remove old images
    - WATCHTOWER_POLL_INTERVAL=300         # Check every 5 minutes
    - WATCHTOWER_LABEL_ENABLE=true         # Only update labeled containers
    - WATCHTOWER_ROLLING_RESTART=true      # One at a time
    - DOCKER_API_VERSION=1.45              # API compatibility
  command: --debug --label-enable
```

## Setup Guide

### Initial Setup (Already Completed)

This setup is already complete for your environment, but documented for reference:

#### 1. GHCR Authentication on VM

```bash
# Authenticate Docker to GHCR (one-time)
echo 'YOUR_GITHUB_TOKEN' | docker login ghcr.io -u zeus-delacruz --password-stdin
```

**Token Requirements:**
- Scope: `read:packages`
- Generate at: https://github.com/settings/tokens
- Token used: `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`

#### 2. Environment Variables

On the VM, set `GITHUB_REPOSITORY_OWNER` in `.env`:
```bash
echo "GITHUB_REPOSITORY_OWNER=zeus-delacruz" >> ~/OneMind-OS/.env
```

#### 3. Start Services with Watchtower

```bash
cd ~/OneMind-OS
docker compose -f infra/deployment/docker/compose/docker-compose.core.ghcr.yml up -d
```

### Verification

Check that everything is working:

```bash
# Check service status
ssh zeus@20.121.38.186 'docker ps'

# Check Watchtower logs
ssh zeus@20.121.38.186 'docker logs onemind-watchtower --tail 50'

# Verify images
ssh zeus@20.121.38.186 'docker images | grep onemind'
```

## Monitoring

### Watch Real-Time Deployment

```bash
# Follow Watchtower logs
ssh zeus@20.121.38.186 'docker logs -f onemind-watchtower'

# Watch for these key messages:
# - "Found new ghcr.io/zeus-delacruz/onemind-agno:latest image"
# - "Stopping /onemind-agno"
# - "Creating /onemind-agno"
# - "Session done" with "Updated=2"
```

### Check GitHub Actions

View build progress:
- Dashboard: https://github.com/Zeus-Delacruz/OneMind-OS/actions
- Status: Look for "Build and Push to GHCR" workflow
- Typical duration: 4-8 minutes

### Service Health

```bash
# Check all services
ssh zeus@20.121.38.186 'cd ~/OneMind-OS && docker compose -f infra/deployment/docker/compose/docker-compose.core.ghcr.yml ps'

# Check specific service
ssh zeus@20.121.38.186 'docker inspect onemind-agno | grep -A 5 Health'
```

### Image Versions

```bash
# See image digests and timestamps
ssh zeus@20.121.38.186 'docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedSince}}" | grep onemind'
```

## Troubleshooting

### Issue: Watchtower Can't Pull Images

**Symptoms:**
```
Error response from daemon: error from registry: unauthorized
Unable to update container "/onemind-agno"
```

**Cause:** Watchtower container doesn't have GHCR credentials

**Fix:**
```bash
# 1. Verify Docker login on VM
ssh zeus@20.121.38.186 'cat ~/.docker/config.json'

# 2. Ensure Watchtower mounts the config
# Check docker-compose.core.ghcr.yml has:
#   volumes:
#     - /home/zeus/.docker/config.json:/config.json:ro

# 3. Restart Watchtower
ssh zeus@20.121.38.186 'cd ~/OneMind-OS && docker compose -f infra/deployment/docker/compose/docker-compose.core.ghcr.yml restart watchtower'
```

### Issue: GitHub Actions Build Fails

**Symptoms:** Workflow shows red X, build error in logs

**Common Causes:**

1. **Flutter compilation error**
   ```
   Error: No named parameter with the name 'initialValue'
   ```
   **Fix:** Flutter 3.27.1 breaking change - use `value:` instead of `initialValue:` in DropdownButtonFormField

2. **Icon tree shaking error**
   ```
   This application cannot tree shake icons fonts
   ```
   **Fix:** Add `--no-tree-shake-icons` flag to `flutter build web` command

3. **Docker context issue**
   ```
   failed to read dockerfile: open Dockerfile: no such file or directory
   ```
   **Fix:** Verify `context: .` (project root) and correct `file:` path in workflow

### Issue: Services Not Updating

**Symptoms:** Code changes pushed but VM still runs old version

**Diagnosis:**
```bash
# 1. Check GitHub Actions completed successfully
# Visit: https://github.com/Zeus-Delacruz/OneMind-OS/actions

# 2. Check Watchtower logs for errors
ssh zeus@20.121.38.186 'docker logs onemind-watchtower --tail 100 | grep -i error'

# 3. Verify image digests differ
# Compare local vs remote
ssh zeus@20.121.38.186 'docker images --digests | grep onemind'
# vs GHCR UI: https://github.com/Zeus-Delacruz?tab=packages
```

**Fix:**
```bash
# Manually trigger Watchtower check
ssh zeus@20.121.38.186 'docker exec onemind-watchtower /watchtower --run-once --debug --label-enable'
```

### Issue: Docker API Version Mismatch

**Symptoms:**
```
Error response from daemon: client version 1.25 is too old
Minimum supported API version is 1.44
```

**Fix:** Already configured in `docker-compose.core.ghcr.yml`:
```yaml
watchtower:
  environment:
    - DOCKER_API_VERSION=1.45
```

If missing, add this environment variable and restart Watchtower.

### Issue: Out of Disk Space

**Symptoms:**
```
no space left on device
```

**Cause:** Old Docker images accumulating

**Prevention:** Watchtower automatically cleans up with `WATCHTOWER_CLEANUP=true`

**Manual cleanup:**
```bash
ssh zeus@20.121.38.186 'docker system prune -a -f'
```

### Issue: Service Won't Start After Update

**Symptoms:** Container in "Restarting" state

**Diagnosis:**
```bash
# Check container logs
ssh zeus@20.121.38.186 'docker logs onemind-agno --tail 100'

# Check exit code
ssh zeus@20.121.38.186 'docker inspect onemind-agno | grep -A 5 State'
```

**Rollback:**
```bash
# 1. Find previous working image digest (from GHCR or local backup)
# 2. Temporarily update compose file to pin that digest
# 3. Restart services
ssh zeus@20.121.38.186 'cd ~/OneMind-OS && docker compose -f infra/deployment/docker/compose/docker-compose.core.ghcr.yml up -d agno'
```

## Advanced Configuration

### Adjust Polling Interval

Edit `WATCHTOWER_POLL_INTERVAL` in [docker-compose.core.ghcr.yml](../../../infra/deployment/docker/compose/docker-compose.core.ghcr.yml):

```yaml
watchtower:
  environment:
    - WATCHTOWER_POLL_INTERVAL=300  # Seconds (300 = 5 minutes)
```

Options:
- `180` = 3 minutes (faster deployments, more frequent checks)
- `300` = 5 minutes (balanced, current setting)
- `600` = 10 minutes (less frequent, conserves resources)

### Enable Deployment Notifications

Add Discord/Slack webhook to receive deployment alerts:

1. **Create webhook**:
   - Discord: Server Settings → Integrations → Webhooks
   - Slack: App Directory → Incoming Webhooks

2. **Update compose file**:
```yaml
watchtower:
  environment:
    - WATCHTOWER_NOTIFICATION_URL=discord://webhook_id/webhook_token
```

3. **Restart Watchtower**:
```bash
ssh zeus@20.121.38.186 'cd ~/OneMind-OS && docker compose -f infra/deployment/docker/compose/docker-compose.core.ghcr.yml restart watchtower'
```

### Add Staging Environment

To test changes before production:

1. **Create staging branch**: `staging`
2. **Add workflow trigger**:
```yaml
on:
  push:
    branches: [main, staging]
```
3. **Tag images with branch**:
```yaml
tags: ghcr.io/zeus-delacruz/onemind-agno:${{ github.ref_name }}
```
4. **Run separate Watchtower** with staging image tag

### Pin Image Versions

For more control, use version tags instead of `latest`:

```yaml
agno:
  image: ghcr.io/zeus-delacruz/onemind-agno:v1.2.3
```

Benefits:
- Explicit version control
- Easier rollbacks
- Audit trail

Trade-off: Requires updating compose file for each deployment

## Performance Optimization

### GitHub Actions Cache

Already optimized:
- Docker layer caching via `docker/build-push-action@v5`
- Flutter pub cache via `subosito/flutter-action@v2`

### Parallel Builds

Current setup builds backend and frontend **sequentially**. To parallelize:

```yaml
jobs:
  build-backend:
    # Backend build steps

  build-frontend:
    # Frontend build steps
```

**Benefit:** Reduces build time from ~8 minutes to ~4 minutes

**Trade-off:** Uses more GitHub Actions minutes

### Reduce Image Size

Backend image: 2.92GB (can be optimized)
Frontend image: 135MB (already optimized)

**Backend optimization** (future):
- Use multi-stage builds
- Remove dev dependencies
- Use Alpine base images

## Security Considerations

### GHCR Access Token

- **Current token**: Read-only (`read:packages` scope)
- **Storage**: VM at `/home/zeus/.docker/config.json` (permissions: 600)
- **Rotation**: Regenerate token every 90 days at https://github.com/settings/tokens

### Private Repository

- GHCR images are private (repo is private)
- Requires authentication to pull
- GitHub Actions uses automatic `GITHUB_TOKEN`
- VM uses personal access token

### Container Security

- Watchtower runs with Docker socket access (required)
- Containers run as non-root where possible
- Health checks ensure service availability
- Auto-cleanup prevents disk exhaustion

## Cost Analysis

### GitHub Actions

- **Build time**: ~8 minutes per deployment
- **Cost**: $0 (included in GitHub Free tier for public repos)
- **Private repo cost**: 2,000 minutes/month free, then $0.008/minute

### GHCR Storage

- **Current usage**: ~3GB (2 images)
- **Cost**: Free (500MB free, then $0.25/GB/month)
- **With cleanup**: Watchtower removes old images automatically

### Azure VM

- **VM size**: Standard B2s (2 vCPU, 4GB RAM)
- **Estimated cost**: ~$30/month
- **No additional deployment costs**

## Quick Reference

### Essential Commands

```bash
# Monitor deployment
ssh zeus@20.121.38.186 'docker logs -f onemind-watchtower'

# Check service status
ssh zeus@20.121.38.186 'docker ps'

# Force Watchtower check now
ssh zeus@20.121.38.186 'docker exec onemind-watchtower /watchtower --run-once --debug --label-enable'

# View service logs
ssh zeus@20.121.38.186 'docker logs onemind-agno -f'
ssh zeus@20.121.38.186 'docker logs onemind-flutter-web -f'

# Restart services
ssh zeus@20.121.38.186 'cd ~/OneMind-OS && docker compose -f infra/deployment/docker/compose/docker-compose.core.ghcr.yml restart agno flutter-web'

# Check image versions
ssh zeus@20.121.38.186 'docker images | grep onemind'

# Clean up disk space
ssh zeus@20.121.38.186 'docker system prune -a -f'
```

### URLs

- **GitHub Actions**: https://github.com/Zeus-Delacruz/OneMind-OS/actions
- **GHCR Packages**: https://github.com/Zeus-Delacruz?tab=packages
- **Frontend**: http://20.121.38.186:3000
- **Backend**: http://20.121.38.186:8081
- **Portainer**: http://20.121.38.186:9000 (container management UI)
- **Netdata**: http://20.121.38.186:19999 (system monitoring)

### Local Scripts

Helpful scripts in the repository:

- [scripts/setup_watchtower.sh](../../../scripts/setup_watchtower.sh) - Initial Watchtower setup
- [scripts/make_ghcr_public.sh](../../../scripts/make_ghcr_public.sh) - Guide to make packages public (if needed)

## Related Documentation

- [CI/CD Workflows](./WORKFLOWS-CICD.md) - General CI/CD strategy
- [Docker Compose Configuration](../../../infra/deployment/docker/compose/) - All compose files
- [Infrastructure Reference](./INFRA-REFERENCE.md) - Complete infrastructure overview
- [Deployment Quickstart](./DEPLOYMENT-QUICKSTART.md) - Quick deployment guide

## Changelog

### 2026-01-28 - Initial Implementation

- ✅ Created [.github/workflows/deploy-ghcr.yml](../../../.github/workflows/deploy-ghcr.yml)
- ✅ Created [docker-compose.core.ghcr.yml](../../../infra/deployment/docker/compose/docker-compose.core.ghcr.yml)
- ✅ Configured Watchtower with GHCR authentication
- ✅ Fixed Flutter 3.27.1 compatibility issues
- ✅ Fixed Docker API version mismatch
- ✅ Tested end-to-end deployment pipeline
- ✅ Documented complete system

### Future Enhancements

- [ ] Add deployment notifications (Discord/Slack)
- [ ] Implement staging environment
- [ ] Add version tagging strategy
- [ ] Optimize backend image size
- [ ] Parallelize GitHub Actions builds
- [ ] Add deployment rollback automation
- [ ] Implement blue-green deployment strategy
