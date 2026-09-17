---
title: "Kubernetes Deployment"
section: "07-DEPLOYMENT-OPS"
type: note
tags: [nats, kubernetes, helm, nack, operator, k8s]
source: https://docs.nats.io/running-a-nats-service/nats-kubernetes
created: 2026-06-14
---

# Kubernetes Deployment

Running NATS on Kubernetes via the **official Helm chart**, plus **NACK** for declarative JetStream resources.

See also: [[NATS/07-DEPLOYMENT-OPS/installation|Installation]] · [[NATS/06-TOPOLOGY-SCALING/clustering|Clustering]] · [[NATS/12-ECOSYSTEM/README|Ecosystem]]

---

## The Helm Chart

The supported install path is the **`nats`** Helm chart from `nats-io/k8s` (repo `https://nats-io.github.io/k8s/helm/charts/`).

```shell
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo update
helm install nats nats/nats -f values.yaml
```

### Example `values.yaml` (JetStream cluster, R3)
```yaml
config:
  cluster:
    enabled: true
    replicas: 3              # 3-node cluster (odd, for quorum)
  jetstream:
    enabled: true
    fileStore:
      pvc:
        size: 50Gi          # PersistentVolume per pod
    memoryStore:
      enabled: true
      maxSize: 1Gi
  monitor:
    enabled: true           # :8222
```

The chart provisions a **StatefulSet** (stable network IDs + per-pod PVCs for JetStream), a headless Service for clustering, and configures routes between pods automatically. JetStream needs StatefulSet + PVCs so file stores survive pod restarts.

---

## NACK — declarative JetStream (CRDs)

**NACK** (`nats-io/nack`, "NATS Action/Controllers for Kubernetes") provides **Custom Resource Definitions** so you manage Streams, Consumers, and KV as Kubernetes objects (GitOps-friendly):

```yaml
apiVersion: jetstream.nats.io/v1beta2
kind: Stream
metadata: { name: orders }
spec:
  name: ORDERS
  subjects: ["orders.>"]
  storage: file
  replicas: 3
  retention: limits
  maxAge: 168h            # 7 days
---
apiVersion: jetstream.nats.io/v1beta2
kind: Consumer
metadata: { name: dispatch }
spec:
  streamName: ORDERS
  durableName: dispatch
  ackPolicy: explicit
  filterSubject: orders.new
```

The NACK controller reconciles these into actual JetStream streams/consumers on the cluster — declarative, version-controlled messaging infrastructure.

---

## Key K8s Considerations

| Concern | Guidance |
|---------|----------|
| **StatefulSet, not Deployment** | JetStream needs stable identity + per-pod PVCs |
| **Anti-affinity** | Spread pods across nodes/zones so one node loss ≠ quorum loss |
| **Resources** | Set CPU/memory requests; JetStream wants fast disk (SSD storage class) |
| **Probes** | `/healthz` liveness + JS-aware readiness (see [[NATS/07-DEPLOYMENT-OPS/monitoring|Monitoring]]) |
| **Service exposure** | Internal ClusterIP for in-cluster clients; LoadBalancer/Ingress only if external clients needed |
| **Leaf nodes** | Edge devices connect *in* via leafnode port — expose 7422 carefully |
| **Upgrades** | Rolling update the StatefulSet; JetStream RAFT tolerates one pod at a time (R3) |

---

## OneMind Relevance

OneMind deliberately runs **launchd-supervised binary** on the MBP hub (no Docker/K8s) for locality and simplicity. Kubernetes becomes relevant only if a **cloud cluster** (omos-cloud / Lightsail → k8s) is stood up for HA or to host an internet-facing edge. The Helm chart + NACK would then give declarative, GitOps-managed JetStream — worth knowing for the growth path, not the current homestead hub.

## Sources
- NATS on Kubernetes — https://docs.nats.io/running-a-nats-service/nats-kubernetes
- Helm charts — https://github.com/nats-io/k8s
- NACK (JetStream controllers) — https://github.com/nats-io/nack
- Production deployment — https://docs.nats.io/running-a-nats-service/nats-kubernetes/helm-charts
