# Contributing to OneMind OS

This repository contains both application source and production GitOps configuration. Changes under `platform/` and `infra/` can affect the live cluster, so keep them small, reviewable, and tied to the resource they manage.

## Branch and review workflow

- Treat `main` as the production-tracked branch.
- Make changes on a feature branch and submit a pull request for review.
- Do not assume GitHub branch protection is enabled; verify the repository settings when that matters.
- For workload changes, check the corresponding Argo CD Application and sync wave under `platform/argocd/`.

## Repository layout

- `infra/`: Terraform for the DigitalOcean foundation.
- `platform/`: Kubernetes resources and Argo CD deployment configuration.
- `services/`: runtime and application source (TAK image builds and the NATS bridge plugin are in `services/takserver/`); `services/schemas/` holds the versioned event contracts and generated code.
- `docs/`: architecture, operating procedures, and infrastructure guides.

See the root [`README.md`](README.md) for the full directory map.

## Code and configuration standards

### YAML and Kubernetes

- Use two-space indentation and avoid trailing whitespace.
- Keep resources with the platform layer or service that owns them.
- Follow the existing namespace, labels, probes, and resource request/limit pattern for that workload.
- Update the relevant Argo CD application or Kustomize configuration when adding resources; avoid expanding the manual bootstrap script for routine application deployments.

### Terraform

- Keep Terraform under `infra/` and format it with `terraform fmt`.
- Declare inputs in `variables.tf` and outputs in `outputs.tf`.
- Explain non-obvious resource settings in comments.
- Review the plan before applying infrastructure changes.

### Shell

- Use `#!/usr/bin/env bash` and `set -euo pipefail` for Bash scripts.
- Keep scripts compatible with ShellCheck.

### NATS contracts

Follow the subject grammar and compatibility rules in [`docs/NATS-DOCTRINE.md`](docs/NATS-DOCTRINE.md). Do not introduce a new subject prefix without updating the doctrine and affected schemas/consumers.

## Secrets

Never commit plaintext credentials, private keys, certificates, or local secret overrides. Use Bitnami SealedSecrets as described in the root `CLAUDE.md` ("Secrets") and follow the pattern used by the owning platform application. [`docs/secrets-eso.md`](docs/secrets-eso.md) is historical. Do not create or rotate live cluster secrets as part of a code change unless the task explicitly includes that operation.

## Local validation

Run the repository quality gate before opening a pull request:

```bash
make validate
```

For infrastructure changes, also review:

```bash
make plan
```

Do not run `make apply` or `make install` as a validation step; both can change live infrastructure or cluster resources.

## Pull request checklist

- [ ] The change follows the owning directory's existing structure.
- [ ] Related docs, manifests, schemas, or Argo CD configuration are updated.
- [ ] No plaintext secrets or generated local state are included.
- [ ] `make validate` passes, or the PR explains the exact blocker.
- [ ] Any Terraform plan or cluster impact is described for review.
