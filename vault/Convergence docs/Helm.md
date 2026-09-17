# Helm: The Package Manager for Kubernetes

**Link:** https://helm.sh/

**Project Type:** Open-source package manager for Kubernetes  
**Current Version:** Helm 4.2.2 (Helm 4 recently released)  
**Maintained by:** CNCF + Helm community  
**License:** Apache 2.0 (open source)  
**Repository:** GitHub (Helm project)  
**Status:** Graduated CNCF project, production-ready  
**First Release:** 2015  

---

## What It Is

**Helm** is a package manager for Kubernetes that simplifies the deployment, management, and lifecycle of Kubernetes applications through **Helm Charts** — templated, versioned, shareable packages.

**Core Purpose:** Replace copy-paste YAML deployment manifests with reusable, version-controlled, templated packages. Enable repeatable, auditable, and shareable application deployment.

**Tagline:** "The package manager for Kubernetes"

**Philosophy:** "Stop the copy-and-paste" — Standardize how Kubernetes applications are packaged, distributed, and managed.

---

## What It Can Do: Core Capabilities

### 1. **Package Kubernetes Applications**
```
A Helm Chart contains:
  - Kubernetes manifests (YAML files)
  - Values (configurable parameters)
  - Templates (templating engine for customization)
  - Metadata (name, version, description)
  - Hooks (lifecycle scripts: pre-install, post-install, etc.)

Result: Complete application package (not scattered YAML files)
```

**Example: Deploy a Database**
```
Without Helm:
  - Write 50+ YAML files manually
  - Secrets YAML, ConfigMap YAML, PVC YAML, StatefulSet YAML, Service YAML
  - Every environment: copy + modify values
  - No versioning, no rollback

With Helm:
  1. helm install postgres bitnami/postgresql
  2. Pass values: --set password=secret, --set replicas=3
  3. Chart handles all manifests
  4. Versioned, repeatable, auditable
```

### 2. **Manage Complexity**
```
Charts describe even the most complex applications:

Example: Multi-tier application
  - Web frontend (nginx Deployment)
  - API backend (Node.js StatefulSet)
  - Database (PostgreSQL StatefulSet)
  - Cache (Redis Deployment)
  - Message queue (RabbitMQ)
  - Monitoring (Prometheus)
  - Logging (ELK stack)

Single Helm Chart handles:
  - All manifests
  - All dependencies
  - All configuration
  - All versioning
  - Single point of authority

Result: Complex apps become manageable
```

### 3. **Easy Updates**
```
Traditional: Modify YAML, reapply to cluster
  - Risk of breaking changes
  - No rollback capability
  - Manual testing needed

Helm: In-place upgrades with hooks
  - Pre-upgrade hooks (validation, backup)
  - Rolling update (gradual replacement)
  - Post-upgrade hooks (migration scripts, tests)
  - Automatic rollback on failure
  
Example: helm upgrade postgres bitnami/postgresql
  → Helm compares old vs new
  → Applies only changed manifests
  → Runs upgrade hooks
  → Tests new version
  → Rolls back if needed
```

### 4. **Rollbacks**
```
Production incident: New version has bug

With Helm:
  helm rollback <release> <revision>
  
Example: helm rollback postgres 2
  → Instantly reverts to previous version
  → All manifests, data, configuration restored
  → No data loss (rolling back app code, not data)
  
Result: Instant incident recovery
```

### 5. **Simple Sharing**
```
Create chart once → Share across organization:

Methods:
  1. Public repositories (Artifact Hub)
  2. Private repositories (internal Artifactory, ChartMuseum)
  3. OCI registries (Docker Hub, Quay, etc.)
  4. Git repositories
  5. Local charts

Benefits:
  ✓ No copy-paste of YAML files
  ✓ Single source of truth
  ✓ Version history
  ✓ Dependency management
  ✓ Integration with CI/CD
```

