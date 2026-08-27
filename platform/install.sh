#!/usr/bin/env bash
# =============================================================================
# OneMind Hub — Platform Install Script
# Reproduces the full DOKS cloud hub from omos-infra git alone.
# Run from repo root. Requires: helm, kubectl, doctl (DO CLI), op (1Password CLI)
# =============================================================================
set -euo pipefail

CTX="${KUBE_CONTEXT:-tailscale-operator.tail717ba0.ts.net}"
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PLATFORM="$REPO_ROOT/platform"

echo "==> Using cluster context: $CTX"

# -----------------------------------------------------------------------------
# 0. SEALED SECRETS CONTROLLER
# All K8s secrets are encrypted SealedSecret YAMLs committed to this repo.
# The controller (kube-system) decrypts them using the cluster private key.
# Must run before any step that applies SealedSecret manifests.
#
# To reseal a rotated secret:
#   kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
#     --controller-namespace=kube-system > /tmp/pub-cert.pem
#   kubectl create secret generic <name> -n <ns> --from-literal=key=value \
#     --dry-run=client -o json \
#     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
#     > platform/<ns>/<name>.sealedsecret.yaml
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/sealed-secrets/controller.yaml"
kubectl --context "$CTX" wait --for=condition=available deployment/sealed-secrets-controller \
  -n kube-system --timeout=120s

# -----------------------------------------------------------------------------
# 1. CORE PLATFORM — cert-manager, ingress-nginx, ArgoCD
# -----------------------------------------------------------------------------
helm repo add jetstack https://charts.jetstack.io --force-update
helm upgrade --install cert-manager jetstack/cert-manager \
  --kube-context "$CTX" -n cert-manager --create-namespace \
  -f "$PLATFORM/values/cert-manager.yaml"

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx --force-update
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --kube-context "$CTX" -n ingress-nginx --create-namespace \
  -f "$PLATFORM/values/ingress-nginx.yaml"

helm repo add argo https://argoproj.github.io/argo-helm --force-update
helm upgrade --install argocd argo/argo-cd \
  --kube-context "$CTX" -n argocd --create-namespace
# ArgoCD uses all defaults — no custom values needed

# -----------------------------------------------------------------------------
# 2. TLS ISSUERS (cert-manager ClusterIssuers)
# cloudflare-api-token is a SealedSecret in platform/cert-manager/
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/issuer-dns.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/ingress/issuers.yaml"

# -----------------------------------------------------------------------------
# 3. CLOUDFLARE TUNNEL (zero public inbound)
# Routes are declared in platform/cloudflare/config.yaml — edit + apply to add new routes.
# Requires tunnel-creds-json secret (SealedSecret — see platform/cloudflare/tunnel-creds-json.sealedsecret.yaml):
#   TOKEN=$(op read "op://00-24 SO/tunnel-credentials/tunnel-token")
#   python3 -c "
#     import sys,json,base64; t='$TOKEN'; p=t+'='*(4-len(t)%4)
#     d=json.loads(base64.urlsafe_b64decode(p))
#     print(json.dumps({'AccountTag':d['a'],'TunnelID':d['t'],'TunnelSecret':d['s']}))
#   " | kubectl create secret generic tunnel-creds-json -n cloudflare --from-file=credentials.json=/dev/stdin
# -----------------------------------------------------------------------------
kubectl --context "$CTX" create namespace cloudflare --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/cloudflare/config.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudflare/cloudflared.yaml"

# -----------------------------------------------------------------------------
# 4. NATS — 3-node HA JetStream cluster
# Requires: nats/values.secret.yaml (has ONEMIND/SYS passwords — inject from 1Password)
# -----------------------------------------------------------------------------
helm repo add nats https://nats-io.github.io/k8s/helm/charts --force-update
kubectl --context "$CTX" create namespace nats --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/nats/certs.yaml"
helm upgrade --install nats nats/nats \
  --kube-context "$CTX" -n nats \
  -f "$PLATFORM/nats/values.secret.yaml"  # inject from 1Password on fresh deploy
kubectl --context "$CTX" apply -f "$PLATFORM/nats/omos-nats-tailnet-service.yaml"
# Note: NATS WebSocket must be enabled in nats/values.secret.yaml:
#   websocket:
#     enabled: true
#     port: 8080
# NATS WebSocket is exposed as wss://ws.onemindos.dev via Cloudflare tunnel (platform/cloudflare/config.yaml)

# -----------------------------------------------------------------------------
# 4b. NACK — JetStream CRD controller (Stream + KeyValue declarative mgmt)
# Requires: nack-values.secret.yaml (gitignored) — see nack-values.template.yaml
# Generate:
#   ONEMIND_PW=$(grep password "$PLATFORM/nats/values.secret.yaml" | awk '{print $2}' | head -1)
#   # or: ONEMIND_PW=$(op read "op://00-24 SO/nats-onemind/password")
#   cat > "$PLATFORM/nats/nack-values.secret.yaml" <<EOF
#   jetstream:
#     nats:
#       url: "nats://onemind:${ONEMIND_PW}@nats.nats.svc.cluster.local:4222"
#   EOF
# -----------------------------------------------------------------------------
helm upgrade --install nack nats/nack \
  --kube-context "$CTX" -n nats \
  -f "$PLATFORM/nats/nack-values.yaml" \
  -f "$PLATFORM/nats/nack-values.secret.yaml"  # inject ONEMIND_PW (see template)

