# Hermes secrets (out-of-band — never committed)

Per the GitOps doctrine, manifests live in git but **secret values do not**. Create these
four secrets in the `hermes` namespace before `kubectl apply -k platform/hermes/`.

## 1. `hermes-env` — full runtime config (from the Mac `.env`)
Loaded via `envFrom`. Carries provider keys + `HERMES_NATS_OWNER=zeus`,
`HERMES_NATS_SESSION_NAME=legacy`, etc. NATS_URL and AWS_* are overridden in the Deployment.
```bash
kubectl create secret generic hermes-env -n hermes \
  --from-env-file="$HOME/.hermes/.env"
```
> Before creating: consider stripping `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` from the
> file so the only copy of the (rotated) keys is the `hermes-aws` secret below.

## 2. `hermes-aws` — ROTATED Bedrock creds (🔴 gate)
Rotate the IAM keys in the AWS console first (the current ones are plaintext in the Mac
launchd plist — do not reuse them).
```bash
kubectl create secret generic hermes-aws -n hermes \
  --from-literal=AWS_ACCESS_KEY_ID='<rotated-id>' \
  --from-literal=AWS_SECRET_ACCESS_KEY='<rotated-secret>' \
  --from-literal=AWS_REGION='us-east-1'
```

## 3. `hermes-nats-ca` — cluster NATS CA (for TLS to the bus)
Copy the CA the rest of the cluster uses to talk to NATS.
```bash
kubectl get secret nats-ca-cert -n nats -o jsonpath='{.data.ca\.crt}' | base64 -d > /tmp/nats-ca.crt
kubectl create secret generic hermes-nats-ca -n hermes --from-file=ca.crt=/tmp/nats-ca.crt
rm /tmp/nats-ca.crt
```
(Verify the source secret name/key — adjust `nats-ca-cert` / `ca.crt` to match the cluster.)

## 4. `ghcr` — GitHub Container Registry pull secret
Hermes lives on ghcr.io (DOCR was at its 5-repo cap; ghcr is free for the org). The image is
`ghcr.io/onemindos/hermes:legacy-v1`. Create the pull secret with a GitHub token that has
`read:packages`:
```bash
kubectl create secret docker-registry ghcr -n hermes \
  --docker-server=ghcr.io \
  --docker-username=Zeus-Delacruz \
  --docker-password="$(gh auth token)"
```
> If the ghcr package is private, ensure the token can read it (or make the package internal
> to the `onemindos` org).

## Verify (step-1 proof of done)
```bash
kubectl -n hermes get pods
nats req '$SRV.INFO.agents' ''      # expect TWO legacy/zeus instances (Mac + cloud)
nats sub 'agents.hb.*.*.*'          # expect both heartbeats
# failover: kill one body, prompt agents.prompt.legacy.zeus — other answers
```
