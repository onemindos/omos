# KubeSolo: Single-Node Kubernetes at the Edge

**Link:** https://www.kubesolo.io/

**License:** MIT (Open Source)  
**Maintained by:** Portainer  
**Status:** v1.0.0+  
**Repository:** GitHub

---

## What It Is

Single-node Kubernetes distribution that removes—not disables—the clustering machinery from Kubernetes. Designed as a drop-in replacement for Docker/Podman on single-node deployments.

**Core Philosophy:** Start from single-node and add clustering if needed, rather than start with a full cluster and strip it down.

---

## What It Can Do

### Resource Efficiency
- **Memory:** <200 MB RAM at idle
- **CPU:** Zero cluster overhead (no etcd quorum, leader election, multi-node networking)
- **Footprint:** Minimal—single process runs API server, scheduler, controller manager, and kubelet together

### Platform Support
- **Architectures:** ARM, ARM64, x86_64, RISC-V 64
- **Target environments:** Edge, IoT, factory floors, developer laptops, appliances

### Kubernetes Compatibility
- **Full Kubernetes API** — not a subset
- **Existing manifests** work without modification
- **Helm charts** supported out of the box
- **Custom Resource Definitions (CRDs)** fully functional
- **Same OCI images** — no container changes needed

### Feature Parity with Full Kubernetes
✓ Declarative workload management via standard manifests  
✓ Full RBAC, ServiceAccounts, namespace isolation  
✓ Liveness/readiness probes, restart policies, resource limits  
✓ Ecosystem compatibility: operators, Flux, ArgoCD, Prometheus  
✓ Service discovery, storage, ingress (simplified)

---

## Key Capabilities

### Installation & Operation
- **One-liner install:** `curl -sfL https://get.kubesolo.io | sudo sh`
- **Automated detection** of architecture
- **Kubectl compatible** — standard K8s tooling works immediately
- **Zero configuration** to get started

### What's Removed vs. Multi-Node K8s
- etcd quorum logic and consistency protocols
- Leader election machinery
- Multi-node networking overlays
- Control plane distribution
- Inter-node heartbeats and health checks

### What Remains
- Single unified control loop running in one process
- Full admission control
- Complete resource scheduling
- Pod networking (local bridge)
- Volume management
- RBAC enforcement

---

## Use Cases

| Use Case | Why KubeSolo | Benefit |
|----------|-------------|---------|
| **Edge Devices** | Kubernetes API at edge without cluster overhead | Standard tooling + local autonomy |
| **Factory/OT** | Air-gapped, self-contained, declarative | Secure, offline-capable, auditable |
| **IoT Gateways** | ARM64 + <200MB RAM | Portable, efficient, production-ready |
| **Developer Workstations** | Local K8s without Docker Desktop complexity | True K8s locally, test before cloud |
| **Remote/Branch Sites** | Single-node autonomy | Operate independently, sync when able |
| **Kiosk/Appliances** | Minimal footprint, embedded | Run containerized logic on appliances |

---

## KubeSolo vs. Alternatives

### vs. Docker/Podman
| Feature | Docker/Podman | KubeSolo |
|---------|--------------|----------|
| Declarative state mgmt | ✗ | ✓ |
| RBAC & multi-tenant | ✗ | ✓ |
| Helm support | ✗ | ✓ |
| Health probes & restarts | Limited | Full |
| Ecosystem (operators, mesh) | ✗ | ✓ |
| Production ready manifests | Compose only | K8s native |

### vs. Lightweight K8s (k3s, microK8s, etc.)
- **KubeSolo** = purpose-built for single-node, zero cluster logic
- **k3s** = lightweight multi-node cluster (can run single-node, but retains clustering overhead)
- **microK8s** = Ubuntu-focused, similar scope to k3s

**When to choose KubeSolo:** Single-node is your deployment model and you want pure single-node optimization.

---

## Getting Started

### Installation
```bash
# Install (one-liner)
curl -sfL https://get.kubesolo.io | sudo sh

# Verify
kubectl get nodes
# Output: kubesolo-node   Ready   <none>   8s   v1.30.0+

# Set context
kubectl config use-context kubernetes-admin@kubesolo

# Deploy a workload
kubectl apply -f mosquitto.yaml
# namespace/mosquitto created
```

### Example Workflow
```bash
# Standard K8s manifests work as-is
kubectl create namespace myapp
kubectl apply -f deployment.yaml -n myapp
kubectl get pods -n myapp
kubectl logs -f deployment/myapp -n myapp
```

---

## Architecture & Design

### Process Model
Single process running:
- API server
- Scheduler
- Controller manager
- Kubelet

All in one address space, no inter-process communication overhead.

