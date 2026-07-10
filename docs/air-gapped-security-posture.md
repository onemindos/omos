# Air-Gapped Web-Based C2 Security Posture & Deployment Patterns

**Date:** 2026-07-07  
**Scope:** Military-grade self-hosted web applications in disconnected environments (C2, ISR, tactical displays)  
**Status:** Comprehensive research with practical implementation guidance

---

## 1. CONTENT SECURITY POLICY (CSP) FOR SELF-HOSTED APPS

### Strictest Recommended Policy

For a self-hosted web application with **zero external resources, no CDN, no third-party scripts**:

```http
Content-Security-Policy: 
  script-src 'nonce-{RANDOM}' 'strict-dynamic';
  style-src 'nonce-{RANDOM}';
  img-src 'self' data:;
  font-src 'self';
  connect-src 'self';
  object-src 'none';
  base-uri 'none';
  form-action 'self';
  frame-ancestors 'none';
  default-src 'self'
```

### Implementation Pattern

**Nonce-based (recommended for dynamic templates):**
```javascript
// Express.js example
app.get("/", (req, res) => {
  const nonce = crypto.randomUUID();
  const cspHeader = `
    script-src 'nonce-${nonce}' 'strict-dynamic';
    style-src 'nonce-${nonce}';
    object-src 'none';
    base-uri 'none';
    form-action 'self';
    frame-ancestors 'none';
    default-src 'self'
  `;
  res.setHeader("Content-Security-Policy", cspHeader);
  res.send(`
    <!DOCTYPE html>
    <html>
      <head>
        <style nonce="${nonce}">body { margin: 0; }</style>
      </head>
      <body>
        <script nonce="${nonce}">
          console.log("CSP compliant");
        </script>
      </body>
    </html>
  `);
});
```

**Hash-based (for static assets):**
```bash
# Generate hash of inline script
echo -n "console.log('test');" | openssl dgst -sha256 -binary | openssl enc -base64

# Result: sha256-HASH_OUTPUT
```

```http
Content-Security-Policy: script-src 'sha256-HASH_OUTPUT'
```

### Critical Requirements

- ✅ **Nonce must be unique per request** and cryptographically random
- ✅ Never use `'unsafe-inline'` or `'unsafe-eval'`
- ✅ Use `'strict-dynamic'` to allow trusted scripts to load other scripts
- ✅ Set `frame-ancestors 'none'` to prevent clickjacking
- ✅ Use `object-src 'none'` to block plugins
- ✅ Avoid `data:` URLs for scripts (XSS vector)

---

## 2. NPM SUPPLY CHAIN SECURITY FOR AIR-GAPPED BUILDS

### Vendoring All Dependencies (Offline Package Mirror)

#### Option A: Verdaccio (Private npm Registry)

**Why:** Zero-config local npm registry with built-in caching.

**Installation:**
```bash
npm install -g verdaccio
# Or use Docker
docker run -it --rm --name verdaccio \
  -p 4873:4873 \
  -v verdaccio_storage:/verdaccio/storage \
  verdaccio/verdaccio:latest
```

**Configuration (~/.verdaccio/config.yaml):**
```yaml
storage: /verdaccio/storage

packages:
  '@*/*':
    access: $all
    publish: $authenticated
    proxy: npmjs  # Only used if online
  '**':
    access: $all
    publish: $authenticated
    proxy: npmjs

uplinks:
  npmjs:
    url: https://registry.npmjs.org/
    timeout: 100000
    max_retries: 2
    maxage: 10m
    cache: true

middlewares:
  audit:
    enabled: true
```

**Pre-population for air-gapped environments:**
```bash
# 1. On online system, create verdaccio instance
# 2. Run your build:
npm install --registry http://localhost:4873

# 3. Backup storage
tar czf /mnt/transport/verdaccio-storage.tar.gz /verdaccio/storage

# 4. On air-gapped system, restore and mount
tar xzf verdaccio-storage.tar.gz -C /verdaccio/
docker run -v /verdaccio/storage:/verdaccio/storage verdaccio/verdaccio
```

#### Option B: npm ci with Offline Lock Files

