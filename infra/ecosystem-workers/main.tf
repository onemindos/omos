terraform {
  required_version = ">= 1.9, < 2.0"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
  # Configure this existing private Spaces backend through the release helper.
  # Credentials are environment-only and are never written into Git.
  backend "s3" {}
}

provider "digitalocean" {}

variable "cluster_id" {
  type        = string
  description = "Existing cloud cluster; this root does not manage or replace it."
  default     = "051f1134-ff05-451c-a52a-aab459feec6e"
}

variable "node_size" {
  type        = string
  description = "CPU worker size verified available in NYC3 on 2026-10-08."
  default     = "s-8vcpu-32gb-amd"
}

resource "digitalocean_kubernetes_node_pool" "ecosystem" {
  cluster_id = var.cluster_id
  name       = "onemind-ecosystem-cpu"
  size       = var.node_size
  node_count = 2
  auto_scale = true
  min_nodes  = 2
  max_nodes  = 4

  labels = {
    "onemind.io/compute-role" = "ecosystem"
    "onemind.io/site"         = "cloud"
  }

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [node_count]
  }
}

output "node_pool_id" {
  value = digitalocean_kubernetes_node_pool.ecosystem.id
}
