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

Docker binds its TCP API only to 127.0.0.1 inside the runner pod. Other pods cannot
operate its privileged Docker service. A namespace NetworkPolicy denies ingress
and permits DNS plus public HTTP/HTTPS for GitHub, registries and package downloads.
Private-cluster connections require a reviewed explicit exception; the separate
observational inventory runner is not selected by this policy. Verify Docker and
job completion with the smoke workflow after changing images or network policy.
