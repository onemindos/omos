# OneMind CoT Tap (TAK → NATS bus, b3 outbound leg)

A lightweight TAK client that connects to the TAK Server CoT streaming port
(`8089/ssl`) as a mutual-TLS EUD, receives the full CoT firehose for its
groups, and republishes each `<event>…</event>` to a NATS subject
(`tak.cot.out`).

## Why this exists

TAK's cluster plugin manager is **send-only** for plugins: `PluginMessenger`
and `PluginClusterMessenger` implement `send`/`networkSend` (plugin → TAK
inject) but have **no receive/subscribe path**. In cluster mode the CoT
firehose (TAK → plugin `onMessage`) is never delivered to the plugins pod — the
plugin's cluster-NATS connection shows `subs=[]`. So the NATS bridge plugin
handles the **inbound** leg (bus → TAK via `tak.cot.inject`), and this tap
handles the **outbound** leg (TAK → bus via `tak.cot.out`), using the stable
client CoT-streaming protocol instead of fragile plugin-cluster internals.

## Components

- `tap.py` — the tap (raw asyncio TLS reader → `nats-py` publisher). No pytak
  dependency; the server streams plain XML CoT to a silent client.
- `deployment.yaml` — ConfigMap (script) + Deployment. Runs `python:3.12-slim`,
  `pip install nats-py` at start.

## Certs / secret (out-of-band, NOT in git)

The tap authenticates with a client cert (`CN=omos-cot-tap`) signed by the DOKS
TAK CA. Mint + register:

```sh
# from tak-cluster/b2-enrollment (has ca.pem + ca.key, key pw: the CA key password, see Secret tak-ca-credentials in ns takserver)
openssl genrsa -out tap.key 2048
openssl req -new -key tap.key -out tap.csr \
  -subj "/C=US/ST=FL/L=Tampa/O=TAK/OU=OneMind/CN=omos-cot-tap"
openssl x509 -req -in tap.csr -CA ca.pem -CAkey ca.key -passin env:CAPASS \
  -CAcreateserial -out tap.pem -days 365 -sha256 -extfile client-ext.cnf   # clientAuth EKU

# register the user so CN maps (cluster DB)
kubectl -n takserver exec <takserver-api-pod> -- \
  java -jar /UserManager.jar usermod -p "OmosCotTap-2026#Bus" omos-cot-tap
```

Create the secret (private key never committed):

```sh
kubectl -n takserver create secret generic cot-tap-certs \
  --from-file=tap.pem --from-file=tap.key --from-file=ca.pem
```

## Deploy

```sh
kubectl -n takserver create configmap cot-tap-script --from-file=tap.py \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl -n takserver apply -f deployment.yaml
# NOTE: deployment.yaml carries a placeholder ConfigMap; re-apply the real
# script cm AFTER applying the Deployment (the placeholder overwrites it):
kubectl -n takserver create configmap cot-tap-script --from-file=tap.py \
  --dry-run=client -o yaml | kubectl apply -f -
kubectl -n takserver rollout restart deploy/cot-tap
```

## Verify

```sh
kubectl -n takserver logs deploy/cot-tap -f          # "published CoT #N"
# on the bus, the COT JetStream stream (subjects tak.cot.>) gains tak.cot.out:
nats stream get COT --last-for="tak.cot.out"
```

## Config (env)

| var            | default                      |
|----------------|------------------------------|
| `TAK_HOST`     | `takserver-messaging-service`|
| `TAK_PORT`     | `8089`                       |
| `NATS_URL`     | `nats://tak-nats-leaf:4222`  |
| `NATS_SUBJECT` | `tak.cot.out`                |
| `VERIFY_SERVER`| `false`                      |
