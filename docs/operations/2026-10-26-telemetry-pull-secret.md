# DOKS telemetry pull-secret cleanup — 2026-10-26

The four `kube-system/doks-telemetry-config-reloader` pods were Ready but repeatedly
reported `FailedToRetrieveImagePullSecret` for `onemind-hub`. Their DigitalOcean-managed
DaemonSet uses public `ghcr.io/digitalocean-packages/busybox:1.36` and has no pull-secret
reference in its pod template. Admission copied the stale reference from
`kube-system/default`, the only ServiceAccount still referencing `onemind-hub`.
That Secret was absent from `kube-system`; no other current pod there used this account.

The DigitalOcean API confirmed `omos-cloud` had `registry_enabled: false`,
`registries: null`, and the account had no container registries. The account was created July 4 and retained resource version `2316` before repair.
DigitalOcean's registry integration
[populates default ServiceAccount pull-secret references](https://docs.digitalocean.com/products/container-registry/how-to/use-registry-docker-kubernetes/),
and Kubernetes [copies those references into newly admitted pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account).
This explains why the provider's unchanged template produced the warnings.

The authorized repair removed only `kube-system/default.imagePullSecrets`, after JSON
Patch tests for its exact UID, resource version, and old one-entry list. The guarded
server dry-run passed first. Each of the four existing DaemonSet-owned pods was then
replaced sequentially, waiting for its replacement to become Ready and checking that it
had no inherited `imagePullSecrets` before replacing the next pod. No other workload was
restarted. The provider DaemonSet specification, service-account permissions, images,
resources, policies, and secrets were unchanged. No registry credential was created,
copied, printed, or deleted. Existing events were retained.

## Recheck or repeat safely

1. Confirm the DigitalOcean cluster registry integration state and any current consumers
   before changing an account. If integration is enabled or private-registry consumers
   exist, repair the provider configuration rather than removing needed credentials.
2. Check the provider DaemonSet's template, default account, missing Secret metadata,
   pod owners, and current Ready status. Preserve a non-secret account baseline.
3. Remove only the verified stale list with compare-and-test JSON Patch; abort if its UID,
   resource version, or exact list changed. A ServiceAccount update affects newly created
   pods, so existing pods retain their old copied reference until replaced.
4. Replace only the affected DaemonSet-owned pods, one at a time. Wait for each replacement
   to become Ready, with the same image and no stale reference. Stop if a replacement fails.
5. Observe the replacement pod UIDs for at least two minutes with zero new
   `FailedToRetrieveImagePullSecret` events. Old-object events can remain until normal TTL.
   Verify all provider DaemonSet replicas are Ready and its specification is unchanged.

The account cleanup is a one-time repair to obsolete integration metadata; the provider
DaemonSet remains under DigitalOcean ownership. There is no new Argo workload to maintain.
The preserved account baseline can restore the old reference if required, although doing
so without an actual registry integration would reintroduce the warning. A future registry
integration should be configured through DigitalOcean's supported integration settings.

## Completion receipt

The account patch passed at `2026-10-26T07:25:42.188614Z`. All four
replacements were Ready by `2026-10-26T07:27:58.679003Z`. The follow-up
observation completed at `2026-10-26T07:30:03.685054Z`: **4/4 Ready, zero new
pull-secret warnings for 125 seconds**, no inherited `imagePullSecrets`, and zero
container restarts. Every replacement uses the same existing public image. The provider
DaemonSet specification remained unchanged.

| Previous pod / UID | Replacement pod / UID | Ready (UTC) |
|---|---|---|
| `doks-telemetry-config-reloader-8g449`<br>`f74efd35-0555-4a61-ab6f-76da07e708f3` | `doks-telemetry-config-reloader-jhl57`<br>`2a31dee8-0384-4931-9fc4-e7199cfaa335` | `2026-10-26T07:26:16.962440Z` |
| `doks-telemetry-config-reloader-dxj6x`<br>`d5046aa0-34ed-4a43-bb8b-4957bb8441aa` | `doks-telemetry-config-reloader-bbvdf`<br>`82cd75b3-3b14-4d01-93f4-cbb134b4f5d3` | `2026-10-26T07:26:49.515098Z` |
| `doks-telemetry-config-reloader-l9jh6`<br>`6114a180-7736-4656-8d37-0975a8eb5b27` | `doks-telemetry-config-reloader-wzzsc`<br>`d75a3032-3fc2-48d0-9e73-79163d151fb4` | `2026-10-26T07:27:23.711028Z` |
| `doks-telemetry-config-reloader-zttkr`<br>`4352e5da-0b14-4c64-b3b4-d05b1992ebb7` | `doks-telemetry-config-reloader-rh59b`<br>`8aefd706-e02d-4562-8343-c44f9ed45681` | `2026-10-26T07:27:58.575267Z` |