# Apply all Stream + KeyValue CRD manifests — controller reconciles them to NATS
kubectl --context "$CTX" apply -f "$PLATFORM/nats/streams/"

# -----------------------------------------------------------------------------
# 5. TAILSCALE OPERATOR (kubectl-over-tailnet + service exposure)
# Requires: tailscale-operator.yaml with real clientSecret from 1Password
# op read "op://00-24 SO/tailscale-oauth-operator/credential"
# -----------------------------------------------------------------------------
helm repo add tailscale https://pkgs.tailscale.com/helmcharts --force-update
helm upgrade --install tailscale-operator tailscale/tailscale-operator \
  --kube-context "$CTX" -n tailscale --create-namespace \
  -f "$PLATFORM/values/tailscale-operator.yaml"  # clientSecret injected at deploy time
kubectl --context "$CTX" apply -f "$PLATFORM/tailscale/apiserver-rbac.yaml"

# -----------------------------------------------------------------------------
# 6. OBSERVABILITY — kube-prometheus-stack, Loki, Tempo, OTel Collector
# -----------------------------------------------------------------------------
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
helm upgrade --install monitoring prometheus-community/kube-prometheus-stack \
  --kube-context "$CTX" -n monitoring --create-namespace \
  -f "$PLATFORM/values/kube-prometheus-stack.yaml"

helm repo add grafana https://grafana.github.io/helm-charts --force-update
helm upgrade --install loki grafana/loki \
  --kube-context "$CTX" -n monitoring \
  -f "$PLATFORM/values/loki.yaml"

helm upgrade --install tempo grafana/tempo \
  --kube-context "$CTX" -n monitoring \
  -f "$PLATFORM/values/tempo.yaml"

helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts --force-update
helm upgrade --install otel-collector open-telemetry/opentelemetry-collector \
  --kube-context "$CTX" -n monitoring \
  -f "$PLATFORM/values/otel-collector.yaml"

# Apply custom datasource ConfigMap (Loki + Tempo linked in Grafana)
kubectl --context "$CTX" apply -f "$PLATFORM/grafana-datasources.yaml"

# -----------------------------------------------------------------------------
# 6b. NATS SURVEYOR — Prometheus exporter for the NATS control plane
# Requires: nats-server-tls secret (CA cert from cert-manager) and $SYS creds
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/nats-surveyor/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/nats-surveyor/servicemonitor.yaml"

# -----------------------------------------------------------------------------
# 8. OMOS SERVICES — shared omos namespace (tak-projector, feed-world, etc.)
# Secrets (omos-nats-creds) are SealedSecrets in platform/omos/
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/omos/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/tak-projector.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos-mission-engine/deployment.yaml"

# -----------------------------------------------------------------------------
# 8b. OMOS WORLD FEEDS — live global data → NATS → TAK → OMOS ATOC map
# Sources: AIS ships, ADS-B aircraft, APRS trackers, satellites, hazards
# Requires:
#   omos-nats-creds secret in omos namespace (already created above)
#   omos-feed-secrets (AISSTREAM_API_KEY, OPENSKY_USERNAME, OPENSKY_PASSWORD):
#     kubectl create secret generic omos-feed-secrets -n omos \
#       --from-literal=AISSTREAM_API_KEY='<key from aisstream.io — free signup>' \
#       --from-literal=OPENSKY_USERNAME='<opensky.org free account>' \
#       --from-literal=OPENSKY_PASSWORD='<opensky password>'
# Toggle feeds: kubectl scale deployment feed-ships feed-aircraft feed-aprs \
#   feed-satellites feed-hazards -n omos --replicas=0|1
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/omos/feed-scripts.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/omos/feeds.yaml"

# -----------------------------------------------------------------------------
# 9. LEGACY AGENT (Hermes — mother) + TAK-AGENT + OneMind AI Bridge
# Secrets (hermes-env, hermes-leaf-conf, hermes-nats-ca) are SealedSecrets in platform/hermes/
# cloudtak-agent-token: out-of-band — generate in CloudTAK UI → Settings → Tokens
#   kubectl create secret generic cloudtak-agent-token -n hermes \
#     --from-literal=token='<token from CloudTAK UI: Settings → Tokens → New>'
# -----------------------------------------------------------------------------
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/agent/namespace.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/agent/pvc.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/agent/deployment.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/agent/tak-agent.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/hermes/bridge/onemind-ai-bridge.yaml"

# -----------------------------------------------------------------------------
# 10. CLOUDTAK — API + Events + Media + Tiles + MinIO + PostGIS
# Secrets (cloudtak-secrets) are SealedSecrets in platform/cloudtak/
# -----------------------------------------------------------------------------
kubectl --context "$CTX" create namespace cloudtak --dry-run=client -o yaml | kubectl --context "$CTX" apply -f -
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/api.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/events.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/media.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/tiles.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/minio.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/postgis.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/cloudflare-route.yaml"
kubectl --context "$CTX" apply -f "$PLATFORM/cloudtak/ingress.yaml"

echo ""
echo "==> Platform install complete."
echo "    TAK cluster:  see tak-cluster/DEPLOY-RUNBOOK.md"
echo "    Streams:      managed as CRDs in platform/nats/streams/ (nack controller)"
echo "    tak-agent:    needs cloudtak-agent-token secret (see step 9 above)"
echo "    feeds:        needs omos-feed-secrets (see step 8b above)"
echo "    AISSTREAM:    sign up free at https://aisstream.io to get API key"