### Storage
- **etcd replaced with:** Local SQLite-backed state (no consensus needed)
- **Default volume backend:** Host filesystem or local volumes
- **Persistent volumes:** Supported via local storage classes

### Networking
- CNI plugin support (but single-node doesn't need complex overlays)
- Service DNS via CoreDNS
- Ingress supported (via NGINX or similar)

---

## Resource Requirements

### Minimal
- **CPU:** 1 core (1-2 cores recommended for real workloads)
- **RAM:** 512 MB minimum, <200 MB idle, 1-2 GB for workloads
- **Disk:** 1-2 GB for OS + container images

### Example Edge Device
- **Raspberry Pi 4 (ARM64, 4GB RAM):** Fully supported and tested
- **Industrial SBC (x86_64, 2GB RAM):** Fully supported
- **RISC-V board (1GB RAM):** Supported (RISC-V 64)

---

## Deployment Patterns

### Pattern 1: Autonomous Edge Node
Deploy KubeSolo on edge device → runs independently → syncs state back to cloud on connectivity.

```yaml
# Edge deployment (fully declarative)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gateway-agent
  namespace: edge
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: agent
        image: myorg/edge-agent:latest
        env:
        - name: CLOUD_ENDPOINT
          value: "https://cloud.example.com"
```

### Pattern 2: Development → Production
Develop locally on KubeSolo (same manifests as cloud), deploy to multi-node K8s without changes.

### Pattern 3: Air-Gapped Operations
Run KubeSolo offline, manage via local API, push updates when connectivity returns.

---

## Integration Points

### NATS/Event Bus
KubeSolo can integrate with NATS for:
- Cluster orchestration (if running multiple KubeSolo nodes)
- State synchronization
- Event streaming from edge nodes

### Observability
- **Prometheus** scraping works via standard K8s metrics
- **Logs** via kubectl/container runtime
- **Metrics** via metrics-server (lightweight for single-node)

### GitOps
- **Flux** supported (simplified for single-node)
- **ArgoCD** supported
- **Kustomize** / **Helm** fully compatible

---

## Known Limitations

- **Single point of failure** — by design (edge devices often are single-points anyway)
- **No HA etcd** — not applicable (uses local state)
- **No federation** — each KubeSolo node is independent
- **Network overlays not optimized** — not needed for single-node

---

## Comparison Matrix

| Aspect | Docker | k3s | KubeSolo | Full K8s |
|--------|--------|-----|----------|----------|
| **Setup** | 5 min | 10 min | 2 min | 30+ min |
| **Memory** | ~100 MB | ~400 MB | <200 MB | 2+ GB |
| **API Compat** | None | Full | Full | Full |
| **Multi-node** | No | Yes | No | Yes |
| **RBAC** | No | Yes | Yes | Yes |
| **Helm** | No | Yes | Yes | Yes |
| **Use Case** | Dev/single | Edge/small | Edge/single | Enterprise |

---

## Learning & Resources

### Official
- **Docs:** https://www.kubesolo.io/ (Docs link)
- **GitHub:** KubeSolo on GitHub
- **Get Started:** https://www.kubesolo.io/ (GET STARTED button)

### Prerequisites
- Basic Kubernetes concepts (pods, deployments, services)
- kubectl CLI familiarity
- Container/OCI image knowledge

### Recommended Learning Path
1. Learn Kubernetes basics (standard K8s docs)
2. Install KubeSolo locally
3. Deploy existing Helm charts
4. Migrate from Docker Compose or Docker swarm
5. Deploy to edge devices (Raspberry Pi, industrial boards)

---

## When to Use KubeSolo

### ✓ Good Fit
- Single-node edge devices (RPi, industrial boards)
- Air-gapped / offline-first operations
- Developer workstations wanting K8s parity
- IoT gateways with containerized logic
- Remote sites needing autonomous operations
- Appliances / kiosk machines

### ✗ Not a Good Fit
- Multi-node clusters (use k3s or full K8s)
- High availability requirements (use full K8s cluster)
- Complex service mesh needs (still on roadmap)
- Heavy load balancing (use full K8s)

---

## Related Technologies

- **NATS**: Event substrate for edge orchestration
- **OneMind OS**: Local-first sovereign OS (can run KubeSolo for container workloads)
- **Kubernetes (full)**: Migration target for scaled deployments
- **k3s**: Lightweight multi-node alternative
- **Docker/Podman**: Container runtimes KubeSolo builds on

---

## Roadmap & Future

- Multi-node capability (optional federation)
- Enhanced networking for edge meshes
- Integration with Kubernetes federation API
- Storage snapshot/restore for edge mobility

---

**Captured:** 2026-07-09  
**Status:** Active & recommended for edge deployments  
**Relevance:** Sovereign local-first infrastructure, NATS integration potential
