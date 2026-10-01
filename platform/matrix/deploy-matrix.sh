#!/usr/bin/env bash
# Matrix Stack Deployment Script for OneMind OS
# Deploys: Synapse + Element Web + mautrix-telegram + NATS-Matrix Bridge

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MATRIX_DIR="${SCRIPT_DIR}"
OVERLAY_DIR="${MATRIX_DIR}/overlays/onemind"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[MATRIX]${NC} $*"; }
success() { echo -e "${GREEN}[OK]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

check_prereqs() {
    log "Checking prerequisites..."
    
    command -v kubectl >/dev/null || error "kubectl not found"
    command -v kustomize >/dev/null || error "kustomize not found"
    command -v doctl >/dev/null || warn "doctl not found (needed for DNS)"
    
    # Check cluster access
    kubectl cluster-info >/dev/null 2>&1 || error "Cannot access Kubernetes cluster"
    
    success "Prerequisites OK"
}

check_secrets() {
    log "Checking 1Password secrets..."
    
    REQUIRED_SECRETS=(
        "MATRIX_POSTGRES_PASSWORD"
        "MATRIX_REDIS_PASSWORD"
        "SYNAPSE_REGISTRATION_SECRET"
        "SYNAPSE_SIGNING_KEY"
        "MAUTRIX_TELEGRAM_LOGIN_SECRET"
        "TELEGRAM_API_ID"
        "TELEGRAM_API_HASH"
        "TELEGRAM_BOT_TOKEN"
        "MATRIX_BOT_ACCESS_TOKEN"
        "MATRIX_BOT_DEVICE_ID"
    )
    
    for secret in "${REQUIRED_SECRETS[@]}"; do
        if ! op item get "Matrix Secrets" --field "$secret" >/dev/null 2>&1; then
            warn "Secret not found in 1Password: $secret"
        fi
    done
}

generate_secrets() {
    log "Generating secrets from 1Password..."
    
    cat > /tmp/matrix-secrets.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: matrix-secrets
  namespace: matrix
type: Opaque
stringData:
  postgres-password: "$(op item get "Matrix Secrets" --field "MATRIX_POSTGRES_PASSWORD" --reveal)"
  redis-password: "$(op item get "Matrix Secrets" --field "MATRIX_REDIS_PASSWORD" --reveal)"
  synapse-registration-secret: "$(op item get "Matrix Secrets" --field "SYNAPSE_REGISTRATION_SECRET" --reveal)"
  synapse-signing-key: "$(op item get "Matrix Secrets" --field "SYNAPSE_SIGNING_KEY" --reveal)"
  mautrix-telegram-login-secret: "$(op item get "Matrix Secrets" --field "MAUTRIX_TELEGRAM_LOGIN_SECRET" --reveal)"
  telegram-api-id: "$(op item get "Matrix Secrets" --field "TELEGRAM_API_ID" --reveal)"
  telegram-api-hash: "$(op item get "Matrix Secrets" --field "TELEGRAM_API_HASH" --reveal)"
  telegram-bot-token: "$(op item get "Matrix Secrets" --field "TELEGRAM_BOT_TOKEN" --reveal)"
  matrix-bot-access-token: "$(op item get "Matrix Secrets" --field "MATRIX_BOT_ACCESS_TOKEN" --reveal)"
  matrix-bot-device-id: "$(op item get "Matrix Secrets" --field "MATRIX_BOT_DEVICE_ID" --reveal)"
EOF
    
    success "Secrets generated at /tmp/matrix-secrets.yaml"
}

deploy() {
    log "Deploying Matrix stack to DOKS..."
    
    # Apply secrets first
    kubectl apply -f /tmp/matrix-secrets.yaml
    
    # Build and apply kustomization
    kustomize build "${OVERLAY_DIR}" | kubectl apply -f -
    
    success "Matrix stack deployed"
}

wait_for_ready() {
    log "Waiting for deployments to be ready..."
    
    DEPLOYMENTS=("synapse" "element-web" "mautrix-telegram" "nats-matrix-bridge" "postgres" "redis")
    
    for deploy in "${DEPLOYMENTS[@]}"; do
        log "Waiting for $deploy..."
        kubectl rollout status deployment/"$deploy" -n matrix --timeout=300s || error "$deploy failed to become ready"
        success "$deploy ready"
    done
}

verify_deployment() {
    log "Verifying deployment..."
    
    # Check pods
    kubectl get pods -n matrix -o wide
    
    # Check services
    kubectl get svc -n matrix
    
    # Check ingress
    kubectl get ingress -n matrix
    
    # Test Synapse health
    log "Testing Synapse health..."
    kubectl exec -n matrix deployment/synapse -- curl -sf http://localhost:8008/health || warn "Synapse health check failed"
    
    # Test Element
    log "Testing Element Web..."
    kubectl exec -n matrix deployment/element-web -- curl -sf http://localhost/ >/dev/null || warn "Element health check failed"
    
    success "Verification complete"
}

show_access_info() {
    echo
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║           MATRIX STACK DEPLOYED SUCCESSFULLY                 ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║ Access URLs:                                                  ║"
    echo "║   • Element Web:  https://element.onemindos.com               ║"
    echo "║   • Matrix API:   https://matrix.onemindos.com                ║"
    echo "║                                                              ║"
    echo "║ Next Steps:                                                   ║"
    echo "║   1. Register @hermes:onemindos.com on your homeserver       ║"
    echo "║   2. Get access token from Element (Settings → Security)     ║"
    echo "║   3. Store token in 1Password as MATRIX_BOT_ACCESS_TOKEN     ║
    echo "║   4. Configure Hermes with hermes-matrix-config.env          ║"
    echo "║   5. Restart Hermes to connect to Matrix                     ║"
    echo "║                                                              ║"
    echo "║ Rooms to Create:                                              ║"
    echo "║   • #alerts:onemindos.com                                     ║"
    echo "║   • #infra:onemindos.com                                      ║"
    echo "║   • #missions:onemindos.com                                   ║"
    echo "║   • #legacy:onemindos.com                                     ║"
    echo "║   • #voice:onemindos.com                                      ║"
    echo "║   • #sensors:onemindos.com                                    ║"
    echo "║   • #agents:onemindos.com                                     ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
}

main() {
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║      OneMind OS Matrix Stack Deployment                      ║"
    echo "║      Synapse + Element + mautrix-telegram + Bridge           ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo
    
    check_prereqs
    check_secrets
    generate_secrets
    deploy
    wait_for_ready
    verify_deployment
    show_access_info
}

main "$@"