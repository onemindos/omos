# CloudTAK Deployment — DOKS (onemind-hub nyc3)

**Goal:** Stand up CloudTAK as the browser COP next to our TAK Server 5.7 cluster.
CloudTAK connects to TAK Server via a Connection (mTLS cert pair) configured in the admin UI.

## Architecture on DOKS

```
namespace: cloudtak              namespace: takserver
┌───────────────────────┐       ┌──────────────────────────┐
│ cloudtak-api (5000)   │─cert─▶│ takserver-api (8443/8446)│
│ cloudtak-tiles (5002) │       │ takserver-messaging(8089)│
│ cloudtak-media (9997) │       └──────────────────────────┘
│ cloudtak-events       │
│ cloudtak-postgis      │
│ cloudtak-minio        │
└───────────────────────┘
```

## Prerequisites

- TAK Server 5.7 live on DOKS (`takserver` namespace) ✅
- DOCR registry `registry.digitalocean.com/onemind-hub` ✅
- `doctl`, `kubectl` (context `do-nyc3-onemind-hub`), `docker buildx` ✅

## Step 1: Build & push CloudTAK images (amd64)

```bash
export PATH="/opt/homebrew/bin:$PATH"
cd /Users/zeus/CloudTAK
REG=registry.digitalocean.com/onemind-hub
doctl registry login

# API (main app — includes web client)
docker buildx build --platform linux/amd64 \
  -f api/Dockerfile -t $REG/cloudtak-api:latest --push ./api/

# Tiles (PMTiles server)
docker buildx build --platform linux/amd64 \
  -f tasks/pmtiles/Dockerfile.compose -t $REG/cloudtak-tiles:latest --push .

# Events (webhook processor)
docker buildx build --platform linux/amd64 \
  -f tasks/events/Dockerfile -t $REG/cloudtak-events:latest --push .
```

Note: media uses a public GHCR image (`ghcr.io/dfpc-coe/media-infra:v9.1.1`) — no build needed.
PostGIS and MinIO use public images — no build needed.

## Step 2: Deploy to DOKS

```bash
CTX=do-nyc3-onemind-hub

# Create namespace + pull secret
kubectl --context $CTX apply -f namespace.yaml
doctl registry kubernetes-manifest --namespace cloudtak | kubectl --context $CTX apply -f -

# Generate a real signing secret
SIGNING=$(openssl rand -hex 32)
sed -i '' "s/CHANGE-ME-RANDOM-64/$SIGNING/" secrets.yaml

# Apply all manifests
kubectl --context $CTX apply -f secrets.yaml
kubectl --context $CTX apply -f postgis.yaml
kubectl --context $CTX apply -f minio.yaml
kubectl --context $CTX apply -f media.yaml
kubectl --context $CTX apply -f tiles.yaml
kubectl --context $CTX apply -f events.yaml
kubectl --context $CTX apply -f api.yaml
```

## Step 3: Expose via Cloudflare tunnel (or ingress)

Create a Cloudflare tunnel route:
- `cloudtak.onemindos.dev` → `http://cloudtak-api.cloudtak.svc.cluster.local:5000`
- `tiles.cloudtak.onemindos.dev` → `http://cloudtak-tiles.cloudtak.svc.cluster.local:5002`
- `video.onemindos.dev` → `http://cloudtak-media.cloudtak.svc.cluster.local:9997`

Or use ingress-nginx with appropriate annotations.

## Step 4: Connect CloudTAK to TAK Server

1. Open `https://cloudtak.onemindos.dev` → login (first user becomes admin)
2. Admin panel → Connections → New Connection
3. Upload a client cert pair from our TAK CA (`certs-out/admin.pem` + key)
   - Point to: `ssl://takserver-api.takserver.svc.cluster.local:8089`
   - (Or external LB if cross-namespace TLS routing needed)
4. Once connected, CloudTAK will stream all CoT from TAK Server and display on the map.

## Step 5: Deploy AIS ETL connector

See `etl-aishub/` directory (separate deployment that pushes CoT into TAK Server via Connection).

## Verify

- CloudTAK shows the globe with connected TAK Server contacts
- Same contacts visible on ATAK/iTAK field devices
- `cot-tap` (already running in takserver namespace) bridges all CoT → NATS
