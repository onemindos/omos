---
title: "Kubernetes & IaC"
section: "12-ECOSYSTEM"
type: note
tags: [nats, kubernetes, helm, nack, terraform, iac, operator, crd]
source: https://github.com/nats-io/k8s
created: 2026-06-14
---

# Kubernetes & IaC

Running and managing NATS declaratively: the **Helm charts** (`k8s`), the **NACK** controllers that turn streams/consumers into Kubernetes CRDs, and the **Terraform provider** for JetStream. This is the path when the fleet outgrows a single hand-managed server.

See also: [[NATS/07-DEPLOYMENT-OPS/kubernetes|Deployment: Kubernetes]] · [[NATS/12-ECOSYSTEM/repo-map|Repo Map]] · [[NATS/02-JETSTREAM/streams|Streams]]

---

## `k8s` — the Helm charts (supported path)

Repo: `nats-io/k8s` (★635, Go). The **official, supported** way to run NATS on Kubernetes. Deploys a NATS StatefulSet (with JetStream PVCs), headless service, monitoring, and optional surveyor/box.

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm install nats nats/nats \
  --set config.jetstream.enabled=true \
  --set config.cluster.enabled=true \
  --set config.cluster.replicas=3
```

Key knobs: JetStream storage (file/PVC size), cluster replicas (3 or 5 for RAFT quorum), TLS, leafnodes, auth. The chart is the **first thing to reach for** — not the deprecated operator.

> See [[NATS/07-DEPLOYMENT-OPS/kubernetes|Section 07: Kubernetes]] for the deployment-side detail.

---

## `nack` — NATS Controllers for Kubernetes

Repo: `nats-io/nack` (★244, Go). NACK lets you manage **JetStream assets as Kubernetes CRDs** — define a `Stream` or `Consumer` in YAML and a controller reconciles it against the server. GitOps for your streams.

```yaml
apiVersion: jetstream.nats.io/v1beta2
kind: Stream
metadata:
  name: orders
spec:
  name: ORDERS
  subjects: ["orders.>"]
  storage: file
  replicas: 3
  retention: limits
---
apiVersion: jetstream.nats.io/v1beta2
kind: Consumer
metadata:
  name: orders-dispatch
spec:
  streamName: ORDERS
  durableName: dispatch
  filterSubject: orders.new
  ackPolicy: explicit
```

`kubectl apply` and the stream/consumer exist. Drift is reconciled. This is how you keep JetStream topology **in version control** alongside the apps that use it.

---

## `terraform-provider-jetstream` — IaC for JetStream

Repo: `nats-io/terraform-provider-jetstream` (★82, Go). Manage JetStream **streams, consumers, and KV buckets** from Terraform — for shops standardized on Terraform rather than k8s CRDs.

```hcl
provider "jetstream" {
  servers     = "nats://hub:4222"
  credentials = "/etc/nats/admin.creds"
}

resource "jetstream_stream" "orders" {
  name     = "ORDERS"
  subjects = ["orders.>"]
  storage  = "file"
  replicas = 3
}

resource "jetstream_consumer" "dispatch" {
  stream_id      = jetstream_stream.orders.id
  durable_name   = "dispatch"
  filter_subject = "orders.new"
  ack_policy     = "explicit"
}
```

Same idea as NACK, different toolchain: declarative, reviewable, reproducible stream topology.

---

## `nats-operator` † — deprecated

Repo: `nats-io/nats-operator` (★568, Go) — **archived/deprecated**. The original CRD-based operator for running the *server*. **Do not use for new deployments.** Its responsibilities split into: **Helm `k8s`** (run the server) + **NACK** (manage JetStream assets). High star count is legacy gravity, not a recommendation.

---

## Choosing your management plane

| You want… | Use |
|-----------|-----|
| Run NATS on k8s | **Helm `k8s`** |
| Streams/consumers as k8s YAML (GitOps) | **NACK** |
| Streams/consumers via Terraform | **terraform-provider-jetstream** |
| Streams/consumers imperatively | `nats` CLI / `jsm.go` |
| (old tutorials say `nats-operator`) | ❌ deprecated — use Helm + NACK |

All of these ultimately call the same JetStream `$JS.API.>` management surface — pick the declarative skin that matches your stack.

---

## OneMind Relevance
- **Today:** the hub is a single `nats-server` on the Mac under launchd — no k8s needed. Stream topology is managed via `nats` CLI / our ops skill.
- **Scaling path:** if the fleet ever runs a Kubernetes hub (cloud Lightsail cluster, omos-cloud), **Helm `k8s` + NACK** is the declarative target — JetStream topology in git, reconciled automatically.
- **Either way, avoid `nats-operator`** — it's the deprecated trap that old guides point to.
- Even without k8s, the **GitOps principle** (streams defined as reviewable config, not ad-hoc CLI) is worth adopting via Terraform or checked-in `nats` scripts.

## Sources
- k8s Helm charts — https://github.com/nats-io/k8s
- NACK — https://github.com/nats-io/nack
- Terraform provider — https://github.com/nats-io/terraform-provider-jetstream
- Running on Kubernetes — https://docs.nats.io/running-a-nats-service/nats-kubernetes
