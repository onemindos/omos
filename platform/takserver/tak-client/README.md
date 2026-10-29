# platform/takserver/tak-client

Deployment of the TAK Server streaming client bridge (`services/bridge/bridges/tak-client`,
FEATURES A21) in namespace `takserver`, applied only by the Argo Application
`platform-takserver-tak-client` with **manual sync**. `platform-takserver` does not read this
directory (it syncs `platform/takserver/deploy` only).

| File | What |
|---|---|
| `deployment.yaml` | `omos-tak-client`, one replica, ecosystem pool, requests 25m/32Mi, limits 250m/128Mi |
| `networkpolicy.yaml` | Adds `app: tak-client` as a source of `takserver-messaging:8089` |
| `tak-client-tls.sealedsecret.yaml` | **owner adds**: TAK client certificate (keys `tls.crt`, `tls.key`, `ca.crt`) |
| `tak-client-nats.sealedsecret.yaml` | **owner adds**: NATS credentials (keys `user`, `password`) |

## Before the first sync

### 1. TAK client certificate

TAK Server authenticates streaming clients by certificate. Mint one for CN `omos-tak-client`
with the TAK Server CA, the same way EUD certificates are made (`platform/takserver/DEPLOY-RUNBOOK.md`),
then give the user the groups whose CoT should reach the bus (today the plugin injects into
`__ANON__` and `Red`). Placeholders below; run where the CA key is, and never commit the output:

```sh
openssl req -new -newkey rsa:2048 -nodes -keyout tls.key -out client.csr \
  -subj "/O=TAK/OU=OneMind/CN=omos-tak-client"
printf 'extendedKeyUsage=clientAuth\n' > client.ext
openssl x509 -req -in client.csr -CA <tak-ca.pem> -CAkey <tak-ca.key> -passin env:TAK_CA_KEY_PASS \
  -CAcreateserial -days 365 -sha256 -extfile client.ext -out tls.crt
cp <tak-ca.pem> ca.crt     # the CA that signed the TAK Server certificate
```

(TAK Server can also sign it through enrollment, `POST /Marti/api/tls/signClient/v2` on 8446.)
Assign groups to `omos-tak-client` in TAK Server (LDAP or the user manager), as for any EUD.

### 2. Seal both Secrets (beside this file)

```sh
kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system > /tmp/pub-cert.pem

kubectl create secret generic tak-client-tls -n takserver \
  --from-file=tls.crt --from-file=tls.key --from-file=ca.crt \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/takserver/tak-client/tak-client-tls.sealedsecret.yaml

# Today: the shared ONEMIND user (the same values as omos-nats-creds). A dedicated user with
# only tak.cot.out publish and tak.cot.inject subscribe is proposal A21-P2.
kubectl create secret generic tak-client-nats -n takserver \
  --from-literal=user=<nats-user> --from-literal=password=<nats-password> \
  --dry-run=client -o json \
  | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
  > platform/takserver/tak-client/tak-client-nats.sealedsecret.yaml
shred -u tls.key client.csr 2>/dev/null || rm -f tls.key client.csr
```

Commit the two `.sealedsecret.yaml` files in a PR.

### 3. Pin the image

After `build-tak-client` runs on `main`, download its `omos-tak-client-image` artifact and
replace `ghcr.io/onemindos/omos-tak-client:0.1.0` in `deployment.yaml` with
`ghcr.io/onemindos/omos-tak-client:0.1.0@sha256:<digest>`.

## Sync and verify

```sh
argocd app sync platform-takserver-tak-client
kubectl -n takserver logs deploy/omos-tak-client
#   tak-client: connected to takserver-messaging-service...:8089 (TLS ...)
#   tak-client: TAK Server version 5.8-RELEASE-69          (if the server offers v1)
#   tak-client: switched to TAK Protocol v1 streaming      (or it stays on XML)
#   tak-client: injecting from tak.cot.inject (queue tak-cot-inject)
nats stream subjects COT                     # tak.cot.out appears next to tak.cot.inject
nats stream get COT --last-for tak.cot.out   # header Tak-Client: omos-tak-client
nats server request leafz                    # unchanged: tak-client is a direct client, not a leaf
```

Rollback: `argocd app delete platform-takserver-tak-client` (finalizer removes the
Deployment and NetworkPolicy; nothing else depends on them).
