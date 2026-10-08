# OneMind portal and GPU deployment split

Recorded 2026-10-08. The user has no GPU now and expects GPU capacity later in the
home cluster. GPU-dependent deployment work belongs in separate future PRs, not
the current cloud deployment wave.

## Current CPU cluster

Proceed with the shared model gateway, operational data/API services, durable
missions, document processing, agent integrations and execution observability.
Use available model endpoints through the gateway. CPU-compatible local OCR and
speech/model workloads can be added with measured resource budgets. CPU-only does
not prevent agent coordination, application actions or the world model.

The future serving resources are under `platform/future/home-gpu-serving/`, with
no Argo or bootstrap reference and zero runtime replicas. Home hardware selection
and activation belong to a later PR. No GPU capacity is purchased by this change.

## Use our Backstage as the long-term portal

Keep `apps/onemind-portal` and integrate OpenChoreo's published frontend/backend,
catalog and scaffolder modules. Preserve the OneMind catalog, branding, TechDocs,
private mesh routing, Authentik subject mapping and existing integrations.

Use the bundled OpenChoreo portal as a disposable reference environment when
comparing integration behavior. It should not become a second permanent user
directory, identity system or main workspace.

The current source uses Backstage 1.55.3 with its new frontend system. OpenChoreo
1.3.x documents an existing-app path, but its tested compatibility matrix names
Backstage 1.51.0. Implement the integration in a portal branch; resolve package
compatibility there rather than downgrading the live portal by copying instructions.

Configure Authentik as the common identity authority with appropriate clients,
audiences and claims for OpenChoreo's actual requirements. The OpenChoreo UI needs
an API-authorized user token, not merely a Backstage identity token. Preserve
stable user mapping and deliberately compose OpenChoreo authorization with our
existing policy. Upstream examples allow non-OpenChoreo permissions by default;
do not replace OneMind's existing restrictions with that fallback.

Core integration adds catalog synchronization and component deployment/definition
views. Build/workflow/observability packs extend it as their APIs are connected.
Argo owns the assigned declarations; OpenChoreo owns generated resources. The
OneMind live operations application and TAK remain distinct mission interfaces.

Sources: [existing Backstage integration](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/installing-into-existing-backstage/),
[compatibility matrix](https://openchoreo.dev/docs/platform-engineer-guide/backstage-plugins/compatibility-matrix/),
[OneMind portal configuration](../../platform/backstage/README.md).
