# Contributing to omos-infra

## Branch Strategy

- `main` — production-deployed state
- Feature branches → PR → quality gate → merge
- No direct pushes to main

## Code Standards

### YAML
- 2-space indentation
- No trailing whitespace
- yamllint config: `.yamllint.yml`

### Terraform
- `terraform fmt` (enforced by CI + pre-commit)
- Variables in `variables.tf`, outputs in `outputs.tf`
- Comments for non-obvious resource configs

### Shell Scripts
- `#!/usr/bin/env bash` shebang
- `set -euo pipefail` at top
- Pass ShellCheck with `-x -s bash`

### Kubernetes Manifests
- One resource per file (unless tightly coupled)
- Namespace declared in its own `namespace.yaml`
- Labels: `app.kubernetes.io/name`, `app.kubernetes.io/part-of: onemind`
- Resource requests + limits on all Deployments

## NATS Subject Grammar (LOCKED)

```
tel.{node}.{class}.{metric}          telemetry
det.{domain}.{type}.{id}             detections
evt.{domain}.{entity}.{event}        semantic events
cmd.{domain}.{entity}.{action}       commands
intel.{domain}.{source}.{type}       intelligence
ent.{entity_type}.{entity_id}.state  fused entity state
svc.{service}.{endpoint}             NATS RPC
agents.{verb}.{agent}.{owner}.*      agent fabric
tak.cot.inject                       CoT → TAK inject
tak.cot.>                            CoT from TAK
```

Do NOT invent new subject patterns without updating this grammar.

## Adding a New Service

1. Create `platform/{service}/` directory
2. Add required files:
   - `namespace.yaml`
   - `deployment.yaml` (Deployment + Service)
   - `external-secret.yaml` (ESO secret sync from 1Password)
   - `cloudflare-route.yaml` (if public-facing)
3. Image: `ghcr.io/onemindos/{service}:latest`
4. Add `kubectl apply` lines to `platform/install.sh`
5. Add ArgoCD Application in `platform/argocd/`
6. Update CLAUDE.md service table

## Secrets

**NEVER commit secrets.** All secrets flow through:
```
1Password vault "00-24 SO" → ESO → K8s Secret → Pod env
```

Bootstrap secrets (one-time on fresh cluster) are documented in `platform/install.sh`.

## Testing Changes

```bash
# Validate everything locally
make validate

# Terraform dry-run
make plan

# Deploy to cluster
make install
```

## Review Checklist

- [ ] `make lint` passes
- [ ] No hardcoded secrets/tokens
- [ ] Resource requests/limits set
- [ ] Namespace declared
- [ ] ExternalSecret for any credentials
- [ ] install.sh updated (if new service)
- [ ] CLAUDE.md service table updated
