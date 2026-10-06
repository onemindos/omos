# Availability and home migration

## Application availability

Two Backstage replicas share PostgreSQL and the same identity/session configuration. Rolling updates keep existing ready replicas while starting one extra replica. A disruption budget preserves one replica during voluntary eviction; preferred node separation reduces shared node risk. Reserve 500m CPU and 2 GiB memory for the two app replicas, plus 250m CPU and 1 GiB during rollout. Limits permit higher actual usage. PostgreSQL and identity remain single-instance dependencies. This is not full high availability.

Before each version upgrade, review plugin/database migration compatibility with overlapping versions. For an incompatible migration, schedule downtime, stop both replicas, back up and use a coordinated update. Rollback of an image does not reverse a database migration.

## Recovery

Daily encrypted backups contain all Backstage plugin databases. `platform/backstage/recovery/` provides an isolated restore exercise. Preserve the backup encryption key, object-store credentials, identity configuration, registry access and SealedSecrets controller key in independent secure storage. Git does not contain decrypted credentials. Authentik's database and its recovery are separate and are not covered by a Backstage archive.

## Moving home

1. Prepare Kubernetes, sufficient capacity for three app pods during rollouts, PostgreSQL 17, retained persistent storage and reliable power/networking.
2. Install mesh routing and private HTTPS, preserving or deliberately replacing portal and identity hostnames. Coordinate any hostname changes with OIDC redirect URIs, app base URLs and TLS.
3. Supply registry pull credentials and reseal secrets for the destination controller; old SealedSecrets cannot simply be decrypted by a new controller.
4. Keep backups outside the home host. Verify access and isolated restore before cutover.
5. Deploy PostgreSQL with a suitable home storage class and restore all original plugin databases while the destination app is stopped.
6. Deploy the pinned portal image and Git catalog with destination database/network configuration. Validate identity, permissions, documentation, catalog and live observations.
7. Review destination cluster credentials, inventory context guards and monitoring endpoints. The inventory script currently guards the DigitalOcean context; change it deliberately for the new cluster.
8. Cut over routing only after review. Retain the cloud deployment and its volume until destination operation and a fresh backup are verified.

## Offline limits

Home-hosted pages need a reachable home server. Cloud links, provider APIs, remote identity and cloud monitoring need their upstream networks. Downloaded runbooks can be read offline. No service worker or offline portal is claimed. Home hosting has not been deployed or tested on future hardware.
