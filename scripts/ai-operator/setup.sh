#!/usr/bin/env bash
# Give AI operators (Claude Code, Codex) read-only access to the cluster and the portal.
# Runs on the admin Mac. See docs/operations/ai-operator-access.md.
#
# 1. Builds ~/.kube/ai-operator.kubeconfig from the `kube/ai-operator` ServiceAccount
#    (platform/kube/ai-operator.yaml) and checks it cannot read Secrets or write.
# 2. Registers two MCP servers for Claude Code (user scope) and Codex (~/.codex/config.toml):
#    onemind-kubernetes (kubernetes-mcp-server --read-only) and onemind-backstage.
#
# Re-run it to refresh tokens. Credentials stay on this machine; nothing is printed.
set -euo pipefail

ADMIN_CONTEXT="${ADMIN_CONTEXT:-do-nyc3-omos-cloud}"
KUBECONFIG_OUT="${KUBECONFIG_OUT:-$HOME/.kube/ai-operator.kubeconfig}"
BACKSTAGE_TOKEN_FILE="${BACKSTAGE_TOKEN_FILE:-$HOME/.config/onemind/backstage-mcp.token}"
BACKSTAGE_MCP_URL="https://onemind-preview.tail717ba0.ts.net/api/mcp-actions/v1"
K8S_MCP_VERSION="0.0.67"
CODEX_CONFIG="${CODEX_CONFIG:-$HOME/.codex/config.toml}"

log() { printf '==> %s\n' "$*"; }
fail() { printf 'error: %s\n' "$*" >&2; exit 1; }

build_kubeconfig() {
  local cluster server ca token
  cluster=$(kubectl config view -o jsonpath="{.contexts[?(@.name==\"$ADMIN_CONTEXT\")].context.cluster}")
  [[ -n "$cluster" ]] || fail "kube context $ADMIN_CONTEXT not found"
  server=$(kubectl config view --raw -o jsonpath="{.clusters[?(@.name==\"$cluster\")].cluster.server}")
  ca=$(kubectl config view --raw -o jsonpath="{.clusters[?(@.name==\"$cluster\")].cluster.certificate-authority-data}")
  token=$(kubectl --context "$ADMIN_CONTEXT" -n kube get secret ai-operator-token -o jsonpath='{.data.token}' | base64 -d)
  [[ -n "$token" ]] || fail "secret kube/ai-operator-token has no token yet (is platform-kube synced?)"

  mkdir -p "$(dirname "$KUBECONFIG_OUT")"
  (
    umask 077
    cat >"$KUBECONFIG_OUT" <<EOF
apiVersion: v1
kind: Config
clusters:
  - name: omos
    cluster:
      server: $server
      certificate-authority-data: $ca
users:
  - name: ai-operator
    user:
      token: $token
contexts:
  - name: omos-ai-operator
    context:
      cluster: omos
      user: ai-operator
current-context: omos-ai-operator
EOF
  )
  log "wrote $KUBECONFIG_OUT"
}

verify_kubeconfig() {
  local k=(kubectl --kubeconfig "$KUBECONFIG_OUT")
  [[ $("${k[@]}" auth can-i list pods --all-namespaces) == yes ]] || fail "cannot list pods"
  [[ $("${k[@]}" auth can-i list applications.argoproj.io -n argocd) == yes ]] || fail "cannot list Argo applications"
  for denied in "get secrets --all-namespaces" "delete pods --all-namespaces" "create pods/exec --all-namespaces" "patch deployments --all-namespaces"; do
    # shellcheck disable=SC2086
    if [[ $("${k[@]}" auth can-i $denied 2>/dev/null || true) == yes ]]; then
      fail "ai-operator is allowed to: $denied"
    fi
  done
  log "kubeconfig verified: reads allowed; secrets, delete, exec and patch denied"
}

register_claude() {
  command -v claude >/dev/null || { log "claude CLI not found, skipping Claude Code"; return; }
  local token
  token=$(<"$BACKSTAGE_TOKEN_FILE")
  claude mcp remove --scope user onemind-kubernetes >/dev/null 2>&1 || true
  claude mcp remove --scope user onemind-backstage >/dev/null 2>&1 || true
  claude mcp add --scope user onemind-kubernetes -- \
    npx -y "kubernetes-mcp-server@$K8S_MCP_VERSION" --read-only --disable-multi-cluster \
    --kubeconfig "$KUBECONFIG_OUT" >/dev/null
  claude mcp add --scope user --transport http onemind-backstage "$BACKSTAGE_MCP_URL" \
    --header "Authorization: Bearer $token" >/dev/null
  log "registered onemind-kubernetes and onemind-backstage in Claude Code (user scope)"
}

register_codex() {
  [[ -f "$CODEX_CONFIG" ]] || { log "no $CODEX_CONFIG, skipping Codex"; return; }
  python3 - "$CODEX_CONFIG" "$KUBECONFIG_OUT" "$BACKSTAGE_TOKEN_FILE" "$BACKSTAGE_MCP_URL" "$K8S_MCP_VERSION" <<'PY'
import re, sys, pathlib
path, kubeconfig, token_file, url, version = sys.argv[1:]
token = pathlib.Path(token_file).read_text().strip()
cfg = pathlib.Path(path)
text = cfg.read_text()
blocks = {
    "onemind-kubernetes": (
        '[mcp_servers.onemind-kubernetes]\n'
        'command = "npx"\n'
        f'args = ["-y", "kubernetes-mcp-server@{version}", "--read-only", '
        f'"--disable-multi-cluster", "--kubeconfig", "{kubeconfig}"]\n'
        'startup_timeout_sec = 60\n'
    ),
    "onemind-backstage": (
        '[mcp_servers.onemind-backstage]\n'
        f'url = "{url}"\n'
        f'http_headers = {{ Authorization = "Bearer {token}" }}\n'
    ),
}
for name, block in blocks.items():
    pattern = re.compile(r'^\[mcp_servers\.' + re.escape(name) + r'\]\n(?:(?!^\[).*\n?)*', re.M)
    text = pattern.sub('', text).rstrip('\n') + '\n\n' + block
cfg.write_text(text)
PY
  chmod 600 "$CODEX_CONFIG"
  log "registered onemind-kubernetes and onemind-backstage in $CODEX_CONFIG"
}

[[ -r "$BACKSTAGE_TOKEN_FILE" ]] || fail "missing $BACKSTAGE_TOKEN_FILE"
build_kubeconfig
verify_kubeconfig
register_claude
register_codex
log "done. Start a new Claude Code or Codex session to load the servers."
