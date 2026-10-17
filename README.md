# OneMind OS (`omos`)

This repository contains the source and GitOps configuration for OneMind OS: the cloud infrastructure, Kubernetes platform, service code, edge integrations, schemas, and operating documentation.

## Start here

- Read [`docs/README.md`](docs/README.md) for the architecture documents and infrastructure runbooks.
- Read [`docs/STACK.md`](docs/STACK.md) for what is running in the cluster, by tier, with backup, exposure and GitOps status.
- Read [`docs/REPO-LAYOUT.md`](docs/REPO-LAYOUT.md) for what each top-level directory is and which are legacy.
- Read [`CONTRIBUTING.md`](CONTRIBUTING.md) before changing code or deployment configuration.
- Use `make help` to see the repository's local commands.

## Repository map

| Path | Purpose |
|------|---------|
| `.github/` | CI quality gate (`quality.yml`) and service image build workflows |
| `infra/` | Terraform for the DigitalOcean foundation (VPC and DOKS cluster), plus `edge/` for hardware-bound and Mac edge configuration |
| `platform/` | Kubernetes resources, Helm values, and Argo CD applications for cluster workloads |
| `tak-server/` | TAK Server chart, image build files, plugins, deployment overlays, and edge tools |
| `services/` | First-party source code built into images (portals, OneMind runtime services, bridges), the bus contracts in `schemas/`, and supporting scripts |
| `docs/` | Architecture, platform, security, and infrastructure documentation |
| `references/` | External project and ecosystem research |
| `scripts/` | CI validation helpers and ecosystem tooling |
The two `*-PICTURE.md` files at the repository root are dated architecture snapshots. For current architecture and deployment guidance, use `docs/` and the manifests under `platform/`.

## Local quality checks

```bash
make help
make lint           # YAML, Terraform, and shell lint
make validate       # Full local validation target
```

CI (`.github/workflows/quality.yml`) additionally runs kubeconform on plain manifests, `helm template` on Helm-based Argo Applications and local charts, and gitleaks. To reproduce locally, install `kubeconform`, `helm` and `gitleaks`, then run:

```bash
python scripts/validate_manifests.py --all        # schema-check every plain manifest
python scripts/validate_helm_renders.py --all     # render every Helm chart (needs network, PyYAML)
gitleaks dir . --config .gitleaks.toml --baseline-path .gitleaks-baseline.json --redact
```

Terraform commands use `infra/` and require DigitalOcean credentials for `plan` or `apply`:

```bash
make plan
make apply
```

`make install` runs the cluster bootstrap script. Routine workload changes should go through the Argo CD configuration under `platform/`; review `platform/install.sh` before using it against a cluster.

NATS JetStream stream and KV resources are declared under `platform/nats/streams/` and reconciled through the NATS platform deployment. Do not use the older Python provisioning scripts as the source of truth for the cluster.

## Deployment model

- Terraform declares the DigitalOcean VPC and DOKS foundation in `infra/`.
- Argo CD applications in `platform/argocd/` (47 Applications live) reconcile Kubernetes resources from `platform/` and related directories. Eleven Helm releases, including the TAK Server core and NATS, are still installed by hand; see "GitOps gaps" in [`docs/STACK.md`](docs/STACK.md).
- Keep credentials out of Git. Secrets are committed as Bitnami SealedSecrets (see `CLAUDE.md`, "Secrets"). [`docs/secrets-eso.md`](docs/secrets-eso.md) describes the earlier External Secrets design and is historical.
- Dependency updates are proposed weekly by Renovate (`renovate.json`) and GitHub Actions are pinned by commit SHA.
- The NATS subject grammar and bus rules are documented in [`docs/NATS-DOCTRINE.md`](docs/NATS-DOCTRINE.md).

## Related repositories

- [OneMind UI](https://github.com/onemindos/omos-ui)
- [OneMind Legacy](https://github.com/onemindos/omos-legacy)
- [OneMind Flutter](https://github.com/onemindos/omos-flutter)
