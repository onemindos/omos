#!/usr/bin/env bash
# ---------------------------------------------------------------------------
# DEPRECATED — kept as fallback/reference only.
#
# Streams and KV buckets are now managed declaratively as Kubernetes CRDs
# in platform/nats/streams/ (38 Stream + KeyValue YAML files).  The NATS
# JetStream controller (nack) reconciles them continuously — no manual
# provisioning needed.
#
# Use this script ONLY if you need to bootstrap streams on a cluster that
# doesn't have the nack controller installed yet, or for emergency recovery.
#
# To install the controller (done by platform/install.sh step 4b):
#   helm upgrade --install nack nats/nack -n nats \
#     -f platform/nats/nack-values.secret.yaml
#   kubectl apply -f platform/nats/streams/
# ---------------------------------------------------------------------------
# Provision all captured Mac streams + KV buckets into the cloud NATS cluster,
# faithfully (same subjects/retention/limits) but bumped to num_replicas=3 (HA).
# Idempotent: `nats stream add` on an existing stream with matching config is a no-op.
set -euo pipefail
cd "$(dirname "$0")"
source ../platform/nats/creds.secret.env

CA=/tmp/nats-ca.crt
kubectl get secret -n nats nats-server-tls -o jsonpath='{.data.ca\.crt}' | base64 -d > "$CA"

# Port-forward the client service (cert SANs cover localhost/127.0.0.1).
kubectl port-forward -n nats svc/nats 4222:4222 >/tmp/nats-pf.log 2>&1 &
PF=$!
trap 'kill $PF 2>/dev/null || true' EXIT
sleep 4

NATS=(nats -s tls://localhost:4222 --user onemind --password "$ONEMIND_PW" --tlsca "$CA")

ok=0; fail=0
for f in streams/*.json; do
  name=$(python3 -c "import json;print(json.load(open('$f'))['config']['name'])")
  python3 - "$f" <<'PY' > /tmp/sc.json
import json,sys
d = json.load(open(sys.argv[1]))['config']
d['num_replicas'] = 3                       # HA across the 3 nodes
md = d.get('metadata')
if isinstance(md, dict):                     # strip internal _nats.* markers
    md = {k: v for k, v in md.items() if not k.startswith('_nats')}
    if md: d['metadata'] = md
    else: d.pop('metadata', None)
print(json.dumps(d))
PY
  if "${NATS[@]}" stream add "$name" --config /tmp/sc.json </dev/null >/tmp/add.log 2>&1; then
    echo "OK   $name"; ok=$((ok+1))
  else
    echo "FAIL $name"; sed 's/^/       /' /tmp/add.log | tail -4; fail=$((fail+1))
  fi
done
echo "=== provisioned OK=$ok FAIL=$fail (of 38) ==="