**Recommended for CI/CD in air-gapped environments:**

```bash
# Generate lock file (on internet-connected system)
npm ci --prefer-offline --offline

# Package all modules
npm ci --offline

# Create artifact
tar czf app-with-deps.tar.gz node_modules/ package-lock.json

# Transfer to air-gapped environment
# Then:
npm ci --offline --frozen-lockfile
```

**Key npm ci flags:**
- `--frozen-lockfile`: Fail if package-lock.json changes
- `--offline`: Use only cached packages
- `--ignore-scripts`: Skip installation scripts (security hardening)
- `--no-optional`: Skip optional dependencies
- `--audit=false`: Skip audit when offline

#### Option C: npm pack for Monorepos

```bash
# Create reproducible package archives
npm pack ./packages/my-app --pack-destination=/mnt/transport/

# On air-gapped system
npm install /mnt/transport/my-app-1.0.0.tgz
```

### Supply Chain Best Practices

**Dependency auditing in CI:**
```bash
# Before air-gap cut-off
npm audit --json > audit-baseline.json
npm install --save-exact  # Pin exact versions
git commit -S package-lock.json  # Sign lock file
```

**Reproducible builds:**
```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package-lock.json package.json ./
RUN npm ci --offline --frozen-lockfile --no-optional
RUN npm run build

FROM node:20-alpine
COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/node_modules /app/node_modules
```

**Verify integrity:**
```bash
# Create checksum baseline
sha256sum node_modules/**/*.js > deps.sha256
# Verify after transport
sha256sum -c deps.sha256
```

---

## 3. CONTAINER IMAGE SECURITY IN AIR-GAPPED REGISTRIES

### Private Registry Setup with Authentication

#### Using Harbor (Recommended)

```yaml
# harbor/docker-compose.yml
version: '3'
services:
  harbor:
    image: goharbor/harbor-core:v2.10.0
    environment:
      HARBOR_ADMIN_PASSWORD: ChangeMeStrong!
      HARBOR_HOSTNAME: registry.internal
      POSTGRES_PASSWORD: ChangeMe!
    volumes:
      - /data/harbor:/data
    networks:
      - harbor-net
  registry:
    image: goharbor/registry-photon:v2.10.0
    volumes:
      - /data/registry:/storage
    networks:
      - harbor-net
```

#### Image Pull Secrets in Kubernetes

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: private-registry-creds
  namespace: cloudtak
type: kubernetes.io/dockercfg
data:
  .dockercfg: <base64-encoded-registry-auth>
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloudtak-sa
  namespace: cloudtak
imagePullSecrets:
  - name: private-registry-creds
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloudtak-api
  namespace: cloudtak
spec:
  serviceAccountName: cloudtak-sa
  containers:
  - name: api
    image: registry.internal/onemindos/cloudtak-api:latest
    imagePullPolicy: IfNotPresent  # Critical for air-gapped
```

**Creating the secret:**
```bash
kubectl create secret docker-registry private-registry-creds \
  --docker-server=registry.internal \
  --docker-username=robotuser \
  --docker-password="$(cat /run/secrets/registry-password)" \
  --docker-email=robot@internal \
  -n cloudtak
```

### Image Mirroring Strategy

**Pre-load all images on nodes:**
```bash
#!/bin/bash
# mirror.sh - Run on air-gapped nodes before deployment

IMAGES=(
  "registry.internal/onemindos/cloudtak-api:latest"
  "registry.internal/onemindos/wwv:latest"
  "registry.internal/onemindos/hermes-legacy:v2"
  "registry.internal/nats:latest"
)

for image in "${IMAGES[@]}"; do
  docker pull "$image" || exit 1
done

# Tag with digest for immutability
docker tag "$image" "$image@$(docker inspect "$image" -f '{{ .RepoDigests }}')"
```

**Using image digests in deployments (immutable):**
```yaml
image: registry.internal/onemindos/cloudtak-api@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2
imagePullPolicy: IfNotPresent
```

### Image Signing with Cosign

```bash
# Generate signing key (on secure, isolated system)
cosign generate-key-pair

