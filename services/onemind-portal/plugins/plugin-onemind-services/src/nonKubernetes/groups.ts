import { hostingKind, HOSTING_KIND_LABEL, HostingKind } from '@internal/plugin-onemind-common';

export type AppRow = {
  ref: string;
  name: string;
  title: string;
  kind: HostingKind;
  /** The Gatus endpoint name, when the app has a live check. */
  check?: string;
  owner?: string;
};

type EntityLike = {
  kind: string;
  metadata: { name: string; namespace?: string; title?: string; annotations?: Record<string, string> };
  spec?: Record<string, unknown>;
};

/** Order the page shows the groups in: what is running outside the cluster first, unknown (never hidden) last but one. */
export const GROUP_ORDER: HostingKind[] = ['saas', 'home', 'digitalocean', 'device', 'desktop', 'mobile', 'web', 'planned', 'unknown'];

/** Pure: Components that are not on Kubernetes, as rows. Kubernetes apps are left out by definition. */
export function nonKubernetesRows(entities: EntityLike[]): AppRow[] {
  return entities
    .filter(e => e.kind === 'Component')
    .map(e => ({
      ref: `component:${e.metadata.namespace ?? 'default'}/${e.metadata.name}`,
      name: e.metadata.name,
      title: e.metadata.title ?? e.metadata.name,
      kind: hostingKind(e),
      check: e.metadata.annotations?.['onemind.io/health-check'],
      owner: typeof e.spec?.owner === 'string' ? e.spec.owner.replace(/^group:(default\/)?/, '') : undefined,
    }))
    .filter(row => row.kind !== 'kubernetes')
    .sort((a, b) => a.title.localeCompare(b.title));
}

export type Group = { kind: HostingKind; label: string; rows: AppRow[] };

/** Pure: rows grouped by hosting kind in GROUP_ORDER; empty groups are dropped. */
export function groupByHosting(rows: AppRow[]): Group[] {
  return GROUP_ORDER.map(kind => ({ kind, label: HOSTING_KIND_LABEL[kind], rows: rows.filter(r => r.kind === kind) })).filter(g => g.rows.length);
}