### 6. **Templating & Customization**
```
Single chart, multiple environments:

Chart template (values.yaml):
  replicas: 3
  image: "{{ .Values.image }}"
  resources:
    limits:
      memory: "{{ .Values.memory }}"

Deploy dev (1 replica, small resources):
  helm install app . -f values-dev.yaml

Deploy prod (3 replicas, large resources):
  helm install app . -f values-prod.yaml

Same chart, different configurations
```

### 7. **Dependency Management**
```
Chart can depend on other charts:

Example: Web app needs PostgreSQL + Redis

Chart.yaml:
  dependencies:
    - name: postgresql
      version: "15.0"
      repository: "https://charts.bitnami.com/bitnami"
    - name: redis
      version: "18.0"
      repository: "https://charts.bitnami.com/bitnami"

Install command:
  helm dependency update   # Fetch dependencies
  helm install myapp .     # Deploy app + all dependencies

Result: One command deploys entire stack
```

---

## Architecture: How Helm Works

### The Three Components

**1. Chart** (package)
```
Directory structure:
  my-chart/
    Chart.yaml              # Chart metadata
    values.yaml             # Default configuration
    templates/
      deployment.yaml       # Kubernetes manifest template
      service.yaml
      configmap.yaml
      ...
    README.md
    LICENSE
```

**2. Repository** (storage)
```
Collections of charts:

Public repositories:
  - Artifact Hub (800+ charts)
  - Bitnami Charts (official)
  - Stable (official, deprecated)
  - Incubator (community)

Private repositories:
  - Artifactory
  - ChartMuseum
  - Harbor
  - Private OCI registry
```

**3. Release** (deployed instance)
```
A chart deployed into a cluster becomes a release.

Example:
  Chart: postgresql (template)
  Release: my-postgres-prod (deployed instance in production)
  Release: my-postgres-dev (another instance in development)

Same chart → Multiple releases → Different configurations
```

### The Workflow

```
1. Create/Download Chart
   helm create myapp        # Create new chart
   helm pull bitnami/nginx  # Download existing chart

2. Customize Chart
   Edit values.yaml
   Edit templates/

3. Validate Chart
   helm lint myapp          # Check for errors
   helm template myapp      # Dry-run, see final YAML

4. Deploy to Cluster
   helm install myapp .               # Deploy with defaults
   helm install myapp . -f prod.yaml  # Deploy with custom values

5. Manage Release
   helm list                          # Show all releases
   helm status myapp                  # Check release status
   helm get values myapp              # See values used
   helm get manifest myapp            # See deployed YAML

6. Update/Upgrade
   helm upgrade myapp .               # Update release

7. Rollback if Needed
   helm rollback myapp                # Revert to previous version

8. Cleanup
   helm uninstall myapp               # Delete release
```

---

## Key Helm Concepts

### Templates (Jinja-like Templating)

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}-app
spec:
  replicas: {{ .Values.replicas }}
  template:
    spec:
      containers:
      - name: app
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        resources:
          limits:
            memory: "{{ .Values.resources.memory }}"
            cpu: "{{ .Values.resources.cpu }}"
        env:
        {{- range .Values.env }}
        - name: {{ .name }}
          value: "{{ .value }}"
        {{- end }}
```

**Templating features:**
- Variable substitution: `{{ .Values.key }}`
- Loops: `{{- range .Values.items }}`
- Conditionals: `{{- if .Values.enabled }}`
- Filters: `{{ .name | upper }}`
- Defaults: `{{ .Values.key | default "fallback" }}`

### Values (Configuration)

```yaml
# values.yaml
replicas: 3
image:
  repository: myapp
  tag: "1.0.0"
resources:
  memory: "512Mi"
  cpu: "250m"
env:
  - name: ENV
    value: "production"
  - name: LOG_LEVEL
    value: "info"
config:
  database_url: "postgres://db:5432"
```

### Hooks (Lifecycle Scripts)

```yaml
# Pre-install hook: Validate prerequisites
apiVersion: batch/v1
kind: Job
metadata:
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
  name: pre-install-job
