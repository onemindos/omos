# omos-infra

Infrastructure-as-code for the **OneMind cloud hub** — the always-on control
plane. Terraform is the single source of truth; the same definitions can later
target a home rack (no cloud lock-in).

Part of the OneMind four-plane infra model:
`Portainer (mgmt) / K8s cloud hub / NEX+KubeSolo edge / NATS bus`.

## Posture (decided 2026-07-04)
- **Cloud is the authority** for now. NATS runs as an HA cluster in-cluster;
  main services live in the cloud, automated and self-healing.
- **Rebuild clean, don't copy the Mac.** Only services with a real job get
  deployed. Clean streams/auth/observability from first principles.
- **Hybrid / local-first comes later**, after production is proven and
  dedicated hardware is in place.
- **Always reversible.** Node pools resize/scale (even down); every service is
  NATS-native so nothing is trapped in the cloud.

## Layout
```
terraform/
  versions.tf    provider + TF version pins
  providers.tf   DigitalOcean provider (token via env)
  variables.tf   all inputs (region, node size/count, HA, registry tier)
  main.tf        VPC + DOKS cluster + container registry + DO project
  outputs.tf     cluster id/endpoint/version, VPC, registry endpoint
```

## Step 1 — provision the empty cluster
Requires: `terraform`, `doctl` authenticated, `kubectl`.

```bash
cd terraform
export DIGITALOCEAN_TOKEN="$(doctl auth token)"   # or your PAT
terraform init
terraform plan
terraform apply
```

After apply, wire up kubeconfig and verify:
```bash
doctl kubernetes cluster kubeconfig save onemind-hub
kubectl get nodes -o wide          # expect 3 Ready nodes
```

## Current spec (Step 1)
| Item | Value |
|------|-------|
| Region | nyc3 |
| K8s | 1.36.0-do.2 |
| Control plane | HA (99.95% SLA) |
| Node pool | 3 × s-4vcpu-8gb, autoscale 3→5 |
| Registry | basic (5 GB) |
| Deferred | Postgres, load balancer, Spaces (added when needed) |

Roughly ~$190/mo running; covered by DO credit. Per-second billing — resize or
destroy anytime.
