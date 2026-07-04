# Outputs. Sensitive values (kubeconfig) are marked so they don't print to
# the console; fetch the kubeconfig via doctl after apply instead.

output "cluster_id" {
  description = "DOKS cluster ID."
  value       = digitalocean_kubernetes_cluster.hub.id
}

output "cluster_name" {
  description = "DOKS cluster name."
  value       = digitalocean_kubernetes_cluster.hub.name
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint."
  value       = digitalocean_kubernetes_cluster.hub.endpoint
}

output "cluster_version" {
  description = "Running Kubernetes version."
  value       = digitalocean_kubernetes_cluster.hub.version
}

output "node_pool" {
  description = "Core node pool shape."
  value = {
    size       = var.node_size
    min        = var.node_min
    max        = var.node_max
    baseline   = var.node_count
  }
}

output "vpc_id" {
  description = "Private VPC ID."
  value       = digitalocean_vpc.hub.id
}

output "registry_endpoint" {
  description = "Container registry endpoint for image pushes."
  value       = digitalocean_container_registry.hub.endpoint
}
