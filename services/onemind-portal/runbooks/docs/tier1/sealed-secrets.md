# Sealed Secrets

Every secret committed to Git is a Bitnami SealedSecret. The controller `sealed-secrets-controller` in `kube-system` decrypts each one into a Secret after Argo CD applies the ciphertext. Its private sealing keys are irreplaceable: without them, no committed SealedSecret can be decrypted again.

Sources: root `CLAUDE.md` ("Secrets"), `platform/core/controller.yaml`, `platform/backups/README.md`, `platform/backups/tools/README.md`, `docs/operations/home-migration.md`, [Availability and home migration](../resilience.md).

## Ownership

- `platform-core` applies `platform/core`: the controller Deployment (`--key-prefix=sealed-secrets-key`, `--update-status`), the self-hosted GitHub runner and the Tailscale policy.
- Seal a new secret against the controller's public certificate with namespace-wide scope, and commit the file beside the workload that uses it; the commands are in `CLAUDE.md`. Never commit plaintext, even test values: CI runs gitleaks on every pull request.
- Sealing is not rotation. A credential that ever appeared in Git still needs provider-side rotation.

## Investigate (read-only)

```sh
kubectl --context do-nyc3-omos-cloud -n argocd get application platform-core
kubectl --context do-nyc3-omos-cloud -n kube-system get deployment sealed-secrets-controller
kubectl --context do-nyc3-omos-cloud get sealedsecrets -A
kubectl --context do-nyc3-omos-cloud -n kube-system logs deployment/sealed-secrets-controller --tail=100
```

Because the controller runs with `--update-status`, each SealedSecret reports a `Synced` condition that shows whether it was decrypted. Never print the generated Secrets or the sealing keys.

## Key backup and recovery

- `backup-cluster-recovery` (namespace `ops`, daily at 04:00 UTC) writes an encrypted export of every Secret and ConfigMap, including the active sealing keys. It holds no other kind of object. Its check is authenticated decryption and equality, not a replay onto a new cluster.
- Keep the sealing keys, the logical-archive key and the Restic password independent of the cluster and of the operator's laptop.
- On a rebuilt cluster, restore the original sealing-key Secrets into `kube-system` before expecting existing ciphertext to decrypt. A new controller with new keys cannot decrypt old SealedSecrets; every secret would have to be resealed.
- `platform/backups/tools/decrypt.py` authenticates and decrypts an archive offline into a new private file. It does not apply anything.

## Verify completion

The controller is Ready, the affected SealedSecrets report `Synced`, and the workloads that consume those Secrets start and authenticate.

Generated workload list: [Platform core manifests app runbook](../apps/platform-core.md).
