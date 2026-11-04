import { buildAppLinks, exploreUrl, tabFilters } from './links';

const grafana = 'https://grafana.example.ts.net/';
const entity = (annotations: Record<string, string>, links: any[] = []) => ({ metadata: { annotations, links } });

test('every link is built from the generated annotations', () => {
  const links = buildAppLinks(
    entity({
      'argocd/app-name': 'platform-nats',
      'onemind.io/argocd-apps': 'platform-nats,edge-nack',
      'onemind.io/grafana-dashboard-path': 'd/onemind-workload/onemind-workload?var-namespace=nats',
      'onemind.io/logs-query': '{k8s_namespace_name="nats"}',
      'onemind.io/kubernetes-namespaces': 'nats',
      'onemind.io/recovery-runbook': 'apps/nats/',
      'github.com/project-slug': 'onemindos/omos',
    }),
    grafana,
  );
  expect(links.argo.map(a => a.name)).toEqual(['platform-nats', 'edge-nack']);
  expect(links.argo[0].url).toBe('https://argocd.tail717ba0.ts.net/applications/argocd/platform-nats');
  expect(links.grafanaDashboard).toBe('https://grafana.example.ts.net/d/onemind-workload/onemind-workload?var-namespace=nats');
  expect(links.runbook).toBe('/docs/default/component/onemind-runbooks/apps/nats/');
  expect(links.source).toBe('https://github.com/onemindos/omos');
  expect(links.opencost).toContain('opencost');
});

test('Explore opens the Loki data source with the app query and a one hour range', () => {
  const url = new URL(exploreUrl(grafana, '{k8s_namespace_name="nats"}'));
  expect(url.origin + url.pathname).toBe('https://grafana.example.ts.net/explore');
  const left = JSON.parse(url.searchParams.get('left')!);
  expect(left.datasource).toBe('Loki');
  expect(left.queries[0].expr).toBe('{k8s_namespace_name="nats"}');
  expect(left.range).toEqual({ from: 'now-1h', to: 'now' });
});

test('an app without the data has no such link, and unsafe URLs are dropped', () => {
  const none = buildAppLinks(entity({ 'onemind.io/launcher-url': 'javascript:alert(1)' }), grafana);
  expect(none).toEqual({ argo: [], openApp: undefined, grafanaDashboard: undefined, logsQuery: undefined, logsExplore: undefined, opencost: undefined, runbook: undefined, source: undefined });
  expect(buildAppLinks(entity({ 'github.com/project-slug': 'not a slug' })).source).toBeUndefined();
  expect(buildAppLinks(entity({}, [{ title: 'Open application', url: 'https://example.com/app' }])).openApp).toBe('https://example.com/app');
  expect(buildAppLinks(entity({ 'onemind.io/logs-query': '{a="b"}' })).logsExplore).toBeUndefined();
});

test('tabs are shown only when their data exists', () => {
  const k8s = entity({ 'onemind.io/logs-query': 'q', 'onemind.io/kubernetes-namespaces': 'a', 'onemind.io/recovery-runbook': 'apps/a/', 'grafana/dashboard-selector': 'x', 'github.com/project-slug': 'o/r' });
  const saas = entity({});
  for (const filter of Object.values(tabFilters)) {
    expect(filter(k8s)).toBe(true);
    expect(filter(saas)).toBe(false);
  }
  expect(tabFilters.metrics(entity({ 'prometheus.io/service-name': 'nats' }))).toBe(true);
});
