# DigitalOcean provider.
# The API token is supplied via the DIGITALOCEAN_TOKEN environment variable
# (never committed). We source it from the authenticated doctl config at apply
# time so the secret lives in exactly one place.
#
# The `spaces_*` keys are only needed if/when we move Terraform state to a
# Spaces (S3) remote backend — deferred until Phase 2 hardening.

provider "digitalocean" {
  # token read from DIGITALOCEAN_TOKEN env var
}
