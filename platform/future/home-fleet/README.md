# Future home fleet governance

Part of the deferred home GPU/fleet profile. No active Argo parent references these
Applications and none has automated sync. Their target is the future registered
cluster `onemind-home`, not the current cloud cluster.

- `../home-crossplane/application.yaml`: pinned Crossplane Helm controller.
  No provider credentials, managed resources or provisioning actions are declared.
- `../home-ocm/application.yaml`: pinned Open Cluster Management Helm hub.
  Register each actual site and define placement/policy through site-specific PRs.
- Cluster API: provider-specific cluster lifecycle belongs in a follow-up once
  the home hardware/virtualization provider is chosen. Do not install conflicting
  cluster lifecycle providers or create placeholder clusters.
- GPU serving: see `../home-gpu-serving/README.md`. KServe's native Kubernetes
  mode and the matching GPU plugin/operator are selected during home activation.

OCM owns fleet inventory/policy/placement; Argo owns assigned application
declarations; Crossplane owns only explicitly assigned infrastructure resources;
Cluster API owns provider-backed cluster lifecycle. Keep these ownership boundaries
separate and preserve autonomous local identity, NATS and critical apps at sites.

Home activation verifies partition/reconnect, local credentials, selected event
replication, artifact availability and the expiry of queued physical commands.
Adding fleet controllers does not make hardware transferable between locations.
