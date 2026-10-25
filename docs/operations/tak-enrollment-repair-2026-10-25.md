# TAK certificate enrollment repair — 2026-10-25

CloudTAK's supported first-run setup reached TAK successfully with strict CA and hostname
verification, but TAK certificate enrollment returned HTTP 500. The installed API
`CertManagerService.signClient` received a null CA chain. The configured
`/certs-configmap/ca-signing.jks` was absent; existing client authentication and messaging
continued working. No human identity, global policy, CA identity, or certificate is changed.

The installed `SubscriptionManagerProxyFactory` resolves the signing provider in the
**messaging** Ignite cluster. Its `DistributedSubscriptionManager.getSigningCertChain` and
`getSigningKey` dynamically read the configured JKS and select `aliases().nextElement()`.
This repair therefore projects one encrypted JKS into messaging only. It has exactly one
`PrivateKeyEntry`, alias `tak-ca`, with the existing CA and its matching private key.
The enrollment V2 path uses this JKS; the separate legacy `CAkey` PEM path remains absent.

The CA SHA-256 fingerprint is
`72:93:0A:5B:51:50:D5:B3:66:C3:DC:F4:CF:F0:FB:E8:BC:E2:20:50:49:7F:B1:4B:07:64:C0:30:0F:B3:3A:15`.
The existing encrypted CA key was checked against that public certificate before conversion.
Only the encrypted JKS is committed, inside the namespace/name-bound
`tak-enrollment-signer` SealedSecret. No plaintext key or keystore password is in Git.

Validation used the exact currently running messaging image digest,
`sha256:b4f5d8c522bb9620e87de444b32fc0824d4b6c0d58622138f1730bbff046ba77`,
in an isolated amd64 container with no network. Its Java 17.0.20.1 loaded the JKS with the
existing configured password, found one alias and a one-certificate CA chain, checked the
public fingerprint, and signed/verified synthetic bytes with SHA256withRSA. This verifies
the installed loader's keystore and crypto semantics; live Ignite enrollment still needs
post-rollout verification. Both that container and live messaging run as UID 0, which can
read the projected JKS at mode 0400. The existing certificate ConfigMap has no colliding
`ca-signing.jks` key.

The adopted Deployment preserves the live replica count, **Recreate** strategy, commands,
resources, probes, service account, and placement. Its only effective changes are the
ConfigMap-plus-Secret projection and pinning the image to the digest already running.
The deprecated duplicate `serviceAccount` alias is omitted; `serviceAccountName` is retained.
Every existing ConfigMap file stays at its original path. No API/config-tier restart or
plaintext private-key mount is required. Recreate causes a brief messaging reconnect.

Argo owns this messaging Deployment through `platform-takserver` after adoption; the old
Helm release still records it. Do not independently upgrade, uninstall, or roll back that
release: Helm could overwrite or remove the adopted workload. The chart and OneMind values
also retain the signer projection for a future coordinated ownership transition. That
projection renders identically; the chart's other historical defaults do **not** reproduce
the entire adopted live specification.

Roll out the SealedSecret first (outer sync wave -1), verify its generated Secret exists
and contains only `ca-signing.jks`, then sync messaging at wave 0 and wait for its new pod
to become Ready. Verify the expected image digest and JKS readability without printing file
contents. Retry CloudTAK's supported first-run setup, then verify authenticated own-profile,
configured server, and live TAK connection metadata. Do not send CoT test messages or create
assets. Record completion only after those checks pass.

To roll back the mount, restore the Deployment's original `configMap: cert-migration`
volume through GitOps and wait for messaging readiness. Keep the existing CA, certificate
ConfigMap, identities, and sealed recovery material. The adopted Deployment has Prune=false and Delete=false protections. Coordinate ownership
before removing its manifest; removal must not implicitly retire the production service.
This rollback restores the earlier enrollment failure while preserving existing clients.
