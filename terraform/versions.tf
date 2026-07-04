# Provider + Terraform version pinning.
# Pinned for reproducibility — the same code must produce the same infra
# on this Mac today or a home rack later (no lock-in, portable IaC).

terraform {
  required_version = ">= 1.9"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.43"
    }
  }
}