spec:
  template:
    spec:
      containers:
      - name: validation
        image: myapp:validation
        command: ["./validate.sh"]
```

**Hook types:**
- `pre-install` — Run before chart installation
- `post-install` — Run after chart installation
- `pre-upgrade` — Run before upgrade
- `post-upgrade` — Run after upgrade
- `pre-delete` — Run before deletion
- `post-delete` — Run after deletion
- `pre-rollback` — Run before rollback
- `post-rollback` — Run after rollback

---

## Integration with OneMind OS / Tactical Operations

### Scenario 1: Deploy TAK Server on KubeSolo

```
TAK Server (containerized) deployment via Helm

Helm Chart for TAKServer:
  - Deployment (TAK server container)
  - StatefulSet (PostgreSQL database)
  - ConfigMaps (TAK server config, certs)
  - Secrets (API keys, passwords)
  - Services (internal + external endpoints)
  - PersistentVolumes (data storage)

Install:
  helm repo add myrepo https://charts.example.com
  helm install takserver myrepo/takserver \
    --set postgresql.password=secret \
    --set ssl.enabled=true \
    --set replicas=3

Result:
  - TAK Server deployed on KubeSolo
  - Database initialized
  - SSL configured
  - Service exposed
  - No YAML copy-paste needed
```

### Scenario 2: Mumble Server Deployment

```
Helm Chart for Mumble (encrypted comms server)

Deploy to KubeSolo edge cluster:
  helm install mumble myrepo/mumble \
    --set domain=tactical-ops.local \
    --set ssl.certificate=/certs/cert.pem \
    --set users.admin=commander \
    --namespace comms

Manages:
  - Mumble server configuration
  - SSL/TLS certificates
  - User authentication
  - Channel structure
  - Persistent voice data

Benefits:
  ✓ Declarative configuration
  ✓ Easy to version control
  ✓ Simple updates (add users, change channels)
  ✓ Rollback if config breaks
  ✓ Reproducible across sites
```

### Scenario 3: Multi-Cluster Deployment (Tactical Network)

```
Deploy OneMind services across 3 tactical locations:

Each location has KubeSolo cluster:
  - Rear headquarters (full stack)
  - Forward operating base (edge)
  - Forward observer team (minimal)

Single Helm Chart with environment values:

values-hq.yaml:
  replicas: 3
  storage: 100Gi
  resources.cpu: "2000m"

values-fob.yaml:
  replicas: 1
  storage: 10Gi
  resources.cpu: "500m"

values-fwd.yaml:
  replicas: 1
  storage: 1Gi
  resources.cpu: "100m"

Deploy to all sites:
  helm install services myrepo/tactical-ops -f values-hq.yaml
  helm install services myrepo/tactical-ops -f values-fob.yaml
  helm install services myrepo/tactical-ops -f values-fwd.yaml

Result:
  - Same configuration across organization
  - Resource-appropriate for each location
  - Easy to manage, update, rollback
```

### Scenario 4: Automated CI/CD Pipeline

```
GitOps workflow with Helm:

1. Developer commits Chart changes to Git
2. CI pipeline triggers:
   helm lint          # Validate chart syntax
   helm template      # Generate YAML
   helm test          # Run chart tests
3. Approved? CD pipeline deploys:
   helm upgrade prod . --values prod.yaml
4. If failed: Automatic rollback
   helm rollback prod 1

Result:
  - Infrastructure as code
  - Auditable deployments
  - Automatic rollbacks
  - GitOps workflow
```

### Scenario 5: Cross-Cluster Service Mesh

```
Deploy TAK + Mumble + Analytics across multiple clusters:

Helm Chart with dependencies:
  dependencies:
    - name: takserver
      version: "1.0"
    - name: mumble
      version: "1.0"
    - name: prometheus
      version: "15.0"
    - name: grafana
      version: "6.0"

