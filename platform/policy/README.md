# platform/policy — policy and audit service (FEATURES I4)

Applied by the Argo Application `platform-policy` (`platform/argocd/platform-policy.yaml`),
**manual sync only**, into the existing namespace `onemind-core` (layer `make-sense`, where
OpenFGA already runs). Source and rules: `services/policy/README.md`.

| File | What |
|---|---|
| `deployment.yaml` | One replica, `Recreate`, ecosystem pool (`onemind.io/compute-role: ecosystem`), requests 10m/32Mi, limits 500m/128Mi, non-root, read-only root filesystem |
| `configmap.yaml` | `principals.json` (NATS identity to principal) and the OpenFGA URL, store id and model id. No secrets |
| `network-policy.yaml` | No ingress: the service talks only over NATS (outbound) |

State: none on disk. Action records are in memory (see the service README); the audit trail is
the `EVT` stream (`evt.so.>`, 90 days) and the commands are in `CMD`. No PVC, so no VolSync
source. Exposure: none (no Service, no ingress).

## Before the first sync (all owner steps)

1. **Image.** `.github/workflows/build-policy.yml` builds `ghcr.io/onemindos/policy:sha-<commit>`
   on main. Replace the placeholder tag in `deployment.yaml` with that tag and pin the digest
   from the run's `image-digest` artifact (`...:sha-<commit>@sha256:<digest>`).
2. **NATS user `policy`.** The per-service users of
   `docs/proposals/2026-10-29-nats-per-service-permissions.md` must be live (NATS Helm values,
   out of band today). The policy user's permissions are listed there.
3. **OpenFGA model.** Adopt `platform/security/openfga/` and fill `OPENFGA_STORE_ID` and
   `OPENFGA_MODEL_ID` in `configmap.yaml`. With either empty the service exits (fail closed).
4. **SealedSecrets** (commit both files here; never commit the plaintext):

   ```bash
   kubeseal --fetch-cert --controller-name=sealed-secrets-controller \
     --controller-namespace=kube-system > /tmp/pub-cert.pem

   # NATS password of the `policy` user and the OpenFGA preshared key, from your password
   # manager, exported in your shell first (never typed into a file in the repo).
   kubectl create secret generic policy-service-runtime -n onemind-core \
     --from-literal=NATS_USER=policy \
     --from-literal=NATS_PASSWORD="$POLICY_NATS_PASSWORD" \
     --from-literal=OPENFGA_TOKEN="$OPENFGA_PRESHARED_KEY" \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/policy/policy-service-runtime.sealedsecret.yaml

   # Ed25519 command-signing key. Keep the private key only in the SealedSecret; give the
   # public key (openssl pkey -in policy-signing.pem -pubout) to the bridges' config.
   openssl genpkey -algorithm ed25519 -out /tmp/policy-signing.pem
   kubectl create secret generic policy-signing-key -n onemind-core \
     --from-file=policy-signing.pem=/tmp/policy-signing.pem \
     --dry-run=client -o json \
     | kubeseal --cert /tmp/pub-cert.pem --format yaml --scope namespace-wide \
     > platform/policy/policy-signing-key.sealedsecret.yaml
   rm /tmp/policy-signing.pem
   ```

5. **Streams.** `CMD` (`cmd.>`) and `EVT` (`evt.so.>`, `evt.drone.>`, `evt.robot.>`) exist
   (`platform/nats/streams/cmd.yaml`, `evt.yaml`). BUS-DESIGN proposes `CMD` `discard: new`
   and 365 d so a full stream refuses commands instead of dropping the audit trail (ONT-P7).

Until all five are done, leave the Application unsynced. It deletes nothing if removed (no
resources-finalizer).
