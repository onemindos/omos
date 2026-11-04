import { safeUrl } from '@internal/plugin-onemind-common';

type EntityLike = {
  metadata: {
    annotations?: Record<string, string>;
    links?: { title?: string; url: string }[];
  };
};

// Tailnet hosts from platform/argocd and platform/ecosystem/resources/opencost-tailnet-ingress.yaml.
export const ARGOCD_URL = 'https://argocd.tail717ba0.ts.net';
export const OPENCOST_URL = 'https://opencost.tail717ba0.ts.net';
export const RUNBOOK_SITE = '/docs/default/component/onemind-runbooks/';
// The Loki data source name provisioned in platform/monitoring/grafana-datasources.yaml.
const LOKI_DATASOURCE = 'Loki';

export type AppLinks = {
  openApp?: string;
  argo: { name: string; url: string }[];
  grafanaDashboard?: string;
  logsQuery?: string;
  logsExplore?: string;
  opencost?: string;
  runbook?: string;
  source?: string;
};

/** Grafana Explore URL for a LogQL query (legacy `left` form, accepted by every Grafana 9+). */
export function exploreUrl(grafanaDomain: string, query: string) {
  const left = { datasource: LOKI_DATASOURCE, queries: [{ refId: 'A', expr: query }], range: { from: 'now-1h', to: 'now' } };
  return `${grafanaDomain.replace(/\/$/, '')}/explore?orgId=1&left=${encodeURIComponent(JSON.stringify(left))}`;
}

/** Every link on the standard page comes from the generated annotations; nothing is composed per app by hand. */
export function buildAppLinks(entity: EntityLike, grafanaDomain?: string): AppLinks {
  const a = entity.metadata.annotations ?? {};
  const raw = a['onemind.io/launcher-url'] ?? entity.metadata.links?.find(l => l.title === 'Open application')?.url;
  const argoNames = [a['argocd/app-name'], ...(a['onemind.io/argocd-apps'] ?? '').split(',')].map(x => x?.trim()).filter((x): x is string => !!x);
  const dashboard = a['onemind.io/grafana-dashboard-path'];
  const query = a['onemind.io/logs-query'];
  const runbook = a['onemind.io/recovery-runbook'];
  const slug = a['github.com/project-slug'];
  return {
    openApp: raw ? safeUrl(raw) : undefined,
    argo: [...new Set(argoNames)].map(name => ({ name, url: `${ARGOCD_URL}/applications/argocd/${encodeURIComponent(name)}` })),
    grafanaDashboard: dashboard && grafanaDomain ? `${grafanaDomain.replace(/\/$/, '')}/${dashboard}` : undefined,
    logsQuery: query,
    logsExplore: query && grafanaDomain ? exploreUrl(grafanaDomain, query) : undefined,
    opencost: a['onemind.io/kubernetes-namespaces'] ? OPENCOST_URL : undefined,
    runbook: runbook ? `${RUNBOOK_SITE}${runbook}` : undefined,
    source: slug && /^[\w.-]+\/[\w.-]+$/.test(slug) ? `https://github.com/${slug}` : undefined,
  };
}

/** Tabs without data are hidden, so every tab has one predicate on the annotations. */
const has = (e: EntityLike, key: string) => !!e.metadata.annotations?.[key];
export const tabFilters = {
  logs: (e: EntityLike) => has(e, 'onemind.io/logs-query'),
  cost: (e: EntityLike) => has(e, 'onemind.io/kubernetes-namespaces'),
  runbook: (e: EntityLike) => has(e, 'onemind.io/recovery-runbook'),
  metrics: (e: EntityLike) => has(e, 'grafana/dashboard-selector') || has(e, 'grafana/alert-label-selector') || has(e, 'onemind.io/prometheus-scopes') || has(e, 'prometheus.io/service-name'),
  code: (e: EntityLike) => has(e, 'github.com/project-slug'),
};