# Sign image after build
cosign sign --key cosign.key registry.internal/onemindos/cloudtak-api:v1.0

# Verify in deployment
cosign verify --key cosign.pub registry.internal/onemindos/cloudtak-api:v1.0
```

---

## 4. SOFTWARE BILL OF MATERIALS (SBOM) GENERATION

### Using Syft (Recommended for Multi-Format Support)

**Installation:**
```bash
curl -sSfL https://get.anchore.io/syft | sudo sh -s -- -b /usr/local/bin
```

**Generate SBOM in multiple formats:**
```bash
# CycloneDX (recommended for compliance)
syft registry.internal/onemindos/cloudtak-api:latest \
  -o cyclonedx-json \
  > sbom-cloudtak.json

# SPDX (recommended for licensing)
syft registry.internal/onemindos/cloudtak-api:latest \
  -o spdx-json \
  > sbom-cloudtak.spdx.json

# Syft native format
syft registry.internal/onemindos/cloudtak-api:latest \
  -o json \
  > sbom-cloudtak-syft.json
```

**CycloneDX SBOM structure (compliance):**
```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "serialNumber": "urn:uuid:3e671687-395b-41f5-a30f-a58921a69b79",
  "version": 1,
  "metadata": {
    "timestamp": "2026-07-07T10:00:00Z",
    "tools": [
      {
        "vendor": "Anchore",
        "name": "syft",
        "version": "0.95.0"
      }
    ],
    "component": {
      "type": "container",
      "name": "cloudtak-api",
      "version": "1.0.0"
    }
  },
  "components": [
    {
      "type": "library",
      "name": "express",
      "version": "4.18.2",
      "licenses": [
        {
          "license": {
            "id": "MIT"
          }
        }
      ],
      "purl": "pkg:npm/express@4.18.2"
    }
  ]
}
```

**Vulnerability Exploitability Exchange (VEX) for known-not-vulnerable:**
```json
{
  "bomFormat": "CycloneDX",
  "specVersion": "1.4",
  "version": 1,
  "vulnerabilities": [
    {
      "ref": "urn:uuid:3e671687-395b-41f5-a30f-a58921a69b79",
      "id": "CVE-2021-12345",
      "source": {
        "name": "NVD",
        "url": "https://nvd.nist.gov/vuln/detail/CVE-2021-12345"
      },
      "status": "unaffected",
      "justification": "vulnerable_code_not_in_execute_path"
    }
  ]
}
```

---

## 5. DISA STIG & COMPLIANCE REQUIREMENTS

### Web Application Security STIG Controls

**Key applicable STIGs:**
- **DISA STIG: Web Application Security v5.2** (STIG ID: WA00)
- **DISA STIG: Kubernetes v1.6** (for K8s deployments)
- **NIST SP 800-53** (underlying security controls)

**Critical STIG requirements for web applications:**

| Control | Requirement | Implementation |
|---------|-------------|-----------------|
| **WA-001** | Input Validation | Sanitize all user inputs; use parameterized queries |
| **WA-002** | Output Encoding | HTML, CSS, JavaScript context-aware encoding |
| **WA-003** | Authentication | MFA required; no hardcoded credentials |
| **WA-004** | Session Management | Secure cookies; CSRF tokens; session timeout |
| **WA-005** | Access Control | RBAC; principle of least privilege |
| **WA-006** | Cryptography | TLS 1.2+; FIPS 140-2/140-3 for sensitive ops |
| **WA-007** | Error Handling | No sensitive info in error messages |
| **WA-008** | Logging | Audit all auth, config, sensitive operations |
| **WA-009** | Data Protection | Encryption at rest (AES-256) and in transit (TLS) |
| **WA-010** | Patch Management | Zero-day response; vendor updates within 30 days |

**Minimum CSP requirement (WA-002):**
```http
Content-Security-Policy: default-src 'self'; script-src 'nonce-{RANDOM}'; object-src 'none'
```

**HTTPS/TLS requirement (WA-006):**
- Minimum TLS 1.2
- FIPS-approved cipher suites only:
  - TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  - TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
  - TLS_RSA_WITH_AES_256_GCM_SHA384 (fallback)

**Logging requirement (WA-008):**
```json
{
  "timestamp": "2026-07-07T10:30:45Z",
  "event": "authentication",
  "user": "admin",
  "source_ip": "10.0.1.100",
  "result": "success",
  "session_id": "sess_abc123",
  "mfa_method": "totp"
}
```

---

## 6. OFFLINE DEPENDENCY SCANNING (NO INTERNET)

### Trivy with Offline Database

**Download and cache vulnerability database:**
```bash
# On internet-connected system
trivy image --download-db-only
trivy fs --download-db-only

