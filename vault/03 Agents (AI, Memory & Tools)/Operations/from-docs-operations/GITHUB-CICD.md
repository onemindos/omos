# GitHub Actions CI/CD Pipeline

This repository uses GitHub Actions for automated continuous integration and deployment.

## Workflows

### CI Workflow (`.github/workflows/ci.yml`)

Triggers on:
- Push to `main` or `develop` branches
- Pull requests to `main`
- Manual workflow dispatch

**Jobs:**

1. **Detect Changes**
   - Uses path filters to determine which checks to run
   - Checks: `ui/**`, `agents/**`, `infra/containers/**`, `infra/compose/**`

2. **TypeScript Type Check & Lint**
   - Type checks Next.js UI
   - Runs ESLint
   - Verifies production build

3. **Python Type Check & Lint**
   - Lints Python code with `ruff`
   - Type checks with `mypy`

4. **Python Tests**
   - Runs pytest on `tests/` directory

5. **Docker Build**
   - Builds all Docker images to verify they compile
   - Services: agno, gateway, events, actions

### Deploy Workflow (`.github/workflows/deploy.yml`)

Triggers on:
- Push to `main` branch (paths: `agents/**`, `infra/**`, `ui/**`)
- Manual workflow dispatch

**Jobs:**

1. **Notify Start**
   - Sends Telegram notification when deployment begins

2. **Deploy to AWS**
   - SSH into Azure VM server
   - Pulls latest code from main branch
   - Rebuilds Docker images on server
   - Restarts services with docker compose
   - Performs health check
   - Auto-rollback on failure

3. **Notify Complete**
   - Sends Telegram notification with deployment result

## Required GitHub Secrets

Configure these secrets in your GitHub repository settings:

### Repository Secrets

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `AWS_SSH_PRIVATE_KEY` | SSH private key for Azure VM server | Contents of `~/.ssh/onemind-key.pem` |
| `AWS_HOST` | Azure VM server IP or hostname | `20.121.38.186` |
| `AWS_USER` | SSH username for AWS server | `ubuntu` |
| `TELEGRAM_BOT_TOKEN` | Telegram bot token for notifications | `123456:ABC...` |
| `TELEGRAM_CHAT_ID` | Telegram chat ID for notifications | `-1001234567890` |
| `LIVEKIT_URL` | LiveKit server URL | `wss://livekit.example.com` |
| `LIVEKIT_API_KEY` | LiveKit API key | `APIKey123` |
| `LIVEKIT_API_SECRET` | LiveKit API secret | `secret...` |
| `DEEPGRAM_API_KEY` | Deepgram API key for voice | `dg-...` |

### How to Add Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret with the name and value from the table above

### SSH Key Setup

The SSH private key should be the same key used to access your Azure VM server:

```bash
# Display your SSH private key (copy this value)
cat ~/.ssh/onemind-key.pem
```

**Important:** Ensure the key includes the `-----BEGIN RSA PRIVATE KEY-----` and `-----END RSA PRIVATE KEY-----` lines.

## Docker Images

Docker images are built locally on the server. Previously pushed images to GHCR:
- `ghcr.io/zeus-delacruz/onemind-agno:latest`
- `ghcr.io/zeus-delacruz/onemind-gateway:latest`
- `ghcr.io/zeus-delacruz/onemind-events:latest`
- `ghcr.io/zeus-delacruz/onemind-actions:latest`

## Manual Deployment

You can manually trigger a deployment:

1. Go to **Actions** tab in GitHub
2. Select **Deploy** workflow
3. Click **Run workflow**
4. Select the `main` branch
5. Click **Run workflow**

## Deployment Architecture

```
GitHub Push (main branch)
    ↓
SSH to Azure VM
    ↓
Pull latest code (git reset --hard origin/main)
    ↓
Build images on server (docker compose build)
    ↓
Start services (docker compose up -d)
    ↓
Health check (essential services)
    ↓
Telegram notification
```

## Health Checks

The deployment checks these essential services:
- `onemind-nats`
- `onemind-postgres`
- `onemind-agno`
- `onemind-caddy`
- `onemind-gateway`

Services must be running and healthy before deployment is considered successful.

## Troubleshooting

### CI Failures

**TypeScript errors:**
```bash
cd ui
pnpm install
pnpm exec tsc --noEmit
pnpm lint
```

**Python linting errors:**
```bash
pip install ruff mypy
ruff check agents/
mypy agents/ --ignore-missing-imports
```

**Python test failures:**
```bash
pip install -r requirements-dev.txt
pytest tests/ -v
```

### Deployment Failures

**SSH connection issues:**
- Verify `AWS_SSH_PRIVATE_KEY` secret is correct
- Ensure `AWS_HOST` and `AWS_USER` are correct
- Check that the server's security group allows SSH from GitHub Actions IPs

**Docker build failures:**
- Check Docker logs on the server
- Verify all required environment variables are set in `.env` on the server

**Service health check failures:**
- SSH into the server and check service logs:
  ```bash
  ssh ubuntu@20.121.38.186
  cd ~/onemind-os/infra/compose/cloud
  docker compose logs -f
  ```

**Rollback:**
The deployment automatically rolls back to the previous commit if health checks fail.

## Local Testing

Test the workflows locally before pushing:

**TypeScript:**
```bash
cd ui
pnpm install
pnpm exec tsc --noEmit
pnpm lint
pnpm build
```

**Python:**
```bash
pip install ruff mypy
ruff check agents/
mypy agents/ --ignore-missing-imports
```

**Docker:**
```bash
# Test building individual images
docker build -f infra/containers/Dockerfile.agno -t onemind/agno:test .
docker build -f infra/containers/Dockerfile.gateway -t onemind/gateway:test .
docker build -f infra/containers/Dockerfile.events -t onemind/events:test .
docker build -f infra/containers/Dockerfile.actions -t onemind/actions:test .
```

## Status Badge

Add this to your README.md to show CI status:

```markdown
![CI](https://github.com/Zeus-Delacruz/onemind-os/workflows/CI/badge.svg)
![Deploy](https://github.com/Zeus-Delacruz/onemind-os/workflows/Deploy/badge.svg)
```
