# infra/ — DigitalOcean foundation (Terraform)

> **Do not run `terraform apply` (or `make apply`) from this root.** The code does
> not match the live cluster, and there is no backend or state. A plan from this root
> proposes creating a second cluster. If state were imported, it would propose
> changing the live cluster, possibly by replacement. `prevent_destroy` on the cluster
> and the VPC blocks destroy and replace plans, but it does not block creating a
> duplicate cluster.

## Code versus live (read-only `doctl`, 2026-10-27)

| Item | `infra/main.tf` (defaults) | Live |
|---|---|---|
| VPC | `onemind-hub-vpc`, nyc3, `10.10.0.0/16` | `onemind-hub-vpc`, nyc3, `10.10.0.0/16` (matches) |
| Cluster name | `var.project_name` = `onemind-hub` | `omos-cloud` |
| DOKS version | `1.36.0-do.2` | `1.36.5-do.1` (auto-upgrade on) |
| HA / surge / auto-upgrade | on / on / on | on / on / on |
| Maintenance window | sunday 07:00 UTC | sunday 07:00 UTC (4 h) |
| Default node pool | `onemind-hub-core`, `s-4vcpu-8gb`, autoscale 3–5, label `role=core` | `onemind-hub-core-v2`, `s-8vcpu-16gb`, 3 nodes, no autoscale, label `role=core` |
| Second pool | not declared here | `onemind-ecosystem-cpu`, `s-8vcpu-32gb-amd`, 1 node, autoscale 1–4 (root `ecosystem-workers/`) |
| DO project | `onemind-hub` | `OMOS-CLOUD` (default project, Production) |

`providers.tf` defers the remote backend, and this root keeps no state in the
repository. `ecosystem-workers/` is a separate root with its own Spaces backend (see
its README) and manages only the ecosystem pool.

## Path to a safe plan

1. Configure a remote backend for this root (private Spaces bucket, the same pattern
   as `ecosystem-workers/`). Never commit state or plan files.
2. Import the live resources with `terraform import` (VPC, cluster, project). Treat
   any plan that shows a replace as a stop sign.
3. Align the code with live: cluster name, version (or `ignore_changes` for
   auto-upgrades), default pool name, size and scaling, and project. Repeat until
   `terraform plan` reports no changes.
4. Only then remove the warning in `main.tf` and in this file, and use `make plan`
   for reviews.