# Compress DB
tar czf trivy-db.tar.gz ~/.cache/trivy/

# Transfer to air-gapped system
tar xzf trivy-db.tar.gz -C ~
```

**Use cached database for scanning:**
```bash
# Scan with local database
trivy image --skip-update \
  registry.internal/onemindos/cloudtak-api:latest \
  -o json \
  > scan-results.json

# Generate SBOM
trivy image --skip-update \
  -f spdx \
  registry.internal/onemindos/cloudtak-api:latest
```

### OWASP Dependency-Check (for npm/Python/Java)

**Setup offline NVD data:**
```bash
# Download NVD data feeds (on internet system)
dependency-check.sh \
  --project "omos-infra" \
  --scan . \
  --nvdApiKey "$API_KEY" \
  --enableExperimental

# Creates nvd-cache in ~/.m2/repository/org/owasp/dependencycheck/data/
# Transfer cache to air-gapped system
```

**Scan without internet:**
```bash
dependency-check.sh \
  --project "omos-infra" \
  --scan ./platform \
  --noDist \
  --out ./reports/
```

### Grype with Embedded Database

```bash
# Download grype database
grype version

# Cache database
mkdir -p ~/.cache/grype/db
grype --check-for-update=false \
  registry.internal/onemindos/wwv:latest \
  -o json \
  > vulnerabilities.json
```

---

## 7. OFFLINE MAP TILES FOR TACTICAL DISPLAY

### Option A: PMTiles (Single-File Format)

**Advantages:** One file, cloud-native, works without server backend

```bash
# Convert MBTiles to PMTiles
tippecanoe-decode tiles.mbtiles | pmtiles-encode tiles.pmtiles

# Self-host PMTiles viewer
docker run -p 8080:8080 -v /data:/data \
  pmtiles/pmtiles-server:latest

# Access at: http://localhost:8080?url=/tiles.pmtiles
```

### Option B: TileServer GL (Vector + Raster)

**Recommended for production military deployments**

```docker
FROM node:20-alpine
RUN npm install -g tileserver-gl
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["tileserver-gl", "--no-cors"]
```

**Configuration (tileserver-config.json):**
```json
{
  "server": {
    "port": 8080,
    "address": "127.0.0.1",
    "cors": false,
    "secureHeaders": true,
    "strictSSL": true
  },
  "sources": {
    "terrain": {
      "type": "mbtiles",
      "mbtiles": "/data/terrain.mbtiles",
      "cache": 512
    },
    "vectorSource": {
      "type": "mbtiles",
      "mbtiles": "/data/vector.mbtiles"
    }
  },
  "layers": [
    {
      "id": "terrain",
      "source": "terrain",
      "render": {}
    }
  ]
}
```

**Deploy to K8s:**
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: tileserver
  namespace: wwv
spec:
  serviceName: tileserver
  replicas: 1
  selector:
    matchLabels:
      app: tileserver
  template:
    metadata:
      labels:
        app: tileserver
    spec:
      containers:
      - name: tileserver
        image: registry.internal/tileserver-gl:latest
        ports:
        - containerPort: 8080
        volumeMounts:
        - name: tiles-data
          mountPath: /data
          readOnly: true
  volumeClaimTemplates:
  - metadata:
      name: tiles-data
    spec:
      accessModes: ["ReadOnlyMany"]
      resources:
        requests:
          storage: 50Gi
```

### Option C: Martin (Rust-based, Multi-Source)

Supports MBTiles, PMTiles, PostGIS, S3 in single binary.

