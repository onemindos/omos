---
title: "CLI — railway (Railway CLI)"
type: reference
status: active
created: 2026-04-11
domain: SO/01
tags: [type/reference, domain/so, topic/cli, topic/railway, topic/deploy]
ai_summary: "Railway CLI quick reference — deploy services, manage environments, run databases, view logs, and inject secrets. Alternative to DigitalOcean/Heroku for fast deployments."
---

# railway — Railway CLI

> Deploy apps and databases instantly. No Dockerfiles required for simple deploys.
> Install: `npm install -g @railway/cli` or `brew install railway`

---

## Auth

```bash
# Login (browser OAuth)
railway login

# Login with API token (headless/VPS)
railway login --browserless

# Set token directly
export RAILWAY_TOKEN=<token>

# Check auth
railway whoami
```

---

## Projects

```bash
# Link current directory to a Railway project
railway link

# Create a new project
railway init

# List all projects
railway list

# Open project in browser
railway open
```

---

## Deploy

```bash
# Deploy from current directory
railway up

# Deploy without detaching (shows logs)
railway up --detach=false

# Deploy a specific service
railway up --service <service-name>

# Deploy from a specific directory
railway up --dir ./backend
```

---

## Services

```bash
# List services in project
railway service list

# Create a new service
railway service create

# Link to a specific service
railway service

# View service metrics
railway status
```

---

## Environment & Variables

```bash
# List environment variables
railway variables

# Set a variable
railway variables --set "KEY=value"
railway variables --set "DATABASE_URL=postgresql://..."

# Set multiple variables
railway variables --set "KEY1=val1" --set "KEY2=val2"

# Import .env file
railway variables --from .env.production

# Run a command with Railway env injected
railway run node index.js
railway run python app.py
```

---

## Environments

```bash
# List environments (production, staging, etc.)
railway environment list

# Switch environment
railway environment
```

---

## Logs

```bash
# View live logs
railway logs

# View logs for specific service
railway logs --service <service-name>

# Stream logs with timestamps
railway logs --tail
```

---

## Database (PostgreSQL / MySQL / Redis / MongoDB)

```bash
# Add a database plugin
# (done via dashboard or railway.toml)

# Connect to PostgreSQL shell
railway connect postgres

# Connect to MySQL
railway connect mysql

# Connect to Redis
railway connect redis

# Get database connection string
railway variables | grep DATABASE_URL
```

---

## Domains

```bash
# Generate a Railway domain
railway domain

# Add a custom domain (via dashboard after CLI setup)
```

---

## railway.toml (Config File)

```toml
[build]
builder = "DOCKERFILE"
dockerfilePath = "./Dockerfile"

[deploy]
startCommand = "npm run start"
healthcheckPath = "/health"
healthcheckTimeout = 100
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 10

[[services]]
name = "web"
```

---

## Run Locally with Railway Env

```bash
# Run local dev with production env injected
railway run npm run dev

# Run migrations against Railway DB
railway run npx prisma migrate deploy
```
