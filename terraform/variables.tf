# Inputs. Defaults encode the Step 1 decision (nyc3, HA on, 3x 8GB autoscale
# 3->5). Everything is a variable so we can resize/relocate without editing
# resource blocks — this is how "we can always adjust and lower" stays true.

variable "project_name" {
  description = "DO project + name prefix for all OneMind hub resources."
  type        = string
  default     = "onemind-hub"
}

variable "region" {
  description = "DO region slug. nyc3 = closest to the operator (EDT)."
  type        = string
  default     = "nyc3"
}

variable "k8s_version" {
  description = "DOKS version slug. Latest stable at provision time."
  type        = string
  default     = "1.36.0-do.2"
}

variable "node_size" {
  description = "Worker node droplet slug. s-4vcpu-8gb = 4 vCPU / 8 GB shared."
  type        = string
  default     = "s-4vcpu-8gb"
}

variable "node_count" {
  description = "Baseline node count. 3 = minimum for HA NATS JetStream (R3)."
  type        = number
  default     = 3
}

variable "node_min" {
  description = "Autoscaler floor."
  type        = number
  default     = 3
}

variable "node_max" {
  description = "Autoscaler ceiling. Headroom under load; scales back down when idle."
  type        = number
  default     = 5
}

variable "ha_control_plane" {
  description = "HA control plane (99.95% SLA, +$40/mo). On for production posture."
  type        = bool
  default     = true
}
