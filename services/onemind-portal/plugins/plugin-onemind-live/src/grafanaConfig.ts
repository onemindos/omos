// Panels embedded on the Live page. UIDs and panel ids come from the provisioned
// dashboards (kube-prometheus-stack ConfigMaps and platform/monitoring dashboards in Git).
// No datasource credentials: Grafana renders them itself.
export type GrafanaPanel = { title: string; uid: string; slug: string; panelId: number };

export const GRAFANA_PANELS: GrafanaPanel[] = [
  { title: 'Cluster CPU', uid: 'efa86fd1d0c121a26444b636a3f509a8', slug: 'kubernetes-compute-resources-cluster', panelId: 7 },
  { title: 'Cluster memory', uid: 'efa86fd1d0c121a26444b636a3f509a8', slug: 'kubernetes-compute-resources-cluster', panelId: 9 },
  { title: 'Network receive', uid: 'ff635a025bcfea7bc3dd4f508990a3e9', slug: 'kubernetes-networking-cluster', panelId: 6 },
  { title: 'NATS JetStream messages', uid: 'omos-nats-observ', slug: 'nats-observability-clickhouse', panelId: 4 },
  { title: 'Postgres (portal) CPU', uid: 'onemind-backstage-postgres', slug: 'onemind-backstage-postgresql', panelId: 2 },
];

export function panelPageUrl(origin: string, panel: GrafanaPanel) {
  return `${origin}/d/${panel.uid}/${panel.slug}?orgId=1&viewPanel=${panel.panelId}`;
}

export function panelEmbedUrl(origin: string, panel: GrafanaPanel, dark: boolean) {
  const url = new URL(`${origin}/d-solo/${panel.uid}/${panel.slug}`);
  url.searchParams.set('orgId', '1');
  url.searchParams.set('panelId', String(panel.panelId));
  url.searchParams.set('from', 'now-1h');
  url.searchParams.set('to', 'now');
  url.searchParams.set('refresh', '1m');
  url.searchParams.set('theme', dark ? 'dark' : 'light');
  return url.toString();
}