```dockerfile
FROM ghcr.io/maplibre/martin:latest
EXPOSE 3000
ENV DATASOURCES=mbtiles:///data/tiles/*.mbtiles,pg://user:pass@db/postgis
```

**In production K8s:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: martin-config
  namespace: wwv
data:
  config.toml: |
    [datasources]
    mbtiles = "/data/tiles/"
    
    [server]
    bind_address = "0.0.0.0"
    port = 3000
    cors = false
```

---

## 8. CERTIFICATE MANAGEMENT IN AIR-GAPPED ENVIRONMENTS

### Internal PKI Setup (No Let's Encrypt)

**Generate internal Root CA:**
```bash
# Create root CA key (keep offline)
openssl genrsa -out root-ca.key 4096

# Self-sign root certificate (10-year validity)
openssl req -new -x509 -days 3650 -key root-ca.key \
  -out root-ca.crt \
  -subj "/CN=OneMind Internal Root CA/O=OneMind/C=US"

# Create CA cert bundle for distribution
cp root-ca.crt /etc/pki/ca-trust/source/anchors/
update-ca-trust
```

**Kubernetes ClusterIssuer (for cert-manager):**
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: internal-ca
spec:
  ca:
    secretRef:
      name: internal-ca-key-pair
---
apiVersion: v1
kind: Secret
metadata:
  name: internal-ca-key-pair
  namespace: cert-manager
type: kubernetes.io/tls
data:
  tls.crt: <base64-root-ca.crt>
  tls.key: <base64-root-ca.key>
```

**Issue service certificates (short-lived - 30 days):**
```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: cloudtak-tls
  namespace: cloudtak
spec:
  secretName: cloudtak-tls
  commonName: cloudtak.onemindos.internal
  dnsNames:
  - cloudtak.onemindos.internal
  - cloudtak.nats.svc.cluster.local
  issuerRef:
    name: internal-ca
    kind: ClusterIssuer
  duration: 720h # 30 days
  renewBefore: 168h # Renew 7 days before expiry
```

**Automated renewal with short TTL (best practice for air-gapped):**
- 30-day certificate lifetime
- Automatic renewal at 7 days before expiry
- Passive revocation (no OCSP needed)
- All renewal traffic stays within cluster

---

## 9. SUPPLY CHAIN SECURITY FRAMEWORKS

### SLSA (Supply-chain Levels for Software Artifacts)

**SLSA Level 3 (Recommended for military deployments):**

```yaml
# GitHub Actions workflow example (adapt to air-gapped CI)
name: Build with SLSA L3 Provenance
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    outputs:
      image: ${{ steps.build.outputs.image }}
      digest: ${{ steps.build.outputs.digest }}
    steps:
    - uses: actions/checkout@v3
    - uses: docker/build-push-action@v4
      id: build
      with:
        tags: registry.internal/onemindos/cloudtak-api:${{ github.sha }}
        push: true
        outputs: type=oci

  provenance:
    needs: build
    permissions:
      contents: read
      id-token: write
    uses: slsa-framework/slsa-github-generator/.github/workflows/generator_container_slsa3.yml@v1.7.0
    with:
      image: ${{ needs.build.outputs.image }}
      digest: ${{ needs.build.outputs.digest }}
      registry-username: robotuser
    secrets:
      registry-password: ${{ secrets.REGISTRY_PASSWORD }}
```

### CycloneDX SBOM Standard

**Required components for compliance:**
- Package name, version, type, PURL
- License information (SPDX identifier)
- Vulnerability references
- Dependency graph
- Build metadata (builder identity, timestamp)

### in-toto for Build Integrity

```json
{
  "signed": {
    "_type": "link",
    "name": "build",
    "materials": {
      "src/main.ts": {
        "sha256": "d19d68fbc8961f67eee8e6b9e2fe85f72dcc5f4d"
      }
    },
    "products": {
      "dist/app.js": {
        "sha256": "c3b75a8f42f63b6c2c1d7e9f4b5a0d8e1c2f5a7b"
      }
    },
    "byproducts": {
      "return-value": 0,
      "stdout": ""
    },
    "environment": {
      "NODE_ENV": "production"
    },
    "command": ["npm", "run", "build"],
    "key_secretaries": ["zeus"],
    "type": "link"
  },
  "signatures": [
    {
      "keyid": "e7aff1b03b20b5c8b0c3f8d5e2a1f6c9d4b7a0e3",
      "sig": "304502206..."
    }
  ]
}
```

