# Parked services

Services that were removed from the cluster on purpose but kept so they can come back, for example on the home cluster.

**Nothing in this folder is deployed.** No ArgoCD Application reads `platform/parked/`, so a file here is inert.

Each parked service keeps, in a folder named after it:

- the ArgoCD Application it ran under (chart, version, destination), moved here unchanged;
- its values file, if it had one (moved from `platform/ecosystem/values/` or `platform/values/`);
- anything else needed to bring it back (namespace, network policy, secrets as SealedSecrets).

Its data backup and the steps to restore it are recorded in `docs/operations/parked-services.md`.

To bring a service back: move its Application back to where it came from (the register says where), restore its data from the backup listed there, and let ArgoCD sync.