Deploy to multiple clusters:
  kubectl config use-context cluster-hq
  helm install services ./tactical-ops

  kubectl config use-context cluster-field
  helm install services ./tactical-ops

  kubectl config use-context cluster-rear
  helm install services ./tactical-ops

Benefits:
  ✓ Consistent deployment across sites
  ✓ All dependent services included
  ✓ Single point of authority
  ✓ Easy to manage at scale
```

---

## Helm vs. Traditional Kubernetes Deployment

| Aspect | Traditional (kubectl + YAML) | Helm |
|--------|------------------------------|------|
| **Configuration** | Multiple YAML files | Single Chart (templated) |
| **Reusability** | Copy-paste YAML | Reference chart from repo |
| **Versioning** | Manual, no tracking | Automatic, chart versions |
| **Updates** | Reapply all YAML | `helm upgrade` (selective) |
| **Rollback** | Manual recovery | `helm rollback` (instant) |
| **Customization** | Edit YAML files | Pass values at install time |
| **Sharing** | Manual sharing | Chart repositories |
| **Scaling** | Edit replica counts manually | Change value + upgrade |
| **Dependencies** | Manual install order | Declared dependencies |
| **Testing** | Manual testing | `helm test`, `helm template` |

---

## Helm Repositories & Charts

### Popular Public Repositories

**Artifact Hub** (Centralized hub)
```
URL: https://artifacthub.io/
Contains: 800+ public charts
Browse: Web interface to search/filter charts
Install: helm repo add <name> <url>
```

**Bitnami Charts** (Official)
```
URL: https://charts.bitnami.com/bitnami
Provides: PostgreSQL, MySQL, Redis, Nginx, Apache, etc.
Quality: High (well-maintained, security-focused)
```

**Stable** (Official, deprecated)
```
Note: Moved to Bitnami (no longer curated)
Recommendation: Use Bitnami instead
```

### Finding Charts

```bash
# Search for chart in Artifact Hub
helm search hub postgresql

# List all repos
helm repo list

# Add repository
helm repo add bitnami https://charts.bitnami.com/bitnami

# Update repo
helm repo update

# Search repo
helm search repo bitnami | grep postgresql

# Inspect chart
helm show chart bitnami/postgresql
helm show values bitnami/postgresql
helm show readme bitnami/postgresql
```

---

## Helm 4: Latest Release

### What's New in Helm 4

**Recent improvements:**
- Enhanced security (better secret handling)
- Improved template debugging
- Better error messages
- Faster chart installation
- Extended hook lifecycle
- New template functions

**Backward compatibility:**
- Helm 3 charts still work
- Gradual migration path
- No forced upgrades

### Version History

```
Helm 3 (2019-2026)
  - Major rewrite (Helm 2 was different)
  - Removed server-side component (Tiller)
  - Security improvements
  - Simplified design

Helm 4 (2026+)
  - Evolutionary improvements
  - Enhanced features
  - Better testing
  - Production hardening
```

---

## Installation & Getting Started

### Install Helm

**macOS (Homebrew):**
```bash
brew install helm
helm version
```

**Linux:**
```bash
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

**Windows (Chocolatey):**
```bash
choco install kubernetes-helm
helm version
```

**Other methods:**
- Download binary directly
- Scoop (Windows)
- Snap (Linux)
- Docker image

### Quick Start

```bash
# 1. Add repository
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# 2. Install chart
helm install nginx bitnami/nginx

# 3. Check release
helm list
helm status nginx
helm get values nginx

# 4. Update/upgrade
helm upgrade nginx bitnami/nginx --set replicas=3

# 5. Rollback if needed
helm rollback nginx 1

# 6. Uninstall
helm uninstall nginx
```

---

## Helm in Production

### Best Practices

**✓ Do This**
- Version control chart changes (Git)
- Use semantic versioning for charts
- Separate values per environment
- Tag releases with Git tags
- Test charts before deploying
- Use namespaces to isolate releases
- Monitor release health
- Document chart usage
- Automate with CI/CD