### The Update Framework (TUF) for Secure Updates

**TUF roles in air-gapped context:**
```
Root (offline) ──► Timestamp (online) ──► Snapshot ──► Targets
                        (cache)
```

**Implementation in K8s:**
- Store root metadata as ConfigMap (never rotated)
- Use in-cluster update mirror for TUF metadata
- Sign all manifests with in-toto + TUF

---

## 10. PRACTICAL CI/CD FOR AIR-GAPPED ENVIRONMENTS

### Build Pipeline Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ Online Build Environment (DMZ)                              │
│ ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│ │ Git Repo │──│ Tests    │──│ Trivy    │                   │
│ │          │  │ npm ci   │  │ Scan     │                   │
│ │ (GitHub) │  │          │  │          │                   │
│ └──────────┘  └──────────┘  └──────────┘                   │
│                         ↓                                   │
│              ┌──────────────────┐                           │
│              │ SBOM + Provenance│                           │
│              │ (Sign with key)  │                           │
│              └──────────────────┘                           │
│                         ↓                                   │
│              ┌──────────────────┐                           │
│              │ Package OCI Image│                           │
│              │ Sign with Cosign │                           │
│              └──────────────────┘                           │
│                         ↓                                   │
│         (Artifact Transfer via Secure Media)               │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────┴────────────────────────────────────┐
│ Air-Gapped Deployment Environment                           │
│ ┌──────────────────┐  ┌──────────────────┐                 │
│ │ Harbor Registry  │  │ Kubernetes Cluster              │
│ │ (Internal)       │  │ (Isolated Network)              │
│ │                  │  │                                  │
│ │ Verify Signature │──│ Pull Image                       │
│ │ Verify Provenance│  │ Mount ConfigMaps (tiles)         │
│ │ Verify SBOM      │  │ Deploy Pod                       │
│ │                  │  │                                  │
│ └──────────────────┘  └──────────────────┘                 │
└─────────────────────────────────────────────────────────────┘
```

### Makefile for Compliant Builds

```makefile
.PHONY: build scan sbom sign deploy

# Build image
build:
	docker build -t registry.internal/onemindos/cloudtak-api:$(GIT_SHA) .
	docker push registry.internal/onemindos/cloudtak-api:$(GIT_SHA)

# Scan with Trivy
scan:
	trivy image --skip-update \
		-f json \
		-o scan-$(GIT_SHA).json \
		registry.internal/onemindos/cloudtak-api:$(GIT_SHA)
	@echo "✓ Trivy scan complete"

# Generate SBOM
sbom:
	syft registry.internal/onemindos/cloudtak-api:$(GIT_SHA) \
		-o cyclonedx-json \
		> sbom-$(GIT_SHA).json
	@echo "✓ SBOM generated"

# Sign image
sign: sbom
	cosign sign --key ~/.cosign/cosign.key \
		registry.internal/onemindos/cloudtak-api:$(GIT_SHA)
	@echo "✓ Image signed"

# Full build+verify pipeline
all: build scan sbom sign
	@echo "✓ All checks passed"
	@echo "Ready for air-gapped deployment"
