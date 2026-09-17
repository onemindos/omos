terraform {
  required_providers {
    jetstream = {
      source  = "nats-io/jetstream"
      version = "~> 0.1"
    }
  }
  required_version = ">= 1.0"
}

provider "jetstream" {
  servers = var.nats_url
}

# ─────────────────────────────────────────────
# Streams
# ─────────────────────────────────────────────

resource "jetstream_stream" "evt" {
  name     = "EVT"
  subjects = ["evt.>"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 30 # 30 days in seconds
}

resource "jetstream_stream" "tel" {
  name     = "TEL"
  subjects = ["tel.>"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 7 # 7 days in seconds
}

resource "jetstream_stream" "cmd" {
  name     = "CMD"
  subjects = ["cmd.>"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 1 # 1 day in seconds
}

resource "jetstream_stream" "agent" {
  name     = "AGENT"
  subjects = ["agent.>"]
  storage  = "file"
  max_age  = 60 * 60 * 24 * 30 # 30 days in seconds
}

# ─────────────────────────────────────────────
# KV Buckets
# ─────────────────────────────────────────────

resource "jetstream_kv_bucket" "om_config" {
  name    = "om-config"
  history = 10
}

resource "jetstream_kv_bucket" "om_state" {
  name    = "om-state"
  history = 5
}

resource "jetstream_kv_bucket" "device_registry" {
  name    = "device-registry"
  history = 10
}