**✗ Don't Do This**
- Manually edit deployed manifests (breaks Helm)
- Hardcode secrets in charts (use Secrets object)
- Deploy without testing (always run `helm template` first)
- Share credentials in chart values
- Ignore chart dependencies
- Deploy without namespace isolation

### Production Deployment

```bash
# Validate chart
helm lint ./mychart

# Dry-run (see what will be deployed)
helm template myrelease ./mychart -f prod.yaml

# Verify once more
helm template myrelease ./mychart -f prod.yaml | kubectl apply --dry-run=client -f -

# Finally deploy
helm upgrade --install myrelease ./mychart -f prod.yaml

# Monitor
helm status myrelease
kubectl get pods -l app=myapp
```

---

## Helm for Developers

### Create a Chart

```bash
# Generate chart scaffolding
helm create mychart

# Chart structure:
# mychart/
#   Chart.yaml              # Chart metadata
#   values.yaml             # Default values
#   templates/
#     deployment.yaml       # K8s manifests
#     service.yaml
#     configmap.yaml
#     _helpers.tpl          # Template helpers
#   README.md
#   LICENSE
#   charts/                 # Dependency charts
```

### Package Chart

```bash
# Create chart package
helm package mychart

# Upload to repository
helm push mychart-1.0.0.tgz myrepo
```

### Test Chart

```bash
# Lint for errors
helm lint mychart

# Generate manifests (dry-run)
helm template mychart

# Install to test cluster
helm install test mychart --namespace test

# Run chart tests
helm test test

# Cleanup
helm uninstall test
```

---

## Community & Support

**Official Resources:**
- **Docs:** https://helm.sh/docs/
- **Community:** https://helm.sh/community/
- **GitHub:** https://github.com/helm/helm
- **Blog:** https://helm.sh/blog/

**Community Support:**
- **Slack:** Helm Users, Helm Development, Charts channels
- **SIG Apps:** Weekly meetings (Kubernetes SIG)
- **Developer Standups:** Thursdays 9:30-10am PT
- **Events:** KubeCon presentations, tutorials

**Maintainers:**
- Supported by CNCF
- Maintained by Helm community
- Backers: Bitnami, Google, Microsoft, Red Hat, IBM, etc.

---

## Helm Ecosystem

### Related Tools

| Tool | Purpose | Integration |
|------|---------|-------------|
| **Artifact Hub** | Chart discovery & sharing | Central registry for charts |
| **ChartMuseum** | Private chart repository | Self-hosted chart storage |
| **Helmfile** | Declarative Helm releases | Orchestrate multiple releases |
| **Helm Diff** | Preview upgrade changes | See what `helm upgrade` will do |
| **Helm Secrets** | Manage secrets in charts | Encrypt values with sops/age |
| **Kustomize** | YAML templating (alternative) | Alternative to Helm templates |
| **ArgoCD** | GitOps CD | Deploy Helm charts via Git |
| **Flux** | GitOps CD | Deploy Helm charts via Git |

---

## Use Cases

### Use Case 1: SaaS Platform Deployment
```
Deploy same SaaS to multiple customers:
  - One chart template
  - One values file per customer
  - helm install customer-1 . -f customers/customer-1.yaml
  - helm install customer-2 . -f customers/customer-2.yaml
  
Result: Consistent platform, customer-specific configs
```

### Use Case 2: Microservices Orchestration
```
Deploy 10 microservices at once:
  - One chart per service
  - Chart with dependencies (links all services)
  - helm install platform . (installs all services)
  
Result: Entire platform deployed in one command
```

### Use Case 3: Multi-Environment Promotion
```
Dev → Staging → Production:
  - Chart in Git
  - Different values per environment
  - CI/CD promotes chart through environments
  
Result: Same code, environment-specific configs
```

### Use Case 4: Vendor Software Deployment
```
Deploy third-party software (databases, monitoring, etc.):
  - Use official Helm chart from vendor
  - Customize with values
  - Automatic updates (new chart versions)
  
Result: One-command deployment of complex software
```

---