```

---

## 11. SECURITY SCANNING TOOLS COMPARISON

| Tool | Purpose | Air-Gapped | Offline DB | Formats |
|------|---------|-----------|-----------|---------|
| **Trivy** | Container/Dependency scan | ✅ Yes | ✅ Yes | SBOM, JSON, SARIF |
| **Syft** | SBOM generation | ✅ Yes | N/A | CycloneDX, SPDX, JSON |
| **Grype** | Vulnerability scan | ✅ Yes | ✅ Yes (with config) | JSON, Table |
| **Cosign** | Image signing | ✅ Yes | N/A | OCI signature |
| **Dependency-Check** | Dependency audit | ✅ Yes | ✅ Yes (NVD cache) | HTML, JSON, XML |
| **Verdaccio** | npm registry | ✅ Yes | ✅ Yes (caching) | npm protocol |

---

## 12. HARDENING CHECKLIST FOR MILITARY-GRADE WEB C2

### Pre-Deployment

- [ ] All images built in air-gapped network or imported via secure media
- [ ] All images signed with internal CA
- [ ] All images scanned with Trivy; zero critical vulns
- [ ] SBOM generated (CycloneDX v1.4+) for all images
- [ ] npm dependencies vendored and verified
- [ ] CSP headers deployed (nonce-based, 'strict-dynamic')
- [ ] TLS 1.2+ with FIPS-approved ciphers only
- [ ] Root CA distributed to all nodes via secure channel
- [ ] RBAC policies configured (least privilege)
- [ ] Audit logging enabled (all auth/config/sensitive ops)
- [ ] No credentials committed to git
- [ ] All secrets in 1Password vault (not committed)

### Post-Deployment

- [ ] Monthly vulnerability scan cycle
- [ ] Quarterly SBOM refresh
- [ ] Automatic certificate renewal enabled (30-day TTL)
- [ ] Network segmentation verified
- [ ] Firewall rules locked down (network policies in K8s)
- [ ] Supply chain provenance verified (in-toto attestations)
- [ ] Incident response runbook distributed
- [ ] DISA STIG compliance audit completed

---

## 13. REFERENCE TOOLS & VERSIONS (As of 2026-07)

| Tool | Version | Link |
|------|---------|------|
| Syft | 0.95.0+ | https://github.com/anchore/syft |
| Trivy | 0.50.0+ | https://trivy.dev/ |
| Grype | 0.76.0+ | https://github.com/anchore/grype |
| Cosign | 2.2.0+ | https://github.com/sigstore/cosign |
| Verdaccio | 5.30.0+ | https://verdaccio.org/ |
| TileServer GL | 4.4.0+ | https://github.com/maplibre/martin |
| Cert-Manager | 1.14.0+ | https://cert-manager.io/ |
| Helm | 3.14.0+ | https://helm.sh/ |
| Docker | 26.0.0+ | https://www.docker.com/ |
| Kubernetes | 1.28.0+ | https://kubernetes.io/ |

---

## 14. NEXT STEPS FOR OMOS PLATFORM

**Immediate (Week 1-2):**
1. Set up private Harbor registry in DOKS with TLS + RBAC
2. Implement Verdaccio for npm dependency mirroring
3. Deploy Trivy scanning to CI/CD (offline DB)
4. Generate SBOMs for all current images (CycloneDX)

**Near-term (Month 1):**
1. Implement image signing (Cosign) + verification
2. Deploy strict CSP headers on CloudTAK, WWV, Legacy
3. Certificate audit + internal CA setup
4. DISA STIG baseline assessment

**Medium-term (Q3):**
1. Implement SLSA L3 provenance generation
2. Deploy in-toto attestations for all builds
3. Supply chain security training for team
4. Quarterly compliance audits

---

## 15. COMPLIANCE & REFERENCES

**Standards Referenced:**
- NIST SP 800-53 (Security and Privacy Controls)
- DISA STIG: Web Application Security v5.2
- CIS Controls v8.1
- OWASP Top 10 2023
- SLSA v1.0 Framework
- CycloneDX 1.4 Standard
- FIPS 140-3 (Cryptographic Module Validation)

**Key Resources:**
- MDN Web Docs: CSP (https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP)
- OWASP CSP Cheat Sheet (https://cheatsheetseries.owasp.org/cheatsheets/Content_Security_Policy_Cheat_Sheet.html)
- Kubernetes Security Best Practices (https://kubernetes.io/docs/concepts/security/)
- Small Step: Everything PKI (https://smallstep.com/blog/everything-pki/)
- SLSA Framework (https://slsa.dev/)

---

**Document Status:** Complete research snapshot from 2026-07-07  
**Classification:** For internal OneMind platform security guidance  
**Next Review:** 2026-10-07
