# OneMind OS (`omos`)

This repository contains the source and GitOps configuration for OneMind OS: the cloud infrastructure, Kubernetes platform, service code, edge integrations, schemas, and operating documentation.

## Start here

- Read [`docs/README.md`](docs/README.md) for the architecture documents and infrastructure runbooks.
- Read [`CONTRIBUTING.md`](CONTRIBUTING.md) before changing code or deployment configuration.
- Use `make help` to see the repository's local commands.

## Repository map

| Path | Purpose |
|------|---------|
| `.github/` | CI and service image build workflows |
| `infra/` | Terraform for the DigitalOcean foundation, including the VPC and DOKS cluster |
| `platform/` | Kubernetes resources, Helm values, and Argo CD applications for cluster workloads |
| `tak-server/` | TAK Server chart, image build files, plugins, deployment overlays, and edge tools |
| `services/` | OneMind runtime services, bridges, and supporting scripts |
| `apps/` | Application code, including the TAK portal |
| `edge/` | Hardware-bound and Mac edge configuration |
| `schemas/` | Versioned event contracts and code generators |
| `plugins/` | OneMind platform plugins |
| `rads/` | Remote Autonomous Deployment System simulation and robotics work |
| `legacy/` | Legacy agent source, configuration, and operating material |
| `docs/` | Architecture, platform, security, and infrastructure documentation |
| `blueprints/` | Product, interface, and capability blueprints |
| `OneMind-Platform/` | Platform capability and integration notes |
| `references/` | External project and ecosystem research |
The two `*-PICTURE.md` files at the repository root are dated architecture snapshots. For current architecture and deployment guidance, use `docs/` and the manifests under `platform/`.

## Local quality checks

```bash
make help
make lint           # YAML, Terraform, and shell lint
make validate       # Full local validation target
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
- Argo CD applications in `platform/argocd/` reconcile Kubernetes resources from `platform/` and related directories.
- Keep credentials out of Git. Follow [`docs/secrets-eso.md`](docs/secrets-eso.md) and the established secret pattern for the service being changed.
- The NATS subject grammar and bus rules are documented in [`docs/NATS-DOCTRINE.md`](docs/NATS-DOCTRINE.md).

## Related repositories

- [OneMind UI](https://github.com/onemindos/omos-ui)
- [OneMind Legacy](https://github.com/onemindos/omos-legacy)
- [OneMind Flutter](https://github.com/onemindos/omos-flutter)
