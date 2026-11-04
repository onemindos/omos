# !!! DO NOT APPLY: this root does NOT match the live cluster (checked 2026-10-27).
# There is no backend and no state, so `terraform plan` proposes a second cluster
# named "onemind-hub"; against imported state it would propose changes to the
# live cluster (name, version and default node pool all differ), which may
# include replacement. Live is `omos-cloud` (DOKS 1.36.5-do.1) with pools
# onemind-hub-core-v2 (s-8vcpu-16gb x3, no autoscale) and onemind-ecosystem-cpu
# (managed by infra/ecosystem-workers/).
# See infra/README.md for live values and the import/alignment path.
# prevent_destroy below blocks destroy/replace plans; it does not block creation.

# OneMind cloud hub — Step 1 foundation.
# Provisions ONLY the primitives an empty production cluster needs:
#   VPC (network isolation) + DOKS cluster.
# Deliberately NOT here yet (added when a workload needs them):
#   managed Postgres, load balancer, Spaces buckets.
# Images are on GHCR (ghcr.io/onemindos/) — DOCR not used.

# --- Private network -------------------------------------------------------
# Isolates all hub traffic on a private VPC. Node<->node and future
# managed-DB traffic stays off the public internet.
resource "digitalocean_vpc" "hub" {
  name     = "${var.project_name}-vpc"
  region   = var.region
  ip_range = "10.10.0.0/16"

  lifecycle {
    prevent_destroy = true
  }
}

# --- Kubernetes cluster ----------------------------------------------------
# The always-on hub. HA control plane for production. Auto-upgrade of patch
# releases on; surge upgrades on for zero-downtime node rollouts.
resource "digitalocean_kubernetes_cluster" "hub" {
  name          = var.project_name
  region        = var.region
  version       = var.k8s_version
  vpc_uuid      = digitalocean_vpc.hub.id
  ha            = var.ha_control_plane
  auto_upgrade  = true
  surge_upgrade = true

  node_pool {
    name       = "${var.project_name}-core"
    size       = var.node_size
    auto_scale = true
    min_nodes  = var.node_min
    max_nodes  = var.node_max
    node_count = var.node_count

    labels = {
      role = "core"
    }
  }

  maintenance_policy {
    start_time = "07:00" # UTC ~ 03:00 EDT, low-activity window
    day        = "sunday"
  }

  # Live cluster: a destroy or replace would take down every workload on it.
  lifecycle {
    prevent_destroy = true
  }
}

# --- Project grouping ------------------------------------------------------
# Group all hub resources under one DO Project for a clean console view.
resource "digitalocean_project" "hub" {
  name        = var.project_name
  description = "OneMind cloud control-plane hub (K8s). Managed by Terraform."
  purpose     = "Operational / Developer tools"
  environment = "Production"

  resources = [
    digitalocean_kubernetes_cluster.hub.urn,
  ]
}
