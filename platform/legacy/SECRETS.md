# Hermes secrets (out-of-band — never committed)

Per the GitOps doctrine, manifests live in git but **secret values do not**. Create these in
the `hermes` namespace before `kubectl apply -k platform/hermes/`.

Context: `kubectl --context tailscale-operator.tail717ba0.ts.net ...` (kubectl-over-tailnet).

## 1. `hermes-env` — full runtime config (from the Mac `.env`)
Carries provider keys + `HERMES_NATS_OWNER=zeus`, `HERMES_NATS_SESSION_NAME=legacy`, and
(for now) the AWS Bedrock keys. `NATS_URL` is overridden in the Deployment to the sidecar.
```bash
kubectl -n hermes create secret generic hermes-env \
  --from-env-file="$HOME/.hermes/.env"
```

## 2. `hermes-leaf-conf` — NATS leaf sidecar config (contains leaf_tak password)
The sidecar gives Hermes a local plaintext `localhost:4222` and carries the TLS+auth uplink
into the cluster ONEMIND account — exactly like the Mac leaf. Create `leaf.conf`:
```hocon
port: 4222
http: 8222
server_name: hermes-legacy-leaf
accounts { OMOS { users: [ { user: local } ] } }
no_auth_user: local
leafnodes {
  remotes: [
    {
      urls: [ "tls://leaf_tak:<LEAF_TAK_PW>@nats.nats.svc.cluster.local:7422" ]
      account: OMOS
      tls { ca_file: "/etc/nats-ca/ca.crt" }
    }
  ]
}
```
`<LEAF_TAK_PW>` = the leaf_tak password (same one the Mac leaf uses). Then:
```bash
kubectl -n hermes create secret generic hermes-leaf-conf --from-file=leaf.conf=./leaf.conf
rm ./leaf.conf
```

## 3. `hermes-nats-ca` — cluster NATS CA (for the leaf's TLS uplink)
```bash
kubectl -n nats get secret nats-server-tls -o jsonpath='{.data.ca\.crt}' | base64 -d > /tmp/nats-ca.crt
kubectl -n hermes create secret generic hermes-nats-ca --from-file=ca.crt=/tmp/nats-ca.crt
rm /tmp/nats-ca.crt
```

## 4. `ghcr` — GitHub Container Registry pull secret
Image is `ghcr.io/onemindos/hermes:legacy-v1` (DOCR was at its 5-repo cap; ghcr is free).
```bash
kubectl -n hermes create secret docker-registry ghcr \
  --docker-server=ghcr.io \
  --docker-username=Zeus-Delacruz \
  --docker-password="$(gh auth token)"
```
> If the ghcr package is private, ensure the token can read it, or make the package internal
> to the `onemindos` org.

## Later: `hermes-aws` — ROTATED Bedrock creds
For now the AWS keys ride along in `hermes-env`. After rotation, create `hermes-aws` and
uncomment the `AWS_*` env overrides in `deployment.yaml`:
```bash
kubectl -n hermes create secret generic hermes-aws \
  --from-literal=AWS_ACCESS_KEY_ID='<rotated-id>' \
  --from-literal=AWS_SECRET_ACCESS_KEY='<rotated-secret>' \
  --from-literal=AWS_REGION='us-east-1'
```

## Go live
```bash
kubectl create namespace hermes --dry-run=client -o yaml | kubectl apply -f -
# ...create the 4 secrets above...
kubectl apply -k platform/hermes/
kubectl -n hermes get pods -w
```

## Proof of done (step 1)
```bash
# from nats-box or Mac with the onemind context:
nats req '$SRV.INFO.agents' ''      # expect TWO legacy/zeus instances (Mac + cloud)
nats sub 'agents.hb.*.*.*'          # expect both heartbeats
# failover: stop one body → prompt agents.prompt.legacy.zeus, the other answers
```
