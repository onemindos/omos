/**
 * Catalog labels shared by the catalog "Apps" view (packages/app catalogViewFilters) and the
 * home Pillar tiles, so both count the same entities. F2 labels every Component with a role:
 * `app` for the primary entity and `workload` for the Kubernetes workloads that hang under it.
 */
export const PILLAR_LABEL_KEY = 'onemind.io/pillar';
export const ROLE_LABEL_KEY = 'onemind.io/catalog-role';
export const APP_ROLE = 'app';
export const WORKLOAD_ROLE = 'workload';

/** The same test the catalog Apps view applies (role label equals `app`). */
export function isAppComponent(labels?: Record<string, string>): boolean {
  return labels?.[ROLE_LABEL_KEY] === APP_ROLE;
}
