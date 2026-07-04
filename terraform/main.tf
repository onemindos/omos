# OneMind cloud hub — Step 1 foundation.
# Provisions ONLY the primitives an empty production cluster needs:
#   VPC (network isolation) + DOKS cluster + container registry.
# Deliberately NOT here yet (added when a workload needs them):
#   managed Postgres, load balancer, Spaces buckets.

# --- Private network -------------------------------------------------------
# Isolates all hub traffic on a private VPC. Node<->node and future
# managed-DB traffic stays off the public internet.
resource "digitalocean_vpc" "hub" {
  name     = "${var.project_name}-vpc"
  region   = var.region
  ip_range = "10.10.0.0/16"
}

# --- Kubernetes cluster ----------------------------------------------------
# The always-on hub. HA control plane for production. Auto-upgrade of patch
# releases on; surge upgrades on for zero-downtime node rollouts.
resource "digitalocean_kubernetes_cluster" "hub" {
  name         = var.project_name
  region       = var.region
  version      = var.k8s_version
  vpc_uuid     = digitalocean_vpc.hub.id
  ha           = var.ha_control_plane
  auto_upgrade = true
  surge_upgrade = true

  # Free container-registry integration so nodes can pull our images.
  registry_integration = true

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
}

# --- Container registry ----------------------------------------------------
# Holds our own service images (om-control-plane, etc.). One registry per
# account; basic tier gives 5 GB and multiple repositories.
resource "digitalocean_container_registry" "hub" {
  name                   = var.project_name
  subscription_tier_slug = var.registry_tier
  region                 = var.region
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
