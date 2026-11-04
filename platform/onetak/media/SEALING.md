# Sealing the onemind-media Secrets (owner)

No password is in Git. The owner generates them and seals them with the cluster's public
certificate (omos `CLAUDE.md` "Secrets"). Run from the repository root on a machine with
`kubectl` access; nothing is applied to the cluster.

```bash
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem

# 1. MediaMTX users: one key per user in mediamtx.yml (URL-safe random values).
args=()
for u in mavlink-bridge edge-home onemind-client tak-viewer takml video-registry; do
  args+=(--from-literal="$u=$(openssl rand -hex 24)")
done
kubectl create secret generic onemind-media-users -n onetak "${args[@]}" \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/onetak/media/onemind-media-users.sealedsecret.yaml
unset args

# 2. video-registry's NATS identity (keys user, password). Until the per-service NATS users of
#    BUS-DESIGN D6 exist, use the same user as omos-nats-creds.
read -rp "NATS user: " NU; read -rsp "NATS password: " NP; echo
kubectl create secret generic video-registry-nats -n onetak \
  --from-literal=user="$NU" --from-literal=password="$NP" \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/onetak/media/video-registry-nats.sealedsecret.yaml
unset NU NP
```

Then add both files to `resources:` in `kustomization.yaml` and commit them.

Hand out the plaintext only through the owner's password manager:

| User | Who gets it |
|---|---|
| `edge-home` | the site edge MediaMTX (`infra/edge/mac-media/edge-media.sh keychain-set`) |
| `mavlink-bridge` | companion computers that push straight to the hub |
| `onemind-client` | OneMind clients (Mac app, L0 web, headset) until JWT auth (phase B) |
| `tak-viewer` | entered once in ATAK/iTAK video aliases (CoT URLs never carry credentials) |
| `takml` | the takml worker, as its own takml-namespace SealedSecret at cut-over |
| `video-registry` | nobody; only the registry pod reads it |
