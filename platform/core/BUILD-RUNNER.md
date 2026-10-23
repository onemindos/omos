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

## Pull-request checks run on this runner

GitHub-hosted runners are unavailable while the account billing lock lasts, so the
read-only check workflows (`quality.yml`, `ecosystem-charts.yml`) run on
`[self-hosted, linux, amd64, zeus-v2]`, the persistent runner in this namespace. The
second pod, `github-runner-hermes`, is repository-scoped, ephemeral, has only the label
`onemind-hermes-readonly` and no `zeus-v2`, so check jobs cannot land on it.

- **Private repositories only.** Never register this runner for a public repository or
  run fork pull requests on it: a pull request job executes its branch's code inside a
  pod that has a privileged Docker sidecar. `omos` is private and only collaborators
  can open branches.
- Jobs keep `permissions: contents: read`, checkouts use `persist-credentials: false`
  (the runner is persistent), and tools are installed into `$RUNNER_TEMP`, not with `sudo`.
- The pod has no Kubernetes API token, no ingress, and egress only to public 80/443
  (see `github-runner-network.yaml`). A job cannot reach the cluster.
- Capacity: one runner, `limits: 2 CPU / 1Gi` for the runner container. Jobs run one at a
  time. The portal `yarn tsc` / `yarn test` need more memory than that; they are not in
  any pull-request workflow.
- To go back to hosted runners, set `runs-on: ubuntu-latest` again in those two files.
