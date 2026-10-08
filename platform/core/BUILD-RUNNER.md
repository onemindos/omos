# Build runner and GitOps boundary

The persistent runner image treats EPHEMERAL and DISABLE_AUTO_UPDATE as enabled
when present, even if their value is "false". Leave them unset for the current
persistent runner. The functional smoke workflow checks Docker without deploying.
Both runner and Docker images are pinned to the inspected running digests.

Image workflows publish immutable commit tags. They do not directly restart cluster
Deployments. Deploy a new image through a reviewed digest change in the owning
manifest and let Argo sync it. Default runner service-account API token mounting
and the obsolete cross-namespace deployment mutation role are removed. The separate
inventory runner uses its own observational credentials and explicit opt-in.

Existing manifests still contain mutable images and some Helm/manual ownership;
this boundary change does not establish complete version pinning or GitOps parity.
