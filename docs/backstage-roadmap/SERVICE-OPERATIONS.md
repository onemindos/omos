# Entity Operations integration

The new catalog entity Operations tab appears for inventory entities with `backstage.io/kubernetes-namespace`. It combines read-only Argo Application observations, the existing bounded Prometheus namespace summaries, Grafana navigation and the monitoring runbook. No cluster writes, sync/rollback buttons, cloud mutation credential or generic proxy is introduced.

## Selection decision

The maintained community Argo CD and Grafana frontend packages both support the new frontend system. They are not rejected as incompatible. They are deferred here because they add provider-specific backend/proxy setup and credentials, while OneMind already has an observational Kubernetes role for Argo Applications and a bounded fixed-target metrics adapter. The small native adapter reuses these existing permissions and avoids installing a broad provider proxy or token merely to render service pages. Community plugin adoption remains a future reviewed option when provider-specific history/alerts and credential scopes are ready.

Sources reviewed:
- https://github.com/backstage/community-plugins/blob/main/workspaces/argocd/plugins/argocd/README.md
- https://github.com/backstage/community-plugins/blob/main/workspaces/grafana/plugins/grafana/README.md
- https://github.com/backstage/community-plugins/blob/main/docs/compatibility/new-frontend-system.md
- Installed `@backstage/plugin-catalog-react` 3.2.3 EntityContentBlueprint and useEntity types.

## Behavior and scope

`GET /api/onemind-services/overview?namespace=nats` requires a signed-in user; machine credentials are denied. Namespace input is DNS-label validated. The destination is fixed Kubernetes API, TLS-verified by its service CA. Read-only role remains unchanged. No Secrets, ConfigMaps, provider tokens, operation messages or source configuration are returned.

Argo destinations are matched by namespace and intentionally labeled namespace-level association, not exclusive service ownership. Multi-namespace apps may not match; absence means no matching destination, not proof of no deployment. The response returns bounded selected Application fields: health, sync, revision, last-operation phase/time and tracked resource count. It is not a full revision history. The fetch timestamp records when Backstage retrieved Kubernetes data; each Application exposes its separate controller `status.reconciledAt`. A missing controller timestamp is unknown freshness; older than ten minutes is visibly stale, a display threshold rather than an Argo SLA. A recent fetch cannot override this warning. Metrics with both values null are no-data; one missing value or a failed query is partial. Metric sample timestamps remain separate from fetch time. Observations are cached for a minute, with shared concurrent refresh, an eight-second wall deadline and 20 MB response cap. Missing sources display unavailable, not healthy.

Metrics reuse the existing configured namespace scope (backstage, auth, nats, takserver, onetak). Other namespaces explicitly display not-configured. NATS has a verified source-defined dashboard UID `omos-nats-observ` from platform/monitoring/grafana-dashboards-nats.yaml. Other namespaces open the Grafana dashboard browser and explicitly say mapping-required. Grafana login stays in Grafana; dashboard embedding, Grafana alert retrieval, and per-pod metrics are not claimed. A dashboard link does not verify its current runtime availability.

## Release review

TypeScript compilation passed with Node 22 and installed dependencies. No tests were added or run; no merge/deployment performed. Before release, review browser tab visibility, signed-in/anonymous/machine behavior, Argo results and explicit absent-mapping states in staging. The live deployment remains unchanged by this PR.
