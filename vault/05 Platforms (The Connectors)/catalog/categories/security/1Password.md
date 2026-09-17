# 1Password

> **Secrets & Password Management**

## Overview

1Password is the central secrets management solution for OneMind OS, storing API keys, credentials, and sensitive configuration.

## Status

- **Active**: ✅ Yes
- **Plan**: Personal/Business
- **CLI**: ✅ Installed
- **Last Updated**: January 25, 2026

## Use Cases

1. **API Keys** - Store service credentials
2. **Database Passwords** - Secure database access
3. **SSH Keys** - Secure server access
4. **Environment Variables** - Inject secrets into apps
5. **Shared Team Secrets** - Secure collaboration

## Vaults Organization

| Vault              | Purpose              | Access  |
| ------------------ | -------------------- | ------- |
| **Personal**       | Personal accounts    | Private |
| **OneMind**        | Project secrets      | Private |
| **Development**    | Dev/test credentials | Private |
| **Infrastructure** | Server/cloud access  | Private |

## CLI Setup

### Installation

```bash
# macOS
brew install 1password-cli

# Verify
op --version
```

### Authentication

```bash
# Sign in (first time)
op account add

# Sign in (subsequent)
eval $(op signin)

# Verify
op whoami
```

## Retrieving Secrets

### Single Secret

```bash
# Get API key
op read "op://OneMind/Anthropic API/api_key"

# Get password
op read "op://OneMind/Database/password"
```

### In Environment Variables

```bash
# Export to environment
export ANTHROPIC_API_KEY=$(op read "op://OneMind/Anthropic API/api_key")

# In a script
API_KEY="$(op read 'op://OneMind/Anthropic API/api_key')"
```

### Full Item

```bash
# Get full item as JSON
op item get "Anthropic API" --vault OneMind --format json
```

## Secret References

### Format

```
op://<vault>/<item>/<field>
```

### Examples

```bash
op://OneMind/Anthropic API/api_key
op://OneMind/Database/password
op://OneMind/AWS/access_key_id
op://OneMind/AWS/secret_access_key
```

## Environment File Integration

### .env Template

```bash
# .env.template
ANTHROPIC_API_KEY=op://OneMind/Anthropic API/api_key
DATABASE_URL=op://OneMind/Database/connection_string
AWS_ACCESS_KEY_ID=op://OneMind/AWS/access_key_id
AWS_SECRET_ACCESS_KEY=op://OneMind/AWS/secret_access_key
```

### Inject Secrets

```bash
# Create actual .env file
op inject -i .env.template -o .env

# Run command with injected secrets
op run --env-file=.env.template -- python main.py
```

## Python Integration

### 1Password SDK

```python
import os
import subprocess
import json

def get_secret(reference: str) -> str:
    """Get a secret from 1Password."""
    result = subprocess.run(
        ["op", "read", reference],
        capture_output=True,
        text=True
    )
    if result.returncode != 0:
        raise Exception(f"Failed to get secret: {result.stderr}")
    return result.stdout.strip()

# Usage
api_key = get_secret("op://OneMind/Anthropic API/api_key")
```

### With Environment Loading

```python
import subprocess

def load_secrets(template_path: str = ".env.template"):
    """Load secrets from 1Password into environment."""
    result = subprocess.run(
        ["op", "inject", "-i", template_path],
        capture_output=True,
        text=True
    )

    for line in result.stdout.strip().split("\n"):
        if "=" in line and not line.startswith("#"):
            key, value = line.split("=", 1)
            os.environ[key] = value
```

## Docker Integration

### Docker Compose

```yaml
services:
  legacy-agent:
    build: .
    environment:
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - DATABASE_URL=${DATABASE_URL}
```

Run with:

```bash
op run --env-file=.env.template -- docker compose up
```

### BuildKit Secrets

```dockerfile
# Dockerfile
RUN --mount=type=secret,id=api_key \
    API_KEY=$(cat /run/secrets/api_key) && \
    do_something_with_api_key
```

## GitHub Actions Integration

### Workflow

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Configure 1Password
        uses: 1password/load-secrets-action@v2
        with:
          export-env: true
        env:
          OP_SERVICE_ACCOUNT_TOKEN: ${{ secrets.OP_SERVICE_ACCOUNT_TOKEN }}
          ANTHROPIC_API_KEY: op://OneMind/Anthropic API/api_key
          DATABASE_URL: op://OneMind/Database/connection_string

      - name: Deploy
        run: ./deploy.sh
```

## SSH Key Management

### SSH Agent Integration

```bash
# Enable SSH agent
op plugin init ssh

# Add to ~/.zshrc
export SSH_AUTH_SOCK=~/.1password/agent.sock

# Use SSH keys from 1Password
ssh user@server  # Uses 1Password SSH key
```

### Configure SSH

```
# ~/.ssh/config
Host *
    IdentityAgent "~/.1password/agent.sock"
```

## Creating Items

### Via CLI

```bash
# Create login item
op item create \
  --category login \
  --vault OneMind \
  --title "New Service" \
  username=admin \
  password=secret123

# Create API credential
op item create \
  --category "API Credential" \
  --vault OneMind \
  --title "Service API" \
  api_key=sk-123456789
```

## Best Practices

### Organization

- ✅ Use descriptive item names
- ✅ Organize by vault (Personal, Work, Project)
- ✅ Add tags for categorization
- ✅ Include URLs for context

### Security

- ✅ Never hardcode secrets in code
- ✅ Use secret references in config templates
- ✅ Rotate secrets regularly
- ✅ Use service accounts for automation
- ✅ Enable 2FA on your account

### Development Workflow

```bash
# Local development
op run --env-file=.env.template -- python main.py

# Testing
op run --env-file=.env.template -- pytest

# Docker
op run --env-file=.env.template -- docker compose up
```

## Secret Categories

| Category      | Items                                 |
| ------------- | ------------------------------------- |
| **APIs**      | Anthropic, OpenAI, Twilio, Cloudflare |
| **Databases** | PostgreSQL, Redis                     |
| **Cloud**     | AWS, Azure, GCP                       |
| **Services**  | GitHub, Vercel, Stripe                |
| **Servers**   | SSH keys, server passwords            |

## Troubleshooting

### Session Expired

```bash
# Re-authenticate
eval $(op signin)
```

### Can't Find Item

```bash
# List items in vault
op item list --vault OneMind

# Search
op item list --vault OneMind | grep -i "anthropic"
```

## Resources

- [1Password CLI Documentation](https://developer.1password.com/docs/cli/)
- [Secret References](https://developer.1password.com/docs/cli/secret-references/)
- [GitHub Actions Integration](https://developer.1password.com/docs/ci-cd/github-actions/)

---

**Related**: [GitHub](./GitHub.md) | [Docker](./Docker.md) | [AWS](./AWS.md)


---

## Related

- [[00-24 SO (Sovereign Operations)/02 Infrastructure (VPS, Servers & Networking) (The Body)/02 Infrastructure|02 Infrastructure]]
