# Deployment verification — 2026-10-06

Verified on `do-nyc3-omos-cloud` at the private preview endpoint.

- Backstage 1.55.3: immutable dependency install, TypeScript checking, backend build and published container passed.
- Application image: `ghcr.io/onemindos/onemind-portal@sha256:587cefd427672ee850f2122b522a48e2a69cf7af57b243c48050b186b1fbfd3a`.
- Server-side validation accepted the complete Kustomize output; the portal and persistent PostgreSQL were ready.
- HTTPS returned 200 with certificate validation. Browser sign-in correctly redirected to the private Authentik route.
- Complete OIDC authorization-code login and refresh succeeded. The resolved identity was `user:default/akadmin`, with `group:default/onemind-admins` membership.
- Catalog API returned 230 entities (including three source locations), with zero processing errors. Inventory contained 67 app links and 122 workload components.
- Authenticated operations API returned 122 workloads, four backup schedules, and `partial: false`.
- The Kubernetes plugin queried Homepage and returned 12 resource groups with no errors, including workload metrics.
- Anonymous catalog and operations API requests returned 401.
- The service account could observe pods and Argo applications; Secrets and workload deletion were denied. ConfigMaps are excluded from Kubernetes fetches.
- The release backup uploaded 14 encrypted plugin databases to Spaces, downloaded and authenticated the archive, restored the catalog into a temporary database, verified 18 tables, and completed successfully. The first recovery exercise also succeeded.
- Homepage remained Synced and Healthy during preview deployment.

Browser review of the authenticated home layout and Homepage cutover are outstanding. This is a private preview deployment, not a claim of high availability or offline support. See README.md for limitations and recovery requirements.