## Limitations & Considerations

### Limitations

- **Learning curve:** Templating syntax takes time to master
- **Debugging:** Generated YAML can be hard to debug
- **State management:** Helm doesn't track system state (only releases)
- **Rollback limits:** Only rolls back app code, not data
- **Secrets management:** Requires additional tooling for production

### When NOT to Use Helm

```
Don't use Helm if:
  - Deploying simple 1-container app (overkill)
  - Need strict isolation (use Kustomize + kubectl)
  - Third-party templating required (use values instead)
  - Team unfamiliar with Helm (training needed)
```

---

## Helm 4 vs. Helm 3

### Key Differences

| Feature | Helm 3 | Helm 4 |
|---------|--------|--------|
| **Server component** | None (removed) | None |
| **Security** | Good | Improved |
| **Template debugging** | Basic | Enhanced |
| **Error messages** | OK | Better |
| **Performance** | Good | Faster |
| **Backward compat** | N/A | Helm 3 charts work |

### Upgrade Path

```
From Helm 3 to Helm 4:
  1. Backup current releases: helm get values <release> > backup.yaml
  2. Upgrade Helm binary: brew upgrade helm (or similar)
  3. Test with dry-run: helm template <release> (should work)
  4. Upgrade releases as needed: helm upgrade <release>
```

---

## Practical Command Reference

```bash
# Repository Management
helm repo add myrepo https://charts.example.com
helm repo update
helm repo list
helm repo remove myrepo

# Search & Discovery
helm search hub postgresql
helm search repo postgresql
helm show chart bitnami/postgresql
helm show values bitnami/postgresql

# Installation
helm install myapp bitnami/postgresql
helm install myapp . -f values.yaml -n namespace

# Release Management
helm list
helm list --all-namespaces
helm status myapp
helm history myapp
helm get values myapp
helm get manifest myapp

# Updates & Upgrades
helm upgrade myapp bitnami/postgresql
helm upgrade myapp . -f new-values.yaml
helm upgrade myapp . --reuse-values

# Rollbacks & Cleanup
helm rollback myapp 1
helm uninstall myapp
helm uninstall myapp --keep-history

# Chart Development
helm create mychart
helm lint mychart
helm template mychart
helm package mychart
helm dependency update mychart
helm test myapp
```

---

**Captured:** 2026-07-09  
**Status:** Production-ready, actively maintained  
**Relevance:** Essential for Kubernetes deployment automation, critical for OneMind OS tactical operations deployment  
**Learning curve:** Moderate (2-4 weeks to proficiency)  
**Adoption:** Industry-standard (de facto Kubernetes package manager)  
**Integration priority:** Very High (needed for KubeSolo deployment, multi-site ops, CI/CD)  
**Key differentiator:** Combines templating, versioning, and deployment management in one tool  
**Latest version:** Helm 4.2.2 (Helm 4 recently released)  

---

## Related Resources in Convergence

- **KubeSolo.md** — Kubernetes distribution (Helm deploys to it)
- **NATS** — Can be deployed via Helm chart
- **TAK Server** — Containerizable, deployable via Helm
- **Mumble** — Deployable via Helm chart
- **CI/CD integration** — Helm integrates with GitOps (ArgoCD, Flux)

---

## Quick Integration: Deploying OneMind Services with Helm

```bash
# Create Helm chart for tactical operations stack
helm create oneMind-tactical

# Chart deploys:
#   - TAK Server (position sync)
#   - Mumble (encrypted comms)
#   - NATS (message bus)
#   - Prometheus (monitoring)
#   - Grafana (dashboards)

# Deploy to KubeSolo edge clusters:
kubectl config use-context fob-cluster
helm install oneMind ./oneMind-tactical -f values-fob.yaml -n tactical

kubectl config use-context hq-cluster
helm install oneMind ./oneMind-tactical -f values-hq.yaml -n tactical

# One chart, multiple deployments, environment-specific configs
# Easy updates, automatic rollbacks, full auditability
```
